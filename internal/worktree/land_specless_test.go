// Spec-less landing tests: publication, gate refusal, base resolution, and resume without a spec.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// WL1, WL2, WL5, and WL10: the spec-less landing publishes and releases like a spec
// landing, keeps both reviewed parents and the marker, transitions no spec, and prints
// the same record.
func TestLandCommandSpecLessLandsPublishesAndReleases(t *testing.T) {
	t.Parallel()
	request := "spec-less-land"
	root, creation, base, tip, tally, home := specLessLandingFixture(t, request)
	specsBefore := gitOutput(t, creation.Path, "rev-parse", tip+":specs")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", specLessLandArgs(request, base, tip, creation.Path), &stdout, &stderr)
	published := gitOutput(t, root, "rev-parse", "main")
	tree := gitOutput(t, root, "rev-parse", published+"^{tree}")
	want := "landed{source_base=" + base + ",source_tip=" + tip + ",destination_base=" + base + ",published_commit=" + published + ",tree=" + tree + ",worktree=released}\n"
	if code != 0 || stdout.String() != want {
		t.Fatalf("spec-less land = (%d, %q, %q), want (0, %q)", code, stdout.String(), stderr.String(), want)
	}
	parents := strings.Fields(gitOutput(t, root, "rev-list", "--parents", "-n", "1", published))
	if len(parents) != 3 || parents[1] != base || parents[2] != tip {
		t.Fatalf("published parents = %q, want destination %s and source %s", parents, base, tip)
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got := gitOutput(t, root, "rev-parse", published+":specs"); got != specsBefore {
		t.Fatalf("published specs tree = %s, want the composed %s", got, specsBefore)
	}
	if got := gitOutput(t, root, "show", published+":specs/x/spec.md"); strings.Contains(got, "Status: implemented") {
		t.Fatalf("spec-less landing transitioned a spec: %q", got)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
	if _, err := os.Stat(creation.Path); !os.IsNotExist(err) {
		t.Fatalf("spec-less landing retained the worktree: %v", err)
	}
}

// WL3 and WL22: the gate still owns the spec-less path, and its refusal exits 1 with
// nothing published.
func TestLandCommandSpecLessGateRefusalPublishesNothing(t *testing.T) {
	t.Parallel()
	request := "spec-less-gate-red"
	root, creation, _, _, tally, home := specLessLandingFixture(t, request)
	mustWrite(t, filepath.Join(root, ".bench", "gate-prospective.sh"), []byte("#!/bin/sh\nset -eu\nruntime=$1\nprintf g >> '"+tally+"'\nexit 1\n"), 0o755)
	gitRun(t, root, "add", ".bench/gate-prospective.sh")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "red prospective gate")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", specLessLandArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.HasPrefix(stdout.String(), "refused{detail=prospective authorization refused") {
		t.Fatalf("spec-less gate refusal = (%d, %q, %q), want exit 1 and an authorization refusal", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != base {
		t.Fatalf("gate refusal published main=%s, want %s", got, base)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
}

// WL6: the identity proofs still run without a spec, and the refusal names both sides.
func TestLandCommandSpecLessRefusesSourceTipMismatch(t *testing.T) {
	t.Parallel()
	request := "spec-less-tip-mismatch"
	root, creation, base, tip, tally, home := specLessLandingFixture(t, request)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", specLessLandArgs(request, base, base, creation.Path), &stdout, &stderr)
	want := "refused{detail=worktree source tip mismatch,observed=" + base + ",wanted=" + tip + "}\n"
	if code != 1 || stdout.String() != want || stderr.Len() != 0 {
		t.Fatalf("spec-less tip mismatch = (%d, %q, %q), want (1, %q, empty)", code, stdout.String(), stderr.String(), want)
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("identity refusal ran the gate: %v", err)
	}
}

// WL23: the range proof still runs without a spec, and it refuses before the gate.
func TestLandCommandSpecLessRefusesNonAncestorBaseBeforeTheGate(t *testing.T) {
	t.Parallel()
	request := "spec-less-nonancestor-base"
	root, creation, _, tip, tally, home := specLessLandingFixture(t, request)
	commitInWorktree(t, root, "destination-only", "destination\n", "destination movement")
	unrelated := gitOutput(t, root, "rev-parse", "HEAD")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", specLessLandArgs(request, unrelated, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.Contains(stdout.String(), "reviewed source range is invalid") || !strings.Contains(stdout.String(), "not an ancestor") {
		t.Fatalf("spec-less non-ancestor base = (%d, %q, %q), want an invalid-range refusal", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("range refusal ran the gate: %v", err)
	}
}

// WL24: the landed record carries the resolved review base, never the flag's spelling.
func TestLandCommandSpecLessLandedSourceBaseIsTheResolvedBase(t *testing.T) {
	t.Parallel()
	request := "spec-less-resolved-base"
	root, creation, base, tip, _, home := specLessLandingFixture(t, request)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", specLessLandArgs(request, base[:12], tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "landed{source_base="+base+",") || strings.Contains(stdout.String(), base[:12]+",") {
		t.Fatalf("abbreviated spec-less base = (%d, %q, %q), want the resolved base in the record", code, stdout.String(), stderr.String())
	}
}

// WL7: a spec-less landing interrupted after publication resumes without a spec.
func TestResumeLandCommandSpecLessCompletesAnInterruptedLanding(t *testing.T) {
	t.Parallel()
	request := "spec-less-resume"
	root, creation, base, tip, tally, home := specLessLandingFixture(t, request)
	working := defaultJoins()
	broken := working
	broken.advanceLandingMarker = func(context.Context, string, string, string, string) error {
		return errors.New("injected marker interruption")
	}
	var stdout, stderr bytes.Buffer
	if code := landWith(broken, root, home, "", specLessLandArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:marker") {
		t.Fatalf("interrupted spec-less landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if strings.Contains(stdout.String(), "--spec") {
		t.Fatalf("spec-less resume instruction named a spec: %q", stdout.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, creation.Path}
	if code := landWith(working, root, home, "", args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "worktree=released}") || stderr.Len() != 0 {
		t.Fatalf("spec-less resume = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("spec-less resume reran the gate: tally=%q error=%v", got, err)
	}
}

// WL25: a resume without a spec completes a published spec-backed landing's marker and
// release, and publishes nothing a second time.
func TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding(t *testing.T) {
	t.Parallel()
	request := "spec-backed-spec-less-resume"
	root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
	working := defaultJoins()
	broken := working
	broken.advanceLandingMarker = func(context.Context, string, string, string, string) error {
		return errors.New("injected marker interruption")
	}
	var stdout, stderr bytes.Buffer
	if code := landWith(broken, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:marker") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, creation.Path}
	if code := landWith(working, root, home, "", args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "worktree=released}") || stderr.Len() != 0 {
		t.Fatalf("spec-less resume of a spec landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != published {
		t.Fatalf("resume republished: main=%s, want %s", got, published)
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("project-green = %s, want %s", got, published)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran the gate: tally=%q error=%v", got, err)
	}
}

// The edge under WL1: the flag is optional, but an empty value stays a usage error.
func TestLandCommandRefusesAnEmptySpecValue(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"first run", []string{"--request", "r", "--base", "b", "--source-tip", "s", "--spec", "", "-m", "m", "path"}},
		{"resume", []string{"--resume", "p", "--request", "r", "--base", "b", "--source-tip", "s", "--spec", "", "path"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := LandCommand("", Home(), "", tc.args, &stdout, &stderr)
			if code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), `--spec ""`) {
				t.Fatalf("empty --spec = (%d, %q, %q), want (2, empty, a usage line naming the empty value)", code, stdout.String(), stderr.String())
			}
		})
	}
}
