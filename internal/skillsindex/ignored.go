package skillsindex

import (
	"sort"

	"github.com/gibbonmi/bench/internal/visibletree"
)

// sortedVisibleSkills sorts the skill directory names and drops the ones git's standard
// excludes ignore. An ignored skill is a local install no commit carries, so it neither
// renders a line nor earns a diagnostic. If the lookup fails, every name stays: a lookup
// fault must not hide a skill from grading.
func sortedVisibleSkills(root string, names []string) []string {
	sort.Strings(names)
	rels := make([]string, len(names))
	for i, name := range names {
		rels[i] = ".agents/skills/" + name
	}
	ignored, err := visibletree.Ignored(root, rels...)
	if err != nil || len(ignored) == 0 {
		return names
	}
	kept := names[:0:0]
	for i, name := range names {
		if !ignored[rels[i]] {
			kept = append(kept, name)
		}
	}
	return kept
}
