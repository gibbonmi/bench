package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// commandsGrammar is the declared argument shape usage.Parse enforces for `bench
// commands`. Arity, flag recognition, `--`, and help all come from there rather than a
// local switch. `--brief` is the only form that lists anything, so an invocation without
// it keeps its exit-2 usage answer.
var commandsGrammar = usage.Grammar{
	Cmd:   "bench commands",
	Help:  "usage: bench commands --brief",
	Flags: []usage.Flag{{Name: "--brief"}},
}

func commandsCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(commandsGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	if _, brief := parsed.Flags["--brief"]; !brief {
		return commandsGrammar.Help + "\n", 2
	}
	if root, err := git.Root(); err == nil && commandsProbeIsStale(root) {
		return toon.Errorf("bench binary does not match this checkout's sources", freshness.RebuildAction(root)) + "\n", 1
	}
	return "version\ncommands --brief\nstatus\n", 0
}

// commandsProbeIsStale grades the executable answering the probe against root's current
// build inputs. A session runs this verb to learn whether Bench answers at all, so an
// answer from last week's bytes is the one false pass no other check in the loop reaches.
// The verdict comes from freshness, so the probe derives no digest and reads no mtime.
//
// Only an executable inside root is graded. root's sources describe the binary root
// published, and they describe no other: a release artifact, an installed platform
// binary, and a private run selection all live outside the checkout and are built from
// their own trees. Grading one of those against the cwd's repository would refuse a
// sound binary. Outside a repository, and in a repository that declares no build inputs,
// there is likewise nothing to grade, and the probe keeps its answer.
func commandsProbeIsStale(root string) bool {
	if !freshness.DeclaresBuildInputs(root) {
		return false
	}
	executable, err := os.Executable()
	if err != nil {
		return false
	}
	// A shim on PATH points at the published bytes, so resolve both sides before the
	// containment test compares them.
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return false
	}
	checkout, err := filepath.EvalSymlinks(root)
	if err != nil {
		return false
	}
	relative, err := filepath.Rel(checkout, executable)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return false
	}
	return freshness.Verify(root, executable) != nil
}
