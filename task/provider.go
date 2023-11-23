package task

import (
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/yuanjiecloud/fire/log"
)

var mapFromPipelineToLocation = make(map[string]string)
var pipelineMapper = make(map[string]*Pipeline)

func AddPipeline(pipelineWithVersion string, dir string) (*Pipeline, error) {
	workdirBackup := Getwd()
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("enter dir: ", dir)
	defer func() {
		err = os.Chdir(workdirBackup)
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("goback: ", workdirBackup)
	}()
	configFile := path.Join(dir, DefaultConfigFile)
	pipeline, err := Parse(configFile)
	if err != nil {
		return nil, errors.Errorf("invalid fire project: %s", dir)
	}
	log.Debug("add pipeline: ", pipelineWithVersion, " => ", dir)
	mapFromPipelineToLocation[pipelineWithVersion] = dir
	pipelineMapper[pipelineWithVersion] = pipeline
	return pipeline, nil
}

func FindPipeline(pipelineWithVersion string) (pipeline *Pipeline, found bool) {
	pipeline, found = pipelineMapper[pipelineWithVersion]
	return
}

func CheckIfContainPipeline(pipelineWithVersion string) bool {
	_, found := pipelineMapper[pipelineWithVersion]
	return found
}

func FindPipelineReposDir(pipelineWithVersion string) (dir string, found bool) {
	dir, found = mapFromPipelineToLocation[pipelineWithVersion]
	return
}
