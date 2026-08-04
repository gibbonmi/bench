package specbuild

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/jsonfile"
)

// The fixtures below drive the real assign and checkpoint paths, whose owner returns no
// branch name, so every record they produce carries the pre-Branch shape the enumeration
// has to reach by assignment identity alone.
const assignmentNamespaceLiteral = "refs/heads/bench/assign/"

func residueDispositions(refs []provisionalRef) map[string]provisionalDisposition {
	classified := map[string]provisionalDisposition{}
	for _, ref := range refs {
		classified[ref.Name] = ref.Disposition
	}
	return classified
}

func namespaceBranchFor(t *testing.T, root, id string) string {
	t.Helper()
	var found []string
	for _, name := range strings.Split(git(t, root, "for-each-ref", "--format=%(refname)", assignmentNamespaceLiteral), "\n") {
		if strings.HasSuffix(name, "/"+id) {
			found = append(found, name)
		}
	}
	if len(found) != 1 {
		t.Fatalf("assignment %s has %d namespace refs, want 1: %v", id, len(found), found)
	}
	return found[0]
}

func residueSubject(t *testing.T, fixture checkpointFixture) (record, assignment, string) {
	t.Helper()
	run := loadRun(t, fixture.service)
	_, assigned, ok := assignmentFor(run, fixture.assigned.ID)
	if !ok || assigned.Branch != "" || assigned.CheckpointRef == "" {
		t.Fatalf("residue subject assignment = %#v", assigned)
	}
	return run, assigned, namespaceBranchFor(t, fixture.root, assigned.ID)
}

func TestReclaimEnumerationClassifiesATerminalRunWithNoStoredBranchName(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	run, assigned, branch := residueSubject(t, fixture)
	run.Terminal = true

	refs, err := fixture.service.provisionalResidue(run)
	if err != nil {
		t.Fatalf("provisionalResidue: %v", err)
	}
	want := map[string]provisionalDisposition{
		branch:                 refReclaimable,
		assigned.CheckpointRef: refReclaimable,
		run.Candidate:          refReclaimable,
	}
	if got := residueDispositions(refs); !reflect.DeepEqual(got, want) {
		t.Fatalf("inventory = %#v, want %#v", got, want)
	}
	for _, ref := range refs {
		if object := git(t, fixture.root, "rev-parse", ref.Name); ref.Object != object {
			t.Fatalf("%s recorded object %q, want %q", ref.Name, ref.Object, object)
		}
	}
}

func TestReclaimEnumerationReportsEveryNonReclaimableClass(t *testing.T) {
	for _, test := range []struct {
		name     string
		terminal bool
		owned    provisionalDisposition
	}{
		{"terminal", true, refReclaimable},
		{"non-terminal", false, refActive},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := checkpointedReleaseFixture(t)
			run, assigned, branch := residueSubject(t, fixture)
			commit := git(t, fixture.root, "rev-parse", "HEAD")
			twin := intent.AssignmentBranchRef(strings.Repeat("a", 32), assigned.ID)
			orphan := intent.AssignmentBranchRef(strings.Repeat("b", 32), strings.Repeat("c", 32))
			for _, ref := range []string{twin, orphan} {
				git(t, fixture.root, "update-ref", ref, commit)
			}
			run.Terminal = test.terminal

			refs, err := fixture.service.provisionalResidue(run)
			if err != nil {
				t.Fatalf("provisionalResidue: %v", err)
			}
			want := map[string]provisionalDisposition{
				branch:                 refAmbiguous,
				twin:                   refAmbiguous,
				orphan:                 refUnclassified,
				assigned.CheckpointRef: test.owned,
				run.Candidate:          test.owned,
			}
			if got := residueDispositions(refs); !reflect.DeepEqual(got, want) {
				t.Fatalf("inventory = %#v, want %#v", got, want)
			}
		})
	}
}

