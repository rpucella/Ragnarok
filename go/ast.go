package main

import "fmt"
import "errors"

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

type LetRec struct {
	names []string
	params [][]string
	bodies []AST
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

func (e *LetRec) eval(env *Env) (Value, error) {
	if len(e.names) != len(e.params) || len(e.names) != len(e.bodies) {
		return nil, errors.New("malformed letrec (names, params, bodies)")
	}
	// create the environment that we'll share across the definitions
	// all names initially allocated #nil
	new_env := env.layer(e.names, nil)
	for i, name := range e.names {
		new_env.update(name, &VFunction{e.params[i], e.bodies[i], new_env})
	}
	return e.body.eval(new_env)
}

func (e *LetRec) str() string {
	return fmt.Sprintf("LetRec[? %s]", e.body)
}

