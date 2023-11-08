package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type listCommand struct {
	*contextCommand
}

func (t *listCommand) Prepare(cmd *cobra.Command) {
}

func (t *listCommand) InitFlag(cmd *cobra.Command) {
	t.contextCommand.InitFlag(cmd)
}

func (t *listCommand) BeforeRun(cmd *cobra.Command) {
	t.contextCommand.BeforeRun(cmd)
}

func (t *listCommand) Run(cmd *cobra.Command, args []string) {
	list := t.pipeline.GetAllowTaskList()
	for _, item := range list {
		fmt.Println(item)
	}
}

func init() {
	rootCommand.AddCommand(NewCommandBuilder("list", &listCommand{&contextCommand{}}).Build())
}
