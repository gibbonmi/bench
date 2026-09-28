package releasepreflight

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// An external phase's git starts no auto-maintenance. The race phase runs the kit's
// tests, and a background prune can remove a worktree admin entry a test fixture planted.
func TestExternalPhaseGitStartsNoAutoMaintenance(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	argv := gittest.KitRunProbe(t, os.Environ())
	for i := range argv {
		argv[i] = sanitize.ShellQuote(argv[i])
	}
	cache, err := gocache.Dir(os.Environ())
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "home")
	probe := filepath.Join(t.TempDir(), "phase")
	body := "#!/bin/sh\nset -eu\n" + strings.Join(argv, " ") + "\n" +
		"test \"$GOCACHE\" = " + sanitize.ShellQuote(cache) + "\n" +
		"printf '%s' \"$HOME\" > " + sanitize.ShellQuote(marker) + "\n"
	if err := os.WriteFile(probe, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_PREFLIGHT_RACE", probe)
	version, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		t.Fatal(err)
	}
	toolchain := strings.TrimPrefix(strings.TrimSpace(string(version)), "go")
	var stderr bytes.Buffer
	r := &runner{root: t.TempDir(), stderr: &stderr, identity: Identity{Toolchain: &toolchain}}

	code, err := r.runExternal(context.Background(), "race")

	if err != nil || code != 0 {
		t.Fatalf("maintenance probe = (%d, %v), want (0, nil); stderr=%q", code, err, stderr.String())
	}
	home, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Dir(string(home))); !os.IsNotExist(err) {
		t.Fatalf("private run remains after phase exit: %v", err)
	}
}
