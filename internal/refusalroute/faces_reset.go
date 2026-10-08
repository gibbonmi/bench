package refusalroute

// resetFaces are the reset's refusal faces, in registry order.
var resetFaces = []Face{
	{
		// A conflicted index has no recoverable envelope, so the clean plans the checkout's
		// retirement. The id is the ledger's own, and the route prints it as the ledger holds it.
		Verb:      Reset,
		Name:      "reset-checkout-conflicted",
		Sentence:  "checkout is conflicted",
		Authority: Agent,
		Route:     []Step{Command(Text("bench worktree clean"), Composed(FactAssignmentID))},
	},
	{
		// The checkout moved after the plan, so the plan of the current checkout prints the
		// apply that it takes.
		Verb:      Reset,
		Name:      "reset-plan-stale",
		Sentence:  "reset plan is stale",
		Authority: Agent,
		Route:     []Step{Command(Composed(FactResetPlan))},
	},
	{
		// The tree is gone, so the route is the one that clears the record: the clean of the
		// landed assignments or the release of this one.
		Verb:      Reset,
		Name:      "reset-tree-missing",
		Sentence:  "worktree tree is missing",
		Authority: Agent,
		Route:     []Step{Command(Composed(FactRecovery))},
	},
	{
		Verb:      Reset,
		Name:      "reset-handback",
		Authority: Reviewer,
		Route:     handback,
	},
}
