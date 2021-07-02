package main

import (
       "rpucella.net/ragnarok/internal/lisp"
)

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

func valueIgnore(val lisp.Value, err error) lisp.Value {
	return val
}
