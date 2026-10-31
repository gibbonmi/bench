package preflighttest

import (
	"strings"
	"testing"
)

// CommitAll commits every change of the working tree and returns the new head.
func CommitAll(t *testing.T, message string) string {
	t.Helper()
	RunGit(t, "add", "-A")
	RunGit(t, "commit", "-q", "-m", message)
	return RunGit(t, "rev-parse", "HEAD")
}

// ResponseTable decodes one table of a response. A table opens at the start of a line, and
// its rows are indented, so the first unindented line after the header closes it.
func ResponseTable(t *testing.T, out, name string) []map[string]any {
	t.Helper()
	var block strings.Builder
	for _, line := range strings.SplitAfter(out, "\n") {
		if block.Len() == 0 && strings.HasPrefix(line, name+"[") || block.Len() > 0 && strings.HasPrefix(line, "  ") {
			block.WriteString(line)
		} else if block.Len() > 0 {
			break
		}
	}
	if block.Len() == 0 {
		t.Fatalf("response holds no %s table:\n%s", name, out)
	}
	return TableRows(t, DecodeMap(t, block.String()), name)
}

// StalenessRow is the one row of the staleness table of an audit response.
func StalenessRow(t *testing.T, out string) map[string]any {
	t.Helper()
	return ResponseTable(t, out, "staleness")[0]
}
