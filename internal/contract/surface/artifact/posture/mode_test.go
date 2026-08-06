package posture

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestArtifactCallersSelectUnsealedBuilderMode(t *testing.T) {
	root := contract.SubjectRoot(t)
	for _, script := range []string{"scripts/build-artifacts.sh", "scripts/native-proof.sh"} {
		data, err := os.ReadFile(filepath.Join(root, script))
		if err != nil {
			t.Fatal(err)
		}
		if got := strings.Count(string(data), "go-build.sh\" --mode artifact"); got != 2 {
			t.Fatalf("%s artifact builder mode calls = %d, want 2", script, got)
		}
		if strings.Contains(string(data), "rm -f \"$binary.seal\"") {
			t.Fatalf("%s deletes a seal after building instead of requesting an unsealed artifact", script)
		}
	}
}

func TestArtifactModeRefusesMalformedSelectorsWithoutChangingPriorPair(t *testing.T) {
	root := contract.SubjectRoot(t)
	script := filepath.Join(root, "scripts", "go-build.sh")
	for _, args := range [][]string{
		{"--mode", "unknown", root},
		{"--mode", "artifact"},
		{"--mode", "artifact", "--mode", "artifact", root, "output"},
	} {
		t.Run(strings.Join(args[:2], " "), func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "existing output")
			if err := os.WriteFile(out, []byte("prior executable"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(out+".seal", []byte("prior seal"), 0o644); err != nil {
				t.Fatal(err)
			}
			args = append(args, out)
			if output, err := exec.Command("bash", append([]string{script}, args...)...).CombinedOutput(); err == nil {
				t.Fatalf("malformed artifact selector succeeded: %q", output)
			}
			if data, err := os.ReadFile(out); err != nil || string(data) != "prior executable" {
				t.Fatalf("malformed selector changed prior output: %q, %v", data, err)
			}
			if data, err := os.ReadFile(out + ".seal"); err != nil || string(data) != "prior seal" {
				t.Fatalf("malformed selector changed prior seal: %q, %v", data, err)
			}
		})
	}
}

func TestArtifactModeCompileFailurePreservesPriorPair(t *testing.T) {
	root := contract.SubjectRoot(t)
	out := filepath.Join(t.TempDir(), "existing output")
	if err := os.WriteFile(out, []byte("prior executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(out+".seal", []byte("prior seal"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, out)
	cmd.Env = append(os.Environ(), "GOOS=not-a-go-target", "GOARCH=amd64")
	if output, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("invalid artifact target succeeded: %q", output)
	}
	if data, err := os.ReadFile(out); err != nil || string(data) != "prior executable" {
		t.Fatalf("compile failure changed prior output: %q, %v", data, err)
	}
	if data, err := os.ReadFile(out + ".seal"); err != nil || string(data) != "prior seal" {
		t.Fatalf("compile failure changed prior seal: %q, %v", data, err)
	}
}

func TestArtifactModeHandlesDashLedRelativeOutputLiterally(t *testing.T) {
	root := contract.SubjectRoot(t)
	dir := fmt.Sprintf("-artifact-output-%d", os.Getpid())
	out := filepath.Join(dir, "bench [*]")
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Join(root, dir)) })
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, out)
	cmd.Env = fakeBuilderEnv(t, "complete")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("dash-led relative output failed: %v\n%s", err, output)
	}
	if info, err := os.Stat(filepath.Join(root, out)); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("dash-led output was not published literally: %v, %v", info, err)
	}
}

