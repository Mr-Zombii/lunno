package lexer

import "fmt"

type Lexer struct {
	Source      []rune
	position    uint16
	line        uint16
	column      uint16
	startPos    uint16
	startLine   uint16
	startColumn uint16
	fileName    string
	sourceLen   uint16
}

func NewLexer(source, filename string) *Lexer {
	return &Lexer{
		Source:    []rune(source),
		line:      1,
		column:    1,
		fileName:  filename,
		sourceLen: uint16(len(source)),
	}
}

func Tokenize(source, filename string) (*Lexer, []Token, error) {
	if len(source) == 0 {
		return nil, nil, fmt.Errorf("%s: empty source file", filename)
	}
	lexer := NewLexer(source, filename)
	var tokens []Token
	for {
		tok := lexer.Next()
		tokens = append(tokens, tok)
		if tok.Type == EndOfFile {
			break
		}
	}
	return lexer, tokens, nil
}

func (lexer *Lexer) Next() Token {
	for lexer.position < lexer.sourceLen {
		ch := lexer.Source[lexer.position]
		cc := classify(ch)
		if cc == CC_Whitespace || cc == CC_Newline {
			lexer.advance()
			continue
		}
		if ch == '#' {
			for lexer.position < lexer.sourceLen && lexer.Source[lexer.position] != '\n' {
				lexer.advance()
			}
			continue
		}
		break
	}
	lexer.startPos = lexer.position
	lexer.startLine = lexer.line
	lexer.startColumn = lexer.column

	if lexer.position >= lexer.sourceLen {
		return lexer.makeToken(EndOfFile, "")
	}
	for opLen := uint16(2); opLen > 0; opLen-- {
		if lexer.position+opLen <= lexer.sourceLen {
			sub := string(lexer.Source[lexer.position : lexer.position+opLen])
			if tt, ok := multiCharOperators[sub]; ok {
				lexer.position += opLen
				lexer.column += opLen
				return lexer.makeToken(tt, sub)
			}
		}
	}
	ch := lexer.peek()
	if tt, ok := singleCharTokens[byte(ch)]; ok {
		lexer.advance()
		return lexer.makeToken(tt, string(ch))
	}
	state := S_Start
	for lexer.position < lexer.sourceLen {
		ch := lexer.peek()
		cc := classify(ch)
		next, ok := dfa[state][cc]
		if !ok {
			break
		}
		state = next
		lexer.advance()
		if lexer.position >= lexer.sourceLen {
			switch state {
			case S_String, S_StringEsc:
				return lexer.errorAt(
					"unterminated string literal",
					string(lexer.Source[lexer.startPos:lexer.position]),
					lexer.startLine,
					lexer.column-1,
				)
			case S_Char, S_CharEsc, S_CharDone:
				return lexer.errorAt(
					"unterminated character literal",
					string(lexer.Source[lexer.startPos:lexer.position]),
					lexer.startLine,
					lexer.column-1,
				)
			default:
			}
		}
	}
	lex := string(lexer.Source[lexer.startPos:lexer.position])
	if state == S_Done && lexer.startPos == lexer.position {
		return lexer.makeToken(EndOfFile, "")
	}
	if state == S_Float && len(lex) > 0 && lex[len(lex)-1] == '.' {
		return lexer.errorAt(
			"malformed float literal",
			lex,
			lexer.startLine,
			lexer.column-1,
		)
	}
	if emit, ok := accepting[state]; ok {
		tokType := emit(lexer, lex)
		if tokType == Illegal {
			return lexer.handleIllegal(lex)
		}
		return lexer.makeToken(tokType, lex)
	}
	if lexer.position == lexer.startPos {
		ch := lexer.advance()
		return lexer.errorToken(
			fmt.Sprintf("unexpected character '%c'", ch),
			string(ch),
		)
	}
	return lexer.errorToken("invalid token", lex)
}

func (lexer *Lexer) handleIllegal(lex string) Token {
	if len(lex) == 0 {
		return lexer.errorToken("invalid token", lex)
	}
	switch lex[0] {
	case '\'':
		content := lex[1 : len(lex)-1]
		if len(content) == 0 {
			return lexer.errorAt("empty characters not allowed", lex, lexer.startLine, lexer.startColumn+1)
		} else if content[0] == '\\' {
			if len(content) != 2 || !isValidEscape(content[1]) {
				return lexer.errorAt("invalid escape sequence in character literal", lex, lexer.startLine, lexer.startColumn+1)
			}
		} else if len([]rune(content)) != 1 {
			return lexer.errorAt("character literal must contain exactly one character", lex, lexer.startLine, lexer.startColumn+1)
		}
	case '"':
		content := lex[1 : len(lex)-1]
		if len(content) == 0 {
			return lexer.errorAt("empty strings not allowed", lex, lexer.startLine, lexer.startColumn+1)
		}
		for i := uint16(1); i < uint16(len(lex)-1); i++ {
			if lex[i] == '\\' && !isValidEscape(lex[i+1]) {
				errCol := lexer.startColumn + i
				return lexer.errorAt("invalid escape sequence in string literal", lex, lexer.startLine, errCol)
			}
		}
	}
	return lexer.errorToken("invalid token", lex)
}

func (lexer *Lexer) peek() rune {
	if lexer.position >= lexer.sourceLen {
		panic(fmt.Sprintf(
			"%s:%d:%d:  attempted to peek past end of file",
			lexer.fileName,
			lexer.line,
			lexer.column,
		))
	}
	return lexer.Source[lexer.position]
}

func (lexer *Lexer) advance() rune {
	if lexer.position >= lexer.sourceLen {
		panic(fmt.Sprintf(
			"%s:%d:%d: attempted to advance past end of file",
			lexer.fileName,
			lexer.line,
			lexer.column,
		))
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

func (lexer *Lexer) makeToken(tt TokenType, lex string) Token {
	if tt < 0 {
		panic(fmt.Sprintf(
			"%s:%d:%d:  attempted to make token with invalid TokenType %d",
			lexer.fileName,
			lexer.startLine,
			lexer.startColumn,
			tt,
		))
	}
	if lex == "" && tt != EndOfFile {
		panic(fmt.Sprintf(
			"%s:%d:%d: attempt to make token with empty lexeme for type %v",
			lexer.fileName,
			lexer.startLine,
			lexer.startColumn,
			tt,
		))
	}
	return Token{
		Type:   tt,
		Lexeme: lex,
		Line:   lexer.startLine,
		Column: lexer.startColumn,
		File:   lexer.fileName,
	}
}

func lookupIdentifier(lex string) TokenType {
	if lex == "" {
		panic("lookupIdentifier: empty string provided")
	}
	if token, ok := keywords[lex]; ok {
		return token
	}
	return Identifier
}

func isValidEscape(ch byte) bool {
	switch ch {
	case 'n', 't', 'r',
		'u', 'x', '0',
		'\\', '\'', '"':
		return true
	default:
		return false
	}
}
