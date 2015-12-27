package googletakeout

import (
	"datatrack/config"
	"datatrack/handler"
)

const baseURL = config.APIURL + "/google"

// Handlers are handlers.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:        "Google Takeout",
		Method:      "post",
		Url:         baseURL,
		Description: "adds the data from a Google Takeout (zip file)",
		Handle:      takeoutHandler}}
