package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"ucal/server"
)

type stdio struct {
	io.Reader
	io.Writer
}

func (stdio) Close() error { return nil }

func main() {
	serve := flag.Bool("serve", false, "Run in server mode")
	flag.Parse()
	if *serve {
		if err := server.Serve(context.Background(), stdio{os.Stdin, os.Stdout}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	expression := strings.Join(os.Args[1:], " ")
	if err := printHistogram(expression); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
