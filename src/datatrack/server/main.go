package main

import (
	"datatrack/config"
	"datatrack/database"
	"datatrack/database/ephemeral"
	"datatrack/handler"
	"datatrack/handler/category"
	"datatrack/handler/facebook"
	"datatrack/handler/googletakeout"
	"datatrack/handler/local"
	"datatrack/handler/testdata"
	"datatrack/handler/user"
	"net/http"
	"os"
	"runtime"

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
	err = database.Start(config.Env.Databasepath)
	if err != nil {
		panic(err)
	}

	defer func() {
		database.Close()
		os.Remove(config.Env.Databasepath)
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
