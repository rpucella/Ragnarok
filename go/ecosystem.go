package main

import "fmt"

// An ecosystem is a global set of environments associated with "modules"

type Ecosystem struct {
	modulesEnv map[string]*Env
	activesEnv map[string]*Env
}

func mkEcosystem() *Ecosystem {
	return &Ecosystem{map[string]*Env{}, map[string]*Env{}}
}

func (eco *Ecosystem) get(name string) (*Env, error) {
	env, ok := eco.activesEnv[name]
	if !ok {
		return nil, fmt.Errorf("Cannot switch to module %s", name)
	}
	return env, nil
}

func (eco *Ecosystem) add(name string, env *Env) {
	eco.modulesEnv[name] = env
	eco.activesEnv[name] = env.layer([]string{}, []Value{})
}
