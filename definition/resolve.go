package definition

import (
	"fmt"
	"github.com/project-flogo/core/data/path"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
)

var defResolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{
	".":         &resolve.ScopeResolver{},
	"env":       &resolve.EnvResolver{},
	"property":  &resolve.PropertyResolver{},
	"loop":      &resolve.LoopResolver{},
	"iteration": &IteratorResolver{}, //todo should we create a separate resolver to use in iterations?
	"activity":  &ActivityResolver{},
	"error":     &ErrorResolver{},
	"flow":      &FlowResolver{}})

func GetDataResolver() resolve.CompositeResolver {
	return defResolver
}

var resolverInfo = resolve.NewResolverInfo(false, false)

type FlowResolver struct {
}

func (r *FlowResolver) GetResolverInfo() *resolve.ResolverInfo {
	return resolverInfo
}

func (r *FlowResolver) Resolve(scope data.Scope, itemName, valueName string) (interface{}, error) {

	value, exists := scope.GetValue(valueName)
	if !exists {
		return nil, fmt.Errorf("failed to resolve flow attr: '%s', not found in flow", valueName)
	}

	return value, nil
}

var dynamicItemResolver = resolve.NewResolverInfo(false, true)

type ActivityResolver struct {
}

func (r *ActivityResolver) GetResolverInfo() *resolve.ResolverInfo {
	return dynamicItemResolver
}

func (r *ActivityResolver) Resolve(scope data.Scope, itemName, valueName string) (interface{}, error) {

	value, exists := scope.GetValue("_A." + itemName + "." + valueName)
	if !exists {
		return nil, fmt.Errorf("failed to resolve activity attr: '%s', not found in activity '%s'", valueName, itemName)
	}

	return value, nil
}

type ErrorResolver struct {
}

func (r *ErrorResolver) GetResolverInfo() *resolve.ResolverInfo {
	return resolverInfo
}

func (r *ErrorResolver) Resolve(scope data.Scope, itemName, valueName string) (interface{}, error) {

	err, exists := scope.GetValue("_E")
	if !exists {
		return nil, fmt.Errorf("failed to resolve error, not found in flow")
	}

	errObj, ok := err.(map[string]interface{})
	if ok {
		value, ok := errObj[valueName]
		if !ok {
			return nil, nil
		}
		return value, nil
	} else {
		return nil, fmt.Errorf("invalid error object: %v", errObj)
	}
}

type IteratorResolver struct {
}

func (*IteratorResolver) GetResolverInfo() *resolve.ResolverInfo {
	return dynamicItemResolver
}

//Resolve resolved iterator value using  the following syntax:  $iteration[key], or $iteration[value]
func (*IteratorResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	value, exists := scope.GetValue("_W.iteration")
	if !exists {
		return nil, fmt.Errorf("failed to resolve iteration value, not in an iterator")
	}
	return path.GetValue(value, "."+item)
}
