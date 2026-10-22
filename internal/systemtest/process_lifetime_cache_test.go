//go:build system

package systemtest

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// TestProcessLifetimeCacheSecondProcessClean is PL55. The gate CLI dies while its oracle
// group lives on, so the shared cache lock drops with the CLI. A second-process clean must
// refuse the unresolved user, and it succeeds once that group is gone.
func TestProcessLifetimeCacheSecondProcessClean(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	dir, err := gocache.Dir([]string{"HOME=" + home})
	if err != nil {
		t.Fatal(err)
	}
	shard := filepath.Join(dir, "aa", "aa0123456789-a")
	if err := os.MkdirAll(filepath.Dir(shard), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(shard, []byte("archive"), 0o644); err != nil {
		t.Fatal(err)
	}
	oracle := filepath.Join(t.TempDir(), "oracle")
	root := gittest.RepoOnBranch(t, "main")
	f := testrepo.NewGateFixture(t.TempDir(), "HOME")
	f.MustWrite(t, root, `echo "$$" > `+sanitize.ShellQuote(oracle+".partial")+` && `+f.Command("mv")+` `+sanitize.ShellQuote(oracle+".partial")+` `+sanitize.ShellQuote(oracle)+`
exec `+f.Command("sleep")+` 300
`, "")
	gittest.Output(t, root, "add", ".")
	gittest.Output(t, root, "commit", "-qm", "fixture")
	env := []string{"HOME=" + home, "TMPDIR=" + t.TempDir(), "TEST_TELEMETRY_DIR=" + os.DevNull, protection.Env}

	gate, _, stderr := systemStartSelected(t, root, env, "gate", "--fresh")
	pgid := processtest.WaitPID(t, oracle)
	if err := gate.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = gate.Wait()
	if refused := systemSelected(t, root, env, "cache", "clean"); refused.code != 1 || !strings.Contains(refused.stdout+refused.stderr, "cache protected") {
		t.Fatalf("clean after the gate CLI died = (%d, %q, %q), want the protection refusal; gate stderr %q", refused.code, refused.stdout, refused.stderr, stderr)
	}
	if _, err := os.Stat(shard); err != nil {
		t.Fatalf("shard after the refused clean = %v, want it intact", err)
	}
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	processtest.WaitAbsent(t, pgid)
	if cleaned := systemSelected(t, root, env, "cache", "clean"); cleaned.code != 0 {
		t.Fatalf("clean after the oracle group ended = (%d, %q, %q), want exit 0", cleaned.code, cleaned.stdout, cleaned.stderr)
	}
	if _, err := os.Stat(shard); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("shard after the clean = %v, want it removed", err)
	}
}
