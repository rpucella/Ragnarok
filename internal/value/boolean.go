package value

import (
	"fmt"
)

type Boolean struct {
	val bool
}

func NewBoolean(val bool) *Boolean {
	return &Boolean{val}
}

func (v *Boolean) Display() string {
	if v.val {
		return "#t"
	} else {
		return "#f"
	}
}

func (v *Boolean) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Boolean) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Boolean) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *Boolean) Str() string {
	if v.val {
		return "Boolean[true]"
	} else {
		return "Boolean[false]"
	}
}

func (v *Boolean) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Boolean) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Boolean) IsTrue() bool {
	return v.val
}

func (v *Boolean) IsEqual(vv Value) bool {
	return IsBool(vv) && v.IsTrue() == vv.IsTrue()
}

func (v *Boolean) Kind() Kind {
	return V_BOOLEAN
}

func (v *Boolean) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Boolean) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Boolean) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Boolean) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

