package coordinate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/pylls/datatrack/database"
	"github.com/pylls/datatrack/model"
	"github.com/zenazn/goji/web"
)

type sortCoordinates []coordinateReply

func (a sortCoordinates) Len() int      { return len(a) }
func (a sortCoordinates) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortCoordinates) Less(i, j int) bool {
	return strings.Compare(a[i].Timestamp, a[j].Timestamp) == -1
}

type coordinateReply struct {
	model.Coordinate
	Next model.Coordinate
	Prev model.Coordinate
}

func getCoordinates(sortChrono bool, op ...operation) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		// get coordinates within the area
		// ASSUMPTION: this will get us a manageable number of coordinates, so later operations can be less efficient
		coords, err := database.GetCoordinates(c.URLParams["neLat"], c.URLParams["neLng"], c.URLParams["swLat"], c.URLParams["swLng"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// create reply struct
		reply := make([]coordinateReply, len(coords))
		for i := 0; i < len(coords); i++ {
			reply[i].DisclosureID = coords[i].DisclosureID
			reply[i].ID = coords[i].ID
			reply[i].Latitude = coords[i].Latitude
			reply[i].Longitude = coords[i].Longitude
			reply[i].Timestamp = coords[i].Timestamp
		}

		if sortChrono {
			// first we sort by time
			sort.Sort(sortCoordinates(reply))

			// then we look for coordinates in the path outside of the area
			includePrev := true
			for i := 0; i < len(reply); i++ {
				// include the previous (by time) coordinate if needed
				if includePrev {
					exists, prev, err := database.GetPrevCoordinateChrono(reply[i].ID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					if exists {
						if !inArea(prev, c.URLParams["neLat"], c.URLParams["neLng"], c.URLParams["swLat"], c.URLParams["swLng"]) {
							reply[i].Prev = prev
						}
					}
					includePrev = false
				}

				exists, next, err := database.GetNextCoordinateChrono(reply[i].ID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if exists {
					// the next coordinate (by time) is not in the area, add it and flag that the next coordinate (in our reply)
					// has to get its previous coordinate (by time) included
					if !inArea(next, c.URLParams["neLat"], c.URLParams["neLng"], c.URLParams["swLat"], c.URLParams["swLng"]) {
						reply[i].Next = next
						includePrev = true
					}
				}
			}
		}

		// go over each operation
		countOutput := false
		for i := 0; i < len(op); i++ {
			switch op[i] {
			case subset:
				first, err := strconv.Atoi(c.URLParams["first"])
				if err != nil {
					panic(err)
				}
				last, err := strconv.Atoi(c.URLParams["last"])
				if err != nil {
					panic(err)
				}
				if first < 0 || len(reply) <= last || last > first {
					http.Error(w, "invalid range", http.StatusBadRequest)
					return
				}
				reply = reply[first:last]
			case reverse:
				for i, j := 0, len(reply)-1; i < j; i, j = i+1, j-1 {
					reply[i], reply[j] = reply[j], reply[i]
				}
			case count:
				countOutput = true
			}
		}

		if countOutput {
			fmt.Fprintf(w, "%d", len(reply))
			return
		}

		j, err := json.Marshal(reply)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-type", "application/json")
		fmt.Fprintf(w, "%s", j)
	}
}

func inArea(next model.Coordinate, neLat, neLng, swLat, swLng string) bool {
	minLng := []byte(database.PadCoordinate(swLng))
	maxLng := []byte(database.PadCoordinate(neLng))
	minLat := []byte(database.PadCoordinate(swLat))
	maxLat := []byte(database.PadCoordinate(neLat))
	targetLng := []byte(next.Longitude)
	targetLat := []byte(next.Latitude)

	// minLng <= nxtLng <= maxLng && minLat <= nxtLat <= maxLat
	return bytes.Compare(minLng, targetLng) <= 0 && bytes.Compare(targetLng, maxLng) <= 0 &&
		bytes.Compare(minLat, targetLat) <= 0 && bytes.Compare(targetLat, maxLat) <= 0
}
