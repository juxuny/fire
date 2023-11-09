package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yuanjiecloud/fire/log"
	"github.com/yuanjiecloud/fire/task"
)

type contextCommand struct {
	configFile string
	workdir    string
	verbose    bool

	pipeline *task.Pipeline
}

func (t *contextCommand) InitFlag(cmd *cobra.Command) {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd.PersistentFlags().StringVarP(&t.configFile, "config", "c", "fire.yaml", "config file name")
	cmd.PersistentFlags().StringVarP(&t.workdir, "workdir", "w", workdir, "working directory")
	cmd.PersistentFlags().BoolVarP(&t.verbose, "verbose", "v", false, "display debug log")
}

func (t *contextCommand) BeforeRun(cmd *cobra.Command) {
	log.Verbose = t.verbose
	var err error
	err = os.Chdir(t.workdir)
	if err != nil {
		log.Fatal(err)
		return
	}
	t.pipeline, err = task.Parse(t.configFile)
	if err != nil {
		log.Fatal(err)
	}
	globalConfigPath, err := task.GetGlobalCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("global config path: ", globalConfigPath)
	globalReposDir, err := task.GetGlobalReposDir()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("global repos dir: ", globalReposDir)
	log.Debug(t.pipeline.ToJson())
}
