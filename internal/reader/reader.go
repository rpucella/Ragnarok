package reader

import (
	"errors"
	"regexp"
	"unicode"
	"unicode/utf8"
	"rpucella.net/ragnarok/internal/value"
	"strconv"
	"strings"
)

func trimLeftSpace(s string) string {
	// stolen from the Go standard library (first half of TrimSpace)
	// Fast path for ASCII: look for the first ASCII non-space byte
	asciiSpace := [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return strings.TrimLeftFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}
	return s[start:]
}

func trimLeftSpaceComment(s string) string {
	ss := trimLeftSpace(s)
	for len(ss) > 1 && ss[:2] == ";;" {
		// while we got leading comments
		// remove leading comment
		nl := strings.Index(ss, "\n")
		if nl < 0 {
			// no NL, so really we have nothing - bail
			return ""
		}
		// bump the content to the NL
		ss = trimLeftSpace(ss[nl:])
	}
	return ss
}	

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
	// fmt.Printf("Trying to read `%s`\n", s)
	// remove leading spaces/comments
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
