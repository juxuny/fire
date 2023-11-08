package task

import "strings"

type Replacement struct {
	Package    string  `json:"package" yaml:"package"`
	Version    Version `json:"version" yaml:"version"`
	Repository string  `json:"repository" yaml:"repository"`
}

func (t Replacement) IsLocal() bool {
	if strings.Index(t.Repository, "http://") == 0 || strings.Index(t.Repository, "https://") == 0 || strings.Index(t.Repository, "git@") == 0 {
		return false
	}
	return true
}
