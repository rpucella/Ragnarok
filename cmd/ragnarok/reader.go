package main

import "strconv"
import "strings"
import "regexp"
import "errors"

func readToken(token string, s string) (string, string) {
	r, _ := regexp.Compile(`^` + token)
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

func readChar(c byte, s string) (bool, string) {
	ss := strings.TrimSpace(s)
	if len(ss) > 0 && ss[0] == c {
		return true, ss[1:]
	}
	return false, s
}

// maybe these should return Values, with simply a Boolean for "token checkers"
// like LP or RP?

func readLP(s string) (bool, string) {
	//fmt.Println("Trying to read as LP")
	return readChar('(', s)
}

func readRP(s string) (bool, string) {
	//fmt.Println("Trying to read as RP")
	return readChar(')', s)
}

func readQuote(s string) (bool, string) {
	return readChar('\'', s)
}

func readSymbol(s string) (Value, string) {
	//fmt.Println("Trying to read as symbol")
	result, rest := readToken(`[^"'()#\s]+`, s)
	if result == "" {
		return nil, s
	}
	return &VSymbol{result}, rest
}

func readString(s string) (Value, string) {
	//fmt.Println("Trying to read as symbol")
	result, rest := readToken(`"[^\n"]+"`, s)
	if result == "" {
		return nil, s
	}
	return &VString{result[1:len(result) - 1]}, rest
}

func readInteger(s string) (Value, string) {
	//fmt.Println("Trying to read as integer")
	result, rest := readToken(`-?[0-9]+`, s)
	if result == "" {
		return nil, s
	}
	num, _ := strconv.Atoi(result)
	return &VInteger{num}, rest
}

func readBoolean(s string) (Value, string) {
	// TODO: read all characters after # and then process
	//       or treat # as a reader macro in some way?
	result, rest := readToken(`#(?:t|T)`, s)
	if result != "" {
		return &VBoolean{true}, rest
	}
	result, rest = readToken(`#(?:f|F)`, s)
	if result != "" {
		return &VBoolean{false}, rest
	}
	return nil, s
}

func readList(s string) (Value, string, error) {
	var current *VCons
	var result *VCons
	expr, rest, err := read(s)
	for err == nil {
		if current == nil {
			result = &VCons{head: expr, tail: &VEmpty{}}
			current = result
		} else {
			temp := &VCons{head: expr, tail: current.tail}
			current.tail = temp
			current = temp
		}
		expr, rest, err = read(rest)
	}
	if current == nil {
		return &VEmpty{}, rest, nil
	}
	return result, rest, nil
}

func read(s string) (Value, string, error) {
	//fmt.Println("Trying to read string", s)
	var resultB bool
	var rest string
	var result Value
	var err error
	result, rest = readInteger(s)
	if result != nil {
		return result, rest, nil
	}
	result, rest = readSymbol(s)
	if result != nil {
		return result, rest, nil
	}
	result, rest = readString(s)
	if result != nil {
		return result, rest, nil
	}
	result, rest = readBoolean(s)
	if result != nil {
		return result, rest, nil
	}
	resultB, rest = readQuote(s)
	if resultB {
		var expr Value
		expr, rest, err = read(rest)
		if err != nil {
			return nil, s, err
		}
		return &VCons{head: &VSymbol{"quote"}, tail: &VCons{head: expr, tail: &VEmpty{}}}, rest, nil
	}
	resultB, rest = readLP(s)
	if resultB {
		var exprs Value
		exprs, rest, err = readList(rest)
		if err != nil {
			return nil, s, err
		}
		if exprs == nil {
			return nil, s, nil
		}
		resultB, rest = readRP(rest)
		if !resultB {
			return nil, s, errors.New("missing closing parenthesis")
		}
		return exprs, rest, nil
	}
	//return nil, s, nil
	return nil, s, errors.New("Cannot read input")
}
