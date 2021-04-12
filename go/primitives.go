package main

func primitivePlus(args []Value) Value {
	var result int
	for _, val := range args {
		result += val.intValue()
	}
	return &VInteger{result}
}
