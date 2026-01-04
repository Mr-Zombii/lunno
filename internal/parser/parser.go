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
			continue
		}
		program.Expressions = append(program.Expressions, expression)
	}
	return program, parser.errors
}

func (parser *Parser) parseExpression(minPrecedence int) Expression {
	token := parser.cur()
	if token.Type == lexer.KwLet {
		return parser.parseLetExpression()
	}
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
	var expr Expression
	switch token.Type {
	case lexer.Integer:
		parser.advance()
		value, _ := strconv.ParseInt(token.Lexeme, 10, 64)
		expr = &IntegerLiteral{
			Value:    value,
			Raw:      token.Lexeme,
			Position: token}
	case lexer.Float:
		parser.advance()
		value, _ := strconv.ParseFloat(token.Lexeme, 64)
		expr = &FloatLiteral{
			Value:    value,
			Raw:      token.Lexeme,
			Position: token}
	case lexer.LeftParen:
		parser.advance()
		expr = parser.parseExpression(0)
		parser.expect(lexer.RightParen)
	case lexer.Identifier:
		parser.advance()
		expr = &Identifier{
			Name:     token.Lexeme,
			Position: token}
	case lexer.KwFn:
		expr = parser.parseFunctionLiteral()
	case lexer.KwIf:
		expr = parser.parseIfExpression()
	case lexer.KwImport:
		parser.advance()
		mod := parser.expect(lexer.Identifier)
		expr = &ImportExpression{
			Module:   mod.Lexeme,
			Position: token}
	case lexer.LeftBracket:
		parser.advance()
		var elements []Expression
		for parser.cur().Type != lexer.RightBracket && parser.cur().Type != lexer.EndOfFile {
			elem := parser.parseExpression(0)
			if elem != nil {
				elements = append(elements, elem)
			}
			if parser.cur().Type == lexer.Comma {
				parser.advance()
			} else {
				break
			}
		}
		parser.expect(lexer.RightBracket)
		expr = &ListExpression{
			Elements: elements,
			Position: token}
	default:
		return nil
	}
	return parser.parsePostfix(expr)
}

func (parser *Parser) parsePostfix(expr Expression) Expression {
	for {
		switch parser.cur().Type {
		case lexer.LeftParen:
			callToken := parser.cur()
			parser.advance()
			var args []Expression
			for parser.cur().Type != lexer.RightParen && parser.cur().Type != lexer.EndOfFile {
				arg := parser.parseExpression(0)
				if arg != nil {
					args = append(args, arg)
				}
				if parser.cur().Type == lexer.Comma {
					parser.advance()
				} else {
					break
				}
			}
			parser.expect(lexer.RightParen)
			expr = &CallExpression{
				Callee:    expr,
				Arguments: args,
				Position:  callToken}
		case lexer.LeftBracket:
			indexToken := parser.cur()
			parser.advance()
			idx := parser.parseExpression(0)
			parser.expect(lexer.RightBracket)
			expr = &IndexExpression{
				Target:   expr,
				Index:    idx,
				Position: indexToken}
		default:
			return expr
		}
	}
}

func (parser *Parser) parseLetExpression() Expression {
	token := parser.cur()
	parser.advance()
	recursive := false
	if parser.cur().Type == lexer.KwRec {
		recursive = true
		parser.advance()
	}
	name := parser.expect(lexer.Identifier)
	if name.Type != lexer.Identifier {
		return nil
	}
	if parser.cur().Type != lexer.Colon {
		return nil
	}
	parser.advance()
	typ := parser.parseType()
	if typ == nil {
		return nil
	}
	if parser.cur().Type != lexer.Assign {
		return nil
	}
	parser.advance()
	for parser.cur().Type == lexer.Newline {
		parser.advance()
	}
	value := parser.parseExpression(0)
	if value == nil {
		return nil
	}
	if fn, ok := value.(*FunctionLiteralExpression); ok {
		return &FunctionDeclarationExpression{
			Name:      name,
			Recursive: recursive,
			Signature: typ,
			Function:  fn,
			Position:  token,
		}
	}
	return &VariableDeclarationExpression{
		Name:      name,
		Type:      typ,
		Value:     value,
		Recursive: recursive,
		Position:  token,
	}
}

func (parser *Parser) parseFunctionLiteral() Expression {
	fnToken := parser.cur()
	parser.advance()
	var parameters []Parameter
	parser.expect(lexer.LeftParen)
	for parser.cur().Type != lexer.RightParen && parser.cur().Type != lexer.EndOfFile {
		paramName := parser.expect(lexer.Identifier)
		var paramType TypeNode
		if parser.cur().Type == lexer.Colon {
			parser.advance()
			paramType = parser.parseType()
		}
		parameters = append(parameters, Parameter{
			Name:     paramName,
			Type:     paramType,
			Position: paramName,
		})
		if parser.cur().Type == lexer.Comma {
			parser.advance()
		} else {
			break
		}
	}
	parser.expect(lexer.RightParen)
	var returnType TypeNode
	if parser.cur().Type == lexer.Arrow {
		parser.advance()
		returnType = parser.parseType()
	}
	for parser.cur().Type == lexer.Newline {
		parser.advance()
	}
	body := parser.parseExpression(0)
	if body == nil {
		// parser.errors = append(parser.errors, "expected function body expression")
		return nil
	}
	return &FunctionLiteralExpression{
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
		Position:   fnToken,
	}
}

func (parser *Parser) parseIfExpression() Expression {
	ifToken := parser.cur()
	parser.advance()
	condition := parser.parseExpression(0)
	if parser.cur().Type != lexer.KwThen {
		// parser.errors = append(parser.errors, "expected 'then' after if condition")
		return nil
	}
	parser.advance()
	thenBranch := parser.parseExpression(0)
	if parser.cur().Type != lexer.KwElse {
		// parser.errors = append(parser.errors, "expected 'else' after then branch")
		return nil
	}
	parser.advance()
	elseBranch := parser.parseExpression(0)
	return &IfExpression{
		Condition: condition,
		Then:      thenBranch,
		Else:      elseBranch,
		Position:  ifToken,
	}
}

func (parser *Parser) parseType() TypeNode {
	token := parser.cur()
	switch token.Type {
	case lexer.KwFn:
		parser.advance()
		var params []TypeNode
		parser.expect(lexer.LeftParen)
		for parser.cur().Type != lexer.RightParen && parser.cur().Type != lexer.EndOfFile {
			paramType := parser.parseType()
			params = append(params, paramType)
			if parser.cur().Type == lexer.Comma {
				parser.advance()
			}
		}
		parser.expect(lexer.RightParen)
		var returnType TypeNode
		if parser.cur().Type == lexer.Arrow {
			parser.advance()
			returnType = parser.parseType()
		}
		return &FunctionType{
			Parameters: params,
			Return:     returnType,
			Position:   token}
	case lexer.LeftBracket:
		parser.advance()
		elemType := parser.parseType()
		parser.expect(lexer.RightBracket)
		return &ListType{
			Element:  elemType,
			Position: token}
	case lexer.Identifier:
		parser.advance()
		return &SimpleType{
			Name: token.Lexeme,
			Pos:  token}
	default:
		parser.advance()
		return &SimpleType{
			Name: token.Lexeme,
			Pos:  token}
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
