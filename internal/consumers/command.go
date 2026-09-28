package consumers

import (
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/diff"
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

// Suffix is the symbol form's argument grammar after the command name. The usage line and
// the `bench help` row both read it, so the two surfaces show one grammar.
const Suffix = " <qualified-symbol>... [--production|--test] [--full]"

const usageLine = "usage: bench consumers" + Suffix + " | bench consumers --changed [--base <commit> [--source-tip <commit>]] [--full]"

// capLine tells the reader which block an over-cap default emits and how to get the rows.
// It reads rowCap rather than restating the number, so the help cannot promise one cap
// while the response applies another.
func capLine() string {
	return "over " + strconv.Itoa(rowCap) + " rows the default emits consumers_packages[N]{dir,rows} instead; --full always emits every row."
}

// vendorLine names the enumeration's package scope, so an empty answer for a vendored
// declaration reads as a scope limit rather than as an absence.
const vendorLine = "packages under vendor/ sit outside ./... and are not enumerated."

// viaLine names the three edge classes and the ambiguous-name answer, so an agent reads
// the whole result vocabulary before it runs the command.
const viaLine = "via is call, reference, or implements; a bare name with several matches answers consumers_candidates[N]{qualified,file,line,kind} at exit 0 with one re-query action per row."

// citationLine tells the reader that every success response discloses the run that
// produced it, so a reviewer knows a replay identity is available without one to read.
const citationLine = "every success response ends with citation[1]{sha,state,version,cmd,hash}: the checkout, its clean or dirty state, and a sha256 over the answer above the row."

// emptyInterfaceLine and genericImplementerLine name the two types the implements pass
// leaves out, so a missing row reads as a stated limit.
const (
	emptyInterfaceLine     = "an empty interface emits no implements rows."
	genericImplementerLine = "a generic type is not listed as an implementer."
)

// changedLine names the blast mode's answer and its one frozen-pair rule, so an agent
// reads what --changed enumerates and which revision the rows are positioned in.
const changedLine = "--changed answers blast[N]{changed_symbol,file,line,touched} for every declaration the pair's diff touched, enumerated at the tip; a declaration the pair deleted answers blast_deleted[N]{changed_symbol,base_file,base_line} instead, and touched says the consumer file is itself inside the diff."

// sourceTipLine names the blast mode's revision limit, so a stale tip refuses rather than
// grading a pair the checkout cannot reproduce.
const sourceTipLine = "--changed refuses a --source-tip that is not the checkout's HEAD."

// refusalLine names the three unsound inputs the command refuses, so an agent knows a
// refusal is a stated outcome rather than a crash, and knows it carries no citation.
const refusalLine = "three inputs refuse at exit 1 with no citation row: an ill-typed tree names its first error position, a missing go binary names itself, and a name only a non-Go file declares names that file's language."

func helpText() string {
	return usageLine + "\n" + promise + "\n" + soundness + "\n" + vendorLine + "\n" + viaLine + "\n" + emptyInterfaceLine + "\n" + genericImplementerLine + "\n" + changedLine + "\n" + sourceTipLine + "\n" + capLine() + "\n" + scopeLine + "\n" + severalLine + "\n" + citationLine + "\n" + refusalLine + "\n"
}

// grammar is the declared argument shape usage.Parse enforces for this subcommand. Arity,
// flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:  "bench consumers",
	Help: strings.TrimSuffix(helpText(), "\n"),
	Flags: []usage.Flag{
		{Name: "--full"},
		{Name: "--changed"},
		{Name: "--base", HasValue: true, NoEmptyValue: true},
		{Name: "--source-tip", HasValue: true, NoEmptyValue: true},
		{Name: string(production)},
		{Name: string(testOnly)},
	},
	MaxArgs: -1,
}

// candidateFields is the ambiguous-name schema: one row per declaration the bare name
// reached. qualified is the exact re-query argument, so the agent retypes nothing.
var candidateFields = []string{"qualified", "file", "line", "kind"}

