package runtime

import "fmt"

func (r *Runtime) registerBuiltins() {

	r.env.Define("project", func(name string, fn func()) {
		r.ctx.Project = &Project{
			Name:     name,
			Contexts: []*BoundedContext{},
		}

		fn()
	})

	r.env.Define("bounded_context", func(name string, fn func()) {

		if r.ctx.Project == nil {
			panic("bounded_context must be inside a project")
		}

		bc := &BoundedContext{
			Name:        name,
			Subcontexts: []string{},
		}

		// Guardar contexto anterior (por si luego anidas en el futuro)
		prev := r.ctx.CurrentContext
		r.ctx.CurrentContext = bc

		fn()

		r.ctx.Project.Contexts = append(r.ctx.Project.Contexts, bc)

		// Restaurar contexto anterior
		r.ctx.CurrentContext = prev
	})

	r.env.Define("subcontext", func(name string) {

		if r.ctx.CurrentContext == nil {
			panic("subcontext must be inside a bounded_context")
		}

		r.ctx.CurrentContext.Subcontexts = append(
			r.ctx.CurrentContext.Subcontexts,
			name,
		)
	})

	r.env.Define("hexagonal", func() {

		if r.ctx.CurrentContext == nil {
			panic("hexagonal must be inside a bounded_context")
		}

		r.ctx.CurrentContext.Hexagonal = true
	})

	// 🔥 BONUS: lenguaje (te sirve luego para templates)
	r.env.Define("language", func(lang string) {
		if r.ctx.Project == nil {
			panic("language must be inside a project")
		}
		r.ctx.Project.Language = lang
	})

	// 🔥 BONUS: debug (muy útil ahora)
	r.env.Define("debug", func() {
		fmt.Printf("%+v\n", r.ctx.Project)
	})
}
