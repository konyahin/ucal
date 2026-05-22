package parser

import (
	"context"
	"math"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expression string
		wantLeft   float64
		wantRight  float64
		errWant    string
	}{
		{expression: "2 + 2", wantLeft: 4, wantRight: 4},
		{expression: "3 - 2", wantLeft: 1, wantRight: 1},
		{expression: "3 * 2", wantLeft: 6, wantRight: 6},
		{expression: "100 / 2", wantLeft: 50, wantRight: 50},
		{expression: "0.1 + 0.2", wantLeft: 0.3, wantRight: 0.3},
		{expression: "(10 + 6) / ((1 - 3) * 2)", wantLeft: -4, wantRight: -4},
		{expression: "- - (- 4)", wantLeft: -4, wantRight: -4},
		{expression: "2 + 3 * 4", wantLeft: 14, wantRight: 14},
		{expression: "0 ~ 10", wantLeft: 0, wantRight: 10},
		{expression: "100 / 4~6", wantLeft: 16.7, wantRight: 25},
		{expression: "10 + (4~6 - 5) * 5~10", wantLeft: 2.28, wantRight: 17.7},
		{expression: "1 ~ (2~10)", errWant: "range cannot contain another range at position 6"},
		{expression: "4 / 0", errWant: errDivisionByZero.Error()},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			node, err := New(test.expression).Parse()

			var result Result
			if err == nil {
				result, err = Eval(context.Background(), node)
			}

			errGot := ""
			if err != nil {
				errGot = err.Error()
			}

			if errGot != test.errWant {
				t.Fatalf("eval return '%v', but we want '%v'", err, test.errWant)
			}

			if test.errWant != "" {
				return
			}

			result25 := result.Percentile(2.5)
			if !closeEnough(result25, test.wantLeft) {
				t.Fatalf("2.5 percentile is %g, but we want %g", result25, test.wantLeft)
			}

			result975 := result.Percentile(97.5)
			if !closeEnough(result975, test.wantRight) {
				t.Fatalf("97.5 percentile is %g, but we want %g", result975, test.wantRight)
			}
		})
	}
}

func closeEnough(a, b float64) bool {
	return math.Abs(a-b) <= 0.1
}
