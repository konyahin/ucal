package server

import (
	"context"
	"ucal/parser"
)

type EvaluateParams struct {
	Expression string `json:"expression"`
}

type EvaluateResult struct {
	Low  float64 `json:"low"`
	High float64 `json:"high"`
}

func evaluate(ctx context.Context, p EvaluateParams) (EvaluateResult, error) {
	node, err := parser.New(p.Expression).Parse()
	if err != nil {
		return EvaluateResult{}, err
	}
	result, err := parser.Eval(ctx, node)
	if err != nil {
		return EvaluateResult{}, err
	}
	return EvaluateResult{
		Low:  result.Percentile(2.5),
		High: result.Percentile(97.5),
	}, nil
}
