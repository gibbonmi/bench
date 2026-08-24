package intent

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const AssignmentRecordSchema = "bench-assignment/v1"

const assignmentBranchNamespace = "refs/heads/bench/assign/"

// RecoveryRefNamespace is the one name for the namespace preserved work lives under.
// Both the writer that puts refs there and the standing cleaner that sweeps it read
// this constant. Neither can address a namespace the other does not.
const RecoveryRefNamespace = "refs/bench/recovery/"

func AssignmentBranchRef(ownerID, assignmentID string) string {
	return assignmentBranchNamespace + ownerID + "/" + assignmentID
}

func RecoveryRefPrefix(ownerID, assignmentID string) string {
	return RecoveryRefNamespace + ownerID + "/" + assignmentID + "/"
}

func validAssignmentBranchRef(ref string) bool {
	return strings.HasPrefix(ref, assignmentBranchNamespace)
}

type AssignmentState string

const (
	StateActive         AssignmentState = "active"
	StateCleanupPending AssignmentState = "cleanup-pending"
	StateRecovered      AssignmentState = "recovered"
	StateComplete       AssignmentState = "complete"
)

type Recovery struct {
	Ref      string   `json:"ref"`
	Root     string   `json:"root"`
	Payloads []string `json:"payloads"`
}

// Assignment is the persisted half of a Bench-owned registration. The worktree
// marker proves immutable owner identity. This record binds that owner to exactly
// one caller request, branch, start commit, path, lifecycle state, and recovery set.
type Assignment struct {
	Schema  string `json:"schema"`
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Request string `json:"request"`
	// RequestToken is the plain caller token the digest above derives from. The digest
	// stays the authorization identity; the token is persisted so `bench worktree list`
	// can hand a resumed landing the exact value to pass. Records written before the
	// field existed carry none and serialize without the key.
	RequestToken string          `json:"request_token,omitempty"`
	Label        string          `json:"label"`
	Start        string          `json:"start"`
	Branch       string          `json:"branch"`
	Worktree     string          `json:"worktree"`
	State        AssignmentState `json:"state"`
	Recovery     []Recovery      `json:"recovery"`
	// CreatedAt is an RFC3339 creation time. A nil stamp is absence, which stays
	// valid because records written before the field existed carry none and
	// serialize without the key. A present stamp must parse, so the pointer is
	// what keeps a hand-written empty string distinguishable from absence.
	CreatedAt *string `json:"created_at,omitempty"`
}

const (
	CleanupReceiptSchema  = "bench-cleanup-receipt/v1"
	ReceiptInFlight       = "in-flight"
	ReceiptComplete       = "complete"
	MaxCleanupReceipts    = 256
	ReceiptPhasePlanned   = "planned"
	ReceiptPhasePreserved = "preserved"
	ReceiptPhaseRemoving  = "removing"
	ReceiptPhaseRemoved   = "removed"
	ReceiptPhaseBranch    = "branch-removed"
	ReceiptPhaseTerminal  = "terminal"
)

