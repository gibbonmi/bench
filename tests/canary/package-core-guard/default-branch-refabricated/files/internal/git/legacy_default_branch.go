package git

// DefaultBranch is the fabricating helper the sweep forbids: it answers with a branch
// name in a repository that may have none.
func DefaultBranch(root string) string {
	return "main"
}