func TestArtifactModeRefusesHostileOutputTypes(t *testing.T) {
	root := contract.SubjectRoot(t)
	for _, test := range []struct {
		name  string
		setup func(*testing.T) (string, func())
		check func(*testing.T, string)
	}{
		{
			name: "directory",
			setup: func(t *testing.T) (string, func()) {
				out := filepath.Join(t.TempDir(), "destination")
				if err := os.Mkdir(out, 0o755); err != nil {
					t.Fatal(err)
				}
				return out, func() {}
			},
			check: func(t *testing.T, out string) {
				if info, err := os.Stat(out); err != nil || !info.IsDir() {
					t.Fatalf("directory target changed: %v, %v", info, err)
				}
			},
		},
		{
			name: "fifo",
			setup: func(t *testing.T) (string, func()) {
				out := filepath.Join(t.TempDir(), "destination")
				if err := syscall.Mkfifo(out, 0o600); err != nil {
					t.Fatalf("FIFO capability unavailable: %v", err)
				}
				return out, func() {}
			},
			check: func(t *testing.T, out string) {
				if info, err := os.Lstat(out); err != nil || info.Mode()&os.ModeNamedPipe == 0 {
					t.Fatalf("FIFO target changed: %v, %v", info, err)
				}
			},
		},
		{
			name: "socket",
			setup: func(t *testing.T) (string, func()) {
				out := filepath.Join(t.TempDir(), "destination")
				listener, err := net.Listen("unix", out)
				if err != nil {
					t.Fatalf("Unix socket capability unavailable: %v", err)
				}
				return out, func() { _ = listener.Close() }
			},
			check: func(t *testing.T, out string) {
				if info, err := os.Lstat(out); err != nil || info.Mode()&os.ModeSocket == 0 {
					t.Fatalf("socket target changed: %v, %v", info, err)
				}
			},
		},
		{
			name: "device",
			setup: func(t *testing.T) (string, func()) {
				if info, err := os.Stat("/dev/null"); err != nil || info.Mode()&os.ModeDevice == 0 {
					t.Fatalf("device capability unavailable: %v", err)
				}
				return "/dev/null", func() {}
			},
			check: func(t *testing.T, out string) {
				if info, err := os.Stat(out); err != nil || info.Mode()&os.ModeDevice == 0 {
					t.Fatalf("device target changed: %v, %v", info, err)
				}
			},
		},
		{
			name: "dangling-symlink",
			setup: func(t *testing.T) (string, func()) {
				out := filepath.Join(t.TempDir(), "destination")
				if err := os.Symlink("missing-target", out); err != nil {
					t.Fatal(err)
				}
				return out, func() {}
			},
			check: func(t *testing.T, out string) {
				if target, err := os.Readlink(out); err != nil || target != "missing-target" {
					t.Fatalf("dangling symlink changed: %q, %v", target, err)
				}
			},
		},
		{
			name: "symlink-component",
			setup: func(t *testing.T) (string, func()) {
				base := t.TempDir()
				real := filepath.Join(base, "real")
				if err := os.Mkdir(real, 0o755); err != nil {
					t.Fatal(err)
				}
				link := filepath.Join(base, "linked")
				if err := os.Symlink(real, link); err != nil {
					t.Fatal(err)
				}
				return filepath.Join(link, "destination"), func() {}
			},
			check: func(t *testing.T, out string) {
				if _, err := os.Lstat(out); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("symlink component received output: %v", err)
				}
			},
		},
		{
			name: "unwritable",
			setup: func(t *testing.T) (string, func()) {
				out := filepath.Join(t.TempDir(), "destination")
				if err := os.WriteFile(out, []byte("prior executable"), 0o400); err != nil {
					t.Fatal(err)
				}
				return out, func() { _ = os.Chmod(out, 0o600) }
			},
			check: func(t *testing.T, out string) {
				data, err := os.ReadFile(out)
				if err != nil || string(data) != "prior executable" {
					t.Fatalf("unwritable target changed: %q, %v", data, err)
				}
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			out, cleanup := test.setup(t)
			defer cleanup()
			if test.name != "device" {
				if err := os.WriteFile(out+".seal", []byte("prior seal"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, out)
			cmd.Env = fakeBuilderEnv(t, "complete")
			if output, err := cmd.CombinedOutput(); err == nil {
				t.Fatalf("hostile %s target succeeded: %q", test.name, output)
			}
			test.check(t, out)
			if test.name != "device" {
				if data, err := os.ReadFile(out + ".seal"); err != nil || string(data) != "prior seal" {
					t.Fatalf("hostile %s target changed prior seal: %q, %v", test.name, data, err)
				}
			}
		})
	}
}

func TestGoBuildInterruptionsPreservePriorPairAndRemoveStaging(t *testing.T) {
	for _, test := range []struct {
		name     string
		modeArgs []string
		block    string
	}{
		{name: "artifact-compilation", modeArgs: []string{"--mode", "artifact"}, block: "compile"},
		{name: "subject-compilation", block: "compile"},
		{name: "subject-publication", block: "publish"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := contract.SubjectRoot(t)
			dir := t.TempDir()
			out := filepath.Join(dir, "prior output")
			if err := os.WriteFile(out, []byte("prior executable"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(out+".seal", []byte("prior seal"), 0o644); err != nil {
				t.Fatal(err)
			}
			ready := filepath.Join(dir, "blocked")
			args := append([]string{filepath.Join(root, "scripts", "go-build.sh")}, test.modeArgs...)
			args = append(args, root, out)
			cmd := exec.Command("bash", args...)
			cmd.Env = fakeBuilderEnv(t, test.block, "BENCH_TEST_BLOCK_READY="+ready)
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			awaitFile(t, ready, cmd)
			if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { done <- cmd.Wait() }()
			select {
			case err := <-done:
				if err == nil {
					t.Fatal("interrupted builder exited successfully")
				}
			case <-time.After(5 * time.Second):
				_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
				t.Fatal("interrupted builder did not exit within five seconds")
			}
			assertPriorPair(t, out)
			if matches, err := filepath.Glob(filepath.Join(dir, ".bench.*")); err != nil || len(matches) != 0 {
				t.Fatalf("interruption left staged residue: %v, %v", matches, err)
			}
		})
	}
}

// TestArtifactInstallNeverExposesAnAbsentPriorOutput holds the builder still at the one
// step that replaces the output and reads the destination from outside. A move-aside
// install passes every before-and-after check while the prior pair is missing for the
// whole compile-to-promote span, so only an observation taken mid-install can see it.
func TestArtifactInstallNeverExposesAnAbsentPriorOutput(t *testing.T) {
	root := contract.SubjectRoot(t)
	dir := t.TempDir()
	out := filepath.Join(dir, "prior output")
	writePriorPair(t, out)
	ready := filepath.Join(dir, "installing")
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, out)
	cmd.Env = fakeBuilderEnv(t, "complete", "BENCH_TEST_BLOCK_PROMOTION=1", "BENCH_TEST_BLOCK_READY="+ready)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	awaitFile(t, ready, cmd)
	assertPriorPair(t, out)
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("builder interrupted during its install exited successfully")
		}
	case <-time.After(5 * time.Second):
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		t.Fatal("builder interrupted during its install did not exit within five seconds")
	}
	assertPriorPair(t, out)
	if residue, err := filepath.Glob(filepath.Join(dir, ".bench.*")); err != nil || len(residue) != 0 {
		t.Fatalf("interrupted install left staging residue: %v, %v", residue, err)
	}
}

