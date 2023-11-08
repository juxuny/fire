package task

type Version string

func (t Version) String() string {
	return string(t)
}
