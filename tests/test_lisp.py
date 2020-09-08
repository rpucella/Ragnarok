import pytest
from src import lisp

# FIXME - share testing function between the various pieces that have commonality
#  e.g., parse_sexp_int <- used to test both parse_sexp() and to test engine.read()

_CONTEXT = {}

def test_environment ():
    
    # basic add + lookup
    e = lisp.Environment()
    e.add('alice', 42)
    assert e.lookup('alice') == 42
    assert e.lookup('ALICE') == 42
    assert list(e.bindings()) == [('ALICE',42)]
    assert e.previous() == None
    
    # initial bindings
    e = lisp.Environment(bindings=[('alice', 42), ('bob', 84)])
    assert e.lookup('alice') == 42
    assert e.lookup('bob') == 84
    assert list(e.bindings()) == [('ALICE', 42), ('BOB', 84)]
    assert e.previous() == None
    
    # linked environments
    e = lisp.Environment(bindings=[('alice', 42)])
    e2 = lisp.Environment(bindings=[('bob', 84)], previous=e)
    assert e2.lookup('alice') == 42
    assert e2.lookup('bob') == 84
    assert e2.previous() == e
    
    # add overwrites existing
    e = lisp.Environment(bindings=[('alice', 42)])
    e.add('alice', 84)
    assert list(e.bindings()) == [('ALICE', 84)]
    e2 = lisp.Environment(previous=e)
    e2.add('alice', 42)
    assert list(e2.bindings()) == [('ALICE', 42)]
    assert list(e2.previous().bindings()) == [('ALICE', 84)]
    
    # updates
    e = lisp.Environment(bindings=[('alice', 42), ('bob', 84)])
    e.update('alice', 168)
    assert list(e.bindings()) == [('ALICE', 168), ('BOB', 84)]
    e = lisp.Environment(bindings=[('alice', 42)])
    e2 = lisp.Environment(bindings=[('bob', 84)], previous=e)
    e2.update('alice', 168)
    assert list(e2.bindings()) == [('ALICE', 168), ('BOB', 84)]
    assert list(e2.previous().bindings()) == [('ALICE', 168)]


