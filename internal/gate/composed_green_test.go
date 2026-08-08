package gate

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

func TestComposedGreenReadsVerdictFromLinkedWorktreeGitDir(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "seed")

	linked := filepath.Join(t.TempDir(), "linked")
	gitRun(t, root, "worktree", "add", "-q", "--detach", linked, "HEAD")
	t.Cleanup(func() { gitRun(t, root, "worktree", "remove", "--force", linked) })

	plan, err := buildSubject(linked)
	if err != nil || !plan.Closed {
		t.Fatalf("linked subject = %+v, %v; want closed", plan, err)
	}
	linkedGitDir := gitOutput(t, linked, "rev-parse", "--absolute-git-dir")
	commonGitDir := commonGitDirOf(t, linked)
	if linkedGitDir == commonGitDir {
		t.Fatalf("linked git dir = common git dir %q, want checkout-local cache directory", linkedGitDir)
	}

	record := verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: time.Now().UTC().Format(time.RFC3339)}
	if err := durableReplace(linkedGitDir, record); err != nil {
		t.Fatal(err)
	}
	evidenceDir := filepath.Join(commonGitDir, "bench-gate-evidence")
	if err := os.MkdirAll(evidenceDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := durableReplaceAt(evidenceDir, evidenceName(plan), record); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(commonGitDir, benchgit.GateCacheFile)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("common cache = %v, want absent while linked cache is authoritative", err)
	}
	if _, err := os.Stat(evidencePath(commonGitDir, plan)); err != nil {
		t.Fatalf("common retained evidence = %v, want present", err)
	}

	if !ComposedGreen(linked) {
		t.Fatal("linked checkout green verdict was not composed")
	}
}
