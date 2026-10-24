package preprelease

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/releaseevidence"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

const (
	// rootEnv, groupEnv, and readyEnv name the checkout root, the descendant group, and the
	// readiness file to the resistant step helper.
	rootEnv  = "BENCH_PREPRELEASE_TEST_ROOT"
	groupEnv = "BENCH_PREPRELEASE_TEST_GROUP"
	readyEnv = "BENCH_PREPRELEASE_TEST_READY"
	// sentinel is the evidence file that the resistant step promotes.
	sentinel = "resistant-child"
)

// TestMain installs the release-evidence row of the composition root's table, so prep-release
// protects its evidence as the CLI does. The test owns the only vector, so an outer run's
// resources stay out of its launches. A step helper restores the vector that its launch gave
// it, because a step is a user of the evidence that the test process owns.
func TestMain(m *testing.M) {
	vector := os.Getenv(protection.Env)
	if err := os.Unsetenv(protection.Env); err != nil {
		os.Exit(3)
	}
	protection.Install(protection.Protocol{Validators: protection.Validators{protection.KindReleaseEvidence: releaseevidence.ValidateProtection}, ReadLimit: bounds.ControlRecordLimit})
	processtest.Main(m, map[string]func() int{"resistant-step": func() int {
		if os.Setenv(protection.Env, vector) != nil {
			return 3
		}
		return resistantStep()
	}})
}

// resistantStep is a nested preflight that ignores TERM. It promotes evidence from inside its
// registered launch, registers a descendant whose separate group groupEnv names, and then
// waits for KILL. That group outlives the step's shutdown.
func resistantStep() int {
	signal.Ignore(syscall.SIGTERM)
	if err := releaseevidence.PromoteEvidenceFiles(os.Getenv(rootEnv), releaseevidence.ModeVerify, nil, releaseevidence.Manifest{SchemaVersion: 1}, map[string][]byte{sentinel: []byte("held\n")}); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		return 1
	}
	group, err := strconv.Atoi(os.Getenv(groupEnv))
	if err != nil {
		return 2
	}
	inherited, err := protection.Inherited()
	if err != nil {
		return 2
	}
	descendant, err := protection.Installed().Prepare(inherited)
	if err != nil || descendant.Started(group, group) != nil || os.WriteFile(os.Getenv(readyEnv), nil, 0o600) != nil {
		return 2
	}
	for {
		time.Sleep(bounds.FixedWindow(time.Hour))
	}
}

// TestPrepIncomplete is PL61. A real prep-release step promotes evidence, and its descendant
// group outlives the TERM, the grace, and the KILL of the step. The run refuses ship green with
// retained evidence, and a later promotion cannot replace that evidence while the descendant
// stays unresolved.
func TestPrepIncomplete(t *testing.T) {
	root := t.TempDir()
	if err := releaseevidence.PromoteEvidence(root, releaseevidence.ModeVerify, nil, releaseevidence.Manifest{SchemaVersion: 1}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { discardEvidenceStore(t, root) })
	ready := filepath.Join(t.TempDir(), "ready")
	// The test's own group is a live group, so it stands for the descendant's group.
	step := processtest.Helper(t, "resistant-step", rootEnv+"="+root, groupEnv+"="+strconv.Itoa(syscall.Getpgrp()), readyEnv+"="+ready)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stdout, stderr bytes.Buffer
	done := make(chan int, 1)
	go func() {
		done <- run(ctx, root, []Step{{Name: "resistant", Argv: step.Args, Env: step.Env}}, &stdout, &stderr)
	}()
	tick := time.NewTicker(bounds.FixedWindow(10 * time.Millisecond))
	defer tick.Stop()
	code := -1
	for code < 0 {
		select {
		case code = <-done:
		case <-tick.C:
			if _, err := os.Stat(ready); err == nil {
				cancel()
			}
		}
	}
	if code != 1 || stdout.Len() != 0 {
		t.Fatalf("run = exit %d, stdout %q, stderr %q; want exit 1 and no ship green", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{"prep-release: resistant failed:", "signal: killed", protection.ErrRetained.Error()} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("stderr = %q, want %q", stderr.String(), want)
		}
	}
	requireSentinel(t, root)

	err := releaseevidence.PromoteEvidence(root, releaseevidence.ModeVerify, nil, releaseevidence.Manifest{SchemaVersion: 1})
	if !errors.Is(err, protection.ErrSealed) {
		t.Fatalf("later promotion = %v, want the sealed-generation refusal", err)
	}
	requireSentinel(t, root)
}

// TestPrepStepZeroUncertified pins that a step's exit zero cannot certify the step while a
// member of its group outlives the leader. The step fails on its own protection.
func TestPrepStepZeroUncertified(t *testing.T) {
	root := t.TempDir()
	t.Cleanup(func() { discardEvidenceStore(t, root) })
	pidFile := filepath.Join(t.TempDir(), "pid")
	step := Step{Name: "lingering", Argv: []string{"sh", "-c", `sleep 300 >/dev/null 2>&1 & echo $! >"$0"`, pidFile}}

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), root, []Step{step}, &stdout, &stderr)
	processtest.WaitPID(t, pidFile)
	if code != 1 || stdout.Len() != 0 {
		t.Fatalf("run = exit %d, stdout %q, stderr %q; want exit 1 and no ship green", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{"prep-release: lingering failed:", "process protection " + string(protection.Protected)} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("stderr = %q, want %q", stderr.String(), want)
		}
	}
}

// TestPrepRefusesEarlierUser pins that an unresolved user of an open generation that an earlier
// run left refuses the next run before its first step.
func TestPrepRefusesEarlierUser(t *testing.T) {
	root := t.TempDir()
	earlier, err := releaseevidence.Protect(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { discardEvidenceStore(t, root) })
	user, err := protection.Installed().Prepare(earlier.Resources())
	// The test's own group is a live group, so it stands for the earlier run's descendant.
	if err != nil || user.Started(os.Getpid(), syscall.Getpgrp()) != nil {
		t.Fatalf("register the earlier user: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), root, []Step{{Name: "later", Argv: []string{"true"}}}, &stdout, &stderr)
	if code != 1 || stdout.Len() != 0 || !strings.Contains(stderr.String(), protection.ErrRetained.Error()) || strings.Contains(stderr.String(), "prep-release: later") {
		t.Fatalf("run = exit %d, stdout %q, stderr %q; want exit 1 with the retention before any step", code, stdout.String(), stderr.String())
	}
}

// requireSentinel fails unless the resistant step's promoted evidence is the current target.
func requireSentinel(t *testing.T, root string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(EvidenceDir), sentinel)); err != nil {
		t.Fatalf("resistant step evidence = %v, want it retained", err)
	}
}

// discardEvidenceStore removes the protection store of root's evidence, which lives under the
// shared store root and outlives the temporary checkout.
func discardEvidenceStore(t *testing.T, root string) {
	t.Helper()
	d, err := releaseevidence.Descriptor(root)
	if err == nil {
		err = protection.Generation{Resource: d}.Discard()
	}
	if err != nil {
		t.Error(err)
	}
}
