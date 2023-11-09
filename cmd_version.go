package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "v0.1.0"
)

type versionCommand struct {
}

func (t *versionCommand) Prepare(cmd *cobra.Command) {
}

func (t *versionCommand) InitFlag(cmd *cobra.Command) {
}

func (t *versionCommand) BeforeRun(cmd *cobra.Command) {
}

func (t *versionCommand) Run(cmd *cobra.Command, args []string) {
	fmt.Println(Version)
}

func init() {
	rootCommand.AddCommand(NewCommandBuilder("version", &versionCommand{}).Build())
}
