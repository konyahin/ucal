package parser

import (
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
		errWant    string
	}{
		{expression: "2 + 2", want: 4},
		{expression: "3 - 2", want: 1},
		{expression: "3 * 2", want: 6},
		{expression: "100 / 2", want: 50},
		{expression: "0.1 + 0.2", want: 0.3},
		{expression: "(10 + 6) / ((1 - 3) * 2)", want: -4},
		{expression: "- - (- 4)", want: -4},
		{expression: "2 + 3 * 4", want: 14},
		{expression: "0 ~ 10", want: 5},
		{expression: "100/4~6", want: 20},
		{expression: "1 ~ (2~10)", errWant: "range cannot contain another range at position 6"},
		{expression: "4 / 0", errWant: errDivisionByZero.Error()},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			node, err := New(test.expression).Parse()

			var result float64
			if err == nil {
				result, err = Eval(node)
			}

			errGot := ""
			if err != nil {
				errGot = err.Error()
			}

			if errGot != test.errWant {
				t.Fatalf("eval return '%v', but we want '%v'", err, test.errWant)
			}

			if test.errWant == "" && !floatEqual(result, test.want) {
				t.Fatalf("eval compute %g, but we want %g", result, test.want)
			}
		})
	}
}
