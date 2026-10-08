package refusalroute

// commitRerun re-runs the caller's own commit in the worktree that the label fact
// addresses. The raising site composes the arguments, because only it holds how many paths
// the caller named.
var commitRerun = TreeCommand("bench commit", Fact(FactLabel), Composed(FactArguments))

// commitFaces are the commit's refusal faces, in registry order.
var commitFaces = []Face{
	{
		// The commit published, and its checkout did not reconcile. The reset verb's plan at
		// the published commit keeps the dirty layer in a recoverable envelope and prints its
		// own apply. The destructive-git guard allows the plan, and it denies a raw restore.
		Verb:      Commit,
		Name:      "commit-published-unreconciled",
		Authority: Agent,
		Route:     []Step{Command(Text("bench worktree reset --to"), Fact(FactPublishedCommit), Fact(FactCheckout))},
	},
	{
		// Within Bench, main receives writes only through landings, so the commit moves to a
		// worktree of the operator's own. The usage package owns the sentence that every
		// write verb prints from the primary checkout, so the face declares none.
		Verb:      Commit,
		Name:      "commit-primary-checkout",
		Authority: Agent,
		Route:     []Step{Command(Text("bench worktree create --request"), Operator("opaque-id"), Text("--label"), Operator("work-item"))},
	},
	{
		// The gate or the lane ran red on the composed tree and printed its failures. The
		// authorization policy owns the sentence.
		Verb:      Commit,
		Name:      "commit-red",
		Authority: Agent,
		Route:     []Step{Instruction(Text("repair each failure that the run reports")), commitRerun},
	},
	{
		// The authorization policy owns the sentence.
		Verb:      Commit,
		Name:      "commit-infrastructure",
		Authority: Agent,
		Route:     []Step{doctor, commitRerun},
	},
	{
		Verb:      Commit,
		Name:      "commit-handback",
		Authority: Reviewer,
		Route:     []Step{clearCause, commitRerun},
	},
}
