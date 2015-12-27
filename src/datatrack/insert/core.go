// Package insert inserts data disclosures into the Data Track database.
package insert

import (
	"datatrack/database"
	"datatrack/model"
	"time"
)

// DoInsert performs an insert.
func DoInsert(ins *model.Insert) (err error) {
	err = database.SetUser(ins.User)
	if err != nil {
		return
	}

	for _, disc := range ins.Disclosures {
		err = addActionDisclose(disc, "")
		if err != nil {
			return
		}
	}

	return nil
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
		err = database.AddAttribute(a)
		if err != nil {
			return err
		}
		attributeIDs = append(attributeIDs, a.ID)
	}

	// sender should be?
	if len(originID) == 0 && len(action.Sender.ID) == 0 {
		senderID = database.Self
	} else {
		err = addOrUseOrg(action.Sender)
		if err != nil {
			return
		}
		senderID = action.Sender.ID
	}

	// add recipient
	err = addOrUseOrg(action.Recipient)
	if err != nil {
		return
	}

	// add disclosure
	thetime := time.Now().String()
	if len(action.Disclosure.Timestamp) > 0 {
		thetime = action.Disclosure.Timestamp
	}
	// we use make to derive our own identifier
	disc, err := model.MakeDisclosure(senderID, action.Recipient.ID, thetime,
		action.Disclosure.PrivacyPolicyHuman, action.Disclosure.PrivacyPolicyMachine,
		action.Disclosure.DataLocation, action.Disclosure.API)
	if err != nil {
		return
	}
	err = database.AddDisclosure(disc)
	if err != nil {
		return
	}

	// add disclosed (link between attributes and disclosure)
	err = database.AddDisclosed(model.Disclosed{
		Disclosure: disc.ID,
		Attribute:  attributeIDs,
	})
	if err != nil {
		return
	}

	// add downstream link
	if len(originID) != 0 {
		err = database.AddDownstream(model.Downstream{
			Origin: originID,
			Result: disc.ID,
		})
	}

	// deal with downstream
	for _, down := range action.Downstream {
		err = addActionDisclose(down, disc.ID)
		if err != nil {
			return
		}
	}

	return nil
}

// addOrUseOrg adds on organization if unknown, otherwise
func addOrUseOrg(org model.Organization) (err error) {
	if len(org.Name) == 0 {
		_, err = database.GetOrganization(org.ID)
		if err != nil {
			if err.Error() == database.NoSuchOrgError {
				return database.AddOrganization(org)
			}
		}
		return // already exists or error to forward
	}
	return database.AddOrganization(org)
}
