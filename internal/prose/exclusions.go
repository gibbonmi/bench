package prose

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// ExclusionFile is the repository-relative path of the row list that names what the
// walk does not grade.
const ExclusionFile = ".bench/prose-exclusions"

// exclusions is the parsed row list. A row that ends with a slash names a directory
// prefix; every other row names one file.
type exclusions struct {
	prefixes []string
	files    map[string]bool
}

// excluded reports whether the subject at the repository-relative path rel is outside
// the grade.
func (e *exclusions) excluded(rel string) bool {
	if e == nil {
		return false
	}
	if e.files[rel] {
		return true
	}
	for _, prefix := range e.prefixes {
		if strings.HasPrefix(rel, prefix) {
			return true
		}
	}
	return false
}

// loadExclusions reads and grades the row list under root. It returns one diagnostic for
// each broken state and refuses the whole list on any of them: a list the parser cannot
// trust would silently widen or narrow the grade. The row grammar is a repository-relative
// path, one space, and a one-clause reason, with a `#` comment and a blank line ignored.
func loadExclusions(root string) (*exclusions, []string) {
	path := filepath.Join(root, filepath.FromSlash(ExclusionFile))
	c := bounds.ClassifyNoFollow(path)
	switch c.State {
	case bounds.StateAbsent:
		return nil, []string{fmt.Sprintf("prose: %q: the exclusion file is absent", ExclusionFile)}
	case bounds.StateEmpty:
		return &exclusions{files: map[string]bool{}}, nil
	case bounds.StateParsed:
	default:
		return nil, []string{fmt.Sprintf("prose: %q: refused unreadable exclusion file: %s", ExclusionFile, c.Reason)}
	}

	out := &exclusions{files: map[string]bool{}}
	var diags []string
	seen := map[string]int{}
	for i, line := range strings.Split(string(c.Data), "\n") {
		number := i + 1
		if idx := strings.IndexByte(line, '#'); idx >= 0 {
			line = line[:idx]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		subject, reason := trimmed, ""
		if idx := strings.IndexAny(trimmed, " \t"); idx >= 0 {
			subject, reason = trimmed[:idx], strings.TrimSpace(trimmed[idx+1:])
		}
		if reason == "" {
			diags = append(diags, fmt.Sprintf("prose: %q line %d: malformed exclusion row: the reason is absent", ExclusionFile, number))
			continue
		}
		if strings.ContainsAny(subject, "*?[") {
			diags = append(diags, fmt.Sprintf("prose: %q line %d: exclusion row %q uses a glob character", ExclusionFile, number, subject))
			continue
		}
		if first, dup := seen[subject]; dup {
			diags = append(diags, fmt.Sprintf("prose: %q line %d: duplicate exclusion row %q first named on line %d", ExclusionFile, number, subject, first))
			continue
		}
		seen[subject] = number
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(strings.TrimSuffix(subject, "/")))); err != nil {
			diags = append(diags, fmt.Sprintf("prose: %q line %d: exclusion row %q names an absent path", ExclusionFile, number, subject))
			continue
		}
		if strings.HasSuffix(subject, "/") {
			out.prefixes = append(out.prefixes, subject)
			continue
		}
		out.files[subject] = true
	}
	if len(diags) > 0 {
		return nil, diags
	}
	return out, nil
}
