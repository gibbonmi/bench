// Package capability owns skipping for the gate's test suite. A bare t.Skip is
// invisible under non-verbose `go test` — a skip and a pass print the same nothing,
// and `go test` without -v discards a package's stdout and stderr entirely — so
// every skip here writes a structured line to the file named by BENCH_SKIP_LOG
// first and only then calls t.Skip. A later gate slice sets that variable to a
// run-scoped path, reads it back once the phase's `go test` invocations exit, and
// tallies skips by kind and class; the line shape here is that consumer's only
// contract. With BENCH_SKIP_LOG unset — a developer's hand-run `go test -v` — the
// line falls back to stdout.
//
// Two kinds cover every skip in the suite: Capability, for a security assertion the
// host cannot run (missing symlink support, no privilege to drop, ...), and
// Environment, for everything else (an absent subject binary, an unset conformance
// root, an unmaterialized fixture).
package capability

import (
	"fmt"
	"io"
	"os"
	"testing"
)

// Class enumerates the capability categories a security test can require from its
// host. The set is closed: a strict release mode will count skips by class, and an
// open vocabulary would let a typo silently mint a class nothing counts. Capability
// rejects any Class outside this list instead of formatting it.
type Class string

const (
	Symlink   Class = "symlink"
	Fifo      Class = "fifo"
	PID       Class = "pid"
	CPU       Class = "cpu"
	Privilege Class = "privilege"
	Signal    Class = "signal"
	Tool      Class = "tool"
)

var validClasses = map[Class]bool{
	Symlink:   true,
	Fifo:      true,
	PID:       true,
	CPU:       true,
	Privilege: true,
	Signal:    true,
	Tool:      true,
}

func (c Class) valid() bool { return validClasses[c] }

// linePrefix opens every line this package writes, so a downstream collector can
// recognize one with a prefix match before parsing the key=value tokens that follow.
const linePrefix = "bench-skip"

// skipLogEnv names the environment variable the gate sets to a run-scoped absolute
// file path. Many test binaries append to that one file concurrently from separate
// processes, so writers must go through appendSkipLog rather than opening it by hand.
const skipLogEnv = "BENCH_SKIP_LOG"

// stdout is the stdout-fallback destination, held behind a var so this package's own
// tests can swap it for a buffer and assert on the emitted bytes without touching the
// process's real stdout.
var stdout io.Writer = os.Stdout

// renderCapabilityLine is the one place a capability skip line is built. reason takes
// the remainder of the line verbatim (no escaping), so the collector need only split
// the leading kind and class tokens off the front. class outside the enumerated set
// is refused rather than formatted.
func renderCapabilityLine(class Class, reason string) (string, error) {
	if !class.valid() {
		return "", fmt.Errorf("capability: unknown class %q", class)
	}
	return fmt.Sprintf("%s kind=capability class=%s reason=%s\n", linePrefix, class, reason), nil
}

// renderEnvironmentLine is the one place an environment skip line is built.
func renderEnvironmentLine(reason string) string {
	return fmt.Sprintf("%s kind=environment reason=%s\n", linePrefix, reason)
}

// appendSkipLog lands one skip line in the shared log with a single Write call. The
// gate runs many `go test` binaries concurrently, each appending to the same
// BENCH_SKIP_LOG path; a write under 4096 bytes made through one Write to an
// O_APPEND-opened file descriptor is atomic on Linux, so one call per line is what
// keeps two writers' lines from interleaving or clobbering each other. Splitting the
// line across multiple Write calls would reopen that race.
func appendSkipLog(path, line string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(line))
	return err
}

// writeSkipLine delivers one rendered line to BENCH_SKIP_LOG, or to w when that
// variable is unset. A failure to deliver the line fails the test loudly rather than
// silently swallowing the evidence a skip is supposed to leave behind.
func writeSkipLine(t testing.TB, w io.Writer, line string) {
	t.Helper()
	if path := os.Getenv(skipLogEnv); path != "" {
		if err := appendSkipLog(path, line); err != nil {
			t.Fatalf("capability: append skip log %s: %v", path, err)
		}
		return
	}
	if _, err := io.WriteString(w, line); err != nil {
		t.Fatalf("capability: write skip line: %v", err)
	}
}

// Capability skips t for a host that cannot run a security assertion in the named
// class, after delivering a structured line so the skip still shows once
// BENCH_SKIP_LOG is read back (or, unset, under a hand-run `go test -v`). class must
// be one of the enumerated Classes; an unrecognized class fails the test rather than
// vanishing from the line silently.
func Capability(t testing.TB, class Class, reason string) {
	t.Helper()
	line, err := renderCapabilityLine(class, reason)
	if err != nil {
		t.Fatalf("capability.Capability: %v", err)
	}
	writeSkipLine(t, stdout, line)
	t.Skip(reason)
}

// Environment skips t for a non-capability reason: an absent subject binary, an
// unset conformance root, an unmaterialized fixture. It delivers its structured line
// before skipping, for the same reason Capability does.
func Environment(t testing.TB, reason string) {
	t.Helper()
	writeSkipLine(t, stdout, renderEnvironmentLine(reason))
	t.Skip(reason)
}
