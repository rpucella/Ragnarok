package main

import (
	"rpucella.net/ragnarok/internal/primitives"
	"rpucella.net/ragnarok/internal/ragnarok"
)

func initializeShell(eco ragnarok.Ecosystem) {

	env := eco.CreateModule("shell")
	env.Update("quit", primitives.PrimQuit)
	env.Update("env", primitives.PrimEnv)
	env.Update("go", primitives.PrimGo)
	env.Update("modules", primitives.PrimModules)
	env.Update("help", primitives.PrimHelp)
	env.Update("print", primitives.PrimPrint)
	env.Update("load", primitives.PrimLoad)
	env.Update("timed-apply", primitives.PrimTimedApply)
}
