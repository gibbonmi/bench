package outline

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
// Help is helpText without its trailing newline, because the caller appends one.
var grammar = usage.Grammar{
	Cmd:  commandName,
	Help: strings.TrimSuffix(helpText(), "\n"),
	Flags: []usage.Flag{{Name: "--full"}, {Name: ProductionFlag}, {Name: TestFlag},
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

// commandName is the one spelling of the command name. The grammar, both usage lines, and
// every derived continuation read it.
const commandName = "bench outline"

// ProductionFlag and TestFlag are the two test-scope filters. filterUsage is their grammar
// group, and filterFlags lists them in the order a continuation repeats them.
const (
	ProductionFlag = "--production"
	TestFlag       = "--test"
	filterUsage    = "[" + ProductionFlag + "|" + TestFlag + "]"
)

var filterFlags = []string{ProductionFlag, TestFlag}

// Suffix is the argument grammar after the command name. The usage line and the
// `bench help` row both read it, so the two surfaces show one grammar.
const Suffix = " [path] [--full] " + filterUsage

const usageLine = "usage: " + commandName + Suffix

// SymbolFlag selects one complete Go declaration. SymbolSuffix is the grammar of that form
// after the command name, and its usage line and its `bench help` row both read it.
const (
	SymbolFlag   = "--symbol"
	SymbolSuffix = " [path] " + SymbolFlag + " <name-or-receiver.method> " + filterUsage + " " + responsebound.CursorUsage
)

const symbolUsageLine = "usage: " + commandName + SymbolSuffix

func helpText() string {
	return usageLine + "\n" + symbolUsageLine + "\n" + promise + "\n"
}

// usageRefusal renders a usage refusal: the usage line of the refused form with the
// reason, at exit 2.
func usageRefusal(line, reason string) (string, int) {
	return line + " (" + reason + ")\n", 2
}

// call is one outline call that both grammars accept: the parsed arguments, the --full
// flag, the two filters, and the --symbol selection when the call selects a declaration.
type call struct {
	parsed                           usage.Result
	full, production, test, symbolic bool
	selection                        symbolRead
}

// parseCall reads one outline call, or answers its help or usage refusal line with the exit
// code. Command and SymbolHead read every call through it, so one parse decides both the
// response and the census head.
func parseCall(args []string) (call, string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return call{}, line + "\n", code
	}
	accepted := call{parsed: parsed}
	_, accepted.full = parsed.Flags["--full"]
	_, accepted.production = parsed.Flags[ProductionFlag]
	_, accepted.test = parsed.Flags[TestFlag]
	// The two filters name opposite scopes, so a call with both is a mistyped invocation.
	// It refuses before repository discovery rather than letting one flag win silently.
	if accepted.production && accepted.test {
		line, code := usageRefusal(usageLine, ProductionFlag+" and "+TestFlag+" are mutually exclusive")
		return call{}, line, code
	}
	selection, symbolic, refusal := parseSymbolCall(parsed.Flags, accepted.full)
	if refusal != "" {
		line, code := usageRefusal(symbolUsageLine, refusal)
		return call{}, line, code
	}
	accepted.selection, accepted.symbolic = selection, symbolic
	return accepted, "", 0
}

// SymbolHead answers the census output head of an outline call that selects one declaration:
// the command name with SymbolFlag. A call of another form, and a call that the grammar
// refuses, answers false and keeps the command's root head.
func SymbolHead(args []string) (string, bool) {
	accepted, line, _ := parseCall(args)
	return commandName + " " + SymbolFlag, line == "" && accepted.symbolic
}

// Command implements both outline grammars. `bench outline [path] [--full]
// [--production|--test]` walks the tracked files, scoped to an optional path and an
// optional filter, and dispatches each through walkSymbols. It drops any row a control byte
// would make unrepresentable, and renders the symbol table for a path or `--full`, and the
// per-directory summary for the bare form. Each of these forms ends with derived
// continuation help; see legacyHelp. `bench outline [path] --symbol <name> [--production|--test]
// [--offset <n> --source <id>]` pages one complete Go declaration of the same scope; see
// symbolRead. Each form has the definitive empty state when nothing matches, and a
// structured stdout error with exit 1 outside a repo or on a git failure. Usage goes to
// stdout with exit 2 for an unknown flag, a second positional argument, or both filter
// flags. The --symbol grammar also refuses with exit 2 for a cursor flag without --symbol,
// --symbol with --full, a malformed cursor, or a selector that is not a Go name.
func Command(args []string) (string, int) {
	accepted, line, code := parseCall(args)
	if line != "" {
		return line, code
	}
	parsed, full, production, test := accepted.parsed, accepted.full, accepted.production, accepted.test
	selection, symbolic := accepted.selection, accepted.symbolic
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
	help, err := axi.RenderHelp(legacyHelp(root, havePath || full, parsed, rows, dirOrder, dirCount))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return tbl + meta + skipTable + help, 0
}

// legacyHelp derives the continuations of a symbol table or a directory summary. A symbol
// table that lists a Go declaration offers one --symbol read of a listed name in the same
// scope. A summary offers one narrowing call per listed top-level directory that holds
// symbols, named from the caller's directory. Root-level files have no directory to narrow
// into. Any other result is terminal and offers nothing.
func legacyHelp(root string, symbolTable bool, parsed usage.Result, rows [][]string, dirOrder []string, dirCount map[string]int) []axi.Action {
	if symbolTable {
		for _, row := range rows {
			if inSymbolScope(row[0]) {
				arguments := symbolCall(parsed.Positionals, parsed.Flags, axi.FutureInput("name"), axi.KnownArgument)
				return []axi.Action{axi.ExecutableInvocation("read one listed Go declaration", arguments...)}
			}
		}
		return nil
	}
	filters := filterArgs(parsed.Flags)
	var actions []axi.Action
	for _, dir := range dirOrder {
		if dir == "." || dirCount[dir] == 0 {
			continue
		}
		argv := commandArgv(append([]string{fromCaller(root, dir)}, filters...)...)
		actions = append(actions, axi.ExecutableInvocation("list the symbols under this directory", knownArguments(argv)...))
	}
	return actions
}

// fromCaller names the root-relative path rel from the caller's directory, so a derived
// command runs where the caller ran this one.
func fromCaller(root, rel string) string {
	cwd, _ := os.Getwd()
	operand, err := filepath.Rel(cwd, filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return rel
	}
	return operand
}

// knownArguments declares each value of argv as a known continuation argument.
func knownArguments(argv []string) []axi.InvocationArgument {
	arguments := make([]axi.InvocationArgument, 0, len(argv))
	for _, argument := range argv {
		arguments = append(arguments, axi.KnownArgument(argument))
	}
	return arguments
}

// filterArgs answers the filter flags of the caller's call, so each continuation keeps the
// caller's test scope.
func filterArgs(flags map[string]string) []string {
	var present []string
	for _, filter := range filterFlags {
		if _, ok := flags[filter]; ok {
			present = append(present, filter)
		}
	}
	return present
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
