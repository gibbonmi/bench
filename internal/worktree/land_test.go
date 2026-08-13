package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
)

func TestLandCommandPublicRealGitJourney(t *testing.T) {
	binary := buildLandingBinary(t)
	for _, tc := range []struct {
		name, ignored, declaration, wantState string
	}{
		{name: "clean", wantState: "released"},
		{name: "declared-output", ignored: "dist/output", declaration: "dist/", wantState: "released"},
		{name: "unknown-ignored", ignored: "private/output", declaration: "dist/", wantState: "incomplete:release"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "public-land-" + tc.name
			root, creation, base, tip, tally := publicLandingFixture(t, request, tc.ignored, tc.declaration)
			disclosure := "landing source{review_base=" + base + ",assignment_start=" + creation.Assignment.Start + "}\n"
			var stdout, stderr bytes.Buffer
			cmd := exec.Command(binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path)
			cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
			err := cmd.Run()
			wantExit := 0
			if tc.wantState != "released" {
				wantExit = 1
			}
			if exitCode(err) != wantExit || !strings.Contains(stdout.String(), "worktree="+tc.wantState+"}") {
				t.Fatalf("land exit=%d stdout=%q stderr=%q", exitCode(err), stdout.String(), stderr.String())
			}
			published := gitOutput(t, root, "rev-parse", "main")
			parents := strings.Fields(gitOutput(t, root, "rev-list", "--parents", "-n", "1", published))
			if len(parents) != 3 || parents[1] != base || parents[2] != tip {
				t.Fatalf("published parents = %q, want destination %s and source %s", parents, base, tip)
			}
			if got := gitOutput(t, root, "show", published+":specs/x/spec.md"); !strings.Contains(got, "Status: implemented") {
				t.Fatalf("published spec = %q", got)
			}
			if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != published {
				t.Fatalf("project-green = %s, want %s", got, published)
			}
			if got, readErr := os.ReadFile(tally); readErr != nil || string(got) != "g" {
				t.Fatalf("gate tally = %q, %v", got, readErr)
			}
			if tc.name == "clean" {
				tree := gitOutput(t, root, "rev-parse", published+"^{tree}")
				if got := authorization.Authorize(t.Context(), root, tree); got.Kind != authorization.Green {
					t.Fatalf("identical-tree authorization = %+v", got)
				}
				if got, readErr := os.ReadFile(tally); readErr != nil || string(got) != "g" {
					t.Fatalf("identical tree reran gate: tally=%q error=%v", got, readErr)
				}
				commitInWorktree(t, root, "destination-only", "destination\n", "destination movement")
				changedTree := gitOutput(t, root, "rev-parse", "HEAD^{tree}")
				if got := authorization.Authorize(t.Context(), root, changedTree); got.Kind != authorization.Green {
					t.Fatalf("changed-tree authorization = %+v", got)
				}
				if got, readErr := os.ReadFile(tally); readErr != nil || string(got) != "gg" {
					t.Fatalf("changed tree did not rerun gate: tally=%q error=%v", got, readErr)
				}
			}
			assignments, readErr := intent.Assignments(root)
			if readErr != nil {
				t.Fatal(readErr)
			}
			_, statErr := os.Stat(creation.Path)
			if tc.wantState == "released" {
				if len(assignments) != 0 || !os.IsNotExist(statErr) || stderr.String() != disclosure {
					t.Fatalf("released state assignments=%#v stat=%v stderr=%q", assignments, statErr, stderr.String())
				}
			} else {
				if len(assignments) != 1 || statErr != nil || !strings.HasPrefix(stderr.String(), disclosure) || !strings.Contains(stderr.String(), "worktree retained (ignored)") || !strings.Contains(stderr.String(), "bench worktree release") {
					t.Fatalf("retained state assignments=%#v stat=%v stderr=%q", assignments, statErr, stderr.String())
				}
				if got, readErr := os.ReadFile(filepath.Join(creation.Path, filepath.FromSlash(tc.ignored))); readErr != nil || string(got) != "residue\n" {
					t.Fatalf("ignored residue = %q, %v", got, readErr)
				}
			}
		})
	}
}

