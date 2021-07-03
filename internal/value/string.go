package value

import (
	"fmt"
)

type String struct {
	val string
}

func NewString(val string) *String {
	return &String{val}
}

func (v *String) Display() string {
	return "\"" + v.val + "\""
}

func (v *String) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *String) GetString() string {
	return v.val
}

func (v *String) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *String) Str() string {
	return fmt.Sprintf("String[%s]", v.GetString())
}

func (v *String) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *String) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *String) IsTrue() bool {
	return (v.val != "")
}

func (v *String) IsEqual(vv Value) bool {
	return IsString(vv) && v.GetString() == vv.GetString()
}

func (v *String) Kind() Kind {
	return V_STRING
}

func (v *String) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *String) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *String) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *String) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
