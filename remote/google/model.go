package google

import "github.com/pylls/datatrack/model"

// WatchHistory is a history of watched videos.
type WatchHistory struct {
	Videos []Video
}

// Video is a watched video.
type Video struct {
	ContentDetails ContentDetails `json:"contentDetails"`
	ETag           string         `json:"etag"`
	ID             string         `json:"id"`
	Kind           string         `json:"kind"`
	Snippet        Snippet        `json:"snippet"`
	Status         Status         `json:"status"`
}

// ContentDetails ahs details about a video.
type ContentDetails struct {
	VideoID string `json:"videoId"`
}

// Snippet is a snippet.
type Snippet struct {
	ChannelID    string     `json:"channelId"`
	ChannelTitle string     `json:"channelTitle"`
	Description  string     `json:"description"`
	PlaylistID   string     `json:"playlistId"`
	Position     int        `json:"position"`
	PublishedAt  string     `json:"publishedAt"`
	ResourceID   ResourceID `json:"resourceId"`
	Thumbnails   Thumbnails `json:"thumbnails"`
	Title        string     `json:"title"`
}

// ResourceID is a resource.
type ResourceID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

// Thumbnails is a collection of thumbnails.
type Thumbnails struct {
	Default  Thumbnail `json:"default"`
	High     Thumbnail `json:"high"`
	MaxRes   Thumbnail `json:"maxres"`
	Medium   Thumbnail `json:"medium"`
	Standard Thumbnail `json:"standard"`
}

// Thumbnail is a video thumbnail.
type Thumbnail struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

// Status is a privacy status.
type Status struct {
	PrivacyStatus string `json:"privacyStatus"`
}

// ParsedVideo is a parsed video.
type ParsedVideo struct {
	Disclosure model.Disclosure
	Attributes []model.Attribute
	Disclosed  model.Disclosed
}

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
	// NOTE: typo above for "activitys" is as-is in Google's format, do not change
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
