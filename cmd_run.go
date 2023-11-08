package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type runCommand struct {
	*contextCommand
	taskName string
}

func (t *runCommand) Run(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		t.taskName = args[0]
	}
	if t.taskName != "" {
		err := t.pipeline.RunTask(t.taskName, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	} else {
		err := t.pipeline.RunAll(nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

}

func (t *runCommand) Prepare(cmd *cobra.Command) {
}

func (t *runCommand) InitFlag(cmd *cobra.Command) {
	t.contextCommand.InitFlag(cmd)
	cmd.PersistentFlags().StringVarP(&t.taskName, "task", "t", "", "task name")
}

func init() {
	cmd := NewCommandBuilder("run", &runCommand{contextCommand: &contextCommand{}}).Build()
	rootCommand.AddCommand(cmd)
}
