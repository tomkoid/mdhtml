package httpserver

import (
	"net/http"

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

	app.GET("/", Page)
	app.GET("/info", Info)

	app.Any("/ws", WSEndpoint)

	/// assets
	app.GET("/default.css", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/css")
		return c.String(http.StatusOK, string(defaultCSSData))
	})

	app.GET("/reload.js", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/javascript")
		return c.String(http.StatusOK, string(reloadData))
	})

	// syntax highlighting
	app.GET("/prism.js", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/javascript")
		return c.String(http.StatusOK, string(prismJSData))
	})

	app.GET("/prism.css", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/css")
		return c.String(http.StatusOK, string(prismCSSData))
	})
}
