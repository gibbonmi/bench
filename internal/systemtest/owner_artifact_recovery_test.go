//go:build system

package systemtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
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
	checkout := filepath.Join(bundle, "checkout")
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
	if result := systemSelected(t, root, artifactLandEnv(root, home, tally, trees, ready, release), systemLandArgs(fresh, base)...); result.code != 0 {
		t.Fatalf("second prospective authorization = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}
	if repeatedTally, err := os.ReadFile(tally); err != nil || string(repeatedTally) != string(gateTally) {
		t.Fatalf("reused exact-green tally = %q, %v, want unchanged %q", repeatedTally, err, gateTally)
	}
	if bundles := prospectiveBundles(t, private); len(bundles) != 0 {
		t.Fatalf("prospective bundles after idempotent retry = %v, want none", bundles)
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
	path := filepath.Join(bundle, "owner.json")
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		t.Fatalf("published owner record mode = %v, want regular 0600", info.Mode())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record struct {
		Schema    int    `json:"schema"`
		OwnerPID  int    `json:"owner_pid"`
		CommonDir string `json:"common_dir"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&record); err != nil {
		t.Fatalf("decode published owner record: %v", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatalf("published owner record trailing data = %v, want EOF", err)
	}
	common := systemGitOutput(t, repository, "rev-parse", "--path-format=absolute", "--git-common-dir")
	common, err = filepath.EvalSymlinks(common)
	if err != nil {
		t.Fatal(err)
	}
	if record.Schema != 1 || record.OwnerPID != pid || record.CommonDir != filepath.Clean(common) {
		t.Fatalf("published owner record = %#v, want schema 1, pid %d, common directory %q", record, pid, filepath.Clean(common))
	}
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
		if strings.HasPrefix(entry.Name(), "bench-prospective-artifact-") {
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
