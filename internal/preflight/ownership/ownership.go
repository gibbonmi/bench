// Package ownership collects the canonical ownership facts of a spec's tickets. Each
// requirement is a path that a ticket's write requires through the fixture, registry, or
// anchor owner, with every source reason that requires it. The snapshot also carries the
// authored planning-claim inventory and every diagnostic that keeps proof incomplete.
//
// The package reads the canonical owners and copies none of their tables. It imports no
// preflight package: the preflight gatherer calls Collect, and its pure Decide grades the
// returned Snapshot without I/O.
package ownership

import (
	"github.com/gibbonmi/bench/internal/tickets"
)

// Source is the committed source identity that one collection reads: the explicit base and
// tip, the spec path from the canonical selected-spec read, the committed path status
// projection, and the caller's verified kit context.
type Source struct {
	Base, Tip string
	SpecPath  string
	Paths     []PathStatus
	Kit       Kit
}

// PathStatus is one committed path of the source range and its change status.
type PathStatus struct {
	Path   string
	Status Status
}

// Status is the committed change status of one path.
type Status string

const (
	StatusChanged Status = "changed"
	StatusDeleted Status = "deleted"
)

// Kit is the caller's authenticated kit context. Identity is the verified source digest of
// the kit build that the caller runs, and it is empty when no verified kit source exists.
// The collector never discovers or runs a kit of its own.
type Kit struct {
	Identity string
}

// Kind is the canonical owner family that requires a path. A caller requirement is a file
// that statically uses a declared changed declaration. An owner requirement is the exact
// rule source that a declared reused owner pins; the ticket reads it and need not own it.
type Kind string

const (
	KindFixture  Kind = "fixture"
	KindRegistry Kind = "registry"
	KindAnchor   Kind = "anchor"
	KindCaller   Kind = "caller"
	KindOwner    Kind = "owner"
)

// Scope tells tree evidence, which binds committed owner bytes, from kit evidence, which
// binds the verified kit source of a compiled table.
type Scope string

const (
	ScopeTree Scope = "tree"
	ScopeKit  Scope = "kit"
)

// Provenance is one source reason for one requirement: the canonical owner file, the line
// there, the rule it applies, a reason that names the path it was applied to, and the
// source identity it binds. Line is zero when the owner names the path as a whole, or when
// no verified source locates it; Source is then empty for kit evidence.
type Provenance struct {
	Owner     string
	Line      int
	Rule      string
	Reason    string
	Base, Tip string
	Source    string
	Scope     Scope
}

// Requirement is one path that a ticket entry requires, directly or through the fixed
// point, with every independent reason that requires it. Entry is the originating
// `Writes:` entry, so a transitive requirement names the write that started its chain. An
// owner requirement starts at no write, so its Entry is the premise that declares it.
type Requirement struct {
	Ticket, Entry, Path string
	Kind                Kind
	Provenance          []Provenance
}

// State is the proof state of one premise row. Only StateGreen is complete supported
// proof, and no other state becomes green.
type State string

const (
	StateGreen   State = "green"
	StateUnknown State = "unknown"
	StatePlanned State = "planned"
	StateManual  State = "manual"
)

// Premise is one premise row: its identity, proof state, the source tip it was read at,
// and the reason for that state.
type Premise struct {
	ID     string
	State  State
	Source string
	Reason string
}

// Diagnostic is one fact that keeps the collection incomplete without refusing it.
type Diagnostic struct {
	ID, Reason string
}

// ChargeReadiness is the charge-specific projection of the premise inventory.
// InventoryComplete means a declared-complete inventory passed its structural checks; it
// is no review and no proof. RequiredCurrent holds every required current-source premise,
// and CurrentComplete means each one has complete supported proof. PendingNative holds the
// validated native-future claims.
type ChargeReadiness struct {
	InventoryComplete bool
	CurrentComplete   bool
	RequiredCurrent   []string
	PendingNative     []ManualClaim
}

// Snapshot is one immutable collection. Complete means full supported premise
// certification of a declared-complete inventory with no diagnostic; it is never reviewer
// approval, and a partial inventory or any manual claim keeps it false.
type Snapshot struct {
	Requirements    []Requirement
	Premises        []Premise
	Diagnostics     []Diagnostic
	Complete        bool
	ChargeReadiness ChargeReadiness
}

