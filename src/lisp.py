# Python 3

import sys
import re
import functools


class LispError (Exception):
    def __init__ (self, msg, type='error'):
        super(LispError, self).__init__(f'{type.upper()}: {msg}')

class LispWrongArgNoError (LispError):
    pass

class LispWrongArgTypeError (LispError):
    pass

class LispReadError (LispError):
    pass

class LispParseError (LispError):
    pass

class LispQuit (Exception):
    pass


class Environment (object):
    def __init__ (self, bindings=[], previous=None):
        self._previous = previous
        self._bindings = {}
        for (name, value) in bindings:
            self.add(name, value)

    def add (self, symbol, value):
        '''
        Add a binding to the current environment
        Replaces old binding if one exists
        '''
        symbol = symbol.upper()
        self._bindings[symbol] = value

    def update (self, symbol, value):
        '''
        Update an existing binding
        Look back into higher environments to see
        if it is there
        If not, add it to current environment
        '''
        symbol = symbol.upper()
        if symbol in self._bindings:
            self._bindings[symbol] = value
            return True
        updated = self._previous and self._previous.update(symbol, value)
        if not updated:
            # if the binding doesn't exist, add it locally
            self.add(symbol, value)

    def find (self, symbol):
        '''
        Look for a binding up the environment chain.
        '''
        symbol = symbol.upper()
        if symbol in self._bindings:
            return self._bindings[symbol]
        if self._previous:
            return self._previous.find(symbol)
        return None
            
    def lookup (self, symbol):
        '''
        Look for a binding up the environment chain.
        '''
        result = self.find(symbol)
        if result:
            return result
        raise LispError(f'Cannot find binding for {symbol}')

    def bindings (self, as_dict=False):
        if self._previous:
            result = self._previous.bindings(as_dict=True)
        else:
            result = {}
        for n in self._bindings:
            result[n] = self._bindings[n]
        if as_dict:
            return result
        else:
            return result.items()

    def names (self):
        return [n for (n, _) in self.bindings()]

    def modules (self):
        return [n for (n, v) in self.bindings() if v.is_module()]

    def previous (self):
        return self._previous

    
# class RootEnvironment (object):
#     def __init__ (self):
#         self._modules = {}
#
#     def modules (self):
#         return self._modules.keys()
#        
#     def add_module (self, name, bindings=[]):
#         '''
#         Add a new module to the root of the environment
#         '''
#         name = name.upper()
#         if name in self._modules:
#             raise LispError('Module {} already exists'.format(name))
#         env = Environment(root=self, bindings=bindings)
#         self._modules[name] = env
#         return env
#
#        
#   # def lookup_module (self, module):
    #     '''
    #     Look for a module 
    #     '''
    #     module = module.upper()
    #     if module in self._modules:
    #         return self._modules[module]
    #     raise LispError('Cannot find module {}'.format(module))


class Value (object):

    def to_list (self):
        raise LispError('Cannot convert to a python list of values: {}'.format(self))

    def _str_cdr (self):
        raise LispError('Cannot use value as list terminator: {}'.format(self))

    def display (self):
        return str(self)

    def type (self):
        pass

    def __str__ (self):
        pass

    def is_number (self):
        return self.type() == 'number'

    def is_boolean (self):
        return self.type() == 'boolean'

    def is_string (self):
        return self.type() == 'string'

    def is_symbol (self):
        return self.type() == 'symbol'

    def is_nil (self):
        return self.type() == 'nil'

    def is_empty (self):
        return self.type() == 'empty-list'
    
    def is_cons (self):
        return self.type() == 'cons-list'

    def is_function (self):
        return self.type() in ('primitive', 'function')

    def is_reference (self):
        return self.type() == 'ref'

    def is_atom (self):
        return self.type() in ['number', 'primitive', 'function', 'symbol', 'string', 'boolean']

    def is_list (self):
        return self.type() in ['empty-list', 'cons-list']

    def is_dict (self):
        return self.type() in ['dict']

    # TODO: add unit tests for modules
    def is_module (self):
        return self.type() in ['module']

    def is_true (self):
        return True

    def is_equal (self, v):
        # by default, do is_eq
        return self.is_eq(v)

    def is_eq (self, v):
        # "pointer" equality
        ##self.type() == v.type() and self.value() == v.value()
        return id(self) == id(v)

    def apply (self, ctxt, args):
        raise LispError('Cannot apply value {}'.format(self))

    
class VBoolean (Value):
    def __init__ (self, b):
        self._value = b

    def __repr__ (self):
        return 'VBoolean({})'.format(self._value)

    def __str__ (self):
        return '#T' if self._value else '#F'

    def type (self):
        return 'boolean'

    def value (self):
        return self._value

    def is_true (self):
        return self._value
        
    def is_eq (self, v):
        return v.is_boolean() and self.value() == v.value()


class VReference (Value):
    def __init__ (self, v):
        self._value = v

    def __repr__ (self):
        return 'VReference({})'.format(self._value)

    def __str__ (self):
        return '#<REF {}>'.format(self._value)

    def type (self):
        return 'ref'

    def value (self):
        return self._value

    def set_value (self, v):
        self._value = v

    def is_equal (self, v):
        return v.is_reference() and self.value().is_equal(v.value())


