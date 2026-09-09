// The staged prose form's fixtures and row tests. This file pins rows DG17 through DG25,
// DG34, DG37, and DG38 over real temporary repositories, because the index is the subject.

package gate

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// runGitIn runs one git command in root through the adapter form and fails the test on a
// nonzero exit. The staged form's subject is the index, so every fixture below is a real
// repository rather than a plain directory.
func runGitIn(t *testing.T, root string, args ...string) {
	t.Helper()
	if _, err := benchgit.Output(append([]string{"-C", root}, args...)...); err != nil {
		t.Fatalf("git %v in %s: %v", args, root, err)
	}
}

// policyRepo makes a repository whose one commit carries the exclusion file with body, so
// the index holds a policy before the test stages anything else.
func policyRepo(t *testing.T, body string) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	write(t, root, ".bench/prose-exclusions", body)
	runGitIn(t, root, "add", "-A")
	runGitIn(t, root, "commit", "-q", "-m", "policy")
	return root
}

// gradeStaged runs the staged form and returns its exit code with stdout. The staged form
// answers on stdout at every exit, so an empty stderr is part of the contract.
func gradeStaged(t *testing.T, root string) (int, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{root, "--staged"}, &stdout, &stderr)
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty: the staged form answers on stdout", stderr.String())
	}
	return code, stdout.String()
}

// TestGateProseStagedGradesTheIndex is DG17: the staged form grades the staged Markdown
// and prints its pass table for a clean index. The exclusion file is staged beside the
// subject and is not Markdown, so it carries no row.
func TestGateProseStagedGradesTheIndex(t *testing.T) {
	root := policyRepo(t, "")
	write(t, root, "docs/notes.md", "Short prose.\n")
	runGitIn(t, root, "add", "-A")

	code, out := gradeStaged(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q", code, out)
	}
	if want := wantProseTable(t, "docs/notes.md"); out != want {
		t.Fatalf("stdout = %q, want the pass table %q", out, want)
	}

	write(t, root, "docs/notes.md", words(27))
	runGitIn(t, root, "add", "-A")
	code, out = gradeStaged(t, root)
	if code != 1 {
		t.Fatalf("exit = %d on a staged over-long sentence, want 1; stdout=%q", code, out)
	}
	if !strings.Contains(out, `"docs/notes.md"`) {
		t.Fatalf("stdout = %q, want it to name docs/notes.md", out)
	}
}

// TestGateProseStagedIgnoresTheWorkingFile is DG18: the graded bytes come from the index
// blob. A long index blob with a short working file reds, and the inverse passes.
func TestGateProseStagedIgnoresTheWorkingFile(t *testing.T) {
	long := policyRepo(t, "")
	write(t, long, "docs/notes.md", words(27))
	runGitIn(t, long, "add", "-A")
	write(t, long, "docs/notes.md", "Short prose.\n")

	code, out := gradeStaged(t, long)
	if code != 1 {
		t.Fatalf("exit = %d with a long index blob, want 1; stdout=%q", code, out)
	}
	if !strings.Contains(out, `"docs/notes.md"`) {
		t.Fatalf("stdout = %q, want the index bytes graded", out)
	}

	short := policyRepo(t, "")
	write(t, short, "docs/notes.md", "Short prose.\n")
	runGitIn(t, short, "add", "-A")
	write(t, short, "docs/notes.md", words(27))

	code, out = gradeStaged(t, short)
	if code != 0 {
		t.Fatalf("exit = %d with a long working file, want 0; stdout=%q", code, out)
	}
	if want := wantProseTable(t, "docs/notes.md"); out != want {
		t.Fatalf("stdout = %q, want the pass table %q", out, want)
	}
}

// TestGateProseStagedReadsTheStagedPolicy is DG19: the policy is the index blob of the
// exclusion file. An index policy that excludes the subject passes while the working
// policy would red, and the reverse reds.
func TestGateProseStagedReadsTheStagedPolicy(t *testing.T) {
	excluding := policyRepo(t, "")
	write(t, excluding, "docs/notes.md", words(27))
	write(t, excluding, ".bench/prose-exclusions", "docs/notes.md a staged reason\n")
	runGitIn(t, excluding, "add", "-A")
	write(t, excluding, ".bench/prose-exclusions", "")

	code, out := gradeStaged(t, excluding)
	if code != 0 {
		t.Fatalf("exit = %d with a staged exclusion, want 0; stdout=%q", code, out)
	}

	including := policyRepo(t, "")
	write(t, including, "docs/notes.md", words(27))
	runGitIn(t, including, "add", "-A")
	write(t, including, ".bench/prose-exclusions", "docs/notes.md a working reason\n")

	code, out = gradeStaged(t, including)
	if code != 1 {
		t.Fatalf("exit = %d with the exclusion only in the working tree, want 1; stdout=%q", code, out)
	}
	if !strings.Contains(out, `"docs/notes.md"`) {
		t.Fatalf("stdout = %q, want the subject graded", out)
	}
}

