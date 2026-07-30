package main

import (
	"fmt"
	"os"

	"github.com/gibbonmi/bench/internal/freshness"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: freshness-publish <root> <staged-executable> <output-path>")
		os.Exit(2)
	}
	if err := freshness.Publish(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
