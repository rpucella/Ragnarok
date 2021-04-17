package main

import "fmt"

func main() {
	fmt.Println("Ragnarok/go 0.1.0")
	eco := initialize()
	fmt.Print("Modules ")
	for k := range eco.modulesEnv {
		fmt.Print(k, " ")
	}
	fmt.Println()
	shell(eco)
}

