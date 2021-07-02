package main

import (
	   "fmt"
	   "rpucella.net/ragnarok/internal/lisp"
)

// An ecosystem is a global set of environments associated with "modules"

type Ecosystem struct {
	modulesEnv map[string]*lisp.Env
	activesEnv map[string]*lisp.Env
}

func mkEcosystem() *Ecosystem {
	return &Ecosystem{map[string]*lisp.Env{}, map[string]*lisp.Env{}}
}

func (eco *Ecosystem) get(name string) (*lisp.Env, error) {
	env, ok := eco.activesEnv[name]
	if !ok {
		return nil, fmt.Errorf("Cannot switch to module %s", name)
	}
	return env, nil
}

func (eco *Ecosystem) mkEnv(name string, bindings map[string]lisp.Value) *lisp.Env {
	env := lisp.NewEnv(bindings, nil, eco.modulesEnv)
	eco.modulesEnv[name] = env
	eco.activesEnv[name] = env.Layer([]string{}, []lisp.Value{})
	return env
}
