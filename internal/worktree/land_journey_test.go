// Real-git journey tests for the landing command: end-to-end publish, release, and destination-edit retention.
package worktree

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/handoffdoc"
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
			root, creation, base, tip, tally, _ := publicLandingFixture(t, request, tc.ignored, tc.declaration)
			if tc.emptyDeclaration || tc.foreignIgnored != "" {
				if tc.emptyDeclaration {
					mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[]}\n"), 0o644)
				}
				if tc.foreignIgnored != "" {
					mustWrite(t, filepath.Join(root, ".gitignore"), []byte(".logs/\nprivate/\n"), 0o644)
				}
				gitRun(t, root, "add", "-A")
				gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "configure ignored residue")
				base = gitOutput(t, root, "rev-parse", "HEAD")
				gitRun(t, creation.Path, "rebase", "main")
				tip = gitOutput(t, creation.Path, "rev-parse", "HEAD")
				if tc.foreignIgnored != "" {
					mustMkdirAll(t, filepath.Dir(filepath.Join(creation.Path, filepath.FromSlash(tc.foreignIgnored))), 0o755)
					mustWrite(t, filepath.Join(creation.Path, filepath.FromSlash(tc.foreignIgnored)), []byte("residue\n"), 0o600)
				}
			}
			disclosure := "landing source{review_base=" + base + ",assignment_start=" + creation.Assignment.Start + "}\n"
			var stdout, stderr bytes.Buffer
			cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path)
			cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
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
				if len(assignments) != 0 || !os.IsNotExist(statErr) || (!tc.runtime && stderr.String() != disclosure) || (tc.runtime && !strings.HasPrefix(stderr.String(), disclosure)) {
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
	markProof(t, "landing/journey/publish-release")
}

func TestLandCommandPublicPreservesHistoricalRuntimeLogs(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-runtime-logs"
	root, creation, _, _, tally, _ := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte(".logs/\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore runtime logs")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	mustMkdirAll(t, filepath.Join(root, ".logs"), 0o700)
	history := filepath.Join(root, ".logs", "history.jsonl")
	mustWrite(t, history, []byte("historical progress\n"), 0o600)

	var stdout, stderr bytes.Buffer
	cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path)
	cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
	if code := exitCode(cmd.Run()); code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("runtime-log landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if got, err := os.ReadFile(history); err != nil || string(got) != "historical progress\n" {
		t.Fatalf("historical progress log = %q, %v", got, err)
	}
	logs, err := filepath.Glob(filepath.Join(root, ".logs", "gate-*.jsonl"))
	if err != nil || len(logs) != 1 {
		t.Fatalf("gate progress logs = %q, %v", logs, err)
	}
	if got, err := os.ReadFile(logs[0]); err != nil || !strings.Contains(string(got), `"event":"gate.start"`) || !strings.Contains(string(got), `"event":"gate.finish"`) {
		t.Fatalf("durable gate progress log = %q, %v", got, err)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
}

func TestLandCommandRefusesPostGateUnknownIgnoredMutation(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-post-gate-ignored"
	root, creation, _, _, tally, _ := publicLandingFixture(t, request, "foreign-generated/output", "")
	mustWrite(t, filepath.Join(root, ".bench", "gate-prospective.sh"), []byte("#!/bin/sh\nset -eu\nruntime=$1\ngrep -q '^Status: implemented$' specs/x/spec.md\n[ -f owned.txt ]\nprintf g >> '"+tally+"'\nmkdir -p \"$LAND_DESTINATION/foreign-generated\"\nprintf injected > \"$LAND_DESTINATION/foreign-generated/output\"\n"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "gate-inputs.json"), []byte("{\"schema\":1,\"closure\":\"local\",\"environment\":[\"LAND_DESTINATION\"],\"paths\":[],\"tools\":[]}\n"), 0o644)
	gitRun(t, root, "add", ".bench/gate-prospective.sh", ".bench/gate-inputs.json")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "inject post-gate ignored mutation")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path)
	cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
	cmd.Env = append(os.Environ(), "LAND_DESTINATION="+root)
	if code := exitCode(cmd.Run()); code != 1 || !strings.Contains(stdout.String(), "landing destination checkout changed") {
		t.Fatalf("post-gate ignored mutation = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != base {
		t.Fatalf("post-gate mutation published main=%s, want %s", got, base)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("post-gate mutation gate tally = %q, %v", got, err)
	}
}

