package anchors

import "github.com/gibbonmi/bench/internal/retros"

// calibrationAnchors forbid copies of the scaffold's calibration columns and header in the
// final check.
var calibrationAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "surface, claim, status, confidence, label, and model, effort, and role", Diagnostic: "calibration: the final check restored a copy of the scaffold's calibration columns"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: retros.CalibrationHeader, Diagnostic: "calibration: the final check restored a copy of the scaffold's calibration table header"},
}
