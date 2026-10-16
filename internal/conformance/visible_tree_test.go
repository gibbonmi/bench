package conformance

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/visibletree"
)

// withoutIgnoredEntries drops the children of the guidance tree rel that git's standard
// excludes ignore. Both budget checks list their subjects through listingTreeEntries,
// which calls this, so an ignored local skill is a subject for neither.
func withoutIgnoredEntries(root, rel string, found []os.DirEntry) ([]os.DirEntry, error) {
	children := make([]string, len(found))
	for i, entry := range found {
		children[i] = path.Join(rel, entry.Name())
	}
	ignored, err := visibletree.Ignored(root, children...)
	if err != nil {
		return nil, err
	}
	kept := found[:0:0]
	for i, entry := range found {
		if !ignored[children[i]] {
			kept = append(kept, entry)
		}
	}
	return kept, nil
}

// TestGuidanceProseBudgetsGradeOnlyTheVisibleTree is the bite proof for the shared skills
// enumeration's scope. An oversized skill under an ignored path is a local install no
// commit carries, so it is not a subject. The same skill at a visible untracked path is.
func TestGuidanceProseBudgetsGradeOnlyTheVisibleTree(t *testing.T) {
	files := healthyProseBudgetFiles()
	files[proseBudgetProfile] = proseBudgetTable(proseBudgetHeader, proseBudgetRows...)
	files[".gitignore"] = ".agents/skills/local-skill/\n"
	files[".agents/skills/local-skill/SKILL.md"] = proseBudgetLines(121, true)
	root := throwawayRoot{gitInit: true, files: files}.build(t)
	if diags := checkGuidanceProseBudgets(root); len(diags) != 0 {
		t.Fatalf("an ignored oversized skill got diagnostics:\n%s", strings.Join(diags, "\n"))
	}

	files[".gitignore"] = "\n"
	root = throwawayRoot{gitInit: true, files: files}.build(t)
	if diags := checkGuidanceProseBudgets(root); !containsDiagnostic(diags, "prose-budget exceeded: .agents/skills/local-skill/SKILL.md is 121 lines, over its 120-line budget") {
		t.Fatalf("a visible oversized skill was accepted:\n%s", strings.Join(diags, "\n"))
	}
}
