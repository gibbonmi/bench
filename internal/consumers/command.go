package consumers

import (
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// promise and soundness are the two clauses the help text states verbatim. They are
// constants so the help and every other reader share one source. promise keeps the
// identification boundary: the tool names resolved edges, and the project profile keeps
// seam blessing. soundness names the whole static-analysis limit, so an agent never reads
// an empty table as proof that no edge exists.
const (
	promise   = "consumers identifies resolved reference edges; it does not bless a seam — projects/<name>.md owns that."
	soundness = "sound for static Go references only: reflection, go:linkname, plugin, and exec edges are invisible; the default build context is the graded one."
)

// rowCap is the default response's row budget. Past it the default answers with the
// per-directory aggregate instead, so one hot symbol never floods an agent's context.
const rowCap = 200

const usageLine = "usage: bench consumers <qualified-symbol> [--full]"

// capLine tells the reader which block an over-cap default emits and how to get the rows.
// It reads rowCap rather than restating the number, so the help cannot promise one cap
// while the response applies another.
func capLine() string {
	return "over " + strconv.Itoa(rowCap) + " rows the default emits consumers_packages[N]{dir,rows} instead; --full always emits every row."
}

// viaLine names the three edge classes and the ambiguous-name answer, so an agent reads
// the whole result vocabulary before it runs the command.
const viaLine = "via is call, reference, or implements; a bare name with several matches answers consumers_candidates[N]{qualified,file,line,kind} at exit 0 with one re-query action per row."

// citationLine tells the reader that every success response discloses the run that
// produced it, so a reviewer knows a replay identity is available without one to read.
const citationLine = "every success response ends with citation{sha,state,version,cmd,hash}: the checkout, its clean or dirty state, and a sha256 over the answer above the row."

func helpText() string {
	return usageLine + "\n" + promise + "\n" + soundness + "\n" + viaLine + "\n" + capLine() + "\n" + citationLine + "\n"
}

// grammar is the declared argument shape usage.Parse enforces for this subcommand. Arity,
// flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:     "bench consumers",
	Help:    strings.TrimSuffix(helpText(), "\n"),
	Flags:   []usage.Flag{{Name: "--full"}},
	MinArgs: 1,
	MaxArgs: 1,
}

// candidateFields is the ambiguous-name schema: one row per declaration the bare name
// reached. qualified is the exact re-query argument, so the agent retypes nothing.
var candidateFields = []string{"qualified", "file", "line", "kind"}

// aggregateFields is the over-cap schema: one row per consumer directory, which is one Go
// package, with the row count that directory contributed. The queried symbol is constant
// across the block, so it gets no column.
var aggregateFields = []string{"dir", "rows"}

// metaFields is the response accounting every form carries.
var metaFields = []string{"packages", "files", "matches", "rows", "truncated"}

// CommandWithVersion implements `bench consumers <qualified-symbol> [--full]` for one
// bench version. The version is a cell of every success response's citation row, and it
// lives in package main, so the registration injects it here rather than the package
// reading a second copy. The command resolves the symbol over the repository's packages
// and emits the consumers table, the meta accounting, the citation row, and the terminal
// help envelope. A symbol result is a terminal read, so its envelope is empty unless the
// default truncated.
func CommandWithVersion(version string) func([]string) (string, int) {
	return func(args []string) (string, int) { return command(version, args) }
}

func command(version string, args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, full := parsed.Flags["--full"]
	symbol := parsed.Positionals[0]

	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	pkgs, err := load(root, "./...")
	if err != nil {
		return toon.Errorf("package load failed: "+err.Error(), "fix the tree so it type-checks, then retry") + "\n", 1
	}
	matches, err := Resolve(pkgs, symbol)
	if err != nil {
		return toon.Errorf(err.Error(), "pass a qualified symbol such as outline.Command") + "\n", 1
	}
	source := citation{root: root, version: version, args: args}
	if len(matches) > 1 {
		return candidates(source, pkgs, matches)
	}
	return response(source, symbol, len(pkgs), len(matches), Rows(pkgs, matches[0].Obj, root), full)
}

