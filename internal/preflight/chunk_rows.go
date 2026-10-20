package preflight

import (
	"slices"
	"strings"
)

// chunkRowsName is the verdict row that grades the chunk table against the coverage map.
const chunkRowsName = "chunk-rows"

// chunkRowsCheck grades the Implementation chunks table against the coverage map in both
// directions: every declared row sits in an acceptance cell, and every listed ID is a declared
// row. An absent table or column reds with its own detail. Otherwise the detail names each
// token fault, each unlisted row, and each undeclared ID. A row that two chunks list is green.
func chunkRowsCheck(f Facts) CheckResult {
	table := f.ChunkTable
	if table.Missing != "" {
		return red(chunkRowsName, table.Missing)
	}
	var faults []string
	if len(table.Faults) > 0 {
		faults = append(faults, "chunk table fault(s): "+strings.Join(table.Faults, ", "))
	}
	if unlisted := missingFrom(f.DeclaredRowIDs, table.Listed); len(unlisted) > 0 {
		faults = append(faults, "declared row(s) that no acceptance cell lists: "+strings.Join(unlisted, ", "))
	}
	if undeclared := missingFrom(table.Listed, f.DeclaredRowIDs); len(undeclared) > 0 {
		faults = append(faults, "acceptance cell ID(s) that the map does not declare: "+strings.Join(undeclared, ", "))
	}
	if len(faults) > 0 {
		return red(chunkRowsName, strings.Join(faults, "; "))
	}
	return green(chunkRowsName)
}

// missingFrom is each distinct ID of ids that set does not hold, in ids order.
func missingFrom(ids, set []string) []string {
	var missing []string
	for _, id := range ids {
		if !slices.Contains(set, id) && !slices.Contains(missing, id) {
			missing = append(missing, id)
		}
	}
	return missing
}
