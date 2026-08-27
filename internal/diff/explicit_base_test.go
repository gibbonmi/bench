package diff

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestExplicitBaseReportsPinnedSourceAndCompleteLiveSnapshot(t *testing.T) {
	root, base, feature := seedCompatibilityRepo(t)
	configBefore, refsBefore := gitState(t, root)
	out, code := runProductionDiff(t, root, "--base", base)
	if code != 0 {
		t.Fatalf("explicit diff exit = %d:\n%s", code, out)
	}
	for _, want := range []string{base, feature, "committed.txt", "tracked.txt", "untracked.txt"} {
		if !strings.Contains(out, want) {
			t.Fatalf("explicit diff omitted %q:\n%s", want, out)
		}
	}
	assertGitState(t, root, configBefore, refsBefore)
}

func TestExplicitBaseRefusalsLeaveConfigAndRefsUntouched(t *testing.T) {
	root, base, feature := seedCompatibilityRepo(t)
	other := unrelatedCommit(t, root)
	ambiguous := ambiguousObjectPrefix(t, root)
	for _, args := range [][]string{{"--base", "missing"}, {"--base", "bad^{commit}"}, {"--base", ambiguous}, {"--base", other}, {"--base", base, "--commit", feature}} {
		config, refs := gitState(t, root)
		out, code := Command(args)
		if code == 0 || !strings.HasPrefix(out, "error:") && !strings.HasPrefix(out, "usage:") {
			t.Fatalf("Command(%v) = (%d,%q), want refusal", args, code, out)
		}
		assertGitState(t, root, config, refs)
	}
}

func gitState(t *testing.T, root string) ([]byte, string) {
	t.Helper()
	config, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	return config, runGit(t, "show-ref", "--head")
}

func assertGitState(t *testing.T, root string, config []byte, refs string) {
	t.Helper()
	after, err := os.ReadFile(filepath.Join(root, ".git", "config"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(config) || runGit(t, "show-ref", "--head") != refs {
		t.Fatal("explicit-base operation changed Git config or refs")
	}
}

func unrelatedCommit(t *testing.T, root string) string {
	t.Helper()
	runGit(t, "checkout", "-q", "--orphan", "unrelated")
	runGit(t, "rm", "-q", "-rf", ".")
	mustWriteFile(t, "unrelated.txt", "unrelated\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "unrelated")
	sha := runGit(t, "rev-parse", "HEAD")
	runGit(t, "checkout", "-q", "feature")
	return sha
}

func ambiguousObjectPrefix(t *testing.T, root string) string {
	t.Helper()
	objects := filepath.Join(root, "objects")
	if err := os.MkdirAll(objects, 0o755); err != nil {
		t.Fatal(err)
	}
	seen := map[string]string{}
	for i := 0; i < 4096; i++ {
		path := filepath.Join(objects, strconv.Itoa(i))
		mustWriteFile(t, path, strconv.Itoa(i))
		sha := runGit(t, "hash-object", "-w", path)
		prefix := sha[:4]
		if previous, ok := seen[prefix]; ok && previous != sha {
			if len(strings.Fields(runGit(t, "rev-parse", "--disambiguate="+prefix))) >= 2 {
				return prefix
			}
		}
		seen[prefix] = sha
	}
	t.Fatal("could not create an ambiguous Git object prefix")
	return ""
}

func TestExplicitAndCommitAreMutuallyExclusive(t *testing.T) {
	out, code := Command([]string{"--base", "HEAD", "--commit", "HEAD"})
	if code != 2 || out == "" {
		t.Fatalf("mutual exclusion = (%d,%q), want usage refusal", code, out)
	}
}

func TestExplicitBaseRetriesThenRefusesMovement(t *testing.T) {
	root, base, _ := seedCompatibilityRepo(t)
	old := snapshotAfterRead
	defer func() { snapshotAfterRead = old }()
	calls := 0
	snapshotAfterRead = func() { calls++; mustWriteFile(t, filepath.Join(root, "tracked.txt"), string(rune('a'+calls))+"\n") }
	t.Chdir(root)
	out, code := Command([]string{"--base", base})
	if code != 1 || calls != 2 {
		t.Fatalf("explicit movement = (%d,%d): %s", code, calls, out)
	}
	if len(out) < 20 || out[:len("error: snapshot drift")] != "error: snapshot drift" {
		t.Fatalf("missing structured drift: %s", out)
	}
}

func TestChangedSubjectDefaultLiveComposition(t *testing.T) {
	root, base, _ := seedCompatibilityRepo(t)
	subject, kind, hint := ResolveChangedSubject(root, "", "")
	if kind != "" {
		t.Fatalf("ResolveChangedSubject = (%q, %q)", kind, hint)
	}
	if !subject.Live || subject.Base != base {
		t.Fatalf("subject identity = %#v, want live subject from %s", subject, base)
	}
	for _, want := range []string{"committed.txt", "tracked.txt", "untracked.txt"} {
		found := false
		for _, path := range subject.Paths {
			found = found || path == want
		}
		if !found {
			t.Errorf("paths = %v, want %q", subject.Paths, want)
		}
	}
}

func TestChangedSubjectRetriesThenRefusesMovement(t *testing.T) {
	root, base, _ := seedCompatibilityRepo(t)
	old := snapshotAfterRead
	defer func() { snapshotAfterRead = old }()
	calls := 0
	snapshotAfterRead = func() {
		calls++
		mustWriteFile(t, filepath.Join(root, "tracked.txt"), "moved "+strconv.Itoa(calls)+"\n")
	}
	_, kind, hint := ResolveChangedSubject(root, base, "")
	if kind != "changed subject drift" || calls != 2 {
		t.Fatalf("movement = (%q, %q), calls=%d, want second-drift refusal", kind, hint, calls)
	}
}
