package refusalroute

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
}
