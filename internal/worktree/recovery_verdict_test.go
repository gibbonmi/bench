package worktree

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// TestDiscardRecoveryRefusesAnUnresolvablePayload drives the real command surface with
// the exact fingerprint its own plan reported, because that is where the authority is
// spent: a payload a history rewrite garbage-collected plans as unclassifiable, and the
// discard must refuse rather than delete the one ref still naming the work.
func TestDiscardRecoveryRefusesAnUnresolvablePayload(t *testing.T) {
	root, assignment, recovery := preserveRecovery(t, "recovery-discard-unresolvable", "dirty.txt")
	// A well-formed OID nothing resolves; the record stays valid because a payload only
	// has to match the OID pattern, so no hand-edited state file is involved.
	const ghost = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	assignment.Recovery[0].Payloads = []string{ghost}
	mustNoError(t, intent.PutAssignment(root, assignment))

	plan, err := PlanRecovery(root, recovery.Ref)
	requireTest(t, err == nil && plan.Action == RecoveryRetain,
		"unresolvable-payload plan = %#v, %v; want the unclassifiable %q", plan, err, RecoveryRetain)

	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	code := RecoveryCommand([]string{recovery.Ref, "--discard", plan.Fingerprint}, &stdout, &stderr)
	requireTest(t, code == 1, "unclassifiable discard exit = %d, stdout %q stderr %q; want a refusal",
		code, stdout.String(), stderr.String())
	requireTest(t, !strings.Contains(stdout.String(), string(RecoveryDiscarded)),
		"refused discard receipt claims the work is gone: %q", stdout.String())
	requireTest(t, refExists(root, recovery.Ref) &&
		gitOutput(t, root, "rev-parse", "--verify", recovery.Ref+"^{commit}") == recovery.Root,
		"refused discard disturbed the recovery ref")
	current, err := assignmentByID(root, assignment.ID)
	requireTest(t, err == nil && current.State == intent.StateRecovered && len(current.Recovery) == 1,
		"refused discard changed the assignment row = %#v, %v", current, err)
}

// TestDiscardRecoveryRefusesEveryUnclassifiableExit enumerates the exits PlanRecovery
// takes without completing verification, because each returns through the same retain
// initialiser: covering one would leave the others free to regain destructive authority
// one at a time.
func TestDiscardRecoveryRefusesEveryUnclassifiableExit(t *testing.T) {
	for _, exit := range []struct {
		name   string
		strand func(t *testing.T, root string, assignment intent.Assignment, recovery intent.Recovery)
	}{
		{
			// Two row entries naming one ref leave the plan unable to say which envelope
			// the fingerprint would answer for.
			name: "ambiguous rows",
			strand: func(t *testing.T, root string, assignment intent.Assignment, recovery intent.Recovery) {
				assignment.Recovery = append(assignment.Recovery, recovery)
				mustNoError(t, intent.PutAssignment(root, assignment))
			},
		},
		{
			name: "assignment is not recovered",
			strand: func(t *testing.T, root string, assignment intent.Assignment, _ intent.Recovery) {
				markPending(t, root, assignment)
			},
		},
		{
			// The ref no longer resolves to the recorded root, so verifyRecovery refuses
			// and the landedness proof never runs.
			name: "verification fails",
			strand: func(t *testing.T, root string, _ intent.Assignment, recovery intent.Recovery) {
				gitRun(t, root, "update-ref", recovery.Ref, gitOutput(t, root, "rev-parse", "--verify", "refs/heads/main"))
			},
		},
		{
			// With two surviving local branches and main gone, ResolvedDefault has no
			// answer — one survivor alone would win the sole-local-branch fallback — so
			// the proof has no branch to judge the payloads against.
			name: "default branch does not resolve",
			strand: func(t *testing.T, root string, _ intent.Assignment, _ intent.Recovery) {
				head := gitOutput(t, root, "rev-parse", "--verify", "refs/heads/main")
				gitRun(t, root, "update-ref", "refs/heads/decoy-one", head)
				gitRun(t, root, "update-ref", "refs/heads/decoy-two", head)
				gitRun(t, root, "update-ref", "-d", "refs/heads/main")
			},
		},
	} {
		t.Run(exit.name, func(t *testing.T) {
			root, assignment, recovery := preserveRecovery(t, "recovery-unclassifiable-"+strings.ReplaceAll(exit.name, " ", "-"), "dirty.txt")
			exit.strand(t, root, assignment, recovery)

			plan, err := PlanRecovery(root, recovery.Ref)
			requireTest(t, err == nil && plan.Action == RecoveryRetain,
				"%s plan = %#v, %v; want the unclassifiable %q", exit.name, plan, err, RecoveryRetain)
			for _, verb := range []recoveryVerb{recoveryDiscard, recoveryRetire} {
				applied, err := applyRecoveryVerb(root, recovery.Ref, plan.Fingerprint, verb)
				requireTest(t, errors.Is(err, errRecoveryUnauthorized),
					"%s over %s = %#v, %v; want an unauthorized refusal", verb, exit.name, applied, err)
			}
			requireTest(t, refExists(root, recovery.Ref), "%s: a refused verb deleted the recovery ref", exit.name)
			current, err := assignmentByID(root, assignment.ID)
			requireTest(t, err == nil && len(current.Recovery) >= 1,
				"%s: a refused verb changed the assignment row = %#v, %v", exit.name, current, err)
		})
	}
}
