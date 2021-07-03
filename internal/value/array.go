package value

import (
	"fmt"
	"strings"
)

type Array struct {
	content []Value
}

func NewArray(content []Value) *Array {
	return &Array{content}
}

func (v *Array) Display() string {
	s := make([]string, len(v.content))
	for i, vv := range v.content {
		s[i] = vv.Display()
	}
	return fmt.Sprintf("#[%s]", strings.Join(s, " "))
}

func (v *Array) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Array) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Array) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

/* old code to make it possible to use application to get/set arrays

        if len(args) < 1 || !IsNumber(args[0]) {
		return nil, fmt.Errorf("array indexing requires an index")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("too many arguments %d to array update", len(args))
	}
	idx := args[0].GetInt()
	if idx < 0 || idx >= len(v.content) {
		return nil, fmt.Errorf("array index out of bounds %d", idx)
	}
	if len(args) == 2 {
		v.content[idx] = args[1]
		return NewNil(), nil
	}
	return v.content[idx], nil
*/

func (v *Array) Str() string {
	s := make([]string, len(v.content))
	for i, vv := range v.content {
		s[i] = vv.Str()
	}
	return fmt.Sprintf("Array[%s]", strings.Join(s, " "))
}

func (v *Array) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Array) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Array) IsTrue() bool {
	return false
}

func (v *Array) IsEqual(vv Value) bool {
	return v == vv // pointer equality because arrays will be mutable
	/*
		if !IsArray(vv) || len(v.content) != len(vv.GetArray()) {
			return false}
		vvcontent := vv.GetArray()
		for i := range(v.content) {
			if !v.content[i].IsEqual(vvcontent[i]) {
				return false
			}
		}
		return true
	*/
}

func (v *Array) Kind() Kind {
	return V_ARRAY
}

func (v *Array) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Array) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Array) GetArray() []Value {
	return v.content
}

func (v *Array) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