// The stages a collection refusal names.
const (
	StageFixtures = "fixture inventory"
	StageAnchors  = "anchor registry"
	StagePath     = "ownership path"
)

// Error is a collection refusal: the stage that failed, the subject it read, and the cause.
type Error struct {
	Stage, Subject string
	Cause          error
}

func (e *Error) Error() string { return e.Stage + " " + e.Subject + ": " + e.Cause.Error() }

func (e *Error) Unwrap() error { return e.Cause }

// Collect reads the canonical closure of every parsed ticket at root, with each Grammar
// binding classified by the committed change it covers and each declared changed
// declaration's static callers, and grades the manifest's premise inventory. A failed owner
// read returns an *Error and no snapshot, so a partial collection never reads as a
// complete one.
func Collect(root string, source Source, manifest Manifest, parsed []tickets.Ticket) (Snapshot, error) {
	census := &census{root: root, source: source}
	graded, owners := census.premises(manifest.Premises)
	requirements, diagnostics, err := closure(root, source, parsed, manifest.Changes, census, owners)
	if err != nil {
		return Snapshot{}, err
	}
	snapshot := Snapshot{
		Requirements: requirements,
		Premises:     manifest.premises(source.Tip, graded),
		Diagnostics:  append(manifest.diagnostics(), diagnostics...),
	}
	snapshot.Complete = manifest.Present && manifest.Completeness == completenessDeclared && len(snapshot.Diagnostics) == 0
	for _, premise := range snapshot.Premises {
		snapshot.Complete = snapshot.Complete && premise.State == StateGreen
	}
	snapshot.ChargeReadiness = manifest.readiness()
	return snapshot, nil
}

// premises projects the inventory's premise and manual items onto premise rows. graded
// holds the row of each current premise that a supported collector graded. Any other
// current premise stays unknown; a planned premise stays planned, and a manual claim stays
// manual.
func (m Manifest) premises(tip string, graded map[string]Premise) []Premise {
	var rows []Premise
	for _, item := range m.Premises {
		row, ok := graded[item.ID]
		switch {
		case item.State == "planned":
			row = Premise{ID: item.ID, State: StatePlanned, Source: tip, Reason: "a planned " + item.Kind + " premise certifies no current source"}
		case !ok:
			row = Premise{ID: item.ID, State: StateUnknown, Source: tip, Reason: "no " + item.Kind + " premise collector certifies current source"}
		}
		rows = append(rows, row)
	}
	for _, claim := range m.ManualClaims {
		rows = append(rows, Premise{ID: claim.ID, State: StateManual, Source: tip, Reason: string(claim.Applicability) + " manual claim: " + claim.Reason})
	}
	return rows
}

// diagnostics is the inventory's own diagnostic list. An absent fence is one diagnostic,
// because the closure still proposes without an inventory.
func (m Manifest) diagnostics() []Diagnostic {
	if !m.Present {
		return []Diagnostic{{ID: metadataID, Reason: "no " + MetadataLabel + " fence: the premise inventory is incomplete"}}
	}
	return m.Diagnostics
}

// readiness is the charge projection of the inventory. It reads no collector proof, so
// CurrentComplete holds only when no current premise is required.
func (m Manifest) readiness() ChargeReadiness {
	ready := ChargeReadiness{InventoryComplete: m.Present && len(m.Diagnostics) == 0 && m.Completeness == completenessDeclared}
	for _, item := range m.Premises {
		if item.State == "current" {
			ready.RequiredCurrent = append(ready.RequiredCurrent, item.ID)
		}
	}
	for _, claim := range m.ManualClaims {
		switch claim.Applicability {
		case CurrentPremise:
			ready.RequiredCurrent = append(ready.RequiredCurrent, claim.ID)
		case NativeFuture:
			ready.PendingNative = append(ready.PendingNative, claim)
		}
	}
	ready.CurrentComplete = ready.InventoryComplete && len(ready.RequiredCurrent) == 0
	return ready
}
