package main

import (
	"fmt"
	"os"

	"github.com/aquilax/slide2marpit"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: "+os.Args[0]+" <slide_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(2)
	}
	defer input.Close()

	output := os.Stdout

	err = slide2marpit.Convert(input, output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