// aggregateFields is the over-cap schema: one row per consumer directory, which is one Go
// package, with the row count that directory contributed. A one-symbol query holds the
// symbol constant across the block, so it gets no column; a several-symbol query prepends one.
var aggregateFields = []string{"dir", "rows"}

// metaFields is the response accounting every form carries.
var metaFields = []string{"packages", "files", "matches", "rows", "truncated"}

// CommandWithVersion implements `bench consumers <qualified-symbol>... [--production |
// --test] [--full]` for one bench version. The version is a cell of every success
// response's citation row, and it lives in package main, so the registration injects it
// here rather than the package reading a second copy. The command resolves each symbol
// over the repository's packages and emits the result blocks, the meta accounting, the
// citation row, and the terminal help envelope. A symbol result is a terminal read, so
// its envelope is empty unless the default truncated or a name was ambiguous.
func CommandWithVersion(version string) func([]string) (string, int) {
	return func(args []string) (string, int) { return command(version, args) }
}

func command(version string, args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	if line, code := checkModes(parsed); line != "" {
		return line + "\n", code
	}
	_, full := parsed.Flags["--full"]
	if _, changed := parsed.Flags["--changed"]; changed {
		return changedCommand(version, args, parsed.Flags["--base"], parsed.Flags["--source-tip"], full)
	}
	symbols := parsed.Positionals

	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	pkgs, err := load(root, "./...")
	if err != nil {
		// The non-Go sweep names the first operand; the load refusal is the same for all.
		return refuseLoadForQuery(root, symbols[0], err) + "\n", 1
	}
	q := query{scope: scopeOf(parsed.Flags), full: full, several: len(symbols) > 1}
	answers := make([]answer, 0, len(symbols))
	for _, symbol := range symbols {
		matches, err := Resolve(pkgs, symbol)
		if err != nil {
			return refuseUnresolved(root, symbol, err) + "\n", 1
		}
		answers = append(answers, q.answer(pkgs, root, symbol, matches))
	}
	return q.response(citation{root: root, version: version, args: args}, len(pkgs), answers)
}

// representableFiles projects each row's file cell and drops the rows TOON cannot carry.
// Every consumers response form routes through it, so a control byte in a git-sourced path
// drops only its own row and reports truncated, the way `bench outline` drops a poisoned
// row. The third result says a row was dropped, which is what the meta flag reports.
func representableFiles[T any](rows []T, file func(T) string) ([]T, []string, bool) {
	kept := make([]T, 0, len(rows))
	files := make([]string, 0, len(rows))
	for _, row := range rows {
		name := file(row)
		if !toon.Representable(name) {
			continue
		}
		kept = append(kept, row)
		files = append(files, name)
	}
	return kept, files, len(kept) != len(rows)
}

// envelope closes every response shape: the result block, the meta accounting, the
// citation row, then the terminal help block. The symbol answer and the blast answer
// render through it, so neither can grow a second accounting, lose its citation, or drop
// the terminal envelope.
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

