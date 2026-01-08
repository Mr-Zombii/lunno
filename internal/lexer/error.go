package lexer

import "lunno/internal/diagnostics"

func (lexer *Lexer) errorToken(msg, lex string) Token {
	_ = diagnostics.Report(
		lexer.Source,
		diagnostics.Span{
			File:   lexer.fileName,
			Line:   lexer.line,
			Column: lexer.column,
		},
		msg,
	)

	return lexer.makeToken(Illegal, lex)
}

func (lexer *Lexer) errorAt(msg, lex string, line, column uint16) Token {
	_ = diagnostics.Report(
		lexer.Source,
		diagnostics.Span{
			File:   lexer.fileName,
			Line:   line,
			Column: column,
		},
		msg,
	)

	return lexer.makeToken(Illegal, lex)
}
