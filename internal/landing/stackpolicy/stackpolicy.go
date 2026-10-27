// Package stackpolicy owns the stack owner rule: which members of a landing stack
// change each path, and which of those paths overlap. The caller reads the trees and
// passes their entries here. The package reads no Git and performs no effects. It
// imports only the standard library and the settle policy. The source census test
// enforces that boundary.
package stackpolicy

import "github.com/gibbonmi/bench/internal/landing/settlepolicy"

// Entries maps each path of one tree to one metadata string, such as the mode and the
// object id. The caller passes every tree in one metadata format, and the rule
// compares the whole string, so a mode-only change differs.
type Entries map[string]string

// Owners returns, for each path whose entry differs between destination and a member,
// the 1-based indexes of each member that changes that path, in member order. An entry
// on one side only differs, so a rename owns its old path and its new path.
func Owners(destination Entries, members []Entries) map[string][]int {
	owners := map[string][]int{}
	for i, member := range members {
		for path, metadata := range member {
			if previous, ok := destination[path]; !ok || previous != metadata {
				owners[path] = append(owners[path], i+1)
			}
		}
		for path := range destination {
			if _, ok := member[path]; !ok {
				owners[path] = append(owners[path], i+1)
			}
		}
	}
	return owners
}

// Overlaps returns each path of owners that two or more members change, with those
// member indexes. It leaves out each path the settle policy's capture rule table names,
// because the composition settles that path.
func Overlaps(owners map[string][]int) map[string][]int {
	overlaps := map[string][]int{}
	for path, indexes := range owners {
		if _, _, capture := settlepolicy.CaptureSide(path); len(indexes) >= 2 && !capture {
			overlaps[path] = indexes
		}
	}
	return overlaps
}
