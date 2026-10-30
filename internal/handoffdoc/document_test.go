package handoffdoc

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/prose"
)

// twoSectionDocument is the fixture the round-trip and prose rows share: a header,
// main, one request section carrying every pin with a repeated spec group, and a
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
					{LabelSpecBinding, "`FT395`"},
					{LabelSpec, "`specs/binary-freshness/spec.md`"},
					{LabelSpecStatus, "`implemented`"},
					{LabelSpecBinding, "`unbound`"},
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

// TestOpenFenceNamesTheOpeningLine proves the export a writer calls before its text
// reaches a section answers the same fence rule Parse applies. The reopened case is the
// sharp one: a later pair closes the fence before it, so only the last unmatched opener
// is the line to close.
func TestOpenFenceNamesTheOpeningLine(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
		line int
		open bool
	}{
		{"no fence", "The build is live.\n\n## main", 0, false},
		{"closed fence", "Pasted:\n\n```console\n## main\n```\n\nDone.", 0, false},
		{"open fence", "Pasted:\n\n```console\n## main\n", 3, true},
		{"reopened fence", "```\nfirst\n```\n\n~~~\n## main\n", 5, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			line, open := OpenFence(tc.text)
			if line != tc.line || open != tc.open {
				t.Fatalf("OpenFence(%q) = (%d, %t), want (%d, %t)", tc.text, line, open, tc.line, tc.open)
			}
		})
	}
}

// TestParseKeepsANestedFenceInsideState proves that a closer uses the opener's character and
// a run at least as long, so the inner three-backtick line and the heading under it stay
// State. A toggle on any marker line closes the block there and splits a section at "## X".
func TestParseKeepsANestedFenceInsideState(t *testing.T) {
	state := "Pasted:\n\n````\n```\n## X\n````"
	want := &Document{Sections: []Section{{Key: MainKey, State: state}}}

	got, err := Parse("handoff.md", []byte(want.Render()))
	if err != nil {
		t.Fatalf("nested fence: Parse refused at line %d, want one main section", refusalLine(err))
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("nested fence: round trip gives %d section(s) or a changed State, want one main section with the State unchanged", len(got.Sections))
	}
}

// TestParseReadsAFrontmatterHeadingAsHeader proves that Parse reads the whole file as a
// document. A heading line inside a leading frontmatter block stays header, so it does not
// declare main a second time.
func TestParseReadsAFrontmatterHeadingAsHeader(t *testing.T) {
	want := &Document{Header: "---\n" + MainHeading + "\n---", Sections: []Section{{Key: MainKey}}}

	got, err := Parse("handoff.md", []byte(want.Render()))
	if err != nil {
		t.Fatalf("frontmatter header: Parse refused at line %d, want the header kept", refusalLine(err))
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("frontmatter header: round trip gives %d section(s) or a changed header, want the header and main alone", len(got.Sections))
	}
}

// TestParseSkipsAFencedStateHeading proves that a section reads the reader's fence class.
// A label-shaped fence opener reads as a field, so without that class the State heading
// inside its block would open State.
func TestParseSkipsAFencedStateHeading(t *testing.T) {
	content := MainHeading + "\n\n```Note: x\n" + StateHeading + "\n```\n\n" + StateHeading + "\n"

	_, err := Parse("handoff.md", []byte(content))
	refusal, ok := err.(*ParseError)
	if !ok || refusal.Line != 4 || !strings.Contains(refusal.Reason, "neither a label line") {
		t.Fatalf("fenced State heading: refusal at line %d (a refusal: %t), want the label refusal at line 4", refusalLine(err), ok)
	}
}

// TestParseRefusesAFenceBehindAFrontmatterFault proves that an unclosed first-line marker
// hides no fence. The reader reports only the frontmatter fault there, so Parse reads the
// fence fault from the fragment reading. The marker line alone stays header. A closed
// frontmatter block is no fault, so a fence marker line inside it opens no fence.
func TestParseRefusesAFenceBehindAFrontmatterFault(t *testing.T) {
	body := MainHeading + "\n\nNext command: `bench status`\n\n" + StateHeading + "\n\n"
	for _, tc := range []struct {
		name, content string
		line          int
	}{
		{"open fence", "---\n" + body + "```\n" + RequestHeadingPrefix + "abc\n", 8},
		{"no fence", "---\n" + body + "The worktree is live.\n", 0},
		{"fence marker in closed frontmatter", "---\n```\n---\n" + body + "The worktree is live.\n", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse("handoff.md", []byte(tc.content))
			refusal, _ := err.(*ParseError)
			if refusalLine(err) != tc.line || tc.line > 0 && !strings.Contains(refusal.Reason, "fenced block") {
				t.Fatalf("%s: Parse refusal at line %d, want the fence refusal at line %d (0 is none)", tc.name, refusalLine(err), tc.line)
			}
		})
	}
}