// response renders the whole answer for one resolved symbol: the result block, the meta
// accounting, and the help envelope. Row rendering itself belongs to the core, so this
// function composes Render rather than restating the row schema.
func response(source citation, symbol string, pkgCount, matchCount int, rows []Row, full bool) (string, int) {
	truncated := !full && len(rows) > rowCap
	var block string
	var err error
	if truncated {
		block, err = toon.TableTyped("consumers_packages", aggregateFields, aggregate(rows))
	} else {
		block, err = Render(rows)
	}
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	var actions []axi.Action
	if truncated {
		actions = append(actions, axi.ExecutableInvocation("emit every consumer row",
			axi.KnownArgument("consumers"), axi.KnownArgument(symbol), axi.KnownArgument("--full")))
	}
	return envelope(source, block, pkgCount, countFiles(rows), matchCount, len(rows), truncated, actions)
}

// envelope closes every response shape: the result block, the meta accounting, the
// citation row, then the terminal help block. Both the symbol answer and the candidates
// answer render through it, so neither can grow a second accounting, lose its citation,
// or drop the terminal envelope.
func envelope(source citation, block string, pkgCount, fileCount, matchCount, rowCount int, truncated bool, actions []axi.Action) (string, int) {
	// The counts and the flag are typed cells, so an integer stays bare and the flag stays
	// a boolean through a TOON round-trip.
	meta, err := toon.TableTyped("meta", metaFields, [][]any{{
		pkgCount, fileCount, matchCount, rowCount, truncated,
	}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	// The citation hashes every byte before it, so it renders after the answer is complete
	// and before the help block the AXI contract pins as terminal.
	cited, err := source.row(block + meta)
	if err != nil {
		return toon.Errorf("citation failed: "+err.Error(), "run the command inside a git checkout with a resolvable HEAD") + "\n", 1
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return block + meta + cited + help, 0
}

// candidates answers an ambiguous bare name. It is an answer, not a refusal: the table
// names every declaration the name reached, the meta accounting states zero consumer
// rows, and the envelope carries one literal re-query per row in table order. Every
// argument is known, so no row offers a slot the agent must fill.
func candidates(source citation, pkgs []*Package, matches []Match) (string, int) {
	rows := make([][]any, 0, len(matches))
	actions := make([]axi.Action, 0, len(matches))
	for _, m := range candidateOrder(pkgs, matches, source.root) {
		rows = append(rows, []any{m.qualified, m.file, m.line, m.kind})
		actions = append(actions, axi.ExecutableInvocation("re-query the qualified symbol",
			axi.KnownArgument("consumers"), axi.KnownArgument(m.qualified)))
	}
	block, err := toon.TableTyped("consumers_candidates", candidateFields, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return envelope(source, block, len(pkgs), 0, len(matches), 0, false, actions)
}

// candidateRow is one match rendered for the candidates table, in the sort order the
// table prints: by file, then by line.
type candidateRow struct {
	qualified string
	file      string
	line      int
	kind      string
}

// candidateOrder positions every match at its declaration. The position comes from the
// declaring package's own file set, because a match names the package it was found in.
func candidateOrder(pkgs []*Package, matches []Match, root string) []candidateRow {
	fsets := map[string]*token.FileSet{}
	for _, pkg := range pkgs {
		fsets[pkg.PkgPath] = pkg.Fset
	}
	out := make([]candidateRow, 0, len(matches))
	for _, m := range matches {
		row := candidateRow{qualified: m.Qualified, kind: m.Kind}
		if fset := fsets[m.PkgPath]; fset != nil {
			pos := fset.Position(m.Obj.Pos())
			row.file, row.line = relPath(root, pos.Filename), pos.Line
		}
		out = append(out, row)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].file != out[j].file {
			return out[i].file < out[j].file
		}
		return out[i].line < out[j].line
	})
	return out
}

// aggregate collapses rows to one row per consumer directory, in directory order. Rows
// arrive sorted by file, so the collapse is a single pass over a stable order.
func aggregate(rows []Row) [][]any {
	counts := map[string]int{}
	for _, r := range rows {
		counts[path.Dir(r.File)]++
	}
	dirs := make([]string, 0, len(counts))
	for dir := range counts {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	out := make([][]any, 0, len(dirs))
	for _, dir := range dirs {
		out = append(out, []any{dir, counts[dir]})
	}
	return out
}

// countFiles is the meta table's files cell: the number of distinct files the rows name.
func countFiles(rows []Row) int {
	seen := map[string]bool{}
	for _, r := range rows {
		seen[r.File] = true
	}
	return len(seen)
}
