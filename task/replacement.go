package task

type Replacement struct {
	Package    string  `json:"package" yaml:"package"`
	Version    Version `json:"version" yaml:"version"`
	Repository string  `json:"repository" yaml:"repository"`
}
