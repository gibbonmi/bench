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

// These independent expectations make removal of a finding-confidence rule fail.
func TestCalibrationFindingAnchors(t *testing.T) {
	const discipline = ".agents/skills/bench-craft-review/references/finding-discipline.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: discipline, section: "What a confidence states", needle: "A finding carries a stated confidence as an integer from 0 to 10.", want: "calibration: a finding needs its stated confidence"},
		{file: discipline, section: "What a confidence states", needle: "The confidence never changes whether a finding blocks.", want: "calibration: the confidence must never change whether a finding blocks"},
		{file: discipline, section: "What a confidence states", needle: "`auto-fix` and `ask-user` label a finding `held`, and `no-op` labels it `refuted`.", want: "calibration: the dispositions need their label mapping"},
		{file: discipline, section: "What a confidence states", needle: "Optional advice carries no confidence.", want: "calibration: optional advice must carry no confidence"},
		{file: ".agents/skills/bench-craft-review/SKILL.md", section: "What a finding must cite", needle: "A finding also states its confidence as an integer from 0 to 10.", want: "calibration: the review skill must point at the finding confidence"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 6, needle: "Each actionable finding line carries its stated confidence.", want: "calibration: the pickup line must carry its stated confidence"},
	}}.check(t)
}

// These independent expectations make removal of a calibration score rule fail.
func TestCalibrationScoreAnchors(t *testing.T) {
	const score = ".agents/skills/bench-craft-line/references/calibration-score.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: score, section: "What the declaration states", needle: "The line declaration states the expected repair-round count and a stated confidence as an integer from 0 to 10.", want: "calibration: the line declaration needs its expected rounds and stated confidence"},
		{file: score, section: "How one claim scores", needle: "One claim's calibration score is `(p - label)^2`, with `p = n / 10` and label 1 for `held` or 0 for `refuted`.", want: "calibration: the Brier rule needs its score expression"},
		{file: score, section: "How one claim scores", needle: "An abstention scores 0, stays out of the Brier mean, and is counted apart.", want: "calibration: an abstention needs its own scoring rule"},
		{file: score, section: "Who writes a label", needle: "A label source is the gate, the coordinator's probe of the exact tree, or the reviewer's disposition.", want: "calibration: the three label sources must stay named"},
		{file: score, section: "Who writes a label", needle: "A model judgment is never a label source.", want: "calibration: a model judgment must never label a claim"},
		{file: score, section: "Who writes a label", needle: "The repair-attribution table's actual round count labels the expectation `held` when it equals the expected count and `refuted` otherwise.", want: "calibration: the round count must label the expectation"},
		{file: ".agents/skills/bench-craft-line/SKILL.md", section: "The declaration", needle: "Expected repair rounds: <count> / confidence <0-10>. `references/calibration-score.md` owns the score.", want: "calibration: the declaration must state expected repair rounds"},
	}}.check(t)
}

// These independent expectations make removal of a retro calibration duty fail.
func TestCalibrationRetroDutyAnchors(t *testing.T) {
	const command = ".agents/commands/bench-final-check.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: command, section: "Capture the implementation retro", needle: "The retro fills the calibration table with one row per labeled claim: surface, claim, status, confidence, label, and model, effort, and role.", want: "calibration: the retro must fill the calibration table"},
		{file: command, section: "Capture the implementation retro", needle: "The retro states the Brier mean, the pair count, and the abstention count below the table, with `unknown` for a mean over zero pairs.", want: "calibration: the retro must state the Brier mean and the counts"},
	}}.check(t)
}
