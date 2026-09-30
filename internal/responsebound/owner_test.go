package responsebound

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
)

// The projection counts below are authored apart from the owner: 4 head lines, one spill
// line, and 5 tail lines for any response over 10 lines. The spec fixes these counts, and
// an expectation derived from the owner's constants would move with a changed constant and
// stay green. The spec does not fix the store directory names, so the tests read them from
// the owner's constants.

// tagged is one write of a canned response, to stdout or to stderr.
type tagged struct {
	stderr bool
	data   string
}

// privateHome points the Bench home at a private directory, so no owner test writes the
// fallback home.
func privateHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	return home
}

// outsideRepository is the root of a process outside any repository, so the spill store
// is the `primary` scope below the `none` repository key.
func outsideRepository() string { return "" }

// assertHomeEmpty checks that a response wrote nothing under the Bench home.
func assertHomeEmpty(t *testing.T, home string) {
	t.Helper()
	if entries, err := os.ReadDir(home); err != nil || len(entries) != 0 {
		t.Fatalf("Bench home entries = %v (%v), want none", entries, err)
	}
}

// respond sends writes through a new owner and answers the joined printed bytes. Both
// streams share one sink, so the printed order is the order the owner wrote.
func respond(t *testing.T, writes []tagged) string {
	t.Helper()
	var sink bytes.Buffer
	return respondWith(t, New(benchhome.Dir(), &sink, &sink, outsideRepository, false), &sink, writes)
}

func respondWith(t *testing.T, owner *Owner, sink *bytes.Buffer, writes []tagged) string {
	t.Helper()
	for _, write := range writes {
		stream := owner.Stdout()
		if write.stderr {
			stream = owner.Stderr()
		}
		if n, err := io.WriteString(stream, write.data); err != nil || n != len(write.data) {
			t.Fatalf("write %q = (%d, %v), want every byte accepted", write.data, n, err)
		}
	}
	owner.Finish()
	return sink.String()
}

func stdoutLines(lines ...string) []tagged {
	writes := make([]tagged, len(lines))
	for i, line := range lines {
		writes[i] = tagged{data: line}
	}
	return writes
}

// numbered answers lines first to last, each "line NN\n".
func numbered(first, last int) []string {
	var lines []string
	for i := first; i <= last; i++ {
		lines = append(lines, fmt.Sprintf("line %02d\n", i))
	}
	return lines
}

func joined(writes []tagged) string {
	var all strings.Builder
	for _, write := range writes {
		all.WriteString(write.data)
	}
	return all.String()
}

// onlySpill answers the absolute path of the one spill file in the store.
func onlySpill(t *testing.T, home string) string {
	t.Helper()
	entries, err := os.ReadDir(storeDir(home))
	if err != nil || len(entries) != 1 {
		t.Fatalf("spill store entries = %v (%v), want exactly one file", entries, err)
	}
	return filepath.Join(storeDir(home), entries[0].Name())
}

func readSpill(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// BO3: an 11-line response prints lines 1 to 4, the spill line, and lines 7 to 11.
func TestOwnerProjectsHeadAndTail(t *testing.T) {
	privateHome(t)
	got := strings.SplitAfter(respond(t, stdoutLines(numbered(1, 11)...)), "\n")
	if len(got) != 11 || got[10] != "" {
		t.Fatalf("response lines = %q, want 10 terminated lines", got)
	}
	want := append(append(numbered(1, 4), "spilled{"), numbered(7, 11)...)
	for i := range 10 {
		if i == 4 {
			if !strings.HasPrefix(got[i], want[i]) {
				t.Errorf("line 5 = %q, want the spill line", got[i])
			}
			continue
		}
		if got[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i+1, got[i], want[i])
		}
	}
}

