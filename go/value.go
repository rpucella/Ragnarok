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
	apply([]Value) (Value, error)
	str() string
	isAtom() bool
	isSymbol() bool
	isCons() bool
	isEmpty() bool
	isTrue() bool
}

type VInteger struct {
	val int
}

type VBoolean struct {
	val bool
}

type VPrimitive struct {
	name      string
	primitive func([]Value) Value
	numArgs   int
}

type VEmpty struct {
}

type VCons struct {
	head   Value
	tail   Value
	length int
}

type VSymbol struct {
	name string
}

type VFunction struct {
	params []string
	body AST
	env *Env
}

type VString struct {
	val string
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

func (v *VInteger) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
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

func (v *VInteger) isEmpty() bool {
	return false
}

func (v *VInteger) isTrue() bool {
	return v.val != 0
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

func (v *VBoolean) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
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

func (v *VBoolean) isEmpty() bool {
	return false
}

func (v *VBoolean) isTrue() bool {
	return v.val
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

func (v *VPrimitive) apply(args []Value) (Value, error) {
	// check length?
	return v.primitive(args), nil
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

func (v *VPrimitive) isEmpty() bool {
	return false
}

func (v *VPrimitive) isTrue() bool {
	return true
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

func (v *VEmpty) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
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

func (v *VEmpty) isEmpty() bool {
	return true
}

func (v *VEmpty) isTrue() bool {
	return false
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

func (v *VCons) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
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

func (v *VCons) isEmpty() bool {
	return false
}

func (v *VCons) isTrue() bool {
	return true
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

func (v *VSymbol) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
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
	return false
}

func (v *VSymbol) isEmpty() bool {
	return false
}

func (v *VSymbol) isTrue() bool {
	return true
}

func (v *VFunction) display() string {
	return fmt.Sprintf("#<FUN ...>")
}

func (v *VFunction) displayCDR() string {
	panic("Boom!")
}

func (v *VFunction) intValue() int {
	panic("Boom!")
}

func (v *VFunction) strValue() string {
	panic("Boom!")
}

func (v *VFunction) boolValue() bool {
	panic("Boom!")
}

func (v *VFunction) apply(args []Value) (Value, error) {
	if len(v.params) != len(args) {
		return nil, fmt.Errorf("Wrong number of arguments to application to %s", v.str())
	}
	new_env := newEnv(v.env)
	for i := range v.params {
		new_env.update(v.params[i], args[i])
	}
	return v.body.eval(new_env)
}

func (v *VFunction) str() string {
	return fmt.Sprintf("VFunction[?]")
}

func (v *VFunction) headValue() Value {
	panic("Boom!")
}

func (v *VFunction) tailValue() Value {
	panic("Boom!")
}

func (v *VFunction) isAtom() bool {
	return false
}

func (v *VFunction) isSymbol() bool {
	return false
}

func (v *VFunction) isCons() bool {
	return false
}

func (v *VFunction) isEmpty() bool {
	return false
}

func (v *VFunction) isTrue() bool {
	return true
}

func (v *VString) display() string {
	return "\"" + v.val + "\""
}

func (v *VString) displayCDR() string {
	panic("Boom!")
}

func (v *VString) intValue() int {
	panic("Boom!")
}

func (v *VString) strValue() string {
	return v.val
}

func (v *VString) boolValue() bool {
	panic("Boom!")
}

func (v *VString) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *VString) str() string {
	return fmt.Sprintf("VString[%s]", v.str)
}

func (v *VString) headValue() Value {
	panic("Boom!")
}

func (v *VString) tailValue() Value {
	panic("Boom!")
}

func (v *VString) isAtom() bool {
	return true
}

func (v *VString) isSymbol() bool {
	return false
}

func (v *VString) isCons() bool {
	return false
}

func (v *VString) isEmpty() bool {
	return false
}

func (v *VString) isTrue() bool {
	return (v.val != "")
}
