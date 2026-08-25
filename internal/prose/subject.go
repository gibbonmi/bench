package prose

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
)

// Grader grades one named subject against the shared exclusion list, byte classifier,
// and finding renderer. Grade composes it for every subject the whole-tree walk finds;
// GradeNamed composes it for a caller-selected list, so the prose rule keeps one source.
type Grader struct {
	root string
	ex   *exclusions
}

// NewGrader loads the exclusion list under root once, so a caller that grades many
// subjects pays the load cost once.
func NewGrader(root string) (*Grader, []string) {
	ex, diags := loadExclusions(root)
	if len(diags) > 0 {
		return nil, diags
	}
	return &Grader{root: root, ex: ex}, nil
}

// GradeSubject grades one repository-relative path and returns its findings. It returns
// nil for a subject that is excluded, empty, or clean.
func (g *Grader) GradeSubject(rel string) []string {
	if g.ex.excluded(rel) {
		return nil
	}
	c := bounds.ClassifyNoFollow(filepath.Join(g.root, filepath.FromSlash(rel)))
	switch c.State {
	case bounds.StateEmpty:
		return nil
	case bounds.StateParsed:
		var out []string
		for _, f := range Findings(string(c.Data)) {
			out = append(out, Render(rel, f))
		}
		return out
	case bounds.StateWrongType:
		return []string{fmt.Sprintf("prose: %q: refused subject: %s", rel, c.Reason)}
	default:
		return []string{fmt.Sprintf("prose: %q: refused unreadable subject: %s", rel, c.Reason)}
	}
}

// GradeNamed grades a caller-selected list of repository-relative paths through the same
// per-subject grader the whole-tree walk composes, so the prose rule keeps one source. A
// path absent from the composed tree is skipped: a commit that deletes a named file grades
// what it commits, not what it once named. A symbolic link is not followed and is not
// graded, matching the whole-tree walk's own rule for a linked directory.
func GradeNamed(root string, rels []string) []string {
	g, diags := NewGrader(root)
	if len(diags) > 0 {
		return diags
	}
	var out []string
	for _, rel := range rels {
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		out = append(out, g.GradeSubject(rel)...)
	}
	return out
}
