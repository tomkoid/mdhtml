package httpserver

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

//go:embed scripts/reload.js
var f embed.FS
var data, err = f.ReadFile("scripts/reload.js")

var reload = false

func SetReload() {
	reload = true
}

func wsEndpoint(c echo.Context) error {
	fmt.Printf("> Access from %s to WS server\n", c.Request().RemoteAddr)
	websocket.Server{Handler: websocket.Handler(websocket.Handler(func(ws *websocket.Conn) {
		websocket.Message.Send(ws, "hello")
		fmt.Println("> Client connected")
		defer ws.Close()
		fmt.Println("> Sending message to websocket to reload client...")

		for {
			if reload {
				fmt.Println("> Sending message to websocket to reload client...")
				err := websocket.Message.Send(ws, "reload")

				if err != nil {
					continue
				}

				reload = false
			}
		}
	}))}.ServeHTTP(c.Response(), c.Request())

	return nil
}

func setupRoutes(app *echo.Echo, args models.Args) {
	app.GET("/", func(c echo.Context) error {
		fmt.Printf("> Access from %s to HTTP server\n", c.Request().RemoteAddr)

		// read file contents
		contents, err := os.ReadFile(args.Out)

		if err != nil {
			log.Fatalf("Error reading file: %s", err)

			return c.String(http.StatusOK, fmt.Sprintf("Error reading file: %s", err))
		}

		// return html content type
		c.Response().Header().Set("Content-Type", "text/html")

		// return file contents
		return c.String(http.StatusOK, string(contents))
	})

	app.Any("/ws", wsEndpoint)

	app.GET("/reload.js", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/javascript")
		return c.String(http.StatusOK, string(data))
	})
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

	setupRoutes(app, args)

	app.Logger.Fatal(app.Start(":8080"))
}
