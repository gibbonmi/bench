package anchors

import "testing"

// These independent expectations make removal of a check-floor guidance rule fail.
func TestCheckFloorAnchors(t *testing.T) {
	const (
		implement  = ".agents/commands/bench-implement-spec.md"
		mapRules   = ".agents/skills/bench-craft-spec/references/map-discipline.md"
		discipline = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	)
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: implement, section: "Build", needle: "The completion plan carries the check floor of each ticket, and the `check-floor` row of `bench preflight build` grades the plan against that floor.", want: "check floor: the completion plan carries each ticket's floor, and check-floor grades it"},
		{file: implement, section: "Build", needle: "The ticket commit rewrites each covered `planned` seam cell into the citation form, so that `bench record chunk` accepts the chunk.", want: "check floor: the ticket commit rewrites each covered planned seam cell"},
		{file: mapRules, section: "Per row", needle: "The `planned` marker lasts until the covering ticket commits.", want: "check floor: the planned marker lasts until the covering ticket commits"},
		{file: discipline, section: "In the charge", needle: "The root conformance pass is `bench test --check root-conformance`.", want: "check floor: the root conformance pass is bench test --check root-conformance"},
	}}.check(t)
}
