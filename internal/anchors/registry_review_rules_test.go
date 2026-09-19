package anchors

import "testing"

// These independent expectations make removal of a finding reader rule or the return rule fail.
func TestReviewRuleAnchors(t *testing.T) {
	const discipline = ".agents/skills/bench-craft-review/references/finding-discipline.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: discipline, section: "What a citation points at", needle: "A review of gate-anchored prose names the anchor state of each sentence that it proposes to change.", want: "review rules: an anchored-prose review names each sentence's anchor state"},
		{file: discipline, section: "Where an axis under-reads", needle: "A Coverage finding describes the tree before the probes of the axis.", want: "review rules: a Coverage finding describes the pre-probe tree"},
		{file: discipline, section: "Where an axis under-reads", needle: "An axis reads the seam cell of a row before it judges a review-owned row unmet.", want: "review rules: an axis reads the seam cell before it judges a review-owned row"},
		{file: discipline, section: "What an axis return carries", needle: "An axis return lists the evidence cursors that it fetched.", want: "review rules: an axis return lists its fetched cursors"},
	}}.check(t)
}
