package main

import (
	"fmt"
	"os"
	"strconv"
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

	fmt.Printf("Result: %s\n", strconv.FormatFloat(result, 'f', -1, 64))
}
