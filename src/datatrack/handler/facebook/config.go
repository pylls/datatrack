package facebook

import (
	"datatrack/config"
	"datatrack/handler"
)

const baseURL = config.APIURL + "/facebook"

// Handlers are handlers.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:        "Facebook",
		Method:      "post",
		Url:         baseURL,
		Description: "adds the data from Facebook (zip file)",
		Handle:      facebookHandler}}
