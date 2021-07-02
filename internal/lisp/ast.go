package lisp

import "fmt"
import "errors"
import "strings"

const DEF_VALUE = 0
const DEF_FUNCTION = 1

type Def struct {
	Name   string
	Type   int
	Params []string
	Body   AST
}

type AST interface {
	Eval(*Env, interface{}) (Value, error)
	evalPartial(*Env, interface{}) (*PartialResult, error)
	Str() string
}

type PartialResult struct {
	exp AST
	env *Env
	val Value // val is null when the result is still partial
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

func NewLiteral(val Value) *Literal {
	return &Literal{val}
}

func NewQuote(val Value) *Quote {
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

func defaultEvalPartial(e AST, env *Env, ctxt interface{}) (*PartialResult, error) {
	// Partial evaluation
	// Sometimes return an expression to evaluate next along
	// with an environment for evaluation.
	// val is null when the result is in fact a value.

	v, err := e.Eval(env, ctxt)
	if err != nil {
		return nil, err
	}
	return &PartialResult{nil, nil, v}, nil
}

func defaultEval(e AST, env *Env, ctxt interface{}) (Value, error) {
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

func (e *Literal) Eval(env *Env, ctxt interface{}) (Value, error) {
	return e.val, nil
}

func (e *Literal) evalPartial(env *Env, ctxt interface{}) (*PartialResult, error) {
	return defaultEvalPartial(e, env, ctxt)
}

func (e *Literal) Str() string {
	return fmt.Sprintf("Literal[%s]", e.val.Str())
}

func (e *Id) Eval(env *Env, ctxt interface{}) (Value, error) {
	return env.find(e.name)
}

func (e *Id) evalPartial(env *Env, ctxt interface{}) (*PartialResult, error) {
	return defaultEvalPartial(e, env, ctxt)
}

func (e *Id) Str() string {
	return fmt.Sprintf("Id[%s]", e.name)
}

func (e *If) Eval(env *Env, ctxt interface{}) (Value, error) {
	return defaultEval(e, env, ctxt)
}

func (e *If) evalPartial(env *Env, ctxt interface{}) (*PartialResult, error) {
	c, err := e.cnd.Eval(env, ctxt)
	if err != nil {
		return nil, err
	}
	if c.IsTrue() {
		return &PartialResult{e.thn, env, nil}, nil
	} else {
		return &PartialResult{e.els, env, nil}, nil
	}
}

func (e *If) Str() string {
	return fmt.Sprintf("If[%s %s %s]", e.cnd.Str(), e.thn.Str(), e.els.Str())
}

func (e *Apply) Eval(env *Env, ctxt interface{}) (Value, error) {
	return defaultEval(e, env, ctxt)
}

func (e *Apply) evalPartial(env *Env, ctxt interface{}) (*PartialResult, error) {
	f, err := e.fn.Eval(env, ctxt)
	if err != nil {
		return nil, err
	}
	args := make([]Value, len(e.args))
	for i := range args {
		args[i], err = e.args[i].Eval(env, ctxt)
		if err != nil {
			return nil, err
		}
	}
	if ff, ok := f.(*VFunction); ok {
		if len(ff.params) != len(args) {
			return nil, fmt.Errorf("Wrong number of arguments to application to %s", ff.Str())
		}
		newEnv := ff.env.Layer(ff.params, args)
		return &PartialResult{ff.body, newEnv, nil}, nil
	}
	v, err := f.Apply(args, ctxt)
	if err != nil {
		return nil, err
	}
	return &PartialResult{nil, nil, v}, nil
}

func (e *Apply) Str() string {
	strArgs := ""
	for _, item := range e.args {
		strArgs += " " + item.Str()
	}
	return fmt.Sprintf("Apply[%s%s]", e.fn.Str(), strArgs)
}

func (e *Quote) Eval(env *Env, ctxt interface{}) (Value, error) {
	return e.val, nil
}

func (e *Quote) evalPartial(env *Env, ctxt interface{}) (*PartialResult, error) {
	return defaultEvalPartial(e, env, ctxt)
}

func (e *Quote) Str() string {
	return fmt.Sprintf("Quote[%s]", e.val.Str())
}

func (e *LetRec) Eval(env *Env, ctxt interface{}) (Value, error) {
	return defaultEval(e, env, ctxt)
}

func (e *LetRec) evalPartial(env *Env, ctxt interface{}) (*PartialResult, error) {
	if len(e.names) != len(e.params) || len(e.names) != len(e.bodies) {
		return nil, errors.New("malformed letrec (names, params, bodies)")
	}
	// create the environment that we'll share across the definitions
	// all names initially allocated #nil
	newEnv := env.Layer(e.names, nil)
	for i, name := range e.names {
		newEnv.Update(name, &VFunction{e.params[i], e.bodies[i], newEnv})
	}
	return &PartialResult{e.body, newEnv, nil}, nil
}

func (e *LetRec) Str() string {
	bindings := make([]string, len(e.names))
	for i := range e.names {
		params := strings.Join(e.params[i], " ")
		bindings[i] = fmt.Sprintf("[%s [%s] %s]", e.names[i], params, e.bodies[i].Str())
	}
	return fmt.Sprintf("LetRec[%s %s]", strings.Join(bindings, " "), e.body.Str())
}
