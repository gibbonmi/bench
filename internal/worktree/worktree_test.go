package worktree

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// cksumGolden pins the Go cksum against values produced by the coreutils `cksum`
// tool. Each was derived once with, e.g.:
//
//	printf '%s\n' "/home/mgibs/workspace/bench" | cksum   -> 2826441890 28
//	printf '%s\n' "/tmp/a b/c"                  | cksum   -> 889650394  11
//
// The `\n` is intentional: Pool checksums `root + "\n"` because the shell used
// `echo "$root" | cksum`. The second vector carries a space to exercise a path the
// shell would otherwise word-split.
var cksumGolden = []struct {
	root string
	sum  uint32
}{
	{"/home/mgibs/workspace/bench", 2826441890},
	{"/tmp/a b/c", 889650394},
}

func TestCksumMatchesGolden(t *testing.T) {
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
	if _, err := exec.LookPath("cksum"); err != nil {
		t.Skip("cksum not available")
	}
	for _, g := range cksumGolden {
		// printf '%s\n' "<root>" | cksum
		printf := exec.Command("printf", "%s\n", g.root)
		ck := exec.Command("cksum")
		pipe, err := printf.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		ck.Stdin = pipe
		out, err := ck.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		if err := ck.Start(); err != nil {
			t.Fatal(err)
		}
		if err := printf.Run(); err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, 64)
		n, _ := out.Read(buf)
		if err := ck.Wait(); err != nil {
			t.Fatal(err)
		}
		field := strings.Fields(string(buf[:n]))
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

func TestPool(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := "/home/mgibs/workspace/bench"
	want := filepath.Join(home, "worktrees", "bench-2826441890")
	if got := Pool(root); got != want {
		t.Errorf("Pool(%q) = %q, want %q", root, got, want)
	}
}

func TestPoolDefaultBenchHome(t *testing.T) {
	// With BENCH_HOME unset, Pool falls back to <home>/.bench.
	t.Setenv("BENCH_HOME", "")
	root := "/tmp/a b/c"
	got := Pool(root)
	suffix := filepath.Join(".bench", "worktrees", "c-889650394")
	if !strings.HasSuffix(got, suffix) {
		t.Errorf("Pool(%q) = %q, want suffix %q", root, got, suffix)
	}
}

func TestLeaseFile(t *testing.T) {
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	lease, err := LeaseFile(dir)
	if err != nil {
		t.Fatalf("LeaseFile: %v", err)
	}
	if !strings.HasSuffix(lease, "bench-lease") {
		t.Errorf("LeaseFile = %q, want suffix bench-lease", lease)
	}
	if !filepath.IsAbs(lease) {
		t.Errorf("LeaseFile = %q, want absolute — a relative path resolves against the caller's CWD, not the worktree", lease)
	}
}

func TestLeaseFileCommandMissingArg(t *testing.T) {
	out, code := LeaseFileCommand(nil)
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if !strings.HasPrefix(out, "usage:") {
		t.Errorf("out = %q, want usage line", out)
	}
}

func TestPoolCommandExplicitRoot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	out, code := PoolCommand([]string{"/home/mgibs/workspace/bench"})
	if code != 0 {
		t.Errorf("exit = %d, want 0", code)
	}
	want := filepath.Join(home, "worktrees", "bench-2826441890") + "\n"
	if out != want {
		t.Errorf("out = %q, want %q", out, want)
	}
}
