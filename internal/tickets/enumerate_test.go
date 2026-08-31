package tickets

import "testing"

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
