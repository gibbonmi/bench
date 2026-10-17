package anchors

import "github.com/gibbonmi/bench/internal/retros"

// calibrationAnchors pin the step-scoped pickup-line confidence of the review phase and forbid
// copies of the scaffold's calibration columns and header in the final check.
var calibrationAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: RequireInStep, Section: "Process", Step: 6, Needle: "Each actionable finding line carries its stated confidence.", Diagnostic: "calibration: the pickup line must carry its stated confidence"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "surface, claim, status, confidence, label, and model, effort, and role", Diagnostic: "calibration: the final check restored a copy of the scaffold's calibration columns"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: retros.CalibrationHeader, Diagnostic: "calibration: the final check restored a copy of the scaffold's calibration table header"},
}
