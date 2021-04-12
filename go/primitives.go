package main

func primitiveAdd(args []Value) Value {
	var result int
	for _, val := range args {
		result += val.intValue()
	}
	return &VInteger{result}
}

func primitiveMult(args []Value) Value {
	var result int = 1
	for _, val := range args {
		result *= val.intValue()
	}
	return &VInteger{result}
}
