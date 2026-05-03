package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func collectAll(l *Lexer) []Token {
	tokens := make([]Token, 0, 10)
	for {
		token := l.Next()
		tokens = append(tokens, token)
		if token.Type == EOF || token.Type == Error {
			break
		}
	}
	return tokens
}

func TestLexer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "happy path",
			input: "4.5 * (2~5 + 16) - 32.45 / 3",
			want: []Token{
				{Type: Number, Literal: "4.5", Value: 4.5, Position: 0},
				{Type: Asterisk, Literal: "*", Position: 4},
				{Type: LeftParen, Literal: "(", Position: 6},
				{Type: Number, Literal: "2", Value: 2, Position: 7},
				{Type: Tilde, Literal: "~", Position: 8},
				{Type: Number, Literal: "5", Value: 5, Position: 9},
				{Type: Plus, Literal: "+", Position: 11},
				{Type: Number, Literal: "16", Value: 16, Position: 13},
				{Type: RightParen, Literal: ")", Position: 15},
				{Type: Minus, Literal: "-", Position: 17},
				{Type: Number, Literal: "32.45", Value: 32.45, Position: 19},
				{Type: Slash, Literal: "/", Position: 25},
				{Type: Number, Literal: "3", Value: 3, Position: 27},
				{Type: EOF, Position: 28},
			},
		},
		{
			name:  "empty input",
			input: "",
			want: []Token{
				{Type: EOF, Position: 0},
			},
		},
		{
			name:  "space only",
			input: " \t\r",
			want: []Token{
				{Type: EOF, Position: 3},
			},
		},
		{
			name:  "new line handling",
			input: "1\n2",
			want: []Token{
				{Type: Number, Literal: "1", Value: 1, Position: 0},
				{Type: Number, Literal: "2", Value: 2, Position: 2},
				{Type: EOF, Position: 3},
			},
		},
		{
			name:  "unknown symbol",
			input: "4.5 ,",
			want: []Token{
				{Type: Number, Literal: "4.5", Value: 4.5, Position: 0},
				{Type: Error, Literal: ",", Position: 4, Err: ErrUnknownChar},
			},
		},
		{
			name:  "wrong number",
			input: "4.5.3",
			want: []Token{
				{Type: Error, Literal: "4.5.3", Position: 0, Err: ErrWrongNumber},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := collectAll(newLexer(tc.input))
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
