package main

type Env struct {
	current  map[string]Value
	previous *Env
}

func (env *Env) lookup(name string) Value {
	val, ok := env.current[name]
	if !ok {
		panic("Boom!")
	}
	return val
}
