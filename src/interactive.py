from .lisp import *


def prim_quit (ctxt, args):
    raise LispQuit


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


def prim_load (ctxt, args):
    check_arg_type('load', args[0], lambda v:v.is_string())
    filename = args[0].value()
    with open(filename, 'rt') as fp:
        content = fp.read()
        ctxt['read_file'](content)
    return VNil()


INTERACTIVE = [
    ('quit', VPrimitive('quit', prim_quit, 0, 0)),
    ('module', VPrimitive('module', prim_module, 0, 1)),
    ('env', VPrimitive('env', prim_env, 0, 1)),
    ('load', VPrimitive('load', prim_load, 1, 1))
]
