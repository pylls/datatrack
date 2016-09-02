package user

import (
	"github.com/pylls/datatrack/config"
	"github.com/pylls/datatrack/handler"
)

const baseURL = config.APIURL + "/user"

// Handlers are handlers.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:        "user",
		Method:      "get",
		URL:         baseURL,
		Description: "retrieve user name and picture",
		Handle:      userHandler},
	handler.Handler{
		Method:      "put",
		URL:         baseURL,
		Description: "create or update user entry",
		Handle:      updateUserHandler}}
