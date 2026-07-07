// Package outline implements `bench outline [path]` — an on-demand repo seam map.
// It walks the tracked files (optionally scoped to a path), runs a hand-rolled
// per-language pattern scan, and emits an AXI-conformant `outline[N]{file,line,kind,name}:`
// TOON table so an agent can locate a candidate seam by name and jump to `file:line`.
// It is regenerated on every call and writes nothing to the tree.
//
// The tool LOCATES candidate seams; it does not IDENTIFY which are the project's
// blessed seams — `projects/<name>.md` owns that. It is a line-regex indexer, not a
// language parser: it does not track comments, string literals, or Markdown code
// fences, so a commented-out or fenced declaration is indexed as a benign candidate
// the agent confirms by reading the line.
package outline

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// promise is the one-line expectation-setting clause the surface repeats: outline
// LOCATES seams, it does not IDENTIFY the project's blessed seams. It lands verbatim
// in the help text below and, as prose a shell surface cannot import, is restated
// with the same verbs in bin/bench.sh's help block.
const promise = "outline locates candidate seams (file:line); it does not identify which are the project's blessed seams — projects/<name>.md owns that."

const usageLine = "usage: bench outline [path]"

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
// order. With no path argument it is the whole repo; with one, the argument is
// resolved from the process cwd to a root-relative `:(literal,top)` pathspec so a glob
// character or space in the path is matched literally, a directory scopes to the files
// beneath it, and a path outside the repo scopes to nothing. The -z framing is
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

// Command implements `bench outline [path]`: it walks the tracked files (scoped to an
// optional path), dispatches each by extension through Symbols, drops any row a control
// byte would make unrepresentable, and renders the `outline[N]{file,line,kind,name}:`
// TOON table — the definitive empty state when nothing matches, a structured stdout
// error with exit 1 outside a repo or on a git failure, and usage on stdout with exit 2
// for an unknown flag or a second positional argument.
func Command(args []string) (string, int) {
	var path string
	var havePath bool
	for _, a := range args {
		switch {
		case a == "-h" || a == "--help":
			return helpText(), 0
		case strings.HasPrefix(a, "-"):
			return toon.Usage("bench outline", a) + "\n", 2
		default:
			if havePath {
				return toon.Usage("bench outline", a) + "\n", 2
			}
			path, havePath = a, true
		}
	}

	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	files, err := listFiles(root, path, havePath)
	if err != nil {
		return toon.Errorf("git ls-files failed", err.Error()) + "\n", 1
	}

	var rows [][]string
	for _, rel := range files {
		abs := filepath.Join(root, filepath.FromSlash(rel))
		if info, err := os.Lstat(abs); err != nil || !info.Mode().IsRegular() {
			continue // a symlink's tracked content is its target string, not the
			// target's declarations — indexing through it would emit file:line
			// anchors that don't hold; non-regular entries contribute no rows
		}
		content, err := os.ReadFile(abs)
		if err != nil {
			continue // an absent/unreadable tracked path contributes no rows
		}
		if bytes.IndexByte(content, 0) >= 0 {
			continue // a NUL byte means binary → skip
		}
		for _, s := range Symbols(rel, content) {
			if !toon.Representable(rel) || !toon.Representable(s.Name) {
				continue // one poisoned path or name drops only its own row
			}
			rows = append(rows, []string{rel, strconv.Itoa(s.Line), s.Kind, s.Name})
		}
	}

	tbl, err := toon.Table("outline", []string{"file", "line", "kind", "name"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return tbl, 0
}
