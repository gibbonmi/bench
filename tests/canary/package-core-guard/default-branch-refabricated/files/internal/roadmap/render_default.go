package roadmap

// repoFacts stands in for git.RepoFacts, whose DefaultBranch field survives the deletion.
type repoFacts struct {
	DefaultBranch string
}

// renderDefault reads the surviving struct field, which the sweep must not report.
func renderDefault(facts repoFacts) string {
	return facts.DefaultBranch
}
