// Flag and terminal-table tests for the landing command: operand parsing, help text, and pre-gate state refusals.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/landing"
)

func TestLandCommandNeverTreatsFlagValuesAsSourcePath(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"--request", "would-be-path", "--base", "b", "--source-tip", "s", "--spec", "x", "-m", "m"},
		{"--request", "r", "--base", "would-be-path", "--source-tip", "s", "--spec", "x", "-m", "m"},
		{"--request", "r", "--base", "b", "--source-tip", "would-be-path", "--spec", "x", "-m", "m"},
		{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "would-be-path", "-m", "m"},
		{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x", "-m", "would-be-path"},
	} {
		var stdout, stderr bytes.Buffer
		if code := LandCommand("", Home(), "", args, &stdout, &stderr); code != 2 {
			t.Fatalf("LandCommand(%q) = %d, stdout=%q stderr=%q", args, code, stdout.String(), stderr.String())
		}
	}
}

func TestLandCommandRequiredFlagsKeepDeclaredHelp(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args []string
		help string
	}{
		{"first request", []string{"--base", "b", "--source-tip", "s", "--spec", "x", "-m", "m", "path"}, landGrammar.Help},
		{"first base", []string{"--request", "r", "--source-tip", "s", "--spec", "x", "-m", "m", "path"}, landGrammar.Help},
		{"first source tip", []string{"--request", "r", "--base", "b", "--spec", "x", "-m", "m", "path"}, landGrammar.Help},
		{"first message", []string{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x", "path"}, landGrammar.Help},
		{"resume request", []string{"--resume", "p", "--base", "b", "--source-tip", "s", "--spec", "x", "path"}, resumeLandGrammar.Help},
		{"resume base", []string{"--resume", "p", "--request", "r", "--source-tip", "s", "--spec", "x", "path"}, resumeLandGrammar.Help},
		{"resume source tip", []string{"--resume", "p", "--request", "r", "--base", "b", "--spec", "x", "path"}, resumeLandGrammar.Help},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := LandCommand("", Home(), "", tc.args, &stdout, &stderr); code != 2 || stdout.Len() != 0 || stderr.String() != tc.help+"\n" {
				t.Fatalf("LandCommand(%q) = (%d, %q, %q), want (2, empty, %q)", tc.args, code, stdout.String(), stderr.String(), tc.help+"\n")
			}
		})
	}

	var stdout, stderr bytes.Buffer
	args := []string{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x", "path"}
	if code := ResumeLandCommand("", Home(), args, &stdout, &stderr); code != 2 || stdout.Len() != 0 || stderr.String() != resumeLandGrammar.Help+"\n" {
		t.Fatalf("ResumeLandCommand(%q) = (%d, %q, %q), want (2, empty, %q)", args, code, stdout.String(), stderr.String(), resumeLandGrammar.Help+"\n")
	}
}

func TestResumeLandCommandNeverTreatsFlagValuesAsSourcePath(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"--resume", "would-be-path", "--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x"},
		{"--resume", "p", "--request", "would-be-path", "--base", "b", "--source-tip", "s", "--spec", "x"},
		{"--resume", "p", "--request", "r", "--base", "would-be-path", "--source-tip", "s", "--spec", "x"},
		{"--resume", "p", "--request", "r", "--base", "b", "--source-tip", "would-be-path", "--spec", "x"},
		{"--resume", "p", "--request", "r", "--base", "b", "--source-tip", "s", "--spec", "would-be-path"},
	} {
		var stdout, stderr bytes.Buffer
		if code := LandCommand("", Home(), "", args, &stdout, &stderr); code != 2 {
			t.Fatalf("LandCommand(%q) = %d, stdout=%q stderr=%q", args, code, stdout.String(), stderr.String())
		}
	}
	args := []string{"--resume", "p", "--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x", "--", "-path"}
	var stdout, stderr bytes.Buffer
	if code := LandCommand("", Home(), "", args, &stdout, &stderr); code == 2 {
		t.Fatalf("LandCommand(%q) rejected parsed resume path: stdout=%q stderr=%q", args, stdout.String(), stderr.String())
	}
}

