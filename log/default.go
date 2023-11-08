package log

import (
	"fmt"
	"os"
)

func Debug(data ...interface{}) {
	fmt.Println(data...)
}

func Fatal(data ...interface{}) {
	fmt.Println(data...)
	os.Exit(-1)
}

func Error(data ...interface{}) {
	os.Stderr.WriteString(fmt.Sprint(data...) + "\n")
}
