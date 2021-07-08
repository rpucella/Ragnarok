package reader

import (
	"errors"
	"regexp"
	"rpucella.net/ragnarok/internal/value"
	"strconv"
	"strings"
)

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

func readSymbol(s string) (value.Value, string) {
	//fmt.Println("Trying to read as symbol")
	result, rest := readToken(`[^"'()#\s]+`, s)
	if result == "" {
		return nil, s
	}
	return value.NewSymbol(result), rest
}

func readString(s string) (value.Value, string) {
	//fmt.Println("Trying to read as symbol")
	result, rest := readToken(`"[^\n"]+"`, s)
	if result == "" {
		return nil, s
	}
	return value.NewString(result[1 : len(result)-1]), rest
}

func readInteger(s string) (value.Value, string) {
	//fmt.Println("Trying to read as integer")
	result, rest := readToken(`-?[0-9]+`, s)
	if result == "" {
		return nil, s
	}
	num, _ := strconv.Atoi(result)
	return value.NewInteger(num), rest
}

func readBoolean(s string) (value.Value, string) {
	// TODO: read all characters after # and then process
	//       or treat # as a reader macro in some way?
	result, rest := readToken(`#(?:t|T)`, s)
	if result != "" {
		return value.NewBoolean(true), rest
	}
	result, rest = readToken(`#(?:f|F)`, s)
	if result != "" {
		return value.NewBoolean(false), rest
	}
	return nil, s
}

func readList(s string) (value.Value, string, error) {
	var current *value.Cons
	var result *value.Cons
	expr, rest, err := Read(s)
	for err == nil {
		if expr == nil {
			// incomplete, so abort
			return nil, s, nil
		}
		if current == nil {
			result = value.NewCons(expr, value.NewEmpty())
			current = result
		} else {
			temp := value.NewCons(expr, current.GetTail())
			current.SetTail(temp)
			current = temp
		}
		expr, rest, err = Read(rest)
	}
	if current == nil {
		return value.NewEmpty(), rest, nil
	}
	return result, rest, nil
}

func Read(s string) (value.Value, string, error) {
	// returns the value read and any leftover string
	// returns nil + the original string if the value read is incomplete
	var resultB bool
	var rest string
	var result value.Value
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
		var expr value.Value
		expr, rest, err = Read(rest)
		if err != nil {
			return nil, s, err
		}
		if expr == nil {
			// incomplete, so abort
			return nil, s, nil
		}
		return value.NewCons(value.NewSymbol("quote"), value.NewCons(expr, value.NewEmpty())), rest, nil
	}
	resultB, rest = readLP(s)
	if resultB {
		var exprs value.Value
		exprs, rest, err = readList(rest)
		if err != nil {
			return nil, s, err
		}
		if exprs == nil {
			return nil, s, nil
		}
		resultB, rest = readRP(rest)
		if !resultB {
			// return nil, s, errors.New("missing closing parenthesis")
			// there's still stuff to be read!
			return nil, s, nil
		}
		return exprs, rest, nil
	}
	//return nil, s, nil
	return nil, s, errors.New("Cannot read input")
}
