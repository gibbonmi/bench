// Package spec owns spec-file addressing and the `bench spec implemented <slug>` flip.
// Resolve is the one source of the spec-argument convention (path-first, then a
// specs/<slug>.md fallback) that both `bench coverage` and `bench commit --spec` take
// their argument through. Flip is the single source of the status-line flip: it turns
// exactly one line-start `Status: staged` into the retirement-detector form
// `Status: implemented`, preserving every other byte, and is composed by `bench commit`.
package spec

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// stagedRe matches the flip's input: a `Status:` line whose sole value is `staged`,
// tab/space separated, only whitespace trailing — the `staged` twin of the retirement
// detector's `^Status:[ \t]+implemented[ \t]*$`. Swapping `staged` for `implemented` in
// place therefore yields exactly the detector's accepted form by construction.
var stagedRe = regexp.MustCompile(`^Status:[ \t]+staged[ \t]*$`)

// Resolve finds the readable file backing a spec argument: the argument as given
// (path-first, so a same-named readable file shadows the fallback), then — for a
// separator-free argument only — a specs/<slug>.md fallback, appending .md only when the
// argument doesn't already end in it. base anchors that fallback: pass the repo root to
// resolve it repo-root-relative (so a cwd deeper than the root still finds it), or "" to
// resolve it relative to the process cwd. ok is false when no form resolves; tried holds
// every form attempted, for the not-found error. A non-nil err is a read failure on an
// existing file (e.g. permissions), reported instead of not-found.
func Resolve(base, arg string) (content []byte, resolved string, tried []string, ok bool, err error) {
	tried = []string{arg}
	if b, err := readCandidate(arg); err != nil || b != nil {
		return b, arg, tried, err == nil, err
	}
	if !strings.ContainsRune(arg, '/') {
		slug := arg
		if !strings.HasSuffix(slug, ".md") {
			slug += ".md"
		}
		fallback := filepath.Join(base, "specs", slug)
		tried = append(tried, fallback)
		if b, err := readCandidate(fallback); err != nil || b != nil {
			return b, fallback, tried, err == nil, err
		}
	}
	return nil, "", tried, false, nil
}

// readCandidate reads path as a candidate spec. An absent path or a directory is not a
// candidate (nil, nil — try the next form); any other read failure is a real error to
// surface, never masked as not-found.
func readCandidate(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err == nil {
		return b, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if fi, statErr := os.Stat(path); statErr == nil && fi.IsDir() {
		return nil, nil
	}
	return nil, err
}

// locateStaged resolves arg (base anchors the specs/<slug>.md fallback) and requires exactly
// one line-start `Status: staged`, returning the resolved path, the file split into lines, and
// the index of that one line. It never writes — it is the shared core of CheckStaged
// (validate only) and Flip (validate then rewrite), so the resolution + single-staged-line
// rule has one source. The error names the file and the reason on not-found, not-readable, no
// `Status: staged` line (missing or already implemented), or more than one.
func locateStaged(base, arg string) (resolved string, lines [][]byte, idx int, err error) {
	content, resolved, tried, ok, readErr := Resolve(base, arg)
	if readErr != nil {
		return "", nil, -1, fmt.Errorf("spec not readable: %s: %v", resolved, readErr)
	}
	if !ok {
		return "", nil, -1, fmt.Errorf("spec not found: %s", strings.Join(tried, ", "))
	}
	lines = bytes.Split(content, []byte("\n"))
	matches := 0
	idx = -1
	for i, line := range lines {
		if stagedRe.Match(line) {
			matches++
			idx = i
		}
	}
	if matches == 0 {
		return "", nil, -1, fmt.Errorf("no `Status: staged` line in %s (already implemented, or missing)", resolved)
	}
	if matches > 1 {
		return "", nil, -1, fmt.Errorf("%d `Status: staged` lines in %s (expected exactly one)", matches, resolved)
	}
	return resolved, lines, idx, nil
}

// CheckStaged resolves arg and confirms it carries exactly one line-start `Status: staged`,
// returning the resolved path. It never writes: it is the fail-fast validation `bench commit
// --spec` runs before the gate, so a bad or already-implemented spec is rejected before the
// gate burns rather than after Flip runs on a green tree. Flip re-checks the (possibly changed)
// file before it rewrites, so the two share locateStaged's one validation rule.
func CheckStaged(base, arg string) (resolved string, err error) {
	resolved, _, _, err = locateStaged(base, arg)
	return resolved, err
}

// Flip resolves arg (base anchors the specs/<slug>.md fallback), requires exactly one
// line-start `Status: staged`, rewrites that one line to `Status: implemented` in place —
// every other byte, including a missing final newline, preserved — writes the file, and
// returns the resolved path. It edits the file only; it never stages. The error names the
// file and the reason on not-found, not-readable, no `Status: staged` line (missing or
// already implemented), or more than one — so a typo or a re-run is non-destructive.
func Flip(base, arg string) (resolved string, err error) {
	resolved, lines, idx, err := locateStaged(base, arg)
	if err != nil {
		return "", err
	}
	lines[idx] = bytes.Replace(lines[idx], []byte("staged"), []byte("implemented"), 1)
	out := bytes.Join(lines, []byte("\n"))
	mode := os.FileMode(0o644)
	if fi, statErr := os.Stat(resolved); statErr == nil {
		mode = fi.Mode().Perm()
	}
	if err := os.WriteFile(resolved, out, mode); err != nil {
		return "", fmt.Errorf("write %s: %v", resolved, err)
	}
	return resolved, nil
}

// Command implements `bench spec implemented <slug>`. The slug's specs/<slug>.md fallback
// is anchored at the repo root, so it resolves from any cwd inside the repo; a path
// argument stays cwd-relative. Usage errors (missing subcommand, unknown subcommand, no
// argument) exit 2; a resolve/validate/write failure exits 1 naming the file and reason.
func Command(args []string) (string, int) {
	if len(args) == 0 {
		return toon.Usage("bench spec", "expected a subcommand: implemented") + "\n", 2
	}
	if args[0] != "implemented" {
		return toon.Usage("bench spec", args[0]) + "\n", 2
	}
	rest := args[1:]
	arg := ""
	for _, a := range rest {
		switch {
		case a == "-h" || a == "--help":
			return "usage: bench spec implemented <spec.md | slug>\n", 0
		case strings.HasPrefix(a, "-"):
			return toon.Usage("bench spec implemented", a) + "\n", 2
		default:
			if arg != "" {
				return toon.Usage("bench spec implemented", a) + "\n", 2
			}
			arg = a
		}
	}
	if arg == "" {
		return toon.Usage("bench spec implemented", "<spec.md | slug> is required") + "\n", 2
	}
	base := ""
	if root, err := git.Root(); err == nil {
		base = root
	}
	resolved, err := Flip(base, arg)
	if err != nil {
		return toon.Errorf(err.Error(), "pass a spec with a single `Status: staged` line") + "\n", 1
	}
	return fmt.Sprintf("spec implemented: %s\n", resolved), 0
}
