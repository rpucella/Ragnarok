package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"rpucella.net/ragnarok/internal/lisp"
	"rpucella.net/ragnarok/internal/parser"
	"rpucella.net/ragnarok/internal/reader"
	"strings"
)

// A context contains anything interesting to the execution

// Right now, it's a global variable

// maybe we want to make it available via the ecosystem (thus during evaluation of forms)
// and passing it to primitives (so that they can use it to access, well, the context)

type Context struct {
	currentModule     string
	nextCurrentModule string // to switch modules, set nextCurrentModule != nil
	ecosystem         Ecosystem
	report            func(string)
}

func shell(eco Ecosystem) {
	env := lisp.NewEnv(map[string]lisp.Value{}, nil, eco.modules)
	eco.addShell("*scratch*", env)
	context := &Context{"*scratch*", "", eco, func(str string) { fmt.Println(";;", str) }}
	stdInReader := bufio.NewReader(os.Stdin)
	showModules(env)
	for {
		if context.nextCurrentModule != "" {
			current := context.currentModule
			context.currentModule = context.nextCurrentModule
			context.nextCurrentModule = ""
			new_env, err := eco.get(context.currentModule)
			if err != nil {
				// reset the module names
				context.currentModule = current
				fmt.Println("ERROR -", err.Error())
			} else {
				env = new_env
			}
		}
		fmt.Printf("%s | ", context.currentModule)
		text, err := stdInReader.ReadString('\n')
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
		v, _, err := reader.Read(text)
		if err != nil {
			fmt.Println("READ ERROR -", err.Error())
			continue
		}
		// check if it's a declaration
		d, err := parser.ParseDef(v)
		if err != nil {
			fmt.Println("PARSE ERROR -", err.Error())
			continue
		}
		if d != nil {
			// we have a declaration
			if d.Type == lisp.DEF_FUNCTION {
				env.Update(d.Name, lisp.NewVFunction(d.Params, d.Body, env))
				fmt.Println(";;", d.Name)
				continue
			}
			if d.Type == lisp.DEF_VALUE {
				v, err := d.Body.Eval(env, context)
				if err != nil {
					fmt.Println("EVAL ERROR -", err.Error())
					continue
				}
				env.Update(d.Name, v)
				fmt.Println(";;", d.Name)
				continue
			}
			fmt.Println("DECLARE ERROR - unknow declaration type", d.Type)
			continue
		}
		// check if it's an expression
		e, err := parser.ParseExpr(v)
		if err != nil {
			fmt.Println("PARSE ERROR -", err.Error())
			continue
		}
		///fmt.Println("expr =", e.str())
		v, err = e.Eval(env, context)
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
	fmt.Println("")    // tada, arrivederci, auf wiedersehen, hasta la vista, goodbye, au revoir
	os.Exit(0)
}

func showModules(env *lisp.Env) {
	modulesFn, err := env.Lookup("shell", "modules")
	if err != nil {
		return
	}
	v, err := modulesFn.Apply([]lisp.Value{}, nil)
	if err != nil {
		return
	}
	fmt.Println("Modules", v.Display())
}
