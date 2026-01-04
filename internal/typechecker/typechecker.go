package typechecker

import (
	"fmt"
	"lunno/internal/parser"
)

func Check(parser *parser.Program) []error {
	checker := &Checker{
		env: newEnv(nil),
	}
	registerBuiltins(checker.env)
	for _, e := range parser.Expressions {
		checker.checkExpr(e)
	}
	return checker.errors
}

func (checker *Checker) checkExpr(expr parser.Expression) Type {
	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		return &IntType{}
	case *parser.FloatLiteral:
		return &FloatType{}
	case *parser.BooleanLiteral:
		return &BoolType{}
	case *parser.StringLiteral:
		return &StringType{}
	case *parser.CharacterLiteral:
		return &CharType{}
	case *parser.UnitLiteral:
		return &UnitType{}
	case *parser.Identifier:
		if s, ok := checker.env.get(e.Name); ok {
			return instantiate(s, checker)
		}
		checker.errors = append(checker.errors, fmt.Errorf("undefined identifier %s", e.Name))
		return checker.freshVar()
	case *parser.ListExpression:
		elem := checker.freshVar()
		subst := Subst{}
		for _, el := range e.Elements {
			t := checker.checkExpr(el)
			if err := unify(elem, t, subst); err != nil {
				checker.errors = append(checker.errors, err)
			}
		}
		return apply(&ListType{
			Element: elem,
		}, subst)
	case *parser.FunctionLiteralExpression:
		fnEnv := newEnv(checker.env)
		params := make([]Type, len(e.Parameters))
		for i, p := range e.Parameters {
			var pt Type
			if p.Type != nil {
				pt = checker.resolveType(p.Type)
			} else {
				pt = checker.freshVar()
			}
			params[i] = pt
			fnEnv.set(p.Name.Lexeme, generalize(checker.env, pt))
		}
		old := checker.env
		checker.env = fnEnv
		body := checker.checkExpr(e.Body)
		checker.env = old
		return &FunctionType{
			Parameters: params,
			Return:     body,
		}
	case *parser.VariableDeclarationExpression:
		valType := checker.checkExpr(e.Value)
		declType := checker.resolveType(e.Type)
		subst := Subst{}
		if err := unify(declType, valType, subst); err != nil {
			checker.errors = append(checker.errors, err)
		}
		generalized := generalize(checker.env, apply(declType, subst))
		checker.env.set(e.Name.Lexeme, generalized)
		return &UnitType{}
	case *parser.IfExpression:
		cond := checker.checkExpr(e.Condition)
		if _, ok := cond.(*BoolType); !ok {
			checker.errors = append(checker.errors,
				fmt.Errorf("if condition must be bool"))
		}
		t1 := checker.checkExpr(e.Then)
		t2 := checker.checkExpr(e.Else)
		subst := Subst{}
		if err := unify(t1, t2, subst); err != nil {
			checker.errors = append(checker.errors, err)
		}
		return apply(t1, subst)
	}
	return checker.freshVar()
}
