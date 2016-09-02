package user

import (
	"github.com/marcelfarres/datatrack/config"
	"github.com/marcelfarres/datatrack/handler"
)

const baseURL = config.APIURL + "/user"

// Handlers are handlers.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:        "user",
		Method:      "get",
		Url:         baseURL,
		Description: "retrieve user name and picture",
		Handle:      userHandler},
	handler.Handler{
		Method:      "put",
		Url:         baseURL,
		Description: "create or update user entry",
		Handle:      updateUserHandler}}