func TestLandCommandRetainsJustInTimeTrackedDestinationEdit(t *testing.T) {
	request := "land-last-moment-tracked-edit"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	commitInWorktree(t, root, "victim.txt", "saved\n", "track victim")
	base = gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip = gitOutput(t, creation.Path, "rev-parse", "HEAD")
	victim := filepath.Join(root, "victim.txt")
	injectLandingResetEdit(t, root, victim)

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	published := gitOutput(t, root, "rev-parse", "main")
	if code != 3 || !strings.Contains(stdout.String(), "published_commit="+published+",") || !strings.Contains(stdout.String(), "worktree=incomplete:reconcile") {
		t.Fatalf("last-moment tracked edit landing = (%d, %q, %q), want published incomplete reconciliation", code, stdout.String(), stderr.String())
	}
	if got, err := os.ReadFile(victim); err != nil || string(got) != "caller bytes\n" {
		t.Fatalf("last-moment tracked edit = %q, %v, want caller bytes", got, err)
	}
	assignments, err := intent.Assignments(root)
	if err != nil || len(assignments) != 1 || assignments[0].ID != creation.Assignment.ID || assignments[0].State != intent.StateActive {
		t.Fatalf("incomplete reconciliation retained assignments = %#v, %v", assignments, err)
	}
	if _, err := os.Stat(creation.Path); err != nil {
		t.Fatalf("incomplete reconciliation removed source assignment: %v", err)
	}
}

