import uuid
import tempfile
import os

from .lisp import *

INTERACTIVE = []

def int_primitive(name, min, max=None):
    name = name.upper()
    def decorator(func):
        INTERACTIVE.append((name, VPrimitive(name, func, min, max)))
        return func
    return decorator


@int_primitive('quit', 0, 0)
def prim_quit (ctxt, args):
    raise LispQuit


@int_primitive('env', 0, 1)
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
            show_env(env.lookup(name)['value'].env())
        else:
            raise LispError('No module {}'.format(name))
    else:
        show_env(env)
    return VNil()


@int_primitive('module', 0, 1)
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


@int_primitive('load', 1, 1)
def prim_load (ctxt, args):
    check_arg_type('load', args[0], lambda v:v.is_string())
    filename = args[0].value()
    with open(filename, 'rt') as fp:
        content = fp.read()
        ctxt['read_file'](content)
    return VNil()


@int_primitive('new-function', 1, 1)
def prim_define (ctxt, args):
    check_arg_type('define', args[0], lambda v:v.is_symbol())
    # what if it's qualified?
    name = args[0].value()
    uid = uuid.uuid1().hex
    path = os.path.join(tempfile.gettempdir(), f'def_{uuid.uuid1().hex}')
    with open(path, 'wt') as fp:
        fp.write(f"(def ({name} )\n  'nothing\n)")
    print(path)
    os.system(f'emacs -nw {path}')
    with open(path, 'rt') as fp:
        content = fp.read()
    print('Removing file')
    os.remove(path)
    # need to first check that what is being defined is just what's being defined
    s = ctxt['shell']._engine.read(content, strict=True)
    (type, result) = ctxt['shell']._engine.parse_sexp(ctxt, s)
    if type != 'def' or result[0].upper() != name:
        raise LispError('Not defining function {}'.format(name))
    ctxt['shell']._engine.eval_parsed_sexp(ctxt, type, result, source=content)
    return VNil()



@int_primitive('new-constant', 1, 1)
def prim_const (ctxt, args):
    check_arg_type('define', args[0], lambda v:v.is_symbol())
    # what if it's qualified?
    name = args[0].value()
    uid = uuid.uuid1().hex
    path = os.path.join(tempfile.gettempdir(), f'def_{uuid.uuid1().hex}')
    with open(path, 'wt') as fp:
        fp.write(f"(const {name}\n  'nothing\n)")
    print(path)
    os.system(f'emacs -nw {path}')
    with open(path, 'rt') as fp:
        content = fp.read()
    print('Removing file')
    os.remove(path)
    # need to first check that what is being defined is just what's being defined
    s = ctxt['shell']._engine.read(content, strict=True)
    (type, result) = ctxt['shell']._engine.parse_sexp(ctxt, s)
    if type != 'const' or result[0].upper() != name:
        raise LispError('Not defining function {}'.format(name))
    ctxt['shell']._engine.eval_parsed_sexp(ctxt, type, result, source=content)
    return VNil()


@int_primitive('show-source', 1, 1)
def prim_show_source (ctxt, args):
    def show_env (env):
        all_bindings = env.bindings()
        width = max(len(b[0]) for b in all_bindings) + 1
        for b in sorted(all_bindings, key=lambda x: x[0]):
            ctxt['print'](f';; {(b[0] + " " * width)[:width]} {b[1]}')

    env = ctxt['env']
    check_arg_type('env', args[0], lambda v:v.is_symbol())
    name = args[0].value()
    content = ctxt['env'].find(name)
    if not content:
        raise LispError(f'Cannot find symbol {name} in environment')
    else:
        if content['source']:
            for row in content['source'].split('\n'):
                ctxt['print'](f';; {row}')
        else:
            ctxt['print'](';; No source available')
    return VNil()
    
