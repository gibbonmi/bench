// Stable-owner tests for the prospective gate: the graded tree supplies the gate
// executable, the landing baseline supplies the phase schedule, and prospective
// evidence keys to the graded tree together with the baseline runner identity.
package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// stubProspectiveBuilder replaces the prospective gate's build and verification seams
// with recording stubs, and returns the source roots each authored selection was built
// from, in call order.
func stubProspectiveBuilder(t *testing.T) *[]string {
	t.Helper()
	sources := []string{}
	old := prospectiveRunBinary
	prospectiveRunBinary = runbinary.Factory{
		Build: func(_ context.Context, sourceRoot, output string) error {
			sources = append(sources, sourceRoot)
			return os.WriteFile(output, []byte("#!/bin/sh\nexit 0\n"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
	t.Cleanup(func() { prospectiveRunBinary = old })
	return &sources
}

// TestProspectiveGateAuthorsItsExecutableFromTheGradedTree is SOL09 and the executable
// half of SOL15. An inherited selection is present and was built from another tree, so
// honoring it would record a source digest the graded subject never produced. The gate
// authors its own executable from the graded checkout instead, and that private
// executable is removable, so no run leaves it behind.
func TestProspectiveGateAuthorsItsExecutableFromTheGradedTree(t *testing.T) {
	checkout := t.TempDir()
	kit := t.TempDir()
	inherited := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(inherited, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(runbinary.Env, inherited)
	sources := stubProspectiveBuilder(t)
	owner := prospectiveRunBinaryOwner(checkout)

	graded, err := owner(t.Context(), checkout)
	if err != nil {
		t.Fatal(err)
	}
	if len(*sources) != 1 || (*sources)[0] != checkout {
		t.Fatalf("graded-tree sources = %#v, want one selection built from %q", *sources, checkout)
	}
	if graded.Path == inherited {
		t.Fatalf("graded subject ran under the inherited selection %q", inherited)
	}
	if err := graded.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Dir(graded.Path)); !os.IsNotExist(err) {
		t.Fatalf("private gate executable survived its selection: %v", err)
	}

	// A source that is not the graded tree is the baseline's own kit, whose inherited
	// selection is already the baseline runner.
	baseline, err := owner(t.Context(), kit)
	if err != nil {
		t.Fatal(err)
	}
	if len(*sources) != 1 || baseline.Path != inherited {
		t.Fatalf("baseline selection = %q with sources %#v, want the inherited %q", baseline.Path, *sources, inherited)
	}
}

// TestProspectiveGateKeepsTheBaselinePhaseScheduleAfterACandidateOmission is SOL10. The
// candidate tree declares a phase manifest that omits a phase the baseline declares.
// The baseline's schedule is the one the run gets.
func TestProspectiveGateKeepsTheBaselinePhaseScheduleAfterACandidateOmission(t *testing.T) {
	baseline := t.TempDir()
	candidate := t.TempDir()
	writeGateFixtureFile(t, baseline, ".bench/phases.json", `{"phases":[{"name":"planted","argv":["true"]},{"name":"kept","argv":["true"]}]}`, 0o644)
	writeGateFixtureFile(t, candidate, ".bench/phases.json", `{"phases":[{"name":"kept","argv":["true"]}]}`, 0o644)
	t.Setenv(baselinePolicyEnv, baseline)

	phases, err := phaseTable(candidate, candidate)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"planted", "kept"}; strings.Join(phaseNames(phases), ",") != strings.Join(want, ",") {
		t.Fatalf("phase names = %v, want the baseline's %v", phaseNames(phases), want)
	}
}

// TestProspectivePhaseScheduleRefusesAnUnusableBaseline keeps the baseline selection
// fail-closed. A named baseline the run cannot read refuses; a fall back to the
// candidate tree would reinstate the omission SOL10 excludes.
func TestProspectivePhaseScheduleRefusesAnUnusableBaseline(t *testing.T) {
	candidate := t.TempDir()
	writeGateFixtureFile(t, candidate, ".bench/phases.json", `{"phases":[{"name":"kept","argv":["true"]}]}`, 0o644)
	t.Setenv(baselinePolicyEnv, filepath.Join(t.TempDir(), "absent"))

	if _, err := phaseTable(candidate, candidate); err == nil {
		t.Fatal("unusable baseline fell back to the candidate schedule")
	}
}

// TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity is SOL13. Retained
// prospective evidence is addressed by the graded tree together with the baseline
// runner identity. A key that binds only one of them reuses evidence after the other
// half of the subject changes.
func TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity(t *testing.T) {
	root := t.TempDir()
	writeGateFixtureFile(t, root, ".bench/gate.sh", "#!/bin/sh\nexit 0\n", 0o755)
	writeGateFixtureFile(t, root, ".bench/gate-prospective.sh", "#!/bin/sh\nexit 0\n", 0o755)
	writeGateFixtureFile(t, root, ".bench/gate-inputs.json", "{\"schema\":1,\"closure\":\"local\",\"environment\":[],\"paths\":[],\"tools\":[]}\n", 0o644)
	// One identity root answers for both runners in turn, so only the runner identity
	// differs between the two keys. Two directories would differ by their own paths,
	// which the subject already binds, and the assertion would pass without the runner.
	identityRoot := t.TempDir()
	tree := strings.Repeat("a", 40)
	other := strings.Repeat("b", 40)

	key := func(tree string) string {
		t.Helper()
		plan, err := buildSubjectForTree(root, identityRoot, policyVersion, tree)
		if err != nil {
			t.Fatal(err)
		}
		return evidenceName(plan)
	}
	unbound := key(tree)
	if unbound != key(tree) {
		t.Fatal("an unchanged subject did not address the same evidence")
	}
	if unbound == key(other) {
		t.Fatal("evidence key ignored the graded tree")
	}
	if got := baselineRunnerIdentity(identityRoot); got != unboundBaselineRunner {
		t.Fatalf("undeclared baseline runner identity = %q, want %q", got, unboundBaselineRunner)
	}
	writeGateFixtureFile(t, identityRoot, "scripts/go-build.inputs", "build_script=scripts/go-build.sh\n", 0o644)
	if got := baselineRunnerIdentity(identityRoot); got == unboundBaselineRunner {
		t.Fatal("a declared build recipe still answered the unbound runner identity")
	}
	if unbound == key(tree) {
		t.Fatal("evidence key ignored the baseline runner identity")
	}
}

func writeGateFixtureFile(t *testing.T, root, rel, body string, mode os.FileMode) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}

// TestProspectiveGateReportsAFailedBuildWithNoResidue is the fail-closed half of the
// graded-tree selection. Authoring the private gate executable fails, so the owner
// returns the build cause instead of a selection, and the private directory it opened
// for that executable does not survive the failure.
func TestProspectiveGateReportsAFailedBuildWithNoResidue(t *testing.T) {
	checkout := t.TempDir()
	tempRoot := t.TempDir()
	old := prospectiveRunBinary
	prospectiveRunBinary = runbinary.Factory{
		TempRoot: tempRoot,
		Build:    func(context.Context, string, string) error { return errors.New("build refused") },
		Verify:   func(string, string) error { return nil },
	}
	t.Cleanup(func() { prospectiveRunBinary = old })

	selection, err := prospectiveRunBinaryOwner(checkout)(t.Context(), checkout)
	if err == nil || selection != nil {
		t.Fatalf("failed prospective build = (%#v, %v), want no selection and an error", selection, err)
	}
	if !strings.Contains(err.Error(), "build refused") {
		t.Fatalf("prospective build error = %v, want the build cause", err)
	}
	entries, readErr := os.ReadDir(tempRoot)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("failed prospective build left %d entries under its temporary root", len(entries))
	}
}

func TestTimedOutProspectiveProducerLeavesNoBundle(t *testing.T) {
	root, tree, tempRoot := prospectiveProducerFixture(t, "#!/bin/sh\nset -eu\nwhile :; do sleep 1; done\n")
	oldTimeout := gateTimeout
	gateTimeout = 100 * time.Millisecond
	t.Cleanup(func() { gateTimeout = oldTimeout })

	var stdout, stderr bytes.Buffer
	result := executeTreeWithOwner(context.Background(), root, tree, &stdout, &stderr, inertProspectiveSelection)
	if result.ActionExit != 124 {
		t.Fatalf("timed-out prospective result = %#v, stderr=%q", result, stderr.String())
	}
	requireNoProspectiveBundles(t, tempRoot)
}

func TestCancelledProspectiveProducerLeavesNoBundle(t *testing.T) {
	root, tree, tempRoot := prospectiveProducerFixture(t, "#!/bin/sh\nset -eu\nprintf ready > \"$PAR_READY\"\nwhile :; do sleep 1; done\n")
	ready := filepath.Join(t.TempDir(), "ready")
	t.Setenv("PAR_READY", ready)
	writeGateFixtureFile(t, root, ".bench/gate-inputs.json", "{\"schema\":1,\"closure\":\"local\",\"environment\":[\"HOME\",\"PAR_READY\"],\"paths\":[],\"tools\":[]}\n", 0o644)
	outcomeCommit(t, root, "declare cancellation barrier")
	tree = outcomeGit(t, root, "write-tree")

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan Result, 1)
	go func() {
		result <- executeTreeWithOwner(ctx, root, tree, &bytes.Buffer{}, &bytes.Buffer{}, inertProspectiveSelection)
	}()
	waitForProspectiveBarrier(t, ready, result)
	cancel()
	select {
	case got := <-result:
		if got.ActionExit == 0 {
			t.Fatalf("cancelled prospective result = %#v, want refusal", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancelled prospective producer did not return")
	}
	requireNoProspectiveBundles(t, tempRoot)
}

func TestProspectiveBuildRefusalLeavesNoBundle(t *testing.T) {
	root, _, tempRoot := prospectiveProducerFixture(t, "#!/bin/sh\nexit 0\n")
	writeGateFixtureFile(t, root, "go.mod", "module prospectivefixture\n\ngo 1.24\n", 0o644)
	writeGateFixtureFile(t, root, "scripts/go-build.sh", "#!/bin/sh\nexit 0\n", 0o755)
	writeGateFixtureFile(t, root, "scripts/go-build.inputs", "build_script=scripts/go-build.sh\n", 0o644)
	outcomeCommit(t, root, "declare prospective build")
	tree := outcomeGit(t, root, "write-tree")
	old := prospectiveRunBinary
	prospectiveRunBinary = runbinary.Factory{
		Build:  func(context.Context, string, string) error { return errors.New("build refused") },
		Verify: func(string, string) error { return nil },
	}
	t.Cleanup(func() { prospectiveRunBinary = old })

	var stderr bytes.Buffer
	result := executeTreeWithOwner(context.Background(), root, tree, &bytes.Buffer{}, &stderr, nil)
	if result.ActionExit == 0 || !strings.Contains(stderr.String(), "gate Bench executable unavailable") {
		t.Fatalf("build-refused prospective result = %#v, stderr=%q", result, stderr.String())
	}
	requireNoProspectiveBundles(t, tempRoot)
}

func TestProspectiveBundleCloseRetainsInheritedBaselineBinary(t *testing.T) {
	root, tree, tempRoot := prospectiveProducerFixture(t, "")
	baseline := filepath.Join(t.TempDir(), "baseline-bench")
	if err := os.WriteFile(baseline, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(runbinary.Env, baseline)
	old := prospectiveRunBinary
	prospectiveRunBinary = runbinary.Factory{Verify: func(string, string) error { return nil }}
	t.Cleanup(func() { prospectiveRunBinary = old })

	var stderr bytes.Buffer
	result := executeTreeWithOwner(context.Background(), root, tree, &bytes.Buffer{}, &stderr, nil)
	if result.ActionExit != 0 {
		t.Fatalf("inherited-binary prospective result = %#v, stderr=%q", result, stderr.String())
	}
	if data, err := os.ReadFile(baseline); err != nil || string(data) != "#!/bin/sh\nexit 0\n" {
		t.Fatalf("inherited baseline binary = %q, %v, want retained", data, err)
	}
	requireNoProspectiveBundles(t, tempRoot)
}

func prospectiveProducerFixture(t *testing.T, prospectiveScript string) (string, string, string) {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	tempRoot := t.TempDir()
	t.Setenv("TMPDIR", tempRoot)
	writeGateFixtureFile(t, root, ".bench/gate-inputs.json", "{\"schema\":1,\"closure\":\"local\",\"environment\":[\"HOME\"],\"paths\":[],\"tools\":[]}\n", 0o644)
	writeGateFixtureFile(t, root, ".bench/gate.sh", "#!/bin/sh\nset -eu\nbench=${BENCH_RUN_BINARY:?}\nif false; then \"$bench\" gate-phases; fi\nexit 0\n", 0o755)
	if prospectiveScript != "" {
		writeGateFixtureFile(t, root, prospectiveGatePath, prospectiveScript, 0o755)
	}
	writeGateFixtureFile(t, root, "tracked.txt", "prospective fixture\n", 0o644)
	outcomeCommit(t, root, "prospective producer fixture")
	return root, outcomeGit(t, root, "write-tree"), tempRoot
}

func inertProspectiveSelection(_ context.Context, source string) (*runbinary.Selection, error) {
	return &runbinary.Selection{Path: "/bin/true", SourceRoot: source}, nil
}

func waitForProspectiveBarrier(t *testing.T, path string, result <-chan Result) {
	t.Helper()
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case got := <-result:
			t.Fatalf("prospective producer returned before cancellation barrier: %#v", got)
		case <-deadline.C:
			t.Fatal("prospective producer did not reach cancellation barrier")
		case <-ticker.C:
			if data, err := os.ReadFile(path); err == nil && string(data) == "ready" {
				return
			}
		}
	}
}

func requireNoProspectiveBundles(t *testing.T, tempRoot string) {
	t.Helper()
	entries, err := os.ReadDir(tempRoot)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "bench-prospective-artifact-") {
			t.Fatalf("prospective bundle %q survived terminal outcome", entry.Name())
		}
	}
}

