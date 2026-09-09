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

// targetKind is what an exclusion target names in the source the policy validates
// against. The row grammar reads the same three answers from a working tree and from a
// Git index, so the validation core keys on them rather than on a stat.
type targetKind int

const (
	targetAbsent targetKind = iota
	targetFile
	targetDirectory
)

// targetLookup answers what one exclusion target names. The target arrives with any
// trailing slash already removed, so a directory row and a file row ask the same question.
type targetLookup func(target string) targetKind

// absentExclusionDiagnostic is the one sentence an absent policy answers. The working-tree
// loader and the index loader both state it, so the diagnostic has one source.
func absentExclusionDiagnostic() string {
	return fmt.Sprintf("prose: %q: the exclusion file is absent", ExclusionFile)
}

// UnreadableExclusionDiagnostic is the one sentence a policy that cannot be read answers,
// with reason naming what stopped the read. The working-tree loader states it for a
// classified file, and a staged caller states it for a classified index blob, so the two
// forms refuse a policy in the same words.
func UnreadableExclusionDiagnostic(reason string) string {
	return fmt.Sprintf("prose: %q: refused unreadable exclusion file: %s", ExclusionFile, reason)
}

// worktreeLookup validates a target against the files under root.
func worktreeLookup(root string) targetLookup {
	return func(target string) targetKind {
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(target)))
		switch {
		case err != nil:
			return targetAbsent
		case info.IsDir():
			return targetDirectory
		default:
			return targetFile
		}
	}
}

// indexLookup validates a target against a Git index entry list. An index holds no
// directory record, so a target that is the prefix of some entry is the directory.
func indexLookup(entries []string) targetLookup {
	files := make(map[string]bool, len(entries))
	for _, entry := range entries {
		files[entry] = true
	}
	return func(target string) targetKind {
		if files[target] {
			return targetFile
		}
		prefix := target + "/"
		for _, entry := range entries {
			if strings.HasPrefix(entry, prefix) {
				return targetDirectory
			}
		}
		return targetAbsent
	}
}

// loadExclusions reads and grades the row list under root. The row grammar is a
// repository-relative path, one space, and a one-clause reason, with a `#` comment and a
// blank line ignored.
func loadExclusions(root string) (*exclusions, []string) {
	c := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(ExclusionFile)))
	switch c.State {
	case bounds.StateAbsent:
		return nil, []string{absentExclusionDiagnostic()}
	case bounds.StateEmpty:
		return &exclusions{files: map[string]bool{}}, nil
	case bounds.StateParsed:
		return gradeExclusionRows(c.Data, worktreeLookup(root))
	default:
		return nil, []string{UnreadableExclusionDiagnostic(c.Reason)}
	}
}

// gradeExclusionRows grades a policy body against the caller's target source. It returns
// one diagnostic for each broken state and refuses the whole list on any of them: a list
// the parser cannot trust would silently widen or narrow the grade.
func gradeExclusionRows(data []byte, lookup targetLookup) (*exclusions, []string) {
	out := &exclusions{files: map[string]bool{}}
	var diags []string
	seen := map[string]int{}
	for _, row := range splitExclusionRows(string(data)) {
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
		kind := lookup(strings.TrimSuffix(subject, "/"))
		if kind == targetAbsent {
			diags = append(diags, fmt.Sprintf("prose: %q line %d: exclusion row %q names an absent path", ExclusionFile, number, subject))
			continue
		}
		if strings.HasSuffix(subject, "/") {
			out.prefixes = append(out.prefixes, subject)
			continue
		}
		// A directory row with no trailing slash excludes nothing, because the walk compares
		// a prefix. The row reds rather than passing as a file row that never matches.
		if kind == targetDirectory {
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
