package main

const kw_DEF string = "def"
const kw_CONST string = "const"
const kw_VAR string = "var"
const kw_MACRO string = "macro"
const kw_LET string = "let"
const kw_LETS string = "let*"
const kw_LETREC string = "letrec"
const kw_LOOP string = "let"
const kw_IF string = "if"
const kw_FUN string = "fun"
const kw_FUNREC string = "rec"
const kw_DO string = "do"
const kw_QUOTE string = "quote"
const kw_DICT string = "dict"
const kw_AND string = "and"
const kw_OR string = "or"

func parseDecl(sexp Value) (string, string, AST, Value) {
	panic("Boom!")
}

func parseExpr(sexp Value) AST {

	expr := parseAtom(sexp)
	if expr != nil {
		return expr
	}

	expr = parseApply(sexp)
	if expr != nil {
		return expr
	}
	return nil
}

func parseAtom(sexp Value) AST {
	if sexp.isSymbol() {
		return &Id{sexp.strValue()}
	}
	if sexp.isAtom() {
		return &Literal{sexp}
	}
	return nil
}

func parseApply(sexp Value) AST {
	if !sexp.isCons() {
		return nil
	}
	fun := parseExpr(sexp.headValue())
	if fun == nil {
		return nil
	}
	args := parseExprs(sexp.tailValue())
	return &Apply{fun, args}
}

func parseExprs(sexp Value) []AST {
	args := make([]AST, 0)
	current := sexp
	for current.isCons() {
		next := parseExpr(current.headValue())
		if next == nil {
			return nil
		}
		args = append(args, next)
		current = current.tailValue()
	}
	// check that current is actually empty!
	return args
}

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