// aggregate collapses the rows' file cells to one row per consumer directory, in directory
// order. It is the one derivation of the over-cap block, and both the symbol form and the
// blast form pass their own rows' files to it.
func aggregate(files []string) [][]any {
	counts := map[string]int{}
	for _, f := range files {
		counts[path.Dir(f)]++
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
func countFiles(files []string) int {
	seen := map[string]bool{}
	for _, f := range files {
		seen[f] = true
	}
	return len(seen)
}

// checkModes separates the two invocation shapes usage.Parse cannot: the symbol query and
// the blast query. The revision rules are the rules `bench test --changed` applies, so one
// revision grammar serves both surfaces — a positional with --changed names two subjects,
// a revision flag without --changed grades a pair nothing selected, and a source tip
// without a base grades the wrong pair against a defaulted base. The two row filters name
// disjoint row sets, so both together select nothing, and neither filters the blast rows.
func checkModes(parsed usage.Result) (string, int) {
	_, changed := parsed.Flags["--changed"]
	_, hasBase := parsed.Flags["--base"]
	_, hasTip := parsed.Flags["--source-tip"]
	if len(parsed.Positionals) > 0 && changed {
		return toon.Usage(grammar.Cmd, parsed.Positionals[0]), 2
	}
	_, hasProduction := parsed.Flags[string(production)]
	if _, hasTest := parsed.Flags[string(testOnly)]; hasTest && hasProduction {
		return toon.Usage(grammar.Cmd, string(testOnly)), 2
	}
	if filter := scopeOf(parsed.Flags); filter != everyRow && changed {
		return toon.Usage(grammar.Cmd, string(filter)), 2
	}
	if (hasBase || hasTip) && !changed {
		flag := "--base"
		if hasTip {
			flag = "--source-tip"
		}
		return toon.Usage(grammar.Cmd, flag), 2
	}
	if hasTip && !hasBase {
		return toon.Usage(grammar.Cmd, "--source-tip"), 2
	}
	if !changed && len(parsed.Positionals) == 0 {
		return toon.MissingArg(grammar.Cmd, "argument"), 2
	}
	return "", 0
}

// changedCommand answers the blast query over one frozen pair. The pair, the changed
// paths, the hunks, and the base-side sources are all read here, at the rim; everything
// below is a pure function of those bytes and the tip packages.
//
// Enumeration runs at the tip, and the loader reads the checkout in place, so a tip that
// is not this checkout's HEAD refuses rather than answering from the wrong tree. Building
// a temporary checkout for a historical tip is a separate priced path.
func changedCommand(version string, args []string, base, sourceTip string, full bool) (string, int) {
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	source := citation{root: root, version: version, args: args}
	// Enumeration reads the checkout in place, so a dirty checkout would position rows in
	// bytes the pair does not contain. The rows must come only from the frozen pair, so a
	// dirty checkout refuses rather than answering from a tree nothing froze.
	if source.state() != "clean" {
		return toon.Errorf("checkout is dirty; blast rows come only from the frozen pair",
			"commit or clean the checkout, then rerun the exact invocation") + "\n", 1
	}
	subject, kind, hint := diff.ResolveChangedSubject(root, base, sourceTip)
	if kind != "" {
		return toon.Errorf("changed selection failed", kind+": "+hint) + "\n", 1
	}
	head, err := git.Output("-C", root, "rev-parse", "HEAD")
	if err != nil {
		return toon.Errorf("cannot resolve HEAD", "run the command inside a git checkout with a resolvable HEAD") + "\n", 1
	}
	if subject.Tip != head {
		return toon.Errorf("tip "+subject.Tip+" is not this checkout's HEAD "+head,
			"blast enumerates at the tip in this checkout; check out the tip commit, then retry") + "\n", 1
	}
	// A pair that changed no Go file has the definitive empty answer, and no package the
	// loader could name would change it. The answer therefore precedes the load: a tree the
	// loader cannot load still gets its empty table rather than a refusal.
	if len(goPaths(subject.Paths)) == 0 {
		return blastResponse(source, 0, 0, nil, nil, full, args)
	}
	hunks, err := readHunks(root, subject.Base, subject.Tip, subject.Paths)
	if err != nil {
		return toon.Errorf("diff read failed: "+err.Error(), "the pair must be readable in this repository; retry the exact invocation") + "\n", 1
	}
	pkgs, err := load(root, "./...")
	if err != nil {
		return refuseLoad(err) + "\n", 1
	}
	added := map[string][]lineSpan{}
	for _, fh := range hunks {
		if fh.TipPath != "" {
			added[fh.TipPath] = append(added[fh.TipPath], fh.Added...)
		}
	}
	changed := map[string]bool{}
	for _, p := range subject.Paths {
		changed[p] = true
	}
	decls := touchedDecls(pkgs, root, added)
	rows := blastRows(pkgs, root, decls, changed)
	deleted := deletedRows(pkgs, root, hunks, readBaseSources(root, subject.Base, hunks))
	return blastResponse(source, len(pkgs), len(decls), rows, deleted, full, args)
}
