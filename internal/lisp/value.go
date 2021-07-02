package lisp

import "fmt"
import "strings"

type Value interface {
	Display() string
	displayCDR() string
	IntValue() int
	boolValue() bool
	StrValue() string
	HeadValue() Value
	TailValue() Value
	Apply([]Value, interface{}) (Value, error)
	Str() string
	IsAtom() bool
	IsSymbol() bool
	IsCons() bool
	IsEmpty() bool
	IsNumber() bool
	IsBool() bool
	IsRef() bool
	IsString() bool
	IsFunction() bool
	IsTrue() bool
	IsNil() bool
	//isEq() bool    -- don't think we need pointer equality for now - = is enough?
	IsEqual(Value) bool
	Kind() string
	getValue() Value
	setValue(Value)
	IsArray() bool
	getArray() []Value
	IsDict() bool
	getDict() map[string]Value
}

type VInteger struct {
	val int
}

func NewVInteger(val int) *VInteger {
	return &VInteger{val}
}

type VBoolean struct {
	val bool
}

func NewVBoolean(val bool) *VBoolean {
	return &VBoolean{val}
}

type VPrimitive struct {
	name      string
	primitive func([]Value, interface{}) (Value, error)
}

func NewVPrimitive(name string, prim func([]Value, interface{}) (Value, error)) *VPrimitive {
	return &VPrimitive{name, prim}
}

type VEmpty struct {
}

type VCons struct {
	head   Value
	tail   Value
	length int
}

func NewVCons(head Value, tail Value) *VCons {
	return &VCons{head, tail, 0} // length ignored for now...
}

func (v *VCons) SetTail(val Value) {
	v.tail = val
}

type VSymbol struct {
	name string
}

func NewVSymbol(name string) *VSymbol {
	return &VSymbol{name}
}

type VFunction struct {
	params []string
	body   AST
	env    *Env
}

func NewVFunction(params []string, body AST, env *Env) *VFunction {
	return &VFunction{params, body, env}
}

type VString struct {
	val string
}

func NewVString(val string) *VString {
	return &VString{val}
}

type VNil struct {
}

type VReference struct {
	content Value
}

func NewVReference(content Value) *VReference {
	return &VReference{content}
}

type VArray struct {
	content []Value
}

func NewVArray(content []Value) *VArray {
	return &VArray{content}
}

type VDict struct {
	content map[string]Value
}

func NewVDict(content map[string]Value) *VDict {
	return &VDict{content}
}

func (v *VInteger) Display() string {
	return fmt.Sprintf("%d", v.val)
}

func (v *VInteger) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) IntValue() int {
	return v.val
}

func (v *VInteger) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VInteger) Str() string {
	return fmt.Sprintf("VInteger[%d]", v.val)
}

func (v *VInteger) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) IsAtom() bool {
	return true
}

func (v *VInteger) IsSymbol() bool {
	return false
}

func (v *VInteger) IsCons() bool {
	return false
}

func (v *VInteger) IsEmpty() bool {
	return false
}

func (v *VInteger) IsNumber() bool {
	return true
}

func (v *VInteger) IsBool() bool {
	return false
}

func (v *VInteger) IsRef() bool {
	return false
}

func (v *VInteger) IsString() bool {
	return false
}

func (v *VInteger) IsFunction() bool {
	return false
}

func (v *VInteger) IsTrue() bool {
	return v.val != 0
}

func (v *VInteger) IsNil() bool {
	return false
}

func (v *VInteger) IsEqual(vv Value) bool {
	return vv.IsNumber() && v.IntValue() == vv.IntValue()
}

func (v *VInteger) Kind() string {
	return "int"
}

func (v *VInteger) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) IsArray() bool {
	return false
}

func (v *VInteger) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) IsDict() bool {
	return false
}

func (v *VInteger) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) Display() string {
	if v.val {
		return "#t"
	} else {
		return "#f"
	}
}

func (v *VBoolean) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) boolValue() bool {
	return v.val
}

func (v *VBoolean) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VBoolean) Str() string {
	if v.val {
		return "VBoolean[true]"
	} else {
		return "VBoolean[false]"
	}
}

func (v *VBoolean) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) IsAtom() bool {
	return true
}

