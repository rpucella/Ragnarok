package evaluator

import (
	"fmt"
	"rpucella.net/ragnarok/internal/value"
	"strings"
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
	Display() string
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

func (e *Literal) Str() string {
	return fmt.Sprintf("Literal[%s]", e.val.Str())
}

func (e *Literal) Display() string {
	return e.val.Display()
}

func (e *Id) Str() string {
	return fmt.Sprintf("Id[%s]", e.name)
}

func (e *Id) Display() string {
	return e.name
}

func (e *If) Str() string {
	return fmt.Sprintf("If[%s %s %s]", e.cnd.Str(), e.thn.Str(), e.els.Str())
}

func (e *If) Display() string {
	return fmt.Sprintf("(if %s %s %s)", e.cnd.Display(), e.thn.Display(), e.els.Display())
}

func (e *Apply) Str() string {
	strArgs := ""
	for _, item := range e.args {
		strArgs += " " + item.Str()
	}
	return fmt.Sprintf("Apply[%s%s]", e.fn.Str(), strArgs)
}

func (e *Apply) Display() string {
	strArgs := ""
	for _, item := range e.args {
		strArgs += " " + item.Display()
	}
	return fmt.Sprintf("(%s%s)", e.fn.Display(), strArgs)
}

func (e *Quote) Str() string {
	return fmt.Sprintf("Quote[%s]", e.val.Str())
}

func (e *Quote) Display() string {
	return fmt.Sprintf("(quote %s)", e.val.Display())
}

func (e *LetRec) Str() string {
	bindings := make([]string, len(e.names))
	for i := range e.names {
		params := strings.Join(e.params[i], " ")
		bindings[i] = fmt.Sprintf("[%s [%s] %s]", e.names[i], params, e.bodies[i].Str())
	}
	return fmt.Sprintf("LetRec[%s %s]", strings.Join(bindings, " "), e.body.Str())
}

func (e *LetRec) Display() string {
	bindings := make([]string, len(e.names))
	for i := range e.names {
		params := strings.Join(e.params[i], " ")
		bindings[i] = fmt.Sprintf("(%s (%s) %s)", e.names[i], params, e.bodies[i].Display())
	}
	return fmt.Sprintf("(letrec %s %s)", strings.Join(bindings, " "), e.body.Display())
}
