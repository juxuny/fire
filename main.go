package main

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCommand = &cobra.Command{}

func main() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
