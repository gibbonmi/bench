// Package staleness is the mechanical audit of a reviewed spec against its current source.
// Collect reads the exact committed members of the audit scope at a baseline and a tip.
// Decide classifies those members into drift, check states, diagnostics, completeness, and
// the manual-review decision. Completeness says that every required collector ran; it says
// nothing about semantic equivalence.
//
// The package consumes the values of the ownership owner. It imports neither its parent
// preflight package nor another workflow state owner.
package staleness

import "github.com/gibbonmi/bench/internal/preflight/ownership"

// Absent is the cell of the side of a drift entry that holds no member. It is never a path
// or an object ID, so an absent side cannot read as an empty one.
const Absent = "absent"

// State is the state of one mechanical check. A supported failed predicate is red, and an
// input that no supported collector grades is diagnostic.
type State string

const (
	StateGreen      State = "green"
	StateRed        State = "red"
	StateDiagnostic State = "diagnostic"
)

// The three mechanical checks. checkNames holds them in render order.
const (
	CheckSpecDrift    = "spec-drift"
	CheckCitedSymbols = "cited-symbols"
	CheckCitedTests   = "cited-tests"
)

var checkNames = []string{CheckSpecDrift, CheckCitedSymbols, CheckCitedTests}

// Object is the committed identity of one tracked member: its six-digit mode and its object
// ID. The zero value is a side that holds no member.
type Object struct {
	Mode, ID string
}

// PathState is one audit-scope path and its object at the baseline and at the tip.
type PathState struct {
	Path          string
	Before, After Object
}

// Diagnostic is one audit input that no supported collector grades. It keeps the audit
// incomplete and manual review required, but it does not refuse the collection. Subject is
// the check or the input that it holds open.
type Diagnostic struct {
	Subject, Reason string
}

// AuditFacts is one immutable collection over one exact committed pair. Paths holds every
// member of the audit scope at either identity, in path order, so a deleted member keeps its
// row. Symbols holds each cited symbol in inventory order. Decide grades these values and
// reads nothing else.
type AuditFacts struct {
	Baseline, Tip string
	SpecPath      string
	Paths         []PathState
	Symbols       []SymbolState
	Diagnostics   []Diagnostic
}

// SymbolState is the audit of one cited symbol: its premise and cited owner, the declaration
// that the baseline source and the tip types hold, and the first discrepancy with its reason.
// An empty discrepancy is an exact match at both identities.
type SymbolState struct {
	Premise           string
	Owner             ownership.Owner
	Baseline, Current ownership.Declaration
	Discrepancy       ownership.Discrepancy
	Reason            string
}

// DriftEntry is one changed member. Before and After are the paths at the two identities,
// and the object cells are their exact object IDs; Absent fills the side that holds no
// member. Kind is added, deleted, modified, renamed, or mode, and a mode difference between
// two present sides adds its before and after modes to the kind.
type DriftEntry struct {
	Before, After             string
	Kind                      string
	BeforeObject, AfterObject string
}

// Check is the state of one named mechanical check.
type Check struct {
	Name  string
	State State
}

// Audit is the decision over one AuditFacts. Complete means that every required collector
// ran successfully. ManualRequired is false only for a complete audit with empty drift and
// every check green.
type Audit struct {
	Baseline, Tip  string
	Checks         []Check
	Drift          []DriftEntry
	Diagnostics    []Diagnostic
	Complete       bool
	ManualRequired bool
}

// The stages that a collection refusal names.
const (
	StageBaselineSpec = "baseline spec"
	StageTipSpec      = "tip spec"
	StageMembers      = "audit members"
	StageSpecial      = "special source"
	StageSymbol       = "cited symbol source"
)

// Error is a collection refusal: the stage that failed, the subject it read, and the cause.
// A refused collection returns no facts, so a partial enumeration never reads as complete.
type Error struct {
	Stage, Subject string
	Cause          error
}

func (e *Error) Error() string { return e.Stage + " " + e.Subject + ": " + e.Cause.Error() }

func (e *Error) Unwrap() error { return e.Cause }
