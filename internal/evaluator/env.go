package evaluator

import (
	"fmt"
	"rpucella.net/ragnarok/internal/value"
	"strings"
)

type Env struct {
	bindings map[string]value.Value
	previous *Env
	modules  map[string](*Env)
}

const moduleSep = "::"

func NewEnv(bindings map[string]value.Value, previous *Env, modules map[string]*Env) *Env {
	return &Env{bindings, previous, modules}
}

func (env *Env) find(name string) (value.Value, error) {
	if strings.Contains(name, moduleSep) {
		subnames := strings.Split(name, moduleSep)
		if len(subnames) > 2 {
			return nil, fmt.Errorf("multiple qualifiers in %s", name)
		}
		return env.Lookup(subnames[0], subnames[1])
	}
	current := env
	for current != nil {
		val, ok := current.bindings[name]
		if ok {
			return val, nil
		}
		current = current.previous
	}
	// can't find it, so look for it in the search path modules
	lookup_path, err := env.Lookup("config", "lookup-path")
	if err != nil || !value.IsRef(lookup_path) {
		return nil, fmt.Errorf("no such identifier %s", name)
	}
	modules := lookup_path.GetValue()
	for value.IsCons(modules) {
		if value.IsSymbol(modules.GetHead()) {
			result, err := env.Lookup(modules.GetHead().GetString(), name)
			if err == nil {
				return result, nil
			}
		}
		modules = modules.GetTail()
	}
	return nil, fmt.Errorf("no such identifier %s", name)
}

func (env *Env) Lookup(module string, name string) (value.Value, error) {
	moduleEnv, ok := env.modules[module]
	if !ok {
		return nil, fmt.Errorf("no such module %s", module)
	}
	v, ok := moduleEnv.bindings[name]
	if !ok {
		return nil, fmt.Errorf("no such identifier %s", name)
	}
	return v, nil
}

func (env *Env) Update(name string, v value.Value) {
	env.bindings[name] = v
}

func (env *Env) Layer(names []string, values []value.Value) *Env {
	// if values is nil or smaller than names, then
	// remaining names are bound to #nil
	bindings := map[string]value.Value{}
	for i, name := range names {
		if values != nil && i < len(values) {
			bindings[name] = values[i]
		} else {
			bindings[name] = &value.VNil{}
		}
	}
	return &Env{bindings: bindings, previous: env, modules: env.modules}
}
