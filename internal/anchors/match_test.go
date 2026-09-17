package anchors

import (
	"slices"
	"strings"
	"testing"
)

func TestEvaluatePathRejectsMixedNestedEmphasis(t *testing.T) {
	anchor := pathTestAnchor(t, ForbidCaseFoldedEmphasis)
	for _, test := range []struct {
		name string
		text string
		line int
	}{
		{"bold around italic", "An executable **_red_** is mandatory before specification.", 3},
		{"italic around bold", "An executable _**red**_ is mandatory before specification.", 3},
		{"triple around bold", "An executable ***__red__*** is mandatory before specification.", 3},
		{"mapped multiline", "intro\nAn **_EXECUTABLE\nRED_** IS MANDATORY before specification.", 4},
		{"repeated bold", "**An _executable **red** is mandatory_ before specification.**", 3},
		{"repeated italic", "*An _executable *red* is mandatory_ before specification.*", 3},
		{"repeated triple", "***An _executable ***red*** is mandatory_ before specification.***", 3},
		{"repeated underscore bold", "__An *executable __red__ is mandatory* before specification.__", 3},
		{"repeated underscore italic", "_An *executable _red_ is mandatory* before specification._", 3},
		{"repeated underscore triple", "___An *executable ___red___ is mandatory* before specification.___", 3},
		{"shared star two then one", "An ***executable** red is mandatory* before specification.", 3},
		{"shared star one then two", "An ***executable* red is mandatory** before specification.", 3},
		{"shared underscore two then one", "An ___executable__ red is mandatory_ before specification.", 3},
		{"shared underscore one then two", "An ___executable_ red is mandatory__ before specification.", 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			h := anchorHarness{rules: []anchorRule{{file: anchor.File, needle: test.text}}}
			result := EvaluatePath(h.write(t, -1), anchor.File)
			if !slices.Contains(result.Diagnostics, anchor.Diagnostic) {
				t.Errorf("missing %q for %q", anchor.Diagnostic, test.text)
			}
			for _, location := range result.Locations {
				if location.Anchor == anchor {
					if location.Line != test.line {
						t.Errorf("violation line = %d, want %d", location.Line, test.line)
					}
					return
				}
			}
			t.Fatal("registered DG15 location missing")
		})
	}
}

func TestUnpairedEmphasisRunsStayVisible(t *testing.T) {
	for _, marker := range []string{"*", "_", "**", "__", "***", "___"} {
		t.Run(marker, func(t *testing.T) {
			text := strings.Repeat(marker+"a ", 4096)
			if !Satisfied(ForbidCaseFoldedEmphasis, text, "a a") {
				t.Error("unpaired markers disappeared between words")
			}
			if Satisfied(ForbidCaseFoldedEmphasis, text, marker+"a") {
				t.Error("unpaired marker is no longer searchable")
			}
		})
	}
}

// TestStepOpenerReadsEveryLeadingDigitAtColumnZero pins the opener shape the reader sees.
// The multi-digit rows keep `10.` apart from `1.`, and the indented rows keep a continuation
// line inside the step above it. The rows without a period keep a line of digits and spaces,
// such as a table cell or a measurement, out of the opener shape.
func TestStepOpenerReadsEveryLeadingDigitAtColumnZero(t *testing.T) {
	for _, test := range []struct {
		line   string
		number int
		opens  bool
	}{
		{"1. the first step", 1, true},
		{"10. the tenth step", 10, true},
		{"06. a padded sixth step", 6, true},
		{"7.\ta tab after the period", 7, true},
		{"   6. an indented line", 0, false},
		{"\t6. a tab-indented line", 0, false},
		{"6.no space after the period", 0, false},
		{"6 . a space before the period", 0, false},
		{"6  two spaces and no period", 0, false},
		{"6\t\ttwo tabs and no period", 0, false},
		{"step 6. a line that opens with a word", 0, false},
	} {
		t.Run(test.line, func(t *testing.T) {
			number, opens := stepOpener([]rune(test.line))
			if number != test.number || opens != test.opens {
				t.Fatalf("stepOpener(%q) = (%d, %t), want (%d, %t)", test.line, number, opens, test.number, test.opens)
			}
		})
	}
}

