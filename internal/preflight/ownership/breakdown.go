package ownership

import (
	"slices"

	"github.com/gibbonmi/bench/internal/tickets"
)

// UncitedIDs is every row ID of ids that no parsed ticket's `Covers:` line cites, in ids
// order. A row ID in ticket prose cites nothing. The rows-owned check reads this one rule
// over the one declared row-ID set.
func UncitedIDs(ids []string, parsed []tickets.Ticket) []string {
	var open []string
	for _, id := range ids {
		if !slices.ContainsFunc(parsed, func(ticket tickets.Ticket) bool { return slices.Contains(ticket.Covers, id) }) {
			open = append(open, id)
		}
	}
	return open
}
