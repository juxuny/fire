package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuanjiecloud/fire/log"
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
	err := t.pipeline.Preload()
	log.CheckAndFatal(err)
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
