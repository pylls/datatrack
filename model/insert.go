package model

// Insert is an insert into the Data Track of new disclosures.
type Insert struct {
	User        User
	Disclosures []ActionDisclose
}

// ActionDisclose models one disclosure and its potential downstream
// disclosures that occurred as a consequence of the disclosure.
type ActionDisclose struct {
	Disclosure Disclosure
	Attribute  []Attribute
	Sender     Organization // if empty on first level = user disclosed it
	Recipient  Organization
	Downstream []ActionDisclose
}
