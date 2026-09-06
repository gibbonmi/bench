package handoffdoc

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// legacyDocumentPath holds a handoff written before the document carried one
// section per assignment. It is a real file rather than a built fixture, because
// the row this covers is the first run of the verb over a repo's existing handoff.
const legacyDocumentPath = "testdata/legacy-session-handoff.md"

// legacyDocument reads the fixture. A read failure is the test's own fault, so it
// fails the test rather than the row.
func legacyDocument(t *testing.T) []byte {
	t.Helper()
	content, err := os.ReadFile(legacyDocumentPath)
	if err != nil {
		t.Fatalf("read the legacy fixture: %v", err)
	}
	return content
}

// TestLegacyDocumentReadsAsMainAndRendersBack is the HS7 conversion row. Without
// the conversion the first close in a repo exits on the file it was asked to
// rewrite, and the reviewer's State is what it refuses over.
func TestLegacyDocumentReadsAsMainAndRendersBack(t *testing.T) {
	doc, err := Parse(DocumentPath, legacyDocument(t))
	if err != nil {
		t.Fatalf("parse the legacy document: %v", err)
	}
	if len(doc.Sections) != 1 || doc.Sections[0].Key != MainKey {
		t.Fatalf("the legacy document read as %d section(s), want main alone", len(doc.Sections))
	}
	main := doc.Sections[0]
	if len(main.Fields) != 0 {
		t.Fatalf("the converted main carries pins %#v, want none", main.Fields)
	}
	if !strings.Contains(main.State, "### Closed decisions") {
		t.Fatalf("the converted State carries no %q heading:\n%s", "### Closed decisions", main.State)
	}
	if !strings.HasPrefix(main.State, "`/bench-implement-spec") || !strings.Contains(main.State, "The user owns the push of `main`.") {
		t.Fatalf("the converted State dropped the legacy State body:\n%s", main.State)
	}
	if !strings.HasPrefix(doc.Shape, "Rewritten in full") {
		t.Fatalf("Shape = %q, want the legacy Shape body", doc.Shape)
	}
	if !strings.HasPrefix(doc.Header, "# Session handoff") {
		t.Fatalf("Header = %q, want the legacy header block", doc.Header)
	}

	rendered := doc.Render()
	again, err := Parse(DocumentPath, []byte(rendered))
	if err != nil {
		t.Fatalf("parse the rendered conversion: %v", err)
	}
	if !reflect.DeepEqual(again, doc) {
		t.Fatalf("the rendered conversion parsed back as %#v, want %#v", again, doc)
	}
}

// TestLegacyNextCommandCarriesOverByteForByte is the HS7 carry-over row. The verb
// preserves a non-empty Next command, so a conversion that reworded or requoted
// the line would hand the next session a command it cannot paste.
func TestLegacyNextCommandCarriesOverByteForByte(t *testing.T) {
	doc, err := Parse(DocumentPath, legacyDocument(t))
	if err != nil {
		t.Fatalf("parse the legacy document: %v", err)
	}
	main, _ := doc.Section(MainKey)
	const want = "`/bench-implement-spec specs/handoff-sections/spec.md --full`"
	if main.Next != want {
		t.Fatalf("the converted Next = %q, want %q", main.Next, want)
	}
}

// TestParseRefusesALegacyStateBesideASection is the HS7 ambiguity row. A document
// that carries both shapes has two candidate homes for the legacy State, and
// guessing one writes over the other.
func TestParseRefusesALegacyStateBesideASection(t *testing.T) {
	section := "\n\n" + labelLine(LabelNextCommand, "`bench status`") + "\n\n" + StateHeading + "\n\nWork is live.\n"
	content := legacyStateHeading + "\n\nThe tree is clean.\n\n" + MainHeading + section

	_, err := Parse(DocumentPath, []byte(content))
	refusal, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("parse the mixed document: error %v is not a *ParseError", err)
	}
	if refusal.Path != DocumentPath || refusal.Line != 1 {
		t.Fatalf("refusal names %s:%d, want %s:1", refusal.Path, refusal.Line, DocumentPath)
	}
	if !strings.Contains(refusal.Reason, "do not mix") {
		t.Fatalf("refusal reason %q does not say the shapes do not mix", refusal.Reason)
	}
}
