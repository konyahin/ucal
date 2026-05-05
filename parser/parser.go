package parser

import "fmt"

type Parser struct {
	lexer   lexer
	current token
}

type Node interface {
	position() int
}

type numberNode struct {
	value float64
	pos   int
}

type binaryOperationNode struct {
	operation tokenType
	left      Node
	right     Node
}

type unaryMinusNode struct {
	value Node
	pos   int
}

func (n *numberNode) position() int {
	return n.pos
}

func (n *binaryOperationNode) position() int {
	return n.left.position()
}

func (n *unaryMinusNode) position() int {
	return n.pos
}

func New(input string) *Parser {
	p := &Parser{lexer: newLexer(input)}
	p.advance()
	return p
}

func (p *Parser) advance() {
	p.current = p.lexer.next()
}

func (p *Parser) Parse() (Node, error) {
	node, err := p.parseAddSub()
	if err != nil {
		return nil, err
	}
	if !p.is(eof) {
		return nil, fmt.Errorf("unexpected token after expression: %s", p.current.literal)
	}
	return node, nil
}

func (p *Parser) parseAddSub() (Node, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}

	for p.is(Plus) || p.is(Minus) {
		op := p.current
		p.advance()
		right, err := p.parseMulDiv()
		if err != nil {
			return nil, err
		}
		left = &binaryOperationNode{operation: op.kind, left: left, right: right}
	}
	return left, nil
}

func (p *Parser) parseMulDiv() (Node, error) {
	left, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	for p.is(Asterisk) || p.is(Slash) {
		op := p.current
		p.advance()
		right, err := p.parseAtom()
		if err != nil {
			return nil, err
		}
		left = &binaryOperationNode{operation: op.kind, left: left, right: right}
	}
	return left, nil
}

func (p *Parser) parseAtom() (Node, error) {
	if p.is(Minus) {
		pos := p.current.position
		p.advance()
		operand, err := p.parseAtom()
		if err != nil {
			return nil, err
		}
		return &unaryMinusNode{value: operand, pos: pos}, nil
	}

	if p.is(LeftParen) {
		p.advance()
		node, err := p.parseAddSub()
		if err != nil {
			return nil, err
		}
		if !p.is(RightParen) {
			return nil, fmt.Errorf("expected ')', got: %s", p.current.literal)
		}
		p.advance()
		return node, nil
	}

	if !p.is(Number) {
		return nil, fmt.Errorf("expected number, got: %s", p.current.literal)
	}
	node := &numberNode{p.current.value, p.current.position}
	p.advance()
	return node, nil
}

func (p *Parser) is(t tokenType) bool {
	return p.current.kind == t
}
