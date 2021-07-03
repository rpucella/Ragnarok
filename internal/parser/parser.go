package parser

import (
	"errors"
	"fmt"
	"rpucella.net/ragnarok/internal/evaluator"
	"rpucella.net/ragnarok/internal/value"
)

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

var fresh = (func(init int) func(string) string {
	id := init
	return func(prefix string) string {
		result := fmt.Sprintf("%s_%d", prefix, id)
		id += 1
		return result
	}
})(0)

func ParseDef(sexp value.Value) (*evaluator.Def, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	isDef := parseKeyword(kw_DEF, sexp.GetHead())
	if !isDef {
		return nil, nil
	}
	next := sexp.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to def")
	}
	defBlock := next.GetHead()
	if value.IsSymbol(defBlock) {
		name := defBlock.GetString()
		next = next.GetTail()
		if !value.IsCons(next) {
			return nil, errors.New("too few arguments to def")
		}
		v, err := ParseExpr(next.GetHead())
		if err != nil {
			return nil, err
		}
		if !value.IsEmpty(next.GetTail()) {
			return nil, errors.New("too many arguments to def")
		}
		return evaluator.NewDef(name, evaluator.DEF_VALUE, nil, v), nil
	}
	if value.IsCons(defBlock) {
		if !value.IsSymbol(defBlock.GetHead()) {
			return nil, errors.New("definition name not a symbol")
		}
		name := defBlock.GetHead().GetString()
		params, err := parseSymbols(defBlock.GetTail())
		if err != nil {
			return nil, err
		}
		next = next.GetTail()
		if !value.IsCons(next) {
			return nil, errors.New("too few arguments to def")
		}
		body, err := ParseExpr(next.GetHead())
		if err != nil {
			return nil, err
		}
		if !value.IsEmpty(next.GetTail()) {
			return nil, errors.New("too many arguments to def")
		}
		return evaluator.NewDef(name, evaluator.DEF_FUNCTION, params, body), nil
	}
	return nil, errors.New("malformed def")
}

func ParseExpr(sexp value.Value) (evaluator.AST, error) {
	expr := parseAtom(sexp)
	if expr != nil {
		return expr, nil
	}
	expr, err := parseQuote(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseIf(sexp)
	if err != nil || expr != nil {
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

func parseAtom(sexp value.Value) evaluator.AST {
	if value.IsSymbol(sexp) {
		return evaluator.NewId(sexp.GetString())
	}
	if value.IsAtom(sexp) {
		return evaluator.NewLiteral(sexp)
	}
	return nil
}

func parseKeyword(kw string, sexp value.Value) bool {
	if !value.IsSymbol(sexp) {
		return false
	}
	return (sexp.GetString() == kw)
}

func parseQuote(sexp value.Value) (evaluator.AST, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	isQ := parseKeyword(kw_QUOTE, sexp.GetHead())
	if !isQ {
		return nil, nil
	}
	next := sexp.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("malformed quote")
	}
	if !value.IsEmpty(next.GetTail()) {
		return nil, errors.New("too many arguments to quote")
	}
	return evaluator.NewQuote(next.GetHead()), nil
}

func parseIf(sexp value.Value) (evaluator.AST, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	isIf := parseKeyword(kw_IF, sexp.GetHead())
	if !isIf {
		return nil, nil
	}
	next := sexp.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to if")
	}
	cnd, err := ParseExpr(next.GetHead())
	if err != nil {
		return nil, err
	}
	next = next.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to if")
	}
	thn, err := ParseExpr(next.GetHead())
	if err != nil {
		return nil, err
	}
	next = next.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to if")
	}
	els, err := ParseExpr(next.GetHead())
	if err != nil {
		return nil, err
	}
	if !value.IsEmpty(next.GetTail()) {
		return nil, errors.New("too many arguments to if")
	}
	return evaluator.NewIf(cnd, thn, els), nil
}

func parseFunction(sexp value.Value) (evaluator.AST, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	isFun := parseKeyword(kw_FUN, sexp.GetHead())
	if !isFun {
		return nil, nil
	}
	next := sexp.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to fun")
	}
	if value.IsSymbol(next.GetHead()) {
		// we need to parse as a recursive function
		// restart from scratch
		return parseRecFunction(sexp)
	}
	params, err := parseSymbols(next.GetHead())
	if err != nil {
		return nil, err
	}
	next = next.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to fun")
	}
	body, err := ParseExpr(next.GetHead())
	if err != nil {
		return nil, err
	}
	if !value.IsEmpty(next.GetTail()) {
		return nil, errors.New("too many arguments to fun")
	}
	return MakeFunction(params, body), nil
}

