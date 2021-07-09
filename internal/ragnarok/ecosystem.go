package ragnarok

import (
	"fmt"
	"rpucella.net/ragnarok/internal/evaluator"
	"rpucella.net/ragnarok/internal/value"
)

// An ecosystem is a global set of environments associated with modules, shells, or buffers

type Ecosystem struct {
	modules map[string]*evaluator.Env
	shells  map[string]*evaluator.Env
	buffers map[string]*evaluator.Env
}

func (eco Ecosystem) Get(name string) (*evaluator.Env, error) {
	env, ok := eco.modules[name]
	if ok {
		return env, nil
	}
	env, ok = eco.shells[name]
	if ok {
		return env, nil
	}
	env, ok = eco.buffers[name]
	if ok {
		return env, nil
	}
	return nil, fmt.Errorf("Cannot switch to environment %s", name)
}

func NewEcosystem() Ecosystem {
	return Ecosystem{map[string]*evaluator.Env{}, map[string]*evaluator.Env{}, map[string]*evaluator.Env{}}
}

func (eco Ecosystem) Modules() map[string]*evaluator.Env {
	return eco.modules
}

func (eco Ecosystem) AddModule(name string, bindings map[string]value.Value) {
	eco.modules[name] = evaluator.NewEnv(bindings, nil, eco.modules)
}

func (eco Ecosystem) AddShell(name string, bindings map[string]value.Value) {
	eco.shells[name] = evaluator.NewEnv(bindings, nil, eco.modules)
}

func (eco Ecosystem) AddBuffer(name string, bindings map[string]value.Value) {
	eco.buffers[name] = evaluator.NewEnv(bindings, nil, eco.modules)

}
