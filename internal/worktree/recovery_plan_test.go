package worktree

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// preserveRecovery drives the real preservation path — a dirty owned checkout released
// through plan and apply — so every recovery assertion below reads a ref and an envelope
// production actually writes rather than one the test invented.
func preserveRecovery(t *testing.T, request string, dirty ...string) (string, intent.Assignment, intent.Recovery) {
	t.Helper()
	root, creation := newOwnedAssignment(t, request)
	for _, name := range dirty {
		mustWrite(t, filepath.Join(creation.Path, name), []byte("preserve "+name+"\n"), 0o644)
	}
	plan, err := PlanExplicit(root, creation.Path)
	mustNoError(t, err)
	requireTest(t, plan.Action == ActionRecoverRemove, "preservation plan = %#v", plan)
	result, err := ApplyExplicit(root, creation.Path, plan.Fingerprint)
	requireTest(t, err == nil && result.Action == ActionRemoved, "preservation apply = %#v, %v", result, err)
	assignment, err := assignmentByID(root, creation.Assignment.ID)
	mustNoError(t, err)
	requireTest(t, assignment.State == intent.StateRecovered && len(assignment.Recovery) == 1,
		"preserved assignment = %#v", assignment)
	return root, assignment, assignment.Recovery[0]
}

func planRecoveryReceipt(t *testing.T, plan RecoveryPlan) string {
	t.Helper()
	var out bytes.Buffer
	mustNoError(t, renderRecovery(&out, plan))
	return out.String()
}

func TestPlanRecoveryReportsOrphanedRefWithNoOwningRow(t *testing.T) {
	root, assignment, recovery := preserveRecovery(t, "recovery-orphaned-row", "dirty.txt")
	mustNoError(t, intent.DeleteAssignment(root, assignment.ID))

	plan, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && plan.Action == RecoveryOrphaned,
		"row-less recovery plan = %#v, %v; want action %q", plan, err, RecoveryOrphaned)
	receipt := planRecoveryReceipt(t, plan)
	requireTest(t, strings.Contains(receipt, string(RecoveryOrphaned)) &&
		!strings.Contains(receipt, string(RecoveryRetain)) && !strings.Contains(receipt, string(RecoveryAbsent)),
		"orphaned receipt = %q", receipt)
	requireTest(t, gitOutput(t, root, "rev-parse", "--verify", recovery.Ref+"^{commit}") == recovery.Root,
		"planning an orphaned ref disturbed it")
}

// TestPlanRecoveryConfinesTheOrphanedVerdictToTheRecoveryNamespace grades the row-less
// path over refs that are not recovery refs at all, because that is the only path where
// the orphaned verdict is decided by ref existence alone: every ref outside the namespace
// exists just as convincingly as one inside it.
func TestPlanRecoveryConfinesTheOrphanedVerdictToTheRecoveryNamespace(t *testing.T) {
	root, assignment, recovery := preserveRecovery(t, "recovery-namespace-confinement", "dirty.txt")
	mustNoError(t, intent.DeleteAssignment(root, assignment.ID))
	head := gitOutput(t, root, "rev-parse", "--verify", "refs/heads/main")

	for _, foreign := range []struct{ name, ref string }{
		{"the default branch", "refs/heads/main"},
		{"an ordinary branch", "refs/heads/ordinary-work"},
		{"an assignment branch", "refs/heads/bench/assign/owner/assignment"},
		{"a namespace lookalike", "refs/bench/recovery-shadow/owner/assignment/1"},
		{"a namespace prefix", "refs/bench/recover/owner/assignment/1"},
	} {
		t.Run(foreign.name, func(t *testing.T) {
			if foreign.ref != "refs/heads/main" {
				gitRun(t, root, "update-ref", foreign.ref, head)
			}
			plan, err := PlanRecovery(root, foreign.ref)
			requireTest(t, err == nil && plan.Action == RecoveryForeign,
				"plan for %s = %#v, %v; want action %q", foreign.ref, plan, err, RecoveryForeign)
			receipt := planRecoveryReceipt(t, plan)
			requireTest(t, !strings.Contains(receipt, string(RecoveryOrphaned)),
				"receipt for %s reads as preserved work: %q", foreign.ref, receipt)
			requireTest(t, gitOutput(t, root, "rev-parse", "--verify", foreign.ref) == head,
				"planning %s disturbed it", foreign.ref)
		})
	}

	// The capability the confinement must leave intact: a genuine recovery ref whose row is
	// gone is still the operator's to discard, and it has no other route out.
	plan, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && plan.Action == RecoveryOrphaned,
		"genuine row-less recovery plan = %#v, %v; want action %q", plan, err, RecoveryOrphaned)
}