// TestGateProseStagedSkipsDeletionsAndLinks is DG20: a staged deletion and a staged
// symbolic link are not subjects. Only the regular file beside them carries a row, and
// only it is reported when it is long.
func TestGateProseStagedSkipsDeletionsAndLinks(t *testing.T) {
	root := policyRepo(t, "")
	write(t, root, "docs/gone.md", "Short prose.\n")
	runGitIn(t, root, "add", "-A")
	runGitIn(t, root, "commit", "-q", "-m", "subjects")
	runGitIn(t, root, "rm", "-q", "docs/gone.md")
	write(t, root, "docs/kept.md", "Short prose.\n")
	if err := os.Symlink("kept.md", filepath.Join(root, "docs", "link.md")); err != nil {
		t.Fatalf("plant a symbolic link: %v", err)
	}
	runGitIn(t, root, "add", "-A")

	code, out := gradeStaged(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q", code, out)
	}
	if want := wantProseTable(t, "docs/kept.md"); out != want {
		t.Fatalf("stdout = %q, want only the regular file %q", out, want)
	}

	write(t, root, "docs/kept.md", words(27))
	runGitIn(t, root, "add", "-A")
	code, out = gradeStaged(t, root)
	if code != 1 {
		t.Fatalf("exit = %d on the long regular file, want 1; stdout=%q", code, out)
	}
	if strings.Contains(out, "link.md") || strings.Contains(out, "gone.md") {
		t.Fatalf("stdout = %q, want neither the link nor the deletion reported", out)
	}
}

// TestGateProseStagedHostilePaths is DG21: a staged path with a space, a double quote, a
// tab, or a newline is graded whole. The expected table derives through the encoder, so a
// split or dropped path reds and a quoted path does not red for the wrong reason.
func TestGateProseStagedHostilePaths(t *testing.T) {
	root := policyRepo(t, "")
	// Git orders the index by path bytes, so the rows arrive in this order.
	paths := []string{"docs/a b.md", "docs/new\nline.md", "docs/q\"uote.md", "docs/tab\t.md"}
	for _, rel := range paths {
		write(t, root, rel, "Short prose.\n")
	}
	runGitIn(t, root, "add", "-A")

	code, out := gradeStaged(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q", code, out)
	}
	if want := wantProseTable(t, paths...); out != want {
		t.Fatalf("stdout = %q, want the pass table %q", out, want)
	}
}

// TestGateProseStagedEmptyIndex is DG22: an empty staged set passes with the empty pass
// table at exit 0.
func TestGateProseStagedEmptyIndex(t *testing.T) {
	root := policyRepo(t, "")

	code, out := gradeStaged(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q", code, out)
	}
	if want := wantProseTable(t); out != want {
		t.Fatalf("stdout = %q, want the empty pass table %q", out, want)
	}
}

// TestGateProseStagedRefusesAPathList is DG23: `--staged` selects its own subjects, so a
// path list or a `--` separator beside it is a usage error at exit 2 with no stdout.
func TestGateProseStagedRefusesAPathList(t *testing.T) {
	root := policyRepo(t, "")
	write(t, root, "docs/notes.md", words(27))
	runGitIn(t, root, "add", "-A")

	for _, args := range [][]string{
		{root, "--staged", "docs/notes.md"},
		{root, "--staged", "--"},
		{root, "--staged", "--", "docs/notes.md"},
		{root, "--", "--staged"},
	} {
		var stdout, stderr bytes.Buffer
		code := GateProseCommand(args, &stdout, &stderr)
		if code != 2 {
			t.Errorf("exit for %q = %d, want 2; stdout=%q", args, code, stdout.String())
		}
		if stdout.Len() != 0 {
			t.Errorf("stdout for %q = %q, want empty on a usage error", args, stdout.String())
		}
		if want := gateProseUsage + "\n"; stderr.String() != want {
			t.Errorf("stderr for %q = %q, want usage %q", args, stderr.String(), want)
		}
	}
}

