package local

import (
	"datatrack/handler"
	"datatrack/handler/local/attribute"
	"datatrack/handler/local/coordinate"
	"datatrack/handler/local/disclosure"
	"datatrack/handler/local/organization"
)

// Handlers are handlers.
var Handlers = handler.Concat(
	disclosure.Handlers,
	coordinate.Handlers,
	attribute.Handlers,
	organization.Handlers)
