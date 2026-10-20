package worktree

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/usage"
)

// rebuiltMain is the body that a restore row's rebuild publishes as the primary executable.
const rebuiltMain = "rebuilt-main-executable"

// rebuildRecorder is the buildSubject join of a restore row. It publishes rebuiltMain with a
// seal over the tree's own sources and a manifest beside the tree's wrapper, as the primary
// build script does, and counts its calls. A recorder with err fails each call instead.
type rebuildRecorder struct {
	t     *testing.T
	calls int
	err   error
}

func (r *rebuildRecorder) join(_ context.Context, root, executable string) error {
	r.calls++
	if r.err != nil {
		return r.err
	}
	staged := filepath.Join(r.t.TempDir(), "staged-bench")
	mustWrite(r.t, staged, []byte(rebuiltMain), 0o755)
	return freshness.Publish(root, staged, executable, wrapperDir(root), "1.2.3")
}

// restore runs the restore form of the build verb on the spec at slug, with the rebuild and
// the repair recorded.
func restore(t *testing.T, f ownedAssignment, slug string, build *rebuildRecorder, doctor *doctorRecorder) verbResult {
	t.Helper()
	j := defaultJoins()
	j.buildSubject = build.join
	j.doctorChild = doctor.join
	return runVerb(t, verbBuild, f.callWith(j, "--restore", slug))
}

// installedRestore is an install fixture whose spec declares the line declaration, with its
// candidate installed through the install form. The doctor recorder holds the install's call.
func installedRestore(t *testing.T, request, declaration string) (ownedAssignment, *doctorRecorder) {
	t.Helper()
	f := newInstallAssignment(t, request, installShape{inputs: true, declaration: declaration, deliverable: installSpec})
	doctor := &doctorRecorder{}
	r := install(t, f, &buildRecorder{body: "installed-candidate", version: "1.2.3"}, doctor, freshness.PublishedExecutable(f.root))
	requireTest(t, r.exit == 0, "install = (%d, %q), want exit 0", r.exit, r.stderr)
	return f, doctor
}

// putGaps records an open chunk gap for each outcome in the ledger at root.
func putGaps(t *testing.T, root string, outcomes ...string) {
	t.Helper()
	for _, outcome := range outcomes {
		commitmenttest.PutGap(t, root, intent.ChunkGap{Milestone: "fixture", Outcome: outcome, Chunk: "C1", Open: true})
	}
}

// ledgerCommitment is the raw commitment object of the ledger at root.
func ledgerCommitment(t *testing.T, root string) map[string]any {
	t.Helper()
	var ledger map[string]any
	mustNoError(t, json.Unmarshal([]byte(fixtureFileText(t, commitmenttest.LedgerPath(t, root))), &ledger))
	state, ok := ledger["commitment"].(map[string]any)
	requireTest(t, ok, "ledger holds no commitment object: %v", ledger)
	return state
}

// restoreCopies are the ledger copies under the Bench home of f.
func restoreCopies(t *testing.T, f ownedAssignment) []string {
	t.Helper()
	copies, err := filepath.Glob(filepath.Join(f.home, "restore", poolkey.Key(f.root), "*"))
	mustNoError(t, err)
	return copies
}

