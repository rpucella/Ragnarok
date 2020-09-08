from .lisp import *

def prim_env (ctxt, args):
    def show_env (env):
        all_bindings = env.bindings()
        width = max(len(b[0]) for b in all_bindings) + 1
        for b in sorted(all_bindings, key=lambda x: x[0]):
            ctxt['print'](f';; {(b[0] + " " * width)[:width]} {b[1]}')

    env = ctxt['env']
    if len(args) > 0:
        check_arg_type('env', args[0], lambda v:v.is_symbol())
        name = args[0].value().upper()
        if name == 'SCRATCH':
            show_env(env)
        elif name in env.modules():
            show_env(env.lookup(name).env())
        else:
            raise LispError('No module {}'.format(name))
    else:
        show_env(env)
    return VNil()

        
def prim_module (ctxt, args):
    if len(args) > 0:
        check_arg_type('module', args[0], lambda v:v.is_symbol())
        name = args[0].value().upper()
        if name == 'SCRATCH':
            ctxt['set_module'](None)
        elif name in ctxt['env'].modules():
            ctxt['set_module'](name)
        else:
            raise LispError('No module {}'.format(name))
    else:
        for name in ctxt['env'].modules():
            ctxt['print'](';; ' + name)
    return VNil()
        

def prim_quit (ctxt, args):
    raise LispQuit


_INTERACTIVE = [
    ('quit', VPrimitive('quit', prim_quit, 0, 0)),
    ('module', VPrimitive('module', prim_module, 0, 1)),
    ('env', VPrimitive('env', prim_env, 0, 1))
]

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
        interactive = Environment(previous=self._root, bindings=_INTERACTIVE)
        self._root.add('interactive', VModule(interactive))
        self._parser = Parser()

    def read (self, s):
        if not s.strip():
            return None
        result = parse_sexp(s)
        if result:
            return result[0]
        raise LispReadError('Cannot read {}'.format(s))
        
    def eval (self, ctxt, s):
        sexp = self.read(s)
        return self.eval_sexp(ctxt, sexp)

    def eval_sexp (self, ctxt, env, sexp):
        (type, result) = self._parser.parse(sexp)
        if type == 'define':
            (name, expr) = result
            name = name.upper()
            v = expr.eval(ctxt, env)
            env.add(name, v)
            return { 'result': VNil(), 'report': ';; {}'.format(name)}
        if type == 'defun':
            (name, params, expr) = result
            params = [ p.upper() for p in params ]
            v = VFunction(params, expr, env)
            env.add(name, v)
            return { 'result': VNil(), 'report': ';; {}'.format(name)}
        if type == 'exp':
            return { 'result': result.eval(ctxt, env) }
        raise LispError('Cannot recognize top level type {}'.format(type))

    def root (self):
        return self._root