class VDict (Value):
    def __init__ (self, entries):
        self._value = entries

    def __repr__ (self):
        return 'VDict({})'.format(self._value)

    def __str__ (self):
        entries = ['({})'.format(' '.join([ str(x) for x in v])) for v in self._value]
        return '#<DICT {}>'.format(' '.join(entries))

    def type (self):
        return 'dict'

    def is_equal (self, v):
        # TODO - fix this comparions!
        return v.is_dict() and self.value().is_equal(v.value())

    def lookup (self, v):
        for (key, value) in self._value:
            if key.is_equal(v):
                return value
        raise LispError('Cannot find key {} in dictionary'.format(v))

    def update (self, k, v):
        result = []
        for (key, value) in self._value:
            if key.is_equal(k):
                result.append((key, v))
            else:
                result.append((key, value))
        else:
            result.append((k, v))
        return VDict(result)

    def set (self, k, v):
        for (i, (key, value)) in enumerate(self._value):
            if key.is_equal(k):
                self._value[i] = (key, v)
        else:
            self._value.append((k,v))
        return VNil()
    
    
class VString (Value):
    def __init__ (self, s):
        self._value = s

    def __repr__ (self):
        return 'VString({})'.format(self._value)

    def __str__ (self):
        return '"{}"'.format(self._value)

    def type (self):
        return 'string'

    def display (self):
        return self._value.replace('\\"', '"').replace('\\t', '\t').replace('\\n', '\n').replace('\\\\','\\')

    def value (self):
        return self._value
        
    def is_true (self):
        return not (not self._value)
        
    def is_equal (self, v):
        return v.is_string() and self.value() == v.value()
    
    
class VInteger (Value):
    def __init__ (self, v):
        self._value = v

    def __repr__ (self):
        return 'VInteger({})'.format(self._value)

    def __str__ (self):
        return str(self._value)

    def type (self):
        return 'number'

    def value (self):
        return self._value

    def is_true (self):
        return not (not self._value)
        
    def is_eq (self, v):
        return v.is_number() and self.value() == v.value()


class VNil (Value):
    def __repr__ (self):
        return 'VNil()'

    def __str__ (self):
        return 'NIL'

    def type (self):
        return 'nil'

    def is_true (self):
        return False

    def value (self):
        return None

    def is_eq (self, v):
        return v.is_nil()


class VEmpty (Value):
    def __repr__ (self):
        return 'VEmpty()'

    def to_list (self):
        return []
    
    def __str__ (self):
        return '()'

    def _str_cdr (self):
        return ')'

    def type (self):
        return 'empty-list'

    def is_true (self):
        return False

    def value (self):
        return None

    def is_eq (self, v):
        return v.is_empty()

    
class VCons (Value):
    def __init__ (self, car, cdr):
        if not cdr.is_list():
            raise LispError('List required as second cons argument')
        self._car = car
        self._cdr = cdr

    def to_list (self):
        lst = self._cdr.to_list()
        return [self._car] + lst
    
    def __repr__ (self):
        return 'VCons({},{})'.format(repr(self._car), repr(self._cdr))

    def __str__ (self):
        return '({}{}'.format(self._car, self._cdr._str_cdr())

    def _str_cdr (self):
        return ' {}{}'.format(self._car, self._cdr._str_cdr())

    def type (self):
        return 'cons-list'

    def value (self):
        return (self._car, self._cdr)

    def car (self):
        return self._car

    def cdr (self):
        return self._cdr

    def is_equal (self, v):
        return v.is_cons() and self.car().is_equal(v.car()) and self.cdr().is_equal(v.cdr())
    

class VPrimitive (Value):
    def __init__ (self, name, primitive, min, max=None):
        self._name = name
        self._primitive = primitive
        self._min = min
        self._max = max
        
    def __repr__ (self):
        return 'VPrimitive({})'.format(self._primitive.__name__)

    def __str__ (self):
        h = id(self)
        return '#<PRIMITIVE {}>'.format(hex(h))

    def type (self):
        return 'primitive'

    def value (self):
        return self._primitive

    def apply (self, ctxt, values):
        if len(values) < self._min:
            raise LispWrongArgNoError('Too few arguments {} to primitive {}'.format(len(values), self._name))
        if self._max and len(values) > self._max:
            raise LispWrongArgNoError('Too many arguments {} to primitive {}'.format(len(values), self._name))
        result = self._primitive(ctxt, values)
        return (result or VNil())
    
    
class VSymbol (Value):
    def __init__ (self, sym):
        self._symbol = sym.upper()

    def __repr__ (self):
        return 'VSymbol({})'.format(self._symbol)

    def __str__ (self):
        return self._symbol

    def type (self):
        return 'symbol'

    def value (self):
        return self._symbol

    def is_eq (self, v):
        return v.is_symbol() and self.value() == v.value()
    
    
class VFunction (Value):
    def __init__ (self, params, body, env):
        self._params = params
        self._body = body
        self._env = env

    def __repr__ (self):
        return 'VFunction({}, {})'.format(self._params, repr(self._body))

    def __str__ (self):
        h = id(self)
        return '#<FUNCTION {}>'.format(hex(h))

    def binding_env (self, values):
        if len(self._params) != len(values):
            raise LispWrongArgNoError('Wrong number of arguments to {}'.format(self))
        params_bindings = list(zip(self._params, values))
        new_env = Environment(previous=self._env)
        for (x, y) in params_bindings:
            new_env.add(x, y)
        return new_env

    def type (self):
        return 'function'

    def value (self):
        return (self._params, self._body, self._env)


class VModule (Value):
    def __init__ (self, env):
        self._env = env

    def __repr__ (self):
        return 'VModule({})'.format(', '.join(self._env.names()))

    def __str__ (self):
        return '#<MODULE {}>'.format(' '.join(self._env.names()))

    def type (self):
        return 'module'

    def env (self):
        return self._env

    def is_true (self):
        return True

    
    

