package main

import "fmt"

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
	str() string
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

type Function struct {
	params []string
	body AST
}

func (e *Literal) eval(env *Env) (Value, error) {
	return e.val, nil
}

func (e *Literal) str() string {
	return fmt.Sprintf("Literal[%s]", e.val.str())
}

func (e *Id) eval(env *Env) (Value, error) {
	return env.lookup(e.name)
}

func (e *Id) str() string {
	return fmt.Sprintf("Id[%s]", e.name)
}

func (e *If) eval(env *Env) (Value, error) {
	c, err := e.cnd.eval(env)
	if err != nil {
		return nil, err
	}
	if c.isTrue() {
		return e.thn.eval(env)
	} else {
		return e.els.eval(env)
	}
}

func (e *If) str() string {
	return fmt.Sprintf("If[%s %s %s]", e.cnd.str(), e.thn.str(), e.els.str())
}

func (e *Apply) eval(env *Env) (Value, error) {
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
	return f.apply(args)
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

func (e *Quote) str() string {
	return fmt.Sprintf("Quote[%s]", e.val.str())
}

func (e *Function) eval(env *Env) (Value, error) {
	return &VFunction{e.params, e.body, env}, nil
}

func (e *Function) str() string {
	return fmt.Sprintf("Function[?]")
}

