// Package inventorytest renders the planning-facts inventories that the preflight fixtures
// declare. Only tests import it. It imports no Bench package but the block reader, a
// standard-library leaf: the shared preflight fixture imports it, and the tests of packages
// that ownership itself imports use that fixture.
package inventorytest

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/markdown"
)

// Manifest is one version 1 inventory of completeness that holds each given item in its
// array, in order.
func Manifest(completeness string, premises, seams, decisions, manual []string) string {
	return fmt.Sprintf(`{"version":1,"authority":"authored-spec","completeness":%q,"baseline":"source-base","scope":[],"changes":[],"premises":[%s],"seams":[%s],"decisions":[%s],"manual_claims":[%s]}`,
		completeness, strings.Join(premises, ","), strings.Join(seams, ","), strings.Join(decisions, ","), strings.Join(manual, ","))
}

// CompleteEmpty is the declared-complete inventory with no item.
var CompleteEmpty = Manifest("declared-complete", nil, nil, nil, nil)

// Section is one planning-facts section of a spec that holds payload.
func Section(payload string) string {
	return "\n## Planning claim inventory\n\n" + markdown.FenceMarker + "bench-planning-facts\n" + payload + "\n" + markdown.FenceMarker + "\n"
}
