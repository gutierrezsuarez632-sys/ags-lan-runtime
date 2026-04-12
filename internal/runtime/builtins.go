package runtime

import (
	"fmt"
	"strings"
)

func (r *Runtime) registerBuiltins() {

	// Legacy DSL support
	r.env.Define("project", func(name string, fn func()) {
		r.ctx.Project = &Project{
			Name:     name,
			Contexts: []*BoundedContext{},
		}

		fn()
	})

	r.env.Define("bounded_context", func(name string, fn func()) {
		if r.ctx.Project == nil {
			r.raise(ErrScopeInvalid, "bounded_context must be inside a project", "legacy/bounded_context")
		}

		bc := &BoundedContext{
			Name:        name,
			Subcontexts: []string{},
		}

		prev := r.ctx.CurrentContext
		r.ctx.CurrentContext = bc
		fn()
		r.ctx.Project.Contexts = append(r.ctx.Project.Contexts, bc)
		r.ctx.CurrentContext = prev
	})

	r.env.Define("subcontext", func(name string) {
		if r.ctx.CurrentContext == nil {
			r.raise(ErrScopeInvalid, "subcontext must be inside a bounded_context", "legacy/subcontext")
		}

		r.ctx.CurrentContext.Subcontexts = append(r.ctx.CurrentContext.Subcontexts, name)
	})

	r.env.Define("hexagonal", func() {
		if r.ctx.CurrentContext == nil {
			r.raise(ErrScopeInvalid, "hexagonal must be inside a bounded_context", "legacy/hexagonal")
		}

		r.ctx.CurrentContext.Hexagonal = true
	})

	r.env.Define("language", func(lang string) {
		if r.ctx.Project == nil {
			r.raise(ErrScopeInvalid, "language must be inside a project", "legacy/language")
		}
		r.ctx.Project.Language = lang
	})

	// AGS v0.1 fluent DSL
	r.env.Define("createContext", func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "context name cannot be empty", "context")
		}
		if r.findContext(name) != nil {
			r.raise(ErrDuplicateSymbol, "duplicate context: "+name, "context/"+name)
		}

		ctxSpec := &ContextSpec{
			Name:     name,
			Shared:   &SharedSpec{ValueObjects: []*ValueObjectSpec{}},
			Modules:  []*ModuleSpec{},
			Commands: []*CommandSpec{},
		}

		r.ctx.Spec.Contexts = append(r.ctx.Spec.Contexts, ctxSpec)
		if r.ctx.Project == nil {
			r.ctx.Project = &Project{Name: name, Contexts: []*BoundedContext{}}
		}

		return r.contextObject(ctxSpec)
	})

	r.env.Define("debug", func() {
		fmt.Printf("legacy=%+v\n", r.ctx.Project)
		fmt.Printf("spec=%+v\n", r.ctx.Spec)
	})
}

func (r *Runtime) contextObject(ctxSpec *ContextSpec) map[string]interface{} {
	obj := map[string]interface{}{}

	obj["shared"] = func() map[string]interface{} {
		if ctxSpec.Shared == nil {
			ctxSpec.Shared = &SharedSpec{ValueObjects: []*ValueObjectSpec{}}
		}
		return r.sharedObject(ctxSpec.Shared)
	}

	obj["module"] = func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "module name cannot be empty", "context/"+ctxSpec.Name+"/module")
		}
		if r.findModule(ctxSpec, name) != nil {
			r.raise(ErrDuplicateSymbol, "duplicate module in context "+ctxSpec.Name+": "+name, "context/"+ctxSpec.Name+"/module/"+name)
		}

		module := &ModuleSpec{
			Name:         name,
			ValueObjects: []*ValueObjectSpec{},
			Entities:     []*EntitySpec{},
			Aggregates:   []*AggregateSpec{},
			Providers:    []*ProviderSpec{},
		}
		ctxSpec.Modules = append(ctxSpec.Modules, module)

		if r.ctx.Project != nil {
			r.ctx.Project.Contexts = append(r.ctx.Project.Contexts, &BoundedContext{Name: name, Subcontexts: []string{}})
		}

		return r.moduleObject(ctxSpec, module)
	}

	obj["command"] = func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "command name cannot be empty", "context/"+ctxSpec.Name+"/command")
		}
		if r.findCommand(ctxSpec, name) != nil {
			r.raise(ErrDuplicateSymbol, "duplicate command in context "+ctxSpec.Name+": "+name, "context/"+ctxSpec.Name+"/command/"+name)
		}

		cmd := &CommandSpec{Name: name, Params: []string{}, Middleware: []string{}}
		ctxSpec.Commands = append(ctxSpec.Commands, cmd)
		return r.commandObject(cmd)
	}

	return obj
}

