package main

/*
func parseToken(token string, s string) (string, string) {
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

func parseLP(s string) (string, string) {
	//fmt.Println("Trying to parse as LP")
	return parseToken(`^\(`, s)
}

func parseRP(s string) (string, string) {
	//fmt.Println("Trying to parse as RP")
	return parseToken(`^\)`, s)
}

func parseSymbol(s string) (string, string) {
	//fmt.Println("Trying to parse as symbol")
	return parseToken(`^[^'()#\s]+`, s)
}

func parseInteger(s string) (string, string) {
	//fmt.Println("Trying to parse as integer")
	return parseToken(`^-?[0-9]+`, s)
}

func parseASTs(s string) ([]AST, string) {
	result := make([]AST, 0, 10)
	var rest string
	var expr AST
	expr, rest = parseAST(s)
	for expr != nil {
		result = append(result, expr)
		expr, rest = parseAST(rest)
	}
	return result, rest
}

func parseAST(s string) (AST, string) {
	//fmt.Println("Trying to parse string", s)
	var result, rest string
	result, rest = parseInteger(s)
	if result != "" {
		num, _ := strconv.Atoi(result)
		return &Literal{&VInteger{num}}, rest
	}
	result, rest = parseSymbol(s)
	if result != "" {
		return &Symbol{result}, rest
	}
	result, rest = parseLP(s)
	if result != "" {
		var expr AST
		var exprs []AST
		expr, rest = parseAST(rest)
		if expr != nil {
			exprs, rest = parseASTs(rest)
			result, rest = parseRP(rest)
			if result != "" {
				return &Apply{expr, exprs}, rest
			}
		}
	}
	return nil, s
}
*/
