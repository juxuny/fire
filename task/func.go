package task

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/yuanjiecloud/fire/log"
)

var (
	globalCacheDir   string
	globalReposDir   string
	globalFireConfig string
)

func GetGlobalCacheDir() (fireCacheDir string, err error) {
	if globalCacheDir != "" {
		return globalCacheDir, nil
	}
	fireCacheDir, err = os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fireCacheDir = path.Join(fireCacheDir, ".config", "fire")
	_, err = os.Stat(fireCacheDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(fireCacheDir, 0775)
		if err != nil {
			log.Fatal(err)
			return "", errors.Errorf("get cache dir failed")
		}
	}
	globalCacheDir = fireCacheDir
	return fireCacheDir, nil
}

func GetGlobalReposDir() (reposDir string, err error) {
	if globalReposDir != "" {
		return globalReposDir, nil
	}
	reposDir, err = GetGlobalCacheDir()
	if err != nil {
		log.Error(err)
		return "", errors.Errorf("get repos dir failed")
	}
	reposDir = path.Join(reposDir, "repos")
	_, err = os.Stat(reposDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(reposDir, 0775)
		if err != nil {
			log.Fatal(err)
			return "", errors.Errorf("create repository dir failed")
		}
	}
	globalReposDir = reposDir
	return
}

func GetGlobalFireConfig() (configFile string, err error) {
	if globalFireConfig != "" {
		return globalFireConfig, nil
	}
	cacheDir, err := GetGlobalCacheDir()
	if err != nil {
		return "", err
	}
	configFile = path.Join(cacheDir, "fire.yaml")
	globalFireConfig = configFile
	return
}

func SplitPackageName(packageName string) (namespace, name, version string, err error) {
	if packageName == "" {
		err = errors.Errorf("invalid package: %v", packageName)
		return
	}
	l := strings.Split(packageName, "@")
	if len(l) == 0 {
		err = errors.Errorf("invalid package: %v", packageName)
	}
	if len(l) == 2 {
		version = l[1]
	}
	l = strings.Split(l[0], "/")
	if len(l) == 1 {
		name = l[0]
	} else if len(l) == 2 {
		namespace = l[0]
		name = l[1]
	} else {
		err = errors.Errorf("invalid package: %v", packageName)
	}
	return
}

func CheckIfExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
