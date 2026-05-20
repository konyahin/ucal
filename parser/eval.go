package parser

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"ucal/distribution"
	"ucal/montecarlo"
)

var errDivisionByZero = errors.New("division by zero")

type evaluationResult float64

type Result interface {
	Percentile(percentile float64) float64
}

// deterministic value — same for every percentile
func (value evaluationResult) Percentile(_ float64) float64 {
	return float64(value)
}

func Eval(node Node) (Result, error) {
	tilde := findOperation(node, Tilde)
	if tilde == nil {
		res, err := evaluate(node, nil)
		if err != nil {
			return nil, err
		}
		return evaluationResult(res), nil
	}
	mc := montecarlo.New(func(r *rand.Rand) (float64, error) {
		return evaluate(node, r)
	})
	return mc.Run()
}

func evaluate(node Node, r *rand.Rand) (float64, error) {
	switch n := node.(type) {
	case *numberNode:
		return n.value, nil
	case *unaryMinusNode:
		v, err := evaluate(n.value, r)
		if err != nil {
			return 0, err
		}
		return -v, nil
	case *binaryOperationNode:
		left, err := evaluate(n.left, r)
		if err != nil {
			return 0, err
		}
		right, err := evaluate(n.right, r)
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
			return distribution.Normal(left, right, r), nil
		default:
			panic(fmt.Sprintf("unknown binary operation: %v", n.operation))
		}
	default:
		panic(fmt.Sprintf("unknown node type: %T", node))
	}
}
