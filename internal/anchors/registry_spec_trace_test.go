package anchors

import "testing"

// These independent expectations make removal of a spec grader trace rule fail.
func TestSpecGraderTraceAnchors(t *testing.T) {
	const discipline = ".agents/skills/bench-craft-spec/references/map-discipline.md"
	const section = "Before the map locks"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: discipline, section: section, needle: "`Pin operators` quotes the comparison operator that the grader of each pin row applies.", want: "map discipline: the pre-review checklist quotes each pin row's grader operator"},
		{file: discipline, section: section, needle: "`Entry reads` lists each unexported read below an entry without an internal form, and names its grader.", want: "map discipline: the pre-review checklist traces each unexported entry read to its grader"},
		{file: discipline, section: section, needle: "`Derived expectations` names the grader of each derived expectation.", want: "map discipline: the pre-review checklist names the grader of each derived expectation"},
		{file: discipline, section: section, needle: "No expectation comes from the code under test.", want: "map discipline: no expectation comes from the code under test"},
		{file: discipline, section: section, needle: "`Consolidated rules`, when a spec consolidates repeated rules, gives a consolidation table of each site's old rule and new rule.", want: "map discipline: a rule consolidation gives each site's old and new rule"},
		{file: discipline, section: section, needle: "Each changed cell of the consolidation table maps to an acceptance row, a flagged addition, or a Won't handle line.", want: "map discipline: each changed consolidation cell takes a row, a flagged addition, or a Won't handle line"},
		{file: discipline, section: section, needle: "`Quantified obligations` checks every quantified obligation at each affected site and across all tickets.", want: "map discipline: quantified obligations hold at each site and across all tickets"},
		{file: discipline, section: section, needle: "The changed-function caller sweep runs `bench consumers` for each changed function, unexported functions included.", want: "map discipline: the caller sweep runs bench consumers on unexported functions too"},
		{file: discipline, section: section, needle: "`Workflow-step writes`, for each new workflow step, traces every write through the digests that later checkpoints compare.", want: "map discipline: a new workflow step traces its writes through checkpoint digests"},
	}}.check(t)
}
