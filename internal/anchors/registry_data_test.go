package anchors

import "testing"

// TestFinalCommunicationMarkerTuples keeps the final required tuples and
// retired-marker policy independent of production values, so omissions cannot
// re-derive the expectation green.
func TestFinalCommunicationMarkerTuples(t *testing.T) {
	const finalRoles = "Never assume the reviewer's decisions, and never assume a claim the gate could check instead"
	const retiredRoles = "NEVER assume, always verify"
	if ReviewerDecisionBoundaryMarker != finalRoles {
		t.Errorf("ReviewerDecisionBoundaryMarker = %q; want %q", ReviewerDecisionBoundaryMarker, finalRoles)
	}
	wanted := []Anchor{
		{File: ".bench/BENCH.md", Kind: RequireInSection, Section: "Roles", Needle: finalRoles},
		{File: ".bench/BENCH.md", Kind: RequireInSection, Section: "How to talk to me", Needle: "Clear beats dense"},
		{File: ".bench/BENCH.md", Kind: RequireInSection, Section: "Workflow", Needle: "Right-size the process"},
		{File: ".bench/BENCH.md", Kind: Forbid, Needle: retiredRoles},
	}
	entries := Entries()
	previous := -1
	for _, want := range wanted {
		at := -1
		matches := 0
		for i, entry := range entries {
			if entry.File == want.File && entry.Kind == want.Kind && entry.Section == want.Section && entry.Needle == want.Needle {
				if at < 0 {
					at = i
				}
				matches++
			}
		}
		if matches != 1 {
			t.Errorf("registry has %d rows matching %+v; want exactly one", matches, want)
		}
		if at < 0 {
			continue
		}
		if at <= previous {
			t.Errorf("registry places communication marker row %+v at index %d, out of final policy order", want, at)
		}
		previous = at
	}
}
