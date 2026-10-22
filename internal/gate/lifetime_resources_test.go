package gate

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/testrepo"
)

const (
	kitRunGateHelper = "kit-run-gate"
	kitRunUserHelper = "kit-run-user"
	cacheGateHelper  = "cache-gate"
	kitRunRootEnv    = "BENCH_TEST_KIT_RUN_ROOT"
)

// TestMain lets this test binary serve as the gates and the nested user below.
func TestMain(m *testing.M) {
	processtest.Main(m, map[string]func() int{kitRunGateHelper: kitRunGate, kitRunUserHelper: kitRunUser, cacheGateHelper: cacheGate})
}

// cacheGate is a fresh gate CLI whose oracle registers on the cache generation.
func cacheGate() int {
	protection.Install(protection.Protocol{Validators: protection.Validators{protection.KindCache: gocache.ValidateProtection}, ReadLimit: bounds.ControlRecordLimit})
	return RunCommand([]string{"--fresh", os.Getenv(kitRunRootEnv)}, os.Stdout, os.Stderr)
}

// TestGateTerminalFailure is PL78. The oracle exits zero after it denies writes to the cache
// store, so its terminal protection cannot publish. The pending record stays, with both causes.
func TestGateTerminalFailure(t *testing.T) {
	requireDirectoryWriteDenied(t)
	home := t.TempDir()
	dir, err := gocache.Dir([]string{"HOME=" + home})
	if err != nil {
		t.Fatal(err)
	}
	store := filepath.Join(dir, gocache.ProtectionStore)
	root := outcomeFixture(t, func(f *testrepo.GateFixture, body string) string {
		f.Environment = []string{"HOME"}
		return body + f.Command("chmod") + " -R a-w " + sanitize.ShellQuote(store) + "\n"
	})
	t.Cleanup(func() {
		_ = filepath.WalkDir(store, func(path string, _ fs.DirEntry, _ error) error { return os.Chmod(path, 0o700) })
	})
	gate := processtest.Helper(t, cacheGateHelper, kitRunRootEnv+"="+root, "HOME="+home)
	gate.Env = protection.Environ(gate.Env, nil)
	out, err := gate.CombinedOutput()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 1 {
		t.Fatalf("gate with an unwritable cache store = %v; output:\n%s\nwant exit 1", err, out)
	}
	for _, cause := range []string{"process exit 0", "process protection uncertain"} {
		if !strings.Contains(string(out), cause) {
			t.Fatalf("gate output lacks %q:\n%s", cause, out)
		}
	}
	if inspection := Inspect(root); inspection.State != Pending || inspection.ReusableGreen {
		t.Fatalf("inspection after the failed terminal protection = %#v, want pending and non-reusable", inspection)
	}
}

// installKitRunProtocol admits the private test run alone, as the composition root's kit-test
// row does. Only a helper process installs it.
func installKitRunProtocol() {
	protection.Install(protection.Protocol{Validators: protection.Validators{protection.KindKitTest: env.ValidateProtection}, ReadLimit: bounds.ControlRecordLimit})
}

// kitRunGate runs the fixture phases of its root through the gate's own kit-run composition.
func kitRunGate() int {
	installKitRunProtocol()
	root := os.Getenv(kitRunRootEnv)
	return phasesCommandAtKitWithSelection(context.Background(), root, root, fixtureSelection(root), os.Stdout, os.Stderr)
}

// kitRunUser is a nested CLI in the phase. It writes into the private home, registers on its
// inherited vector, and exits zero before it resolves its startup identity.
func kitRunUser() int {
	installKitRunProtocol()
	inherited, err := protection.Inherited()
	if err != nil || len(inherited) == 0 || os.WriteFile(filepath.Join(os.Getenv("HOME"), "owned"), []byte("held"), 0o600) != nil {
		return 2
	}
	if _, err := protection.Installed().Prepare(inherited); err != nil {
		return 2
	}
	return 0
}

// TestPhaseCloseError is PL38. The phase exits zero, but its nested user keeps the private run
// protected, so the gate must consume the Close failure and keep the file.
func TestPhaseCloseError(t *testing.T) {
	user := processtest.Helper(t, kitRunUserHelper)
	phaseEnv := map[string]string{}
	for _, entry := range user.Env[len(os.Environ()):] {
		name, value, _ := strings.Cut(entry, "=")
		phaseEnv[name] = value
	}
	phase, err := json.Marshal(map[string]any{"name": "nested-user", "argv": []string{user.Path}, "env": phaseEnv})
	if err != nil {
		t.Fatal(err)
	}
	root := fixturePhaseRoot(t, `{"phases":[`+string(phase)+`]}`)
	tmp := t.TempDir()
	gate := processtest.Helper(t, kitRunGateHelper, kitRunRootEnv+"="+root, "TMPDIR="+tmp)
	// The run under test owns the only vector, so an outer run's resources stay out of it.
	gate.Env = protection.Environ(gate.Env, nil)
	out, err := gate.CombinedOutput()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 1 || !strings.Contains(string(out), "stays retained") {
		t.Fatalf("gate with a protected private run = %v; output:\n%s\nwant exit 1 with the retained run", err, out)
	}
	owned, err := filepath.Glob(filepath.Join(tmp, "*", "h", "owned"))
	if err != nil || len(owned) != 1 {
		t.Fatalf("child-owned files after the gate = %v, %v; want the one retained file", owned, err)
	}
}
