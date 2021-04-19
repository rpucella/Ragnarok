package main

import "fmt"
import "errors"
import "strings"

const DEF_VALUE = 0
const DEF_FUNCTION = 1

type Def struct {
	name string
	typ int
	params []string
	body AST
}

type AST interface {
	eval(*Env) (Value, error)
	evalPartial(*Env) (*PartialResult, error)
	str() string
}

type PartialResult struct {
	exp AST
	env *Env
	val Value  // val is null when the result is still partial
}

type Literal struct {
	val Value
}

type Id struct {
	name string
}

type If struct {
	cnd AST
	thn AST
	els AST
}

type Apply struct {
	fn   AST
	args []AST
}

type Quote struct {
	val Value
}

type LetRec struct {
	names []string
	params [][]string
	bodies []AST
	body AST
}

func defaultEvalPartial(e AST, env *Env) (*PartialResult, error) {
        // Partial evaluation
        // Sometimes return an expression to evaluate next along 
        // with an environment for evaluation.
        // val is null when the result is in fact a value.
	
	v, err := e.eval(env)
	if err != nil {
		return nil, err
	}
	return &PartialResult{nil, nil, v}, nil
}

func defaultEval(e AST, env *Env) (Value, error) {
	// evaluation with tail call optimization
	var currExp AST = e
	currEnv := env
	for {
		partial, err := currExp.evalPartial(currEnv)
		if err != nil {
			return nil, err
		}
		if partial.val != nil {
			return partial.val, nil
		}
		currExp = partial.exp
		currEnv = partial.env
	}
}

func (e *Literal) eval(env *Env) (Value, error) {
	return e.val, nil
}

func (e *Literal) evalPartial(env *Env) (*PartialResult, error) {
	return defaultEvalPartial(e, env)
}

func (e *Literal) str() string {
	return fmt.Sprintf("Literal[%s]", e.val.str())
}

func (e *Id) eval(env *Env) (Value, error) {
	return env.find(e.name)
}

func (e *Id) evalPartial(env *Env) (*PartialResult, error) {
	return defaultEvalPartial(e, env)
}

func (e *Id) str() string {
	return fmt.Sprintf("Id[%s]", e.name)
}

func (e *If) eval(env *Env) (Value, error) {
	return defaultEval(e, env)
}

func (e *If) evalPartial(env *Env) (*PartialResult, error) { 
	c, err := e.cnd.eval(env)
	if err != nil {
		return nil, err
	}
	if c.isTrue() {
		return &PartialResult{e.thn, env, nil}, nil
	} else {
		return &PartialResult{e.els, env, nil}, nil
	}
}

func (e *If) str() string {
	return fmt.Sprintf("If[%s %s %s]", e.cnd.str(), e.thn.str(), e.els.str())
}

func (e *Apply) eval(env *Env) (Value, error) {
	return defaultEval(e, env)
}

func (e *Apply) evalPartial(env *Env) (*PartialResult, error) {
	f, err := e.fn.eval(env)
	if err != nil {
		return nil, err
	}
	args := make([]Value, len(e.args))
	for i := range args {
		args[i], err = e.args[i].eval(env)
		if err != nil {
			return nil, err
		}
	}
	if ff, ok := f.(*VFunction); ok {
		if len(ff.params) != len(args) {
			return nil, fmt.Errorf("Wrong number of arguments to application to %s", ff.str())
		}
		newEnv := ff.env.layer(ff.params, args)
		return &PartialResult{ff.body, newEnv, nil}, nil
	}
	v, err := f.apply(args)
	if err != nil {
		return nil, err
	}
	return &PartialResult{nil, nil, v}, nil
}

func (e *Apply) str() string {
	strArgs := ""
	for _, item := range e.args {
		strArgs += " " + item.str()
	}
	return fmt.Sprintf("Apply[%s%s]", e.fn.str(), strArgs)
}

func (e *Quote) eval(env *Env) (Value, error) {
	return e.val, nil
}

func (e *Quote) evalPartial(env *Env) (*PartialResult, error) {
	return defaultEvalPartial(e, env)
}

func (e *Quote) str() string {
	return fmt.Sprintf("Quote[%s]", e.val.str())
}

func (e *LetRec) eval(env *Env) (Value, error) {
	return defaultEval(e, env)
}

func (e *LetRec) evalPartial(env *Env) (*PartialResult, error) {
	if len(e.names) != len(e.params) || len(e.names) != len(e.bodies) {
		return nil, errors.New("malformed letrec (names, params, bodies)")
	}
	// create the environment that we'll share across the definitions
	// all names initially allocated #nil
	newEnv := env.layer(e.names, nil)
	for i, name := range e.names {
		newEnv.update(name, &VFunction{e.params[i], e.bodies[i], newEnv})
	}
	return &PartialResult{e.body, newEnv, nil}, nil
}

func (e *LetRec) str() string {
	bindings := make([]string, len(e.names))
	for i := range e.names {
		params := strings.Join(e.params[i], " ")
		bindings[i] = fmt.Sprintf("[%s [%s] %s]", e.names[i], params, e.bodies[i].str())
	}
	return fmt.Sprintf("LetRec[%s %s]", strings.Join(bindings, " "), e.body.str())
}
