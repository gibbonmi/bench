// The worktree merge verb's declared lane: the fixture that commits a manifest lane, and
// the caller-root row, where the lane grades the composed tree whatever directory the
// caller runs the verb from.
package worktree

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// mergeSetAt commits a manifest lane that appends to a tally file on root's checked-out
// branch, and answers the merge set over root. The merge resolves its lane from the
// target's manifest only under a kit apart from the target, so the set holds its own kit.
// A lane that ran leaves a byte in the tally, and one that did not leaves no file.
func mergeSetAt(t *testing.T, root, home string) mergeSet {
	t.Helper()
	tally := filepath.Join(t.TempDir(), "lane-tally")
	commitLaneManifest(t, root, gate.Phase{Name: "unit", Argv: []string{"sh", "-c", "printf g >> " + sanitize.ShellQuote(tally)}})
	return mergeSet{repoHome: repoHome{root: root, home: home}, joins: defaultJoins(), kit: t.TempDir(), tally: tally}
}

// commitLaneManifest commits a phase manifest at dir whose lane declares check alone.
func commitLaneManifest(t *testing.T, dir string, check gate.Phase) {
	t.Helper()
	body, err := json.Marshal(map[string]any{"lane": []map[string]any{{"name": check.Name, "argv": check.Argv}}})
	mustNoError(t, err)
	mustMkdirAll(t, filepath.Join(dir, ".bench"), 0o755)
	commitInWorktree(t, dir, ".bench/phases.json", string(body), "declare the lane")
}

// FT335: a merge run from the primary checkout grades the prose the sibling brings in
// from the composed tree. The lane resolves against the target worktree, so its prose
// check names the target as its file root, as the built-in lane anchors `gate-prose`.
// An authority rooted at the caller's checkout leaves that anchor unmoved, and the check
// reads the incoming Markdown from the target checkout, where the merge has not put it.
func TestMergeGradesIncomingProseFromTheComposedTreeWhateverTheCallerRoot(t *testing.T) {
	t.Parallel()
	f := mergeFixture(t, "integration", "sibling")
	target, sibling := f.created[0], f.created[1]
	commitInWorktree(t, target.Path, "target-only.md", "target prose\n", "target prose")
	commitInWorktree(t, sibling.Path, "sibling-only.md", "sibling prose\n", "sibling prose")
	commitLaneManifest(t, target.Path, gate.Phase{Name: "prose", Argv: []string{"sh", "-c",
		`anchor=$1; shift; for path; do test -f "$anchor/$path" || { echo "unreadable $path"; exit 1; }; done`,
		"prose", target.Path, gate.LaneNamedMarkdownToken}})
	if f.root == target.Path {
		t.Fatal("the fixture runs the merge from the target, so the caller root never differs")
	}

	r := runVerb(t, verbMerge, f.merge("--from", sibling.Assignment.Label, target.Assignment.ID))
	if r.exit != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	mergedRecord(t, r.stdout)
}
