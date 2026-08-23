package worktree

import (
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CleanupAction and CleanupReason are the policy package's action and reason
// types; internal/worktree/lifecyclepolicy owns their values and semantics.
type CleanupAction = lifecyclepolicy.Action

type refusal struct {
	detail, observed, wanted, next string
	paths                          []string
}
type refusalError struct{ refusal }

func (r refusal) fields() string {
	fields := []string{"detail=" + sanitize.Controls(r.detail)}
	for _, pair := range [][2]string{{"observed", r.observed}, {"wanted", r.wanted}, {"next", r.next}} {
		if pair[1] != "" {
			fields = append(fields, pair[0]+"="+sanitize.Controls(pair[1]))
		}
	}
	return strings.Join(fields, ",")
}
func (e refusalError) Error() string {
	text := sanitize.Controls(e.detail)
	if fields := strings.TrimPrefix(e.fields(), "detail="+sanitize.Controls(e.detail)); fields != "" {
		text += "; " + strings.TrimPrefix(fields, ",")
	}
	if table := e.table(); table != "" {
		text += "\n" + table
	}
	return text
}
func (r refusal) table() string {
	if len(r.paths) == 0 {
		return ""
	}
	shown := len(r.paths)
	if shown > ignoredEntryLimit {
		shown = ignoredEntryLimit
	}
	rows := make([][]string, 0, shown)
	for _, path := range r.paths[:shown] {
		rows = append(rows, []string{sanitize.Controls(path)})
	}
	out, err := toon.Table("refusal_paths", []string{"path"}, rows)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("paths_total=%d\n%s", len(r.paths), out)
}

func renderCleanup(stdout io.Writer, plan CleanupPlan) error {
	return renderCleanups(stdout, []CleanupPlan{plan})
}

var cleanupFields = []string{"target", "action", "tracked", "ignored", "recovery", "fingerprint", "detail"}

func renderCleanups(stdout io.Writer, plans []CleanupPlan) error {
	rows := make([][]string, 0, len(plans))
	for _, plan := range plans {
		rows = append(rows, cleanupRow(plan))
	}
	out, err := toon.Table("worktree_cleanup", cleanupFields, rows)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprint(stdout, out); err != nil {
		return err
	}
	for _, plan := range plans {
		if err := renderIgnoredPreview(stdout, plan); err != nil {
			return err
		}
	}
	return nil
}

func cleanupRow(plan CleanupPlan) []string {
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
	// nothing else in the row implies. The plan half names the exact ref it spends that
	// assertion on. The derived case stays silent: its deletion follows from the landedness
	// the tool proved for itself, and its row is a settled output contract.
	if plan.discardBranch && plan.deleteBranch && plan.Action.Removes() {
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
	return []string{target, string(plan.Action), tracked, ignored, recovery, plan.Fingerprint, detail}
}

func renderIgnoredPreview(stdout io.Writer, plan CleanupPlan) error {
	if plan.Ignored.Shown == 0 {
		return nil
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

type CleanupReason = lifecyclepolicy.Reason

const (
	ActionRetain         = lifecyclepolicy.ActionRetain
	ActionRemove         = lifecyclepolicy.ActionRemove
	ActionRecoverRemove  = lifecyclepolicy.ActionRecoverRemove
	ActionDiscardRemove  = lifecyclepolicy.ActionDiscardRemove
	ActionRemoved        = lifecyclepolicy.ActionRemoved
	ActionError          = lifecyclepolicy.ActionError
	ReasonForeign        = lifecyclepolicy.ReasonForeign
	ReasonActive         = lifecyclepolicy.ReasonActive
	ReasonLiveLease      = lifecyclepolicy.ReasonLiveLease
	ReasonUnmerged       = lifecyclepolicy.ReasonUnmerged
	ReasonIgnored        = lifecyclepolicy.ReasonIgnored
	ReasonMalformed      = lifecyclepolicy.ReasonMalformed
	ReasonUncertain      = lifecyclepolicy.ReasonUncertain
	ReasonUnexpectedLock = lifecyclepolicy.ReasonUnexpectedLock
	ReasonOrphaned       = lifecyclepolicy.ReasonOrphaned
	ReasonDirty          = lifecyclepolicy.ReasonDirty
	ReasonLanded         = lifecyclepolicy.ReasonLanded
)
const actionReleaseRemove = lifecyclepolicy.ActionReleaseRemove
const actionReleaseLeftover = lifecyclepolicy.ActionReleaseLeftover

// preserves is the policy Preserves predicate read through this plan's own
// action, tracked state, and registration shape.
func (plan CleanupPlan) preserves() bool {
	return lifecyclepolicy.Preserves(plan.Action, plan.Tracked, plan.registration.Detached)
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
	// regular file (a gitfile worktree pointer) nor a directory (an embedded repository).
	// The entry is a FIFO, device, socket, or a symlinked .git — something git itself
	// refuses to treat as ordinary metadata. No consumer may invoke git against this path:
	// a FIFO with no writer would block the invocation forever. This shape fails closed
	// rather than joining ShapeCheckoutDirectory or the decayed set.
	ShapeSpecialMetadata PathShape = "special-git-metadata"
	ShapeUnknown         PathShape = "unknown"
)

// ClassifyPathShape is the single source for the decayed-shape policy. Every consumer
// asking whether an assignment's checkout is still live, or whether an abandon releases
// residue rather than removing a checkout, decides it here. The two can never answer
// differently for the same bytes.
//
// The path is never opened: a FIFO at an assignment path has no writer and would block a
// reader forever. The verdict rests on lstat and stat shape alone. Only
// ShapeCheckoutDirectory licenses a caller to run git against the path. Symlinks are
// followed, because a resolvable one is already resolved away by the time a canonical
// target names it. One surviving here resolves to nothing.
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
	// landedTyped carries the same landedness evidence as landed, but as the typed value
	// decideExplicit produced. A production reader can branch on its kind and proof fields
	// rather than parse the wire string. landed stays alongside it, since subshell.go's
	// fingerprint hashes the string as evidence. A hand-built stub (clean_landed.go) may
	// populate only the typed field when it need not fabricate the string.
	landedTyped landedness
	// leftover names the present bytes a release-leftover plan hands on rather than
	// removes. It is empty for every plan that answers for a checkout.
	leftover string
	// unresolved marks a plan whose operand named nothing this repository can act on,
	// as opposed to a checkout it resolved and then declined to remove. Only the
	// path-addressed command reads it, to keep a destructive call from reporting success
	// it never earned. The automatic sweep addresses registrations it already holds and
	// never produces one.
	unresolved bool
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

func (inventory IgnoredInventory) Summary() string {
	count := fmt.Sprintf("%d", inventory.Count)
	if inventory.AtLeast {
		count = fmt.Sprintf("at-least=%d", ignoredEntryLimit+1)
	}
	return fmt.Sprintf("count=%s bytes=%d shown=%d truncated=%t", count, inventory.Bytes, inventory.Shown, inventory.Truncated)
}

// orphaned is the policy age decision over the bounds.AssignmentStale window
// this boundary supplies; internal/worktree/lifecyclepolicy owns its edge
// readings.
func orphaned(a intent.Assignment, now time.Time) bool {
	return lifecyclepolicy.Orphaned(a, now, bounds.AssignmentStale)
}

// PlanAutomatic decides the automatic, unattended reading of eligibility. It calls
// PlanExplicit for its own result, then gathers every automatic-specific fact
// decideAutomatic needs. A fact stays ungathered where an earlier one already made it
// inapplicable. It calls decideAutomatic exactly once, then projects the returned
// verdict onto the plan.
func PlanAutomatic(root, path string) (CleanupPlan, error) {
	return planAutomaticAt(root, path, currentTime())
}

// planAutomaticAt is PlanAutomatic with the instant resolved explicitly at the
// caller's effect boundary; PlanAutomatic is its temporary compatibility form.
func planAutomaticAt(root, path string, now time.Time) (CleanupPlan, error) {
	explicitPlan, explicitErr := PlanExplicit(root, path)
	facts := automaticFacts{ExplicitErr: explicitErr, Explicit: explicitOutcome(explicitPlan)}

	var missingBranchAssignment *intent.Assignment
	if explicitErr != nil {
		if assignment, missing := activeAssignmentWithMissingBranch(root, path); missing {
			salvage := retainedPlan(path, ReasonActive, "assignment landedness is unknown")
			salvage.owned, salvage.assignment = true, &assignment
			missingBranchAssignment = &assignment
			facts.MissingBranchAssignmentID = assignment.ID
			facts.MissingBranchLiveLease = planHasLiveLease(salvage)
		}
	} else if explicitPlan.assignment != nil {
		facts.LiveLease = planHasLiveLease(explicitPlan)
		facts.Landed = assignmentLanded(*explicitPlan.assignment, explicitPlan)
		if explicitPlan.assignment.State == intent.StateActive {
			facts.OrphanedActive = orphaned(*explicitPlan.assignment, now)
		}
		if explicitPlan.owned && explicitPlan.assignment.State == intent.StateCleanupPending {
			facts.RecoveryMatches = recoveryMetadataMatches(root, *explicitPlan.assignment)
		}
	}

	verdict := decideAutomatic(facts)

	if explicitErr != nil {
		if missingBranchAssignment != nil {
			plan := retainedPlan(path, verdict.ReasonCode, verdict.Reason)
			plan.Assignment = verdict.AssignmentID
			plan.owned, plan.assignment = true, missingBranchAssignment
			return automaticFingerprint(plan), nil
		}
		return retainedPlan(path, verdict.ReasonCode, verdict.Reason), nil
	}

	plan := explicitPlan
	plan.Action, plan.ReasonCode, plan.Reason = verdict.Action, verdict.ReasonCode, verdict.Reason
	if verdict.AssignmentID != "" {
		plan.Assignment = verdict.AssignmentID
	}
	return automaticFingerprint(plan), nil
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
		if candidate.Worktree == target && candidate.Label == "foreign exact cleanup" && candidate.Request == intent.RequestDigest("foreign:"+target) && len(candidate.Recovery) > 0 {
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

// unresolvedPlan is retainedPlan for an operand this repository could not resolve to one
// of its own checkouts.
func unresolvedPlan(target string, reason CleanupReason, detail string) CleanupPlan {
	plan := retainedPlan(target, reason, detail)
	plan.unresolved = true
	return plan
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
	common, err := git.CommonDir(root)
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