func TestLandCommandRetainsJustInTimeOverlappingDestinationEdit(t *testing.T) {
	request := "land-last-moment-overlapping-edit"
	root, creation, _, _, _, home := publicLandingFixture(t, request, "", "")
	commitInWorktree(t, root, "victim.txt", "saved\n", "track victim")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	mustWrite(t, filepath.Join(creation.Path, "victim.txt"), []byte("reviewed bytes\n"), 0o600)
	specPath := filepath.Join(creation.Path, "specs", "x", "spec.md")
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, specPath, append(specBytes, []byte("- `victim.txt`\n")...), 0o644)
	gitRun(t, creation.Path, "add", "victim.txt", "specs/x/spec.md")
	gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "review victim change")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	victim := filepath.Join(root, "victim.txt")
	injectLandingResetEdit(t, root, victim)

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	published := gitOutput(t, root, "rev-parse", "main")
	if code != 3 || !strings.Contains(stdout.String(), "published_commit="+published+",") || !strings.Contains(stdout.String(), "worktree=incomplete:reconcile") {
		t.Fatalf("last-moment overlapping edit landing = (%d, %q, %q), want published incomplete reconciliation", code, stdout.String(), stderr.String())
	}
	if got, err := os.ReadFile(victim); err != nil || string(got) != "caller bytes\n" {
		t.Fatalf("last-moment overlapping edit = %q, %v, want caller bytes", got, err)
	}
	assignments, err := intent.Assignments(root)
	if err != nil || len(assignments) != 1 || assignments[0].ID != creation.Assignment.ID || assignments[0].State != intent.StateActive {
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
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "private/output", "dist/")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 3 {
		t.Fatalf("published release exit = %d, want 3; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	tree := gitOutput(t, root, "rev-parse", published+"^{tree}")
	wantNext := "bench worktree land --resume '" + published + "' --request <request> --base '" + base + "' --source-tip '" + tip + "' --spec 'x' '" + creation.Path + "'"
	want := "landed{source_base=" + base + ",source_tip=" + tip + ",destination_base=" + base + ",published_commit=" + published + ",tree=" + tree + ",worktree=incomplete:release,next=" + wantNext + ",census=0}\n"
	if stdout.String() != want || strings.Contains(stdout.String(), request) {
		t.Fatalf("published release stdout = %q, want %q without caller token", stdout.String(), want)
	}
}

func TestLandCommandPublicConflictRepairRequiresNewReviewedTip(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-conflict-repair"
	root, creation, base, reviewedTip, tally, _ := publicLandingFixture(t, request, "", "")
	disclosure := "landing source{review_base=" + base + ",assignment_start=" + creation.Assignment.Start + "}\n"
	commitInWorktree(t, root, "owned.txt", "destination bytes\n", "destination conflict")
	destination := gitOutput(t, root, "rev-parse", "HEAD")
	run := func(tip string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land repaired source", creation.Path)
		cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}
	code, stdout, stderr := run(reviewedTip)
	if code != 1 || !strings.Contains(stdout, "refused{detail=composition conflict: textual,next=") || stderr != disclosure {
		t.Fatalf("conflict result = (%d, %q, %q)", code, stdout, stderr)
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) || gitOutput(t, root, "rev-parse", "HEAD") != destination || gitOutput(t, creation.Path, "rev-parse", "HEAD") != reviewedTip || gitOutput(t, root, "status", "--porcelain=v1") != "" || gitOutput(t, creation.Path, "status", "--porcelain=v1") != "" {
		t.Fatalf("conflict changed state or ran gate: tally=%v", err)
	}
	// LRS6: a source worktree that holds MERGE_HEAD is mid-merge, so the route names the
	// continuation of that merge and not a second one.
	mergeHead := filepath.Join(gitOutput(t, creation.Path, "rev-parse", "--absolute-git-dir"), "MERGE_HEAD")
	mustWrite(t, mergeHead, []byte(destination+"\n"), 0o644)
	code, stdout, stderr = run(reviewedTip)
	if code != 1 || !strings.Contains(stdout, "next=git -C '"+creation.Path+"' merge --continue") || strings.Contains(stdout, "then bench commit") {
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
	code, stdout, stderr = run(reviewedTip)
	if code != 1 || !strings.Contains(stdout, "; then bench commit; then /bench-review-implementation; then ") || strings.Contains(stdout, "merge --continue") {
		t.Fatalf("undecided merge state route = (%d, %q, %q), want the commit-and-review form", code, stdout, stderr)
	}
	if err := os.Remove(mergeHead); err != nil {
		t.Fatal(err)
	}
	merge := descendant(t, "git", "-C", creation.Path, "merge", "--no-commit", "main")
	if got, err := merge.CombinedOutput(); err == nil || !strings.Contains(string(got), "CONFLICT") {
		t.Fatalf("repair setup merge = %v, %s", err, got)
	}
	mustWrite(t, filepath.Join(creation.Path, "owned.txt"), []byte("destination bytes\nreviewed repair\n"), 0o644)
	gitRun(t, creation.Path, "add", "owned.txt")
	gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "repair conflict")
	repairedTip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	code, stdout, stderr = run(reviewedTip)
	// LRS17: the repair moved the source tip, so the refusal names both tips and routes the
	// operator to the caller's own command re-pointed at the tip the worktree now holds.
	want := "refused{detail=worktree source tip mismatch,observed=" + reviewedTip + ",wanted=" + repairedTip +
		",next=" + landingRerun(request, base, repairedTip, "x", creation.Path, creation.Assignment.ID) + "}\n"
	if code != 1 || stdout != want || stderr != "" {
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
	markProof(t, "landing/journey/conflict-refusal")
}

// OG14: a source the fast lane committed still pays the whole-project gate at the
// landing. The lane pass authorizes the worktree commit alone, so the tally records
// exactly one gate run, and that run is the landing's.
func TestLandGradesASourceCommittedByALanePass(t *testing.T) {
	binary := testRunBinary(t)
	request := "public-land-lane-source"
	root, creation, base, _, tally, _ := publicLandingFixture(t, request, "", "")
	// The kit-root selection must answer something other than this fixture, which is a
	// linked project and declares its lane in its own phase manifest.
	bindEnv(t, "BENCH_KIT", t.TempDir())
	manifest := filepath.Join(creation.Path, ".bench", "phases.json")
	mustWrite(t, manifest, []byte(`{"phases":[{"name":"build","argv":["true"]}],"lane":[{"name":"unit","argv":["true"]}]}`), 0o644)
	mustWrite(t, filepath.Join(creation.Path, "owned.txt"), []byte("lane bytes\n"), 0o644)

	var commitOut, commitErr bytes.Buffer
	commit := descendant(t, binary, "commit", "-m", "commit through the lane", "--", "owned.txt")
	commit.Dir, commit.Stdout, commit.Stderr = creation.Path, &commitOut, &commitErr
	if err := commit.Run(); err != nil || !strings.Contains(commitOut.String(), "lane{outcome=pass") {
		t.Fatalf("lane commit exit=%d stdout=%q stderr=%q", exitCode(err), commitOut.String(), commitErr.String())
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("the lane commit ran the whole-project gate (stat err %v)", err)
	}
	// The manifest is a fixture input, not landed bytes; the release refuses residue.
	if err := os.Remove(manifest); err != nil {
		t.Fatal(err)
	}
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	land := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land the laned source", creation.Path)
	land.Dir, land.Stdout, land.Stderr = root, &stdout, &stderr
	if err := land.Run(); err != nil || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("land exit=%d stdout=%q stderr=%q", exitCode(err), stdout.String(), stderr.String())
	}
	if recorded, err := os.ReadFile(tally); err != nil || string(recorded) != "g" {
		t.Fatalf("gate tally = %q, %v; want the landing to be the one whole-project gate run", recorded, err)
	}
	published := gitOutput(t, root, "rev-parse", "main")
	if gitOutput(t, root, "rev-parse", "refs/bench/green/main") != published {
		t.Fatal("the landing published without advancing the project-green marker")
	}
	if got := gitOutput(t, root, "show", published+":owned.txt"); got != "lane bytes" {
		t.Fatalf("published owned.txt = %q, want the lane-committed bytes", got)
	}
}

// BG22: the gate's bounded green shape reaches `bench worktree land`'s stdout byte for
// byte, ahead of the landing's own record. The landing relays the gate's writers, so it
// filters nothing out of the shape and puts nothing between it and the record.
func TestLandRelaysTheBoundedGreenShapeBeforeTheLandedRecord(t *testing.T) {
	t.Parallel()
	request := "public-land-canned-shape"
	root, creation, _, _, _, home := publicLandingFixture(t, request, "", "")
	base := commitCannedShapeGate(t, root, cannedGreenShape)
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)

	if code != 0 {
		t.Fatalf("land exit = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), cannedGreenShape) {
		t.Fatalf("stdout = %q, want it to open with the gate's bytes unchanged %q", stdout.String(), cannedGreenShape)
	}
	if rest := strings.TrimPrefix(stdout.String(), cannedGreenShape); !strings.HasPrefix(rest, "landed{") {
		t.Errorf("stdout after the shape = %q, want the landed record next", rest)
	}
}

// recordRawCalls appends n raw-call records for the assignment the pool path names.
// The fixture calls the recorder rather than write the file, so the test and the
// production writer keep one shape.
func recordRawCalls(t *testing.T, home, root, path string, n int) {
	t.Helper()
	recordRawCallsWithHead(t, home, root, path, "sed -i s/a/b/", n)
}

// recordRawCallsWithHead appends n raw-call records that one command text makes, which
// lets a test state a breakdown over more than one verb head.
func recordRawCallsWithHead(t *testing.T, home, root, path, command string, n int) {
	t.Helper()
	for range n {
		if err := census.Record(command+" "+filepath.Join(path, "owned.txt"), root, home, time.Now()); err != nil {
			t.Fatal(err)
		}
	}
}

// censusRecordPath names one assignment's record file.
func censusRecordPath(home, root, assignment string) string {
	return filepath.Join(census.Dir(home, root), assignment)
}

// TestLandCommandStatesTheCensusCountAndDropsTheRecords is EC20 and the landing half
// of EC24. The landed record carries the count as its last key, and the release step
// the landing runs leaves no record file for the retired assignment.
func TestLandCommandStatesTheCensusCountAndDropsTheRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-count"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	recordRawCalls(t, home, root, creation.Path, 3)
	survivor := seedHandoffSections(t, root, creation.Assignment)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.HasSuffix(stdout.String(), ",census=3}\n") {
		t.Fatalf("landed record = (%d, %q, %q), want census=3 as the last key", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(censusRecordPath(home, root, creation.Assignment.ID)); !os.IsNotExist(err) {
		t.Fatalf("the released landing kept the census record: %v", err)
	}
	requireHandoffSections(t, root, handoffdoc.MainKey, survivor)
}

// TestLandCommandStatesZeroForAnAssignmentWithNoRecords is EC21. Zero is a stated
// fact, and an absent record file is not a landing failure.
func TestLandCommandStatesZeroForAnAssignmentWithNoRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-zero"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.HasSuffix(stdout.String(), ",census=0}\n") {
		t.Fatalf("landed record = (%d, %q, %q), want census=0", code, stdout.String(), stderr.String())
	}
}

