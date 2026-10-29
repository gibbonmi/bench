package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// redKeyWitness names a file in the Git common directory that turns the gate of
// gradeKeyFixture red. Every checkout of the repository reads it, and no tree holds it.
const redKeyWitness outcomeWitness = ".gate-red-key"

// reuseLine is the line that a reused green prints.
const reuseLine = "gate: green (fresh verdict reused for this tree)\n"

// gradeKeyFixture answers an outcome fixture, a linked worktree of its tip, and the tip's
// tree. The fixture gate copies the gate cache of the common directory, so the primary
// checkout holds a green of another tree before a run of the linked worktree.
func gradeKeyFixture(t *testing.T) (primary, linked, tree string) {
	t.Helper()
	primary = outcomeFixture(t, func(f *testrepo.GateFixture, body string) string {
		return body + `if [ -e "$common/` + string(redKeyWitness) + `" ]; then exit 7; fi
`
	})
	linked = filepath.Join(t.TempDir(), "linked checkout")
	gittest.MustRun(t, primary, "worktree", "add", "-q", "--detach", linked)
	writeInspectRecord(t, primary, inspectFullRecord(strings.Repeat("0", 40), strings.Repeat("0", 64), time.Now().UTC()))
	return primary, linked, gittest.Output(t, primary, "rev-parse", "HEAD^{tree}")
}

// gradeKeyRun runs the prospective gate of tree from root in mode, requires action exit want,
// and answers the standard output.
func gradeKeyRun(t *testing.T, root, tree string, mode runMode, want int) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if result := executeTreeWithOwner(context.Background(), root, tree, &stdout, &stderr, nil, mode, nil); result.ActionExit != want {
		t.Fatalf("prospective run from %s = %#v, want action exit %d; stderr=%q", root, result, want, stderr.String())
	}
	return stdout.String()
}

// requireRuns requires that the fixture gate of root ran want times.
func requireRuns(t *testing.T, root string, want int, when string) {
	t.Helper()
	if got := outcomeRuns(t, root); got != want {
		t.Fatalf("runs %s = %d, want %d", when, got, want)
	}
}

// GR10, GR11: a green of one tree in a linked worktree answers a run of that tree from the
// primary checkout, which prints the reuse line and runs no phase.
func TestGradeKeyCrossesCheckouts(t *testing.T) {
	primary, linked, tree := gradeKeyFixture(t)
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)

	if got := gradeKeyRun(t, primary, tree, reuseFreshGreen, 0); got != reuseLine {
		t.Fatalf("primary stdout = %q, want the reuse line alone", got)
	}
	requireRuns(t, primary, 1, "after a primary run of the linked worktree's green tree")
}

// GR14: a green under the relaxed posture does not answer a strict run of the same tree.
func TestGradeKeyBindsThePosture(t *testing.T) {
	root, _, tree := gradeKeyFixture(t)
	t.Setenv(requireCapabilitiesEnv, "0")
	gradeKeyRun(t, root, tree, reuseFreshGreen, 0)
	t.Setenv(requireCapabilitiesEnv, "1")
	gradeKeyRun(t, root, tree, reuseFreshGreen, 0)
	requireRuns(t, root, 2, "after a strict run of a relaxed green")
}

// GR15: a primary checkout whose baseline phase manifest differs runs the gate again.
func TestGradeKeyBindsTheBaselineSchedule(t *testing.T) {
	primary, linked, tree := gradeKeyFixture(t)
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)
	writeGateFixtureFile(t, primary, canary.PhaseManifestPath, `{"phases":[{"name":"kept","argv":["true"]}]}`, 0o644)
	gradeKeyRun(t, primary, tree, reuseFreshGreen, 0)
	requireRuns(t, primary, 2, "after a run under another baseline schedule")
}

