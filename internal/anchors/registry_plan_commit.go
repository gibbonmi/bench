package anchors

// planCommitAnchors pin the two halves of the plan-commit route in the implement phase: a plan
// commit that a dispatch follows uses plain `bench commit`, and a plan commit that no dispatch
// follows uses the commit chain.
var planCommitAnchors = []Anchor{
	{Group: AfterImplementSpec, File: implementPhase, Kind: Require, Class: Route, Needle: "a dispatch follows uses plain `bench commit`", Diagnostic: "plan commit: a plan commit that a dispatch follows uses plain bench commit"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Require, Class: Command, Needle: "a plan commit that no dispatch follows uses `bench commit --preflight-build <slug>`", Diagnostic: "plan commit: a plan commit that no dispatch follows uses bench commit --preflight-build <slug>"},
}
