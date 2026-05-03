package parser

import "fmt"

func Eval(node Node) (float64, error) {
	switch n := node.(type) {
	case *NumberNode:
		return n.value, nil
	case *BinaryOperationNode:
		left, err := Eval(n.left)
		if err != nil {
			return 0, err
		}
		right, err := Eval(n.right)
		if err != nil {
			return 0, err
		}
		switch n.operation {
		case Plus:
			return left + right, nil
		case Minus:
			return left - right, nil
		default:
			panic(fmt.Sprintf("Unknown binary operation: %v", n.operation))
		}
	default:
		panic(fmt.Sprintf("Unknown node type: %T", node))
	}
}
