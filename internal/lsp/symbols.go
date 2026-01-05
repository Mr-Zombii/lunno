package lsp

import (
	"fmt"
	"lunno/internal/lexer"
	"lunno/internal/parser"
)

type SymbolKind int

var (
	lx      *lexer.Lexer
	symbols = map[string][]Symbol{}
)

const (
	FunctionSymbol SymbolKind = iota
)

type Symbol struct {
	Name string
	Kind SymbolKind
}

func updateSymbols(uri, text string, lx *lexer.Lexer) {
	_, tokens, err := lexer.Tokenize(text, uri)
	if err != nil {
		return
	}
	program, errs := parser.ParseProgram(tokens, lx)
	var symbol []Symbol
	var walk func(expr parser.Expression)
	walk = func(expr parser.Expression) {
		switch expr := expr.(type) {
		case *parser.VariableDeclarationExpression:
			symbol = append(symbol, Symbol{
				Name: expr.Name.Lexeme,
				Kind: 0,
			})
		case *parser.FunctionDeclarationExpression:
			symbol = append(symbol, Symbol{
				Name: expr.Name.Lexeme,
				Kind: 1,
			})
		case *parser.BlockExpression:
			for _, sub := range expr.Expressions {
				walk(sub)
			}
		case *parser.CallExpression:
			walk(expr.Callee)
			for _, arg := range expr.Arguments {
				walk(arg)
			}
		case *parser.IfExpression:
			walk(expr.Condition)
			walk(expr.Then)
			walk(expr.Else)
		case *parser.ListExpression:
			for _, element := range expr.Elements {
				walk(element)
			}
		}
	}
	for _, expr := range program.Expressions {
		walk(expr)
	}

	symbols[uri] = symbol
	var diagnostics []Diagnostic
	for _, e := range errs {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: 1,
			Message:  fmt.Sprintf("%s:%d:%d: %s", uri, 1, 1, e),
			Range: Range{
				Start: Position{Line: 0, Character: 0},
				End:   Position{Line: 0, Character: 1},
			},
		})
	}

	send(map[string]any{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnostics,
		},
	})
}
