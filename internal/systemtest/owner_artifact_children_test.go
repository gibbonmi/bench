//go:build system

package systemtest

import (
	"bytes"
	"os/exec"
	"syscall"
	"testing"
)

// artifactChildEnv is the landing fixture's environment with private as TMPDIR. A review
// or landing child opens its prospective bundle under TMPDIR, so each such child takes the
// root that the row reads. The value goes to each child, because parallel tests share the
// process environment.
func artifactChildEnv(root, home, tally, trees, ready, release, private string) []string {
	return append(systemLandEnv(root, home, tally, trees, ready, release), "TMPDIR="+private)
}

func artifactLandEnv(root, home, tally, trees, ready, release, private string) []string {
	return append(artifactChildEnv(root, home, tally, trees, ready, release, private), "BENCH_RUN_BINARY", "HOME="+home, "LAND_BASELINE_BENCH="+owner.selected.path)
}

func startArtifactLand(t *testing.T, root, home, tally, trees, ready, release, private string, source systemLandingWorktree, base string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer, chan error) {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(owner.selected.path, systemLandArgs(source, base)...)
	cmd.Dir = root
	cmd.Env = mergeEnvironment(systemBaseEnvironment(), artifactLandEnv(root, home, tally, trees, ready, release, private))
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
