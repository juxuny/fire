package main

import (
	"github.com/spf13/cobra"
)

type updateCommand struct {
}

func (t *updateCommand) Prepare(cmd *cobra.Command) {
}

func (t *updateCommand) InitFlag(cmd *cobra.Command) {
}

func (t *updateCommand) BeforeRun(cmd *cobra.Command) {
}

func (t *updateCommand) Run(cmd *cobra.Command, args []string) {
}

func init() {
	rootCommand.AddCommand(NewCommandBuilder("update", &updateCommand{}).Build())
}
