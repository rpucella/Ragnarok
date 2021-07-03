package evaluator

import (
	"fmt"
	"errors"
	"strings"
	"rpucella.net/ragnarok/internal/value"
)

const DEF_VALUE = 0
const DEF_FUNCTION = 1

type Def struct {
	Name   string
	Type   int
	Params []string
	Body   AST
}

type AST interface {
	Eval(*Env, interface{}) (value.Value, error)
	evalPartial(*Env, interface{}) (*partialResult, error)
	Str() string
}

type partialResult struct {
	exp AST
	env *Env
	val value.Value // val is null when the result is still partial
}

type Literal struct {
	val value.Value
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
	val value.Value
}

type LetRec struct {
	names  []string
	params [][]string
	bodies []AST
	body   AST
}

func NewDef(name string, typ int, params []string, body AST) *Def {
	return &Def{name, typ, params, body}
}

func NewId(name string) *Id {
	return &Id{name}
}

func NewLiteral(val value.Value) *Literal {
	return &Literal{val}
}

func NewQuote(val value.Value) *Quote {
	return &Quote{val}
}

func NewIf(cnd AST, thn AST, els AST) *If {
	return &If{cnd, thn, els}
}

func NewLetRec(names []string, params [][]string, bodies []AST, body AST) *LetRec {
	return &LetRec{names, params, bodies, body}
}

func NewApply(fn AST, args []AST) *Apply {
	return &Apply{fn, args}
}

func defaultEvalPartial(e AST, env *Env, ctxt interface{}) (*partialResult, error) {
	// Partial evaluation
	// Sometimes return an expression to evaluate next along
	// with an environment for evaluation.
	// val is null when the result is in fact a value.

	v, err := e.Eval(env, ctxt)
	if err != nil {
		return nil, err
	}
	return &partialResult{nil, nil, v}, nil
}

func defaultEval(e AST, env *Env, ctxt interface{}) (value.Value, error) {
	// evaluation with tail call optimization
	var currExp AST = e
	currEnv := env
	for {
		partial, err := currExp.evalPartial(currEnv, ctxt)
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

func (e *Literal) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return e.val, nil
}

func (e *Literal) evalPartial(env *Env, ctxt interface{}) (*partialResult, error) {
	return defaultEvalPartial(e, env, ctxt)
}

func (e *Literal) Str() string {
	return fmt.Sprintf("Literal[%s]", e.val.Str())
}

func (e *Id) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return env.find(e.name)
}

func (e *Id) evalPartial(env *Env, ctxt interface{}) (*partialResult, error) {
	return defaultEvalPartial(e, env, ctxt)
}

func (e *Id) Str() string {
	return fmt.Sprintf("Id[%s]", e.name)
}

func (e *If) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return defaultEval(e, env, ctxt)
}

func (e *If) evalPartial(env *Env, ctxt interface{}) (*partialResult, error) {
	c, err := e.cnd.Eval(env, ctxt)
	if err != nil {
		return nil, err
	}
	if value.IsTrue(c) {
		return &partialResult{e.thn, env, nil}, nil
	} else {
		return &partialResult{e.els, env, nil}, nil
	}
}

func (e *If) Str() string {
	return fmt.Sprintf("If[%s %s %s]", e.cnd.Str(), e.thn.Str(), e.els.Str())
}

func (e *Apply) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return defaultEval(e, env, ctxt)
}

func (e *Apply) evalPartial(env *Env, ctxt interface{}) (*partialResult, error) {
	f, err := e.fn.Eval(env, ctxt)
	if err != nil {
		return nil, err
	}
	args := make([]value.Value, len(e.args))
	for i := range args {
		args[i], err = e.args[i].Eval(env, ctxt)
		if err != nil {
			return nil, err
		}
	}
	// // this doesn't work if the VFunction doesn't hold the AST of the body!???
	// if ff, ok := f.(*VFunction); ok {
	// 	if len(ff.params) != len(args) {
	// 		return nil, fmt.Errorf("Wrong number of arguments to application to %s", ff.Str())
	// 	}
	// 	newEnv := ff.env.Layer(ff.params, args)
	// 	return &partialResult{ff.body, newEnv, nil}, nil
	// }
	v, completed, err := f.Apply(args, ctxt)
	if err != nil {
		return nil, err
	}
	if completed {
		// we're done
		return &partialResult{nil, nil, v}, nil
	}
	result := NewApply(NewLiteral(v), []AST{})
	// environment returned is effectively ignored since we're going to evaluate an application immediately
	return &partialResult{result, env, nil}, nil
}

func (e *Apply) Str() string {
	strArgs := ""
	for _, item := range e.args {
		strArgs += " " + item.Str()
	}
	return fmt.Sprintf("Apply[%s%s]", e.fn.Str(), strArgs)
}

func (e *Quote) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return e.val, nil
}

func (e *Quote) evalPartial(env *Env, ctxt interface{}) (*partialResult, error) {
	return defaultEvalPartial(e, env, ctxt)
}

func (e *Quote) Str() string {
	return fmt.Sprintf("Quote[%s]", e.val.Str())
}

func (e *LetRec) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return defaultEval(e, env, ctxt)
}

/*
func processPartial (partial *partialResult) (value.Value, bool) {
	if partial.val != nil {
		return partial.val, true
	}
	f := func(args []value.Value, context interface{}) (value.Value, bool, error) {
		// same exact env? I _think_ so
		newPartial, err := partial.exp.evalPartial(partial.env, context)
		if err != nil {
			return nil, true, err
		}
		result, compl := processPartial(newPartial)
		return result, compl, nil
	}
	return NewVFunction([]string{}, f), false
}
*/

func (e *LetRec) evalPartial(env *Env, ctxt interface{}) (*partialResult, error) {
	if len(e.names) != len(e.params) || len(e.names) != len(e.bodies) {
		return nil, errors.New("malformed letrec (names, params, bodies)")
	}
	// create the environment that we'll share across the definitions
	// all names initially allocated #nil
	newEnv := env.Layer(e.names, nil)
	for i, name := range e.names {
		newEnv.Update(name, value.NewVPrimitive("__letrec__", func(args []value.Value, context interface{}) (value.Value, error) {
			newNewEnv := newEnv.Layer(e.params[i], args)
			return e.bodies[i].Eval(newNewEnv, context)
		}))      //e.bodies[i], newEnv})
	}
	return &partialResult{e.body, newEnv, nil}, nil
}

func (e *LetRec) Str() string {
	bindings := make([]string, len(e.names))
	for i := range e.names {
		params := strings.Join(e.params[i], " ")
		bindings[i] = fmt.Sprintf("[%s [%s] %s]", e.names[i], params, e.bodies[i].Str())
	}
	return fmt.Sprintf("LetRec[%s %s]", strings.Join(bindings, " "), e.body.Str())
}
