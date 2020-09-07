import pytest
from src import shell

def test_engine_read ():
    engine = shell.Engine()
    # integer
    inp = '42'
    s = engine.read(inp)
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == '42'
    assert s.as_value().is_number()
    assert s.as_value().value() == 42
    # cons
    inp = '(42 Alice Bob)'
    s = engine.read(inp)
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
    engine = shell.Engine()
    inp = _make_sexp('42')
    v = engine.eval_sexp(inp)
    assert v.is_number()
    assert v.value() == 42
    # application
    engine = shell.Engine()
    inp = _make_sexp([['fun', ['a', 'b'], 'a'], '42', '0'])
    v = engine.eval_sexp(inp)
    assert v.is_number()
    assert v.value() == 42
    # define
    engine = shell.Engine()
    inp = _make_sexp(['def', 'a', '42'])
    engine.eval_sexp(inp)
    v = engine.eval_sexp(_make_sexp('a'))
    assert v.is_number()
    assert v.value() == 42
    # defun
    engine = shell.Engine()
    inp = _make_sexp(['def', ['foo', 'a'], 'a'])
    engine.eval_sexp(inp)
    v = engine.eval_sexp(_make_sexp(['foo', '42']))
    assert v.is_number()
    assert v.value() == 42
    
    
def test_engine_eval ():
    # integer
    engine = shell.Engine()
    inp = '42'
    v = engine.eval(inp)
    assert v.is_number()
    assert v.value() == 42
    # application
    engine = shell.Engine()
    inp = '((fun (a b) a) 42 0)'
    v = engine.eval(inp)
    assert v.is_number()
    assert v.value() == 42
    # define
    engine = shell.Engine()
    inp = '(def a 42)'
    engine.eval(inp)
    v = engine.eval('a')
    assert v.is_number()
    assert v.value() == 42
    # defun
    engine = shell.Engine()
    inp = '(def (foo a) a)'
    engine.eval(inp)
    v = engine.eval('(foo 42)')
    assert v.is_number()
    assert v.value() == 42


def test_engine_bindings ():
    # no init bindings
    engine = shell.Engine()
    v = engine.eval('type')
    assert v.is_function()
    v = engine.eval('empty')
    assert v.is_empty()
    # init bindings
    engine = shell.Engine(bindings=[('a', lisp.VInteger(42)), ('b', lisp.VString('Alice'))])
    v = engine.eval('a')
    assert v.is_number()
    assert v.value() == 42
    v = engine.eval('b')
    assert v.is_string()
    assert v.value() == 'Alice'

    
def test_engine_add_env ():
    # no init bindings
    engine = shell.Engine()
    engine.add_env([('a', lisp.VInteger(42)), ('b', lisp.VString('Alice'))])
    v = engine.eval('a')
    assert v.is_number()
    assert v.value() == 42
    v = engine.eval('empty')
    assert v.is_empty()
    # init bindings
    engine = shell.Engine(bindings=[('x', lisp.VInteger(42)), ('y', lisp.VString('Alice'))])
    engine.add_env([('a', lisp.VInteger(42)), ('b', lisp.VString('Alice'))])
    v = engine.eval('a')
    assert v.is_number()
    assert v.value() == 42
    v = engine.eval('b')
    assert v.is_string()
    assert v.value() == 'Alice'
    v = engine.eval('empty')
    assert v.is_empty()
    v = engine.eval('y')
    assert v.is_string()
    assert v.value() == 'Alice'
    

def test_engine_balance ():
    engine = shell.Engine()
    assert engine.balance('hello')
    assert engine.balance('hello ()')
    assert engine.balance('()')
    assert engine.balance('(1 2 (3 4) 5)')
    assert engine.balance('(()()())')
    assert engine.balance('())')
    assert engine.balance('(1 2 3 (4 5) 6))()')
    assert engine.balance(u'(\u00e9 \u00ea 3 (4 5) 6))()')
    assert not engine.balance('(')
    assert not engine.balance('hello (')
    assert not engine.balance('(()')
    assert not engine.balance('( 1 2 (4)')
    assert not engine.balance('( 1 2 (()(()((')
    
