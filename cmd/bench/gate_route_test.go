package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

const gateRouteProcessEnv = "BENCH_GATE_ROUTE_PROCESS"

// gateRouteOwnerProcess runs the owner process that the gateRoute wrapper starts from
// this test binary, and reports whether this process is that owner.
func gateRouteOwnerProcess() (int, bool) {
	if os.Getenv(gateRouteProcessEnv) != "1" {
		return 0, false
	}
	for i, arg := range os.Args {
		if arg == "--" {
			return Command{Stdout: os.Stdout, Stderr: os.Stderr}.Run(os.Args[i+1:]), true
		}
	}
	fmt.Fprintln(os.Stderr, "missing owner process arguments")
	return 2, true
}

// gateRoute runs the public shell wrapper's gate route with this test binary as the
// owner process. The owner's temporary directory must stay empty: an owner that exits
// before its cleanup leaves an entry there.
func gateRoute(t *testing.T, root string, args ...string) (string, string, int) {
	t.Helper()
	kit := t.TempDir()
	temporary := t.TempDir()
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(kit, "bin", "bench.sh"))
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	writeExecutable(t, filepath.Join(kit, "dist", "bench"), "#!/bin/sh\nexec \"$BENCH_GATE_ROUTE_EXECUTABLE\" -- \"$@\"\n")
	cmd := exec.Command("bash", append([]string{filepath.Join(kit, "bin", "bench.sh"), "gate"}, args...)...)
	if root != "" {
		cmd.Dir = root
	}
	environment := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	environment = capability.WithoutEnvironment(capability.WithoutEnvironment(environment, "BENCH_KIT"), "BENCH_WRAPPER")
	cmd.Env = append(environment, gateRouteProcessEnv+"=1", "BENCH_GATE_ROUTE_EXECUTABLE="+executable, "BENCH_HOME="+filepath.Join(kit, "home"), "BENCH_KIT="+kit, "TMPDIR="+temporary)
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
	left, err := os.ReadDir(temporary)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range left {
		t.Errorf("gate %v: owner process left %s in its temporary directory", args, entry.Name())
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
	f := recordtest.AttachAt(t, gittest.RepoOnBranch(t, "main"), 1, "specs/example [*] space/spec.md")
	f.Write(".bench/gate.sh", "#!/bin/sh\nexit 0\n")
	if err := os.Chmod(filepath.Join(f.Root, ".bench/gate.sh"), 0755); err != nil {
		t.Fatal(err)
	}
	f.Write(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.AddChunk()
	f.Complete()
	f.Save()
	f.Commit("record evidence")
	_, diagnostic, code := gateRoute(t, f.Root, "--checkpoint", f.Record.Spec, "--chunk", "1")
	if code != 0 {
		t.Fatalf("public chunk checkpoint: %d %s", code, diagnostic)
	}
	_, diagnostic, code = gateRoute(t, f.Root, "--checkpoint", f.Record.Spec, "--complete")
	if code != 0 {
		t.Fatalf("public complete checkpoint: %d %s", code, diagnostic)
	}
	f.Record.Chunks[0].Reviews = f.Record.Chunks[0].Reviews[:2]
	f.Save()
	_, diagnostic, code = gateRoute(t, f.Root, "--checkpoint", f.Record.Spec, "--chunk", "1")
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
