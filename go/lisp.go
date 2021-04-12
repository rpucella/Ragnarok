package main

import "fmt"
import "bufio"
import "os"

func main() {

	test_value_10()
	test_value_plus()

	test_literal()
	test_lookup()
	test_apply()
	test_if()

	test_lists()

	test_read()

	fmt.Println("------------------------------------------------------------")

	env := sampleEnv()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		v, _ := read(text)
		if v == nil {
			fmt.Println("Cannot read", text)
			continue
		}
		e := parseExpr(v)
		if e == nil {
			fmt.Println("Cannot parse", v.str())
			continue
		}
		fmt.Println(e.eval(env).display())
	}
}

func sampleEnv() *Env {
	current := map[string]Value{
		"a": &VInteger{10},
		"b": &VInteger{20},
		"+": &VPrimitive{"+", primitiveAdd, 0},
		"*": &VPrimitive{"*", primitiveMult, 0},
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
	var vp Value = &VPrimitive{"+", primitiveAdd, 0}
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
	e1 := &Id{"a"}
	fmt.Println(e1.str(), "->", e1.eval(env).display())
	e2 := &Id{"+"}
	fmt.Println(e2.str(), "->", e2.eval(env).display())
}

func test_apply() {
	env := sampleEnv()
	e1 := &Id{"a"}
	e2 := &Id{"b"}
	args := []AST{e1, e2}
	e3 := &Apply{&Id{"+"}, args}
	fmt.Println(e3.str(), "->", e3.eval(env).display())
}

func test_if() {
	env := sampleEnv()
	e1 := &If{&Id{"t"}, &Id{"a"}, &Id{"b"}}
	fmt.Println(e1.str(), "->", e1.eval(env).display())
	e2 := &If{&Id{"f"}, &Id{"a"}, &Id{"b"}}
	fmt.Println(e2.str(), "->", e2.eval(env).display())
}

func test_read() {
	v1, _ := read("123")
	fmt.Println(v1.str(), "->", v1.display())
	v2, _ := read("a")
	fmt.Println(v2.str(), "->", v2.display())
	v3, _ := read("(+ 33 a)")
	fmt.Println(v3.str(), "->", v3.display())
	v4, _ := read("(+ 33 (+ a b))")
	fmt.Println(v4.str(), "->", v4.display())
	v5, _ := read("(this is my life)")
	fmt.Println(v5.str(), "->", v5.display())
}

func test_lists() {
	var v Value = &VEmpty{}
	v = &VCons{head: &VInteger{33}, tail: v}
	v = &VCons{head: &VInteger{66}, tail: v}
	v = &VCons{head: &VInteger{99}, tail: v}
	fmt.Println(v.str(), "->", v.display())
}