func parseRecFunction(sexp value.Value) (evaluator.AST, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	isFun := parseKeyword(kw_FUN, sexp.GetHead())
	if !isFun {
		return nil, nil
	}
	next := sexp.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to fun")
	}
	recName := next.GetHead().GetString()
	next = next.GetTail()
	params, err := parseSymbols(next.GetHead())
	if err != nil {
		return nil, err
	}
	next = next.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to fun")
	}
	body, err := ParseExpr(next.GetHead())
	if err != nil {
		return nil, err
	}
	if !value.IsEmpty(next.GetTail()) {
		return nil, errors.New("too many arguments to fun")
	}
	return makeRecFunction(recName, params, body), nil
}

func parseLet(sexp value.Value) (evaluator.AST, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	isLet := parseKeyword(kw_LET, sexp.GetHead())
	if !isLet {
		return nil, nil
	}
	next := sexp.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to let")
	}
	params, bindings, err := parseBindings(next.GetHead())
	if err != nil {
		return nil, err
	}
	next = next.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to let")
	}
	body, err := ParseExpr(next.GetHead())
	if err != nil {
		return nil, err
	}
	if !value.IsEmpty(next.GetTail()) {
		return nil, errors.New("too many arguments to let")
	}
	return makeLet(params, bindings, body), nil
}

func parseLetStar(sexp value.Value) (evaluator.AST, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	isLet := parseKeyword(kw_LETSTAR, sexp.GetHead())
	if !isLet {
		return nil, nil
	}
	next := sexp.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to let*")
	}
	params, bindings, err := parseBindings(next.GetHead())
	if err != nil {
		return nil, err
	}
	next = next.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to let*")
	}
	body, err := ParseExpr(next.GetHead())
	if err != nil {
		return nil, err
	}
	if !value.IsEmpty(next.GetTail()) {
		return nil, errors.New("too many arguments to let*")
	}
	return makeLetStar(params, bindings, body), nil
}

func parseLetRec(sexp value.Value) (evaluator.AST, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	isLetRec := parseKeyword(kw_LETREC, sexp.GetHead())
	if !isLetRec {
		return nil, nil
	}
	next := sexp.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to letrec")
	}
	names, params, bodies, err := parseFunBindings(next.GetHead())
	if err != nil {
		return nil, err
	}
	next = next.GetTail()
	if !value.IsCons(next) {
		return nil, errors.New("too few arguments to letrec")
	}
	body, err := ParseExpr(next.GetHead())
	if err != nil {
		return nil, err
	}
	if !value.IsEmpty(next.GetTail()) {
		return nil, errors.New("too many arguments to letrec")
	}
	return evaluator.NewLetRec(names, params, bodies, body), nil
}

func parseBindings(sexp value.Value) ([]string, []evaluator.AST, error) {
	params := make([]string, 0)
	bindings := make([]evaluator.AST, 0)
	current := sexp
	for value.IsCons(current) {
		if !value.IsCons(current.GetHead()) {
			return nil, nil, errors.New("expected binding (name expr)")
		}
		if !value.IsSymbol(current.GetHead().GetHead()) {
			return nil, nil, errors.New("expected name in binding")
		}
		params = append(params, current.GetHead().GetHead().GetString())
		if !value.IsCons(current.GetHead().GetTail()) {
			return nil, nil, errors.New("expected expr in binding")
		}
		if !value.IsEmpty(current.GetHead().GetTail().GetTail()) {
			return nil, nil, errors.New("too many elements in binding")
		}
		binding, err := ParseExpr(current.GetHead().GetTail().GetHead())
		if err != nil {
			return nil, nil, err
		}
		bindings = append(bindings, binding)
		current = current.GetTail()
	}
	if !value.IsEmpty(current) {
		return nil, nil, errors.New("malformed binding list")
	}
	return params, bindings, nil
}

