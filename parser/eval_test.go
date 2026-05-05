package parser

import (
	"errors"
	"math"
	"testing"
)

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

func TestEval(t *testing.T) {
	tests := []struct {
		expression string
		want       float64
		errWant    error
	}{
		{expression: "2 + 2", want: 4},
		{expression: "3 - 2", want: 1},
		{expression: "3 * 2", want: 6},
		{expression: "100 / 2", want: 50},
		{expression: "0.1 + 0.2", want: 0.3},
		{expression: "(10 + 6) / ((1 - 3) * 2)", want: -4},
		{expression: "- - (- 4)", want: -4},
		{expression: "2 + 3 * 4", want: 14},
		{expression: "4 / 0", errWant: errDivisionByZero},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			node, err := New(test.expression).Parse()
			if err != nil {
				t.Fatal("parser return error:", err)
			}

			result, err := Eval(node)
			if !errors.Is(err, test.errWant) {
				t.Fatalf("eval return %v, but we want %v", err, test.errWant)
			}

			if test.errWant == nil && !floatEqual(result, test.want) {
				t.Fatalf("eval compute %g, but we want %g", result, test.want)
			}
		})
	}
}
