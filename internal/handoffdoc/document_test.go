package handoffdoc

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/prose"
)

// twoSectionDocument is the fixture the round-trip and prose rows share: a header,
// main, one request section carrying all six pins with a repeated spec pair, and a
// Shape body. It is built through the exported types, so a render the parser cannot
// read back fails here rather than in a consumer.
func twoSectionDocument() *Document {
	return &Document{
		Header: "# Session handoff\n\nRepository: `bench`\nMain HEAD: `8252957`\nGate: `green`",
		Sections: []Section{
			{
				Key:   MainKey,
				Next:  "`bench status`",
				State: "The primary checkout is clean.\n\nNo phase is live from here.",
			},
			{
				Key: "9f2ab77",
				Fields: []Field{
					{LabelRequestToken, "`spec-handoff-sections-20260904`"},
					{LabelLabel, "`handoff-sections`"},
					{LabelWorktreeTip, "`c0ffee1`"},
					{LabelRecordedBase, "`8252957`"},
					{LabelSpec, "`specs/handoff-sections/spec.md`"},
					{LabelSpecStatus, "`staged`"},
					{LabelSpec, "`specs/binary-freshness/spec.md`"},
					{LabelSpecStatus, "`implemented`"},
				},
				Next:  "`bench worktree exec handoff-sections -- bench gate`",
				State: "The leaf package is in the tree.\n\nThe verb still owns the whole file.",
			},
		},
		Shape: "The document holds one section per live assignment.",
	}
}

// TestRenderedDocumentParsesBackToTheSameSections is the HS1 round-trip row. A
// render the parser reads back differently would let one phase close rewrite a
// sibling's pins under the sibling's own heading.
func TestRenderedDocumentParsesBackToTheSameSections(t *testing.T) {
	want := twoSectionDocument()
	rendered := want.Render()

	got, err := Parse("capture/session-handoff.md", []byte(rendered))
	if err != nil {
		t.Fatalf("parse the rendered document: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip = %#v, want %#v", got, want)
	}
	if second := got.Render(); second != rendered {
		t.Fatalf("second render = %q, want the first render %q", second, rendered)
	}
}

// TestRenderedDocumentPassesTheProseLane is the HS1 prose row. A label line that
// carries a sentence terminator reads as prose, and the lane then grades the pins
// as paragraphs the reviewer wrote.
func TestRenderedDocumentPassesTheProseLane(t *testing.T) {
	rendered := twoSectionDocument().Render()
	if findings := prose.Findings(rendered); len(findings) != 0 {
		t.Fatalf("prose findings on the rendered document = %v, want none", findings)
	}
	for _, label := range []string{LabelRequestToken, LabelWorktreeTip, LabelNextCommand} {
		line := label + ":"
		if !strings.Contains(rendered, line) {
			t.Fatalf("rendered document carries no %q line", line)
		}
	}
}

// TestParseRefusesAnAmbiguousSectionHeading walks the three ambiguity states the
// edge inventory names. Each refusal names the file and the line, because the
// reviewer fixes the document by hand.
func TestParseRefusesAnAmbiguousSectionHeading(t *testing.T) {
	body := "Request token: `t`\n\n" + StateHeading + "\n\nWork is live.\n"
	cases := []struct {
		name    string
		content string
		line    int
		want    string
	}{
		{
			name:    "unknown heading",
			content: MainHeading + "\n\n" + body + "\n## Notes\n\n" + body,
			line:    9,
			want:    "is not",
		},
		{
			name:    "repeated key",
			content: MainHeading + "\n\n" + body + "\n" + RequestHeadingPrefix + "abc\n\n" + body + "\n" + RequestHeadingPrefix + "abc\n\n" + body,
			line:    17,
			want:    "declared twice",
		},
		{
			name:    "no main",
			content: RequestHeadingPrefix + "abc\n\n" + body,
			line:    8,
			want:    "carries no",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse("capture/session-handoff.md", []byte(tc.content))
			var refusal *ParseError
			if err == nil {
				t.Fatalf("parse %s: want a refusal, got none", tc.name)
			}
			refusal, ok := err.(*ParseError)
			if !ok {
				t.Fatalf("parse %s: error %v is not a *ParseError", tc.name, err)
			}
			if refusal.Path != "capture/session-handoff.md" || refusal.Line != tc.line {
				t.Fatalf("refusal names %s:%d, want capture/session-handoff.md:%d", refusal.Path, refusal.Line, tc.line)
			}
			if !strings.Contains(refusal.Reason, tc.want) {
				t.Fatalf("refusal reason %q does not carry %q", refusal.Reason, tc.want)
			}
		})
	}
}