func parseFunBindings(sexp value.Value) ([]string, [][]string, []evaluator.AST, error) {
	names := make([]string, 0)
	params := make([][]string, 0)
	bodies := make([]evaluator.AST, 0)
	current := sexp
	for value.IsCons(current) {
		if !value.IsCons(current.GetHead()) {
			return nil, nil, nil, errors.New("expected binding (name params expr)")
		}
		if !value.IsSymbol(current.GetHead().GetHead()) {
			return nil, nil, nil, errors.New("expected name in binding")
		}
		names = append(names, current.GetHead().GetHead().GetString())
		if !value.IsCons(current.GetHead().GetTail()) {
			return nil, nil, nil, errors.New("expected params in binding")
		}
		these_params, err := parseSymbols(current.GetHead().GetTail().GetHead())
		if err != nil {
			return nil, nil, nil, err
		}
		params = append(params, these_params)
		if !value.IsCons(current.GetHead().GetTail().GetTail()) {
			return nil, nil, nil, errors.New("expected expr in binding")
		}
		if !value.IsEmpty(current.GetHead().GetTail().GetTail().GetTail()) {
			return nil, nil, nil, errors.New("too many elements in binding")
		}
		body, err := ParseExpr(current.GetHead().GetTail().GetTail().GetHead())
		if err != nil {
			return nil, nil, nil, err
		}
		bodies = append(bodies, body)
		current = current.GetTail()
	}
	if !value.IsEmpty(current) {
		return nil, nil, nil, errors.New("malformed binding list")
	}
	return names, params, bodies, nil
}

func makeLet(params []string, bindings []evaluator.AST, body evaluator.AST) evaluator.AST {
	return evaluator.NewApply(MakeFunction(params, body), bindings)
}

func makeLetStar(params []string, bindings []evaluator.AST, body evaluator.AST) evaluator.AST {
	result := body
	for i := len(params) - 1; i >= 0; i-- {
		result = makeLet([]string{params[i]}, []evaluator.AST{bindings[i]}, result)
	}
	return result
}

func MakeFunction(params []string, body evaluator.AST) evaluator.AST {
	name := fresh("__temp")
	return evaluator.NewLetRec([]string{name}, [][]string{params}, []evaluator.AST{body}, evaluator.NewId(name))
}

func makeRecFunction(recName string, params []string, body evaluator.AST) evaluator.AST {
	return evaluator.NewLetRec([]string{recName}, [][]string{params}, []evaluator.AST{body}, evaluator.NewId(recName))
}

func parseApply(sexp value.Value) (evaluator.AST, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	fun, err := ParseExpr(sexp.GetHead())
	if err != nil {
		return nil, err
	}
	if fun == nil {
		return nil, nil
	}
	args, err := parseExprs(sexp.GetTail())
	if err != nil {
		return nil, err
	}
	return evaluator.NewApply(fun, args), nil
}

func parseExprs(sexp value.Value) ([]evaluator.AST, error) {
	args := make([]evaluator.AST, 0)
	current := sexp
	for value.IsCons(current) {
		next, err := ParseExpr(current.GetHead())
		if err != nil {
			return nil, err
		}
		if next == nil {
			return nil, nil
		}
		args = append(args, next)
		current = current.GetTail()
	}
	if !value.IsEmpty(current) {
		return nil, errors.New("malformed expression list")
	}
	return args, nil
}

func parseSymbols(sexp value.Value) ([]string, error) {
	params := make([]string, 0)
	current := sexp
	for value.IsCons(current) {
		if !value.IsSymbol(current.GetHead()) {
			return nil, errors.New("expected symbol in list")
		}
		params = append(params, current.GetHead().GetString())
		current = current.GetTail()
	}
	if !value.IsEmpty(current) {
		return nil, errors.New("malformed symbol list")
	}
	return params, nil
}

func parseDo(sexp value.Value) (evaluator.AST, error) {
	if !value.IsCons(sexp) {
		return nil, nil
	}
	isDo := parseKeyword(kw_DO, sexp.GetHead())
	if !isDo {
		return nil, nil
	}
	exprs, err := parseExprs(sexp.GetTail())
	if err != nil {
		return nil, err
	}
	return makeDo(exprs), nil
}

func makeDo(exprs []evaluator.AST) evaluator.AST {
	if len(exprs) > 0 {
		result := exprs[len(exprs)-1]
		for i := len(exprs) - 2; i >= 0; i-- {
			result = makeLet([]string{fresh("__temp")}, []evaluator.AST{exprs[i]}, result)
		}
		return result
	}
	return evaluator.NewLiteral(&value.VNil{})
}
