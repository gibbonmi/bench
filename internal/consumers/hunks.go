package consumers

import (
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// lineSpan is one inclusive run of lines in one revision of one file. A hunk with a zero
// count contributes no span, so an empty span never reaches an intersection test.
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

// parseHunkHeader reads `@@ -a,b +c,d @@` into the base run and the tip run. An omitted
// count is one line, which is git's own default. A count of zero yields an empty span:
// git points a zero-count side at the line before the change, and treating that as a
// one-line run would touch a declaration the hunk never edited.
func parseHunkHeader(line string) (base, tip lineSpan, ok bool) {
	fields := strings.Fields(line)
	if len(fields) < 3 || !strings.HasPrefix(fields[1], "-") || !strings.HasPrefix(fields[2], "+") {
		return lineSpan{}, lineSpan{}, false
	}
	base, ok = parseSpan(fields[1][1:])
	if !ok {
		return lineSpan{}, lineSpan{}, false
	}
	tip, ok = parseSpan(fields[2][1:])
	if !ok {
		return lineSpan{}, lineSpan{}, false
	}
	return base, tip, true
}

func parseSpan(field string) (lineSpan, bool) {
	startText, countText, hasCount := strings.Cut(field, ",")
	start, err := strconv.Atoi(startText)
	if err != nil {
		return lineSpan{}, false
	}
	count := 1
	if hasCount {
		count, err = strconv.Atoi(countText)
		if err != nil {
			return lineSpan{}, false
		}
	}
	return lineSpan{Start: start, End: start + count - 1}, true
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

// readHunks is the diff half of the git rim. It asks for zero context lines, so every
// reported run is an edit rather than its neighborhood, and it restricts the diff to the
// Go paths the pair changed, so an unrelated file never reaches the parse.
func readHunks(root, base, tip string, paths []string) ([]fileHunks, error) {
	goPaths := make([]string, 0, len(paths))
	for _, p := range paths {
		if strings.HasSuffix(p, ".go") {
			goPaths = append(goPaths, p)
		}
	}
	if len(goPaths) == 0 {
		return nil, nil
	}
	args := append([]string{"-C", root, "diff", "-U0", "--no-color", base, tip, "--"}, goPaths...)
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
