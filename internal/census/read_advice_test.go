package census

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReaderCensusAttribution is ER42. An eligible read under a verified assignment appends
// one record in the existing format with the resolved head. A read that also names a pool
// path keeps the pool-path attribution and still appends one record. An unowned read, a
// transformed read, and a malformed assignment id append none.
func TestReaderCensusAttribution(t *testing.T) {
	t.Parallel()
	pooled := func(pool string) string { return "sed -n 1p " + filepath.Join(pool, ownerID+"-"+unknownID, "x") }
	cases := []struct {
		name       string
		command    func(pool string) string
		assignment string
		want       map[string]string
		wantErr    bool
	}{
		{name: "owned eligible read", command: func(string) string { return "cat notes.md" }, assignment: knownID, want: map[string]string{knownID: "cat"}},
		{name: "owned read behind a routine prefix", command: func(string) string { return "timeout 5 head -n 3 notes.md" }, assignment: knownID, want: map[string]string{knownID: "head"}},
		{name: "owned read that names a pool path", command: pooled, assignment: knownID, want: map[string]string{unknownID: "sed"}},
		{name: "unowned read that names a pool path", command: pooled, want: map[string]string{unknownID: "sed"}},
		{name: "unowned eligible read", command: func(string) string { return "cat notes.md" }},
		{name: "owned transformed read", command: func(string) string { return "cat -n notes.md" }, assignment: knownID},
		{name: "malformed assignment id", command: func(string) string { return "cat notes.md" }, assignment: "x", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home, root, pool := fixtureHome(t)
			err := RecordRead(tc.command(pool), root, home, tc.assignment, fixedTime)
			if (err != nil) != tc.wantErr {
				t.Fatalf("RecordRead error = %v, want error %t", err, tc.wantErr)
			}
			entries, _ := os.ReadDir(Dir(home, root))
			if len(entries) != len(tc.want) {
				t.Fatalf("record files = %v, want %d", entries, len(tc.want))
			}
			for id, head := range tc.want {
				lines := records(t, home, root, id)
				if len(lines) != 1 {
					t.Fatalf("records under %s = %v, want one", id, lines)
				}
				if got, ok := parseRecord(lines[0]); !ok || got != head {
					t.Fatalf("record head = %q, want %q", got, head)
				}
			}
		})
	}
}
