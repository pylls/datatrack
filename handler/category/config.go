package category

import (
	"github.com/marcelfarres/datatrack/config"
	"github.com/marcelfarres/datatrack/handler"
)

const baseURL = config.APIURL + "/type"
const withCategory = "/:categoryId"

// Handlers contain handlers.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:        "attribute type categories",
		Method:      "get",
		URL:         baseURL + withCategory,
		Description: "retrieve all types belonging to a (sub)category",
		Handle:      categoryHandler}}
