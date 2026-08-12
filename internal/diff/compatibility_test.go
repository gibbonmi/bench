package diff

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreMigrationResponsesDifferOnlyByNamedAXIDelta(t *testing.T) {
	sourceDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root, base, feature := seedCompatibilityRepo(t)
	for _, tc := range []struct {
		name    string
		fixture string
		args    []string
		full    bool
	}{
		{name: "bare", fixture: "pre-migration-bare.toon"},
		{name: "full", fixture: "pre-migration-full.toon", args: []string{"--full"}, full: true},
		{name: "commit", fixture: "pre-migration-commit.toon", args: []string{"--commit", feature}},
		{name: "commit full", fixture: "pre-migration-commit-full.toon", args: []string{"--commit", feature, "--full"}, full: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			baseline, err := os.ReadFile(filepath.Join(sourceDir, "testdata", tc.fixture))
			if err != nil {
				t.Fatal(err)
			}
			candidate, code := runProductionDiff(t, root, tc.args...)
			if code != 0 {
				t.Fatalf("candidate exit = %d:\n%s", code, candidate)
			}
			assertPreservedTrackedRows(t, string(baseline), candidate)
			if tc.full {
				if got, want := tableBlock(t, candidate, "log"), tableBlock(t, string(baseline), "log"); got != want {
					t.Fatalf("log bytes changed outside the named delta:\ngot:\n%swant:\n%s", got, want)
				}
				oldBody := strings.SplitN(string(baseline), "diff_body:\n", 2)[1]
				newBody := strings.SplitN(candidate, "diff_body:\n", 2)[1]
				if !strings.HasPrefix(newBody, oldBody) {
					t.Fatalf("tracked diff_body bytes changed outside the named delta:\ngot:\n%swant prefix:\n%s", newBody, oldBody)
				}
			}
		})
	}

	for _, tc := range []struct {
		name string
		args []string
		code int
		out  string
	}{
		{name: "unknown flag", args: []string{"--bogus"}, code: 2, out: "usage: bench diff (unknown argument: --bogus)\n"},
		{name: "missing commit value", args: []string{"--commit"}, code: 2, out: "usage: bench diff (missing argument: --commit)\n"},
		{name: "unresolvable commit", args: []string{"--commit", "missing"}, code: 1, out: "error: cannot resolve --commit — 'missing' does not name a commit reachable in this repository\n"},
		{name: "root commit", args: []string{"--commit", base}, code: 1, out: "error: --commit has no parent — '" + base + "' is a root commit — there is no first parent to diff against\n"},
		{name: "repeated flag", args: []string{"--full", "--full"}, code: 2, out: "usage: bench diff (unknown argument: --full)\n"},
	} {
		t.Run("argv "+tc.name, func(t *testing.T) {
			out, code := runProductionDiff(t, root, tc.args...)
			if code != tc.code || out != tc.out {
				t.Fatalf("response = (%d,%q), want (%d,%q)", code, out, tc.code, tc.out)
			}
		})
	}

	left, leftCode := runProductionDiff(t, root, "--commit", feature, "--full")
	right, rightCode := runProductionDiff(t, root, "--full", "--commit", feature)
	if leftCode != 0 || rightCode != 0 || left != right {
		t.Fatalf("accepted flag orders diverged: left=(%d,%q) right=(%d,%q)", leftCode, left, rightCode, right)
	}
}

func seedCompatibilityRepo(t *testing.T) (root, base, feature string) {
	t.Helper()
	root = t.TempDir()
	t.Chdir(root)
	runGit(t, "init", "-q", "-b", "main")
	runGit(t, "config", "user.email", "fixture@example.com")
	runGit(t, "config", "user.name", "Fixture")
	runGit(t, "config", "core.autocrlf", "false")
	mustWriteFile(t, "tracked.txt", "base line\n")
	mustWriteFile(t, "delete.txt", "delete me\n")
	runGit(t, "add", "tracked.txt", "delete.txt")
	runGitEnv(t, root, []string{"GIT_AUTHOR_DATE=2001-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2001-01-01T00:00:00Z"}, "commit", "-q", "-m", "base")
	base = runGit(t, "rev-parse", "HEAD")
	runGit(t, "switch", "-q", "-c", "feature")
	mustWriteFile(t, "tracked.txt", "base line\ncommitted line\n")
	mustWriteFile(t, "committed.txt", "committed file\n")
	runGit(t, "add", "tracked.txt", "committed.txt")
	runGitEnv(t, root, []string{"GIT_AUTHOR_DATE=2001-01-02T00:00:00Z", "GIT_COMMITTER_DATE=2001-01-02T00:00:00Z"}, "commit", "-q", "-m", "feature commit")
	feature = runGit(t, "rev-parse", "HEAD")
	if base != "06d53f18d44494ead344d1f720fdb978128684a7" || feature != "f590a43d52b3b03f6d897375ee8f530f48efb596" {
		t.Fatalf("compatibility fixture identity drifted: base=%s feature=%s", base, feature)
	}
	runGit(t, "config", "branch.feature.benchBase", base)
	mustWriteFile(t, "tracked.txt", "base line\ncommitted line\nstaged line\n")
	runGit(t, "add", "tracked.txt")
	mustWriteFile(t, "tracked.txt", "base line\ncommitted line\nstaged line\nunstaged line\n")
	mustWriteFile(t, "untracked.txt", "new bytes\n")
	return root, base, feature
}

func runGitEnv(t *testing.T, root string, env []string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	cmd.Env = append(os.Environ(), env...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func assertPreservedTrackedRows(t *testing.T, baseline, candidate string) {
	t.Helper()
	oldRows := tableRows(t, baseline, "files")
	newRows := tableRows(t, candidate, "files")
	var preserved [][]string
	for _, row := range newRows {
		if len(row) == 3 && row[0] != "?" {
			preserved = append(preserved, row[:2])
		}
	}
	if !bytes.Equal([]byte(strings.TrimSpace(rowsText(oldRows))), []byte(strings.TrimSpace(rowsText(preserved)))) {
		t.Fatalf("tracked files rows changed outside the named kind-column delta: old=%v new=%v", oldRows, preserved)
	}
}

func rowsText(rows [][]string) string {
	var b strings.Builder
	for _, row := range rows {
		b.WriteString(strings.Join(row, "\x00"))
		b.WriteByte('\n')
	}
	return b.String()
}

func tableBlock(t *testing.T, out, name string) string {
	t.Helper()
	lines := strings.SplitAfter(out, "\n")
	prefix := name + "["
	for i, line := range lines {
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		rows := tableRows(t, out, name)
		return strings.Join(lines[i:i+1+len(rows)], "")
	}
	t.Fatalf("missing %s table", name)
	return ""
}
