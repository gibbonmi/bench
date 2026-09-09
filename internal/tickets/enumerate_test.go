package tickets

import (
	"os"
	"path/filepath"
	"testing"
)

// TestTagOfDerivesTheSpecTag pins the one derivation both the sweep and the
// preflight rows read. The degenerate case matters most: a digit-leading row ID
// answers the empty tag, which stands the Covers tag rule down rather than
// grading every declared token as foreign.
func TestTagOfDerivesTheSpecTag(t *testing.T) {
	for _, testCase := range []struct{ rowID, want string }{
		{"TG1", "TG"},
		{"TG12", "TG"},
		{"FT93", "FT"},
		{"1TG", ""},
		{"", ""},
		{"TG", "TG"},
	} {
		if got := TagOf(testCase.rowID); got != testCase.want {
			t.Errorf("TagOf(%q) = %q, want %q", testCase.rowID, got, testCase.want)
		}
	}
}

// TestUnrepresentableValueNamesABlockerControlByte is the representability
// guard's blocker half. A blocker basename reaches a rendered detail cell the
// same way a `Writes:` entry does, so a control byte in either is refused before
// a verdict renders.
func TestUnrepresentableValueNamesABlockerControlByte(t *testing.T) {
	ticket := Ticket{Name: "one.md", Blockers: []string{"two\x01.md"}, Writes: []string{"a.go (new)"}}
	field, value, found := UnrepresentableValue(ticket)
	if !found || field != "Blocked by" || value != "two\x01.md" {
		t.Fatalf("UnrepresentableValue(blocker control byte) = (%q, %q, %v), want the Blocked by entry named", field, value, found)
	}
	if _, _, found := UnrepresentableValue(Ticket{Name: "one.md", Blockers: []string{"two.md"}, Writes: []string{"a.go"}}); found {
		t.Error("UnrepresentableValue(clean ticket) reported a value, want none")
	}
}

func TestEnumerateNoFollowRefusesLinkAndPreservesDefault(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.txt")
	if err := os.WriteFile(target, []byte("ticket bytes"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "one.md")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}

	got, _, refusal := Enumerate(root, entries)
	if refusal != nil || len(got) != 1 || string(got[0].Data) != "ticket bytes" {
		t.Fatalf("Enumerate() = %#v, refusal=%v", got, refusal)
	}
	_, _, refusal = EnumerateNoFollow(root, entries)
	if refusal == nil || refusal.State != "wrong-type" {
		t.Fatalf("EnumerateNoFollow() refusal = %#v, want wrong-type", refusal)
	}
}
