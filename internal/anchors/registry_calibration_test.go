package anchors

import "testing"

// These independent expectations make the return of the calibration column copy fail.
func TestCalibrationAnchors(t *testing.T) {
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: ".agents/commands/bench-final-check.md", needle: "surface, claim, status, confidence, label, and model, effort, and role", want: "calibration: the final check restored a copy of the scaffold's calibration columns", forbidden: true},
	}}.check(t)
}
