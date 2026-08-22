// Package outline implements `bench outline [path] [--full]`, an on-demand repo seam
// map. It walks the tracked files, optionally scoped to a path, and runs a hand-rolled
// per-language pattern scan. It emits an AXI-conformant TOON table, so an agent can
// locate a candidate seam by name and jump to `file:line`. It regenerates on every call
// and writes nothing to the tree.
//
// The form decides the table. A path argument or `--full` emits the symbol rows
// `outline[N]{file,line,kind,name}:`, scoped to the path, or repository-wide. The bare
// invocation emits the summary `outline_dirs[N]{dir,symbols}:` instead. This summary
// carries one row per scanned top-level directory with that directory's total symbol
// count, so a cold probe costs a screen. Every form emits its complete answer; no form
// truncates to a silent prefix.
//
// The tool LOCATES candidate seams; it does not IDENTIFY the project's blessed seams.
// `projects/<name>.md` owns that identification. The tool is a line-regex indexer, not
// a language parser. It does not track comments, string literals, or Markdown code
// fences, so a commented-out or fenced declaration indexes as a benign candidate. The
// agent confirms the candidate by reading the line.
package outline

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
// Help is helpText without its trailing newline, because the caller appends one.
var grammar = usage.Grammar{
	Cmd:     "bench outline",
	Help:    strings.TrimSuffix(helpText(), "\n"),
	Flags:   []usage.Flag{{Name: "--full"}},
	MaxArgs: 1,
}

// promise is the one-line expectation-setting clause the surface repeats: outline
// LOCATES seams, it does not IDENTIFY the project's blessed seams. It lands verbatim
// in the help text below and, as prose a shell surface cannot import, is restated
// with the same verbs in bin/bench.sh's help block.
const promise = "outline locates candidate seams (file:line); it does not identify which are the project's blessed seams — projects/<name>.md owns that."

const usageLine = "usage: bench outline [path] [--full]"

var openOutlineFile = os.Open

func helpText() string {
	return usageLine + "\n" + promise + "\n"
}

// Symbol is one located declaration: the 1-based line it sits on, the kind the
// language pattern emits (func, type, function, heading, def, class, const), and the
// declared name.
type Symbol struct {
	Line int
	Kind string
	Name string
}

// pattern is one ordered declaration form for a language: a compiled line-anchored
// regex, the kind it emits, and the capture group holding the name. Adding a language
// or a form is adding a table entry — never a new code path.
type pattern struct {
	re    *regexp.Regexp
	kind  string
	group int
}

// langTable is the one source of the per-language pattern fact: file extension →
// ordered declaration patterns. Per line the first matching pattern wins, so more
// specific forms (a Go method) precede the general ones (a plain Go func).
var langTable = map[string][]pattern{
	".sh": {
		{re: regexp.MustCompile(`^[ \t]*([A-Za-z_][A-Za-z0-9_-]*)[ \t]*\(\)`), kind: "function", group: 1},
		{re: regexp.MustCompile(`^[ \t]*function[ \t]+([A-Za-z_][A-Za-z0-9_-]*)`), kind: "function", group: 1},
	},
	".go": {
		{re: regexp.MustCompile(`^func[ \t]+\([^)]*\)[ \t]*([A-Za-z_][A-Za-z0-9_]*)`), kind: "func", group: 1},
		{re: regexp.MustCompile(`^func[ \t]+([A-Za-z_][A-Za-z0-9_]*)`), kind: "func", group: 1},
		{re: regexp.MustCompile(`^type[ \t]+([A-Za-z_][A-Za-z0-9_]*)`), kind: "type", group: 1},
	},
	".md": {
		{re: regexp.MustCompile(`^#{1,6}[ \t]+(.+?)[ \t]*#*[ \t]*$`), kind: "heading", group: 1},
	},
	".py": {
		{re: regexp.MustCompile(`^[ \t]*def[ \t]+([A-Za-z_][A-Za-z0-9_]*)`), kind: "def", group: 1},
		{re: regexp.MustCompile(`^[ \t]*class[ \t]+([A-Za-z_][A-Za-z0-9_]*)`), kind: "class", group: 1},
	},
}

