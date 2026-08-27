package gate

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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
	if !reflect.DeepEqual(lane, want) {
		t.Fatalf("kit lane = %+v, want %+v", lane, want)
	}
	for _, name := range []string{"gofmt", "prose"} {
		check := laneCheck(t, lane, name)
		if check.Argv[0] != runBinaryArgvToken {
			t.Errorf("lane check %s argv[0] = %q, want the run binary token %q", name, check.Argv[0], runBinaryArgvToken)
		}
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
// path-keyed archive. The literals are independent of the flag owner on purpose: a lane
// table that reads the flag from somewhere else, or drops it, reds here.
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
