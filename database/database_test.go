package database

import (
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/boltdb/bolt"

	"github.com/pylls/datatrack/ephemeral"
	"github.com/pylls/datatrack/model"
)

func TestBoltDatabase(t *testing.T) {
	ephemeral.Setup()
	testFile := path.Join(os.TempDir(), "dtmptestfile.db")
	var err error
	DB, err = bolt.Open(testFile, 0600, nil)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	err = Setup()
	if err != nil {
		t.Fatalf("failed to setup database: %s", err)
	}
	defer os.Remove(testFile)
	defer DB.Close()

	// adds 2 attributes with the same type
	testCoreAttribute(t)
	// adds 2 organizations: Facebook and Spotify
	testCoreOrganization(t)
	// adds 2 disclosures to Facebook at different times
	testCoreDisclosure(t)
	// links the 2 disclosures to having disclosued all attributes twice
	testCoreDisclosed(t)
	// test explicitly disclosed distinction for disclosures
	testExplicitlyDisclosed(t)
	testExplicitlyDisclosedAttributes(t)
	// adds 1 explicit disclosure to Spotify with 2 attributes, then 2 implicit
	// disclosures with a new attribute each
	testImplicitDisclosures(t)
	// direct downstream disclosure function
	testDownstreamDisclosures(t)
	// the user
	testUser(t)
	// WIP categories
	testCategories(t)
	// extended organization and attribute tests (what did not fit above)
	testExtended(t)
	// creates some coordinates
	testCoordinates(t)
}