func (v *VBoolean) IsSymbol() bool {
	return false
}

func (v *VBoolean) IsCons() bool {
	return false
}

func (v *VBoolean) IsEmpty() bool {
	return false
}

func (v *VBoolean) IsNumber() bool {
	return false
}

func (v *VBoolean) IsBool() bool {
	return true
}

func (v *VBoolean) IsRef() bool {
	return false
}

func (v *VBoolean) IsString() bool {
	return false
}

func (v *VBoolean) IsFunction() bool {
	return false
}

func (v *VBoolean) IsTrue() bool {
	return v.val
}

func (v *VBoolean) IsNil() bool {
	return false
}

func (v *VBoolean) IsEqual(vv Value) bool {
	return vv.IsBool() && v.boolValue() == vv.boolValue()
}

func (v *VBoolean) Kind() string {
	return "bool"
}

func (v *VBoolean) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) IsArray() bool {
	return false
}

func (v *VBoolean) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) IsDict() bool {
	return false
}

func (v *VBoolean) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) Display() string {
	return fmt.Sprintf("#<prim %s>", v.name)
}

func (v *VPrimitive) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) Apply(args []Value, ctxt interface{}) (Value, error) {
	return v.primitive(args, ctxt)
}

func (v *VPrimitive) Str() string {
	return fmt.Sprintf("VPrimitive[%s]", v.name)
}

func (v *VPrimitive) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) IsAtom() bool {
	return false
}

func (v *VPrimitive) IsSymbol() bool {
	return false
}

func (v *VPrimitive) IsCons() bool {
	return false
}

func (v *VPrimitive) IsEmpty() bool {
	return false
}

func (v *VPrimitive) IsNumber() bool {
	return false
}

func (v *VPrimitive) IsBool() bool {
	return false
}

func (v *VPrimitive) IsRef() bool {
	return false
}

func (v *VPrimitive) IsString() bool {
	return false
}

func (v *VPrimitive) IsFunction() bool {
	return true
}

func (v *VPrimitive) IsTrue() bool {
	return true
}

func (v *VPrimitive) IsNil() bool {
	return false
}

func (v *VPrimitive) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *VPrimitive) Kind() string {
	return "fun"
}

func (v *VPrimitive) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) IsArray() bool {
	return false
}

func (v *VPrimitive) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) IsDict() bool {
	return false
}

func (v *VPrimitive) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) Display() string {
	return "()"
}

func (v *VEmpty) displayCDR() string {
	return ")"
}

func (v *VEmpty) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VEmpty) Str() string {
	return fmt.Sprintf("VEmpty")
}

func (v *VEmpty) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) IsAtom() bool {
	return false
}

func (v *VEmpty) IsSymbol() bool {
	return false
}

func (v *VEmpty) IsCons() bool {
	return false
}

func (v *VEmpty) IsEmpty() bool {
	return true
}

func (v *VEmpty) IsNumber() bool {
	return false
}

func (v *VEmpty) IsBool() bool {
	return false
}

func (v *VEmpty) IsRef() bool {
	return false
}

func (v *VEmpty) IsString() bool {
	return false
}

func (v *VEmpty) IsFunction() bool {
	return false
}

func (v *VEmpty) IsTrue() bool {
	return false
}

func (v *VEmpty) IsNil() bool {
	return false
}

func (v *VEmpty) IsEqual(vv Value) bool {
	return vv.IsEmpty()
}

func (v *VEmpty) Kind() string {
	return "list"
}

func (v *VEmpty) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) IsArray() bool {
	return false
}

func (v *VEmpty) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) IsDict() bool {
	return false
}

func (v *VEmpty) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) Display() string {
	return "(" + v.head.Display() + v.tail.displayCDR()
}

func (v *VCons) displayCDR() string {
	return " " + v.head.Display() + v.tail.displayCDR()
}

func (v *VCons) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VCons) Str() string {
	return fmt.Sprintf("VCons[%s %s]", v.head.Str(), v.tail.Str())
}

func (v *VCons) HeadValue() Value {
	return v.head
}

func (v *VCons) TailValue() Value {
	return v.tail
}

func (v *VCons) IsAtom() bool {
	return false
}

