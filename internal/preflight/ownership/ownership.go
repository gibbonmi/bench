// Package ownership collects the authored planning-claim inventory of a spec: each premise
// with its proof state, and every diagnostic that keeps proof incomplete.
//
// The package reads the canonical owners and copies none of their tables. It imports no
// preflight package: the preflight gatherer calls Collect for the staleness audit, which
// grades the returned Snapshot.
package ownership

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"

	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/tickets"
)

// Source is the committed source identity that one collection reads: the explicit base and
// tip, and the spec path from the canonical selected-spec read.
type Source struct {
	Base, Tip string
	SpecPath  string
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
// is no review and no proof. PendingNative holds the validated native-future claims.
// Blocker names the first unmet condition, and it is empty when the inventory is complete
// and each required current-source premise has complete supported proof.
type ChargeReadiness struct {
	InventoryComplete bool
	PendingNative     []ManualClaim
	Blocker           string
}

// Snapshot is one immutable collection. Complete means full supported premise
// certification of a declared-complete inventory with no diagnostic; it is never reviewer
// approval, and a partial inventory or any manual claim keeps it false. Symbols holds the
// tip identity of each current symbol premise, in inventory order.
type Snapshot struct {
	Premises        []Premise
	Diagnostics     []Diagnostic
	Complete        bool
	ChargeReadiness ChargeReadiness
	Symbols         []SymbolIdentity
}

// Declaration is one declaration as one source holds it: the path that declares it, its
// spelling inside its package, its consumers declaration kind, and its normalized
// declaration syntax. The zero value is a declaration that the source does not hold.
type Declaration struct {
	Path, Name, Kind, Signature string
}

// Discrepancy is the class of the first difference between a cited owner identity and the
// declaration that a source holds. The empty class is an exact match.
type Discrepancy string

const (
	DiscrepancyAbsent    Discrepancy = "absent"
	DiscrepancyMoved     Discrepancy = "moved"
	DiscrepancyKind      Discrepancy = "kind"
	DiscrepancySignature Discrepancy = "signature"
	// DiscrepancyAmbiguous is a cited spelling that names more than one declaration.
	DiscrepancyAmbiguous Discrepancy = "ambiguous"
	// DiscrepancyContext is a declaration that only a file outside the default build
	// context holds, which the default typed load cannot grade.
	DiscrepancyContext Discrepancy = "unsupported-context"
	// DiscrepancyUnread is an identity that the census cannot read: a premise that does not
	// decode, a tip that does not load, or an owner path that holds no Go source.
	DiscrepancyUnread Discrepancy = "unread"
)

// SymbolIdentity is the tip identity of one current symbol premise: the cited owner, the
// declaration that the typed packages resolve it to, and the first discrepancy with its
// reason. A bare cited spelling that names more than one declaration is ambiguous in place
// of any other class, although the premise verdict can still pass.
type SymbolIdentity struct {
	Premise     string
	Owner       Owner
	Current     Declaration
	Discrepancy Discrepancy
	Reason      string
}

// Mismatch is the first part of d that differs from o, in the order path, kind, and
// signature, and its reason. A zero d is absent. at names the source that holds d.
func (o Owner) Mismatch(d Declaration, at string) (Discrepancy, string) {
	switch {
	case d == Declaration{}:
		return DiscrepancyAbsent, fmt.Sprintf("%s declares no %s at the %s", o.Package, localName(o), at)
	case d.Path != o.Path:
		return DiscrepancyMoved, fmt.Sprintf("%s is declared in %s, not %s", o.Symbol, d.Path, o.Path)
	case d.Kind != o.Kind:
		return DiscrepancyKind, fmt.Sprintf("%s is a %s, not a %s", o.Symbol, d.Kind, o.Kind)
	case d.Signature != o.Signature:
		return DiscrepancySignature, fmt.Sprintf("%s declares %q, not %q", o.Symbol, d.Signature, o.Signature)
	}
	return "", ""
}

// Collect grades the manifest's premise inventory at root against the parsed tickets. A
// premise that its collector cannot read stays unknown with that fault as its reason, so a
// failed read never reads as complete proof.
func Collect(root string, source Source, manifest Manifest, parsed []tickets.Ticket) Snapshot {
	census := &census{root: root, source: source}
	graded := census.premises(manifest.Premises)
	decisionPremises := census.decisions(manifest)
	manualPremises, pendingNative, nativeValid := census.manualClaims(manifest, parsed)
	allPremises := append(append(manifest.premises(source.Tip, graded), decisionPremises...), manualPremises...)
	allDiag := append(manifest.diagnostics(), census.checkSeams(manifest)...)
	snapshot := Snapshot{
		Premises:    allPremises,
		Diagnostics: allDiag,
		Symbols:     census.symbols,
	}
	snapshot.Complete = manifest.Present && manifest.Completeness == completenessDeclared && len(snapshot.Diagnostics) == 0 && nativeValid
	for _, premise := range snapshot.Premises {
		snapshot.Complete = snapshot.Complete && premise.State == StateGreen
	}
	snapshot.ChargeReadiness = chargeReadiness(manifest, snapshot.Premises, snapshot.Diagnostics, pendingNative, nativeValid)
	return snapshot
}

// premises projects the inventory's premise items onto premise rows. graded holds the row
// of each current premise that a supported collector graded. Any other current premise
// stays unknown, and a planned premise stays planned.
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
	return rows
}

