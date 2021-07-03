package evaluator

import (
	"testing"
)

func TestFunction(t *testing.T) {
}

/*
def test_value_function ():
    # no arguments
    i = lisp.Integer(42)
    e = lisp.Environment()
    b = lisp.VFunction([], i, e)
    assert str(b).startswith('#<FUNCTION')
    assert b.display().startswith('#<FUNCTION')
    assert b.type() == 'function'
    assert not b.is_number()
    assert not b.is_boolean()
    assert not b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert not b.is_cons()
    assert b.is_function()
    assert b.is_atom()
    assert not b.is_list()
    assert b.is_true()
    assert not b.is_equal(lisp.VFunction([], lisp.Integer(42), lisp.Environment()))
    assert not b.is_equal(lisp.VFunction([], lisp.Integer(84), lisp.Environment()))
    assert not b.is_equal(lisp.VInteger(42))
    assert not b.is_eq(lisp.VFunction([], lisp.Integer(42), lisp.Environment()))
    assert not b.is_eq(lisp.VFunction([], lisp.Integer(84), lisp.Environment()))
    assert not b.is_eq(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == ([], i, e)
    result = b.apply(_CONTEXT, [])
    assert result.is_number() and result.value() == 42

    # 2 arguments
    i = lisp.Integer(42)
    e = lisp.Environment()
    b = lisp.VFunction(['x', 'y'], i, e)
    assert b.value() == (['x', 'y'], i, e)
    result = b.apply(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(0)])
    assert result.is_number() and result.value() == 42

    # 2 arguments, one used
    i = lisp.Symbol('x')
    e = lisp.Environment()
    b = lisp.VFunction(['x', 'y'], i, e)
    assert b.value() == (['x', 'y'], i, e)
    result = b.apply(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(0)])
    assert result.is_number() and result.value() == 42

    # 2 arguments, using environment
    i = lisp.Symbol('z')
    e = lisp.Environment(bindings=[('z', lisp.VInteger(42))])
    b = lisp.VFunction(['x', 'y'], i, e)
    assert b.value() == (['x', 'y'], i, e)
    result = b.apply(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(0)])
    assert result.is_number() and result.value() == 42

*/
