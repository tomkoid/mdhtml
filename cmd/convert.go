/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
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

var sourceFile string = ""

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [FILE]",
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

		// print("sourceFile: ", sourceFile)

		var out string

		if out == "" {
			split := strings.Split(sourceFile, ".md")
			out = fmt.Sprintf("%s.html", split[0])
		}

		transformArgs := models.Args{
			File: sourceFile,
			Out:  out,
		}

		// fmt.Println("file flag: ", cmd.Flags().Lookup("file").Value)

		transform.Transform(transformArgs, true)
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
	convertCmd.Flags().BoolP("file", "f", false, "Source file")
}
