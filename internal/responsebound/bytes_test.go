package responsebound

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/bounds"
)

// The byte fixtures here read the byte value and the line cut from the registry and the
// owner, so each boundary moves with its one source. The response shape stays authored
// apart, as in owner_test.go.

// filled answers one line of size bytes, its newline included.
func filled(size int) string { return strings.Repeat("x", size-1) + "\n" }

// cutLine is the printed form of a line longer than the cut: its first lineCut bytes and
// its newline.
func cutLine() string { return strings.Repeat("x", lineCut) + "\n" }

// BO63: a 3-line response of 5000 bytes spills on its bytes alone. It prints every line,
// cut, and its printed bytes stay within the byte value.
func TestOwnerAppliesByteBound(t *testing.T) {
	home := privateHome(t)
	writes := stdoutLines(filled(1667), filled(1667), filled(1666))
	output := joined(writes)
	if len(output) != 5000 {
		t.Fatalf("fixture bytes = %d, want 5000", len(output))
	}
	got := respond(t, writes)
	path := onlySpill(t, home)
	want := strings.Repeat(cutLine(), 3) + "spilled{lines=3,bytes=5000,omitted_lines=0,cut_lines=3,path=" + path + "}\n"
	if got != want {
		t.Fatalf("response = %q, want %q", got, want)
	}
	if len(got) > bounds.ResponseBytes {
		t.Fatalf("response holds %d bytes, want at most %d", len(got), bounds.ResponseBytes)
	}
	if file := readSpill(t, path); file != output {
		t.Fatalf("spill file holds %d bytes, want the complete %d", len(file), len(output))
	}
}

// byteValueLines answers ResponseLines lines of exactly the byte value plus extra bytes.
// The last line takes the rest after the others, so it alone passes the cut.
func byteValueLines(extra int) []string {
	lines := make([]string, bounds.ResponseLines-1)
	for i := range lines {
		lines[i] = filled(lineCut)
	}
	return append(lines, filled(bounds.ResponseBytes-len(lines)*lineCut+extra))
}

// BO75: a 10-line response of exactly the byte value prints unchanged and creates no spill
// file. One more byte spills.
func TestOwnerByteBoundary(t *testing.T) {
	home := privateHome(t)
	at := stdoutLines(byteValueLines(0)...)
	if got, want := respond(t, at), joined(at); len(want) != bounds.ResponseBytes || got != want {
		t.Fatalf("response of %d bytes = %q, want it unchanged", len(want), got)
	}
	assertHomeEmpty(t, home)

	home = privateHome(t)
	lines := byteValueLines(1)
	over := stdoutLines(lines...)
	got := respond(t, over)
	path := onlySpill(t, home)
	spill := fmt.Sprintf("spilled{lines=10,bytes=%d,omitted_lines=1,cut_lines=1,path=%s}\n", bounds.ResponseBytes+1, path)
	want := strings.Join(lines[:4], "") + spill + strings.Join(lines[5:9], "") + cutLine()
	if got != want {
		t.Fatalf("response = %q, want %q", got, want)
	}
	if file := readSpill(t, path); file != joined(over) {
		t.Fatalf("spill file holds %d bytes, want the complete %d", len(file), len(joined(over)))
	}
}

// BO65: a head line of 999 bytes of 3-byte runes prints only whole runes within the cut,
// and the spill line counts one cut line.
func TestOwnerCutsLongLine(t *testing.T) {
	home := privateHome(t)
	const wide = "語"
	long := strings.Repeat(wide, 333)
	writes := append(stdoutLines(long+"\n"), stdoutLines(numbered(2, 11)...)...)
	if len(long) != 999 {
		t.Fatalf("fixture line = %d bytes, want 999", len(long))
	}
	got := respond(t, writes)
	path := onlySpill(t, home)
	kept := long[:lineCut-lineCut%len(wide)]
	spill := fmt.Sprintf("spilled{lines=11,bytes=%d,omitted_lines=2,cut_lines=1,path=%s}\n", len(joined(writes)), path)
	want := kept + "\n" + strings.Join(numbered(2, 4), "") + spill + strings.Join(numbered(7, 11), "")
	if got != want {
		t.Fatalf("response = %q, want %q", got, want)
	}
	if !utf8.ValidString(got) {
		t.Fatalf("response %q is not valid UTF-8", got)
	}
}

// heapBytes answers the live heap after a full collection.
func heapBytes() int64 {
	runtime.GC()
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return int64(stats.HeapAlloc)
}

// BO76: one 64 MiB write with no newline spills, and the owner retains at most the byte
// value plus one cut line for each head and tail line. An owner that keeps the current
// line retains the whole write.
func TestOwnerRetainsBoundedLongLine(t *testing.T) {
	home := privateHome(t)
	const size = 64 << 20
	data := bytes.Repeat([]byte("x"), size)
	// A first spill in a separate home pays the process's one-time costs of a spill, such
	// as the random source, so the measure below holds the owner's own bytes alone.
	warm := New(t.TempDir(), io.Discard, io.Discard, outsideRepository, false)
	_, _ = io.WriteString(warm.Stdout(), strings.Join(numbered(1, 11), ""))
	warm.Finish()
	var sink bytes.Buffer
	owner := New(benchhome.Dir(), &sink, &sink, outsideRepository, false)
	before := heapBytes()
	if n, err := owner.Stdout().Write(data); err != nil || n != size {
		t.Fatalf("write = (%d, %v), want every byte accepted", n, err)
	}
	retained := heapBytes() - before
	limit := int64(bounds.ResponseBytes + (bounds.ResponseLines-1)*lineCut)
	if retained > limit {
		t.Fatalf("owner retains %d bytes after one %d-byte line, want at most %d", retained, size, limit)
	}
	owner.Finish()
	runtime.KeepAlive(data)
	path := onlySpill(t, home)
	if info, err := os.Stat(path); err != nil || info.Size() != size {
		t.Fatalf("spill file = (%v, %v), want %d bytes", info, err, size)
	}
	want := cutLine() + fmt.Sprintf("spilled{lines=1,bytes=%d,omitted_lines=0,cut_lines=1,path=%s}\n", size, path)
	if sink.String() != want {
		t.Fatalf("response = %q, want %q", sink.String(), want)
	}
}
