package local

import (
	"github.com/pylls/datatrack/handler"
	"github.com/pylls/datatrack/handler/local/attribute"
	"github.com/pylls/datatrack/handler/local/coordinate"
	"github.com/pylls/datatrack/handler/local/disclosure"
	"github.com/pylls/datatrack/handler/local/organization"
)

// Handlers are handlers.
var Handlers = handler.Concat(
	disclosure.Handlers,
	coordinate.Handlers,
	attribute.Handlers,
	organization.Handlers)
