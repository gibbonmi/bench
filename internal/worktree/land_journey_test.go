// Real-git journey tests for the landing command: end-to-end publish, release, and destination-edit retention.
package worktree

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestLandCommandPublicRealGitJourney(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	for _, tc := range []struct {
		name, ignored, declaration, foreignIgnored, wantState string
		runtime, emptyDeclaration                             bool
	}{
		{name: "clean", wantState: "released"},
		{name: "declared-output", ignored: "dist/output", declaration: "dist/", wantState: "released"},
		{name: "runtime-log", ignored: ".logs/gate.jsonl", wantState: "released", runtime: true},
		{name: "runtime-log-empty-declaration", ignored: ".logs/gate.jsonl", wantState: "released", runtime: true, emptyDeclaration: true},
		{name: "runtime-and-unknown-ignored", ignored: ".logs/gate.jsonl", foreignIgnored: "private/output", wantState: "incomplete:release", runtime: true},
		{name: "unknown-ignored", ignored: "private/output", declaration: "dist/", wantState: "incomplete:release"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "public-land-" + tc.name
			f := publicLandingFixture(t, request, tc.ignored, tc.declaration)
			if tc.emptyDeclaration || tc.foreignIgnored != "" {
				if tc.emptyDeclaration {
					mustWrite(t, filepath.Join(f.root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[]}\n"), 0o644)
				}
				if tc.foreignIgnored != "" {
					mustWrite(t, filepath.Join(f.root, ".gitignore"), []byte(".logs/\nprivate/\n"), 0o644)
				}
				gitRun(t, f.root, "add", "-A")
				gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "configure ignored residue")
				f.base = gitOutput(t, f.root, "rev-parse", "HEAD")
				gitRun(t, f.creation.Path, "rebase", "main")
				refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
				f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
				if tc.foreignIgnored != "" {
					mustMkdirAll(t, filepath.Dir(filepath.Join(f.creation.Path, filepath.FromSlash(tc.foreignIgnored))), 0o755)
					mustWrite(t, filepath.Join(f.creation.Path, filepath.FromSlash(tc.foreignIgnored)), []byte("residue\n"), 0o600)
				}
			}
			disclosure := "landing source{review_base=" + f.base + ",assignment_start=" + f.creation.Assignment.Start + "}\n"
			var stdout, stderr bytes.Buffer
			cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", "-m", "land reviewed source", f.creation.Path)
			cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
			err := cmd.Run()
			wantExit := 0
			wantState := "worktree=" + tc.wantState + ",census=0}"
			if tc.wantState != "released" {
				wantExit = 3
				wantState = "worktree=" + tc.wantState + ",next="
			}
			if exitCode(err) != wantExit || !strings.Contains(stdout.String(), wantState) {
				t.Fatalf("land exit=%d stdout=%q stderr=%q", exitCode(err), stdout.String(), stderr.String())
			}
			published := gitOutput(t, f.root, "rev-parse", "main")
			parents := strings.Fields(gitOutput(t, f.root, "rev-list", "--parents", "-n", "1", published))
			if len(parents) != 3 || parents[1] != f.base || parents[2] != f.tip {
				t.Fatalf("published parents = %q, want destination %s and source %s", parents, f.base, f.tip)
			}
			if got := gitOutput(t, f.root, "show", published+":specs/x/spec.md"); !strings.Contains(got, "Status: implemented") {
				t.Fatalf("published spec = %q", got)
			}
			if got := gitOutput(t, f.root, "rev-parse", "refs/bench/green/main"); got != published {
				t.Fatalf("project-green = %s, want %s", got, published)
			}
			if got, readErr := os.ReadFile(f.tally); readErr != nil || string(got) != "g" {
				t.Fatalf("gate tally = %q, %v", got, readErr)
			}
			if tc.name == "clean" {
				tree := gitOutput(t, f.root, "rev-parse", published+"^{tree}")
				if got := authorization.Authorize(gate.WithCompletion(t.Context(), "specs/x/spec.md", f.tip), f.root, tree); got.Kind != authorization.Green {
					t.Fatalf("identical-tree authorization = %+v", got)
				}
				if got, readErr := os.ReadFile(f.tally); readErr != nil || string(got) != "g" {
					t.Fatalf("identical tree reran gate: tally=%q error=%v", got, readErr)
				}
				commitInWorktree(t, f.root, "destination-only", "destination\n", "destination movement")
				changedTree := gitOutput(t, f.root, "rev-parse", "HEAD^{tree}")
				if got := authorization.Authorize(t.Context(), f.root, changedTree); got.Kind != authorization.Green {
					t.Fatalf("changed-tree authorization = %+v", got)
				}
				if got, readErr := os.ReadFile(f.tally); readErr != nil || string(got) != "gg" {
					t.Fatalf("changed tree did not rerun gate: tally=%q error=%v", got, readErr)
				}
			}
			assignments, readErr := intent.Assignments(f.root)
			if readErr != nil {
				t.Fatal(readErr)
			}
			_, statErr := os.Stat(f.creation.Path)
			if tc.wantState == "released" {
				if len(assignments) != 0 || !os.IsNotExist(statErr) || (!tc.runtime && stderr.String() != disclosure) || (tc.runtime && !strings.HasPrefix(stderr.String(), disclosure)) {
					t.Fatalf("released state assignments=%#v stat=%v stderr=%q", assignments, statErr, stderr.String())
				}
			} else {
				if len(assignments) != 1 || statErr != nil || !strings.HasPrefix(stderr.String(), disclosure) || !strings.Contains(stderr.String(), "worktree retained (ignored)") || !strings.Contains(stderr.String(), "bench worktree release") {
					t.Fatalf("retained state assignments=%#v stat=%v stderr=%q", assignments, statErr, stderr.String())
				}
				if got, readErr := os.ReadFile(filepath.Join(f.creation.Path, filepath.FromSlash(tc.ignored))); readErr != nil || string(got) != "residue\n" {
					t.Fatalf("ignored residue = %q, %v", got, readErr)
				}
			}
		})
	}
	markProof(t, "landing/journey/publish-release")
}

