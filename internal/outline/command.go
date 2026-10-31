package outline

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
// Help is helpText without its trailing newline, because the caller appends one.
var grammar = usage.Grammar{
	Cmd:  "bench outline",
	Help: strings.TrimSuffix(helpText(), "\n"),
	Flags: []usage.Flag{{Name: "--full"}, {Name: "--production"}, {Name: "--test"},
		{Name: SymbolFlag, HasValue: true, NoEmptyValue: true},
		{Name: responsebound.OffsetFlag, HasValue: true, NoEmptyValue: true},
		{Name: responsebound.SourceFlag, HasValue: true, NoEmptyValue: true}},
	MaxArgs: 1,
}

// promise is the one-line expectation-setting clause the surface repeats: outline
// LOCATES seams, it does not IDENTIFY the project's blessed seams. It lands verbatim
// in the help text below and, as prose a shell surface cannot import, is restated
// with the same verbs in bin/bench.sh's help block.
const promise = "outline locates candidate seams (file:line); it does not identify which are the project's blessed seams — projects/<name>.md owns that."

// Suffix is the argument grammar after the command name. The usage line and the
// `bench help` row both read it, so the two surfaces show one grammar.
const Suffix = " [path] [--full] [--production|--test]"

const usageLine = "usage: bench outline" + Suffix

// SymbolFlag selects one complete Go declaration. SymbolSuffix is the grammar of that form
// after the command name, and its usage line and its `bench help` row both read it.
const (
	SymbolFlag   = "--symbol"
	SymbolSuffix = " [path] " + SymbolFlag + " <name-or-receiver.method> [--production|--test] " + responsebound.CursorUsage
)

const symbolUsageLine = "usage: bench outline" + SymbolSuffix

func helpText() string {
	return usageLine + "\n" + symbolUsageLine + "\n" + promise + "\n"
}

// Command implements `bench outline [path] [--full] [--production|--test]`. It walks the
// tracked files, scoped to an optional path and an optional filter, and dispatches each
// through walkSymbols. It drops any row a control byte would make unrepresentable, and
// renders the symbol table for a path or `--full`, and the per-directory summary for
// the bare form. Each form has the
// definitive empty state when nothing matches, a structured stdout error with exit 1
// outside a repo or on a git failure, and usage on stdout with exit 2 for an unknown
// flag, a second positional argument, or both filter flags. The --symbol form pages one
// complete Go declaration of the same scope; see symbolRead.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, full := parsed.Flags["--full"]
	_, production := parsed.Flags["--production"]
	_, test := parsed.Flags["--test"]
	// The two filters name opposite scopes, so a call with both is a mistyped invocation.
	// It refuses before repository discovery rather than letting one flag win silently.
	if production && test {
		return usageLine + " (--production and --test are mutually exclusive)\n", 2
	}
	selection, symbolic, refusal := parseSymbolCall(parsed.Flags, full)
	if refusal != "" {
		return symbolUsageLine + " (" + refusal + ")\n", 2
	}
	var path string
	havePath := len(parsed.Positionals) == 1
	if havePath {
		path = parsed.Positionals[0]
	}

	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	files, err := listFiles(root, path, havePath)
	if err != nil {
		return toon.Errorf("git ls-files failed", err.Error()) + "\n", 1
	}
	if production || test {
		files = testScope(files, test)
	}
	if symbolic {
		return selection.page(root, files, symbolCommand(parsed.Positionals, parsed.Flags))
	}

	var rows, skips [][]string
	// dirOrder/dirCount form the bare form's ledger: every scanned top-level directory in
	// first-seen order, so one holding no declarations still reports a zero, and the
	// symbols its whole subtree contributed. They index the same rows the symbol table
	// carries, so both forms report one accounting rather than two.
	var dirOrder []string
	dirCount := map[string]int{}
	scanned := 0
	totalSymbols := 0
	for _, rel := range files {
		abs := filepath.Join(root, filepath.FromSlash(rel))
		symbols, reason, ok := walkSymbols(rel, abs)
		if !ok {
			skips = appendSkip(skips, rel, reason)
			continue
		}
		scanned++
		dir := topLevel(rel)
		if toon.Representable(rel) {
			if _, seen := dirCount[dir]; !seen {
				dirOrder = append(dirOrder, dir)
				dirCount[dir] = 0
			}
		}
		for _, s := range symbols {
			totalSymbols++
			if !toon.Representable(rel) || !toon.Representable(s.Name) {
				continue // one poisoned path or name drops only its own row
			}
			rows = append(rows, []string{rel, strconv.Itoa(s.Line), s.Kind, s.Name})
			dirCount[dir]++
		}
	}

	emitted := len(rows)
	omitted := totalSymbols - emitted
	truncated := omitted > 0 || len(skips) > 0
	tbl, err := resultTable(havePath || full, rows, dirOrder, dirCount)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	meta, err := toon.Table("outline_meta", []string{"tracked_files", "scanned_files", "skipped_files", "total_symbols", "emitted_symbols", "omitted_symbols", "truncated"}, [][]string{{strconv.Itoa(len(files)), strconv.Itoa(scanned), strconv.Itoa(len(skips)), strconv.Itoa(totalSymbols), strconv.Itoa(emitted), strconv.Itoa(omitted), strconv.FormatBool(truncated)}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	skipTable, err := toon.Table("outline_skips", []string{"file", "reason"}, skips)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return tbl + meta + skipTable, 0
}

// testScope keeps the files whose Go test status equals test: the Go test files for
// `--test`, and every other file for `--production`. It reuses the helper form's path
// predicate, so the filter and the helper kind agree on what a test file is.
func testScope(files []string, test bool) []string {
	var kept []string
	for _, rel := range files {
		if IsGoTestFile(rel) == test {
			kept = append(kept, rel)
		}
	}
	return kept
}

// topLevel is the bare summary's grouping key: the first segment of git's
// slash-separated path, or "." for a file at the repository root. The summary collapses
// to this depth because a row per scanned directory is larger than the symbol rows it
// replaces once a tree is deep, which defeats the one-screen probe it exists to give.
func topLevel(rel string) string {
	if before, _, ok := strings.Cut(rel, "/"); ok {
		return before
	}
	return "."
}

// resultTable renders the block the invoked form answers with: the symbol rows when a
// path or --full asked for them, the top-level summary otherwise. Both go through the
// shared flat-table emitter, so the block contract has one source.
func resultTable(wantSymbols bool, rows [][]string, dirOrder []string, dirCount map[string]int) (string, error) {
	if wantSymbols {
		return toon.Table("outline", []string{"file", "line", "kind", "name"}, rows)
	}
	summary := make([][]string, 0, len(dirOrder))
	for _, dir := range dirOrder {
		summary = append(summary, []string{dir, strconv.Itoa(dirCount[dir])})
	}
	return toon.Table("outline_dirs", []string{"dir", "symbols"}, summary)
}

func appendSkip(rows [][]string, file, reason string) [][]string {
	if !toon.Representable(file) {
		return rows
	}
	return append(rows, []string{file, reason})
}