func testCoreAttribute(t *testing.T) {
	attr, err := model.MakeAttribute("Artist", "music", "Shakira")
	if err != nil {
		t.Fatalf("failed to make attribute: %s", err)
	}
	err = AddAttribute(attr)
	if err != nil {
		t.Fatalf("failed to add attribute: %s", err)
	}

	IDs, err := GetAttributeIDs()
	if err != nil {
		t.Fatalf("failed to get attribute IDs: %s", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("unexpected number of attribute IDs, expected %d, got %d", 1, len(IDs))
	}
	for _, id := range IDs {
		att, errr := GetAttribute(id)
		if errr != nil {
			t.Fatalf("failed to get attribute: %s", errr)
		}
		if attr.Name != att.Name {
			t.Fatalf("invalid attribute name, expected %s, got %s", attr.Name, att.Name)
		}
	}
	attr2, err := model.MakeAttribute("Artist", "music", "Dylan")
	if err != nil {
		t.Fatalf("failed to make attribute: %s", err)
	}
	err = AddAttribute(attr2)
	if err != nil {
		t.Fatalf("failed to add attribute: %s", err)
	}
	IDs, err = GetAttributeIDs()
	if err != nil {
		t.Fatalf("failed to get attribute IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of attribute IDs, expected %d, got %d", 2, len(IDs))
	}
	for _, id := range IDs {
		att, errr := GetAttribute(id)
		if errr != nil {
			t.Fatalf("failed to get attribute: %s", errr)
		}
		if attr.Name != att.Name {
			t.Fatalf("invalid attribute name, expected %s, got %s", attr.Name, att.Name)
		}
	}

	typeIDs, err := GetAttributeTypeIDs()
	if err != nil {
		t.Fatalf("failed to get attribute type IDs: %s", err)
	}
	if len(typeIDs) != 1 {
		t.Fatalf("unexpected number of attribute type IDs, expected %d, got %d", 1, len(typeIDs))
	}
	for _, id := range typeIDs {
		values, err := GetAttributeTypeValues(id)
		if err != nil {
			t.Fatalf("failed to get attribute type values: %s", err)
		}
		if len(values) != 2 {
			t.Fatalf("unexpected number of attribute type values, expected %d, got %d", 2, len(values))
		}
	}
}

func testCoreOrganization(t *testing.T) {
	err := AddOrganization(model.Organization{
		ID:          "Facebook",
		Name:        "Facebook",
		Street:      "Facebook road",
		Zip:         "Facebook zip",
		City:        "Facebook city",
		State:       "Facebook state",
		Country:     "Facebook country",
		Telephone:   "Facebook telephone",
		Email:       "Facebook email",
		URL:         "http://www.facebook.com",
		Description: "Facebook description",
	})
	if err != nil {
		t.Fatalf("failed to add organization: %s", err)
	}

	IDs, err := GetOrganizationIDs()
	if err != nil {
		t.Fatalf("failed to get organization IDs: %s", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("unexpected number of organization IDs, expected %d, got %d", 1, len(IDs))
	}
	for _, id := range IDs {
		org, errr := GetOrganization(id)
		if errr != nil {
			t.Fatalf("failed to get organization: %s", errr)
		}
		if org.Name != "Facebook" {
			t.Fatalf("invalid organization name, expected %s, got %s", "Facebook", org.Name)
		}
	}

	err = AddOrganization(model.Organization{
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
	})
	if err != nil {
		t.Fatalf("failed to add organization: %s", err)
	}

	IDs, err = GetOrganizationIDs()
	if err != nil {
		t.Fatalf("failed to get organization IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of organization IDs, expected %d, got %d", 2, len(IDs))
	}
}

func testCoreDisclosure(t *testing.T) {
	disc, err := model.MakeDisclosure(Self, "Facebook", time.Now().String(), "policyHuman", "policyMachine", "in a cubicle close to you", "")
	if err != nil {
		t.Fatalf("failed to make disclosure: %s", err)
	}
	err = AddDisclosure(disc)
	if err != nil {
		t.Fatalf("failed to add disclosure: %s", err)
	}

	IDs, err := GetDisclosureIDs()
	if err != nil {
		t.Fatalf("failed to get disclosure IDs: %s", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("unexpected number of disclosure IDs, expected %d, got %d", 1, len(IDs))
	}
	for _, id := range IDs {
		d, errr := GetDisclosure(id)
		if errr != nil {
			t.Fatalf("failed to get disclosre: %s", errr)
		}
		if d.Sender != disc.Sender {
			t.Fatalf("invalid disclosure name, expected %s, got %s", disc.Sender, d.Sender)
		}
	}

	IDs, err = GetDisclosureIDsToOrg("Facebook")
	if err != nil {
		t.Fatalf("failed to get disclosure IDs: %s", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("unexpected number of disclosure IDs, expected %d, got %d", 1, len(IDs))
	}

	disc2, err := model.MakeDisclosure(Self, "Facebook", time.Now().String(), "policyHuman", "policyMachine", "in a cubicle close to you", "")
	if err != nil {
		t.Fatalf("failed to make disclosure: %s", err)
	}
	err = AddDisclosure(disc2)
	if err != nil {
		t.Fatalf("failed to add disclosure: %s", err)
	}
	IDs, err = GetDisclosureIDs()
	if err != nil {
		t.Fatalf("failed to get disclosure IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of disclosure IDs, expected %d, got %d", 2, len(IDs))
	}
	IDs, err = GetDisclosureIDsToOrg("Facebook")
	if err != nil {
		t.Fatalf("failed to get disclosure IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of disclosure IDs, expected %d, got %d", 2, len(IDs))
	}
	IDs, err = GetDisclosureIDsChrono()
	if err != nil {
		t.Fatalf("failed to get disclosure IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of disclosure IDs, expected %d, got %d", 2, len(IDs))
	}
	if IDs[1] != disc2.ID {
		t.Fatalf("unexpected ID, expected %s, got %s", disc2.ID, IDs[0])
	}
	IDs, err = GetDisclosureIDsToOrgChrono("Facebook")
	if err != nil {
		t.Fatalf("failed to get disclosure IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of disclosure IDs, expected %d, got %d", 2, len(IDs))
	}
	if IDs[1] != disc2.ID {
		t.Fatalf("unexpected ID, expected %s, got %s", disc2.ID, IDs[0])
	}
}

func testCoreDisclosed(t *testing.T) {
	discIDs, _ := GetDisclosureIDs()
	attrIDs, _ := GetAttributeIDs()

	for _, id := range discIDs {
		err := AddDisclosed(model.Disclosed{
			Disclosure: id,
			Attribute:  attrIDs,
		})
		if err != nil {
			t.Fatalf("failed to add disclosed: %s", err)
		}
	}

	disclosed, err := GetDisclosed(discIDs[0])
	if err != nil {
		t.Fatalf("failed to get disclosed: %s", err)
	}
	if len(disclosed.Attribute) != len(attrIDs) {
		t.Fatalf("unexpected number of disclosed attributes, expected %d, got %d",
			len(attrIDs), len(disclosed.Attribute))
	}
}

func testExplicitlyDisclosed(t *testing.T) {
	IDs, err := GetExplicitDisclosureIDsToOrg("Facebook")
	if err != nil {
		t.Fatalf("failed to get explicit disclosures to org: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of IDs, expected %d, got %d", 2, len(IDs))
	}
}

func testExplicitlyDisclosedAttributes(t *testing.T) {
	allAttrIDs, _ := GetAttributeIDs()
	expAttrIDs, err := GetExplicitlyDisclosedAttributeIDs()
	if err != nil {
		t.Fatalf("failed to get explicitly disclosed attribute IDs: %s", err)
	}
	if len(allAttrIDs) != len(expAttrIDs) {
		t.Fatalf("unexpected number of IDs, expected %d, got %d", len(allAttrIDs), len(expAttrIDs))
	}
	expAttrIDsToOrg, err := GetExplicitlyDisclosedAttributeIDsToOrg("Facebook")
	if err != nil {
		t.Fatalf("failed to get explicitly disclosed attribute IDs: %s", err)
	}
	if len(expAttrIDsToOrg) != len(expAttrIDs) {
		t.Fatalf("unexpected number of IDs, expected %d, got %d", len(expAttrIDsToOrg), len(expAttrIDs))
	}
	expAttrTypeIDs, err := GetExplicitlyDisclosedAttributeTypeIDs()
	if err != nil {
		t.Fatalf("failed to get explicitly disclosed attribute type IDs: %s", err)
	}
	if len(expAttrTypeIDs) != 1 {
		t.Fatalf("unexpected number of IDs, expected %d, got %d", 1, len(expAttrTypeIDs))
	}

	expAttrTypeIDsToOrg, err := GetExplicitlyDisclosedAttributeTypeIDsToOrg("Facebook")
	if err != nil {
		t.Fatalf("failed to get explicitly disclosed attribute type IDs to org: %s", err)
	}
	if len(expAttrTypeIDsToOrg) != 1 {
		t.Fatalf("unexpected number of IDs, expected %d, got %d", 1, len(expAttrTypeIDs))
	}

	expAttrTypeValues, err := GetExplicitlyDisclosedAttributeTypeValues("music")
	if err != nil {
		t.Fatalf("failed to get explicitly disclosed attribute type values: %s", err)
	}
	if len(expAttrTypeValues) != 2 {
		t.Fatalf("unexpected number of IDs, expected %d, got %d", 2, len(expAttrTypeValues))
	}
	expAttrTypeValuesToOrg, err := GetExplicitlyDisclosedAttributeTypeValuesToOrg("Facebook", "music")
	if err != nil {
		t.Fatalf("failed to get explicitly disclosed attribute type values to org: %s", err)
	}
	if len(expAttrTypeValuesToOrg) != 2 {
		t.Fatalf("unexpected number of IDs, expected %d, got %d", 2, len(expAttrTypeValuesToOrg))
	}
}

func testImplicitDisclosures(t *testing.T) {
	// add explicit disclosure to Spotify of our attributes
	discExplicit, _ := model.MakeDisclosure(Self, "Spotify", time.Now().String(),
		"policyHuman", "policyMachine", "in a cubicle close to you", "")
	AddDisclosure(discExplicit)
	attrIDs, _ := GetAttributeIDs()
	AddDisclosed(model.Disclosed{
		Disclosure: discExplicit.ID,
		Attribute:  attrIDs,
	})

	IDs, err := GetImplicitDisclosureIDs(discExplicit.ID)
	if err != nil {
		t.Fatalf("failed to get implicit disclosure IDs: %s", err)
	}
	if len(IDs) > 0 {
		t.Fatalf("unexpected number of implicit disclosure IDs, expected %d, got %d", 0, len(IDs))
	}

	// add first implicit disclosure to Spotify
	attr, _ := model.MakeAttribute("Music Taste", "star", "confused")
	AddAttribute(attr)
	discImplicit, _ := model.MakeDisclosure("Spotify", "Spotify", time.Now().String(),
		"policyHuman", "policyMachine", "Spotify", "")
	AddDisclosure(discImplicit)
	AddDisclosed(model.Disclosed{
		Disclosure: discImplicit.ID,
		Attribute:  []string{attr.ID},
	})
	err = AddDownstream(model.Downstream{
		Origin: discExplicit.ID,
		Result: discImplicit.ID,
	})
	if err != nil {
		t.Fatalf("failed to add downstream disclosure: %s", err)
	}

	IDs, err = GetImplicitDisclosureIDs(discExplicit.ID)
	if err != nil {
		t.Fatalf("failed to get implicit disclosure IDs: %s", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("unexpected number of implicit disclosure IDs, expected %d, got %d", 1, len(IDs))
	}

	// add second implicit disclosure to Spotify
	attr2, _ := model.MakeAttribute("Coordinates", "location-arrow", "69.3111° N, 13.5333° E")
	AddAttribute(attr2)
	discImplicit2, _ := model.MakeDisclosure("Spotify", "Spotify", time.Now().String(),
		"policyHuman", "policyMachine", "Spotify", "")
	AddDisclosure(discImplicit2)
	AddDisclosed(model.Disclosed{
		Disclosure: discImplicit2.ID,
		Attribute:  []string{attr2.ID},
	})
	err = AddDownstream(model.Downstream{
		Origin: discExplicit.ID,
		Result: discImplicit2.ID,
	})
	if err != nil {
		t.Fatalf("failed to add downstream disclosure: %s", err)
	}

	IDs, err = GetImplicitDisclosureIDsChrono(discExplicit.ID)
	if err != nil {
		t.Fatalf("failed to get implicit disclosure IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of implicit disclosure IDs, expected %d, got %d", 2, len(IDs))
	}
	if IDs[0] != discImplicit.ID {
		t.Fatal("unexpected first chronologically sorted ID")
	}

	IDs, err = GetImplicitlyDisclosedAttributeIDs()
	if err != nil {
		t.Fatalf("failed to get implicitly disclosed attribute IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of implicitly disclosed attributes, expected %d, got %d", 2, len(IDs))
	}
	IDs, err = GetImplicitlyDisclosedAttributeIDsToOrg("Spotify")
	if err != nil {
		t.Fatalf("failed to get implicitly disclosed attribute IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of implicitly disclosed attributes, expected %d, got %d", 2, len(IDs))
	}

	IDs, err = GetImplicitlyDisclosedAttributeTypeIDs()
	if err != nil {
		t.Fatalf("failed to get implicitly disclosed attribute type IDs: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of implicitly disclosed attribute types, expected %d, got %d", 2, len(IDs))
	}
	IDs, err = GetImplicitlyDisclosedAttributeTypeIDsToOrg("Spotify")
	if err != nil {
		t.Fatalf("failed to get implicitly disclosed attribute type IDs to org: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of implicitly disclosed attribute types to org, expected %d, got %d", 2, len(IDs))
	}

	IDs, err = GetImplicitlyDisclosedAttributeTypeValues("location-arrow")
	if err != nil {
		t.Fatalf("failed to get implicitly disclosed attribute type IDs: %s", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("unexpected number of implicitly disclosed attribute values, expected %d, got %d", 1, len(IDs))
	}

	IDs, err = GetImplicitlyDisclosedAttributeTypeValuesToOrg("Spotify", "location-arrow")
	if err != nil {
		t.Fatalf("failed to get implicitly disclosed attribute type IDs to org: %s", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("unexpected number of implicitly disclosed attribute values, expected %d, got %d", 1, len(IDs))
	}
}

func testDownstreamDisclosures(t *testing.T) {
	IDs, _ := GetDisclosureIDsToOrg("Spotify")
	if len(IDs) != 1 {
		t.Fatalf("unexpected number of disclosures to Spotify, expected %d, got %d", 1, len(IDs))
	}
	spotifyID := IDs[0]
	IDs, err := GetDownstreamDisclosureIDs(spotifyID)
	if err != nil {
		t.Fatalf("failed to get downstream disclosures: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of downstream disclosures by Spotify, expected %d, got %d", 2, len(IDs))
	}
	IDs, err = GetDownstreamDisclosureIDsChrono(spotifyID)
	if err != nil {
		t.Fatalf("failed to get downstream disclosures: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected number of downstream disclosures by Spotify, expected %d, got %d", 2, len(IDs))
	}
}

func testCategories(t *testing.T) {
	err := AddCategory(model.Category{
		ID:   "test",
		Type: "testtype",
	})
	if err != nil {
		t.Fatalf("failed to add category: %s", err)
	}
	cat, err := GetCategory("test")
	if err != nil {
		t.Fatalf("failed to get category: %s", err)
	}
	if len(cat) != 1 {
		t.Fatalf("unexpected number of categories, expected %d, got %d", 1, len(cat))
	}

}

func testUser(t *testing.T) {
	err := SetUser(model.User{
		Name:    "Alice",
		Picture: "aliceselfie",
	})
	if err != nil {
		t.Fatalf("failed to set user: %s", err)
	}
	alice, err := GetUser()
	if err != nil {
		t.Fatalf("failed to get user: %s", err)
	}
	if alice.Name != "Alice" {
		t.Fatalf("unexpected name, expected %s, got %s", "Alice", alice.Name)
	}
	err = SetUser(model.User{
		Name:    "Bob",
		Picture: "bobselfie",
	})
	if err != nil {
		t.Fatalf("failed to set user: %s", err)
	}
	bob, err := GetUser()
	if err != nil {
		t.Fatalf("failed to get user: %s", err)
	}
	if bob.Name != "Bob" {
		t.Fatalf("unexpected name, expected %s, got %s", "Bob", bob.Name)
	}
}

func testExtended(t *testing.T) {
	attr, _ := GetAttributeIDs()
	if len(attr) == 0 {
		t.Fatalf("unexpected attribute ID count, expected at least %d, got %d", 1, len(attr))
	}

	IDs, err := GetReceivingOrgIDs(attr[0])
	if err != nil {
		t.Fatalf("failed to get receiving org IDs: %s", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("unexpected org ID count, expected at least %d, got %d", 1, len(IDs))
	}

	IDs, err = GetAttributeIDsToOrg("Spotify")
	if err != nil {
		t.Fatalf("failed to get attribute IDs to org: %s", err)
	}
	if len(IDs) != 2 {
		t.Fatalf("unexpected attribute IDs count, expected %d, got %d", 2, len(IDs))
	}

	IDs, err = GetAttributeTypeIDsToOrg("Spotify")
	if err != nil {
		t.Fatalf("failed to get attribute type IDs to org: %s", err)
	}
	if len(IDs) != 1 {
		t.Fatalf("unexpected attribute type IDs count, expected %d, got %d", 1, len(IDs))
	}
	values, err := GetAttributeTypeValuesToOrg("Spotify", IDs[0])
	if err != nil {
		t.Fatalf("failed to get attribute type values to org: %s", err)
	}
	if len(values) != 2 {
		t.Fatalf("unexpected values count, expected %d, got %d", 2, len(values))
	}
}

func testCoordinates(t *testing.T) {
	// the area we are testing around
	neLat := "59.428665021112984"
	neLng := "13.675060272216797"
	swLat := "59.37676383741801"
	swLng := "13.35336685180664"

	cords := make([]model.Coordinate, 2)
	// the first two cords are within the area
	cords[0] = model.MakeCoordinate("59.408665021112984", "13.425060272216797", "0", "0")
	cords[1] = model.MakeCoordinate("59.418665021112984", "13.525060272216797", "1", "1")

	wg := new(sync.WaitGroup)
	wg.Add(1)
	errChan := make(chan error, 1)
	AddCoordinates(cords, wg, errChan)
	close(errChan)
	for err := range errChan {
		t.Fatalf("failed to add coordinates: %s", err)
	}

	reply, err := GetCoordinates(neLat, neLng, swLat, swLng)
	if err != nil {
		t.Fatalf("failed to get coordinates: %s", err)
	}
	if len(reply) != 2 {
		t.Fatalf("expected %d coordinates as reply, got %d", 2, len(reply))
	}

	// we add two more coordinate outside our area of interest
	cords[0] = model.MakeCoordinate("60.418665021112984",
		"13.525060272216797", "2", "2")
	cords[1] = model.MakeCoordinate("59.418665021112984",
		"14.525060272216797", "3", "3")

	// Adding the coordinates make the test fail
	wg = new(sync.WaitGroup)
	wg.Add(1)
	errChan = make(chan error, 1)
	AddCoordinates(cords, wg, errChan)
	close(errChan)
	for err := range errChan {
		t.Fatalf("failed to add coordinates: %s", err)
	}

	reply, err = GetCoordinates(neLat, neLng, swLat, swLng)
	if err != nil {
		t.Fatalf("failed to get coordinates: %s", err)
	}
	if len(reply) != 2 {
		t.Fatalf("expected %d coordinates as reply, got %d", 2, len(reply))
	}
}
