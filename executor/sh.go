package executor

import (
	"bytes"
	"fmt"
	"os"
)

func NewShExecutor(env map[string]string, scripts []string) IExecutor {
	in := bytes.NewBuffer(nil)
	if len(env) > 0 {
		for k, v := range env {
			in.WriteString(fmt.Sprintf("export %v=%v\n", k, v))
		}
	}
	for _, line := range scripts {
		if len(line) == 0 {
			continue
		}
		in.WriteString(line + "\n")
	}
	result := &bashExecutor{
		Binary: "/bin/sh",
		in:     in,
		out:    os.Stdout,
	}
	return result
}
