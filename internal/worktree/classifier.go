package worktree

import (
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CleanupAction string

func renderCleanup(stdout io.Writer, plan CleanupPlan) error {
	tracked, recovery := plan.Tracked, plan.Recovery
	if tracked == "" {
		tracked = "unknown"
	}
	if recovery == "" {
		recovery = "none"
	}
	detail := plan.Reason
	if detail == "" && plan.Action != ActionRemoved {
		detail = "apply with exact fingerprint"
	}
	// An operator-asserted branch deletion is the one removal this command performs that
	// nothing else in the row implies, so the plan half names the exact ref it will spend
	// that assertion on. The derived case stays silent: its deletion follows from the
	// landedness the tool proved for itself, and its row is a settled output contract.
	if plan.discardBranch && plan.deleteBranch && plan.Action.removes() {
		detail = "discards branch " + plan.branchRef + "; " + detail
	}
	ignored := plan.ignoredSummary
	if ignored == "" {
		ignored = plan.Ignored.Summary()
	}
	target := cleanupOutputValue(plan.Target)
	if !cleanupOutputSafe(detail) {
		detail = "unsafe detail " + cleanupOutputValue(detail)
	}
	out, err := toon.Table("worktree_cleanup", []string{"target", "action", "tracked", "ignored", "recovery", "fingerprint", "detail"}, [][]string{{target, string(plan.Action), tracked, ignored, recovery, plan.Fingerprint, detail}})
	if err != nil {
		return err
	}
	if _, err = fmt.Fprint(stdout, out); err != nil || plan.Ignored.Shown == 0 {
		return err
	}
	rows := make([][]string, 0, plan.Ignored.Shown)
	for _, path := range plan.Ignored.Paths[:plan.Ignored.Shown] {
		rows = append(rows, []string{path})
	}
	preview, err := toon.Table("ignored_paths", []string{"path"}, rows)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(stdout, preview)
	return err
}

type CleanupReason string

const (
	ActionRetain         CleanupAction = "retain"
	ActionRemove         CleanupAction = "remove"
	ActionRecoverRemove  CleanupAction = "recover-remove"
	ActionDiscardRemove  CleanupAction = "discard-remove"
	ActionRemoved        CleanupAction = "removed"
	ActionError          CleanupAction = "error"
	ReasonForeign        CleanupReason = "foreign"
	ReasonActive         CleanupReason = "active"
	ReasonLiveLease      CleanupReason = "live-lease"
	ReasonUnmerged       CleanupReason = "unmerged"
	ReasonIgnored        CleanupReason = "ignored"
	ReasonMalformed      CleanupReason = "malformed"
	ReasonUncertain      CleanupReason = "uncertain"
	ReasonUnexpectedLock CleanupReason = "unexpected-lock"
	ReasonOrphaned       CleanupReason = "orphaned"
)
const actionReleaseRemove CleanupAction = "release-remove"

// actionReleaseLeftover releases one assignment's registration and ledger entry while the
// bytes at its path stay exactly where they are. It is deliberately outside removes():
// nothing proves what those bytes are — no checkout answers for them and no recovery ref
// holds them — so disposing of them stays with the path-addressed clean surface, whose
// inventory is size-bounded.
const actionReleaseLeftover CleanupAction = "release-leftover"

// removes reports whether an action still has a removal ahead of it, as opposed to
// reporting a refusal, an invocation error, or a transaction that already completed.
func (action CleanupAction) removes() bool {
	return action == ActionRemove || action == ActionRecoverRemove || action == ActionDiscardRemove || action == actionReleaseRemove
}

var ignoredLstat = os.Lstat

// PathShape names what the entry at an assignment path is.
type PathShape string

const (
	ShapeAbsent            PathShape = "absent"
	ShapeDanglingSymlink   PathShape = "dangling-symlink"
	ShapeNonDirectory      PathShape = "non-directory"
	ShapeDecayedDirectory  PathShape = "directory-without-git-metadata"
	ShapeCheckoutDirectory PathShape = "directory-with-git-metadata"
	// ShapeSpecialMetadata names a directory whose .git entry exists but is neither a
	// regular file (a gitfile worktree pointer) nor a directory (an embedded repository):
	// a FIFO, device, socket, or a symlinked .git, any of which git itself refuses to
	// treat as ordinary metadata. No consumer may invoke git against this path — a FIFO
	// with no writer would block the invocation forever — so this shape fails closed
	// rather than joining ShapeCheckoutDirectory or the decayed set.
	ShapeSpecialMetadata PathShape = "special-git-metadata"
	ShapeUnknown         PathShape = "unknown"
)

// ClassifyPathShape is the single source for the decayed-shape policy: every consumer
// asking whether an assignment's checkout is still live, or whether an abandon is
// releasing residue rather than removing a checkout, decides it here, so the two can
// never answer differently for the same bytes.
//
// The path is never opened — a FIFO at an assignment path has no writer and would block a
// reader forever — so the verdict rests on lstat and stat shape alone, and only
// ShapeCheckoutDirectory licenses a caller to run git against the path. Symlinks are
// followed, because a resolvable one is already resolved away by the time a canonical
// target names it, so one surviving here resolves to nothing.
//
// ShapeUnknown is the only shape carrying an error, the stat failure that left the shape
// undecided. It is never absence: an unreadable live checkout stats exactly this way.
func ClassifyPathShape(path string) (PathShape, error) {
	if _, err := os.Lstat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ShapeAbsent, nil
		}
		return ShapeUnknown, err
	}
	entry, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ShapeDanglingSymlink, nil
		}
		return ShapeUnknown, err
	}
	if !entry.IsDir() {
		return ShapeNonDirectory, nil
	}
	gitEntry, err := os.Lstat(filepath.Join(path, ".git"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ShapeDecayedDirectory, nil
		}
		return ShapeUnknown, err
	}
	// Only a regular file (a gitfile worktree pointer) or a directory (an ordinary or
	// embedded repository) is metadata git itself will read; everything else — including
	// a symlinked .git, which git refuses on principle — is classified without opening it.
	if gitEntry.Mode().IsRegular() || gitEntry.IsDir() {
		return ShapeCheckoutDirectory, nil
	}
	return ShapeSpecialMetadata, nil
}

