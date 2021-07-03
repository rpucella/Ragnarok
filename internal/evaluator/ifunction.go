package evaluator

import (
	"fmt"
	"rpucella.net/ragnarok/internal/value"
	"strings"
)

// interpreted function

type IFunction struct {
	params []string
	body   AST
	env    *Env
}

func NewIFunction(params []string, body AST, env *Env) *IFunction {
	return &IFunction{params, body, env}
}

func (v *IFunction) Display() string {
	return fmt.Sprintf("#<fun %s>", strings.Join(v.params, " "))
}

func (v *IFunction) GetInt() int {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *IFunction) GetString() string {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *IFunction) Apply(args []value.Value, ctxt interface{}) (value.Value, error) {
	if len(v.params) != len(args) {
		return nil, fmt.Errorf("Wrong number of arguments in application to %s", v.Str())
	}
	newEnv := v.env.Layer(v.params, args)
	return v.body.Eval(newEnv, ctxt)
}

func (v *IFunction) Str() string {
	return fmt.Sprintf("IFunction[%s]", strings.Join(v.params, " "))
}

func (v *IFunction) GetHead() value.Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *IFunction) GetTail() value.Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *IFunction) IsTrue() bool {
	return true
}

func (v *IFunction) IsEqual(vv value.Value) bool {
	return v == vv // pointer equality
}

func (v *IFunction) Kind() value.Kind {
	return value.V_FUNCTION
}

func (v *IFunction) GetValue() value.Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *IFunction) SetValue(cv value.Value) {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *IFunction) GetArray() []value.Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}

func (v *IFunction) GetMap() map[string]value.Value {
	panic(fmt.Sprintf("unchecked access to %s", v.Str()))
}
