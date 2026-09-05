package intent

import "github.com/gibbonmi/bench/internal/intent/ledger"

// The intent package re-exports every name the leaf ledger package declares, so an
// importing file keeps its `intent.` spelling and its import line. A type is an
// alias, and a function and a constant are a declared value.

type (
	Kind            = ledger.Kind
	Entry           = ledger.Entry
	Ledger          = ledger.Ledger
	AssignmentState = ledger.AssignmentState
	Recovery        = ledger.Recovery
	Assignment      = ledger.Assignment
	CleanupReceipt  = ledger.CleanupReceipt
)

const (
	LegacySchema = ledger.LegacySchema
	Schema       = ledger.Schema

	KindShift       = ledger.KindShift
	KindWorktree    = ledger.KindWorktree
	KindClaudeAgent = ledger.KindClaudeAgent

	AssignmentRecordSchema = ledger.AssignmentRecordSchema
	RecoveryRefNamespace   = ledger.RecoveryRefNamespace

	StateActive         = ledger.StateActive
	StateCleanupPending = ledger.StateCleanupPending
	StateRecovered      = ledger.StateRecovered
	StateComplete       = ledger.StateComplete

	CleanupReceiptSchema  = ledger.CleanupReceiptSchema
	ReceiptInFlight       = ledger.ReceiptInFlight
	ReceiptComplete       = ledger.ReceiptComplete
	MaxCleanupReceipts    = ledger.MaxCleanupReceipts
	ReceiptPhasePlanned   = ledger.ReceiptPhasePlanned
	ReceiptPhasePreserved = ledger.ReceiptPhasePreserved
	ReceiptPhaseRemoving  = ledger.ReceiptPhaseRemoving
	ReceiptPhaseRemoved   = ledger.ReceiptPhaseRemoved
	ReceiptPhaseBranch    = ledger.ReceiptPhaseBranch
	ReceiptPhaseTerminal  = ledger.ReceiptPhaseTerminal
)

var (
	AssignmentBranchPrefix = ledger.AssignmentBranchPrefix
	AssignmentBranchRef    = ledger.AssignmentBranchRef
	RecoveryRefPrefix      = ledger.RecoveryRefPrefix
	RequestDigest          = ledger.RequestDigest
	ValidIdentity          = ledger.ValidIdentity
	ValidateAssignment     = ledger.ValidateAssignment

	validAssignmentBranchRef = ledger.ValidAssignmentBranchRef
	validEntry               = ledger.ValidateEntry
	validateCleanupReceipts  = ledger.ValidateCleanupReceipts
)
