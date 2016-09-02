// Package insert inserts data disclosures into the Data Track database.
package insert

import (
	"time"

	"github.com/marcelfarres/datatrack/database"
	"github.com/marcelfarres/datatrack/model"
)

// DoInsert performs an insert.
func DoInsert(ins *model.Insert) (err error) {
	if err = database.SetUser(ins.User); err != nil {
		return
	}

	for _, disc := range ins.Disclosures {
		if err = addActionDisclose(disc, ""); err != nil {
			return
		}
	}

	return
}

func addActionDisclose(action model.ActionDisclose, originID string) (err error) {
	var senderID string

	// add attributes
	attributeIDs := make([]string, 0, len(action.Attribute))
	for _, attr := range action.Attribute {
		// we use make to derive our own identifier
		a, err := model.MakeAttribute(attr.Name, attr.Type, attr.Value)
		if err != nil {
			return err
		}
		if err = database.AddAttribute(a); err != nil {
			return err
		}
		attributeIDs = append(attributeIDs, a.ID)
	}

	// sender should be?
	if len(originID) == 0 && len(action.Sender.ID) == 0 {
		senderID = database.Self
	} else {
		if err = addOrUseOrg(action.Sender); err != nil {
			return
		}
		senderID = action.Sender.ID
	}

	// add recipient
	if err = addOrUseOrg(action.Recipient); err != nil {
		return
	}

	// add disclosure
	t := time.Now().String()
	if len(action.Disclosure.Timestamp) > 0 {
		t = action.Disclosure.Timestamp
	}
	// we use make to derive our own identifier
	disc, err := model.MakeDisclosure(senderID, action.Recipient.ID, t,
		action.Disclosure.PrivacyPolicyHuman, action.Disclosure.PrivacyPolicyMachine,
		action.Disclosure.DataLocation, action.Disclosure.API)
	if err != nil {
		return
	}
	if err = database.AddDisclosure(disc); err != nil {
		return
	}

	// add disclosed (link between attributes and disclosure)
	if err = database.AddDisclosed(model.Disclosed{
		Disclosure: disc.ID,
		Attribute:  attributeIDs,
	}); err != nil {
		return
	}

	// add downstream link
	if len(originID) != 0 {
		err = database.AddDownstream(model.Downstream{
			Origin: originID,
			Result: disc.ID,
		})
		return
	}

	// deal with downstream
	for _, down := range action.Downstream {
		if err = addActionDisclose(down, disc.ID); err != nil {
			return
		}
	}

	return
}

// addOrUseOrg adds on organization if unknown, otherwise
func addOrUseOrg(org model.Organization) (err error) {
	if len(org.Name) != 0 {
		return database.AddOrganization(org)
	}
	if _, err = database.GetOrganization(org.ID); err != nil {
		if err.Error() == database.NoSuchOrgError {
			return database.AddOrganization(org)
		}
	}
	return // already exists or error to forward
}
