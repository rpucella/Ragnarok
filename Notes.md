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

Though we can "save" functions into the module after we've defined them, but they must bring with it everything they depend on.

(Maybe)


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

