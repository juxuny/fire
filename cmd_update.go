package main

import (
	"github.com/spf13/cobra"
	"github.com/yuanjiecloud/fire/log"
	"github.com/yuanjiecloud/fire/task"
)

type updateCommand struct {
	*contextCommand
}

func (t *updateCommand) Prepare(cmd *cobra.Command) {
}

func (t *updateCommand) InitFlag(cmd *cobra.Command) {
	t.contextCommand.InitFlag(cmd)
}

func (t *updateCommand) BeforeRun(cmd *cobra.Command) {
	t.contextCommand.BeforeRun(cmd)
}

func (t *updateCommand) Run(cmd *cobra.Command, args []string) {
	mapper, err := t.pipeline.CreateRepositoryProvider().GetRepositoryMapper()
	if err != nil {
		log.Fatal(err)
	}
	keys := mapper.GetKeys()
	for _, k := range keys {
		repositoryDir, b := mapper[k]
		if !b {
			continue
		}
		if task.CheckIfGitRepository(repositoryDir) {
			log.Info("update: ", repositoryDir)
			err = task.GitFetchAndUpdate(repositoryDir)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Info("ignore local repository dir: ", repositoryDir)
		}
	}
}

func init() {
	rootCommand.AddCommand(NewCommandBuilder("update", &updateCommand{&contextCommand{}}).Build())
}
