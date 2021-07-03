package value

import (
	"fmt"
	"strings"
)

type Kind int

const (
	V_INTEGER Kind = iota
	V_BOOLEAN
	V_FUNCTION
	V_EMPTY
	V_CONS
	V_SYMBOL
	V_STRING
	V_NIL
	V_REFERENCE
	V_ARRAY
	V_DICT
)

type Value interface {
	Kind() Kind
	Display() string
	Str() string
	GetInt() int
	GetString() string
	GetHead() Value
	GetTail() Value
	GetValue() Value
	SetValue(Value)
	GetArray() []Value
	GetMap() map[string]Value
	Apply([]Value, interface{}) (Value, error)
	// internal methods
	IsTrue() bool
	IsEqual(Value) bool
	//isEq() bool    -- don't think we need pointer equality for now - = is enough?
}

func IsEqual(v1 Value, v2 Value) bool {
	return v1.IsEqual(v2)
}

func IsTrue(v Value) bool {
	return v.IsTrue()
}

func IsAtom(v Value) bool {
	k := v.Kind()
	return k == V_INTEGER || k == V_BOOLEAN ||
		k == V_SYMBOL || k == V_STRING
}

func IsSymbol(v Value) bool {
	return v.Kind() == V_SYMBOL
}

func IsCons(v Value) bool {
	return v.Kind() == V_CONS
}

func IsEmpty(v Value) bool {
	return v.Kind() == V_EMPTY
}

func IsNumber(v Value) bool {
	return v.Kind() == V_INTEGER
}

func IsBool(v Value) bool {
	return v.Kind() == V_BOOLEAN
}

func IsRef(v Value) bool {
	return v.Kind() == V_REFERENCE
}

func IsString(v Value) bool {
	return v.Kind() == V_STRING
}

func IsFunction(v Value) bool {
	k := v.Kind() 
	return k == V_FUNCTION
}

func IsNil(v Value) bool {
	return v.Kind() == V_NIL
}

func IsArray(v Value) bool {
	return v.Kind() == V_ARRAY
}

func IsDict(v Value) bool {
	return v.Kind() == V_DICT
}

func Classify(v Value) string {
	switch v.Kind() {
	case V_INTEGER:
		return "int"
	case V_BOOLEAN:
		return "bool"
	case V_EMPTY:
		return "list"
	case V_CONS:
		return "list"
	case V_SYMBOL:
		return "symbol"
	case V_FUNCTION:
		return "fun"
	case V_STRING:
		return "string"
	case V_NIL:
		return "nil"
	case V_REFERENCE:
		return "reference"
	case V_ARRAY:
		return "array"
	case V_DICT:
		return "dict"
	default:
		return "?"
	}
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

/*
type VFunction struct {
	params []string
	function func([]Value, interface{}) (Value, bool, error)
}

func NewVFunction(params []string, function func([]Value, interface{}) (Value, bool, error)) *VFunction {
	return &VFunction{params, function}
}
*/

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

func (v *VInteger) GetInt() int {
	return v.val
}

func (v *VInteger) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VInteger) Str() string {
	return fmt.Sprintf("VInteger[%d]", v.val)
}

func (v *VInteger) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) IsTrue() bool {
	return v.val != 0
}

func (v *VInteger) IsEqual(vv Value) bool {
	return IsNumber(vv) && v.GetInt() == vv.GetInt()
}

func (v *VInteger) Kind() Kind {
	return V_INTEGER
}

func (v *VInteger) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VInteger) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) Display() string {
	if v.val {
		return "#t"
	} else {
		return "#f"
	}
}

func (v *VBoolean) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
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

func (v *VBoolean) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) IsTrue() bool {
	return v.val
}

func (v *VBoolean) IsEqual(vv Value) bool {
	return IsBool(vv) && v.IsTrue() == vv.IsTrue()
}

func (v *VBoolean) Kind() Kind {
	return V_BOOLEAN
}

func (v *VBoolean) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VBoolean) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) Display() string {
	return fmt.Sprintf("#<prim %s>", v.name)
}

func (v *VPrimitive) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) Apply(args []Value, ctxt interface{}) (Value, error) {
	result, err := v.primitive(args, ctxt)
	// primitives always run to completion
	return result, err
}

func (v *VPrimitive) Str() string {
	return fmt.Sprintf("VPrimitive[%s]", v.name)
}

func (v *VPrimitive) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) IsTrue() bool {
	return true
}

func (v *VPrimitive) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *VPrimitive) Kind() Kind {
	return V_FUNCTION
}

func (v *VPrimitive) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VPrimitive) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) Display() string {
	return "()"
}

func (v *VEmpty) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VEmpty) Str() string {
	return fmt.Sprintf("VEmpty")
}

func (v *VEmpty) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) IsTrue() bool {
	return false
}

func (v *VEmpty) IsEqual(vv Value) bool {
	return IsEmpty(vv)
}

func (v *VEmpty) Kind() Kind {
	return V_EMPTY
}

func (v *VEmpty) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VEmpty) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) Display() string {
	if IsEmpty(v.tail) {
		return "(" + v.head.Display() + ")"
	}
	result := "(" + v.head.Display()
	curr := v.tail
	for IsCons(curr) {
		result += " " + curr.GetHead().Display()
		curr = curr.GetTail()
	}
	if !IsEmpty(curr) {
		return result + " <?>)"
	}
	return result + ")"
}