func TestLandCommandPublicPreservesHistoricalRuntimeLogs(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-runtime-logs"
	f := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(f.root, ".gitignore"), []byte(".logs/\n"), 0o644)
	gitRun(t, f.root, "add", ".gitignore")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore runtime logs")
	base := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	mustMkdirAll(t, filepath.Join(f.root, ".logs"), 0o700)
	history := filepath.Join(f.root, ".logs", "history.jsonl")
	mustWrite(t, history, []byte("historical progress\n"), 0o600)

	var stdout, stderr bytes.Buffer
	cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", f.creation.Path)
	cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
	if code := exitCode(cmd.Run()); code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("runtime-log landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if got, err := os.ReadFile(history); err != nil || string(got) != "historical progress\n" {
		t.Fatalf("historical progress log = %q, %v", got, err)
	}
	logs, err := filepath.Glob(filepath.Join(f.root, ".logs", "gate-*.jsonl"))
	if err != nil || len(logs) != 1 {
		t.Fatalf("gate progress logs = %q, %v", logs, err)
	}
	if got, err := os.ReadFile(logs[0]); err != nil || !strings.Contains(string(got), `"event":"gate.start"`) || !strings.Contains(string(got), `"event":"gate.finish"`) {
		t.Fatalf("durable gate progress log = %q, %v", got, err)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
}