// TestDiscardRecoveryRefusesARefOutsideTheRecoveryNamespace drives the verb with the exact
// fingerprint its own plan reported, so the refusal is the namespace and not a stale value.
func TestDiscardRecoveryRefusesARefOutsideTheRecoveryNamespace(t *testing.T) {
	root, assignment, recovery := preserveRecovery(t, "recovery-discard-foreign", "dirty.txt")
	mustNoError(t, intent.DeleteAssignment(root, assignment.ID))
	const branch = "refs/heads/work-nobody-preserved"
	gitRun(t, root, "update-ref", branch, gitOutput(t, root, "rev-parse", "--verify", "refs/heads/main"))
	before := gitOutput(t, root, "rev-parse", "--verify", branch)

	plan, err := PlanRecovery(root, branch)
	requireTest(t, err == nil && plan.Action == RecoveryForeign,
		"branch plan = %#v, %v; want action %q", plan, err, RecoveryForeign)
	for _, verb := range []recoveryVerb{recoveryDiscard, recoveryRetire} {
		applied, err := applyRecoveryVerb(root, branch, plan.Fingerprint, verb)
		requireTest(t, errors.Is(err, errRecoveryUnauthorized),
			"%s over an ordinary branch = %#v, %v; want an unauthorized refusal", verb, applied, err)
	}
	requireTest(t, refExists(root, branch) && gitOutput(t, root, "rev-parse", "--verify", branch) == before,
		"a refused verb deleted or moved the ordinary branch %s", branch)
	requireTest(t, refExists(root, recovery.Ref), "a refused verb deleted the recovery ref")
}

// TestRecoveryFingerprintSealsTheChangeSummary mutates the reported summary and nothing
// else. The count is what story 3 offers the operator before a destructive act, so a
// fingerprint blind to it authorizes a discard the operator agreed to under a different
// number.
func TestRecoveryFingerprintSealsTheChangeSummary(t *testing.T) {
	root, _, recovery := preserveRecovery(t, "recovery-fingerprint-changes", "one.txt", "two.txt")
	canonical := canonicalRoot(root)
	ledger, err := intent.LifecycleEvidence(canonical)
	mustNoError(t, err)
	plan, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && plan.Changes == "paths=2", "planned changes = %#v, %v", plan, err)

	// The control: without it every comparison below would pass on a derivation that had
	// simply stopped answering for this plan.
	requireTest(t, fingerprintRecovery(canonical, ledger, plan).Fingerprint == plan.Fingerprint,
		"the derivation no longer reproduces the planned fingerprint")

	for _, changes := range []string{"paths=999", "paths=1", recoveryUnknownChanges, ""} {
		t.Run(changes, func(t *testing.T) {
			mutated := plan
			mutated.Changes = changes
			requireTest(t, fingerprintRecovery(canonical, ledger, mutated).Fingerprint != plan.Fingerprint,
				"a plan reporting changes %q fingerprints to the value planned for %q", changes, plan.Changes)
		})
	}
}

// TestRecoveryNamespaceIsReadFromIntent holds the namespace to one definition. The literal
// is never written here either: it is quoted from the value under test, so this file cannot
// become the second copy it exists to forbid.
func TestRecoveryNamespaceIsReadFromIntent(t *testing.T) {
	namespace := recoveryNamespace()
	prefix := intent.RecoveryRefPrefix("owner", "assignment")
	requireTest(t, namespace != "" && strings.TrimPrefix(prefix, namespace) == "owner/assignment/",
		"recovery namespace %q is not the head of intent's recovery prefix %q", namespace, prefix)

	literal := strconv.Quote(namespace)
	for _, name := range packageSourceFiles(t) {
		source, err := os.ReadFile(name)
		mustNoError(t, err)
		requireTest(t, !strings.Contains(string(source), literal),
			"%s restates the recovery namespace as %s; read intent.RecoveryRefPrefix instead", name, literal)
	}
}

