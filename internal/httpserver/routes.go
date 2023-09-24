package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"github.com/labstack/echo/v4"
)

func setupRoutes(app *echo.Group, args models.Args) {
	// MIDDLEWARE: pass args to every route
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("args", args)
			return next(c)
		}
	})

	app.GET("/", func(c echo.Context) error {
		if args.Debug {
			fmt.Printf("> Access from %s to HTTP server\n", c.Request().RemoteAddr)
		}

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

	app.Any("/ws", WSEndpoint)

	app.GET("/reload.js", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/javascript")
		return c.String(http.StatusOK, string(data))
	})
}