func TestLandCommandRefusesPostGateUnknownIgnoredMutation(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-post-gate-ignored"
	fixture := publicLandingFixture(t, request, "foreign-generated/output", "")
	f := landingGateFixture(t, "LAND_DESTINATION")
	f.MustWrite(t, fixture.root, "", "set -eu\nruntime=$1\n"+f.Command("grep")+" -q '^Status: implemented$' specs/x/spec.md\n[ -f owned.txt ]\nprintf g >> '"+fixture.tally+"'\n"+f.Command("mkdir")+" -p \"$LAND_DESTINATION/foreign-generated\"\nprintf injected > \"$LAND_DESTINATION/foreign-generated/output\"\n")
	gitRun(t, fixture.root, "add", ".bench/gate-prospective.sh", ".bench/gate-inputs.json")
	gitRun(t, fixture.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "inject post-gate ignored mutation")
	base := gitOutput(t, fixture.root, "rev-parse", "HEAD")
	gitRun(t, fixture.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, fixture.creation.Path, gitOutput(t, fixture.root, "rev-parse", "HEAD"))
	tip := gitOutput(t, fixture.creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", fixture.creation.Path)
	cmd.Dir, cmd.Stdout, cmd.Stderr = fixture.root, &stdout, &stderr
	cmd.Env = append(os.Environ(), "LAND_DESTINATION="+fixture.root)
	if code := exitCode(cmd.Run()); code != 1 || !strings.Contains(stdout.String(), "landing destination checkout changed") {
		t.Fatalf("post-gate ignored mutation = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, fixture.root, "rev-parse", "main"); got != base {
		t.Fatalf("post-gate mutation published main=%s, want %s", got, base)
	}
	if got, err := os.ReadFile(fixture.tally); err != nil || string(got) != "g" {
		t.Fatalf("post-gate mutation gate tally = %q, %v", got, err)
	}
}

func TestLandCommandRetainsJustInTimeTrackedDestinationEdit(t *testing.T) {
	request := "land-last-moment-tracked-edit"
	f := publicLandingFixture(t, request, "", "")
	commitInWorktree(t, f.root, "victim.txt", "saved\n", "track victim")
	f.base = gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	victim := filepath.Join(f.root, "victim.txt")
	injectLandingResetEdit(t, f.root, victim)

	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	published := gitOutput(t, f.root, "rev-parse", "main")
	if r.exit != 3 || !strings.Contains(r.stdout, "published_commit="+published+",") || !strings.Contains(r.stdout, "worktree=incomplete:reconcile") {
		t.Fatalf("last-moment tracked edit landing = (%d, %q, %q), want published incomplete reconciliation", r.exit, r.stdout, r.stderr)
	}
	if got, err := os.ReadFile(victim); err != nil || string(got) != "caller bytes\n" {
		t.Fatalf("last-moment tracked edit = %q, %v, want caller bytes", got, err)
	}
	assignments, err := intent.Assignments(f.root)
	if err != nil || len(assignments) != 1 || assignments[0].ID != f.creation.Assignment.ID || assignments[0].State != intent.StateActive {
		t.Fatalf("incomplete reconciliation retained assignments = %#v, %v", assignments, err)
	}
	if _, err := os.Stat(f.creation.Path); err != nil {
		t.Fatalf("incomplete reconciliation removed source assignment: %v", err)
	}
}

func TestLandCommandRetainsJustInTimeOverlappingDestinationEdit(t *testing.T) {
	request := "land-last-moment-overlapping-edit"
	f := publicLandingFixture(t, request, "", "")
	commitInWorktree(t, f.root, "victim.txt", "saved\n", "track victim")
	base := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	mustWrite(t, filepath.Join(f.creation.Path, "victim.txt"), []byte("reviewed bytes\n"), 0o600)
	specPath := filepath.Join(f.creation.Path, "specs", "x", "spec.md")
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, specPath, withFenceEntry(specBytes, "victim.txt"), 0o644)
	gitRun(t, f.creation.Path, "add", "victim.txt", "specs/x/spec.md")
	gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "review victim change")
	refreshLandingEvidence(t, f.creation.Path, base)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	victim := filepath.Join(f.root, "victim.txt")
	injectLandingResetEdit(t, f.root, victim)

	r := runVerb(t, verbLand, f.call(landArgs(request, base, tip, f.creation.Path)...))
	published := gitOutput(t, f.root, "rev-parse", "main")
	if r.exit != 3 || !strings.Contains(r.stdout, "published_commit="+published+",") || !strings.Contains(r.stdout, "worktree=incomplete:reconcile") {
		t.Fatalf("last-moment overlapping edit landing = (%d, %q, %q), want published incomplete reconciliation", r.exit, r.stdout, r.stderr)
	}
	if got, err := os.ReadFile(victim); err != nil || string(got) != "caller bytes\n" {
		t.Fatalf("last-moment overlapping edit = %q, %v, want caller bytes", got, err)
	}
	assignments, err := intent.Assignments(f.root)
	if err != nil || len(assignments) != 1 || assignments[0].ID != f.creation.Assignment.ID || assignments[0].State != intent.StateActive {
		t.Fatalf("incomplete overlapping reconciliation retained assignments = %#v, %v", assignments, err)
	}
}