// BO4: a 25-line response of 300 bytes names its totals, its omitted lines, and the
// absolute spill path on its fifth line.
func TestOwnerSpillLineFields(t *testing.T) {
	home := privateHome(t)
	var lines []string
	for i := 1; i <= 25; i++ {
		lines = append(lines, fmt.Sprintf("row %02d fill\n", i))
	}
	writes := stdoutLines(lines...)
	if size := len(joined(writes)); size != 300 {
		t.Fatalf("fixture bytes = %d, want 300", size)
	}
	got := strings.Split(respond(t, writes), "\n")
	path := onlySpill(t, home)
	if !filepath.IsAbs(path) {
		t.Fatalf("spill path %q is not absolute", path)
	}
	want := "spilled{lines=25,bytes=300,omitted_lines=16,cut_lines=0,path=" + path + "}"
	if len(got) < 5 || got[4] != want {
		t.Fatalf("response = %q, want fifth line %q", got, want)
	}
}

// BO5: the spill file holds the complete input bytes.
func TestOwnerSpillHoldsCompleteOutput(t *testing.T) {
	home := privateHome(t)
	writes := stdoutLines(numbered(1, 25)...)
	respond(t, writes)
	if got, want := readSpill(t, onlySpill(t, home)), joined(writes); got != want {
		t.Fatalf("spill file = %q, want the complete output %q", got, want)
	}
}

// BO6: alternate stdout and stderr writes reach the spill file in write order.
func TestOwnerKeepsArrivalOrder(t *testing.T) {
	home := privateHome(t)
	var writes []tagged
	for i := 1; i <= 12; i++ {
		writes = append(writes, tagged{stderr: i%2 == 0, data: fmt.Sprintf("stream %02d\n", i)})
	}
	var stdout, stderr bytes.Buffer
	owner := New(benchhome.Dir(), &stdout, &stderr, outsideRepository, false)
	respondWith(t, owner, &stdout, writes)
	if got, want := readSpill(t, onlySpill(t, home)), joined(writes); got != want {
		t.Fatalf("spill file = %q, want arrival order %q", got, want)
	}
}

// A response within the bound replays each write to its own stream.
func TestOwnerReplaysWithinBound(t *testing.T) {
	home := privateHome(t)
	var stdout, stderr bytes.Buffer
	owner := New(benchhome.Dir(), &stdout, &stderr, outsideRepository, false)
	writes := []tagged{{data: "out 1\n"}, {stderr: true, data: "err 1\n"}, {data: "out 2\n"}}
	respondWith(t, owner, &stdout, writes)
	if stdout.String() != "out 1\nout 2\n" || stderr.String() != "err 1\n" {
		t.Fatalf("streams = (%q, %q), want each write on its own stream", stdout.String(), stderr.String())
	}
	assertHomeEmpty(t, home)
}

// faultAfter accepts limit bytes into the real file, then refuses every later byte.
type faultAfter struct {
	file  io.WriteCloser
	limit int
}

var errInjected = errors.New("injected fault")

func (f *faultAfter) Write(p []byte) (int, error) {
	if len(p) <= f.limit {
		f.limit -= len(p)
		return f.file.Write(p)
	}
	n, err := f.file.Write(p[:f.limit])
	f.limit = 0
	if err != nil {
		return n, err
	}
	return n, errInjected
}

func (f *faultAfter) Close() error { return f.file.Close() }

// BO19: a write fault after 100 spill bytes leaves a 100-byte file, and the response
// holds every later byte after the spill-failed line.
func TestOwnerMidwayFailureKeepsEveryByte(t *testing.T) {
	home := privateHome(t)
	var sink bytes.Buffer
	owner := New(benchhome.Dir(), &sink, &sink, outsideRepository, false)
	owner.create = func(path string) (io.WriteCloser, error) {
		file, err := exclusiveCreate(path)
		if err != nil {
			return nil, err
		}
		return &faultAfter{file: file, limit: 100}, nil
	}
	writes := stdoutLines(numbered(1, 25)...)
	output := joined(writes)
	got := respondWith(t, owner, &sink, writes)
	path := onlySpill(t, home)
	if file := readSpill(t, path); file != output[:100] {
		t.Fatalf("spill file = %q, want the first 100 bytes %q", file, output[:100])
	}
	head := strings.Join(numbered(1, 4), "")
	rest, found := strings.CutPrefix(got, head)
	if !found {
		t.Fatalf("response = %q, want the head %q first", got, head)
	}
	line, rest, _ := strings.Cut(rest, "\n")
	prefix := "spill-failed{path=" + path + ",written_bytes=100,reason="
	if !strings.HasPrefix(line, prefix) || !strings.HasSuffix(line, "}") {
		t.Fatalf("failure line = %q, want prefix %q", line, prefix)
	}
	if rest != output[100:] {
		t.Fatalf("response after the failure line = %q, want every unwritten byte %q", rest, output[100:])
	}
}