// TestEvidenceInspectionPublishesTheOwnerRecord is PAR29. Evidence inspection creates a
// private checkout of its own, so it needs the same bundle owner the full gate uses. A
// producer still on the defer-only helper publishes no record at all.
func TestEvidenceInspectionPublishesTheOwnerRecord(t *testing.T) {
	root, tree, tempRoot := prospectiveProducerFixture(t, "#!/bin/sh\nexit 0\n")
	snapshot := observeProspectiveOwnerRecord(t, root)

	if inspection := InspectTree(root, tree); inspection.Tree != tree {
		t.Fatalf("inspection = %#v, want the graded tree %q", inspection, tree)
	}
	requireProspectiveOwnerRecord(t, snapshot, root)
	requireNoProspectiveBundles(t, tempRoot)
}

// observeProspectiveOwnerRecord returns the path a copy of the next prospective owner
// record appears at. Git runs a post-checkout hook inside the private checkout while the
// bundle is still open, which is the one point a producer's caller reaches a record the
// producer removes before it returns. The hook copies nothing when the record is absent
// or is not a regular file, so the absent copy is itself the refusal.
func observeProspectiveOwnerRecord(t *testing.T, root string) string {
	t.Helper()
	common, err := benchgit.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	snapshot := filepath.Join(t.TempDir(), "published-owner-record")
	hook := "#!/bin/sh\nrecord=\"$PWD/../owner.json\"\nif [ -h \"$record\" ] || [ ! -f \"$record\" ]; then exit 0; fi\ncp -p \"$record\" \"" + snapshot + "\"\n"
	writeGateFixtureFile(t, common, filepath.Join("hooks", "post-checkout"), hook, 0o755)
	return snapshot
}

// requireProspectiveOwnerRecord grades one observed record against the published shape:
// a private regular file carrying schema 1, this process, and this repository's
// canonical common directory. Every producer row asks for exactly this shape.
func requireProspectiveOwnerRecord(t *testing.T, snapshot, root string) {
	t.Helper()
	info, err := os.Lstat(snapshot)
	if err != nil {
		t.Fatalf("published owner record = %v, want one observed record", err)
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		t.Fatalf("published owner record mode = %v, want regular 0600", info.Mode())
	}
	data, err := os.ReadFile(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	var record struct {
		Schema    int    `json:"schema"`
		OwnerPID  int    `json:"owner_pid"`
		CommonDir string `json:"common_dir"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&record); err != nil {
		t.Fatalf("decode published owner record: %v", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatalf("published owner record trailing data = %v, want EOF", err)
	}
	common, err := benchgit.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(common)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Clean(resolved); record.Schema != 1 || record.OwnerPID != os.Getpid() || record.CommonDir != want {
		t.Fatalf("published owner record = %#v, want schema 1, pid %d, common directory %q", record, os.Getpid(), want)
	}
}
