package diagnostics

import "fmt"

func Report(source []byte, span Span, msg string) error {
	lineText := getLineText(source, span.Line)
	caret := makeCaret(span.Column)

	fmt.Printf(
		"error: %s\n  --> %s:%d:%d\n   |\n%2d | %s\n   | %s\n",
		msg,
		span.File, span.Line, span.Column,
		span.Line, lineText, caret,
	)

	return fmt.Errorf(
		"%s:%d:%d: %s",
		span.File, span.Line, span.Column, msg,
	)
}
