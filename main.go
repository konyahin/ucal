package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"ucal/parser"
)

func main() {
	expression := strings.Join(os.Args[1:], " ")
	node, err := parser.New(expression).Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error from parser:\n%s\n", err)
		os.Exit(1)
	}

	result, err := parser.Eval(node)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error from evaluator:\n%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Result: %.0f~%.0f\n", math.Round(result.Percentile(2.5)), math.Round(result.Percentile(97.5)))
	for _, percentile := range []float64{2.5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 97.5} {
		fmt.Printf("%4.1f%%\t%.2f\n", percentile, result.Percentile(percentile))
	}
}
