//go:build system

package systemtest

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/gate/prospectiveartifact"
)

func TestProspectiveArtifactRecoveryAfterKilledLanding(t *testing.T) {
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	configureArtifactLandingFixture(t, root)
	private := t.TempDir()
	t.Setenv("TMPDIR", private)
	first := systemCreateLandingWorktree(t, root, home, "artifact-dead", "artifact dead")
	base := systemGitOutput(t, root, "rev-parse", "main")
	systemCommit(t, first.path, "loser.txt", "dead\n", "dead prospective subject")
	first.tip = systemGitOutput(t, first.path, "rev-parse", "HEAD")
	if review := systemSelected(t, first.path, systemLandEnv(root, home, tally, trees, ready, release), "preflight", "review", "x", "--base", base); review.code != 0 {
		t.Fatalf("first review = (%d, %q, %q)", review.code, review.stdout, review.stderr)
	}

	cmd, childOut, childErr, done := startArtifactLand(t, root, home, tally, trees, ready, release, first, base)
	awaitArtifactBarrier(t, ready, done, childOut, childErr)

	bundle := oneProspectiveBundle(t, private)
	requirePublishedArtifactOwnerRecord(t, bundle, cmd.Process.Pid, root)
	checkout := filepath.Join(bundle, prospectiveartifact.CheckoutName)
	oldBinary := oneOwnerBinary(t, bundle)
	marker := filepath.Join(private, "old-candidate-ran")
	plantedBinary := filepath.Join(bundle, "bench-run-planted", "bench")
	if err := os.MkdirAll(filepath.Dir(plantedBinary), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(plantedBinary, []byte("#!/bin/sh\nprintf old > "+marker+"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	if code := systemExitCode(<-done); code == 0 {
		t.Fatal("killed prospective owner exited green")
	}
	releaseArtifactBarrier(t, release)

	fresh := systemCreateLandingWorktree(t, root, home, "artifact-fresh", "artifact fresh")
	systemCommit(t, fresh.path, "winner.txt", "fresh\n", "fresh prospective subject")
	fresh.tip = systemGitOutput(t, fresh.path, "rev-parse", "HEAD")
	if review := systemSelected(t, fresh.path, systemLandEnv(root, home, tally, trees, ready, release), "preflight", "review", "x", "--base", base); review.code != 0 {
		t.Fatalf("fresh review = (%d, %q, %q)", review.code, review.stdout, review.stderr)
	}
	_, freshOut, freshErr, freshDone := startArtifactLand(t, root, home, tally, trees, ready, release, fresh, base)
	awaitArtifactBarrier(t, ready, freshDone, freshOut, freshErr)
	if _, err := os.Stat(bundle); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("dead bundle = %v, want absent", err)
	}
	if _, err := os.Stat(oldBinary); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("old owner binary = %v, want absent", err)
	}
	if _, err := os.Stat(plantedBinary); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("planted old candidate = %v, want absent", err)
	}
	if got := systemGitOutput(t, root, "worktree", "list", "--porcelain"); strings.Contains(got, checkout) {
		t.Fatalf("dead checkout registration %q survived", checkout)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("recovery executed old candidate bytes: %v", err)
	}
	systemGit(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "--allow-empty", "-qm", "same-tree destination move")
	releaseArtifactBarrier(t, release)
	if code := systemExitCode(<-freshDone); code != 1 || !strings.Contains(freshOut.String(), "landing destination checkout changed") {
		t.Fatalf("same-tree moved destination = (%d, %q, %q)", code, freshOut.String(), freshErr.String())
	}
	gateTally, err := os.ReadFile(tally)
	if err != nil || string(gateTally) != "lw" {
		t.Fatalf("recovery gate tally = %q, %v, want one killed and one green run", gateTally, err)
	}
	// PAR04: the second fresh-process sweep of the removed bundle changes no path, so the
	// row records the shared temporary root before it runs and asks for it back byte for
	// byte. The landing moves its own worktrees either way, so the registration half asks
	// only that no registration names a prospective bundle.
	plantedRoot := artifactRootEntries(t, private)
	if result := systemSelected(t, root, artifactLandEnv(root, home, tally, trees, ready, release), systemLandArgs(fresh, base)...); result.code != 0 {
		t.Fatalf("second prospective authorization = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}
	if repeatedTally, err := os.ReadFile(tally); err != nil || string(repeatedTally) != string(gateTally) {
		t.Fatalf("reused exact-green tally = %q, %v, want unchanged %q", repeatedTally, err, gateTally)
	}
	if after := artifactRootEntries(t, private); after != plantedRoot {
		t.Fatalf("temporary root after the second sweep =\n%s\nwant\n%s", after, plantedRoot)
	}
	if registrations := systemGitOutput(t, root, "worktree", "list", "--porcelain"); strings.Contains(registrations, prospectiveartifact.BundlePrefix) {
		t.Fatalf("a prospective checkout registration survived the second sweep:\n%s", registrations)
	}
	owner.markTerminal("green")
}

func configureArtifactLandingFixture(t *testing.T, root string) {
	t.Helper()
	gate := "#!/bin/sh\nset -eu\nroot=${1:-$(git rev-parse --show-toplevel)}\nkit=${BENCH_KIT:?}\nbench=${BENCH_RUN_BINARY:?}\nexec env BENCH_KIT=\"$kit\" BENCH_RUN_BINARY=\"$bench\" \"$bench\" gate-phases \"$root\"\n"
	for _, file := range []string{"gate.sh", "gate-prospective.sh"} {
		if err := os.WriteFile(filepath.Join(root, ".bench", file), []byte(gate), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	phase := "#!/bin/sh\nset -eu\nruntime=$(git rev-parse --show-toplevel)\ngrep -q '^Status: implemented$' \"$runtime/specs/x/spec.md\"\ntree=$(git -C \"$runtime\" write-tree)\nprintf '%s\\n' \"$tree\" >> \"$LAND_GATE_TREES\"\nif [ -f \"$runtime/loser.txt\" ]; then\n  printf l >> \"$LAND_GATE_TALLY\"\nelse\n  printf w >> \"$LAND_GATE_TALLY\"\nfi\nprintf r > \"$LAND_RACE_READY\"\nIFS= read -r _ < \"$LAND_RACE_RELEASE\"\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "landing-race-phase.sh"), []byte(phase), 0o755); err != nil {
		t.Fatal(err)
	}
	phases := "{\"phases\":[{\"name\":\"landing-race\",\"argv\":[\".bench/landing-race-phase.sh\"]}]}\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "phases.json"), []byte(phases), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "cmd", "bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module landingrace\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cmd", "bench", "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	build := "#!/bin/sh\nset -eu\nroot=$1\nout=$2\nstaged=$out.staged\ncp \"$LAND_BASELINE_BENCH\" \"$staged\"\nchmod 0700 \"$staged\"\n\"$staged\" freshness-publish \"$root\" \"$out\"\n"
	if err := os.WriteFile(filepath.Join(root, "scripts", "go-build.sh"), []byte(build), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "go-build.inputs"), []byte("build_script=scripts/go-build.sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	inputs := "{\"schema\":1,\"closure\":\"local\",\"environment\":[\"HOME\",\"LAND_BASELINE_BENCH\",\"LAND_GATE_TALLY\",\"LAND_GATE_TREES\",\"LAND_RACE_READY\",\"LAND_RACE_RELEASE\"],\"paths\":[],\"tools\":[]}\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(inputs), 0o644); err != nil {
		t.Fatal(err)
	}
	systemGit(t, root, "add", ".")
	systemGit(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "artifact recovery fixture")
	base := systemGitOutput(t, root, "rev-parse", "HEAD")
	systemGit(t, root, "update-ref", "refs/bench/green/main", base)
}

func requirePublishedArtifactOwnerRecord(t *testing.T, bundle string, pid int, repository string) {
	t.Helper()
	record, err := prospectiveartifact.ReadPublished(filepath.Join(bundle, prospectiveartifact.RecordName))
	if err != nil {
		t.Fatalf("published owner record: %v", err)
	}
	common, err := prospectiveartifact.CanonicalCommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	want := prospectiveartifact.Record{
		Schema:    prospectiveartifact.RecordSchema,
		OwnerPID:  pid,
		CommonDir: common,
	}
	if record != want {
		t.Fatalf("published owner record = %#v, want %#v", record, want)
	}
}

// artifactRootEntries names what the shared temporary root holds, entry by entry with its
// mode. A row that must prove a sweep changed no path compares two of these.
func artifactRootEntries(t *testing.T, root string) string {
	t.Helper()
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	described := make([]string, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			t.Fatal(err)
		}
		described = append(described, entry.Name()+" "+info.Mode().String())
	}
	return strings.Join(described, "\n")
}

func releaseArtifactBarrier(t *testing.T, release string) {
	t.Helper()
	releaseWriter, err := os.OpenFile(release, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := releaseWriter.Write([]byte("go\n")); err != nil {
		t.Fatal(err)
	}
	if err := releaseWriter.Close(); err != nil {
		t.Fatal(err)
	}
}

func artifactLandEnv(root, home, tally, trees, ready, release string) []string {
	return append(systemLandEnv(root, home, tally, trees, ready, release), "BENCH_RUN_BINARY", "HOME="+home, "LAND_BASELINE_BENCH="+owner.selected.path)
}

func startArtifactLand(t *testing.T, root, home, tally, trees, ready, release string, source systemLandingWorktree, base string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer, chan error) {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(owner.selected.path, systemLandArgs(source, base)...)
	cmd.Dir = root
	cmd.Env = mergeEnvironment(os.Environ(), artifactLandEnv(root, home, tally, trees, ready, release))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	owner.mu.Lock()
	owner.starts++
	owner.mu.Unlock()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
	})
	return cmd, stdout, stderr, done
}

func awaitArtifactBarrier(t *testing.T, ready string, done <-chan error, stdout, stderr *bytes.Buffer) {
	t.Helper()
	observed := make(chan error, 1)
	go func() {
		file, err := os.Open(ready)
		if err != nil {
			observed <- err
			return
		}
		defer file.Close()
		byteRead := make([]byte, 1)
		_, err = file.Read(byteRead)
		if err == nil && string(byteRead) != "r" {
			err = errors.New("unexpected prospective gate barrier")
		}
		observed <- err
	}()
	select {
	case err := <-observed:
		if err != nil {
			t.Fatal(err)
		}
	case err := <-done:
		t.Fatalf("prospective owner ended before the artifact barrier: %v, stdout=%q, stderr=%q", err, stdout.String(), stderr.String())
	case <-time.After(20 * time.Second):
		t.Fatal("prospective owner did not reach the artifact barrier")
	}
}

func oneProspectiveBundle(t *testing.T, root string) string {
	t.Helper()
	bundles := prospectiveBundles(t, root)
	if len(bundles) == 1 {
		return bundles[0]
	}
	t.Fatalf("prospective bundles = %v, want one", bundles)
	return ""
}

func prospectiveBundles(t *testing.T, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	var bundles []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), prospectiveartifact.BundlePrefix) {
			bundles = append(bundles, filepath.Join(root, entry.Name()))
		}
	}
	return bundles
}

