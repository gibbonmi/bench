package anchors

import "testing"

// These independent expectations make removal of a claim-schema rule fail.
func TestCalibrationAnchors(t *testing.T) {
	const reference = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: reference, section: "Claim schema", needle: "A done-claim row carries a `status` of `verified`, `claimed`, or `abstained` and a stated confidence as an integer from 0 to 10.", want: "calibration: a done-claim row needs its status and stated confidence"},
		{file: reference, section: "Claim schema", needle: "`verified` means the author ran the named check and returns its red-to-green log, and `claimed` means an assertion with no executed check.", want: "calibration: verified and claimed need their evidence definitions"},
		{file: reference, section: "Claim schema", needle: "A claim carries no free-text field.", want: "calibration: a claim must carry no free-text field"},
		{file: reference, section: "Claim schema", needle: "A delegate that cannot state a confidence returns the row `abstained` with no confidence.", want: "calibration: an unstatable confidence needs an abstained row"},
		{file: reference, section: "Claim schema", needle: "No later probe turns an abstention into a refuted claim.", want: "calibration: an abstention must never become a refuted claim"},
		{file: reference, section: "Claim schema", needle: "The coordinator's probe of the exact tree labels a done-claim row `held` or `refuted`, whatever its status.", want: "calibration: the coordinator's tree probe labels each done-claim row"},
		{file: ".agents/skills/bench-craft-delegate/SKILL.md", section: "The charge", needle: "The charge names the `Claim schema` section of `references/delegation-discipline.md` as the return shape.", want: "calibration: the charge must name the claim schema section"},
	}}.check(t)
}