func TestLandCommandPublicConflictRepairRequiresNewReviewedTip(t *testing.T) {
	binary := buildLandingBinary(t)
	request := "public-land-conflict-repair"
	root, creation, base, reviewedTip, tally := publicLandingFixture(t, request, "", "")
	disclosure := "landing source{review_base=" + base + ",assignment_start=" + creation.Assignment.Start + "}\n"
	commitInWorktree(t, root, "owned.txt", "destination bytes\n", "destination conflict")
	destination := gitOutput(t, root, "rev-parse", "HEAD")
	run := func(tip string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := exec.Command(binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land repaired source", creation.Path)
		cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}
	code, stdout, stderr := run(reviewedTip)
	if code != 1 || !strings.Contains(stdout, "refused{detail=composition conflict: textual}") || stderr != disclosure {
		t.Fatalf("conflict result = (%d, %q, %q)", code, stdout, stderr)
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) || gitOutput(t, root, "rev-parse", "HEAD") != destination || gitOutput(t, creation.Path, "rev-parse", "HEAD") != reviewedTip || gitOutput(t, root, "status", "--porcelain=v1") != "" || gitOutput(t, creation.Path, "status", "--porcelain=v1") != "" {
		t.Fatalf("conflict changed state or ran gate: tally=%v", err)
	}
	merge := exec.Command("git", "-C", creation.Path, "merge", "--no-commit", "main")
	if got, err := merge.CombinedOutput(); err == nil || !strings.Contains(string(got), "CONFLICT") {
		t.Fatalf("repair setup merge = %v, %s", err, got)
	}
	mustWrite(t, filepath.Join(creation.Path, "owned.txt"), []byte("destination bytes\nreviewed repair\n"), 0o644)
	gitRun(t, creation.Path, "add", "owned.txt")
	gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "repair conflict")
	repairedTip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	code, stdout, stderr = run(reviewedTip)
	if code != 1 || !strings.Contains(stdout, "refused{detail=worktree source tip mismatch}") || stderr != "" {
		t.Fatalf("old review after repair = (%d, %q, %q)", code, stdout, stderr)
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("old review ran gate: %v", err)
	}
	code, stdout, stderr = run(repairedTip)
	if code != 0 || !strings.Contains(stdout, "worktree=released") || stderr != disclosure {
		t.Fatalf("repaired landing = (%d, %q, %q)", code, stdout, stderr)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("repaired gate tally = %q, %v", got, err)
	}
}

func buildLandingBinary(t *testing.T) string {
	t.Helper()
	source := gitOutput(t, ".", "rev-parse", "--show-toplevel")
	output := filepath.Join(t.TempDir(), "bench")
	cmd := exec.Command("bash", filepath.Join(source, "scripts", "go-build.sh"), source, output)
	cmd.Dir = source
	if got, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build exact landing binary: %v\n%s", err, got)
	}
	return output
}

