 package main

import "fmt"
import "strings"

type PrimitiveDesc struct {
	name string
	min int
	max int      // <0 for no max #
	prim func(string, []Value)(Value, error)
}

func listLength (v Value) int {
	current := v
	result := 0
	for current.isCons() {
		result += 1
		current = current.tailValue()
	}
	return result
}

func listAppend (v1 Value, v2 Value) Value {
	current := v1
	var result Value = nil
	var current_result *VCons = nil
	for current.isCons() {
		cell := &VCons{head: current.headValue(), tail: nil}
		current = current.tailValue()
		if current_result == nil {
			result = cell
		} else {
			current_result.tail = cell
		}
		current_result = cell
	}
	current_result.tail = v2
	return result
}

func corePrimitives() map[string]Value {
	bindings := map[string]Value{}
	for _, d := range PRIMITIVES {
		bindings[d.name] = &VPrimitive{d.name, mkPrimitive(d)}
	}
	return bindings
}

func mkPrimitive(d PrimitiveDesc) func([]Value)(Value, error) {
	return func(args []Value) (Value, error) {
		if err := checkMinArgs(d.name, args, d.min); err != nil { 
			return nil, err
		}
		if d.max >= 0 {
			if err := checkMaxArgs(d.name, args, d.max); err != nil { 
				return nil, err
			}
		}
		return d.prim(d.name, args)
	}
}
	
func checkArgType(name string, arg Value, pred func(Value)bool) error {
	if !pred(arg) {
		return fmt.Errorf("%s - wrong argument type %s", name, arg.typ())
	}
	return nil
}

func checkMinArgs(name string, args []Value, n int) error {
	if len(args) < n {
		return fmt.Errorf("%s - too few arguments %d", name, len(args))
	}
	return nil
}

func checkMaxArgs(name string, args []Value, n int) error {
	if len(args) > n {
		return fmt.Errorf("%s - too many arguments %d", name, len(args))
	}
	return nil
}

func checkExactArgs(name string, args []Value, n int) error {
	if len(args) != n {
		return fmt.Errorf("%s - wrong number of arguments %d", name, len(args))
	}
	return nil
}

func isInt(v Value) bool {
	return v.typ() == "int"
}

func isString(v Value) bool {
	return v.typ() == "string"
}

func isFunction(v Value) bool {
	return v.typ() == "fun"
}

func isList(v Value) bool {
	return v.typ() == "list"
}

func mkNumPredicate(pred func(int, int)bool) func(string, []Value)(Value, error) {
	return func(name string, args []Value) (Value, error) {
		if err := checkExactArgs(name, args, 2); err != nil {
			return nil, err
		}
		if err := checkArgType(name, args[0], isInt); err != nil {
			return nil, err
		}
		if err := checkArgType(name, args[1], isInt); err != nil {
			return nil, err
		}
		return &VBoolean{pred(args[0].intValue(), args[1].intValue())}, nil
	}
}