// TT61: when the spill file cannot open, the owner passes the output through, so it prints
// no lead after that output. A regular file in place of the Bench home makes the spill
// store unwritable for every user.
func TestOwnerPrintsNoRowWhenSpillCannotOpen(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	var sink bytes.Buffer
	owner := New(filepath.Join(blocker, "home"), &sink, &sink, outsideRepository, false)
	owner.Lead("tree[1]{target,head,dirty}:\n  primary,none,false\n")
	writes := stdoutLines(numbered(1, 25)...)
	got := respondWith(t, owner, &sink, writes)
	if !strings.HasPrefix(got, joined(writes)+"spill-failed{reason=") || strings.Contains(got, "tree[") {
		t.Fatalf("response = %q, want the verb output, the spill-failed line, and no row", got)
	}
}

// BO22: an empty response prints nothing and creates no spill store.
func TestOwnerEmptyResponse(t *testing.T) {
	home := privateHome(t)
	if got := respond(t, nil); got != "" {
		t.Fatalf("empty response printed %q", got)
	}
	assertHomeEmpty(t, home)
}

// BO23: a last line without a newline counts as the eleventh line, and the response and
// the spill file keep its exact final bytes.
func TestOwnerCountsUnterminatedLine(t *testing.T) {
	home := privateHome(t)
	writes := append(stdoutLines(numbered(1, 10)...), tagged{data: "line 11 unterminated"})
	got := respond(t, writes)
	path := onlySpill(t, home)
	if file := readSpill(t, path); file != joined(writes) {
		t.Fatalf("spill file = %q, want %q", file, joined(writes))
	}
	if !strings.HasSuffix(got, strings.Join(numbered(7, 10), "")+"line 11 unterminated") {
		t.Fatalf("response = %q, want lines 7 to 10 and the exact unterminated line", got)
	}
}

// BO24: a NUL byte and invalid UTF-8 reach the spill file byte for byte.
func TestOwnerCopiesBinaryBytes(t *testing.T) {
	home := privateHome(t)
	var writes []tagged
	for i := 1; i <= 12; i++ {
		writes = append(writes, tagged{data: fmt.Sprintf("bin \x00\xff\xfe\xc3 %02d\n", i)})
	}
	respond(t, writes)
	if got, want := readSpill(t, onlySpill(t, home)), joined(writes); got != want {
		t.Fatalf("spill file = %q, want %q", got, want)
	}
}

// TestOwnerSerializesConcurrentWrites drives both streams from two goroutines, as os/exec
// copies a child's two pipes. The race phase runs it, so an unserialized owner reds there,
// and every line must reach the spill file whole.
func TestOwnerSerializesConcurrentWrites(t *testing.T) {
	home := privateHome(t)
	var sink bytes.Buffer
	owner := New(benchhome.Dir(), &sink, &sink, outsideRepository, false)
	var wait sync.WaitGroup
	for _, stream := range []struct {
		name   string
		writer io.Writer
	}{{"out", owner.Stdout()}, {"err", owner.Stderr()}} {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for i := range 200 {
				fmt.Fprintf(stream.writer, "%s %03d\n", stream.name, i)
			}
		}()
	}
	wait.Wait()
	owner.Finish()
	lines := strings.Split(strings.TrimSuffix(readSpill(t, onlySpill(t, home)), "\n"), "\n")
	if len(lines) != 400 {
		t.Fatalf("spill file holds %d lines, want 400", len(lines))
	}
	for _, line := range lines {
		if len(line) != len("out 000") {
			t.Fatalf("spill line %q is torn", line)
		}
	}
}
