package main

import (
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
)

var testHelpSuffix = strings.TrimPrefix(testreport.Usage, "bench test")

func testCommand(args []string) (string, int) {
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	return testreport.Command(root, args)
}
