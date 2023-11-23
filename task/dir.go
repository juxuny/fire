package task

import (
	"container/list"
	"os"

	"github.com/yuanjiecloud/fire/log"
)

var history = list.New()

func enter(dir string) {
	wd := Getwd()
	err := os.Chdir(dir)
	log.CheckAndFatal(err)
	history.PushBack(wd)
	log.Debug("enter dir: ", dir)
}

func goback() {
	if history.Len() == 0 {
		return
	}
	back := history.Back()
	history.Remove(back)
	err := os.Chdir(back.Value.(string))
	log.CheckAndFatal(err)
	log.Debug("goback: ", back.Value)
}