// packageSourceFiles lists this package's non-test sources by their own location rather
// than by the working directory, which a test elsewhere in the package may have moved.
func packageSourceFiles(t *testing.T) []string {
	t.Helper()
	_, self, _, ok := runtime.Caller(0)
	requireTest(t, ok, "the test binary does not report its own source location")
	dir := filepath.Dir(self)
	entries, err := os.ReadDir(dir)
	mustNoError(t, err)
	var sources []string
	for _, entry := range entries {
		if name := entry.Name(); strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			sources = append(sources, filepath.Join(dir, name))
		}
	}
	requireTest(t, len(sources) > 0, "no package sources found under %s", dir)
	return sources
}

func TestPlanRecoveryReportsAbsentRefAsDistinctFromOrphaned(t *testing.T) {
	root, assignment, recovery := preserveRecovery(t, "recovery-absent-ref", "dirty.txt")
	mustNoError(t, intent.DeleteAssignment(root, assignment.ID))
	gitRun(t, root, "update-ref", "-d", recovery.Ref)

	plan, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && plan.Action == RecoveryAbsent,
		"nonexistent recovery plan = %#v, %v; want action %q", plan, err, RecoveryAbsent)
	receipt := planRecoveryReceipt(t, plan)
	requireTest(t, strings.Contains(receipt, string(RecoveryAbsent)) && !strings.Contains(receipt, string(RecoveryOrphaned)),
		"absent receipt = %q", receipt)
}

func TestPlanRecoveryKeepsRetireAndDiscardVerdicts(t *testing.T) {
	t.Run("unlanded payload plans discard", func(t *testing.T) {
		root, _, recovery := preserveRecovery(t, "recovery-discard-unlanded", "dirty.txt")
		plan, err := PlanRecovery(root, recovery.Ref)
		requireTest(t, err == nil && plan.Action == RecoveryDiscard && strings.Contains(plan.Landed, "unlanded"),
			"unlanded recovery plan = %#v, %v", plan, err)
	})
	t.Run("landed payloads retire", func(t *testing.T) {
		root, _, recovery := preserveRecovery(t, "recovery-retire-landed", "dirty.txt")
		// The recovery root's parents are exactly its payloads, so landing the root lands
		// every payload at once however many layers the checkout produced.
		gitRun(t, root, "update-ref", "refs/heads/main", recovery.Root)
		plan, err := PlanRecovery(root, recovery.Ref)
		requireTest(t, err == nil && plan.Action == RecoveryRetire,
			"landed recovery plan = %#v, %v", plan, err)
	})
}

func TestPlanRecoveryReportsChangedPathCount(t *testing.T) {
	root, assignment, recovery := preserveRecovery(t, "recovery-change-count", "one.txt", "two.txt")
	const want = "paths=2"

	unproven, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && unproven.Action == RecoveryDiscard && unproven.Changes == want,
		"discard-eligible plan changes = %#v, %v; want %q", unproven, err, want)
	requireTest(t, strings.Contains(planRecoveryReceipt(t, unproven), want),
		"discard-eligible receipt = %q", planRecoveryReceipt(t, unproven))

	gitRun(t, root, "update-ref", "refs/heads/main", recovery.Root)
	retired, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && retired.Action == RecoveryRetire && retired.Changes == want,
		"retire plan changes = %#v, %v; want %q", retired, err, want)

	mustNoError(t, intent.DeleteAssignment(root, assignment.ID))
	orphaned, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && orphaned.Action == RecoveryOrphaned && orphaned.Changes == want,
		"orphaned plan changes = %#v, %v; want %q", orphaned, err, want)

	gitRun(t, root, "update-ref", "-d", recovery.Ref)
	absent, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && absent.Action == RecoveryAbsent && absent.Changes == recoveryUnknownChanges,
		"absent plan changes = %#v, %v; want %q", absent, err, recoveryUnknownChanges)
	requireTest(t, strings.Contains(planRecoveryReceipt(t, absent), recoveryUnknownChanges),
		"absent receipt = %q", planRecoveryReceipt(t, absent))
}

