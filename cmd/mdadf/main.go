package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	adf "github.com/ericmason/mdadf"
)

func main() {
	compact := flag.Bool("c", false, "compact output (no indentation)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [file]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Convert Markdown to Atlassian Document Format (ADF).\n\n")
		fmt.Fprintf(os.Stderr, "If no file is specified, reads from stdin.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var input []byte
	var err error

	if flag.NArg() > 0 {
		input, err = os.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	doc, err := adf.ConvertToDoc(string(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting markdown: %v\n", err)
		os.Exit(1)
	}

	var output []byte
	if *compact {
		output, err = json.Marshal(doc)
	} else {
		output, err = json.MarshalIndent(doc, "", "  ")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