// decayed reports whether a shape is one an abandon releases as leftover bytes: present,
// but nothing a checkout answers for.
func (shape PathShape) decayed() bool {
	return shape == ShapeDanglingSymlink || shape == ShapeNonDirectory || shape == ShapeDecayedDirectory
}

type CleanupPlan struct {
	Target               string
	Action               CleanupAction
	ReasonCode           CleanupReason
	Reason               string
	Assignment           string
	Tracked              string
	Ignored              IgnoredInventory
	Recovery             string
	Fingerprint          string
	owned                bool
	assignment           *intent.Assignment
	registration         git.Worktree
	discardIgnored       bool
	discardBranch        bool
	deleteBranch         bool
	branchRef, branchOID string
	ignoredSummary       string
	landed               string
	// leftover names the present bytes a release-leftover plan hands on rather than
	// removes; it is empty for every plan that answers for a checkout.
	leftover string
}
type IgnoredInventory struct {
	Count     int
	Bytes     int64
	Shown     int
	Truncated bool
	Paths     []string
	AtLeast   bool
	OverLimit bool
	Uncertain bool
	Digest    string
}
type RecoveryAction string

const (
	RecoveryRetain                 RecoveryAction = "retain"
	RecoveryRetire                 RecoveryAction = "retire"
	RecoveryRetired, RecoveryError RecoveryAction = "retired", "error"
	// RecoveryOrphaned names a ref that exists with no assignment row claiming it, and
	// RecoveryAbsent a ref name nothing resolves. They are separate verdicts because one
	// still holds preserved work and the other is already gone; collapsing them onto
	// retain leaves the first unreachable and makes a re-run of the second look actionable.
	RecoveryOrphaned RecoveryAction = "orphaned"
	RecoveryAbsent   RecoveryAction = "absent"
	// RecoveryForeign names a ref outside the recovery namespace. It is neither orphaned
	// nor absent, because both of those are claims about preserved work and this ref holds
	// none: no verb may spend a deletion on it, whether or not it resolves.
	RecoveryForeign RecoveryAction = "foreign"
	// RecoveryDiscarded names work an operator dropped without the proof accepting it.
	// It stays separate from RecoveryRetired because the receipt is the only durable
	// record of which of the two claims was made about the same disappearance.
	RecoveryDiscarded RecoveryAction = "discarded"
	// RecoveryDiscard is the one verdict that makes a ref the operator's to drop:
	// verification succeeded, the default branch resolved, and the landedness proof
	// refused the payloads, so the preserved work is real and only unproven.
	// RecoveryRetain is every exit that never completed that chain, and it authorizes
	// nothing: while retain doubled as the discard-eligible verdict, a payload a
	// history rewrite had garbage-collected planned retain and --discard deleted the
	// one ref still naming the work.
	RecoveryDiscard RecoveryAction = "discard"
)

