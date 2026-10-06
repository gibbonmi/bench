package testreport

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/prose"
)

// cleanProseTree writes an empty exclusion file and the two clean subjects `b.md` and
// `b/x.md`. The walk visits `b/x.md` first, so a sorted list differs from the walk order.
func cleanProseTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeProseCheckFile(t, root, ".bench/prose-exclusions", "")
	writeProseCheckFile(t, root, "b/x.md", "Short prose.\n")
	writeProseCheckFile(t, root, "b.md", "Short prose.\n")
	return root
}

// TestProseGreenPrintsOnlyCheckRow grades that a green prose run prints the check row
// and nothing else.
func TestProseGreenPrintsOnlyCheckRow(t *testing.T) {
	output, code := Command(cleanProseTree(t), []string{"--check", "prose"})
	want := checkHeader + "  prose,prose,0,2\n"
	if output != want || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, %q)", code, output, want)
	}
}

// TestProseFullListsSubjects grades that `--full` adds the graded subjects in sorted order.
func TestProseFullListsSubjects(t *testing.T) {
	output, code := Command(cleanProseTree(t), []string{"--full", "--check", "prose"})
	want := checkHeader + "  prose,prose,0,2\nsubjects[2]{path}:\n  b.md\n  b/x.md\n"
	if output != want || code != 0 {
		t.Fatalf("Command = (%d, %q), want (0, %q)", code, output, want)
	}
}

// TestProseZeroSubjectsExitsOne grades the zero rule for the prose check: a tree whose
// only document is excluded graded nothing, so the run is no evidence.
func TestProseZeroSubjectsExitsOne(t *testing.T) {
	root := t.TempDir()
	writeProseCheckFile(t, root, ".bench/prose-exclusions", "skip.md planted fixture text\n")
	writeProseCheckFile(t, root, "skip.md", "Short prose.\n")
	output, code := Command(root, []string{"--check", "prose"})
	want := checkHeader + "  prose,prose,0,0\nerror: named check ran nothing"
	if !strings.HasPrefix(output, want) || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, prefix %q)", code, output, want)
	}
}

// TestProseRedKeepsFindingsAfterCheckRow grades that a red prose run prints the check row
// and then each finding line of the grader.
func TestProseRedKeepsFindingsAfterCheckRow(t *testing.T) {
	root := t.TempDir()
	writeProseCheckFile(t, root, ".bench/prose-exclusions", "")
	writeProseCheckFile(t, root, "first.md", strings.Repeat("word ", 27)+".\n")
	writeProseCheckFile(t, root, "second.md", strings.Repeat("word ", 27)+".\n")
	findings := prose.Grade(root)
	if len(findings) != 2 {
		t.Fatalf("prose.Grade() = %q, want two findings", findings)
	}
	output, code := Command(root, []string{"--check", "prose"})
	want := checkHeader + "  prose,prose,0,2\n" + strings.Join(findings, "\n") + "\n"
	if output != want || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, %q)", code, output, want)
	}
}

// TestProseGraderRefusalPrintsCheckRow grades that a grader refusal prints the check row
// with no subject, then the refusal diagnostic, and not the zero-rule title.
func TestProseGraderRefusalPrintsCheckRow(t *testing.T) {
	root := t.TempDir()
	writeProseCheckFile(t, root, "only.md", "Short prose.\n")
	diagnostics := prose.Grade(root)
	if len(diagnostics) != 1 {
		t.Fatalf("prose.Grade() = %q, want one refusal diagnostic", diagnostics)
	}
	output, code := Command(root, []string{"--check", "prose"})
	want := checkHeader + "  prose,prose,0,0\n" + diagnostics[0] + "\n"
	if output != want || code != 1 {
		t.Fatalf("Command = (%d, %q), want (1, %q)", code, output, want)
	}
}
