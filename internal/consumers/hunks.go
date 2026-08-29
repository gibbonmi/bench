package consumers

import (
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// lineSpan is one inclusive run of lines in one revision of one file. A zero-count hunk
// side is the one case the two sides read differently, and this is the one statement of
// that rule. A zero-count tip side is a deletion-only edit: git points it at the tip line
// the removed run sat after, so the run is that line and the next one, and a declaration
// that survived the deletion is still marked touched. A zero-count base side is an
// insertion, which dropped no base line, so it names no removed run. An empty span never
// reaches an intersection test.
type lineSpan struct {
	Start, End int
}

// fileHunks is one file's contribution to a unified diff: the two paths the diff named
// and the line runs each side carries. Added holds the tip-side runs a hunk introduced,
// and Removed holds the base-side runs it dropped. A pure rename with no edit yields
// neither, so it names no declaration.
type fileHunks struct {
	BasePath, TipPath string
	Added, Removed    []lineSpan
}

// parseHunks reads `git diff -U0` text into per-file line runs. It is a pure function of
// the diff bytes, so the blast derivation is testable without a repository. The parse
// reads only the two file headers and the hunk headers: with -U0 every remaining line is
// content, and content never starts a new file or a new run.
func parseHunks(text string) []fileHunks {
	var out []fileHunks
	var current *fileHunks
	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			out = append(out, fileHunks{})
			current = &out[len(out)-1]
		case current == nil:
			// Text before the first file header belongs to no file.
		case strings.HasPrefix(line, "--- "):
			current.BasePath = diffPath(line[4:])
		case strings.HasPrefix(line, "+++ "):
			current.TipPath = diffPath(line[4:])
		case strings.HasPrefix(line, "@@ "):
			base, tip, ok := parseHunkHeader(line)
			if !ok {
				continue
			}
			if base.End >= base.Start {
				current.Removed = append(current.Removed, base)
			}
			if tip.End >= tip.Start {
				current.Added = append(current.Added, tip)
			}
		}
	}
	return out
}

// diffPath strips the a/ or b/ prefix git prints. A missing side reads as /dev/null and
// keeps the empty string, so an added or deleted file has one nameless side.
func diffPath(field string) string {
	field = strings.TrimSpace(field)
	if field == "/dev/null" {
		return ""
	}
	if len(field) > 2 && (field[:2] == "a/" || field[:2] == "b/") {
		return field[2:]
	}
	return field
}

// parseHunkHeader reads `@@ -a,b +c,d @@` into the base run and the tip run. Each side is
// widened by the rule the lineSpan comment states, so the header parse holds no second
// reading of a zero count.
func parseHunkHeader(line string) (base, tip lineSpan, ok bool) {
	fields := strings.Fields(line)
	if len(fields) < 3 || !strings.HasPrefix(fields[1], "-") || !strings.HasPrefix(fields[2], "+") {
		return lineSpan{}, lineSpan{}, false
	}
	baseStart, baseCount, ok := parseSpan(fields[1][1:])
	if !ok {
		return lineSpan{}, lineSpan{}, false
	}
	tipStart, tipCount, ok := parseSpan(fields[2][1:])
	if !ok {
		return lineSpan{}, lineSpan{}, false
	}
	return removedRun(baseStart, baseCount), addedRun(tipStart, tipCount), true
}

// parseSpan reads one `start[,count]` hunk side. An omitted count is one line, which is
// git's own default.
func parseSpan(field string) (start, count int, ok bool) {
	startText, countText, hasCount := strings.Cut(field, ",")
	start, err := strconv.Atoi(startText)
	if err != nil {
		return 0, 0, false
	}
	count = 1
	if hasCount {
		if count, err = strconv.Atoi(countText); err != nil {
			return 0, 0, false
		}
	}
	return start, count, true
}

// emptyRun is the span an intersection test never meets.
var emptyRun = lineSpan{Start: 1, End: 0}

// removedRun is the base-side run the hunk dropped.
func removedRun(start, count int) lineSpan {
	if count <= 0 {
		return emptyRun
	}
	return lineSpan{Start: start, End: start + count - 1}
}

// addedRun is the tip-side run the hunk introduced. A zero start names a tip that holds no
// lines at all, which is the deleted-file case, so it names no run.
func addedRun(start, count int) lineSpan {
	switch {
	case count > 0:
		return lineSpan{Start: start, End: start + count - 1}
	case start < 1:
		return emptyRun
	}
	return lineSpan{Start: start, End: start + 1}
}

// intersects reports whether the inclusive run [start,end] meets any span. It is the one
// derivation of "the diff touched these lines", and both the declaration test and the
// deleted-declaration test read it.
func intersects(spans []lineSpan, start, end int) bool {
	for _, s := range spans {
		if start <= s.End && s.Start <= end {
			return true
		}
	}
	return false
}

// goPaths keeps the Go files of a changed path set. It is the one derivation of "this pair
// changed Go source", so the empty-answer shortcut and the diff read never disagree.
func goPaths(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if strings.HasSuffix(p, ".go") {
			out = append(out, p)
		}
	}
	return out
}

// readHunks is the diff half of the git rim. It asks for zero context lines, so every
// reported run is an edit rather than its neighborhood, and it restricts the diff to the
// Go paths the pair changed, so an unrelated file never reaches the parse.
//
// core.quotePath is disabled for the run, so a path carrying a non-ASCII byte arrives raw
// and matches the spelling the loader reports for the same file. Git still C-quotes the
// narrower classes it always quotes, and such a path reaches the response path and drops
// there as an unrepresentable row rather than failing the whole answer.
func readHunks(root, base, tip string, paths []string) ([]fileHunks, error) {
	goPaths := goPaths(paths)
	if len(goPaths) == 0 {
		return nil, nil
	}
	args := append([]string{"-c", "core.quotePath=false", "-C", root, "diff", "-U0", "--no-color", base, tip, "--"}, goPaths...)
	out, err := git.Raw(args...)
	if err != nil {
		return nil, err
	}
	return parseHunks(string(out)), nil
}

// readBaseSources reads the base-side bytes of every file the diff removed lines from.
// A path with no base-side name was added by the pair, so it has nothing to read and can
// declare nothing the tip deleted.
func readBaseSources(root, base string, hunks []fileHunks) map[string]string {
	sources := map[string]string{}
	for _, fh := range hunks {
		if fh.BasePath == "" || len(fh.Removed) == 0 {
			continue
		}
		out, err := git.Raw("-C", root, "show", base+":"+fh.BasePath)
		if err != nil {
			continue
		}
		sources[fh.BasePath] = string(out)
	}
	return sources
}