// TestParseKeepsAFencedHeadingInsideState proves a heading the reviewer fenced in
// their own State stays State. Ending a section at it is the one silent data-loss
// path a splitter of this shape has.
func TestParseKeepsAFencedHeadingInsideState(t *testing.T) {
	state := "The verb renders this:\n\n```\n## main\n\n### State\n```"
	want := &Document{Sections: []Section{{Key: MainKey, State: state}}}

	got, err := Parse("handoff.md", []byte(want.Render()))
	if err != nil {
		t.Fatalf("parse the fenced document: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("fenced round trip = %#v, want %#v", got, want)
	}
}

// HS28, story 28. A fence the file never closes makes every heading below it prose, so
// the parser would otherwise read one section holding the whole document and the next
// writer would append a duplicate of the section it could not see. The refusal names the
// file and the line that opened the fence, because that line is the one to close.
func TestParseRefusesAFenceLeftOpenAtEndOfFile(t *testing.T) {
	content := "# Session handoff\n\n## main\n\nNext command: `bench status`\n\n### State\n\n```console\n$ bench gate\n\n## request 9f2ab77\n\nNext command: `bench gate`\n\n### State\n\nThe worktree is live.\n"

	_, err := Parse("capture/session-handoff.md", []byte(content))
	if err == nil {
		t.Fatal("Parse accepted a document whose fence is still open at the end of the file")
	}
	refusal, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("error %v is not a *ParseError", err)
	}
	if refusal.Path != "capture/session-handoff.md" || refusal.Line != 9 {
		t.Fatalf("refusal names %s:%d, want capture/session-handoff.md:9", refusal.Path, refusal.Line)
	}
	if !strings.Contains(refusal.Reason, "```console") {
		t.Fatalf("refusal reason %q does not print the line that opened the fence", refusal.Reason)
	}
}

// The two State shapes ValidateState refuses, and the two it must not. The fenced heading
// row is the one that separates the rule from a bare `## ` prefix test: a writer pasting
// the rendered document into their own State is the case the parser already honours.
func TestValidateStateRefusesTheShapesThatBreakTheGrammar(t *testing.T) {
	for _, tc := range []struct {
		name   string
		state  string
		reason string
		line   int
		text   string
	}{
		{name: "a closed fence passes", state: "Resume here:\n\n```\n## main\n```\n\nThat is all."},
		{name: "an indented heading is prose", state: "  ## Open questions"},
		{name: "an empty State passes"},
		{
			name:   "an unterminated fence",
			state:  "Resume here:\n\n```console\n$ bench gate",
			reason: StateFaultOpenFence, line: 3, text: "```console",
		},
		{
			name:   "an unfenced level-two heading",
			state:  "The build is live.\n\n## Open questions\n\nNone.",
			reason: StateFaultHeading, line: 3, text: "## Open questions",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateState(tc.state)
			if tc.reason == "" {
				if err != nil {
					t.Fatalf("ValidateState refused a State the parser accepts: %v", err)
				}
				return
			}
			fault, ok := err.(*StateError)
			if !ok {
				t.Fatalf("ValidateState = %v, want a *StateError naming %q", err, tc.reason)
			}
			if fault.Reason != tc.reason || fault.Line != tc.line || fault.Text != tc.text {
				t.Fatalf("fault = %+v, want reason %q at line %d, %q", *fault, tc.reason, tc.line, tc.text)
			}
		})
	}
}

// TestEnsureMainAndRemoveKeepTheFallbackSection covers the retirement path's two
// document-level moves: the removal drops one key, and main survives the last one.
func TestEnsureMainAndRemoveKeepTheFallbackSection(t *testing.T) {
	doc := twoSectionDocument()
	if !doc.Remove("9f2ab77") {
		t.Fatal("Remove reported no section under the request key")
	}
	if _, found := doc.Section("9f2ab77"); found {
		t.Fatal("the removed section is still in the document")
	}

	doc.Remove(MainKey)
	doc.EnsureMain()
	main, found := doc.Section(MainKey)
	if !found || len(doc.Sections) != 1 {
		t.Fatalf("after EnsureMain the document holds %d section(s), want main alone", len(doc.Sections))
	}
	if main.State != "" || main.Next != "" {
		t.Fatalf("EnsureMain wrote %#v, want an empty section", main)
	}
}
