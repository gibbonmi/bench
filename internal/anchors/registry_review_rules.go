package anchors

// reviewRuleAnchors pin the finding discipline's reader rules and its axis-return rule.
var reviewRuleAnchors = []Anchor{
	{Group: AfterImplementSpec, File: findingDiscipline, Kind: RequireInSection, Section: "What a citation points at", Needle: "A review of gate-anchored prose names the anchor state of each sentence that it proposes to change.", Diagnostic: "review rules: an anchored-prose review names each sentence's anchor state"},
	{Group: AfterImplementSpec, File: findingDiscipline, Kind: RequireInSection, Section: "Where an axis under-reads", Needle: "A Coverage finding describes the tree before the probes of the axis.", Diagnostic: "review rules: a Coverage finding describes the pre-probe tree"},
	{Group: AfterImplementSpec, File: findingDiscipline, Kind: RequireInSection, Section: "Where an axis under-reads", Needle: "An axis reads the seam cell of a row before it judges a review-owned row unmet.", Diagnostic: "review rules: an axis reads the seam cell before it judges a review-owned row"},
	{Group: AfterImplementSpec, File: findingDiscipline, Kind: RequireInSection, Section: "What an axis return carries", Needle: "An axis return lists the evidence cursors that it fetched.", Diagnostic: "review rules: an axis return lists its fetched cursors"},
}

// findingDiscipline is the craft-review reference that states what a finding carries.
const findingDiscipline = ".agents/skills/bench-craft-review/references/finding-discipline.md"
