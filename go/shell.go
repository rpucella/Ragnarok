package main

import "fmt"
import "bufio"
import "os"
import "strings"

func shell() {
	bindings := primitivesBindings()
	bindings["true"] = &VBoolean{true}
	bindings["false"] = &VBoolean{false}
	reader := bufio.NewReader(os.Stdin)
	env := makeEnv(bindings)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(text) == "" {
			continue
		}
		v, _, err := read(text)
		if err != nil {
			fmt.Println("READ -", err.Error())
			continue
		}
		// check if it's a declaration
		d, err := parseDef(v)
		if err != nil { 
			fmt.Println("PARSE -", err.Error())
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
					fmt.Println("EVAL -", err.Error())
					continue
				}
				env.update(d.name, v)
				fmt.Println(d.name)
				continue
			}
			fmt.Println("DECLARE - unknow declaration type", d.typ)
			continue
		}
		// check if it's an expression
		e, err := parseExpr(v)
		if err != nil { 
			fmt.Println("PARSE -", err.Error())
			continue
		}
		///fmt.Println("expr =", e.str())
		v, err = e.eval(env)
		if err != nil {
			fmt.Println("EVAL -", err.Error())
			continue
		}
		fmt.Println(v.display())
	}
}