// GR16: the run purpose is not in the grade key, so an unpurposed green of the published tree
// answers its complete checkpoint. The completion evidence pass still runs first, so a source
// commit that no review covers refuses with no run.
func TestCompleteCheckpointReusesAnUnpurposedGreenOfItsTree(t *testing.T) {
	publishedGreen := func(t *testing.T, f *recordtest.Fixture) {
		t.Helper()
		gradeKeyRun(t, f.Root, publishedTree(t, f), reuseFreshGreen, 0)
	}
	f := closureSource(t, false)
	publishedGreen(t, f)
	if code, out := completeCheckpoint(t, f); code != 0 || !strings.Contains(out, reuseLine) {
		t.Fatalf("complete checkpoint after an unpurposed green = (%d, %q), want the reuse line", code, out)
	}
	requireRuns(t, f.Root, 1, "after the complete checkpoint")

	stale := closureSource(t, false)
	correctCommentGap(stale)
	publishedGreen(t, stale)
	if code, out := completeCheckpoint(t, stale); code == 0 || !strings.Contains(out, "completion evidence: ") {
		t.Fatalf("complete checkpoint of a stale completion = (%d, %q), want a completion evidence refusal", code, out)
	}
	requireRuns(t, stale.Root, 1, "after the refused complete checkpoint")
}

// GR17: a green answers only inside the window of its own recorded time.
func TestGradeKeyKeepsTheSourceWindow(t *testing.T) {
	primary, linked, tree := gradeKeyFixture(t)
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)
	expireEvidence(t, primary)
	gradeKeyRun(t, primary, tree, reuseFreshGreen, 0)
	requireRuns(t, primary, 2, "after a primary run of an expired green")
}

// GR18: a red run of a key removes the green that another checkout recorded under it.
func TestGradeKeyRedRemovesAnotherCheckoutsGreen(t *testing.T) {
	primary, linked, tree := gradeKeyFixture(t)
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)
	red := redKeyWitness.path(t, primary)
	if err := os.WriteFile(red, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	gradeKeyRun(t, primary, tree, forceRun, 7)
	if err := os.Remove(red); err != nil {
		t.Fatal(err)
	}
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)
	requireRuns(t, primary, 3, "after a green, a red of its key, and a linked-worktree run")
}

// GR20: a record whose grade field differs from its file name is changed evidence.
func TestGradeKeyRefusesAMisnamedRecord(t *testing.T) {
	primary, linked, tree := gradeKeyFixture(t)
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)
	common, err := benchgit.CommonDir(primary)
	if err != nil {
		t.Fatal(err)
	}
	store := filepath.Dir(evidencePath(common, subject{}))
	entries, err := os.ReadDir(store)
	if err != nil || len(entries) != 1 {
		t.Fatalf("evidence store = %v %v, want one record", entries, err)
	}
	var record verdictRecord
	if err := json.Unmarshal(outcomeRead(t, filepath.Join(store, entries[0].Name())), &record); err != nil {
		t.Fatal(err)
	}
	record.Grade = strings.Repeat("0", 64)
	if err := durableReplaceAt(durablefile.Native(), store, entries[0].Name(), record); err != nil {
		t.Fatal(err)
	}
	gradeKeyRun(t, primary, tree, reuseFreshGreen, 0)
	requireRuns(t, primary, 2, "after a primary run beside a misnamed record")
}

// A malformed evidence record never answers a lookup. The well-formed row proves that the
// record bytes otherwise load as a ready green.
func TestMalformedEvidenceAnswersNoLookup(t *testing.T) {
	now := time.Now().UTC()
	at := now.Truncate(time.Second).Format(time.RFC3339)
	key, tree := strings.Repeat("a", 64), strings.Repeat("b", 40)
	record := func(schema, tree, keyField string) string {
		return `{` + keyField + `"recorded_at":"` + at + `","schema":` + schema + `,"state":"ready","status":"green","tree":"` + tree + `"}`
	}
	grade := `"grade":"` + key + `",`
	for _, row := range []struct {
		name, body string
		mode       os.FileMode
		want       State
	}{
		{"well formed", record("1", tree, grade), 0o600, Ready},
		{"missing grade", record("1", tree, ""), 0o600, Invalid},
		{"legacy oracle", record("1", tree, `"oracle":"`+key+`",`), 0o600, Invalid},
		{"duplicate key", record("1", tree, grade+grade), 0o600, Invalid},
		{"wrong schema", record("2", tree, grade), 0o600, Invalid},
		{"wrong tree", record("1", "not-a-tree", grade), 0o600, Invalid},
		{"wrong mode", record("1", tree, grade), 0o644, Invalid},
	} {
		t.Run(row.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), key)
			if err := os.WriteFile(path, []byte(row.body), row.mode); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(path, row.mode); err != nil {
				t.Fatal(err)
			}
			if got := loadEvidence(path, now); got.state != row.want {
				t.Fatalf("loadEvidence(%s) = %#v, want state %q", row.body, got, row.want)
			}
		})
	}
}