// TestArtifactSignalAfterInstallLeavesNoStaleSeal holds the builder still in the one
// window where the destination carries the new artifact and the retired subject's seal at
// once, and terminates it there. Two directory entries cannot change together, so only an
// exit taken inside that window shows whether the builder finishes the install on its way
// out or abandons the pair half-converted.
func TestArtifactSignalAfterInstallLeavesNoStaleSeal(t *testing.T) {
	root := contract.SubjectRoot(t)
	dir := t.TempDir()
	out := filepath.Join(dir, "prior output")
	writePriorPair(t, out)
	ready := filepath.Join(dir, "promoted")
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, out)
	cmd.Env = fakeBuilderEnv(t, "complete", "BENCH_TEST_BLOCK_AFTER_PROMOTION=1", "BENCH_TEST_BLOCK_READY="+ready, "BENCH_TEST_BUILD_BODY=signalled-artifact")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	awaitFile(t, ready, cmd)
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("builder terminated after its install exited successfully")
		}
	case <-time.After(5 * time.Second):
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		t.Fatal("builder terminated after its install did not exit within five seconds")
	}
	if data, err := os.ReadFile(out); err != nil || !bytes.Contains(data, []byte("signalled-artifact")) {
		t.Fatalf("terminated install did not leave the new artifact: %q, %v", data, err)
	}
	if _, err := os.Stat(out + ".seal"); !os.IsNotExist(err) {
		t.Fatalf("terminated install left the retired subject's seal: %v", err)
	}
	if residue, err := filepath.Glob(filepath.Join(dir, ".bench.*")); err != nil || len(residue) != 0 {
		t.Fatalf("terminated install left staging residue: %v, %v", residue, err)
	}
}

