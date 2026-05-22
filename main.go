package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"

	"ucal/parser"
)

const (
	bins     = 20
	barWidth = 40
)

func main() {
	expression := strings.Join(os.Args[1:], " ")
	node, err := parser.New(expression).Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error from parser:\n%s\n", err)
		os.Exit(1)
	}

	result, err := parser.Eval(context.Background(), node)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error from evaluator:\n%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Result: %.0f~%.0f\n", math.Round(result.Percentile(2.5)), math.Round(result.Percentile(97.5)))

	low := result.Percentile(1)
	high := result.Percentile(99)
	if high-low < 1e-12 {
		return
	}

	step := (high - low) / bins
	masses := make([]float64, bins)
	maxMass := 0.0

	current := low
	prev := result.CDF(current)
	for i := range bins {
		current += step
		cdf := result.CDF(current)
		mass := cdf - prev
		if mass < 0 {
			mass = 0
		}
		masses[i] = mass
		if mass > maxMass {
			maxMass = mass
		}
		prev = cdf
	}

	for i, mass := range masses {
		edge := low + step*float64(i)
		width := int(math.Round(mass / maxMass * barWidth))
		fmt.Printf("%8.2f │%s\n", edge, strings.Repeat("█", width))
	}
}
