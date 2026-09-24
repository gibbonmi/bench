package otelrecord

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// The retention tests drive the writer through its unexported constructor with a small
// segment limit and a small retained count, then read the segment files back off disk.

// retentionLine is one record line of about 400 bytes that names its index.
func retentionLine(index int) []byte {
	return []byte(fmt.Sprintf(`{"resourceSpans":[],"line":%d,"pad":%q}`, index, strings.Repeat("x", 360)))
}

// appendRetentionLines appends count lines, numbered from first, through writer.
func appendRetentionLines(t *testing.T, writer *Writer, first, count int) {
	t.Helper()
	for index := first; index < first+count; index++ {
		if err := writer.Append(retentionLine(index)); err != nil {
			t.Fatalf("append line %d: %v", index, err)
		}
	}
}

// segmentLines returns the lines of one segment file below root's record directory.
func segmentLines(t *testing.T, home, root, name string) []string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(Dir(home, root), name))
	if err != nil {
		t.Fatalf("read segment %s: %v", name, err)
	}
	return strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n")
}

// sealedNames returns the sealed segment names below root's record directory, sorted.
func sealedNames(t *testing.T, home, root string) []string {
	t.Helper()
	sequences, err := sealedSequences(Dir(home, root))
	if err != nil {
		t.Fatalf("list the record directory: %v", err)
	}
	var names []string
	for _, sequence := range sequences {
		names = append(names, sealedName(sequence))
	}
	return names
}

// Each segment of a 1 KiB limit holds two 400-byte lines, so every second append after
// the first two seals the live segment.
const testSegmentLimit = 1024

// TestAnAppendPastTheLimitSealsTheLiveSegment holds row LE13: a writer without rotation
// keeps one growing file, so the sealed-segment read reds.
func TestAnAppendPastTheLimitSealsTheLiveSegment(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	appendRetentionLines(t, newWriter(home, root, testSegmentLimit, 8), 1, 3)

	if got := sealedNames(t, home, root); !slices.Equal(got, []string{sealedName(1)}) {
		t.Fatalf("sealed segments = %v, want only %s", got, sealedName(1))
	}
	if sealed := segmentLines(t, home, root, sealedName(1)); len(sealed) != 2 {
		t.Fatalf("the sealed segment holds %d lines, want the two it sealed", len(sealed))
	}
	if live := segmentLines(t, home, root, recordFile); len(live) != 1 || !strings.Contains(live[0], `"line":3`) {
		t.Fatalf("the live segment = %v, want only line 3", live)
	}
}

// TestTwoRotationsSealConsecutiveSequences holds row LE92: a name built from a coarse
// clock stamp repeats in one stamp unit, so the rename overwrites the first segment and
// the line count reds.
func TestTwoRotationsSealConsecutiveSequences(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	appendRetentionLines(t, newWriter(home, root, testSegmentLimit, 8), 1, 5)

	if got := sealedNames(t, home, root); !slices.Equal(got, []string{sealedName(1), sealedName(2)}) {
		t.Fatalf("sealed segments = %v, want sequences 1 and 2", got)
	}
	first := segmentLines(t, home, root, sealedName(1))
	if len(first) != 2 || !strings.Contains(first[0], `"line":1`) || !strings.Contains(first[1], `"line":2`) {
		t.Fatalf("the first sealed segment = %v, want lines 1 and 2", first)
	}
}

