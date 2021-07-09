package value

import (
	"fmt"
)

type Reference struct {
	content Value
}

func NewReference(content Value) *Reference {
	return &Reference{content}
}

func (v *Reference) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Reference) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Reference) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *Reference) Display() string {
	return v.content.Display()
}

func (v *Reference) Str() string {
	return fmt.Sprintf("#(ref %s)", v.content.Str())
}

func (v *Reference) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Reference) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Reference) IsTrue() bool {
	return true
}

func (v *Reference) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *Reference) Kind() Kind {
	return V_REFERENCE
}

func (v *Reference) GetValue() Value {
	return v.content
}

func (v *Reference) SetValue(cv Value) {
	v.content = cv
}

func (v *Reference) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Reference) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
