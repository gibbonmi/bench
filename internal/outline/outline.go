// Package outline implements `bench outline [path] [--full] [--production|--test]`, an
// on-demand repo seam map. It walks the tracked files, optionally scoped to a path, and
// runs a hand-rolled per-language pattern scan. It emits an AXI-conformant TOON table,
// so an agent can locate a candidate seam by name and jump to `file:line`. It
// regenerates on every call and writes nothing to the tree.
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
//
// `--production` drops every Go test file from the walk, and `--test` keeps only the Go
// test files. A filter scopes the walk as the path argument does, so every form and its
// metadata count only the files in scope.
package outline

import (
	"bytes"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// Symbol is one located declaration: the 1-based line it sits on, the kind the
// language pattern emits (func, type, function, heading, def, class, const), the
// candidate-inventory kind a name or a path refines it to (helper, double, fixture),
// and the declared name.
type Symbol struct {
	Line int
	Kind string
	Name string
}

// pattern is one ordered declaration form for a language: a compiled line-anchored
// regex, the kind it emits, the capture group holding the name, and an optional path
// predicate that limits the form to the paths it accepts. Adding a language or a form
// is adding a table entry — never a new code path.
type pattern struct {
	re    *regexp.Regexp
	kind  string
	group int
	path  func(string) bool // nil means every path the language table scans
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
		{re: regexp.MustCompile(`^func[ \t]+((?:new|make|with)[A-Z][A-Za-z0-9_]*)`), kind: "helper", group: 1, path: IsGoTestFile},
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
			if p.path != nil && !p.path(path) {
				continue
			}
			m := p.re.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			name := strings.TrimSpace(m[p.group])
			if name == "" {
				continue
			}
			out = append(out, Symbol{Line: i + 1, Kind: kindFor(p.kind, name), Name: name})
			break // one row per line: the first matching form wins
		}
	}
	return out
}

// GoTestFileSuffix names a Go test file: the go tool builds a file whose name ends in it
// only under go test.
const GoTestFileSuffix = "_test.go"

// IsGoTestFile is the one Go test-file predicate. It is the helper form's path predicate:
// a helper is a construction function a test file owns, so the same name in production
// code stays a plain func. The outline and consumers filters read it too.
func IsGoTestFile(path string) bool {
	return strings.HasSuffix(path, GoTestFileSuffix)
}

// doublePrefixes are the test-double name forms, matched case-insensitively.
var doublePrefixes = []string{"fake", "stub", "mock", "spy"}

// kindFor refines a matched row's kind: a name carrying a double prefix reports
// "double" whatever form matched it, in any file the language tables scan. The
// refinement rides the matched name rather than a double entry per language, because a
// per-language entry would restate that language's declaration grammar N times.
// Precedence: the pattern table's kind (helper included) resolves first, and the double
// prefix overrides it. The helper and double prefix sets are disjoint, so no name
// reaches both and the precedence is unobservable.
func kindFor(kind, name string) string {
	lower := strings.ToLower(name)
	for _, prefix := range doublePrefixes {
		if strings.HasPrefix(lower, prefix) {
			return "double"
		}
	}
	return kind
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

// walkSymbols is the walk's per-file dispatch. A fixture path answers with its one row
// before any content is read, so a binary, oversized, or otherwise unscannable fixture is
// still listed; only the stat policy can skip it. Every other file passes the whole read
// policy and goes through the language table.
func walkSymbols(rel, abs string) ([]Symbol, string, bool) {
	if IsFixture(rel) {
		if reason, ok := statPolicy(abs); !ok {
			return nil, reason, false
		}
		return []Symbol{{Line: 1, Kind: "fixture", Name: path.Base(rel)}}, "", true
	}
	content, reason, ok := ReadScannable(abs)
	if !ok {
		return nil, reason, false
	}
	return Symbols(rel, content), "", true
}

// TrackedFiles is the tracked-file census other packages read: git's ls-files order for
// the whole repository rooted at root. It is the listing half of the outline walk, and it
// is exported so a caller sweeping the same files scans one census rather than a second.
func TrackedFiles(root string) ([]string, error) {
	return listFiles(root, "", false)
}
