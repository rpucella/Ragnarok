import pytest
from src import engine
from src import lisp

_CONTEXT = {'env': lisp.Environment() }

def _make_sexp (struct):
    if type(struct) == type([]):
        result = lisp.SEmpty()
        for r in reversed(struct):
            result = lisp.SCons(_make_sexp(r), result)
        return result
    else:
        return lisp.SAtom(struct)

def test_engine_read ():
    eng = engine.Engine()
    # integer
    inp = '42'
    s = eng.read(inp)
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == '42'
    assert s.as_value().is_number()
    assert s.as_value().value() == 42
    # cons
    inp = '(42 Alice Bob)'
    s = eng.read(inp)
    assert not s.is_atom()
    assert not s.is_empty()
    assert s.is_cons()
    assert s.as_value().is_cons()
    assert s.as_value().car().is_number()
    assert s.as_value().car().value() == 42
    assert s.as_value().cdr().is_cons()
    assert s.as_value().cdr().car().is_symbol()
    assert s.as_value().cdr().car().value() == 'ALICE'
    assert s.as_value().cdr().cdr().is_cons()
    assert s.as_value().cdr().cdr().car().is_symbol()
    assert s.as_value().cdr().cdr().car().value() == 'BOB'
    assert s.as_value().cdr().cdr().cdr().is_empty()


def test_engine_eval_sexp ():
    # integer
    eng = engine.Engine()
    inp = _make_sexp('42')
    v = eng.eval_sexp(_CONTEXT, inp)
    assert v['result'].is_number()
    assert v['result'].value() == 42
    # application
    eng = engine.Engine()
    inp = _make_sexp([['fun', ['a', 'b'], 'a'], '42', '0'])
    v = eng.eval_sexp(_CONTEXT, inp)
    assert v['result'].is_number()
    assert v['result'].value() == 42
    # define
    eng = engine.Engine()
    inp = _make_sexp(['def', 'a', '42'])
    eng.eval_sexp(_CONTEXT, inp)
    v = eng.eval_sexp(_CONTEXT, _make_sexp('a'))
    assert v['result'].is_number()
    assert v['result'].value() == 42
    # defun
    eng = engine.Engine()
    inp = _make_sexp(['def', ['foo', 'a'], 'a'])
    eng.eval_sexp(_CONTEXT, inp)
    v = eng.eval_sexp(_CONTEXT, _make_sexp(['foo', '42']))
    assert v['result'].is_number()
    assert v['result'].value() == 42
    
    
def test_engine_eval ():
    # integer
    eng = engine.Engine()
    inp = '42'
    v = eng.eval(_CONTEXT, inp)
    assert v['result'].is_number()
    assert v['result'].value() == 42
    # application
    eng = engine.Engine()
    inp = '((fun (a b) a) 42 0)'
    v = eng.eval(_CONTEXT, inp)
    assert v['result'].is_number()
    assert v['result'].value() == 42
    # define
    eng = engine.Engine()
    inp = '(def a 42)'
    eng.eval(_CONTEXT, inp)
    v = eng.eval(_CONTEXT, 'a')
    assert v['result'].is_number()
    assert v['result'].value() == 42
    # defun
    eng = engine.Engine()
    inp = '(def (foo a) a)'
    eng.eval(_CONTEXT, inp)
    v = eng.eval(_CONTEXT, '(foo 42)')
    assert v['result'].is_number()
    assert v['result'].value() == 42
