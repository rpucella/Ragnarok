package value

import (
	"fmt"
)

type Integer struct {
	val int
}

func NewInteger(val int) *Integer {
	return &Integer{val}
}

func (v *Integer) GetInt() int {
	return v.val
}

func (v *Integer) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Integer) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *Integer) Display() string {
	return fmt.Sprintf("%d", v.val)
}

func (v *Integer) Str() string {
	return fmt.Sprintf("%d", v.val)
}

func (v *Integer) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Integer) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Integer) IsTrue() bool {
	return v.val != 0
}

func (v *Integer) IsEqual(vv Value) bool {
	return IsNumber(vv) && v.GetInt() == vv.GetInt()
}

func (v *Integer) Kind() Kind {
	return V_INTEGER
}

func (v *Integer) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Integer) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Integer) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Integer) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
