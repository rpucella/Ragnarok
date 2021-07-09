package primitives

import (
	"fmt"
	"rpucella.net/ragnarok/internal/value"
)

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

var PrimType = mkPrimitive(
	"type", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewSymbol(value.Classify(args[0])), nil
	})

var PrimPlus = mkPrimitive(
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
	})

var PrimTimes = mkPrimitive(
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
	})

var PrimMinus = mkPrimitive(
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
	})

var PrimEqual = mkPrimitive("=", 2, -1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		var reference value.Value = args[0]
		for _, v := range args[1:] {
			if !value.IsEqual(reference, v) {
				return value.NewBoolean(false), nil
			}
		}
		return value.NewBoolean(true), nil
	})

var PrimLess = mkPrimitive("<", 2, 2,
	mkNumPredicate(func(n1 int, n2 int) bool { return n1 < n2 }))

var PrimLessEqual = mkPrimitive("<=", 2, 2,
	mkNumPredicate(func(n1 int, n2 int) bool { return n1 <= n2 }))

var PrimMore = mkPrimitive(">", 2, 2,
	mkNumPredicate(func(n1 int, n2 int) bool { return n1 > n2 }))

var PrimMoreEqual = mkPrimitive(">=", 2, 2,
	mkNumPredicate(func(n1 int, n2 int) bool { return n1 >= n2 }))

var PrimNot = mkPrimitive("not", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(!value.IsTrue(args[0])), nil
	})

var PrimApply = mkPrimitive("apply", 2, 2,
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
	})

var PrimCons = mkPrimitive("cons", 2, 2,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		if err := checkArgType(name, args[1], isList); err != nil {
			return nil, err
		}
		return value.NewCons(args[0], args[1]), nil
	})

var PrimAppend = mkPrimitive(
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
	})

var PrimReverse = mkPrimitive("reverse", 1, 1,
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
	})

var PrimHead = mkPrimitive("head", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		if err := checkArgType(name, args[0], isList); err != nil {
			return nil, err
		}
		if value.IsEmpty(args[0]) {
			return nil, fmt.Errorf("%s: empty list argument", name)
		}
		return args[0].GetHead(), nil
	})

var PrimTail = mkPrimitive("tail", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		if err := checkArgType(name, args[0], isList); err != nil {
			return nil, err
		}
		if value.IsEmpty(args[0]) {
			return nil, fmt.Errorf("%s: empty list argument", name)
		}
		return args[0].GetTail(), nil
	})

var PrimList = mkPrimitive("list", 0, -1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		var result value.Value = value.NewEmpty()
		for i := len(args) - 1; i >= 0; i -= 1 {
			result = value.NewCons(args[i], result)
		}
		return result, nil
	})

var PrimLength = mkPrimitive("length", 1, 1,
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
	})

var PrimNth = mkPrimitive("nth", 2, 2,
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
	})

var PrimMap = mkPrimitive("map", 2, -1,
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
	})

var PrimFor = mkPrimitive("for", 2, -1,
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
	})

var PrimFilter = mkPrimitive("filter", 2, 2,
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
	})

var PrimFoldr = mkPrimitive("foldr", 3, 3,
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
	})

var PrimFoldl = mkPrimitive("foldl", 3, 3,
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
	})

var PrimRef = mkPrimitive("ref", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewReference(args[0]), nil
	})

// set should probably be a special form
// (set (x) 10)
// (set (arr 1) 10)
// (set (dict 'a) 10)
// like setf in CLISP

var PrimGet = mkPrimitive("get", 1, 2,
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
	})

var PrimSet = mkPrimitive("set!", 2, 3,
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
	})

var PrimEmptyP = mkPrimitive("empty?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsEmpty(args[0])), nil
	})

var PrimConsP = mkPrimitive("cons?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsCons(args[0])), nil
	})

var PrimListP = mkPrimitive("list?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsCons(args[0]) || value.IsEmpty(args[0])), nil
	})

var PrimNumberP = mkPrimitive("number?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsNumber(args[0])), nil
	})

var PrimRefP = mkPrimitive("ref?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsRef(args[0])), nil
	})

var PrimBooleanP = mkPrimitive("boolean?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsBool(args[0])), nil
	})

var PrimStringP = mkPrimitive("string?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsString(args[0])), nil
	})

var PrimSymbolP = mkPrimitive("symbol?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsSymbol(args[0])), nil
	})

var PrimFunctionP = mkPrimitive("function?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsFunction(args[0])), nil
	})

var PrimNilP = mkPrimitive("nil?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsNil(args[0])), nil
	})

var PrimArray = mkPrimitive("array", 0, -1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		content := make([]value.Value, len(args))
		for i, v := range args {
			content[i] = v
		}
		return value.NewArray(content), nil
	})

var PrimArrayP = mkPrimitive("array?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsArray(args[0])), nil
	})

var PrimDict = mkPrimitive("dict", 0, -1,
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
	})

var PrimDictP = mkPrimitive("dict?", 1, 1,
	func(name string, args []value.Value, ctxt interface{}) (value.Value, error) {
		return value.NewBoolean(value.IsDict(args[0])), nil
	})
