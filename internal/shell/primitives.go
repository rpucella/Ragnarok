package shell

import (
	"fmt"
	"io/ioutil"
	"rpucella.net/ragnarok/internal/util"
	"rpucella.net/ragnarok/internal/value"
	"sort"
	"strings"
	"time"
)

type PrimitiveDesc struct {
	name string
	min  int
	max  int // <0 for no max #
	prim func(string, []value.Value, interface{}) (value.Value, error)
}

func listLength(v value.Value) int {
	current := v
	result := 0
	for value.IsCons(current) {
		result += 1
		current = current.GetTail()
	}
	return result
}

func listAppend(v1 value.Value, v2 value.Value) value.Value {
	current := v1
	var result value.Value = nil
	var current_result *value.Cons = nil
	for value.IsCons(current) {
		cell := value.NewCons(current.GetHead(), nil)
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

func allConses(vs []value.Value) bool {
	for _, v := range vs {
		if !value.IsCons(v) {
			return false
		}
	}
	return true
}

func corePrimitives() map[string]value.Value {
	bindings := map[string]value.Value{}
	for _, d := range CORE_PRIMITIVES {
		bindings[d.name] = value.NewPrimitive(d.name, mkPrimitive(d))
	}
	return bindings
}

func shellPrimitives() map[string]value.Value {
	bindings := map[string]value.Value{}
	for _, d := range SHELL_PRIMITIVES {
		bindings[d.name] = value.NewPrimitive(d.name, mkPrimitive(d))
	}
	return bindings
}

func mkPrimitive(d PrimitiveDesc) func([]value.Value, interface{}) (value.Value, error) {
	return func(args []value.Value, ctxt interface{}) (value.Value, error) {
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

func checkArgType(name string, arg value.Value, pred func(value.Value) bool) error {
	if !pred(arg) {
		return fmt.Errorf("%s: wrong argument type %s", name, value.Classify(arg))
	}
	return nil
}

func checkMinArgs(name string, args []value.Value, n int) error {
	if len(args) < n {
		return fmt.Errorf("%s: too few arguments %d", name, len(args))
	}
	return nil
}

func checkMaxArgs(name string, args []value.Value, n int) error {
	if len(args) > n {
		return fmt.Errorf("%s: too many arguments %d", name, len(args))
	}
	return nil
}

func checkExactArgs(name string, args []value.Value, n int) error {
	if len(args) != n {
		return fmt.Errorf("%s: wrong number of arguments %d", name, len(args))
	}
	return nil
}

func isInt(v value.Value) bool {
	return value.IsNumber(v)
}

func IsString(v value.Value) bool {
	return value.IsString(v)
}

func IsSymbol(v value.Value) bool {
	return value.IsSymbol(v)
}

func IsFunction(v value.Value) bool {
	return value.IsFunction(v)
}

func isList(v value.Value) bool {
	return value.IsCons(v) || value.IsEmpty(v)
}

func IsReference(v value.Value) bool {
	return value.IsRef(v)
}

func mkNumPredicate(pred func(int, int) bool) func(string, []value.Value, interface{}) (value.Value, error) {
	return func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		if err := checkExactArgs(name, args, 2); err != nil {
			return nil, err
		}
		if err := checkArgType(name, args[0], isInt); err != nil {
			return nil, err
		}
		if err := checkArgType(name, args[1], isInt); err != nil {
			return nil, err
		}
		return value.NewBoolean(pred(args[0].GetInt(), args[1].GetInt())), nil
	}
}

var CORE_PRIMITIVES = []PrimitiveDesc{

	PrimitiveDesc{
		"type", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewSymbol(value.Classify(args[0])), nil
		},
	},

	PrimitiveDesc{
		"+", 0, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			v := 0
			for _, arg := range args {
				if err := checkArgType(name, arg, isInt); err != nil {
					return nil, err
				}
				v += arg.GetInt()
			}
			return value.NewInteger(v), nil
		},
	},

	PrimitiveDesc{
		"*", 0, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			v := 1
			for _, arg := range args {
				if err := checkArgType(name, arg, isInt); err != nil {
					return nil, err
				}
				v *= arg.GetInt()
			}
			return value.NewInteger(v), nil
		},
	},

	PrimitiveDesc{
		"-", 1, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
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
			return value.NewInteger(v), nil
		},
	},

	PrimitiveDesc{"=", 2, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			var reference value.Value = args[0]
			for _, v := range args[1:] {
				if !value.IsEqual(reference, v) {
					return value.NewBoolean(false), nil
				}
			}
			return value.NewBoolean(true), nil
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
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(!value.IsTrue(args[0])), nil
		},
	},

	PrimitiveDesc{
		"string-append", 0, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			v := ""
			for _, arg := range args {
				if err := checkArgType(name, arg, IsString); err != nil {
					return nil, err
				}
				v += arg.GetString()
			}
			return value.NewString(v), nil
		},
	},

	PrimitiveDesc{"string-length", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsString); err != nil {
				return nil, err
			}
			return value.NewInteger(len(args[0].GetString())), nil
		},
	},

	PrimitiveDesc{"string-lower", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsString); err != nil {
				return nil, err
			}
			return value.NewString(strings.ToLower(args[0].GetString())), nil
		},
	},

	PrimitiveDesc{"string-upper", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsString); err != nil {
				return nil, err
			}
			return value.NewString(strings.ToUpper(args[0].GetString())), nil
		},
	},

	PrimitiveDesc{"string-substring", 1, 3,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsString); err != nil {
				return nil, err
			}
			start := 0
			end := len(args[0].GetString())
			if len(args) > 2 {
				if err := checkArgType(name, args[2], isInt); err != nil {
					return nil, err
				}
				end = util.MinInt(args[2].GetInt(), end)
			}
			if len(args) > 1 {
				if err := checkArgType(name, args[1], isInt); err != nil {
					return nil, err
				}
				start = util.MaxInt(args[1].GetInt(), start)
			}
			// or perhaps raise an exception
			if end < start {
				return value.NewString(""), nil
			}
			return value.NewString(args[0].GetString()[start:end]), nil
		},
	},

	PrimitiveDesc{"apply", 2, 2,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			arguments := make([]value.Value, listLength(args[1]))
			current := args[1]
			for i := range arguments {
				arguments[i] = current.GetHead()
				current = current.GetTail()
			}
			if !value.IsEmpty(current) {
				return nil, fmt.Errorf("%s: malformed list", name)
			}
			return args[0].Apply(arguments, ctxt)
		},
	},

	PrimitiveDesc{"cons", 2, 2,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			return value.NewCons(args[0], args[1]), nil
		},
	},

	PrimitiveDesc{
		"append", 0, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if len(args) == 0 {
				return value.NewEmpty(), nil
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
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			var result value.Value = value.NewEmpty()
			current := args[0]
			for value.IsCons(current) {
				result = value.NewCons(current.GetHead(), result)
				current = current.GetTail()
			}
			if !value.IsEmpty(current) {
				return nil, fmt.Errorf("%s: malformed list", name)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"head", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if value.IsEmpty(args[0]) {
				return nil, fmt.Errorf("%s: empty list argument", name)
			}
			return args[0].GetHead(), nil
		},
	},

	PrimitiveDesc{"tail", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if value.IsEmpty(args[0]) {
				return nil, fmt.Errorf("%s: empty list argument", name)
			}
			return args[0].GetTail(), nil
		},
	},

	PrimitiveDesc{"list", 0, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			var result value.Value = value.NewEmpty()
			for i := len(args) - 1; i >= 0; i -= 1 {
				result = value.NewCons(args[i], result)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"length", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			count := 0
			current := args[0]
			for value.IsCons(current) {
				count += 1
				current = current.GetTail()
			}
			if !value.IsEmpty(current) {
				return nil, fmt.Errorf("%s: malformed list", name)
			}
			return value.NewInteger(count), nil
		},
	},

	PrimitiveDesc{"nth", 2, 2,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isInt); err != nil {
				return nil, err
			}
			idx := args[1].GetInt()
			if idx >= 0 {
				current := args[0]
				for value.IsCons(current) {
					if idx == 0 {
						return current.GetHead(), nil
					} else {
						idx -= 1
						current = current.GetTail()
					}
				}
			}
			return nil, fmt.Errorf("%s: index %d out of bound", name, args[1].GetInt())
		},
	},

	PrimitiveDesc{"map", 2, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			for i := range args[1:] {
				if err := checkArgType(name, args[i+1], isList); err != nil {
					return nil, err
				}
			}
			var result value.Value = nil
			var current_result *value.Cons = nil
			currents := make([]value.Value, len(args)-1)
			firsts := make([]value.Value, len(args)-1)
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
				cell := value.NewCons(v, nil)
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
				return value.NewEmpty(), nil
			}
			current_result.SetTail(value.NewEmpty())
			return result, nil
		},
	},

	PrimitiveDesc{"for", 2, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			// TODO - allow different types in the same iteration!
			for i := range args[1:] {
				if err := checkArgType(name, args[i+1], isList); err != nil {
					return nil, err
				}
			}
			currents := make([]value.Value, len(args)-1)
			firsts := make([]value.Value, len(args)-1)
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
			return value.NewNil(), nil
		},
	},

	PrimitiveDesc{"filter", 2, 2,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			var result value.Value = nil
			var current_result *value.Cons = nil
			current := args[1]
			for value.IsCons(current) {
				v, err := args[0].Apply([]value.Value{current.GetHead()}, ctxt)
				if err != nil {
					return nil, err
				}
				if value.IsTrue(v) {
					cell := value.NewCons(current.GetHead(), nil)
					if current_result == nil {
						result = cell
					} else {
						current_result.SetTail(cell)
					}
					current_result = cell
				}
				current = current.GetTail()
			}
			if !value.IsEmpty(current) {
				return nil, fmt.Errorf("%s: malformed list", name)
			}
			if current_result == nil {
				return value.NewEmpty(), nil
			}
			current_result.SetTail(value.NewEmpty())
			return result, nil
		},
	},

	PrimitiveDesc{"foldr", 3, 3,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			var temp value.Value = value.NewEmpty()
			// first reverse the list
			current := args[1]
			for value.IsCons(current) {
				temp = value.NewCons(current.GetHead(), temp)
				current = current.GetTail()
			}
			if !value.IsEmpty(current) {
				return nil, fmt.Errorf("%s: malformed list", name)
			}
			// then fold it
			result := args[2]
			current = temp
			for value.IsCons(current) {
				v, err := args[0].Apply([]value.Value{current.GetHead(), result}, ctxt)
				if err != nil {
					return nil, err
				}
				result = v
				current = current.GetTail()
			}
			if !value.IsEmpty(current) {
				return nil, fmt.Errorf("%s: malformed list", name)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"foldl", 3, 3,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			result := args[2]
			current := args[1]
			for value.IsCons(current) {
				v, err := args[0].Apply([]value.Value{result, current.GetHead()}, ctxt)
				if err != nil {
					return nil, err
				}
				result = v
				current = current.GetTail()
			}
			if !value.IsEmpty(current) {
				return nil, fmt.Errorf("%s: malformed list", name)
			}
			return result, nil
		},
	},

	PrimitiveDesc{"ref", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewReference(args[0]), nil
		},
	},

	// set should probably be a special form
	// (set (x) 10)
	// (set (arr 1) 10)
	// (set (dict 'a) 10)
	// like setf in CLISP

	PrimitiveDesc{"get", 1, 2,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if value.IsRef(args[0]) {
				if err := checkExactArgs(name, args, 1); err != nil {
					return nil, err
				}
				return args[0].GetValue(), nil
			}
			if value.IsArray(args[0]) {
				if err := checkExactArgs(name, args, 2); err != nil {
					return nil, err
				}
				if err := checkArgType(name, args[1], isInt); err != nil {
					return nil, err
				}
				arr := args[0].GetArray()
				idx := args[1].GetInt()
				if idx < 0 || idx >= len(arr) {
					return nil, fmt.Errorf("%s: index %d out of bounds", name, idx)
				}
				return arr[idx], nil
			}
			if value.IsDict(args[0]) {
				if err := checkExactArgs(name, args, 2); err != nil {
					return nil, err
				}
				if err := checkArgType(name, args[1], IsSymbol); err != nil {
					return nil, err
				}
				m := args[0].GetMap()
				key := args[1].GetString()
				result, ok := m[key]
				if !ok {
					return nil, fmt.Errorf("%s: key %s not in dictionary", name, key)
				}
				return result, nil
			}
			return nil, fmt.Errorf("%s: wrong argument type %s", name, value.Classify(args[0]))
		},
	},

	PrimitiveDesc{"set!", 2, 3,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			if value.IsRef(args[0]) {
				if err := checkExactArgs(name, args, 2); err != nil {
					return nil, err
				}
				args[0].SetValue(args[1])
				return value.NewNil(), nil
			}
			if value.IsArray(args[0]) {
				if err := checkExactArgs(name, args, 3); err != nil {
					return nil, err
				}
				if err := checkArgType(name, args[1], isInt); err != nil {
					return nil, err
				}
				arr := args[0].GetArray()
				idx := args[1].GetInt()
				if idx < 0 || idx >= len(arr) {
					return nil, fmt.Errorf("%s: index %d out of bounds", name, idx)
				}
				arr[idx] = args[2]
				return value.NewNil(), nil
			}
			if value.IsDict(args[0]) {
				if err := checkExactArgs(name, args, 3); err != nil {
					return nil, err
				}
				if err := checkArgType(name, args[1], IsSymbol); err != nil {
					return nil, err
				}
				m := args[0].GetMap()
				key := args[1].GetString()
				m[key] = args[2]
				return value.NewNil(), nil
			}
			return nil, fmt.Errorf("%s: wrong argument type %s", name, value.Classify(args[0]))
		},
	},

	PrimitiveDesc{"empty?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsEmpty(args[0])), nil
		},
	},

	PrimitiveDesc{"cons?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsCons(args[0])), nil
		},
	},

	PrimitiveDesc{"list?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsCons(args[0]) || value.IsEmpty(args[0])), nil
		},
	},

	PrimitiveDesc{"number?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsNumber(args[0])), nil
		},
	},

	PrimitiveDesc{"ref?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsRef(args[0])), nil
		},
	},

	PrimitiveDesc{"boolean?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsBool(args[0])), nil
		},
	},

	PrimitiveDesc{"string?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsString(args[0])), nil
		},
	},

	PrimitiveDesc{"symbol?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsSymbol(args[0])), nil
		},
	},

	PrimitiveDesc{"function?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsFunction(args[0])), nil
		},
	},

	PrimitiveDesc{"nil?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsNil(args[0])), nil
		},
	},

	PrimitiveDesc{"array", 0, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			content := make([]value.Value, len(args))
			for i, v := range args {
				content[i] = v
			}
			return value.NewArray(content), nil
		},
	},

	PrimitiveDesc{"array?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsArray(args[0])), nil
		},
	},

	PrimitiveDesc{"dict", 0, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			content := make(map[string]value.Value, len(args))
			for _, v := range args {
				if !value.IsCons(v) || !value.IsCons(v.GetTail()) || !value.IsEmpty(v.GetTail().GetTail()) {
					return nil, fmt.Errorf("dict item not a pair - %s", v.Display())
				}
				if !value.IsSymbol(v.GetHead()) {
					return nil, fmt.Errorf("dict key is not a symbol - %s", v.GetHead().Display())
				}
				content[v.GetHead().GetString()] = v.GetTail().GetHead()
			}
			return value.NewDict(content), nil
		},
	},

	PrimitiveDesc{"dict?", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			return value.NewBoolean(value.IsDict(args[0])), nil
		},
	},
}

