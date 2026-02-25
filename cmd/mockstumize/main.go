package main

import (
	"flag"
	"fmt"
	"os"

	"mockstumize/pkg/builder"
	"mockstumize/pkg/printer"
)

func main() {
	path := flag.String("path", "examples/base", "path to kustomization directory")
	flag.Parse()

	result, err := builder.Build(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	output, err := printer.Print(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(output)
}