func injectLandingResetEdit(t *testing.T, root, victim string) {
	t.Helper()
	shimDir := t.TempDir()
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(shimDir, "git"), []byte("#!/bin/sh\nset -eu\nwanted_mode=${LAND_RESET_MODE:---merge}\nsaw_reset=false\nsaw_mode=false\nsaw_destination=false\nfor arg in \"$@\"; do\n  [ \"$arg\" = reset ] && saw_reset=true\n  [ \"$arg\" = \"$wanted_mode\" ] && saw_mode=true\n  [ \"$arg\" = \"$LAND_RESET_DESTINATION\" ] && saw_destination=true\ndone\nif [ \"$saw_reset\" = true ] && [ \"$saw_mode\" = true ] && [ \"$saw_destination\" = true ]; then\n  printf 'caller bytes\\n' > \"$LAND_RESET_VICTIM\"\nfi\nexec "+realGit+" \"$@\"\n"), 0o755)
	bindEnv(t, "LAND_RESET_DESTINATION", root)
	bindEnv(t, "LAND_RESET_VICTIM", victim)
	bindEnv(t, "PATH", shimDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestLandCommandPublishedReleaseFailureExitsIncomplete(t *testing.T) {
	t.Parallel()
	request := "published-release-incomplete"
	f := publicLandingFixture(t, request, "private/output", "dist/")
	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 3 {
		t.Fatalf("published release exit = %d, want 3; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	tree := gitOutput(t, f.root, "rev-parse", published+"^{tree}")
	wantNext := "bench worktree land --resume '" + published + "' --request <request> --base '" + f.base + "' --source-tip '" + f.tip + "' --spec 'x' '" + f.creation.Path + "'"
	want := "landed{source_base=" + f.base + ",source_tip=" + f.tip + ",destination_base=" + f.base + ",published_commit=" + published + ",tree=" + tree + ",worktree=incomplete:release,next=" + wantNext + ",census=0}\n"
	if r.stdout != want || strings.Contains(r.stdout, request) {
		t.Fatalf("published release stdout = %q, want %q without caller token", r.stdout, want)
	}
}

func TestLandCommandPublicConflictRepairRequiresNewReviewedTip(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-conflict-repair"
	f := publicLandingFixture(t, request, "", "")
	disclosure := "landing source{review_base=" + f.base + ",assignment_start=" + f.creation.Assignment.Start + "}\n"
	commitInWorktree(t, f.root, "owned.txt", "destination bytes\n", "destination conflict")
	destination := gitOutput(t, f.root, "rev-parse", "HEAD")
	run := func(tip string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", f.base, "--source-tip", tip, "--spec", "x", "-m", "land repaired source", f.creation.Path)
		cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}
	code, stdout, stderr := run(f.tip)
	if code != 1 || !strings.Contains(stdout, "refused{detail=composition conflict: textual,next=") || stderr != disclosure {
		t.Fatalf("conflict result = (%d, %q, %q)", code, stdout, stderr)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) || gitOutput(t, f.root, "rev-parse", "HEAD") != destination || gitOutput(t, f.creation.Path, "rev-parse", "HEAD") != f.tip || gitOutput(t, f.root, "status", "--porcelain=v1") != "" || gitOutput(t, f.creation.Path, "status", "--porcelain=v1") != "" {
		t.Fatalf("conflict changed state or ran gate: tally=%v", err)
	}
	// LRS6: a source worktree that holds MERGE_HEAD is mid-merge, so the route names the
	// continuation of that merge and not a second one.
	mergeHead := filepath.Join(gitOutput(t, f.creation.Path, "rev-parse", "--absolute-git-dir"), "MERGE_HEAD")
	mustWrite(t, mergeHead, []byte(destination+"\n"), 0o644)
	code, stdout, stderr = run(f.tip)
	if code != 1 || !strings.Contains(stdout, "next=git -C '"+f.creation.Path+"' merge --continue") || strings.Contains(stdout, "then bench commit") {
		t.Fatalf("pending-merge conflict route = (%d, %q, %q), want the merge continuation", code, stdout, stderr)
	}
	// LRS22: an unreadable source Git directory leaves the merge state undecided, so the
	// route falls back to the commit-and-review form the committed resolution needs.
	if err := os.Remove(mergeHead); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(mergeHead, 0o755); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr = run(f.tip)
	if code != 1 || !strings.Contains(stdout, "; then bench commit; then /bench-review-implementation; then ") || strings.Contains(stdout, "merge --continue") {
		t.Fatalf("undecided merge state route = (%d, %q, %q), want the commit-and-review form", code, stdout, stderr)
	}
	if err := os.Remove(mergeHead); err != nil {
		t.Fatal(err)
	}
	merge := descendant(t, "git", "-C", f.creation.Path, "merge", "--no-commit", "main")
	if got, err := merge.CombinedOutput(); err == nil || !strings.Contains(string(got), "CONFLICT") {
		t.Fatalf("repair setup merge = %v, %s", err, got)
	}
	mustWrite(t, filepath.Join(f.creation.Path, "owned.txt"), []byte("destination bytes\nreviewed repair\n"), 0o644)
	gitRun(t, f.creation.Path, "add", "owned.txt")
	gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "repair conflict")
	refreshLandingEvidence(t, f.creation.Path, f.base)
	repairedTip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	code, stdout, stderr = run(f.tip)
	// LRS17: the repair moved the source tip, so the refusal names both tips and routes the
	// operator to the caller's own command re-pointed at the tip the worktree now holds.
	want := "refused{detail=worktree source tip mismatch,observed=" + f.tip + ",wanted=" + repairedTip +
		",next=" + landingRerun(request, f.base, repairedTip, "x", f.creation.Path, f.creation.Assignment.ID) + "}\n"
	if code != 1 || stdout != want || stderr != "" {
		t.Fatalf("old review after repair = (%d, %q, %q)", code, stdout, stderr)
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("old review ran gate: %v", err)
	}
	code, stdout, stderr = run(repairedTip)
	if code != 0 || !strings.Contains(stdout, "worktree=released") || stderr != disclosure {
		t.Fatalf("repaired landing = (%d, %q, %q)", code, stdout, stderr)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("repaired gate tally = %q, %v", got, err)
	}
	markProof(t, "landing/journey/conflict-refusal")
}

// OG14: a source the fast lane committed still pays the whole-project gate at the
// landing. The lane pass authorizes the worktree commit alone, so the tally records
// exactly one gate run, and that run is the landing's.
func TestLandGradesASourceCommittedByALanePass(t *testing.T) {
	binary := testRunBinary(t)
	request := "public-land-lane-source"
	f := publicLandingFixture(t, request, "", "")
	// The kit-root selection must answer something other than this fixture, which is a
	// linked project and declares its lane in its own phase manifest.
	bindEnv(t, "BENCH_KIT", t.TempDir())
	manifest := filepath.Join(f.creation.Path, ".bench", "phases.json")
	mustWrite(t, manifest, []byte(`{"phases":[{"name":"build","argv":["true"]}],"lane":[{"name":"unit","argv":["true"]}]}`), 0o644)
	mustWrite(t, filepath.Join(f.creation.Path, "owned.txt"), []byte("lane bytes\n"), 0o644)

	var commitOut, commitErr bytes.Buffer
	commit := descendant(t, binary, "commit", "-m", "commit through the lane", "--", "owned.txt")
	commit.Dir, commit.Stdout, commit.Stderr = f.creation.Path, &commitOut, &commitErr
	if err := commit.Run(); err != nil || !strings.Contains(commitOut.String(), "lane{outcome=pass") {
		t.Fatalf("lane commit exit=%d stdout=%q stderr=%q", exitCode(err), commitOut.String(), commitErr.String())
	}
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("the lane commit ran the whole-project gate (stat err %v)", err)
	}
	// The manifest is a fixture input, not landed bytes; the release refuses residue.
	if err := os.Remove(manifest); err != nil {
		t.Fatal(err)
	}
	refreshLandingEvidence(t, f.creation.Path, f.base)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	land := descendant(t, binary, "worktree", "land", "--request", request, "--base", f.base, "--source-tip", tip, "--spec", "x", "-m", "land the laned source", f.creation.Path)
	land.Dir, land.Stdout, land.Stderr = f.root, &stdout, &stderr
	if err := land.Run(); err != nil || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("land exit=%d stdout=%q stderr=%q", exitCode(err), stdout.String(), stderr.String())
	}
	if recorded, err := os.ReadFile(f.tally); err != nil || string(recorded) != "g" {
		t.Fatalf("gate tally = %q, %v; want the landing to be the one whole-project gate run", recorded, err)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	if gitOutput(t, f.root, "rev-parse", "refs/bench/green/main") != published {
		t.Fatal("the landing published without advancing the project-green marker")
	}
	if got := gitOutput(t, f.root, "show", published+":owned.txt"); got != "lane bytes" {
		t.Fatalf("published owned.txt = %q, want the lane-committed bytes", got)
	}
}

// BG22: the gate's bounded green shape reaches `bench worktree land`'s stdout byte for
// byte, ahead of the landing's own record. The landing relays the gate's writers, so it
// filters nothing out of the shape and puts nothing between it and the record.
func TestLandRelaysTheBoundedGreenShapeBeforeTheLandedRecord(t *testing.T) {
	t.Parallel()
	request := "public-land-canned-shape"
	f := publicLandingFixture(t, request, "", "")
	base := commitCannedShapeGate(t, f.root, cannedGreenShape)
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	r := runVerb(t, verbLand, f.call(landArgs(request, base, tip, f.creation.Path)...))

	if r.exit != 0 {
		t.Fatalf("land exit = %d, want 0; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	if !strings.HasPrefix(r.stdout, cannedGreenShape) {
		t.Fatalf("stdout = %q, want it to open with the gate's bytes unchanged %q", r.stdout, cannedGreenShape)
	}
	if rest := strings.TrimPrefix(r.stdout, cannedGreenShape); !strings.HasPrefix(rest, wantEffects("not-applicable")+"landed{") {
		t.Errorf("stdout after the shape = %q, want the effects table and then the landed record", rest)
	}
}

// recordRawCalls appends n raw-call records for the assignment the pool path names.
// The fixture calls the recorder rather than write the file, so the test and the
// production writer keep one shape.
