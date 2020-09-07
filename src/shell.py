from .lisp import *

class Engine (object):
    
    def __init__ (self):
        root = Environment()
        core = Environment(bindings=PRIMITIVES, previous=root)
        core.add('empty', VEmpty())
        core.add('nil', VNil())
        test = Environment(previous=root)
        test.add('test', VPrimitive('test', prim_test, 0, 0))
        interactive = Environment(previous=root)
        interactive.add('quit', VPrimitive('quit', prim_quit, 0, 0))
        interactive.add('module', VPrimitive('module', prim_module, 0, 1))
        interactive.add('env', VPrimitive('env', prim_env, 0, 1))
        root.add('core', VModule(core))
        root.add('test', VModule(test))
        root.add('interactive', VModule(interactive))
        self._root = root
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

    def shell (self):
        return Shell(self)

    def root (self):
        return self._root

    
class Shell:
    def __init__ (self, engine):
        self._engine = engine
        self._module = None   
        self._env = Environment(previous=engine.root())

    def prompt (self):
        name = self._module or 'scratch'
        return name.upper() + '> '

    def cont_prompt (self):
        name = self._module or 'scratch'
        return '.' * len(name) + '  '

    def current_env (self):
        if self._module:
            return self._env.lookup(self._module).env()
        else:
            return self._env

    def balance (self, str):
        state = 'normal'
        count = 0
        pos = 0
        while pos < len(str):
            if state == 'normal':
                if str[pos] == '(':
                    pos += 1
                    count += 1
                elif str[pos] == ')':
                    pos += 1
                    count -= 1
                elif str[pos] == '"':
                    pos += 1
                    state = 'string'
                else:
                    pos += 1
            elif state == 'string':
                if str[pos] == '"':
                    pos += 1
                    state = 'normal'
                elif str[pos] == '\\':
                    pos += 1
                    state = 'escape'
                elif str[pos] == '\n':
                    raise LispParseError('Unterminated string')
                else:
                    pos += 1
            elif state == 'escape':
                pos += 1
                state = 'string'
        # this will ignore inputs past the end of the first expression
        return count <= 0

    def context (self):
        return {
            'print': self.emit,
            'env': self._env,
            'set_module': self.set_module
        }

    def set_module (self, name):
        self._module = name
                    
    def process_line (self, full_input):
        try:
            sexp = self._engine.read(full_input)
            if sexp:
                env = self.current_env()
                result = self._engine.eval_sexp(self.context(), env, sexp)
                if 'report' in result:
                    self.emit(result['report'])
                self.emit_result(result['result'])
        except LispError as e:
            self.emit_error(e)

    def emit (self, s):
        print(s)

    def emit_error (self, e):
        self.emit(';; ' + str(e))

    def emit_result (self, v):
        if not v.is_nil():
            self.emit(str(v))
            
    def repl (self, on_error=None):
        '''
        A simple read-eval-print loop 
        '''
        done = False
        while not done:
            try:
                # to deal with win_unicode_console flushing problem
                full_input = ""
                pr = self.prompt()
                while True:
                    print(pr, end='')
                    sys.stdout.flush()
                    s = input()
                    full_input += s + ' '
                    if self.balance(full_input):
                        break
                    pr = self.cont_prompt()   # use continuation prompt after first iteration
                self.process_line(full_input)
            except EOFError:
                done = True
            except LispQuit:
                done = True
        self.emit(';; tada')
        

if __name__ == '__main__':
    Engine().shell().repl()
