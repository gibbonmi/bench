//go:build system

package systemtest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
)

// candidateSpec is the spec that the candidate build binds. Its default-branch copy declares
// candidateField in the ledger and in the review record, and the candidate tree adds the
// field to both record types, so the main executable knows neither.
const (
	candidateSpec  = "specs/candidate/spec.md"
	candidateField = "fixture_install"
	candidateLabel = "candidate"
)

// candidateInstall is a copy of the kit whose main executable is built at its primary
// `dist/bench`, and one bound build worktree whose tree adds candidateField and whose own
// `dist/bench` is built. Each command runs through the copy's real wrapper, under a private
// Bench home and a PATH whose first directory takes the doctor repair's shim.
type candidateInstall struct {
	primary, home, shims, worktree string
}

// newCandidateInstall builds the fixture with the main executable, which binds the build as
// a start of the declared outcome records it.
func newCandidateInstall(t *testing.T) candidateInstall {
	t.Helper()
	declaration := "Persisted fields: intent-ledger:commitment." + candidateField + ", review-record:" + candidateField + "\n"
	primary := copyKitRepository(t, "candidate-install [primary]-")
	if err := os.MkdirAll(filepath.Join(primary, filepath.Dir(candidateSpec)), 0o755); err != nil {
		t.Fatal(err)
	}
	systemCommit(t, primary, candidateSpec, "# Candidate\n\nStatus: staged\n"+declaration, "stage the declared spec")
	c := candidateInstall{primary: primary}
	for _, dir := range []*string{&c.home, &c.shims} {
		created, err := os.MkdirTemp(owner.root, "candidate-install [private]-")
		if err != nil {
			t.Fatal(err)
		}
		*dir = created
	}
	if built := owner.runAt(primary, c.env(), "bash", filepath.Join(primary, "scripts", "go-build.sh"), primary, freshness.PublishedExecutable(primary)); built.code != 0 {
		t.Fatalf("main build = (%d, %q)", built.code, built.stderr)
	}
	c.require(t, primary, "worktree", "create", "--request", candidateLabel, "--label", candidateLabel)
	assignment, found, err := intent.FindAssignmentForRequest(primary, candidateLabel)
	if err != nil || !found {
		t.Fatalf("candidate assignment = %t, %v", found, err)
	}
	c.worktree = assignment.Worktree
	identity, err := commitrepo.SourceIdentity(primary, "main", candidateSpec)
	if err != nil {
		t.Fatal(err)
	}
	commitmenttest.EditCommitment(t, primary, func(state *intent.CommitmentState) {
		state.Bindings = append(state.Bindings, intent.DeliveryBinding{Assignment: assignment.ID, Request: assignment.Request, Milestone: "fixture", Outcome: candidateLabel, Deliverable: candidateSpec, Identity: identity})
		state.Claims = append(state.Claims, intent.OutcomeClaim{Milestone: "fixture", Outcome: candidateLabel})
	})
	for _, edit := range []struct{ file, anchor, field string }{
		{"internal/intent/ledger/commitment.go", "type CommitmentState struct {\n", "\tFixtureInstall string `json:\"" + candidateField + ",omitempty\"`\n"},
		{"internal/reviewrecord/record.go", "type Record struct {\n", "\tFixtureInstall string `json:\"" + candidateField + ",omitempty\"`\n"},
	} {
		source, err := os.ReadFile(filepath.Join(c.worktree, edit.file))
		if err != nil || !strings.Contains(string(source), edit.anchor) {
			t.Fatalf("candidate source %s = %v, want the anchor %q", edit.file, err, edit.anchor)
		}
		systemCommit(t, c.worktree, edit.file, strings.Replace(string(source), edit.anchor, edit.anchor+edit.field, 1), "add "+candidateField)
	}
	c.require(t, primary, "worktree", "build", candidateLabel)
	return c
}

// env strips every inherited routing override, so the wrapper selects the executable that
// its own resolution names from the working directory.
func (c candidateInstall) env() []string {
	return []string{"BENCH_HOME=" + c.home, "BENCH_KIT", "BENCH_RUN_BINARY", "BENCH_WRAPPER", "PATH=" + c.shims + string(os.PathListSeparator) + os.Getenv("PATH")}
}

