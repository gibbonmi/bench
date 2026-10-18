package staleness

import (
	"slices"

	"github.com/gibbonmi/bench/internal/preflight/ownership"
)

// completeness is each diagnostic that keeps the audit incomplete although no collector of
// the audit names it. manifest is the tip inventory and snapshot its closure collection.
//
// A premise that the snapshot holds unproved stays open with its own reason, unless the
// cited-symbols audit grades it. So a manual claim, a planned premise, and a premise that no
// collector supports each keep manual review required. An inventory that the snapshot does
// not hold complete stays open with the snapshot's blocker, unless a snapshot diagnostic
// already names the fault. A present inventory that declares no claim cannot certify its own
// scope, so an empty container never reads as a complete one.
func completeness(manifest ownership.Manifest, snapshot ownership.Snapshot) []Diagnostic {
	var open []Diagnostic
	for _, premise := range snapshot.Premises {
		graded := slices.ContainsFunc(snapshot.Symbols, func(s ownership.SymbolIdentity) bool { return s.Premise == premise.ID })
		if premise.State != ownership.StateGreen && !graded {
			open = append(open, Diagnostic{Subject: premise.ID, Reason: premise.Reason})
		}
	}
	if len(snapshot.Diagnostics) == 0 && !snapshot.ChargeReadiness.InventoryComplete {
		open = append(open, Diagnostic{Subject: ownership.MetadataLabel, Reason: "the planning inventory is not complete: " + snapshot.ChargeReadiness.Blocker})
	}
	claims := len(manifest.Changes) + len(manifest.Premises) + len(manifest.Seams) + len(manifest.Decisions) + len(manifest.ManualClaims)
	if manifest.Present && len(manifest.Diagnostics) == 0 && claims == 0 {
		open = append(open, Diagnostic{Subject: ownership.MetadataLabel, Reason: "the planning inventory declares no claim, so it cannot certify its own scope"})
	}
	return open
}
