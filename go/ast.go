package main

import "fmt"

type AST interface {
	eval(*Env) Value
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

func (e *Literal) eval(env *Env) Value {
	return e.val
}

func (e *Literal) str() string {
	return fmt.Sprintf("Literal[%s]", e.val.str())
}

func (e *Id) eval(env *Env) Value {
	return env.lookup(e.name)
}

func (e *Id) str() string {
	return fmt.Sprintf("Id[%s]", e.name)
}

func (e *If) eval(env *Env) Value {
	c := e.cnd.eval(env)
	if c.boolValue() {
		return e.thn.eval(env)
	} else {
		return e.els.eval(env)
	}
}

func (e *If) str() string {
	return fmt.Sprintf("If[%s %s %s]", e.cnd.str(), e.thn.str(), e.els.str())
}

type Apply struct {
	fn   AST
	args []AST
}

func (e *Apply) eval(env *Env) Value {
	f := e.fn.eval(env)
	args := make([]Value, len(e.args))
	for i := range args {
		args[i] = e.args[i].eval(env)
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
