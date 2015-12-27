package handler

import (
	"datatrack/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
)

var CommonHandlers = Handlers{
	Handler{
		Name:        "API reference",
		Method:      "get",
		Url:         config.APIURL,
		Description: "print the API reference",
		Handle:      getHandler}}

func getHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	html := "<html><head>"
	html += "<title>API reference</title>"
	html += "</head><body>"
	html += "<table>"
	html += toHtmlTable(activeHandlers)
	html += "</table>"
	html += "</body></html>"
	fmt.Fprintf(w, "%s", html)
}

func toHtmlTable(hs Handlers) string {
	var html string
	for _, h := range hs {
		var format string
		if h.Handle == nil {
			format = "<tr><td><b>%s<b></td><td><tt>[%s]</tt></td><td><font color=\"red\"><b><tt>%s</tt></b></color></td></tr><tr><td></td><td colspan=\"2\"><i>%s</i></td></tr>"
		} else {
			format = "<tr><td><b>%s<b></td><td><tt>[%s]</tt></td><td><font color=\"blue\"><b><tt>%s</tt></b></color></td></tr><tr><td></td><td colspan=\"2\"><i>%s</i></td></tr>"
		}
		html += fmt.Sprintf(format,
			h.Name, strings.ToUpper(h.Method), h.Url, h.Description)
	}
	return html
}
