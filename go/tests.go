package main

import "fmt"

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

func primitiveAdd(args []Value) (Value, error) {
	var result int
	for _, val := range args {
		result += val.intValue()
	}
	return &VInteger{result}, nil
}

func primitiveMult(args []Value) (Value, error) {
	var result int = 1
	for _, val := range args {
		result *= val.intValue()
	}
	return &VInteger{result}, nil
}

func sampleEnv() *Env {
	current := map[string]Value{
		"a": &VInteger{10},
		"b": &VInteger{20},
		"+": &VPrimitive{"+", primitiveAdd},
		"*": &VPrimitive{"*", primitiveMult},
		"t": &VBoolean{true},
		"f": &VBoolean{false},
	}
	env := &Env{bindings: current}
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
	var vp Value = &VPrimitive{"+", primitiveAdd}
	var args []Value = []Value{v1, v2, v3}
	vr, _ := vp.apply(args)
	fmt.Println(vp.str(), "->", vr.intValue())
}

func evalDisplay(e AST, env *Env) string {
	v, _ := e.eval(env)
	return v.display()
}

func test_literal() {
	v1 := &VInteger{10}
	e1 := &Literal{v1}
	fmt.Println(e1.str(), "->", evalDisplay(e1, nil))
	v2 := &VBoolean{true}
	e2 := &Literal{v2}
	fmt.Println(e2.str(), "->", evalDisplay(e2, nil))
}

func test_lookup() {
	env := sampleEnv()
	e1 := &Id{"a"}
	fmt.Println(e1.str(), "->", evalDisplay(e1, env))
	e2 := &Id{"+"}
	fmt.Println(e2.str(), "->", evalDisplay(e2, env))
}

func test_apply() {
	env := sampleEnv()
	e1 := &Id{"a"}
	e2 := &Id{"b"}
	args := []AST{e1, e2}
	e3 := &Apply{&Id{"+"}, args}
	fmt.Println(e3.str(), "->", evalDisplay(e3, env))
}

func test_if() {
	env := sampleEnv()
	e1 := &If{&Id{"t"}, &Id{"a"}, &Id{"b"}}
	fmt.Println(e1.str(), "->", evalDisplay(e1, env))
	e2 := &If{&Id{"f"}, &Id{"a"}, &Id{"b"}}
	fmt.Println(e2.str(), "->", evalDisplay(e2, env))
}

func test_read() {
	v1, _, _ := read("123")
	fmt.Println(v1.str(), "->", v1.display())
	v2, _, _ := read("a")
	fmt.Println(v2.str(), "->", v2.display())
	v6, _, _ := read("+")
	fmt.Println(v6.str(), "->", v6.display())
	v3, _, _ := read("(+ 33 a)")
	fmt.Println(v3.str(), "->", v3.display())
	v4, _, _ := read("(+ 33 (+ a b))")
	fmt.Println(v4.str(), "->", v4.display())
	v5, _, _ := read("(this is my life)")
	fmt.Println(v5.str(), "->", v5.display())
}

func test_lists() {
	var v Value = &VEmpty{}
	v = &VCons{head: &VInteger{33}, tail: v}
	v = &VCons{head: &VInteger{66}, tail: v}
	v = &VCons{head: &VInteger{99}, tail: v}
	fmt.Println(v.str(), "->", v.display())
}