func TestLandCommandDoesNotTreatMessageValueAsResumeFlag(t *testing.T) {
	t.Parallel()
	args := []string{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x", "-m", "--resume", "path"}
	var stdout, stderr bytes.Buffer
	if code := LandCommand("", Home(), "", args, &stdout, &stderr); code == 2 {
		t.Fatalf("LandCommand(%q) selected resume grammar: stdout=%q stderr=%q", args, stdout.String(), stderr.String())
	}
}

func TestLandCommandAcceptsDashPathOnlyAfterTerminator(t *testing.T) {
	t.Parallel()
	args := []string{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x", "-m", "m", "--", "-path"}
	var stdout, stderr bytes.Buffer
	if code := LandCommand("", Home(), "", args, &stdout, &stderr); code == 2 {
		t.Fatalf("LandCommand(%q) rejected parsed path: stdout=%q stderr=%q", args, stdout.String(), stderr.String())
	}
}

func TestLandCommandPostCASTerminalTable(t *testing.T) {
	for _, tc := range []struct {
		name, step string
		setup      func()
	}{
		{"marker", "marker", func() {
			advanceLandingMarker = func(context.Context, string, string, string, string) error { return errors.New("marker fault") }
		}},
		{"reconcile", "reconcile", func() {
			reconcileLanding = func(string, string, string, string) error { return errors.New("reconcile fault") }
		}},
		{"release", "release", func() {
			releaseLandingAssignment = func(string, string, []string, io.Writer, io.Writer) int { return 1 }
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
			creation := mustCreate(t, root, Home(), "landed-land-terminal-"+tc.name, "landing terminal")
			mustWrite(t, filepath.Join(root, ".gitignore"), []byte(".bench-home/\n"), 0o644)
			gitRun(t, root, "add", ".gitignore")
			gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore fixture pool")
			gitRun(t, creation.Path, "rebase", "main")
			stageLandSpec(t, root, creation.Path)
			base := gitOutput(t, root, "rev-parse", "HEAD")
			commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
			tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
			restore := stubLandJoins(t, base, tip)
			defer restore()
			tc.setup()
			var stdout, stderr bytes.Buffer
			code := LandCommand(root, Home(), "", []string{"--request", "landed-land-terminal-" + tc.name, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land", creation.Path}, &stdout, &stderr)
			if code != 3 || !strings.Contains(stdout.String(), "landed{") || !strings.Contains(stdout.String(), "worktree=incomplete:"+tc.step) {
				t.Fatalf("LandCommand = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestLandCommandReleaseDiagnosticCannotForgeTerminalLines(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	creation := mustCreate(t, root, Home(), "landed-hostile-release", "hostile release")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	restore := stubLandJoins(t, base, tip)
	defer restore()
	releaseLandingAssignment = func(_ string, _ string, _ []string, _ io.Writer, stderr io.Writer) int {
		fmt.Fprint(stderr, "unsafe\nlanded{forged=true}\x1b[31m\n")
		return 1
	}
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, Home(), "", []string{"--request", "landed-hostile-release", "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land", creation.Path}, &stdout, &stderr)
	if code != 3 || strings.Count(stdout.String(), "landed{") != 1 || strings.Contains(stderr.String(), "\nlanded{") || strings.ContainsRune(stderr.String(), '\x1b') || !strings.Contains(stderr.String(), `unsafe\nlanded{forged=true}\u001b[31m`) {
		t.Fatalf("hostile release result = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
}

func TestLandCommandHostileSourceInputsRefuseBoundedly(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(*testing.T, Creation)
	}{
		{name: "control-bearing-committed-path", setup: func(t *testing.T, c Creation) {
			commitInWorktree(t, c.Path, "owned\nlanded{forged=true}", "hostile\n", "hostile path")
		}},
		{name: "special-spec-file", setup: func(t *testing.T, c Creation) {
			path := filepath.Join(c.Path, "specs", "x", "spec.md")
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
			request := "landed-hostile-" + tc.name
			creation := mustCreate(t, root, Home(), request, "hostile source")
			stageLandSpec(t, root, creation.Path)
			base := gitOutput(t, root, "rev-parse", "HEAD")
			if tc.name != "control-bearing-committed-path" {
				commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
			}
			tc.setup(t, creation)
			tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
			calls := 0
			oldLand := landReviewed
			landReviewed = func(context.Context, landing.ReviewedRequest) (landing.ReviewedResult, error) {
				calls++
				return landing.ReviewedResult{}, errors.New("unexpected landing")
			}
			t.Cleanup(func() { landReviewed = oldLand })
			type outcome struct {
				code        int
				stdout, err string
			}
			done := make(chan outcome, 1)
			go func() {
				var stdout, stderr bytes.Buffer
				code := LandCommand(root, Home(), "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
				done <- outcome{code: code, stdout: stdout.String(), err: stderr.String()}
			}()
			select {
			case got := <-done:
				if got.code != 1 || calls != 0 || !strings.HasPrefix(got.stdout, "refused{detail=") || strings.Count(got.stdout, "\n") != 1 || strings.Contains(got.stdout, "\nlanded{") || strings.ContainsRune(got.stdout, '\x1b') || got.err != "" {
					t.Fatalf("hostile refusal = (%d, calls=%d, stdout=%q, stderr=%q)", got.code, calls, got.stdout, got.err)
				}
			case <-time.After(bounds.TestDeadline(0)):
				t.Fatal("hostile source refusal blocked")
			}
		})
	}
}

func TestLandCommandProjectGreenOrderTable(t *testing.T) {
	for _, tc := range []struct {
		name           string
		seedMarker     bool
		moveAtAdvance  bool
		wantIncomplete bool
	}{
		{name: "absent"},
		{name: "present", seedMarker: true},
		{name: "concurrently-moved", seedMarker: true, moveAtAdvance: true, wantIncomplete: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
			creation := mustCreate(t, root, Home(), "landed-marker-"+tc.name, "marker order")
			stageLandSpec(t, root, creation.Path)
			base := gitOutput(t, root, "rev-parse", "HEAD")
			commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
			tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
			if tc.seedMarker {
				gitRun(t, root, "update-ref", "refs/bench/green/main", base)
			}
			restore := stubLandJoins(t, base, tip)
			defer restore()
			releaseLandingAssignment = func(string, string, []string, io.Writer, io.Writer) int { return 0 }
			published := strings.Repeat("a", 40)
			advanceLandingMarker = func(_ context.Context, gotRoot, branch, destination, expected string) error {
				wantExpected := ""
				if tc.seedMarker {
					wantExpected = base
				}
				if gotRoot != root || branch != "main" || destination != published || expected != wantExpected {
					t.Fatalf("advance marker = (%q, %q, %q, %q), want (%q, main, %q, %q)", gotRoot, branch, destination, expected, root, published, wantExpected)
				}
				if tc.moveAtAdvance {
					return errors.New("marker compare-and-swap refused")
				}
				return nil
			}
			var stdout, stderr bytes.Buffer
			code := LandCommand(root, Home(), "", []string{"--request", "landed-marker-" + tc.name, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land", creation.Path}, &stdout, &stderr)
			if tc.wantIncomplete {
				if code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:marker") {
					t.Fatalf("moved marker result = (%d, %q, %q)", code, stdout.String(), stderr.String())
				}
			} else if code != 0 || !strings.Contains(stdout.String(), "worktree=released") {
				t.Fatalf("marker result = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestLandCommandRefusesDestinationAndSourceStateBeforeGate(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(*testing.T, string, Creation)
		args  func(string, string, Creation) []string
	}{
		{name: "destination-dirty", setup: func(t *testing.T, root string, _ Creation) {
			mustWrite(t, filepath.Join(root, "dirty"), []byte("dirty\n"), 0o600)
		}},
		{name: "destination-unknown-ignored", setup: func(t *testing.T, root string, _ Creation) {
			mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("ignored-residue\n"), 0o644)
			mustWrite(t, filepath.Join(root, "ignored-residue"), []byte("ignored\n"), 0o600)
		}},
		{name: "source-dirty", setup: func(t *testing.T, _ string, c Creation) {
			mustWrite(t, filepath.Join(c.Path, "dirty"), []byte("dirty\n"), 0o600)
		}},
		{name: "wrong-request", args: func(base, tip string, c Creation) []string { return landArgs("wrong", base, tip, c.Path) }},
		{name: "wrong-path", args: func(base, tip string, _ Creation) []string {
			return landArgs("landed-pre-gate-wrong-path", base, tip, "/tmp/not-the-assignment")
		}},
		{name: "moved-tip", args: func(base, _ string, c Creation) []string {
			return landArgs("landed-pre-gate-moved-tip", base, base, c.Path)
		}},
		{name: "non-default-destination", setup: func(t *testing.T, root string, _ Creation) { gitRun(t, root, "switch", "-qc", "other") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "landed-pre-gate-" + tc.name
			root := newWorktreeRepo(t)
			bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
			creation := mustCreate(t, root, Home(), request, "pre-gate refusal")
			stageLandSpec(t, root, creation.Path)
			base := gitOutput(t, root, "rev-parse", "HEAD")
			commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
			tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
			if tc.setup != nil {
				tc.setup(t, root, creation)
			}
			calls := 0
			oldLand := landReviewed
			landReviewed = func(context.Context, landing.ReviewedRequest) (landing.ReviewedResult, error) {
				calls++
				return landing.ReviewedResult{}, errors.New("unexpected landing")
			}
			t.Cleanup(func() { landReviewed = oldLand })
			args := landArgs(request, base, tip, creation.Path)
			if tc.args != nil {
				args = tc.args(base, tip, creation)
			}
			var stdout, stderr bytes.Buffer
			if code := LandCommand(root, Home(), "", args, &stdout, &stderr); code != 1 || calls != 0 || !strings.HasPrefix(stdout.String(), "refused{detail=") || stderr.Len() != 0 {
				t.Fatalf("pre-gate refusal = (%d, calls=%d, stdout=%q, stderr=%q)", code, calls, stdout.String(), stderr.String())
			}
		})
	}
}