func TestPlanRecoveryReportsUnknownChangesWhenEnvelopeDoesNotResolve(t *testing.T) {
	const missing = "cafe0000cafe0000cafe0000cafe0000cafe0000"
	for _, spec := range []struct {
		name     string
		envelope func(layer string) (base, payload string)
	}{
		{
			name:     "base-does-not-resolve",
			envelope: func(layer string) (string, string) { return missing, layer },
		},
		{
			name:     "payload-does-not-resolve",
			envelope: func(layer string) (string, string) { return layer, missing },
		},
	} {
		t.Run(spec.name, func(t *testing.T) {
			root, assignment, recovery := preserveRecovery(t, "recovery-unknown-"+spec.name, "dirty.txt")
			mustNoError(t, intent.DeleteAssignment(root, assignment.ID))
			base, payload := spec.envelope(recovery.Payloads[0])
			rewriteRecoveryEnvelope(t, root, recovery, base, payload)

			plan, err := PlanRecovery(root, recovery.Ref)
			requireTest(t, err == nil && plan.Changes == recoveryUnknownChanges,
				"unresolvable envelope plan = %#v, %v; want changes %q and no error", plan, err, recoveryUnknownChanges)
		})
	}
}

func TestDiscardRecoveryDropsAnUnprovenPayload(t *testing.T) {
	t.Run("last ref closes the assignment row", func(t *testing.T) {
		root, assignment, recovery := preserveRecovery(t, "recovery-discard-unproven", "dirty.txt")
		plan, err := PlanRecovery(root, recovery.Ref)
		requireTest(t, err == nil && plan.Action == RecoveryDiscard, "unproven plan = %#v, %v", plan, err)

		applied, err := applyRecoveryVerb(root, recovery.Ref, plan.Fingerprint, recoveryDiscard)
		requireTest(t, err == nil && applied.Action == RecoveryDiscarded,
			"discard = %#v, %v; want action %q", applied, err, RecoveryDiscarded)
		receipt := planRecoveryReceipt(t, applied)
		requireTest(t, strings.Contains(receipt, string(RecoveryDiscarded)) && !strings.Contains(receipt, string(RecoveryRetired)),
			"discard receipt = %q; want it separable from the retire claim", receipt)
		requireTest(t, !refExists(root, recovery.Ref), "discarded recovery ref survived")
		if _, err := assignmentByID(root, assignment.ID); err == nil {
			t.Fatal("last-ref discard did not compact the recovered assignment")
		}
	})
	t.Run("a surviving sibling leaves the row recovered", func(t *testing.T) {
		root, assignment, recovery := preserveRecovery(t, "recovery-discard-sibling", "dirty.txt")
		sibling := addSiblingRecovery(t, root, assignment, recovery)

		plan, err := PlanRecovery(root, recovery.Ref)
		mustNoError(t, err)
		applied, err := applyRecoveryVerb(root, recovery.Ref, plan.Fingerprint, recoveryDiscard)
		requireTest(t, err == nil && applied.Action == RecoveryDiscarded, "sibling discard = %#v, %v", applied, err)
		requireTest(t, !refExists(root, recovery.Ref) && refExists(root, sibling.Ref),
			"sibling discard removed the wrong refs")
		current, err := assignmentByID(root, assignment.ID)
		requireTest(t, err == nil && current.State == intent.StateRecovered && len(current.Recovery) == 1 &&
			current.Recovery[0].Ref == sibling.Ref, "intermediate discarded state = %#v, %v", current, err)
	})
}

