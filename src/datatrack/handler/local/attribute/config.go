package attribute

import (
	"datatrack/config"
	"datatrack/handler"
)

const baseURL = config.APIURL + "/attribute"
const withID = "/:attributeId"
const withType = "/type"
const withTypeID = "/:typeId/value"
const organizationURL = "/toOrganization/:organizationId"

type mode int

const (
	attribute mode = iota
	explicitAttribute
	implicitAttribute
	organization
	explicitOrganization
	implicitOrganization
)

type operation int

const (
	subset operation = iota
	reverse
	count
)

type field int

const (
	thetype = iota
	thevalue
)

// Handlers contains all handlers for attributes.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:   "local attribute data",
		Method: "get",
		Url:    baseURL,
		Handle: attributeHandler(attribute)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.CountURL,
		Handle: attributeHandler(attribute, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.RangeURL,
		Description: "retrieve all attribute ids",
		Handle:      attributeHandler(attribute, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL,
		Handle: attributeHandler(explicitAttribute)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + config.CountURL,
		Handle: attributeHandler(explicitAttribute, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ExplicitURL + config.RangeURL,
		Description: "retrieve all explicit attribute ids",
		Handle:      attributeHandler(explicitAttribute, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL,
		Handle: attributeHandler(implicitAttribute)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + config.CountURL,
		Handle: attributeHandler(implicitAttribute, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ImplicitURL + config.RangeURL,
		Description: "retrieve all implicit attribute ids",
		Handle:      attributeHandler(implicitAttribute, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + organizationURL,
		Handle: attributeHandler(organization)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + organizationURL + config.CountURL,
		Handle: attributeHandler(organization, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + organizationURL + config.RangeURL,
		Description: "retrieve all attribute ids that have been disclosed to :organizationId",
		Handle:      attributeHandler(organization, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + organizationURL,
		Handle: attributeHandler(explicitOrganization)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + organizationURL + config.CountURL,
		Handle: attributeHandler(explicitOrganization, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ExplicitURL + organizationURL + config.RangeURL,
		Description: "retrieve all explicit attribute ids that have been disclosed to :organizationId",
		Handle:      attributeHandler(explicitOrganization, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + organizationURL,
		Handle: attributeHandler(implicitOrganization)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + organizationURL + config.CountURL,
		Handle: attributeHandler(implicitOrganization, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ImplicitURL + organizationURL + config.RangeURL,
		Description: "retrieve all implicit attribute ids that have been disclosed to :organizationId",
		Handle:      attributeHandler(implicitOrganization, subset)},

	handler.Handler{
		Method: "get",
		Url:    baseURL + withType,
		Handle: attributeFieldHandler(attribute, thetype)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + withType + config.CountURL,
		Handle: attributeFieldHandler(attribute, thetype, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + withType + config.RangeURL,
		Description: "retrieve all attribute types",
		Handle:      attributeFieldHandler(attribute, thetype, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + withType + withTypeID,
		Handle: attributeFieldHandler(attribute, thevalue)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + withType + withTypeID + config.CountURL,
		Handle: attributeFieldHandler(attribute, thevalue, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + withType + withTypeID + config.RangeURL,
		Description: "retrieve all attribute values",
		Handle:      attributeFieldHandler(attribute, thevalue, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + withType,
		Handle: attributeFieldHandler(explicitAttribute, thetype)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + withType + config.CountURL,
		Handle: attributeFieldHandler(explicitAttribute, thetype, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ExplicitURL + withType + config.RangeURL,
		Description: "retrieve explicit attribute types",
		Handle:      attributeFieldHandler(explicitAttribute, thetype, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + withType + withTypeID,
		Handle: attributeFieldHandler(explicitAttribute, thevalue)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + withType + withTypeID + config.CountURL,
		Handle: attributeFieldHandler(explicitAttribute, thevalue, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ExplicitURL + withType + withTypeID + config.RangeURL,
		Description: "retrieve explicit attribute values",
		Handle:      attributeFieldHandler(explicitAttribute, thevalue, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + withType,
		Handle: attributeFieldHandler(implicitAttribute, thetype)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + withType + config.CountURL,
		Handle: attributeFieldHandler(implicitAttribute, thetype, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ImplicitURL + withType + config.RangeURL,
		Description: "retrieve all implicit attribute types",
		Handle:      attributeFieldHandler(implicitAttribute, thetype, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + withType + withTypeID,
		Handle: attributeFieldHandler(implicitAttribute, thevalue)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + withType + withTypeID + config.CountURL,
		Handle: attributeFieldHandler(implicitAttribute, thevalue, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ImplicitURL + withType + withTypeID + config.RangeURL,
		Description: "retrieve all implicit attribute values",
		Handle:      attributeFieldHandler(implicitAttribute, thevalue, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + organizationURL + withType,
		Handle: attributeFieldHandler(organization, thetype)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + organizationURL + withType + config.CountURL,
		Handle: attributeFieldHandler(organization, thetype, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + organizationURL + withType + config.RangeURL,
		Description: "retrieve all attribute types that have been disclosed to :organizationId",
		Handle:      attributeFieldHandler(organization, thetype, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + organizationURL + withType + withTypeID,
		Handle: attributeFieldHandler(organization, thevalue)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + organizationURL + withType + withTypeID + config.CountURL,
		Handle: attributeFieldHandler(organization, thevalue, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + organizationURL + withType + withTypeID + config.RangeURL,
		Description: "retrieve all attribute values that have been disclosed to :organizationId",
		Handle:      attributeFieldHandler(organization, thevalue, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + organizationURL + withType,
		Handle: attributeFieldHandler(explicitOrganization, thetype)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + organizationURL + withType + config.CountURL,
		Handle: attributeFieldHandler(explicitOrganization, thetype, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ExplicitURL + organizationURL + withType + config.RangeURL,
		Description: "retrieve explicit attribute types that have been disclosed to :organizationId",
		Handle:      attributeFieldHandler(explicitOrganization, thetype, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + organizationURL + withType + withTypeID,
		Handle: attributeFieldHandler(explicitOrganization, thevalue)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ExplicitURL + organizationURL + withType + withTypeID + config.CountURL,
		Handle: attributeFieldHandler(explicitOrganization, thevalue, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ExplicitURL + organizationURL + withType + withTypeID + config.RangeURL,
		Description: "retrieve explicit attribute values that have been disclosed to :organizationId",
		Handle:      attributeFieldHandler(explicitOrganization, thevalue, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + organizationURL + withType,
		Handle: attributeFieldHandler(implicitOrganization, thetype)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + organizationURL + withType + config.CountURL,
		Handle: attributeFieldHandler(implicitOrganization, thetype, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ImplicitURL + organizationURL + withType + config.RangeURL,
		Description: "retrieve implicit attribute types that have been disclosed to :organizationId",
		Handle:      attributeFieldHandler(implicitOrganization, thetype, subset)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + organizationURL + withType + withTypeID,
		Handle: attributeFieldHandler(implicitOrganization, thevalue)},
	handler.Handler{
		Method: "get",
		Url:    baseURL + config.ImplicitURL + organizationURL + withType + withTypeID + config.CountURL,
		Handle: attributeFieldHandler(implicitOrganization, thevalue, count)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + config.ImplicitURL + organizationURL + withType + withTypeID + config.RangeURL,
		Description: "retrieve implicit attribute values that have been disclosed to :organizationId",
		Handle:      attributeFieldHandler(implicitOrganization, thevalue, subset)},
	handler.Handler{
		Method:      "get",
		Url:         baseURL + withID,
		Description: "retrieve attribute",
		Handle:      detailsHandler},
}
