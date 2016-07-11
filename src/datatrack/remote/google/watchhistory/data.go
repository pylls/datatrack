package watchhistory

import "datatrack/model"

type WatchHistory struct {
	Videos []Video
}

type Video struct {
	ContentDetails ContentDetails `json:"contentDetails"`
	Etag           string         `json:"etag"`
	Id             string         `json:"id"`
	Kind           string         `json:"kind"`
	Snippet        Snippet        `json:"snippet"`
	Status         Status         `json:"status"`
}

type ContentDetails struct {
	VideoId string `json:"videoId"`
}

type Snippet struct {
	ChannelId    string     `json:"channelId"`
	ChannelTitle string     `json:"channelTitle"`
	Description  string     `json:"description"`
	PlaylistId   string     `json:"playlistId"`
	Position     int        `json:"position"`
	PublishedAt  string     `json:"publishedAt"`
	ResourceId   ResourceId `json:"resourceId"`
	Thumbnails   Thumbnails `json:"thumbnails"`
	Title        string     `json:"title"`
}

type ResourceId struct {
	Kind    string `json:"kind"`
	VideoId string `json:"videoId"`
}

type Thumbnails struct {
	Default  Thumbnail `json:"default"`
	High     Thumbnail `json:"high"`
	Maxres   Thumbnail `json:"maxres"`
	Medium   Thumbnail `json:"medium"`
	Standard Thumbnail `json:"standard"`
}

type Thumbnail struct {
	Height int    `json:"height"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
}

type Status struct {
	PrivacyStatus string `json:"privacyStatus"`
}

var org = model.Organization{
	ID:   "Google",
	Name: "Google",
}


type ParsedVideo struct {
	Disclosure model.Disclosure
	Attributes []model.Attribute
	Disclosed  model.Disclosed
}
