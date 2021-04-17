package main

import "fmt"
import "strings"

type Env struct {
	bindings map[string]Value
	previous *Env
	ecosystem *Ecosystem
}

func (env *Env) lookup(name string) (Value, error) {
	if strings.Contains(name, ":") {
		subnames := strings.Split(name, ":")
		if len(subnames) > 2 {
			return nil, fmt.Errorf("multiple qualifiers in %s", name)
		}
		moduleEnv, ok := env.ecosystem.modulesEnv[subnames[0]]
		if !ok {
			return nil, fmt.Errorf("unknown module %s", subnames[0])
		}
		return moduleEnv.lookup(subnames[1])
	}
	current := env
	for current != nil {
		val, ok := current.bindings[name]
		if ok {
			return val, nil
		}
		current = current.previous
	}
	return nil, fmt.Errorf("no such identifier %s", name)
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