// reclamationSubject is a terminal run whose namespace holds every class the enumeration
// can report at once: the assignment's own branch and a twin sharing its ID are ambiguous,
// a ref no record claims is unclassified, and the run's checkpoint and candidate refs are
// reclaimable. It returns the classification each ref must receive.
func reclamationSubject(t *testing.T) (checkpointFixture, map[string]provisionalDisposition) {
	t.Helper()
	fixture := checkpointedReleaseFixture(t)
	run, assigned, branch := residueSubject(t, fixture)
	twin := intent.AssignmentBranchRef(strings.Repeat("a", 32), assigned.ID)
	orphan := intent.AssignmentBranchRef(strings.Repeat("b", 32), strings.Repeat("c", 32))
	for _, ref := range []string{twin, orphan} {
		git(t, fixture.root, "update-ref", ref, git(t, fixture.root, "rev-parse", "HEAD"))
	}
	setRunTerminal(t, fixture.service, true)
	return fixture, map[string]provisionalDisposition{
		branch: refAmbiguous, twin: refAmbiguous, orphan: refUnclassified,
		assigned.CheckpointRef: refReclaimable, run.Candidate: refReclaimable,
	}
}

// setRunTerminal drives the durable record across the terminal transition the enumeration
// reads its verdicts from, touching no ref, so a plan taken either side of the call spans
// the same refs.
func setRunTerminal(t *testing.T, service *Service, terminal bool) {
	t.Helper()
	run := loadRun(t, service)
	run.Terminal = terminal
	saveRun(t, service, run)
}

func refSnapshot(t *testing.T, root string) string {
	t.Helper()
	return git(t, root, "for-each-ref", "--format=%(refname) %(objectname)")
}

func refIdentities(plan ReclamationPlan) []string {
	identities := make([]string, 0, len(plan.Refs))
	for _, ref := range plan.Refs {
		identities = append(identities, ref.Name+" "+ref.Object)
	}
	return identities
}

// driftRefs moves the namespace the plan was derived over, so a fingerprint planned before
// the call no longer describes the tree.
func driftRefs(t *testing.T, root string) {
	t.Helper()
	git(t, root, "update-ref", intent.AssignmentBranchRef(strings.Repeat("d", 32), strings.Repeat("e", 32)), git(t, root, "rev-parse", "HEAD"))
}

func TestReclaimPlansWithoutMutatingAnyRef(t *testing.T) {
	fixture, _ := reclamationSubject(t)
	before := refSnapshot(t, fixture.root)

	plan, err := fixture.service.Reclaim("build demo")
	if err != nil {
		t.Fatalf("Reclaim: %v", err)
	}
	if after := refSnapshot(t, fixture.root); after != before {
		t.Fatalf("plan mutated refs:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if !canonicalDigest(plan.Fingerprint) {
		t.Fatalf("fingerprint = %q", plan.Fingerprint)
	}
	repeat, err := fixture.service.Reclaim("build demo")
	if err != nil || repeat.Fingerprint != plan.Fingerprint {
		t.Fatalf("re-plan = %q, %v, want %q", repeat.Fingerprint, err, plan.Fingerprint)
	}
	driftRefs(t, fixture.root)
	moved, err := fixture.service.Reclaim("build demo")
	if err != nil || moved.Fingerprint == plan.Fingerprint {
		t.Fatalf("fingerprint survived a changed ref: %q, %v", moved.Fingerprint, err)
	}
}

func TestReclaimAppliesToTheReclaimableClassAlone(t *testing.T) {
	fixture, classified := reclamationSubject(t)
	plan, err := fixture.service.Reclaim("build demo")
	if err != nil {
		t.Fatalf("Reclaim: %v", err)
	}

	applied, err := fixture.service.ApplyReclaim("build demo", plan.Fingerprint)
	if err != nil || !applied.Applied || applied.Fingerprint != plan.Fingerprint {
		t.Fatalf("ApplyReclaim = %#v, %v", applied, err)
	}
	for name, disposition := range classified {
		present := git(t, fixture.root, "for-each-ref", "--format=%(refname)", name) != ""
		if want := disposition != refReclaimable; present != want {
			t.Fatalf("%s classified %s present = %v, want %v", name, disposition, present, want)
		}
	}
}

// A plan-time refusal — an empty, malformed, or drifted fingerprint — never reaches the
// deletion loop, so it has no spent set to report: the refusal text is unchanged and the
// returned plan is the zero value rather than a receipt (SS3).
func TestReclaimRefusesEveryFingerprintItDidNotJustPlan(t *testing.T) {
	for _, test := range []struct {
		name    string
		offered func(planned string) string
	}{
		{"empty", func(string) string { return "" }},
		{"malformed", func(string) string { return "not-a-digest" }},
		{"drifted", func(planned string) string { return planned }},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture, _ := reclamationSubject(t)
			plan, err := fixture.service.Reclaim("build demo")
			if err != nil {
				t.Fatalf("Reclaim: %v", err)
			}
			driftRefs(t, fixture.root)
			before := refSnapshot(t, fixture.root)

			applied, err := fixture.service.ApplyReclaim("build demo", test.offered(plan.Fingerprint))
			if err == nil {
				t.Fatalf("ApplyReclaim accepted an %s fingerprint: %#v", test.name, applied)
			}
			if !strings.Contains(err.Error(), "fresh plan") {
				t.Fatalf("%s refusal = %q, want the fresh-plan action named", test.name, err)
			}
			if !reflect.DeepEqual(applied, ReclamationPlan{}) {
				t.Fatalf("%s refusal returned %#v, want no receipt", test.name, applied)
			}
			if after := refSnapshot(t, fixture.root); after != before {
				t.Fatalf("%s refusal mutated refs:\nbefore:\n%s\nafter:\n%s", test.name, before, after)
			}
		})
	}
}

