# Things to investigate

## Modules

A module is a wrapper around an environment.

    (module-lookup 'symbol)
    
    (with-modules '(foo bar baz)
       expression)

Do we have a concept of "going into a module"?

Only really makes sense if there is a definition of private definitions local to the module that cannot be accessed from 
outside the module and that we may want to play with easily.

But then we are working in a shell local to that module, that doesn't get saved as part of the source of the module.

The source is a protected part of the environment that we can't go into. 

Though we can "save" functions into the module after we've defined them, but they must bring with it everything they depend on. (Maybe)


    (source:new-module 'name)       <-- create module
    (source:delete-module 'name)    <-- drop module - complain if dangling module references?
    (source:rename-module 'name)    <-- replicate to all module names in symbols? [won't catch clever non-hygienic macros]
    
    (source:edit 'name)


A module is either a library [reusable code across projects] or an application [self-contained code base with a specific functionality]


## Primitives

Create a `primitive` primitive taking the primitive symbol name and that evaluates to a `VPrimitive`

This refers to a big dictionary of primitives

When defining the core module, you get:

    (const list (primitive 'load))
    (const cons (primitive 'cons))
    (const true #t)
    (const false #f)
    (macro and args (if args `(let ((x ,(first args))) (if x ,(and (rest args)) false)) true))
    (macro or args (if args `(let ((x ,(first args))) (if x true ,(or (rest args)))) false))
    ...
    
Each of them presumably editable - so you can add defintions to core if you so wish.

But then how do you share code?
    

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


## Persistence

Sqlite3 file with rows:

    CREATE TABLE source (
      module TEXT,
      name TEXT,
      uuid TEXT
    )
    
This schema does not allow nested modules.

Source proper is stored in files:

    source/{uuid}/code.rg
    
Maybe? Or do we persist in a different database, like DBM, with key = uuid, and value = source string?

At startup, you need to reconstruct the environment: for every UUID, load the corresponding source and place it in the right module.