func oneOwnerBinary(t *testing.T, bundle string) string {
	t.Helper()
	entries, err := os.ReadDir(bundle)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "bench-run-") {
			return filepath.Join(bundle, entry.Name(), "bench")
		}
	}
	t.Fatal("prospective owner published no run binary")
	return ""
}

// TestConcurrentAuthorizationRetainsALiveOwner is PAR11. A live owner holds a registered
// checkout and an owner-authored run binary under the shared temporary root while a
// second authorization runs from end to end. An unconditional prefix sweep takes those
// resources with it.
//
// The live owner is a live process beside a planted bundle rather than a second gate
// run: the gate's own execution lock refuses a second concurrent authorization in one
// repository, so the sweep never sees two gate runs. The classification input -- a
// recognized bundle for this repository whose recorded PID answers -- is the same one a
// second gate run would present.
func TestConcurrentAuthorizationRetainsALiveOwner(t *testing.T) {
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	configureArtifactLandingFixture(t, root)
	private := t.TempDir()
	t.Setenv("TMPDIR", private)
	base := systemGitOutput(t, root, "rev-parse", "main")
	live, liveCheckout, liveBinary := plantLiveProspectiveBundle(t, private, root)

	_, out, errOut, done := startArtifactAuthorization(t, root, home, tally, trees, ready, release, base, "artifact-second", "winner.txt", "second prospective subject")
	awaitArtifactBarrier(t, ready, done, out, errOut)
	releaseArtifactBarrier(t, release)
	<-done

	requireLiveBundleRetained(t, root, live, liveCheckout, liveBinary)
}

