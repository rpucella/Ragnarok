package evaluator

import (
	"errors"
	"fmt"
	"rpucella.net/ragnarok/internal/value"
)

func defaultEvalPartial(e AST, env *Env, ctxt interface{}) (AST, *Env, value.Value, error) {
	// Partial evaluation
	// Sometimes return an expression to evaluate next along
	// with an environment for evaluation.
	// val is null when the result is in fact a value.

	v, err := e.Eval(env, ctxt)
	if err != nil {
		return nil, nil, nil, err
	}
	return nil, nil, v, nil
}

func defaultEval(e AST, env *Env, ctxt interface{}) (value.Value, error) {
	// evaluation with tail call optimization
	var currExp AST = e
	currEnv := env
	for {
		partialExp, partialEnv, partialVal, err := currExp.evalPartial(currEnv, ctxt)
		if err != nil {
			return nil, err
		}
		if partialVal != nil {
			return partialVal, nil
		}
		currExp = partialExp
		currEnv = partialEnv
	}
}

func (e *Literal) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return e.val, nil
}

func (e *Literal) evalPartial(env *Env, ctxt interface{}) (AST, *Env, value.Value, error) {
	return defaultEvalPartial(e, env, ctxt)
}

func (e *Id) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return env.find(e.name)
}

func (e *Id) evalPartial(env *Env, ctxt interface{}) (AST, *Env, value.Value, error) {
	return defaultEvalPartial(e, env, ctxt)
}

func (e *If) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return defaultEval(e, env, ctxt)
}

func (e *If) evalPartial(env *Env, ctxt interface{}) (AST, *Env, value.Value, error) {
	c, err := e.cnd.Eval(env, ctxt)
	if err != nil {
		return nil, nil, nil, err
	}
	if value.IsTrue(c) {
		return e.thn, env, nil, nil
	} else {
		return e.els, env, nil, nil
	}
}

func (e *Apply) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return defaultEval(e, env, ctxt)
}

func (e *Apply) evalPartial(env *Env, ctxt interface{}) (AST, *Env, value.Value, error) {
	f, err := e.fn.Eval(env, ctxt)
	if err != nil {
		return nil, nil, nil, err
	}
	args := make([]value.Value, len(e.args))
	for i := range args {
		args[i], err = e.args[i].Eval(env, ctxt)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	if ff, ok := f.(*IFunction); ok {
		// do something special if we have an interpreted function!
		if len(ff.params) != len(args) {
			return nil, nil, nil, fmt.Errorf("Wrong number of arguments in application to %s", ff.Str())
		}
		newEnv := ff.env.Layer(ff.params, args)
		return ff.body, newEnv, nil, nil
	}
	v, err := f.Apply(args, ctxt)
	if err != nil {
		return nil, nil, nil, err
	}
	return nil, nil, v, nil
}

func (e *Quote) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return e.val, nil
}

func (e *Quote) evalPartial(env *Env, ctxt interface{}) (AST, *Env, value.Value, error) {
	return defaultEvalPartial(e, env, ctxt)
}

func (e *LetRec) Eval(env *Env, ctxt interface{}) (value.Value, error) {
	return defaultEval(e, env, ctxt)
}

func (e *LetRec) evalPartial(env *Env, ctxt interface{}) (AST, *Env, value.Value, error) {
	if len(e.names) != len(e.params) || len(e.names) != len(e.bodies) {
		return nil, nil, nil, errors.New("malformed letrec (names, params, bodies)")
	}
	// create the environment that we'll share across the definitions
	// all names initially allocated #nil
	newEnv := env.Layer(e.names, nil)
	for i, name := range e.names {
		newEnv.Update(name, NewIFunction(e.params[i], e.bodies[i], newEnv))
	}
	return e.body, newEnv, nil, nil
}
