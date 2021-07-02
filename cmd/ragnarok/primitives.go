package main

import (
	"fmt"
	"rpucella.net/ragnarok/internal/lisp"
	"strings"
)

type PrimitiveDesc struct {
	name string
	min  int
	max  int // <0 for no max #
	prim func(string, []lisp.Value, interface{}) (lisp.Value, error)
}

func listLength(v lisp.Value) int {
	current := v
	result := 0
	for lisp.IsCons(current) {
		result += 1
		current = current.GetTail()
	}
	return result
}

func listAppend(v1 lisp.Value, v2 lisp.Value) lisp.Value {
	current := v1
	var result lisp.Value = nil
	var current_result *lisp.VCons = nil
	for lisp.IsCons(current) {
		cell := lisp.NewVCons(current.GetHead(), nil)
		current = current.GetTail()
		if current_result == nil {
			result = cell
		} else {
			current_result.SetTail(cell)
		}
		current_result = cell
	}
	if current_result == nil {
		return v2
	}
	current_result.SetTail(v2)
	return result
}

func allConses(vs []lisp.Value) bool {
	for _, v := range vs {
		if !lisp.IsCons(v) {
			return false
		}
	}
	return true
}

func corePrimitives() map[string]lisp.Value {
	bindings := map[string]lisp.Value{}
	for _, d := range CORE_PRIMITIVES {
		bindings[d.name] = lisp.NewVPrimitive(d.name, mkPrimitive(d))
	}
	return bindings
}

func shellPrimitives() map[string]lisp.Value {
	bindings := map[string]lisp.Value{}
	for _, d := range SHELL_PRIMITIVES {
		bindings[d.name] = lisp.NewVPrimitive(d.name, mkPrimitive(d))
	}
	return bindings
}

func mkPrimitive(d PrimitiveDesc) func([]lisp.Value, interface{}) (lisp.Value, error) {
	return func(args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
		if err := checkMinArgs(d.name, args, d.min); err != nil {
			return nil, err
		}
		if d.max >= 0 {
			if err := checkMaxArgs(d.name, args, d.max); err != nil {
				return nil, err
			}
		}
		return d.prim(d.name, args, ctxt)
	}
}

func checkArgType(name string, arg lisp.Value, pred func(lisp.Value) bool) error {
	if !pred(arg) {
		return fmt.Errorf("%s - wrong argument type %s", name, lisp.Classify(arg))
	}
	return nil
}

func checkMinArgs(name string, args []lisp.Value, n int) error {
	if len(args) < n {
		return fmt.Errorf("%s - too few arguments %d", name, len(args))
	}
	return nil
}

func checkMaxArgs(name string, args []lisp.Value, n int) error {
	if len(args) > n {
		return fmt.Errorf("%s - too many arguments %d", name, len(args))
	}
	return nil
}

func checkExactArgs(name string, args []lisp.Value, n int) error {
	if len(args) != n {
		return fmt.Errorf("%s - wrong number of arguments %d", name, len(args))
	}
	return nil
}

func isInt(v lisp.Value) bool {
	return lisp.IsNumber(v)
}

func IsString(v lisp.Value) bool {
	return lisp.IsString(v)
}

func IsSymbol(v lisp.Value) bool {
	return lisp.IsSymbol(v)
}

func IsFunction(v lisp.Value) bool {
	return lisp.IsFunction(v)
}

func isList(v lisp.Value) bool {
	return lisp.IsCons(v) || lisp.IsEmpty(v)
}

func IsReference(v lisp.Value) bool {
	return lisp.IsRef(v)
}

func mkNumPredicate(pred func(int, int) bool) func(string, []lisp.Value, interface{}) (lisp.Value, error) {
	return func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
		if err := checkExactArgs(name, args, 2); err != nil {
			return nil, err
		}
		if err := checkArgType(name, args[0], isInt); err != nil {
			return nil, err
		}
		if err := checkArgType(name, args[1], isInt); err != nil {
			return nil, err
		}
		return lisp.NewVBoolean(pred(args[0].GetInt(), args[1].GetInt())), nil
	}
}

