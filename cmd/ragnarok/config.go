package main

import (
	"rpucella.net/ragnarok/internal/ragnarok"
)

var defaultLookupPath = vref(vsymlist([]string{"shell", "core"}))

var defaultEditor = vref(vstr("emacs"))

func initializeConfig(eco ragnarok.Ecosystem) {

	config := eco.CreateModule("config")
	config.Update("lookup-path", defaultLookupPath)
	config.Update("editor", defaultEditor)
}
