package main

import (
	"fmt"
	"rpucella.net/ragnarok/internal/primitives"
	"rpucella.net/ragnarok/internal/value"
)

func CoreBindings() map[string]value.Value {
	bindings := primitives.CorePrimitives()
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
	bindings := primitives.ShellPrimitives()
	return bindings
}

func ConfigBindings() map[string]value.Value {
	bindings := map[string]value.Value{
		"lookup-path": value.NewReference(value.NewCons(value.NewSymbol("shell"), value.NewCons(value.NewSymbol("core"), value.NewEmpty()))),
		"editor":      value.NewReference(value.NewString("emacs")),
	}
	return bindings
}