// TestOneAuthorizationRecoversTheDeadBundleOfADeadAndLivePair is PAR20 across the
// process boundary. One killed owner's bundle and one live owner's bundle share the
// temporary root, so a sweep that classified the pair together would expose the live
// owner's resources to the dead owner's recovery. The live owner is a live process
// beside a planted bundle for the reason PAR11 records.
func TestOneAuthorizationRecoversTheDeadBundleOfADeadAndLivePair(t *testing.T) {
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	configureArtifactLandingFixture(t, root)
	private := t.TempDir()
	t.Setenv("TMPDIR", private)
	base := systemGitOutput(t, root, "rev-parse", "main")
	live, liveCheckout, liveBinary := plantLiveProspectiveBundle(t, private, root)

	deadCommand, deadOut, deadErr, deadDone := startArtifactAuthorization(t, root, home, tally, trees, ready, release, base, "artifact-dead", "loser.txt", "dead prospective subject")
	awaitArtifactBarrier(t, ready, deadDone, deadOut, deadErr)
	dead := otherProspectiveBundle(t, private, live)
	deadCheckout := filepath.Join(dead, prospectiveartifact.CheckoutName)
	requireCheckoutRegistered(t, root, deadCheckout, true)
	if err := syscall.Kill(-deadCommand.Process.Pid, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	if code := systemExitCode(<-deadDone); code == 0 {
		t.Fatal("killed prospective owner exited green")
	}
	releaseArtifactBarrier(t, release)

	_, freshOut, freshErr, freshDone := startArtifactAuthorization(t, root, home, tally, trees, ready, release, base, "artifact-fresh", "winner.txt", "fresh prospective subject")
	awaitArtifactBarrier(t, ready, freshDone, freshOut, freshErr)

	if _, err := os.Stat(dead); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("dead bundle = %v, want absent", err)
	}
	requireCheckoutRegistered(t, root, deadCheckout, false)
	requireLiveBundleRetained(t, root, live, liveCheckout, liveBinary)

	releaseArtifactBarrier(t, release)
	<-freshDone
}

