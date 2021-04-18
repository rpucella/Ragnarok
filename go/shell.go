package main

import "fmt"
import "bufio"
import "os"
import "strings"
import "io"

var context = Context{"", "", nil}

func shell(eco *Ecosystem) {
	env := eco.mkEnv("*scratch*", map[string]Value{})
	context.currentModule = "*scratch*"
	context.nextCurrentModule = "*scratch*"
	context.ecosystem = eco
	reader := bufio.NewReader(os.Stdin)
	showModules(env)
	for {
		if context.nextCurrentModule != context.currentModule {
			current := context.currentModule
			context.currentModule = context.nextCurrentModule
			new_env, err := eco.get(context.currentModule)
			if err != nil {
				// reset the module names
				context.currentModule = current
				context.nextCurrentModule = current
				fmt.Println("ERROR -", err.Error())
			} else {
				env = new_env
			}
		}
		fmt.Printf("%s> ", context.currentModule)
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
				return nil, fmt.Errorf("argument to square should be int")
			}
			return &VInteger{args[0].intValue() * args[0].intValue()}, nil
		}},
	}
	eco.mkEnv("test", testBindings)
	shellBindings := shellPrimitives()
	eco.mkEnv("shell", shellBindings)
	configBindings := map[string]Value{
		"lookup-path": &VReference{&VCons{head: &VSymbol{"shell"}, tail: &VCons{head: &VSymbol{"core"}, tail: &VEmpty{}}}},
		"editor": &VReference{&VString{"emacs"}},
	}
	eco.mkEnv("config", configBindings)
	return eco
}

func showModules(env *Env) {
	modulesFn, err := env.lookup("shell", "modules")
	if err != nil {
		return
	}
	v, err := modulesFn.apply([]Value{})
	if err != nil {
		return}
	fmt.Println("Modules", v.display())
}
