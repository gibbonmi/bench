// Package probe mutates one file once, runs one focused test or named check, restores
// the file, and reports whether the mutation bit. The verb owns the sequence a
// coordinator otherwise spells as five raw shell calls, so a done-claim's proof is one
// bounded command whose restore is verified rather than assumed.
package probe

import (
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar declares the one argument shape the verb parses through usage.Parse. Exactly
// one mutation form and exactly one selection form are required, because a probe that
// defaulted either would mutate or run something the caller did not name.
var grammar = usage.Grammar{
	Cmd:  "bench probe",
	Help: "usage: bench probe <file> (--swap <old> --with <new> | --omit <old>) (--package <expr> [--run <go-regex>] | --check <name>) [--full]",
	Flags: []usage.Flag{
		{Name: "--swap", HasValue: true, NoEmptyValue: true},
		{Name: "--with", HasValue: true, NoEmptyValue: true},
		{Name: "--omit", HasValue: true, NoEmptyValue: true},
		{Name: "--package", HasValue: true, NoEmptyValue: true},
		{Name: "--run", HasValue: true, NoEmptyValue: true},
		{Name: "--check", HasValue: true, NoEmptyValue: true},
		{Name: "--full"},
	},
	MinArgs: 1,
	MaxArgs: 1,
}

// Command parses one probe invocation, resolves the repository root, and runs the probe.
// The refusal order is fixed: usage, not in a repository, the subject, the mutation
// count, the selection, the unsupported check, then the gate lock. Every one of those
// runs before the preserved copy is written and before any run child starts.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(probeGrammar(), args)
	if line != "" {
		return line + "\n", code
	}
	if line, code := gradeForms(parsed); line != "" {
		return line, code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	return run(root, parsed)
}

// probeGrammar answers the grammar a help request reads: the usage line and the focused-run
// owner's own notes. The notes come from the owner rather than from a copy here, so the
// help states what the run does instead of a second description of it. A form refusal keeps
// the bare usage line, because a refused argv needs the shape it broke, not the notes.
func probeGrammar() usage.Grammar {
	withNotes := grammar
	withNotes.Help = grammar.Help + "\n" + testreport.ProbeNotes()
	return withNotes
}

// gradeForms refuses every argv that names no mutation form, no selection form, or two
// of either. The shared parser grades one flag at a time, so the mutual exclusions
// between flags are graded here, and each answers with the grammar's own usage line.
func gradeForms(parsed usage.Result) (string, int) {
	_, swap := parsed.Flags["--swap"]
	_, with := parsed.Flags["--with"]
	_, omit := parsed.Flags["--omit"]
	_, pkg := parsed.Flags["--package"]
	_, run := parsed.Flags["--run"]
	_, check := parsed.Flags["--check"]
	switch {
	case swap == omit:
		return usageLine()
	case swap && !with, omit && with:
		return usageLine()
	case pkg == check:
		return usageLine()
	case run && !pkg:
		return usageLine()
	}
	return "", 0
}

func usageLine() (string, int) { return grammar.Help + "\n", 2 }
