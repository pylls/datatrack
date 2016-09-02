package facebook

import "time"

type LatLong struct {
	Latitude  float64
	Longitude float64
}
type ActDate struct {
	Act       string
	Date      time.Time
	UserAgent string
}
type IPAddr string