func (v *VCons) IsSymbol() bool {
	return false
}

func (v *VCons) IsCons() bool {
	return true
}

func (v *VCons) IsEmpty() bool {
	return false
}

func (v *VCons) IsNumber() bool {
	return false
}

func (v *VCons) IsBool() bool {
	return false
}

func (v *VCons) IsRef() bool {
	return false
}

func (v *VCons) IsString() bool {
	return false
}

func (v *VCons) IsFunction() bool {
	return false
}

func (v *VCons) IsTrue() bool {
	return true
}

func (v *VCons) IsNil() bool {
	return false
}

func (v *VCons) IsEqual(vv Value) bool {
	if !vv.IsCons() {
		return false
	}
	var curr1 Value = v
	var curr2 Value = vv
	for curr1.IsCons() {
		if !curr2.IsCons() {
			return false
		}
		if !curr1.HeadValue().IsEqual(curr2.HeadValue()) {
			return false
		}
		curr1 = curr1.TailValue()
		curr2 = curr2.TailValue()
	}
	return curr1.IsEqual(curr2) // should both be empty at the end
}

func (v *VCons) Kind() string {
	return "list"
}

func (v *VCons) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) IsArray() bool {
	return false
}

func (v *VCons) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) IsDict() bool {
	return false
}

func (v *VCons) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) Display() string {
	return v.name
}

func (v *VSymbol) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) StrValue() string {
	return v.name
}

func (v *VSymbol) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VSymbol) Str() string {
	return fmt.Sprintf("VSymbol[%s]", v.name)
}

func (v *VSymbol) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) IsAtom() bool {
	return true
}

func (v *VSymbol) IsSymbol() bool {
	return true
}

func (v *VSymbol) IsCons() bool {
	return false
}

func (v *VSymbol) IsEmpty() bool {
	return false
}

func (v *VSymbol) IsNumber() bool {
	return false
}

func (v *VSymbol) IsBool() bool {
	return false
}

func (v *VSymbol) IsRef() bool {
	return false
}

func (v *VSymbol) IsString() bool {
	return false
}

func (v *VSymbol) IsFunction() bool {
	return false
}

func (v *VSymbol) IsTrue() bool {
	return true
}

func (v *VSymbol) IsNil() bool {
	return false
}

func (v *VSymbol) IsEqual(vv Value) bool {
	return vv.IsSymbol() && v.StrValue() == vv.StrValue()
}

func (v *VSymbol) Kind() string {
	return "symbol"
}

func (v *VSymbol) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) IsArray() bool {
	return false
}

func (v *VSymbol) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) IsDict() bool {
	return false
}

func (v *VSymbol) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) Display() string {
	return fmt.Sprintf("#<fun %s ...>", strings.Join(v.params, " "))
}

func (v *VFunction) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) Apply(args []Value, ctxt interface{}) (Value, error) {
	if len(v.params) != len(args) {
		return nil, fmt.Errorf("Wrong number of arguments to application to %s", v.Str())
	}
	newEnv := v.env.Layer(v.params, args)
	return v.body.Eval(newEnv, ctxt)
}

func (v *VFunction) Str() string {
	return fmt.Sprintf("VFunction[[%s] %s]", strings.Join(v.params, " "), v.body.Str())
}

func (v *VFunction) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) IsAtom() bool {
	return false
}

func (v *VFunction) IsSymbol() bool {
	return false
}

func (v *VFunction) IsCons() bool {
	return false
}

func (v *VFunction) IsEmpty() bool {
	return false
}

func (v *VFunction) IsNumber() bool {
	return false
}

func (v *VFunction) IsBool() bool {
	return false
}

func (v *VFunction) IsRef() bool {
	return false
}

func (v *VFunction) IsString() bool {
	return false
}

func (v *VFunction) IsFunction() bool {
	return true
}

func (v *VFunction) IsTrue() bool {
	return true
}

func (v *VFunction) IsNil() bool {
	return false
}

func (v *VFunction) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *VFunction) Kind() string {
	return "fun"
}

func (v *VFunction) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) IsArray() bool {
	return false
}

func (v *VFunction) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) IsDict() bool {
	return false
}

func (v *VFunction) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) Display() string {
	return "\"" + v.val + "\""
}

