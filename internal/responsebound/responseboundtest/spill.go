// Package responseboundtest reads the spill line of an over-bound response for consumer
// tests. The line has one parser here, so a test in any package reads the owner's fields
// the same way.
package responseboundtest

import (
	"fmt"
	"strings"
	"testing"
)

// Spill is the parsed spill line of one over-bound response.
type Spill struct {
	Lines, Bytes, OmittedLines, CutLines int64
	// Path is the absolute path of the spill file. It is the last field, so it can hold a
	// comma.
	Path string
}

const (
	linePrefix   = "spilled{"
	pathField    = ",path="
	lineSuffix   = "}"
	countsFormat = "lines=%d,bytes=%d,omitted_lines=%d,cut_lines=%d"
)

// ParseLine parses one spill line, with or without its newline. A line that is not the
// complete spill line answers false.
func ParseLine(line string) (Spill, bool) {
	body, ok := strings.CutPrefix(strings.TrimSuffix(line, "\n"), linePrefix)
	if !ok {
		return Spill{}, false
	}
	body, ok = strings.CutSuffix(body, lineSuffix)
	if !ok {
		return Spill{}, false
	}
	counts, path, ok := strings.Cut(body, pathField)
	if !ok || path == "" {
		return Spill{}, false
	}
	spill := Spill{Path: path}
	if _, err := fmt.Sscanf(counts, countsFormat, &spill.Lines, &spill.Bytes, &spill.OmittedLines, &spill.CutLines); err != nil {
		return Spill{}, false
	}
	if fmt.Sprintf(countsFormat, spill.Lines, spill.Bytes, spill.OmittedLines, spill.CutLines) != counts {
		return Spill{}, false
	}
	return spill, true
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