func publicLandingFixture(t *testing.T, request, ignored, declaration string) (string, Creation, string, string, string) {
	t.Helper()
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	common := gitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	tally := filepath.Join(common, "bench-land-gate-tally")
	t.Setenv("LAND_GATE_TALLY", tally)
	mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate.sh"), []byte("#!/bin/sh\nset -eu\nIFS= read -r status < specs/x/spec.md\n[ \"$status\" = \"Status: implemented\" ]\n[ -f owned.txt ]\nprintf g >> \"$LAND_GATE_TALLY\"\n"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-prospective.sh"), []byte("#!/bin/sh\nset -eu\nruntime=$1\nrg -q '^Status: implemented$' specs/x/spec.md\n[ -f owned.txt ]\nprintf g >> \"$LAND_GATE_TALLY\"\n"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-inputs.json"), []byte("{\"schema\":1,\"closure\":\"local\",\"environment\":[\"LAND_GATE_TALLY\"],\"paths\":[],\"tools\":[]}\n"), 0o644)
	if declaration != "" {
		mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[\""+declaration+"\"]}\n"), 0o644)
	}
	ignore := ""
	if ignored != "" {
		ignore = strings.Split(ignored, "/")[0] + "/\n"
		mustWrite(t, filepath.Join(root, ".gitignore"), []byte(ignore), 0o644)
	}
	specBody := "# x\n\nStatus: staged\n\n## User stories\n1. Land source.\n\n### Acceptance coverage map\n| row | story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|---|\n| LX1 | 1 | lands | command | command fails | catches failure |\n\n## Ownership fences\n\n- `owned.txt`\n"
	mustMkdirAll(t, filepath.Join(root, "specs", "x", "tickets"), 0o755)
	mustWrite(t, filepath.Join(root, "specs", "x", "spec.md"), []byte(specBody), 0o644)
	mustWrite(t, filepath.Join(root, "specs", "x", "tickets", "one.md"), []byte("Ticket covers LX1.\n"), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "landing base")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	creation := mustCreate(t, root, request, "public landing")
	commitInWorktree(t, creation.Path, "owned.txt", "reviewed bytes\n", "reviewed source")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	if ignored != "" {
		mustMkdirAll(t, filepath.Dir(filepath.Join(creation.Path, filepath.FromSlash(ignored))), 0o755)
		mustWrite(t, filepath.Join(creation.Path, filepath.FromSlash(ignored)), []byte("residue\n"), 0o600)
	}
	return root, creation, base, tip, tally
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}

