package runtime

import "github.com/mattn/anko/env"

type Runtime struct {
	env *env.Env
	ctx *Context
}

func NewRuntime() *Runtime {
	r := &Runtime{
		env: env.NewEnv(),
		ctx: &Context{
			Spec: &AGSSpec{Contexts: []*ContextSpec{}},
		},
	}

	r.registerBuiltins()

	return r
}

func (r *Runtime) GetProject() *Project {
	return r.ctx.Project
}

func (r *Runtime) GetSpec() *AGSSpec {
	return r.ctx.Spec
}