// TestParseRefusesAnUnterminatedComment proves that an open comment, which hides every
// section below it, refuses at the line of its opener, behind an unclosed first-line
// marker too.
func TestParseRefusesAnUnterminatedComment(t *testing.T) {
	body := MainHeading + "\n\nNext command: `bench status`\n\n" + StateHeading + "\n\nThe note <!-- draft\n\n" + RequestHeadingPrefix + "abc\n"
	for _, tc := range []struct {
		name, content string
		line          int
	}{
		{"open comment", body, 7},
		{"behind a frontmatter fault", "---\n" + body, 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse("handoff.md", []byte(tc.content))
			refusal, ok := err.(*ParseError)
			if !ok || refusal.Line != tc.line || !strings.Contains(refusal.Reason, "HTML comment") {
				t.Fatalf("%s: Parse refusal at line %d (a refusal: %t), want the comment refusal at line %d", tc.name, refusalLine(err), ok, tc.line)
			}
		})
	}
}

// TestParseReadsHeadingsOnCommentFreeText proves that a section heading is read on its
// comment-free text. A heading inside a closed comment opens no section and stays State,
// and a heading behind a leading comment opens its section under its own key.
func TestParseReadsHeadingsOnCommentFreeText(t *testing.T) {
	state := "Kept:\n\n<!--\n## Notes\n-->"
	want := &Document{Sections: []Section{{Key: MainKey, State: state}}}

	got, err := Parse("handoff.md", []byte("<!-- lead -->"+MainHeading+"\n\n"+StateHeading+"\n\n"+state+"\n"))
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("commented headings: Parse refused at line %d or changed the sections, want main alone with the commented heading in State", refusalLine(err))
	}
}

// TestParseQuotesTheRawHeading proves that a heading refusal and a converted legacy heading
// keep the line as the file spells it, comment included.
func TestParseQuotesTheRawHeading(t *testing.T) {
	const notes = "## Notes <!-- c -->"
	const request = RequestHeadingPrefix + "a b <!-- c -->"
	main := MainHeading + "\n\nNext command: `bench status`\n\n" + StateHeading + "\n\n"
	for _, tc := range []struct{ name, content, quote string }{
		{"unknown heading", main + notes + "\n", notes},
		{"request digest", main + request + "\n", request},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse("handoff.md", []byte(tc.content))
			refusal, ok := err.(*ParseError)
			if !ok || !strings.HasPrefix(refusal.Reason, "heading "+strconv.Quote(tc.quote)+" ") {
				t.Fatalf("%s: Parse refusal at line %d (a refusal: %t), want the heading refusal quoting the raw line", tc.name, refusalLine(err), ok)
			}
		})
	}

	got, err := Parse("handoff.md", []byte(legacyStateHeading+"\n\nx\n\n"+notes+"\n\ny\n"))
	if err != nil || got.Sections[0].State != "x\n\n### Notes <!-- c -->\n\ny" {
		t.Fatalf("legacy heading: Parse refused at line %d or changed the converted heading, want it with its comment", refusalLine(err))
	}
}

// refusalLine is the line that a Parse refusal names, or 0 for no refusal. A failure message
// names the line, not the refusal text, because that text quotes a source line.
func refusalLine(err error) int {
	if refusal, ok := err.(*ParseError); ok {
		return refusal.Line
	}
	return 0
}

// TestOpenFenceUsesTheCloseRule proves that a State body is a fragment, so a leading
// marker line opens no frontmatter, and only a closer of the opener's character and at
// least its run closes a block. The comment case is the other side of the fault
// kind: the reader reports the comment ahead of the fence, so no fence is reported.
func TestOpenFenceUsesTheCloseRule(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
		line int
		open bool
	}{
		{"other character", "~~~\n```", 1, true},
		{"shorter run", "````\n```", 1, true},
		{"longer run closes", "````\n```\n````", 0, false},
		{"leading marker line", "---\n```\nx\n---", 2, true},
		{"open comment", "```\n<!--\nx", 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			line, open := OpenFence(tc.text)
			if line != tc.line || open != tc.open {
				t.Fatalf("%s: OpenFence = (%d, %t), want (%d, %t)", tc.name, line, open, tc.line, tc.open)
			}
		})
	}

	state := "---\nbefore \t\n````\n```\ninside\n````\nafter\n---"
	var got []string
	for line := range UnfencedLines(state) {
		got = append(got, line)
	}
	if want := []string{"---", "before", "after", "---"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("nested block: UnfencedLines yields %d line(s) or other lines, want the %d unfenced lines", len(got), len(want))
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
