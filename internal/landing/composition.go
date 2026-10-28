package landing

import "github.com/gibbonmi/bench/internal/landing/published"

// The composition owner is published, which imports no gate, so the gate's background
// landing grade composes the same tree that a landing composes. These names keep each
// landing caller unchanged.
type (
	// CompositionRequest identifies two immutable commits to merge without checkout state.
	CompositionRequest = published.CompositionRequest
	// CompositionResult is either a prospective tree or one bounded conflict kind.
	CompositionResult = published.CompositionResult
	// Conflict describes why Git could not produce one prospective tree.
	Conflict = published.Conflict
)

// ConflictError is the refusal a conflicted reviewed landing returns. Its message is the
// message of published.ConflictError.
type ConflictError struct{ Conflict }

func (e ConflictError) Error() string { return published.ConflictError(e).Error() }

// Compose performs Git's three-way tree merge using the repository's real merge base.
// ReviewBase is metadata only and is never used as the merge base.
func (o Owner) Compose(r CompositionRequest) (CompositionResult, error) {
	return published.Compose(r)
}

func compositionCommit(root, value, role string) (string, error) {
	return published.CompositionCommit(root, value, role)
}
