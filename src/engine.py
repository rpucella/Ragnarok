import re
from .lisp import *
from src import persistence, interactive


_VERSION = '0.0.1'

class Engine (object):
    
    def __init__ (self, persist=True):
        print(f';; Ragnarok Engine {_VERSION}')
        self._parser = Parser()
        self._root = Environment()
        # core
        if persist:
            core = persistence.load_module('CORE', self, self._root)
            interactive = persistence.load_module('INTERACTIVE', self, self._root)
        else:
            # obsolete, but kept for giggles
            prims = PRIMITIVES.items()
            core = Environment(bindings=prims, previous=self._root)
            core.add('empty', VEmpty())
            core.add('nil', VNil())
            # # interactive
            interactive = Environment(previous=self._root)
        self._root.add('core', VModule(core))
        self._root.add('interactive', VModule(interactive))

    def read (self, s, strict=True):
        ss = re.sub(r';[^\n]*', '', s)   # remove comments
        if not ss.strip():
            return None
        result = parse_sexp(ss)
        if result:
            if strict and result[1].strip():
                raise LispReadError('Input past end of expression: {}'.format(result[1]))
            if strict:
                # strict = return only the one result
                return result[0]
            else:
                # otherwise, return the result and the rest of the input
                return result
        raise LispReadError('Cannot read {}'.format(ss))
        
    def eval (self, ctxt, s):
        sexp = self.read(s)
        return self.eval_sexp(ctxt, sexp)

    def parse_sexp (self, ctxt, sexp):
        return self._parser.parse(ctxt, sexp)

    def eval_parsed_sexp (self, ctxt, type, result, source=None):
        env = ctxt['env']
        if type == 'var':
            (name, expr) = result
            name = name.upper()
            v = expr.eval(ctxt, env)
            ctxt['def_env'].add(name, v, mutable=True, source=source)
            return { 'result': VNil(), 'report': ';; {}'.format(name)}
        if type == 'const':
            (name, expr) = result
            name = name.upper()
            v = expr.eval(ctxt, env)
            ctxt['def_env'].add(name, v, source=source)
            return { 'result': VNil(), 'report': ';; {}'.format(name)}
        if type == 'def':
            (name, params, expr) = result
            params = [ p.upper() for p in params ]
            v = VFunction(params, expr, env)
            ctxt['def_env'].add(name, v, source=source)
            return { 'result': VNil(), 'report': ';; {}'.format(name)}
        if type == 'macro':
            (name, params, expr) = result
            params = [ p.upper() for p in params ]
            v = VFunction(params, expr, env)
            self._parser.add_macro(name, v)
            return { 'result': VNil(), 'report': ';; {}'.format(name)}
        if type == 'exp':
            return { 'result': result.eval(ctxt, env) }
        raise LispError('Cannot recognize top level type ({})'.format(type))
        
    def eval_sexp (self, ctxt, sexp, source=None):
        # need to pass an environment to the parser for evaluating macros
        (type, result) = self._parser.parse(ctxt, sexp)
        return self.eval_parsed_sexp(ctxt, type, result, source=source)
           
    def root (self):
        return self._root