// TestLandCommandPrintsTheCensusHeadBreakdown proves the landing states the raw-call
// count for each verb head before the release step drops the records, so the retro
// reads the breakdown from the run. The heaviest head prints first.
func TestLandCommandPrintsTheCensusHeadBreakdown(t *testing.T) {
	t.Parallel()
	request := "census-landed-heads"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	recordRawCallsWithHead(t, home, root, creation.Path, "sed -i s/a/b/", 2)
	recordRawCallsWithHead(t, home, root, creation.Path, "awk -f x", 1)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stderr.String(), "census heads{sed=2,awk=1}\n") {
		t.Fatalf("landing evidence = (%d, %q, %q), want the head breakdown on stderr", code, stdout.String(), stderr.String())
	}
	if !strings.HasSuffix(stdout.String(), ",census=3}\n") {
		t.Fatalf("landed record = %q, want census=3 beside the breakdown", stdout.String())
	}
}

// TestLandCommandPrintsNoHeadsLineWithoutRecords proves an assignment that made no raw
// call prints no breakdown at all, and still states the zero count in its record.
func TestLandCommandPrintsNoHeadsLineWithoutRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-no-heads"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || strings.Contains(stderr.String(), "census heads{") || !strings.HasSuffix(stdout.String(), ",census=0}\n") {
		t.Fatalf("empty census landing = (%d, %q, %q), want no heads line and census=0", code, stdout.String(), stderr.String())
	}
}

// TestLandCommandRefusalKeepsTheCensusRecords proves a landing that refuses before
// its gate prints no landed record and drops nothing, so the operator can repair the
// invocation and land with the evidence intact.
func TestLandCommandRefusalKeepsTheCensusRecords(t *testing.T) {
	t.Parallel()
	request := "census-landed-refusal"
	root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
	recordRawCalls(t, home, root, creation.Path, 2)
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs("no-such-request", base, tip, creation.Path), &stdout, &stderr)
	if code == 0 || !strings.Contains(stdout.String(), "refused{") || strings.Contains(stdout.String(), "landed{") {
		t.Fatalf("refused landing = (%d, %q, %q), want a refusal and no landed record", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("the refusal ran the gate: %v", err)
	}
	if _, err := os.Stat(censusRecordPath(home, root, creation.Assignment.ID)); err != nil {
		t.Fatalf("the refusal dropped the census records: %v", err)
	}
}
