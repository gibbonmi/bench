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

// exclusionRow is one row the grammar keeps: the subject, the reason, and the 1-based
// physical line. A comment row and a blank row hold no subject and are dropped.
type exclusionRow struct {
	subject string
	reason  string
	number  int
}

// splitExclusionRows applies the row grammar to the file body. It is the one reader of
// that grammar: the engine grades these rows, and the conformance test reads the same
// subjects through ExclusionRows.
func splitExclusionRows(data string) []exclusionRow {
	var rows []exclusionRow
	for i, line := range strings.Split(data, "\n") {
		if idx := strings.IndexByte(line, '#'); idx >= 0 {
			line = line[:idx]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		row := exclusionRow{subject: trimmed, number: i + 1}
		if idx := strings.IndexAny(trimmed, " \t"); idx >= 0 {
			row.subject = trimmed[:idx]
			row.reason = strings.TrimSpace(trimmed[idx+1:])
		}
		rows = append(rows, row)
	}
	return rows
}

// ExclusionRows returns the subject of every row of the exclusion file under root. It
// is the read seam for a caller that grades the row set rather than the documents, so
// the grammar stays in this package.
func ExclusionRows(root string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(ExclusionFile)))
	if err != nil {
		return nil, err
	}
	var subjects []string
	for _, row := range splitExclusionRows(string(data)) {
		subjects = append(subjects, row.subject)
	}
	return subjects, nil
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
	for _, row := range splitExclusionRows(string(c.Data)) {
		subject, number := row.subject, row.number
		if row.reason == "" {
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
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(strings.TrimSuffix(subject, "/"))))
		if err != nil {
			diags = append(diags, fmt.Sprintf("prose: %q line %d: exclusion row %q names an absent path", ExclusionFile, number, subject))
			continue
		}
		if strings.HasSuffix(subject, "/") {
			out.prefixes = append(out.prefixes, subject)
			continue
		}
		// A directory row with no trailing slash excludes nothing, because the walk compares
		// a prefix. The row reds rather than passing as a file row that never matches.
		if info.IsDir() {
			diags = append(diags, fmt.Sprintf("prose: %q line %d: exclusion row %q names a directory: a directory row needs a trailing slash", ExclusionFile, number, subject))
			continue
		}
		out.files[subject] = true
	}
	if len(diags) > 0 {
		return nil, diags
	}
	return out, nil
}
