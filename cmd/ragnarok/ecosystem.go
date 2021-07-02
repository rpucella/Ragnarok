package main

import (
	"fmt"
	"rpucella.net/ragnarok/internal/lisp"
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
	env *lisp.Env
}

type Ecosystem struct {
	modules map[string]*lisp.Env
	shells map[string]*lisp.Env
	buffers map[string]*lisp.Env
}

func (eco Ecosystem) get(name string) (*lisp.Env, error) {
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
	return Ecosystem{map[string]*lisp.Env{}, map[string]*lisp.Env{}, map[string]*lisp.Env{}}
}

func (eco Ecosystem) addModule(name string, env *lisp.Env) {
	eco.modules[name] = env
}

func (eco Ecosystem) addShell(name string, env *lisp.Env) {
	eco.shells[name] = env
}

func (eco Ecosystem) addBuffer(name string, env *lisp.Env) {
	eco.buffers[name] = env
}

func initialize() Ecosystem {
	eco := NewEcosystem()
	coreBindings := corePrimitives()
	coreBindings["true"] = lisp.NewVBoolean(true)
	coreBindings["false"] = lisp.NewVBoolean(false)
	coreEnv := lisp.NewEnv(coreBindings, nil, eco.modules)
	eco.addModule("core", coreEnv)
	testBindings := map[string]lisp.Value{
		"a": lisp.NewVInteger(99),
		"square": lisp.NewVPrimitive("square", func(args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if len(args) != 1 || !args[0].IsNumber() {
				return nil, fmt.Errorf("argument to square should be int")
			}
			return lisp.NewVInteger(args[0].IntValue() * args[0].IntValue()), nil
		}),
	}
	testEnv := lisp.NewEnv(testBindings, nil, eco.modules)
	eco.addModule("test", testEnv)
	shellBindings := shellPrimitives()
	shellEnv := lisp.NewEnv(shellBindings, nil, eco.modules)
	eco.addModule("shell", shellEnv)
	configBindings := map[string]lisp.Value{
		"lookup-path": lisp.NewVReference(lisp.NewVCons(lisp.NewVSymbol("shell"), lisp.NewVCons(lisp.NewVSymbol("core"), &lisp.VEmpty{}))),
		"editor":      lisp.NewVReference(lisp.NewVString("emacs")),
	}
	configEnv := lisp.NewEnv(configBindings, nil, eco.modules)
	eco.addModule("config", configEnv)
	return eco
}