// addSiblingRecovery gives an assignment a second recovery ref pointing at the same
// envelope, so a discard of the first has something left to keep the row open.
func addSiblingRecovery(t *testing.T, root string, assignment intent.Assignment, recovery intent.Recovery) intent.Recovery {
	t.Helper()
	sibling := recovery
	sibling.Ref = strings.TrimSuffix(recovery.Ref, "/1") + "/2"
	gitRun(t, root, "update-ref", sibling.Ref, sibling.Root)
	assignment.Recovery = append(assignment.Recovery, sibling)
	mustNoError(t, intent.PutAssignment(root, assignment))
	return sibling
}

func TestDiscardRecoveryRefusesAStaleFingerprint(t *testing.T) {
	root, assignment, recovery := preserveRecovery(t, "recovery-discard-stale", "dirty.txt")
	plan, err := PlanRecovery(root, recovery.Ref)
	mustNoError(t, err)
	// The stale value differs only in its last character, so a comparison that reads a
	// prefix of the fingerprint rather than the whole of it still accepts it.
	stale := plan.Fingerprint[:len(plan.Fingerprint)-1] + "0"
	if strings.HasSuffix(plan.Fingerprint, "0") {
		stale = plan.Fingerprint[:len(plan.Fingerprint)-1] + "1"
	}

	applied, err := applyRecoveryVerb(root, recovery.Ref, stale, recoveryDiscard)
	requireTest(t, errors.Is(err, errStaleFingerprint), "stale discard = %#v, %v; want a stale-fingerprint refusal", applied, err)
	requireTest(t, gitOutput(t, root, "rev-parse", "--verify", recovery.Ref+"^{commit}") == recovery.Root,
		"stale discard disturbed the recovery ref")
	current, err := assignmentByID(root, assignment.ID)
	requireTest(t, err == nil && current.State == intent.StateRecovered && len(current.Recovery) == 1,
		"stale discard changed the assignment row = %#v, %v", current, err)
}

func TestDiscardRecoveryRefusesWhatTheProofAccepts(t *testing.T) {
	t.Run("a retire plan keeps --apply as its route", func(t *testing.T) {
		root, assignment, recovery := preserveRecovery(t, "recovery-discard-landed", "dirty.txt")
		gitRun(t, root, "update-ref", "refs/heads/main", recovery.Root)
		plan, err := PlanRecovery(root, recovery.Ref)
		requireTest(t, err == nil && plan.Action == RecoveryRetire, "landed plan = %#v, %v", plan, err)

		applied, err := applyRecoveryVerb(root, recovery.Ref, plan.Fingerprint, recoveryDiscard)
		requireTest(t, err != nil, "discard accepted a proven-landed payload: %#v", applied)
		requireTest(t, refExists(root, recovery.Ref), "refused discard deleted the ref")
		current, err := assignmentByID(root, assignment.ID)
		requireTest(t, err == nil && current.State == intent.StateRecovered && len(current.Recovery) == 1,
			"refused discard changed the assignment row = %#v, %v", current, err)
	})
	t.Run("a row that is not recovered is not the discard's to close", func(t *testing.T) {
		root, assignment, recovery := preserveRecovery(t, "recovery-discard-pending", "dirty.txt")
		markPending(t, root, assignment)
		plan, err := PlanRecovery(root, recovery.Ref)
		mustNoError(t, err)

		applied, err := applyRecoveryVerb(root, recovery.Ref, plan.Fingerprint, recoveryDiscard)
		requireTest(t, err != nil, "discard closed a row mid-release: %#v", applied)
		requireTest(t, refExists(root, recovery.Ref), "refused discard deleted the ref")
	})
}

