package value

import (
	"fmt"
)

type Primitive struct {
	name      string
	primitive func([]Value, interface{}) (Value, error)
}

func NewPrimitive(name string, prim func([]Value, interface{}) (Value, error)) *Primitive {
	return &Primitive{name, prim}
}

func (v *Primitive) Display() string {
	return fmt.Sprintf("#[prim %s]", v.name)
}

func (v *Primitive) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Primitive) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Primitive) Apply(args []Value, ctxt interface{}) (Value, error) {
	result, err := v.primitive(args, ctxt)
	// primitives always run to completion
	return result, err
}

func (v *Primitive) Str() string {
	return fmt.Sprintf("Primitive[%s]", v.name)
}

func (v *Primitive) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Primitive) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Primitive) IsTrue() bool {
	return true
}

func (v *Primitive) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *Primitive) Kind() Kind {
	return V_FUNCTION
}

func (v *Primitive) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Primitive) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Primitive) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Primitive) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
