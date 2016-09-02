package coordinate

import (
	"github.com/pylls/datatrack/config"
	"github.com/pylls/datatrack/handler"
)

const baseURL = config.APIURL + "/coordinate"
const coordURL = "/area/:neLat/:neLng/:swLat/:swLng"

type operation int

const (
	subset operation = iota
	reverse
	count
)

// Handlers contains all handlers for disclosures.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:   "coordinate attributes in a specific area",
		Method: "get",
		Url:    baseURL + coordURL,
		Handle: getCoordinates(false)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + coordURL + config.CountURL,
		Handle: getCoordinates(false, count)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + coordURL + config.RangeURL,
		Handle: getCoordinates(false, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + coordURL + config.ChronologicalURL,
		Handle: getCoordinates(true)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + coordURL + config.ChronologicalURL + config.RangeURL,
		Handle: getCoordinates(true, subset)},
}
