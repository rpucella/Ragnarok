package value

import (
	"fmt"
)

type Symbol struct {
	name string
}

func NewSymbol(name string) *Symbol {
	return &Symbol{name}
}

func (v *Symbol) Display() string {
	return v.name
}

func (v *Symbol) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Symbol) GetString() string {
	return v.name
}

func (v *Symbol) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *Symbol) Str() string {
	return fmt.Sprintf("Symbol[%s]", v.name)
}

func (v *Symbol) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Symbol) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Symbol) IsTrue() bool {
	return true
}

func (v *Symbol) IsEqual(vv Value) bool {
	return IsSymbol(vv) && v.GetString() == vv.GetString()
}

func (v *Symbol) Kind() Kind {
	return V_SYMBOL
}

func (v *Symbol) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Symbol) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Symbol) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Symbol) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

