package shell

import (
	"fmt"
	"github.com/peterh/liner"
	"io"
	"os"
	"rpucella.net/ragnarok/internal/evaluator"
	"rpucella.net/ragnarok/internal/parser"
	"rpucella.net/ragnarok/internal/reader"
	"rpucella.net/ragnarok/internal/value"
	"strings"
)

// A context contains anything interesting to the execution

// Right now, it's a global variable

// maybe we want to make it available via the ecosystem (thus during evaluation of forms)
// and passing it to primitives (so that they can use it to access, well, the context)

type Context struct {
	homeModule        string
	currentModule     string
	nextCurrentModule string // to switch modules, set nextCurrentModule != nil
	ecosystem         Ecosystem
	currentEnv        *evaluator.Env
	report            func(string)
	bail              func()
	readAll           func(string, *Context) error
}

func Shell(eco Ecosystem) {
	name := "*1*"
	eco.AddShell(name, map[string]value.Value{})
	env, _ := eco.get(name)
	line := liner.NewLiner()
	defer line.Close()
	report := func(str string) {
		fmt.Println(";;", str)
	}
	bail := func() {
		line.Close()
		fmt.Println("") // tada, arrivederci, auf wiedersehen, hasta la vista, goodbye, au revoir
		os.Exit(0)
	}
	readAll := func(str string, context *Context) error {
		vLines := []value.Value{}
		// do something better
		curr := str
		for curr != "" {
			vLine, rest, err := reader.Read(curr)
			if err != nil {
				return err
			}
			if vLine == nil {
				if rest == "" {
					// we're done...
					break
				}
				// we have an incomplete term, but we're done reading
				// so we must fail
				return fmt.Errorf("READ ERROR - incomplete form")
			}
			//fmt.Println("Got a form:", vLine.Display())
			vLines = append(vLines, vLine)
			curr = rest
		}
		// all good, so process all inputs
		for _, vLine := range vLines {
			_, err := processInput(vLine, context.currentEnv, context)
			if err != nil {
				return err
			}
		}
		return nil
	}
	context := &Context{name, name, "", eco, env, report, bail, readAll}
	//stdInReader := bufio.NewReader(os.Stdin)
	//showModules(env, context)
REPL:
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
				context.currentEnv = env
			}
		}
		// fmt.Printf("\n%s | ", context.currentModule)
		// text, err := stdInReader.ReadString('\n')
		// if err != nil {
		// 	if err == io.EOF {
		// 		fmt.Println()
		// 		bail()
		// 	}
		// 	fmt.Println("IO ERROR - ", err.Error())
		// }
		vText, err := readInput(line, context)
		if err != nil {
			fmt.Println(err.Error())
			continue REPL
		}
		v, err := processInput(vText, env, context)
		if err != nil {
			fmt.Println(err.Error())
			continue REPL
		}
		if v != nil && !value.IsNil(v) {
			fmt.Println(v.Str())
		}
	}
}

func readInput(line *liner.State, context *Context) (value.Value, error) {
	currText := ""
	var vText value.Value = nil
	prompt := fmt.Sprintf("%s | ", context.currentModule)
	for vText == nil {
		text, err := line.Prompt(prompt)
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return nil, fmt.Errorf("Use (quit) to bail out.")
			}
			return nil, fmt.Errorf("IO ERROR - %s", err)
		}
		if strings.TrimSpace(text) == "" {
			continue
		}
		line.AppendHistory(text)
		currText = currText + "\n" + text
		vText, _, err = reader.Read(currText)
		if err != nil {
			return nil, fmt.Errorf("READ ERROR - %s", err)
		}
		prompt = fmt.Sprintf("%*s | ", len(context.currentModule), " ")
	}
	return vText, nil
}

func processInput(v value.Value, env *evaluator.Env, context *Context) (value.Value, error) {
	// check if it's a declaration
	d, err := parser.ParseDef(v)
	if err != nil {
		return nil, fmt.Errorf("PARSE ERROR - %s", err)
	}
	if d != nil {
		// we have a declaration
		if d.Type == evaluator.DEF_FUNCTION {
			f, err := parser.MakeFunction(d.Params, d.Body).Eval(env, context)
			if err != nil {
				fmt.Println("EVAL ERROR -", err.Error())
			}
			env.Update(d.Name, f)
			fmt.Println(";;", d.Name)
			return nil, nil
		}
		if d.Type == evaluator.DEF_VALUE {
			v, err := d.Body.Eval(env, context)
			if err != nil {
				fmt.Println("EVAL ERROR -", err.Error())
			}
			env.Update(d.Name, v)
			fmt.Println(";;", d.Name)
			return nil, nil
		}
		return nil, fmt.Errorf("DECLARE ERROR - unknow declaration type %d", d.Type)
	}
	// check if it's an expression
	e, err := parser.ParseExpr(v)
	if err != nil {
		return nil, fmt.Errorf("PARSE ERROR - %s", err)
	}
	///fmt.Println("expr =", e.str())
	v, err = e.Eval(env, context)
	if err != nil {
		return nil, fmt.Errorf("EVAL ERROR - %s", err)
	}
	return v, nil
}

func showModules(env *evaluator.Env, ctxt *Context) {
	modulesFn, err := env.Lookup("shell", "modules")
	if err != nil {
		fmt.Println("Problem in showModules():", err)
		return
	}
	v, err := modulesFn.Apply([]value.Value{}, ctxt)
	if err != nil {
		fmt.Println("Problem in showModules():", err)
		return
	}
	fmt.Println("Modules", v.Display())
}

//func proces
