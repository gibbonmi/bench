package preflight

import (
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/tickets"
)

// kitPinCheck grades the kit pin into the ticket. A system-tagged test file reads
// BENCH_KIT, so an ambient value flips the fixture verdict under composition
// unless the ticket states the variable.
func kitPinCheck(f Facts) CheckResult {
	var unpinned []string
	for _, ticket := range f.Tickets {
		if f.TicketPinsKit[ticket.Name] {
			continue
		}
		for _, entry := range ticket.Writes {
			if f.WritesSystemTagged[entry] {
				unpinned = append(unpinned, ticket.Name+": "+entry)
			}
		}
	}
	if len(unpinned) > 0 {
		return red("kit-pin", "ticket writes a system-tagged test file without stating BENCH_KIT: "+strings.Join(unpinned, ", "))
	}
	return green("kit-pin")
}

// coversTokens is every row ID the parsed tickets cite, in enumeration order.
// It is the one ownership evidence both row checks read.
func coversTokens(f Facts) []string {
	var tokens []string
	for _, ticket := range f.Tickets {
		tokens = append(tokens, ticket.Covers...)
	}
	return tokens
}

func rowsOwnedCheck(f Facts) CheckResult {
	cited := coversTokens(f)
	var uncited []string
	for _, id := range f.DeclaredRowIDs {
		if !slices.Contains(cited, id) {
			uncited = append(uncited, id)
		}
	}
	if len(uncited) > 0 {
		return red("rows-owned", "declared row(s) cited by no ticket file: "+strings.Join(uncited, ", "))
	}
	return green("rows-owned")
}

func rowsMembershipCheck(f Facts) CheckResult {
	seen := map[string]bool{}
	var phantom []string
	for _, tok := range coversTokens(f) {
		if tickets.TagOf(tok) != f.SpecTag {
			continue
		}
		if slices.Contains(f.DeclaredRowIDs, tok) {
			continue
		}
		if !seen[tok] {
			seen[tok] = true
			phantom = append(phantom, tok)
		}
	}
	if len(phantom) > 0 {
		return red("rows-membership", "ticket token(s) under this spec's tag name no declared row: "+strings.Join(phantom, ", "))
	}
	return green("rows-membership")
}
