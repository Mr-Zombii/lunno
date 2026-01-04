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

func Tokenize(source, filename string) (*Lexer, []Token, error) {
	lexer := NewLexer(source, filename)
	var tokens []Token
	for {
		token := lexer.Next()
		tokens = append(tokens, token)
		if token.Type == Illegal {
			err := lexer.error("illegal token")
			if err != nil {
				return nil, nil, err
			}
		}
		if token.Type == EndOfFile {
			break
		}
	}
	return lexer, tokens, nil
}

func (lexer *Lexer) Next() Token {
	lexer.skipWhiteSpaceAndComments()
	startLine, startColumn := lexer.line, lexer.column
	ch := lexer.peek()
	switch {
	case ch == 0:
		return lexer.makeToken(EndOfFile, "", startLine, startColumn)
	case isIdentifierStart(ch):
		lex := lexer.scanIdentifier()
		typ := LookupIdentifier(lex)
		return lexer.makeToken(typ, lex, startLine, startColumn)
	case unicode.IsDigit(rune(ch)):
		return lexer.scanNumber(startLine, startColumn)
	case ch == '"' || ch == '\'':
		typ := String
		if ch == '\'' {
			typ = Char
		}
		return lexer.makeToken(typ, lexer.scanDelimited(ch, ch == '\''), startLine, startColumn)
	default:
		if token, ok := lexer.matchMultiCharOperator(); ok {
			return lexer.makeToken(token.Type, token.Lexeme, startLine, startColumn)
		}
		if tokenType, ok := singleCharTokens[ch]; ok {
			lexer.advance()
			return lexer.makeToken(tokenType, string(ch), startLine, startColumn)
		}
		return lexer.makeToken(Illegal, string(lexer.advance()), startLine, startColumn)
	}
}

func (lexer *Lexer) peek() byte {
	if lexer.position >= len(lexer.Source) {
		return 0
	}
	return lexer.Source[lexer.position]
}

func (lexer *Lexer) peekNext() byte {
	if lexer.position+1 >= len(lexer.Source) {
		return 0
	}
	return lexer.Source[lexer.position+1]
}

func (lexer *Lexer) advance() byte {
	if lexer.position >= len(lexer.Source) {
		return 0
	}
	ch := lexer.Source[lexer.position]
	lexer.position++
	if ch == '\n' {
		lexer.line++
		lexer.column = 1
	} else {
		lexer.column++
	}
	return ch
}

func (lexer *Lexer) makeToken(typ TokenType, lex string, line, column int) Token {
	return Token{
		Type:   typ,
		Lexeme: lex,
		Line:   line,
		Column: column,
		File:   lexer.fileName,
	}
}

func (lexer *Lexer) skipWhiteSpaceAndComments() {
	for {
		ch := lexer.peek()
		if ch == 0 {
			return
		}
		if ch == '#' {
			lexer.skipLineComment()
		} else if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			lexer.advance()
		} else {
			break
		}
	}
}
func (lexer *Lexer) skipLineComment() {
	lexer.advance()
	for {
		ch := lexer.peek()
		if ch == 0 || ch == '\n' {
			break
		}
		lexer.advance()
	}
	if lexer.peek() == '\n' {
		lexer.advance()
	}
}

func (lexer *Lexer) scanIdentifier() string {
	start := lexer.position
	for {
		ch := lexer.peek()
		if ch == 0 || !isIdentifierPart(ch) {
			break
		}
		lexer.advance()
	}
	return string(lexer.Source[start:lexer.position])
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

func (lexer *Lexer) scanNumber(line, col int) Token {
	start := lexer.position
	hasDot := false
	for {
		ch := lexer.peek()
		if unicode.IsDigit(rune(ch)) {
			lexer.advance()
			continue
		}
		if ch == '.' && !hasDot {
			hasDot = true
			lexer.advance()
			continue
		}
		break
	}
	typ := Integer
	if hasDot {
		typ = Float
	}
	return lexer.makeToken(typ, string(lexer.Source[start:lexer.position]), line, col)
}

func (lexer *Lexer) scanDelimited(delim byte, isChar bool) string {
	lexer.advance()
	start := lexer.position
	for {
		ch := lexer.peek()
		if ch == 0 {
			err := lexer.error("unterminated literal")
			if err != nil {
				return ""
			}
			return ""
		}
		if ch == delim {
			lexer.advance()
			break
		}
		if ch == '\\' {
			lexer.advance()
			if lexer.peek() != 0 {
				lexer.advance()
			}
		} else {
			lexer.advance()
		}
	}
	lex := string(lexer.Source[start : lexer.position-1])
	if isChar && len([]rune(lex)) != 1 {
		err := lexer.error("char literal must be 1 character")
		if err != nil {
			return ""
		}
		return ""
	}
	return lex
}

func (lexer *Lexer) matchMultiCharOperator() (Token, bool) {
	for op, typ := range multiCharOperators {
		if strings.HasPrefix(string(lexer.Source[lexer.position:]), op) {
			for range op {
				lexer.advance()
			}
			return Token{
				Type:   typ,
				Lexeme: op,
				Line:   lexer.line,
				Column: lexer.column,
				File:   lexer.fileName}, true
		}
	}
	return Token{}, false
}
