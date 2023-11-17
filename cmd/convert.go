/*
Copyright © 2023 Tomkoid <tomkoid@proton.me>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"codeberg.org/Tomkoid/mdhtml/internal/models"
	"codeberg.org/Tomkoid/mdhtml/internal/transform"
	"github.com/spf13/cobra"
)

var (
	sourceFile         string = ""
	out                string = ""
	stylesheet         string = ""
	watch              bool   = false
	httpserver         bool   = false
	httpserverPort     int    = 8080
	httpserverHostname string = "localhost"
	noExternalLibs     bool   = false
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [file]",
	Short: "Convert a Markdown file to HTML",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for index, arg := range args {
			if index == 0 {
				sourceFile = arg
			}
		}

		if sourceFile == "" {
			fmt.Println("No source file specified")

			cmd.Help()
		}

		if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
			fmt.Printf("Source file `%s` does not exist\n", sourceFile)

			os.Exit(1)
		}

		if out == "" {
			split := strings.Split(sourceFile, ".md")
			out = fmt.Sprintf("%s.html", split[0])
		}

		if httpserver {
			watch = true
		}

		transformArgs := models.Args{
			File:           sourceFile,
			Style:          stylesheet,
			NoExternalLibs: noExternalLibs,
			Watch:          watch,
			Debug:          debug,
			HttpServer:     httpserver,
			ServerPort:     httpserverPort,
			ServerHostname: httpserverHostname,
			Out:            out,
		}

		// fmt.Println("file flag: ", cmd.Flags().Lookup("file").Value)

		if httpserver {
			transformArgs.HttpServer = true
			transform.TransformWatch(transformArgs, httpserver)
		} else {
			transform.Transform(transformArgs, true)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	convertCmd.Flags().StringVarP(&out, "output", "o", "", "The destination file to write the HTML to")
	convertCmd.Flags().StringVarP(&stylesheet, "stylesheet", "s", "", "Apply extra styling to the HTML using a CSS file")
	convertCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch the source file for changes and reconvert when changes are detected")
	convertCmd.Flags().BoolVarP(&httpserver, "httpserver", "H", false, "Start a HTTP server to serve the HTML file and reload the page when changes are detected")
	convertCmd.Flags().IntVarP(&httpserverPort, "port", "P", 8080, "The port to use for the HTTP server")
	convertCmd.Flags().StringVarP(&httpserverHostname, "hostname", "S", "localhost", "The hostname to use for the HTTP server")
	convertCmd.Flags().BoolVarP(&noExternalLibs, "no-external-libs", "n", false, "Don't use external libraries for CSS and JS")
}
