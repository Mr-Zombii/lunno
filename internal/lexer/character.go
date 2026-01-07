package lexer

type CharClass uint8

const (
	CC_EOF CharClass = iota
	CC_Whitespace
	CC_Newline
	CC_Letter
	CC_Digit
	CC_Underscore
	CC_Dot
	CC_Quote
	CC_Apostrophe
	CC_Backslash
	CC_Operator
	CC_Other
)
