package primitives

import (
	"fmt"
	"rpucella.net/ragnarok/internal/value"
)

func mkPrimitive(name string, min int, max int, prim func(string, []value.Value, interface{}) (value.Value, error)) value.Value {

	return value.NewPrimitive(name, func(args []value.Value, ctxt interface{}) (value.Value, error) {
		if err := checkMinArgs(name, args, min); err != nil {
			return nil, err
		}
		if max >= 0 {
			if err := checkMaxArgs(name, args, max); err != nil {
				return nil, err
			}
		}
		return prim(name, args, ctxt)
	})
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

