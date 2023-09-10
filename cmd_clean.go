package main

import "github.com/spf13/cobra"

type cleanCommand struct {
}

func (t *cleanCommand) Prepare(cmd *cobra.Command) {
}

func (t *cleanCommand) InitFlag(cmd *cobra.Command) {
}

func (t *cleanCommand) BeforeRun(cmd *cobra.Command) {
}

func (t *cleanCommand) Run(cmd *cobra.Command, args []string) {
}

func init() {
	rootCommand.AddCommand(NewCommandBuilder("clean", &cleanCommand{}).Build())
}
