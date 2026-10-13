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

// Usage is the one owner of the bench probe grammar line.
const Usage = "bench probe <file> (--swap <old> (--with <new> | --with-file <path>) [--occurrence <n>] | --omit <old> [--occurrence <n>] | --unwrap <call> [--occurrence <n>] | --insert-before <anchor> (--with <text> | --with-file <path>) [--occurrence <n>] | --omit-file <path> | --delete | (--swap <old> --with <new> | --omit <old>)...) (--package <expr> [--package <expr>]... [--run <go-regex>] | --check <name> | --check system --run <go-regex>) [--expect <test>] [--full]"

// grammar declares the one argument shape the verb parses through usage.Parse. Exactly
// one mutation group and exactly one selection form are required, because a probe that
// defaulted either would mutate or run something the caller did not name. --swap, --with,
// and --omit repeat, because their repetitions make one combined mutant. --package repeats,
// because each run passes every expression to one Go start.
var grammar = usage.Grammar{
	Cmd:  usage.CommandName(Usage),
	Help: "usage: " + Usage,
	Flags: []usage.Flag{
		{Name: "--swap", HasValue: true, NoEmptyValue: true, Repeatable: true},
		{Name: "--with", HasValue: true, NoEmptyValue: true, Repeatable: true},
		{Name: "--with-file", HasValue: true, NoEmptyValue: true},
		{Name: "--omit", HasValue: true, NoEmptyValue: true, Repeatable: true},
		{Name: "--unwrap", HasValue: true, NoEmptyValue: true},
		{Name: "--omit-file", HasValue: true, NoEmptyValue: true},
		{Name: "--insert-before", HasValue: true, NoEmptyValue: true},
		{Name: "--occurrence", HasValue: true, NoEmptyValue: true},
		{Name: "--delete"},
		{Name: "--package", HasValue: true, NoEmptyValue: true, Repeatable: true},
		{Name: "--run", HasValue: true, NoEmptyValue: true},
		{Name: "--check", HasValue: true, NoEmptyValue: true},
		{Name: "--expect", HasValue: true, NoEmptyValue: true},
		{Name: "--full"},
	},
	MinArgs: 1,
	MaxArgs: 1,
}

// Command parses one probe invocation, resolves the repository root, and runs the probe.
// The refusal order is fixed: usage, not in a repository, the subject, the mutation input,
// the mutation match, the selection, a subject the selected grade never reads, the
// substring miss, the gate lock, the subject lock, the subject re-read, then the baseline.
// Every one of those runs before the preserved copy is written, and only the baseline
// starts a run child.
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

// mutationGroup is one mutation group of the closed form matrix.
type mutationGroup int

const (
	// editGroup is --swap and --omit: one edit, or several that make one combined mutant.
	editGroup mutationGroup = iota
	unwrapGroup
	omitFileGroup
	insertGroup
	deleteGroup
)

// groupFlags names the flags that select each mutation group.
var groupFlags = [...][]string{
	editGroup:     {"--swap", "--omit"},
	unwrapGroup:   {"--unwrap"},
	omitFileGroup: {"--omit-file"},
	insertGroup:   {"--insert-before"},
	deleteGroup:   {"--delete"},
}

// namedGroup answers the mutation group that the argv names, and ok is false when the argv
// names no group or several. gradeForms and mutationInput both read the group here, so the
// form that the grader accepts is the form that the input builds.
func namedGroup(parsed usage.Result) (mutationGroup, bool) {
	group, count := editGroup, 0
	for g, flags := range groupFlags {
		for _, flag := range flags {
			if _, ok := parsed.Flags[flag]; ok {
				group, count = mutationGroup(g), count+1
				break
			}
		}
	}
	return group, count == 1
}

// gradeForms grades the closed form matrix: an argv names exactly one mutation group and
// exactly one selection form, and each form takes only its partner flags. The shared parser
// grades one flag at a time, so the rules between flags are graded here, and each answers
// with the grammar's own usage line.
func gradeForms(parsed usage.Result) (string, int) {
	has := func(flag string) bool { _, ok := parsed.Flags[flag]; return ok }
	group, oneGroup := namedGroup(parsed)
	swaps, withs, omits := len(parsed.Repeated["--swap"]), len(parsed.Repeated["--with"]), len(parsed.Repeated["--omit"])
	withFile, insert := has("--with-file"), group == insertGroup
	combined := swaps+omits > 1
	texts := withs
	if withFile {
		texts++
	}
	_, validOccurrence := occurrence(parsed)
	_, pkg := parsed.Flags["--package"]
	_, run := parsed.Flags["--run"]
	_, check := parsed.Flags["--check"]
	switch {
	case !oneGroup:
		return usageLine()
	// --with and --with-file are two spellings of one text.
	case withs > 0 && withFile:
		return usageLine()
	// An insertion takes exactly one text.
	case insert && texts != 1:
		return usageLine()
	// Outside an insertion, each text pairs with one --swap, so a text without a swap and a
	// swap without a text both refuse, --delete and the omission forms included.
	case !insert && texts != swaps:
		return usageLine()
	// A combined mutant takes scalar texts only.
	case combined && withFile:
		return usageLine()
	case has("--occurrence") && (!validOccurrence || group == omitFileGroup || group == deleteGroup || combined):
		return usageLine()
	case pkg == check:
		return usageLine()
	// A run pattern selects within a package expression or a check that takes one.
	case run && !pkg && !testreport.CheckTakesRun(parsed.Flags["--check"]):
		return usageLine()
	// A prose run names no test, so an expected witness there could never match.
	case has("--expect") && parsed.Flags["--check"] == testreport.ProseCheckName:
		return usageLine()
	}
	return "", 0
}

func usageLine() (string, int) { return grammar.Help + "\n", 2 }
