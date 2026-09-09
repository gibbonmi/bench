// Package anchors provides the shared matching mechanism for conformance anchors.
// It stays below the conformance import edge so other anchor consumers can use the
// same semantics without importing the conformance suite.
package anchors

import "strings"

// Kind selects an anchor's normalization and presence requirement.
type Kind uint8

const (
	// Require requires text across an entire file.
	Require Kind = iota
	// Forbid forbids text across an entire file.
	Forbid
	// RequireInSection requires text inside an H2 section.
	RequireInSection
	// ForbidInSection forbids text inside an H2 section.
	ForbidInSection
)

// Satisfied reports whether text satisfies kind's presence requirement after
// applying its normalization rules.
func Satisfied(kind Kind, text, needle string) bool {
	text = CollapseSpace(text)
	needle = CollapseSpace(needle)
	switch kind {
	case Require:
		return strings.Contains(text, needle)
	case Forbid:
		return !strings.Contains(text, needle)
	case RequireInSection:
		text = strings.ToLower(text)
		needle = strings.ToLower(needle)
		return strings.Contains(text, needle)
	case ForbidInSection:
		text = strings.ToLower(text)
		needle = strings.ToLower(needle)
		return !strings.Contains(text, needle)
	default:
		return false
	}
}

// MarkdownH2Section returns the first matching H2 section body.
func MarkdownH2Section(text, title string) string {
	body, _ := MarkdownH2Sections(text, title)
	return body
}

// MarkdownH2Sections returns the first matching H2 section body and occurrence
// count. Headings inside backtick fences do not delimit or count as sections. It
// projects sectionRunesMapped, the package's one section resolution, so the evaluator
// and Locate read the same body from the same walk.
func MarkdownH2Sections(text, title string) (string, int) {
	runes := []rune(text)
	body, _, count := sectionRunesMapped(runes, identityOrigin(len(runes)), title)
	if count == 0 {
		return "", 0
	}
	return string(body), count
}

// CollapseSpace replaces each whitespace run with one ASCII space. It projects
// collapseSpaceMapped, the package's one collapse.
func CollapseSpace(text string) string {
	runes := []rune(text)
	collapsed, _ := collapseSpaceMapped(runes, identityOrigin(len(runes)))
	return string(collapsed)
}
