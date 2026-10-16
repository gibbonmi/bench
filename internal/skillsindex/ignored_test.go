package skillsindex

import (
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// TestEntriesSkipIgnoredSkills is the bite proof for the enumeration's scope. A skill
// under an ignored path is a local install no commit carries, so it renders no line and
// earns no missing-trigger diagnostic. The same skill at a visible untracked path earns
// that diagnostic.
func TestEntriesSkipIgnoredSkills(t *testing.T) {
	root := gittest.RepoIgnoring(t, ".agents/skills/local-skill/")
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
	writeFile(t, root, ".agents/skills/local-skill/SKILL.md", "---\nname: local-skill\n---\n")
	writeFile(t, root, referenceRel, reference("- doing alpha things → `.agents/skills/alpha/SKILL.md`\n"))
	if diags := Check(root); len(diags) != 0 {
		t.Fatalf("ignored skill: diagnostics %v, want none", diags)
	}

	writeFile(t, root, ".gitignore", "\n")
	want := "skill 'local-skill' missing index: frontmatter (the skills index is generated)"
	if diags := Check(root); len(diags) != 1 || diags[0] != want {
		t.Fatalf("visible skill: diagnostics %v, want only %q", diags, want)
	}
}
