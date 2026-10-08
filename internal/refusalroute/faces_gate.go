package refusalroute

// gateRerun re-runs the caller's own gate in the worktree that the label fact addresses.
// The raising site composes the arguments from the run it refuses, and a face whose
// recovery forces a fresh run gets them with --fresh.
var gateRerun = TreeCommand("bench gate", Fact(FactLabel), Composed(FactArguments))

// gateFaces are the gate checkpoint's refusal faces, in registry order. The raising site
// supplies each sentence: the checkpoint names the cause it observed.
var gateFaces = []Face{
	{
		// The checkpoint cannot show what the completion evidence lacks, and one read
		// reports the plan row and the record state. The slug comes from the spec-path
		// grammar's one owner.
		Verb:      Gate,
		Name:      "checkpoint-completion-evidence",
		Authority: Agent,
		Route:     []Step{Command(Text("bench preflight review"), Fact(FactSlug))},
	},
	{
		// The landing publishes only committed bytes, so the caller commits its own
		// uncommitted work and grades the committed source again.
		Verb:      Gate,
		Name:      "checkpoint-dirty-checkout",
		Authority: Agent,
		Route:     []Step{commitAt(FactLabel), gateRerun},
	},
	{
		// The landing's closure transform refuses the spec: its status or its delivery
		// closure is the reviewer's to change.
		Verb:      Gate,
		Name:      "checkpoint-composition",
		Authority: Reviewer,
		Route:     []Step{Instruction(Text("the delivery closure of the spec does not compose; hand back"))},
	},
	{
		// The source tip moved between the accepted subject and the run, so the rerun
		// grades the moved tip.
		Verb:      Gate,
		Name:      "checkpoint-tip-moved",
		Authority: Agent,
		Route:     []Step{gateRerun},
	},
	{
		// A fault captured no subject. The diagnosis comes first, and the rerun does not
		// reuse a verdict.
		Verb:      Gate,
		Name:      "checkpoint-subject-unavailable",
		Authority: Agent,
		Route:     []Step{doctor, gateRerun},
	},
	{
		// The graded completion is not the exact transform of its reviewed source, which
		// no repair in the caller's worktree reaches.
		Verb:      Gate,
		Name:      "gate-handback",
		Authority: Reviewer,
		Route:     []Step{clearCause, gateRerun},
	},
}
