package anchors

import "testing"

// These independent expectations make removal of the pickup-line confidence, or the return
// of the calibration column copy, fail.
func TestCalibrationAnchors(t *testing.T) {
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 6, needle: "Each actionable finding line carries its stated confidence.", want: "calibration: the pickup line must carry its stated confidence"},
		{file: ".agents/commands/bench-final-check.md", needle: "surface, claim, status, confidence, label, and model, effort, and role", want: "calibration: the final check restored a copy of the scaffold's calibration columns", forbidden: true},
	}}.check(t)
}
