package main

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/yuanjiecloud/fire/log"
	"github.com/yuanjiecloud/fire/task"
)

type contextCommand struct {
	workdir string
	verbose bool
	global  bool

	pipeline *task.Pipeline
}

func (t *contextCommand) InitFlag(cmd *cobra.Command) {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd.PersistentFlags().StringVarP(&t.workdir, "workdir", "w", workdir, "working directory")
	cmd.PersistentFlags().BoolVarP(&t.verbose, "verbose", "v", false, "display debug log")
	cmd.PersistentFlags().BoolVarP(&t.global, "global", "g", false, "use global config")
}

func (t *contextCommand) BeforeRun(cmd *cobra.Command) {
	log.Verbose = t.verbose
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

	t.prepareGlobalConfig()
	configFile := path.Join(t.workdir, task.DefaultConfigFile)
	if task.CheckIfExists(configFile) {
		// found fire.yaml in working dir
		log.Debug(fmt.Sprintf("found %s in working directory: %s", task.DefaultConfigFile, t.workdir))
		err = os.Chdir(t.workdir)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		// change to global config dir, and use default global config file
		globalConfigDir, err := task.GetGlobalConfigDir()
		if err != nil {
			log.Fatal(err)
		}
		err = os.Chdir(globalConfigDir)
		if err != nil {
			log.Fatal("use global config error:", err)
		}
		configFile = path.Join(globalConfigDir, task.DefaultConfigFile)
	}
	t.pipeline, err = task.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(t.pipeline.ToJson())
}

func (t *contextCommand) prepareGlobalConfig() {
	dir, err := task.GetGlobalCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	if !task.CheckIfExists(dir) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal("init global config failed:", err)
		}
	}
	globalConfigFile, err := task.GetGlobalFireConfig()
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stat(globalConfigFile)
	if os.IsNotExist(err) {
		var f *os.File
		// initialize global config file
		f, err = os.Create(globalConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err = f.WriteString("")
		if err != nil {
			log.Fatal("initialize global config file error:", err)
		}
	} else if err != nil {
		log.Fatal(err)
	}
}
