package structure

import (
	"strings"
	"testing"
)

func TestTouchedPreservesPathBytesAndScope(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "1")
	t.Setenv("BENCH_MAX_DIR_FILES", "100")
	root := initRepo(t)
	run(t, root, "config", "core.quotePath", "true")
	write(t, root, "unchanged.go", lines(2))
	write(t, root, "deleted.go", lines(2))
	commit(t, root, "base")
	base := headSha(t, root)

	wants := []string{
		"café.go",
		"line\nbreak.go",
		"quote\"name.go",
		"space name.go",
		"tab\tname.go",
		"before \t.go",
	}
	for _, path := range wants {
		write(t, root, path, lines(2))
	}
	write(t, root, "after.go ", lines(2))
	run(t, root, "rm", "deleted.go")
	commit(t, root, "hostile source paths")

	write(t, root, "café.go", lines(3))
	write(t, root, "uncommitted.go", lines(2))

	report, violations, err := Touched(root, base)
	if err != nil {
		t.Fatal(err)
	}
	if violations != len(wants) {
		t.Errorf("violations = %d, want %d:\n%s", violations, len(wants), report)
	}
	for _, path := range wants {
		if !strings.Contains(report, "   "+path) {
			t.Errorf("report does not contain exact path bytes %q:\n%s", path, report)
		}
	}
	if !strings.Contains(report, "3 lines (max 1)   café.go") {
		t.Errorf("selected path was not read from the working tree:\n%s", report)
	}
	excluded := []string{"unchanged.go", "deleted.go", "uncommitted.go", "after.go "}
	for _, path := range excluded {
		if strings.Contains(report, path) {
			t.Errorf("excluded path %q reached the report:\n%s", path, report)
		}
	}

	t.Chdir(root)
	commandReport, code := Command([]string{"--since", base})
	if code != 1 {
		t.Errorf("Command exit = %d, want 1:\n%s", code, commandReport)
	}
	if got := strings.Count(commandReport, "FILE TOO LONG"); got != len(wants) {
		t.Errorf("Command violation rows = %d, want %d:\n%s", got, len(wants), commandReport)
	}
	for _, path := range wants {
		if !strings.Contains(commandReport, "   "+path) {
			t.Errorf("Command report does not contain exact path bytes %q:\n%s", path, commandReport)
		}
	}
	for _, path := range excluded {
		if strings.Contains(commandReport, path) {
			t.Errorf("excluded path %q reached the Command report:\n%s", path, commandReport)
		}
	}

	if failedReport, failedViolations, failedErr := Touched(root, "missing-base"); failedErr == nil || failedReport != "" || failedViolations != 0 {
		t.Errorf("failed Touched query = (%q, %d, %v), want empty report, zero violations, and error", failedReport, failedViolations, failedErr)
	}
	if failedReport, failedCode := Command([]string{"--since", "missing-base"}); failedCode != 1 || failedReport != "" {
		t.Errorf("failed Command query = (%q, %d), want empty report and exit 1", failedReport, failedCode)
	}
}
