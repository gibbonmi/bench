package anchors

// checkFloorAnchors pin the one root conformance command of the check floor.
var checkFloorAnchors = []Anchor{
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "In the charge", Needle: "The root conformance pass is `bench test --check root-conformance`.", Diagnostic: "check floor: the root conformance pass is bench test --check root-conformance"},
}
