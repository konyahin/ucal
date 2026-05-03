package parser

import "fmt"

type Parser struct {
	lexer Lexer
}

type Node interface {
	position() int
}

type NumberNode struct {
	value float64
	pos   int
}

type BinaryOperationNode struct {
	operation TokenType
	left      Node
	right     Node
}

func (n *NumberNode) position() int {
	return n.pos
}

func (n *BinaryOperationNode) position() int {
	return n.left.position()
}

func New(input string) *Parser {
	p := &Parser{
		lexer: *newLexer(input),
	}
	return p
}

func (p *Parser) Parse() (Node, error) {
	left, err := p.parseNumber()
	if err != nil {
		return nil, err
	}

	op := p.lexer.Next()
	if op.Type == EOF {
		return left, nil
	}
	if op.Type != Plus && op.Type != Minus {
		return nil, fmt.Errorf("Unexpected token at %d: %s", op.Position, op.Literal)
	}

	right, err := p.parseNumber()
	if err != nil {
		return nil, err
	}

	end := p.lexer.Next()
	if end.Type != EOF {
		return nil, fmt.Errorf("Unexpected token after expression: %s", end.Literal)
	}

	return &BinaryOperationNode{
		operation: op.Type,
		left:      left,
		right:     right,
	}, nil
}

func (p *Parser) parseNumber() (Node, error) {
	token := p.lexer.Next()
	if token.Type != Number {
		return nil, fmt.Errorf("Expected number, got: %s", token.Literal)
	}
	return &NumberNode{token.Value, token.Position}, nil
}