// TestGateProseStagedOutsideARepository is DG25: a root that is not a working-tree top
// answers a `prose:` refusal on stdout at exit 1. A plain temporary directory must not
// answer an empty pass.
func TestGateProseStagedOutsideARepository(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".bench/prose-exclusions", "")

	code, out := gradeStaged(t, root)
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q", code, out)
	}
	if !strings.HasPrefix(out, "prose: ") {
		t.Fatalf("stdout = %q, want a prose refusal line", out)
	}
	if strings.Contains(out, "prose[") {
		t.Fatalf("stdout = %q, want no pass table outside a repository", out)
	}
}

// TestGateProseStagedValidatesTargetsAgainstTheIndex is DG34: an exclusion target present
// in the index and absent from the working tree is honored, a directory prefix of index
// entries is a directory row, and a target absent from the index reds the policy.
func TestGateProseStagedValidatesTargetsAgainstTheIndex(t *testing.T) {
	honored := policyRepo(t, "")
	write(t, honored, "docs/notes.md", words(27))
	write(t, honored, ".bench/prose-exclusions", "docs/notes.md an index-only target\n")
	runGitIn(t, honored, "add", "-A")
	if err := os.Remove(filepath.Join(honored, "docs", "notes.md")); err != nil {
		t.Fatalf("remove the working file: %v", err)
	}

	code, out := gradeStaged(t, honored)
	if code != 0 {
		t.Fatalf("exit = %d for an index-only target, want 0; stdout=%q", code, out)
	}

	directory := policyRepo(t, "")
	write(t, directory, "docs/notes.md", words(27))
	write(t, directory, ".bench/prose-exclusions", "docs/ an index directory prefix\n")
	runGitIn(t, directory, "add", "-A")
	if err := os.RemoveAll(filepath.Join(directory, "docs")); err != nil {
		t.Fatalf("remove the working directory: %v", err)
	}

	code, out = gradeStaged(t, directory)
	if code != 0 {
		t.Fatalf("exit = %d for an index directory prefix, want 0; stdout=%q", code, out)
	}

	absent := policyRepo(t, "")
	write(t, absent, "docs/notes.md", "Short prose.\n")
	write(t, absent, ".bench/prose-exclusions", "docs/working-only.md a working-tree target\n")
	runGitIn(t, absent, "add", "-A")
	write(t, absent, "docs/working-only.md", "Short prose.\n")

	code, out = gradeStaged(t, absent)
	if code != 1 {
		t.Fatalf("exit = %d for a working-only target, want 1; stdout=%q", code, out)
	}
	if !strings.Contains(out, "names an absent path") {
		t.Fatalf("stdout = %q, want the absent-path diagnostic", out)
	}
}

// TestGateProseStagedUnbornBranch is DG37: with no HEAD to diff against, every index entry
// is staged. A form that diffs against HEAD refuses or grades nothing here.
func TestGateProseStagedUnbornBranch(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	write(t, root, ".bench/prose-exclusions", "")
	write(t, root, "docs/notes.md", "Short prose.\n")
	runGitIn(t, root, "add", "-A")

	code, out := gradeStaged(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q", code, out)
	}
	if want := wantProseTable(t, "docs/notes.md"); out != want {
		t.Fatalf("stdout = %q, want the pass table %q", out, want)
	}

	write(t, root, "docs/notes.md", words(27))
	runGitIn(t, root, "add", "-A")
	code, out = gradeStaged(t, root)
	if code != 1 {
		t.Fatalf("exit = %d on an unborn branch with a long subject, want 1; stdout=%q", code, out)
	}
	if !strings.Contains(out, `"docs/notes.md"`) {
		t.Fatalf("stdout = %q, want the subject graded", out)
	}
}

// TestGateProseStagedAbsentPolicy is DG38: an index without the exclusion file answers the
// absent-policy diagnostic at exit 1. A form that read an absent policy as empty would
// pass here.
func TestGateProseStagedAbsentPolicy(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	runGitIn(t, root, "commit", "-q", "--allow-empty", "-m", "base")
	write(t, root, "docs/notes.md", "Short prose.\n")
	runGitIn(t, root, "add", "-A")

	code, out := gradeStaged(t, root)
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q", code, out)
	}
	if !strings.Contains(out, "the exclusion file is absent") {
		t.Fatalf("stdout = %q, want the absent-policy diagnostic", out)
	}
}
