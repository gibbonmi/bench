package worktree

import (
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
	"io"
	"os"
	"path/filepath"
	"strings"
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
)

var ignoredLstat = os.Lstat

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
	deleteBranch         bool
	branchRef, branchOID string
	ignoredSummary       string
	landed               string
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
)

type RecoveryPlan struct {
	Ref, Root, Payloads, Landed string
	Action                      RecoveryAction
	Fingerprint, Detail         string
	assignment                  *intent.Assignment
	recovery                    intent.Recovery
}

func (inventory IgnoredInventory) Summary() string {
	count := fmt.Sprintf("%d", inventory.Count)
	if inventory.AtLeast {
		count = fmt.Sprintf("at-least=%d", ignoredEntryLimit+1)
	}
	return fmt.Sprintf("count=%s bytes=%d shown=%d truncated=%t", count, inventory.Bytes, inventory.Shown, inventory.Truncated)
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
	plan := RecoveryPlan{Ref: ref, Action: RecoveryRetain, Root: "none", Payloads: "none", Landed: "unknown"}
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
	if plan.assignment == nil || plan.assignment.State != intent.StateRecovered {
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
		plan.Detail = "every recovery payload must be proven landed"
	}
	return fingerprintRecovery(root, ledgerBytes, plan), nil
}
func fingerprintRecovery(root string, ledger []byte, plan RecoveryPlan) RecoveryPlan {
	common, _ := git.Output("-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	def, _ := git.ResolvedDefault(root)
	defaultOID, _ := git.Output("-C", root, "rev-parse", "--verify", def+"^{commit}")
	refOID, _ := git.Output("-C", root, "rev-parse", "--verify", plan.Ref+"^{commit}")
	parts := [][]byte{
		[]byte("bench-recovery-retire/v1"), []byte(common), []byte(def), []byte(defaultOID), ledger,
		[]byte(plan.Ref), []byte(refOID), []byte(plan.Root), []byte(plan.Payloads), []byte(plan.Landed), []byte(plan.Action), []byte(plan.Detail),
		[]byte("delete-exact-ref,update-assignment,compact-if-last"),
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
