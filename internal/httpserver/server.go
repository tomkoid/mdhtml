package httpserver

import (
	"embed"
	"fmt"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed scripts/reload.js
var f embed.FS
var data, err = f.ReadFile("scripts/reload.js")

var Reload = false

func SetReload() {
	Reload = true
}

func HttpServer(args models.Args) {
	app := echo.New()
	app.HideBanner = true
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())

	if args.Debug {
		app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "${method}: ${uri}, status=${status}\n",
		}))
	}

	// Create a new router
	router := app.Group("") // root group

	setupRoutes(router, args)

	app.Logger.Fatal(app.Start(fmt.Sprintf("%s:%d", args.ServerHostname, args.ServerPort)))
}
