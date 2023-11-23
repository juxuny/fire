package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yuanjiecloud/fire/log"
	"github.com/yuanjiecloud/fire/task"
)

type cleanCommand struct {
	*contextCommand
}

func (t *cleanCommand) Prepare(cmd *cobra.Command) {
}

func (t *cleanCommand) InitFlag(cmd *cobra.Command) {
	t.contextCommand.InitFlag(cmd)
}

func (t *cleanCommand) BeforeRun(cmd *cobra.Command) {
	t.contextCommand.BeforeRun(cmd)
	err := t.pipeline.Preload()
	log.CheckAndFatal(err)
}

func (t *cleanCommand) Run(cmd *cobra.Command, args []string) {
	log.Info("cleaning...")
	if t.global {
		dir, err := task.GetGlobalReposDir()
		log.CheckAndFatal(err)
		err = os.RemoveAll(dir)
		log.CheckAndFatal(err)
	} else {
		err := t.pipeline.CleanDependencies()
		log.CheckAndFatal(err)
	}
}

func init() {
	rootCommand.AddCommand(NewCommandBuilder("clean", &cleanCommand{&contextCommand{}}).Build())
}
