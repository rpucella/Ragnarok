package main

func min (a int, b int) int {
	if (a > b) {
		return b
	}
	return a
}

func max (a int, b int) int {
	if (a < b) {
		return b
	}
	return a
}

func valueIgnore(val Value, err error) Value {
	return val
}
