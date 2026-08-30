package gate

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
	"github.com/gibbonmi/bench/internal/gocache/cleanprobe"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// OG21, OG13: the kit root's lane is exactly the four declared checks, and the two
// Bench-owned ones name the run-binary token rather than an installed `bench`.
func TestBenchkitLaneTable(t *testing.T) {
	lane := BenchkitLane("/repo", "/repo")
	want := []Phase{
		{Name: "gofmt", Argv: []string{runBinaryArgvToken, "gate-go", "gofmt"}},
		{Name: "prose", Argv: []string{runBinaryArgvToken, "gate-prose", "/repo", "--", LaneNamedMarkdownToken}},
		{Name: "vet", Argv: []string{"go", "vet", "-trimpath", "./..."}},
		{Name: "build", Argv: []string{"go", "build", "-trimpath", disableBuildVCS, "./..."}},
	}
	if len(lane) < len(want) || !reflect.DeepEqual(lane[:len(want)], want) {
		t.Fatalf("kit lane = %+v, want it to open with %+v", lane, want)
	}
	// PL34: every check but the two toolchain ones is Bench-owned, so it names the token
	// rather than an installed `bench`. The document rows are Bench-owned too.
	for _, check := range lane {
		if check.Name == "vet" || check.Name == "build" {
			continue
		}
		if check.Argv[0] != runBinaryArgvToken {
			t.Errorf("lane check %s argv[0] = %q, want the run binary token %q", check.Name, check.Argv[0], runBinaryArgvToken)
		}
	}
}