func (v *VString) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) StrValue() string {
	return v.val
}

func (v *VString) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VString) Str() string {
	return fmt.Sprintf("VString[%s]", v.StrValue())
}

func (v *VString) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) IsAtom() bool {
	return true
}

func (v *VString) IsSymbol() bool {
	return false
}

func (v *VString) IsCons() bool {
	return false
}

func (v *VString) IsEmpty() bool {
	return false
}

func (v *VString) IsNumber() bool {
	return false
}

func (v *VString) IsBool() bool {
	return false
}

func (v *VString) IsRef() bool {
	return false
}

func (v *VString) IsString() bool {
	return true
}

func (v *VString) IsFunction() bool {
	return false
}

func (v *VString) IsTrue() bool {
	return (v.val != "")
}

func (v *VString) IsNil() bool {
	return false
}

func (v *VString) IsEqual(vv Value) bool {
	return vv.IsString() && v.StrValue() == vv.StrValue()
}

func (v *VString) Kind() string {
	return "string"
}

func (v *VString) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) IsArray() bool {
	return false
}

func (v *VString) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) IsDict() bool {
	return false
}

func (v *VString) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) Display() string {
	// figure out if this is the right thing?
	return "#nil"
}

func (v *VNil) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VNil) Str() string {
	return fmt.Sprintf("VNil")
}

func (v *VNil) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) IsAtom() bool {
	return false
}

func (v *VNil) IsSymbol() bool {
	return false
}

func (v *VNil) IsCons() bool {
	return false
}

func (v *VNil) IsEmpty() bool {
	return false
}

func (v *VNil) IsNumber() bool {
	return false
}

func (v *VNil) IsBool() bool {
	return false
}

func (v *VNil) IsRef() bool {
	return false
}

func (v *VNil) IsString() bool {
	return false
}

func (v *VNil) IsFunction() bool {
	return false
}

func (v *VNil) IsTrue() bool {
	return false
}

func (v *VNil) IsNil() bool {
	return true
}

func (v *VNil) IsEqual(vv Value) bool {
	return vv.IsNil()
}

func (v *VNil) Kind() string {
	return "nil"
}

func (v *VNil) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) IsArray() bool {
	return false
}

func (v *VNil) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) IsDict() bool {
	return false
}

func (v *VNil) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) Display() string {
	return fmt.Sprintf("#<ref %s>", v.content.Display())
}

func (v *VReference) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) Apply(args []Value, ctxt interface{}) (Value, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("too many arguments %d to ref update", len(args))
	}
	if len(args) == 1 {
		v.content = args[0]
		return &VNil{}, nil
	}
	return v.content, nil
}

func (v *VReference) Str() string {
	return fmt.Sprintf("VReference[%s]", v.content.Str())
}

func (v *VReference) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) IsAtom() bool {
	return false // ?
}

func (v *VReference) IsSymbol() bool {
	return false
}

func (v *VReference) IsCons() bool {
	return false
}

func (v *VReference) IsEmpty() bool {
	return false
}

func (v *VReference) IsNumber() bool {
	return false
}

func (v *VReference) IsBool() bool {
	return false
}

func (v *VReference) IsRef() bool {
	return true
}

func (v *VReference) IsString() bool {
	return false
}

func (v *VReference) IsFunction() bool {
	return false
}

func (v *VReference) IsTrue() bool {
	return false
}

func (v *VReference) IsNil() bool {
	return false
}

func (v *VReference) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *VReference) Kind() string {
	return "reference"
}

func (v *VReference) getValue() Value {
	return v.content
}

func (v *VReference) setValue(cv Value) {
	v.content = cv
}

func (v *VReference) IsArray() bool {
	return false
}

func (v *VReference) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) IsDict() bool {
	return false
}

func (v *VReference) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) Display() string {
	s := make([]string, len(v.content))
	for i, vv := range v.content {
		s[i] = vv.Display()
	}
	return fmt.Sprintf("#[%s]", strings.Join(s, " "))
}

