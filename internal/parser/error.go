package parser

import (
	"fmt"
	"lunno/internal/lexer"
	"strings"
)

func (parser *Parser) error(token lexer.Token, msg string) error {
	lineText := getLineText(parser.source, token.Line)
	caret := makeCaret(token.Column)

	fmt.Printf(
		"error: %s\n  --> %s:%d:%d\n   |\n%2d | %s\n   | %s\n",
		msg,
		token.File, token.Line, token.Column,
		token.Line, lineText, caret,
	)

	return fmt.Errorf("%s:%d:%d: %s", token.File, token.Line, token.Column, msg)
}

func getLineText(source []byte, lineNum int) string {
	if lineNum < 1 {
		return ""
	}
	start := 0
	curLine := 1
	for i, r := range source {
		if curLine == lineNum {
			start = i
			break
		}
		if r == '\n' {
			curLine++
		}
	}
	end := len(source)
	for i := start; i < len(source); i++ {
		if source[i] == '\n' {
			end = i
			break
		}
	}
	return string(source[start:end])
}

func makeCaret(column int) string {
	spaces := column - 1
	if spaces < 0 {
		spaces = 0
	}
	return fmt.Sprintf("%s^", strings.Repeat(" ", spaces))
}
