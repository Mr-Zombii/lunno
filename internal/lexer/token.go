package lexer

import (
	"fmt"
	"lunno/internal/diagnostics"
	"strings"
)

type TokenType int

type Token struct {
	Type        TokenType
	Lexeme      string
	Line        int
	Column      int
	Indentation int
	File        string
}

const (
	Illegal TokenType = iota
	Newline

	Identifier
	Integer
	Float
	String
	Char
	Bool

	LeftParen
	RightParen
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket

	Plus
	Minus
	Asterisk
	Slash

	Assign
	Equal
	NotEqual
	LessThan
	GreaterThan
	LessThanOrEqual
	GreaterThanOrEqual

	Comma
	Colon
	Arrow

	KwLet
	KwRec
	KwFn
	KwIf
	KwThen
	KwElse
	KwImport
	KwInt
	KwFloat
	KwString
	KwChar
	KwBool
	KwList
	KwUnit

	EndOfFile
)

var keywords = map[string]TokenType{
	"let":    KwLet,
	"rec":    KwRec,
	"fn":     KwFn,
	"if":     KwIf,
	"then":   KwThen,
	"else":   KwElse,
	"import": KwImport,
	"int":    KwInt,
	"float":  KwFloat,
	"string": KwString,
	"char":   KwChar,
	"bool":   KwBool,
	"list":   KwList,
	"unit":   KwUnit,
	"true":   Bool,
	"false":  Bool,
}

var precedences = map[TokenType]int{
	Colon:    1,
	Equal:    2,
	NotEqual: 2,

	LessThan:           3,
	GreaterThan:        3,
	LessThanOrEqual:    3,
	GreaterThanOrEqual: 3,

	Plus:     4,
	Minus:    4,
	Asterisk: 5,
	Slash:    5,
}

var multiCharOperators = map[string]TokenType{
	"==": Equal,
	"!=": NotEqual,
	"<=": LessThanOrEqual,
	">=": GreaterThanOrEqual,
	"->": Arrow,
}

var singleCharTokens = map[byte]TokenType{
	'(': LeftParen, ')': RightParen,
	'[': LeftBracket, ']': RightBracket,
	'{': LeftBrace, '}': RightBrace,
	'+': Plus, '-': Minus,
	'*': Asterisk, '/': Slash,
	':': Colon, ',': Comma,
	'=': Assign, '\n': Newline,
}

type bracePair struct {
	ClosingChar   string
	ClosingTT     TokenType
	Name          string
	ValidClosings map[TokenType]bool
}

var Braces = map[TokenType]*bracePair{
	LeftBracket: {
		ClosingChar: "]",
		ClosingTT:   RightBracket,
		Name:        "list",
		ValidClosings: map[TokenType]bool{
			RightBracket: true,
			RightBrace:   false,
			RightParen:   false,
		},
	},
	LeftBrace: {
		ClosingChar: "}",
		ClosingTT:   RightBrace,
		Name:        "tuple",
		ValidClosings: map[TokenType]bool{
			RightBracket: false,
			RightBrace:   true,
			RightParen:   false,
		},
	},
	LeftParen: {
		ClosingChar: ")",
		ClosingTT:   RightParen,
		Name:        "argument list",
		ValidClosings: map[TokenType]bool{
			RightBracket: false,
			RightBrace:   false,
			RightParen:   true,
		},
	},
}

func (t Token) String() string {
	return fmt.Sprintf("%d('%s') at %d:%d", t.Type, strings.ReplaceAll(t.Lexeme, "\n", "\\n"), t.Line, t.Column)
}

func (t Token) Span() diagnostics.Span {
	return diagnostics.Span{
		File:   t.File,
		Line:   t.Line,
		Column: t.Column,
	}
}

func (typ TokenType) Precedence() int {
	if precedence, ok := precedences[typ]; ok {
		return precedence
	}
	return -1
}
