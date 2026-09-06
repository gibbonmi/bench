// The gate route family: the wrapper is driven out of process from a synthetic kit
// under t.TempDir(), over the shapes that never reach the oracle.
package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// gateUsageLine is the grammar the wrapper answers a rejected gate shape with. It is
// written here independently of bin/bench.sh, so a reworded or deleted usage arm reds
// this test rather than passing against a re-derived expectation.
const gateUsageLine = "usage: bench gate [--fresh]\n"

// TestShellWrapperRejectsUnknownGateShapes holds the gate route to its three accepted
// spellings: the bare run, --fresh, and help. Every other shape is a usage error, and
// none of these cases reaches run_gate, so the case needs no oracle and no dist/bench.
func TestShellWrapperRejectsUnknownGateShapes(t *testing.T) {
	root := t.TempDir()
	kit := filepath.Join(root, "kit")
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(kit, "bin", "bench.sh"))

	run := func(args ...string) (string, string, int) {
		t.Helper()
		cmd := exec.Command("bash", append([]string{filepath.Join(kit, "bin", "bench.sh"), "gate"}, args...)...)
		env := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
		for _, name := range []string{"BENCH_KIT", "BENCH_WRAPPER"} {
			env = capability.WithoutEnvironment(env, name)
		}
		cmd.Env = append(env, "BENCH_HOME="+filepath.Join(root, "home"), "BENCH_KIT="+kit)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		code := 0
		if err := cmd.Run(); err != nil {
			exit, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatalf("bench.sh gate %v: %v", args, err)
			}
			code = exit.ExitCode()
		}
		return stdout.String(), stderr.String(), code
	}

	for _, shape := range [][]string{{"pin"}, {"--fresh", "unexpected"}, {"--brief"}} {
		stdout, stderr, code := run(shape...)
		if code != 2 || stdout != "" || stderr != gateUsageLine {
			t.Errorf("gate %v = (exit %d, stdout %q, stderr %q), want exit 2 and %q on stderr", shape, code, stdout, stderr, gateUsageLine)
		}
	}
}