func (r *Runtime) sharedObject(shared *SharedSpec) map[string]interface{} {
	obj := map[string]interface{}{}
	obj["vo"] = func(name string, voType string) map[string]interface{} {
		r.addVO(&shared.ValueObjects, name, voType, "shared")
		return obj
	}
	return obj
}

func (r *Runtime) moduleObject(ctxSpec *ContextSpec, module *ModuleSpec) map[string]interface{} {
	obj := map[string]interface{}{}

	obj["vo"] = func(name string, voType string) map[string]interface{} {
		r.addVO(&module.ValueObjects, name, voType, "module")
		return obj
	}

	obj["entity"] = func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "entity name cannot be empty", "module/"+module.Name+"/entity")
		}
		if r.findEntity(module, name) != nil {
			r.raise(ErrDuplicateSymbol, "duplicate entity in module "+module.Name+": "+name, "module/"+module.Name+"/entity/"+name)
		}

		entity := &EntitySpec{Name: name, ValueObjects: []*ValueObjectSpec{}, Abilities: []*AbilitySpec{}}
		module.Entities = append(module.Entities, entity)
		return r.entityObject(entity)
	}

	obj["aggregate"] = func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "aggregate name cannot be empty", "module/"+module.Name+"/aggregate")
		}
		if r.findAggregate(module, name) != nil {
			r.raise(ErrDuplicateSymbol, "duplicate aggregate in module "+module.Name+": "+name, "module/"+module.Name+"/aggregate/"+name)
		}

		agg := &AggregateSpec{
			Name:         name,
			ValueObjects: []*ValueObjectSpec{},
			Collections:  []string{},
			Abilities:    []*AbilitySpec{},
			Behaviors:    []*BehaviorSpec{},
		}
		module.Aggregates = append(module.Aggregates, agg)
		return r.aggregateObject(module, agg)
	}

	obj["provider"] = func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "provider name cannot be empty", "module/"+module.Name+"/provider")
		}
		if r.findProvider(module, name) != nil {
			r.raise(ErrDuplicateSymbol, "duplicate provider in module "+module.Name+": "+name, "module/"+module.Name+"/provider/"+name)
		}

		provider := &ProviderSpec{Name: name, Methods: []*ProviderMethodSpec{}}
		module.Providers = append(module.Providers, provider)
		return r.providerObject(provider)
	}

	_ = ctxSpec
	return obj
}

func (r *Runtime) entityObject(entity *EntitySpec) map[string]interface{} {
	obj := map[string]interface{}{}

	obj["vo"] = func(name string, voType string) map[string]interface{} {
		r.addVO(&entity.ValueObjects, name, voType, "entity")
		return obj
	}

	obj["ability"] = func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "ability name cannot be empty", "entity/"+entity.Name+"/ability")
		}
		ability := &AbilitySpec{Name: name, Params: []string{}}
		entity.Abilities = append(entity.Abilities, ability)
		return r.abilityObject(ability)
	}

	return obj
}

