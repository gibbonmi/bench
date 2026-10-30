package anchors

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/prose"
)

// TestAdmissionClasses walks the live registry and applies the class shape rules of the
// anchor admission. Each table case is one row that a rule must refuse or admit.
func TestAdmissionClasses(t *testing.T) {
	guide, err := os.ReadFile(filepath.Join("..", "..", operatingGuide))
	if err != nil {
		t.Fatal(err)
	}
	for _, anchor := range Entries() {
		if faults := admissionFaults(anchor, string(guide)); len(faults) > 0 {
			t.Errorf("%s needle %q is not admitted: %s", anchor.File, anchor.Needle, strings.Join(faults, "; "))
		}
	}
	// Workflow opens the test guide, so a needle in both sections locates in Workflow.
	const testGuide = "## Workflow\n\nWorkflow text. Shared text.\n\n## Roles\n\nRole text. Shared text.\n"
	words := func(n int) string { return strings.TrimSpace(strings.Repeat("word ", n)) }
	for _, test := range []struct {
		name   string
		anchor Anchor
		admit  bool
	}{
		{"AP01 unclassified", Anchor{File: "a.md", Kind: Require, Needle: "Status: staged"}, false},
		{"forbid needs no class", Anchor{File: "a.md", Kind: Forbid, Needle: "retired form"}, true},
		{"AP12 label with a reader", Anchor{File: "a.md", Kind: Require, Class: Label, Reader: "ticket parser", Needle: "Blocked by:"}, true},
		{"AP12 label with an empty reader", Anchor{File: "a.md", Kind: Require, Class: Label, Needle: "Blocked by:"}, false},
		{"AP13 route at the word limit", Anchor{File: "a.md", Kind: Require, Class: Route, Needle: "craft-gate one two three four five"}, true},
		{"AP13 route of seven words", Anchor{File: "a.md", Kind: Require, Class: Route, Needle: "craft-gate one two three four five six"}, false},
		{"route with no route token", Anchor{File: "a.md", Kind: Require, Class: Route, Needle: "the gate"}, false},
		{"AP14 invariant on another file", Anchor{File: "AGENTS.md", Kind: Require, Class: Invariant, Needle: "Role text."}, false},
		{"invariant located in Roles", Anchor{File: operatingGuide, Kind: Require, Class: Invariant, Needle: "Role text."}, true},
		{"AP15 invariant in the Workflow section", Anchor{File: operatingGuide, Kind: RequireInSection, Class: Invariant, Section: "Workflow", Needle: "Workflow text."}, false},
		{"AP15 invariant located in Workflow", Anchor{File: operatingGuide, Kind: Require, Class: Invariant, Needle: "Workflow text."}, false},
		{"AP15 invariant located in Workflow and also in Roles", Anchor{File: operatingGuide, Kind: Require, Class: Invariant, Needle: "Shared text."}, false},
		{"AP09 route row on a Workflow line", Anchor{File: operatingGuide, Kind: Require, Class: Route, Needle: "Workflow text."}, false},
		{"AP16 command with a flag", Anchor{File: "a.md", Kind: Require, Class: Command, Needle: "Run `bench roadmap --flow` once."}, true},
		{"AP16 command with a placeholder", Anchor{File: "a.md", Kind: Require, Class: Command, Needle: "bench coverage <spec>"}, true},
		{"AP16 command with a quoted value", Anchor{File: "a.md", Kind: Require, Class: Command, Needle: "Run `bench idea \"text\"` once."}, true},
		{"AP16 command with a path", Anchor{File: "a.md", Kind: Require, Class: Command, Needle: "Run `bench lint ./docs` once."}, true},
		{"AP16 command with no bench span", Anchor{File: "a.md", Kind: Require, Class: Command, Needle: "Run the gate --full once."}, false},
		{"AP16 flag after the closing backtick", Anchor{File: "a.md", Kind: Require, Class: Command, Needle: "Run `bench maps` before the --full pass."}, false},
		{"AP17 needle at the bound with a code span", Anchor{File: "a.md", Kind: Require, Class: Command, Needle: "`bench test --check a b c d e f` " + words(prose.MaxSentenceWords-1)}, true},
		{"AP17 needle over the bound", Anchor{File: "a.md", Kind: Require, Class: Command, Needle: "`bench test --check system` " + words(prose.MaxSentenceWords)}, false},
	} {
		if got := len(admissionFaults(test.anchor, testGuide)) == 0; got != test.admit {
			t.Errorf("%s: admitted = %t, want %t", test.name, got, test.admit)
		}
	}
	// Each argument-free verb sentence of the classification snapshot names a verb and its
	// subcommands with no argument, so it is a rule sentence and not a command row.
	for _, needle := range []string{
		"only `bench worktree land` releases it after the reviewed source publishes",
		"It opens with `bench preflight review` in explicit-base mode; a red preflight stops the phase.",
		"a linked repo upgrades with `bench upgrade`",
		"Each fold after the first chunk takes `bench record fold` before the next commit.",
		"The changed-function caller sweep runs `bench consumers` for each changed function, unexported functions included.",
		"Run `bench maps` before declaring readiness",
	} {
		if len(admissionFaults(Anchor{File: "a.md", Kind: Require, Class: Command, Needle: needle}, testGuide)) == 0 {
			t.Errorf("argument-free verb sentence %q is admitted as a command row", needle)
		}
	}
}

