// Package ledger owns the intent ledger's persisted schema. It holds the schema
// constants, the record types with their constants, the identity patterns, the
// branch-ref helpers, and every record validator. The package imports nothing
// under internal/, so no importer of it can close a cycle, and the schema has one
// pure owner. The intent package owns the address, the locking, and the file
// effects, and it re-exports every name declared here.
package ledger

import "time"

const (
	LegacySchema = 1
	Schema       = 2
)

type Kind string

const (
	KindShift       Kind = "shift"
	KindWorktree    Kind = "worktree"
	KindClaudeAgent Kind = "claude-agent"
)

type Entry struct {
	Key       string    `json:"key"`
	Kind      Kind      `json:"kind"`
	CreatedAt time.Time `json:"created_at"`
	Worktree  string    `json:"worktree,omitempty"`
	Branch    string    `json:"branch,omitempty"`
	// Outcome and Recovery record a shift's final result state. Outcome is one of the
	// FT79 taxonomy's outcome names. Recovery is a pointer ("ref:<name>" | "worktree:<path>"
	// | "none") once a later slice adds snapshot machinery. Both fields are optional, so
	// every non-shift writer stays valid, as does every entry created before its writer
	// resolves an outcome.
	Outcome  string `json:"outcome,omitempty"`
	Recovery string `json:"recovery,omitempty"`
}

type Ledger struct {
	Schema          int              `json:"schema"`
	Entries         []Entry          `json:"entries"`
	Assignments     []Assignment     `json:"assignments,omitempty"`
	CleanupReceipts []CleanupReceipt `json:"cleanup_receipts,omitempty"`
}

const AssignmentRecordSchema = "bench-assignment/v1"

const assignmentBranchNamespace = "refs/heads/bench/assign/"

// RecoveryRefNamespace is the one name for the namespace preserved work lives under.
// Both the writer that puts refs there and the standing cleaner that sweeps it read
// this constant. Neither can address a namespace the other does not.
const RecoveryRefNamespace = "refs/bench/recovery/"

// ResetRefNamespace holds envelopes for recoverable assignment resets.
const ResetRefNamespace = "refs/bench/reset/"

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
