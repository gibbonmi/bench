package anchors

// chargeBindingAnchors pin the two repair-template field labels of the delegation discipline,
// and they forbid its retired repair fence and its rebase premise.
var chargeBindingAnchors = []Anchor{
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: Forbid, Needle: "A repair fence is the approved chunk union plus the exact paths that the review names.", Diagnostic: "charge binding: the repair fence restored the chunk union plus review-named paths"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "Repair-charge template", Needle: "Debug route:", Diagnostic: "delegation-discipline.md Repair-charge template dropped the debug route field"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "Repair-charge template", Needle: "Repair context:", Diagnostic: "delegation-discipline.md Repair-charge template dropped the repair context field"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: Forbid, Needle: "After a rebase changes the source digest", Diagnostic: "charge binding: the digest rule restored the rebase premise that the workflow rejects"},
}

// delegationDiscipline is the craft-delegate reference these anchors read.
const delegationDiscipline = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
