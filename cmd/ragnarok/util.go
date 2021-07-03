package main

import (
	"rpucella.net/ragnarok/internal/value"
)

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func valueIgnore(val value.Value, err error) value.Value {
	return val
}