type CleanupReceipt struct {
	Schema      string `json:"schema"`
	Repo        string `json:"repo"`
	Operation   string `json:"operation"`
	Target      string `json:"target"`
	Fingerprint string `json:"fingerprint"`
	State       string `json:"state"`
	Phase       string `json:"phase"`
	Checkpoint  string `json:"checkpoint,omitempty"`
	Action      string `json:"action"`
	Tracked     string `json:"tracked"`
	Ignored     string `json:"ignored"`
	Recovery    string `json:"recovery"`
	Detail      string `json:"detail"`
	Owned       bool   `json:"owned,omitempty"`
	Branch      string `json:"branch,omitempty"`
	BranchOID   string `json:"branch_oid,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Assignment  string `json:"assignment,omitempty"`
	Request     string `json:"request,omitempty"`
}

func validateCleanupReceipts(receipts []CleanupReceipt) error {
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
		if (receipt.Branch == "") != (receipt.BranchOID == "") || receipt.Branch != "" && (!receipt.Owned || !validAssignmentBranchRef(receipt.Branch) || !oidPattern.MatchString(receipt.BranchOID)) {
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

func CleanupReceiptFor(root, repo, operation, target, fingerprint string) (CleanupReceipt, bool, error) {
	ledger, err := Read(root)
	if err != nil {
		return CleanupReceipt{}, false, err
	}
	for _, receipt := range ledger.CleanupReceipts {
		if receipt.Repo == repo && receipt.Operation == operation && receipt.Target == target && receipt.Fingerprint == fingerprint {
			return receipt, true, nil
		}
	}
	return CleanupReceipt{}, false, nil
}

func CleanupReceiptForRequest(root, repo, operation, target, request string) (CleanupReceipt, bool, error) {
	ledger, err := Read(root)
	if err != nil {
		return CleanupReceipt{}, false, err
	}
	var found *CleanupReceipt
	for i := range ledger.CleanupReceipts {
		receipt := &ledger.CleanupReceipts[i]
		if receipt.Repo != repo || receipt.Operation != operation || receipt.Target != target || receipt.Request != request || !receipt.Owned || receipt.Assignment == "" {
			continue
		}
		if found != nil {
			return CleanupReceipt{}, false, errors.New("cleanup receipt request is ambiguous")
		}
		found = receipt
	}
	if found == nil {
		return CleanupReceipt{}, false, nil
	}
	return *found, true, nil
}

func PutCleanupReceipt(root string, receipt CleanupReceipt) error {
	if err := validateCleanupReceipts([]CleanupReceipt{receipt}); err != nil {
		return err
	}
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return err
	}
	defer release()
	ledger, err := readPath(path)
	if err != nil {
		return err
	}
	replaced := false
	for i, current := range ledger.CleanupReceipts {
		if current.Repo == receipt.Repo && current.Operation == receipt.Operation && current.Target == receipt.Target && current.Fingerprint == receipt.Fingerprint {
			ledger.CleanupReceipts[i], replaced = receipt, true
			break
		}
	}
	if !replaced {
		ledger.CleanupReceipts = append(ledger.CleanupReceipts, receipt)
	}
	completed := 0
	for _, current := range ledger.CleanupReceipts {
		if current.State == ReceiptComplete {
			completed++
		}
	}
	drop := completed - MaxCleanupReceipts
	if drop > 0 {
		next := ledger.CleanupReceipts[:0]
		for _, current := range ledger.CleanupReceipts {
			if drop > 0 && current.State == ReceiptComplete {
				drop--
				continue
			}
			next = append(next, current)
		}
		ledger.CleanupReceipts = next
	}
	return writePath(path, ledger)
}

// LifecycleEvidence returns the exact persisted schema and assignment JSON values.
// Cleanup receipts deliberately live in the same ledger but outside this evidence:
// recording an apply transaction must not make its already-approved plan stale.
func LifecycleEvidence(root string) ([]byte, error) {
	path, err := Address(root)
	if err != nil {
		return nil, err
	}
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		schema, _ := json.Marshal(Schema)
		return lifecycleEvidenceValues(schema, json.RawMessage("[]")), nil
	} else if err != nil {
		return nil, err
	}
	if _, err := readPath(path); err != nil {
		return nil, err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Schema      json.RawMessage `json:"schema"`
		Assignments json.RawMessage `json:"assignments"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if len(raw.Assignments) == 0 {
		raw.Assignments = json.RawMessage("[]")
	}
	return lifecycleEvidenceValues(raw.Schema, raw.Assignments), nil
}

func lifecycleEvidenceValues(schema, assignments []byte) []byte {
	return []byte(fmt.Sprintf("%d:%s%d:%s", len(schema), schema, len(assignments), assignments))
}

var (
	idPattern     = regexp.MustCompile(`^[0-9a-f]{32}$`)
	digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
	oidPattern    = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
)

func ValidIdentity(value string) bool { return idPattern.MatchString(value) }

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

func FindAssignmentByRequest(root, requestDigest string) (Assignment, bool, error) {
	ledger, err := Read(root)
	if err != nil {
		return Assignment{}, false, err
	}
	for _, assignment := range ledger.Assignments {
		if assignment.Request == requestDigest {
			return assignment, true, nil
		}
	}
	return Assignment{}, false, nil
}

// FindAssignmentForRequest resolves an opaque caller token through the ledger's
// request-digest owner. Lifecycle commands never compare caller tokens directly.
func FindAssignmentForRequest(root, request string) (Assignment, bool, error) {
	return FindAssignmentByRequest(root, RequestDigest(request))
}

func Assignments(root string) ([]Assignment, error) {
	ledger, err := Read(root)
	if err != nil {
		return nil, err
	}
	return append([]Assignment(nil), ledger.Assignments...), nil
}

// PutAssignment atomically inserts or updates one assignment ID. A request digest
// can name only one assignment, and an identical write preserves the ledger bytes.
func PutAssignment(root string, assignment Assignment) error {
	if err := ValidateAssignment(assignment); err != nil {
		return err
	}
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return err
	}
	defer release()
	ledger, err := readPath(path)
	if err != nil {
		return err
	}
	for i, current := range ledger.Assignments {
		if current.Request == assignment.Request && current.ID != assignment.ID {
			return errors.New("assignment request already belongs to another assignment")
		}
		if current.ID != assignment.ID {
			continue
		}
		if reflect.DeepEqual(current, assignment) {
			return nil
		}
		ledger.Assignments[i] = assignment
		return writePath(path, ledger)
	}
	ledger.Assignments = append(ledger.Assignments, assignment)
	return writePath(path, ledger)
}