class Expression (object):

    def eval_partial (self, ctxt, env):
        ''' 
        Partial evaluation.
        Sometimes return an expression to evaluate next along 
        with an environment for evaluation.
        Environment is None when the result is in fact a value.
        '''
        return (self.eval(ctxt, env), None)

    def eval (self, ctxt, env):
        '''
        Evaluation with tail-call optimization.
        '''
        curr_exp = self
        curr_env = env

        while (True):
            (new_exp, new_env) = curr_exp.eval_partial(ctxt, curr_env)
            if new_env is None:
                # actually a value!
                return new_exp
            curr_exp = new_exp
            curr_env = new_env
    
    
class Symbol (Expression):
    def __init__ (self, sym, qualifiers=[]):
        self._symbol = sym.upper()
        self._qualifiers = qualifiers

    def __repr__ (self):
        if self._qualifiers:
            return 'Symbol({}, {})'.format(':'.join(self._qualifiers), self._symbol)
        return 'Symbol({})'.format(self._symbol)

    def eval (self, ctxt, env):
        if self._qualifiers:
            if len(self._qualifiers) > 1:
                raise LispError('No support for nested modules yet')
            v = env.lookup(self._qualifiers[0])
            if not v.is_module():
                raise LispError('Symbol {} does not represent a module'.format(self._qualifiers[0]))
            v = v.env().lookup(self._symbol)
        else:
            # unqualified name - look in the default environment + opened modules if any
            v = env.find(self._symbol)
            if v is None:
                if 'modules' in ctxt:
                    for m in ctxt['modules']:
                        mv = env.lookup(m)
                        if mv.is_module():
                            v = mv.env().find(self._symbol)
                            if v is not None:
                                break
                    else:
                        raise LispError(f'Cannot find binding for {self._symbol}')
                else:
                    raise LispError(f'Cannot find binding for {self._symbol}')
        if v is None:
            # this can't arise at this point I think...
            raise LispError('Trying to access a non-initialized binding {} in a LETREC'.format(self._symbol))
        return v   #  env.lookup(self._symbol)


class String (Expression):
    def __init__ (self, s):
        self._string = s

    def __repr__ (self):
        return 'String({})'.format(self._string)
                           
    def eval (self, ctxt, env):
        return VString(self._string)
                            
    
class Integer (Expression):
    def __init__ (self, s):
        self._value = int(s)

    def __repr__ (self):
        return 'Integer({})'.format(self._value)
                            
    def eval (self, ctxt, env):
        return VInteger(self._value)

    
class Boolean (Expression):
    def __init__ (self, b):
        self._value = b

    def __repr__ (self):
        return 'Boolean({})'.format(self._value)
                            
    def eval (self, ctxt, env):
        return VBoolean(self._value)

    
class Apply (Expression):
    def __init__ (self, fun, args):
        self._fun = fun
        self._args = args
        
    def __repr__ (self):
        return 'Apply({}, [{}])'.format(repr(self._fun),
                                        ', '.join([ repr(arg) for arg in self._args ]))

    def eval_partial (self, ctxt, env):
        f = self._fun.eval(ctxt, env)
        values = [ arg.eval(ctxt, env) for arg in self._args ]
        if isinstance(f, VPrimitive):
            return (f.apply(ctxt, values), None)
        elif isinstance(f, VFunction):
            (_, body, _) = f.value()
            new_env = f.binding_env(values)
            return (body, new_env)
        else:
            raise LispError('Cannot apply value {}'.format(f))
    
    
class If (Expression):
    def __init__ (self, cnd, thn, els):
        self._cond = cnd
        self._then = thn
        self._else = els

    def __repr__ (self):
        return 'If({}, {}, {})'.format(repr(self._cond),
                                       repr(self._then),
                                       repr(self._else))
        
    def eval_partial (self, ctxt, env):
        c = self._cond.eval(ctxt, env)
        if c.is_true():
            return (self._then, env)
        else:
            return (self._else, env)

        
class Quote (Expression):
    def __init__ (self, sexpr):
        self._sexpr = sexpr

    def __repr__ (self):
        return 'Quote({})'.format(repr(self._sexpr))

    def eval (self, ctxt, env):
        return self._sexpr.as_value()


class Lambda (Expression):
    def __init__ (self, params, expr):
        self._params = [ p.upper() for p in params ]
        self._expr = expr

    def __repr__ (self):
        return 'Lambda({}, {})'.format(self._params, repr(self._expr))
        
    def eval (self, ctxt, env):
        return VFunction(self._params, self._expr, env)

    
class LetRec (Expression):
    def __init__ (self, bindings, expr):
        self._bindings = bindings
        self._expr = expr

    def __repr__ (self):
        return 'LetRec({}, {})'.format([ (x, repr(z)) for (x, z) in self._bindings ], repr(self._expr))

    def eval_partial (self, ctxt, env):
        new_env = Environment(previous=env)
        for (n, e) in self._bindings:
            new_env.add(n, None)
        vs = [ e.eval(ctxt, new_env) for (_, e) in self._bindings ]
        for ((n, _), v) in zip(self._bindings, vs):
            new_env.add(n, v)
        return (self._expr, new_env)
            

class Do (Expression):
    def __init__ (self, exprs):
        self._exprs = exprs
        
    def __repr__ (self):
        return 'Do([{}])'.format(', '.join([ repr(arg) for arg in self._exprs ]))
        
    def eval_partial (self, ctxt, env):
        if not self._exprs:
            return (VNil(), None)
        for expr in self._exprs[:-1]:
            expr.eval(ctxt, env)
        return (self._exprs[-1], env)


