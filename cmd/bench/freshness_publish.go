package main

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/freshness"
)

func freshnessCheck(args []string, invoked string, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: bench freshness-check <root>")
		return 2
	}
	root, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	executable, err := filepath.Abs(invoked)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := freshness.Verify(root, executable); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func freshnessPublish(args []string, invoked string, stderr io.Writer) int {
	if len(args) != 4 {
		fmt.Fprintln(stderr, "usage: bench freshness-publish <root> <output-path> <manifest-dir> <version>")
		return 2
	}
	root, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	executable, err := filepath.Abs(invoked)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := freshness.Publish(root, executable, args[1], args[2], args[3]); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
