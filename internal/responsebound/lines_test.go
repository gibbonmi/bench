package responsebound

import (
	"strings"
	"testing"
)

// assertProjection grades a spilled response against its complete output: the full spill
// line, the printed head and tail lines, and the spill file bytes.
func assertProjection(t *testing.T, home, got, output, spillFields string, head, tail []string) {
	t.Helper()
	path := onlySpill(t, home)
	want := strings.Join(head, "") + "spilled{" + spillFields + ",cut_lines=0,path=" + path + "}\n" + strings.Join(tail, "")
	if got != want {
		t.Fatalf("response = %q, want %q", got, want)
	}
	if file := readSpill(t, path); file != output {
		t.Fatalf("spill file = %q, want %q", file, output)
	}
}

// A line that arrives in two writes counts once, and the head and the tail print it whole.
func TestOwnerJoinsLineSplitAcrossWrites(t *testing.T) {
	home := privateHome(t)
	var writes []tagged
	for _, line := range numbered(1, 12) {
		writes = append(writes, tagged{data: line[:3]}, tagged{data: line[3:]})
	}
	output := joined(writes)
	assertProjection(t, home, respond(t, writes), output, "lines=12,bytes=96,omitted_lines=3", numbered(1, 4), numbered(8, 12))
}

// One write that holds several lines and crosses line 11 counts each line, and the head
// and the tail print each line whole.
func TestOwnerSplitsOneWriteIntoLines(t *testing.T) {
	home := privateHome(t)
	writes := append(stdoutLines(numbered(1, 2)...), tagged{data: strings.Join(numbered(3, 12), "")}, tagged{data: "line 13\n"})
	output := joined(writes)
	assertProjection(t, home, respond(t, writes), output, "lines=13,bytes=104,omitted_lines=4", numbered(1, 4), numbered(9, 13))
}