class SExpression (object):
    def is_atom (self):
        return False

    def is_cons (self):
        return False

    def is_empty (self):
        return False

#    def __str__ (self):
#        return str(self)

    @staticmethod
    def from_value (v):
        if v.is_atom():
            return SAtom(str(v))
        if v.is_empty():
            return SEmpty()
        if v.is_cons():
            return SCons(SExpression.from_value(v.car()), SExpression.from_value(v.cdr()))
        raise LispError('Cannot convert {} back to s-expression'.format(repr(v)))
    
    
class SAtom (SExpression):
    def __init__ (self, s):
        self._content = s

    def is_atom (self):
        return True

    def content (self):
        return self._content

    def __repr__ (self):
        return ('SAtom({})'.format(self._content))

    def __str__ (self):
        return str(self._content)

    def _str_cdr (self):
        return ' . {})'.format(self._content)

    def match_token (self, tok):
        tok = tok.upper()
        s = self._content.upper()
        m = re.match('^{}$'.format(tok), s)
        if m:
            return m.group()
        return None

    def as_value (self):
        content = self._content
        if content[0] == '"' and content[-1] == '"':
            return VString(content[1:-1])
        if self.match_token(r'-?[0-9]+'):
            return VInteger(int(content))
        if self.match_token(r'#t'):
            return VBoolean(True)
        if self.match_token(r'#f'):
            return VBoolean(False)
        return VSymbol(content)

    def to_expression (self): 
        content = self._content
        if content[0] == '"' and content[-1] == '"':
            return String(content[1:-1])
        if self.match_token(r'-?[0-9]+'):
            return Integer(int(content))
        if self.match_token(r'#t'):
            return Boolean(True)
        if self.match_token(r'#f'):
            return Boolean(False)
        if ':' in content:
            subsymbols = content.split(':')
            return Symbol(subsymbols[-1], qualifiers=subsymbols[:-1])
        return Symbol(content)
   
    
class SCons (SExpression):
    def __init__ (self, car, cdr):
        self._car = car
        self._cdr = cdr

    def is_cons (self):
        return True

    def content (self):
        return (self._car, self._cdr)

    def __repr__ (self):
        return 'SCons({}, {})'.format(repr(self._car), repr(self._cdr))

    def __str__ (self):
        return '({}{}'.format(self._car, self._cdr._str_cdr())

    def _str_cdr (self):
        return ' {}{}'.format(self._car, self._cdr._str_cdr())

    def as_value (self):
        return VCons(self._car.as_value(), self._cdr.as_value())
            

class SEmpty (SExpression):
    def __repr__ (self):
        return 'SEmpty()'

    def is_empty (self):
        return True

    def content (self):
        return None
    
    def __str__ (self):
        return '()'

    def _str_cdr (self):
        return ')'

    def as_value (self):
        return VEmpty()



# PARSER COMBINATORS

# a parser is a function String -> Option ('a, String)

def parse_sexp_wrap (p, f):
    def parser (s):
        result = p(s)
        if not result:
            return None
        return (f(result[0]), result[1])
    return parser


def parse_token (token):
    def parser (s):
        ss = s.strip()
        m = re.match(token, ss)
        if m:
            return (m.group(), ss[m.end():])
        return None
    return parser


def parse_success (v):
    def parser (s):
        return (v, s)
    return parser



# SEXPRESSIONS parser
    
def parse_lparen (s):
    return parse_token(r'\(')(s)

def parse_rparen (s):
    return parse_token(r'\)')(s)

def parse_dot (s):
    return parse_token(r'\.')(s)

def parse_atom (s):
    p = parse_token(r"[^'()#\s]+")
    return parse_sexp_wrap(p, lambda x: SAtom(x))(s)

def parse_string (s):
    def clean (s):
        return s.replace('\\"', '"').replace('\\\\', '\\')
    # p = parse_token(r'"[^"]*"')
    p = parse_token(r'"(?:[^"\\]|\\.)*"')
    return parse_sexp_wrap(p, lambda x: SAtom(x))(s)

def parse_boolean (s):
    p = parse_token(r'#(?:t|f|T|F)')
    return parse_sexp_wrap(p, lambda x: SAtom(x))(s)

def parse_sexp (s):
    p = parse_first([parse_string,
                     parse_atom,
                     parse_boolean,
                     parse_sexp_wrap(parse_seq([parse_token(r"'"),
                                                parse_sexp]),
                                lambda x: SCons(SAtom('quote'), SCons(x[1], SEmpty()))),
                                          ##parse_sexp_string,
                     parse_sexp_wrap(parse_seq([parse_lparen,
                                           parse_sexps,
                                           parse_rparen]),
                                lambda x: x[1])])
    return p(s)
    
def parse_sexps (s):
    p = parse_first([parse_sexp_wrap(parse_seq([parse_sexp,
                                           parse_sexps]),
                                lambda x: SCons(x[0], x[1])),
                     parse_success(SEmpty())])
    return p(s)

    

# perhaps create a ParserComponent class acting as a decorator
# to have + and | and > as possible combinators?
# cf: http://tomerfiliba.com/blog/Infix-Operators/

def parse_wrap (p, f):

    def parser (s):
        result = p(s)
        if result is None:
            return None
        return f(result)

    return parser


