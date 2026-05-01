package main

import (
	"fmt"
	"os"
	"strings"

	"ucal/parser"
)

func main() {
	expression := strings.Join(os.Args[1:], " ")
	l := parser.New(expression)

	for {
		token := l.Next()
		fmt.Println(token)

		if token.Type == parser.EOF || token.Type == parser.Error {
			break
		}
	}
}
