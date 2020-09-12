#
# Persistence layer
#

import os
import uuid
import sqlite3

from .lisp import Environment

_PATH = 'source'
_MODULES_DB = 'source.db'


class SqliteConnect:
    def __enter__ (self):
        self._db = sqlite3.connect(os.path.join(_PATH, _MODULES_DB))
        self._curs = self._db.cursor()
        return self._curs

    def __exit__ (self, type, value, traceback):
        if type is None:
            self._db.commit()
        self._curs.close()
        self._db.close()
    

def create_modules ():
    with SqliteConnect() as c:
        c.execute('''CREATE TABLE IF NOT EXISTS source (
                       module TEXT,
                       name TEXT,
                       uuid TEXT,
                       PRIMARY KEY (module, name)
                     )''')


def get_modules ():
    with SqliteConnect() as c:
        c.execute('''SELECT DISTINCT module FROM source''')
        return [r[0] for r in c]

    
def get_names (module):
    with SqliteConnect() as c:
        c.execute('''SELECT name FROM source WHERE module = ?''', (module,))
        return [r[0] for r in c]

def get_uids (module):
    with SqliteConnect() as c:
        c.execute('''SELECT uuid FROM source WHERE module = ?''', (module,))
        return [r[0] for r in c]

    
def create_entry (module, name, source):
    # check doesn't already exists?
    uid = uuid.uuid1().hex
    filename = f's{uid}.rg'
    with SqliteConnect() as c:
        c.execute('''INSERT INTO source (module, name, uuid) VALUES (?, ?, ?)''',
                  (module, name, uid))
    with open(os.path.join(_PATH, filename), 'wt') as fp:
        fp.write(source + '\n')

        
def load_module (module, engine):
    uids = get_uids(module)
    env = Environment()
    context = {
        # we may need more...
        'env': env,
        'def_env': env
    }
    for uid in uids:
        filename = f's{uid}.rg'
        print(f'Loading {filename}...')
        with open(os.path.join(_PATH, filename), 'rt') as fp:
            src = fp.read()
        print(src)
        s = engine.read(src, strict=True)
        print(s)
        engine.eval_sexp(context, s)
    return env

    
    
