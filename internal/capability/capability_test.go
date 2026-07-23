package capability

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestRenderCapabilityLine pins the exact structured line a capability skip writes:
// kind, class, and reason must all be recoverable by a downstream collector that
// prefix-matches bench-skip and reads the leading key=value tokens.
func TestRenderCapabilityLine(t *testing.T) {
	got, err := Render(Skip{Kind: KindCapability, Class: Symlink, Reason: "requires unprivileged symlink support"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	want := "bench-skip kind=capability class=symlink reason=requires unprivileged symlink support\n"
	if got != want {
		t.Fatalf("Render = %q, want %q", got, want)
	}
}

// TestRenderCapabilityLineRejectsUnknownClass proves a class outside the enumerated
// seven is refused rather than formatted — an open vocabulary would let a typo
// silently mint a class the strict count never tallies. Making Class.valid always
// report true is the mutant this pins against: it turns this red.
func TestRenderCapabilityLineRejectsUnknownClass(t *testing.T) {
	got, err := Render(Skip{Kind: KindCapability, Class: Class("network"), Reason: "needs a socket"})
	if err == nil {
		t.Fatalf("Render(%q) = %q, want an error for an unenumerated class", "network", got)
	}
	if got != "" {
		t.Fatalf("Render returned line %q alongside the error; want empty", got)
	}
}

// TestParseLine pins the reader against the writer: every line Render produces must
// come back with its fields intact, and anything else must be refused. A parser that
// accepts a reason-less or unenumerated-class line would let the gate's tally credit
// a class no writer can emit.
func TestParseLine(t *testing.T) {
	for _, want := range []Skip{
		{Kind: KindCapability, Class: Symlink, Reason: "requires unprivileged symlink support"},
		{Kind: KindCapability, Class: Privilege, Reason: "reason with kind=capability class=fifo inside it"},
		{Kind: KindEnvironment, Reason: "subject root has no bin/bench.sh"},
	} {
		line, err := Render(want)
		if err != nil {
			t.Fatalf("Render(%#v): %v", want, err)
		}
		got, ok := ParseLine(line)
		if !ok || got != want {
			t.Fatalf("ParseLine(%q) = %#v, %v; want %#v, true", line, got, ok, want)
		}
	}

	for _, line := range []string{
		"",
		"ok  \tinternal/gate\t0.4s",
		"bench-skip kind=capability class=network reason=needs a socket",
		"bench-skip kind=capability reason=no class token",
		"bench-skip kind=capability class=symlink",
		"bench-skip kind=other reason=unknown kind",
	} {
		if got, ok := ParseLine(line); ok {
			t.Fatalf("ParseLine(%q) = %#v, true; want refused", line, got)
		}
	}
}

// TestCapabilityWritesLineBeforeSkip proves the write lands before t.Skip fires, for
// both transports, by driving the real exported Capability entry point. t.Skip calls
// runtime.Goexit, so a line written after it — or folded only into the skip message
// — would never be delivered; a mutant that reorders Capability's write after its
// t.Skip turns both subtests red because the destination (file or buffer) stays
// empty.
func TestCapabilityWritesLineBeforeSkip(t *testing.T) {
	want := "bench-skip kind=capability class=fifo reason=requires a host fifo\n"

	t.Run("file transport", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "skip.log")
		t.Setenv(LogEnv, path)
		var skipped bool
		t.Run("skips", func(t *testing.T) {
			t.Cleanup(func() { skipped = t.Skipped() })
			Capability(t, Fifo, "requires a host fifo")
			t.Fatal("unreachable: t.Skip must stop this goroutine")
		})
		if !skipped {
			t.Fatal("subtest did not report a skip")
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read skip log: %v", err)
		}
		if string(got) != want {
			t.Fatalf("skip log contents = %q, want %q", got, want)
		}
	})

	t.Run("stdout fallback", func(t *testing.T) {
		t.Setenv(LogEnv, "")
		var buf bytes.Buffer
		prev := stdout
		stdout = &buf
		t.Cleanup(func() { stdout = prev })
		var skipped bool
		t.Run("skips", func(t *testing.T) {
			t.Cleanup(func() { skipped = t.Skipped() })
			Capability(t, Fifo, "requires a host fifo")
			t.Fatal("unreachable: t.Skip must stop this goroutine")
		})
		if !skipped {
			t.Fatal("subtest did not report a skip")
		}
		if got := buf.String(); got != want {
			t.Fatalf("captured skip line = %q, want %q", got, want)
		}
	})
}

// TestEnvironmentWritesLineBeforeSkip mirrors the capability ordering proof for the
// non-capability kind, driving the real exported Environment entry point, since it
// carries the same load-bearing write-then-skip contract through the file
// transport.
func TestEnvironmentWritesLineBeforeSkip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "skip.log")
	t.Setenv(LogEnv, path)
	var skipped bool
	t.Run("skips", func(t *testing.T) {
		t.Cleanup(func() { skipped = t.Skipped() })
		Environment(t, "subject root has no bin/bench.sh")
		t.Fatal("unreachable: t.Skip must stop this goroutine")
	})
	if !skipped {
		t.Fatal("subtest did not report a skip")
	}
	want := "bench-skip kind=environment reason=subject root has no bin/bench.sh\n"
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read skip log: %v", err)
	}
	if string(got) != want {
		t.Fatalf("skip log contents = %q, want %q", got, want)
	}
}

// TestAppendSkipLogConcurrentWriters pins the atomic-single-write requirement: many
// goroutines (standing in for many test binaries sharing one BENCH_SKIP_LOG) append
// at once, and every line must survive intact — none interleaved with another, none
// truncated. Splitting a line across two Write calls reopens exactly this race,
// which is why appendSkipLog issues one Write per line.
func TestAppendSkipLogConcurrentWriters(t *testing.T) {
	path := filepath.Join(t.TempDir(), "skip.log")
	const writers = 40

	var wg sync.WaitGroup
	errs := make(chan error, writers)
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			line := fmt.Sprintf("%s kind=environment reason=writer-%02d\n", linePrefix, i)
			errs <- appendSkipLog(path, line)
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("appendSkipLog: %v", err)
		}
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read skip log: %v", err)
	}
	lines := bytes.Split(bytes.TrimSuffix(got, []byte("\n")), []byte("\n"))
	if len(lines) != writers {
		t.Fatalf("skip log has %d lines, want %d (interleaved or truncated write): %q", len(lines), writers, got)
	}
	seen := make(map[string]bool, writers)
	for _, l := range lines {
		if !bytes.HasPrefix(l, []byte(linePrefix+" kind=environment reason=writer-")) {
			t.Fatalf("malformed line %q in skip log", l)
		}
		if seen[string(l)] {
			t.Fatalf("duplicate line %q in skip log", l)
		}
		seen[string(l)] = true
	}
}
