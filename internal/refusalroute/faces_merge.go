package refusalroute

// mergeFaces are the merge's refusal faces, in registry order.
var mergeFaces = []Face{
	{
		// The target tip alone fails its grade, so the repair is the target's own. The commit
		// grades the repaired tree, and the fold then composes onto a green target. The
		// authorization policy owns the sentence.
		Verb:      Merge,
		Name:      "merge-target-red",
		Authority: Agent,
		Route: []Step{
			Instruction(Text("repair each failing check in"), Fact(FactLabel)),
			commitAt(FactLabel),
			rerun,
		},
	},
	{
		// The target tip alone grades green, so the fold adds the red. FT342 owns the
		// decision of who repairs it. The authorization policy owns the sentence.
		Verb:      Merge,
		Name:      "merge-fold-red",
		Authority: Reviewer,
		Route: []Step{
			Instruction(Text("the fold of"), Fact(FactFrom), Text("adds a red to"), Fact(FactLabel)),
			Instruction(Text("decide who repairs a red that the fold of a moved main adds, the open FT342 question")),
		},
	},
	{
		// The reconcile after the publication resets the target checkout, so its own
		// uncommitted work is committed first.
		Verb:      Merge,
		Name:      "merge-target-not-clean",
		Sentence:  "merge target checkout is not clean",
		Authority: Agent,
		Route:     []Step{commitAt(FactLabel), rerun},
	},
	{
		// A sibling contributes its committed branch tip alone, so its uncommitted work is
		// committed in the sibling before the fold.
		Verb:      Merge,
		Name:      "merge-sibling-not-clean",
		Sentence:  "sibling checkout is not clean",
		Authority: Agent,
		Route:     []Step{commitAt(FactSiblingLabel), rerun},
	},
	{
		// The composition names the conflicted paths in its own sentence, so the face
		// declares none.
		Verb:      Merge,
		Name:      "merge-conflict",
		Authority: Reviewer,
		Route:     append(append([]Step{}, handMerge...), rerun),
	},
	{
		// The authorization policy owns the sentence.
		Verb:      Merge,
		Name:      "merge-infrastructure",
		Authority: Agent,
		Route:     []Step{doctor, rerun},
	},
	{
		// The merge published, and the target checkout did not reconcile. The reset verb's
		// plan at the published tip reconciles it under a preserved envelope.
		Verb:      Merge,
		Name:      "merge-published-unreconciled",
		Authority: Agent,
		Route:     []Step{Command(Composed(FactResetPlan))},
	},
	{
		Verb:      Merge,
		Name:      "merge-handback",
		Authority: Reviewer,
		Route:     handback,
	},
}
