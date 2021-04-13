package main

import "errors"

type Env struct {
	bindings map[string]Value
	previous *Env
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

func newEnv(env *Env) (*Env) {
	return &Env{bindings: map[string]Value{}, previous: env}
}
