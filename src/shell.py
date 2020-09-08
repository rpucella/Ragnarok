import sys

from .lisp import Environment, LispQuit, LispError, LispParseError
from .engine import Engine

class Shell:
    def __init__ (self, engine):
        self._engine = engine
        self._module = None
        self._open_modules = ['interactive', 'core']
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
            'set_module': self.set_module,
            'modules': self._open_modules
        }

    def set_module (self, name):
        self._module = name
                    
    def process_line (self, full_input):
        try:
            sexp = self._engine.read(full_input)
            if sexp:
                env = self.current_env()
                result = self._engine.eval_sexp(self.context(), sexp)
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
        self.emit(';; Ragnarok CLI shell')
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
    e = Engine()
    Shell(e).repl()
