package structure

import (
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// ValidateAcceptGrants grades the accept list of a live tree: every reviewer grant in
// <root>/.bench/structure-accept must still name a scanned subject and must still carry
// its reason. `bench structure` prints those two conditions as report text and leaves the
// exit code alone, so an ambient run is never reddened by list hygiene; the gate's
// conformance registry runs this against the graded root, where a grant that outlived its
// file or lost its reason reds the oracle.
//
// The parse, the staleness rule, and the diagnostic wording come from the same
// loadAccepts/filterSources/staleAcceptWarnings the report uses.
//
// An absent or grant-free accept file is silence, not a finding: the registry grades roots
// that are not the kit, and a tree with no grants has nothing to keep current. Only once a
// grant exists does the tracked-source query run, and a failure of that query is loud —
// unresolvable grants are unknown, never assumed live.
func ValidateAcceptGrants(root string) []string {
	accepts, warnings, err := loadAccepts(filepath.Join(root, ".bench", "structure-accept"))
	if err != nil {
		return []string{"structure-accept: present but unreadable: " + err.Error()}
	}
	if len(accepts) == 0 && len(warnings) == 0 {
		return nil
	}
	out, queryErr := git.Output("-C", root, "ls-files")
	if queryErr != nil {
		return append(warnings, "structure-accept: cannot resolve grants against the tracked source tree: "+queryErr.Error())
	}
	return append(warnings, staleAcceptWarnings(accepts, filterSources(strings.Split(out, "\n")))...)
}
