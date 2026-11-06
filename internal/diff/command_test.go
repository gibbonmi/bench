package diff

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

func TestCommandCoherentSnapshotIncludesAllFilesFacts(t *testing.T) {
	_, _, base, _ := seedDivergedRepo(t)
	if err := os.MkdirAll("nested/new", 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, "nested/new/untracked.txt", "untracked\n")
	mustWriteFile(t, "f.txt", "worktree change \n")
	gittest.MustRun(t, "", "add", "f.txt")
	mustWriteFile(t, "f.txt", "worktree change \nsecond \n")

	out, code := Command(nil)
	if code != 0 {
		t.Fatalf("Command exit = %d, output:\n%s", code, out)
	}
	for _, want := range []string{
		"revision[1]{branch,default,ahead,behind,base,method,head}:",
		"aggregate[1]{commits,files,insertions,deletions,staged,unstaged,untracked}:",
		"files[2]{status,path,kind}:",
		"?,nested/new/untracked.txt,\"\"",
		"checkout[2]{index,worktree,path}:",
		"whitespace[1]{clean,offenses}:",
		"\"false\",\"" + rawWhitespaceOffenses(t, base) + "\"",
		"help[1]{cmd,why}:\n  bench diff --full,inspect the complete patch\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output does not contain %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, base) {
		t.Errorf("revision does not carry raw Git merge-base %q:\n%s", base, out)
	}
}

func TestCommandFullAppendsSortedUntrackedPatch(t *testing.T) {
	seedDivergedRepo(t)
	mustWriteFile(t, "z.txt", "z\n")
	mustWriteFile(t, "a.txt", "a\n")
	out, code := Command([]string{"--full"})
	if code != 0 {
		t.Fatalf("Command --full exit = %d, output:\n%s", code, out)
	}
	a := rawUntrackedPatch(t, "a.txt")
	z := rawUntrackedPatch(t, "z.txt")
	if !strings.Contains(out, a) || !strings.Contains(out, z) {
		t.Fatalf("full output omitted a raw-Git untracked patch:\n%s", out)
	}
	if strings.Index(out, a) > strings.Index(out, z) {
		t.Fatalf("untracked patches are not path sorted:\n%s", out)
	}
	if !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
		t.Fatalf("complete full output does not end in empty help:\n%s", out)
	}
}

func rawUntrackedPatch(t *testing.T, path string) string {
	t.Helper()
	out, err := exec.Command("git", "diff", "--no-index", "--", "/dev/null", path).CombinedOutput()
	if err == nil {
		t.Fatalf("git diff --no-index %s unexpectedly exited 0", path)
	}
	return string(out)
}

func rawWhitespaceOffenses(t *testing.T, base string) string {
	t.Helper()
	out, err := exec.Command("git", "diff", "--check", base).CombinedOutput()
	if err == nil {
		t.Fatal("fixture has no raw-Git whitespace offense")
	}
	return strconv.Itoa(len(strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")))
}

func TestCommandRetriesThenRefusesSnapshotDrift(t *testing.T) {
	seedDivergedRepo(t)
	mustWriteFile(t, "f.txt", "dirty\n")
	previous := snapshotAfterRead
	defer func() { snapshotAfterRead = previous }()
	calls := 0
	snapshotAfterRead = func() {
		calls++
		mustWriteFile(t, "f.txt", "dirty "+string(rune('0'+calls))+"\n")
	}
	out, code := Command(nil)
	if code != 1 || calls != 2 {
		t.Fatalf("Command drift = (code %d, calls %d), want (1, 2); output:\n%s", code, calls, out)
	}
	for _, want := range []string{"error: snapshot drift", "dirty content", "help[1]{cmd,why}:", "bench diff,retry after the repository stopped moving"} {
		if !strings.Contains(out, want) {
			t.Errorf("drift output does not contain %q:\n%s", want, out)
		}
	}
}

func TestDiffPackageDoesNotOwnPorcelainInvocationOrParsing(t *testing.T) {
	for _, file := range []string{"diff.go", "snapshot.go", "range.go"} {
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, forbidden := range []string{`"status", "--porcelain`, "ParsePorcelainZ("} {
			if strings.Contains(string(source), forbidden) {
				t.Errorf("internal/diff production file %s owns forbidden porcelain knowledge %q", file, forbidden)
			}
		}
	}
}

// TestSnapshotOverrideRestoresAfterFailure proves that the movement seam's cleanup runs
// on the path of a fatal assertion. A helper child runs this test again: its subtest
// installs a distinguishable replacement and fails fatally, and the child then reports
// which collaborator answers. The parent requires the failed subtest and the prior one.
func TestSnapshotOverrideRestoresAfterFailure(t *testing.T) {
	if os.Getenv("BENCH_DIFF_RESTORE_HELPER") == "1" {
		prior, replacement := 0, 0
		t.Cleanup(SetSnapshotAfterReadForTest(func() { prior++ }))
		t.Run("fatal path", func(t *testing.T) {
			restore := SetSnapshotAfterReadForTest(func() { replacement++ })
			t.Cleanup(restore)
			if restore == nil {
				t.Fatal("the setter returned no restoration")
			}
			snapshotAfterRead()
			t.Fatal("the fatal path ends the subtest")
		})
		snapshotAfterRead()
		fmt.Printf("observed prior=%d replacement=%d\n", prior, replacement)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestSnapshotOverrideRestoresAfterFailure$")
	cmd.Env = append(os.Environ(), "BENCH_DIFF_RESTORE_HELPER=1")
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "--- FAIL: TestSnapshotOverrideRestoresAfterFailure/fatal_path") {
		t.Fatalf("the child's fatal path did not fail (%v):\n%s", err, out)
	}
	if !strings.Contains(string(out), "observed prior=1 replacement=1\n") {
		t.Fatalf("after the fatal path the child observed, want prior=1 replacement=1:\n%s", out)
	}
}
