package lexer

import (
	"fmt"
	"strings"
)

func (l *Lexer) error(msg string) error {
	lineText := l.getLineText(l.line)
	caret := makeCaret(l.column)

	fmt.Printf(
		"error: %s\n  --> %s:%d:%d\n   |\n%2d | %s\n   | %s\n",
		msg,
		l.fileName, l.line, l.column,
		l.line, lineText, caret,
	)

	return fmt.Errorf("%s:%d:%d: %s", l.fileName, l.line, l.column, msg)
}

func (l *Lexer) getLineText(line int) string {
	if line < 1 {
		return ""
	}
	currentLine := 1
	start := 0
	for i, r := range l.Source {
		if currentLine == line {
			start = i
			break
		}
		if r == '\n' {
			currentLine++
		}
	}
	end := len(l.Source)
	for i := start; i < len(l.Source); i++ {
		if l.Source[i] == '\n' {
			end = i
			break
		}
	}
	return string(l.Source[start:end])
}

func makeCaret(column int) string {
	if column < 1 {
		column = 1
	}
	return fmt.Sprintf("%s^", strings.Repeat(" ", column-2))
}
