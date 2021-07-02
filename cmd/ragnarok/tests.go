package main

import (
	"fmt"
	"rpucella.net/ragnarok/internal/lisp"
	"rpucella.net/ragnarok/internal/reader"
)

func test() {

	test_value_10()
	test_value_plus()
	test_literal()
	test_lookup()
	test_apply()
	test_if()
	test_lists()
	test_read()
}

func primitiveAdd(args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
	var result int
	for _, val := range args {
		result += val.IntValue()
	}
	return lisp.NewVInteger(result), nil
}

func primitiveMult(args []lisp.Value, ctxt interface{}) (lisp.Value, error) {
	var result int = 1
	for _, val := range args {
		result *= val.IntValue()
	}
	return lisp.NewVInteger(result), nil
}

func sampleEnv() *lisp.Env {
	current := map[string]lisp.Value{
		"a": lisp.NewVInteger(10),
		"b": lisp.NewVInteger(20),
		"+": lisp.NewVPrimitive("+", primitiveAdd),
		"*": lisp.NewVPrimitive("*", primitiveMult),
		"t": lisp.NewVBoolean(true),
		"f": lisp.NewVBoolean(false),
	}
	env := lisp.NewEnv(current, nil, nil)
	return env
}

func test_value_10() {
	var v1 lisp.Value = lisp.NewVInteger(10)
	fmt.Println(v1.Str(), "->", v1.IntValue())
}

func test_value_plus() {
	var v1 lisp.Value = lisp.NewVInteger(10)
	var v2 lisp.Value = lisp.NewVInteger(20)
	var v3 lisp.Value = lisp.NewVInteger(30)
	var vp lisp.Value = lisp.NewVPrimitive("+", primitiveAdd)
	var args []lisp.Value = []lisp.Value{v1, v2, v3}
	vr, _ := vp.Apply(args, nil)
	fmt.Println(vp.Str(), "->", vr.IntValue())
}

func evalDisplay(e lisp.AST, env *lisp.Env) string {
	v, _ := e.Eval(env, nil)
	return v.Display()
}

func test_literal() {
	v1 := lisp.NewVInteger(10)
	e1 := lisp.NewLiteral(v1)
	fmt.Println(e1.Str(), "->", evalDisplay(e1, nil))
	v2 := lisp.NewVBoolean(true)
	e2 := lisp.NewLiteral(v2)
	fmt.Println(e2.Str(), "->", evalDisplay(e2, nil))
}

func test_lookup() {
	env := sampleEnv()
	e1 := lisp.NewId("a")
	fmt.Println(e1.Str(), "->", evalDisplay(e1, env))
	e2 := lisp.NewId("+")
	fmt.Println(e2.Str(), "->", evalDisplay(e2, env))
}

func test_apply() {
	env := sampleEnv()
	e1 := lisp.NewId("a")
	e2 := lisp.NewId("b")
	args := []lisp.AST{e1, e2}
	e3 := lisp.NewApply(lisp.NewId("+"), args)
	fmt.Println(e3.Str(), "->", evalDisplay(e3, env))
}

func test_if() {
	env := sampleEnv()
	e1 := lisp.NewIf(lisp.NewId("t"), lisp.NewId("a"), lisp.NewId("b"))
	fmt.Println(e1.Str(), "->", evalDisplay(e1, env))
	e2 := lisp.NewIf(lisp.NewId("f"), lisp.NewId("a"), lisp.NewId("b"))
	fmt.Println(e2.Str(), "->", evalDisplay(e2, env))
}

func test_read() {
	v1, _, _ := reader.Read("123")
	fmt.Println(v1.Str(), "->", v1.Display())
	v2, _, _ := reader.Read("a")
	fmt.Println(v2.Str(), "->", v2.Display())
	v6, _, _ := reader.Read("+")
	fmt.Println(v6.Str(), "->", v6.Display())
	v3, _, _ := reader.Read("(+ 33 a)")
	fmt.Println(v3.Str(), "->", v3.Display())
	v4, _, _ := reader.Read("(+ 33 (+ a b))")
	fmt.Println(v4.Str(), "->", v4.Display())
	v5, _, _ := reader.Read("(this is my life)")
	fmt.Println(v5.Str(), "->", v5.Display())
}

func test_lists() {
	var v lisp.Value = &lisp.VEmpty{}
	v = lisp.NewVCons(lisp.NewVInteger(33), v)
	v = lisp.NewVCons(lisp.NewVInteger(66), v)
	v = lisp.NewVCons(lisp.NewVInteger(99), v)
	fmt.Println(v.Str(), "->", v.Display())
}