// TestRecoveryVerbAuthorizesOnlyItsOwnVerdicts enumerates the whole action vocabulary
// rather than sampling one verdict, so a verb that loses its guard is caught for every
// action the plan can carry, not only the one a behavioral test happens to drive.
func TestRecoveryVerbAuthorizesOnlyItsOwnVerdicts(t *testing.T) {
	for _, spec := range []struct {
		action          RecoveryAction
		retire, discard bool
	}{
		{action: RecoveryRetain, retire: false, discard: false},
		{action: RecoveryRetire, retire: true, discard: false},
		{action: RecoveryDiscard, retire: false, discard: true},
		{action: RecoveryRetired, retire: false, discard: false},
		{action: RecoveryDiscarded, retire: false, discard: false},
		{action: RecoveryError, retire: false, discard: false},
		{action: RecoveryOrphaned, retire: false, discard: true},
		{action: RecoveryAbsent, retire: false, discard: false},
		{action: RecoveryForeign, retire: false, discard: false},
	} {
		t.Run(string(spec.action), func(t *testing.T) {
			for _, verb := range []struct {
				verb recoveryVerb
				want bool
			}{{recoveryRetire, spec.retire}, {recoveryDiscard, spec.discard}} {
				authorized, detail := verb.verb.authorizes(spec.action)
				requireTest(t, authorized == verb.want, "%s authorizes %q = %t; want %t",
					verb.verb, spec.action, authorized, verb.want)
				requireTest(t, authorized == (detail == ""), "%s refusal detail for %q = %q",
					verb.verb, spec.action, detail)
			}
		})
	}
}

// TestRecoveryCommandRefusesAMalformedFingerprintForEveryVerb drives the command surface
// because the format control is an argument-parse control: it must refuse before the plan
// is ever computed, for whichever verb carries the fingerprint.
func TestRecoveryCommandRefusesAMalformedFingerprintForEveryVerb(t *testing.T) {
	const ref = "refs/bench/recovery/owner/assignment/1"
	for _, verb := range []string{"--apply", "--discard"} {
		for _, fingerprint := range []string{"", "bad", strings.Repeat("a", 63), strings.Repeat("a", 65), strings.Repeat("A", 64), strings.Repeat("g", 64)} {
			var stdout, stderr bytes.Buffer
			code := RecoveryCommand([]string{ref, verb, fingerprint}, &stdout, &stderr)
			requireTest(t, code == 2, "recovery %s %q exit = %d; want 2", verb, fingerprint, code)
			requireTest(t, strings.Contains(stdout.String(), string(RecoveryError)) && stderr.Len() == 0,
				"recovery %s %q streams = stdout %q stderr %q", verb, fingerprint, stdout.String(), stderr.String())
		}
	}
	var stdout, stderr bytes.Buffer
	requireTest(t, RecoveryCommand([]string{ref, "--drop", strings.Repeat("a", 64)}, &stdout, &stderr) == 2,
		"an unknown recovery verb was accepted: stdout %q stderr %q", stdout.String(), stderr.String())
}

func TestDiscardRecoveryDropsAnOrphanedRefAndConvergesOnAbsent(t *testing.T) {
	root, assignment, recovery := preserveRecovery(t, "recovery-discard-orphan", "dirty.txt")
	mustNoError(t, intent.DeleteAssignment(root, assignment.ID))
	before, err := intent.Assignments(root)
	mustNoError(t, err)

	plan, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && plan.Action == RecoveryOrphaned, "orphaned plan = %#v, %v", plan, err)
	applied, err := applyRecoveryVerb(root, recovery.Ref, plan.Fingerprint, recoveryDiscard)
	requireTest(t, err == nil && applied.Action == RecoveryDiscarded, "orphaned discard = %#v, %v", applied, err)
	requireTest(t, !refExists(root, recovery.Ref), "discarded orphaned ref survived")
	after, err := intent.Assignments(root)
	requireTest(t, err == nil && len(after) == len(before),
		"orphaned discard touched the intent records: before %#v after %#v", before, after)

	absent, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && absent.Action == RecoveryAbsent, "re-plan = %#v, %v", absent, err)
	rerun, err := applyRecoveryVerb(root, recovery.Ref, absent.Fingerprint, recoveryDiscard)
	requireTest(t, err == nil && rerun.Action == RecoveryAbsent,
		"absent discard = %#v, %v; want a no-op success", rerun, err)
}

