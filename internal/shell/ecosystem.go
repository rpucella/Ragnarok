package shell

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

func (eco Ecosystem) AddModule(name string, bindings map[string]value.Value) {
	eco.modules[name] = evaluator.NewEnv(bindings, nil, eco.modules)
}

func (eco Ecosystem) AddShell(name string, bindings map[string]value.Value) {
	eco.shells[name] = evaluator.NewEnv(bindings, nil, eco.modules)
}

func (eco Ecosystem) AddBuffer(name string, bindings map[string]value.Value) {
	eco.buffers[name] = evaluator.NewEnv(bindings, nil, eco.modules)

}

func CoreBindings() map[string]value.Value {
	bindings := corePrimitives()
	bindings["true"] = value.NewBoolean(true)
	bindings["false"] = value.NewBoolean(false)
	return bindings
}

func TestBindings() map[string]value.Value {
	bindings := map[string]value.Value{
		"a": value.NewInteger(99),
		"square": value.NewPrimitive("square", func(args []value.Value, ctxt interface{}) (value.Value, error) {
			if len(args) != 1 || !value.IsNumber(args[0]) {
				return nil, fmt.Errorf("argument to square should be int")
			}
			return value.NewInteger(args[0].GetInt() * args[0].GetInt()), nil
		}),
	}
	return bindings
}

func ShellBindings() map[string]value.Value {
	bindings := shellPrimitives()
	return bindings
}

func ConfigBindings() map[string]value.Value {
	bindings := map[string]value.Value{
		"lookup-path": value.NewReference(value.NewCons(value.NewSymbol("shell"), value.NewCons(value.NewSymbol("core"), value.NewEmpty()))),
		"editor":      value.NewReference(value.NewString("emacs")),
	}
	return bindings
}
