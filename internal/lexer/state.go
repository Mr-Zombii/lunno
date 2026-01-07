package lexer

type State uint8

const (
	S_Start State = iota
	S_Ident
	S_Int
	S_Float
	S_String
	S_StringEsc
	S_Char
	S_CharEsc
	S_CharDone
	S_Operator
	S_Done
	S_Error
)

var dfa = map[State]map[CharClass]State{
	S_Start: {
		CC_Whitespace: S_Start,
		CC_Newline:    S_Start,
		CC_Letter:     S_Ident,
		CC_Underscore: S_Ident,
		CC_Digit:      S_Int,
		CC_Quote:      S_String,
		CC_Apostrophe: S_Char,
		CC_Operator:   S_Operator,
		CC_EOF:        S_Done,
	},

	S_Ident: {
		CC_Letter:     S_Ident,
		CC_Digit:      S_Ident,
		CC_Underscore: S_Ident,
	},

	S_Int: {
		CC_Digit: S_Int,
		CC_Dot:   S_Float,
	},

	S_Float: {
		CC_Digit: S_Float,
	},

	S_String: {
		CC_Backslash:  S_StringEsc,
		CC_Quote:      S_Done,
		CC_EOF:        S_Error,
		CC_Whitespace: S_String,
		CC_Newline:    S_String,
		CC_Letter:     S_String,
		CC_Digit:      S_String,
		CC_Operator:   S_String,
		CC_Underscore: S_String,
		CC_Dot:        S_String,
		CC_Other:      S_String,
		CC_Apostrophe: S_String,
	},

	S_StringEsc: {
		CC_Other: S_String,
	},

	S_Char: {
		CC_Backslash:  S_CharEsc,
		CC_Apostrophe: S_CharDone,
		CC_EOF:        S_Error,
		CC_Letter:     S_CharDone,
		CC_Digit:      S_CharDone,
		CC_Operator:   S_CharDone,
		CC_Other:      S_CharDone,
	},

	S_CharEsc: {
		CC_Other: S_CharDone,
		CC_EOF:   S_Error,
	},

	S_CharDone: {
		CC_Apostrophe: S_Done,
		CC_EOF:        S_Error,
		CC_Whitespace: S_Done,
		CC_Newline:    S_Done,
	},

	S_Operator: {
		CC_Operator: S_Operator,
	},
}

var accepting = map[State]func(*Lexer, string) TokenType{
	S_Ident: func(_ *Lexer, lex string) TokenType {
		return lookupIdentifier(lex)
	},

	S_Int: func(_ *Lexer, _ string) TokenType {
		return Int
	},

	S_Float: func(_ *Lexer, _ string) TokenType {
		return Float
	},

	S_String: func(_ *Lexer, _ string) TokenType {
		return String
	},

	S_Operator: func(_ *Lexer, lex string) TokenType {
		if tt, ok := multiCharOperators[lex]; ok {
			return tt
		}
		if len(lex) == 1 {
			if tt, ok := singleCharTokens[lex[0]]; ok {
				return tt
			}
		}
		return Illegal
	},

	S_CharDone: func(_ *Lexer, _ string) TokenType {
		return Char
	},

	S_Done: func(_ *Lexer, lex string) TokenType {
		if lex == "" {
			return EndOfFile
		}
		if lex[0] == '"' {
			return String
		}
		if lex[0] == '\'' {
			return Char
		}
		return Illegal
	},
}
