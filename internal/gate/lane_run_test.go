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
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gittest"
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
	gittest.MustRun(t, root, "add", "-A")
	gittest.Commit(t, root, "module")
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")

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

// BG29 and BC51: the lane keeps its relay. A worktree commit still prints its red check's
// own lines, each under that check's prefix, from the lane log. A passing check's lines
// stay in the log, and a lane that hides every report prints nothing here.
func TestLaneRelaysARedCheckOwnLines(t *testing.T) {
	root := outcomeFixture(t)
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	var stdout, stderr bytes.Buffer
	result, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree, Stdout: &stdout, Stderr: &stderr,
		Checks: []Phase{
			{Name: "docs", Argv: []string{"sh", "-c", "echo docs passed"}},
			{Name: "prose", Argv: []string{"sh", "-c", "echo prose refused a marker phrase; echo on stderr too >&2; exit 1"}},
		},
	})
	if err != nil || result.Passed() || result.Check != "prose" {
		t.Fatalf("lane = %+v, %v; want a red prose check", result, err)
	}
	for _, want := range []string{"[prose] prose refused a marker phrase\n", "[prose] on stderr too\n"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("lane stdout = %q, want %q", stdout.String(), want)
		}
	}
	if strings.Contains(stdout.String()+stderr.String(), "[docs]") {
		t.Errorf("lane response = %q %q, want no line of the passing check", stdout.String(), stderr.String())
	}
	// Which stream reached the tap first is the copier goroutines' race, so the
	// diagnostic is pinned to being one of the check's own lines, not to which one.
	switch result.Diagnostic {
	case "prose refused a marker phrase", "on stderr too":
	default:
		t.Errorf("lane diagnostic = %q, want one of the check's own lines", result.Diagnostic)
	}
}

