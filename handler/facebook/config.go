package facebook

import (
	"github.com/pylls/datatrack/config"
	"github.com/pylls/datatrack/handler"
)

const baseURL = config.APIURL + "/facebook"

// Handlers are handlers.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:        "Facebook",
		Method:      "post",
		URL:         baseURL,
		Description: "adds the data from Facebook (zip file)",
		Handle:      facebookHandler}}