// The disposition is the whole authorization, so it has to sit inside the fingerprint: a
// plan taken while a run is live shows every ref as active and promises no deletion, and
// the same refs become reclaimable the moment the run goes terminal. Nothing about the refs
// moves across that transition, so a fingerprint over names and objects alone still matches
// and would authorize deleting refs the operator was shown as safe.
func TestReclaimRefusesAFingerprintPlannedUnderADifferentDisposition(t *testing.T) {
	fixture, _ := reclamationSubject(t)
	setRunTerminal(t, fixture.service, false)

	stale, err := fixture.service.Reclaim("build demo")
	if err != nil {
		t.Fatalf("Reclaim over a live run: %v", err)
	}
	for _, ref := range stale.Refs {
		if ref.Disposition == string(refReclaimable) {
			t.Fatalf("a live run planned %s as reclaimable", ref.Name)
		}
	}

	setRunTerminal(t, fixture.service, true)
	current, err := fixture.service.Reclaim("build demo")
	if err != nil {
		t.Fatalf("Reclaim after the terminal transition: %v", err)
	}
	if got, want := refIdentities(current), refIdentities(stale); !reflect.DeepEqual(got, want) {
		t.Fatalf("the transition moved ref identities, so the verdict is not the isolated change: %v, want %v", got, want)
	}
	if current.Fingerprint == stale.Fingerprint {
		t.Fatalf("fingerprint %q survived the reclassification of %v", stale.Fingerprint, refIdentities(stale))
	}
	before := refSnapshot(t, fixture.root)

	applied, err := fixture.service.ApplyReclaim("build demo", stale.Fingerprint)
	if err == nil {
		t.Fatalf("ApplyReclaim accepted a fingerprint planned under a different disposition: %#v", applied)
	}
	if !strings.Contains(err.Error(), "fresh plan") {
		t.Fatalf("refusal = %q, want the fresh-plan action named", err)
	}
	if after := refSnapshot(t, fixture.root); after != before {
		t.Fatalf("refusal mutated refs:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// restartedResidue is one slug that has run twice: a prior attempt the lifecycle drove to
// terminal and moved into history, and the live attempt that replaced it. It names the refs
// each attempt owns so a plan taken over the live record can be graded per owning record.
type restartedResidue struct {
	fixture                                       checkpointFixture
	live, prior                                   record
	priorBranch, priorCheckpoint, priorAssignment string
	liveBranch, liveAssignment                    string
}

// restartedResidueFixture drives the real lifecycle across a restart — checkpoint, abandon,
// advance the tip, start again, assign again — so the history entry the enumeration reads is
// one the lifecycle itself wrote rather than a struct the test filled in. Abandon preserves
// the prior assignment's payload into a recovery ref and deletes its branch, so the branch is
// restored at the exact identity the lifecycle gave it: a branch stranded beside a terminal
// record is the residue this pass exists to reach, and it is what a run that ended before
// reclamation existed leaves behind.
func restartedResidueFixture(t *testing.T) restartedResidue {
	t.Helper()
	fixture := checkpointedReleaseFixture(t)
	prior, priorAssigned, priorBranch := residueSubject(t, fixture)
	priorObject := git(t, fixture.root, "rev-parse", priorBranch)
	fixture.service.worktrees = &abandonOwner{}
	abandonRun(t, fixture.service)
	advanceWorking(t, fixture.root)
	if _, err := fixture.service.Start(t.Context(), "build demo"); err != nil {
		t.Fatalf("restart: %v", err)
	}
	git(t, fixture.root, "update-ref", priorBranch, priorObject)
	fixture.service.worktrees = realOwner{}
	live, _, err := fixture.service.Assign(t.Context(), "build demo", "one.md", "restarted attempt")
	if err != nil {
		t.Fatalf("Assign into the restarted run: %v", err)
	}
	run := loadRun(t, fixture.service)
	if len(run.History) != 1 || run.Run == prior.Run || run.Terminal {
		t.Fatalf("restarted record = %#v, want one retained prior run and a live current run", run)
	}
	return restartedResidue{
		fixture: fixture, live: run, prior: prior,
		priorBranch: priorBranch, priorCheckpoint: priorAssigned.CheckpointRef, priorAssignment: priorAssigned.ID,
		liveBranch: namespaceBranchFor(t, fixture.root, live.ID), liveAssignment: live.ID,
	}
}

// The whole classification a plan over the restarted slug must produce. Both tests below
// grade against it so neither can pass by answering for its own half alone.
func (r restartedResidue) want(prior provisionalDisposition) map[string]provisionalDisposition {
	return map[string]provisionalDisposition{
		r.priorBranch: prior, r.priorCheckpoint: prior, r.prior.Candidate: prior,
		r.liveBranch: refActive, r.live.Candidate: refActive,
	}
}

// A restart retires the prior attempt into history and leaves its refs in the tree. The
// history entry is the record that owns them and it says terminal, so the enumeration has to
// answer for it: before this, the prior candidate and checkpoint were absent from the plan
// entirely and the prior branch was reported unclassified and retained forever.
func TestReclaimEnumerationClassifiesResidueOfAPriorTerminalRun(t *testing.T) {
	subject := restartedResidueFixture(t)

	refs, err := subject.fixture.service.provisionalResidue(subject.live)
	if err != nil {
		t.Fatalf("provisionalResidue: %v", err)
	}
	if got, want := residueDispositions(refs), subject.want(refReclaimable); !reflect.DeepEqual(got, want) {
		t.Fatalf("inventory = %#v, want %#v", got, want)
	}
	attributed := map[string]string{}
	for _, ref := range refs {
		attributed[ref.Name] = ref.Assignment
	}
	if attributed[subject.priorBranch] != subject.priorAssignment || attributed[subject.priorCheckpoint] != subject.priorAssignment {
		t.Fatalf("prior refs attributed to %q and %q, want the history record's assignment %q", attributed[subject.priorBranch], attributed[subject.priorCheckpoint], subject.priorAssignment)
	}
}

// Terminality is per record, so extending the enumeration must not spill the history
// entry's verdict onto the live attempt whose work is still in flight.
func TestReclaimEnumerationLeavesTheLiveRecordsOwnRefsActive(t *testing.T) {
	subject := restartedResidueFixture(t)

	refs, err := subject.fixture.service.provisionalResidue(subject.live)
	if err != nil {
		t.Fatalf("provisionalResidue: %v", err)
	}
	classified := residueDispositions(refs)
	for _, name := range []string{subject.liveBranch, subject.live.Candidate} {
		if classified[name] != refActive {
			t.Fatalf("live ref %s classified %q, want %q", name, classified[name], refActive)
		}
	}
	attributed := map[string]string{}
	for _, ref := range refs {
		attributed[ref.Name] = ref.Assignment
	}
	if attributed[subject.liveBranch] != subject.liveAssignment {
		t.Fatalf("live branch attributed to %q, want the live record's assignment %q", attributed[subject.liveBranch], subject.liveAssignment)
	}
}

// Sitting in history is not itself a licence to delete: the verdict is the entry's own
// terminal flag. The lifecycle never writes a live entry into history — the loader refuses
// one — so the entry is derived from the genuine one the restart wrote and handed straight
// to the enumeration, which is the only surface that can be asked about a record the durable
// state would reject.
func TestReclaimEnumerationRetainsAPriorRunWhoseRecordIsNotTerminal(t *testing.T) {
	subject := restartedResidueFixture(t)
	live := subject.live
	live.History = []json.RawMessage{reopenedHistoryEntry(t, live.History[0])}

	refs, err := subject.fixture.service.provisionalResidue(live)
	if err != nil {
		t.Fatalf("provisionalResidue: %v", err)
	}
	if got, want := residueDispositions(refs), subject.want(refActive); !reflect.DeepEqual(got, want) {
		t.Fatalf("inventory = %#v, want %#v", got, want)
	}
}

func reopenedHistoryEntry(t *testing.T, entry json.RawMessage) json.RawMessage {
	t.Helper()
	var prior record
	if err := jsonfile.DecodeDocument(entry, &prior); err != nil || !prior.Terminal {
		t.Fatalf("retained history entry = %#v, %v; want a terminal record to reopen", prior, err)
	}
	prior.Terminal = false
	raw, err := json.Marshal(prior)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// Git deletes refs one at a time, so an apply that meets a moved ref partway has already
// spent the deletions before it. Reporting nothing there tells the operator the tree is
// untouched while it is not, and the spent fingerprint then refuses every retry. The
// deleted set is carried once, by the returned receipt, rather than restated into the
// error text: the two would drift apart the moment they disagreed, so the refusal names
// only the ref that drifted and the fresh-plan hint, and leaves the deleted set to the
// receipt alone (SS1, SS2).
func TestReclaimApplyReportsTheDeletionsItSpentBeforeADriftedRef(t *testing.T) {
	fixture, classified := reclamationSubject(t)
	plan, err := fixture.service.Reclaim("build demo")
	if err != nil {
		t.Fatalf("Reclaim: %v", err)
	}
	spent, drifted := reclaimableOrder(t, plan)
	fixture.service.fault = driftAfterPlanning(t, fixture.root, drifted)

	receipt, err := fixture.service.ApplyReclaim("build demo", plan.Fingerprint)
	if err == nil {
		t.Fatalf("ApplyReclaim reported success over a drifted ref: %#v", receipt)
	}
	if !strings.Contains(err.Error(), drifted) {
		t.Fatalf("interrupted apply = %q, want %s named", err, drifted)
	}
	if strings.Contains(err.Error(), spent.Name) {
		t.Fatalf("interrupted apply = %q, restates the deleted set the receipt already carries", err)
	}
	if !strings.Contains(err.Error(), "request a fresh plan") {
		t.Fatalf("interrupted apply = %q, want the convergence hint kept", err)
	}
	if !receipt.Applied || !reflect.DeepEqual(receipt.Refs, []ReclamationRef{spent}) {
		t.Fatalf("interrupted receipt = %#v, want the spent deletion %#v", receipt, spent)
	}
	if receipt.Fingerprint == plan.Fingerprint {
		t.Fatalf("interrupted receipt carried the spent plan's fingerprint %q", plan.Fingerprint)
	}
	if want := digest(reclamationFacts(ReclamationPlan{Slug: receipt.Slug, Refs: receipt.Refs})); receipt.Fingerprint != want {
		t.Fatalf("interrupted receipt fingerprint = %q, want %q over the refs it reports", receipt.Fingerprint, want)
	}
	if present := git(t, fixture.root, "for-each-ref", "--format=%(refname)", spent.Name); present != "" {
		t.Fatalf("%s is reported deleted but still present", spent.Name)
	}
	if present := git(t, fixture.root, "for-each-ref", "--format=%(refname)", drifted); present == "" {
		t.Fatalf("%s drifted yet was deleted", drifted)
	}
	if _, ok := classified[drifted]; !ok {
		t.Fatalf("%s is not a ref this subject classifies", drifted)
	}
}

// The remainder has to be reachable, or an interrupted apply is a dead end: the spent
// fingerprint can never match again, so only a fresh plan can finish the work.
func TestReclaimReplanAfterAnInterruptedApplyFinishesTheRemainder(t *testing.T) {
	fixture, _ := reclamationSubject(t)
	plan, err := fixture.service.Reclaim("build demo")
	if err != nil {
		t.Fatalf("Reclaim: %v", err)
	}
	spent, drifted := reclaimableOrder(t, plan)
	fixture.service.fault = driftAfterPlanning(t, fixture.root, drifted)
	if _, err := fixture.service.ApplyReclaim("build demo", plan.Fingerprint); err == nil {
		t.Fatal("ApplyReclaim reported success over a drifted ref")
	}
	fixture.service.fault = nil

	fresh, err := fixture.service.Reclaim("build demo")
	if err != nil {
		t.Fatalf("re-plan: %v", err)
	}
	if fresh.Fingerprint == plan.Fingerprint {
		t.Fatalf("re-plan reused the spent fingerprint %q", fresh.Fingerprint)
	}
	if _, err := fixture.service.ApplyReclaim("build demo", plan.Fingerprint); err == nil {
		t.Fatal("the spent fingerprint still applied after the interruption")
	}
	applied, err := fixture.service.ApplyReclaim("build demo", fresh.Fingerprint)
	if err != nil || !applied.Applied {
		t.Fatalf("apply of the fresh plan = %#v, %v", applied, err)
	}
	for _, name := range []string{spent.Name, drifted} {
		if present := git(t, fixture.root, "for-each-ref", "--format=%(refname)", name); present != "" {
			t.Fatalf("%s survived the completed reclamation", name)
		}
	}
}

// reclaimableOrder names the first ref the apply will delete and the second one, which the
// deletion order makes the earliest ref whose drift can strand a spent deletion behind it.
func reclaimableOrder(t *testing.T, plan ReclamationPlan) (ReclamationRef, string) {
	t.Helper()
	var order []ReclamationRef
	for _, ref := range plan.Refs {
		if ref.Disposition == string(refReclaimable) {
			order = append(order, ref)
		}
	}
	if len(order) < 2 {
		t.Fatalf("plan holds %d reclaimable refs, want at least 2 to interrupt between", len(order))
	}
	return order[0], order[1].Name
}

// driftAfterPlanning stands in for the concurrent writer the fingerprint cannot exclude: it
// moves one planned ref after the plan is taken and lets the apply proceed, so the deletion
// that meets it fails its own compare-and-swap rather than an injected error.
func driftAfterPlanning(t *testing.T, root, ref string) func(string) error {
	t.Helper()
	return func(point string) error {
		if point == "reclaim/apply" {
			git(t, root, "update-ref", ref, git(t, root, "rev-parse", "HEAD"))
		}
		return nil
	}
}

func TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes(t *testing.T) {
	fixture, _ := reclamationSubject(t)

	plan, err := fixture.service.Reclaim("build demo")
	if err != nil {
		t.Fatalf("Reclaim: %v", err)
	}
	want := []ReclamationClass{
		{string(refReclaimable), 2}, {string(refActive), 0},
		{string(refUnclassified), 1}, {string(refAmbiguous), 2},
	}
	if got := plan.Classes(); !reflect.DeepEqual(got, want) {
		t.Fatalf("classes = %#v, want %#v", got, want)
	}
}

func TestReclaimOnPromotionDropsResidueOfARecordWithNoStoredBranchName(t *testing.T) {
	fixture := reviewedPromotionFixture(t)
	setAssignmentsReleased(t, fixture.service, true)
	run, assigned, branch := residueSubject(t, fixture)
	fixture.service.gate = &promotionGate{accept: true}

	status, err := fixture.service.Promote(t.Context(), "build demo")
	if err != nil || status.State != "terminal" {
		t.Fatalf("Promote = %#v, %v", status, err)
	}
	for _, ref := range []string{branch, assigned.CheckpointRef, run.Candidate} {
		if git(t, fixture.root, "for-each-ref", "--format=%(refname)", ref) != "" {
			t.Fatalf("promotion retained %s", ref)
		}
	}
	if left := git(t, fixture.root, "for-each-ref", "--format=%(refname)", assignmentNamespaceLiteral); left != "" {
		t.Fatalf("promotion left assignment namespace refs: %q", left)
	}
}

// tamperedCheckpointRef writes ref into the subject assignment's CheckpointRef through the
// service's own save, so a refusal on read-back is validCore's verdict on a value its own
// writer durably accepted — not a mangled state file the decoder would refuse anyway.
func tamperedCheckpointRef(t *testing.T, fixture checkpointFixture, ref string) record {
	t.Helper()
	run, assigned, _ := residueSubject(t, fixture)
	key, stored, ok := assignmentFor(run, assigned.ID)
	if !ok {
		t.Fatal("missing assignment")
	}
	stored.CheckpointRef = ref
	run.Assignments[key] = stored
	saveRun(t, fixture.service, run)
	return run
}

// A stored CheckpointRef is authorization to delete the ref it names, so the loader grades
// it with the record's other identity facts: a ref outside the checkpoint namespace is
// incomplete prior state, however real the ref it points at.
func TestReclaimTrustBoundaryRefusesACheckpointRefOutsideTheCheckpointNamespace(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	tamperedCheckpointRef(t, fixture, git(t, fixture.root, "symbolic-ref", "HEAD"))

	if _, found, err := fixture.service.load("build demo"); err == nil || found {
		t.Fatalf("load = found:%v err:%v, want a refusal of a checkpoint ref outside the namespace", found, err)
	}
}

// The namespace alone is not identity: a ref another run or assignment would have written
// sits in the right namespace and still authorizes deleting someone else's checkpoint.
func TestReclaimTrustBoundaryRefusesAnotherAssignmentsCheckpointRef(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	run := loadRun(t, fixture.service)
	tamperedCheckpointRef(t, fixture, checkpointIdentity(run.Run, strings.Repeat("f", 32)))

	if _, found, err := fixture.service.load("build demo"); err == nil || found {
		t.Fatalf("load = found:%v err:%v, want a refusal of another assignment's checkpoint ref", found, err)
	}
}

// The enumeration is the second line behind the loader: a history entry reaches it without
// passing validCore, so it must re-derive the ref it claims rather than trust the stored
// string. The tampered record is handed to it in memory because the durable path now
// refuses to produce one.
func TestReclaimEnumerationNeverClaimsACheckpointRefTheRunDidNotWrite(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	run, assigned, branch := residueSubject(t, fixture)
	key, stored, ok := assignmentFor(run, assigned.ID)
	if !ok {
		t.Fatal("missing assignment")
	}
	foreign := git(t, fixture.root, "symbolic-ref", "HEAD")
	object := git(t, fixture.root, "rev-parse", foreign)
	stored.CheckpointRef = foreign
	run.Assignments[key] = stored
	run.Terminal = true

	refs, err := fixture.service.provisionalResidue(run)
	if err != nil {
		t.Fatalf("provisionalResidue: %v", err)
	}
	want := map[string]provisionalDisposition{
		branch:        refReclaimable,
		run.Candidate: refReclaimable,
	}
	if got := residueDispositions(refs); !reflect.DeepEqual(got, want) {
		t.Fatalf("inventory = %#v, want the foreign ref unplanned: %#v", got, want)
	}
	fixture.service.reclaimProvisionalRefs(run)
	if got := git(t, fixture.root, "rev-parse", foreign); got != object {
		t.Fatalf("%s = %s after reclamation, want the default branch intact at %s", foreign, got, object)
	}
}

// The empty CheckpointRef is every assignment's state between lease and checkpoint, so the
// validation must admit it or the normal path never reaches a checkpoint at all.
func TestReclaimValidationAdmitsAPreCheckpointAssignment(t *testing.T) {
	fixture := newCheckpointFixture(t)
	run := loadRun(t, fixture.service)
	_, stored, ok := assignmentFor(run, fixture.assigned.ID)
	if !ok || stored.CheckpointRef != "" {
		t.Fatalf("pre-checkpoint assignment = %#v, want an empty CheckpointRef to load", stored)
	}
	checkpointAssignment(t, fixture.root, fixture.service, fixture.assigned, []string{"internal/specbuild/checkpoint-change.go"})
	if _, err := fixture.service.Integrate(t.Context(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatalf("Integrate: %v", err)
	}
}

// The checkpoint namespace literal appears exactly once in the package's production
// sources — the checkpointIdentity construction — so the writer and both validators
// cannot drift apart by each stating the recipe themselves.
func TestReclaimCheckpointRefConstructionHasOneSource(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	sources := 0
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		source, err := os.ReadFile(entry.Name())
		if err != nil {
			t.Fatal(err)
		}
		sources += strings.Count(string(source), `"refs/bench/specbuild/checkpoint/"`)
	}
	if sources != 1 {
		t.Fatalf("checkpoint namespace stated %d times in production sources, want only the checkpointIdentity construction", sources)
	}
}