var CORE_PRIMITIVES = []PrimitiveDesc{

	PrimitiveDesc{
		"type", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVSymbol(lisp.Classify(args[0])), nil
		},
	},

	PrimitiveDesc{
		"+", 0, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			v := 0
			for _, arg := range args {
				if err := checkArgType(name, arg, isInt); err != nil {
					return nil, err
				}
				v += arg.GetInt()
			}
			return lisp.NewVInteger(v), nil
		},
	},

	PrimitiveDesc{
		"*", 0, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			v := 1
			for _, arg := range args {
				if err := checkArgType(name, arg, isInt); err != nil {
					return nil, err
				}
				v *= arg.GetInt()
			}
			return lisp.NewVInteger(v), nil
		},
	},

	PrimitiveDesc{
		"-", 1, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			v := args[0].GetInt()
			if len(args) > 1 {
				for _, arg := range args[1:] {
					if err := checkArgType(name, arg, isInt); err != nil {
						return nil, err
					}
					v -= arg.GetInt()
				}
			} else {
				v = -v
			}
			return lisp.NewVInteger(v), nil
		},
	},

	PrimitiveDesc{"=", 2, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			var reference lisp.Value = args[0]
			for _, v := range args[1:] {
				if !lisp.IsEqual(reference, v) {
					return lisp.NewVBoolean(false), nil
				}
			}
			return lisp.NewVBoolean(true), nil
		},
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
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(!lisp.IsTrue(args[0])), nil
		},
	},

	PrimitiveDesc{
		"string-append", 0, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			v := ""
			for _, arg := range args {
				if err := checkArgType(name, arg, IsString); err != nil {
					return nil, err
				}
				v += arg.GetString()
			}
			return lisp.NewVString(v), nil
		},
	},

	PrimitiveDesc{"string-length", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsString); err != nil {
				return nil, err
			}
			return lisp.NewVInteger(len(args[0].GetString())), nil
		},
	},

	PrimitiveDesc{"string-lower", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsString); err != nil {
				return nil, err
			}
			return lisp.NewVString(strings.ToLower(args[0].GetString())), nil
		},
	},

	PrimitiveDesc{"string-upper", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsString); err != nil {
				return nil, err
			}
			return lisp.NewVString(strings.ToUpper(args[0].GetString())), nil
		},
	},

	PrimitiveDesc{"string-substring", 1, 3,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsString); err != nil {
				return nil, err
			}
			start := 0
			end := len(args[0].GetString())
			if len(args) > 2 {
				if err := checkArgType(name, args[2], isInt); err != nil {
					return nil, err
				}
				end = min(args[2].GetInt(), end)
			}
			if len(args) > 1 {
				if err := checkArgType(name, args[1], isInt); err != nil {
					return nil, err
				}
				start = max(args[1].GetInt(), start)
			}
			// or perhaps raise an exception
			if end < start {
				return lisp.NewVString(""), nil
			}
			return lisp.NewVString(args[0].GetString()[start:end]), nil
		},
	},

	PrimitiveDesc{"apply", 2, 2,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			arguments := make([]lisp.Value, listLength(args[1]))
			current := args[1]
			for i := range arguments {
				arguments[i] = current.GetHead()
				current = current.GetTail()
			}
			if !lisp.IsEmpty(current) {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return args[0].Apply(arguments, ctxt)
		},
	},

	PrimitiveDesc{"cons", 2, 2,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			return lisp.NewVCons(args[0], args[1]), nil
		},
	},

	PrimitiveDesc{
		"append", 0, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if len(args) == 0 {
				return &lisp.VEmpty{}, nil
			}
			if err := checkArgType(name, args[len(args)-1], isList); err != nil {
				return nil, err
			}
			result := args[len(args)-1]
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
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			var result lisp.Value = &lisp.VEmpty{}
			current := args[0]
			for lisp.IsCons(current) {
				result = lisp.NewVCons(current.GetHead(), result)
				current = current.GetTail()
			}
			if !lisp.IsEmpty(current) {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"head", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if lisp.IsEmpty(args[0]) {
				return nil, fmt.Errorf("%s - empty list argument", name)
			}
			return args[0].GetHead(), nil
		},
	},

	PrimitiveDesc{"tail", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if lisp.IsEmpty(args[0]) {
				return nil, fmt.Errorf("%s - empty list argument", name)
			}
			return args[0].GetTail(), nil
		},
	},

	PrimitiveDesc{"list", 0, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			var result lisp.Value = &lisp.VEmpty{}
			for i := len(args) - 1; i >= 0; i -= 1 {
				result = lisp.NewVCons(args[i], result)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"length", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			count := 0
			current := args[0]
			for lisp.IsCons(current) {
				count += 1
				current = current.GetTail()
			}
			if !lisp.IsEmpty(current) {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return lisp.NewVInteger(count), nil
		},
	},

	PrimitiveDesc{"nth", 2, 2,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isInt); err != nil {
				return nil, err
			}
			idx := args[1].GetInt()
			if idx >= 0 {
				current := args[0]
				for lisp.IsCons(current) {
					if idx == 0 {
						return current.GetHead(), nil
					} else {
						idx -= 1
						current = current.GetTail()
					}
				}
			}
			return nil, fmt.Errorf("%s - index %d out of bound", name, args[1].GetInt())
		},
	},

	PrimitiveDesc{"map", 2, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			for i := range args[1:] {
				if err := checkArgType(name, args[i+1], isList); err != nil {
					return nil, err
				}
			}
			var result lisp.Value = nil
			var current_result *lisp.VCons = nil
			currents := make([]lisp.Value, len(args)-1)
			firsts := make([]lisp.Value, len(args)-1)
			for i := range args[1:] {
				currents[i] = args[i+1]
			}
			for allConses(currents) {
				for i := range currents {
					firsts[i] = currents[i].GetHead()
				}
				v, err := args[0].Apply(firsts, ctxt)
				if err != nil {
					return nil, err
				}
				cell := lisp.NewVCons(v, nil)
				if current_result == nil {
					result = cell
				} else {
					current_result.SetTail(cell)
				}
				current_result = cell
				for i := range currents {
					currents[i] = currents[i].GetTail()
				}
			}
			if current_result == nil {
				return &lisp.VEmpty{}, nil
			}
			current_result.SetTail(&lisp.VEmpty{})
			return result, nil
		},
	},

	PrimitiveDesc{"for", 2, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			// TODO - allow different types in the same iteration!
			for i := range args[1:] {
				if err := checkArgType(name, args[i+1], isList); err != nil {
					return nil, err
				}
			}
			currents := make([]lisp.Value, len(args)-1)
			firsts := make([]lisp.Value, len(args)-1)
			for i := range args[1:] {
				currents[i] = args[i+1]
			}
			for allConses(currents) {
				for i := range currents {
					firsts[i] = currents[i].GetHead()
				}
				_, err := args[0].Apply(firsts, ctxt)
				if err != nil {
					return nil, err
				}
				for i := range currents {
					currents[i] = currents[i].GetTail()
				}
			}
			return &lisp.VNil{}, nil
		},
	},

	PrimitiveDesc{"filter", 2, 2,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			var result lisp.Value = nil
			var current_result *lisp.VCons = nil
			current := args[1]
			for lisp.IsCons(current) {
				v, err := args[0].Apply([]lisp.Value{current.GetHead()}, ctxt)
				if err != nil {
					return nil, err
				}
				if lisp.IsTrue(v) {
					cell := lisp.NewVCons(current.GetHead(), nil)
					if current_result == nil {
						result = cell
					} else {
						current_result.SetTail(cell)
					}
					current_result = cell
				}
				current = current.GetTail()
			}
			if !lisp.IsEmpty(current) {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			if current_result == nil {
				return &lisp.VEmpty{}, nil
			}
			current_result.SetTail(&lisp.VEmpty{})
			return result, nil
		},
	},

	PrimitiveDesc{"foldr", 3, 3,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			var temp lisp.Value = &lisp.VEmpty{}
			// first reverse the list
			current := args[1]
			for lisp.IsCons(current) {
				temp = lisp.NewVCons(current.GetHead(), temp)
				current = current.GetTail()
			}
			if !lisp.IsEmpty(current) {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			// then fold it
			result := args[2]
			current = temp
			for lisp.IsCons(current) {
				v, err := args[0].Apply([]lisp.Value{current.GetHead(), result}, ctxt)
				if err != nil {
					return nil, err
				}
				result = v
				current = current.GetTail()
			}
			if !lisp.IsEmpty(current) {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"foldl", 3, 3,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			result := args[2]
			current := args[1]
			for lisp.IsCons(current) {
				v, err := args[0].Apply([]lisp.Value{result, current.GetHead()}, ctxt)
				if err != nil {
					return nil, err
				}
				result = v
				current = current.GetTail()
			}
			if !lisp.IsEmpty(current) {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"ref", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVReference(args[0]), nil
		},
	},

	// set should probably be a special form
	// (set (x) 10)
	// (set (arr 1) 10)
	// (set (dict 'a) 10)
	// like setf in CLISP

	// PrimitiveDesc{"set", 2, 2,
	// 	func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
	// 		if err := checkArgType(name, args[0], IsReference); err != nil {
	// 			return nil, err
	// 		}
	// 		args[0].setValue(args[1])
	// 		return &lisp.VNil{}, nil
	// 	},
	// },

	PrimitiveDesc{"empty?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsEmpty(args[0])), nil
		},
	},

	PrimitiveDesc{"cons?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsCons(args[0])), nil
		},
	},

	PrimitiveDesc{"list?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsCons(args[0]) || lisp.IsEmpty(args[0])), nil
		},
	},

	PrimitiveDesc{"number?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsNumber(args[0])), nil
		},
	},

	PrimitiveDesc{"ref?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsRef(args[0])), nil
		},
	},

	PrimitiveDesc{"boolean?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsBool(args[0])), nil
		},
	},

	PrimitiveDesc{"string?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsString(args[0])), nil
		},
	},

	PrimitiveDesc{"symbol?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsSymbol(args[0])), nil
		},
	},

	PrimitiveDesc{"function?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsFunction(args[0])), nil
		},
	},

	PrimitiveDesc{"nil?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsNil(args[0])), nil
		},
	},

	PrimitiveDesc{"array", 0, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			content := make([]lisp.Value, len(args))
			for i, v := range args {
				content[i] = v
			}
			return lisp.NewVArray(content), nil
		},
	},

	PrimitiveDesc{"array?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsArray(args[0])), nil
		},
	},

	PrimitiveDesc{"dict", 0, -1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			content := make(map[string]lisp.Value, len(args))
			for _, v := range args {
				if !lisp.IsCons(v) || !lisp.IsCons(v.GetTail()) || !lisp.IsEmpty(v.GetTail().GetTail()) {
					return nil, fmt.Errorf("dict item not a pair - %s", v.Display())
				}
				if !lisp.IsSymbol(v.GetHead()) {
					return nil, fmt.Errorf("dict key is not a symbol - %s", v.GetHead().Display())
				}
				content[v.GetHead().GetString()] = v.GetTail().GetHead()
			}
			return lisp.NewVDict(content), nil
		},
	},

	PrimitiveDesc{"dict?", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			return lisp.NewVBoolean(lisp.IsDict(args[0])), nil
		},
	},
}

