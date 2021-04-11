package main

import "fmt"
import "bufio"
import "os"
import "regexp"
import "strconv"
import "strings"

func main() {

	test_value_10()
	test_value_plus()

	test_literal()
	test_lookup()
	test_apply()
	test_if()

	test_parse()

	reader := bufio.NewReader(os.Stdin)
	env := sampleEnv()
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		e, _ := parseExpr(text)
		if e != nil {
			fmt.Println(e.eval(env).display())
		}
	}
}

func sampleEnv() *Env {
	current := map[string]Value{
		"a": &VInteger{10},
		"b": &VInteger{20},
		"+": &VPrimitive{"+", primitivePlus, 0},
		"t": &VBoolean{true},
		"f": &VBoolean{false},
	}
	env := &Env{current: current}
	return env
}

func test_value_10() {
	var v1 Value = &VInteger{10}
	fmt.Println(v1.str(), "->", v1.intValue())
}

func test_value_plus() {
	var v1 Value = &VInteger{10}
	var v2 Value = &VInteger{20}
	var v3 Value = &VInteger{30}
	var vp Value = &VPrimitive{"+", primitivePlus, 0}
	var args []Value = []Value{v1, v2, v3}
	fmt.Println(vp.str(), "->", vp.apply(args).intValue())
}

func test_literal() {
	v1 := &VInteger{10}
	e1 := &Literal{v1}
	fmt.Println(e1.str(), "->", e1.eval(nil).display())
	v2 := &VBoolean{true}
	e2 := &Literal{v2}
	fmt.Println(e2.str(), "->", e2.eval(nil).display())
}

func test_lookup() {
	env := sampleEnv()
	e1 := &Symbol{"a"}
	fmt.Println(e1.str(), "->", e1.eval(env).display())
	e2 := &Symbol{"+"}
	fmt.Println(e2.str(), "->", e2.eval(env).display())
}

func test_apply() {
	env := sampleEnv()
	e1 := &Symbol{"a"}
	e2 := &Symbol{"b"}
	args := []Expr{e1, e2}
	e3 := &Apply{&Symbol{"+"}, args}
	fmt.Println(e3.str(), "->", e3.eval(env).display())
}

func test_if() {
	env := sampleEnv()
	e1 := &If{&Symbol{"t"}, &Symbol{"a"}, &Symbol{"b"}}
	fmt.Println(e1.str(), "->", e1.eval(env).display())
	e2 := &If{&Symbol{"f"}, &Symbol{"a"}, &Symbol{"b"}}
	fmt.Println(e2.str(), "->", e2.eval(env).display())
}

func test_parse() {
	env := sampleEnv()
	e1, _ := parseExpr("123")
	fmt.Println(e1.str(), "->", e1.eval(env).display())
	e2, _ := parseExpr("a")
	fmt.Println(e2.str(), "->", e2.eval(env).display())
	e3, _ := parseExpr("(+ 33 a)")
	fmt.Println(e3.str(), "->", e3.eval(env).display())
	e4, _ := parseExpr("(+ 33 (+ a b))")
	fmt.Println(e4.str(), "->", e4.eval(env).display())
}

type Value interface {
	display() string
	intValue() int
	boolValue() bool
	apply([]Value) Value
	str() string
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

func (v *VInteger) intValue() int {
	return v.val
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

func (v *VBoolean) intValue() int {
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

type VPrimitive struct {
	name      string
	primitive func([]Value) Value
	numArgs   int
}

func (v *VPrimitive) display() string {
	return fmt.Sprintf("#<PRIMITIVE %s>", v.name)
}

func (v *VPrimitive) intValue() int {
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

func primitivePlus(args []Value) Value {
	var result int
	for _, val := range args {
		result += val.intValue()
	}
	return &VInteger{result}
}

type Expr interface {
	eval(*Env) Value
	str() string
}

// do we use literal? or integer?

type Literal struct {
	val Value
}

func (e *Literal) eval(env *Env) Value {
	return e.val
}

func (e *Literal) str() string {
	return fmt.Sprintf("Literal[%s]", e.val.str())
}

type Symbol struct {
	name string
}

func (e *Symbol) eval(env *Env) Value {
	return env.lookup(e.name)
}

func (e *Symbol) str() string {
	return fmt.Sprintf("Symbol[%s]", e.name)
}

type If struct {
	condExpr Expr
	thenExpr Expr
	elseExpr Expr
}

func (e *If) eval(env *Env) Value {
	c := e.condExpr.eval(env)
	if c.boolValue() {
		return e.thenExpr.eval(env)
	} else {
		return e.elseExpr.eval(env)
	}
}

func (e *If) str() string {
	return fmt.Sprintf("If[%s %s %s]", e.condExpr.str(), e.thenExpr.str(), e.elseExpr.str())
}

type Apply struct {
	fn   Expr
	args []Expr
}

func (e *Apply) eval(env *Env) Value {
	f := e.fn.eval(env)
	args := make([]Value, len(e.args))
	for i := range args {
		args[i] = e.args[i].eval(env)
	}
	return f.apply(args)
}

func (e *Apply) str() string {
	strArgs := ""
	for _, item := range e.args {
		strArgs += " " + item.str()
	}
	return fmt.Sprintf("Apply[%s%s]", e.fn.str(), strArgs)
}

type Env struct {
	current  map[string]Value
	previous *Env
}

func (env *Env) lookup(name string) Value {
	val, ok := env.current[name]
	if !ok {
		panic("Boom!")
	}
	return val
}

func parseToken(token string, s string) (string, string) {
	r, _ := regexp.Compile(token)
	ss := strings.TrimSpace(s)
	match := r.FindStringIndex(ss)
	if len(match) == 0 {
		// no match
		return "", s
	} else {
		//fmt.Println("Token match", ss, match)
		return ss[:match[1]], ss[match[1]:]
	}
}

func parseLP(s string) (string, string) {
	//fmt.Println("Trying to parse as LP")
	return parseToken(`^\(`, s)
}

func parseRP(s string) (string, string) {
	//fmt.Println("Trying to parse as RP")
	return parseToken(`^\)`, s)
}

func parseSymbol(s string) (string, string) {
	//fmt.Println("Trying to parse as symbol")
	return parseToken(`^[^'()#\s]+`, s)
}

func parseInteger(s string) (string, string) {
	//fmt.Println("Trying to parse as integer")
	return parseToken(`^-?[0-9]+`, s)
}

func parseExprs(s string) ([]Expr, string) {
	result := make([]Expr, 0, 10)
	var rest string
	var expr Expr
	expr, rest = parseExpr(s)
	for expr != nil {
		result = append(result, expr)
		expr, rest = parseExpr(rest)
	}
	return result, rest
}

func parseExpr(s string) (Expr, string) {
	//fmt.Println("Trying to parse string", s)
	var result, rest string
	result, rest = parseInteger(s)
	if result != "" {
		num, _ := strconv.Atoi(result)
		return &Literal{&VInteger{num}}, rest
	}
	result, rest = parseSymbol(s)
	if result != "" {
		return &Symbol{result}, rest
	}
	result, rest = parseLP(s)
	if result != "" {
		var expr Expr
		var exprs []Expr
		expr, rest = parseExpr(rest)
		if expr != nil {
			exprs, rest = parseExprs(rest)
			result, rest = parseRP(rest)
			if result != "" {
				return &Apply{expr, exprs}, rest
			}
		}
	}
	return nil, s
}
