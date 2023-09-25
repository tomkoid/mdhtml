package httpserver

import (
	"fmt"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func WSEndpoint(c echo.Context) error {
	args := c.Get("args").(models.Args)

	if args.Debug {
		fmt.Printf("> Access from %s to WS server\n", c.Request().RemoteAddr)
	}

	websocket.Server{Handler: websocket.Handler(websocket.Handler(func(ws *websocket.Conn) {
		websocket.Message.Send(ws, "hello")

		if args.Debug {
			fmt.Println("> Client connected")
		}

		defer ws.Close()

		for <-Reload {
			if args.Debug {
				fmt.Printf("> Sending message to %s using websocket to reload client...\n", c.Request().RemoteAddr)
			}

			err := websocket.Message.Send(ws, "reload")

			if err != nil {
				continue
			}

			Reload = false
		}
	}))}.ServeHTTP(c.Response(), c.Request())

	return nil
}
