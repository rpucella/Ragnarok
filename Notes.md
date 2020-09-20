# Notes

## Modules

A module is a wrapper around an environment.

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

For now, don't allow submodules 

- Modules are at the root of the global environment
- When we allow submodules, we could define a VModule()

Modules or declarations within could be loaded lazily. The environment could contain:

    {'value': ...,       <- if None or VDelayed(), need to load + evaluate the source
     'source': uuid,
    }


## Symbols and Keywords

Symbols are always interned - always belong to a module - if you don't specify a module, they're part of the current module.

Why? Because a module can be looked up in the environment, and the module tells you where to find it.

Two symbols in two different modules are not equal.

This means that if a module FOO create a dictionary `#dict((a 10) (b 20))` then to look up the value bound to `a` from some other module you would need to write `(dict-get 'foo:a ...)`. That's a bit annoying. Though really what you would use for that would be keywords, so that the dictionary should say `#dict((:a 10) (:b 20))` and you'd look it up with `(dict-get :a ...)`. 

I think the main reason this is needed is to get `eval` to work somewhat nicely. 

There is a special module called `keyword`. Symbols in the `keyword` module evaluate to themselves, and are called keywords. They can be abbreviated `:foo`. Thus, `:foo` evaluates to `:foo`. And of course, two keywords can be compared simply based off their label, since their module is always `keyword`. 


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

