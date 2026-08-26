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