// TestMarkdownH2SectionExcludesItsHeading pins the section body's open boundary. The body
// starts under the heading, so a needle that repeats the heading text is located inside the
// section only when the section carries that text in its own body.
func TestMarkdownH2SectionExcludesItsHeading(t *testing.T) {
	const doc = "# title\n\n## Alpha\n\nfirst line\n\n## Beta\n\nsecond line\n"
	if body := MarkdownH2Section(doc, "Alpha"); strings.Contains(body, "## Alpha") || !strings.Contains(body, "first line") {
		t.Fatalf("MarkdownH2Section(Alpha) = %q, want the body under the heading and not the heading line", body)
	}
	if got := Locate(RequireInSection, "Alpha", "## Alpha", doc); got != 0 {
		t.Fatalf("Locate of the heading text inside its own section = %d, want 0", got)
	}
}

// TestMarkdownNumberedStepsIncludesItsOpener pins the step body's open boundary, which the
// section body's boundary inverts. A step's own first words sit on its opener line, so a
// body that started under that line would lose them. The body still stops at the next
// opener, so the two boundaries are graded together.
func TestMarkdownNumberedStepsIncludesItsOpener(t *testing.T) {
	const section = "1. the first step opens here\n   a continuation line\n2. the second step opens here\n   another continuation\n"
	body, count := MarkdownNumberedSteps(section, 1)
	if count != 1 {
		t.Fatalf("MarkdownNumberedSteps(1) counted %d owning openers, want 1", count)
	}
	if !strings.HasPrefix(body, "1. the first step opens here") {
		t.Fatalf("step body = %q, want it to start with its own opener line", body)
	}
	if !Satisfied(RequireInStep, body, "the first step opens here") {
		t.Fatalf("step body = %q, want a needle on the opener line to satisfy the step-scoped kind", body)
	}
	if strings.Contains(body, "the second step opens here") {
		t.Fatalf("step body = %q, want it to stop at the next opener", body)
	}
}

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
		// A step-scoped kind normalizes as a section-scoped kind does, so the case fold
		// reaches the needle inside the step.
		{"require in step", RequireInStep, "Alpha\n  Beta", "alpha   beta", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Satisfied(test.kind, test.text, test.needle); got != test.want {
				t.Fatalf("Satisfied(%v, %q, %q) = %t, want %t", test.kind, test.text, test.needle, got, test.want)
			}
		})
	}
}

func TestForbidCaseFoldedEmphasisMatchesBoundedForms(t *testing.T) {
	const needle = "executable red is mandatory"
	for _, test := range []struct {
		name string
		text string
	}{
		{"bold", "An executable **red** is mandatory before specification."},
		{"asterisk italic", "An executable *red* is mandatory before specification."},
		{"underscore bold", "An executable __red__ is mandatory before specification."},
		{"underscore italic", "An executable _red_ is mandatory before specification."},
		{"nested asterisks", "An executable ***red*** is mandatory before specification."},
		{"uppercase", "An EXECUTABLE **RED** IS MANDATORY before specification."},
		{"inline code stays searchable", "`An executable red is mandatory before specification.`"},
		{"fenced code stays searchable", "```markdown\nAn executable red is mandatory before specification.\n```"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if Satisfied(ForbidCaseFoldedEmphasis, test.text, needle) {
				t.Errorf("Satisfied(ForbidCaseFoldedEmphasis, %q, %q) = true, want false", test.text, needle)
			}
		})
	}
	for _, test := range []struct {
		name string
		text string
	}{
		{"legitimate negative", "An executable red is not mandatory before specification."},
		{"escaped opening", `An executable \**red** is mandatory before specification.`},
		{"escaped closing", `An executable **red\** is mandatory before specification.`},
		{"unpaired", "An executable **red is mandatory before specification."},
		{"intraword underscore", "An executable r_ed_ is mandatory before specification."},
		{"escaped inner opening", `An executable **\_red_** is mandatory before specification.`},
		{"unpaired inner", "An executable **_red** is mandatory before specification."},
	} {
		t.Run(test.name, func(t *testing.T) {
			if !Satisfied(ForbidCaseFoldedEmphasis, test.text, needle) {
				t.Errorf("Satisfied(ForbidCaseFoldedEmphasis, %q, %q) = false, want true", test.text, needle)
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
