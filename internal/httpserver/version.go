package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type InfoData struct {
	EchoVersion      string `json:"echoVersion"`
	EchoServerAddr   string `json:"echoServerAddr"`
	EchoServerLogger string `json:"echoServerLogger"`
}

func Info(c echo.Context) error {
	return c.JSON(http.StatusOK, InfoData{
		EchoVersion:      echo.Version,
		EchoServerAddr:   c.Echo().Server.Addr,
		EchoServerLogger: c.Echo().Logger.Prefix(),
	})
}
