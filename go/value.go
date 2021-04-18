 package main

import "fmt"
import "strings"

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
	isNumber() bool
	isBool() bool
	isRef() bool
	isString() bool
	isFunction() bool
	isTrue() bool
	isNil() bool
	typ() string
	getValue() Value
	setValue(Value)
}

type VInteger struct {
	val int
}

type VBoolean struct {
	val bool
}

type VPrimitive struct {
	name      string
	primitive func([]Value) (Value, error)
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

type VNil struct {
}

type VReference struct {
	content Value
}

// dictionary?
// what are possible keys?
// anything immutable
// probably need to define hashes

/*
type VDict struct {
	dict map[Value]Value
}
*/
  
func (v *VInteger) display() string {
	return fmt.Sprintf("%d", v.val)
}

func (v *VInteger) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VInteger) intValue() int {
	return v.val
}

func (v *VInteger) strValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VInteger) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VInteger) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *VInteger) str() string {
	return fmt.Sprintf("VInteger[%d]", v.val)
}

func (v *VInteger) headValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VInteger) tailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
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

func (v *VInteger) isNumber() bool {
	return true
}

func (v *VInteger) isBool() bool {
	return false
}

func (v *VInteger) isRef() bool {
	return false
}

func (v *VInteger) isString() bool {
	return false
}

func (v *VInteger) isFunction() bool {
	return false
}

func (v *VInteger) isTrue() bool {
	return v.val != 0
}

func (v *VInteger) isNil() bool {
	return false
}

func (v *VInteger) typ() string {
	return "int"
}

func (v *VInteger) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VInteger) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VBoolean) display() string {
	if v.val {
		return "#t"
	} else {
		return "#f"
	}
}

func (v *VBoolean) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VBoolean) intValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VBoolean) strValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
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
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VBoolean) tailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
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

func (v *VBoolean) isNumber() bool {
	return false
}

func (v *VBoolean) isBool() bool {
	return true
}

func (v *VBoolean) isRef() bool {
	return false
}

func (v *VBoolean) isString() bool {
	return false
}

func (v *VBoolean) isFunction() bool {
	return false
}

func (v *VBoolean) isTrue() bool {
	return v.val
}

func (v *VBoolean) isNil() bool {
	return false
}

func (v *VBoolean) typ() string {
	return "bool"
}

func (v *VBoolean) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VBoolean) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VPrimitive) display() string {
	return fmt.Sprintf("#<prim %s>", v.name)
}

func (v *VPrimitive) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VPrimitive) intValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VPrimitive) strValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VPrimitive) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VPrimitive) apply(args []Value) (Value, error) {
	return v.primitive(args)
}

func (v *VPrimitive) str() string {
	return fmt.Sprintf("VPrimitive[%s]", v.name)
}

func (v *VPrimitive) headValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VPrimitive) tailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
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

func (v *VPrimitive) isNumber() bool {
	return false
}

func (v *VPrimitive) isBool() bool {
	return false
}

func (v *VPrimitive) isRef() bool {
	return false
}

func (v *VPrimitive) isString() bool {
	return false
}

func (v *VPrimitive) isFunction() bool {
	return true
}

func (v *VPrimitive) isTrue() bool {
	return true
}

func (v *VPrimitive) isNil() bool {
	return false
}

func (v *VPrimitive) typ() string {
	return "fun"
}

func (v *VPrimitive) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VPrimitive) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VEmpty) display() string {
	return "()"
}

func (v *VEmpty) displayCDR() string {
	return ")"
}

func (v *VEmpty) intValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VEmpty) strValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VEmpty) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VEmpty) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *VEmpty) str() string {
	return fmt.Sprintf("VEmpty")
}

func (v *VEmpty) headValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VEmpty) tailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
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

func (v *VEmpty) isNumber() bool {
	return false
}

func (v *VEmpty) isBool() bool {
	return false
}

func (v *VEmpty) isRef() bool {
	return false
}

func (v *VEmpty) isString() bool {
	return false
}

func (v *VEmpty) isFunction() bool {
	return false
}

func (v *VEmpty) isTrue() bool {
	return false
}

func (v *VEmpty) isNil() bool {
	return false
}

func (v *VEmpty) typ() string {
	return "list"
}

func (v *VEmpty) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VEmpty) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VCons) display() string {
	return "(" + v.head.display() + v.tail.displayCDR()
}

func (v *VCons) displayCDR() string {
	return " " + v.head.display() + v.tail.displayCDR()
}

func (v *VCons) intValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VCons) strValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VCons) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
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

func (v *VCons) isNumber() bool {
	return false
}

func (v *VCons) isBool() bool {
	return false
}

func (v *VCons) isRef() bool {
	return false
}

func (v *VCons) isString() bool {
	return false
}

func (v *VCons) isFunction() bool {
	return false
}

func (v *VCons) isTrue() bool {
	return true
}

func (v *VCons) isNil() bool {
	return false
}

func (v *VCons) typ() string {
	return "list"
}

func (v *VCons) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VCons) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VSymbol) display() string {
	return v.name
}