func (r *Runtime) aggregateObject(module *ModuleSpec, agg *AggregateSpec) map[string]interface{} {
	obj := map[string]interface{}{}

	obj["vo"] = func(name string, voType string) map[string]interface{} {
		r.addVO(&agg.ValueObjects, name, voType, "aggregate")
		return obj
	}

	obj["id"] = func(voName string) map[string]interface{} {
		if strings.TrimSpace(voName) == "" {
			r.raise(ErrScopeInvalid, "id voName cannot be empty", "aggregate/"+agg.Name+"/id")
		}
		agg.ID = voName
		return obj
	}

	obj["collection"] = func(entityName string) map[string]interface{} {
		if strings.TrimSpace(entityName) == "" {
			r.raise(ErrScopeInvalid, "collection entity name cannot be empty", "aggregate/"+agg.Name+"/collection")
		}
		agg.Collections = append(agg.Collections, entityName)
		return obj
	}

	obj["ability"] = func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "ability name cannot be empty", "aggregate/"+agg.Name+"/ability")
		}
		agg.Abilities = append(agg.Abilities, &AbilitySpec{Name: name, Params: []string{}})
		return r.abilityObject(agg.Abilities[len(agg.Abilities)-1])
	}

	obj["behavior"] = func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "behavior name cannot be empty", "aggregate/"+agg.Name+"/behavior")
		}
		behavior := &BehaviorSpec{Name: name, Params: []string{}, Emits: []string{}}
		agg.Behaviors = append(agg.Behaviors, behavior)
		return r.behaviorObject(behavior)
	}

	_ = module
	return obj
}

func (r *Runtime) abilityObject(ability *AbilitySpec) map[string]interface{} {
	obj := map[string]interface{}{}
	obj["params"] = func(params ...string) map[string]interface{} {
		ability.Params = append(ability.Params, params...)
		return obj
	}
	return obj
}

func (r *Runtime) behaviorObject(behavior *BehaviorSpec) map[string]interface{} {
	obj := map[string]interface{}{}
	obj["params"] = func(params ...string) map[string]interface{} {
		behavior.Params = append(behavior.Params, params...)
		return obj
	}
	obj["emits"] = func(eventName string) map[string]interface{} {
		if strings.TrimSpace(eventName) == "" {
			r.raise(ErrScopeInvalid, "event name cannot be empty", "behavior/"+behavior.Name+"/emits")
		}
		behavior.Emits = append(behavior.Emits, eventName)
		return obj
	}
	return obj
}

func (r *Runtime) providerObject(provider *ProviderSpec) map[string]interface{} {
	obj := map[string]interface{}{}
	obj["method"] = func(name string) map[string]interface{} {
		if strings.TrimSpace(name) == "" {
			r.raise(ErrScopeInvalid, "provider method name cannot be empty", "provider/"+provider.Name+"/method")
		}
		for _, m := range provider.Methods {
			if m.Name == name {
				r.raise(ErrDuplicateSymbol, "duplicate provider method in provider "+provider.Name+": "+name, "provider/"+provider.Name+"/method/"+name)
			}
		}
		method := &ProviderMethodSpec{Name: name, Params: []string{}}
		provider.Methods = append(provider.Methods, method)
		return r.providerMethodObject(method)
	}
	return obj
}

func (r *Runtime) providerMethodObject(method *ProviderMethodSpec) map[string]interface{} {
	obj := map[string]interface{}{}
	finalized := false
	obj["params"] = func(params ...string) map[string]interface{} {
		if finalized {
			r.raise(ErrChainInvalid, "params() cannot be called after returns()", "provider_method/"+method.Name)
		}
		method.Params = append(method.Params, params...)
		return obj
	}
	obj["returns"] = func(returnType string) map[string]interface{} {
		if finalized {
			r.raise(ErrChainInvalid, "returns() cannot be called more than once", "provider_method/"+method.Name)
		}
		if strings.TrimSpace(returnType) == "" {
			r.raise(ErrScopeInvalid, "return type cannot be empty", "provider_method/"+method.Name+"/returns")
		}
		method.Return = returnType
		finalized = true
		return obj
	}
	return obj
}

