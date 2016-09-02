package main

import (
	"net/http"
	"os"
	"runtime"

	"github.com/pylls/datatrack/config"
	"github.com/pylls/datatrack/database"
	"github.com/pylls/datatrack/ephemeral"
	"github.com/pylls/datatrack/handler"
	"github.com/pylls/datatrack/handler/category"
	"github.com/pylls/datatrack/handler/facebook"
	"github.com/pylls/datatrack/handler/googletakeout"
	"github.com/pylls/datatrack/handler/local"
	"github.com/pylls/datatrack/handler/testdata"
	"github.com/pylls/datatrack/handler/user"

	"github.com/toqueteos/webbrowser"
	"github.com/unrolled/secure"
	"github.com/zenazn/goji"
)

var handlers = handler.Concat(
	handler.CommonHandlers,
	user.Handlers,
	local.Handlers,
	category.Handlers,
	googletakeout.Handlers,
	facebook.Handlers,
	testdata.Handlers,
)

func main() {
	// read config data
	err := config.Configure()
	if err != nil {
		panic(err)
	}

	// setup ephemeral database encryption
	ephemeral.Setup()

	// start database
	err = database.Start(config.Env.DatabasePath)
	if err != nil {
		panic(err)
	}

	defer func() {
		database.Close()
		os.Remove(config.Env.DatabasePath)
	}()

	// handlers
	if err := handlers.Register(); err != nil {
		panic(err)
	}
	goji.Handle("/*", http.FileServer(http.Dir(config.StaticPath)))

	// middleware
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentSecurityPolicy: "default-src 'self'; img-src 'self' https://*.openstreetmap.org; style-src 'self' 'unsafe-inline'",
		BrowserXssFilter:      true,
	})
	goji.Use(secureMiddleware.Handler)

	// open browser
	webbrowser.Open("http://localhost:8000/")

	// clear environment if not Windows due to odd reliance on environment
	// for form parsing in Windows
	if runtime.GOOS != "windows" {
		os.Clearenv()
	}

	goji.Serve()
}
