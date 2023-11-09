package task

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/yuanjiecloud/fire/log"
)

type Resolver struct {
	dependencies  []string
	replaceMapper map[string]Replacement
}

func NewResolver(dependencies []string, replace []Replacement) *Resolver {
	ret := &Resolver{
		dependencies:  dependencies,
		replaceMapper: make(map[string]Replacement),
	}
	for _, item := range replace {
		ret.replaceMapper[item.Package] = item
	}
	return ret
}

func (t *Resolver) checkout(repositoryUrl string, namespace, name, branch string, repositoryPath string) error {
	stat, err := os.Stat(repositoryPath)
	if os.IsNotExist(err) {
		return errors.Errorf("repository directory not exists: %v", repositoryPath)
	}
	if !stat.IsDir() {
		return errors.Errorf("%s is not a directory", repositoryPath)
	}
	if branch == "" {
		if strings.Contains(repositoryUrl, "github.com") {
			branch = "main"
		} else {
			branch = "master"
		}
	}
	gitCommand := exec.Command("git", "clone", "-b", branch, repositoryUrl, branch)
	gitCommand.Stderr = os.Stderr
	gitCommand.Stdout = os.Stdout
	gitCommand.Dir = path.Join(repositoryPath, namespace, name)
	if CheckIfExists(path.Join(gitCommand.Dir, branch)) {
		if CheckIfExists(path.Join(gitCommand.Dir, branch, ".git")) {
			log.Info("ignore: ", fmt.Sprintf("%s/%s@%s", namespace, name, branch))
			return nil
		} else {
			err = os.Remove(path.Join(gitCommand.Dir, branch))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	log.Debug(gitCommand.Dir, "=>", repositoryUrl, branch)
	err = os.MkdirAll(gitCommand.Dir, 0775)
	if err != nil {
		log.Error(err)
		return errors.Errorf("init cache dir failed: %v", gitCommand.Dir)
	}
	err = gitCommand.Run()
	if err != nil {
		return err
	}
	return t.resolveDirectory(path.Join(repositoryPath, namespace, name, branch))
}

func (t *Resolver) resolveDirectory(dir string) error {
	wd := Getwd()
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	if path.IsAbs(dir) {
		log.Info("enter dir: ", dir)
	} else {
		log.Info("enter dir: ", path.Join(wd, dir))
	}
	defer func() {
		err = os.Chdir(wd)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("goback dir: ", wd)
	}()
	if !CheckIfExists(DefaultConfigFile) {
		return errors.Errorf("%s is an invalid repository, %s not found", dir, DefaultConfigFile)
	}
	pipeline, err := Parse(DefaultConfigFile)
	if err != nil {
		log.Fatal("invalid repository, ", err)
		return errors.Errorf("invalid repository")
	}
	return pipeline.Resolve()
}

func (t *Resolver) Start() error {
	reposDir, err := GetGlobalReposDir()
	if err != nil {
		log.Fatal(err)
	}
	for _, depend := range t.dependencies {
		var (
			namespace string
			name      string
			version   string
			err       error
		)
		log.Info("resolving", depend)
		replacement, found := t.replaceMapper[depend]
		if found {
			namespace, name, _, err = SplitPackageName(replacement.Package)
			if err != nil {
				log.Error(err)
				return errors.Errorf("resolve replacement failed: %v", depend)
			}
			version = replacement.Version.String()
			if replacement.IsLocal() {
				if CheckIfExists(replacement.Repository) {
					log.Info("ignore local repository: ", depend)
					err = t.resolveDirectory(replacement.Repository)
					if err != nil {
						log.Fatal("resolve repository failed: ", err)
					}
					continue
				}
			}
			err = t.checkout(replacement.Repository, namespace, name, version, reposDir)
			if err != nil {
				log.Error(err)
				return errors.Errorf("checkout replacement failed: %v", depend)
			}
		} else {
			namespace, name, version, err = SplitPackageName(depend)
			if err != nil {
				log.Error(err)
				return errors.Errorf("resolve dependencies failed")
			}
			repositoryUrl := fmt.Sprintf("https://github.com/%s/%s.git", namespace, name)
			err = t.checkout(repositoryUrl, namespace, name, version, reposDir)
			if err != nil {
				log.Error(err)
				return errors.Errorf("checkout package failed: %v", depend)
			}
		}
	}
	return nil
}
