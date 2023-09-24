package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func WSEndpoint(c echo.Context) error {
	fmt.Printf("> Access from %s to WS server\n", c.Request().RemoteAddr)
	websocket.Server{Handler: websocket.Handler(websocket.Handler(func(ws *websocket.Conn) {
		websocket.Message.Send(ws, "hello")
		fmt.Println("> Client connected")
		defer ws.Close()
		fmt.Println("> Sending message to websocket to reload client...")

		for {
			if Reload {
				fmt.Println("> Sending message to websocket to reload client...")
				err := websocket.Message.Send(ws, "reload")

				if err != nil {
					continue
				}

				Reload = false
			}
		}
	}))}.ServeHTTP(c.Response(), c.Request())

	return nil
}