// plantLiveProspectiveBundle publishes one bundle for root whose recorded owner is a
// live process, and gives it the resources a blocked owner holds: a registered private
// checkout and an owner-authored run binary.
func plantLiveProspectiveBundle(t *testing.T, private, root string) (string, string, string) {
	t.Helper()
	live := exec.Command("sleep", "600")
	if err := live.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- live.Wait() }()
	t.Cleanup(func() {
		if live.ProcessState == nil {
			_ = live.Process.Kill()
			<-done
		}
	})
	bundle, err := os.MkdirTemp(private, prospectiveartifact.BundlePrefix)
	if err != nil {
		t.Fatal(err)
	}
	common, err := prospectiveartifact.CanonicalCommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	record := prospectiveartifact.Record{
		Schema:    prospectiveartifact.RecordSchema,
		OwnerPID:  live.Process.Pid,
		CommonDir: common,
	}
	if err := prospectiveartifact.Publish(bundle, record); err != nil {
		t.Fatal(err)
	}
	checkout := filepath.Join(bundle, prospectiveartifact.CheckoutName)
	systemGit(t, root, "worktree", "add", "-q", "--detach", checkout, "HEAD")
	binary := filepath.Join(bundle, "bench-run-live", "bench")
	if err := os.MkdirAll(filepath.Dir(binary), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	return bundle, checkout, binary
}

func requireLiveBundleRetained(t *testing.T, root, bundle, checkout, binary string) {
	t.Helper()
	if _, err := os.Stat(bundle); err != nil {
		t.Fatalf("live bundle = %v, want retained", err)
	}
	if _, err := os.Stat(checkout); err != nil {
		t.Fatalf("live checkout = %v, want retained", err)
	}
	if _, err := os.Stat(binary); err != nil {
		t.Fatalf("live owner binary = %v, want retained", err)
	}
	requireCheckoutRegistered(t, root, checkout, true)
}

// startArtifactAuthorization reviews one new source worktree and starts its landing, so
// a row names the authorizations it wants rather than repeating the review and start
// sequence for each.
func startArtifactAuthorization(t *testing.T, root, home, tally, trees, ready, release, base, name, file, message string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer, chan error) {
	t.Helper()
	source := systemCreateLandingWorktree(t, root, home, name, name)
	systemCommit(t, source.path, file, message+"\n", message)
	source.tip = systemGitOutput(t, source.path, "rev-parse", "HEAD")
	if review := systemSelected(t, source.path, systemLandEnv(root, home, tally, trees, ready, release), "preflight", "review", "x", "--base", base); review.code != 0 {
		t.Fatalf("%s review = (%d, %q, %q)", name, review.code, review.stdout, review.stderr)
	}
	return startArtifactLand(t, root, home, tally, trees, ready, release, source, base)
}

// otherProspectiveBundle answers the one bundle under private that no held owner opened.
// Each owner opens its own bundle, so a row names the bundles it already holds and takes
// the newcomer.
func otherProspectiveBundle(t *testing.T, private string, held ...string) string {
	t.Helper()
	holders := map[string]bool{}
	for _, bundle := range held {
		holders[bundle] = true
	}
	var others []string
	for _, bundle := range prospectiveBundles(t, private) {
		if !holders[bundle] {
			others = append(others, bundle)
		}
	}
	if len(others) != 1 {
		t.Fatalf("unheld prospective bundles = %v, want one", others)
	}
	return others[0]
}

func requireCheckoutRegistered(t *testing.T, root, checkout string, want bool) {
	t.Helper()
	if got := strings.Contains(systemGitOutput(t, root, "worktree", "list", "--porcelain"), checkout); got != want {
		t.Fatalf("registration of %q = %v, want %v", checkout, got, want)
	}
}