// ReauthorizeAssignment swaps one request digest only after its caller's reversible
// external transition succeeds. The rollback keeps that derived state coherent when
// the expected-old write cannot commit.
func ReauthorizeAssignment(root, id, request string, verify func(Assignment) error, transition func(Assignment, Assignment) (func(), error), beforeCAS func(*Assignment)) (Assignment, error) {
	path, err := Address(root)
	if err != nil {
		return Assignment{}, err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return Assignment{}, err
	}
	defer release()
	ledger, err := readPath(path)
	if err != nil {
		return Assignment{}, err
	}
	for i := range ledger.Assignments {
		current := ledger.Assignments[i]
		if current.ID != id {
			continue
		}
		expectedOld := current.Request
		if err := verify(current); err != nil {
			return Assignment{}, err
		}
		if !ValidIdentity(id) || request == "" {
			return Assignment{}, errors.New("invalid reauthorization identity")
		}
		newDigest := RequestDigest(request)
		for j, other := range ledger.Assignments {
			if j != i && other.Request == newDigest {
				return Assignment{}, errors.New("request digest already belongs to another assignment")
			}
		}
		next := current
		next.Request = newDigest
		next.RequestToken = request
		rollback, err := transition(current, next)
		if err != nil {
			return Assignment{}, err
		}
		if beforeCAS != nil {
			beforeCAS(&ledger.Assignments[i])
		}
		if err := compareAndSwapRequestDigest(&ledger.Assignments[i], expectedOld, newDigest); err != nil {
			rollback()
			return Assignment{}, err
		}
		if err := writePath(path, ledger); err != nil {
			rollback()
			return Assignment{}, err
		}
		return current, nil
	}
	return Assignment{}, errors.New("assignment not found")
}

// RequestDigest derives the persisted identity for an opaque caller request.
func RequestDigest(value string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(value))) }

func compareAndSwapRequestDigest(assignment *Assignment, expectedOld, replacement string) error {
	if assignment.Request != expectedOld {
		return errors.New("assignment request changed during reauthorization")
	}
	assignment.Request = replacement
	return nil
}

// PurgeAssignments drops every assignment record keep rejects, plus every record this
// build can no longer read at all, and reports how many it dropped. Read is strict
// because a record it cannot account for must not authorize anything. The purge is
// deliberately not. A ledger written by an older binary is unreadable exactly when the
// standing cleaner is the only thing able to clear it. Refusing there would leave every
// later command reading the same unreadable file. A purge that drops nothing leaves the
// file's bytes untouched, so a re-run over converged state is a no-op.
func PurgeAssignments(root string, keep func(Assignment) bool) (int, error) {
	path, err := Address(root)
	if err != nil {
		return 0, err
	}
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return 0, err
	}
	defer release()
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	// One tolerant pass for the records the purge decides about, and one for everything
	// else in the file. Assignments are read as raw values so a single unreadable record
	// is dropped on its own rather than taking the whole ledger with it.
	var stored struct {
		Schema      int               `json:"schema"`
		Assignments []json.RawMessage `json:"assignments"`
	}
	if err := json.Unmarshal(data, &stored); err != nil {
		return 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	kept := make([]Assignment, 0, len(stored.Assignments))
	ids, requests := map[string]bool{}, map[string]bool{}
	// A record the legacy schema carries was never authorized to be there: Read refuses the
	// whole file over one. Under that schema every record is debris, whatever it says about
	// itself, so the loop that would judge them individually is skipped.
	if stored.Schema != LegacySchema {
		for _, record := range stored.Assignments {
			var assignment Assignment
			if json.Unmarshal(record, &assignment) != nil || ValidateAssignment(assignment) != nil {
				continue
			}
			if ids[assignment.ID] || requests[assignment.Request] || !keep(assignment) {
				continue
			}
			ids[assignment.ID], requests[assignment.Request] = true, true
			kept = append(kept, assignment)
		}
	}
	dropped := len(stored.Assignments) - len(kept)
	if dropped == 0 {
		return 0, nil
	}
	var rest struct {
		Entries         []Entry          `json:"entries"`
		CleanupReceipts []CleanupReceipt `json:"cleanup_receipts"`
	}
	if err := json.Unmarshal(data, &rest); err != nil {
		return 0, fmt.Errorf("purge intent ledger: %w", err)
	}
	ledger := Ledger{Entries: rest.Entries, Assignments: kept, CleanupReceipts: rest.CleanupReceipts}
	return dropped, writePath(path, ledger)
}

func DeleteAssignment(root, id string) error {
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return err
	}
	defer release()
	ledger, err := readPath(path)
	if err != nil {
		return err
	}
	next := ledger.Assignments[:0]
	for _, assignment := range ledger.Assignments {
		if assignment.ID != id {
			next = append(next, assignment)
		}
	}
	if len(next) == len(ledger.Assignments) {
		return nil
	}
	ledger.Assignments = next
	return writePath(path, ledger)
}
