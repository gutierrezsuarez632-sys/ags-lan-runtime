package runtime

import "github.com/mattn/anko/env"

type Runtime struct {
	env *env.Env
	ctx *Context
}

func NewRuntime() *Runtime {
	r := &Runtime{
		env: env.NewEnv(),
		ctx: &Context{},
	}

	r.registerBuiltins()

	return r
}

func (r *Runtime) GetProject() *Project {
	return r.ctx.Project
}
