package preflighttest

import (
	"strings"
	"testing"
)

// DecisionMap renders a decision map whose Sources section holds records in order. Each
// record is one whole Sources entry, so a caller can render a malformed one.
func DecisionMap(records ...string) string {
	return "# Map\n\nStatus: shaping\n\n## Destination\n\nSettle it.\n\n## Sources\n\n" + strings.Join(records, "")
}

// PathRecord renders one well-formed Sources record that cites path.
func PathRecord(path string) string {
	return "- Path: `" + path + "`\n  Supports: the decision.\n  Drift: a change to the decision.\n"
}

// WriteDecisionMap writes at path a decision map whose Sources section holds records.
func WriteDecisionMap(t *testing.T, path string, records ...string) {
	t.Helper()
	MustWriteFile(t, path, DecisionMap(records...))
}

// RestartFeature runs setup on main in the SeedConformant repository and commits every
// change. It then restarts the feature branch at main with one fenced change of slug.
func RestartFeature(t *testing.T, slug string, setup func()) {
	t.Helper()
	RunGit(t, "checkout", "-q", "main")
	setup()
	CommitAll(t, "main setup")
	RunGit(t, "checkout", "-q", "-B", "feature")
	MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n\nfunc Foo(ctx int) int { return ctx }\n")
	CommitAll(t, "fenced change")
}
