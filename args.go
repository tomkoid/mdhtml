package main

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

type Args struct {
	file  string
	out   string
	style string
}

func parseArgs() Args {
	// create and parse flags using the flag package
	file := flag.StringP("file", "f", "", "The markdown file to convert to HTML")
	out := flag.StringP("out", "o", "", "The destination file to write the HTML to")
	style := flag.StringP("style", "s", "", "Apply extra styling to the HTML using a CSS file")

	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *out == "" {
		split := strings.Split(*file, ".md")
		*out = fmt.Sprintf("%s.html", split[0])
	}

	// return instance of Args
	return Args{
		file:  *file,
		out:   *out,
		style: *style,
	}
}
