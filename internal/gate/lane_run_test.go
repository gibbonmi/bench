package gate

// The lane's execution rows: each one drives a real RunLane or a real schedule against a
// fixture repository. The lane's declaration and resolution rows live in lane_test.go,
// and the shared laneCheck helper stays there with them.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gocache/cleanprobe"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// The defect behind the kit lane's build check. Go's VCS discovery treats a linked
// worktree's `.git` file as no root, walks up, and adopts any `.git` directory above the
// temporary checkout. Git refuses that directory, so the build fails with "error
// obtaining VCS status". The build carries -buildvcs=false, so a stray directory above
// TMPDIR grades nothing.
func TestBenchkitLaneBuildIgnoresAStrayGitDirAboveTheCheckout(t *testing.T) {
	stray := t.TempDir()
	if err := os.Mkdir(filepath.Join(stray, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", stray)

	root := outcomeFixture(t)
	outcomeWrite(t, root, "go.mod", "module example.com/x\n\ngo 1.24\n", 0o644)
	outcomeWrite(t, root, "main.go", "package main\n\nfunc main() {}\n", 0o644)
	outcomeGit(t, root, "add", "-A")
	outcomeGit(t, root, "commit", "-q", "-m", "module")
	tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")

	result, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree, Checks: []Phase{laneCheck(t, BenchkitLane(root, root), "build")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed() {
		t.Fatalf("build check %s: %s", result.Outcome, result.Diagnostic)
	}
}

// BG29: the lane keeps its relay. The engine buffers each phase's stream and prints
// failure rows instead, but a worktree commit still prints its red check's own lines,
// each under that check's prefix. A lane routed through the engine's buffer would print
// nothing here.
func TestLaneRelaysARedCheckOwnLines(t *testing.T) {
	diagnostics := &laneDiagnostics{first: map[string]string{}}
	var stdout, stderr bytes.Buffer
	checks := []Phase{{Name: "prose", Argv: []string{"sh", "-c", "echo prose refused a marker phrase; echo on stderr too >&2; exit 1"}}}
	results, cancelled := schedule(context.Background(), t.TempDir(), checks, laneWriters(&stdout, &stderr, diagnostics))

	if cancelled || len(results) != 1 || results[0].Code == 0 {
		t.Fatalf("lane schedule = %+v, cancelled %v; want one red check", results, cancelled)
	}
	if want := "[prose] prose refused a marker phrase\n"; stdout.String() != want {
		t.Errorf("lane stdout = %q, want %q", stdout.String(), want)
	}
	if want := "[prose] on stderr too\n"; stderr.String() != want {
		t.Errorf("lane stderr = %q, want %q", stderr.String(), want)
	}
	// Which stream reached the tap first is the copier goroutines' race, so the
	// diagnostic is pinned to being one of the check's own lines, not to which one.
	switch diagnostics.firstLine("prose", nil) {
	case "prose refused a marker phrase", "on stderr too":
	default:
		t.Errorf("lane diagnostic = %q, want one of the check's own lines", diagnostics.firstLine("prose", nil))
	}
}

// T03: the kit lane's two toolchain checks carry -trimpath, so a lane checkout writes no
// path-keyed archive. The literals are independent of the flag owner.
func TestBenchkitLaneToolchainChecksCarryTrimPath(t *testing.T) {
	lane := BenchkitLane("/repo", "/repo")
	for name, want := range map[string][]string{
		"vet":   {"go", "vet", "-trimpath", "./..."},
		"build": {"go", "build", "-trimpath", "-buildvcs=false", "./..."},
	} {
		if got := laneCheck(t, lane, name).Argv; !reflect.DeepEqual(got, want) {
			t.Errorf("lane %s argv = %v, want %v", name, got, want)
		}
	}
}

// L03: a lane run holds the shared cache lock across its checks, so `bench cache clean`
// exits 1 while a lane check is compiling. The probe is the lane's own check, which is the
// one point inside the lane's span a second process can observe.
func TestLaneRunHoldsTheCacheLockAcrossItsChecks(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	answerPath := filepath.Join(t.TempDir(), "clean-answer")
	t.Setenv(cleanprobe.Env, answerPath)

	root := outcomeFixture(t)
	tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")
	result, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree,
		Checks: []Phase{{Name: "probe", Argv: probeArgv(t)}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed() {
		t.Fatalf("probe check %s: %s", result.Outcome, result.Diagnostic)
	}
	requireRefusedClean(t, answerPath)
}

// TestFastLanePublishesTheOwnerRecord is PAR30. The fast lane materializes its own
// private checkout, so it publishes the same owner record the full gate publishes.
func TestFastLanePublishesTheOwnerRecord(t *testing.T) {
	tempRoot := t.TempDir()
	t.Setenv("TMPDIR", tempRoot)
	root := outcomeFixture(t)
	tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")
	snapshot := observeProspectiveOwnerRecord(t, root)

	result, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree, Checks: []Phase{{Name: "declared", Argv: []string{"true"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed() {
		t.Fatalf("lane check %s: %s", result.Outcome, result.Diagnostic)
	}
	requireProspectiveOwnerRecord(t, snapshot, root)
	requireNoProspectiveBundles(t, tempRoot)
}

// laneRecordsArgv replaces the lane's run-binary factory with a stub whose executable
// appends its own argv to record and exits 0. It is the seam that lets a row observe
// which Bench-owned checks ran without paying for a real build. An inherited selection
// would bypass the stub, so the variable is removed for the row's own span.
func laneRecordsArgv(t *testing.T, record string) {
	t.Helper()
	if _, present := os.LookupEnv(runbinary.Env); present {
		t.Setenv(runbinary.Env, "")
		if err := os.Unsetenv(runbinary.Env); err != nil {
			t.Fatal(err)
		}
	}
	previous := laneRunBinary
	t.Cleanup(func() { laneRunBinary = previous })
	laneRunBinary = runbinary.Factory{
		Build: func(_ context.Context, _, output string) error {
			script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> " + record + "\nexit 0\n"
			return os.WriteFile(output, []byte(script), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
}

// TestBenchkitLaneRunsOnlyTheSelectedChecks is PL19. The fixture holds no `go.mod`, so
// a selected `vet` or `build` reds the lane. The pass and the one recorded `gate-prose`
// invocation together prove the Markdown change selected the prose check alone.
func TestBenchkitLaneRunsOnlyTheSelectedChecks(t *testing.T) {
	root := outcomeFixture(t)
	base := outcomeGit(t, root, "rev-parse", "HEAD^{commit}")
	outcomeWrite(t, root, "note.md", "# Note\n", 0o644)
	outcomeGit(t, root, "add", "-A")
	outcomeGit(t, root, "commit", "-q", "-m", "note")
	tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")
	changes, err := ComposedChanges(root, base, tree)
	if err != nil {
		t.Fatal(err)
	}
	record := filepath.Join(t.TempDir(), "lane-argv")
	laneRecordsArgv(t, record)

	result, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree, Lane: "benchkit", Selective: true,
		Checks: BenchkitLane(root, root), Changes: changes,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed() {
		t.Fatalf("lane %s on check %s: %s", result.Outcome, result.Check, result.Diagnostic)
	}
	if !reflect.DeepEqual(result.Checks, []string{"prose"}) || !reflect.DeepEqual(result.Classes, []string{"markdown"}) {
		t.Fatalf("result checks %v classes %v, want [prose] and [markdown]", result.Checks, result.Classes)
	}
	invocations := strings.Split(strings.TrimSpace(string(outcomeRead(t, record))), "\n")
	if len(invocations) != 1 || !strings.Contains(invocations[0], "gate-prose") {
		t.Fatalf("recorded invocations = %v, want one gate-prose run", invocations)
	}
	if strings.Contains(strings.Join(invocations, " "), "gate-go") {
		t.Errorf("recorded invocations = %v, want no gate-go run: gofmt was not selected", invocations)
	}
}

// laneRefusesArgv is laneRecordsArgv's red twin: the stub exits 1 for the one invocation
// whose argv holds marker, and 0 for every other. It is the seam a row proves a red
// document check through without paying for a real conformance run.
func laneRefusesArgv(t *testing.T, marker string) {
	t.Helper()
	if _, present := os.LookupEnv(runbinary.Env); present {
		t.Setenv(runbinary.Env, "")
		if err := os.Unsetenv(runbinary.Env); err != nil {
			t.Fatal(err)
		}
	}
	previous := laneRunBinary
	t.Cleanup(func() { laneRunBinary = previous })
	laneRunBinary = runbinary.Factory{
		Build: func(_ context.Context, _, output string) error {
			script := "#!/bin/sh\ncase \"$*\" in\n  *'" + marker + "'*) echo 'the check refused' >&2; exit 1 ;;\nesac\nexit 0\n"
			return os.WriteFile(output, []byte(script), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
}

// TestBenchkitLaneDocumentCheckFailsTheLane is PL35. A red document check refuses the
// commit under its own name, so the refusal names the rule the reader must satisfy. A
// document check run as optional would read this red as a pass.
func TestBenchkitLaneDocumentCheckFailsTheLane(t *testing.T) {
	const check = "roadmap-detail-integrity"
	root := outcomeFixture(t)
	base := outcomeGit(t, root, "rev-parse", "HEAD^{commit}")
	outcomeWrite(t, root, "ROADMAP.md", "# Roadmap\n", 0o644)
	outcomeGit(t, root, "add", "-A")
	outcomeGit(t, root, "commit", "-q", "-m", "roadmap")
	tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")
	changes, err := ComposedChanges(root, base, tree)
	if err != nil {
		t.Fatal(err)
	}
	laneRefusesArgv(t, "test --check "+check)

	result, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree, Lane: "benchkit", Selective: true,
		Checks: BenchkitLane(root, root), Changes: changes,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed() || result.Check != check {
		t.Fatalf("lane outcome %s on check %q, want a fail on %s", result.Outcome, result.Check, check)
	}
	if want := []string{"prose", check}; !reflect.DeepEqual(result.Checks, want) {
		t.Errorf("lane checks = %v, want %v", result.Checks, want)
	}
}
