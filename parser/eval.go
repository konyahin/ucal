package parser

import (
	"errors"
	"fmt"
)

var errDivisionByZero = errors.New("division by zero")

func Eval(node Node) (float64, error) {
	switch n := node.(type) {
	case *numberNode:
		return n.value, nil
	case *unaryMinusNode:
		v, err := Eval(n.value)
		if err != nil {
			return 0, err
		}
		return -v, nil
	case *binaryOperationNode:
		left, err := Eval(n.left)
		if err != nil {
			return 0, err
		}
		right, err := Eval(n.right)
		if err != nil {
			return 0, err
		}
		switch n.operation.kind {
		case Plus:
			return left + right, nil
		case Minus:
			return left - right, nil
		case Asterisk:
			return left * right, nil
		case Slash:
			if right == 0 {
				return 0, errDivisionByZero
			}
			return left / right, nil
		case Tilde:
			// todo make real range calculation
			return (left + right) / 2, nil
		default:
			panic(fmt.Sprintf("unknown binary operation: %v", n.operation))
		}
	default:
		panic(fmt.Sprintf("unknown node type: %T", node))
	}
}
