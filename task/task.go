package task

type ExecutorType string

const (
	ExecutorTypeBash = ExecutorType("bash")
	ExecutorTypeSsh  = ExecutorType("ssh")
)

type Task struct {
	Name         string       `json:"name" yaml:"name"`
	Environments Environment  `json:"environments" yaml:"environments"`
	Type         ExecutorType `json:"type" yaml:"type"`
	UseEnv       string       `json:"useEnv" yaml:"use-env"`
	SubTasks     []SubTask    `json:"subTasks,omitempty" yaml:"sub-tasks,omitempty"`
	Scripts      []string     `json:"scripts,omitempty" yaml:"scripts,omitempty"`
}

type SubTask struct {
	Name         string      `json:"name" yaml:"name"`
	Package      string      `json:"package"`
	Environments Environment `json:"environments" yaml:"environments"`
	UseEnv       string      `json:"useEnv" yaml:"use-env"`
}

func (t *Task) Exec(ctx *Context) error {
	// TODO: execute task
	return nil
}
