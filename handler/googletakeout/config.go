package googletakeout

import (
	"github.com/pylls/datatrack/config"
	"github.com/pylls/datatrack/handler"
)

const baseURL = config.APIURL + "/google"

// Handlers are handlers.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:        "Google Takeout",
		Method:      "post",
		URL:         baseURL,
		Description: "adds the data from a Google Takeout (zip file)",
		Handle:      takeoutHandler}}