// manifest is the primary broker manifest, beside the primary wrapper.
func (c candidateInstall) manifest() string {
	return filepath.Join(c.primary, "bin", brokermanifest.Name)
}

// run launches the primary wrapper in dir.
func (c candidateInstall) run(t *testing.T, dir string, args ...string) processResult {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runAt(dir, c.env(), "bash", append([]string{filepath.Join(c.primary, "bin", "bench.sh")}, args...)...)
}

// require runs the primary wrapper in dir and fails t unless the command exits 0.
func (c candidateInstall) require(t *testing.T, dir string, args ...string) processResult {
	t.Helper()
	result := c.run(t, dir, args...)
	if result.code != 0 {
		t.Fatalf("bench %q in %s = (%d, %q, %q)", args, dir, result.code, result.stdout, result.stderr)
	}
	return result
}

// plantLedgerField writes candidateField into the commitment state of the shared ledger, as
// the candidate's own writer would.
func (c candidateInstall) plantLedgerField(t *testing.T) {
	t.Helper()
	path := commitmenttest.LedgerPath(t, c.primary)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var ledger map[string]any
	if err := json.Unmarshal(data, &ledger); err != nil {
		t.Fatal(err)
	}
	commitment, ok := ledger["commitment"].(map[string]any)
	if !ok {
		t.Fatalf("ledger %s holds no commitment object", data)
	}
	commitment[candidateField] = "planted"
	if data, err = json.Marshal(ledger); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestCandidateInstallPublishesTheCandidate is CB15, CB21, and CB23. An install that the
// wrapper runs under the worktree's own executable refuses and writes nothing. From the
// primary checkout, the install publishes the candidate, the doctor repair binds the
// primary executable in the manifest, and `bench status` then reads a ledger that holds the
// field that the main executable refuses.
func TestCandidateInstallPublishesTheCandidate(t *testing.T) {
	c := newCandidateInstall(t)
	primary, built := freshness.PublishedExecutable(c.primary), freshness.PublishedExecutable(c.worktree)
	mainDigest := fileDigest(t, primary)
	mainManifest, err := os.ReadFile(c.manifest())
	if err != nil {
		t.Fatal(err)
	}

	refused := c.run(t, c.worktree, "worktree", "build", "--install", candidateLabel)
	if refused.code != 1 || !strings.Contains(refused.stderr, "the running executable "+built+" is not the primary executable") {
		t.Fatalf("install under the worktree executable = (%d, %q, %q), want the refusal naming %s", refused.code, refused.stdout, refused.stderr, built)
	}
	if after, err := os.ReadFile(c.manifest()); err != nil || string(after) != string(mainManifest) || fileDigest(t, primary) != mainDigest {
		t.Fatalf("the refused install changed the primary set: manifest %q (%v)", after, err)
	}

	// The status row names a ledger that the running executable refuses, so the main
	// executable's refusal of the candidate field is visible before the install.
	const refusedLedger = "intent ledger unavailable"
	if status := c.require(t, c.primary, "status"); strings.Contains(status.stdout, refusedLedger) {
		t.Fatalf("main status before the plant = %q, want a ledger that reads", status.stdout)
	}
	c.plantLedgerField(t)
	if status := c.run(t, c.primary, "status"); !strings.Contains(status.stdout, refusedLedger) {
		t.Fatalf("main status over the candidate field = (%d, %q), want the main executable to refuse the ledger", status.code, status.stdout)
	}
	c.require(t, c.primary, "worktree", "build", "--install", candidateLabel)
	fields, err := brokermanifest.Read(c.manifest())
	if err != nil {
		t.Fatal(err)
	}
	digest := fileDigest(t, primary)
	if fields["path"] != primary || fields["sha256"] != digest || digest != fileDigest(t, built) || digest == mainDigest {
		t.Fatalf("manifest after the install = %v; want path %s and the candidate digest %s", fields, primary, digest)
	}
	if status := c.require(t, c.primary, "status"); strings.Contains(status.stdout, refusedLedger) {
		t.Fatalf("status under the installed candidate = %q, want a ledger that reads", status.stdout)
	}
	if ledger, err := os.ReadFile(commitmenttest.LedgerPath(t, c.primary)); err != nil || !strings.Contains(string(ledger), `"`+candidateField+`"`) {
		t.Fatalf("ledger after the status = %s (%v), want it to keep %s", ledger, err, candidateField)
	}
}
