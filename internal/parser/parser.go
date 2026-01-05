package parser

import (
	"fmt"
	"lunno/internal/lexer"
	"strconv"
)

type Parser struct {
	tokens   []lexer.Token
	position int
	errors   []string
	lexer    *lexer.Lexer
}

func NewParser(tokens []lexer.Token, lexer *lexer.Lexer) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
		errors:   []string{},
		lexer:    lexer,
	}
}

func ParseProgram(tokens []lexer.Token, lx *lexer.Lexer) (*Program, []string) {
	parser := NewParser(tokens, lx)
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
	for parser.cur().Type == lexer.Newline {
		parser.advance()
	}
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
		if right == nil {
			err := parser.error(op, "expected expression on right-hand side of operator")
			parser.errors = append(parser.errors, err.Error())
			return left
		}
		left = &InfixExpression{
			Left:     left,
			Operator: op,
			Right:    right,
			Position: op,
		}
	}
	for parser.cur().Type == lexer.Newline {
		parser.advance()
	}
	return left
}

func (parser *Parser) parseExpressionList(args []Expression) []Expression {
	tk := parser.prev()
	bracePair, valid := lexer.Braces[tk.Type]
	if !valid {
		parser.errors = append(parser.errors, fmt.Sprintf("Expected '(', '[', or '{', not %v", tk))
	}
	if isValidClosing, isClosing := bracePair.ValidClosings[parser.cur().Type]; isClosing {
		if !isValidClosing {
			parser.errors = append(parser.errors, fmt.Sprintf("Expected '%v', not %v", bracePair.ClosingChar, tk))
		}
		return args
	}

	for parser.cur().Type != bracePair.ClosingTT && parser.cur().Type != lexer.EndOfFile {
		arg := parser.parseExpression(0)
		if arg != nil {
			args = append(args, arg)
		} else {
			e := parser.error(parser.cur(), "invalid expression in "+bracePair.Name)
			parser.errors = append(parser.errors, e.Error())
			parser.advance()
			continue
		}
		if parser.cur().Type == lexer.Comma {
			parser.advance()
		} else if parser.cur().Type != bracePair.ClosingTT {
			e := parser.error(parser.cur(), fmt.Sprintf("expected ',' or '%v' in %v", bracePair.ClosingChar, bracePair.Name))
			parser.errors = append(parser.errors, e.Error())
			break
		}
	}
	return args
}