// jsPatterns are shared by every JS/TS extension: a function declaration, a class
// declaration, and an exported const bound to an arrow function.
var jsPatterns = []pattern{
	{re: regexp.MustCompile(`^[ \t]*(?:export[ \t]+)?(?:default[ \t]+)?(?:async[ \t]+)?function[ \t*]+([A-Za-z_$][A-Za-z0-9_$]*)`), kind: "function", group: 1},
	{re: regexp.MustCompile(`^[ \t]*(?:export[ \t]+)?(?:default[ \t]+)?(?:abstract[ \t]+)?class[ \t]+([A-Za-z_$][A-Za-z0-9_$]*)`), kind: "class", group: 1},
	{re: regexp.MustCompile(`^[ \t]*export[ \t]+const[ \t]+([A-Za-z_$][A-Za-z0-9_$]*)[ \t]*=.*=>`), kind: "const", group: 1},
}

func init() {
	for _, ext := range []string{".js", ".jsx", ".ts", ".tsx"} {
		langTable[ext] = jsPatterns
	}
}

// Symbols scans content as the language selected by path's extension and returns its
// declarations in ascending line order. A path whose extension has no table entry
// yields no rows. A final line lacking a trailing newline is still a line and is
// still scanned.
func Symbols(path string, content []byte) []Symbol {
	patterns, ok := langTable[strings.ToLower(filepath.Ext(path))]
	if !ok {
		return nil
	}
	var out []Symbol
	for i, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSuffix(line, "\r")
		for _, p := range patterns {
			m := p.re.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			name := strings.TrimSpace(m[p.group])
			if name == "" {
				continue
			}
			out = append(out, Symbol{Line: i + 1, Kind: p.kind, Name: name})
			break // one row per line: the first matching form wins
		}
	}
	return out
}

// listFiles returns the tracked files git reports, root-relative, in git's ls-files
// order. With no path argument it lists the whole repo. With one, it resolves the
// argument from the process cwd to a root-relative `:(literal,top)` pathspec. A glob
// character or space in the path then matches literally, a directory scopes to the
// files beneath it, and a path outside the repo scopes to nothing. The -z framing is
// NUL-delimited and never C-quotes, so a path with spaces or an embedded newline
// survives whole.
func listFiles(root, path string, havePath bool) ([]string, error) {
	args := []string{"-C", root, "ls-files", "-z"}
	if havePath {
		abs, err := filepath.Abs(path) // resolved against the process cwd
		if err != nil {
			return nil, err
		}
		rel, err := filepath.Rel(root, abs)
		if err != nil {
			return nil, err
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil, nil // outside the repo → scopes to nothing
		}
		if rel != "." {
			args = append(args, "--", ":(literal,top)"+filepath.ToSlash(rel))
		}
	}
	raw, err := git.Raw(args...)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, rec := range bytes.Split(raw, []byte{0}) {
		if len(rec) == 0 {
			continue // trailing empty after the final NUL
		}
		files = append(files, string(rec))
	}
	return files, nil
}

// Command implements `bench outline [path] [--full]`. It walks the tracked files, scoped
// to an optional path, and dispatches each by extension through Symbols. It drops any
// row a control byte would make unrepresentable, and renders the symbol table for a
// path or `--full`, and the per-directory summary for the bare form. Each form has the
// definitive empty state when nothing matches, a structured stdout error with exit 1
// outside a repo or on a git failure, and usage on stdout with exit 2 for an unknown
// flag or a second positional argument.
func Command(args []string) (string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, full := parsed.Flags["--full"]
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
		info, statErr := os.Lstat(abs)
		if statErr != nil {
			skips = appendSkip(skips, rel, "unreadable")
			continue
		}
		if !info.Mode().IsRegular() {
			skips = appendSkip(skips, rel, "nonregular")
			continue
		}
		file, err := openOutlineFile(abs)
		if err != nil {
			skips = appendSkip(skips, rel, "unreadable")
			continue
		}
		read := bounds.Read(file, bounds.OutlineFileLimit)
		closeErr := file.Close()
		if read.Status == bounds.ReadOversized {
			skips = appendSkip(skips, rel, "oversized")
			continue
		}
		if read.Status == bounds.ReadFailed || closeErr != nil {
			skips = appendSkip(skips, rel, "unreadable")
			continue
		}
		content := read.Data
		if bytes.IndexByte(content, 0) >= 0 {
			skips = appendSkip(skips, rel, "binary")
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
		for _, s := range Symbols(rel, content) {
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

// topLevel is the bare summary's grouping key: the first segment of git's
// slash-separated path, or "." for a file at the repository root. The summary collapses
// to this depth because a row per scanned directory is larger than the symbol rows it
// replaces once a tree is deep, which defeats the one-screen probe it exists to give.
func topLevel(rel string) string {
	if i := strings.Index(rel, "/"); i >= 0 {
		return rel[:i]
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
