package main

import "errors"

type Env struct {
	bindings map[string]Value
	previous *Env
}

func min (a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (env *Env) lookup(name string) (Value, error) {
	current := env
	for current != nil {
		val, ok := current.bindings[name]
		if ok {
			return val, nil
		}
		current = current.previous
	}
	return nil, errors.New("no such identifier " + name)
}

func (env *Env) update(name string, v Value) {
	env.bindings[name] = v
}

func (env *Env) layer(names []string, values []Value) *Env {
	bindings := map[string]Value{}
	len := min(len(names), len(values))
	for i := 0; i < len; i++  {
		bindings[names[i]] = values[i]
	}
	return &Env{bindings: bindings, previous: env}
}

