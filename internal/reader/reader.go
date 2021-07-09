package reader

import (
	"errors"
	"regexp"
	"rpucella.net/ragnarok/internal/value"
	"strconv"
)

func readToken(token string, s string) (string, string) {
	r, _ := regexp.Compile(`^` + token)
	match := r.FindStringIndex(s)
	if len(match) == 0 {
		// no match
		return "", s
	} else {
		//fmt.Println("Token match", s, match)
		return s[:match[1]], s[match[1]:]
	}
}

func readChar(c byte, s string) (bool, string) {
	if len(s) > 0 && s[0] == c {
		return true, s[1:]
	}
	return false, s
}

func readLP(s string) (bool, string) {
	//fmt.Println("Trying to read as LP")
	return readChar('(', s)
}

func readRP(s string) (bool, string) {
	///fmt.Println("Trying to read as RP")
	// Trimming needed since readRP is usually not
	// accessed via Read().
	ss := trimLeftSpaceComment(s)
	return readChar(')', ss)
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
	currRest := s
	for {
		// are we done?
		resultB, rest := readRP(currRest)
		if resultB {
			if (current == nil) {
				return value.NewEmpty(), rest, nil
			}
			return result, rest, nil
		}
		// we aren't done
		expr, rest, err := Read(currRest)
		if err != nil {
			return nil, s, err
		}
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
		currRest = rest
	}
}

func Read(s string) (value.Value, string, error) {
	// Returns the value read and any leftover string.
	// Returns nil + the original string if the value read is incomplete.
	var rest string
	var result value.Value
	var err error
	///fmt.Printf("Trying to read `%s`\n", s)
	// Remove leading spaces/comments.
	s = trimLeftSpaceComment(s)
	if s == "" {
		// we're done...
		return nil, "", nil
	}
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
	resultB, rest := readQuote(s)
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
		return exprs, rest, nil
	}
	//return nil, s, nil
	return nil, s, errors.New("Cannot read input")
}