func (v *VSymbol) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VSymbol) intValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VSymbol) strValue() string {
	return v.name
}

func (v *VSymbol) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VSymbol) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *VSymbol) str() string {
	return fmt.Sprintf("VSymbol[%s]", v.name)
}

func (v *VSymbol) headValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VSymbol) tailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
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

func (v *VSymbol) isNumber() bool {
	return false
}

func (v *VSymbol) isBool() bool {
	return false
}

func (v *VSymbol) isRef() bool {
	return false
}

func (v *VSymbol) isString() bool {
	return false
}

func (v *VSymbol) isFunction() bool {
	return false
}

func (v *VSymbol) isTrue() bool {
	return true
}

func (v *VSymbol) isNil() bool {
	return false
}

func (v *VSymbol) typ() string {
	return "symbol"
}

func (v *VSymbol) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VSymbol) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VFunction) display() string {
	return fmt.Sprintf("#<fun %s ...>", strings.Join(v.params, " "))
}

func (v *VFunction) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VFunction) intValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VFunction) strValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VFunction) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VFunction) apply(args []Value) (Value, error) {
	if len(v.params) != len(args) {
		return nil, fmt.Errorf("Wrong number of arguments to application to %s", v.str())
	}
	newEnv := v.env.layer(v.params, args)
	return v.body.eval(newEnv)
}

func (v *VFunction) str() string {
	return fmt.Sprintf("VFunction[[%s] %s]", strings.Join(v.params, " "), v.body.str())
}

func (v *VFunction) headValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VFunction) tailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
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

func (v *VFunction) isNumber() bool {
	return false
}

func (v *VFunction) isBool() bool {
	return false
}

func (v *VFunction) isRef() bool {
	return false
}

func (v *VFunction) isString() bool {
	return false
}

func (v *VFunction) isFunction() bool {
	return true
}

func (v *VFunction) isTrue() bool {
	return true
}

func (v *VFunction) isNil() bool {
	return false
}

func (v *VFunction) typ() string {
	return "fun"
}

func (v *VFunction) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VFunction) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VString) display() string {
	return "\"" + v.val + "\""
}

func (v *VString) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VString) intValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VString) strValue() string {
	return v.val
}

func (v *VString) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VString) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *VString) str() string {
	return fmt.Sprintf("VString[%s]", v.str)
}

func (v *VString) headValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VString) tailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
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

func (v *VString) isNumber() bool {
	return false
}

func (v *VString) isBool() bool {
	return false
}

func (v *VString) isRef() bool {
	return false
}

func (v *VString) isString() bool {
	return true
}

func (v *VString) isFunction() bool {
	return false
}

func (v *VString) isTrue() bool {
	return (v.val != "")
}

func (v *VString) isNil() bool {
	return false
}

func (v *VString) typ() string {
	return "string"
}

func (v *VString) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VString) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VNil) display() string {
	// figure out if this is the right thing?
	return "#nil"
}

func (v *VNil) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VNil) intValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VNil) strValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VNil) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VNil) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *VNil) str() string {
	return fmt.Sprintf("VNil")
}

func (v *VNil) headValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VNil) tailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VNil) isAtom() bool {
	return false
}

func (v *VNil) isSymbol() bool {
	return false
}

func (v *VNil) isCons() bool {
	return false
}

func (v *VNil) isEmpty() bool {
	return false
}

func (v *VNil) isNumber() bool {
	return false
}

func (v *VNil) isBool() bool {
	return false
}

func (v *VNil) isRef() bool {
	return false
}

func (v *VNil) isString() bool {
	return false
}

func (v *VNil) isFunction() bool {
	return false
}

func (v *VNil) isTrue() bool {
	return false
}

func (v *VNil) isNil() bool {
	return true
}

func (v *VNil) typ() string {
	return "nil"
}

func (v *VNil) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VNil) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VReference) display() string {
	return fmt.Sprintf("#<ref %s>", v.content.display())
}

func (v *VReference) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VReference) intValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VReference) strValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VReference) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VReference) apply(args []Value) (Value, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments %d to reference dereference", len(args))
	}
	return v.content, nil
}

func (v *VReference) str() string {
	return fmt.Sprintf("VReference[%s]", v.content.str())
}

func (v *VReference) headValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VReference) tailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *VReference) isAtom() bool {
	return false   // ?
}

func (v *VReference) isSymbol() bool {
	return false
}

func (v *VReference) isCons() bool {
	return false
}

func (v *VReference) isEmpty() bool {
	return false
}

func (v *VReference) isNumber() bool {
	return false
}

func (v *VReference) isBool() bool {
	return false
}

func (v *VReference) isRef() bool {
	return true
}

func (v *VReference) isString() bool {
	return false
}

func (v *VReference) isFunction() bool {
	return false
}

func (v *VReference) isTrue() bool {
	return false
}

func (v *VReference) isNil() bool {
	return false
}

func (v *VReference) typ() string {
	return "ref"
}

func (v *VReference) getValue() Value {
	return v.content
}

func (v *VReference) setValue(cv Value) {
	v.content = cv
}

