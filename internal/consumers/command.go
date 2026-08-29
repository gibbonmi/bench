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

const usageLine = "usage: bench consumers <qualified-symbol> [--full] | bench consumers --changed [--base <commit> [--source-tip <commit>]] [--full]"

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
	return usageLine + "\n" + promise + "\n" + soundness + "\n" + vendorLine + "\n" + viaLine + "\n" + emptyInterfaceLine + "\n" + genericImplementerLine + "\n" + changedLine + "\n" + sourceTipLine + "\n" + capLine() + "\n" + citationLine + "\n" + refusalLine + "\n"
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
	},
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
	if line, code := checkModes(parsed); line != "" {
		return line + "\n", code
	}
	_, full := parsed.Flags["--full"]
	if _, changed := parsed.Flags["--changed"]; changed {
		return changedCommand(version, args, parsed.Flags["--base"], parsed.Flags["--source-tip"], full)
	}
	symbol := parsed.Positionals[0]

	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	pkgs, err := load(root, "./...")
	if err != nil {
		return refuseLoadForQuery(root, symbol, err) + "\n", 1
	}
	matches, err := Resolve(pkgs, symbol)
	if err != nil {
		return refuseUnresolved(root, symbol, err) + "\n", 1
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
	rows, files, dropped := representableFiles(rows, func(r Row) string { return r.File })
	overCap := !full && len(rows) > rowCap
	var block string
	var err error
	if overCap {
		block, err = toon.TableTyped("consumers_packages", aggregateFields, aggregate(files))
	} else {
		block, err = Render(rows)
	}
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	var actions []axi.Action
	if overCap {
		actions = append(actions, axi.ExecutableInvocation("emit every consumer row",
			axi.KnownArgument("consumers"), axi.KnownArgument(symbol), axi.KnownArgument("--full")))
	}
	return envelope(source, block, pkgCount, countFiles(files), matchCount, len(rows), overCap || dropped, actions)
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
// without a base grades the wrong pair against a defaulted base.
func checkModes(parsed usage.Result) (string, int) {
	_, changed := parsed.Flags["--changed"]
	_, hasBase := parsed.Flags["--base"]
	_, hasTip := parsed.Flags["--source-tip"]
	if len(parsed.Positionals) > 0 && changed {
		return toon.Usage(grammar.Cmd, parsed.Positionals[0]), 2
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
