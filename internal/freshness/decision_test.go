package freshness

import (
	"strings"
	"testing"
)

func TestSelectFreshnessFromImmutableDigests(t *testing.T) {
	digestA := strings.Repeat("a", 64)
	digestB := strings.Repeat("b", 64)
	accepted := SelectionInput{StoredSource: digestA, CurrentSource: digestA, StoredExecutable: digestB, CurrentExecutable: digestB}
	for _, tc := range []struct {
		name  string
		input SelectionInput
		want  Selection
	}{
		{name: "accepted", input: accepted, want: Selection{Accepted: true}},
		{name: "source refusal", input: SelectionInput{StoredSource: digestA, CurrentSource: digestB, StoredExecutable: digestB, CurrentExecutable: digestB}, want: Selection{Reason: "seal source digest does not match current build inputs"}},
		{name: "executable refusal", input: SelectionInput{StoredSource: digestA, CurrentSource: digestA, StoredExecutable: digestA, CurrentExecutable: digestB}, want: Selection{Reason: "seal executable digest does not match binary contents"}},
		{name: "empty", input: SelectionInput{}, want: Selection{Reason: "stored source digest is malformed"}},
		{name: "hostile", input: SelectionInput{StoredSource: "\x1b", CurrentSource: digestA, StoredExecutable: digestB, CurrentExecutable: digestB}, want: Selection{Reason: "stored source digest is malformed"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			first := Select(tc.input)
			second := Select(tc.input)
			if first != tc.want || second != first {
				t.Fatalf("Select() = %#v then %#v, want %#v", first, second, tc.want)
			}
		})
	}
}
