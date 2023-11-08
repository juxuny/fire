package executor

type Type string

const (
	TypeBash = "bash"
	TypeSh   = "sh"
	TypeSsh  = "ssh"
)

type IExecutor interface {
	Start(args ...string) error
	StartAndWait(args ...string) error
}
