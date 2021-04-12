package main

import "strconv"
import "strings"
import "regexp"

func readToken(token string, s string) (string, string) {
	r, _ := regexp.Compile(token)
	ss := strings.TrimSpace(s)
	match := r.FindStringIndex(ss)
	if len(match) == 0 {
		// no match
		return "", s
	} else {
		//fmt.Println("Token match", ss, match)
		return ss[:match[1]], ss[match[1]:]
	}
}

func readLP(s string) (string, string) {
	//fmt.Println("Trying to read as LP")
	return readToken(`^\(`, s)
}

func readRP(s string) (string, string) {
	//fmt.Println("Trying to read as RP")
	return readToken(`^\)`, s)
}

func readSymbol(s string) (string, string) {
	//fmt.Println("Trying to read as symbol")
	return readToken(`^[^'()#\s]+`, s)
}

func readInteger(s string) (string, string) {
	//fmt.Println("Trying to read as integer")
	return readToken(`^-?[0-9]+`, s)
}

func readList(s string) (Value, string) {
	var rest string
	var current *VCons
	var result *VCons
	var expr Value
	expr, rest = read(s)
	for expr != nil {
		if current == nil {
			result = &VCons{head: expr, tail: &VEmpty{}}
			current = result
		} else {
			temp := &VCons{head: expr, tail: current.tail}
			current.tail = temp
			current = temp
		}
		expr, rest = read(rest)
	}
	if current == nil {
		return &VEmpty{}, rest
	}
	return result, rest
}

func read(s string) (Value, string) {
	//fmt.Println("Trying to read string", s)
	var result, rest string
	result, rest = readInteger(s)
	if result != "" {
		num, _ := strconv.Atoi(result)
		return &VInteger{num}, rest
	}
	result, rest = readSymbol(s)
	if result != "" {
		return &VSymbol{result}, rest
	}
	result, rest = readLP(s)
	if result != "" {
		var expr Value
		var exprs Value
		expr, rest = read(rest)
		if expr != nil {
			exprs, rest = readList(rest)
			result, rest = readRP(rest)
			if result != "" {
				return &VCons{head: expr, tail: exprs}, rest
			}
		}
	}
	return nil, s
}