func (v *VCons) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VCons) Str() string {
	return fmt.Sprintf("VCons[%s %s]", v.head.Str(), v.tail.Str())
}

func (v *VCons) GetHead() Value {
	return v.head
}

func (v *VCons) GetTail() Value {
	return v.tail
}

func (v *VCons) IsTrue() bool {
	return true
}

func (v *VCons) IsEqual(vv Value) bool {
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

func (v *VCons) Kind() Kind {
	return V_CONS
}

func (v *VCons) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VCons) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) Display() string {
	return v.name
}

func (v *VSymbol) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) GetString() string {
	return v.name
}

func (v *VSymbol) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VSymbol) Str() string {
	return fmt.Sprintf("VSymbol[%s]", v.name)
}

func (v *VSymbol) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) IsTrue() bool {
	return true
}

func (v *VSymbol) IsEqual(vv Value) bool {
	return IsSymbol(vv) && v.GetString() == vv.GetString()
}

func (v *VSymbol) Kind() Kind {
	return V_SYMBOL
}

func (v *VSymbol) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VSymbol) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

/*
func (v *VFunction) Display() string {
	return fmt.Sprintf("#<fun %s>", strings.Join(v.params, " "))
}

func (v *VFunction) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) Apply(args []Value, ctxt interface{}) (Value, error) {
	if len(v.params) != len(args) {
		return nil, fmt.Errorf("Wrong number of arguments in application to %s", v.Str())
	}
	return v.function(args, ctxt)
}

func (v *VFunction) Str() string {
	return fmt.Sprintf("VFunction[%s]", strings.Join(v.params, " "))
}

func (v *VFunction) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) IsTrue() bool {
	return true
}

func (v *VFunction) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *VFunction) Kind() Kind {
	return V_FUNCTION
}

func (v *VFunction) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VFunction) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
*/

func (v *VString) Display() string {
	return "\"" + v.val + "\""
}

func (v *VString) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) GetString() string {
	return v.val
}

func (v *VString) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VString) Str() string {
	return fmt.Sprintf("VString[%s]", v.GetString())
}

func (v *VString) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) IsTrue() bool {
	return (v.val != "")
}

func (v *VString) IsEqual(vv Value) bool {
	return IsString(vv) && v.GetString() == vv.GetString()
}

func (v *VString) Kind() Kind {
	return V_STRING
}

func (v *VString) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VString) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) Display() string {
	// figure out if this is the right thing?
	return "#nil"
}

func (v *VNil) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) Apply(args []Value, ctxt interface{}) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.Str())
}

func (v *VNil) Str() string {
	return fmt.Sprintf("VNil")
}

func (v *VNil) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) IsTrue() bool {
	return false
}

func (v *VNil) IsEqual(vv Value) bool {
	return IsNil(vv)
}

func (v *VNil) Kind() Kind {
	return V_NIL
}

func (v *VNil) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VNil) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) Display() string {
	return fmt.Sprintf("#<ref %s>", v.content.Display())
}

func (v *VReference) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) GetString() string {
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

func (v *VReference) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) IsTrue() bool {
	return false
}

func (v *VReference) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *VReference) Kind() Kind {
	return V_REFERENCE
}

func (v *VReference) GetValue() Value {
	return v.content
}

func (v *VReference) SetValue(cv Value) {
	v.content = cv
}

func (v *VReference) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VReference) GetMap() map[string]Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) Display() string {
	s := make([]string, len(v.content))
	for i, vv := range v.content {
		s[i] = vv.Display()
	}
	return fmt.Sprintf("#[%s]", strings.Join(s, " "))
}

func (v *VArray) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) Apply(args []Value, ctxt interface{}) (Value, error) {
	if len(args) < 1 || !IsNumber(args[0]) {
		return nil, fmt.Errorf("array indexing requires an index")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("too many arguments %d to array update", len(args))
	}
	idx := args[0].GetInt()
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

func (v *VArray) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) IsTrue() bool {
	return false
}

func (v *VArray) IsEqual(vv Value) bool {
	return v == vv // pointer equality because arrays will be mutable
	/*
		if !IsArray(vv) || len(v.content) != len(vv.GetArray()) {
			return false}
		vvcontent := vv.GetArray()
		for i := range(v.content) {
			if !v.content[i].IsEqual(vvcontent[i]) {
				return false
			}
		}
		return true
	*/
}

func (v *VArray) Kind() Kind {
	return V_ARRAY
}

func (v *VArray) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VArray) GetArray() []Value {
	return v.content
}

func (v *VArray) GetMap() map[string]Value {
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

func (v *VDict) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) Apply(args []Value, ctxt interface{}) (Value, error) {
	if len(args) < 1 || !IsSymbol(args[0]) {
		return nil, fmt.Errorf("dict indexing requires a key")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("too many arguments %d to dict update", len(args))
	}
	key := args[0].GetString()
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

func (v *VDict) GetHead() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) GetTail() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) IsTrue() bool {
	return false
}

func (v *VDict) IsEqual(vv Value) bool {
	return v == vv // pointer equality due to mutability
}

func (v *VDict) Kind() Kind {
	return V_DICT
}

func (v *VDict) GetValue() Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) SetValue(cv Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) GetArray() []Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *VDict) GetMap() map[string]Value {
	return v.content
}