func TestLandCommandNeverTreatsFlagValuesAsSourcePath(t *testing.T) {
	for _, args := range [][]string{
		{"--request", "would-be-path", "--base", "b", "--source-tip", "s", "--spec", "x", "-m", "m"},
		{"--request", "r", "--base", "would-be-path", "--source-tip", "s", "--spec", "x", "-m", "m"},
		{"--request", "r", "--base", "b", "--source-tip", "would-be-path", "--spec", "x", "-m", "m"},
		{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "would-be-path", "-m", "m"},
		{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x", "-m", "would-be-path"},
	} {
		var stdout, stderr bytes.Buffer
		if code := LandCommand("", args, &stdout, &stderr); code != 2 {
			t.Fatalf("LandCommand(%q) = %d, stdout=%q stderr=%q", args, code, stdout.String(), stderr.String())
		}
	}
}

func TestLandCommandAcceptsDashPathOnlyAfterTerminator(t *testing.T) {
	args := []string{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x", "-m", "m", "--", "-path"}
	var stdout, stderr bytes.Buffer
	if code := LandCommand("", args, &stdout, &stderr); code == 2 {
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
		{"reconcile", "reconcile", func() { reconcileLanding = func(string, string) error { return errors.New("reconcile fault") } }},
		{"release", "release", func() { releaseLandingAssignment = func(string, []string, io.Writer, io.Writer) int { return 1 } }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
			creation := mustCreate(t, root, "landed-land-terminal-"+tc.name, "landing terminal")
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
			code := LandCommand(root, []string{"--request", "landed-land-terminal-" + tc.name, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land", creation.Path}, &stdout, &stderr)
			if code != 1 || !strings.Contains(stdout.String(), "landed{") || !strings.Contains(stdout.String(), "worktree=incomplete:"+tc.step) {
				t.Fatalf("LandCommand = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestLandCommandReleaseDiagnosticCannotForgeTerminalLines(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	creation := mustCreate(t, root, "landed-hostile-release", "hostile release")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	restore := stubLandJoins(t, base, tip)
	defer restore()
	releaseLandingAssignment = func(_ string, _ []string, _ io.Writer, stderr io.Writer) int {
		fmt.Fprint(stderr, "unsafe\nlanded{forged=true}\x1b[31m\n")
		return 1
	}
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, []string{"--request", "landed-hostile-release", "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land", creation.Path}, &stdout, &stderr)
	if code != 1 || strings.Count(stdout.String(), "landed{") != 1 || strings.Contains(stderr.String(), "\nlanded{") || strings.ContainsRune(stderr.String(), '\x1b') || !strings.Contains(stderr.String(), `unsafe\nlanded{forged=true}\u001b[31m`) {
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
			t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
			request := "landed-hostile-" + tc.name
			creation := mustCreate(t, root, request, "hostile source")
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
				code := LandCommand(root, landArgs(request, base, tip, creation.Path), &stdout, &stderr)
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
			t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
			creation := mustCreate(t, root, "landed-marker-"+tc.name, "marker order")
			stageLandSpec(t, root, creation.Path)
			base := gitOutput(t, root, "rev-parse", "HEAD")
			commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
			tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
			if tc.seedMarker {
				gitRun(t, root, "update-ref", "refs/bench/green/main", base)
			}
			restore := stubLandJoins(t, base, tip)
			defer restore()
			releaseLandingAssignment = func(string, []string, io.Writer, io.Writer) int { return 0 }
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
			code := LandCommand(root, []string{"--request", "landed-marker-" + tc.name, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land", creation.Path}, &stdout, &stderr)
			if tc.wantIncomplete {
				if code != 1 || !strings.Contains(stdout.String(), "worktree=incomplete:marker") {
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
			t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
			creation := mustCreate(t, root, request, "pre-gate refusal")
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
			if code := LandCommand(root, args, &stdout, &stderr); code != 1 || calls != 0 || !strings.HasPrefix(stdout.String(), "refused{detail=") || stderr.Len() != 0 {
				t.Fatalf("pre-gate refusal = (%d, calls=%d, stdout=%q, stderr=%q)", code, calls, stdout.String(), stderr.String())
			}
		})
	}
}

func TestLandingDestinationAllowsOnlyDeclaredIgnoredOutput(t *testing.T) {
	root := newWorktreeRepo(t)
	mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("dist/\n"), 0o644)
	mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[\"dist/\"]}\n"), 0o644)
	gitRun(t, root, "add", ".gitignore", ".bench/build-outputs.json")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "declare output")
	mustMkdirAll(t, filepath.Join(root, "dist"), 0o755)
	mustWrite(t, filepath.Join(root, "dist", "bench"), []byte("output\n"), 0o755)
	tip, branch, marker, fingerprint, err := landingDestination(root)
	if err != nil || tip == "" || branch != "main" || marker != "" || fingerprint == "" {
		t.Fatalf("declared destination output = (%q, %q, %q, %q, %v)", tip, branch, marker, fingerprint, err)
	}
}

func landArgs(request, base, tip, path string) []string {
	return []string{"--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land", path}
}

func stageLandSpec(t *testing.T, root, source string) {
	t.Helper()
	path := filepath.Join(root, "specs", "x", "spec.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, path, []byte("Status: staged\n"), 0o644)
	gitRun(t, root, "add", "specs/x/spec.md")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "stage spec")
	gitRun(t, source, "rebase", "main")
}

func stubLandJoins(t *testing.T, base, tip string) func() {
	t.Helper()
	oldLand, oldMarker, oldReconcile, oldRelease, oldAuthorize := landReviewed, advanceLandingMarker, reconcileLanding, releaseLandingAssignment, authorizeLandingSource
	landReviewed = func(context.Context, landing.ReviewedRequest) (landing.ReviewedResult, error) {
		return landing.ReviewedResult{SourceBase: base, SourceTip: tip, DestinationBase: base, Commit: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40)}, nil
	}
	advanceLandingMarker = func(context.Context, string, string, string, string) error { return nil }
	reconcileLanding = func(string, string) error { return nil }
	releaseLandingAssignment = ReleaseCommand
	authorizeLandingSource = func(string, string, string) (diff.SourceRange, error) {
		return diff.SourceRange{Base: base, Tip: tip}, nil
	}
	return func() {
		landReviewed, advanceLandingMarker, reconcileLanding, releaseLandingAssignment, authorizeLandingSource = oldLand, oldMarker, oldReconcile, oldRelease, oldAuthorize
	}
}
