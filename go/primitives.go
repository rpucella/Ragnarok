package main

import "fmt"

type PrimitiveDesc struct {
	name string
	prim func([]Value)(Value, error)
}

func primitivesBindings() map[string]Value {
	bindings := map[string]Value{}
	for _, d := range PRIMITIVES {
		bindings[d.name] = &VPrimitive{d.name, d.prim}
	}
	return bindings
}

func checkArgType(name string, arg Value, pred func(Value)bool) error {
	if !pred(arg) {
		return fmt.Errorf("wrong argument type %s to primitive %s", arg.typ(), name)
	}
	return nil
}

func checkMinArgs(name string, args []Value, n int) error {
	if len(args) < n {
		return fmt.Errorf("too few arguments %d to primitive %s", len(args), name)
	}
	return nil
}

func checkMaxArgs(name string, args []Value, n int) error {
	if len(args) > n {
		return fmt.Errorf("too many arguments %d to primitive %s", len(args), name)
	}
	return nil
}

func checkExactArgs(name string, args []Value, n int) error {
	if len(args) != n {
		return fmt.Errorf("wrong number of arguments %d to primitive %s", len(args), name)
	}
	return nil
}

func checkInt(v Value) bool {
	return v.typ() == "int"
}

var PRIMITIVES = []PrimitiveDesc{
	
	PrimitiveDesc{
		"type",
		func(args []Value) (Value, error) {
			if err := checkExactArgs("type", args, 1); err != nil { 
				return nil, err
			}
			return &VSymbol{args[0].typ()}, nil
		},
	},

	PrimitiveDesc{
		"+", 
		func(args []Value) (Value, error) {
			v := 0
			for _, arg := range args {
				if err := checkArgType("+", arg, checkInt); err != nil { 
					return nil, err
				}
				v += arg.intValue()
			}
			return &VInteger{v}, nil
		},
	},

	PrimitiveDesc{
		"*", 
		func(args []Value) (Value, error) {
			v := 1
			for _, arg := range args {
				if err := checkArgType("*", arg, checkInt); err != nil { 
					return nil, err
				}
				v *= arg.intValue()
			}
			return &VInteger{v}, nil
		},
	},
	
	PrimitiveDesc{
		"-", 
		func(args []Value) (Value, error) {
			if err := checkMinArgs("-", args, 1); err != nil {
				return nil, err
			}
			v := args[0].intValue()
			if len(args) > 1 { 
				for _, arg := range args[1:] {
					if err := checkArgType("-", arg, checkInt); err != nil { 
						return nil, err
					}
					v -= arg.intValue()
				}
			} else {
				v = -v
			}
			return &VInteger{v}, nil
		},
	},
	
}

