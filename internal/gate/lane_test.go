package gate

import (
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

// OG21, OG13, SR53: the kit root's lane opens with the five declared checks, and the
// Bench-owned ones name the run-binary token rather than an installed `bench`. The
// structure row names the growth flag and the base token, so the lane declares the
// ratchet's operand rather than a resolved revision.
func TestBenchkitLaneTable(t *testing.T) {
	lane := BenchkitLane("/repo", "/repo")
	want := []Phase{
		{Name: "gofmt", Argv: []string{runBinaryArgvToken, "gate-go", "gofmt"}},
		{Name: "prose", Argv: []string{runBinaryArgvToken, "gate-prose", "/repo", "--", LaneNamedMarkdownToken}},
		{Name: "vet", Argv: []string{"go", "vet", "-trimpath", "./..."}},
		{Name: "build", Argv: []string{"go", "build", "-trimpath", DisableBuildVCS, "./..."}},
		{Name: "structure", Argv: []string{runBinaryArgvToken, "structure", "--growth", LaneBaseToken}},
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

// The prose placeholder resolves to the named paths, the base placeholder resolves to the
// request's base (SR54), and a declared root is re-anchored to the private checkout that
// holds the composed tree.
func TestResolveLane(t *testing.T) {
	lane := BenchkitLane("/repo", "/repo")
	lane = append(lane, Phase{Name: "sub", Argv: []string{"make"}, Dir: filepath.Join("/repo", "sub")})
	resolved := resolveLane(lane, "/repo", "/checkout", []string{"a.md", "b.md"}, "abc123")

	// SR54: the base reaches the check as its own operand. An unreplaced token would hand
	// git the literal placeholder, which names no revision.
	structure := laneCheck(t, resolved, "structure")
	wantStructure := []string{runBinaryArgvToken, "structure", "--growth", "abc123"}
	if !reflect.DeepEqual(structure.Argv, wantStructure) {
		t.Errorf("structure argv = %v, want %v", structure.Argv, wantStructure)
	}

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
	empty := laneCheck(t, resolveLane(lane, "/repo", "/checkout", nil, "abc123"), "prose")
	if strings.Contains(strings.Join(empty.Argv, " "), LaneNamedMarkdownToken) {
		t.Errorf("prose argv = %v, want the placeholder resolved", empty.Argv)
	}
	if len(empty.Argv) != 4 {
		t.Errorf("prose argv = %v, want the four leading elements only", empty.Argv)
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