def parse_seq (ps):

    def parser (s):
        acc_result = []
        current = s
        for p in ps:
            result = p(current)
            if result is None:
                return None
            acc_result.append(result[0])
            current = result[1]
        return (acc_result, current)

    return parser


def parse_first (ps):

    def parser (s):
        for p in ps:
            result = p(s)
            if result is not None:
                return result
        return None

    return parser


class Parser (object):
    def __init__ (self):
        self._macros = {}
        self._gensym_count = 0

    def gensym (self, prefix='gsym'):
        c = self._gensym_count
        self._gensym_count += 1
        return ' __{}_{}'.format(prefix, c)

    def parse (self, sexp):
        result = self.parse_define(sexp)
        if result:
            return ('define', result)
        result = self.parse_defun(sexp)
        if result:
            return ('defun', result)
        result = self.parse_exp(sexp)
        if result:
            return ('exp', result)
        raise LispParseError('Cannot parse {}'.format(sexp))
        
    def parse_atom (self, s):
        if not s:
            return None
        if s.is_atom():
            return s.to_expression()
        return None


    def parse_empty (self, s):
        return [] if s.is_empty() else None


    def parse_list (self, ps, tail=None):

        ptail = tail if tail else self.parse_empty

        def parser (s):
            current = s
            acc = []
            for p in ps:
                if current.is_cons():
                    (car, cdr) = current.content()
                    e = p(car)
                    if e is None:
                        return None
                    acc.append(e)
                    current = cdr
                else:
                    return None
            last = ptail(current)
            if last is None:
                return None
            if tail:
                return (acc, last)
            else:
                return acc

        return parser


    def parse_rep (self, p, tail=None):

        ptail = tail if tail else self.parse_empty

        def parser (s):
            current = s
            acc = []
            while current.is_cons():
                (car, cdr) = current.content()
                e = p(car)
                if e is None:
                    return None
                acc.append(e)
                current = cdr
            last = ptail(current)
            if last is None:
                return None
            if tail:
                return (acc, last)
            else:
                return acc

        return parser


    def parse_rep1 (self, p, tail=None):

        ptail = tail if tail else self.parse_empty

        def parser (s):
            # at least 1
            if not s.is_cons():
                return None
            (car, cdr) = s.content()
            e = p(car)
            if e is None:
                return None
            acc = [e]
            current = cdr
            while current.is_cons():
                (car, cdr) = current.content()
                e = p(car)
                if e is None:
                    return None
                acc.append(e)
                current = cdr
            last = ptail(current)
            if last is None:
                return None
            if tail:
                return (acc, last)
            else:
                return acc

        return parser
    
    
    def parse_keyword (self, kw):

        def parser (s):
            if not s:
                return None
            if s.is_atom() and s.content().upper() == kw.upper():
                return kw.upper()
            return None

        return parser


    def parse_qualified_identifier (self, s):

        char = r'A-Za-z-+/*_.?!@$'
        identifier = r'[{c}0-9]*[{c}#][{c}#0-9]*'.format(c=char)
        qidentifier = r'({id}:)?{id}'.format(id=identifier)

        if not s:
            return None
        if s.is_atom():
            m = s.match_token(identifier)
            sm = m.split(':')
            if len(sm) > 1:
                return (sm[0], sm[1])
            else:
                return m
        return None

        
    def parse_identifier (self, s):

        char = r'A-Za-z-+/*_.?!@$'
        identifier = '[{c}0-9]*[{c}#][{c}#0-9]*'.format(c=char)

        if not s:
            return None
        if s.is_atom():
            m = s.match_token(identifier)
            return m
        return None


    def parse_exp (self, s):

        p = parse_first([self.parse_atom,
                         self.parse_quote,
                         self.parse_if,
                         self.parse_lambda,
                         self.parse_do,
                         self.parse_letrec,
                         self.parse_builtin_macros,
                         self.parse_apply])
        return p(s)


    def parse_if (self, s):

        p = self.parse_list([self.parse_keyword('if'),
                             self.parse_exp,
                             self.parse_exp,
                             self.parse_exp])
        p = parse_wrap(p, lambda x: If(x[1], x[2], x[3]))
        return p(s)


    def parse_lambda (self, s):

        p = self.parse_list([self.parse_keyword('fun'),
                             self.parse_rep(self.parse_identifier)],
                            tail=self.parse_exps)
        p = parse_wrap(p, lambda x:Lambda(x[0][1], Do(x[1])))
        return p(s)

    
    def parse_do (self, s):

        p = self.parse_list([self.parse_keyword('do')],
                            tail=self.parse_exps)
        p = parse_wrap(p, lambda x: Do(x[1]))
        return p(s)


    def parse_quote (self, s):

        p = self.parse_list([self.parse_keyword('quote'),
                             lambda s: s])
        p = parse_wrap(p, lambda x: Quote(x[1]))
        return p(s)


    def parse_letrec (self, s):
        p = self.parse_list([self.parse_keyword('letrec'),
                             self.parse_rep(self.parse_binding)],
                            tail=self.parse_exps)
        p = parse_wrap(p, lambda x: LetRec(x[0][1], Do(x[1])))
        return p(s)

    
    def parse_apply (self, s):

        p = self.parse_rep1(self.parse_exp)
        p = parse_wrap(p, lambda x: Apply(x[0], x[1:]))
        return p(s)


    def parse_exps (self, s):

        p = self.parse_rep(self.parse_exp)
        return p(s)


    ############################################################
    #
    # Top level commands
    #

    def parse_define (self, s):
        p = self.parse_list([self.parse_keyword('def'),
                             self.parse_identifier,
                             self.parse_exp])
        p = parse_wrap(p, lambda x: (x[1], x[2]))
        return p(s)


    def parse_defun (self, s):
        p1 = self.parse_list([self.parse_keyword('def'),
                              self.parse_list([self.parse_identifier],
                                              tail=self.parse_rep(self.parse_identifier))],
                             tail=self.parse_exps)
        p1 = parse_wrap(p1, lambda x:(x[0][1][0][0], x[0][1][1], Do(x[1])))
        p2 = self.parse_list([self.parse_keyword('defun'),
                              self.parse_identifier,
                              self.parse_rep(self.parse_identifier)],
                             tail=self.parse_exps)
        p2 = parse_wrap(p2, lambda x:(x[0][1], x[0][2], Do(x[1])))
        p = parse_first([p1, p2])
        return p(s)


    ############################################################
    #
    # Built-in macros
    #

    def parse_builtin_macros (self, s):
        p = parse_first([self.parse_let,
                         self.parse_letstar,
                         self.parse_loop,
                         self.parse_funrec,
                         self.parse_dict,
                         self.parse_and,
                         self.parse_or])
        return p(s)
    
    def parse_binding (self, s):
        p = self.parse_list([self.parse_identifier,
                        self.parse_exp])
        p = parse_wrap(p, lambda x: (x[0], x[1]))
        return p(s)

    def parse_let (self, s):
        p = self.parse_list([self.parse_keyword('let'),
                             self.parse_rep(self.parse_binding)],
                       tail=self.parse_exps)
        p = parse_wrap(p, lambda x: self.mk_Let(x[0][1], Do(x[1])))
        return p(s)

    def parse_loop (self, s):
        p = self.parse_list([self.parse_keyword('loop'),
                             self.parse_identifier,
                             self.parse_rep(self.parse_binding)],
                       tail=self.parse_exps)
        p = parse_wrap(p, lambda x: self.mk_Loop(x[0][1], x[0][2], Do(x[1])))
        return p(s)
    
    def parse_funrec (self, s):
        p = self.parse_list([self.parse_keyword('funrec'),
                             self.parse_identifier,
                             self.parse_rep(self.parse_identifier)],
                            tail=self.parse_exps)
        p = parse_wrap(p, lambda x: self.mk_FunRec(x[0][1], x[0][2], Do(x[1])))
        return p(s)
    
    def parse_letstar (self, s):
        p = self.parse_list([self.parse_keyword('let*'),
                             self.parse_rep(self.parse_binding)],
                            tail=self.parse_exps)
        p = parse_wrap(p, lambda x: self.mk_LetStar(x[0][1], Do(x[1])))
        return p(s)

    def parse_exp_pair (self, s):
        p = self.parse_list([self.parse_exp,
                             self.parse_exp])
        p = parse_wrap(p, lambda x: (x[0], x[1]))
        return p(s)
 
    def parse_dict (self, s):
        p = self.parse_list([self.parse_keyword('dict')],
                            tail=self.parse_rep(self.parse_exp_pair))
        p = parse_wrap(p, lambda x: self.mk_Dict(x[1]))
        return p(s)
    
    def parse_and (self, s):
        p = self.parse_list([self.parse_keyword('and')],
                            tail=self.parse_exps)
        p = parse_wrap(p, lambda x: self.mk_And(x[1]))
        return p(s)

    def parse_or (self, s):
        p = self.parse_list([self.parse_keyword('or')],
                            tail=self.parse_exps)
        p = parse_wrap(p, lambda x: self.mk_Or(x[1]))
        return p(s)


    def mk_Let (self, bindings, body):
        params = [ id for (id, _) in bindings ]
        args = [ e for (_, e) in bindings ]
        return Apply(Lambda(params, body), args)

    def mk_LetStar (self, bindings, body):
        result = body
        for (id, e) in reversed(bindings):
            result = Apply(Lambda([id], result), [e])
        return result

    def mk_And (self, exps):
        if exps:
            result = exps[-1]
            for e in reversed(exps[:-1]):
                n = self.gensym()
                result = self.mk_Let([(n, e)], If(Symbol(n), result, Symbol(n)))
            return result
        return Boolean(True)

    def mk_Or (self, exps):
        if exps:
            result = exps[-1]
            for e in reversed(exps[:-1]):
                n = self.gensym()
                result = self.mk_Let([(n, e)], If(Symbol(n), Symbol(n), result))
            return result
        return Boolean(False)

    def mk_Dict (self, exps):
        exps = [Apply(Symbol('list'), [x, y]) for (x, y) in exps]
        return Apply(Symbol('make-dict'), [Apply(Symbol('list'), exps)])

    def mk_Loop (self, name, bindings, body):
        return Apply(LetRec([(name, Lambda([ n for (n, e) in bindings ], body))],
                            Symbol(name)),
                     [ e for (n, e) in bindings ])

    def mk_FunRec (self, name, params, body):
        return LetRec([(name, Lambda(params, body))], Symbol(name))


