package value

import (
	"fmt"
	"strings"
)

type Dict struct {
	content map[string]Value
}

func NewDict(content map[string]Value) *Dict {
	return &Dict{content}
}

func (v *Dict) Display() string {
	s := make([]string, len(v.content))
	i := 0
	for k, vv := range v.content {
		s[i] = fmt.Sprintf("(%s %s)", k, vv.Display())
		i++
	}
	return fmt.Sprintf("#(%s)", strings.Join(s, " "))
}

func (v *Dict) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Dict) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Dict) Apply(args []Value, ctxt interface{}) (Value, error) {
	if len(args) < 1 || !IsSymbol(args[0]) {
		return nil, fmt.Errorf("dict indexing requires a key")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("too many arguments %d to dict update", len(args))
	}
	key := args[0].GetString()
	if len(args) == 2 {
		v.content[key] = args[1]
		return NewNil(), nil
	}
	result, ok := v.content[key]
	if !ok {
		return nil, fmt.Errorf("key %s not in dict", key)
	}
	return result, nil
}

func (v *Dict) Str() string {
	s := make([]string, len(v.content))
	i := 0
	for k, vv := range v.content {
		s[i] = fmt.Sprintf("[%s %s]", k, vv.Str())
		i++
	}
	return fmt.Sprintf("Dict[%s]", strings.Join(s, " "))
}

func (v *Dict) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Dict) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Dict) IsTrue() bool {
	return false
}

func (v *Dict) IsEqual(vv Value) bool {
	return v == vv // pointer equality due to mutability
}

func (v *Dict) Kind() Kind {
	return V_DICT
}

func (v *Dict) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Dict) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Dict) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Dict) GetMap() map[string]Value {
	return v.content
}
