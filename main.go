package main

import (
	"net/http"
	"os"
	"runtime"

	"github.com/marcelfarres/datatrack/config"
	"github.com/marcelfarres/datatrack/database"
	"github.com/marcelfarres/datatrack/ephemeral"
	"github.com/marcelfarres/datatrack/handler"
	"github.com/marcelfarres/datatrack/handler/category"
	"github.com/marcelfarres/datatrack/handler/facebook"
	"github.com/marcelfarres/datatrack/handler/googletakeout"
	"github.com/marcelfarres/datatrack/handler/local"
	"github.com/marcelfarres/datatrack/handler/testdata"
	"github.com/marcelfarres/datatrack/handler/user"

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
