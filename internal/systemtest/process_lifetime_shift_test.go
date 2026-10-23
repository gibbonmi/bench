//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// TestProcessLifetimeCleanProtectedRecovery is PL85. The first shift CLI dies while its
// registered adapter group lives on, and its worktree holds only scratch. A second shift CLI
// runs the real recovery owner before its own acquire. The recovery must retain the worktree
// with its scratch, because an empty dirty set is no proof that the adapter group cannot
// change the worktree.
func TestProcessLifetimeCleanProtectedRecovery(t *testing.T) {
	t.Parallel()
	fixture := scaffoldRecordedPublicationRepo(t)
	commitmenttest.Register(t, fixture.root, "shift-protected")
	commitmenttest.Admit(t, fixture.root, "shift-protected", "specs/record/spec.md")
	marks := t.TempDir()
	mark, adapter := filepath.Join(marks, "adapter.pid"), filepath.Join(marks, "agent")
	// The adapter changes no worktree file, and it leaves the CLI's streams, so only its
	// group keeps the worktree in use after the CLI dies. It names its worktree, because the
	// killed CLI publishes no bounded response.
	script := "#!/bin/sh\ncat >/dev/null\npwd -P > '" + mark + ".dir'\necho $$ > '" + mark + ".partial' && mv '" + mark + ".partial' '" + mark + "'\nexec sleep 300 </dev/null >/dev/null 2>&1\n"
	if err := os.WriteFile(adapter, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	first, stdout, stderr := systemStartSelected(t, fixture.root, append(fixture.environment(fixture.home),
		"BENCH_AGENT="+adapter, "BENCH_MAX_ITERS=1", "TMPDIR="+t.TempDir()), "shift", "--outcome", "delivery", "protected scratch")
	processtest.WaitPID(t, mark)
	if err := first.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = exitWithin(t, first, 0)
	dir, err := os.ReadFile(mark + ".dir")
	if err != nil {
		t.Fatal(err)
	}
	worktree := strings.TrimSpace(string(dir))
	scratch := filepath.Join(worktree, ".bench-objective")
	if _, err := os.Stat(scratch); worktree == "" || err != nil {
		t.Fatalf("the killed shift left worktree %q with scratch %v; stdout %q stderr %q", worktree, err, stdout.String(), stderr.String())
	}

	second := systemSelected(t, fixture.root, append(fixture.environment(fixture.home),
		"BENCH_AGENT=true", "BENCH_MAX_ITERS=1", "TMPDIR="+t.TempDir()), "shift", "--outcome", "delivery", "recovering shift")
	if !strings.Contains(second.stdout, "bench shift recovery: recovered 1") {
		t.Fatalf("second shift = (%d, %q, %q), want the recovery of the killed shift", second.code, second.stdout, second.stderr)
	}
	if _, err := os.Stat(scratch); err != nil {
		t.Fatalf("the recovery removed the protected scratch: %v", err)
	}
	listing := systemGitOutput(t, fixture.root, "worktree", "list", "--porcelain")
	at := strings.Index(listing, "worktree "+worktree+"\n")
	if at < 0 {
		t.Fatalf("git lists no worktree %s after the recovery: %q", worktree, listing)
	}
	if block, _, _ := strings.Cut(listing[at:], "\n\n"); !strings.Contains(block, "\nlocked") {
		t.Fatalf("worktree %s after the recovery = %q, want it retained and locked", worktree, block)
	}
}
