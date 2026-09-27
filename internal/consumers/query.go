package consumers

import (
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/outline"
	"github.com/gibbonmi/bench/internal/toon"
)

// scope is the production-or-test row filter. A filter's value is the flag that selects
// it, so the grammar, the refusal, and every replay action spell the flag from one source.
type scope string

const (
	everyRow   scope = ""
	production scope = "--production"
	testOnly   scope = "--test"
)

// scopeLine and severalLine state the filter and the several-symbol answer in the help.
// scopeLine reads the flag constants and outline's test-file suffix, so the help cannot
// advertise a filter the grammar and the predicate do not enforce.
const (
	scopeLine = string(production) + " keeps only rows outside " + outline.GoTestFileSuffix + " files and " + string(testOnly) +
		" keeps only rows inside them, before the row cap; the two together, or either with --changed, refuse at exit 2."
	severalLine = "several symbols answer in one response: every block gains a leading symbol column, and the row cap and the candidates form apply to each symbol."
)

// scopeOf reads the filter a parse selected. checkModes refuses both flags together
// before any caller acts on the answer.
func scopeOf(flags map[string]string) scope {
	for _, s := range []scope{production, testOnly} {
		if _, ok := flags[string(s)]; ok {
			return s
		}
	}
	return everyRow
}

// keep drops the rows the filter excludes. An implements row reads the file of the
// implementing declaration.
func (s scope) keep(rows []Row) []Row {
	if s == everyRow {
		return rows
	}
	kept := make([]Row, 0, len(rows))
	for _, r := range rows {
		if outline.IsGoTestFile(r.File) == (s == testOnly) {
			kept = append(kept, r)
		}
	}
	return kept
}

// query is one invocation's answer policy: the filter, the cap switch, and whether the
// blocks carry the symbol column.
type query struct {
	scope   scope
	full    bool
	several bool
}

// answer is one operand's share of the response. An ambiguous operand holds candidates
// and no rows; a resolved operand holds its filtered, representable rows.
type answer struct {
	symbol     string
	matches    int
	candidates []candidateRow
	rows       []Row
	files      []string
	dropped    bool
	overCap    bool
}

// answer resolves one operand's rows, or its candidates when the name was ambiguous. The
// filter applies before the cap, so the cap counts the rows the agent asked for.
func (q query) answer(pkgs []*Package, root, symbol string, matches []Match) answer {
	a := answer{symbol: symbol, matches: len(matches)}
	if len(matches) > 1 {
		a.candidates = candidateOrder(pkgs, matches, root)
		return a
	}
	a.rows, a.files, a.dropped = representableFiles(q.scope.keep(Rows(pkgs, matches[0].Obj, root)),
		func(r Row) string { return r.File })
	a.overCap = !q.full && len(a.rows) > rowCap
	return a
}

// lead prepends the symbol cell when the query named several symbols.
func (q query) lead(symbol string, cells ...any) []any {
	if !q.several {
		return cells
	}
	return append([]any{symbol}, cells...)
}

// fields prepends the symbol column when the query named several symbols.
func (q query) fields(schema []string) []string {
	if !q.several {
		return schema
	}
	return append([]string{"symbol"}, schema...)
}

// invocation is a replay of one operand under this query's filter, with any extra flags.
func (q query) invocation(symbol string, extra ...string) []axi.InvocationArgument {
	args := []axi.InvocationArgument{axi.KnownArgument("consumers"), axi.KnownArgument(symbol)}
	if q.scope != everyRow {
		args = append(args, axi.KnownArgument(string(q.scope)))
	}
	for _, flag := range extra {
		args = append(args, axi.KnownArgument(flag))
	}
	return args
}

// response renders every answer in operand order. An operand answers in exactly one
// block: an ambiguous name in consumers_candidates, an over-cap default in
// consumers_packages, and anything else in consumers. The consumers block renders
// whenever one operand answers there, so a zero-row answer stays the definitive empty
// table. The candidates answer is not a refusal: it offers one literal re-query per
// candidate row. An over-cap answer offers the one --full invocation that returns its
// complete set.
func (q query) response(source citation, pkgCount int, answers []answer) (string, int) {
	var consumerRows, packageRows, candidateRows [][]any
	var files []string
	var actions []axi.Action
	listed, capped, ambiguous, truncated := false, false, false, false
	matches, rowCount := 0, 0
	for _, a := range answers {
		matches += a.matches
		rowCount += len(a.rows)
		files = append(files, a.files...)
		truncated = truncated || a.overCap || a.dropped
		switch {
		case a.candidates != nil:
			ambiguous = true
			for _, c := range a.candidates {
				candidateRows = append(candidateRows, q.lead(a.symbol, c.qualified, c.file, c.line, c.kind))
				actions = append(actions, axi.ExecutableInvocation("re-query the qualified symbol", q.invocation(c.qualified)...))
			}
		case a.overCap:
			capped = true
			for _, row := range aggregate(a.files) {
				packageRows = append(packageRows, q.lead(a.symbol, row...))
			}
			actions = append(actions, axi.ExecutableInvocation("emit every consumer row", q.invocation(a.symbol, "--full")...))
		default:
			listed = true
			for _, r := range a.rows {
				consumerRows = append(consumerRows, q.lead(a.symbol, rowCells(r)...))
			}
		}
	}
	var block strings.Builder
	for _, table := range []struct {
		emit   bool
		name   string
		fields []string
		rows   [][]any
	}{
		{listed, rowBlock, rowFields, consumerRows},
		{capped, "consumers_packages", aggregateFields, packageRows},
		{ambiguous, "consumers_candidates", candidateFields, candidateRows},
	} {
		if !table.emit {
			continue
		}
		rendered, err := toon.TableTyped(table.name, q.fields(table.fields), table.rows)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		block.WriteString(rendered)
	}
	return envelope(source, block.String(), pkgCount, countFiles(files), matches, rowCount, truncated, actions)
}