// TestBuildRestore is CB37 to CB43, CB49, and CB64. After a declared build stops, the restore
// form copies the ledger, clears the declared ledger fields, rebuilds the primary executable
// from the primary sources, and repairs under it. Each refusal comes before any write.
func TestBuildRestore(t *testing.T) {
	t.Parallel()
	const gapsDeclared = "Persisted fields: intent-ledger:commitment.gaps\n"
	rerun := "bench worktree build --restore 'install'"
	install := "next=from the primary checkout, run bench worktree build --install "
	t.Run("help names the restore form", func(t *testing.T) {
		t.Parallel()
		r := runVerb(t, verbBuild, verbCall{root: t.TempDir(), home: t.TempDir(), args: []string{"--help"}})
		requireTest(t, r.exit == 0 && strings.Contains(r.stdout, usage.WorktreeBuildRestore), "help = (%d, %q), want the restore form", r.exit, r.stdout)
	})
	t.Run("the install and restore forms refuse together", func(t *testing.T) {
		t.Parallel()
		f := newInstallAssignment(t, "restore-both-forms", declaredInstall)
		build := &rebuildRecorder{t: t}
		j := defaultJoins()
		j.buildSubject = build.join
		r := runVerb(t, verbBuild, f.callWith(j, "--install", "--restore", "install"))
		requireTest(t, r.exit == 2 && strings.Contains(r.stderr, restoreFlag) && build.calls == 0, "both forms = (%d, %q), build calls %d; want a usage refusal", r.exit, r.stderr, build.calls)
	})
	t.Run("copies the ledger, clears the field, rebuilds, and repairs", func(t *testing.T) {
		t.Parallel()
		f, doctor := installedRestore(t, "restore-top-level", gapsDeclared)
		putGaps(t, f.root, "install")
		before := fixtureFileText(t, commitmenttest.LedgerPath(t, f.root))
		build := &rebuildRecorder{t: t}
		r := restore(t, f, "install", build, doctor)
		requireTest(t, r.exit == 0, "restore = (%d, %q, %q), want exit 0", r.exit, r.stdout, r.stderr)
		copies := restoreCopies(t, f)
		requireTest(t, len(copies) == 1 && fixtureFileText(t, copies[0]) == before && strings.Contains(r.stdout, copies[0]), "copies %v, stdout %q; want one copy of the ledger before the clear", copies, r.stdout)
		state := ledgerCommitment(t, f.root)
		_, gaps := state["gaps"]
		_, bindings := state["bindings"]
		requireTest(t, !gaps && bindings, "commitment = %v, want no gaps key and the bindings kept", state)
		primary := freshness.PublishedExecutable(f.root)
		digest, err := brokermanifest.Digest(primary)
		mustNoError(t, err)
		fields, err := brokermanifest.Read(manifestPath(f.root))
		mustNoError(t, err)
		requireTest(t, build.calls == 1 && freshness.Verify(f.root, primary) == nil && fields["sha256"] == digest, "rebuild calls %d, verify %v, manifest %v, digest %s; want a verified primary that the manifest binds", build.calls, freshness.Verify(f.root, primary), fields, digest)
		requireTest(t, len(doctor.calls) == 2, "doctor calls = %d, want the install's and the restore's", len(doctor.calls))
		installed, restored := doctor.calls[0], doctor.calls[1]
		requireTest(t, restored.dir == installed.dir && slices.Equal(restored.argv, installed.argv) && slices.Equal(restored.env, installed.env) && restored.published == rebuiltMain,
			"restore repair = %+v, want the install's argv, directory, and environment %+v under the rebuilt executable", restored, installed)
	})
	t.Run("clears a nested field from each element", func(t *testing.T) {
		t.Parallel()
		f, doctor := installedRestore(t, "restore-nested", "Persisted fields: intent-ledger:commitment.gaps.chunk\n")
		putGaps(t, f.root, "install", "other")
		r := restore(t, f, "install", &rebuildRecorder{t: t}, doctor)
		requireTest(t, r.exit == 0, "restore = (%d, %q), want exit 0", r.exit, r.stderr)
		gaps, _ := ledgerCommitment(t, f.root)["gaps"].([]any)
		requireTest(t, len(gaps) == 2, "gaps = %v, want both elements kept", gaps)
		for _, element := range gaps {
			gap, _ := element.(map[string]any)
			_, chunk := gap["chunk"]
			_, outcome := gap["outcome"]
			requireTest(t, !chunk && outcome, "gap = %v, want no chunk and the outcome kept", gap)
		}
	})
	t.Run("clears each declared ledger path and no review-record path", func(t *testing.T) {
		t.Parallel()
		f, doctor := installedRestore(t, "restore-several", "Persisted fields: intent-ledger:commitment.gaps.chunk, intent-ledger:commitment.claims, review-record:commitment.bindings\n")
		putGaps(t, f.root, "install")
		commitmenttest.EditCommitment(t, f.root, func(state *intent.CommitmentState) {
			state.Claims = append(state.Claims, intent.OutcomeClaim{Milestone: "fixture", Outcome: "ended"})
		})
		r := restore(t, f, "install", &rebuildRecorder{t: t}, doctor)
		requireTest(t, r.exit == 0, "restore = (%d, %q), want exit 0", r.exit, r.stderr)
		state := ledgerCommitment(t, f.root)
		gaps, _ := state["gaps"].([]any)
		requireTest(t, len(gaps) == 1, "gaps = %v, want the one gap kept", gaps)
		gap, _ := gaps[0].(map[string]any)
		_, chunk := gap["chunk"]
		_, claims := state["claims"]
		_, bindings := state["bindings"]
		requireTest(t, !chunk && !claims && bindings, "commitment = %v, want no gap chunk, no claims key, and the bindings kept", state)
	})
	t.Run("a declared field that the ledger does not hold needs no change", func(t *testing.T) {
		t.Parallel()
		f, doctor := installedRestore(t, "restore-absent-field", installDeclared)
		commitmenttest.RequireLedgerKept(t, f.root, func() {
			r := restore(t, f, "install", &rebuildRecorder{t: t}, doctor)
			requireTest(t, r.exit == 0, "restore = (%d, %q), want exit 0", r.exit, r.stderr)
		})
		primary := freshness.PublishedExecutable(f.root)
		requireTest(t, fixtureFileText(t, primary) == rebuiltMain && freshness.Verify(f.root, primary) == nil, "primary = %q, verify %v; want the verified rebuild", fixtureFileText(t, primary), freshness.Verify(f.root, primary))
	})
	t.Run("a ledger file that does not exist holds no declared field", func(t *testing.T) {
		t.Parallel()
		f, doctor := installedRestore(t, "restore-no-ledger", gapsDeclared)
		ledger := commitmenttest.LedgerPath(t, f.root)
		mustNoError(t, os.Remove(ledger))
		r := restore(t, f, "install", &rebuildRecorder{t: t}, doctor)
		requireTest(t, r.exit == 0, "restore = (%d, %q, %q), want exit 0", r.exit, r.stdout, r.stderr)
		_, err := os.Lstat(ledger)
		requireTest(t, os.IsNotExist(err) && len(restoreCopies(t, f)) == 0, "ledger stat %v, copies %v; want no ledger and no copy", err, restoreCopies(t, f))
		primary := freshness.PublishedExecutable(f.root)
		requireTest(t, fixtureFileText(t, primary) == rebuiltMain && freshness.Verify(f.root, primary) == nil, "primary = %q, verify %v; want the verified rebuild", fixtureFileText(t, primary), freshness.Verify(f.root, primary))
	})
	t.Run("refuses before any write", func(t *testing.T) {
		t.Parallel()
		for _, row := range []struct {
			name    string
			shape   installShape
			prepare func(*testing.T, ownedAssignment)
			want    func(ownedAssignment) []string
		}{
			{name: "no build inputs", shape: installShape{declaration: installDeclared, deliverable: installSpec}, want: func(ownedAssignment) []string { return []string{"declares no Bench build inputs"} }},
			{name: "an undeclared default-branch spec", shape: installShape{inputs: true, deliverable: installSpec}, want: func(ownedAssignment) []string { return []string{"declares no persisted field"} }},
			{name: "an invalid declaration", shape: installShape{inputs: true, declaration: "Persisted fields: gate-evidence:x\n", deliverable: installSpec}, want: func(ownedAssignment) []string { return []string{spec.ErrPersistedFieldsInvalid.Error()} }},
			{name: "a ledger that the strict read refuses", shape: declaredInstall, prepare: func(t *testing.T, f ownedAssignment) { plantCommitmentField(t, f.root) }, want: func(f ownedAssignment) []string {
				return []string{"cannot read the intent ledger", install + "'" + f.creation.Assignment.Label + "'; then " + rerun}
			}},
			{name: "a refused ledger whose target tree is gone", shape: declaredInstall, prepare: func(t *testing.T, f ownedAssignment) {
				plantCommitmentField(t, f.root)
				mustNoError(t, os.RemoveAll(f.creation.Path))
			}, want: func(ownedAssignment) []string {
				return []string{install + "<target>; then " + rerun}
			}},
			{name: "a ledger that is not JSON", shape: declaredInstall, prepare: func(t *testing.T, f ownedAssignment) {
				mustWrite(t, commitmenttest.LedgerPath(t, f.root), []byte("{\n"), 0o600)
			}, want: func(ownedAssignment) []string {
				return []string{install + "<target>; then " + rerun}
			}},
			{name: "an open build between chunks", shape: declaredInstall, prepare: func(t *testing.T, f ownedAssignment) {
				openInstallBuild(t, f.root)
				putGaps(t, f.root, "install")
			}, want: func(ownedAssignment) []string {
				return []string{"holds an open build", "before bench worktree release", "next=reviewer: bench commitment block --outcome <id> --reason <text>; then " + rerun}
			}},
			{name: "an open build of an undeclared outcome", shape: declaredInstall, prepare: func(t *testing.T, f ownedAssignment) { openUndeclaredBuild(t, f) }, want: func(ownedAssignment) []string {
				return []string{"holds an open build", "next=reviewer: bench commitment block"}
			}},
		} {
			t.Run(row.name, func(t *testing.T) {
				t.Parallel()
				f := newInstallAssignment(t, "restore-refuses-"+strings.ReplaceAll(row.name, " ", "-"), row.shape)
				if row.prepare != nil {
					row.prepare(t, f)
				}
				build, doctor := &rebuildRecorder{t: t}, &doctorRecorder{}
				commitmenttest.RequireLedgerKept(t, f.root, func() {
					r := restore(t, f, "install", build, doctor)
					for _, want := range row.want(f) {
						requireTest(t, r.exit == 1 && strings.Contains(r.stderr, want), "restore = (%d, %q), want exit 1 naming %q", r.exit, r.stderr, want)
					}
				})
				requireTest(t, build.calls == 0 && len(doctor.calls) == 0 && len(restoreCopies(t, f)) == 0, "build calls %d, doctor calls %d, copies %v; want none", build.calls, len(doctor.calls), restoreCopies(t, f))
				_, err := os.Lstat(filepath.Dir(freshness.PublishedExecutable(f.root)))
				requireTest(t, os.IsNotExist(err), "refused restore left the primary executable directory: %v", err)
			})
		}
	})
	t.Run("a failed copy refuses before the clear", func(t *testing.T) {
		t.Parallel()
		f, doctor := installedRestore(t, "restore-copy-fails", gapsDeclared)
		putGaps(t, f.root, "install")
		mustWrite(t, filepath.Join(f.home, "restore"), []byte("not a directory"), 0o644)
		build := &rebuildRecorder{t: t}
		commitmenttest.RequireLedgerKept(t, f.root, func() {
			r := restore(t, f, "install", build, doctor)
			requireTest(t, r.exit == 1 && strings.Contains(r.stderr, "preserve intent ledger") && strings.Contains(r.stderr, "next="+rerun), "restore = (%d, %q), want the copy refusal", r.exit, r.stderr)
		})
		requireTest(t, build.calls == 0, "rebuild calls = %d, want none", build.calls)
	})
	t.Run("a failed rebuild stops before the repair", func(t *testing.T) {
		t.Parallel()
		f, doctor := installedRestore(t, "restore-rebuild-fails", installDeclared)
		r := restore(t, f, "install", &rebuildRecorder{t: t, err: errors.New("the primary sources do not compile")}, doctor)
		requireTest(t, r.exit == 1 && strings.Contains(r.stderr, "the primary sources do not compile") && len(doctor.calls) == 1, "restore = (%d, %q), doctor calls %d; want the build failure and no second repair", r.exit, r.stderr, len(doctor.calls))
	})
	t.Run("a failed repair names the rerun", func(t *testing.T) {
		t.Parallel()
		f, doctor := installedRestore(t, "restore-repair-fails", installDeclared)
		doctor.exit = 3
		r := restore(t, f, "install", &rebuildRecorder{t: t}, doctor)
		requireTest(t, r.exit == 1 && strings.Contains(r.stderr, "exited 3") && strings.Contains(r.stderr, "next="+rerun), "restore = (%d, %q), want the repair refusal and the rerun", r.exit, r.stderr)
	})
}
