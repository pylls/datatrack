package locationhistory

import "datatrack/model"

// LocationHistory is the top struct for Google takeout location history.
type LocationHistory struct {
	Locations []Location `json:"locations"`
}

// Location is a location.
type Location struct {
	TimeStampMs     string          `json:"timestampMs"`
	LatitudeE7      int64           `json:"latitudeE7"`
	LongitudeE7     int64           `json:"longitudeE7"`
	Accuracy        int64           `json:"accuracy"`
	Velocity        int64           `json:"velocity"`
	Heading         int             `json:"heading"`  // angle
	Altitude        int             `json:"altitude"` // meter over sea level
	OuterActivities []OuterActivity `json:"activitys"`
}

// OuterActivity is many InnerActivity at a particular point in time.
type OuterActivity struct {
	OTimeStampMs    string          `json:"timestampMs"`
	InnerActivities []InnerActivity `json:"activities"`
}

// InnerActivity is a type of activity with a given confidence (from Google's PoV).
type InnerActivity struct {
	Type       string `json:"type"`
	Confidence int    `json:"confidence"` // percent
}

var org = model.Organization{
	ID:   "Google",
	Name: "Google"}

// ParsedLocation is Location translated to the DT model.
type ParsedLocation struct {
	Position   Position
	Activities []Activity
}

// Position is a position parsed to the DT model.
type Position struct {
	Disclosure model.Disclosure
	Attributes []model.Attribute
	Coordinate model.Coordinate
	Disclosed  model.Disclosed
}

// Activity is an activity parsed to the DT model.
type Activity struct {
	Disclosure model.Disclosure
	Downstream model.Downstream
	Attributes []model.Attribute
	Disclosed  model.Disclosed
}