var PRIMITIVES = []PrimitiveDesc{
	
	PrimitiveDesc{
		"type", 1, 1,
		func(name string, args []Value) (Value, error) {
			return &VSymbol{args[0].typ()}, nil
		},
	},
	
	PrimitiveDesc{
		"+", 0, -1,
		func(name string, args []Value) (Value, error) {
			v := 0
			for _, arg := range args {
				if err := checkArgType(name, arg, isInt); err != nil { 
					return nil, err
				}
				v += arg.intValue()
			}
			return &VInteger{v}, nil
		},
	},

	PrimitiveDesc{
		"*", 0, -1,
		func(name string, args []Value) (Value, error) {
			v := 1
			for _, arg := range args {
				if err := checkArgType(name, arg, isInt); err != nil { 
					return nil, err
				}
				v *= arg.intValue()
			}
			return &VInteger{v}, nil
		},
	},
	
	PrimitiveDesc{
		"-", 1, -1,
		func(name string, args []Value) (Value, error) {
			v := args[0].intValue()
			if len(args) > 1 { 
				for _, arg := range args[1:] {
					if err := checkArgType(name, arg, isInt); err != nil { 
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

	PrimitiveDesc{"=", 2, 2,
		mkNumPredicate(func(n1 int, n2 int) bool { return n1 == n2 }),
	},

	PrimitiveDesc{"<", 2, 2,
		mkNumPredicate(func(n1 int, n2 int) bool { return n1 < n2 }),
	},

	PrimitiveDesc{"<=", 2, 2,
		mkNumPredicate(func(n1 int, n2 int) bool { return n1 <= n2 }),
	},

	PrimitiveDesc{">", 2, 2,
		mkNumPredicate(func(n1 int, n2 int) bool { return n1 > n2 }),
	},

	PrimitiveDesc{">=", 2, 2,
		mkNumPredicate(func(n1 int, n2 int) bool { return n1 >= n2 }),
	},

	PrimitiveDesc{"not", 1, 1,
		func(name string, args []Value) (Value, error) {
			return &VBoolean{!args[0].isTrue()}, nil
		},
	},

	PrimitiveDesc{
		"string-append", 0, -1,
		func(name string, args []Value) (Value, error) {
			v := ""
			for _, arg := range args {
				if err := checkArgType(name, arg, isString); err != nil { 
					return nil, err
				}
				v += arg.strValue()
			}
			return &VString{v}, nil
		},
	},

	PrimitiveDesc{"string-length", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isString); err != nil {
				return nil, err
			}
			return &VInteger{len(args[0].strValue())}, nil
		},
	},

	PrimitiveDesc{"string-lower", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isString); err != nil {
				return nil, err
			}
			return &VString{strings.ToLower(args[0].strValue())}, nil
		},
	},

	PrimitiveDesc{"string-upper", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isString); err != nil {
				return nil, err
			}
			return &VString{strings.ToUpper(args[0].strValue())}, nil
		},
	},

	PrimitiveDesc{"string-substring", 1, 3,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isString); err != nil {
				return nil, err
			}
			start := 0
			end := len(args[0].strValue())
			if len(args) > 2 {
				if err := checkArgType(name, args[2], isInt); err != nil {
					return nil, err
				}
				end = min(args[2].intValue(), end)
			}
			if len(args) > 1 {
				if err := checkArgType(name, args[1], isInt); err != nil {
					return nil, err
				}
				start = max(args[1].intValue(), start)
			}
			// or perhaps raise an exception
			if (end < start) {
				return &VString{""}, nil
			}
			return &VString{args[0].strValue()[start:end]}, nil
		},
	},

	PrimitiveDesc{"apply", 2, 2,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			arguments := make([]Value, listLength(args[1]))
			current := args[1]
			for i := range arguments {
				arguments[i] = current.headValue()
				current = current.tailValue()
			}
			if !current.isEmpty() {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return args[0].apply(arguments)
		},
	},

	PrimitiveDesc{"cons", 2, 2,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			return &VCons{head: args[0], tail: args[1]}, nil
		},
	},

	PrimitiveDesc{
		"append", 0, -1,
		func(name string, args []Value) (Value, error) {
			if len(args) == 0 {
				return &VEmpty{}, nil
			}
			if err := checkArgType(name, args[len(args) - 1], isList); err != nil {
				return nil, err
			}
			result := args[len(args) - 1]
			for i := len(args) - 2; i >= 0; i -= 1 { 
				if err := checkArgType(name, args[i], isList); err != nil { 
					return nil, err
				}
				result = listAppend(args[i], result)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"reverse", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			var result Value = &VEmpty{}
			current := args[0]
			for current.isCons() {
				result = &VCons{head: current.headValue(), tail: result}
				current = current.tailValue()
			}
			if !current.isEmpty() {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"head", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if args[0].isEmpty() {
				return nil, fmt.Errorf("%s - empty list argument", name)
			}
			return args[0].headValue(), nil
		},
	},

	PrimitiveDesc{"tail", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if args[0].isEmpty() {
				return nil, fmt.Errorf("%s - empty list argument", name)
			}
			return args[0].tailValue(), nil
		},
	},

	PrimitiveDesc{"list", 0, -1,
		func(name string, args []Value) (Value, error) {
			var result Value = &VEmpty{}
			for i := len(args) - 1; i >= 0; i -= 1 {
				result = &VCons{head: args[i], tail: result}
			}
			return result, nil
		},
	},

	PrimitiveDesc{"length", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			count := 0
			current := args[0]
			for current.isCons() {
				count += 1
				current = current.tailValue()
			}
			if !current.isEmpty() { 
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return &VInteger{count}, nil
		},
	},

	PrimitiveDesc{"nth", 2, 2,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isInt); err != nil {
				return nil, err
			}
			idx := args[1].intValue()
			if idx >= 0 {
				current := args[0]
				for current.isCons() {
					if idx == 0 {
						return current.headValue(), nil
					} else {
						idx -= 1
						current = current.tailValue()
					}
				}
			}
			return nil, fmt.Errorf("%s - index %d out of bound", name, args[1].intValue())
		},
	},

}
