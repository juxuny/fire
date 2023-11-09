package task

import (
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/yuanjiecloud/fire/log"
)

type Provider struct {
	dependencies  []string
	replaceMapper map[string]Replacement
}

func NewProvider(dependencies []string, replaceList []Replacement) (ret *Provider) {
	ret = &Provider{
		dependencies:  dependencies,
		replaceMapper: make(map[string]Replacement),
	}
	for _, replace := range replaceList {
		ret.replaceMapper[replace.Package] = replace
	}
	return
}

func (t *Provider) FindPipeline(pipelineName string) (ret *Pipeline, err error) {
	return nil, nil
}

func (t *Provider) getRepositoryMapperFromLocal(dir string) (result ReposMapper, err error) {
	wd := Getwd()
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("enter dir: ", dir)
	defer func() {
		err = os.Chdir(wd)
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("goback: ", wd)
	}()
	pipeline, err := Parse(DefaultConfigFile)
	if err != nil {
		log.Error(err)
		return nil, errors.Errorf("invalid repository dir: %v", dir)
	}
	return pipeline.GetRepositoryMapper()
}

func (t *Provider) GetRepositoryMapper() (ret ReposMapper, err error) {
	ret = NewReposMapper()
	for _, depend := range t.dependencies {
		log.Debug("resolving dependency: ", depend)
		var namespace, name, version string
		namespace, name, version, err = SplitPackageName(depend)
		if err != nil {
			log.Fatal("invalid repository:", depend)
			return
		}
		replacement, b := t.replaceMapper[name]
		if b {
			if replacement.IsLocal() {
				repositoryDir := path.Join(Getwd(), replacement.Repository)
				var mapperFromDependency ReposMapper
				mapperFromDependency, err = t.getRepositoryMapperFromLocal(repositoryDir)
				if err != nil {
					return
				}
				ret[depend] = repositoryDir
				ret = ret.MergeIgnoreDuplicated(mapperFromDependency)
			} else {
				repositoryDir := CreateRepositoryLocationSpecificVersion(namespace, name, replacement.Version.String())
				var mapperFromDependency ReposMapper
				mapperFromDependency, err = t.getRepositoryMapperFromLocal(repositoryDir)
				if err != nil {
					return
				}
				ret[depend] = repositoryDir
				ret = ret.MergeIgnoreDuplicated(mapperFromDependency)
			}
		} else {
			log.Debug("no replacement dependency: ", name)
			// use default location
			repositoryDir := CreateRepositoryLocationSpecificVersion(namespace, name, version)
			var mapperFromDependency ReposMapper
			mapperFromDependency, err = t.getRepositoryMapperFromLocal(repositoryDir)
			if err != nil {
				return
			}
			ret[depend] = repositoryDir
			ret = ret.MergeIgnoreDuplicated(mapperFromDependency)
		}
	}
	return
}
