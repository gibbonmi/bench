package poolkey

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// cksumGolden pins POSIX cksum(root+newline), including a spaced path.
var cksumGolden = []struct {
	root string
	sum  uint32
}{
	{"/home/mgibs/workspace/bench", 2826441890},
	{"/tmp/a b/c", 889650394},
}

func TestCksumMatchesGolden(t *testing.T) {
	t.Parallel()
	for _, g := range cksumGolden {
		got := cksum([]byte(g.root + "\n"))
		if got != g.sum {
			t.Errorf("cksum(%q+NL) = %d, want %d", g.root, got, g.sum)
		}
	}
}

// TestCksumMatchesSystemTool cross-checks against the live `cksum` when it is on
// PATH, so the pinned goldens can never silently drift from the real tool. It is
// skipped where `cksum` is unavailable, keeping the suite hermetic there.
func TestCksumMatchesSystemTool(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("cksum"); err != nil {
		capability.Capability(t, capability.Tool, "cksum not available")
	}
	for _, g := range cksumGolden {
		cmd := exec.Command("cksum")
		cmd.Stdin = strings.NewReader(g.root + "\n")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		field := strings.Fields(out.String())
		if len(field) == 0 {
			t.Fatalf("empty cksum output for %q", g.root)
		}
		want, err := strconv.ParseUint(field[0], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		if got := cksum([]byte(g.root + "\n")); uint64(got) != want {
			t.Errorf("cksum(%q) = %d, system tool = %d", g.root, got, want)
		}
	}
}

// TestKeyMatchesPoolPathGolden pins the key the pool path carries, so the pool of
// an existing repository keeps its directory.
func TestKeyMatchesPoolPathGolden(t *testing.T) {
	t.Parallel()
	for _, g := range cksumGolden {
		want := filepath.Base(g.root) + "-" + strconv.FormatUint(uint64(g.sum), 10)
		if got := Key(g.root); got != want {
			t.Errorf("Key(%q) = %q, want %q", g.root, got, want)
		}
	}
}

// TestKeyOfLinkedWorktreeNamesThePrimary proves a call from inside a linked worktree
// keys the primary repository, so one repository has one pool.
func TestKeyOfLinkedWorktreeNamesThePrimary(t *testing.T) {
	t.Parallel()
	primary, linked := repositoryWithLinkedWorktree(t)
	if got, want := Key(linked), Key(primary); got != want {
		t.Errorf("Key(linked) = %q, want the primary key %q", got, want)
	}
	if got, want := Canonical(linked), primary; got != want {
		t.Errorf("Canonical(linked) = %q, want %q", got, want)
	}
}

// repositoryWithLinkedWorktree returns a primary checkout root and the root of a
// linked worktree of it.
func repositoryWithLinkedWorktree(t *testing.T) (string, string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		capability.Capability(t, capability.Tool, "git not available")
	}
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	primary := filepath.Join(base, "primary")
	if err := os.Mkdir(primary, 0o700); err != nil {
		t.Fatal(err)
	}
	run := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run(primary, "init", "--quiet", ".")
	if err := os.WriteFile(filepath.Join(primary, "seed"), []byte("seed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run(primary, "add", "seed")
	run(primary, "commit", "--quiet", "-m", "seed")
	linked := filepath.Join(base, "linked")
	run(primary, "worktree", "add", "--quiet", "--detach", linked)
	return primary, linked
}

// TestPoolJoinsTheHomeAndTheKey pins the pool path, which `bench worktree` and the exec
// census both address and an existing pool on disk already carries.
func TestPoolJoinsTheHomeAndTheKey(t *testing.T) {
	t.Parallel()
	home := filepath.Join(string(filepath.Separator), "home-a", ".bench")
	root := cksumGolden[0].root
	want := filepath.Join(home, "worktrees", Key(root))
	if got := Pool(home, root); got != want {
		t.Errorf("Pool = %q, want %q", got, want)
	}
	if got, wantSuffix := Pool(home, root), Key(root); filepath.Base(got) != wantSuffix {
		t.Errorf("Pool = %q, want the key %q as its last element", got, wantSuffix)
	}
}

// TestPoolSitsDirectlyUnderPools pins that Pool(home, root) sits directly below
// Pools(home), so a caller that tests a command's text for Pools(home) also matches
// every repository's pool without restating the join.
func TestPoolSitsDirectlyUnderPools(t *testing.T) {
	t.Parallel()
	home := filepath.Join(string(filepath.Separator), "home-a", ".bench")
	root := cksumGolden[0].root
	if got, want := filepath.Dir(Pool(home, root)), Pools(home); got != want {
		t.Errorf("Dir(Pool) = %q, want Pools(home) = %q", got, want)
	}
}

// TestAssignmentSegmentRoundTrips proves the writer of a pool directory name and the
// reader of one agree, so a reader never restates the name's shape.
func TestAssignmentSegmentRoundTrips(t *testing.T) {
	t.Parallel()
	owner, id := strings.Repeat("a", 32), strings.Repeat("b", 32)
	got, ok := SplitAssignmentSegment(AssignmentSegment(owner, id))
	if !ok {
		t.Fatalf("SplitAssignmentSegment refused %q", AssignmentSegment(owner, id))
	}
	if got != id {
		t.Errorf("assignment id = %q, want %q", got, id)
	}
}

// TestSplitAssignmentSegmentRejectsOtherNames proves a pool entry that is not an
// assignment is refused, so a reader of the pool needs no ledger.
func TestSplitAssignmentSegmentRejectsOtherNames(t *testing.T) {
	t.Parallel()
	for _, segment := range []string{"scratch", strings.Repeat("a", 32), strings.Repeat("a", 32) + "-" + strings.Repeat("z", 32), ""} {
		if got, ok := SplitAssignmentSegment(segment); ok {
			t.Errorf("SplitAssignmentSegment(%q) = %q, want a refusal", segment, got)
		}
	}
}
