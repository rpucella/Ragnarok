package main

import (
	"rpucella.net/ragnarok/internal/value"
)

func vref(v value.Value) value.Value {
	return value.NewReference(v)
}

func vstr(s string) value.Value {
	return value.NewString(s)
}

func vsymlist(list []string) value.Value {
	var curr value.Value = value.NewEmpty()
	for i := len(list) - 1; i >= 0; i-- {
		curr = value.NewCons(value.NewSymbol(list[i]), curr)
	}
	return curr
}
