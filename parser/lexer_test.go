package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func collectAll(l lexer) []token {
	tokens := make([]token, 0, 10)
	for {
		tok := l.next()
		tokens = append(tokens, tok)
		if tok.kind == eof || tok.kind == Error {
			break
		}
	}
	return tokens
}

func TestLexer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []token
	}{
		{
			name:  "happy path",
			input: "4.5 * (2~5 + 16) - 32.45 / 3",
			want: []token{
				{kind: Number, literal: "4.5", value: 4.5, position: 0},
				{kind: Asterisk, literal: "*", position: 4},
				{kind: LeftParen, literal: "(", position: 6},
				{kind: Number, literal: "2", value: 2, position: 7},
				{kind: Tilde, literal: "~", position: 8},
				{kind: Number, literal: "5", value: 5, position: 9},
				{kind: Plus, literal: "+", position: 11},
				{kind: Number, literal: "16", value: 16, position: 13},
				{kind: RightParen, literal: ")", position: 15},
				{kind: Minus, literal: "-", position: 17},
				{kind: Number, literal: "32.45", value: 32.45, position: 19},
				{kind: Slash, literal: "/", position: 25},
				{kind: Number, literal: "3", value: 3, position: 27},
				{kind: eof, position: 28},
			},
		},
		{
			name:  "empty input",
			input: "",
			want: []token{
				{kind: eof, position: 0},
			},
		},
		{
			name:  "space only",
			input: " \t\r",
			want: []token{
				{kind: eof, position: 3},
			},
		},
		{
			name:  "new line handling",
			input: "1\n2",
			want: []token{
				{kind: Number, literal: "1", value: 1, position: 0},
				{kind: Number, literal: "2", value: 2, position: 2},
				{kind: eof, position: 3},
			},
		},
		{
			name:  "unknown symbol",
			input: "4.5 ,",
			want: []token{
				{kind: Number, literal: "4.5", value: 4.5, position: 0},
				{kind: Error, literal: ",", position: 4, err: errUnknownChar},
			},
		},
		{
			name:  "wrong number",
			input: "4.5.3",
			want: []token{
				{kind: Error, literal: "4.5.3", position: 0, err: errWrongNumber},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := collectAll(newLexer(tc.input))
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateErrors(), cmp.AllowUnexported(token{})); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
