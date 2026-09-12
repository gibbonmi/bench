package learnings

import "testing"

// TestFormatEntryNormalizesEmbeddedControlBytes covers the writer-reader contract
// FormatEntry's own doc comment claims: the bench learning verb, the adopt scaffold, and
// Parse cannot disagree about an entry FormatEntry produced. A title, what, right, or
// rule value that carries a raw newline or another control byte previously reached the
// heading or a body bullet unsanitized. An embedded newline in title split the heading
// line in two; the first half then missed the "[open]" state marker Parse requires, so
// the verb's own entry read back as malformed (Command's own "dated learning heading
// must end with [open]" reason) — the exact class `bench status`'s drain row and `bench
// learnings` then disagreed over, one hiding the good count as "unknown", the other
// still listing it.
func TestFormatEntryNormalizesEmbeddedControlBytes(t *testing.T) {
	entry := FormatEntry("2026-09-11", "bad\ntitle", "line one\nline two", "right\tanswer", "ru\x01le")
	journal := JournalSchemaHeading + "\n\n" + entry
	entries, malformed := Parse([]byte(journal))
	if len(malformed) != 0 {
		t.Fatalf("FormatEntry produced a malformed entry: %+v\n%s", malformed, journal)
	}
	if len(entries) != 1 || entries[0].State != "open" {
		t.Fatalf("entries = %+v, want one open entry\n%s", entries, journal)
	}
	want := "bad title"
	if entries[0].Title != want {
		t.Fatalf("title = %q, want %q", entries[0].Title, want)
	}
	wantBody := "- **What happened:** line one line two\n" +
		"- **Right behavior:** right answer\n" +
		"- **Proposed rule change:** ru le"
	if entries[0].Body != wantBody {
		t.Fatalf("body = %q, want %q", entries[0].Body, wantBody)
	}
}
