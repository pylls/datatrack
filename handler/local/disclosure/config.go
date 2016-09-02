package disclosure

import (
	"github.com/marcelfarres/datatrack/config"
	"github.com/marcelfarres/datatrack/handler"
)

const baseURL = config.APIURL + "/disclosure"
const withID = "/:disclosureId"
const organizationURL = "/toOrganization/:organizationId"
const attributeURL = "/attribute"
const downstreamURL = "/downstream"

type mode int

const (
	disclosure mode = iota
	disclosureChrono
	organization
	organizationChrono
	attribute
	implicit
	implicitChrono
	downstream
	downstreamChrono
)

type operation int

const (
	subset operation = iota
	reverse
	count
)

// Handlers contains all handlers for disclosures.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:   "local disclosure data",
		Method: "get",
		URL:    baseURL,
		Handle: disclosureHandler(disclosure)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + config.CountURL,
		Handle: disclosureHandler(disclosure, count)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + config.RangeURL,
		Handle: disclosureHandler(disclosure, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + config.ChronologicalURL,
		Handle: disclosureHandler(disclosureChrono)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + config.ChronologicalURL + config.ReverseURL,
		Handle: disclosureHandler(disclosureChrono, reverse)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + config.ChronologicalURL + config.RangeURL,
		Handle: disclosureHandler(disclosureChrono, subset)},
	handler.Handler{
		Method:      "get",
		URL:         baseURL + config.ChronologicalURL + config.ReverseURL + config.RangeURL,
		Description: "retrieve all data disclosure ids",
		Handle:      disclosureHandler(disclosureChrono, reverse, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + organizationURL,
		Handle: disclosureHandler(organization)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + organizationURL + config.CountURL,
		Handle: disclosureHandler(organization, count)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + organizationURL + config.RangeURL,
		Handle: disclosureHandler(organization, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + organizationURL + config.ChronologicalURL,
		Handle: disclosureHandler(organizationChrono)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + organizationURL + config.ChronologicalURL + config.ReverseURL,
		Handle: disclosureHandler(organizationChrono, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + organizationURL + config.ChronologicalURL + config.RangeURL,
		Handle: disclosureHandler(organizationChrono, subset)},
	handler.Handler{
		Method:      "get",
		URL:         baseURL + organizationURL + config.ChronologicalURL + config.ReverseURL + config.RangeURL,
		Description: "retrieve all data disclosure ids of disclosure to :organizationId",
		Handle:      disclosureHandler(organizationChrono, reverse, subset)},
	handler.Handler{
		Method:      "get",
		URL:         baseURL + withID,
		Description: "retrieve data disclosure",
		Handle:      detailsHandler},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + attributeURL,
		Handle: disclosureHandler(attribute)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + attributeURL + config.CountURL,
		Handle: disclosureHandler(attribute, count)},
	handler.Handler{
		Method:      "get",
		URL:         baseURL + withID + attributeURL + config.RangeURL,
		Description: "retrieve all attribute ids of the data disclosure",
		Handle:      disclosureHandler(attribute, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + config.ImplicitURL,
		Handle: disclosureHandler(implicit)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + config.ImplicitURL + config.CountURL,
		Handle: disclosureHandler(implicit, count)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + config.ImplicitURL + config.RangeURL,
		Handle: disclosureHandler(implicit, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + config.ImplicitURL + config.ChronologicalURL,
		Handle: disclosureHandler(implicitChrono)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + config.ImplicitURL + config.ChronologicalURL + config.ReverseURL,
		Handle: disclosureHandler(implicitChrono, reverse)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + config.ImplicitURL + config.ChronologicalURL + config.RangeURL,
		Handle: disclosureHandler(implicitChrono, subset)},
	handler.Handler{
		Method:      "get",
		URL:         baseURL + withID + config.ImplicitURL + config.ChronologicalURL + config.ReverseURL + config.RangeURL,
		Description: "retrieve all implicit data disclosure ids",
		Handle:      disclosureHandler(implicitChrono, reverse, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + downstreamURL,
		Handle: disclosureHandler(downstream)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + downstreamURL + config.CountURL,
		Handle: disclosureHandler(downstream, count)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + downstreamURL + config.RangeURL,
		Handle: disclosureHandler(downstream, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + downstreamURL + config.ChronologicalURL,
		Handle: disclosureHandler(downstreamChrono)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + downstreamURL + config.ChronologicalURL + config.ReverseURL,
		Handle: disclosureHandler(downstreamChrono, subset)},
	handler.Handler{
		Method: "get",
		URL:    baseURL + withID + downstreamURL + config.ChronologicalURL + config.RangeURL,
		Handle: disclosureHandler(downstreamChrono, subset)},
	handler.Handler{
		Method:      "get",
		URL:         baseURL + withID + downstreamURL + config.ChronologicalURL + config.ReverseURL + config.RangeURL,
		Description: "retrieve all downstream data disclosure ids",
		Handle:      disclosureHandler(downstreamChrono, reverse, subset)},
}
