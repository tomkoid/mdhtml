package httpserver

import (
	"fmt"
	"log"
	"time"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func difference(global []BroadcastData, local []BroadcastData) []string {
	var diffMap []string = []string{}
	for _, h := range global {
		found := false
		for _, lh := range local {
			if h.Index == lh.Index {
				found = true
				break
			}
		}
		if !found {
			diffMap = append(diffMap, h.Data)
		}
	}
	return diffMap
}

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

		localHistory := History

		for {
			var diff []string = difference(History, localHistory)

			if len(diff) == 0 {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			if args.Debug {
				fmt.Printf("> Sending message to %s using websocket to reload client...\n", c.Request().RemoteAddr)
			}

			for _, value := range diff {
				err := websocket.Message.Send(ws, value)

				if err != nil {
					if args.Debug {
						log.Printf("> Error sending message to %s using websocket to reload client: %s\n", c.Request().RemoteAddr, err)
					}

					localHistory = History
					break
				}
			}

			localHistory = History
		}
	}))}.ServeHTTP(c.Response(), c.Request())

	if args.Debug {
		fmt.Println("> Client disconnected")
	}

	return nil
}
