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
	Starts     []SentenceStart
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
	if len(result.Starts) > 0 {
		items := make([]string, 0, len(result.Starts))
		for _, start := range result.Starts {
			// The start is quoted the way the sentence text is quoted, so a control byte inside
			// it is escaped and one finding stays one line.
			items = append(items, strconv.Itoa(start.Line)+" "+strconv.Quote(start.Text))
		}
		diagnostic += ": sentences " + strings.Join(items, ", ")
	}
	return diagnostic
}

// renderStripped renders results without the document text. The whole-tree grade and the
// `prose` named check state the counts alone, so neither the sentence nor a start reaches
// the gate output.
func renderStripped(results []NamedResult) []string {
	if len(results) == 0 {
		return nil
	}
	out := make([]string, 0, len(results))
	for _, result := range results {
		result.Sentence, result.Starts = "", nil
		out = append(out, RenderNamedResult(result))
	}
	return out
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

// IndexSource is the staged material a policy needs: the exclusion file's index blob,
// whether the index holds that entry at all, and every index entry path. The target
// validation reads Entries, so a policy that names a path present in the index and absent
// from the working tree is honored.
type IndexSource struct {
	Policy        []byte
	PolicyPresent bool
	Entries       []string
}

// NewGraderFromIndex builds a grader whose exclusion policy is the index blob rather than
// the working file, and whose targets validate against the index entry list. An index
// without the exclusion entry answers the same absent-policy diagnostic the working-tree
// loader states, because a staged commit with no policy is as ungradeable as a tree
// without one.
func NewGraderFromIndex(src IndexSource) (*Grader, []string) {
	if !src.PolicyPresent {
		return nil, []string{absentExclusionDiagnostic()}
	}
	if len(src.Policy) == 0 {
		return &Grader{ex: &exclusions{files: map[string]bool{}}}, nil
	}
	ex, diags := gradeExclusionRows(src.Policy, indexLookup(src.Entries))
	if len(diags) > 0 {
		return nil, diags
	}
	return &Grader{ex: ex}, nil
}

// GradeBytes grades caller-supplied bytes as the subject at the repository-relative path
// rel. A staged grade reads the index blob, so the subject never comes from a file the
// grader could open itself. It returns nil for a subject that is excluded, empty, or
// clean.
func (g *Grader) GradeBytes(rel string, data []byte) []NamedResult {
	if g.ex.excluded(rel) || len(data) == 0 {
		return nil
	}
	return gradeDocument(rel, data)
}

// GradeSubject grades one repository-relative path and returns its findings. It returns
// nil for a subject that is excluded, empty, or clean.
func (g *Grader) GradeSubject(rel string) []string {
	if g.ex.excluded(rel) {
		return nil
	}
	return renderStripped(g.gradeSubjectResults(rel))
}

// GradeNamed grades a caller-selected list of repository-relative paths through the same
// per-subject grader the whole-tree walk composes, so the prose rule keeps one source. A
// path absent from the composed tree is skipped: a commit that deletes a named file grades
// what it commits, not what it once named. A symbolic link is not followed and is not
// graded, matching the whole-tree walk's own rule for a linked directory.
func GradeNamed(root string, rels []string) []string {
	return renderStripped(GradeNamedResults(root, rels))
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
		return gradeDocument(rel, classification.Data)
	case bounds.StateWrongType:
		return []NamedResult{{diagnostic: fmt.Sprintf("prose: %q: refused subject: %s", rel, classification.Reason)}}
	default:
		return []NamedResult{{diagnostic: fmt.Sprintf("prose: %q: refused unreadable subject: %s", rel, classification.Reason)}}
	}
}

// gradeDocument applies the one parser to a document's bytes and attaches the offending
// sentence text. Every entry point reaches the rule through here, so the working-file
// grade and the staged grade cannot drift apart.
func gradeDocument(rel string, data []byte) []NamedResult {
	lines := strings.Split(string(data), "\n")
	var out []NamedResult
	for _, finding := range Findings(string(data)) {
		result := NamedResult{Path: rel, Line: finding.Line, Rule: finding.Kind, Count: finding.Count, Starts: finding.Starts}
		if finding.Kind == KindSentence && finding.Line <= len(lines) {
			result.Sentence = strings.TrimSpace(lines[finding.Line-1])
		}
		out = append(out, result)
	}
	return out
}
