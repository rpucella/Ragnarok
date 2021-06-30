package main

import "fmt"
import "strings"

type Env struct {
	bindings map[string]Value
	previous *Env
	ecosystem *Ecosystem
}

const moduleSep = "::"

func (env *Env) find(name string) (Value, error) {
	if strings.Contains(name, moduleSep) {
		subnames := strings.Split(name, moduleSep)
		if len(subnames) > 2 {
			return nil, fmt.Errorf("multiple qualifiers in %s", name)
		}
		return env.lookup(subnames[0], subnames[1])
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
	lookup_path, err := env.lookup("config", "lookup-path")
	if err != nil || !lookup_path.isRef() {
		return nil, fmt.Errorf("no such identifier %s", name)
	}
	modules := lookup_path.getValue()
	for modules.isCons() {
		if modules.headValue().isSymbol() { 
			result, err := env.lookup(modules.headValue().strValue(), name)
			if err == nil {
				return result, nil
			}
		}
		modules = modules.tailValue()
	}
	return nil, fmt.Errorf("no such identifier %s", name)
}

func (env *Env) lookup(module string, name string) (Value, error) {
	moduleEnv, ok := env.ecosystem.modulesEnv[module]
	if !ok {
		return nil, fmt.Errorf("no such module %s", module)
	}
	v, ok := moduleEnv.bindings[name]
	if !ok {
		return nil, fmt.Errorf("no such identifier %s", name)
	}
	return v, nil
}
 
func (env *Env) update(name string, v Value) {
	env.bindings[name] = v
}

func (env *Env) layer(names []string, values []Value) *Env {
	// if values is nil or smaller than names, then
	// remaining names are bound to #nil
	bindings := map[string]Value{}
	for i, name := range names {
		if values != nil && i < len(values) {
			bindings[name] = values[i]
		} else {
			bindings[name] = &VNil{}
		}
	}
	return &Env{bindings: bindings, previous: env, ecosystem: env.ecosystem}
}

