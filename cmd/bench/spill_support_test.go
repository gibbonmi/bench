package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
)

// TestMain runs the package under a private Bench home. A bounded response that spills
// in a test then writes its spill file there, never under the operator's home.
func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "bench-cmd-home-")
	if err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME:", err)
		os.Exit(1)
	}
	code := 1
	if err := os.Setenv(benchhome.Env, home); err != nil {
		fmt.Fprintln(os.Stderr, "private BENCH_HOME:", err)
	} else {
		code = m.Run()
	}
	_ = os.RemoveAll(home)
	os.Exit(code)
}

// spilledResponse answers the complete output behind a response: the spill file that its
// spill line names, or the response itself when it did not spill. The spill file holds
// both streams in arrival order, and a spilled response prints nothing on stderr.
func spilledResponse(t *testing.T, response string) string {
	t.Helper()
	spill, ok := responseboundtest.Find(response)
	if !ok {
		return response
	}
	data, err := os.ReadFile(spill.Path)
	if err != nil {
		t.Fatalf("read spill file: %v", err)
	}
	return string(data)
}

func runAXICommandAt(t *testing.T, cwd string, argv []string) axiCommandResult {
	t.Helper()
	return runAXICommandAsAt(t, cwd, "bench", argv)
}

// runAXICommandAsAt runs argv through the production dispatcher from cwd. The stdout of
// the result is the complete output, read from the spill file when the response spilled.
func runAXICommandAsAt(t *testing.T, cwd, executable string, argv []string) axiCommandResult {
	t.Helper()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	}()
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: executable}.Run(argv)
	return axiCommandResult{stdout: spilledResponse(t, withoutTreeRow(stdout.String())), stderr: stderr.String(), code: code}
}
