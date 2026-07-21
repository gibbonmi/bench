package shift

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoopRefreshResolvesStartAfterFetch(t *testing.T) {
	root, _ := shiftCollisionFixture(t)
	remote := filepath.Join(t.TempDir(), "origin.git")
	run := func(dir string, args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	run(root, "init", "-q", "--bare", remote)
	branch := run(root, "branch", "--show-current")
	run(root, "remote", "add", "origin", remote)
	run(root, "push", "-q", "-u", "origin", branch)
	run(root, "remote", "set-head", "origin", "-a")
	publisher := filepath.Join(t.TempDir(), "publisher")
	run(root, "clone", "-q", remote, publisher)
	if err := os.WriteFile(filepath.Join(publisher, "remote-only"), []byte("refreshed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(publisher, "-c", "user.name=bench", "-c", "user.email=bench@local", "add", "remote-only")
	run(publisher, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "remote advance")
	run(publisher, "push", "-q", "origin", branch)
	remoteHead := run(publisher, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	if code := loop("refresh start ref", true, &stdout, &stderr); code != 4 {
		t.Fatalf("loop = %d, want no-op/4; stdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "log --oneline "+remoteHead+"..") {
		t.Fatalf("refresh did not select fetched start %s:\n%s", remoteHead, stdout.String())
	}
}
