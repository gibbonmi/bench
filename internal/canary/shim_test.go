package canary

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCompatibilityShimDelegatesToBenchCanary(t *testing.T) {
	root := filepath.Join(t.TempDir(), "repo with space")
	mkdir(t, filepath.Join(root, ".bench", "lib"))
	copyFile(t, kitPath(t, ".bench", "lib", "canary-run.sh"), filepath.Join(root, ".bench", "lib", "canary-run.sh"))
	copyFile(t, kitPath(t, ".bench", "lib", "resolve-bench.sh"), filepath.Join(root, ".bench", "lib", "resolve-bench.sh"))
	runGitInit(t, root)

	bin := filepath.Join(t.TempDir(), "bin")
	mkdir(t, bin)
	callPath := filepath.Join(t.TempDir(), "call")
	write(t, filepath.Join(bin, "bench"), "#!/usr/bin/env bash\nprintf '%s\\n' \"$@\" > \"$BENCH_SHIM_CALL\"\nexit 7\n")
	if err := os.Chmod(filepath.Join(bin, "bench"), 0o755); err != nil {
		t.Fatal(err)
	}

	probe := `set -u
root="$1"
fail=0
err() { echo "ERR:$*" >&2; fail=1; }
. "$root/.bench/lib/canary-run.sh"
exit "$fail"
`
	cmd := exec.Command("bash", "-c", probe, "probe", root)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "BENCH_SHIM_CALL="+callPath)
	out, err := cmd.CombinedOutput()
	if err == nil || cmd.ProcessState.ExitCode() != 1 {
		t.Fatalf("shim probe exit err=%v code=%d output=%s, want fail=1", err, cmd.ProcessState.ExitCode(), out)
	}
	if !strings.Contains(string(out), "ERR:canary sweep failed") {
		t.Fatalf("shim did not fold nonzero into err, output:\n%s", out)
	}
	got := read(t, callPath)
	if got != "canary\n"+root+"\n" {
		t.Fatalf("shim call = %q, want canary + root", got)
	}
}

func kitPath(t *testing.T, elem ...string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	root := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	return filepath.Join(append([]string{root}, elem...)...)
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, data, 0o755); err != nil {
		t.Fatal(err)
	}
}

func runGitInit(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