// recoveryUnknownChanges is the plan's answer when the recovery envelope, its recorded
// base, or a payload no longer resolves. A rewritten history must not make a ref
// unplannable, so the summary degrades to this value instead of raising an error.
const recoveryUnknownChanges = "unknown"

type RecoveryPlan struct {
	Ref, Root, Payloads, Landed string
	// Changes summarizes what the payload does to its recorded base, derived at plan time
	// and never stored: the record would hold a second copy of a fact Git already answers,
	// and that copy goes stale the moment the base is rewritten.
	Changes             string
	Action              RecoveryAction
	Fingerprint, Detail string
	assignment          *intent.Assignment
	recovery            intent.Recovery
}

func (inventory IgnoredInventory) Summary() string {
	count := fmt.Sprintf("%d", inventory.Count)
	if inventory.AtLeast {
		count = fmt.Sprintf("at-least=%d", ignoredEntryLimit+1)
	}
	return fmt.Sprintf("count=%s bytes=%d shown=%d truncated=%t", count, inventory.Bytes, inventory.Shown, inventory.Truncated)
}

// orphaned reports whether an assignment has been abandoned by the session that cut it.
// Age is the whole test: nothing records liveness for a request-created worktree, so
// bounds.AssignmentStale is the only thing separating a long-running one from residue,
// and every consumer must ask this one question so the window has a single meaning.
//
// An absent stamp is aged, because a record written before the field existed carries
// none and would otherwise be immortal. A stamp the reading host's clock has not
// reached yet is not aged, so skew cannot manufacture an orphan. An unparseable stamp
// is unknown age rather than infinite age; ValidateAssignment rejects one on every
// ledger read, so a record reaching here with one never came from the ledger.
func orphaned(a intent.Assignment, now time.Time) bool {
	if a.State != intent.StateActive {
		return false
	}
	if a.CreatedAt == nil {
		return true
	}
	created, err := time.Parse(time.RFC3339, *a.CreatedAt)
	if err != nil {
		return false
	}
	return now.Sub(created) > bounds.AssignmentStale
}

