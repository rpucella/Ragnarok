# Notes

Jumble of thoughts and ideas for where to go and what to do next. Very little editing.

## Big picture

An evaluation engine with a Lisp syntax with an associated source code storage infrastructure.

Set of distinct modules (each representing a library or an application).

Modes: 

- CLI: simply the evaluator + the ability to define new definitions in the scratch/immediate environment (not saved as source code?) Can you go into an application module and make changes from the CLI?
- UI: evaluator, and the ability to visit source code

Source code is simply a string - it needs reparsing. (That lets us preserve comments and layout.)

If code requires testing/compilation, when does that phase happen? When a piece of source code is changed? Do we walk down the dependency graph and check all dependent pieces as well?

Built-in git like commit/checking/branches structure for exploring new code features and ideas?

Concurrency? A single evaluator instantiated multiple times? Or a single evaluator all shared with others? If you open multiple evaluation windows, do they communicate?

UI where you can have evaluation windows, where you can see code, where you can navigate code (click on an identifier and see its definition), where you can define watchers for variables (in scope?) or for function invocations, where you can define a band of icons for going to a particular "application" and open it up, all definable within the language. 

Differentiate between applications that read from the console versus those that that define their own UI?

See if the platform can be used to devise better testing/exception handling/error management infrastructure that doesn't pollute code.

Can we export an application as an independent app, perhaps as a Python project executable independently of Ragnarok?


## To Read

https://talks.golang.org/2012/splash.article


## Editing source code

    (source:edit <identifier>)
    (source:dependencies <identifier>)   - what depends on this (transitive closure)
    (source:references <identifier>)     - find all references (same as dependencies?)


## Modules

A module is a wrapper around an environment.

Main question: are modules flat, or hierarchical? To a first approximation, I want to say they are flat: a module is either an application, or a library of code. I don't picture Ragnarok used for describing applications that require local libraries, say. 

Do we have a concept of "going into a module"?

Only really makes sense if there is a definition of private definitions local to the module that cannot be accessed from 
outside the module and that we may want to play with easily.

But then we are working in a shell local to that module, that doesn't get saved as part of the source of the module.

The source is a protected part of the environment that we can't go into - we can't update those values.

