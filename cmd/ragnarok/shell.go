package main

import (
       "fmt"
       "bufio"
       "os"
       "strings"
       "io"
       "rpucella.net/ragnarok/internal/lisp"
)

var context = Context{"", "", nil}

func shell(eco *Ecosystem) {
	env := eco.mkEnv("*scratch*", map[string]lisp.Value{})
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
			if d.Type == lisp.DEF_FUNCTION { 
				env.Update(d.Name, lisp.NewVFunction(d.Params, d.Body, env))
				fmt.Println(d.Name)
				continue
			}
			if d.Type == lisp.DEF_VALUE {
				v, err := d.Body.Eval(env)
				if err != nil {
					fmt.Println("EVAL ERROR -", err.Error())
					continue
				}
				env.Update(d.Name, v)
				fmt.Println(d.Name)
				continue
			}
			fmt.Println("DECLARE ERROR - unknow declaration type", d.Type)
			continue
		}
		// check if it's an expression
		e, err := parseExpr(v)
		if err != nil { 
			fmt.Println("PARSE ERROR -", err.Error())
			continue
		}
		///fmt.Println("expr =", e.str())
		v, err = e.Eval(env)
		if err != nil {
			fmt.Println("EVAL ERROR -", err.Error())
			continue
		}
		if !v.IsNil() { 
			fmt.Println(v.Display())
		}
	}
}

func bail() {
	fmt.Println("tada")
	os.Exit(0)
}

func initialize() *Ecosystem {
	eco := mkEcosystem()
	coreBindings := corePrimitives()
	coreBindings["true"] = lisp.NewVBoolean(true)
	coreBindings["false"] = lisp.NewVBoolean(false)
	eco.mkEnv("core", coreBindings)
	testBindings := map[string]lisp.Value{
		"a": lisp.NewVInteger(99),
		"square": lisp.NewVPrimitive("square", func(args []lisp.Value) (lisp.Value, error) {
			if len(args) != 1 || !args[0].IsNumber() {
				return nil, fmt.Errorf("argument to square should be int")
			}
			return lisp.NewVInteger(args[0].IntValue() * args[0].IntValue()), nil
		}),
	}
	eco.mkEnv("test", testBindings)
	shellBindings := shellPrimitives()
	eco.mkEnv("shell", shellBindings)
	configBindings := map[string]lisp.Value{
		"lookup-path": lisp.NewVReference(lisp.NewVCons(lisp.NewVSymbol("shell"), lisp.NewVCons(lisp.NewVSymbol("core"), &lisp.VEmpty{}))),
		"editor": lisp.NewVReference(lisp.NewVString("emacs")),
	}
	eco.mkEnv("config", configBindings)
	return eco
}

func showModules(env *lisp.Env) {
	modulesFn, err := env.Lookup("shell", "modules")
	if err != nil {
		return
	}
	v, err := modulesFn.Apply([]lisp.Value{})
	if err != nil {
		return}
	fmt.Println("Modules", v.Display())
}