// BC51: a red check whose output ends with no newline still relays only its own line. Each
// check's last partial line ends at its close, so in the shared lane log the next check's
// text starts a line of its own.
func TestLaneRelaysARedCheckPartialLine(t *testing.T) {
	root := outcomeFixture(t)
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	var stdout bytes.Buffer
	result, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree, Stdout: &stdout,
		Checks: []Phase{
			{Name: "docs", Argv: []string{"sh", "-c", "printf 'docs passed'"}},
			{Name: "prose", Argv: []string{"sh", "-c", "printf refused; exit 1"}},
		},
	})
	if err != nil || result.Passed() || result.Check != "prose" {
		t.Fatalf("lane = %+v, %v; want a red prose check", result, err)
	}
	if want := "[prose] refused\n"; stdout.String() != want {
		t.Errorf("lane stdout = %q, want exactly %q", stdout.String(), want)
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
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
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
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
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
// a selected `vet` or `build` reds the lane. The pass and the three recorded invocations,
// `gate-prose`, the docs-currency check, and the prose-budget check, prove the Markdown
// change selected those three checks and no Go check.
func TestBenchkitLaneRunsOnlyTheSelectedChecks(t *testing.T) {
	root := outcomeFixture(t)
	base := gittest.Output(t, root, "rev-parse", "HEAD^{commit}")
	outcomeWrite(t, root, "note.md", "# Note\n", 0o644)
	gittest.MustRun(t, root, "add", "-A")
	gittest.Commit(t, root, "note")
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
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
	if !reflect.DeepEqual(result.Checks, []string{"prose", "docs-currency-workflow", "guidance-prose-budgets"}) || !reflect.DeepEqual(result.Classes, []string{"markdown"}) {
		t.Fatalf("result checks %v classes %v, want [prose docs-currency-workflow guidance-prose-budgets] and [markdown]", result.Checks, result.Classes)
	}
	invocations := strings.Split(strings.TrimSpace(string(outcomeRead(t, record))), "\n")
	if len(invocations) != 3 || !strings.Contains(invocations[0], "gate-prose") || !strings.Contains(invocations[1], "test --check docs-currency-workflow") || !strings.Contains(invocations[2], "test --check guidance-prose-budgets") {
		t.Fatalf("recorded invocations = %v, want one gate-prose run, one docs-currency run, and one prose-budget run", invocations)
	}
	if strings.Contains(strings.Join(invocations, " "), "gate-go") {
		t.Errorf("recorded invocations = %v, want no gate-go run: gofmt was not selected", invocations)
	}
}

// TestLaneRunBinaryHookKeepsFreshnessRefusal installs a build hook and keeps the real
// verifier. The stub executable carries no seal, so the lane selection must refuse it.
func TestLaneRunBinaryHookKeepsFreshnessRefusal(t *testing.T) {
	t.Setenv(runbinary.Env, "")
	if err := os.Unsetenv(runbinary.Env); err != nil {
		t.Fatal(err)
	}
	previous := laneRunBinary
	t.Cleanup(func() { laneRunBinary = previous })
	laneRunBinary = runbinary.Factory{Build: func(_ context.Context, _, output string) error {
		return os.WriteFile(output, []byte("#!/bin/sh\nexit 0\n"), 0o755)
	}}
	checks := []Phase{{Name: "prose", Argv: []string{runBinaryArgvToken, "gate-prose"}}}
	_, _, release, err := selectLaneRunBinary(t.Context(), LaneRequest{}, t.TempDir(), t.TempDir(), checks)
	if err == nil {
		t.Cleanup(func() { _ = release() })
		t.Fatal("lane selection accepted an unsealed executable, want the freshness refusal")
	}
	if !strings.Contains(err.Error(), "lane Bench executable unavailable") || !strings.Contains(err.Error(), "does not match source") {
		t.Fatalf("lane selection error = %q, want the verifier's refusal", err)
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

// TestSelectLaneByCanaryPin states the checks a pinned path adds. A pinned path of any
// content class selects the fixture bite check, and an unpinned Markdown path does not.
func TestSelectLaneByCanaryPin(t *testing.T) {
	pins := TreeFacts{FixturePins: []string{".agents/commands/x.md", "internal/x/y.go"}}
	for _, tc := range []struct {
		path            string
		checks, classes []string
	}{
		{".agents/commands/x.md", []string{"prose", "docs-currency-workflow", "guidance-prose-budgets", fixtureBiteCheck}, []string{"markdown", "canary-pin"}},
		{"internal/x/y.go", []string{"gofmt", "vet", "build", "structure", fixtureBiteCheck}, []string{"go-source", "canary-pin"}},
		{".agents/commands/unpinned.md", []string{"prose", "docs-currency-workflow", "guidance-prose-budgets"}, []string{"markdown"}},
	} {
		selected, classes := SelectLane(BenchkitLane("/repo", "/repo"), []ComposedChange{laneChange(tc.path)}, pins)
		if got := laneCheckNames(selected); !reflect.DeepEqual(got, tc.checks) || !reflect.DeepEqual(classes, tc.classes) {
			t.Errorf("%s selects %v with classes %v, want %v with %v", tc.path, got, classes, tc.checks, tc.classes)
		}
	}
}

// TestSelectLaneByWriteVerbSource states the check that a scanned write-verb file adds. A
// changed or deleted production Go file in each package that the route bypass test scans
// selects that check. A test file there, and a Go file in a subdirectory package that the
// scan does not read, select only the Go source checks.
func TestSelectLaneByWriteVerbSource(t *testing.T) {
	goSource := []string{"gofmt", "vet", "build", "structure"}
	for _, pkg := range registry.RouteBypassPackages {
		for _, tc := range []struct {
			change          ComposedChange
			checks, classes []string
		}{
			{laneChange(pkg + "/x.go"), append(slices.Clone(goSource), routeBypassCheck), []string{"go-source", "write-verb-source"}},
			{laneDeletion(pkg + "/x.go"), append(slices.Clone(goSource), routeBypassCheck), []string{"go-source", "write-verb-source"}},
			{laneChange(pkg + "/x_test.go"), goSource, []string{"go-source"}},
			{laneChange(pkg + "/unscanned/x.go"), goSource, []string{"go-source"}},
		} {
			selected, classes := SelectLane(BenchkitLane("/repo", "/repo"), []ComposedChange{tc.change}, TreeFacts{})
			if got := laneCheckNames(selected); !reflect.DeepEqual(got, tc.checks) || !reflect.DeepEqual(classes, tc.classes) {
				t.Errorf("%s %s selects %v with classes %v, want %v with %v", tc.change.Status, tc.change.Path, got, classes, tc.checks, tc.classes)
			}
		}
	}
}

// TestBenchkitLaneRunsTheFixtureBiteTestForAPinnedPath grades the pin derivation on the
// composed checkout. The fixture inventory in the tree pins `note.md`, so a change to that
// file selects the fixture bite check, and a red bite test refuses the commit under the
// check's own name. A lane that read a copied pin list, or no list, would pass.
func TestBenchkitLaneRunsTheFixtureBiteTestForAPinnedPath(t *testing.T) {
	root := outcomeFixture(t)
	fixture := filepath.Join("tests", "canary", "example-family", "planted")
	outcomeWrite(t, root, filepath.Join(fixture, "EXPECT"), "planted diagnostic\n", 0o644)
	outcomeWrite(t, root, filepath.Join(fixture, "MUTATE.json"), `[{"path":"note.md","old":"a","new":"b"}]`, 0o644)
	outcomeWrite(t, root, "note.md", "# Note\n\na\n", 0o644)
	outcomeCommit(t, root, "pinned")
	base := gittest.Output(t, root, "rev-parse", "HEAD^{commit}")
	outcomeWrite(t, root, "note.md", "# Note\n\na reflowed\n", 0o644)
	outcomeCommit(t, root, "reflow")
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	changes, err := ComposedChanges(root, base, tree)
	if err != nil {
		t.Fatal(err)
	}
	laneRefusesArgv(t, registry.FixtureBiteTest)

	result, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree, Lane: "benchkit", Selective: true,
		Checks: BenchkitLane(root, root), Changes: changes,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed() || result.Check != fixtureBiteCheck {
		t.Fatalf("lane outcome %s on check %q, want a fail on %s", result.Outcome, result.Check, fixtureBiteCheck)
	}
	if want := []string{"markdown", "canary-pin"}; !reflect.DeepEqual(result.Classes, want) {
		t.Errorf("lane classes = %v, want %v", result.Classes, want)
	}
}

// TestBenchkitLaneDocumentCheckFailsTheLane is PL35. A red document check refuses the
// commit under its own name, so the refusal names the rule the reader must satisfy. A
// document check run as optional would read this red as a pass.
func TestBenchkitLaneDocumentCheckFailsTheLane(t *testing.T) {
	const check = "roadmap-detail-integrity"
	root := outcomeFixture(t)
	base := gittest.Output(t, root, "rev-parse", "HEAD^{commit}")
	outcomeWrite(t, root, "ROADMAP.md", "# Roadmap\n", 0o644)
	gittest.MustRun(t, root, "add", "-A")
	gittest.Commit(t, root, "roadmap")
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
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
	if want := []string{"prose", "docs-currency-workflow", "guidance-prose-budgets", check}; !reflect.DeepEqual(result.Checks, want) {
		t.Errorf("lane checks = %v, want %v", result.Checks, want)
	}
}
