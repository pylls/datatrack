package model

import (
	"datatrack/util"
	"encoding/hex"
)

// MakeAttribute creates an attribute.
func MakeAttribute(name, thetype, value string) (Attribute, error) {
	id := util.Hash([]byte(name), []byte(thetype), []byte(value))[:8]
	return Attribute{
		ID:    hex.EncodeToString(id),
		Name:  name,
		Type:  thetype,
		Value: value}, nil
}

// MakeDisclosure creates a Disclosure from provided data.
func MakeDisclosure(sender string, recipient string, time string, policyHuman string,
	policyMachine string, location string, API string) (Disclosure, error) {
	id := util.Hash([]byte(sender), []byte(recipient), []byte(time))[:8]

	return Disclosure{
		ID:                   hex.EncodeToString(id),
		Sender:               sender,
		Recipient:            recipient,
		Timestamp:            time,
		PrivacyPolicyHuman:   policyHuman,
		PrivacyPolicyMachine: policyMachine,
		DataLocation:         location,
		API:                  API}, nil
}

// MakeCoordinate creates a Coordinate.
func MakeCoordinate(latitude, longitude, disclosureID, timestamp string) Coordinate {
	id := util.Hash([]byte(latitude), []byte(longitude), []byte(timestamp))[:8]

	return Coordinate{
		ID:           hex.EncodeToString(id),
		Latitude:     latitude,
		Longitude:    longitude,
		DisclosureID: disclosureID,
		Timestamp:    timestamp,
	}
}
