package task

type Context struct {
	Parent             *Context
	EnvProvider        EnvProvider
	Env                string
	RepositoryProvider *Provider
}

func WrapContext(parent *Context, child *Context) *Context {
	child.Parent = parent
	return child
}

func (t *Context) GetEnv(name string) (env Environment, found bool) {
	env, found = t.EnvProvider[name]
	if !found && t.Parent != nil {
		return t.Parent.GetEnv(name)
	}
	return
}

func (t *Context) RunTask(name string) error {
	return nil
}