func PlanAutomatic(root, path string) (CleanupPlan, error) {
	plan, err := PlanExplicit(root, path)
	if err != nil {
		reason := ReasonUncertain
		if strings.Contains(err.Error(), "assignment") || strings.Contains(err.Error(), "intent ledger") {
			reason = ReasonMalformed
		}
		return retainedPlan(path, reason, err.Error()), nil
	}
	if plan.Action == ActionRetain {
		return automaticFingerprint(plan), nil
	}
	if plan.assignment == nil || !plan.owned {
		return automaticRetain(plan, ReasonForeign, "registration is not a verified owned assignment"), nil
	}
	plan.Assignment = plan.assignment.ID
	if plan.assignment.State != intent.StateCleanupPending {
		reason := ReasonUncertain
		if plan.assignment.State == intent.StateActive {
			reason = ReasonActive
			if orphaned(*plan.assignment, time.Now()) {
				reason = ReasonOrphaned
			}
		}
		return automaticRetain(plan, reason, "assignment is not cleanup-pending"), nil
	}
	if !recoveryMetadataMatches(root, *plan.assignment) {
		return automaticRetain(plan, ReasonMalformed, "assignment recovery metadata does not match refs"), nil
	}
	if strings.HasPrefix(plan.landed, "unknown") {
		return automaticRetain(plan, ReasonUncertain, "assignment landedness is unknown"), nil
	}
	if !strings.HasPrefix(plan.landed, "true:") {
		return automaticRetain(plan, ReasonUnmerged, "assignment branch has not landed"), nil
	}
	return automaticFingerprint(plan), nil
}
func automaticRetain(plan CleanupPlan, reason CleanupReason, detail string) CleanupPlan {
	plan.Action, plan.ReasonCode, plan.Reason = ActionRetain, reason, detail
	return automaticFingerprint(plan)
}
func automaticFingerprint(plan CleanupPlan) CleanupPlan {
	if plan.assignment != nil {
		plan.Fingerprint = fingerprintParts([]byte("bench-automatic-registration/v1"), []byte(plan.Target), []byte(plan.assignment.OwnerID), []byte(plan.assignment.ID), []byte(plan.assignment.Request))
		return plan
	}
	plan.Fingerprint = fingerprintParts([]byte("bench-automatic/v1"), []byte(plan.Fingerprint), []byte(plan.Action), []byte(plan.ReasonCode), []byte(plan.Reason))
	return plan
}
func foreignRecoveryAssignment(root, target string) *intent.Assignment {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return nil
	}
	var found *intent.Assignment
	for i := range assignments {
		candidate := assignments[i]
		if candidate.Worktree == target && candidate.Label == "foreign exact cleanup" && candidate.Request == requestDigest("foreign:"+target) && len(candidate.Recovery) > 0 {
			if found != nil {
				return nil
			}
			found = &candidate
		}
	}
	return found
}
func retainedPlan(target string, reason CleanupReason, detail string) CleanupPlan {
	return CleanupPlan{Target: target, Action: ActionRetain, ReasonCode: reason, Reason: detail}
}
func recoveryMetadataMatches(root string, assignment intent.Assignment) bool {
	prefix := intent.RecoveryRefPrefix(assignment.OwnerID, assignment.ID)
	out, err := git.Output("-C", root, "for-each-ref", "--format=%(refname)", prefix)
	if err != nil {
		return false
	}
	actual := map[string]bool{}
	if out != "" {
		for _, ref := range strings.Split(out, "\n") {
			actual[ref] = true
		}
	}
	expected := map[string]intent.Recovery{}
	for _, recovery := range assignment.Recovery {
		if _, duplicate := expected[recovery.Ref]; duplicate {
			return false
		}
		expected[recovery.Ref] = recovery
	}
	for ref := range actual {
		if _, ok := expected[ref]; !ok {
			return false
		}
	}
	for ref, recovery := range expected {
		if !actual[ref] {
			if !recoveryEnvelopeValid(root, recovery) {
				return false
			}
			continue
		}
		resolved, resolveErr := git.Output("-C", root, "show-ref", "--verify", "--hash", ref)
		if resolveErr != nil || resolved != recovery.Root {
			return false
		}
		for _, payload := range recovery.Payloads {
			if !git.OK("-C", root, "cat-file", "-e", payload+"^{commit}") || !git.OK("-C", root, "merge-base", "--is-ancestor", payload, recovery.Root) {
				return false
			}
		}
	}
	return true
}
func PlanRecovery(root, ref string) (RecoveryPlan, error) {
	root = canonicalRoot(root)
	ledgerBytes, err := intent.LifecycleEvidence(root)
	if err != nil {
		return RecoveryPlan{}, err
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return RecoveryPlan{}, err
	}
	plan := RecoveryPlan{Ref: ref, Action: RecoveryRetain, Root: "none", Payloads: "none", Landed: "unknown", Changes: changeSummary(root, ref)}
	for i := range assignments {
		for _, recovery := range assignments[i].Recovery {
			if recovery.Ref != ref {
				continue
			}
			if plan.assignment != nil {
				plan.Detail = "recovery ref belongs to ambiguous assignments"
				return fingerprintRecovery(root, ledgerBytes, plan), nil
			}
			assignment := assignments[i]
			plan.assignment, plan.recovery = &assignment, recovery
			plan.Root, plan.Payloads = recovery.Root, strings.Join(recovery.Payloads, ";")
		}
	}
	if plan.assignment == nil {
		// This is the one path that never reaches verifyRecovery's envelope check, so the
		// namespace is checked here or nowhere: without it an ordinary branch resolves, reads
		// as an orphan, and a discard deletes it. Existence is not consulted first, because a
		// name outside the namespace carries no authorization either way.
		if !strings.HasPrefix(ref, recoveryNamespace()) {
			plan.Action, plan.Detail = RecoveryForeign, "ref is outside the recovery namespace"
			return fingerprintRecovery(root, ledgerBytes, plan), nil
		}
		// With no row to consult, the ref itself is the only evidence separating work still
		// preserved from work already retired.
		plan.Action, plan.Detail = RecoveryAbsent, "recovery ref does not exist"
		if refExists(root, ref) {
			plan.Action, plan.Detail = RecoveryOrphaned, "recovery ref has no owning assignment"
		}
		return fingerprintRecovery(root, ledgerBytes, plan), nil
	}
	if plan.assignment.State != intent.StateRecovered {
		plan.Detail = "recovery ref has no recovered assignment"
		return fingerprintRecovery(root, ledgerBytes, plan), nil
	}
	if err := verifyRecovery(root, *plan.assignment, plan.recovery); err != nil {
		plan.Detail = err.Error()
		return fingerprintRecovery(root, ledgerBytes, plan), nil
	}
	def, ok := git.ResolvedDefault(root)
	if !ok {
		plan.Detail = "default branch does not resolve"
		return fingerprintRecovery(root, ledgerBytes, plan), nil
	}
	verdicts := make([]string, 0, len(plan.recovery.Payloads))
	all := len(plan.recovery.Payloads) > 0
	for _, payload := range plan.recovery.Payloads {
		landed, byContent, landedErr := git.LandedInDefault(root, payload, def)
		verdict := "unlanded"
		if landedErr != nil {
			verdict = "unknown"
		} else if landed && byContent {
			verdict = "patch"
		} else if landed {
			verdict = "ancestor"
		}
		verdicts = append(verdicts, payload+"="+verdict)
		all = all && landedErr == nil && landed
	}
	plan.Landed = strings.Join(verdicts, ";")
	if all {
		plan.Action, plan.Detail = RecoveryRetire, "apply with exact fingerprint"
	} else {
		// The one assignment of the discard-eligible verdict: every earlier return leaves
		// the retain initialiser in place, so reaching this line is the whole proof chain —
		// envelope verified, default branch resolved, payloads judged and refused.
		plan.Action, plan.Detail = RecoveryDiscard, "discard with exact fingerprint"
	}
	return fingerprintRecovery(root, ledgerBytes, plan), nil
}