// writeFakeGoBuilder installs the package's one fake `go` on dir, ahead of the real
// toolchain on a builder's PATH. It re-derives the builder's own `-o` argv contract to
// learn where the output belongs, then fabricates a bench-like executable there. Every
// leg — the compile and exec logs, the blocking seams, the malformed output, the staged
// freshness-publish copy — is selected by the environment the caller hands the builder,
// so a consumer that needs none of them passes none of them.
func writeFakeGoBuilder(t *testing.T, dir string) {
	t.Helper()
	goStub := `#!/usr/bin/env bash
set -euo pipefail
if [[ -n "${BENCH_TEST_GO_LOG:-}" ]]; then printf 'target=%s/%s go' "${GOOS:-}" "${GOARCH:-}"; printf ' <%s>' "$@"; printf '\n'; fi >> "${BENCH_TEST_GO_LOG:-/dev/null}"
[[ -z "${BENCH_PUBLICATION_TRACE:-}" ]] || printf 'go:%s\n' "$*" >> "$BENCH_PUBLICATION_TRACE"
out=""
while [[ "$#" -gt 0 ]]; do
  if [[ "$1" == -o ]]; then out="$2"; shift 2; continue; fi
  shift
done
if [[ "${BENCH_TEST_BUILD_BEHAVIOR:-}" == compile ]]; then
  : > "$BENCH_TEST_BLOCK_READY"
  while :; do sleep 1; done
fi
if [[ "${BENCH_TEST_BUILD_BEHAVIOR:-}" == fail ]]; then exit 91; fi
if [[ "${BENCH_TEST_BUILD_BEHAVIOR:-}" == invalid ]]; then printf 'not executable\n' > "$out"; exit 0; fi
printf '%s\n' '#!/usr/bin/env bash' 'set -euo pipefail' '[[ -z "${BENCH_TEST_EXEC_LOG:-}" ]] || printf "bench <%s>\\n" "$*" >> "$BENCH_TEST_EXEC_LOG"' '[[ -z "${BENCH_PUBLICATION_TRACE:-}" ]] || printf "bench:%s\\n" "$*" >> "$BENCH_PUBLICATION_TRACE"' 'if [[ "${BENCH_TEST_BUILD_BEHAVIOR:-}" == publish ]]; then : > "$BENCH_TEST_BLOCK_READY"; while :; do sleep 1; done; fi' 'if [[ "${1:-}" == freshness-publish ]]; then cp -- "$0" "$3"; printf "fixture seal\\n" > "$3.seal"; fi' "# ${BENCH_TEST_BUILD_BODY:-fixture}" > "$out"
chmod 0755 "$out"
`
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(goStub), 0o755); err != nil {
		t.Fatal(err)
	}
}

func fakeBuilderEnv(t *testing.T, behavior string, extra ...string) []string {
	t.Helper()
	dir := t.TempDir()
	writeFakeGoBuilder(t, dir)
	mvStub := "#!/usr/bin/env bash\nif [[ \"${BENCH_TEST_FAIL_PROMOTION:-}\" == 1 && \" $* \" == *'/.bench.tmp.'* ]]; then exit 73; fi\nif [[ \"${BENCH_TEST_BLOCK_PROMOTION:-}\" == 1 && \" $* \" == *'/.bench.tmp.'* ]]; then : > \"$BENCH_TEST_BLOCK_READY\"; while :; do sleep 1; done; fi\nif [[ \"${BENCH_TEST_BLOCK_AFTER_PROMOTION:-}\" == 1 && \" $* \" == *'/.bench.tmp.'* ]]; then /usr/bin/mv \"$@\"; : > \"$BENCH_TEST_BLOCK_READY\"; while :; do sleep 1; done; fi\nexec /usr/bin/mv \"$@\"\n"
	if err := os.WriteFile(filepath.Join(dir, "mv"), []byte(mvStub), 0o755); err != nil {
		t.Fatal(err)
	}
	env := append(os.Environ(), "PATH="+dir+string(os.PathListSeparator)+os.Getenv("PATH"), "BENCH_TEST_BUILD_BEHAVIOR="+behavior)
	return append(env, extra...)
}

func awaitFile(t *testing.T, path string, cmd *exec.Cmd) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			t.Fatalf("builder did not reach blocked seam %s", path)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func assertPriorPair(t *testing.T, out string) {
	t.Helper()
	if data, err := os.ReadFile(out); err != nil || string(data) != "prior executable" {
		t.Fatalf("prior output changed: %q, %v", data, err)
	}
	if data, err := os.ReadFile(out + ".seal"); err != nil || string(data) != "prior seal" {
		t.Fatalf("prior seal changed: %q, %v", data, err)
	}
}
