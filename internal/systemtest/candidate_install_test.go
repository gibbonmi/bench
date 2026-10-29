//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/testrepo"
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

// gapChange is the unrelated Go file that a light-path landing publishes during the
// build's open chunk gap, and gapSlug names that landing's request and tickets-only folder.
const (
	gapChange = "internal/sanitize/gap_change.go"
	gapSlug   = "gap"
)

// newCandidateInstall builds the fixture with the main executable. It binds the candidate
// assignment to the declared outcome and claims that outcome, as a start records them. The
// policy approves the declared spec for that outcome, the spec's one ticket fences the
// candidate's sources and gapChange, and the gate is a plain green script, so a landing
// grades the broker and not the kit's checks.
func newCandidateInstall(t *testing.T) candidateInstall {
	t.Helper()
	declaration := "Persisted fields: intent-ledger:commitment." + candidateField + ", review-record:" + candidateField + "\n"
	primary := copyKitRepository(t, "candidate-install [primary]-")
	systemIdentity(t, primary)
	body := landableSpec("Candidate", declaration)
	recordtest.Prepare(t, primary, 1, candidateSpec, body).SetTicketWrites("source.txt, internal/intent/ledger/commitment.go, internal/reviewrecord/record.go, " + gapChange + " (new)")
	staged, err := os.ReadFile(filepath.Join(primary, candidateSpec))
	if err != nil {
		t.Fatal(err)
	}
	approved := commitment.Outcome{ID: candidateLabel, Criteria: []commitment.Criterion{{ID: "candidate.done", Text: "The candidate build lands."}}, Deliverables: []commitment.DeliveryBinding{{Source: commitment.SourceBinding{ID: "spec", Path: candidateSpec, Identity: commitment.Identity(staged)}}}}
	commitmenttest.WritePolicy(t, primary, commitment.Policy{Version: 1, ActiveMilestone: "fixture", Milestones: []commitment.Milestone{{ID: "fixture", Outcomes: []commitment.Outcome{approved}}}})
	testrepo.NewGateFixture(t.TempDir()).MustWrite(t, primary, "exit 0\n", "exit 0\n")
	commitmenttest.Commit(t, primary, "approve the declared spec under a green gate")
	systemGit(t, primary, "update-ref", "refs/bench/green/main", systemGitOutput(t, primary, "rev-parse", "HEAD"))
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

// plantRecordField writes candidateField into the build's review record, as the candidate's
// own writer would, and commits it. The record's JSON object opens on the document's first
// brace line.
func (c candidateInstall) plantRecordField(t *testing.T) {
	t.Helper()
	path, err := rr.RecordPath(candidateSpec)
	if err != nil {
		t.Fatal(err)
	}
	document, err := os.ReadFile(filepath.Join(c.worktree, path))
	if err != nil || !strings.Contains(string(document), "{\n") {
		t.Fatalf("review record %s = %v, want a JSON object", path, err)
	}
	systemCommit(t, c.worktree, path, strings.Replace(string(document), "{\n", "{\n  \""+candidateField+"\": \"planted\",\n", 1), "record "+candidateField)
	if _, err := rr.Read(c.worktree, candidateSpec); err == nil {
		t.Fatal("the strict record read accepts the planted field, so the row grades nothing")
	}
}

// landGapChange opens a chunk gap of the build, plants candidateField in the ledger, and
// lands gapChange from an unbound assignment through the installed broker. The assignment
// is created before the plant, because this test process cannot read the planted ledger.
func (c candidateInstall) landGapChange(t *testing.T) processResult {
	t.Helper()
	commitmenttest.OpenChunkGaps(t, c.primary)
	c.require(t, c.primary, "worktree", "create", "--request", gapSlug, "--label", gapSlug)
	gap, found, err := intent.FindAssignmentForRequest(c.primary, gapSlug)
	if err != nil || !found {
		t.Fatalf("gap assignment = %t, %v", found, err)
	}
	base := systemGitOutput(t, c.primary, "rev-parse", "main")
	commitmenttest.WriteLightTicket(t, gap.Worktree, spec.ClosedFolderPath(gapSlug)+"/tickets/one.md", gapChange+" (new)")
	commitmenttest.Write(t, gap.Worktree, gapChange, "package sanitize\n\n// gapChange marks an unrelated landing.\nconst gapChange = true\n")
	systemGit(t, gap.Worktree, "add", "-A")
	systemGit(t, gap.Worktree, "commit", "-qm", "an unrelated change")
	commitmenttest.PlantCommitmentField(t, c.primary, candidateField)
	return c.run(t, c.primary, "worktree", "land", "--request", gapSlug, "--base", base, "--source-tip", systemGitOutput(t, gap.Worktree, "rev-parse", "HEAD"), "--spec", gapSlug, "-m", "land an unrelated change", gap.Worktree)
}

// TestCandidateInstallSurvivesAGapLanding is CB50. An unrelated landing during the open
// chunk gap of the declared build holds the installed candidate, so the primary executable
// keeps its digest and still reads the ledger that holds the candidate's field.
func TestCandidateInstallSurvivesAGapLanding(t *testing.T) {
	t.Parallel()
	c := newCandidateInstall(t)
	c.require(t, c.primary, "worktree", "build", "--install", candidateLabel)
	primary := freshness.PublishedExecutable(c.primary)
	installed := fileDigest(t, primary)

	landed := c.landGapChange(t)
	if landed.code != 0 || !strings.Contains(landed.stdout, "  refresh,held\n") {
		t.Fatalf("gap landing = (%d, %q, %q), want a held refresh at exit 0", landed.code, landed.stdout, landed.stderr)
	}
	if got := fileDigest(t, primary); got != installed {
		t.Fatalf("primary digest after the gap landing = %s, want the installed candidate %s", got, installed)
	}
	if status := c.run(t, c.primary, "status"); status.code != 0 {
		t.Fatalf("status after the gap landing = (%d, %q, %q), want exit 0", status.code, status.stdout, status.stderr)
	}
}

// TestCandidateInstallLandsItsNewFields is CB24 and CB25. The build folds the gap landing's
// unrelated Go change and records the review-record field that the main broker refuses.
// The installed candidate lands it through the unchanged land route. The landed tree is not
// the installed candidate's source, so the refresh rebuilds the primary executable from the
// landed sources and binds it in the manifest.
func TestCandidateInstallLandsItsNewFields(t *testing.T) {
	t.Parallel()
	c := newCandidateInstall(t)
	start := systemGitOutput(t, c.worktree, "merge-base", "HEAD", "main")
	c.require(t, c.primary, "worktree", "build", "--install", candidateLabel)
	primary := freshness.PublishedExecutable(c.primary)
	installed := fileDigest(t, primary)
	if gap := c.landGapChange(t); gap.code != 0 {
		t.Fatalf("gap landing = (%d, %q, %q)", gap.code, gap.stdout, gap.stderr)
	}
	systemGit(t, c.worktree, "merge", "-q", "--no-edit", "main")
	recordtest.RetainSingleChunk(t, c.worktree, candidateSpec, start)
	c.plantRecordField(t)

	landed := c.run(t, c.primary, "worktree", "land", "--request", candidateLabel, "--base", start, "--source-tip", systemGitOutput(t, c.worktree, "rev-parse", "HEAD"), "--spec", spec.LiveSpecSlug(candidateSpec), "-m", "land the candidate build", c.worktree)
	if landed.code != 0 || !strings.Contains(landed.stdout, "  refresh,complete\n") {
		t.Fatalf("build landing = (%d, %q, %q), want a complete refresh at exit 0", landed.code, landed.stdout, landed.stderr)
	}
	path, err := rr.RecordPath(candidateSpec)
	if err != nil {
		t.Fatal(err)
	}
	if published := systemGitOutput(t, c.primary, "show", "main:"+path); !strings.Contains(published, `"`+candidateField+`"`) {
		t.Fatalf("published review record = %q, want %s", published, candidateField)
	}
	if err := freshness.Verify(c.primary, primary); err != nil {
		t.Fatalf("the primary seal after the landing does not verify against the primary sources: %v", err)
	}
	fields, err := brokermanifest.Read(c.manifest())
	if err != nil {
		t.Fatal(err)
	}
	if digest := fileDigest(t, primary); fields["sha256"] != digest || digest == installed {
		t.Fatalf("manifest after the landing = %v, primary digest %s; want the rebuilt digest, not the installed %s", fields, digest, installed)
	}
}

// TestCandidateInstallPublishesTheCandidate is CB15, CB21, and CB23. An install that the
// wrapper runs under the worktree's own executable refuses and writes nothing. From the
// primary checkout, the install publishes the candidate, the doctor repair binds the
// primary executable in the manifest, and `bench status` then reads a ledger that holds the
// field that the main executable refuses.
func TestCandidateInstallPublishesTheCandidate(t *testing.T) {
	t.Parallel()
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
	commitmenttest.PlantCommitmentField(t, c.primary, candidateField)
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
