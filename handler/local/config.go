package local

import (
	"github.com/marcelfarres/datatrack/handler"
	"github.com/marcelfarres/datatrack/handler/local/attribute"
	"github.com/marcelfarres/datatrack/handler/local/coordinate"
	"github.com/marcelfarres/datatrack/handler/local/disclosure"
	"github.com/marcelfarres/datatrack/handler/local/organization"
)

// Handlers are handlers.
var Handlers = handler.Concat(
	disclosure.Handlers,
	coordinate.Handlers,
	attribute.Handlers,
	organization.Handlers)
