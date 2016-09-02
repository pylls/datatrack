package testdata

import (
	"github.com/marcelfarres/datatrack/config"
	"github.com/marcelfarres/datatrack/handler"
)

const baseAdbokis = config.APIURL + "/testadbokis"
const baseTestdata = config.APIURL + "/testdata"

// Handlers are handlers.
var Handlers = handler.Handlers{
	handler.Handler{
		Name:        "insert test data for Adbokis",
		Method:      "get",
		URL:         baseAdbokis,
		Description: "resets the database and inserts test data for adbokis",
		Handle:      getHandler("adbokis.json")},
	handler.Handler{
		Name:        "insert standard test data",
		Method:      "get",
		URL:         baseTestdata,
		Description: "resets the database and inserts standard test data",
		Handle:      getHandler("standard.json")},
	handler.Handler{
		Name:        "insert test data from local file",
		Method:      "get",
		URL:         baseTestdata + "/:f",
		Description: "resets the database and inserts test data from provided local file",
		Handle:      getHandler("")},
}
