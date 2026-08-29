package consumers

import (
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

func helpText() string {
	return usageLine + "\n" + promise + "\n" + soundness + "\n" + capLine() + "\n"
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

// aggregateFields is the over-cap schema: one row per consumer directory, which is one Go
// package, with the row count that directory contributed. The queried symbol is constant
// across the block, so it gets no column.
var aggregateFields = []string{"dir", "rows"}

// metaFields is the response accounting every form carries.
var metaFields = []string{"packages", "files", "matches", "rows", "truncated"}

// Command implements `bench consumers <qualified-symbol> [--full]`. It resolves the
// symbol over the repository's packages and emits the consumers table, the meta
// accounting, and the terminal help envelope. A symbol result is a terminal read, so its
// envelope is empty unless the default truncated.
func Command(args []string) (string, int) {
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
	if len(matches) != 1 {
		return toon.Errorf("ambiguous symbol "+symbol, "qualify the symbol with its package, such as "+matches[0].Qualified) + "\n", 1
	}
	return response(symbol, len(pkgs), len(matches), Rows(pkgs, matches[0].Obj, root), full)
}

// response renders the whole answer for one resolved symbol: the result block, the meta
// accounting, and the help envelope. Row rendering itself belongs to the core, so this
// function composes Render rather than restating the row schema.
func response(symbol string, pkgCount, matchCount int, rows []Row, full bool) (string, int) {
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
	// The counts and the flag are typed cells, so an integer stays bare and the flag stays
	// a boolean through a TOON round-trip.
	meta, err := toon.TableTyped("meta", metaFields, [][]any{{
		pkgCount, countFiles(rows), matchCount, len(rows), truncated,
	}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	var actions []axi.Action
	if truncated {
		actions = append(actions, axi.ExecutableInvocation("emit every consumer row",
			axi.KnownArgument("consumers"), axi.KnownArgument(symbol), axi.KnownArgument("--full")))
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return block + meta + help, 0
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
