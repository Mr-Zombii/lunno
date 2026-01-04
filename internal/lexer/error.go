package lexer

import (
	"lunno/internal/diagnostics"
)

func (lexer *Lexer) error(msg string) error {
	return diagnostics.Report(
		lexer.Source,
		diagnostics.Span{
			File:   lexer.fileName,
			Line:   lexer.line,
			Column: lexer.column,
		},
		msg,
	)
}
