package parser

import (
	"lunno/internal/lexer"
	"strconv"
)

type Parser struct {
	tokens   []lexer.Token
	position int
	errors   []string
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
		errors:   []string{},
	}
}

func ParseProgram(tokens []lexer.Token) (*Program, []string) {
	parser := NewParser(tokens)
	program := &Program{Expressions: []Expression{}}
	for parser.cur().Type != lexer.EndOfFile {
		expression := parser.parseExpression(0)
		if expression == nil {
			parser.advance()
			break
		}
		program.Expressions = append(program.Expressions, expression)
	}
	return program, parser.errors
}

func (parser *Parser) parseExpression(minPrecedence int) Expression {
	left := parser.parsePrimary()
	if left == nil {
		return nil
	}
	for {
		op := parser.cur()
		precedence := op.Type.Precedence()
		if precedence < minPrecedence {
			break
		}
		parser.advance()
		right := parser.parseExpression(precedence + 1)
		left = &InfixExpression{
			Left:     left,
			Operator: op,
			Right:    right,
			Position: op,
		}
	}
	return left
}

func (parser *Parser) parsePrimary() Expression {
	token := parser.cur()
	switch token.Type {
	case lexer.Integer:
		parser.advance()
		v, _ := strconv.ParseInt(token.Lexeme, 10, 64)
		return &IntegerLiteral{Value: v, Raw: token.Lexeme, Position: token}
	case lexer.Float:
		parser.advance()
		v, _ := strconv.ParseFloat(token.Lexeme, 64)
		return &FloatLiteral{Value: v, Raw: token.Lexeme, Position: token}
	case lexer.LeftParen:
		parser.advance()
		expr := parser.parseExpression(0)
		parser.expect(lexer.RightParen)
		return expr
	default:
		//parser.error(token, "unexpected token")
		return nil
	}
}

func (parser *Parser) cur() lexer.Token {
	if parser.position >= len(parser.tokens) {
		return lexer.Token{
			Type:   lexer.EndOfFile,
			Lexeme: "",
			Line:   0,
			Column: 0}
	}
	return parser.tokens[parser.position]
}

func (parser *Parser) advance() lexer.Token {
	token := parser.cur()
	parser.position++
	return token
}

func (parser *Parser) expect(typ lexer.TokenType) lexer.Token {
	token := parser.cur()
	if token.Type != typ {
		//parser.error(token, fmt.Sprintf("expected %q, got %q", typ, token.Type))
		return lexer.Token{}
	}
	parser.advance()
	return token
}
