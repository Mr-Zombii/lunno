package parser

import (
	"lunno/internal/lexer"
)

type Node interface {
	NodeType() string
}

type Expression interface {
	Node
	exprNode()
}

type TypeNode interface {
	Node
	typeNode()
}

type Program struct {
	Expressions []Expression
}

func (p *Program) NodeType() string {
	return "Program"
}

type Identifier struct {
	Name     string
	Position lexer.Token
}

func (i *Identifier) exprNode() {}
func (i *Identifier) NodeType() string {
	return "Identifier"
}

type IntegerLiteral struct {
	Value    int64
	Raw      string
	Position lexer.Token
}

func (i *IntegerLiteral) exprNode() {}
func (i *IntegerLiteral) NodeType() string {
	return "IntegerLiteral"
}

type FloatLiteral struct {
	Value    float64
	Raw      string
	Position lexer.Token
}

func (f *FloatLiteral) exprNode() {}
func (f *FloatLiteral) NodeType() string {
	return "FloatLiteral"
}

type StringLiteral struct {
	Value    string
	Position lexer.Token
}

func (s *StringLiteral) exprNode() {}
func (s *StringLiteral) NodeType() string {
	return "StringLiteral"
}

type CharacterLiteral struct {
	Value    byte
	Raw      string
	Position lexer.Token
}

func (c *CharacterLiteral) exprNode() {}
func (c *CharacterLiteral) NodeType() string {
	return "CharacterLiteral"
}

type BooleanLiteral struct {
	Value    bool
	Position lexer.Token
}

func (b *BooleanLiteral) exprNode() {}
func (b *BooleanLiteral) NodeType() string {
	return "BooleanLiteral"
}

type UnitLiteral struct {
	Position lexer.Token
}

func (u *UnitLiteral) exprNode() {}
func (u *UnitLiteral) NodeType() string {
	return "UnitLiteral"
}

type ListExpression struct {
	Elements []Expression
	Position lexer.Token
}

func (l *ListExpression) exprNode() {}
func (l *ListExpression) NodeType() string {
	return "ListExpression"
}

type IndexExpression struct {
	Target   Expression
	Index    Expression
	Position lexer.Token
}

func (i *IndexExpression) exprNode() {}
func (i *IndexExpression) NodeType() string {
	return "IndexExpression"
}

type PrefixExpression struct {
	Operator lexer.Token
	Right    Expression
	Position lexer.Token
}

func (p *PrefixExpression) exprNode() {}
func (p *PrefixExpression) NodeType() string {
	return "PrefixExpression"
}

type InfixExpression struct {
	Left     Expression
	Operator lexer.Token
	Right    Expression
	Position lexer.Token
}

func (i *InfixExpression) exprNode() {}
func (i *InfixExpression) NodeType() string {
	return "InfixExpression"
}

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
	Position  lexer.Token
}

func (c *CallExpression) exprNode() {}
func (c *CallExpression) NodeType() string {
	return "CallExpression"
}

type VariableDeclarationExpression struct {
	Name      lexer.Token
	Type      TypeNode
	Value     Expression
	Recursive bool
	Position  lexer.Token
}

func (v *VariableDeclarationExpression) exprNode() {}
func (v *VariableDeclarationExpression) NodeType() string {
	return "VariableDeclarationExpression"
}

type Parameter struct {
	Name     lexer.Token
	Type     TypeNode
	Position lexer.Token
}

type FunctionLiteralExpression struct {
	Parameters []Parameter
	Body       Expression
	Position   lexer.Token
}

func (f *FunctionLiteralExpression) exprNode() {}
func (f *FunctionLiteralExpression) NodeType() string {
	return "FunctionLiteralExpression"
}

type FunctionDeclarationExpression struct {
	Name      lexer.Token
	Recursive bool
	Signature TypeNode
	Function  *FunctionLiteralExpression
	Position  lexer.Token
}

func (f *FunctionDeclarationExpression) exprNode() {}
func (f *FunctionDeclarationExpression) NodeType() string {
	return "FunctionDeclarationExpression"
}

type BlockExpression struct {
	Expressions []Expression
	Position    lexer.Token
}

func (b *BlockExpression) exprNode() {}
func (b *BlockExpression) NodeType() string {
	return "BlockExpression"
}

type IfExpression struct {
	Condition Expression
	Then      Expression
	Else      Expression
	Position  lexer.Token
}

func (i *IfExpression) exprNode() {}
func (i *IfExpression) NodeType() string {
	return "IfExpression"
}

type ImportExpression struct {
	Module   string
	Position lexer.Token
}

func (i *ImportExpression) exprNode() {}
func (i *ImportExpression) NodeType() string {
	return "ImportExpression"
}

type SliceExpression struct {
	Target   Expression
	Start    Expression
	End      Expression
	Position lexer.Token
}

func (s *SliceExpression) exprNode() {}
func (s *SliceExpression) NodeType() string {
	return "SliceExpression"
}

type SimpleType struct {
	Name string
	Pos  lexer.Token
}

func (s *SimpleType) typeNode() {}
func (s *SimpleType) NodeType() string {
	return "SimpleType"
}

type ListType struct {
	Element  TypeNode
	Position lexer.Token
}

func (l *ListType) typeNode() {}
func (l *ListType) NodeType() string {
	return "ListType"
}

type FunctionType struct {
	Parameters []TypeNode
	Return     TypeNode
	Position   lexer.Token
}

func (f *FunctionType) typeNode() {}
func (f *FunctionType) NodeType() string {
	return "FunctionType"
}
