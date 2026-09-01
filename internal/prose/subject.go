package prose

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// Grader grades one named subject against the shared exclusion list, byte classifier,
// and finding renderer. Grade composes it for every subject the whole-tree walk finds;
// GradeNamed composes it for a caller-selected list, so the prose rule keeps one source.
type Grader struct {
	root string
	ex   *exclusions
}

// NamedResult is one result from grading a caller-selected path.
type NamedResult struct {
	Path       string
	Line       int
	Rule       FindingKind
	Count      int
	Sentence   string
	diagnostic string
}

// RenderNamedResult renders one named result without making callers reconstruct the
// prose diagnostic protocol.
func RenderNamedResult(result NamedResult) string {
	if result.diagnostic != "" {
		return result.diagnostic
	}
	diagnostic := Render(result.Path, Finding{Kind: result.Rule, Line: result.Line, Count: result.Count})
	if result.Sentence != "" {
		diagnostic += ": " + strconv.Quote(result.Sentence)
	}
	return diagnostic
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
	results := g.gradeSubjectResults(rel)
	out := make([]string, 0, len(results))
	for _, result := range results {
		result.Sentence = ""
		out = append(out, RenderNamedResult(result))
	}
	return out
}

// GradeNamed grades a caller-selected list of repository-relative paths through the same
// per-subject grader the whole-tree walk composes, so the prose rule keeps one source. A
// path absent from the composed tree is skipped: a commit that deletes a named file grades
// what it commits, not what it once named. A symbolic link is not followed and is not
// graded, matching the whole-tree walk's own rule for a linked directory.
func GradeNamed(root string, rels []string) []string {
	results := GradeNamedResults(root, rels)
	if len(results) == 0 {
		return nil
	}
	out := make([]string, 0, len(results))
	for _, result := range results {
		result.Sentence = ""
		out = append(out, RenderNamedResult(result))
	}
	return out
}

// GradeNamedResults grades caller-selected paths and exposes each prose finding as
// fields. Refusal diagnostics remain prose-owned rendered results.
func GradeNamedResults(root string, rels []string) []NamedResult {
	g, diags := NewGrader(root)
	if len(diags) > 0 {
		out := make([]NamedResult, 0, len(diags))
		for _, diagnostic := range diags {
			out = append(out, NamedResult{diagnostic: diagnostic})
		}
		return out
	}
	var out []NamedResult
	for _, rel := range rels {
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if g.ex.excluded(rel) {
			continue
		}
		out = append(out, g.gradeSubjectResults(rel)...)
	}
	return out
}

func (g *Grader) gradeSubjectResults(rel string) []NamedResult {
	classification := bounds.ClassifyNoFollow(filepath.Join(g.root, filepath.FromSlash(rel)))
	switch classification.State {
	case bounds.StateEmpty:
		return nil
	case bounds.StateParsed:
		lines := strings.Split(string(classification.Data), "\n")
		var out []NamedResult
		for _, finding := range Findings(string(classification.Data)) {
			result := NamedResult{Path: rel, Line: finding.Line, Rule: finding.Kind, Count: finding.Count}
			if finding.Kind == KindSentence && finding.Line <= len(lines) {
				result.Sentence = strings.TrimSpace(lines[finding.Line-1])
			}
			out = append(out, result)
		}
		return out
	case bounds.StateWrongType:
		return []NamedResult{{diagnostic: fmt.Sprintf("prose: %q: refused subject: %s", rel, classification.Reason)}}
	default:
		return []NamedResult{{diagnostic: fmt.Sprintf("prose: %q: refused unreadable subject: %s", rel, classification.Reason)}}
	}
}
