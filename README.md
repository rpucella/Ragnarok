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

Ragnarok is a dialect of
[Lisp](https://en.wikipedia.org/wiki/Lisp_(programming_language)). Why
a new dialect as opposed to either Common Lisp or Scheme? Good
question. It may circle back to being Scheme-compatible. 

### Declarations

`(def _identifier_ _expression_)`

`(def (_identifier_ _identifier_ ...) _expression_)`


### Expressions

`(_if_ _expression1_ _expression2_ _expression3_)`

...

`(_expression1_ _expression2_ ...)` : application of function resulting from evaluating `_expression1_` to values resulting from evaluating `_expression2_`, ...


### Primitive operations

...
