from .lisp import *
from .interactive import INTERACTIVE

############################################################


class Engine (object):
    
    def __init__ (self):
        self._root = Environment()
        # core
        core = Environment(bindings=PRIMITIVES, previous=self._root)
        core.add('empty', VEmpty())
        core.add('nil', VNil())
        self._root.add('core', VModule(core))
        # interactive
        interactive = Environment(previous=self._root, bindings=INTERACTIVE)
        self._root.add('interactive', VModule(interactive))
        self._parser = Parser()

    def read (self, s, strict=True):
        if not s.strip():
            return None
        result = parse_sexp(s)
        if result:
            if strict and result[1].strip():
                raise LispReadError('Input past end of expression: {}'.format(result[1]))
            if strict:
                # strict = return only the one result
                return result[0]
            else:
                # otherwise, return the result and the rest of the input
                return result
        raise LispReadError('Cannot read {}'.format(s))
        
    def eval (self, ctxt, s):
        sexp = self.read(s)
        return self.eval_sexp(ctxt, sexp)

    def eval_sexp (self, ctxt, sexp):
        env = ctxt['env']
        (type, result) = self._parser.parse(sexp)
        if type == 'define':
            (name, expr) = result
            name = name.upper()
            v = expr.eval(ctxt, env)
            ctxt['def_env'].add(name, v)
            return { 'result': VNil(), 'report': ';; {}'.format(name)}
        if type == 'defun':
            (name, params, expr) = result
            params = [ p.upper() for p in params ]
            v = VFunction(params, expr, env)
            ctxt['def_env'].add(name, v)
            return { 'result': VNil(), 'report': ';; {}'.format(name)}
        if type == 'exp':
            return { 'result': result.eval(ctxt, env) }
        raise LispError('Cannot recognize top level type {}'.format(type))

    def root (self):
        return self._root

