package evaluator

import (
	"errors"
	"fmt"
	"rpucella.net/ragnarok/internal/value"
)

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

func (e *Id) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return env.find(e.name)
}

func (e *Id) evalPartial(env *Env, ctxt interface{}) (*partialResult, error) {
	return defaultEvalPartial(e, env, ctxt)
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
	if ff, ok := f.(*IFunction); ok {
		// do something special if we have an interpreted function!
		if len(ff.params) != len(args) {
			return nil, fmt.Errorf("Wrong number of arguments in application to %s", ff.Str())
		}
		newEnv := ff.env.Layer(ff.params, args)
		return &partialResult{ff.body, newEnv, nil}, nil
	}
	v, err := f.Apply(args, ctxt)
	if err != nil {
		return nil, err
	}
	return &partialResult{nil, nil, v}, nil
}

func (e *Quote) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return e.val, nil
}

func (e *Quote) evalPartial(env *Env, ctxt interface{}) (*partialResult, error) {
	return defaultEvalPartial(e, env, ctxt)
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
		/*
			newEnv.Update(name, value.NewVPrimitive("__letrec__", func(args []value.Value, context interface{}) (value.Value, error) {
				newNewEnv := newEnv.Layer(e.params[i], args)
				return e.bodies[i].Eval(newNewEnv, context)
			}))      //e.bodies[i], newEnv})
		*/
		newEnv.Update(name, NewIFunction(e.params[i], e.bodies[i], newEnv))
	}
	return &partialResult{e.body, newEnv, nil}, nil
}
