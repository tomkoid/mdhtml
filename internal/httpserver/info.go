package httpserver

import (
	"fmt"
	"time"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"github.com/fatih/color"
)

func printServerInfo(args models.Args) {
	if !args.NoServerHeaderWait {
		// sleep for 50ms to prevent this message from printing if the server is instantly killed or crashes
		time.Sleep(50 * time.Millisecond)
	}

	prefix := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fmt.Printf("%s: Server running on %s:%d\n", prefix("server"), args.ServerHostname, args.ServerPort)
}