Though we can "save" functions into the module after we've defined them, but they must bring with it everything they depend on. (Maybe)

    (source:new-module 'name)       <-- create module
    (source:delete-module 'name)    <-- drop module - complain if dangling module references?
    (source:rename-module 'name)    <-- replicate to all module names in symbols? [won't catch clever non-hygienic macros]
    
    (source:edit-module 'name)      <-- see the full code for the module?
    
Notion of a private declaration in a module - starts with _?

A module is either a library [reusable code across projects] or an application [self-contained code base with a specific functionality]

(For now, don't allow submodules.)

- Modules are at the root of the global environment
- When we allow submodules, we could define a VModule()

Modules or declarations within could be loaded lazily. The environment could contain:

    {'value': ...,       <- if None or VDelayed(), need to load + evaluate the source
     'source': uuid,
    }

The initial module when starting an evaluator is `*scratch*` à la Emacs. It doesn't get remembered and is not associated with permanent source code. Can you move definitions from a module to another? Are defintiions like files in a file system?

Does the code have access to the current module name? Maybe a global read-only identifier `**module**`? 

Are every identifiers qualified? Perhaps implicitly? Whether to show an identifier qualified or not depends on the qualifier and the current module? If in module `A` and you refer to `A:x`, that should mean `x` in the current module - if you change `A`'s name, then you change the qualifier of `A.x` as well.

Can we lock a module? Not allowing new defintions? When defining something in a module, does it mean that its definition goes in the module, or just in the environment, and you can promote definitions in the environment to be source code? 


## Symbols

A symbol is a symbol - it's only when you evaluate a symbol that you take it apart into module and name. Comparing two symbols for equality is really just comparing their underlying strings. 'a and 'core:a are distinct symbols, not matter which module you're in. But (eval 'a) looks up the value of a in the current module, while (eval 'core:a) looks up the value of a in the core module. I think this makes sense, and it avoids the need for keywords. And it lets us nicely use : for module qualification.

Why do we use symbols? Symbols allow for quick comparison. (Symbols are interned and tested for equality using pointer equality.) They can be used as keys in dictionaries, or as keywords.

When interpreting a list as source code, symbols become identifiers and special forms. That is where qualification comes in - when evaluating an identifier. 

Subtlety: when pulling the list description of source code (do you ever need to do that?), which module you're in may affect the symbols you get out. So if you're in module A, you pull out some source code for module A, identifiers referring within A will not be qualified, but if you then move to module B and look at the list you pulled and try to evaluate it there, then the references to identifiers within A will not refer to B since they were not fully qualified. One possibility is that when you pull source code you specify whether you want all identifiers fully qualified or not!


### CLisp model with symbols and keywords

(Don't make the distinction)

Symbols are always interned - always belong to a module - if you don't specify a module, they're part of the current module.

Why? Because a module can be looked up in the environment, and the module tells you where to find it.

Two symbols in two different modules are not equal.

This means that if a module FOO create a dictionary `#dict((a 10) (b 20))` then to look up the value bound to `a` from some other module you would need to write `(dict-get 'foo:a ...)`. That's a bit annoying. Though really what you would use for that would be keywords, so that the dictionary should say `#dict((:a 10) (:b 20))` and you'd look it up with `(dict-get :a ...)`. 

I think the main reason this is needed is to get `eval` to work somewhat nicely.

There is a special module called `keyword`. Symbols in the `keyword` module evaluate to themselves, and are called keywords. They can be abbreviated `:foo`. Thus, `:foo` evaluates to `:foo`. And of course, two keywords can be compared simply based off their label, since their module is always `keyword`. 


## For macro

    (for <var> <seq> <body>)  ==   (for-each (lambda (<var>) <body>) <seq>)
    
    (for <pat> <seq> <body>)  ==   (for-each (lambda (<var>) (match <var> (<pat> <body>) (else <error>))) <seq>)



## Vectors

    (vector? ...)
    (make-vector n f)
    (vector-length ...)
    (vector-get ...)
    (vector-set ...)

Helpful:

    (def (range n)
      (let loop ((curr (- n 1))
                 (result '()))
        (if (>= curr 0)
          (loop (- curr 1) (cons curr result))
          result)))
          
    (def (vector . items)
      (let ((vec (make-vector (length items) (fn (_) nil))))
        (do (for-each-iter items 
              (fn (x index) (vector-set vec index x)))
            vec)))

Converting a vector to a list:

    (vector-foldr (fn (a r) -> (cons a r)) vec '())
    
Vector map/folds

    (def (vector-map vec f) 
      (let ((n (vector-length vec))
            (new-vec (make-vector n (fn (i) (f (vector-get vec i))))))
        new-vec))
        
    (def (vector-foldr vec f base)
      (let ((n (vector-length vec)))
        (letrec ((foldr (fn (index)
                          (if (= index n)
                            base
                            (f (vector-get vec index) (foldr (+ index 1)))))))
          (foldr 0))))

    (def (vector-foldl vec f base)
      (let ((n (vector-length vec)))
        (letrec ((foldl (fn (index acc)
                          (if (= index n)
                            acc
                            (foldl (+ index 1) (f acc (vector-get vec index)))))))
          (foldl 0 base))))

    (def (vector-for-each vec f)
      (for-each (range (vector-length vec))
        (fn (index) (f (vector-get vec index)))))
        
    (def (vector-for-each vec f)
      (for index (range (vector-length vec))
        (f (vector-get vec index))))
        
    (macro (vector-for var vec . body)
      `(vector-for-each ,vec (fn (,var) (do ,@body))))
      
There's no easy way to reverse a vector:

    (def (vector-reverse vec)
      (let ((n (vector-length vec))
            (new-vec (make-vector n (fn (i) (vector-get vec (- n (+ i 1)))))))
        new-vec))

Reader macro:

    #(v1 v2 v3 v4)
    

## Declarations

    (def (name ...) ...)
    (const name ...)
    (macro (name ...) ...)
    
    (class name ...)    
    
We don't need `var`, since we can define `(const a (ref 10))`, though we could define:

     (macro (var name value)
       `(const ,name (ref ,value)))
       
Note that it means that everything is immutable in the environment

     (macro (set! name v)
       `(ref-set ,name ,v))

Also: 

    (fn (...) ...)   for functions
    (fn args ...)    for varargs functions

## Configuration

Module `config` with a few entries like:

    (const editor "emacs -nw {}")
    
    
## First app

Bookmarking app - kept where? What interface?



## Mutability

(For now, in a module, everything must be immutable - can only define within in a module using source commands such as `source:new-function`
and `source-new-constant`.)

Items in the environment may be mutable or immutable.

A function or a constant is immutable -- can only be modified by changing the source and re-evaluating.

    (def (f x) (+ x 1))
    
    (const pi 3.141591)

Variables are mutable - they can be changed by running code.

    (var x 10)
    
How to update: 

    (set! x 10) 

What happens if you re-declare something? Update if in the same environment layer, hide if in a different environment

In scratch, doesn't really matter. Just update the existing definition. 


## Chaining

The operation:

    (chain v (f1 args1) (f2 args2) (f3 args3) ...)

is short for

    (f3 (f2 (f1 v args1) args2) args3)
    
So that for instance:

    (string-split (string-upper (string-strip s)) " ")
    
Can be written:

    (chain s (string-strip) (string-upper) (string-split " "))
    
Note that `chain` must be a macro:

    (macro (chain init . apps) 
      (if apps
        `(,(first (first apps)) (chain ,init ,@(rest apps)) ,@(rest (first apps)))
        init))


## Classes and objects

    (class point (ix iy)
      (field x ix)      ;; = (const x (ref ix))
      (field y iy)      ;; = (const y (ref iy))
      (def (move dx dy)
        (do
          (field-set this x (+ (field-get this x) dx))
          (field-set this y (+ (field-get this y) dy))))
      (def (distance)
        (sqrt (+ (square (field-get this x)) (square (field-get this y))))))

Note that a ref cell is really a one-element vector. ANYWAY.

To create an instance:

    (const pt (point 10 20))
    
To access a method, use

    (move pt 5 -5)
    
To access a field, use

    (field-get pt x)
    
To update a field, use

    (field-set pt x v)
    
Core functions:

    (lookup-field obj 'name)
    (lookup-method obj 'name)
    
So that

    (macro (field-get obj name)
      `(ref-get (lookup-field ,obj (quote ,name))))
      
    (macro (field-set obj name val)
      `(ref-set (lookup-field ,obj (quote ,name) ,val)))
      
For methods, it's part of the application evaluation - 

    (f-exp arg1-exp arg2-exp ...)
    
First evaluate `arg1-exp` to a value `v` - if `v` is an object and `f-exp` is a symbol, check if `v` has a method called `f-exp`, and if so, lookup the method, finish evaluating arguments, and apply the method to the arguments (passing `v` as `this`). Otherwise, evaluate `f-exp` using the normal evaluation rules.

## UI

Possibly use [PyQt](https://build-system.fman.io/pyqt5-tutorial) and packaged with [fbs](https://github.com/mherrmann/fbs-tutorial).

See also QML

## Generic functions

We have (string-append ...)  and (append ...) for lists, but we could also write (append ...) as a generic
function that depending on the type of arguments appends lists or strings or vectors or whatnot.

(In fact, you can append anything that is sequence-like, suggesting the obvious subtyping)

Do we want a genuine object system?


## Reader macros

    #xyz
    #xyz(args)
    #primitive("name")
    #dict((10 20) (30 40) (50 60))
    

## Generic getters/updaters for data structures

Define generic getters and updater for various interesting operations:

    (get <obj> <args...>)
    (update <obj> <args...> <value>)
    
This works for dictionaries, references, lists.

Data structures are immutable - update reconstructs a new object.

To make a mutable dictionary, you need to make values references.

    (get <reference>)
    (set <reference> <value>)

If we have a dictionary d, we can read key k that holds a mutable value with:

    (get (get d k))
    
which can be written:

    (get! d k)
    
and write to it with:

    (set (get d k) v)
    
Or more simply:

    (set! d k v)
    
where:

    (get! obj arg ...) = (get (get obj arg ...))
    (set! obj args ... val) = (set (get obj arg ...) v)


## Generics

Cf Practical Common Lisp - https://gigamonkeys.com/book/object-reorientation-generic-functions.html


## Loop macro from CL

Cf Practical Common Lisp - https://gigamonkeys.com/book/macros-standard-control-constructs.html


## Condition system from CL

Cf Practical Common Lisp - https://gigamonkeys.com/book/beyond-exception-handling-conditions-and-restarts.html

