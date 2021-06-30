import pytest
from src import shell
from src import engine

def test_shell_balance ():
    e = engine.Engine()
    s = shell.Shell(e)
    assert s.balance('hello')
    assert s.balance('hello ()')
    assert s.balance('()')
    assert s.balance('(1 2 (3 4) 5)')
    assert s.balance('(()()())')
    assert s.balance('())')
    assert s.balance('(1 2 3 (4 5) 6))()')
    assert s.balance(u'(\u00e9 \u00ea 3 (4 5) 6))()')
    assert not s.balance('(')
    assert not s.balance('hello (')
    assert not s.balance('(()')
    assert not s.balance('( 1 2 (4)')
    assert not s.balance('( 1 2 (()(()((')
    