// recoveryNamespace is the ref namespace preserved work lives under, read from the one
// definition in internal/intent rather than restated here: a reader that named the path
// for itself could disagree with the writer that puts refs there. The probe identity is
// arbitrary — only the text ahead of it is the namespace.
func recoveryNamespace() string {
	const probe = "namespace-probe"
	prefix := intent.RecoveryRefPrefix(probe, probe)
	return prefix[:strings.Index(prefix, probe)]
}

// refExists reports whether ref names a ref this repository resolves. Only a fully
// qualified name answers true, so a payload OID or a short name is never mistaken for a
// surviving recovery ref.
func refExists(root, ref string) bool {
	return git.OK("-C", root, "show-ref", "--verify", "--quiet", ref)
}

// changeSummary counts the distinct paths the envelope at commitish changes against the
// base that envelope records. Every layer contributes to one count: the operator is
// deciding about the ref as a whole, and a per-layer breakdown would make a two-file
// leftover look larger than a one-file one.
func changeSummary(root, commitish string) string {
	manifest, ok := readRecoveryManifest(root, commitish)
	if !ok {
		return recoveryUnknownChanges
	}
	paths := map[string]bool{}
	for _, payload := range manifest.Layers {
		out, err := git.Output("-C", root, "diff", "--name-only", manifest.Base, payload)
		if err != nil {
			return recoveryUnknownChanges
		}
		for _, path := range strings.Split(out, "\n") {
			if path != "" {
				paths[path] = true
			}
		}
	}
	return fmt.Sprintf("paths=%d", len(paths))
}

// recoveryFingerprintDomain and recoveryFingerprintEffects name the widest authority a
// recovery fingerprint carries. Both name the destructive discard, because the same
// fingerprint now authorizes dropping work the landedness proof never accepted, and a
// value planned under the older retire-only authority must not be able to authorize that.
const (
	recoveryFingerprintDomain  = "bench-recovery-retire-or-discard/v1"
	recoveryFingerprintEffects = "delete-exact-ref,discard-unproven-payload,update-assignment,compact-if-last"
)

func fingerprintRecovery(root string, ledger []byte, plan RecoveryPlan) RecoveryPlan {
	return fingerprintRecoveryUnder(root, ledger, plan, recoveryFingerprintDomain, recoveryFingerprintEffects)
}

