package task

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
)

type Fire struct {
	Version            Version                `json:"version" yaml:"version"`
	Default            string                 `json:"default,omitempty" yaml:"default,omitempty"`
	EnvironmentDeclare map[string]Environment `json:"environmentDeclare" yaml:"environment-declare"`
	Tasks              []Task                 `json:"tasks" yaml:"tasks"`
	Dependencies       []string
	Replace            []Replacement `json:"replace,omitempty" yaml:"replace,omitempty"`
}

func Parse(file string) (c *Fire, err error) {
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		return nil, err
	}
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c = &Fire{}
	err = yaml.Unmarshal(fileContent, c)
	return
}

func (t *Fire) ToJson() string {
	data, _ := json.Marshal(t)
	return string(data)
}

func (t *Fire) GetAllowTaskList() (list []string) {
	for _, item := range t.Tasks {
		list = append(list, item.Name)
	}
	return
}

func (t *Fire) FindTask(taskName string) (result Task, found bool) {
	for _, item := range t.Tasks {
		if item.Name == taskName {
			return item, true
		}
	}
	return result, false
}

func (t *Fire) CreateContext() *Context {
	return &Context{}
}
