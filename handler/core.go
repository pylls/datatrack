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
	Url         string
	Description string
	Handle      func(web.C, http.ResponseWriter, *http.Request)
}

func (h Handler) Register() error {
	switch strings.ToLower(h.Method) {
	case "options":
		goji.Options(h.Url, h.Handle)
		goji.Options(h.Url+"/", h.Handle)
	case "get":
		goji.Get(h.Url, h.Handle)
		goji.Get(h.Url+"/", h.Handle)
	case "head":
		goji.Head(h.Url, h.Handle)
		goji.Head(h.Url+"/", h.Handle)
	case "post":
		goji.Post(h.Url, h.Handle)
		goji.Post(h.Url+"/", h.Handle)
	case "put":
		goji.Put(h.Url, h.Handle)
		goji.Put(h.Url+"/", h.Handle)
	case "delete":
		goji.Delete(h.Url, h.Handle)
		goji.Delete(h.Url+"/", h.Handle)
	case "trace":
		goji.Trace(h.Url, h.Handle)
		goji.Trace(h.Url+"/", h.Handle)
	case "connect":
		goji.Connect(h.Url, h.Handle)
		goji.Connect(h.Url+"/", h.Handle)
	// non-HTTP 1.1
	case "patch":
		goji.Patch(h.Url, h.Handle)
		goji.Patch(h.Url+"/", h.Handle)
	default:
		return errors.New("unsupported method")
	}
	activeHandlers = append(activeHandlers, h)
	return nil
}

type Handlers []Handler

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
