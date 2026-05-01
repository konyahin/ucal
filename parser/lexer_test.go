package parser

import (
	"reflect"
	"testing"
)

func (l *Lexer) All() []Token {
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

func TestHappyPath(t *testing.T) {
	lexer := New("4.5 * (2~5 + 16) - 32.45 / 3")
	tokens := lexer.All()

	expect := []Token{
		{
			Type:     Number,
			Literal:  "4.5",
			Value:    4.5,
			Position: 0,
			Err:      nil,
		},
		{
			Type:     Asterisk,
			Literal:  "*",
			Value:    0,
			Position: 4,
			Err:      nil,
		},
		{
			Type:     LeftParen,
			Literal:  "(",
			Value:    0,
			Position: 6,
			Err:      nil,
		},
		{
			Type:     Number,
			Literal:  "2",
			Value:    2,
			Position: 7,
			Err:      nil,
		},
		{
			Type:     Tilde,
			Literal:  "~",
			Value:    0,
			Position: 8,
			Err:      nil,
		},
		{
			Type:     Number,
			Literal:  "5",
			Value:    5,
			Position: 9,
			Err:      nil,
		},
		{
			Type:     Plus,
			Literal:  "+",
			Value:    0,
			Position: 11,
			Err:      nil,
		},
		{
			Type:     Number,
			Literal:  "16",
			Value:    16,
			Position: 13,
			Err:      nil,
		},
		{
			Type:     RightParen,
			Literal:  ")",
			Value:    0,
			Position: 15,
			Err:      nil,
		},
		{
			Type:     Minus,
			Literal:  "-",
			Value:    0,
			Position: 17,
			Err:      nil,
		},
		{
			Type:     Number,
			Literal:  "32.45",
			Value:    32.45,
			Position: 19,
			Err:      nil,
		},
		{
			Type:     Slash,
			Literal:  "/",
			Value:    0,
			Position: 25,
			Err:      nil,
		},
		{
			Type:     Number,
			Literal:  "3",
			Value:    3,
			Position: 27,
			Err:      nil,
		},
		{
			Type:     EOF,
			Literal:  "",
			Value:    0,
			Position: 28,
			Err:      nil,
		},
	}

	if !reflect.DeepEqual(expect, tokens) {
		t.Errorf("\nExpect:\n\t%#v\nGot:\n\t%#v", expect, tokens)
	}
}

func TestEmptyInput(t *testing.T) {
	lexer := New("")
	tokens := lexer.All()

	expect := []Token{
		{
			Type:     EOF,
			Literal:  "",
			Value:    0,
			Position: 0,
			Err:      nil,
		},
	}

	if !reflect.DeepEqual(expect, tokens) {
		t.Errorf("\nExpect:\n\t%#v\nGot:\n\t%#v", expect, tokens)
	}

}

func TestUnknownSymbol(t *testing.T) {
	lexer := New("4.5 ,")
	tokens := lexer.All()

	expect := []Token{
		{
			Type:     Number,
			Literal:  "4.5",
			Value:    4.5,
			Position: 0,
			Err:      nil,
		},
		{
			Type:     Error,
			Literal:  ",",
			Value:    0,
			Position: 4,
			Err:      ErrUnknownChar,
		},
	}

	if !reflect.DeepEqual(expect, tokens) {
		t.Errorf("\nExpect:\n\t%#v\nGot:\n\t%#v", expect, tokens)
	}
}

func TestWrongNumber(t *testing.T) {
	lexer := New("4.5.3")
	tokens := lexer.All()

	expect := []Token{
		{
			Type:     Error,
			Literal:  "4.5.3",
			Value:    0,
			Position: 0,
			Err:      ErrWrongNumber,
		},
	}

	if !reflect.DeepEqual(expect, tokens) {
		t.Errorf("\nExpect:\n\t%#v\nGot:\n\t%#v", expect, tokens)
	}
}