PRIMITIVES = []

def check_arg_type (name, v, pred):
    if not pred(v):
        raise LispWrongArgTypeError('Wrong argument type {} to primitive {}'.format(v, name))

def primitive(name, min, max=None):
    name = name.upper()
    def decorator(func):
        PRIMITIVES.append((name, VPrimitive(name, func, min, max)))
        return func
    return decorator


@primitive('type', 1, 1)
def prim_type (ctxt, args):
    return VSymbol(args[0].type())


@primitive('+', 0)
def prim_plus (ctxt, args):
    v = 0
    for arg in args:
        check_arg_type('+', arg, lambda v:v.is_number())
        v += arg.value()
    return VInteger(v)

@primitive('*', 0)
def prim_times (ctxt, args):
    v = 1
    for arg in args:
        check_arg_type('*', arg, lambda v:v.is_number())
        v *= arg.value()
    return VInteger(v)

@primitive('-', 1)
def prim_minus (ctxt, args):
    check_arg_type('-', args[0], lambda v:v.is_number())
    v = args[0].value()
    if args[1:]:
        for arg in args[1:]:
            check_arg_type('-', arg, lambda v:v.is_number())
            v -= arg.value()
        return VInteger(v)
    else:
        return VInteger(-v)

def _num_predicate (arg1, arg2, sym, pred):
    check_arg_type(sym, arg1, lambda v:v.is_number())
    check_arg_type(sym, arg2, lambda v:v.is_number())
    return VBoolean(pred(arg1.value(), arg2.value()))
    
