package diagnostics

import (
	"fmt"
	"strings"
)

func getLineText(source []byte, line int) string {
	if line < 1 {
		return ""
	}
	curLine := 1
	start := 0
	for i, b := range source {
		if curLine == line {
			start = i
			break
		}
		if b == '\n' {
			curLine++
		}
	}
	if curLine != line {
		return ""
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
