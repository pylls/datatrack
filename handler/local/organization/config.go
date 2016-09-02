package organization

import (
	"github.com/marcelfarres/datatrack/config"
	"github.com/marcelfarres/datatrack/handler"
)

const baseURL = config.APIURL + "/organization"
const withID = "/:organizationId"
const attributeURL = "/receivedAttribute/:attributeId"

type mode int

const (
	organization mode = iota
	attribute
)

type operation int

const (
	subset operation = iota
	reverse
	count
)

// Handlers contains handlers.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:   "local organization data",
		Method: "get",
		URL:    baseURL,
		Handle: orgHandler(organization)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + config.CountURL,
		Handle: orgHandler(organization, count)},
	handler.Handler{
		Method:      "get",
		URL:         baseURL + config.RangeURL,
		Description: "retrieve all organization ids",
		Handle:      orgHandler(organization, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + attributeURL,
		Handle: orgHandler(attribute)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + attributeURL + config.CountURL,
		Handle: orgHandler(attribute, count)},
	handler.Handler{
		Method:      "get",
		URL:         baseURL + attributeURL + config.RangeURL,
		Description: "retrieve all organization ids that received :attributeId",
		Handle:      orgHandler(attribute, subset)},
	handler.Handler{
		Method:      "get",
		URL:         baseURL + withID,
		Description: "retrieve organization",
		Handle:      detailsHandler}}
