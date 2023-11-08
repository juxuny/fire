package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{}

func main() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
