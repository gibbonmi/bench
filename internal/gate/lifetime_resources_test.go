package gate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gate/prospectiveartifact"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/testrepo"
)

const (
	kitRunGateHelper = "kit-run-gate"
	kitRunUserHelper = "kit-run-user"
	cacheGateHelper  = "cache-gate"
	laneGateHelper   = "lane-gate"
	nestedUserHelper = "nested-user"
	kitRunRootEnv    = "BENCH_TEST_KIT_RUN_ROOT"
	nestedGroupEnv   = "BENCH_TEST_NESTED_GROUP"
)

// TestMain lets this test binary serve as the gates and the nested users below, and as the
// background grade.
func TestMain(m *testing.M) {
	processtest.Main(m, map[string]func() int{kitRunGateHelper: kitRunGate, kitRunUserHelper: kitRunUser, cacheGateHelper: cacheGate, laneGateHelper: laneGate, nestedUserHelper: nestedUser, backgroundGradeHelper: backgroundGradeChild})
}

// installTestProtocol admits each kind the composition root's table admits. Only a helper
// process installs it.
func installTestProtocol() {
	protection.Install(protection.Protocol{Validators: protection.Validators{
		protection.KindCache:       gocache.ValidateProtection,
		protection.KindProspective: prospectiveartifact.ValidateProtection,
		protection.KindBinary:      runbinary.ValidateProtection,
		protection.KindKitTest:     env.ValidateProtection,
	}, ReadLimit: bounds.ControlRecordLimit})
}

