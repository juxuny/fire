package main

import (
	"github.com/spf13/cobra"
	"github.com/yuanjiecloud/fire/task"
	"log"
)

type contextCommand struct {
	configFile string

	fireInstance *task.Fire
}

func (t *contextCommand) InitFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&t.configFile, "c", "fire.yaml", "config file name")
}

func (t *contextCommand) BeforeRun(cmd *cobra.Command) {
	var err error
	t.fireInstance, err = task.Parse(t.configFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(t.fireInstance.ToJson())
}
