package main

import (
	"fmt"
	"rpucella.net/ragnarok/internal/evaluator"
	"rpucella.net/ragnarok/internal/value"
)

// An ecosystem is a global set of environments associated with modules, shells, or buffers

type envKind int

const (
	ENV_MODULE envKind = iota
	ENV_SHELL
	ENV_BUFFER
)

type NamedEnv struct {
	kind envKind
	name string
	env  *evaluator.Env
}

type Ecosystem struct {
	modules map[string]*evaluator.Env
	shells  map[string]*evaluator.Env
	buffers map[string]*evaluator.Env
}

func (eco Ecosystem) get(name string) (*evaluator.Env, error) {
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

func (eco Ecosystem) addModule(name string, env *evaluator.Env) {
	eco.modules[name] = env
}

func (eco Ecosystem) addShell(name string, env *evaluator.Env) {
	eco.shells[name] = env
}

func (eco Ecosystem) addBuffer(name string, env *evaluator.Env) {
	eco.buffers[name] = env
}

func initialize() Ecosystem {
	eco := NewEcosystem()
	coreBindings := corePrimitives()
	coreBindings["true"] = value.NewVBoolean(true)
	coreBindings["false"] = value.NewVBoolean(false)
	coreEnv := evaluator.NewEnv(coreBindings, nil, eco.modules)
	eco.addModule("core", coreEnv)
	testBindings := map[string]value.Value{
		"a": value.NewVInteger(99),
		"square": value.NewVPrimitive("square", func(args []value.Value, ctxt interface{}) (value.Value, error) {
			if len(args) != 1 || !value.IsNumber(args[0]) {
				return nil, fmt.Errorf("argument to square should be int")
			}
			return value.NewVInteger(args[0].GetInt() * args[0].GetInt()), nil
		}),
	}
	testEnv := evaluator.NewEnv(testBindings, nil, eco.modules)
	eco.addModule("test", testEnv)
	shellBindings := shellPrimitives()
	shellEnv := evaluator.NewEnv(shellBindings, nil, eco.modules)
	eco.addModule("shell", shellEnv)
	configBindings := map[string]value.Value{
		"lookup-path": value.NewVReference(value.NewVCons(value.NewVSymbol("shell"), value.NewVCons(value.NewVSymbol("core"), &value.VEmpty{}))),
		"editor":      value.NewVReference(value.NewVString("emacs")),
	}
	configEnv := evaluator.NewEnv(configBindings, nil, eco.modules)
	eco.addModule("config", configEnv)
	return eco
}
