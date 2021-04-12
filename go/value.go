package main

import "fmt"

type Value interface {
	display() string
	displayCDR() string
	intValue() int
	boolValue() bool
	strValue() string
	headValue() Value
	tailValue() Value
	apply([]Value) Value
	str() string
	isAtom() bool
	isSymbol() bool
	isCons() bool
	/*
		type() string
		isNumber() bool
		isBoolean() bool
		isString() bool
		isSymbol() bool
		isNil() bool
		isEmpty() bool
		isCons() bool
		isFunction() bool
		isMacro() bool
		isReference() bool
		isAtom() bool
		isList() bool
		isDict() bool
		isTrue() bool
		isEqual(Value) bool
		isEq(Value) bool
	*/
}

type VInteger struct {
	val int
}

func (v *VInteger) display() string {
	return fmt.Sprintf("%d", v.val)
}

func (v *VInteger) displayCDR() string {
	panic("Boom!")
}

func (v *VInteger) intValue() int {
	return v.val
}

func (v *VInteger) strValue() string {
	panic("Boom!")
}

func (v *VInteger) boolValue() bool {
	panic("Boom!")
}

func (v *VInteger) apply(args []Value) Value {
	panic("Boom!")
}

func (v *VInteger) str() string {
	return fmt.Sprintf("VInteger[%d]", v.val)
}

func (v *VInteger) headValue() Value {
	panic("Boom!")
}

func (v *VInteger) tailValue() Value {
	panic("Boom!")
}

func (v *VInteger) isAtom() bool {
	return true
}

func (v *VInteger) isSymbol() bool {
	return false
}

func (v *VInteger) isCons() bool {
	return false
}



type VBoolean struct {
	val bool
}

func (v *VBoolean) display() string {
	if v.val {
		return "#t"
	} else {
		return "#f"
	}
}

func (v *VBoolean) displayCDR() string {
	panic("Boom!")
}

func (v *VBoolean) intValue() int {
	panic("Boom!")
}

func (v *VBoolean) strValue() string {
	panic("Boom!")
}

func (v *VBoolean) boolValue() bool {
	return v.val
}

func (v *VBoolean) apply(args []Value) Value {
	panic("Boom!")
}

func (v *VBoolean) str() string {
	if v.val {
		return "VBoolean[true]"
	} else {
		return "VBoolean[false]"
	}
}

func (v *VBoolean) headValue() Value {
	panic("Boom!")
}

func (v *VBoolean) tailValue() Value {
	panic("Boom!")
}

func (v *VBoolean) isAtom() bool {
	return true
}

func (v *VBoolean) isSymbol() bool {
	return false
}

func (v *VBoolean) isCons() bool {
	return false
}

type VPrimitive struct {
	name      string
	primitive func([]Value) Value
	numArgs   int
}

func (v *VPrimitive) display() string {
	return fmt.Sprintf("#<PRIMITIVE %s>", v.name)
}

func (v *VPrimitive) displayCDR() string {
	panic("Boom!")
}

func (v *VPrimitive) intValue() int {
	panic("Boom!")
}

func (v *VPrimitive) strValue() string {
	panic("Boom!")
}

func (v *VPrimitive) boolValue() bool {
	panic("Boom!")
}

func (v *VPrimitive) apply(args []Value) Value {
	// check length?
	return v.primitive(args)
}

func (v *VPrimitive) str() string {
	return fmt.Sprintf("VPrimitive[%s]", v.name)
}

func (v *VPrimitive) headValue() Value {
	panic("Boom!")
}

func (v *VPrimitive) tailValue() Value {
	panic("Boom!")
}

func (v *VPrimitive) isAtom() bool {
	return false
}

func (v *VPrimitive) isSymbol() bool {
	return false
}

func (v *VPrimitive) isCons() bool {
	return false
}

type VEmpty struct {
}

func (v *VEmpty) display() string {
	return "()"
}

func (v *VEmpty) displayCDR() string {
	return ")"
}

func (v *VEmpty) intValue() int {
	panic("Boom!")
}

func (v *VEmpty) strValue() string {
	panic("Boom!")
}

func (v *VEmpty) boolValue() bool {
	panic("Boom!")
}

func (v *VEmpty) apply(args []Value) Value {
	panic("Boom!")
}

func (v *VEmpty) str() string {
	return fmt.Sprintf("VEmpty")
}

func (v *VEmpty) headValue() Value {
	panic("Boom!")
}

func (v *VEmpty) tailValue() Value {
	panic("Boom!")
}

func (v *VEmpty) isAtom() bool {
	return false
}

func (v *VEmpty) isSymbol() bool {
	return false
}

func (v *VEmpty) isCons() bool {
	return false
}

type VCons struct {
	head   Value
	tail   Value
	length int
}

func (v *VCons) display() string {
	return "(" + v.head.display() + v.tail.displayCDR()
}

func (v *VCons) displayCDR() string {
	return " " + v.head.display() + v.tail.displayCDR()
}

func (v *VCons) intValue() int {
	panic("Boom!")
}

func (v *VCons) strValue() string {
	panic("Boom!")
}

func (v *VCons) boolValue() bool {
	panic("Boom!")
}

func (v *VCons) apply(args []Value) Value {
	panic("Boom!")
}

func (v *VCons) str() string {
	return fmt.Sprintf("VCons[%s %s]", v.head.str(), v.tail.str())
}

func (v *VCons) headValue() Value {
	return v.head
}

func (v *VCons) tailValue() Value {
	return v.tail
}

func (v *VCons) isAtom() bool {
	return false
}

func (v *VCons) isSymbol() bool {
	return false
}

func (v *VCons) isCons() bool {
	return true
}

type VSymbol struct {
	name string
}

func (v *VSymbol) display() string {
	return v.name
}

func (v *VSymbol) displayCDR() string {
	panic("Boom!")
}

func (v *VSymbol) intValue() int {
	panic("Boom!")
}

func (v *VSymbol) strValue() string {
	return v.name
}

func (v *VSymbol) boolValue() bool {
	panic("Boom!")
}

func (v *VSymbol) apply(args []Value) Value {
	panic("Boom!")
}

func (v *VSymbol) str() string {
	return fmt.Sprintf("VSymbol[%s]", v.name)
}

func (v *VSymbol) headValue() Value {
	panic("Boom!")
}

func (v *VSymbol) tailValue() Value {
	panic("Boom!")
}

func (v *VSymbol) isAtom() bool {
	return true
}

func (v *VSymbol) isSymbol() bool {
	return true
}

func (v *VSymbol) isCons() bool {
	return true
}
