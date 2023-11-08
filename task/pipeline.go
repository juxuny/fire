package task

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Pipeline struct {
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
	return
}

func (t *Pipeline) ToJson() string {
	data, _ := json.Marshal(t)
	return string(data)
}

func (t *Pipeline) GetAllowTaskList() (list []string) {
	for _, item := range t.Tasks {
		list = append(list, item.Name)
	}
	return
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
		EnvProvider: t.Environments.Clone(),
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
		err = item.Exec(newContext)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Pipeline) RunTask(name string, ctx *Context) error {
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
