package ownership

import "github.com/gibbonmi/bench/internal/tickets"

// Unordered names each other parsed ticket whose owned paths overlap proposed while no
// blocker chain orders it against selected, in parsed order. A proposal grants no order, so
// each named ticket needs an approved blocker edge before the overlap can dispatch.
func Unordered(all []tickets.Ticket, selected tickets.Ticket, proposed []string) []string {
	var names []string
	for _, other := range all {
		if other.Name == selected.Name || ordered(selected, other, all) || !overlaps(proposed, OwnedPaths(other)) {
			continue
		}
		names = append(names, other.Name)
	}
	return names
}

func overlaps(left, right []string) bool {
	for _, a := range left {
		for _, b := range right {
			if tickets.Covers(b, a) || tickets.Covers(a, b) {
				return true
			}
		}
	}
	return false
}

func ordered(left, right tickets.Ticket, all []tickets.Ticket) bool {
	return blocks(left.Name, right.Name, all, map[string]bool{}) || blocks(right.Name, left.Name, all, map[string]bool{})
}

// blocks reports whether a blocker chain leads from the ticket from to the ticket target.
func blocks(from, target string, all []tickets.Ticket, seen map[string]bool) bool {
	if seen[from] {
		return false
	}
	seen[from] = true
	for _, ticket := range all {
		if ticket.Name != from {
			continue
		}
		for _, blocker := range ticket.Blockers {
			if blocker == target || blocks(blocker, target, all, seen) {
				return true
			}
		}
	}
	return false
}
