package responsebound

import (
	"bytes"
	"fmt"
	"math"
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

// BO65: a long head line prints only whole runes within the cut, and the spill line counts
// the cut line. The cases put 1, 2, and 3 kept bytes of a rune across the cut, end a rune
// at the cut, and end the kept bytes with an incomplete rune start that a later ASCII byte
// makes ill-formed. A line of exactly the cut prints whole and is not a cut line.
func TestOwnerCutsLongLine(t *testing.T) {
	const wide, wider = "語", "🙂"
	past := strings.Repeat("x", lineCut)
	pad := func(n int) string { return strings.Repeat("x", n) }
	cases := []struct {
		name     string
		line     string
		printed  int
		cutLines int
	}{
		{"one kept byte of a 3-byte rune", strings.Repeat(wide, 333), lineCut - lineCut%len(wide), 1},
		{"two kept bytes of a 3-byte rune", pad(lineCut-2) + wide + past, lineCut - 2, 1},
		{"three kept bytes of a 4-byte rune", pad(lineCut-3) + wider + past, lineCut - 3, 1},
		{"a 3-byte rune that ends at the cut", pad(lineCut-3) + wide + past, lineCut, 1},
		{"an incomplete rune start before an ASCII byte", pad(lineCut-2) + "\xF0\x9F" + past, lineCut - 2, 1},
		{"exactly the cut", pad(lineCut), lineCut, 0},
		{"one byte past the cut", pad(lineCut + 1), lineCut, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := privateHome(t)
			writes := append(stdoutLines(tc.line+"\n"), stdoutLines(numbered(2, 11)...)...)
			got := respond(t, writes)
			path := onlySpill(t, home)
			spill := fmt.Sprintf("spilled{lines=11,bytes=%d,omitted_lines=2,cut_lines=%d,path=%s}\n", len(joined(writes)), tc.cutLines, path)
			want := tc.line[:tc.printed] + "\n" + strings.Join(numbered(2, 4), "") + spill + strings.Join(numbered(7, 11), "")
			if got != want {
				t.Fatalf("response = %q, want %q", got, want)
			}
			if !utf8.ValidString(got) {
				t.Fatalf("response %q is not valid UTF-8", got)
			}
		})
	}
}

// heapBytes answers the live heap after two full collections. The second one frees the
// objects that the first one moved from a sync.Pool to its victim cache.
func heapBytes() int64 {
	runtime.GC()
	runtime.GC()
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return int64(stats.HeapAlloc)
}

// BO76: a spill retains at most the byte value plus one cut line for each head and tail
// line. One 64 MiB write with no newline catches an owner that keeps the current line.
// Ten long lines, each in two writes, fill every head and tail line, so they catch a line
// buffer that grows past the cut.
func TestOwnerRetainsBoundedLongLine(t *testing.T) {
	const size = 64 << 20
	long := bytes.Repeat([]byte("x"), size)
	start, end := bytes.Repeat([]byte("x"), lineCut-1), []byte("xx\n")
	var tenLines [][]byte
	for range 10 {
		tenLines = append(tenLines, start, end)
	}
	cases := []struct {
		name   string
		writes [][]byte
		want   func(path string) string
	}{
		{"one 64 MiB line", [][]byte{long}, func(path string) string {
			return cutLine() + fmt.Sprintf("spilled{lines=1,bytes=%d,omitted_lines=0,cut_lines=1,path=%s}\n", size, path)
		}},
		{"ten long lines", tenLines, func(path string) string {
			spill := fmt.Sprintf("spilled{lines=10,bytes=%d,omitted_lines=1,cut_lines=9,path=%s}\n", len(bytes.Join(tenLines, nil)), path)
			return strings.Repeat(cutLine(), 4) + spill + strings.Repeat(cutLine(), 5)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// The first of three spills pays the process's one-time costs of a spill, such as
			// the random source, and another goroutine can allocate during one measure. So
			// the smallest measure of the three holds the owner's own bytes alone.
			retained, total := int64(math.MaxInt64), 0
			for range 3 {
				home := privateHome(t)
				var sink bytes.Buffer
				owner := New(benchhome.Dir(), &sink, &sink, outsideRepository, false)
				total = 0
				before := heapBytes()
				for _, p := range tc.writes {
					if n, err := owner.Stdout().Write(p); err != nil || n != len(p) {
						t.Fatalf("write = (%d, %v), want every byte accepted", n, err)
					}
					total += len(p)
				}
				retained = min(retained, heapBytes()-before)
				owner.Finish()
				runtime.KeepAlive(tc.writes)
				path := onlySpill(t, home)
				if info, err := os.Stat(path); err != nil || info.Size() != int64(total) {
					t.Fatalf("spill file = (%v, %v), want %d bytes", info, err, total)
				}
				if want := tc.want(path); sink.String() != want {
					t.Fatalf("response = %q, want %q", sink.String(), want)
				}
			}
			if limit := int64(bounds.ResponseBytes + (bounds.ResponseLines-1)*lineCut); retained > limit {
				t.Fatalf("owner retains %d bytes after %d bytes, want at most %d", retained, total, limit)
			}
		})
	}
}