var SHELL_PRIMITIVES = []PrimitiveDesc{

	PrimitiveDesc{
		"quit", 0, 0,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			context, ok := ctxt.(*Context)
			if !ok {
				return nil, fmt.Errorf("Problem understanding context")
			}
			context.bail()
			return value.NewNil(), nil
		},
	},

	PrimitiveDesc{
		"env", 0, 0,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			context, ok := ctxt.(*Context)
			if !ok {
				return nil, fmt.Errorf("Problem understanding context")
			}
			bindings := context.currentEnv.Bindings()
			maxWidth := 0
			keys := make([]string, len(bindings))
			i := 0
			for name := range bindings {
				if len(name) > maxWidth {
					maxWidth = len(name)
				}
				keys[i] = name
				i += 1
			}
			sort.Strings(keys)
			for _, name := range keys {
				context.report(fmt.Sprintf("%*s %s", -maxWidth-2, name, bindings[name].Display()))
			}
			return value.NewNil(), nil
		},
	},

	PrimitiveDesc{
		"go", 0, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			context, ok := ctxt.(*Context)
			if !ok {
				return nil, fmt.Errorf("Problem understanding context")
			}
			if len(args) == 0 {
				context.nextCurrentModule = context.homeModule
				return value.NewNil(), nil
			}
			if err := checkArgType(name, args[0], IsSymbol); err != nil {
				return nil, err
			}
			context.nextCurrentModule = args[0].GetString()
			return value.NewNil(), nil
		},
	},

	PrimitiveDesc{
		"modules", 0, 0,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			context, ok := ctxt.(*Context)
			if !ok {
				return nil, fmt.Errorf("Problem understanding context")
			}
			var result value.Value = value.NewEmpty()
			for m := range context.ecosystem.modules {
				result = value.NewCons(value.NewSymbol(m), result)
			}
			return result, nil
		},
	},

	PrimitiveDesc{
		"help", 0, 0,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
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
			return value.NewNil(), nil
		},
	},

	PrimitiveDesc{
		"print", 0, -1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			for _, arg := range args {
				fmt.Print(arg.Display(), " ")
			}
			fmt.Println()
			return value.NewNil(), nil
		},
	},

	PrimitiveDesc{
		"load", 1, 1,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			context, ok := ctxt.(*Context)
			if !ok {
				return nil, fmt.Errorf("Problem understanding context")
			}
			if err := checkArgType(name, args[0], IsString); err != nil {
				return nil, err
			}
			filename := args[0].GetString()
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			str := string(b)
			if err := context.readAll(str, context); err != nil {
				return nil, err
			}
			return value.NewNil(), nil
		},
	},

	PrimitiveDesc{
		"timed-apply", 2, 2,
		func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
			context, ok := ctxt.(*Context)
			if !ok {
				return nil, fmt.Errorf("Problem understanding context")
			}
			timeTrack := func(start time.Time) {
				elapsed := time.Since(start)
				context.report(fmt.Sprintf("Time: %s", elapsed))
			}
			defer timeTrack(time.Now())
			if err := checkArgType(name, args[0], IsFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			arguments := make([]value.Value, listLength(args[1]))
			current := args[1]
			for i := range arguments {
				arguments[i] = current.GetHead()
				current = current.GetTail()
			}
			if !value.IsEmpty(current) {
				return nil, fmt.Errorf("%s: malformed list", name)
			}
			return args[0].Apply(arguments, ctxt)
		},
	},
}