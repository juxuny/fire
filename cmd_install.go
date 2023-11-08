package main

import (
	"github.com/spf13/cobra"
)

type installCommand struct {
}

func (t *installCommand) Prepare(cmd *cobra.Command) {
}

func (t *installCommand) InitFlag(cmd *cobra.Command) {
}

func (t *installCommand) BeforeRun(cmd *cobra.Command) {
}

func (t *installCommand) Run(cmd *cobra.Command, args []string) {
}

func init() {
	rootCommand.AddCommand(NewCommandBuilder("install", &installCommand{}).Build())
}
