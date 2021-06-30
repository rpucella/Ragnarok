package main

import "errors"
import "fmt"

const kw_DEF string = "def"
const kw_LET string = "let"
const kw_LETSTAR string = "let*"
const kw_LETREC string = "letrec"
const kw_LOOP string = "let"
const kw_IF string = "if"
const kw_FUN string = "fn"
const kw_QUOTE string = "quote"
const kw_DO string = "do"

const kw_MACRO string = "macro"
const kw_AND string = "and"
const kw_OR string = "or"

var fresh = (func(init int) func(string)string { 
	id := init
	return func(prefix string) string {
		result := fmt.Sprintf("%s_%d", prefix, id)
		id += 1
		return result
	}
})(0)

func parseDef(sexp Value) (*Def, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	isDef := parseKeyword(kw_DEF, sexp.headValue())
	if !isDef {
		return nil, nil
	}
	next := sexp.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to def")
	}
	defBlock := next.headValue()
	if defBlock.isSymbol() {
		name := defBlock.strValue()
		next = next.tailValue()
		if !next.isCons() {
			return nil, errors.New("too few arguments to def")
		}
		value, err := parseExpr(next.headValue())
		if err != nil {
			return nil, err
		}
		if !next.tailValue().isEmpty() {
			return nil, errors.New("too many arguments to def")
		}
		return &Def{name, DEF_VALUE, nil, value}, nil
	}		
	if defBlock.isCons() {
		if !defBlock.headValue().isSymbol() { 
			return nil, errors.New("definition name not a symbol")
		}
		name := defBlock.headValue().strValue()
		params, err := parseSymbols(defBlock.tailValue())
		if err != nil {
			return nil, err
		}
		next = next.tailValue()
		if !next.isCons() {
			return nil, errors.New("too few arguments to def")
		}
		body, err := parseExpr(next.headValue())
		if err != nil {
			return nil, err
		}
		if !next.tailValue().isEmpty() {
			return nil, errors.New("too many arguments to def")
		}
		return &Def{name, DEF_FUNCTION, params, body}, nil
	}
	return nil, errors.New("malformed def")
}

func parseExpr(sexp Value) (AST, error) {
	expr := parseAtom(sexp)
	if expr != nil {
		return expr, nil
	}
	expr, err := parseQuote(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseIf(sexp)
	if err != nil  || expr != nil {
		return expr, err
	}
	expr, err = parseFunction(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseLet(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseLetStar(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseLetRec(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseDo(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseApply(sexp)
	if err != nil || expr != nil {
		return expr, err
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
		return nil, errors.New("malformed quote")
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("too many arguments to quote")
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
		return nil, errors.New("too few arguments to if")
	}
	cnd, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to if")
	}
	thn, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to if")
	}
	els, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("too many arguments to if")
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
		return nil, errors.New("too few arguments to fun")
	}
	if next.headValue().isSymbol() {
		// we need to parse as a recursive function
		// restart from scratch
		return parseRecFunction(sexp)
	}
	params, err := parseSymbols(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to fun")
	}
	body, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("too many arguments to fun")
	}
	return makeFunction(params, body), nil
}

func parseRecFunction(sexp Value) (AST, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	isFun := parseKeyword(kw_FUN, sexp.headValue())
	if !isFun {
		return nil, nil
	}
	next := sexp.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to fun")
	}
	recName := next.headValue().strValue()
	next = next.tailValue()
	params, err := parseSymbols(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to fun")
	}
	body, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("too many arguments to fun")
	}
	return makeRecFunction(recName, params, body), nil
}

func parseLet(sexp Value) (AST, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	isLet := parseKeyword(kw_LET, sexp.headValue())
	if !isLet {
		return nil, nil
	}
	next := sexp.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to let")
	}
	params, bindings, err := parseBindings(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to let")
	}
	body, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("too many arguments to let")
	}
	return makeLet(params, bindings, body), nil
}

func parseLetStar(sexp Value) (AST, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	isLet := parseKeyword(kw_LETSTAR, sexp.headValue())
	if !isLet {
		return nil, nil
	}
	next := sexp.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to let*")
	}
	params, bindings, err := parseBindings(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to let*")
	}
	body, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("too many arguments to let*")
	}
	return makeLetStar(params, bindings, body), nil
}

func parseLetRec(sexp Value) (AST, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	isLetRec := parseKeyword(kw_LETREC, sexp.headValue())
	if !isLetRec {
		return nil, nil
	}
	next := sexp.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to letrec")
	}
	names, params, bodies, err := parseFunBindings(next.headValue())
	if err != nil {
		return nil, err
	}
	next = next.tailValue()
	if !next.isCons() {
		return nil, errors.New("too few arguments to letrec")
	}
	body, err := parseExpr(next.headValue())
	if err != nil {
		return nil, err
	}
	if !next.tailValue().isEmpty() {
		return nil, errors.New("too many arguments to letrec")
	}
	return &LetRec{names, params, bodies, body}, nil
}

