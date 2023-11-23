package task

import (
	"os"

	"github.com/pkg/errors"
	"github.com/yuanjiecloud/fire/executor"
	"github.com/yuanjiecloud/fire/log"
)

type ExecutorType string

const (
	ExecutorTypeBash = ExecutorType("bash")
	ExecutorTypeSsh  = ExecutorType("ssh")
)

type Task struct {
	Name         string               `json:"name,omitempty" yaml:"name,omitempty"`
	Environments Environment          `json:"environments,omitempty" yaml:"environments,omitempty"`
	Type         executor.Type        `json:"type,omitempty" yaml:"type,omitempty"`
	Env          string               `json:"env,omitempty" yaml:"env,omitempty"`
	Pipeline     string               `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`
	Scripts      []string             `json:"scripts,omitempty" yaml:"scripts,omitempty"`
	SshOptions   *executor.SshOptions `json:"sshOptions,omitempty" yaml:"ssh-options,omitempty"`
}

func (t *Task) Exec(ctx *Context) error {
	log.Debug("context env: ", ctx.GetCurrentEnv())
	var err error
	err = t.runPipeline(ctx)
	if err != nil {
		return err
	}
	err = t.runScripts(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) runPipeline(ctx *Context) error {
	if t.Pipeline == "" {
		return nil
	}
	var err error
	log.Debug("start pipeline: ", t.Pipeline)
	pipeline, found := FindPipeline(t.Pipeline)
	if !found {
		return errors.Errorf("pipeline not found: %s", t.Pipeline)
	}
	wd := Getwd()
	err = os.Chdir(pipeline.Getwd())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = os.Chdir(wd)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return pipeline.RunAll(ctx)
}

func (t *Task) getCurrentEnv(ctx *Context) (result Environment, found bool) {
	if t.Environments != nil {
		result = t.Environments
	} else {
		result = make(Environment)
	}
	currentEnv := ctx.GetCurrentEnv()
	if currentEnv != "" {
		env, b := ctx.GetEnv(currentEnv)
		if b {
			result = result.MergeIgnoreDuplicated(env)
			found = true
		}
	}
	if t.Env != "" {
		env, b := ctx.GetEnv(t.Env)
		if b {
			result = result.MergeIgnoreDuplicated(env)
			found = true
		}
	}
	return
}

func (t *Task) runScripts(ctx *Context) (err error) {
	if len(t.Scripts) == 0 {
		return nil
	}
	env, found := t.getCurrentEnv(ctx)
	if !found {
		return errors.Errorf("unset env")
	}
	if t.Type == executor.TypeBash {
		return executor.NewBashExecutor(env, t.Scripts).StartAndWait()
	} else if t.Type == executor.TypeSsh {
		return executor.NewSshExecutor(env, t.Scripts, t.SshOptions).StartAndWait()
	} else {
		return errors.Errorf("unknown executor type: %v", t.Type)
	}
}
