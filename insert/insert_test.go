package insert

import (
	"os"
	"path"
	"testing"

	"github.com/boltdb/bolt"

	"github.com/marcelfarres/datatrack/database"
	"github.com/marcelfarres/datatrack/ephemeral"
	"github.com/marcelfarres/datatrack/model"
)

func TestDoInsert(t *testing.T) {
	ephemeral.Setup()
	testFile := path.Join(os.TempDir(), "dttmptestdoinsert.db")
	var err error
	database.DB, err = bolt.Open(testFile, 0600, nil)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	err = database.Setup()
	if err != nil {
		t.Fatalf("failed to setup database: %s", err)
	}
	defer os.Remove(testFile)
	defer database.DB.Close()

	var ins model.Insert
	var actions []model.ActionDisclose

	var action model.ActionDisclose
	action.Attribute = make([]model.Attribute, 2)
	action.Attribute[0] = model.Attribute{
		Name:  "name0",
		Type:  "type0",
		Value: "value0",
	}
	action.Attribute[1] = model.Attribute{
		Name:  "name1",
		Type:  "type1",
		Value: "value1",
	}

	action.Recipient = model.Organization{
		ID:          "Spotify",
		Name:        "Spotify",
		Street:      "Birger Jarlsgatan 61, 10tr",
		Zip:         "113 56 ",
		City:        "Stockholm",
		State:       "Stockholm",
		Country:     "Sweden",
		Telephone:   "",
		Email:       "info@spotify.com",
		URL:         "http://www.spotify.com",
		Description: "Med Spotify är det enkelt att hitta rätt musik för varje tillfälle – på mobilen, datorn, surfplattan m.fl",
	}

	action.Disclosure = model.Disclosure{
		PrivacyPolicyHuman:   "humanpolicy",
		PrivacyPolicyMachine: "machinepolicy",
		DataLocation:         "location",
		API:                  "api",
	}

	actions = append(actions, action)
	ins.Disclosures = actions
	ins.User = model.User{
		Name:    "Bob",
		Picture: "selfie.png",
	}

	err = DoInsert(&ins)
	if err != nil {
		t.Fatalf("failed to do insert: %s", err)
	}

}
