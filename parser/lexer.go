package parser

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type Lexer struct {
	input []rune
	pos   int
}

type TokenType int

const (
	EOF TokenType = iota
	Error
	Number
	Plus
	Minus
	Asterisk
	Slash
	Tilde
	LeftParen
	RightParen
)

var ErrWrongNumber = errors.New("Wrong number")
var ErrUnknownChar = errors.New("Unknown char")

type Token struct {
	Type     TokenType
	Literal  string
	Value    float64
	Position int
	Err      error
}

var tokenNames = map[TokenType]string{
	EOF:        "EOF",
	Error:      "Error",
	Number:     "Number",
	Plus:       "Plus",
	Minus:      "Minus",
	Asterisk:   "Asterisk",
	Slash:      "Slash",
	Tilde:      "Tilde",
	LeftParen:  "LeftParen",
	RightParen: "RightParen",
}

var singleRuneTokens = map[rune]TokenType{
	'+': Plus,
	'-': Minus,
	'*': Asterisk,
	'/': Slash,
	'~': Tilde,
	'(': LeftParen,
	')': RightParen,
}

func (t *Token) String() string {
	name, ok := tokenNames[t.Type]
	if !ok {
		name = "Unknown"
	}
	if t.Type == Number {
		return fmt.Sprintf("%s[%d] %s", name, t.Position, strconv.FormatFloat(t.Value, 'f', -1, 64))
	}
	if t.Type == Error {
		return fmt.Sprintf("%s[%d] %s: %s", name, t.Position, t.Err, t.Literal)
	}
	return fmt.Sprintf("%s[%d] %s", name, t.Position, t.Literal)
}

func New(expression string) *Lexer {
	return &Lexer{
		input: []rune(expression),
	}
}

func (l *Lexer) Next() Token {
	l.skipWhiteSpace()

	if l.isEnd() {
		return Token{
			Type:     EOF,
			Position: len(l.input),
		}
	}

	r := l.input[l.pos]
	if t, ok := singleRuneTokens[r]; ok {
		return l.singleRuneToken(t)
	}

	if unicode.IsDigit(r) {
		pos := l.pos
		literal := l.readNumber()
		value, err := strconv.ParseFloat(literal, 64)
		if err != nil {
			return l.errorToken(ErrWrongNumber, pos, literal)
		}

		return Token{
			Type:     Number,
			Position: pos,
			Value:    value,
			Literal:  literal,
		}
	}

	tok := l.errorToken(ErrUnknownChar, l.pos, string(l.input[l.pos]))
	l.pos++
	return tok
}

func (l *Lexer) errorToken(err error, pos int, literal string) Token {
	return Token{
		Type: Error,
		Position: pos,
		Err: err,
		Literal: literal,
	}
}

func (l *Lexer) singleRuneToken(t TokenType) Token {
	tok := Token{Type: t, Position: l.pos, Literal: string(l.input[l.pos])}
	l.pos++
	return tok
}

func (l *Lexer) readNumber() string {
	start := l.pos
	for !l.isEnd() && (unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '.') {
		l.pos++
	}
	return string(l.input[start:l.pos])
}

func (l *Lexer) skipWhiteSpace() {
	for !l.isEnd() && unicode.IsSpace(l.input[l.pos]) {
		l.pos++
	}
}

func (l *Lexer) isEnd() bool {
	return l.pos >= len(l.input)
}
