package runtime

type Project struct {
	Name     string
	Language string
	Contexts []*BoundedContext
}

type BoundedContext struct {
	Name        string
	Hexagonal   bool
	Subcontexts []string
}

type AGSSpec struct {
	Contexts []*ContextSpec
}

type ContextSpec struct {
	Name     string
	Shared   *SharedSpec
	Modules  []*ModuleSpec
	Commands []*CommandSpec
}

type SharedSpec struct {
	ValueObjects []*ValueObjectSpec
}

type ModuleSpec struct {
	Name         string
	ValueObjects []*ValueObjectSpec
	Entities     []*EntitySpec
	Aggregates   []*AggregateSpec
	Providers    []*ProviderSpec
}

type EntitySpec struct {
	Name         string
	ValueObjects []*ValueObjectSpec
	Abilities    []*AbilitySpec
}

type AggregateSpec struct {
	Name         string
	ID           string
	ValueObjects []*ValueObjectSpec
	Collections  []string
	Abilities    []*AbilitySpec
	Behaviors    []*BehaviorSpec
}

type ValueObjectSpec struct {
	Name string
	Type string
}

type AbilitySpec struct {
	Name   string
	Params []string
}

type BehaviorSpec struct {
	Name   string
	Params []string
	Emits  []string
}

type ProviderSpec struct {
	Name    string
	Methods []*ProviderMethodSpec
}

type ProviderMethodSpec struct {
	Name   string
	Params []string
	Return string
}

type CommandSpec struct {
	Name       string
	Params     []string
	Middleware []string
	Handler    string
}

type Context struct {
	Project        *Project
	CurrentContext *BoundedContext
	Spec           *AGSSpec
}