#
# Values
#

    
def test_value_boolean ():
    # True
    b = lisp.VBoolean(True)
    assert str(b) == '#T'
    assert b.display() == '#T'
    assert b.type() == 'boolean'
    assert not b.is_number()
    assert b.is_boolean()
    assert not b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert not b.is_cons()
    assert not b.is_function()
    assert b.is_atom()
    assert not b.is_list()
    assert b.is_true()
    assert b.is_equal(lisp.VBoolean(True))
    assert not b.is_equal(lisp.VBoolean(False))
    assert not b.is_equal(lisp.VInteger(42))
    assert b.is_eq(lisp.VBoolean(True))
    assert not b.is_eq(lisp.VBoolean(False))
    assert not b.is_eq(lisp.VInteger(42))
    assert b.value() == True
    # False
    b = lisp.VBoolean(False)
    assert str(b) == '#F'
    assert b.display() == '#F'
    assert b.type() == 'boolean'
    assert not b.is_number()
    assert b.is_boolean()
    assert not b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert not b.is_cons()
    assert not b.is_function()
    assert b.is_atom()
    assert not b.is_list()
    assert not b.is_true()
    assert not b.is_equal(lisp.VBoolean(True))
    assert b.is_equal(lisp.VBoolean(False))
    assert not b.is_equal(lisp.VInteger(42))
    assert not b.is_eq(lisp.VBoolean(True))
    assert b.is_eq(lisp.VBoolean(False))
    assert not b.is_eq(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == False
    

def test_value_string ():
    # empty
    b = lisp.VString(u'')
    assert str(b) == u'""'
    assert b.display() == u''
    assert b.type() == 'string'
    assert not b.is_number()
    assert not b.is_boolean()
    assert b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert not b.is_cons()
    assert not b.is_function()
    assert b.is_atom()
    assert not b.is_list()
    assert not b.is_true()
    assert b.is_equal(lisp.VString(u''))
    assert not b.is_equal(lisp.VString(u'Alice'))
    assert not b.is_equal(lisp.VInteger(42))
    assert not b.is_eq(lisp.VString(u''))
    assert not b.is_eq(lisp.VString(u'Alice'))
    assert not b.is_eq(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == u''
    # non-empty
    b = lisp.VString(u'Alice')
    assert str(b) == u'"Alice"'
    assert b.display() == u'Alice'
    assert b.type() == 'string'
    assert not b.is_number()
    assert not b.is_boolean()
    assert b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert not b.is_cons()
    assert not b.is_function()
    assert b.is_atom()
    assert not b.is_list()
    assert b.is_true()
    assert not b.is_equal(lisp.VString(u''))
    assert b.is_equal(lisp.VString(u'Alice'))
    assert not b.is_equal(lisp.VInteger(42))
    assert not b.is_eq(lisp.VString(u''))
    assert not b.is_eq(lisp.VString(u'Alice'))
    assert not b.is_eq(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == u'Alice'
    # special characters
    b = lisp.VString(u'\\t\\n\\"')
    assert str(b) == u'"\\t\\n\\""'
    assert b.display() == u'\t\n"'
    # accented characters
    b = lisp.VString(u'\u00e9\u00ea\00e8')
    assert str(b) == u'"\u00e9\u00ea\00e8"'
    assert b.display() == u'\u00e9\u00ea\00e8'



def test_value_integer ():
    # zero
    b = lisp.VInteger(0)
    assert str(b) == '0'
    assert b.display() == '0'
    assert b.type() == 'number'
    assert b.is_number()
    assert not b.is_boolean()
    assert not b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert not b.is_cons()
    assert not b.is_function()
    assert b.is_atom()
    assert not b.is_list()
    assert not b.is_true()
    assert b.is_equal(lisp.VInteger(0))
    assert not b.is_equal(lisp.VInteger(42))
    assert not b.is_equal(lisp.VString('Alice'))
    assert b.is_eq(lisp.VInteger(0))
    assert not b.is_eq(lisp.VInteger(42))
    assert not b.is_eq(lisp.VString('Alice'))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == 0
    # non-zero
    b = lisp.VInteger(42)
    assert str(b) == '42'
    assert b.display() == '42'
    assert b.type() == 'number'
    assert b.is_number()
    assert not b.is_boolean()
    assert not b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert not b.is_cons()
    assert not b.is_function()
    assert b.is_atom()
    assert not b.is_list()
    assert b.is_true()
    assert not b.is_equal(lisp.VInteger(0))
    assert b.is_equal(lisp.VInteger(42))
    assert not b.is_equal(lisp.VString('Alice'))
    assert not b.is_eq(lisp.VInteger(0))
    assert b.is_eq(lisp.VInteger(42))
    assert not b.is_eq(lisp.VString('Alice'))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == 42


def test_value_nil ():
    b = lisp.VNil()
    assert str(b) == 'NIL'
    assert b.display() == 'NIL'
    assert b.type() == 'nil'
    assert not b.is_number()
    assert not b.is_boolean()
    assert not b.is_string()
    assert not b.is_symbol()
    assert b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert not b.is_cons()
    assert not b.is_function()
    assert not b.is_atom()
    assert not b.is_list()
    assert not b.is_true()
    assert b.is_equal(lisp.VNil())
    assert not b.is_equal(lisp.VInteger(42))
    assert b.is_eq(lisp.VNil())
    assert not b.is_eq(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() is None


def test_value_reference ():
    # num
    val = lisp.VInteger(42)
    r = lisp.VReference(val)
    assert str(r) == '#<REF 42>'
    assert r.display() == '#<REF 42>'
    assert r.type() == 'ref'
    assert not r.is_number()
    assert not r.is_boolean()
    assert not r.is_string()
    assert not r.is_symbol()
    assert not r.is_nil()
    assert r.is_reference()
    assert not r.is_empty()
    assert not r.is_cons()
    assert not r.is_function()
    assert not r.is_atom()
    assert not r.is_list()
    assert r.is_true()
    assert r.is_equal(lisp.VReference(val))
    assert r.is_equal(lisp.VReference(lisp.VInteger(42)))
    assert not r.is_equal(lisp.VReference(lisp.VInteger(0)))
    assert not r.is_equal(lisp.VReference(lisp.VString("Alice")))
    assert not r.is_eq(lisp.VReference(val))
    assert not r.is_eq(lisp.VReference(lisp.VInteger(42)))
    assert not r.is_eq(lisp.VReference(lisp.VInteger(0)))
    assert not r.is_eq(lisp.VReference(lisp.VString("Alice")))
    assert r.is_equal(r)
    assert r.is_eq(r)
    assert r.value() == val
    # string
    val = lisp.VString("Alice")
    r = lisp.VReference(val)
    assert str(r) == '#<REF "Alice">'
    assert r.display() == '#<REF "Alice">'
    assert r.type() == 'ref'
    assert not r.is_number()
    assert not r.is_boolean()
    assert not r.is_string()
    assert not r.is_symbol()
    assert not r.is_nil()
    assert r.is_reference()
    assert not r.is_empty()
    assert not r.is_cons()
    assert not r.is_function()
    assert not r.is_atom()
    assert not r.is_list()
    assert r.is_true()
    assert not r.is_equal(lisp.VReference(lisp.VInteger(42)))
    assert r.is_equal(lisp.VReference(val))
    assert r.is_equal(lisp.VReference(lisp.VString("Alice")))
    assert not r.is_equal(lisp.VReference(lisp.VString("Bob")))
    assert not r.is_eq(lisp.VReference(lisp.VInteger(42)))
    assert not r.is_eq(lisp.VReference(val))
    assert not r.is_eq(lisp.VReference(lisp.VString("Alice")))
    assert not r.is_eq(lisp.VReference(lisp.VString("Bob")))
    assert r.is_equal(r)
    assert r.is_eq(r)
    assert r.value() == val
    # nested
    val = lisp.VReference(lisp.VInteger(42))
    r = lisp.VReference(val)
    assert str(r) == '#<REF #<REF 42>>'
    assert r.display() == '#<REF #<REF 42>>'
    assert r.type() == 'ref'
    assert not r.is_number()
    assert not r.is_boolean()
    assert not r.is_string()
    assert not r.is_symbol()
    assert not r.is_nil()
    assert r.is_reference()
    assert not r.is_empty()
    assert not r.is_cons()
    assert not r.is_function()
    assert not r.is_atom()
    assert not r.is_list()
    assert r.is_true()
    assert r.is_equal(lisp.VReference(val))
    assert r.is_equal(lisp.VReference(lisp.VReference(lisp.VInteger(42))))
    assert not r.is_equal(lisp.VReference(lisp.VInteger(42)))
    assert not r.is_eq(lisp.VReference(val))
    assert not r.is_eq(lisp.VReference(lisp.VReference(lisp.VInteger(42))))
    assert not r.is_eq(lisp.VReference(lisp.VInteger(42)))
    assert r.is_equal(r)
    assert r.is_eq(r)
    assert r.value() == val
    

def test_value_empty ():
    b = lisp.VEmpty()
    assert str(b) == '()'
    assert b.display() == '()'
    assert b.type() == 'empty-list'
    assert not b.is_number()
    assert not b.is_boolean()
    assert not b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert b.is_empty()
    assert not b.is_cons()
    assert not b.is_function()
    assert not b.is_atom()
    assert b.is_list()
    assert not b.is_true()
    assert b.is_equal(lisp.VEmpty())
    assert not b.is_equal(lisp.VCons(lisp.VInteger(42), lisp.VEmpty()))
    assert not b.is_equal(lisp.VInteger(42))
    assert b.is_eq(lisp.VEmpty())
    assert not b.is_eq(lisp.VCons(lisp.VInteger(42), lisp.VEmpty()))
    assert not b.is_eq(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() is None
    assert b.to_list() == []
    
def test_value_cons_1 ():
    car = lisp.VInteger(42)
    cdr = lisp.VEmpty()
    b = lisp.VCons(car, cdr)
    assert str(b) == '(42)'
    assert b.display() == '(42)'
    assert b.type() == 'cons-list'
    assert not b.is_number()
    assert not b.is_boolean()
    assert not b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert b.is_cons()
    assert not b.is_function()
    assert not b.is_atom()
    assert b.is_list()
    assert b.is_true()
    assert not b.is_equal(lisp.VEmpty())
    assert b.is_equal(lisp.VCons(lisp.VInteger(42), lisp.VEmpty()))
    assert not b.is_equal(lisp.VInteger(42))
    assert not b.is_eq(lisp.VEmpty())
    assert not b.is_eq(lisp.VCons(lisp.VInteger(42), lisp.VEmpty()))
    assert not b.is_eq(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == (car, cdr)
    assert b.to_list() == [car]
    assert b.car() == car
    assert b.cdr() == cdr
    
def test_value_cons_2 ():
    car = lisp.VInteger(42)
    cadr = lisp.VInteger(84)
    cddr = lisp.VEmpty()
    c = lisp.VCons(cadr, cddr)
    b = lisp.VCons(car, c)
    assert str(b) == '(42 84)'
    assert b.display() == '(42 84)'
    assert b.type() == 'cons-list'
    assert not b.is_number()
    assert not b.is_boolean()
    assert not b.is_string()
    assert not b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert b.is_cons()
    assert not b.is_function()
    assert not b.is_atom()
    assert b.is_list()
    assert b.is_true()
    assert not b.is_equal(lisp.VEmpty())
    assert not b.is_equal(lisp.VCons(lisp.VInteger(42), lisp.VEmpty()))
    assert b.is_equal(lisp.VCons(lisp.VInteger(42), lisp.VCons(lisp.VInteger(84), lisp.VEmpty())))
    assert not b.is_equal(lisp.VInteger(42))
    assert not b.is_eq(lisp.VEmpty())
    assert not b.is_eq(lisp.VCons(lisp.VInteger(42), lisp.VEmpty()))
    assert not b.is_eq(lisp.VCons(lisp.VInteger(42), lisp.VCons(lisp.VInteger(84), lisp.VEmpty())))
    assert not b.is_eq(lisp.VCons(lisp.VInteger(42), c))
    assert not b.is_eq(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == (car, c)
    assert b.to_list() == [car, cadr]
    assert b.car() == car
    assert b.cdr() == c
    

def test_value_primitive ():
    def prim (ctxt, args):
        return (args[0], args[1])
    i = lisp.VInteger(42)
    j = lisp.VInteger(0)
    b = lisp.VPrimitive('test', prim, 2)
    assert str(b).startswith('#<PRIMITIVE')
    assert b.display().startswith('#<PRIMITIVE')
    assert b.type() == 'primitive'
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
    assert not b.is_equal(lisp.VPrimitive('test', prim, 2))
    assert not b.is_equal(lisp.VPrimitive('test', lambda args: 0, 2))
    assert not b.is_equal(lisp.VInteger(42))
    assert not b.is_eq(lisp.VPrimitive('test', prim, 2))
    assert not b.is_equal(lisp.VPrimitive('test2', lambda args: 0, 2))
    assert not b.is_equal(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == prim
    assert b.apply(_CONTEXT, [i, j]) == (i, j)

    
def test_value_symbol ():
    b = lisp.VSymbol('Alice')
    assert str(b) == 'ALICE'
    assert b.display() == 'ALICE'
    assert b.type() == 'symbol'
    assert not b.is_number()
    assert not b.is_boolean()
    assert not b.is_string()
    assert b.is_symbol()
    assert not b.is_nil()
    assert not b.is_reference()
    assert not b.is_empty()
    assert not b.is_cons()
    assert not b.is_function()
    assert b.is_atom()
    assert not b.is_list()
    assert b.is_true()
    assert b.is_equal(lisp.VSymbol('Alice'))
    assert b.is_equal(lisp.VSymbol('alice'))
    assert not b.is_equal(lisp.VSymbol('Bob'))
    assert not b.is_equal(lisp.VInteger(42))
    assert b.is_eq(lisp.VSymbol('Alice'))
    assert b.is_eq(lisp.VSymbol('alice'))
    assert not b.is_eq(lisp.VSymbol('Bob'))
    assert not b.is_eq(lisp.VInteger(42))
    assert b.is_equal(b)
    assert b.is_eq(b)
    assert b.value() == 'ALICE'
    # accents
    b = lisp.VSymbol(u'Test\u00e9')
    assert str(b) == u'TEST\u00C9'
    assert b.display() == u'TEST\u00C9'
    assert b.value() == u'TEST\u00C9'
    

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


#
# SExpressions
#
    

def test_sexp_atom_symbol ():
    s = lisp.SAtom(u'Alice')
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == u'Alice'
    assert s.as_value().is_symbol()
    assert s.as_value().value() == u'ALICE'
    # accents
    s = lisp.SAtom(u'Test\u00e9')
    assert s.content() == u'Test\u00e9'
    assert s.as_value().value() == u'TEST\u00c9'

    
def test_sexp_string ():
    s = lisp.SAtom('"Alice"')
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == '"Alice"'
    assert s.as_value().is_string()
    assert s.as_value().value() == 'Alice'
    # accents
    s = lisp.SAtom(u'"Test\u00e9"')
    assert s.content() == u'"Test\u00e9"'
    assert s.as_value().value() == u'Test\u00e9'

    
def test_sexp_integer ():
    s = lisp.SAtom('42')
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == '42'
    assert s.as_value().is_number()
    assert s.as_value().value() == 42

    
def test_sexp_boolean ():
    s = lisp.SAtom('#t')
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == '#t'
    assert s.as_value().is_boolean()
    assert s.as_value().value() == True
    s = lisp.SAtom('#f')
    assert s.is_atom()
    assert s.content() == '#f'
    assert s.as_value().is_boolean()
    assert s.as_value().value() == False

    
def test_sexp_empty ():
    s = lisp.SEmpty()
    assert not s.is_atom()
    assert s.is_empty()
    assert not s.is_cons()
    assert s.content() == None
    assert s.as_value().is_empty()
    
    
def test_sexp_cons ():
    car = lisp.SAtom('42')
    cdr = lisp.SEmpty()
    s = lisp.SCons(car, cdr)
    assert not s.is_atom()
    assert not s.is_empty()
    assert s.is_cons()
    assert s.content() == (car, cdr)
    assert s.as_value().is_cons()
    assert s.as_value().car().is_number()
    assert s.as_value().car().value() == 42
    assert s.as_value().cdr().is_empty()
    
    
def text_sexp_from_value ():
    v = lisp.VBoolean(True)
    assert lisp.SExpression.from_value(v).is_atom()
    assert lisp.SExpression.from_value(v).content() == '#T'
    v = lisp.VBoolean(False)
    assert lisp.SExpression.from_value(v).is_atom()
    assert lisp.SExpression.from_value(v).content() == '#F'
    v = lisp.VString('Alice')
    assert lisp.SExpression.from_value(v).is_atom()
    assert lisp.SExpression.from_value(v).content() == '"Alice"'
    v = lisp.VString(u'Test\u00e9')
    assert lisp.SExpression.from_value(v).is_atom()
    assert lisp.SExpression.from_value(v).content() == u'"Test\u00e9"'
    v = lisp.VInteger(42)
    assert lisp.SExpression.from_value(v).is_atom()
    assert lisp.SExpression.from_value(v).content() == '42'
    v = lisp.VNil()
    assert lisp.SExpression.from_value(v).is_atom()
    assert lisp.SExpression.from_value(v).content() == 'NIL'
    v = lisp.VSymbol('Alice')
    assert lisp.SExpression.from_value(v).is_atom()
    assert lisp.SExpression.from_value(v).content() == 'ALICE'
    v = lisp.VSymbol('Test\u00e9')
    assert lisp.SExpression.from_value(v).is_atom()
    assert lisp.SExpression.from_value(v).content() == 'TEST\u00c9'
    v = lisp.VEmpty()
    assert lisp.SExpression.from_value(v).is_empty()
    v = lisp.VCons(lisp.VInteger(42), lisp.VEmpty())
    assert lisp.SExpression.from_value(v).is_cons()
    assert lisp.SExpression.from_value(v).content()[0].is_number()
    assert lisp.SExpression.from_value(v).content()[0].value() == 42
    assert lisp.SExpression.from_value(v).content()[1].is_empty()
    
    # function? primitive? -- these might be unreadable!?
    


#
# Expressions
#
    

def test_exp_symbol ():
    env = lisp.Environment(bindings=[('Alice', lisp.VInteger(42))])
    e = lisp.Symbol('Alice')
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    e = lisp.Symbol('alice')
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42

    
def test_exp_string ():
    env = lisp.Environment()
    e = lisp.String('')
    v = e.eval(_CONTEXT, env)
    assert v.is_string() and v.value() == ''
    e = lisp.String('Alice')
    v = e.eval(_CONTEXT, env)
    assert v.is_string() and v.value() == 'Alice'
    # accents
    e = lisp.String(u'Test\u00e9')
    v = e.eval(_CONTEXT, env)
    assert v.is_string() and v.value() == u'Test\u00e9'

    
def test_exp_integer ():
    env = lisp.Environment()
    e = lisp.Integer(0)
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 0
    e = lisp.Integer(42)
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42


def test_exp_boolean ():
    env = lisp.Environment()
    e = lisp.Boolean(True)
    v = e.eval(_CONTEXT, env)
    assert v.is_boolean() and v.value() == True
    e = lisp.Boolean(False)
    v = e.eval(_CONTEXT, env)
    assert v.is_boolean() and v.value() == False

    
def test_exp_if ():
    # then branch
    env = lisp.Environment([('a', lisp.VInteger(42))])
    e = lisp.If(lisp.Boolean(True), lisp.Symbol('a'), lisp.Symbol('b'))
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    # else branch
    e = lisp.If(lisp.Boolean(False), lisp.Symbol('b'), lisp.Symbol('a'))
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42


def test_exp_lambda ():
    # simple
    env = lisp.Environment()
    e = lisp.Lambda(['a', 'b'], lisp.Symbol('a'))
    v = e.eval(_CONTEXT, env)
    assert v.is_function()
    v = v.apply(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(0)])
    assert v.is_number() and v.value() == 42
    # environment
    env = lisp.Environment(bindings=[('c', lisp.VInteger(42))])
    e = lisp.Lambda(['a', 'b'], lisp.Symbol('c'))
    v = e.eval(_CONTEXT, env)
    assert v.is_function()
    v = v.apply(_CONTEXT, [lisp.VInteger(1), lisp.VInteger(0)])
    assert v.is_number() and v.value() == 42

    
def test_exp_apply ():
    # simple
    env = lisp.Environment()
    f = lisp.VFunction(['x', 'y'], lisp.Symbol('x'), env)
    env = lisp.Environment(bindings=[('f', f), ('a', lisp.VInteger(42)), ('b', lisp.VInteger(0))])
    e = lisp.Apply(lisp.Symbol('f'),[lisp.Symbol('a'), lisp.Symbol('b')])
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    
    # static vs dynamic binding
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    f = lisp.VFunction(['x', 'y'], lisp.Symbol('a'), env)
    env = lisp.Environment(bindings=[('f', f), ('a', lisp.VInteger(84)), ('b', lisp.VInteger(0))])
    e = lisp.Apply(lisp.Symbol('f'),[lisp.Symbol('a'), lisp.Symbol('b')])
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    

def test_exp_quote ():
    env = lisp.Environment()
    # symbol
    s = lisp.SAtom('Alice')
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_symbol() and v.value() == 'ALICE'
    # symobl (accents)
    s = lisp.SAtom(u'Test\u00e9')
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_symbol() and v.value() == u'TEST\u00c9'
    # string
    s = lisp.SAtom('"Alice"')
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_string() and v.value() == 'Alice'
    # string (accents)
    s = lisp.SAtom(u'"Test\u00e9"')
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_string() and v.value() == u'Test\u00e9'
    # integer
    s = lisp.SAtom('42')
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    # boolean
    s = lisp.SAtom('#t')
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_boolean() and v.value() == True
    s = lisp.SAtom('#f')
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_boolean() and v.value() == False
    # empty
    s = lisp.SEmpty()
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_empty()
    # cons
    s = lisp.SCons(lisp.SAtom('42'), lisp.SEmpty())
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_cons() and v.car().is_number() and v.car().value() == 42 and v.cdr().is_empty()
    # cons 2
    s = lisp.SCons(lisp.SAtom('42'), lisp.SCons(lisp.SAtom('Alice'), lisp.SEmpty()))
    e = lisp.Quote(s)
    v = e.eval(_CONTEXT, env)
    assert v.is_cons()
    assert v.car().is_number()
    assert v.car().value() == 42
    assert v.cdr().is_cons()
    assert v.cdr().car().is_symbol()
    assert v.cdr().car().value() == 'ALICE'
    assert v.cdr().cdr().is_empty()



# maybe use FakeExps for everything?

def test_exp_do ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    # empty
    e = lisp.Do([])
    v = e.eval(_CONTEXT, env)
    assert v.is_nil()
    # single
    e = lisp.Do([lisp.Symbol('a')])
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    # many 
    class FakeExp:   # -- if it quacks like a duck... 
        def __init__ (self, newN):
            self.value = 0
            self.newN = newN
        def eval (self, ctxt, env):
            self.value = self.newN
            return env.lookup('a')
    fe1 = FakeExp(42)
    fe2 = FakeExp(84)
    e = lisp.Do([fe1, fe2, lisp.Symbol('a')])
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    assert fe1.value == 42
    assert fe2.value == 84

    
def test_exp_letrec ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    # empty
    e = lisp.LetRec([], lisp.Symbol('a'))
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    # many / one -> two
    e = lisp.LetRec([('one', lisp.Lambda(['x', 'y'], lisp.Symbol('two'))),
                     ('two', lisp.Lambda(['x'], lisp.Symbol('a')))],
                    lisp.Apply(lisp.Symbol('one'), [lisp.Integer(0), lisp.Integer(0)]))                    
    v = e.eval(_CONTEXT, env)
    assert v.is_function()
    v = v.apply(_CONTEXT, [lisp.VInteger(0)])
    assert v.is_number()
    assert v.value() == 42
    # many / two -> one
    e = lisp.LetRec([('one', lisp.Lambda(['x'], lisp.Symbol('a'))),
                     ('two', lisp.Lambda(['x', 'y'], lisp.Symbol('one')))],
                    lisp.Apply(lisp.Symbol('two'), [lisp.Integer(0), lisp.Integer(0)]))                    
    v = e.eval(_CONTEXT, env)
    assert v.is_function()
    v = v.apply(_CONTEXT, [lisp.VInteger(0)])
    assert v.is_number()
    assert v.value() == 42


def test_sexp_to_exp ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    # symbol
    s = lisp.SAtom('a')
    v = s.to_expression().eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    # string
    s = lisp.SAtom('"Alice"')
    v = s.to_expression().eval(_CONTEXT, env)
    assert v.is_string() and v.value() == 'Alice'
    # integer
    s = lisp.SAtom('42')
    v = s.to_expression().eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    # boolean
    s = lisp.SAtom('#t')
    v = s.to_expression().eval(_CONTEXT, env)
    assert v.is_boolean() and v.value() == True
    s = lisp.SAtom('#f')
    v = s.to_expression().eval(_CONTEXT, env)
    assert v.is_boolean() and v.value() == False
    

#
# SExpression parsing
#

def test_sexp_parse_symbol ():
    inp = u'Alice'
    (s, rest) = lisp.parse_sexp(inp)
    assert rest == u''
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == u'Alice'
    assert s.as_value().is_symbol()
    assert s.as_value().value() == u'ALICE'
    # accents
    inp = u'Test\u00e9'
    (s, rest) = lisp.parse_sexp(inp)
    assert s.content() == u'Test\u00e9'
    assert s.as_value().is_symbol()
    assert s.as_value().value() == u'TEST\u00c9'

    
def test_sexp_parse_string ():
    inp = u'"Alice"'
    (s, rest) = lisp.parse_sexp(inp)
    assert rest == ''
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == u'"Alice"'
    assert s.as_value().is_string()
    assert s.as_value().value() == u'Alice'
    # accents
    inp = u'"Test\u00e9"'
    (s, rest) = lisp.parse_sexp(inp)
    assert s.content() == u'"Test\u00e9"'
    assert s.as_value().is_string()
    assert s.as_value().value() == u'Test\u00e9'

    
def test_sexp_parse_integer ():
    inp = '42'
    (s, rest) = lisp.parse_sexp(inp)
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == '42'
    assert s.as_value().is_number()
    assert s.as_value().value() == 42

    
def test_sexp_parse_boolean ():
    inp = '#t'
    (s, rest) = lisp.parse_sexp(inp)
    assert s.is_atom()
    assert not s.is_empty()
    assert not s.is_cons()
    assert s.content() == '#t'
    assert s.as_value().is_boolean()
    assert s.as_value().value() == True
    inp = '#f'
    (s, rest) = lisp.parse_sexp(inp)
    assert s.is_atom()
    assert s.content() == '#f'
    assert s.as_value().is_boolean()
    assert s.as_value().value() == False

    
def test_sexp_parse_empty ():
    inp = '()'
    (s, rest) = lisp.parse_sexp(inp)
    assert not s.is_atom()
    assert s.is_empty()
    assert not s.is_cons()
    assert s.content() == None
    assert s.as_value().is_empty()
    
    
def test_sexp_parse_cons ():
    inp = '(42 Alice Bob)'
    (s, rest) = lisp.parse_sexp(inp)
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
    
    
def test_sexp_parse_cons_nested ():
    inp = '((42 Alice) ((Bob)))'
    (s, rest) = lisp.parse_sexp(inp)
    assert not s.is_atom()
    assert not s.is_empty()
    assert s.is_cons()
    assert s.as_value().is_cons()
    # (42 Alice)
    assert s.as_value().car().is_cons()
    assert s.as_value().car().car().is_number()
    assert s.as_value().car().car().value() == 42
    assert s.as_value().car().cdr().is_cons()
    assert s.as_value().car().cdr().car().is_symbol()
    assert s.as_value().car().cdr().car().value() == 'ALICE'
    assert s.as_value().car().cdr().cdr().is_empty()
    assert s.as_value().cdr().is_cons()
    # ((Bob))
    assert s.as_value().cdr().car().is_cons()
    assert s.as_value().cdr().car().car().is_cons()
    assert s.as_value().cdr().car().car().car().is_symbol()
    assert s.as_value().cdr().car().car().car().value() == 'BOB'
    assert s.as_value().cdr().car().car().cdr().is_empty()
    assert s.as_value().cdr().car().cdr().is_empty()
    assert s.as_value().cdr().cdr().is_empty()
    

def test_sexp_parse_rest ():
    inp = '42 xyz'
    (s, rest) = lisp.parse_sexp(inp)
    assert rest == ' xyz'
    inp = 'Alice xyz'
    (s, rest) = lisp.parse_sexp(inp)
    assert rest == ' xyz'
    inp = '"Alice" xyz'
    (s, rest) = lisp.parse_sexp(inp)
    assert rest == ' xyz'
    inp = '#t xyz'
    (s, rest) = lisp.parse_sexp(inp)
    assert rest == ' xyz'
    inp = '#f xyz'
    (s, rest) = lisp.parse_sexp(inp)
    assert rest == ' xyz'
    inp = '() xyz'
    (s, rest) = lisp.parse_sexp(inp)
    assert rest == ' xyz'
    inp = '(Alice Bob) xyz'
    (s, rest) = lisp.parse_sexp(inp)
    assert rest == ' xyz'
    

    
#
# Expression parsing
#


def _make_sexp (struct):
    if type(struct) == type([]):
        result = lisp.SEmpty()
        for r in reversed(struct):
            result = lisp.SCons(_make_sexp(r), result)
        return result
    else:
        return lisp.SAtom(struct)

def test_exp_parse_symbol ():
    env = lisp.Environment(bindings=[('Alice', lisp.VInteger(42))])
    inp = _make_sexp('Alice')
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    # accents
    env = lisp.Environment(bindings=[(u'Test\u00e9', lisp.VInteger(42))])
    inp = _make_sexp(u'Test\u00e9')
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    

def test_exp_parse_string ():
    env = lisp.Environment()
    inp = _make_sexp('"Alice"')
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_string() and v.value() == 'Alice'
    # accents
    inp = _make_sexp(u'"Test\u00e9"')
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_string() and v.value() == u'Test\u00e9'

    
def test_exp_parse_integer ():
    env = lisp.Environment()
    inp = _make_sexp('42')
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42


def test_exp_parse_boolean ():
    env = lisp.Environment()
    inp = _make_sexp('#t')
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_boolean() and v.value() == True
    inp = _make_sexp('#f')
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_boolean() and v.value() == False

    
def test_exp_parse_if ():
    # then branch
    env = lisp.Environment([('a', lisp.VInteger(42))])
    inp = _make_sexp(['if', '#t', 'a', '#f'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    # else branch
    inp = _make_sexp(['if', '#f', '#f', 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42


def test_exp_parse_lambda ():
    # simple
    env = lisp.Environment()
    inp = _make_sexp(['fun', ['a', 'b'], 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_function()
    v = v.apply(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(0)])
    assert v.is_number() and v.value() == 42

    
def test_exp_parse_apply ():
    env = lisp.Environment()
    f = lisp.VFunction(['x', 'y'], lisp.Symbol('x'), env)
    env = lisp.Environment(bindings=[('f', f), ('a', lisp.VInteger(42)), ('b', lisp.VInteger(0))])
    inp = _make_sexp(['f', 'a', 'b'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42
    

def test_exp_parse_quote ():
    env = lisp.Environment()
    # symbol
    inp = _make_sexp(['quote', 'Alice'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_symbol() and v.value() == 'ALICE'
    # empty
    inp = _make_sexp(['quote', []])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_empty()
    # cons
    inp = _make_sexp(['quote', ['42']])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_cons() and v.car().is_number() and v.car().value() == 42 and v.cdr().is_empty()


def test_exp_parse_do ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    # empty
    inp = _make_sexp(['do'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_nil()
    # single
    inp = _make_sexp(['do', 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    # many
    inp = _make_sexp(['do', '0', '1', 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42

    
def test_exp_parse_letrec ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    # empty
    inp = _make_sexp(['letrec', [], 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    # many
    inp = _make_sexp(['letrec', [['one', ['fun', ['x', 'y'], 'two']],
                                 ['two', ['fun', ['x'], 'a']]],
                      ['one', '0', '0']])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_function()
    v = v.apply(_CONTEXT, [lisp.VInteger(0)])
    assert v.is_number()
    assert v.value() == 42


def test_exp_parse_let ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    # empty
    inp = _make_sexp(['let', [], 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    # many
    inp = _make_sexp(['let', [['a', '84'], ['b', 'a']], 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 84
    inp = _make_sexp(['let', [['a', '84'], ['b', 'a']], 'b'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42

    
def test_exp_parse_letstar ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    # empty
    inp = _make_sexp(['let*', [], 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    # many
    inp = _make_sexp(['let*', [['a', '84'], ['b', 'a']], 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 84
    inp = _make_sexp(['let*', [['a', '84'], ['b', 'a']], 'b'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 84


def test_exp_parse_and ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    # empty
    inp = _make_sexp(['and'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_boolean()
    assert v.value() == True
    # many
    inp = _make_sexp(['and', 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    inp = _make_sexp(['and', '1', 'a' ])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    inp = _make_sexp(['and', '1', '2', 'a' ])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    inp = _make_sexp(['and', '0', '2', 'a' ])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 0
    inp = _make_sexp(['and', '1', '#f', 'a' ])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_boolean()
    assert v.value() == False

    
def test_exp_parse_or ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    # empty
    inp = _make_sexp(['or'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_boolean()
    assert v.value() == False
    # many
    inp = _make_sexp(['or', 'a'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    inp = _make_sexp(['or', '1', 'a' ])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 1
    inp = _make_sexp(['or', '0', '2', 'a' ])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 2
    inp = _make_sexp(['or', '0', '0', 'a' ])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    inp = _make_sexp(['or', '0', '#f', '0' ])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 0


def test_exp_parse_loop ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42)),
                                     ('=', lisp.VPrimitive('=', lisp.prim_numequal, 2)),
                                     ('+', lisp.VPrimitive('+', lisp.prim_plus, 2))])
    inp = _make_sexp(['loop', 's', [['n', 'a'], ['sum', '0']],
                      ['if', ['=', 'n', '0'], 'sum',
                       ['s', ['+', 'n', '-1'], ['+', 'sum', 'n']]]])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 903


def test_exp_parse_funrec ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42)),
                                     ('=', lisp.VPrimitive('=', lisp.prim_numequal, 2)),
                                     ('+', lisp.VPrimitive('+', lisp.prim_plus, 2))])
    inp = _make_sexp([['funrec', 's', ['n', 'sum'],
                       ['if', ['=', 'n', '0'], 'sum',
                        ['s', ['+', 'n', '-1'], ['+', 'sum', 'n']]]], 'a', '0'])
    e = lisp.Parser().parse_exp(inp)
    v = e.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 903
    

    
#
# Declarations
#

def test_parse_define ():
    env = lisp.Environment()
    inp = _make_sexp(['def', 'a', '42'])
    p = lisp.Parser().parse_define(inp)
    assert type(p) == type((1, 2))
    assert p[0] == 'A'
    v = p[1].eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42

def test_parse_defun ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    inp = _make_sexp(['def', ['foo', 'a', 'b'], 'a'])
    p = lisp.Parser().parse_defun(inp)
    assert type(p) == type((1, 2))
    assert p[0] == 'FOO'
    assert p[1] == ['A', 'B']
    v = p[2].eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42


def test_parse_decl_define ():
    env = lisp.Environment()
    inp = _make_sexp(['def', 'a', '42'])
    r = lisp.Parser().parse(inp)
    assert type(r) == type((1, 2))
    assert r[0] == 'define'
    p = r[1]
    assert type(p) == type((1, 2))
    assert p[0] == 'A'
    v = p[1].eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    

def test_parse_decl_defun ():
    env = lisp.Environment(bindings=[('a', lisp.VInteger(42))])
    inp = _make_sexp(['def', ['foo', 'a', 'b'], 'a'])
    r = lisp.Parser().parse(inp)
    assert type(r) == type((1, 2))
    assert r[0] == 'defun'
    p = r[1]
    assert type(p) == type((1, 2))
    assert p[0] == 'FOO'
    assert p[1] == ['A', 'B']
    v = p[2].eval(_CONTEXT, env)
    assert v.is_number() and v.value() == 42


def test_parse_decl_exp ():
    env = lisp.Environment()
    # int
    inp = _make_sexp('42')
    r = lisp.Parser().parse(inp)
    assert type(r) == type((1, 2))
    assert r[0] == 'exp'
    p = r[1]
    v = p.eval(_CONTEXT, env)
    assert v.is_number()
    assert v.value() == 42
    # lambda
    inp = _make_sexp(['fun', ['a', 'b'], 'a'])
    r = lisp.Parser().parse(inp)
    assert type(r) == type((1, 2))
    assert r[0] == 'exp'
    p = r[1]
    v = p.eval(_CONTEXT, env)
    assert v.is_function()
    v = v.apply(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(0)])
    assert v.is_number() and v.value() == 42
    

#
# Operations
#

def test_prim_type ():
    v = lisp.prim_type(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_symbol() and v.value() == 'BOOLEAN'
    v = lisp.prim_type(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_symbol() and v.value() == 'STRING'
    v = lisp.prim_type(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_symbol() and v.value() == 'NUMBER'
    v = lisp.prim_type(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_symbol() and v.value() == 'REF'
    v = lisp.prim_type(_CONTEXT, [lisp.VNil()])
    assert v.is_symbol() and v.value() == 'NIL'
    v = lisp.prim_type(_CONTEXT, [lisp.VEmpty()])
    assert v.is_symbol() and v.value() == 'EMPTY-LIST'
    v = lisp.prim_type(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_symbol() and v.value() == 'CONS-LIST'
    def prim (ctxt, args):
        return (args[0], args[1])
    v = lisp.prim_type(_CONTEXT, [lisp.VPrimitive('prim', prim, 2)])
    assert v.is_symbol() and v.value() == 'PRIMITIVE'
    v = lisp.prim_type(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_symbol() and v.value() == 'SYMBOL'
    v = lisp.prim_type(_CONTEXT, [lisp.VFunction(['a', 'b'], lisp.Symbol('a'), lisp.Environment())])
    assert v.is_symbol() and v.value() == 'FUNCTION'
    

def test_prim_plus ():
    v = lisp.prim_plus(_CONTEXT, [])
    assert v.is_number() and v.value() == 0
    v = lisp.prim_plus(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_number() and v.value() == 42
    v = lisp.prim_plus(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(84)])
    assert v.is_number() and v.value() == 42 + 84
    v = lisp.prim_plus(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(84), lisp.VInteger(168)])
    assert v.is_number() and v.value() == 42 + 84 + 168

    
def test_prim_times ():
    v = lisp.prim_times(_CONTEXT, [])
    assert v.is_number() and v.value() == 1
    v = lisp.prim_times(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_number() and v.value() == 42
    v = lisp.prim_times(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(84)])
    assert v.is_number() and v.value() == 42 * 84
    v = lisp.prim_times(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(84), lisp.VInteger(168)])
    assert v.is_number() and v.value() == 42 * 84 * 168

    
def test_prim_minus ():
    v = lisp.prim_minus(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_number() and v.value() == -42
    v = lisp.prim_minus(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(84)])
    assert v.is_number() and v.value() == 42 - 84
    v = lisp.prim_minus(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(84), lisp.VInteger(168)])
    assert v.is_number() and v.value() == 42 - 84 - 168


def test_prim_numequal ():
    v = lisp.prim_numequal(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numequal(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numequal(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_numequal(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == True
    

def test_prim_numless ():
    v = lisp.prim_numless(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_numless(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numless(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numless(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    

def test_prim_numlesseq ():
    v = lisp.prim_numlesseq(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_numlesseq(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numlesseq(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_numlesseq(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == True
    

def test_prim_numgreater ():
    v = lisp.prim_numgreater(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numgreater(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_numgreater(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numgreater(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    

def test_prim_numgreatereq ():
    v = lisp.prim_numgreatereq(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numgreatereq(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_numgreatereq(_CONTEXT, [lisp.VInteger(0), lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_numgreatereq(_CONTEXT, [lisp.VInteger(42), lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == True


def test_prim_not ():
    v = lisp.prim_not(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_not(_CONTEXT, [lisp.VBoolean(False)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_not(_CONTEXT, [lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_not(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_not(_CONTEXT, [lisp.VString('')])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_not(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_not(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_not(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == False


def test_prim_string_append ():
    v = lisp.prim_string_append(_CONTEXT, [])
    assert v.is_string() and v.value() == ''
    v = lisp.prim_string_append(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_string() and v.value() == 'Alice'
    v = lisp.prim_string_append(_CONTEXT, [lisp.VString('Alice'), lisp.VString('Bob')])
    assert v.is_string() and v.value() == 'AliceBob'
    v = lisp.prim_string_append(_CONTEXT, [lisp.VString('Alice'), lisp.VString('Bob'), lisp.VString('Charlie')])
    assert v.is_string() and v.value() == 'AliceBobCharlie'
    

def test_prim_string_length ():
    v = lisp.prim_string_length(_CONTEXT, [lisp.VString('')])
    assert v.is_number() and v.value() == 0
    v = lisp.prim_string_length(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_number() and v.value() == 5
    v = lisp.prim_string_length(_CONTEXT, [lisp.VString('Alice Bob')])
    assert v.is_number() and v.value() == 9

    
def test_prim_string_lower ():
    v = lisp.prim_string_lower(_CONTEXT, [lisp.VString('')])
    assert v.is_string() and v.value() == ''
    v = lisp.prim_string_lower(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_string() and v.value() == 'alice'
    v = lisp.prim_string_lower(_CONTEXT, [lisp.VString('Alice Bob')])
    assert v.is_string() and v.value() == 'alice bob'
    

def test_prim_string_upper ():
    v = lisp.prim_string_upper(_CONTEXT, [lisp.VString('')])
    assert v.is_string() and v.value() == ''
    v = lisp.prim_string_upper(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_string() and v.value() == 'ALICE'
    v = lisp.prim_string_upper(_CONTEXT, [lisp.VString('Alice Bob')])
    assert v.is_string() and v.value() == 'ALICE BOB'


def test_prim_string_substring ():
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('')])
    assert v.is_string() and v.value() == ''
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_string() and v.value() == 'Alice'
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('Alice'), lisp.VInteger(0)])
    assert v.is_string() and v.value() == 'Alice'
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('Alice'), lisp.VInteger(1)])
    assert v.is_string() and v.value() == 'lice'
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('Alice'), lisp.VInteger(2)])
    assert v.is_string() and v.value() == 'ice'
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('Alice'), lisp.VInteger(0), lisp.VInteger(5)])
    assert v.is_string() and v.value() == 'Alice'
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('Alice'), lisp.VInteger(0), lisp.VInteger(3)])
    assert v.is_string() and v.value() == 'Ali'
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('Alice'), lisp.VInteger(2), lisp.VInteger(3)])
    assert v.is_string() and v.value() == 'i'
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('Alice'), lisp.VInteger(0), lisp.VInteger(0)])
    assert v.is_string() and v.value() == ''
    v = lisp.prim_string_substring(_CONTEXT, [lisp.VString('Alice'), lisp.VInteger(3), lisp.VInteger(3)])
    assert v.is_string() and v.value() == ''


def _make_list (struct):
    if type(struct) == type([]):
        result = lisp.VEmpty()
        for r in reversed(struct):
            result = lisp.VCons(_make_list(r), result)
        return result
    else:
        return struct

def _unmake_list (lst):
    if lst.is_list():
        current = lst
        result = []
        while current.is_cons():
            result.append(_unmake_list(current.car()))
            current = current.cdr()
        return result
    else:
        return lst
        
    
def test_prim_apply ():
    def prim (ctxt, args):
        return (args[0], args[1])
    v = lisp.prim_apply(_CONTEXT, [lisp.VPrimitive('test', prim, 2, 2),
                         _make_list([lisp.VInteger(42), lisp.VString('Alice')])])
    assert v[0].is_number() and v[0].value() == 42
    assert v[1].is_string() and v[1].value() == 'Alice'
    v = lisp.prim_apply(_CONTEXT, [lisp.VFunction(['a', 'b'], lisp.Symbol('a'), lisp.Environment()),
                         _make_list([lisp.VInteger(42), lisp.VString('Alice')])])
    assert v.is_number() and v.value() == 42


def test_prim_cons ():
    v = lisp.prim_cons(_CONTEXT, [lisp.VInteger(42), lisp.VEmpty()])
    l = _unmake_list(v)
    assert len(l) == 1
    assert l[0].is_number() and l[0].value() == 42
    v = lisp.prim_cons(_CONTEXT, [lisp.VInteger(42), _make_list([lisp.VString('Alice'), lisp.VString('Bob')])])
    l = _unmake_list(v)
    assert len(l) == 3
    assert l[0].is_number() and l[0].value() == 42
    assert l[1].is_string() and l[1].value() == 'Alice'
    assert l[2].is_string() and l[2].value() == 'Bob'


def test_prim_append ():
    v = lisp.prim_append(_CONTEXT, [])
    l = _unmake_list(v)
    assert len(l) == 0
    v = lisp.prim_append(_CONTEXT, [_make_list([lisp.VInteger(1), lisp.VInteger(2)])])
    l = _unmake_list(v)
    assert len(l) == 2
    assert l[0].is_number() and l[0].value() == 1
    assert l[1].is_number() and l[1].value() == 2
    v = lisp.prim_append(_CONTEXT, [_make_list([lisp.VInteger(1), lisp.VInteger(2)]),
                          _make_list([lisp.VInteger(3), lisp.VInteger(4)])])
    l = _unmake_list(v)
    assert len(l) == 4
    assert l[0].is_number() and l[0].value() == 1
    assert l[1].is_number() and l[1].value() == 2
    assert l[2].is_number() and l[2].value() == 3
    assert l[3].is_number() and l[3].value() == 4
    v = lisp.prim_append(_CONTEXT, [_make_list([lisp.VInteger(1), lisp.VInteger(2)]),
                          _make_list([lisp.VInteger(3), lisp.VInteger(4)]),
                          _make_list([lisp.VInteger(5), lisp.VInteger(6)])])
    l = _unmake_list(v)
    assert len(l) == 6
    assert l[0].is_number() and l[0].value() == 1
    assert l[1].is_number() and l[1].value() == 2
    assert l[2].is_number() and l[2].value() == 3
    assert l[3].is_number() and l[3].value() == 4
    assert l[4].is_number() and l[4].value() == 5
    assert l[5].is_number() and l[5].value() == 6


def test_prim_reverse ():
    v = lisp.prim_reverse(_CONTEXT, [_make_list([lisp.VInteger(1),
                                       lisp.VInteger(2),
                                       lisp.VInteger(3),
                                       lisp.VInteger(4)])])
    l = _unmake_list(v)
    assert len(l) == 4
    assert l[0].is_number() and l[0].value() == 4
    assert l[1].is_number() and l[1].value() == 3
    assert l[2].is_number() and l[2].value() == 2
    assert l[3].is_number() and l[3].value() == 1
    

def test_prim_first ():
    v = lisp.prim_first(_CONTEXT, [_make_list([lisp.VInteger(42)])])
    assert v.is_number() and v.value() == 42
    v = lisp.prim_first(_CONTEXT, [_make_list([lisp.VInteger(42),
                                     lisp.VString('Alice'),
                                     lisp.VString('Bob')])])
    assert v.is_number() and v.value() == 42

    
def test_prim_rest ():
    v = lisp.prim_rest(_CONTEXT, [_make_list([lisp.VInteger(42)])])
    l = _unmake_list(v)
    assert len(l) == 0
    v = lisp.prim_rest(_CONTEXT, [_make_list([lisp.VInteger(42),
                                    lisp.VString('Alice'),
                                    lisp.VString('Bob')])])
    l = _unmake_list(v)
    assert len(l) == 2
    assert l[0].is_string() and l[0].value() == 'Alice'
    assert l[1].is_string() and l[1].value() == 'Bob'


def test_prim_list ():
    v = lisp.prim_list(_CONTEXT, [])
    l = _unmake_list(v)
    assert len(l) == 0
    v = lisp.prim_list(_CONTEXT, [lisp.VInteger(42)])
    l = _unmake_list(v)
    assert len(l) == 1
    assert l[0].is_number() and l[0].value() == 42
    v = lisp.prim_list(_CONTEXT, [lisp.VInteger(42),
                                  lisp.VString('Alice')])
    l = _unmake_list(v)
    assert len(l) == 2
    assert l[0].is_number() and l[0].value() == 42
    assert l[1].is_string() and l[1].value() == 'Alice'
    v = lisp.prim_list(_CONTEXT, [lisp.VInteger(42),
                                  lisp.VString('Alice'),
                                  lisp.VString('Bob')])
    l = _unmake_list(v)
    assert len(l) == 3
    assert l[0].is_number() and l[0].value() == 42
    assert l[1].is_string() and l[1].value() == 'Alice'
    assert l[2].is_string() and l[2].value() == 'Bob'
    

def test_prim_length ():
    v = lisp.prim_length(_CONTEXT, [_make_list([])])
    assert v.is_number() and v.value() == 0
    v = lisp.prim_length(_CONTEXT, [_make_list([lisp.VInteger(42)])])
    assert v.is_number() and v.value() == 1
    v = lisp.prim_length(_CONTEXT, [_make_list([lisp.VInteger(42),
                                      lisp.VString('Alice')])])
    assert v.is_number() and v.value() == 2
    v = lisp.prim_length(_CONTEXT, [_make_list([lisp.VInteger(42),
                                      lisp.VString('Alice'),
                                      lisp.VString('Bob')])])
    assert v.is_number() and v.value() == 3


def test_prim_nth ():
    v = lisp.prim_nth(_CONTEXT, [_make_list([lisp.VInteger(42),
                                   lisp.VString('Alice'),
                                   lisp.VString('Bob')]),
                       lisp.VInteger(0)])
    assert v.is_number() and v.value() == 42
    v = lisp.prim_nth(_CONTEXT, [_make_list([lisp.VInteger(42),
                                   lisp.VString('Alice'),
                                   lisp.VString('Bob')]),
                       lisp.VInteger(1)])
    assert v.is_string() and v.value() == 'Alice'
    v = lisp.prim_nth(_CONTEXT, [_make_list([lisp.VInteger(42),
                                   lisp.VString('Alice'),
                                   lisp.VString('Bob')]),
                       lisp.VInteger(2)])
    assert v.is_string() and v.value() == 'Bob'
    

def test_prim_map ():
    def prim1 (ctxt, args):
        return args[0]
    def prim2 (ctxt, args):
        return args[1]
    v = lisp.prim_map(_CONTEXT, [lisp.VPrimitive('test', prim1, 1),
                       _make_list([])])
    l = _unmake_list(v)
    assert len(l) == 0
    v = lisp.prim_map(_CONTEXT, [lisp.VPrimitive('test', prim1, 1),
                       _make_list([lisp.VInteger(42),
                                   lisp.VString('Alice'),
                                   lisp.VString('Bob')])])
    l = _unmake_list(v)
    assert len(l) == 3
    assert l[0].is_number() and l[0].value() == 42
    assert l[1].is_string() and l[1].value() == 'Alice'
    assert l[2].is_string() and l[2].value() == 'Bob'
    v = lisp.prim_map(_CONTEXT, [lisp.VPrimitive('test', prim2, 2),
                       _make_list([]),
                       _make_list([])])
    l = _unmake_list(v)
    assert len(l) == 0
    v = lisp.prim_map(_CONTEXT, [lisp.VPrimitive('test', prim2, 2),
                       _make_list([]),
                       _make_list([lisp.VInteger(42)])])
    l = _unmake_list(v)
    assert len(l) == 0
    v = lisp.prim_map(_CONTEXT, [lisp.VPrimitive('test', prim2, 2),
                       _make_list([lisp.VInteger(42),
                                   lisp.VString('Alice'),
                                   lisp.VString('Bob')]),
                       _make_list([lisp.VInteger(84),
                                   lisp.VString('Charlie'),
                                   lisp.VString('Darlene')])])
    l = _unmake_list(v)
    assert len(l) == 3
    assert l[0].is_number() and l[0].value() == 84
    assert l[1].is_string() and l[1].value() == 'Charlie'
    assert l[2].is_string() and l[2].value() == 'Darlene'


def test_prim_filter ():
    def prim_none (ctxt, args):
        return lisp.VBoolean(False)
    def prim_int (ctxt, args):
        return lisp.VBoolean(args[0].is_number())
    v = lisp.prim_filter(_CONTEXT, [lisp.VPrimitive('test', prim_none, 1),
                          _make_list([])])
    l = _unmake_list(v)
    assert len(l) == 0
    v = lisp.prim_filter(_CONTEXT, [lisp.VPrimitive('test', prim_none, 1),
                          _make_list([lisp.VInteger(42),
                                      lisp.VString('Alice'),
                                      lisp.VString('Bob')])])
    l = _unmake_list(v)
    assert len(l) == 0
    v = lisp.prim_filter(_CONTEXT, [lisp.VPrimitive('test', prim_int, 1),
                          _make_list([lisp.VInteger(42),
                                      lisp.VString('Alice'),
                                      lisp.VString('Bob')])])
    l = _unmake_list(v)
    assert len(l) == 1
    assert l[0].is_number() and l[0].value() == 42


def test_prim_foldr ():
    def prim (ctxt, args):
        return lisp.VString(args[0].value() + '(' + args[1].value() + ')')
    v = lisp.prim_foldr(_CONTEXT, [lisp.VPrimitive('test', prim, 2),
                                   _make_list([]),
                                   lisp.VString('base')])
    assert v.is_string() and v.value() == 'base'
    v = lisp.prim_foldr(_CONTEXT, [lisp.VPrimitive('test', prim, 2),
                                   _make_list([lisp.VString('Alice'),
                                               lisp.VString('Bob'),
                                               lisp.VString('Charlie')]),
                                   lisp.VString('base')])
    assert v.is_string() and v.value() == 'Alice(Bob(Charlie(base)))'


def test_prim_foldl ():
    def prim (ctxt, args):
        return lisp.VString('(' + args[0].value() + ')' + args[1].value())
    v = lisp.prim_foldl(_CONTEXT, [lisp.VPrimitive('test', prim, 2),
                         lisp.VString('base'),
                         _make_list([])])
    assert v.is_string() and v.value() == 'base'
    v = lisp.prim_foldl(_CONTEXT, [lisp.VPrimitive('test', prim, 2),
                         lisp.VString('base'),
                         _make_list([lisp.VString('Alice'),
                                     lisp.VString('Bob'),
                                     lisp.VString('Charlie')])])
    assert v.is_string() and v.value() == '(((base)Alice)Bob)Charlie'


def test_prim_eqp ():
    v = lisp.prim_eqp(_CONTEXT, [lisp.VInteger(42),
                       lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_eqp(_CONTEXT, [lisp.VInteger(42),
                       lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == False
    lst = _make_list([lisp.VInteger(42)])
    v = lisp.prim_eqp(_CONTEXT, [lst, lst])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_eqp(_CONTEXT, [lst, _make_list([lisp.VInteger(42)])])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_eqp(_CONTEXT, [lst, lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_eqp(_CONTEXT, [lst, _make_list([lisp.VInteger(84)])])
    assert v.is_boolean() and v.value() == False


def test_prim_eqlp ():
    v = lisp.prim_eqlp(_CONTEXT, [lisp.VInteger(42),
                       lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_eqlp(_CONTEXT, [lisp.VInteger(42),
                       lisp.VInteger(0)])
    assert v.is_boolean() and v.value() == False
    lst = _make_list([lisp.VInteger(42)])
    v = lisp.prim_eqlp(_CONTEXT, [lst, lst])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_eqlp(_CONTEXT, [lst, _make_list([lisp.VInteger(42)])])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_eqlp(_CONTEXT, [lst, lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_eqlp(_CONTEXT, [lst, _make_list([lisp.VInteger(84)])])
    assert v.is_boolean() and v.value() == False
    ref = lisp.VReference(lisp.VInteger(42))
    v = lisp.prim_eqlp(_CONTEXT, [ref, ref])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_eqlp(_CONTEXT, [ref, lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_eqlp(_CONTEXT, [ref, lisp.VReference(lisp.VInteger(0))])
    assert v.is_boolean() and v.value() == False
    

def test_prim_emptyp ():
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == False 
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_emptyp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == False
   

def test_prim_consp ():
    v = lisp.prim_consp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_consp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_consp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_consp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_consp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_consp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_consp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_consp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_consp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_consp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == False 
    v = lisp.prim_consp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_consp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == False
    

def test_prim_listp ():
    v = lisp.prim_listp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_listp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_listp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_listp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_listp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_listp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_listp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_listp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_listp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_listp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == False 
    v = lisp.prim_listp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_listp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == False


def test_prim_numberp ():
    v = lisp.prim_numberp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numberp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numberp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numberp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_numberp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numberp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numberp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numberp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numberp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numberp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == False 
    v = lisp.prim_numberp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_numberp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == False


def test_prim_booleanp ():
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == False 
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_booleanp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == False
    

def test_prim_stringp ():
    v = lisp.prim_stringp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_stringp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_stringp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_stringp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_stringp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_stringp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_stringp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_stringp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_stringp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_stringp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == False 
    v = lisp.prim_stringp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_stringp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == False


def test_prim_symbolp ():
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == False 
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_symbolp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == False


def test_prim_functionp ():
    v = lisp.prim_functionp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_functionp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_functionp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_functionp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_functionp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_functionp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_functionp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_functionp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_functionp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_functionp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == False 
    v = lisp.prim_functionp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_functionp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == True


def test_prim_nilp ():
    v = lisp.prim_nilp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == True 
    v = lisp.prim_nilp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_nilp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == False


def test_prim_refp ():
    v = lisp.prim_refp(_CONTEXT, [lisp.VEmpty()])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VCons(lisp.VInteger(42), lisp.VEmpty())])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VBoolean(True)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VInteger(42)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VReference(lisp.VInteger(42))])
    assert v.is_boolean() and v.value() == True
    v = lisp.prim_refp(_CONTEXT, [lisp.VString('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VString(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VSymbol('Alice')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VSymbol(u'Test\u00e9')])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VNil()])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VPrimitive('test', lambda args: args[0], 1)])
    assert v.is_boolean() and v.value() == False
    v = lisp.prim_refp(_CONTEXT, [lisp.VFunction(['a'], lisp.VSymbol('a'), lisp.Environment())])
    assert v.is_boolean() and v.value() == False


def test_prim_ref ():
    val = lisp.VInteger(42)
    v = lisp.prim_ref(_CONTEXT, [val])
    assert v.is_reference() and v.value() == val
    val = lisp.VString("Alice")
    v = lisp.prim_ref(_CONTEXT, [val])
    assert v.is_reference() and v.value() == val
    val = lisp.VReference(lisp.VInteger(42))
    v = lisp.prim_ref(_CONTEXT, [val])
    assert v.is_reference() and v.value() == val


def test_prim_ref_get ():
    val = lisp.VReference(lisp.VInteger(42))
    v = lisp.prim_ref_get(_CONTEXT, [val])
    assert v.is_number() and v.value() == 42
    val = lisp.VReference(lisp.VString("Alice"))
    v = lisp.prim_ref_get(_CONTEXT, [val])
    assert v.is_string() and v.value() == "Alice"
    val = lisp.VReference(lisp.VReference(lisp.VInteger(42)))
    v = lisp.prim_ref_get(_CONTEXT, [val])
    assert v.is_reference()
    assert v.value().is_number() and v.value().value() == 42
    

def test_prim_ref_set ():
    val = lisp.VReference(lisp.VInteger(0))
    v = lisp.prim_ref_set(_CONTEXT, [val, lisp.VInteger(42)])
    assert v.is_nil()
    assert val.value().is_number() and val.value().value() == 42
    val = lisp.VReference(lisp.VInteger(0))
    v = lisp.prim_ref_set(_CONTEXT, [val, lisp.VString("Alice")])
    assert v.is_nil()
    assert val.value().is_string() and val.value().value() == "Alice"
    val = lisp.VReference(lisp.VInteger(0))
    v = lisp.prim_ref_set(_CONTEXT, [val, lisp.VReference(lisp.VInteger(42))])
    assert v.is_nil()
    assert val.value().is_reference()
    assert val.value().value().is_number() and val.value().value().value() == 42
    
    
