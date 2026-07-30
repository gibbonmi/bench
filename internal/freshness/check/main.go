package main

import (
	"fmt"
	"os"

	"github.com/gibbonmi/bench/internal/freshness"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: freshness-check <root> <executable>")
		os.Exit(2)
	}
	if err := freshness.Check(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
