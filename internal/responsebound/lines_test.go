package responsebound

import (
	"errors"
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

// The reserved frames are authored apart from Finish, in the shapes that the spec and the
// owner tests above fix. A reader that misses a frame, or accepts a frame that does not
// close, reds a row here.
func TestProjectionLineClassifiesReservedFrames(t *testing.T) {
	for _, row := range []struct {
		line    string
		want    Projection
		framed  bool
		refused bool
	}{
		{"probe[1]{verdict}:\n", Projection{}, false, false},
		{"  spilled{path=/x}\n", Projection{}, false, false},
		{"spilled{lines=12,bytes=96,omitted_lines=3,cut_lines=0,path=/h/a,b.out}\n", Projection{Spilled, "/h/a,b.out"}, true, false},
		{"spill-failed{path=/h/c.out,written_bytes=100,reason=no space}\n", Projection{SpillFailed, "/h/c.out"}, true, false},
		{"spill-failed{reason=mkdir home: not a directory}", Projection{SpillFailed, ""}, true, false},
		{"spilled{lines=12,bytes=96,path=/h/a.out\n", Projection{Kind: Spilled}, true, true},
		{"spilled{lines=12,bytes=96}\n", Projection{Kind: Spilled}, true, true},
		{"spill-failed{path=/h/c.out}\n", Projection{Kind: SpillFailed, Path: "/h/c.out"}, true, true},
		{"spill-failed{}\n", Projection{Kind: SpillFailed}, true, true},
	} {
		got, framed, err := ProjectionLine(row.line)
		if got != row.want || framed != row.framed || (err != nil) != row.refused {
			t.Errorf("ProjectionLine(%q) = (%+v, %v, %v), want (%+v, %v, refused=%v)", row.line, got, framed, err, row.want, row.framed, row.refused)
		}
	}
}

// RequireComplete refuses a frame on any physical line, and a complete output passes.
func TestRequireCompleteRefusesEveryFrame(t *testing.T) {
	complete := "probe[1]{verdict}:\n  bit\n"
	if err := RequireComplete([]byte(complete)); err != nil {
		t.Fatalf("complete output = %v, want nil", err)
	}
	for _, frame := range []string{"spilled{lines=1,bytes=1,omitted_lines=0,cut_lines=0,path=/p}\n", "spill-failed{reason=r}\n", "spilled{\n"} {
		if err := RequireComplete([]byte(complete + frame + "tail\n")); !errors.Is(err, ErrIncomplete) {
			t.Errorf("output with %q = %v, want ErrIncomplete", frame, err)
		}
	}
}