func (parser *Parser) parsePrimary() Expression {
	token := parser.cur()
	if token.Type == lexer.EndOfFile {
		return nil
	}
	var expr Expression
	switch token.Type {
	case lexer.Integer:
		parser.advance()
		value, err := strconv.ParseInt(token.Lexeme, 10, 64)
		if err != nil {
			e := parser.error(token, "invalid integer literal")
			parser.errors = append(parser.errors, e.Error())
			return nil
		}
		expr = &IntegerLiteral{
			Value:    value,
			Raw:      token.Lexeme,
			Position: token}
	case lexer.Float:
		parser.advance()
		value, err := strconv.ParseFloat(token.Lexeme, 64)
		if err != nil {
			e := parser.error(token, "invalid float literal")
			parser.errors = append(parser.errors, e.Error())
			return nil
		}
		expr = &FloatLiteral{
			Value:    value,
			Raw:      token.Lexeme,
			Position: token}
	case lexer.Char:
		token := parser.expect(lexer.Char)
		expr = &CharacterLiteral{
			Value:    token.Lexeme[0],
			Raw:      token.Lexeme,
			Position: token,
		}
	case lexer.String:
		token := parser.expect(lexer.String)
		expr = &StringLiteral{
			Value:    token.Lexeme,
			Position: token,
		}
	case lexer.Bool:
		token := parser.expect(lexer.Bool)
		expr = &BooleanLiteral{
			Value:    token.Lexeme == "true",
			Position: token,
		}
	case lexer.LeftParen:
		parser.advance()
		if parser.cur().Type == lexer.RightParen {
			token := parser.cur()
			parser.advance()
			return &UnitLiteral{
				Position: token,
			}
		}
		expr = parser.parseExpression(0)
		parser.expect(lexer.RightParen)
	case lexer.Identifier:
		parser.advance()
		expr = &Identifier{
			Name:     token.Lexeme,
			Position: token}
	case lexer.KwUnit:
		token := parser.expect(lexer.KwUnit)
		expr = &UnitLiteral{
			Position: token,
		}
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
		elements = parser.parseExpressionList(elements)
		parser.expect(lexer.RightBracket)
		expr = &ListExpression{
			Elements: elements,
			Position: token}
	default:
		e := parser.error(token, fmt.Sprintf("unexpected token %q", token.Lexeme))
		parser.errors = append(parser.errors, e.Error())
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
			args = parser.parseExpressionList(args)
			if parser.cur().Type != lexer.RightParen {
				e := parser.error(parser.cur(), "expected ')' after function call")
				parser.errors = append(parser.errors, e.Error())
			} else {
				parser.advance()
			}
			expr = &CallExpression{
				Callee:    expr,
				Arguments: args,
				Position:  callToken}
		case lexer.LeftBracket:
			startToken := parser.cur()
			parser.advance()
			var startExpr, endExpr Expression
			if parser.cur().Type != lexer.Colon && parser.cur().Type != lexer.RightBracket {
				startExpr = parser.parsePrimary()
			}
			if parser.cur().Type == lexer.Colon {
				parser.advance()
				if parser.cur().Type != lexer.RightBracket {
					endExpr = parser.parsePrimary()
				}
				if parser.cur().Type != lexer.RightBracket {
					e := parser.error(parser.cur(), "expected ']' to close slice expression")
					parser.errors = append(parser.errors, e.Error())
				} else {
					parser.advance()
				}
				expr = &SliceExpression{
					Target:   expr,
					Start:    startExpr,
					End:      endExpr,
					Position: startToken,
				}
			} else {
				if parser.cur().Type != lexer.RightBracket {
					e := parser.error(parser.cur(), "expected ']' after index expression")
					parser.errors = append(parser.errors, e.Error())
				} else {
					parser.advance()
				}
				expr = &IndexExpression{
					Target:   expr,
					Index:    startExpr,
					Position: startToken,
				}
			}
		default:
			return expr
		}
	}
}

func (parser *Parser) parseSliceExpr() Expression {
	switch parser.cur().Type {
	case lexer.RightBracket, lexer.Colon:
		return nil
	default:
		return parser.parsePrimary()
	}
}

func (parser *Parser) parseLetExpression() Expression {
	letToken := parser.cur()
	parser.advance()
	recursive := false
	if parser.cur().Type == lexer.KwRec {
		recursive = true
		parser.advance()
	}
	name := parser.expect(lexer.Identifier)
	if name.Type != lexer.Identifier {
		e := parser.error(name, "expected identifier after 'let'")
		parser.errors = append(parser.errors, e.Error())
		return nil
	}
	if parser.cur().Type != lexer.Colon {
		e := parser.error(parser.cur(), "expected ':' after identifier in let declaration")
		parser.errors = append(parser.errors, e.Error())
		return nil
	}
	parser.advance()
	typ := parser.parseType()
	if typ == nil {
		e := parser.error(parser.cur(), "expected type in let declaration")
		parser.errors = append(parser.errors, e.Error())
		return nil
	}
	if parser.cur().Type != lexer.Assign {
		e := parser.error(parser.cur(), "expected '=' after type in let declaration")
		parser.errors = append(parser.errors, e.Error())
		return nil
	}
	parser.advance()
	assignToken := parser.cur()
	value := parser.parseExpression(0)
	if value == nil {
		e := parser.error(assignToken, "expected value in let declaration")
		parser.errors = append(parser.errors, e.Error())
		return nil
	}
	if fn, ok := value.(*FunctionLiteralExpression); ok {
		return &FunctionDeclarationExpression{
			Name:      name,
			Recursive: recursive,
			Signature: typ,
			Function:  fn,
			Position:  letToken,
		}
	}
	return &VariableDeclarationExpression{
		Name:      name,
		Type:      typ,
		Value:     value,
		Recursive: recursive,
		Position:  letToken,
	}
}

