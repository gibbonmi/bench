// Package responseboundtest reads the spill path of an over-bound response for consumer
// tests. The spill line has one reader here, so a test in any package finds the spill
// file the same way. The reader checks only the frame of the line and the path; the
// owner is the one source of the count fields.
package responseboundtest

import (
	"strings"
	"testing"
)

// Spill is the parsed spill line of one over-bound response.
type Spill struct {
	// Path is the absolute path of the spill file. It is the last field, so it can hold a
	// comma.
	Path string
}

const (
	linePrefix = "spilled{"
	pathField  = ",path="
	lineSuffix = "}"
)

// ParseLine parses one spill line, with or without its newline. A line without the spill
// frame or a path answers false.
func ParseLine(line string) (Spill, bool) {
	body, ok := strings.CutPrefix(strings.TrimSuffix(line, "\n"), linePrefix)
	if !ok {
		return Spill{}, false
	}
	body, ok = strings.CutSuffix(body, lineSuffix)
	if !ok {
		return Spill{}, false
	}
	_, path, ok := strings.Cut(body, pathField)
	if !ok || path == "" {
		return Spill{}, false
	}
	return Spill{Path: path}, true
}

// Find answers the first spill line of response. A response within the bound has none.
func Find(response string) (Spill, bool) {
	for _, line := range strings.SplitAfter(response, "\n") {
		if spill, ok := ParseLine(line); ok {
			return spill, true
		}
	}
	return Spill{}, false
}

// Path answers the spill path that response names, and fails the test without one.
func Path(t testing.TB, response string) string {
	t.Helper()
	spill, ok := Find(response)
	if !ok {
		t.Fatalf("response = %q, want a spill line", response)
	}
	return spill.Path
}
