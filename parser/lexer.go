package parser

import (
	"errors"
	"strconv"
)

type lexer struct {
	input string
	pos   int
}

type tokenType int

const (
	eof tokenType = iota
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

var errWrongNumber = errors.New("wrong number")
var errUnknownChar = errors.New("unknown char")

type token struct {
	kind     tokenType
	literal  string
	value    float64
	position int
	err      error
}

var singleByteTokens = map[byte]tokenType{
	'+': Plus,
	'-': Minus,
	'*': Asterisk,
	'/': Slash,
	'~': Tilde,
	'(': LeftParen,
	')': RightParen,
}

func newLexer(expression string) lexer {
	return lexer{
		input: expression,
	}
}

func (l *lexer) next() token {
	l.skipWhiteSpace()

	if l.isEnd() {
		return token{
			kind:     eof,
			position: len(l.input),
		}
	}

	b := l.input[l.pos]
	if t, ok := singleByteTokens[b]; ok {
		return l.singleByteToken(t)
	}

	if isDigit(b) {
		pos := l.pos
		literal := l.readNumber()
		value, err := strconv.ParseFloat(literal, 64)
		if err != nil {
			return l.errorToken(errWrongNumber, pos, literal)
		}

		return token{
			kind:     Number,
			position: pos,
			value:    value,
			literal:  literal,
		}
	}

	tok := l.errorToken(errUnknownChar, l.pos, string(b))
	l.pos++
	return tok
}

func (l *lexer) errorToken(err error, pos int, literal string) token {
	return token{
		kind:     Error,
		position: pos,
		err:      err,
		literal:  literal,
	}
}

func (l *lexer) singleByteToken(t tokenType) token {
	tok := token{kind: t, position: l.pos, literal: string(l.input[l.pos])}
	l.pos++
	return tok
}

func (l *lexer) readNumber() string {
	start := l.pos
	for !l.isEnd() && (isDigit(l.input[l.pos]) || l.input[l.pos] == '.') {
		l.pos++
	}
	return l.input[start:l.pos]
}

func (l *lexer) skipWhiteSpace() {
	for !l.isEnd() && isSpace(l.input[l.pos]) {
		l.pos++
	}
}

func (l *lexer) isEnd() bool {
	return l.pos >= len(l.input)
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