// routeMaxWords is the word limit of a ROUTE needle.
const routeMaxWords = 6

// admissionFaults applies the class shape rules to one row. guide is the operating guide
// text, where a row with no section field locates the H2 of its needle.
func admissionFaults(anchor Anchor, guide string) []string {
	if anchor.Kind != Require && anchor.Kind != RequireInSection {
		return nil
	}
	var faults []string
	words := prose.SentenceWords(anchor.Needle)
	if words > prose.MaxSentenceWords {
		faults = append(faults, fmt.Sprintf("needle of %d words is over the %d-word bound", words, prose.MaxSentenceWords))
	}
	if anchor.File == operatingGuide && !inCoreSection(anchor, guide) {
		faults = append(faults, "row guards an operating guide line outside its core sections")
	}
	switch anchor.Class {
	case Label:
		if strings.TrimSpace(anchor.Reader) == "" {
			faults = append(faults, "label row names no reader")
		}
	case Route:
		if words > routeMaxWords || !slices.ContainsFunc([]string{"/bench-", "$bench-", "craft-", "bench-craft-", "bench "}, func(token string) bool { return strings.Contains(anchor.Needle, token) }) {
			faults = append(faults, fmt.Sprintf("route needle is over %d words or names no route token", routeMaxWords))
		}
	case Invariant:
		if anchor.File != operatingGuide {
			faults = append(faults, "invariant row guards a file other than the operating guide")
		}
	case Command:
		if !commandArgument(anchor.Needle) {
			faults = append(faults, "command needle holds no `bench ` span with an argument or a flag")
		}
	default:
		faults = append(faults, "require row declares no class")
	}
	return faults
}

// inCoreSection reports whether a row sits in Roles, The Bench CLI contract, or the
// four-invariants section of the operating guide. A row with no section field uses the
// located H2 of its needle.
func inCoreSection(anchor Anchor, guide string) bool {
	core := []string{"Roles", "The Bench CLI contract", invariantsSection}
	if anchor.Section != "" {
		return slices.Contains(core, anchor.Section)
	}
	line := Locate(Require, "", anchor.Needle, guide)
	return line != 0 && slices.ContainsFunc(core, func(title string) bool {
		return Locate(RequireInSection, title, anchor.Needle, guide) == line
	})
}

// commandArgument is the COMMAND argument token rule. A `bench ` span runs from "bench " to
// the next backtick, or to the end of the needle. The span has an argument when one of its
// whitespace-separated tokens opens with "-", "<", "\"", or ".": a flag, a placeholder, a
// quoted value, or a path. A bare word, such as "land" in `bench worktree land`, names a verb
// or a subcommand and is no argument.
func commandArgument(needle string) bool {
	for rest := needle; strings.Contains(rest, "bench "); {
		_, rest, _ = strings.Cut(rest, "bench ")
		span, _, _ := strings.Cut(rest, "`")
		if slices.ContainsFunc(strings.Fields(span), func(token string) bool { return strings.ContainsAny(token[:1], "-<\".") }) {
			return true
		}
	}
	return false
}

// TestLocateHonorsSectionScope pins DG28 and DG39: a section kind searches its own
// section body only, and a duplicated owning section — which the evaluator's section
// resolution refuses — locates nothing rather than searching the first duplicate.
func TestLocateHonorsSectionScope(t *testing.T) {
	doc := "# Title\n" +
		"outside needle here\n" +
		"\n" +
		"## Alpha\n" +
		"before\n" +
		"needle here\n" +
		"## Beta\n" +
		"after\n"

	if got := Locate(RequireInSection, "Alpha", "needle here", doc); got != 6 {
		t.Fatalf("needle inside its section: Locate = %d, want line 6", got)
	}
	if got := Locate(RequireInSection, "Alpha", "outside needle", doc); got != 0 {
		t.Fatalf("needle only outside its section: Locate = %d, want 0", got)
	}
	// A whole-file search finds the same needle outside the section, so the section
	// scope — not an absent needle — is what zeroes the row above.
	if got := Locate(Require, "", "outside needle", doc); got != 2 {
		t.Fatalf("whole-file search for the outside needle: Locate = %d, want line 2", got)
	}

	duplicated := "## Alpha\nneedle here\n## Alpha\nneedle here\n"
	if got := Locate(RequireInSection, "Alpha", "needle here", duplicated); got != 0 {
		t.Fatalf("duplicated owning section: Locate = %d, want 0", got)
	}

	if got := Locate(RequireInSection, "Missing", "needle here", doc); got != 0 {
		t.Fatalf("absent section: Locate = %d, want 0", got)
	}
}

