// Package capability owns skipping for the gate's test suite. A bare t.Skip is
// invisible under non-verbose `go test` — a skip and a pass print the same nothing,
// and `go test` without -v discards a package's stdout and stderr entirely — so
// every skip here writes a structured line to the file named by BENCH_SKIP_LOG
// first and only then calls t.Skip. The gate sets that variable to a run-scoped
// path, reads it back once the phase's `go test` invocations exit, and tallies skips
// by kind and class. With BENCH_SKIP_LOG unset — a developer's hand-run
// `go test -v` — the line falls back to stdout.
//
// Render and ParseLine are the two halves of one line shape, kept in this package so
// the writer and the gate's collector cannot drift apart.
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
	"strings"
)

// TB is the subset of testing.TB the skip helpers use. Depending on this interface
// rather than on testing.TB is what keeps the testing package out of every binary
// that links this one: the gate's collector imports it for Render and ParseLine.
type TB interface {
	Helper()
	Fatalf(format string, args ...any)
	Name() string
	Skip(args ...any)
}

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

// classOrder is the closed vocabulary and the order every consumer reports classes
// in, so a tally's rows are stable run to run rather than map-iteration order.
var classOrder = []Class{Symlink, Fifo, PID, CPU, Privilege, Signal, Tool}

// Classes returns the enumerated capability classes in reporting order.
func Classes() []Class { return append([]Class(nil), classOrder...) }

func (c Class) valid() bool {
	for _, known := range classOrder {
		if c == known {
			return true
		}
	}
	return false
}

// Kind separates the two skip populations. Only KindCapability describes a security
// assertion the host could not make, so only it belongs in a strict-mode count.
type Kind string

const (
	KindCapability  Kind = "capability"
	KindEnvironment Kind = "environment"
)

// Skip is one skip's structured content. Class is empty for KindEnvironment. Name is
// the emitting test, carried because a count cannot tell a reader which assertion went
// unmade — the gate reports skips long after the test binary that emitted them exited,
// so the name has to travel on the line or be lost.
type Skip struct {
	Kind   Kind
	Class  Class
	Name   string
	Reason string
}

// linePrefix opens every line this package writes, so a downstream collector can
// recognize one with a prefix match before parsing the key=value tokens that follow.
const linePrefix = "bench-skip"

// LogEnv names the environment variable the gate sets to a run-scoped absolute file
// path. Many test binaries append to that one file concurrently from separate
// processes, so writers must go through appendSkipLog rather than opening it by hand.
const LogEnv = "BENCH_SKIP_LOG"

// WithoutEnvironment returns env without entries for name, preserving every other entry
// byte-for-byte so a child cannot inherit a caller-owned capability side channel.
func WithoutEnvironment(env []string, name string) []string {
	prefix := name + "="
	filtered := make([]string, 0, len(env))
	for _, entry := range env {
		if !strings.HasPrefix(entry, prefix) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// stdout is the stdout-fallback destination, held behind a var so this package's own
// tests can swap it for a buffer and assert on the emitted bytes without touching the
// process's real stdout.
var stdout io.Writer = os.Stdout

// Render is the one place a skip line is built, newline included. Reason takes the
// remainder of the line verbatim (no escaping), so a reader need only split the
// leading kind, class, and name tokens off the front. A class outside the enumerated
// set is refused rather than formatted: an open vocabulary would let a typo mint a
// class nothing counts. An empty or space-carrying name is refused for a narrower
// reason: name is a fixed-width token in a space-delimited line, and an unnamed skip
// is exactly the unactionable count this line shape exists to replace.
func Render(skip Skip) (string, error) {
	if skip.Name == "" || strings.ContainsAny(skip.Name, " \t") {
		return "", fmt.Errorf("capability: skip name %q must be a single non-empty token", skip.Name)
	}
	switch skip.Kind {
	case KindCapability:
		if !skip.Class.valid() {
			return "", fmt.Errorf("capability: unknown class %q", skip.Class)
		}
		return fmt.Sprintf("%s kind=%s class=%s name=%s reason=%s\n", linePrefix, skip.Kind, skip.Class, skip.Name, skip.Reason), nil
	case KindEnvironment:
		return fmt.Sprintf("%s kind=%s name=%s reason=%s\n", linePrefix, skip.Kind, skip.Name, skip.Reason), nil
	}
	return "", fmt.Errorf("capability: unknown kind %q", skip.Kind)
}

// ParseLine reads back one line Render produced, reporting false for anything else —
// a phase's ordinary output shares the log with nothing, but the gate's collector
// must not mistake a stray line for evidence. A capability line naming a class
// outside the enumerated set is refused for the same reason Render refuses to write
// one.
func ParseLine(line string) (Skip, bool) {
	rest, ok := strings.CutPrefix(strings.TrimRight(line, "\r\n"), linePrefix+" ")
	if !ok {
		return Skip{}, false
	}
	kind, rest, ok := cutField(rest, "kind")
	if !ok {
		return Skip{}, false
	}
	skip := Skip{Kind: Kind(kind)}
	switch skip.Kind {
	case KindCapability:
		class, tail, ok := cutField(rest, "class")
		if !ok || !Class(class).valid() {
			return Skip{}, false
		}
		skip.Class, rest = Class(class), tail
	case KindEnvironment:
	default:
		return Skip{}, false
	}
	name, rest, ok := cutField(rest, "name")
	if !ok || name == "" {
		return Skip{}, false
	}
	skip.Name = name
	reason, ok := strings.CutPrefix(rest, "reason=")
	if !ok {
		return Skip{}, false
	}
	skip.Reason = reason
	return skip, true
}

// cutField takes one space-terminated `name=value` token off the front of rest.
func cutField(rest, name string) (value, tail string, ok bool) {
	rest, ok = strings.CutPrefix(rest, name+"=")
	if !ok {
		return "", "", false
	}
	value, tail, ok = strings.Cut(rest, " ")
	return value, tail, ok
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
func writeSkipLine(t TB, w io.Writer, line string) {
	t.Helper()
	if path := os.Getenv(LogEnv); path != "" {
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
func Capability(t TB, class Class, reason string) {
	t.Helper()
	line, err := Render(Skip{Kind: KindCapability, Class: class, Name: t.Name(), Reason: reason})
	if err != nil {
		t.Fatalf("capability.Capability: %v", err)
	}
	writeSkipLine(t, stdout, line)
	t.Skip(reason)
}

// Environment skips t for a non-capability reason: an absent subject binary, an
// unset conformance root, an unmaterialized fixture. It delivers its structured line
// before skipping, for the same reason Capability does.
func Environment(t TB, reason string) {
	t.Helper()
	line, err := Render(Skip{Kind: KindEnvironment, Name: t.Name(), Reason: reason})
	if err != nil {
		t.Fatalf("capability.Environment: %v", err)
	}
	writeSkipLine(t, stdout, line)
	t.Skip(reason)
}
