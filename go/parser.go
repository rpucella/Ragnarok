package main

import "errors"

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

func parseExpr(sexp Value) (AST, error) {
	expr := parseAtom(sexp)
	if expr != nil {
		return expr, nil
	}
	expr, err := parseQuote(sexp)
	if err != nil {
		return nil, err
	}
	if expr != nil {
		return expr, nil
	}
	expr, err = parseIf(sexp)
	if err != nil {
		return nil, err
	}
	if expr != nil {
		return expr, nil
	}
	expr, err = parseFunction(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseApply(sexp)
	if err != nil {
		return nil, err
	}
	if expr != nil {
		return expr, nil
	}
	return nil, nil
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

func parseKeyword(kw string, sexp Value) bool {
	if !sexp.isSymbol() {
		return false
	}
	return (sexp.strValue() == kw)
}

func parseQuote(sexp Value) (AST, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	isQ := parseKeyword(kw_QUOTE, sexp.headValue())
	if !isQ {
		return nil, nil
	}
	next := sexp.tailValue()
	if !next.isCons() {
		return nil, errors.New("Malformed quote")
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("Too many arguments to quote")
	}
	return &Quote{next.headValue()}, nil
}

func parseIf(sexp Value) (AST, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	isIf := parseKeyword(kw_IF, sexp.headValue())
	if !isIf {
		return nil, nil
	}
	next := sexp.tailValue()
	if !next.isCons() {
		return nil, errors.New("Too few arguments to if")
	}
	cnd, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("Too few arguments to if")
	}
	thn, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("Too few arguments to if")
	}
	els, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("Too many arguments to if")
	}
	return &If{cnd, thn, els}, nil
}

func parseFunction(sexp Value) (AST, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	isFun := parseKeyword(kw_FUN, sexp.headValue())
	if !isFun {
		return nil, nil
	}
	next := sexp.tailValue()
	if !next.isCons() {
		return nil, errors.New("Too few arguments to fun")
	}
	params, err := parseSymbols(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("Too few arguments to fun")
	}
	body, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("Too many arguments to fun")
	}
	return &Function{params, body}, nil
}

func parseApply(sexp Value) (AST, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	fun, err := parseExpr(sexp.headValue())
	if err != nil {
		return nil, err
	}
	if fun == nil {
		return nil, nil
	}
	args, err := parseExprs(sexp.tailValue())
	if err != nil {
		return nil, err
	}
	return &Apply{fun, args}, nil
}

func parseExprs(sexp Value) ([]AST, error) {
	args := make([]AST, 0)
	current := sexp
	for current.isCons() {
		next, err := parseExpr(current.headValue())
		if err != nil {
			return nil, err
		}
		if next == nil {
			return nil, nil
		}
		args = append(args, next)
		current = current.tailValue()
	}
	if !current.isEmpty() {
		return nil, errors.New("Malformed expression list")
	}
	return args, nil
}

func parseSymbols(sexp Value) ([]string, error) {
	if !sexp.isCons() {
		return nil, errors.New("Expected symbols list")
	}
	params := make([]string, 0)
	current := sexp
	for current.isCons() {
		if !current.headValue().isSymbol() {
			return nil, errors.New("Expected symbol in list")
		}
		params = append(params, current.headValue().strValue())
		current = current.tailValue()
	}
	if !current.isEmpty() {
		return nil, errors.New("Malformed symbol list")
	}
	return params, nil
}