// fingerprintRecoveryUnder derives a plan's fingerprint under a named authority. The
// authority is a parameter rather than a literal so the separation is observable: asking
// what the same plan fingerprints to under the retire-only authority is the only way to
// see that the two are different values, which is the whole of the guarantee.
func fingerprintRecoveryUnder(root string, ledger []byte, plan RecoveryPlan, domain, effects string) RecoveryPlan {
	common, _ := git.Output("-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	def, _ := git.ResolvedDefault(root)
	defaultOID, _ := git.Output("-C", root, "rev-parse", "--verify", def+"^{commit}")
	refOID, _ := git.Output("-C", root, "rev-parse", "--verify", plan.Ref+"^{commit}")
	// Every fact the plan reported is sealed, the derived change summary among them: it is
	// what the operator judged the discard by, so a plan whose count has moved since must
	// not still authorize one.
	parts := [][]byte{
		[]byte(domain), []byte(common), []byte(def), []byte(defaultOID), ledger,
		[]byte(plan.Ref), []byte(refOID), []byte(plan.Root), []byte(plan.Payloads), []byte(plan.Landed), []byte(plan.Changes), []byte(plan.Action), []byte(plan.Detail),
		[]byte(effects),
	}
	if plan.assignment != nil {
		for _, recovery := range plan.assignment.Recovery {
			actual, _ := git.Output("-C", root, "rev-parse", "--verify", recovery.Ref+"^{commit}")
			parts = append(parts, []byte(recovery.Ref), []byte(recovery.Root), []byte(strings.Join(recovery.Payloads, ";")), []byte(actual))
			for _, payload := range recovery.Payloads {
				exists := git.OK("-C", root, "cat-file", "-e", payload+"^{commit}")
				reachable := exists && git.OK("-C", root, "merge-base", "--is-ancestor", payload, recovery.Root)
				parts = append(parts, []byte(payload), []byte(fmt.Sprintf("exists=%t reachable=%t", exists, reachable)))
			}
		}
	}
	plan.Fingerprint = fingerprintParts(parts...)
	return plan
}

type ownerEvidence struct {
	path         string
	admin        string
	marker       Marker
	registration git.Worktree
}

func validateOwnerMarker(root, path string) (ownerEvidence, error) {
	root = canonicalRoot(root)
	target, err := canonicalPath(path)
	if err != nil {
		return ownerEvidence{}, fmt.Errorf("canonical target: %w", err)
	}
	if samePath(root, target) {
		return ownerEvidence{}, errors.New("primary checkout is never cleanup-eligible")
	}
	registrations, err := git.Worktrees(root)
	if err != nil {
		return ownerEvidence{}, fmt.Errorf("read worktree registrations: %w", err)
	}
	var registration *git.Worktree
	for i := range registrations {
		registeredPath, pathErr := canonicalPath(registrations[i].Path)
		if pathErr == nil && registeredPath == target {
			if registration != nil {
				return ownerEvidence{}, errors.New("target has ambiguous registrations")
			}
			registration = &registrations[i]
		}
	}
	if registration == nil {
		return ownerEvidence{}, errors.New("target is not the current registered worktree")
	}
	admin, err := git.Output("-C", target, "rev-parse", "--path-format=absolute", "--git-dir")
	if err != nil || !filepath.IsAbs(admin) {
		return ownerEvidence{}, errors.New("private worktree administration directory is unavailable")
	}
	admin = filepath.Clean(admin)
	common, err := git.Output("-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil || !filepath.IsAbs(common) {
		return ownerEvidence{}, errors.New("common Git directory is unavailable")
	}
	privateRoot := filepath.Join(filepath.Clean(common), "worktrees")
	if !insidePool(privateRoot, admin) {
		return ownerEvidence{}, errors.New("marker is not in a private linked-worktree administration directory")
	}
	markerFile := filepath.Join(admin, OwnerMarkerFile)
	info, err := os.Lstat(markerFile)
	if err != nil {
		return ownerEvidence{}, errors.New("owner marker is absent")
	}
	if info.Mode().Type() != 0 {
		return ownerEvidence{}, errors.New("owner marker is not a regular file")
	}
	if info.Mode().Perm() != 0o600 {
		return ownerEvidence{}, errors.New("owner marker permissions are not 0600")
	}
	data, err := os.ReadFile(markerFile)
	if err != nil || len(data) == 0 {
		return ownerEvidence{}, errors.New("owner marker is empty or unreadable")
	}
	marker, err := decodeMarker(data)
	if err != nil {
		return ownerEvidence{}, fmt.Errorf("owner marker is malformed: %w", err)
	}
	if marker.Schema != OwnerMarkerSchema {
		return ownerEvidence{}, errors.New("owner marker schema is unsupported")
	}
	if !intent.ValidIdentity(marker.OwnerID) {
		return ownerEvidence{}, errors.New("owner marker ID is invalid")
	}
	if marker.Path != target || !filepath.IsAbs(marker.Path) || filepath.Clean(marker.Path) != marker.Path {
		return ownerEvidence{}, errors.New("owner marker path does not match the registration")
	}
	return ownerEvidence{path: target, admin: admin, marker: marker, registration: *registration}, nil
}
