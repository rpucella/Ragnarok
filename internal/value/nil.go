package value

import (
	"fmt"
)

type Nil struct {
}

var defaultNil *Nil = &Nil{}

func NewNil() *Nil {
	return defaultNil
}

func (v *Nil) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Nil) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Nil) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *Nil) Display() string {
	return ""
}

func (v *Nil) Str() string {
	return fmt.Sprintf("#nil")
}

func (v *Nil) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Nil) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Nil) IsTrue() bool {
	return false
}

func (v *Nil) IsEqual(vv Value) bool {
	return IsNil(vv)
}

func (v *Nil) Kind() Kind {
	return V_NIL
}

func (v *Nil) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Nil) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Nil) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Nil) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
