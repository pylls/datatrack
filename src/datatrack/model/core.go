// Package model contains the type definitions that make up the Data Track
// disclosure model, with associated functions for creating new objects.
package model

// Attribute is an attribute.
type Attribute struct {
	ID    string
	Name  string
	Type  string
	Value string
}

// Category is a category.
type Category struct {
	ID   string
	Type string
}

// Disclosed represents the disclosure of several attributes.
type Disclosed struct {
	Disclosure string
	Attribute  []string
}

// Disclosure is a data disclosure.
type Disclosure struct {
	ID                   string
	Sender               string
	Recipient            string
	Timestamp            string
	PrivacyPolicyHuman   string
	PrivacyPolicyMachine string
	DataLocation         string
	API                  string
}

// Downstream represents a relationship between an origin data disclosure and
// a resulting data disclosure.
type Downstream struct {
	Origin string
	Result string
}

// Organization is an organization.
type Organization struct {
	ID          string // popular name
	Name        string // legal name
	Street      string
	Zip         string
	City        string
	State       string
	Country     string
	Telephone   string
	Email       string
	URL         string
	Description string
}

// User represents a user.
type User struct {
	Name    string
	Picture string
}