func (v *VArray) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) Apply(args []Value, ctxt interface{}) (Value, error) {
	if len(args) < 1 || !args[0].IsNumber() {
		return nil, fmt.Errorf("array indexing requires an index")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("too many arguments %d to array update", len(args))
	}
	idx := args[0].IntValue()
	if idx < 0 || idx >= len(v.content) {
		return nil, fmt.Errorf("array index out of bounds %d", idx)
	}
	if len(args) == 2 {
		v.content[idx] = args[1]
		return &VNil{}, nil
	}
	return v.content[idx], nil
}

func (v *VArray) Str() string {
	s := make([]string, len(v.content))
	for i, vv := range v.content {
		s[i] = vv.Str()
	}
	return fmt.Sprintf("VArray[%s]", strings.Join(s, " "))
}

func (v *VArray) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) IsAtom() bool {
	return false // ?
}

func (v *VArray) IsSymbol() bool {
	return false
}

func (v *VArray) IsCons() bool {
	return false
}

func (v *VArray) IsEmpty() bool {
	return false
}

func (v *VArray) IsNumber() bool {
	return false
}

func (v *VArray) IsBool() bool {
	return false
}

func (v *VArray) IsRef() bool {
	return false
}

func (v *VArray) IsString() bool {
	return false
}

func (v *VArray) IsFunction() bool {
	return false
}

func (v *VArray) IsTrue() bool {
	return false
}

func (v *VArray) IsNil() bool {
	return false
}

func (v *VArray) IsEqual(vv Value) bool {
	return v == vv // pointer equality because arrays will be mutable
	/*
		if !vv.IsArray() || len(v.content) != len(vv.getArray()) {
			return false}
		vvcontent := vv.getArray()
		for i := range(v.content) {
			if !v.content[i].IsEqual(vvcontent[i]) {
				return false
			}
		}
		return true
	*/
}

func (v *VArray) Kind() string {
	return "array"
}

func (v *VArray) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) IsArray() bool {
	return true
}

func (v *VArray) getArray() []Value {
	return v.content
}

func (v *VArray) IsDict() bool {
	return false
}

func (v *VArray) getDict() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) Display() string {
	s := make([]string, len(v.content))
	i := 0
	for k, vv := range v.content {
		s[i] = fmt.Sprintf("(%s %s)", k, vv.Display())
		i++
	}
	return fmt.Sprintf("#(%s)", strings.Join(s, " "))
}

func (v *VDict) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) IntValue() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) StrValue() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) boolValue() bool {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) Apply(args []Value, ctxt interface{}) (Value, error) {
	if len(args) < 1 || !args[0].IsSymbol() {
		return nil, fmt.Errorf("dict indexing requires a key")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("too many arguments %d to dict update", len(args))
	}
	key := args[0].StrValue()
	if len(args) == 2 {
		v.content[key] = args[1]
		return &VNil{}, nil
	}
	result, ok := v.content[key]
	if !ok {
		return nil, fmt.Errorf("key %s not in dict", key)
	}
	return result, nil
}

func (v *VDict) Str() string {
	s := make([]string, len(v.content))
	i := 0
	for k, vv := range v.content {
		s[i] = fmt.Sprintf("[%s %s]", k, vv.Str())
		i++
	}
	return fmt.Sprintf("VDict[%s]", strings.Join(s, " "))
}

func (v *VDict) HeadValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) TailValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) IsAtom() bool {
	return false // ?
}

func (v *VDict) IsSymbol() bool {
	return false
}

func (v *VDict) IsCons() bool {
	return false
}

func (v *VDict) IsEmpty() bool {
	return false
}

func (v *VDict) IsNumber() bool {
	return false
}

func (v *VDict) IsBool() bool {
	return false
}

func (v *VDict) IsRef() bool {
	return false
}

func (v *VDict) IsString() bool {
	return false
}

func (v *VDict) IsFunction() bool {
	return false
}

func (v *VDict) IsTrue() bool {
	return false
}

func (v *VDict) IsNil() bool {
	return false
}

func (v *VDict) IsEqual(vv Value) bool {
	return v == vv // pointer equality due to mutability
}

func (v *VDict) Kind() string {
	return "dict"
}

func (v *VDict) getValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) setValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) IsArray() bool {
	return false
}

func (v *VDict) getArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) IsDict() bool {
	return true
}

func (v *VDict) getDict() map[string]Value {
	return v.content
}