// TestApplyRecoveryRefusesAPlanThatProvesNothingLanded pins the deliberate change to
// --apply: a caller who supplied a fingerprint asked for an action, and exit-zero silence
// reads as "the work is gone" when it is still there.
func TestApplyRecoveryRefusesAPlanThatProvesNothingLanded(t *testing.T) {
	root, assignment, recovery := preserveRecovery(t, "recovery-apply-unproven", "dirty.txt")
	plan, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && plan.Action == RecoveryDiscard, "unproven plan = %#v, %v", plan, err)

	applied, err := ApplyRecovery(root, recovery.Ref, plan.Fingerprint)
	requireTest(t, err != nil, "discard-eligible apply returned silent success: %#v", applied)
	requireTest(t, refExists(root, recovery.Ref), "refused apply deleted the ref")
	current, err := assignmentByID(root, assignment.ID)
	requireTest(t, err == nil && current.State == intent.StateRecovered && len(current.Recovery) == 1,
		"refused apply changed the assignment row = %#v, %v", current, err)
}

// TestRecoveryFingerprintCommitsToTheDiscardAuthority pins both halves of the authority
// the fingerprint carries. The retire-only values are written here as literals rather than
// read from the subject, so a rename in the subject cannot carry the expectation with it:
// the question this asks is whether a fingerprint planned under the exact authority that
// preceded the discard is still a different value, which is the security property itself.
func TestRecoveryFingerprintCommitsToTheDiscardAuthority(t *testing.T) {
	const (
		retireOnlyDomain  = "bench-recovery-retire/v1"
		retireOnlyEffects = "delete-exact-ref,update-assignment,compact-if-last"
	)
	root, _, recovery := preserveRecovery(t, "recovery-fingerprint-authority", "dirty.txt")
	ledger, err := intent.LifecycleEvidence(canonicalRoot(root))
	mustNoError(t, err)
	plan, err := PlanRecovery(root, recovery.Ref)
	mustNoError(t, err)

	// The control: under the authority the subject actually ships, this plan must
	// reproduce the fingerprint it reported. Without it every comparison below would pass
	// on a derivation that had simply stopped answering for this plan at all.
	shipped := fingerprintRecoveryUnder(canonicalRoot(root), ledger, plan, recoveryFingerprintDomain, recoveryFingerprintEffects)
	requireTest(t, shipped.Fingerprint == plan.Fingerprint,
		"the shipped authority does not reproduce the planned fingerprint: %q vs %q", shipped.Fingerprint, plan.Fingerprint)

	for _, authority := range []struct{ name, domain, effects string }{
		{"the retire-only domain tag", retireOnlyDomain, recoveryFingerprintEffects},
		{"the retire-only effect string", recoveryFingerprintDomain, retireOnlyEffects},
		{"the whole retire-only authority", retireOnlyDomain, retireOnlyEffects},
	} {
		t.Run(authority.name, func(t *testing.T) {
			under := fingerprintRecoveryUnder(canonicalRoot(root), ledger, plan, authority.domain, authority.effects)
			requireTest(t, under.Fingerprint != plan.Fingerprint,
				"%s fingerprints this plan to %q, the value the discard authority reports; a fingerprint planned before the discard existed would authorize one",
				authority.name, under.Fingerprint)
		})
	}
	requireTest(t, strings.Contains(recoveryFingerprintDomain, "discard") && strings.Contains(recoveryFingerprintEffects, "discard"),
		"the recovery fingerprint authority does not name the discard: domain %q effects %q",
		recoveryFingerprintDomain, recoveryFingerprintEffects)
}

// rewriteRecoveryEnvelope repoints a recovery ref at an envelope naming the given base and
// payload, standing in for a history rewrite that leaves either unreachable.
func rewriteRecoveryEnvelope(t *testing.T, root string, recovery intent.Recovery, base, payload string) {
	t.Helper()
	manifest, err := json.Marshal(recoveryManifest{Schema: recoverySchema, Base: base, Layers: map[string]string{"working": payload}})
	mustNoError(t, err)
	blob, err := gitInput(root, nil, append(manifest, '\n'), "hash-object", "-w", "--stdin")
	mustNoError(t, err)
	tree, err := gitInput(root, nil, []byte("100644 blob "+blob+"\tmanifest.json\n"), "mktree")
	mustNoError(t, err)
	rewritten, err := commitTree(root, tree, []string{recovery.Root}, "rewritten recovery root\n")
	mustNoError(t, err)
	gitRun(t, root, "update-ref", recovery.Ref, rewritten)
}
