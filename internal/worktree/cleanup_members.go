// The one ordered traversal every cleanup apply mode runs its members through. It owns the
// visit order, the pass-through of a member the mode does not visit, the accumulation of
// completed outcomes, and the stop at the first failed visit. Each mode adapter keeps what is
// its own: the preflight before the traversal, the predicate that selects a member for a
// visit, the qualification and effect a visit runs, and the rows a stopped apply reports.
//
// The traversal derives no fingerprint, ownership proof, classification, or assignment
// identity, and it takes no run: a mode's visit closure captures the run it answers to.

package worktree

// cleanupTraversal is one mode's adapter to the shared traversal over members of type M.
type cleanupTraversal[M any] struct {
	// visits selects a member for a visit. A member it does not select passes through.
	visits func(M) bool
	// retained is the row a member that is not visited reports, exactly as planned.
	retained func(M) CleanupPlan
	// visit runs the mode's qualification and effect on one reached member. It returns the
	// member's current outcome, and on failure the error that stops the traversal.
	visit func(M) (CleanupPlan, error)
	// stopped projects the rows a stopped apply reports. It receives every row so far, the
	// failed member's outcome last, and the index of the failed member in members.
	stopped func(plans []CleanupPlan, failed int) []CleanupPlan
}

// run visits members in order after the rows plans already holds. The first failed visit
// stops it, and the mode's projection of the stopped rows is the result.
func (c cleanupTraversal[M]) run(members []M, plans []CleanupPlan) ([]CleanupPlan, error) {
	for i, member := range members {
		if !c.visits(member) {
			plans = append(plans, c.retained(member))
			continue
		}
		outcome, err := c.visit(member)
		plans = append(plans, outcome)
		if err != nil {
			return c.stopped(plans, i), err
		}
	}
	return plans, nil
}
