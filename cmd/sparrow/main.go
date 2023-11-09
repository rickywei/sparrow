package main

import (
	"log"

	"github.com/spf13/cobra"

	new "github.com/rickywei/sparrow/cmd/new"
)

const (
	release = "v0.1"
)

var rootCmd = &cobra.Command{
	Use:     "sparrow",
	Short:   "sparrow: An quick toolkit for Go http service.",
	Long:    `github.com/rickywei/sparrow: An quick toolkit for Go http service.`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(new.CmdNew)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
