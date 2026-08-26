// Package testlines classifies the plain-text output lines of a go test run.
package testlines

import "strings"

// RunnerLine reports whether the line is go test's own bookkeeping rather than a diagnostic.
func RunnerLine(line string) bool {
	return strings.HasPrefix(line, "bench-skip ") || strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "=== RUN") || strings.HasPrefix(line, "=== PAUSE") || strings.HasPrefix(line, "=== CONT") || strings.HasPrefix(line, "=== NAME") || strings.HasPrefix(line, "--- FAIL:") || strings.HasPrefix(line, "--- PASS:") || strings.HasPrefix(line, "--- SKIP:") || line == "FAIL" || strings.HasPrefix(line, "FAIL ") || strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "ok\t") || strings.HasPrefix(line, "? ") || strings.HasPrefix(line, "?\t")
}

// FailureRows returns the failure rows of one red plain-text stream, never -json output.
// A "--- FAIL:" line at any indent and a "# <package>" line each open a block whose later
// lines are rows until the next runner line. The caller owns the fallback for an empty result.
func FailureRows(lines []string) []string {
	rows := make([]string, 0)
	inBlock := false
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "--- FAIL:"), strings.HasPrefix(line, "# "):
			inBlock = true
			rows = append(rows, line)
		case strings.HasPrefix(line, "FAIL ") || strings.HasPrefix(line, "FAIL\t"):
			inBlock = false
			rows = append(rows, line)
		case strings.HasPrefix(line, "panic:"):
			rows = append(rows, line)
		case RunnerLine(line):
			inBlock = false
		case inBlock:
			rows = append(rows, line)
		}
	}
	return rows
}
