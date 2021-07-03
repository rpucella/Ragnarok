package main

import (
	"fmt"
	"rpucella.net/ragnarok/internal/evaluator"
	"rpucella.net/ragnarok/internal/value"
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

func primitiveAdd(args []value.Value, ctxt interface{}) (value.Value, error) {
	var result int
	for _, val := range args {
		result += val.GetInt()
	}
	return value.NewVInteger(result), nil
}

func primitiveMult(args []value.Value, ctxt interface{}) (value.Value, error) {
	var result int = 1
	for _, val := range args {
		result *= val.GetInt()
	}
	return value.NewVInteger(result), nil
}

func sampleEnv() *evaluator.Env {
	current := map[string]value.Value{
		"a": value.NewVInteger(10),
		"b": value.NewVInteger(20),
		"+": value.NewVPrimitive("+", primitiveAdd),
		"*": value.NewVPrimitive("*", primitiveMult),
		"t": value.NewVBoolean(true),
		"f": value.NewVBoolean(false),
	}
	env := evaluator.NewEnv(current, nil, nil)
	return env
}

func test_value_10() {
	var v1 value.Value = value.NewVInteger(10)
	fmt.Println(v1.Str(), "->", v1.GetInt())
}

func test_value_plus() {
	var v1 value.Value = value.NewVInteger(10)
	var v2 value.Value = value.NewVInteger(20)
	var v3 value.Value = value.NewVInteger(30)
	var vp value.Value = value.NewVPrimitive("+", primitiveAdd)
	var args []value.Value = []value.Value{v1, v2, v3}
	vr, _ := vp.Apply(args, nil)
	fmt.Println(vp.Str(), "->", vr.GetInt())
}

func evalDisplay(e evaluator.AST, env *evaluator.Env) string {
	v, _ := e.Eval(env, nil)
	return v.Display()
}

func test_literal() {
	v1 := value.NewVInteger(10)
	e1 := evaluator.NewLiteral(v1)
	fmt.Println(e1.Str(), "->", evalDisplay(e1, nil))
	v2 := value.NewVBoolean(true)
	e2 := evaluator.NewLiteral(v2)
	fmt.Println(e2.Str(), "->", evalDisplay(e2, nil))
}

func test_lookup() {
	env := sampleEnv()
	e1 := evaluator.NewId("a")
	fmt.Println(e1.Str(), "->", evalDisplay(e1, env))
	e2 := evaluator.NewId("+")
	fmt.Println(e2.Str(), "->", evalDisplay(e2, env))
}

func test_apply() {
	env := sampleEnv()
	e1 := evaluator.NewId("a")
	e2 := evaluator.NewId("b")
	args := []evaluator.AST{e1, e2}
	e3 := evaluator.NewApply(evaluator.NewId("+"), args)
	fmt.Println(e3.Str(), "->", evalDisplay(e3, env))
}

func test_if() {
	env := sampleEnv()
	e1 := evaluator.NewIf(evaluator.NewId("t"), evaluator.NewId("a"), evaluator.NewId("b"))
	fmt.Println(e1.Str(), "->", evalDisplay(e1, env))
	e2 := evaluator.NewIf(evaluator.NewId("f"), evaluator.NewId("a"), evaluator.NewId("b"))
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
	var v value.Value = &value.VEmpty{}
	v = value.NewVCons(value.NewVInteger(33), v)
	v = value.NewVCons(value.NewVInteger(66), v)
	v = value.NewVCons(value.NewVInteger(99), v)
	fmt.Println(v.Str(), "->", v.Display())
}
