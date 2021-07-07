package main

import (
	"fmt"
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

func main() {
	fmt.Println(banner)
	eco := NewEcosystem()
	eco.addModule("core", coreBindings())
	eco.addModule("test", testBindings())
	eco.addModule("shell", shellBindings())
	eco.addModule("config", configBindings())
	shell(eco)
}
