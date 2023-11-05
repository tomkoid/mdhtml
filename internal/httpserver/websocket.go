package httpserver

import (
	"fmt"
	"log"
	"time"

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

		localHistory := History

		for {
			// if reflect.DeepEqual(localHistory, History) {
			// 	time.Sleep(50 * time.Millisecond)
			// 	continue
			// }

			// println("localHistory: ", localHistory)
			// println("History: ", History)

			// get the difference between the two slices
			// https://stackoverflow.com/a/45485900
			var diff []string = []string{}
			for _, h := range History {
				found := false
				for _, lh := range localHistory {
					if h == lh {
						found = true
						break
					}
				}
				if !found {
					diff = append(diff, h)
				}
			}
			for _, value := range History {
				fmt.Printf("- %s\n", value)
			}

			println(len(diff))
			if len(diff) == 0 {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			println("diff: ", diff)

			// please dont forget to remove this later lol
			// diff = []string{}

			if args.Debug {
				fmt.Printf("> Sending message to %s using websocket to reload client...\n", c.Request().RemoteAddr)
			}

			err := websocket.Message.Send(ws, "reload")

			if err != nil {
				if args.Debug {
					log.Printf("> Error sending message to %s using websocket to reload client: %s\n", c.Request().RemoteAddr, err)
				}

				break
			}

			localHistory = History
		}
	}))}.ServeHTTP(c.Response(), c.Request())

	if args.Debug {
		fmt.Println("> Client disconnected")
	}

	return nil
}
