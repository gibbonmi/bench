// Package anchors provides the shared matching mechanism for conformance anchors.
// It stays below the conformance import edge so other anchor consumers can use the
// same semantics without importing the conformance suite.
package anchors

import "strconv"

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
	// ForbidCaseFoldedEmphasis forbids case-folded text after ordinary emphasis normalization.
	ForbidCaseFoldedEmphasis
	// RequireInStep requires text inside one numbered step of an H2 section. The kind
	// implies section scope: the resolution narrows to the section and then to the step.
	RequireInStep
)

// Satisfied reports whether text satisfies kind's presence requirement after
// applying its normalization rules.
func Satisfied(kind Kind, text, needle string) bool {
	textRunes := []rune(text)
	needleRunes := []rune(needle)
	textRunes, _ = normalizeMatchMapped(kind, textRunes, identityOrigin(len(textRunes)))
	needleRunes, _ = normalizeMatchMapped(kind, needleRunes, identityOrigin(len(needleRunes)))
	found := indexRunes(textRunes, needleRunes) >= 0
	switch kind {
	case Require, RequireInSection, RequireInStep:
		return found
	case Forbid, ForbidInSection, ForbidCaseFoldedEmphasis:
		return !found
	default:
		return false
	}
}

func (kind Kind) sectionScoped() bool {
	return kind == RequireInSection || kind == ForbidInSection || kind == RequireInStep
}

// stepScoped reports whether kind narrows its section body to one numbered step. The kind
// is the one owner of that decision: the evaluator and the locator both ask this predicate,
// so neither can narrow a subject the other reads whole.
func (kind Kind) stepScoped() bool {
	return kind == RequireInStep
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

// stepOpener answers the step number a line opens, the number a reader sees. An opener
// carries every one of its decimal digits at column zero, then a period, then a space or a
// tab. An indented line opens no step, so a continuation line belongs to the step above it.
// A leading zero reads as the number without it, because a markdown reader sees `06.` and
// `6.` as the same step.
func stepOpener(line []rune) (int, bool) {
	digits := 0
	for digits < len(line) && line[digits] >= '0' && line[digits] <= '9' {
		digits++
	}
	if digits == 0 || digits+1 >= len(line) || line[digits] != '.' || line[digits+1] != ' ' && line[digits+1] != '\t' {
		return 0, false
	}
	number, err := strconv.Atoi(string(line[:digits]))
	return number, err == nil
}

// MarkdownNumberedSteps returns the first body of the numbered step inside text and the
// number of lines that open that step. A step opener sits at column zero outside a fenced
// block, and its body runs to the next opener or to the end of text. It projects
// stepRunesMapped, the package's one step resolution, so the evaluator and Locate read the
// same step body from the same walk.
func MarkdownNumberedSteps(text string, step int) (string, int) {
	runes := []rune(text)
	body, _, count := stepRunesMapped(runes, identityOrigin(len(runes)), step)
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