var SHELL_PRIMITIVES = []PrimitiveDesc{

	PrimitiveDesc{
		"quit", 0, 0,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			bail()
			return &lisp.VNil{}, nil
		},
	},

	PrimitiveDesc{
		"go", 1, 1,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			if err := checkArgType(name, args[0], IsSymbol); err != nil {
				return nil, err
			}
			context, ok := ctxt.(*Context)
			if !ok {
				return nil, fmt.Errorf("Problem understanding context")
			}
			context.nextCurrentModule = args[0].GetString()
			return &lisp.VNil{}, nil
		},
	},

	PrimitiveDesc{
		"modules", 0, 0,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			context, ok := ctxt.(*Context)
			if !ok {
				return nil, fmt.Errorf("Problem understanding context")
			}
			var result lisp.Value = &lisp.VEmpty{}
			for m := range context.ecosystem.modules {
				result = lisp.NewVCons(lisp.NewVSymbol(m), result)
			}
			return result, nil
		},
	},

	PrimitiveDesc{
		"help", 0, 0,
		func(name string, args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
			context, ok := ctxt.(*Context)
			if !ok {
				return nil, fmt.Errorf("Problem understanding context")
			}
			context.report("Some help about the system")
			context.report("")
			context.report("      (quit)   bail out")
			context.report("   (modules)   see available modules")
			context.report("  (go 'buff)   navigate to a particular buffer")
			context.report("")
			return &lisp.VNil{}, nil
		},
	},
	

}

// left:
//
// dictionaries #((a 1) (b 2))  (dict '(a 10) '(b 20) '(c 30))  vs (apply dict '((a 10) (b 20) (c 30)))?
// arrays #[a b c]
