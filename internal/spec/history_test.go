package spec

import "testing"

func TestParseHistoryLog(t *testing.T) {
	raw := []byte("aaa1111\x00aaa\x002026-07-01T10:00:00-04:00\x00spec-retire: foo\n" +
		"bbb2222\x00bbb\x002026-06-01T09:00:00-04:00\x00some other commit\n")
	got := parseHistoryLog(raw, "retire")
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].full != "aaa1111" || got[0].short != "aaa" || got[0].iso != "2026-07-01T10:00:00-04:00" || got[0].subject != "spec-retire: foo" || got[0].kind != "retire" {
		t.Errorf("row 0 = %+v", got[0])
	}
}

func TestParseHistoryLogEmpty(t *testing.T) {
	if got := parseHistoryLog([]byte(""), "retire"); got != nil {
		t.Errorf("empty input: got %v, want nil", got)
	}
}

// TestMergeHistoryDedupes pins story 3: a commit present in both the retire list and
// the delete list (the common case — a `bench spec retire` commit both deletes the
// file and carries the message) appears exactly once, keeping the retire-list's kind
// because retire is passed first.
func TestMergeHistoryDedupes(t *testing.T) {
	retire := []historyEntry{
		{full: "shared111", short: "shared1", iso: "2026-07-01T10:00:00-04:00", kind: "retire", subject: "spec-retire: foo"},
	}
	del := []historyEntry{
		{full: "shared111", short: "shared1", iso: "2026-07-01T10:00:00-04:00", kind: "delete", subject: "spec-retire: foo"},
	}
	got := mergeHistory(retire, del)
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1 (deduped)", len(got))
	}
	if got[0].kind != "retire" {
		t.Errorf("kind = %q, want %q (retire-list wins the collision)", got[0].kind, "retire")
	}
}

// TestMergeHistoryDeleteOnlyKeepsDeleteKind pins story 3's edge: a delete-only commit
// (the file was removed without a `spec-retire:` message) is tagged delete, not
// retire — a degenerate always-retire implementation would mislabel it.
func TestMergeHistoryDeleteOnlyKeepsDeleteKind(t *testing.T) {
	del := []historyEntry{
		{full: "deleteonly", short: "delonly", iso: "2026-05-01T00:00:00Z", kind: "delete", subject: "drop old spec"},
	}
	got := mergeHistory(nil, del)
	if len(got) != 1 || got[0].kind != "delete" {
		t.Fatalf("got %+v, want one delete-tagged row", got)
	}
}

// TestMergeHistorySortsNewestFirst pins the newest-first ordering across a merge of
// both lists, using the full ISO-8601 timestamp (not just the date) so same-day
// commits still order correctly.
func TestMergeHistorySortsNewestFirst(t *testing.T) {
	retire := []historyEntry{
		{full: "old", short: "old", iso: "2026-01-01T00:00:00Z", kind: "retire", subject: "old retire"},
	}
	del := []historyEntry{
		{full: "new", short: "new", iso: "2026-06-01T00:00:00Z", kind: "delete", subject: "new delete"},
		{full: "mid", short: "mid", iso: "2026-03-01T00:00:00Z", kind: "delete", subject: "mid delete"},
	}
	got := mergeHistory(retire, del)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	wantOrder := []string{"new", "mid", "old"}
	for i, w := range wantOrder {
		if got[i].full != w {
			t.Errorf("position %d = %q, want %q (order: %v)", i, got[i].full, w, got)
		}
	}
}

// TestRetireTokenMatches pins the exact retire-token cut layered over the coarse
// --grep contains filter: the token must extend to the end of the subject, so a slug
// that is a string prefix of another slug (`dash` vs `dashboard`) never claims the
// longer slug's retirement commit, and a spaced slug still terminates soundly.
func TestRetireTokenMatches(t *testing.T) {
	cases := []struct {
		subject, slug string
		want          bool
	}{
		{"spec-retire: dash", "dash", true},
		{"spec-retire: dashboard", "dash", false},
		{"spec-retire: dashboard", "dashboard", true},
		{"spec-retire: weird name", "weird name", true},
		{"spec-retire: weird name", "weird", false},
		{"spec-retire: dash plus trailing words", "dash", false},
		{"unrelated subject", "dash", false},
	}
	for _, c := range cases {
		if got := retireTokenMatches(c.subject, c.slug); got != c.want {
			t.Errorf("retireTokenMatches(%q, %q) = %v, want %v", c.subject, c.slug, got, c.want)
		}
	}
}

func TestMergeHistoryEmpty(t *testing.T) {
	if got := mergeHistory(nil, nil); len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}
