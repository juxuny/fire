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
}

func (t *cleanCommand) Run(cmd *cobra.Command, args []string) {
	log.Info("cleaning...")
	mapper, err := t.pipeline.GetRepositoryMapper()
	if err != nil {
		log.Fatal(err)
	}
	list := mapper.GetKeys()
	for _, name := range list {
		if task.CheckIfNestedRepository(mapper[name]) {
			log.Debug("ignore nested repository:", mapper[name])
			continue
		}
		log.Info("clean repository: ", name, " => ", mapper[name])
		if task.CheckIfExists(mapper[name]) {
			err = os.RemoveAll(mapper[name])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func init() {
	rootCommand.AddCommand(NewCommandBuilder("clean", &cleanCommand{&contextCommand{}}).Build())
}
