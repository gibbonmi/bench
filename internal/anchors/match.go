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
// count. Headings inside backtick fences do not delimit or count as sections.
func MarkdownH2Sections(text, title string) (string, int) {
	lines := strings.Split(text, "\n")
	heading := "## " + title
	count := 0
	start := -1
	end := -1
	fenced := false
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		if strings.TrimSpace(line) == heading {
			count++
			if count == 1 {
				start = i + 1
			}
			continue
		}
		if start >= 0 && end < 0 && strings.HasPrefix(line, "## ") {
			end = i
		}
	}
	if count == 0 {
		return "", 0
	}
	if end < 0 {
		end = len(lines)
	}
	return strings.Join(lines[start:end], "\n"), count
}

// CollapseSpace replaces each whitespace run with one ASCII space.
func CollapseSpace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
