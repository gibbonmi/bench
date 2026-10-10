package learnings

import (
	"errors"
	"strings"
	"testing"
)

const resolvingCommit = "0123456789abcdef0123456789abcdef01234567"

// twoEntryJournal holds two open entries with distinct titles, the shape MarkResolved
// edits in place.
func twoEntryJournal() string {
	return JournalSchemaHeading + "\n\n" + FormatEntry("2026-10-01", "first", "a", "b", "") + "\n" +
		FormatEntry("2026-10-02", "second", "c", "d", "")
}

// TestMarkResolvedLinksOnlyTheNamedEntry covers the link's whole contract: the named
// entry stays open, reads back with its commit, and every other byte of the journal is
// unchanged. Removing the inserted bullet restores the original bytes exactly.
func TestMarkResolvedLinksOnlyTheNamedEntry(t *testing.T) {
	before := twoEntryJournal()
	after, err := MarkResolved([]byte(before), "second", resolvingCommit)
	if err != nil {
		t.Fatalf("MarkResolved: %v", err)
	}
	entries, malformed := Parse(after)
	if len(malformed) != 0 || len(entries) != 2 {
		t.Fatalf("entries=%+v malformed=%+v\n%s", entries, malformed, after)
	}
	if entries[0].ResolvedBy != "" || entries[1].ResolvedBy != resolvingCommit {
		t.Fatalf("ResolvedBy = %q, %q; want empty then %q", entries[0].ResolvedBy, entries[1].ResolvedBy, resolvingCommit)
	}
	if entries[1].State != "open" {
		t.Fatalf("state = %q, want open", entries[1].State)
	}
	if got := strings.Replace(string(after), resolvedBulletPrefix+resolvingCommit+"\n", "", 1); got != before {
		t.Fatalf("journal changed beyond the one bullet:\n%s", after)
	}
}

// TestMarkResolvedRefusesWithoutChange covers each refusal: the content is not rewritten
// and the error names the cause for the caller to render.
func TestMarkResolvedRefusesWithoutChange(t *testing.T) {
	resolved, err := MarkResolved([]byte(twoEntryJournal()), "first", resolvingCommit)
	if err != nil {
		t.Fatal(err)
	}
	duplicate := twoEntryJournal() + "\n" + FormatEntry("2026-10-03", "first", "e", "f", "")
	for _, tc := range []struct {
		name    string
		journal []byte
		title   string
		want    error
	}{
		{"unknown title", []byte(twoEntryJournal()), "third", ErrNoOpenEntry},
		{"title is a prefix only", []byte(twoEntryJournal()), "fir", ErrNoOpenEntry},
		{"two entries share the title", []byte(duplicate), "first", ErrAmbiguousEntry},
		{"entry already linked", resolved, "first", ErrAlreadyResolved},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := MarkResolved(tc.journal, tc.title, resolvingCommit)
			if !errors.Is(err, tc.want) || out != nil {
				t.Fatalf("MarkResolved = %q, %v; want nil and %v", out, err, tc.want)
			}
		})
	}
}

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
