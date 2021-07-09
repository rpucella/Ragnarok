package util

import (
	"rpucella.net/ragnarok/internal/value"
)

func MinInt(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func MaxInt(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func ValueIgnore(val value.Value, err error) value.Value {
	return val
}