@primitive('=', 2, 2)
def prim_numequal (ctxt, args):
    return _num_predicate(args[0], args[1], '=', lambda v1, v2: v1 == v2)

@primitive('<', 2, 2)
def prim_numless (ctxt, args):
    return _num_predicate(args[0], args[1], '<', lambda v1, v2: v1 < v2)

@primitive('<=', 2, 2)
def prim_numlesseq (ctxt, args):
    return _num_predicate(args[0], args[1], '<=', lambda v1, v2: v1 <= v2)

@primitive('>', 2, 2)
def prim_numgreater (ctxt, args):
    return _num_predicate(args[0], args[1], '>', lambda v1, v2: v1 > v2)

@primitive('>=', 2, 2)
def prim_numgreatereq (ctxt, args):
    return _num_predicate(args[0], args[1], '>=', lambda v1, v2: v1 >= v2)

@primitive('not', 1, 1)
def prim_not (ctxt, args):
    return VBoolean(not args[0].is_true())

@primitive('string-append', 0)
def prim_string_append (ctxt, args):
    v = ''
    for arg in args:
        check_arg_type('string-append', arg, lambda v:v.is_string())
        v += arg.value()
    return VString(v)

@primitive('string-length', 1, 1)
def prim_string_length (ctxt, args):
    check_arg_type('string-length', args[0], lambda v:v.is_string())
    return VInteger(len(args[0].value()))

@primitive('string-lower', 1, 1)
def prim_string_lower (ctxt, args):
    check_arg_type('string-lower', args[0], lambda v:v.is_string())
    return VString(args[0].value().lower())

@primitive('string-upper', 1, 1)
def prim_string_upper (ctxt, args):
    check_arg_type('string-upper', args[0], lambda v:v.is_string())
    return VString(args[0].value().upper())

@primitive('string-substring', 1, 3)
def prim_string_substring (ctxt, args):
    check_arg_type('string-substring', args[0], lambda v:v.is_string())
    if len(args) > 2:
        check_arg_type('string-substring', args[2], lambda v:v.is_number())
        end = args[2].value()
    else:
        end = len(args[0].value())
    if len(args) > 1:
        check_arg_type('string-substring', args[1], lambda v:v.is_number())
        start = args[1].value()
    else:
        start = 0
    return VString(args[0].value()[start:end])

@primitive('apply', 2, 2)
def prim_apply (ctxt, args):
    check_arg_type('apply', args[0], lambda v:v.is_function())
    check_arg_type('apply', args[1], lambda v:v.is_list())
    return args[0].apply(args[1].to_list())
    
@primitive('cons', 2, 2)
def prim_cons (ctxt, args):
    check_arg_type('cons', args[1], lambda v:v.is_list())
    return VCons(args[0], args[1])

@primitive('append', 0)
def prim_append (ctxt, args):
    v = VEmpty()
    for arg in reversed(ctxt, args):
        check_arg_type('append', arg, lambda v:v.is_list())
        curr = arg
        temp = []
        while not curr.is_empty():
            temp.append(curr.car())
            curr = curr.cdr()
        for t in reversed(temp):
            v = VCons(t, v)
    return v

@primitive('reverse', 1, 1)
def prim_reverse (ctxt, args):
    check_arg_type('reverse', args[0], lambda v:v.is_list())
    v = VEmpty()
    curr = args[0]
    while not curr.is_empty():
        v = VCons(curr.car(), v)
        curr = curr.cdr()
    return v

@primitive('first', 1, 1)
def prim_first (ctxt, args):
    check_arg_type('first', args[0], lambda v:v.is_cons())
    return args[0].car()

@primitive('rest', 1, 1)
def prim_rest (ctxt, args):
    check_arg_type('rest', args[0], lambda v:v.is_cons())
    return args[0].cdr()

@primitive('list', 0)
def prim_list (ctxt, args):
    v = VEmpty()
    for arg in reversed(ctxt, args):
        v = VCons(arg, v)
    return v

@primitive('length', 1, 1)
def prim_length (ctxt, args):
    check_arg_type('length', args[0], lambda v:v.is_list())
    count = 0
    curr = args[0]
    while not curr.is_empty():
        count += 1
        curr = curr.cdr()
    return VInteger(count)

@primitive('nth', 2, 2)
def prim_nth (ctxt, args):
    check_arg_type('nth', args[0], lambda v:v.is_list())
    check_arg_type('nth', args[1], lambda v:v.is_number())
    idx = args[1].value()
    curr = args[0]
    while not curr.is_empty():
        if idx:
            idx -= 1
            curr = curr.cdr()
        else:
            return curr.car()
    raise LispError('Index out of range of list')

