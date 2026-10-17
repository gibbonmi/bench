package anchors

import "testing"

// This independent expectation makes removal of the root conformance command fail.
func TestCheckFloorAnchors(t *testing.T) {
	const discipline = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: discipline, section: "In the charge", needle: "The root conformance pass is `bench test --check root-conformance`.", want: "check floor: the root conformance pass is bench test --check root-conformance"},
	}}.check(t)
}