// cacheGate is a fresh gate CLI whose oracle registers on the cache generation.
func cacheGate() int {
	installTestProtocol()
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
	// The store derives from the cache's identity, so the cache exists before the gate holds it.
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	store, err := gocache.ProtectionStore(dir)
	if err != nil {
		t.Fatal(err)
	}
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

// kitRunGate runs the fixture phases of its root through the gate's own kit-run composition.
func kitRunGate() int {
	installTestProtocol()
	root := os.Getenv(kitRunRootEnv)
	return phasesCommandAtKitWithSelection(context.Background(), root, root, fixtureSelection(root), os.Stdout, os.Stderr)
}

// kitRunUser is a nested user in the phase that first writes into the private home.
func kitRunUser() int {
	if os.WriteFile(filepath.Join(os.Getenv("HOME"), "owned"), []byte("held"), 0o600) != nil {
		return 2
	}
	return nestedUser()
}

// nestedUser is a nested CLI. It registers on its inherited vector, publishes the live group
// that nestedGroupEnv names as its started identity, and exits zero without a terminal record,
// so that group keeps the resources protected.
func nestedUser() int {
	installTestProtocol()
	inherited, err := protection.Inherited()
	if err != nil || len(inherited) == 0 {
		return 2
	}
	group, err := strconv.Atoi(os.Getenv(nestedGroupEnv))
	if err != nil {
		return 2
	}
	reg, err := protection.Installed().Prepare(inherited)
	if err != nil || reg.Started(group, group) != nil {
		return 2
	}
	return 0
}

// nestedHelper prepares the named nested-user helper. The test's own group is a live group,
// so it stands for the group that the nested user started.
func nestedHelper(t *testing.T, helper string) *exec.Cmd {
	t.Helper()
	return processtest.Helper(t, helper, nestedGroupEnv+"="+strconv.Itoa(syscall.Getpgrp()))
}

// userPhase declares one check or phase that runs the named nested-user helper.
func userPhase(t *testing.T, helper string) map[string]any {
	t.Helper()
	user := nestedHelper(t, helper)
	phaseEnv := map[string]string{}
	for _, entry := range user.Env[len(os.Environ()):] {
		name, value, _ := strings.Cut(entry, "=")
		phaseEnv[name] = value
	}
	return map[string]any{"name": helper, "argv": []string{user.Path}, "env": phaseEnv}
}

// TestPhaseCloseError is PL38. The phase exits zero, but its nested user keeps the private run
// protected, so the gate must consume the Close failure and keep the file.
func TestPhaseCloseError(t *testing.T) {
	phase, err := json.Marshal(userPhase(t, kitRunUserHelper))
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

// laneGate runs its root's declared lane on the root's tree. The stub build publishes no seal,
// so the lane's selection skips freshness, and an inherited selection would bypass the build.
func laneGate() int {
	installTestProtocol()
	laneRunBinary.Verify = func(string, string) error { return nil }
	root := os.Getenv(kitRunRootEnv)
	checks, err := LaneFor(root, "")
	if err == nil {
		err = os.Unsetenv(runbinary.Env)
	}
	if err == nil {
		_, err = RunLane(context.Background(), LaneRequest{Root: root, Tree: benchgit.TreeHash(root), Checks: checks, Stdout: os.Stdout, Stderr: os.Stderr})
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// requireLaneRetained runs the lane of root and requires the refusal that names the retained
// bundle, no lane record, and a retained path that matches pattern under the lane's TMPDIR.
func requireLaneRetained(t *testing.T, root, pattern string) {
	t.Helper()
	tmp := t.TempDir()
	lane := processtest.Helper(t, laneGateHelper, kitRunRootEnv+"="+root, "HOME="+t.TempDir(), "TMPDIR="+tmp)
	lane.Env = protection.Environ(lane.Env, nil)
	out, err := lane.CombinedOutput()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 1 || !strings.Contains(string(out), "stays retained") {
		t.Fatalf("lane with a nested user = %v; output:\n%s\nwant exit 1 with the retained bundle", err, out)
	}
	if _, err := os.Stat(filepath.Join(gittest.Output(t, root, "rev-parse", "--absolute-git-dir"), laneRecordFile)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("lane record after the retained release = %v, want none", err)
	}
	if retained, err := filepath.Glob(filepath.Join(tmp, prospectiveartifact.BundlePrefix+"*", pattern)); err != nil || len(retained) != 1 {
		t.Fatalf("retained %s = %v, %v; want one", pattern, retained, err)
	}
}

// laneRoot commits a repository whose phase manifest declares lane, beside files.
func laneRoot(t *testing.T, files map[string]string, lane ...map[string]any) string {
	t.Helper()
	root := outcomeFixture(t)
	manifest, err := json.Marshal(map[string]any{"lane": lane})
	if err != nil {
		t.Fatal(err)
	}
	files[canary.PhaseManifestPath] = string(manifest)
	for path, body := range files {
		outcomeWrite(t, root, path, body, 0o755)
	}
	outcomeCommit(t, root, "lane")
	return root
}

// TestLaneIncomplete is PL33. The check exits zero, but its nested user keeps the bundle
// protected, so the lane must consume the bundle release and record no pass.
func TestLaneIncomplete(t *testing.T) {
	root := laneRoot(t, map[string]string{}, userPhase(t, nestedUserHelper))
	requireLaneRetained(t, root, prospectiveartifact.CheckoutName)
}

// nestedUserCommand is the shell command that runs the nested-user helper.
func nestedUserCommand(t *testing.T) string {
	t.Helper()
	user := nestedHelper(t, nestedUserHelper)
	nested := []string{"env"}
	for _, entry := range append(user.Env[len(os.Environ()):], user.Path) {
		nested = append(nested, sanitize.ShellQuote(entry))
	}
	return strings.Join(nested, " ")
}

// selectedLane commits a lane whose one check runs the selected executable that build publishes.
func selectedLane(t *testing.T, build string) string {
	t.Helper()
	return laneRoot(t, map[string]string{"scripts/go-build.sh": "#!/bin/sh\nshift $(($# - 2))\n" + build}, map[string]any{"name": "selected", "argv": []string{runBinaryArgvToken}})
}

// TestLaneSelectionCloseFailure is PL123. The selected build leaves a nested user on the
// selection alone, so only the selection release reports it, and the lane must consume it.
func TestLaneSelectionCloseFailure(t *testing.T) {
	root := selectedLane(t, nestedUserCommand(t)+"\nprintf '#!/bin/sh\\n' > \"$2\"\nchmod 0700 \"$2\"\n")
	requireLaneRetained(t, root, filepath.Join("bench-run-*", "bench"))
}

// TestLaneSelectedExecutableUser is story 28. The check runs the selected executable, which
// leaves a nested user and exits zero, so the selection must retain the executable that user
// can still run.
func TestLaneSelectedExecutableUser(t *testing.T) {
	executable := filepath.Join(t.TempDir(), "selected")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nexec "+nestedUserCommand(t)+"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	root := selectedLane(t, "cp "+sanitize.ShellQuote(executable)+" \"$2\"\n")
	requireLaneRetained(t, root, filepath.Join("bench-run-*", "bench"))
}
