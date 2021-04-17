package main

import "fmt"
import "bufio"
import "os"
import "strings"
import "io"

func shell(eco *Ecosystem) {
	currentModule := "core"
	env, _ := eco.get(currentModule)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s> ", currentModule)
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				bail()
			}
			fmt.Println("IO ERROR - ", err.Error())
		}
		if strings.TrimSpace(text) == "" {
			continue
		}
		v, _, err := read(text)
		if err != nil {
			fmt.Println("READ ERROR -", err.Error())
			continue
		}
		// check if it's a declaration
		d, err := parseDef(v)
		if err != nil { 
			fmt.Println("PARSE ERROR -", err.Error())
			continue
		}
		if d != nil {
			if d.typ == DEF_FUNCTION { 
				env.update(d.name, &VFunction{d.params, d.body, env})
				fmt.Println(d.name)
				continue
			}
			if d.typ == DEF_VALUE {
				v, err := d.body.eval(env)
				if err != nil {
					fmt.Println("EVAL ERROR -", err.Error())
					continue
				}
				env.update(d.name, v)
				fmt.Println(d.name)
				continue
			}
			fmt.Println("DECLARE ERROR - unknow declaration type", d.typ)
			continue
		}
		// check if it's an expression
		e, err := parseExpr(v)
		if err != nil { 
			fmt.Println("PARSE ERROR -", err.Error())
			continue
		}
		///fmt.Println("expr =", e.str())
		v, err = e.eval(env)
		if err != nil {
			fmt.Println("EVAL ERROR -", err.Error())
			continue
		}
		fmt.Println(v.display())
	}
}

func bail() {
	fmt.Println("tada")
	os.Exit(0)
}

func initialize() *Ecosystem {
	eco := mkEcosystem()
	coreBindings := corePrimitives()
	coreBindings["true"] = &VBoolean{true}
	coreBindings["false"] = &VBoolean{false}
	eco.mkEnv("core", coreBindings)
	testBindings := map[string]Value{
		"a": &VInteger{99},
		"square": &VPrimitive{"square", func(args []Value) (Value, error) {
			if len(args) != 1 || !args[0].isNumber() {
				return nil, fmt.Errorf("argument to square missing or not int")
			}
			return &VInteger{args[0].intValue() * args[0].intValue()}, nil
		}},
	}
	eco.mkEnv("test", testBindings)
	shellBindings := map[string]Value{
		"quit": &VPrimitive{"quit", func(args []Value) (Value, error) {
			if len(args) > 0 {
				return nil, fmt.Errorf("too many arguments 0 to primitive quit")
			}
			bail()
			return nil, nil
		}},
	}
	eco.mkEnv("shell", shellBindings)
	return eco
}