@primitive('map', 2)
def prim_map (ctxt, args):
    check_arg_type('map', args[0], lambda v:v.is_function())
    for arg in args[1:]:
        check_arg_type('map', arg, lambda v:v.is_list())
    temp = []
    currs = args[1:]
    while all(curr.is_cons() for curr in currs):
        firsts = [ curr.car() for curr in currs ]
        currs = [ curr.cdr() for curr in currs ]
        temp.append(args[0].apply(firsts))
    v = VEmpty()
    for t in reversed(temp):
        v = VCons(t, v)
    return v

@primitive('filter', 2, 2)
def prim_filter (ctxt, args):
    check_arg_type('filter', args[0], lambda v:v.is_function())
    check_arg_type('filter', args[1], lambda v:v.is_list())
    temp = []
    curr = args[1]
    while not curr.is_empty():
        if args[0].apply([curr.car()]).is_true():
            temp.append(curr.car())
        curr = curr.cdr()
    v = VEmpty()
    for t in reversed(temp):
        v = VCons(t, v)
    return v

@primitive('foldr', 3, 3)
def prim_foldr (ctxt, args):
    check_arg_type('foldr', args[0], lambda v:v.is_function())
    check_arg_type('foldr', args[1], lambda v:v.is_list())
    curr = args[1]
    temp = []
    while not curr.is_empty():
        temp.append(curr.car())
        curr = curr.cdr()
    v = args[2]
    for t in reversed(temp):
        v = args[0].apply([t, v])
    return v

@primitive('foldl', 3, 3)
def prim_foldl (ctxt, args):
    check_arg_type('foldl', args[0], lambda v:v.is_function())
    check_arg_type('foldl', args[2], lambda v:v.is_list())
    curr = args[2]
    v = args[1]
    while not curr.is_empty():
        v = args[0].apply([v, curr.car()])
        curr = curr.cdr()
    return v

@primitive('eq?', 2, 2)
def prim_eqp (ctxt, args):
    return VBoolean(args[0].is_eq(args[1]))

@primitive('eql?', 2, 2)
def prim_eqlp (ctxt, args):
    return VBoolean(args[0].is_equal(args[1]))

@primitive('empty?', 1, 1)
def prim_emptyp (ctxt, args):
    return VBoolean(args[0].is_empty())
    
@primitive('cons?', 1, 1)
def prim_consp (ctxt, args):
    return VBoolean(args[0].is_cons())

@primitive('list?', 1, 1)
def prim_listp (ctxt, args):
    return VBoolean(args[0].is_list())

@primitive('number?', 1, 1)
def prim_numberp (ctxt, args):
    return VBoolean(args[0].is_number())

@primitive('boolean?', 1, 1)
def prim_booleanp (ctxt, args):
    return VBoolean(args[0].is_boolean())

@primitive('string?', 1, 1)
def prim_stringp (ctxt, args):
    return VBoolean(args[0].is_string())

@primitive('symbol?', 1, 1)
def prim_symbolp (ctxt, args):
    return VBoolean(args[0].is_symbol())

@primitive('function?', 1, 1)
def prim_functionp (ctxt, args):
    return VBoolean(args[0].is_function())

@primitive('nil?', 1, 1)
def prim_nilp (ctxt, args):
    return VBoolean(args[0].is_nil())


@primitive('ref?', 1, 1)
def prim_refp (ctxt, args):
    return VBoolean(args[0].is_reference())

@primitive('ref', 1, 1)
def prim_ref (ctxt, args):
    return VReference(args[0])

@primitive('ref-get', 1, 1)
def prim_ref_get (ctxt, args):
    check_arg_type('ref-get', args[0], lambda v: v.is_reference())
    return args[0].value()

@primitive('ref-set', 2, 2)
def prim_ref_set (ctxt, args):
    check_arg_type('ref-set', args[0], lambda v: v.is_reference())
    args[0].set_value(args[1])
    return VNil()


@primitive('dict?', 1, 1)
def prim_dictp (ctxt, args):
    return VBoolean(args[0].is_dict())
    
@primitive('make-dict', 1, 1)
def prim_make_dict (ctxt, args):
    check_arg_type('make-dict', args[0], lambda v:v.is_list())
    entries = [ tuple(v.to_list()) for v in args[0].to_list() ]
    for entry in entries:
        if len(entry) != 2:
            raise LispError('Wrong number of element in entry {}'.format(entry))
    return VDict(entries)

@primitive('dict-get', 2, 2)
def prim_dict_get (ctxt, args):
    check_arg_type('dict-get', args[0], lambda v:v.is_dict())
    check_arg_type('dict-get', args[1], lambda v:v.is_atom())
    return args[0].lookup(args[1])

@primitive('dict-update', 3, 3)
def prim_dict_update (ctxt, args):
    check_arg_type('dict-update', args[0], lambda v:v.is_dict())
    check_arg_type('dict-update', args[1], lambda v:v.is_atom())
    return args[0].update(args[1], args[2])

@primitive('dict-set!', 3, 3)
def prim_dict_set (ctxt, args):
    check_arg_type('dict-set!', args[0], lambda v:v.is_dict())
    check_arg_type('dict-set!', args[1], lambda v:v.is_atom())
    return args[0].set(args[1], args[2])


@primitive('print', 0)
def prim_print (self, ctxt, args):
    result = ' '.join([arg.display() for arg in args])
    ctxt['print'](result)
