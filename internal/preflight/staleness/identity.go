package staleness

import (
	"slices"

	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
)

// diagnosticDiscrepancies are the discrepancy classes that no supported collector grades: an
// ambiguous bare name, a declaration outside the default build context, and an unread
// identity. Each one stays diagnostic. Any other discrepancy is red.
var diagnosticDiscrepancies = []ownership.Discrepancy{ownership.DiscrepancyAmbiguous, ownership.DiscrepancyContext, ownership.DiscrepancyUnread}

// symbols audits each cited symbol that the ownership census identified at the tip against
// the baseline object of its owner path in paths. The baseline reads the declaration that
// the tip resolves to, so an alias keeps its consumers origin, and the cited spelling when
// the tip holds none. A diagnostic class at the tip skips the baseline read. Each diagnostic
// state also answers one diagnostic, and a baseline object that does not read refuses the
// collection.
func symbols(root string, paths []PathState, identities []ownership.SymbolIdentity) ([]SymbolState, []Diagnostic, error) {
	var states []SymbolState
	var diagnostics []Diagnostic
	for _, identity := range identities {
		state := SymbolState{Premise: identity.Premise, Owner: identity.Owner, Current: identity.Current, Discrepancy: identity.Discrepancy, Reason: identity.Reason}
		if !slices.Contains(diagnosticDiscrepancies, state.Discrepancy) {
			baseline, unread, err := baselineDeclaration(root, paths, identity)
			switch {
			case err != nil:
				return nil, nil, err
			case state.Discrepancy != "":
			case unread != "":
				state.Discrepancy, state.Reason = ownership.DiscrepancyUnread, unread
			default:
				state.Discrepancy, state.Reason = identity.Owner.Mismatch(baseline, "source base")
			}
			state.Baseline = baseline
		}
		if slices.Contains(diagnosticDiscrepancies, state.Discrepancy) {
			diagnostics = append(diagnostics, Diagnostic{Subject: state.Premise, Reason: "cited symbol " + state.Owner.Symbol + " is " + string(state.Discrepancy) + ": " + state.Reason})
		}
		states = append(states, state)
	}
	return states, diagnostics, nil
}

// baselineDeclaration is the declaration that the baseline object of the owner path holds.
// A path that the baseline does not hold holds no declaration. unread is the reason of a
// source that does not parse, and an object that does not read is a refusal.
func baselineDeclaration(root string, paths []PathState, identity ownership.SymbolIdentity) (decl ownership.Declaration, unread string, err error) {
	owner := identity.Owner
	at := slices.IndexFunc(paths, func(state PathState) bool { return state.Path == owner.Path })
	if at < 0 || paths[at].Before.ID == "" {
		return ownership.Declaration{}, "", nil
	}
	src, err := git.ReadControlBlob(root, paths[at].Before.ID)
	if err != nil {
		return ownership.Declaration{}, "", &Error{Stage: StageSymbol, Subject: owner.Path, Cause: err}
	}
	name := identity.Current.Name
	if name == "" {
		name, _ = consumers.LocalName(owner.Package, owner.Symbol)
	}
	if decl, _, err = ownership.SourceDeclaration(owner.Path, src, name); err != nil {
		return ownership.Declaration{}, owner.Path + " holds no Go source at the source base: " + err.Error(), nil
	}
	return decl, "", nil
}

// symbolEntry is the drift entry of one red cited symbol: the owner path at the baseline and
// the declaring path at the tip, a kind that names the symbol and its discrepancy, and the
// normalized syntax of each side as its object, or the kind of a constant. Absent fills a
// side that holds no declaration.
func symbolEntry(state SymbolState) DriftEntry {
	e := DriftEntry{Before: Absent, After: Absent, BeforeObject: Absent, AfterObject: Absent, Kind: "symbol " + state.Owner.Symbol + " " + string(state.Discrepancy)}
	if state.Baseline.Path != "" {
		e.Before, e.BeforeObject = state.Baseline.Path, identityCell(state.Baseline)
	}
	if state.Current.Path != "" {
		e.After, e.AfterObject = state.Current.Path, identityCell(state.Current)
	}
	return e
}

// identityCell is the normalized syntax of a declaration, or its kind when the syntax is
// empty, as it is for a constant.
func identityCell(d ownership.Declaration) string {
	if d.Signature == "" {
		return d.Kind
	}
	return d.Signature
}
