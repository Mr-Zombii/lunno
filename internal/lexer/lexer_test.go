package lexer_test

import (
	"lunno/internal/lexer"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []lexer.TokenType
		lexemes   []string
		expectErr bool
	}{
		{
			name:      "empty file",
			input:     "",
			expected:  []lexer.TokenType{},
			lexemes:   []string{},
			expectErr: true,
		},
		{
			name:     "simple identifiers",
			input:    "foo bar baz",
			expected: []lexer.TokenType{lexer.Identifier, lexer.Identifier, lexer.Identifier, lexer.EndOfFile},
			lexemes:  []string{"foo", "bar", "baz", ""},
		},
		{
			name:     "integers and floats",
			input:    "123 45.67",
			expected: []lexer.TokenType{lexer.Int, lexer.Float, lexer.EndOfFile},
			lexemes:  []string{"123", "45.67", ""},
		},
		{
			name:     "malformed float",
			input:    "12.",
			expected: []lexer.TokenType{lexer.Illegal, lexer.EndOfFile},
			lexemes:  []string{"12.", ""},
		},
		{
			name:     "strings",
			input:    `"hello" "a\nb"`,
			expected: []lexer.TokenType{lexer.String, lexer.String, lexer.EndOfFile},
			lexemes:  []string{`"hello"`, `"a\nb"`, ""},
		},
		{
			name:     "empty string",
			input:    `""`,
			expected: []lexer.TokenType{lexer.Illegal, lexer.EndOfFile},
			lexemes:  []string{`""`, ""},
		},
		{
			name:     "characters",
			input:    `'a' '\n'`,
			expected: []lexer.TokenType{lexer.Char, lexer.Char, lexer.EndOfFile},
			lexemes:  []string{`'a'`, `'\n'`, ""},
		},
		{
			name:     "empty char",
			input:    `''`,
			expected: []lexer.TokenType{lexer.Illegal, lexer.EndOfFile},
			lexemes:  []string{`''`, ""},
		},
		{
			name:     "invalid character",
			input:    "@",
			expected: []lexer.TokenType{lexer.Illegal, lexer.EndOfFile},
			lexemes:  []string{"@", ""},
		},
		{
			name:     "unterminated string",
			input:    `"hello`,
			expected: []lexer.TokenType{lexer.Illegal, lexer.EndOfFile},
			lexemes:  []string{`"hello`, ""},
		},
		{
			name:     "unterminated char",
			input:    `'x`,
			expected: []lexer.TokenType{lexer.Illegal, lexer.EndOfFile},
			lexemes:  []string{`'x`, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, tokens, err := lexer.Tokenize(tt.input, "test.ln")
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}
			for i, tok := range tokens {
				if tok.Type != tt.expected[i] {
					t.Errorf("token %d: expected type %v, got %v", i, tt.expected[i], tok.Type)
				}
				if tok.Lexeme != tt.lexemes[i] {
					t.Errorf("token %d: expected lexeme '%s', got '%s'", i, tt.lexemes[i], tok.Lexeme)
				}
				if tok.Line == 0 || tok.Column == 0 {
					t.Errorf("token %d: invalid line/column: %d/%d", i, tok.Line, tok.Column)
				}
			}
		})
	}
}
