package anchors

import "testing"

// These independent expectations make removal of a repair-template field, or the return of a
// retired repair fence or rebase premise, fail.
func TestChargeBindingAnchors(t *testing.T) {
	const discipline = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: discipline, needle: "A repair fence is the approved chunk union plus the exact paths that the review names.", want: "charge binding: the repair fence restored the chunk union plus review-named paths", forbidden: true},
		{file: discipline, section: "Repair-charge template", needle: "Debug route:", want: "delegation-discipline.md Repair-charge template dropped the debug route field"},
		{file: discipline, section: "Repair-charge template", needle: "Repair context:", want: "delegation-discipline.md Repair-charge template dropped the repair context field"},
		{file: discipline, needle: "After a rebase changes the source digest", want: "charge binding: the digest rule restored the rebase premise that the workflow rejects", forbidden: true},
	}}.check(t)
}