func parseBindings(sexp Value) ([]string, []AST, error) {
	params := make([]string, 0)
	bindings := make([]AST, 0)
	current := sexp
	for current.isCons() {
		if !current.headValue().isCons() {
			return nil, nil, errors.New("expected binding (name expr)")
		}
		if !current.headValue().headValue().isSymbol() {
			return nil, nil, errors.New("expected name in binding")
		}
		params = append(params, current.headValue().headValue().strValue())
		if !current.headValue().tailValue().isCons() {
			return nil, nil, errors.New("expected expr in binding")
		}
		if !current.headValue().tailValue().tailValue().isEmpty() {
			return nil, nil, errors.New("too many elements in binding")
		}
		binding, err := parseExpr(current.headValue().tailValue().headValue())
		if err != nil {
			return nil, nil, err
		}
		bindings = append(bindings, binding)
		current = current.tailValue()
	}
	if !current.isEmpty() {
		return nil, nil, errors.New("malformed binding list")
	}
	return params, bindings, nil
}

func parseFunBindings(sexp Value) ([]string, [][]string, []AST, error) {
	names := make([]string, 0)
	params := make([][]string, 0)
	bodies := make([]AST, 0)
	current := sexp
	for current.isCons() {
		if !current.headValue().isCons() {
			return nil, nil, nil, errors.New("expected binding (name params expr)")
		}
		if !current.headValue().headValue().isSymbol() {
			return nil, nil, nil, errors.New("expected name in binding")
		}
		names = append(names, current.headValue().headValue().strValue())
		if !current.headValue().tailValue().isCons() {
			return nil, nil, nil, errors.New("expected params in binding")
		}
		these_params, err := parseSymbols(current.headValue().tailValue().headValue())
		if err != nil {
			return nil, nil, nil, err
		}
		params = append(params, these_params)
		if !current.headValue().tailValue().tailValue().isCons() {
			return nil, nil, nil, errors.New("expected expr in binding")
		}
		if !current.headValue().tailValue().tailValue().tailValue().isEmpty() {
			return nil, nil, nil, errors.New("too many elements in binding")
		}
		body, err := parseExpr(current.headValue().tailValue().tailValue().headValue())
		if err != nil {
			return nil, nil, nil, err
		}
		bodies = append(bodies, body)
		current = current.tailValue()
	}
	if !current.isEmpty() {
		return nil, nil, nil, errors.New("malformed binding list")
	}
	return names, params, bodies, nil
}

func makeLet(params []string, bindings []AST, body AST) AST {
	return &Apply{makeFunction(params, body), bindings}
}

func makeLetStar(params []string, bindings []AST, body AST) AST {
	result := body
	for i := len(params) - 1; i >= 0; i-- {
		result = makeLet([]string{params[i]}, []AST{bindings[i]}, result)
	}
	return result
}

func makeFunction(params []string, body AST) AST {
	name := fresh("__temp")
	return &LetRec{[]string{name}, [][]string{params}, []AST{body}, &Id{name}}
}

func makeRecFunction(recName string, params []string, body AST) AST {
	return &LetRec{[]string{recName}, [][]string{params}, []AST{body}, &Id{recName}}
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
		return nil, errors.New("malformed expression list")
	}
	return args, nil
}

func parseSymbols(sexp Value) ([]string, error) {
	params := make([]string, 0)
	current := sexp
	for current.isCons() {
		if !current.headValue().isSymbol() {
			return nil, errors.New("expected symbol in list")
		}
		params = append(params, current.headValue().strValue())
		current = current.tailValue()
	}
	if !current.isEmpty() {
		return nil, errors.New("malformed symbol list")
	}
	return params, nil
}

func parseDo(sexp Value) (AST, error) {
	if !sexp.isCons() {
		return nil, nil
	}
	isDo := parseKeyword(kw_DO, sexp.headValue())
	if !isDo {
		return nil, nil
	}
	exprs, err := parseExprs(sexp.tailValue())
	if err != nil {
		return nil, err
	}
	return makeDo(exprs), nil
}

func makeDo(exprs []AST) AST {
	if len(exprs) > 0 {
		result := exprs[len(exprs) - 1]
		for i := len(exprs) - 2; i >= 0; i-- {
			result = makeLet([]string{fresh("__temp")}, []AST{exprs[i]}, result)
		}
		return result
	}
	return &Literal{&VNil{}}
}
