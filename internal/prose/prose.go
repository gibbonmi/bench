// Package prose grades the mechanics of authored Markdown against two bounds: a
// sentence of at most 25 words and a paragraph of at most 6 sentences. It owns the
// document parser, the tree walk, the exclusion grammar, and the subject
// classification, so one caller passes a root and receives diagnostics.
//
// The package is fail-closed. Every state that stops a grade — a refused subject, an
// unterminated delimiter, or a broken exclusion file — returns its own diagnostic
// instead of silence. A diagnostic names the path and the line and the counts. It never
// quotes the document, so a control byte in the text cannot reach the gate output, and
// it renders the path with %q, so one diagnostic stays one line.
package prose

import "fmt"

// MaxSentenceWords and MaxParagraphSentences are the two bounds the rule file states.
// A count over a bound reds; a count at a bound passes.
const (
	MaxSentenceWords      = 25
	MaxParagraphSentences = 6
)

// FindingKind names which bound or which broken delimiter one finding reports.
type FindingKind string

const (
	KindSentence    FindingKind = "sentence"
	KindParagraph   FindingKind = "paragraph"
	KindFence       FindingKind = "fence"
	KindComment     FindingKind = "comment"
	KindFrontmatter FindingKind = "frontmatter"
)

// Finding is one graded fault in one document. Line is the 1-based physical line of the
// sentence, of the first line of the paragraph, or of the opening delimiter. Count is
// the observed word or sentence count, and it is zero for a delimiter fault.
type Finding struct {
	Kind  FindingKind
	Line  int
	Count int
}

// Render states one finding as the diagnostic the gate prints for the subject at rel.
// Each kind carries its own wording, so a reader tells the states apart without the
// document text.
func Render(rel string, f Finding) string {
	switch f.Kind {
	case KindSentence:
		return fmt.Sprintf("prose: %q line %d: sentence of %d words is over the %d-word bound", rel, f.Line, f.Count, MaxSentenceWords)
	case KindParagraph:
		return fmt.Sprintf("prose: %q line %d: paragraph of %d sentences is over the %d-sentence bound", rel, f.Line, f.Count, MaxParagraphSentences)
	case KindFence:
		return fmt.Sprintf("prose: %q line %d: unterminated fenced code block opens here", rel, f.Line)
	case KindComment:
		return fmt.Sprintf("prose: %q line %d: unterminated HTML comment opens here", rel, f.Line)
	case KindFrontmatter:
		return fmt.Sprintf("prose: %q line %d: unterminated frontmatter block opens here", rel, f.Line)
	}
	return fmt.Sprintf("prose: %q line %d: unknown finding %q", rel, f.Line, string(f.Kind))
}
