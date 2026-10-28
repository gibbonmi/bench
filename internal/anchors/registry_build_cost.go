package anchors

// buildCostAnchors forbid the retired run of each verification requirement after a ticket
// commit and the retired check-floor row in the build guidance and the completion-plan
// reference. They also forbid a second copy of a cost budget in the drain guidance, because
// the budget table of the cost renderer is the one source. The guidance rows require the
// `--in` form of the evidence check, the worktree probe, and the batch gate, and they forbid
// the retired exec forms and the learning for each `Writes:` widen.
var buildCostAnchors = []Anchor{
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "then runs each verification requirement once at that commit", Diagnostic: "build cost: the implement phase restored the run of each verification requirement"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "check-floor", Diagnostic: "build cost: the implement phase restored the retired check-floor row"},
	{Group: AfterImplementSpec, File: completionPlanReference, Kind: Forbid, Needle: "check-floor", Diagnostic: "build cost: the completion-plan reference restored the retired check-floor row"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "60 min", Diagnostic: "build cost: the drain repeats the wall budget of the cost renderer"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "7.1", Diagnostic: "build cost: the drain repeats the bookkeeping budget of the cost renderer"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Require, Class: Command, Needle: "Run it as `bench preflight --in <target> evidence <id> --check-current`.", Diagnostic: "build cost: the known issues dropped the --in form of the evidence check"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Forbid, Needle: "names the exec form", Diagnostic: "build cost: the known-issues retire column restored the exec form"},
	{Group: AfterImplementSpec, File: benchReference, Kind: Require, Class: Command, Needle: "For a worktree, run `bench probe --in <target> ...`.", Diagnostic: "build cost: the command notes dropped the --in form of the worktree probe"},
	{Group: AfterImplementSpec, File: benchReference, Kind: Forbid, Needle: "`bench worktree exec <target> -- bench probe ...`", Diagnostic: "build cost: the command notes restored the exec probe form"},
	{Group: AfterImplementSpec, File: synthesisSkill, Kind: Require, Class: Command, Needle: "Take the green verdict from the whole-project gate. That gate is the landing's gate, or `bench gate --in <target>` for a batch that waits for approval.", Diagnostic: "build cost: craft-synthesis dropped the whole-project gate or its --in batch form"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "record every plan or gate expansion with `bench learning`", Diagnostic: "build cost: the operating guide restored the learning for each Writes: widen"},
}

const (
	// completionPlanReference is the craft-spec reference that states the completion plan schema.
	completionPlanReference = ".agents/skills/bench-craft-spec/references/completion-plan.md"
	// synthesisSkill is the skill that folds a candidate change into the kit.
	synthesisSkill = ".agents/skills/bench-craft-synthesis/SKILL.md"
)
