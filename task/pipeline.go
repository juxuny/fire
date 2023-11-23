package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/yuanjiecloud/fire/datatype"
	"github.com/yuanjiecloud/fire/log"
	"gopkg.in/yaml.v3"
)

type Pipeline struct {
	configfile string

	Version      Version       `json:"version,omitempty" yaml:"version,omitempty"`
	Environments EnvProvider   `json:"environments,omitempty" yaml:"environments,omitempty"`
	Tasks        []Task        `json:"tasks,omitempty" yaml:"tasks,omitempty"`
	Dependencies []string      `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Replace      []Replacement `json:"replace,omitempty" yaml:"replace,omitempty"`
}

func Parse(file string) (c *Pipeline, err error) {
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		return nil, err
	}
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c = &Pipeline{}
	err = yaml.Unmarshal(fileContent, c)
	if err != nil {
		return
	}
	c.configfile = file
	return
}

func (t *Pipeline) ToJson() string {
	data, _ := json.Marshal(t)
	return string(data)
}

func (t *Pipeline) FindTask(taskName string) (result Task, found bool) {
	for _, item := range t.Tasks {
		if item.Name == taskName {
			return item, true
		}
	}
	return result, false
}

func (t *Pipeline) CreateContext(ctx *Context) *Context {
	result := &Context{
		EnvProvider:        t.Environments.Clone(),
		RepositoryProvider: t.CreateRepositoryProvider(),
	}
	if ctx != nil && ctx.EnvProvider != nil {
		result.EnvProvider = result.EnvProvider.MergeIgnoreDuplicated(ctx.EnvProvider)
	}
	return result
}

func (t *Pipeline) RunAll(ctx *Context) error {
	var err error
	newContext := t.CreateContext(ctx)
	for _, item := range t.Tasks {
		showTitle(fmt.Sprintf("start task: %s", item.Name))
		err = item.Exec(newContext)
		if err != nil {
			return err
		}
		showTitle(fmt.Sprintf("end(%s)", item.Name))
	}
	return nil
}

func (t *Pipeline) RunTask(name string, ctx *Context) error {
	showTitle(fmt.Sprintf("start task: %s", name))
	defer func() {
		showTitle(fmt.Sprintf("end(%s)", name))
	}()
	selected, found := t.FindTask(name)
	if !found {
		return errors.Errorf("task not found: %s", name)
	}
	return selected.Exec(t.CreateContext(ctx))
}

func (t *Pipeline) Resolve() error {
	resolver := NewResolver(t.Dependencies, t.Replace)
	return resolver.Start()
}

func (t *Pipeline) CreateRepositoryProvider() (repositoryProvider *Provider) {
	repositoryProvider = NewProvider(t.Dependencies, t.Replace)
	return
}

func (t *Pipeline) GetRepositoryMapper() (mapper ReposMapper, err error) {
	return t.CreateRepositoryProvider().GetRepositoryMapper()
}

func (t *Pipeline) GetAllowTaskList() (taskList datatype.SortableStringList) {
	if t == nil {
		return make(datatype.SortableStringList, 0)
	}
	repositoryProvider := t.CreateRepositoryProvider()
	for _, item := range t.Tasks {
		taskList = append(taskList, item.Name)
		if item.Pipeline != "" {
			pipeline, err := repositoryProvider.FindPipeline(item.Pipeline)
			if err != nil {
				log.Error(err)
				continue
			}
			taskList = append(taskList, pipeline.GetAllowTaskList()...)
		}
	}
	return
}

func (t *Pipeline) Getwd() string {
	return path.Dir(t.configfile)
}