func (r *Runtime) commandObject(cmd *CommandSpec) map[string]interface{} {
	obj := map[string]interface{}{}
	finalized := false
	obj["params"] = func(params ...string) map[string]interface{} {
		if finalized {
			r.raise(ErrChainInvalid, "params() cannot be called after handler()", "command/"+cmd.Name)
		}
		cmd.Params = append(cmd.Params, params...)
		return obj
	}
	obj["middleware"] = func(middleware ...string) map[string]interface{} {
		if finalized {
			r.raise(ErrChainInvalid, "middleware() cannot be called after handler()", "command/"+cmd.Name)
		}
		cmd.Middleware = append(cmd.Middleware, middleware...)
		return obj
	}
	obj["handler"] = func(handler string) map[string]interface{} {
		if finalized {
			r.raise(ErrChainInvalid, "handler() cannot be called more than once", "command/"+cmd.Name)
		}
		if strings.TrimSpace(handler) == "" {
			r.raise(ErrScopeInvalid, "handler cannot be empty", "command/"+cmd.Name+"/handler")
		}
		cmd.Handler = handler
		finalized = true
		return obj
	}
	return obj
}

func (r *Runtime) addVO(target *[]*ValueObjectSpec, name string, voType string, scope string) {
	if strings.TrimSpace(name) == "" {
		r.raise(ErrScopeInvalid, "vo name cannot be empty", scope+"/vo")
	}
	if strings.TrimSpace(voType) == "" {
		r.raise(ErrScopeInvalid, "vo type cannot be empty", scope+"/vo/"+name)
	}
	if !isValidAGSBaseType(voType) {
		r.raise(ErrTypeInvalid, "invalid vo type in "+scope+": "+voType, scope+"/vo/"+name)
	}
	for _, existing := range *target {
		if existing.Name == name {
			r.raise(ErrDuplicateSymbol, "duplicate vo in "+scope+": "+name, scope+"/vo/"+name)
		}
	}
	*target = append(*target, &ValueObjectSpec{Name: name, Type: voType})
}

func (r *Runtime) raise(code ErrorCode, message string, semanticPath string) {
	panic(newAGSError(code, message, semanticPath))
}

func isValidAGSBaseType(voType string) bool {
	baseTypes := map[string]bool{
		"uuid":    true,
		"string":  true,
		"integer": true,
		"decimal": true,
		"email":   true,
		"json":    true,
	}

	if baseTypes[voType] {
		return true
	}

	if strings.HasPrefix(voType, "enum:") {
		values := strings.Split(strings.TrimPrefix(voType, "enum:"), ",")
		if len(values) < 2 {
			return false
		}
		for _, value := range values {
			if strings.TrimSpace(value) == "" {
				return false
			}
		}
		return true
	}

	return false
}

func (r *Runtime) findContext(name string) *ContextSpec {
	for _, c := range r.ctx.Spec.Contexts {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (r *Runtime) findModule(ctxSpec *ContextSpec, name string) *ModuleSpec {
	for _, m := range ctxSpec.Modules {
		if m.Name == name {
			return m
		}
	}
	return nil
}

func (r *Runtime) findEntity(module *ModuleSpec, name string) *EntitySpec {
	for _, e := range module.Entities {
		if e.Name == name {
			return e
		}
	}
	return nil
}

func (r *Runtime) findAggregate(module *ModuleSpec, name string) *AggregateSpec {
	for _, a := range module.Aggregates {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func (r *Runtime) findProvider(module *ModuleSpec, name string) *ProviderSpec {
	for _, p := range module.Providers {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (r *Runtime) findCommand(ctxSpec *ContextSpec, name string) *CommandSpec {
	for _, cmd := range ctxSpec.Commands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}
