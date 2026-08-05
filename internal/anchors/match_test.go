package anchors

import (
	"strings"
	"testing"
)

func TestSatisfiedNormalizesByKind(t *testing.T) {
	tests := []struct {
		name   string
		kind   Kind
		text   string
		needle string
		want   bool
	}{
		{"require", Require, "alpha\n  beta", "alpha   beta", true},
		{"forbid", Forbid, "alpha\n  beta", "alpha   beta", false},
		{"require in section", RequireInSection, "Alpha\n  Beta", "alpha   beta", true},
		{"forbid in section", ForbidInSection, "Alpha\n  Beta", "alpha   beta", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Satisfied(test.kind, test.text, test.needle); got != test.want {
				t.Fatalf("Satisfied(%v, %q, %q) = %t, want %t", test.kind, test.text, test.needle, got, test.want)
			}
		})
	}
}

func TestStripHTMLComments(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"no comment", "alpha beta", "alpha beta"},
		{"complete", "alpha <!-- hidden --> beta", "alpha  beta"},
		{"multiple", "a<!-- one -->b<!-- two -->c", "abc"},
		{"multiline", "before<!-- hidden\ntext -->after", "beforeafter"},
		{"unterminated", "before<!-- hidden", "before"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := StripHTMLComments(test.text); got != test.want {
				t.Fatalf("StripHTMLComments(%q) = %q, want %q", test.text, got, test.want)
			}
		})
	}
}

// TestMarkdownH2SectionsSkipsFencedHeadings records that section-scoped anchors
// resolve past quoted templates instead of treating their headings as boundaries.
func TestMarkdownH2SectionsSkipsFencedHeadings(t *testing.T) {
	const doc = "# Doc\n" +
		"\n" +
		"## Write one file per ticket\n" +
		"\n" +
		"```markdown\n" +
		"## What to build\n" +
		"\n" +
		"## Acceptance\n" +
		"```\n" +
		"\n" +
		"below the fence\n" +
		"\n" +
		"## Draft the breakdown\n" +
		"\n" +
		"a later section\n"

	body, count := MarkdownH2Sections(doc, "Write one file per ticket")
	if count != 1 {
		t.Fatalf("occurrence count = %d, want 1", count)
	}
	if !strings.Contains(body, "below the fence") {
		t.Fatalf("body stops above the prose after the fence:\n%s", body)
	}
	if strings.Contains(body, "a later section") {
		t.Fatalf("body runs past the next real heading:\n%s", body)
	}

	if _, count := MarkdownH2Sections(doc, "Acceptance"); count != 0 {
		t.Fatalf("fenced heading occurrence count = %d, want 0", count)
	}

	later, count := MarkdownH2Sections(doc, "Draft the breakdown")
	if count != 1 {
		t.Fatalf("later section occurrence count = %d, want 1", count)
	}
	if !strings.Contains(later, "a later section") {
		t.Fatalf("later section body = %q", later)
	}

	unclosed := "## One\n\n```\n## Two\n\nstill fenced\n"
	body, count = MarkdownH2Sections(unclosed, "One")
	if count != 1 {
		t.Fatalf("unclosed fence occurrence count = %d, want 1", count)
	}
	if !strings.Contains(body, "still fenced") {
		t.Fatalf("unclosed fence body stops before end of text:\n%s", body)
	}
}
