# Ragnarok

_Warning: this repo is pretty consistently in flux. When it stabilizes
a bit, I'll add the usual `dev` and feature branches. But for now, all
work is being done in the `master` branch._ 

## Context

I've been idly twiddling mentally with what I consider my ideal
development environment. It started with the idea that source code
should really be a persistent database of definitions, and that the
idea of putting code in files is really a presentation-layer issue. 

I intend Ragnarok to be my exploration of this idea. It sounds
grandiose, I know. I fully expect that I'll end replicating the early
[Smalltalk](https://en.wikipedia.org/wiki/Smalltalk) or [Lisp
Machine](https://en.wikipedia.org/wiki/Genera_(operating_system))
programming envronments, or more modern attempts at IDEs like
[LightTable](http://lighttable.com/). So be it. Learn by doing, right?

The programming language supported by the development will be a
variant of Lisp for various reasons. But if my idea works like I want
it to, it shouldn't be too difficult to add a python-like interface to
the language. Again, reinventing stuff smarter people have explored in
the past. Again, so be it.

Inspirations include Plan 9's [Acme](http://doc.cat-v.org/plan_9/4th_edition/papers/acme/), [mux](http://doc.cat-v.org/bell_labs/transparent_wsys/transparent_wsys.pdf), and [Oberon](https://people.inf.ethz.ch/wirth/ProjectOberon/UsingOberon.pdf)a


## Running the code

The code is currently written in Python 3.7. You'll need at least that
version of Python for things to work out of the box. It may well work
with 3.6, but not anything lower.

1. Install [`pipenv`](https://github.com/pypa/pipenv) - it should be a simple matter of running `pipx install pipenv`, but your OS might have a dedicated package

2. Run `pipenv install` to install the packages needed.

3. Run `pipenv run shell` to fire up a simple Ragnarok shell.

If you make modifications, you can run the linter for error messages
with `pipenv run lint`, and run the unit tests in `tests/` with
`pipenv run test`.


## The language

Ragnarok is a dialect of [Lisp](https://en.wikipedia.org/wiki/Lisp_(programming_language)). Whya new dialect as opposed to either Common Lisp or Scheme? Good question. It may circle back to being Scheme-compatible.

Links of interest include Scheme's [R5RS standard](https://schemers.org/Documents/Standards/R5RS/r5rs.pdf) and Paul Graham's [BEL](http://paulgraham.com/bel.html)


### Values

Primitive values include integers such as `42` or `-1`, booleans such as `#t` and `#f`, string such as `"hello world"`, symbols such as `'foo`, lists (built using `cons`), and functions (built using `fun` or function declarations).


### Declarations

**(`def` (_name_ _arg_ ...) _body_)** : Define a function _name_ with parameters _arg_, ... with expression _body_ as a body.

**(`const` _name_ _expression_)** : Define a constant _name_ with value the result of evaluating _expression_. Constants are immutable during execution,

**(`var` _name_ _expression_)** : Define a variable _name_ with value the result of evaluating _expression_. Variables are mutable during execution. 



### Special forms

**(`if` _expression1_ _expression2_ _expression3_)** : Conditional - if _expression1_ evaluates to true, evaluate _expression2_ otherwise _expression3_.

**(`let` ((_name_ _expression_) ...) _body_)** : Local declaration - evaluate _expression_ and bind it to _name_ before evaluating _body_.

**(`let*` ((_name_ _expression_) ...) _body_)**

**(`letrec` ((_name_ _expression_) ...) _body_)**

**(`loop` _loop-name_ ((_name_ _expression_) ...) _body)**

**(`fun` (_name_ ...) _body_)**

**(`funrec` _name_ (_name_ ...) _body_)**

**(`quote` _expression_)**

**(`do` _expression_ ...)**

**(`and` _expression_ ...)**

**(`or` _expression_ ...)**

**(_expression1_ _expression2_ ...)** : Application - Evaluate _expression1_ to a function, evaluate _expression2_, ... to values, then apply the function to the values.


### Primitive operations

...