// TestLocateMapsCollapsedMatchesToLines pins DG29: a needle split across lines, a needle
// after a multi-line comment, and a case fold that changes a rune's byte length all
// locate the match's first character on its real physical line — never shifted by the
// collapse, the strip, or the fold that found it.
func TestLocateMapsCollapsedMatchesToLines(t *testing.T) {
	t.Run("split across lines", func(t *testing.T) {
		doc := "prefix\none\ntwo suffix\n"
		if got := Locate(Require, "", "one two", doc); got != 2 {
			t.Fatalf("Locate = %d, want line 2 (the first character of \"one\")", got)
		}
	})

	t.Run("after a multi-line comment", func(t *testing.T) {
		doc := "before\n<!-- hidden\nacross lines -->needle right after\nend\n"
		if got := Locate(Require, "", "needle right after", doc); got != 3 {
			t.Fatalf("Locate = %d, want line 3", got)
		}
		if got := Locate(Require, "", "hidden", doc); got != 0 {
			t.Fatalf("needle inside the comment: Locate = %d, want 0", got)
		}
	})

	t.Run("non-ASCII uppercase text", func(t *testing.T) {
		// U+212A KELVIN SIGN lowercases to the ASCII letter k: 3 bytes fold to 1. It sits
		// on the line before the match, not inside it, so a byte-offset implementation runs out
		// of its rune budget two bytes early — before the newline that starts the match's
		// own line — and reports the wrong (earlier) line. A rune-mapped one does not.
		doc := "# Doc\n\n## Notes\nKelvin first\nkelvin note\n## End\n"
		if got := Locate(RequireInSection, "Notes", "kelvin note", doc); got != 5 {
			t.Fatalf("Locate = %d, want line 5", got)
		}
	})

	t.Run("absent needle and absent data", func(t *testing.T) {
		if got := Locate(Require, "", "never appears", "one\ntwo\n"); got != 0 {
			t.Fatalf("absent needle: Locate = %d, want 0", got)
		}
		if got := Locate(Require, "", "anything", ""); got != 0 {
			t.Fatalf("absent data: Locate = %d, want 0", got)
		}
	})

	t.Run("forbid kind locates the violation", func(t *testing.T) {
		if got := Locate(Forbid, "", "banned phrase", "one\nbanned phrase here\n"); got != 2 {
			t.Fatalf("Locate = %d, want line 2 (the forbidden text's own line)", got)
		}
		if got := Locate(Forbid, "", "banned phrase", "one\ntwo\n"); got != 0 {
			t.Fatalf("clean forbid text: Locate = %d, want 0", got)
		}
	})
}

func TestLocateMapsCaseFoldedEmphasisToPhysicalLine(t *testing.T) {
	doc := "<!-- An executable **red** is mandatory in a comment. -->\n" +
		"KELVIN appears before the match.\n" +
		"An EXECUTABLE **RED** IS MANDATORY before specification.\n" +
		"An executable *red* is mandatory later.\n"
	if got := Locate(ForbidCaseFoldedEmphasis, "", "executable red is mandatory", doc); got != 3 {
		t.Fatalf("Locate = %d, want line 3", got)
	}
}

// TestLocateStripsRejoinedComments pins DG43: the removal of one comment can join the text
// on either side of it into a second comment, and the evaluator's strip rescans and removes
// that one too. Locate rescans with it, so a needle the evaluator reads as absent locates
// nothing here either. A single-pass strip prints a line for that needle.
func TestLocateStripsRejoinedComments(t *testing.T) {
	const needle = "never by polling"
	rejoined := "before\n<!<!-- z -->-- " + needle + " -->\ntail needle\n"
	if got := Locate(Require, "", needle, rejoined); got != 0 {
		t.Fatalf("needle inside a rejoined comment: Locate = %d, want 0", got)
	}
	if StripHTMLComments(rejoined) != "before\n\ntail needle\n" {
		t.Fatalf("the evaluator's strip left %q, so the fixture does not rejoin", StripHTMLComments(rejoined))
	}
	// The text after the rejoined comment keeps its physical line, so the repeat maps its
	// surviving runes back to data rather than shifting them.
	if got := Locate(Require, "", "tail needle", rejoined); got != 3 {
		t.Fatalf("text after a rejoined comment: Locate = %d, want line 3", got)
	}

	// The same needle as plain text locates its line, so the strip — not an absent needle —
	// is what zeroes the row above. A plain comment around it keeps that behavior.
	plain := "before\n" + needle + " here\n"
	if got := Locate(Require, "", needle, plain); got != 2 {
		t.Fatalf("plain needle: Locate = %d, want line 2", got)
	}
	commented := "before\n<!-- " + needle + " -->\n" + needle + " here\n"
	if got := Locate(Require, "", needle, commented); got != 3 {
		t.Fatalf("needle after a plain comment: Locate = %d, want line 3", got)
	}
}

