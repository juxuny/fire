package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type runCommand struct {
	*contextCommand
	taskName string
}

func (t *runCommand) Run(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		t.taskName = args[0]
	}
	selectedTask, found := t.fireInstance.FindTask(t.taskName)
	if !found {
		fmt.Println("task not found: ", t.taskName)
	}
	err := selectedTask.Exec(t.fireInstance.CreateContext())
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func (t *runCommand) Prepare(cmd *cobra.Command) {
}

func (t *runCommand) InitFlag(cmd *cobra.Command) {
	t.contextCommand.InitFlag(cmd)

}

func init() {
	cmd := NewCommandBuilder("run", &runCommand{contextCommand: &contextCommand{}}).Build()
	cmd.Args = cobra.ExactArgs(1)
	rootCommand.AddCommand(cmd)
}
