package runtime

import (
	"fmt"
	"strings"
)

func (r *Runtime) validateSpec() error {
	if r.ctx == nil || r.ctx.Spec == nil {
		return nil
	}

	errList := make([]*AGSError, 0)

	for _, ctxSpec := range r.ctx.Spec.Contexts {
		aggregatesByName := map[string]*AggregateSpec{}

		for _, module := range ctxSpec.Modules {
			entityNames := map[string]bool{}
			for _, e := range module.Entities {
				entityNames[e.Name] = true
			}

			for _, agg := range module.Aggregates {
				aggPath := fmt.Sprintf("context/%s/module/%s/aggregate/%s", ctxSpec.Name, module.Name, agg.Name)
				if _, exists := aggregatesByName[agg.Name]; exists {
					errList = append(errList, newAGSError(ErrDuplicateSymbol, fmt.Sprintf("duplicate aggregate name %q in context %s", agg.Name, ctxSpec.Name), aggPath))
				} else {
					aggregatesByName[agg.Name] = agg
				}

				if strings.TrimSpace(agg.ID) == "" {
					errList = append(errList, newAGSError(ErrSymbolNotFound, "aggregate does not define id(voName)", aggPath+"/id"))
				} else if !r.voVisibleInAggregate(ctxSpec, module, agg, agg.ID) {
					errList = append(errList, newAGSError(ErrSymbolNotFound, fmt.Sprintf("id vo %q not found", agg.ID), aggPath+"/id"))
				}

				for _, entityName := range agg.Collections {
					if !entityNames[entityName] {
						errList = append(errList, newAGSError(ErrSymbolNotFound, fmt.Sprintf("collection entity %q not found", entityName), aggPath+"/collection"))
					}
				}

				seenEvents := map[string]bool{}
				for _, behavior := range agg.Behaviors {
					for _, eventName := range behavior.Emits {
						if seenEvents[eventName] {
							errList = append(errList, newAGSError(ErrDuplicateSymbol, fmt.Sprintf("duplicate emitted event %q", eventName), aggPath+"/behavior/"+behavior.Name+"/emits"))
						}
						seenEvents[eventName] = true
					}
				}
			}

			for _, provider := range module.Providers {
				for _, method := range provider.Methods {
					if strings.TrimSpace(method.Return) == "" {
						continue
					}
					if !isValidAGSBaseType(method.Return) && !r.voVisibleInModule(ctxSpec, module, method.Return) {
						errList = append(errList, newAGSError(ErrTypeInvalid, fmt.Sprintf("return type %q is not a known AGS type or VO", method.Return), fmt.Sprintf("context/%s/module/%s/provider/%s/method/%s/returns", ctxSpec.Name, module.Name, provider.Name, method.Name)))
					}
				}
			}
		}

		for _, cmd := range ctxSpec.Commands {
			path := fmt.Sprintf("context/%s/command/%s", ctxSpec.Name, cmd.Name)
			if strings.TrimSpace(cmd.Handler) == "" {
				errList = append(errList, newAGSError(ErrHandlerUnresolved, "command has empty handler", path+"/handler"))
				continue
			}

			parts := strings.Split(cmd.Handler, ".")
			if len(parts) != 2 {
				errList = append(errList, newAGSError(ErrHandlerUnresolved, fmt.Sprintf("handler %q must be Aggregate.method", cmd.Handler), path+"/handler"))
				continue
			}

			aggName := parts[0]
			methodName := parts[1]
			agg := aggregatesByName[aggName]
			if agg == nil {
				errList = append(errList, newAGSError(ErrHandlerUnresolved, fmt.Sprintf("aggregate %q not found for command", aggName), path+"/handler"))
				continue
			}

			if !aggregateHasMethod(agg, methodName) {
				errList = append(errList, newAGSError(ErrHandlerUnresolved, fmt.Sprintf("method %q not found in aggregate %s", methodName, aggName), path+"/handler"))
			}
		}
	}

	if len(errList) > 0 {
		return &AGSMultiError{Errors: errList}
	}

	return nil
}

func (r *Runtime) voVisibleInAggregate(ctxSpec *ContextSpec, module *ModuleSpec, agg *AggregateSpec, voName string) bool {
	for _, vo := range agg.ValueObjects {
		if vo.Name == voName {
			return true
		}
	}
	return r.voVisibleInModule(ctxSpec, module, voName)
}

func (r *Runtime) voVisibleInModule(ctxSpec *ContextSpec, module *ModuleSpec, voName string) bool {
	for _, vo := range module.ValueObjects {
		if vo.Name == voName {
			return true
		}
	}
	if ctxSpec.Shared != nil {
		for _, vo := range ctxSpec.Shared.ValueObjects {
			if vo.Name == voName {
				return true
			}
		}
	}
	return false
}

func aggregateHasMethod(agg *AggregateSpec, methodName string) bool {
	for _, behavior := range agg.Behaviors {
		if behavior.Name == methodName {
			return true
		}
	}
	for _, ability := range agg.Abilities {
		if ability.Name == methodName {
			return true
		}
	}
	return false
}