// TestLocateReadsTheBlockReader pins the block reader's rules under the section walk and the
// comment strip. A fenced or an indented heading line neither opens nor closes a section. The
// walk and the strip read the whole document, so a comment opener or a fence marker in the
// frontmatter hides no section, a marker line in a section body stays body, and a comment in
// the frontmatter stays. The walk reads the file as stored, in Locate and in the evaluator, so
// a removed comment cannot make a frontmatter marker, and a comment in a section still hides its
// needle. Multi-byte text before a section or a comment moves neither the section body nor the
// source line of a match.
func TestLocateReadsTheBlockReader(t *testing.T) {
	dir := t.TempDir()
	evaluatorBody := func(t *testing.T, doc string) string {
		t.Helper()
		path := filepath.Join(dir, "doc.md")
		if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
			t.Fatal(err)
		}
		return resolveSection("doc.md", "S", read(path, "doc.md")).body
	}
	for _, test := range []struct {
		name, doc   string
		count, line int
	}{
		{"a tilde-fenced heading leaves the section open", "## S\n~~~\n## Other\n~~~\nneedle\n## T\n", 1, 5},
		{"a tilde-fenced heading counts as no section", "~~~\n## S\n~~~\n## S\nneedle\n", 1, 5},
		{"an indented heading opens no section", "  ## S\nneedle\n", 0, 0},
		{"an indented heading closes no section", "## S\n  ## T\nneedle\n", 1, 3},
		{"frontmatter that holds a comment opener keeps the section", "---\nk: <!--\n---\n## S\nneedle\n", 1, 5},
		{"frontmatter that holds a tilde marker keeps the section and its marker body line", "---\n~~~\n---\n## S\n---\nneedle\n", 1, 6},
		{"a comment before a marker on the first line opens no frontmatter", "<!-- c -->---\n## S\nneedle\n---\n", 1, 3},
		{"a comment before a later marker closes no frontmatter", "---\n## S\n<!-- c -->---\nneedle\n", 1, 4},
		{"an empty section holds no needle", "## S\n## T\nneedle\n", 1, 0},
		{"a section after multi-byte text keeps its body", "€€€€\n## S\nneedle\n## T\nafter\n", 1, 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, count := MarkdownH2Sections(test.doc, "S"); count != test.count {
				t.Errorf("%s: section count = %d, want %d", test.name, count, test.count)
			}
			if got := strings.Contains(MarkdownH2Section(test.doc, "S"), "needle"); got != (test.line > 0) {
				t.Errorf("%s: section body holds the needle = %t, want %t", test.name, got, test.line > 0)
			}
			if got := Locate(RequireInSection, "S", "needle", test.doc); got != test.line {
				t.Errorf("%s: Locate = %d, want line %d", test.name, got, test.line)
			}
			if got := strings.Contains(evaluatorBody(t, test.doc), "needle"); got != (test.line > 0) {
				t.Errorf("%s: evaluator section body holds the needle = %t, want %t", test.name, got, test.line > 0)
			}
		})
	}
	const frontmatter = "---\n<!-- kept -->\n---\n<!-- hidden -->\n"
	const section = "## S\n<!-- hidden -->\nkept\n"
	for _, test := range []struct {
		name, doc, needle string
		kind              Kind
		line              int
	}{
		{"a match after a comment keeps its line after multi-byte text", "€€€€<!-- x -->\nab\nneedle\n", "needle", Require, 3},
		{"a comment in the frontmatter stays", frontmatter, "kept", Require, 2},
		{"a comment below the frontmatter hides its needle", frontmatter, "hidden", Require, 0},
		{"a comment in a section hides its needle", section, "hidden", RequireInSection, 0},
		{"a match after a comment in a section keeps its line", section, "kept", RequireInSection, 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := Locate(test.kind, "S", test.needle, test.doc); got != test.line {
				t.Errorf("%s: Locate = %d, want line %d", test.name, got, test.line)
			}
			if !test.kind.sectionScoped() {
				return
			}
			if got := strings.Contains(evaluatorBody(t, test.doc), test.needle); got != (test.line > 0) {
				t.Errorf("%s: evaluator section body holds the needle = %t, want %t", test.name, got, test.line > 0)
			}
		})
	}
}
