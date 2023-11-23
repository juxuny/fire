package main

import (
	"github.com/spf13/cobra"
	"github.com/yuanjiecloud/fire/log"
)

type installCommand struct {
	*contextCommand
}

func (t *installCommand) Prepare(cmd *cobra.Command) {
}

func (t *installCommand) InitFlag(cmd *cobra.Command) {
	t.contextCommand.InitFlag(cmd)
}

func (t *installCommand) BeforeRun(cmd *cobra.Command) {
	t.contextCommand.BeforeRun(cmd)
}

func (t *installCommand) Run(cmd *cobra.Command, args []string) {
	err := t.pipeline.Resolve()
	if err != nil {
		log.Fatal("resolve error: ", err)
	}
	log.Info("Success")
}

func init() {
	rootCommand.AddCommand(NewCommandBuilder("install", &installCommand{contextCommand: &contextCommand{}}).Build())
}
