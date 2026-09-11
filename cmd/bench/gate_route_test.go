package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

func TestGateRouteOwnerProcess(t *testing.T) {
	if os.Getenv("BENCH_GATE_ROUTE_PROCESS") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg == "--" {
			os.Exit((Command{Stdout: os.Stdout, Stderr: os.Stderr}).Run(os.Args[i+1:]))
		}
	}
	t.Fatal("missing owner process arguments")
}

func gateRoute(t *testing.T, root string, args ...string) (string, string, int) {
	t.Helper()
	kit := t.TempDir()
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(kit, "bin", "bench.sh"))
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	writeExecutable(t, filepath.Join(kit, "dist", "bench"), "#!/bin/sh\nexec \"$BENCH_GATE_ROUTE_EXECUTABLE\" -test.run '^TestGateRouteOwnerProcess$' -- \"$@\"\n")
	cmd := exec.Command("bash", append([]string{filepath.Join(kit, "bin", "bench.sh"), "gate"}, args...)...)
	if root != "" {
		cmd.Dir = root
	}
	environment := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	environment = capability.WithoutEnvironment(capability.WithoutEnvironment(environment, "BENCH_KIT"), "BENCH_WRAPPER")
	cmd.Env = append(environment, "BENCH_GATE_ROUTE_PROCESS=1", "BENCH_GATE_ROUTE_EXECUTABLE="+executable, "BENCH_HOME="+filepath.Join(kit, "home"), "BENCH_KIT="+kit)
	var out, diagnostic bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &diagnostic
	code := 0
	if err := cmd.Run(); err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatal(err)
		}
		code = exit.ExitCode()
	}
	return out.String(), diagnostic.String(), code
}

func TestShellWrapperRejectsUnknownGateShapes(t *testing.T) {
	for _, args := range [][]string{{"pin"}, {"--fresh", "unexpected"}, {"--brief"}, {"--checkpoint"}, {"--checkpoint", "specs/x/spec.md"}, {"--chunk", "1"}, {"--checkpoint", "specs/x/spec.md", "--chunk", "1", "--complete"}} {
		out, diagnostic, code := gateRoute(t, "", args...)
		if code != 2 || out != "" || diagnostic != gate.CommandUsage+"\n" {
			t.Fatalf("gate %v: %d %q %q", args, code, out, diagnostic)
		}
	}
}

func TestGateCheckpointRoute(t *testing.T) {
	f := recordtest.New(t, 1)
	f.Write(".bench/gate.sh", "#!/bin/sh\nexit 0\n")
	if err := os.Chmod(filepath.Join(f.Root, ".bench/gate.sh"), 0755); err != nil {
		t.Fatal(err)
	}
	f.Write(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.AddChunk()
	f.Save()
	f.Commit("record evidence")
	_, diagnostic, code := gateRoute(t, f.Root, "--checkpoint", recordtest.Spec, "--chunk", "1")
	if code != 0 {
		t.Fatalf("public chunk checkpoint: %d %s", code, diagnostic)
	}
	f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2]
	f.Save()
	_, diagnostic, code = gateRoute(t, f.Root, "--checkpoint", recordtest.Spec, "--chunk", "1")
	if code == 0 || !strings.Contains(diagnostic, "missing Coverage") {
		t.Fatalf("public wrapper lost checkpoint: %d %s", code, diagnostic)
	}
}

func TestRunGateRejectsBriefUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"gate", "--brief"})
	if code != 2 || stdout.Len() != 0 || stderr.String() != gate.CommandUsage+"\n" {
		t.Fatalf("gate --brief = (stdout %q, stderr %q, exit %d), want usage on stderr and exit 2", stdout.String(), stderr.String(), code)
	}
}
