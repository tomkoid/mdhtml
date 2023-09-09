package main

import (
	"flag"
	"os"
)

type Args struct {
	file  string
	out   string
	style string
}

func parseArgs() Args {
	// create and parse flags using the flag package
	file := flag.String("file", "", "The markdown file to convert to HTML")
	out := flag.String("out", "", "The destination file to write the HTML to")
	style := flag.String("style", "", "Apply extra styling to the HTML using a CSS file")
	// test := flag.String("test", "", "The destination file to write the HTML to")

	flag.Parse()

	if *file == "" || *out == "" {
		flag.Usage()
		os.Exit(1)
	}

	// return instance of Args
	return Args{
		file:  *file,
		out:   *out,
		style: *style,
	}
}