func (parser *Parser) parseFunctionLiteral() Expression {
	fnToken := parser.cur()
	parser.advance()
	var parameters []Parameter
	parser.expect(lexer.LeftParen)
	for parser.cur().Type != lexer.RightParen && parser.cur().Type != lexer.EndOfFile {
		paramName := parser.expect(lexer.Identifier)
		if paramName.Type != lexer.Identifier {
			e := parser.error(parser.cur(), "expected parameter name")
			parser.errors = append(parser.errors, e.Error())
			return nil
		}
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
	if parser.cur().Type == lexer.Arrow {
		parser.advance()
	}
	for parser.cur().Type == lexer.Newline {
		parser.advance()
	}

	return &FunctionLiteralExpression{
		Parameters: parameters,
		Body:       parser.parseBlock(fnToken, "function body"),
		Position:   fnToken,
	}
}

func (parser *Parser) parseBlock(token lexer.Token, name string) Expression {
	targetIndention := token.Indentation + 1

	var exprs []Expression
	firstExprTk := parser.cur()
	first := parser.parseExpression(0)
	if first == nil {
		e := parser.error(token, "expected "+name+" expression")
		parser.errors = append(parser.errors, e.Error())
		return nil
	}
	exprs = append(exprs, first)
	if firstExprTk.Indentation == targetIndention {
		for parser.cur().Indentation == targetIndention && parser.prev().Type == lexer.Newline {
			if parser.cur().Type == lexer.EndOfFile || parser.cur().Type == lexer.KwLet {
				break
			}

			next := parser.parseExpression(0)
			if next == nil {
				break
			}
			exprs = append(exprs, next)
		}
	}
	var body Expression
	if len(exprs) == 1 {
		body = exprs[0]
	} else {
		body = &BlockExpression{
			Expressions: exprs,
			Position:    token,
		}
	}
	return body
}

func (parser *Parser) parseIfExpression() Expression {
	ifToken := parser.cur()
	parser.advance()
	condition := parser.parseExpression(0)
	if parser.cur().Type != lexer.KwThen {
		e := parser.error(parser.cur(), "expected 'then' after if condition")
		parser.errors = append(parser.errors, e.Error())
		return nil
	}
	parser.advance()
	thenBranch := parser.parseBlock(ifToken, "if body")
	if parser.cur().Type != lexer.KwElse {
		e := parser.error(parser.cur(), "expected 'else' after then branch")
		parser.errors = append(parser.errors, e.Error())
		return nil
	}
	elseToken := parser.advance()
	elseBranch := parser.parseBlock(elseToken, "else body")
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
	case lexer.KwList:
		parser.advance()
		parser.expect(lexer.LeftParen)
		elemType := parser.parseType()
		parser.expect(lexer.RightParen)
		return &ListType{
			Element:  elemType,
			Position: token}
	case lexer.Identifier, lexer.KwInt, lexer.KwFloat, lexer.KwBool, lexer.KwString, lexer.KwChar, lexer.KwUnit:
		parser.advance()
		return &SimpleType{
			Name: token.Lexeme,
			Pos:  token}
	case lexer.EndOfFile:
		return nil
	default:
		e := parser.error(token, fmt.Sprintf("unexpected token in type: %q", token.Lexeme))
		parser.errors = append(parser.errors, e.Error())
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

func (parser *Parser) prev() lexer.Token {
	if parser.position != 0 {
		return parser.tokens[parser.position-1]
	}
	return lexer.Token{}
}

func (parser *Parser) advance() lexer.Token {
	token := parser.cur()
	parser.position++
	return token
}

func (parser *Parser) expect(typ lexer.TokenType) lexer.Token {
	token := parser.cur()
	if token.Type != typ {
		err := parser.error(token, fmt.Sprintf("expected %q, got %q", typ, token.Type))
		parser.errors = append(parser.errors, err.Error())
		return token
	}
	parser.advance()
	return token
}
