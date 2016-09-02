package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var activeHandlers Handlers

type Handler struct {
	Name        string
	Method      string
	URL         string
	Description string
	Handle      func(web.C, http.ResponseWriter, *http.Request)
}

type Handlers []Handler

func (h Handler) Register() error {
	switch strings.ToLower(h.Method) {
	case "options":
		goji.Options(h.URL, h.Handle)
		goji.Options(h.URL+"/", h.Handle)
	case "get":
		goji.Get(h.URL, h.Handle)
		goji.Get(h.URL+"/", h.Handle)
	case "head":
		goji.Head(h.URL, h.Handle)
		goji.Head(h.URL+"/", h.Handle)
	case "post":
		goji.Post(h.URL, h.Handle)
		goji.Post(h.URL+"/", h.Handle)
	case "put":
		goji.Put(h.URL, h.Handle)
		goji.Put(h.URL+"/", h.Handle)
	case "delete":
		goji.Delete(h.URL, h.Handle)
		goji.Delete(h.URL+"/", h.Handle)
	case "trace":
		goji.Trace(h.URL, h.Handle)
		goji.Trace(h.URL+"/", h.Handle)
	case "connect":
		goji.Connect(h.URL, h.Handle)
		goji.Connect(h.URL+"/", h.Handle)
	// non-HTTP 1.1
	case "patch":
		goji.Patch(h.URL, h.Handle)
		goji.Patch(h.URL+"/", h.Handle)
	default:
		return errors.New("unsupported method")
	}
	activeHandlers = append(activeHandlers, h)
	return nil
}

func (hs Handlers) Register() error {
	for _, handler := range hs {
		if err := handler.Register(); err != nil {
			return err
		}
	}
	return nil
}

func Concat(hss ...Handlers) Handlers {
	var hs Handlers
	for _, hxs := range hss {
		hs = append(hs, hxs...)
	}
	return hs
}
