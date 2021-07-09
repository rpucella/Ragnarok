package main

import (
	"fmt"
	"rpucella.net/ragnarok/internal/shell"
)

const banner = `
8888888b.                                                      888      
888   Y88b                                                     888      
888    888                                                     888      
888   d88P  8888b.   .d88b.  88888b.   8888b.  888d888 .d88b.  888  888 
8888888P"      "88b d88P"88b 888 "88b     "88b 888P"  d88""88b 888 .88P 
888 T88b   .d888888 888  888 888  888 .d888888 888    888  888 888888K  
888  T88b  888  888 Y88b 888 888  888 888  888 888    Y88..88P 888 "88b 
888   T88b "Y888888  "Y88888 888  888 "Y888888 888     "Y88P"  888  888 
                         888                                            
                    Y8b d88P                                            
                     "Y88P"
`

const help = `
Use (modules) to see installed modules
Use (go <module>) to navigate
Use (env) to show bindings
`

func main() {
	fmt.Print(banner)
	fmt.Print(help)
	eco := shell.NewEcosystem()
	eco.AddModule("core", shell.CoreBindings())
	eco.AddModule("test", shell.TestBindings())
	eco.AddModule("shell", shell.ShellBindings())
	eco.AddModule("config", shell.ConfigBindings())
	shell.Shell(eco)
}
