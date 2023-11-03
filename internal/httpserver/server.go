package httpserver

import (
	"embed"
	"fmt"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed scripts/reload.js
//go:embed scripts/prism.js
//go:embed scripts/prism.css
var f embed.FS
var reloadData, reloadDataErr = f.ReadFile("scripts/reload.js")
var prismJSData, prismJSErr = f.ReadFile("scripts/prism.js")
var prismCSSData, prismCSSErr = f.ReadFile("scripts/prism.css")

var History = 0

func SetReload() {
	History++
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
