// Package probetest installs the canned `go` that probe tests run under. A probe starts one
// baseline `go test` child and then one mutated child for each mutant, so a test supplies one
// answer for each start. A test outside the probe package runs the probe over the clamp
// repository here. Only tests import this package; it holds no verdict, row, or record schema.
package probetest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/treetarget/kittest"
)

// ClampPath is the subject of the clamp repository, and ClampSource is its start bytes.
const (
	ClampPath   = "clamp/clamp.go"
	ClampSource = "package clamp\n\n// Clamp keeps n at or above zero.\nfunc Clamp(n int) int {\n\tif n < 0 {\n\t\treturn 0\n\t}\n\treturn n\n}\n"
)

// ClampSelection selects TestClampNegative of the clamp repository.
var ClampSelection = []string{"--package", "./clamp", "--run", "^TestClampNegative$"}

// ClampPass is a run of the clamp package that passes TestClampNegative.
var ClampPass = Start{Events: `{"Action":"run","Package":"benchkit/clamp","Test":"TestClampNegative"}
{"Action":"pass","Package":"benchkit/clamp","Test":"TestClampNegative","Elapsed":0}
{"Action":"pass","Package":"benchkit/clamp","Elapsed":0.01}`}

// ClampFailure is a run of the clamp package that fails TestClampNegative at exit 1, with
// diagnostic as the test's output line.
func ClampFailure(diagnostic string) Start {
	return Start{Events: fmt.Sprintf(`{"Action":"run","Package":"benchkit/clamp","Test":"TestClampNegative"}
{"Action":"output","Package":"benchkit/clamp","Test":"TestClampNegative","Output":"    %s\n"}
{"Action":"fail","Package":"benchkit/clamp","Test":"TestClampNegative","Elapsed":0}
{"Action":"fail","Package":"benchkit/clamp","Elapsed":0.01}`, diagnostic), Exit: 1}
}

// Repository writes a kit tree and the clamp package into the Git checkout dir and publishes a
// sealed run binary for it. For the rest of the test, it binds the kit, the run binary, a
// private Bench home, and the working directory to that root. It answers the root with every
// link resolved, and the home.
func Repository(t testing.TB, dir string) (string, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	kittest.WriteTree(t, root)
	if err := os.MkdirAll(filepath.Join(root, "clamp"), 0o755); err != nil {
		t.Fatal(err)
	}
	staged, executable := filepath.Join(t.TempDir(), "staged"), filepath.Join(t.TempDir(), "bench")
	for _, file := range []struct {
		path, body string
		mode       os.FileMode
	}{
		{filepath.Join(root, filepath.FromSlash(ClampPath)), ClampSource, 0o644},
		{filepath.Join(root, "clamp", "clamp_test.go"), "package clamp\n\nimport \"testing\"\n\nfunc TestClampNegative(t *testing.T) {\n\tif Clamp(-1) != 0 {\n\t\tt.Fatal(\"negative\")\n\t}\n}\n", 0o644},
		{staged, "#!/bin/sh\nexit 0\n", 0o755},
	} {
		if err := os.WriteFile(file.path, []byte(file.body), file.mode); err != nil {
			t.Fatal(err)
		}
	}
	if err := freshness.Publish(root, staged, executable, filepath.Dir(executable), "probetest"); err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	t.Setenv(runbinary.Env, executable)
	t.Setenv("BENCH_KIT", root)
	t.Setenv(benchhome.Env, home)
	t.Chdir(root)
	return root, home
}

// Start is what the canned `go` does on one `go test` start: the shell it runs first, the
// event stream it prints, and the exit it takes.
type Start struct {
	Hook   string
	Events string
	Exit   int
}

// unansweredExit is the exit of a `go test` start that InstallSequence has no answer for.
const unansweredExit = 97

// InstallStarts puts a canned `go` ahead of the real one for the single form. An odd test
// start answers baseline, and the even start after it answers mutated.
func InstallStarts(t testing.TB, marker string, baseline, mutated Start) {
	t.Helper()
	install(t, marker, "if [ $((n % 2)) = 1 ]; then\n"+answer(baseline)+"fi\n"+answer(mutated))
}

// InstallSequence puts a canned `go` ahead of the real one for a mutant set. The first test
// start answers baseline, and the k-th later start answers mutated[k-1]. A start past the list
// prints no event and exits unansweredExit.
func InstallSequence(t testing.TB, marker string, baseline Start, mutated ...Start) {
	t.Helper()
	arms := "1)\n" + answer(baseline) + ";;\n"
	for k, start := range mutated {
		arms += strconv.Itoa(k+2) + ")\n" + answer(start) + ";;\n"
	}
	install(t, marker, "case $n in\n"+arms+"esac\nexit "+strconv.Itoa(unansweredExit)+"\n")
}

// install writes the canned `go` for the rest of the test. Its prelude records every start in
// marker and passes each non-test call to the real toolchain. It reads the test start's
// ordinal into `$n` for body and its hooks. The marker is removed first, so each install
// restarts the count.
func install(t testing.TB, marker, body string) {
	t.Helper()
	real, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(marker); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	quoted := sanitize.ShellQuote(marker)
	script := "#!/usr/bin/env bash\n" +
		"printf '%s\\n' \"$1\" >> " + quoted + "\n" +
		"if [ \"$1\" != test ]; then exec " + sanitize.ShellQuote(real) + " \"$@\"; fi\n" +
		"n=$(grep -c '^test$' " + quoted + ")\n" + body
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func answer(start Start) string {
	return start.Hook + "cat <<'PROBEEOF'\n" + start.Events + "\nPROBEEOF\nexit " + strconv.Itoa(start.Exit) + "\n"
}
