package main

import (
	"fmt"

	"codeberg.org/Tomkoid/mdhtml/internal/transform"
	"codeberg.org/Tomkoid/mdhtml/internal/utils"
)

func main() {
	// parse arguments
	args := utils.ParseArgs()

	if args.Watch {
		fmt.Printf("== Watching %s for changes...\n", args.File)
		transform.TransformWatch(args, args.HttpServer)
		return
	}

	utils.CheckFilesExist(args)

	transform.Transform(args, true)
}
