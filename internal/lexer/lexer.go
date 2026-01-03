package lexer

import (
	"strings"
	"unicode"
)

type Lexer struct {
	Source   []byte
	position int
	line     int
	column   int
	fileName string
}

func NewLexer(source, filename string) *Lexer {
	return &Lexer{
		Source:   []byte(source),
		position: 0,
		line:     1,
		column:   1,
		fileName: filename,
	}
}

func Tokenize(source, filename string) ([]Token, error) {
	lexer := NewLexer(source, filename)
	var tokens []Token
	for {
		token := lexer.Next()
		tokens = append(tokens, token)
		if token.Type == Illegal {
			lexer.error("illegal token")
		}
		if token.Type == EndOfFile {
			break
		}
	}
	return tokens, nil
}

func (l *Lexer) Next() Token {
	l.skipWhiteSpaceAndComments()
	startLine, startColumn := l.line, l.column
	ch := l.peek()

	if ch == 0 {
		return l.makeToken(EndOfFile, "", startLine, startColumn)
	}

	if isIdentifierStart(ch) {
		lex := l.scanIdentifier()
		typ := LookupIdentifier(lex)
		return l.makeToken(typ, lex, startLine, startColumn)
	}

	if unicode.IsDigit(rune(ch)) {
		lex, typ := l.scanNumber()
		return l.makeToken(typ, lex, startLine, startColumn)
	}

	if ch == '"' {
		lex := l.scanString()
		return l.makeToken(String, lex, startLine, startColumn)
	}

	if ch == '\'' {
		lex := l.scanChar()
		return l.makeToken(Char, lex, startLine, startColumn)
	}

	if token, ok := l.matchMultiCharOperator(); ok {
		return l.makeToken(token.Type, token.Lexeme, startLine, startColumn)
	}

	if tokenType, ok := singleCharTokens[ch]; ok {
		l.advance()
		return l.makeToken(tokenType, string(ch), startLine, startColumn)
	}

	r := l.advance()
	return l.makeToken(Illegal, string(r), startLine, startColumn)
}

func (l *Lexer) peek() byte {
	if l.position >= len(l.Source) {
		return 0
	}
	return l.Source[l.position]
}

func (l *Lexer) peekNext() byte {
	if l.position+1 >= len(l.Source) {
		return 0
	}
	return l.Source[l.position+1]
}

func (l *Lexer) advance() byte {
	if l.position >= len(l.Source) {
		return 0
	}
	ch := l.Source[l.position]
	l.position++
	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return ch
}

func (l *Lexer) makeToken(typ TokenType, lex string, line, column int) Token {
	return Token{
		Type:   typ,
		Lexeme: lex,
		Line:   line,
		Column: column,
		File:   l.fileName,
	}
}

func (l *Lexer) skipWhiteSpaceAndComments() {
	for {
		ch := l.peek()
		if ch == 0 {
			return
		}
		if ch == '#' {
			l.skipLineComment()
		} else if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.advance()
		} else {
			break
		}
	}
}
func (l *Lexer) skipLineComment() {
	l.advance()
	for {
		ch := l.peek()
		if ch == 0 || ch == '\n' {
			break
		}
		l.advance()
	}
	if l.peek() == '\n' {
		l.advance()
	}
}

func (l *Lexer) scanIdentifier() string {
	start := l.position
	for {
		ch := l.peek()
		if ch == 0 || !isIdentifierPart(ch) {
			break
		}
		l.advance()
	}
	return string(l.Source[start:l.position])
}

func isIdentifierStart(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isIdentifierPart(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_'
}

func LookupIdentifier(lex string) TokenType {
	if token, ok := keywords[lex]; ok {
		return token
	}
	return Identifier
}

func (l *Lexer) scanNumber() (string, TokenType) {
	start := l.position
	hasDot := false

	for {
		ch := l.peek()

		if unicode.IsDigit(rune(ch)) {
			l.advance()
			continue
		}

		if ch == '.' && !hasDot {
			hasDot = true
			l.advance()
			continue
		}

		break
	}

	lex := string(l.Source[start:l.position])
	if hasDot {
		return lex, Float
	}
	return lex, Integer
}

func (l *Lexer) scanString() string {
	quote := l.advance()
	start := l.position
	for {
		ch := l.peek()
		if ch == 0 {
			l.error("unterminated string literal")
		}
		if ch == quote {
			l.advance()
			break
		}
		if ch == '\\' {
			l.advance()
			if l.peek() != 0 {
				l.advance()
			}
		} else {
			l.advance()
		}
	}
	return string(l.Source[start : l.position-1])
}

func (l *Lexer) scanChar() string {
	l.advance()
	start := l.position
	ch := l.advance()
	if ch == 0 {
		l.error("unterminated char literal")
	}
	if l.peek() != '\'' {
		l.error("char literal must be 1 character")
	}
	l.advance()
	return string(l.Source[start : l.position-1])
}

var multiCharOperators = map[string]TokenType{
	"==": Equal,
	"!=": NotEqual,
	"<=": LessThanOrEqual,
	">=": GreaterThanOrEqual,
	"->": Arrow,
}

func (l *Lexer) matchMultiCharOperator() (Token, bool) {
	for op, typ := range multiCharOperators {
		if strings.HasPrefix(string(l.Source[l.position:]), op) {
			for range op {
				l.advance()
			}
			return Token{
				Type:   typ,
				Lexeme: op,
				Line:   l.line,
				Column: l.column,
				File:   l.fileName}, true
		}
	}
	return Token{}, false
}

var singleCharTokens = map[byte]TokenType{
	'(': LeftParen, ')': RightParen,
	'[': LeftBracket, ']': RightBracket,
	'{': LeftBrace, '}': RightBrace,
	'+': Plus, '-': Minus,
	'*': Asterisk, '/': Slash,
	':': Colon, ',': Comma,
	'=': Assign,
}
