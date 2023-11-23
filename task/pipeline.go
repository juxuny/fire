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
	log.Debug("parsing file: ", file, " wd:", Getwd())
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
		EnvProvider: t.Environments.Clone(),
	}
	if ctx != nil && ctx.EnvProvider != nil {
		result.Parent = ctx
		result.EnvProvider = result.EnvProvider.MergeIgnoreDuplicated(ctx.EnvProvider)
	}
	return result
}

func (t *Pipeline) RunAll(ctx *Context) error {
	var err error
	newContext := t.CreateContext(ctx)
	for _, item := range t.Tasks {
		showTitle(fmt.Sprintf("start task: %s", item.Name))
		log.Debug("task env: ", item.Env)
		log.Debug("context current env: ", ctx.GetCurrentEnv())
		err = item.Exec(newContext.UseEnv(item.Env))
		if err != nil {
			return err
		}
		showTitle(fmt.Sprintf("end(%s)", item.Name))
	}
	return nil
}

func (t *Pipeline) RunTask(name string, ctx *Context) error {
	selected, found := t.FindTask(name)
	if !found {
		log.Debug("task not found: ", name)
		pipeline, b := FindPipeline(name)
		if !b {
			log.Debug("pipeline not found: ", name)
			return errors.Errorf("task not found: %s", name)
		}
		return pipeline.RunAll(t.CreateContext(ctx))
	}
	showTitle(fmt.Sprintf("start task: %s", name))
	defer func() {
		showTitle(fmt.Sprintf("end(%s)", name))
	}()
	return selected.Exec(t.CreateContext(ctx).UseEnv(selected.Env))
}

func (t *Pipeline) Resolve() error {
	resolver := NewResolver(t.Dependencies, t.Replace)
	return resolver.Start()
}

func (t *Pipeline) GetAllowTaskList() (taskList datatype.SortableStringList) {
	if t == nil {
		return make(datatype.SortableStringList, 0)
	}
	filter := make(map[string]bool)
	for _, item := range t.Tasks {
		if _, b := filter[item.Name]; !b {
			taskList = append(taskList, item.Name)
			filter[item.Name] = true
			if item.Pipeline != "" {
				pipeline, pipelineExists := FindPipeline(item.Pipeline)
				if !pipelineExists {
					continue
				}
				for _, tn := range pipeline.GetAllowTaskList() {
					if _, b := filter[tn]; !b {
						filter[tn] = true
						taskList = append(taskList, tn)
					}
				}
			}
		}
	}
	return
}

func (t *Pipeline) Preload() error {
	var err error
	enter(t.Getwd())
	defer goback()
	replaceMapper := make(map[string]Replacement)
	for _, replace := range t.Replace {
		replaceMapper[replace.Package] = replace
	}
	for _, depend := range t.Dependencies {
		replacement, found := replaceMapper[depend]
		log.Debug("resolving dependency: ", depend)
		var namespace, name, version string
		namespace, name, version, err = SplitPackageName(depend)
		if err != nil {
			log.Fatal("invalid repository:", depend)
		}
		if found {
			if replacement.Repository == "" {
				log.Fatal("repository in replacement is empty")
			}
			if replacement.IsLocal() {
				repositoryDir := path.Join(Getwd(), replacement.Repository)
				pipeline, err := AddPipeline(depend, repositoryDir)
				if err != nil {
					log.Fatal(err)
				}
				err = pipeline.Preload()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				repositoryDir := CreateRepositoryLocationSpecificVersion(namespace, name, replacement.Version.String())
				pipeline, err := AddPipeline(depend, repositoryDir)
				if err != nil {
					log.Fatal(err)
				}
				err = pipeline.Preload()
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			log.Debug("no replacement dependency: ", name)
			repositoryDir := CreateRepositoryLocationSpecificVersion(namespace, name, version)
			pipeline, err := AddPipeline(depend, repositoryDir)
			if err != nil {
				log.Fatal(err)
			}
			err = pipeline.Preload()
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	return nil
}

func (t *Pipeline) Getwd() string {
	return path.Dir(t.configfile)
}

func (t *Pipeline) CleanDependencies() error {
	for _, depend := range t.Dependencies {
		reposDir, found := FindPipelineReposDir(depend)
		if !found {
			continue
		}
		pipeline, found := FindPipeline(depend)
		if !found {
			continue
		}
		err := pipeline.CleanDependencies()
		log.CheckAndFatal(err)
		log.Debug("found dependency: ", reposDir)
		if !CheckIfNestedRepository(reposDir) {
			log.Info("rm ", reposDir)
			err = os.RemoveAll(reposDir)
			log.CheckAndFatal(err)
		}
	}
	return nil
}

func (t *Pipeline) UpdateDependencies() {
	for _, depend := range t.Dependencies {
		reposDir, found := FindPipelineReposDir(depend)
		if !found {
			continue
		}
		pipeline, found := FindPipeline(depend)
		if !found {
			continue
		}
		pipeline.UpdateDependencies()
		if CheckIfGitRepository(reposDir) {
			log.Info("update: ", reposDir)
			err := GitFetchAndUpdate(reposDir)
			log.CheckAndFatal(err)
		} else {
			log.Info("ignore local repository dir: ", reposDir)
		}
	}
}
