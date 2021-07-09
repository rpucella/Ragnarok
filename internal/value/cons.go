package value

import (
	"fmt"
)

type Cons struct {
	head   Value
	tail   Value
	length int
}

func NewCons(head Value, tail Value) *Cons {
	return &Cons{head, tail, 0} // length ignored for now...
}

func (v *Cons) SetTail(val Value) {
	v.tail = val
}

func (v *Cons) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Cons) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Cons) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *Cons) Display() string {
	if IsEmpty(v.tail) {
		return v.head.Display()
	}
	result := v.head.Display()
	curr := v.tail
	for IsCons(curr) {
		result += " " + curr.GetHead().Display()
		curr = curr.GetTail()
	}
	if !IsEmpty(curr) {
		return result + " <?>"
	}
	return result
}

func (v *Cons) Str() string {
	if IsEmpty(v.tail) {
		return "(" + v.head.Str() + ")"
	}
	result := "(" + v.head.Str()
	curr := v.tail
	for IsCons(curr) {
		result += " " + curr.GetHead().Str()
		curr = curr.GetTail()
	}
	if !IsEmpty(curr) {
		return result + " <?>)"
	}
	return result + ")"
}

func (v *Cons) GetHead() Value {
	return v.head
}

func (v *Cons) GetTail() Value {
	return v.tail
}

func (v *Cons) IsTrue() bool {
	return true
}

func (v *Cons) IsEqual(vv Value) bool {
	if !IsCons(vv) {
		return false
	}
	var curr1 Value = v
	var curr2 Value = vv
	for IsCons(curr1) {
		if !IsCons(curr2) {
			return false
		}
		if !curr1.GetHead().IsEqual(curr2.GetHead()) {
			return false
		}
		curr1 = curr1.GetTail()
		curr2 = curr2.GetTail()
	}
	return curr1.IsEqual(curr2) // should both be empty at the end
}

func (v *Cons) Kind() Kind {
	return V_CONS
}

func (v *Cons) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Cons) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Cons) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *Cons) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
