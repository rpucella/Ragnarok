package main

import (
	"rpucella.net/ragnarok/internal/primitives"
	"rpucella.net/ragnarok/internal/ragnarok"
)

func initializeString(eco ragnarok.Ecosystem) {

	env := eco.CreateModule("string")
	env.Update("append", primitives.PrimStringAppend)
	env.Update("length", primitives.PrimStringLength)
	env.Update("lower", primitives.PrimStringLower)
	env.Update("upper", primitives.PrimStringUpper)
	env.Update("substring", primitives.PrimStringSubstring)
}
