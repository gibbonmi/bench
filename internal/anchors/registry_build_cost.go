package anchors

// buildCostAnchors forbid the retired run of each verification requirement after a ticket
// commit and the retired check-floor row in the build guidance and the completion-plan
// reference. They also forbid a second copy of a cost budget in the drain guidance, because
// the budget table of the cost renderer is the one source.
var buildCostAnchors = []Anchor{
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "then runs each verification requirement once at that commit", Diagnostic: "build cost: the implement phase restored the run of each verification requirement"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "check-floor", Diagnostic: "build cost: the implement phase restored the retired check-floor row"},
	{Group: AfterImplementSpec, File: completionPlanReference, Kind: Forbid, Needle: "check-floor", Diagnostic: "build cost: the completion-plan reference restored the retired check-floor row"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "60 min", Diagnostic: "build cost: the drain repeats the wall budget of the cost renderer"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "7.1", Diagnostic: "build cost: the drain repeats the bookkeeping budget of the cost renderer"},
}

// completionPlanReference is the craft-spec reference that states the completion plan schema.
const completionPlanReference = ".agents/skills/bench-craft-spec/references/completion-plan.md"