// TestARotationSkipsAPlantedSequenceName holds row LE94: a rotation that ignores the
// present names renames over the planted file, so the byte comparison reds.
func TestARotationSkipsAPlantedSequenceName(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	if err := os.MkdirAll(Dir(home, root), 0o700); err != nil {
		t.Fatal(err)
	}
	planted := filepath.Join(Dir(home, root), sealedName(1))
	if err := os.WriteFile(planted, []byte("PLANTED\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	appendRetentionLines(t, newWriter(home, root, testSegmentLimit, 8), 1, 3)

	if raw, err := os.ReadFile(planted); err != nil || string(raw) != "PLANTED\n" {
		t.Fatalf("the planted file = %q, %v, want it unchanged", raw, err)
	}
	if sealed := segmentLines(t, home, root, sealedName(2)); len(sealed) != 2 {
		t.Fatalf("the rotation sealed %d lines under sequence 2, want 2", len(sealed))
	}
}

// TestThePruneKeepsTheRetainedCount holds row LE14: a writer that never prunes keeps
// every segment, so the count reds.
func TestThePruneKeepsTheRetainedCount(t *testing.T) {
	const retained = 3
	home, root := t.TempDir(), t.TempDir()
	// Two lines fill the first segment, and each further pair of lines rotates once.
	appendRetentionLines(t, newWriter(home, root, testSegmentLimit, retained), 1, 2+2*(retained+2))

	want := []string{sealedName(3), sealedName(4), sealedName(5)}
	if got := sealedNames(t, home, root); !slices.Equal(got, want) {
		t.Fatalf("sealed segments = %v, want exactly %v", got, want)
	}
}

// TestThePruneRemovesTheLowestSequenceWhateverItsTime holds row LE93: a prune by
// modification time removes a newer sequence, so the survivor set reds.
func TestThePruneRemovesTheLowestSequenceWhateverItsTime(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	writer := newWriter(home, root, testSegmentLimit, 2)
	appendRetentionLines(t, writer, 1, 6)
	future := time.Now().Add(time.Hour)
	if err := os.Chtimes(filepath.Join(Dir(home, root), sealedName(1)), future, future); err != nil {
		t.Fatal(err)
	}
	appendRetentionLines(t, writer, 7, 2)

	want := []string{sealedName(2), sealedName(3)}
	if got := sealedNames(t, home, root); !slices.Equal(got, want) {
		t.Fatalf("sealed segments = %v, want the two highest sequences %v", got, want)
	}
}

// TestTwoWritersAcrossRotationsLeaveOnlyWholeLines holds row LE18: a rotation that copies
// or truncates the live segment loses or splits a line, so the parse or the count reds.
func TestTwoWritersAcrossRotationsLeaveOnlyWholeLines(t *testing.T) {
	const perWriter = 200
	home, root := t.TempDir(), t.TempDir()

	var group sync.WaitGroup
	for writer := 0; writer < 2; writer++ {
		group.Add(1)
		go func(writer int) {
			defer group.Done()
			appender := newWriter(home, root, 4*testSegmentLimit, 1000)
			for line := 0; line < perWriter; line++ {
				if err := appender.Append(retentionLine(writer*perWriter + line)); err != nil {
					t.Errorf("append: %v", err)
					return
				}
			}
		}(writer)
	}
	group.Wait()

	names := append(sealedNames(t, home, root), recordFile)
	if len(names) < 3 {
		t.Fatalf("the writers left segments %v, want forced rotations", names)
	}
	total := 0
	for _, name := range names {
		for _, line := range segmentLines(t, home, root, name) {
			var decoded map[string]any
			if err := json.Unmarshal([]byte(line), &decoded); err != nil {
				t.Fatalf("a line in %s does not parse: %v", name, err)
			}
			total++
		}
	}
	if total != 2*perWriter {
		t.Fatalf("the segments hold %d lines, want the %d appended", total, 2*perWriter)
	}
}

// TestAHeldRotationLockSkipsTheRotation holds row LE19: a blocking lock wait hangs past
// the deadline, and a lock-free rename seals the segment, so the absence read reds.
func TestAHeldRotationLockSkipsTheRotation(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	if err := os.MkdirAll(Dir(home, root), 0o700); err != nil {
		t.Fatal(err)
	}
	holder, err := os.OpenFile(filepath.Join(Dir(home, root), rotationLockFile), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer holder.Close()
	if err := syscall.Flock(int(holder.Fd()), syscall.LOCK_EX); err != nil {
		t.Fatalf("hold the rotation lock: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		appendRetentionLines(t, newWriter(home, root, testSegmentLimit, 8), 1, 3)
	}()
	window := bounds.TestDeadline(0)
	select {
	case <-done:
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("the appends under a held rotation lock", window))
	}

	if got := sealedNames(t, home, root); len(got) != 0 {
		t.Fatalf("sealed segments = %v under a held lock, want none", got)
	}
	if live := segmentLines(t, home, root, recordFile); len(live) != 3 {
		t.Fatalf("the live segment holds %d lines, want all 3", len(live))
	}
}

// TestARotationRefusesTheLastSequence holds row LE106: a rotation whose next sequence
// wraps to 0 seals under sequence 0 and renames over it on the next rotation, so the
// sealed-name read reds.
func TestARotationRefusesTheLastSequence(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	if err := os.MkdirAll(Dir(home, root), 0o700); err != nil {
		t.Fatal(err)
	}
	last := filepath.Join(Dir(home, root), sealedName(math.MaxUint64))
	if err := os.WriteFile(last, []byte("PLANTED\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	writer := newWriter(home, root, testSegmentLimit, 8)
	var refused error
	for index := 1; index <= 3; index++ {
		if err := writer.Append(retentionLine(index)); err != nil {
			refused = err
		}
	}

	if refused == nil {
		t.Fatal("the append past the limit reported no refused rotation")
	}
	if raw, err := os.ReadFile(last); err != nil || string(raw) != "PLANTED\n" {
		t.Fatalf("the planted last sequence = %q, %v, want it unchanged", raw, err)
	}
	if got := sealedNames(t, home, root); !slices.Equal(got, []string{sealedName(math.MaxUint64)}) {
		t.Fatalf("sealed segments = %v, want only the planted last sequence", got)
	}
	if live := segmentLines(t, home, root, recordFile); len(live) != 3 {
		t.Fatalf("the live segment holds %d lines, want all 3", len(live))
	}
}

// retainedMemoryNames returns the memory file names below root's record directory.
func retainedMemoryNames(t *testing.T, home, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(MemoryDir(home, root))
	if err != nil {
		t.Fatalf("list the memory directory: %v", err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

// TestTheMemoryStoreKeepsTheRetainedCount holds row LE57: a store that never prunes keeps
// every file, so the count reds.
func TestTheMemoryStoreKeepsTheRetainedCount(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	start := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)
	for index := range 4 {
		if _, err := retainMemory(home, root, fmt.Sprintf("trace%d", index), []byte("notes\n"), start.Add(time.Duration(index)*time.Second), 3); err != nil {
			t.Fatalf("retain memory %d: %v", index, err)
		}
	}
	names := retainedMemoryNames(t, home, root)
	if len(names) != 3 || slices.ContainsFunc(names, func(name string) bool { return strings.Contains(name, "trace0") }) {
		t.Fatalf("memory files = %v, want the newest 3", names)
	}
}

// TestTheMemoryStoreRefusesASymlinkedDirectory holds row LE58: a store that follows the
// link writes outside the Bench home, so the refusal read reds.
func TestTheMemoryStoreRefusesASymlinkedDirectory(t *testing.T) {
	home, root, outside := t.TempDir(), t.TempDir(), t.TempDir()
	if err := os.MkdirAll(Dir(home, root), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, MemoryDir(home, root)); err != nil {
		t.Fatal(err)
	}
	if _, err := retainMemory(home, root, "trace", []byte("notes\n"), time.Now(), 3); err == nil {
		t.Fatal("the memory write through a symlinked directory returned no error")
	}
	if entries, _ := os.ReadDir(outside); len(entries) != 0 {
		t.Fatalf("the memory write left %d files outside the home", len(entries))
	}
}
