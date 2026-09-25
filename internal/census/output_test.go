package census

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// outputLines returns the lines of one assignment's output record file.
func outputLines(t *testing.T, home, root, id string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(Dir(home, root), id+".output"))
	if err != nil {
		t.Fatalf("read output records: %v", err)
	}
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
}

// BO57: one bounded response appends one tab-separated line: the time, the head, the
// lines, the bytes, and the disposition. The row fixes the field order.
func TestOutputRecordAppendsLine(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	output := Output{Head: "bench worktree list", Lines: 55, Bytes: 16938, Spilled: true}
	if err := RecordOutput(home, root, knownID, output, fixedTime); err != nil {
		t.Fatal(err)
	}
	lines := outputLines(t, home, root, knownID)
	want := []string{fixedTime.Format(time.RFC3339), "bench worktree list", "55", "16938", "spilled"}
	if len(lines) != 1 || !reflect.DeepEqual(strings.Split(lines[0], "\t"), want) {
		t.Fatalf("output records = %q, want one line with the fields %q", lines, want)
	}
}

// BO59: output records sit apart from the raw-call records, so the raw-call count of the
// assignment does not change.
func TestOutputRecordsLeaveRawCount(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	if err := Record("sed -i x "+filepath.Join(pool, ownerID+"-"+knownID, "x"), root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	before, err := Counts(home, root)
	if err != nil {
		t.Fatal(err)
	}
	for range 3 {
		if err := RecordOutput(home, root, knownID, Output{Head: "bench gate", Lines: 2, Bytes: 40}, fixedTime); err != nil {
			t.Fatal(err)
		}
	}
	after, err := Counts(home, root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("Counts after 3 output records = %v, want %v", after, before)
	}
}

// BO61: the retirement drop removes the raw-call file and the output file.
func TestDropRemovesOutputRecord(t *testing.T) {
	t.Parallel()
	home, root, pool := fixtureHome(t)
	if err := Record("sed -i x "+filepath.Join(pool, ownerID+"-"+knownID, "x"), root, home, fixedTime); err != nil {
		t.Fatal(err)
	}
	if err := RecordOutput(home, root, knownID, Output{Head: "bench gate", Lines: 1, Bytes: 5}, fixedTime); err != nil {
		t.Fatal(err)
	}
	if err := Drop(home, root, knownID); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(Dir(home, root))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("Drop left %v in the census directory", entries)
	}
}

// The breakdown sums each head's calls and bytes, and the heaviest head prints first.
func TestOutputBreakdownSumsEachHead(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	for _, output := range []Output{
		{Head: "bench gate", Lines: 3, Bytes: 90},
		{Head: "bench worktree list", Lines: 10, Bytes: 400, Spilled: true},
		{Head: "bench gate", Lines: 3, Bytes: 90},
	} {
		if err := RecordOutput(home, root, knownID, output, fixedTime); err != nil {
			t.Fatal(err)
		}
	}
	if got, want := OutputBreakdown(home, root, knownID), "bench worktree list=1/400,bench gate=2/180"; got != want {
		t.Fatalf("OutputBreakdown = %q, want %q", got, want)
	}
	if got := OutputBreakdown(home, root, unknownID); got != "" {
		t.Fatalf("OutputBreakdown with no records = %q, want no text", got)
	}
}

// An identifier that is not an assignment id never becomes a path.
func TestRecordOutputRefusesAMalformedAssignment(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	if err := RecordOutput(home, root, "../escape", Output{Head: "bench gate"}, fixedTime); err == nil {
		t.Fatal("RecordOutput accepted a malformed assignment id")
	}
	if _, err := os.Stat(Dir(home, root)); !os.IsNotExist(err) {
		t.Fatalf("a refused record created the census directory: %v", err)
	}
}