// TestBenchkitLaneDocumentRowsFollowTheRegistry is PL28. The expectation is enumerated
// from registry.Checks, so a hand-written row list that misses a check the registry adds,
// or that carries one the registry binds to no document family, reds here.
func TestBenchkitLaneDocumentRowsFollowTheRegistry(t *testing.T) {
	documents := map[registry.InputSource]bool{
		registry.InputRoadmapBoard:      true,
		registry.InputDecisionDocuments: true,
		registry.InputCaptureRetros:     true,
		registry.InputBenchkitProfile:   true,
	}
	var want []Phase
	for _, check := range registry.Checks {
		if documents[check.Inputs] && check.RunsAt(registry.Dev) {
			want = append(want, Phase{Name: check.Name, Argv: []string{runBinaryArgvToken, "test", "--check", check.Name}})
		}
	}
	if len(want) == 0 {
		t.Fatal("the registry binds no dev-tier check to a document family")
	}
	var got []Phase
	for _, check := range BenchkitLane("/repo", "/repo") {
		if slices.Contains(check.Argv, "--check") {
			got = append(got, check)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("lane document rows = %+v, want the registry's own %+v", got, want)
	}
}

// OG32: the kit root's gate phase table keeps the whole-project test argv. The lane is
// a second table beside it, never a narrowing of this one.
func TestBenchkitPhasesKeepWholeProjectTestArgv(t *testing.T) {
	root := t.TempDir()
	writeLaneFile(t, filepath.Join(root, "go.mod"), "module example.com/x\n\ngo 1.24\n")
	for _, phase := range BenchkitPhases(root, root) {
		if phase.Name != "test" {
			continue
		}
		if want := []string{"go", "test", "-trimpath", "-count=1", "./..."}; !reflect.DeepEqual(phase.Argv, want) {
			t.Fatalf("test phase argv = %v, want %v", phase.Argv, want)
		}
		return
	}
	t.Fatal("the kit phase table declares no test phase")
}

// OG22: a manifest lane array is the root's lane, entry for entry.
func TestLaneForReadsTheManifestLane(t *testing.T) {
	root := t.TempDir()
	writeLaneFile(t, filepath.Join(root, ".bench", "phases.json"), `{
	  "phases": [{"name": "build", "argv": ["go", "build", "./..."]}],
	  "lane": [
	    {"name": "fmt", "argv": ["make", "fmt"]},
	    {"name": "unit", "argv": ["make", "unit"], "needs": ["fmt"], "dir": "sub"}
	  ]
	}`)
	lane, err := LaneFor(root, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	want := []Phase{
		{Name: "fmt", Argv: []string{"make", "fmt"}, Dir: root},
		{Name: "unit", Argv: []string{"make", "unit"}, Needs: []string{"fmt"}, Dir: filepath.Join(root, "sub")},
	}
	if !reflect.DeepEqual(lane, want) {
		t.Fatalf("manifest lane = %+v, want %+v", lane, want)
	}
}

// A manifest with no lane array, and a non-kit root with no manifest, declare no lane.
func TestLaneForWithoutADeclaration(t *testing.T) {
	withoutLane := t.TempDir()
	writeLaneFile(t, filepath.Join(withoutLane, ".bench", "phases.json"),
		`{"phases": [{"name": "build", "argv": ["go", "build", "./..."]}]}`)
	for _, root := range []string{withoutLane, t.TempDir()} {
		lane, err := LaneFor(root, t.TempDir())
		if err != nil {
			t.Fatalf("%s: %v", root, err)
		}
		if lane != nil {
			t.Fatalf("%s declared a lane %+v, want none", root, lane)
		}
	}
}

// The manifest defect behind OG24: a lane entry with an empty argv is refused with the
// loader's three-part diagnostic, exactly as a phase entry is.
func TestLaneForRefusesAMalformedEntry(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".bench", "phases.json")
	writeLaneFile(t, path, `{
	  "phases": [{"name": "build", "argv": ["go", "build", "./..."]}],
	  "lane": [{"name": "fmt", "argv": []}]
	}`)
	_, err := LaneFor(root, t.TempDir())
	if err == nil {
		t.Fatal("a lane entry with an empty argv was accepted")
	}
	if want := "gate: " + path + ": empty argv: fmt"; err.Error() != want {
		t.Fatalf("diagnostic = %q, want %q", err.Error(), want)
	}
}

// The prose placeholder resolves to the named paths, and a declared root is re-anchored
// to the private checkout that holds the composed tree.
func TestResolveLane(t *testing.T) {
	lane := BenchkitLane("/repo", "/repo")
	lane = append(lane, Phase{Name: "sub", Argv: []string{"make"}, Dir: filepath.Join("/repo", "sub")})
	resolved := resolveLane(lane, "/repo", "/checkout", []string{"a.md", "b.md"})

	prose := laneCheck(t, resolved, "prose")
	want := []string{runBinaryArgvToken, "gate-prose", "/checkout", "--", "a.md", "b.md"}
	if !reflect.DeepEqual(prose.Argv, want) {
		t.Errorf("prose argv = %v, want %v", prose.Argv, want)
	}
	if dir := laneCheck(t, resolved, "sub").Dir; dir != filepath.Join("/checkout", "sub") {
		t.Errorf("sub dir = %q, want the checkout counterpart", dir)
	}
	if vet := laneCheck(t, resolved, "vet"); !reflect.DeepEqual(vet.Argv, []string{"go", "vet", "-trimpath", "./..."}) {
		t.Errorf("vet argv = %v, want it unchanged", vet.Argv)
	}

	// An empty named list leaves the placeholder resolved away rather than passed on.
	empty := laneCheck(t, resolveLane(lane, "/repo", "/checkout", nil), "prose")
	if strings.Contains(strings.Join(empty.Argv, " "), LaneNamedMarkdownToken) {
		t.Errorf("prose argv = %v, want the placeholder resolved", empty.Argv)
	}
	if len(empty.Argv) != 4 {
		t.Errorf("prose argv = %v, want the four leading elements only", empty.Argv)
	}
}

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

func laneCheck(t *testing.T, lane []Phase, name string) Phase {
	t.Helper()
	for _, check := range lane {
		if check.Name == name {
			return check
		}
	}
	t.Fatalf("lane declares no check named %s", name)
	return Phase{}
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

// TestLaneForCommitMarksOnlyTheKitLaneSelective is PL22. The selection has one switch,
// and it follows the built-in-versus-manifest decision rather than a second rule.
func TestLaneForCommitMarksOnlyTheKitLaneSelective(t *testing.T) {
	kit := t.TempDir()
	t.Setenv("BENCH_KIT", kit)

	built, err := LaneForCommit(kit)
	if err != nil {
		t.Fatal(err)
	}
	if built == nil || !built.Selective {
		t.Fatalf("kit lane = %+v, want a selective lane", built)
	}

	project := t.TempDir()
	writeLaneFile(t, filepath.Join(project, ".bench", "phases.json"), `{
	  "phases": [{"name": "build", "argv": ["go", "build", "./..."]}],
	  "lane": [{"name": "fmt", "argv": ["make", "fmt"]}]
	}`)
	declared, err := LaneForCommit(project)
	if err != nil {
		t.Fatal(err)
	}
	if declared == nil || declared.Selective {
		t.Fatalf("manifest lane = %+v, want a lane that is not selective", declared)
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
