//go:build system

package systemtest

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
)

func TestDetectedProjectGateRejectsIgnoredDeclaredInput(t *testing.T) {
	repo := t.TempDir()
	home := filepath.Join(t.TempDir(), "bench-home")
	environment := []string{"BENCH_HOME=" + home, "BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}
	launch := func(program string, args ...string) processResult {
		if err := owner.observeSelected(); err != nil {
			t.Fatal(err)
		}
		return owner.runAt(repo, environment, program, args...)
	}
	if initialized := launch("git", "init", "-q"); initialized.code != 0 {
		t.Fatalf("git init = (%d, %q, %q)", initialized.code, initialized.stdout, initialized.stderr)
	}
	if err := os.WriteFile(filepath.Join(repo, "go.mod"), []byte("module example.com/detected\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	setup := launch(owner.selected.path, "setup", "--yes")
	if setup.code != 0 || !strings.Contains(setup.stdout, "go.mod detected") {
		t.Fatalf("detected-project setup = (%d, %q, %q)", setup.code, setup.stdout, setup.stderr)
	}
	if err := os.WriteFile(filepath.Join(repo, ".gitignore"), []byte("ignored file\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(repo, ".bench", "gate-inputs.json")
	declaration, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	declared := strings.Replace(string(declaration), `"paths": []`, `"paths": ["ignored file"]`, 1)
	if declared == string(declaration) {
		t.Fatalf("seeded gate input manifest has no empty paths list:\n%s", declaration)
	}
	if err := os.WriteFile(manifest, []byte(declared), 0o644); err != nil {
		t.Fatal(err)
	}

	wrapper := filepath.Join(repo, ".bench", "bin", "bench.sh")
	gate := launch("bash", wrapper, "gate", "--fresh")
	if gate.code != 1 || !strings.Contains(gate.stdout+gate.stderr, "gate input path ignored file is gitignored") {
		t.Fatalf("detected-project gate = (%d, %q, %q)", gate.code, gate.stdout, gate.stderr)
	}
	assertPrivateHomeEmpty(t, home)
}

// TestAdoptionSmokeJourney adopts one disposable repository with `bench setup --yes` and
// drives its scaffolded gate through the installed wrapper. The kit's own oracle observes
// the adopter's side of adoption, not an audit. Every launch binds one private BENCH_HOME
// under the test's temporary directory. runSelected and runWrapper carry fixed
// environments with no BENCH_HOME override. This journey composes observeSelected plus
// runAt the way TestWorktreeReauthorizeJourney does. It asserts after every leg that the
// private home stayed empty.
func TestAdoptionSmokeJourney(t *testing.T) {
	repo := owner.repos[1]
	home := filepath.Join(t.TempDir(), "bench-home")
	environment := []string{"BENCH_HOME=" + home, "BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}
	launch := func(program string, args ...string) processResult {
		if err := owner.observeSelected(); err != nil {
			t.Fatal(err)
		}
		return owner.runAt(repo, environment, program, args...)
	}

	setup := launch(owner.selected.path, "setup", "--yes")
	sentinelRow := ".bench/gate.sh is still the unconfigured fail-closed stub (replace the " + adopt.SentinelMarker + " sentinel with real checks)"
	if setup.code != 3 || !strings.Contains(setup.stdout+setup.stderr, sentinelRow) {
		t.Fatalf("bench setup --yes = (%d, %q, %q)", setup.code, setup.stdout, setup.stderr)
	}
	for _, rel := range []string{".bench/gate.sh", ".bench/gate-inputs.json", ".bench/bin/bench.sh", ".bench/dist/bench"} {
		if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("bench setup left %s unwritten: %v", rel, err)
		}
	}
	assertPrivateHomeEmpty(t, home)

	wrapper := filepath.Join(repo, ".bench", "bin", "bench.sh")
	gate := func(args ...string) processResult {
		return launch("bash", append([]string{wrapper}, args...)...)
	}

	stub := gate("gate")
	if stub.code != 1 || !strings.Contains(stub.stderr, "configure .bench/gate.sh - replace this sentinel with real checks") {
		t.Fatalf("untouched stub gate = (%d, %q, %q)", stub.code, stub.stdout, stub.stderr)
	}
	assertPrivateHomeEmpty(t, home)

	gateScript := filepath.Join(repo, ".bench", "gate.sh")
	retireSentinelLine(t, gateScript)
	canaryDir := filepath.Join(repo, "tests", "canary")
	if _, err := os.Stat(canaryDir); !os.IsNotExist(err) {
		t.Fatalf("tests/canary exists before the fixture leg: %v", err)
	}
	retired := gate("gate", "--fresh")
	if retired.code != 0 || !hasExactLine(retired.stdout, "gate: green") {
		t.Fatalf("gate with the sentinel retired = (%d, %q, %q)", retired.code, retired.stdout, retired.stderr)
	}
	assertPrivateHomeEmpty(t, home)

	fixture := filepath.Join(canaryDir, "adoption-smoke", "seeded-fixture")
	if err := os.MkdirAll(fixture, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixture, "input.txt"), []byte("seeded\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	withFixture := gate("gate", "--fresh")
	if withFixture.code != 0 || !strings.Contains(withFixture.stdout, "canary inventory ok (1 fixture bindings)") || !hasExactLine(withFixture.stdout, "gate: green") {
		t.Fatalf("gate with one project fixture = (%d, %q, %q)", withFixture.code, withFixture.stdout, withFixture.stderr)
	}
	assertPrivateHomeEmpty(t, home)

	manifest := filepath.Join(repo, ".bench", "gate-inputs.json")
	declaration, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(manifest); err != nil {
		t.Fatal(err)
	}
	// The test removes the manifest so HOME cannot reach the gate's environment. The
	// wrapper's own pool-home refusal fires in response. The conformance wrapper-pool-home
	// test pins the exact wording once. This leg asserts only that an undeclared input
	// reds the gate on HOME.
	unbound := gate("gate", "--fresh")
	if unbound.code != 1 || !strings.Contains(unbound.stderr, "HOME:") {
		t.Fatalf("gate without the seeded manifest = (%d, %q, %q)", unbound.code, unbound.stdout, unbound.stderr)
	}
	assertPrivateHomeEmpty(t, home)

	if err := os.WriteFile(manifest, declaration, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(canaryDir, "adoption-smoke")); err != nil {
		t.Fatal(err)
	}
	empty := gate("gate", "--fresh")
	if empty.code != 1 || !strings.Contains(boundedStderr(t, empty), "canary fixture inventory is empty") {
		t.Fatalf("gate with an empty tests/canary = (%d, %q, %q)", empty.code, empty.stdout, empty.stderr)
	}
	assertPrivateHomeEmpty(t, home)
}

// TestCommitmentBootstrapInstall links a project before any commitment policy exists. The
// installation needs no policy, and absent adoption refuses new delivery at its start and
// at the installed broker. A policy with no approval cannot publish itself; the planned and
// approved policy publishes and then admits only its approved deliverable.
func TestCommitmentBootstrapInstall(t *testing.T) {
	project := newLinkedProject(t)
	if _, err := os.Stat(filepath.Join(project.root, filepath.FromSlash(commitment.PolicyPath))); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("setup wrote a commitment policy: %v", err)
	}
	show := project.run(t, project.root, "commitment", "show")
	if show.code != 0 || !strings.Contains(show.stdout, commitment.OutlookAdoptionRequired) {
		t.Fatalf("installed commitment show = (%d, %q, %q)", show.code, show.stdout, show.stderr)
	}
	base := systemGitOutput(t, project.root, "rev-parse", "HEAD")
	delivery := systemCreateLandingWorktree(t, project.root, project.home, "delivery", "delivery")
	planning := systemCreateLandingWorktree(t, project.root, project.home, "planning", "planning")

	early := project.run(t, delivery.path, "commitment", "start", "--outcome", commitmenttest.DeliveryOutcome, "--request", delivery.request, "--deliverable", linkedDeliverable)
	if early.code != 1 || !strings.Contains(early.stdout, "commitment adoption required") {
		t.Fatalf("start before adoption = (%d, %q, %q)", early.code, early.stdout, early.stderr)
	}
	commitmenttest.Write(t, delivery.path, "delivery.txt", "delivered\n")
	commitSource(t, &delivery, "deliver before adoption")
	project.refuseLanding(t, delivery, base, "commitment adoption required")

	project.propose(t, &planning)
	project.refuseLanding(t, planning, base, "candidate policy has no exact approval")
	project.approve(t, &planning)
	project.publish(t, planning, base)
	// The linked project has no board, so the adoption publishes the policy alone.
	if got := systemGitOutput(t, project.root, "ls-tree", "--name-only", "main", "--", commitment.PolicyPath, "ROADMAP.md"); got != commitment.PolicyPath {
		t.Fatalf("published adoption tree = %q, want the policy alone", got)
	}
	if _, err := os.Lstat(filepath.Join(project.root, "ROADMAP.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("adoption wrote ROADMAP.md in the project root: %v", err)
	}

	unapproved := project.run(t, delivery.path, "commitment", "start", "--outcome", commitmenttest.DeliveryOutcome, "--request", delivery.request, "--deliverable", "specs/other/spec.md")
	if unapproved.code != 1 || !strings.Contains(unapproved.stdout, `deliverable "specs/other/spec.md" is not approved`) {
		t.Fatalf("start of an unapproved deliverable = (%d, %q, %q)", unapproved.code, unapproved.stdout, unapproved.stderr)
	}
	approved := project.run(t, delivery.path, "commitment", "start", "--outcome", commitmenttest.DeliveryOutcome, "--request", delivery.request, "--deliverable", linkedDeliverable)
	if approved.code != 0 {
		t.Fatalf("start of the approved deliverable = (%d, %q, %q)", approved.code, approved.stdout, approved.stderr)
	}
}

// TestCommitmentLinkedAdoption drives an adopted linked project through its installed
// wrapper. A start of an outcome that the commitment does not hold refuses, and a
// production commit with no delivery binding refuses until the committed outcome starts.
func TestCommitmentLinkedAdoption(t *testing.T) {
	project := newLinkedProject(t)
	delivery := systemCreateLandingWorktree(t, project.root, project.home, "delivery", "delivery")
	project.adopt(t)
	// The run predates adoption, so it takes the published policy before it commits.
	systemGit(t, delivery.path, "merge", "--no-edit", "main")

	uncommitted := project.run(t, delivery.path, "commitment", "start", "--outcome", "uncommitted", "--request", delivery.request, "--deliverable", linkedDeliverable)
	if uncommitted.code != 1 || !strings.Contains(uncommitted.stdout, `outcome "uncommitted" is not committed`) {
		t.Fatalf("start of an uncommitted outcome = (%d, %q, %q)", uncommitted.code, uncommitted.stdout, uncommitted.stderr)
	}
	commitmenttest.Write(t, delivery.path, "delivery.txt", "delivered\n")
	unbound := project.run(t, delivery.path, "commit", "-m", "deliver", "--", "delivery.txt")
	if unbound.code == 0 || !strings.Contains(unbound.stdout+unbound.stderr, "assignment has no current delivery binding") {
		t.Fatalf("production commit with no delivery binding = (%d, %q, %q)", unbound.code, unbound.stdout, unbound.stderr)
	}
	started := project.run(t, delivery.path, "commitment", "start", "--outcome", commitmenttest.DeliveryOutcome, "--request", delivery.request, "--deliverable", linkedDeliverable)
	if started.code != 0 {
		t.Fatalf("start of the committed outcome = (%d, %q, %q)", started.code, started.stdout, started.stderr)
	}
	bound := project.run(t, delivery.path, "commit", "-m", "deliver", "--", "delivery.txt")
	if bound.code != 0 {
		t.Fatalf("production commit with its delivery binding = (%d, %q, %q)", bound.code, bound.stdout, bound.stderr)
	}
}

// TestCommitmentInstalledAuthority publishes a candidate that bypasses every Bench guard
// in its own checkout. The installed broker refuses its displaced production work before
// publication, and the same tip lands only after the candidate's admission call.
func TestCommitmentInstalledAuthority(t *testing.T) {
	project := newLinkedProject(t)
	base := systemGitOutput(t, project.root, "rev-parse", "HEAD")
	candidate := systemCreateLandingWorktree(t, project.root, project.home, "candidate", "candidate")
	project.adopt(t)
	commitmenttest.Write(t, candidate.path, "displaced.txt", "displaced\n")
	commitSource(t, &candidate, "displace the committed outcome")
	project.refuseLanding(t, candidate, base, "assignment has no current delivery binding")

	admitted := project.run(t, candidate.path, "commitment", "start", "--outcome", commitmenttest.DeliveryOutcome, "--request", candidate.request, "--deliverable", linkedDeliverable)
	if admitted.code != 0 {
		t.Fatalf("candidate admission = (%d, %q, %q)", admitted.code, admitted.stdout, admitted.stderr)
	}
	project.publish(t, candidate, base)
	if got := systemGitOutput(t, project.root, "show", "main:displaced.txt"); got != "displaced" {
		t.Fatalf("published candidate content = %q", got)
	}
}

// retireSentinelLine performs the one documented operator step. It removes exactly the
// line that carries the sentinel marker and leaves the rest of the scaffolded gate intact.
func retireSentinelLine(t *testing.T, path string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")
	kept := make([]string, 0, len(lines))
	removed := 0
	for _, line := range lines {
		if strings.Contains(line, adopt.SentinelMarker) {
			removed++
			continue
		}
		kept = append(kept, line)
	}
	if removed != 1 {
		t.Fatalf("sentinel lines in %s = %d, want 1", path, removed)
	}
	if err := os.WriteFile(path, []byte(strings.Join(kept, "\n")), 0o755); err != nil {
		t.Fatal(err)
	}
}

// hasExactLine matches a whole output line, never a substring. A reused verdict prints
// "gate: green (fresh verdict reused for this tree)". A substring check on "gate: green"
// would wrongly accept that line as a real green run.
func hasExactLine(output, want string) bool {
	for _, line := range strings.Split(output, "\n") {
		if line == want {
			return true
		}
	}
	return false
}

// boundedStderr answers the stderr of a gate call. Each nested Bench call prints its own
// identity row, so a gate response can pass the bound and spill. The spill file then holds
// both streams in arrival order, and this answers that complete output.
func boundedStderr(t *testing.T, result processResult) string {
	t.Helper()
	if spill, ok := responseboundtest.Find(result.stdout); ok {
		return readSpillFile(t, spill.Path)
	}
	return result.stderr
}

func assertPrivateHomeEmpty(t *testing.T, home string) {
	t.Helper()
	entries, err := os.ReadDir(home)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		t.Fatal(err)
	}
	// The seam record lives at <home>/otel/<key>/traces.jsonl by the FT274 spec's
	// decision, so a gate run writes it even in a private home. A gate response that
	// spills writes its spill file under <home>/responses. Every other entry stays a leak.
	for _, entry := range entries {
		if entry.Name() == "otel" || entry.Name() == "responses" {
			continue
		}
		t.Fatalf("private BENCH_HOME %s is not empty: %v", home, entries)
	}
}
