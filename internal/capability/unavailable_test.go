package capability

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runSkip drives call in a subtest over the file transport. call is the skip entry
// point under test, and it must end the subtest's goroutine. runSkip returns the raw
// skip log, the emitting subtest's name, and whether that subtest reported a skip.
func runSkip(t *testing.T, call func(t *testing.T)) (log, name string, skipped bool) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "skip.log")
	t.Setenv(LogEnv, path)
	t.Run("skips", func(t *testing.T) {
		name = t.Name()
		t.Cleanup(func() { skipped = t.Skipped() })
		call(t)
		t.Fatal("unreachable: the skip entry point must end this goroutine")
	})
	got, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("read skip log: %v", err)
	}
	return string(got), name, skipped
}

// runUnavailable drives the exported Unavailable entry point through runSkip.
func runUnavailable(t *testing.T, class Class, detail string) (log, name string, skipped bool) {
	t.Helper()
	return runSkip(t, func(t *testing.T) { Unavailable(t, class, detail) })
}

// parseOne reads the one skip record that the log must hold.
func parseOne(t *testing.T, log string) Skip {
	t.Helper()
	lines := strings.SplitAfter(log, "\n")
	if len(lines) != 2 || lines[1] != "" {
		t.Fatalf("skip log = %q, want exactly one line", log)
	}
	skip, ok := ParseLine(lines[0])
	if !ok {
		t.Fatalf("ParseLine(%q) refused the Unavailable record", lines[0])
	}
	return skip
}

// TestUnavailableReason pins the one reason shape for each migrated class: the class
// token, the fixed prefix, and the caller's operation detail verbatim. A
// caller-specific prefix, a prefix fixed to one class token, or a prefix that drops
// the word capability turns this test red.
func TestUnavailableReason(t *testing.T) {
	for _, tc := range []struct {
		class  Class
		detail string
		want   string
	}{
		{Fifo, "mkfifo a: operation not permitted", "fifo capability unavailable: mkfifo a: operation not permitted"},
		{Symlink, "symlink a b: operation not permitted", "symlink capability unavailable: symlink a b: operation not permitted"},
	} {
		t.Run(string(tc.class), func(t *testing.T) {
			log, _, _ := runUnavailable(t, tc.class, tc.detail)
			if got := parseOne(t, log).Reason; got != tc.want {
				t.Fatalf("Unavailable(%s) reason = %q, want %q", tc.class, got, tc.want)
			}
		})
	}
}

// TestUnavailableClassification proves the record keeps the caller's class and the
// capability kind for both migrated families. An environment skip, or a class fixed
// to one family, changes the parsed record of at least one case.
func TestUnavailableClassification(t *testing.T) {
	for _, class := range []Class{Fifo, Symlink} {
		t.Run(string(class), func(t *testing.T) {
			log, _, _ := runUnavailable(t, class, "operation not permitted")
			got := parseOne(t, log)
			if got.Kind != KindCapability || got.Class != class {
				t.Fatalf("Unavailable(%s) record kind=%q class=%q, want kind=%q class=%q", class, got.Kind, got.Class, KindCapability, class)
			}
		})
	}
}

// TestUnavailableWritesBeforeSkip proves the record lands before the skip ends the
// test. A skip that fires before the write leaves the transport empty.
func TestUnavailableWritesBeforeSkip(t *testing.T) {
	log, name, skipped := runUnavailable(t, Fifo, "operation not permitted")
	if !skipped {
		t.Fatal("subtest did not report a skip")
	}
	if log == "" {
		t.Fatal("skip log is empty, so the skip ended the test before the write")
	}
	if got := parseOne(t, log).Name; got != name {
		t.Fatalf("Unavailable record name = %q, want the emitting subtest %q", got, name)
	}
}
