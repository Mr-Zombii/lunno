package lexer

type Lexer struct {
	Source      []byte
	position    int
	line        int
	column      int
	startPos    int
	startLine   int
	startColumn int
	fileName    string
}

func NewLexer(source, filename string) *Lexer {
	return &Lexer{
		Source:   []byte(source),
		line:     1,
		column:   1,
		fileName: filename,
	}
}

func Tokenize(source, filename string) (*Lexer, []Token, error) {
	lexer := NewLexer(source, filename)
	var tokens []Token

	for {
		tok := lexer.Next()
		tokens = append(tokens, tok)

		//if tok.Type == Illegal {
		//	return nil, nil, lexer.error("illegal token")
		//}
		if tok.Type == EndOfFile {
			break
		}
	}
	return lexer, tokens, nil
}

func (lexer *Lexer) Next() Token {
	for {
		for {
			ch := lexer.peek()
			cc := classify(ch)

			if cc == CC_Whitespace || cc == CC_Newline {
				lexer.advance()
			} else {
				break
			}
		}

		lexer.startPos = lexer.position
		lexer.startLine = lexer.line
		lexer.startColumn = lexer.column

		state := S_Start

		for {
			ch := lexer.peek()
			cc := classify(ch)

			next, ok := dfa[state][cc]
			if !ok {
				break
			}

			state = next
			lexer.advance()
		}

		if state == S_Done && lexer.startPos == lexer.position {
			return lexer.makeToken(EndOfFile, "")
		}

		if emit, ok := accepting[state]; ok {
			lex := string(lexer.Source[lexer.startPos:lexer.position])
			tokType := emit(lexer, lex)
			return lexer.makeToken(tokType, lex)
		}
	}
}

func (lexer *Lexer) peek() byte {
	if lexer.position >= len(lexer.Source) {
		return 0
	}
	return lexer.Source[lexer.position]
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

func (lexer *Lexer) makeToken(tt TokenType, lex string) Token {
	return Token{
		Type:   tt,
		Lexeme: lex,
		Line:   lexer.startLine,
		Column: lexer.startColumn,
		File:   lexer.fileName,
	}
}

func lookupIdentifier(lex string) TokenType {
	if token, ok := keywords[lex]; ok {
		return token
	}
	return Identifier
}
