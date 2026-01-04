package parser

import (
	"lunno/internal/diagnostics"
	"lunno/internal/lexer"
)

func (parser *Parser) error(token lexer.Token, msg string) error {
	return diagnostics.Report(parser.lexer.Source, token.Span(), msg)
}
