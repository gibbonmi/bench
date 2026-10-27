package gate

import "context"

// checkoutBaselineKey selects the private checkout of a prospective run as its baseline root.
type checkoutBaselineKey struct{}

// withCheckoutBaseline marks a run whose baseline root is its private checkout. A chunk
// checkpoint grades the HEAD commit of its worktree without the review record, so that checkout
// holds the committed baseline members. An edit or a commit in the worktree during the run then
// moves neither member.
func withCheckoutBaseline(ctx context.Context) context.Context {
	return context.WithValue(ctx, checkoutBaselineKey{}, true)
}

// baselineRoot answers the root whose baseline schedule and baseline runner a prospective run of
// checkout for root reads: checkout when ctx selects it, and root otherwise.
func baselineRoot(ctx context.Context, checkout, root string) string {
	if selected, _ := ctx.Value(checkoutBaselineKey{}).(bool); selected {
		return checkout
	}
	return root
}
