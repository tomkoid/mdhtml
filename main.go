package main

import (
	"fmt"
)

func main() {
	// parse arguments
	args := parseArgs()

	if args.watch {
		fmt.Printf("== Watching %s for changes...\n", args.file)
		transformWatch(args, true)
		return
	}

	checkFilesExist(args)

	transform(args, true)
}
