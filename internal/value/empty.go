package value

import (
	"fmt"
)

type Empty struct {
}

var defaultEmpty *Empty = &Empty{}

func NewEmpty() *Empty {
	return defaultEmpty
}

func (v *Empty) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Empty) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Empty) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *Empty) Display() string {
	return ""
}

func (v *Empty) Str() string {
	return fmt.Sprintf("()")
}

func (v *Empty) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Empty) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Empty) IsTrue() bool {
	return false
}

func (v *Empty) IsEqual(vv Value) bool {
	return IsEmpty(vv)
}

func (v *Empty) Kind() Kind {
	return V_EMPTY
}

func (v *Empty) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Empty) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Empty) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Empty) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
