package anchors

// debugLoopAnchors refuse each retired debug-loop form in the debug phase and the
// mandatory executable-red forms in the spec skill, and they keep the delegation skill's
// route to the debug integration.
var debugLoopAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-debug.md", Kind: ForbidInSection, Section: "How it meets the rest of Bench", Needle: "The reviewer invokes this phase; a write delegate never charges it.", Diagnostic: "debug loop: DG3 forbids the blanket write-delegate debug ban"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-debug.md", Kind: ForbidInSection, Section: "How it meets the rest of Bench", Needle: "It discards its in-fence work from its owned worktree.", Diagnostic: "debug loop: DG5 forbids discarded dirty work"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-debug.md", Kind: ForbidInSection, Section: "How it meets the rest of Bench", Needle: "When a repro proves an out-of-fence cause, the delegate continues implementation edits.", Diagnostic: "debug loop: DG5 forbids continued out-of-fence implementation"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-debug.md", Kind: ForbidInSection, Section: "How it meets the rest of Bench", Needle: "For an out-of-fence defect, hand off before running Phases 1 through 3.", Diagnostic: "debug loop: DG4 forbids handoff before Phases 1 through 3"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-debug.md", Kind: ForbidInSection, Section: "How it meets the rest of Bench", Needle: "For an in-fence defect, the write delegate must not run debug through Phase 6.", Diagnostic: "debug loop: DG3 forbids refusal of in-fence Phase 6"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-debug.md", Kind: ForbidInSection, Section: "How it meets the rest of Bench", Needle: "Delegate implementation and repair work.", Diagnostic: "debug loop: DG2 forbids diagnostic write delegation"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-debug.md", Kind: ForbidInSection, Section: "How it meets the rest of Bench", Needle: "The session that owns the debug loop delegates its in-scope repair.", Diagnostic: "debug loop: DG1 forbids transferred repair authorship"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: RequireInSection, Section: "The charge", Needle: "[debug integration](../../commands/bench-debug.md#how-it-meets-the-rest-of-bench)", Diagnostic: "debug loop: delegation requires the debug integration reference"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: Forbid, Needle: "Require an executable red before you specify a new feature.", Diagnostic: "debug loop: DG15 forbids an executable red requirement for new-feature specification"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: ForbidCaseFoldedEmphasis, Needle: "executable red is mandatory", Diagnostic: "debug loop: DG15 forbids mandatory executable-red variants"},
}
