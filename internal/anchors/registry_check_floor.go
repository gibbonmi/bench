package anchors

// checkFloorAnchors pin the check floor that the completion plan carries, the planned-cell
// rewrite at the ticket commit, the lifetime of the `planned` marker, and the one root
// conformance command.
var checkFloorAnchors = []Anchor{
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "The completion plan carries the check floor of each ticket, and the `check-floor` row of `bench preflight build` grades the plan against that floor.", Diagnostic: "check floor: the completion plan carries each ticket's floor, and check-floor grades it"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "The ticket commit rewrites each covered `planned` seam cell into the citation form, so that `bench record chunk` accepts the chunk.", Diagnostic: "check floor: the ticket commit rewrites each covered planned seam cell"},
	{Group: AfterImplementSpec, File: mapDiscipline, Kind: RequireInSection, Section: "Per row", Needle: "The `planned` marker lasts until the covering ticket commits.", Diagnostic: "check floor: the planned marker lasts until the covering ticket commits"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "In the charge", Needle: "The root conformance pass is `bench test --check root-conformance`.", Diagnostic: "check floor: the root conformance pass is bench test --check root-conformance"},
}