// SubjectPaths is every repository path that one premise subject names, nested paths
// included, read through the subject decoder of the kind's collector. A subject that does not
// decode, or a kind that no collector grades, answers an error and no path.
func SubjectPaths(item PremiseItem) ([]string, error) {
	if shape, member := memberCensuses[item.Kind]; member {
		subject, err := decodeCensusSubject(item, shape)
		if err != nil {
			return nil, err
		}
		return subject.paths(), nil
	}
	project, ok := subjectPaths[item.Kind]
	if !ok {
		return nil, fmt.Errorf("no %s premise collector reads its subject", item.Kind)
	}
	paths, err := project(item)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

// subjectPaths projects the subject of each premise kind outside the member censuses onto
// its repository paths. An import premise's caller and target are package import paths, so
// its owner path is its one repository path.
var subjectPaths = map[string]func(PremiseItem) ([]string, error){
	"lines": func(item PremiseItem) ([]string, error) {
		subject, err := decodeLinesSubject(item)
		return []string{subject.Path}, err
	},
	"glob-members": func(item PremiseItem) ([]string, error) {
		subject, err := decodeGlobSubject(item)
		return []string{subject.Scope}, err
	},
	"headroom": func(item PremiseItem) ([]string, error) {
		subject, err := decodeHeadroomSubject(item)
		return []string{subject.Directory}, err
	},
	"symbol": ownerSubjectPath,
	"bound":  ownerSubjectPath,
	"import": func(item PremiseItem) ([]string, error) {
		subject, err := decodeImportSubject(item)
		return []string{subject.Owner.Path}, err
	},
	"anchor-kind": func(item PremiseItem) ([]string, error) {
		subject, err := decodeAnchorSubject(item)
		return []string{subject.Owner.Path}, err
	},
	"sentinel": ruleSubjectPaths,
	"error":    ruleSubjectPaths,
	"reused-test": func(item PremiseItem) ([]string, error) {
		subject, err := testSubject(item.Subject, "subject: ")
		return []string{subject.Path}, err
	},
	"reused-owner": func(item PremiseItem) ([]string, error) {
		owner, err := reusedOwnerSubject(item)
		return []string{owner.Path}, err
	},
}

// ownerSubjectPath is the owner path of a subject that is one owner identity.
func ownerSubjectPath(item PremiseItem) ([]string, error) {
	owner, err := ownerSubject(item)
	return []string{owner.Path}, err
}

// ruleSubjectPaths is the owner path and each producer path of a sentinel or error subject.
func ruleSubjectPaths(item PremiseItem) ([]string, error) {
	owner, producers, err := ruleSubject(item)
	paths := []string{owner.Path}
	for _, producer := range producers {
		paths = append(paths, producer.Path)
	}
	return paths, err
}

// diagnostics is the inventory's own diagnostic list. An absent fence is one diagnostic, so
// a spec with no inventory never reads as complete.
func (m Manifest) diagnostics() []Diagnostic {
	if !m.Present {
		return []Diagnostic{{ID: metadataID, Reason: "no " + MetadataLabel + " fence: the premise inventory is incomplete"}}
	}
	return m.Diagnostics
}

// SourceDeclaration projects one committed Go source without types: the top-level
// declaration of src that declares name by the consumers declaration-name forms, as path
// holds it, with the one normalization of the typed side. No types load here, so its kind
// is the declaration's syntax form. found is false when src declares no such name, and a
// source that does not parse is an error.
func SourceDeclaration(path string, src []byte, name string) (decl Declaration, found bool, err error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.SkipObjectResolution)
	if err != nil {
		return Declaration{}, false, err
	}
	for _, top := range topLevel(file) {
		if names, _ := consumers.DeclNames(top.node); !slices.Contains(names, name) {
			continue
		}
		kind := map[token.Token]string{token.FUNC: consumers.KindFunc, token.TYPE: consumers.KindType, token.CONST: consumers.KindConst, token.VAR: consumers.KindVar}[top.tok]
		if fn, ok := top.node.(*ast.FuncDecl); ok && fn.Recv != nil {
			kind = consumers.KindMethod
		}
		signature, err := normalize(fset, top.node, top.tok == token.CONST, name)
		return Declaration{Path: path, Name: name, Kind: kind, Signature: signature}, err == nil, err
	}
	return Declaration{}, false, nil
}
