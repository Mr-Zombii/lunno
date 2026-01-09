package lexer

import "strings"

func classify(ch rune) CharClass {
	switch {
	case ch == 0:
		return CC_EOF
	case ch == '\n':
		return CC_Newline
	case ch == ' ' || ch == '\t' || ch == '\r':
		return CC_Whitespace
	case ch == '_':
		return CC_Underscore
	case ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z'):
		return CC_Letter
	case '0' <= ch && ch <= '9':
		return CC_Digit
	case ch == '.':
		return CC_Dot
	case ch == '"':
		return CC_Quote
	case ch == '\'':
		return CC_Apostrophe
	case ch == '\\':
		return CC_Backslash
	case strings.ContainsRune("+-*/=<>!:,()[]{}", ch):
		return CC_Operator
	default:
		return CC_Other
	}
}
