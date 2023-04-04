package main

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/merj/tangled/pkg/tangled"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-h] target source\n", path.Base(os.Args[0]))
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	targetURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	sourceURL, err := url.Parse(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	target, err := tangled.DefaultTarget(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	source, err := tangled.DefaultSource(sourceURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = source.ListenAndServe(target, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}
}
