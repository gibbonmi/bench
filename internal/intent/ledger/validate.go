package ledger

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	idPattern     = regexp.MustCompile(`^[0-9a-f]{32}$`)
	digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
	oidPattern    = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
)

// AssignmentBranchPrefix names the local-ref namespace reserved for assignment branches.
func AssignmentBranchPrefix() string { return assignmentBranchNamespace }

const shiftBranchNamespace = "refs/heads/bench/shift-"

// ShiftBranchPrefix names the local-ref namespace reserved for shift branches. The shift
// loop derives its branch name from it, and the unclaimed cleanup selects on it.
func ShiftBranchPrefix() string { return shiftBranchNamespace }

func AssignmentBranchRef(ownerID, assignmentID string) string {
	return assignmentBranchNamespace + ownerID + "/" + assignmentID
}

func RecoveryRefPrefix(ownerID, assignmentID string) string {
	return RecoveryRefNamespace + ownerID + "/" + assignmentID + "/"
}

// ValidAssignmentBranchRef reports whether ref sits in the assignment namespace. The
// name is exported because the intent package is a sibling and cannot reach an
// unexported declaration here.
func ValidAssignmentBranchRef(ref string) bool {
	return strings.HasPrefix(ref, assignmentBranchNamespace)
}

// RequestDigest derives the persisted identity for an opaque caller request.
func RequestDigest(value string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(value))) }

func ValidIdentity(value string) bool { return idPattern.MatchString(value) }

// ValidateEntry grades one writer entry. The name is exported because the intent
// package is a sibling and cannot reach an unexported declaration here.
func ValidateEntry(entry Entry) error {
	if entry.Key == "" || entry.CreatedAt.IsZero() {
		return errors.New("entry requires key and creation time")
	}
	switch entry.Kind {
	case KindShift, KindWorktree, KindClaudeAgent:
		return nil
	default:
		return fmt.Errorf("entry %q has unknown writer kind %q", entry.Key, entry.Kind)
	}
}

func ValidateAssignment(a Assignment) error {
	if a.Schema != AssignmentRecordSchema {
		return fmt.Errorf("assignment %q has unsupported schema %q", a.ID, a.Schema)
	}
	if !ValidIdentity(a.ID) || !ValidIdentity(a.OwnerID) {
		return errors.New("assignment has invalid owner or assignment ID")
	}
	if !digestPattern.MatchString(a.Request) || a.Label == "" || !oidPattern.MatchString(a.Start) {
		return fmt.Errorf("assignment %q has invalid request, label, or start", a.ID)
	}
	wantBranch := AssignmentBranchRef(a.OwnerID, a.ID)
	if a.Branch != wantBranch {
		return fmt.Errorf("assignment %q has non-canonical branch", a.ID)
	}
	if !filepath.IsAbs(a.Worktree) || filepath.Clean(a.Worktree) != a.Worktree {
		return fmt.Errorf("assignment %q has non-canonical worktree", a.ID)
	}
	// A stamp ahead of the reading host's clock stays valid here: skew is the age
	// predicate's to interpret. Rejecting it would make one skewed write unreadable
	// to every command, since this runs on every ledger read.
	if a.CreatedAt != nil {
		if _, err := time.Parse(time.RFC3339, *a.CreatedAt); err != nil {
			return fmt.Errorf("assignment %q has unparseable created_at", a.ID)
		}
	}
	switch a.State {
	case StateActive, StateCleanupPending, StateComplete:
	case StateRecovered:
		if len(a.Recovery) == 0 {
			return fmt.Errorf("assignment %q is recovered without recovery metadata", a.ID)
		}
	default:
		return fmt.Errorf("assignment %q has unknown state %q", a.ID, a.State)
	}
	if a.State == StateActive && len(a.Recovery) != 0 {
		return fmt.Errorf("assignment %q is active with recovery metadata", a.ID)
	}
	if a.State == StateComplete && len(a.Recovery) != 0 {
		return fmt.Errorf("assignment %q is complete with recovery metadata", a.ID)
	}
	for _, recovery := range a.Recovery {
		prefix := RecoveryRefPrefix(a.OwnerID, a.ID)
		if !strings.HasPrefix(recovery.Ref, prefix) || !oidPattern.MatchString(recovery.Root) || len(recovery.Payloads) == 0 {
			return fmt.Errorf("assignment %q has invalid recovery metadata", a.ID)
		}
		for _, payload := range recovery.Payloads {
			if !oidPattern.MatchString(payload) {
				return fmt.Errorf("assignment %q has invalid recovery payload", a.ID)
			}
		}
	}
	return nil
}

// ValidateCleanupReceipts grades the receipt set. The name is exported because the
// intent package is a sibling and cannot reach an unexported declaration here.
func ValidateCleanupReceipts(receipts []CleanupReceipt) error {
	seen := map[string]bool{}
	for _, receipt := range receipts {
		if receipt.Schema != CleanupReceiptSchema || !filepath.IsAbs(receipt.Repo) || filepath.Clean(receipt.Repo) != receipt.Repo || receipt.Operation == "" || !filepath.IsAbs(receipt.Target) || filepath.Clean(receipt.Target) != receipt.Target || !digestPattern.MatchString(receipt.Fingerprint) {
			return errors.New("cleanup receipt has invalid identity")
		}
		if receipt.State != ReceiptInFlight && receipt.State != ReceiptComplete {
			return errors.New("cleanup receipt has invalid state")
		}
		switch receipt.Phase {
		case ReceiptPhasePlanned, ReceiptPhasePreserved, ReceiptPhaseRemoving, ReceiptPhaseRemoved, ReceiptPhaseBranch, ReceiptPhaseTerminal:
		default:
			return errors.New("cleanup receipt has invalid phase")
		}
		if receipt.Checkpoint != "" && !digestPattern.MatchString(receipt.Checkpoint) {
			return errors.New("cleanup receipt has invalid checkpoint")
		}
		if (receipt.Branch == "") != (receipt.BranchOID == "") || receipt.Branch != "" && (!receipt.Owned || !ValidAssignmentBranchRef(receipt.Branch) || !oidPattern.MatchString(receipt.BranchOID)) {
			return errors.New("cleanup receipt has invalid branch CAS")
		}
		if (receipt.Owner == "") != (receipt.Assignment == "") || (receipt.Assignment == "") != (receipt.Request == "") || receipt.Assignment != "" && (!receipt.Owned || !ValidIdentity(receipt.Owner) || !ValidIdentity(receipt.Assignment) || !digestPattern.MatchString(receipt.Request) || receipt.Branch != "" && receipt.Branch != AssignmentBranchRef(receipt.Owner, receipt.Assignment)) {
			return errors.New("cleanup receipt has invalid owned assignment")
		}
		if receipt.State == ReceiptComplete && receipt.Phase != ReceiptPhaseTerminal {
			return errors.New("completed cleanup receipt is not terminal")
		}
		key := receipt.Repo + "\x00" + receipt.Operation + "\x00" + receipt.Target + "\x00" + receipt.Fingerprint
		if seen[key] {
			return errors.New("cleanup receipt has duplicate identity")
		}
		seen[key] = true
	}
	return nil
}
