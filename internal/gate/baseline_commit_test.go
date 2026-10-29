package gate

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/gittest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// emptySchedule is the bytes of a phase manifest that declares no phase.
const emptySchedule = `{"phases":[]}`

// buildProspectiveSubjectForGeneration builds the prospective subject of root with its identity
// root as its baseline root, as every caller off the chunk-subject path builds it.
func buildProspectiveSubjectForGeneration(root, identityRoot string, generation *treeGeneration) (subject, error) {
	return buildSubjectInMode(root, identityRoot, identityRoot, policyVersion, generation.tree, true)
}

// committedChunkIdentity commits what write leaves in a new repository, removes the working
// phase manifest, and answers the grade identity of the chunk-subject run of the HEAD tree.
func committedChunkIdentity(t *testing.T, write func(root string)) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	write(root)
	outcomeCommit(t, root, "retain the phase manifest")
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(canary.PhaseManifestPath))); err != nil {
		t.Fatal(err)
	}
	tree := gittest.Output(t, root, "rev-parse", "HEAD^{tree}")
	artifacts, err := openProspectiveArtifacts(root, tree)
	if err != nil {
		t.Fatal(err)
	}
	defer artifacts.Close()
	checkout := artifacts.Checkout()
	plan, err := newProspectiveTreeEvaluation(checkout, root, baselineRoot(withCheckoutBaseline(t.Context()), checkout, root), tree).acceptPre()
	if err != nil {
		t.Fatal(err)
	}
	return plan.Identity
}

// baselineFixture is a checkpoint fixture whose committed tree declares Go build inputs and
// holds a prospective gate that runs the fixture gate. So each chunk-subject run frames a real
// baseline runner digest. The prepare steps run before the fixture commits its chunk.
func baselineFixture(t *testing.T, prepare ...func(*recordtest.Fixture)) *recordtest.Fixture {
	t.Helper()
	t.Setenv("BENCH_KIT", "")
	return attachedCheckpointFixture(t, recordtest.Attach, append([]func(*recordtest.Fixture){func(f *recordtest.Fixture) {
		writeRunBinarySource(t, f.Root, "")
		outcomeWrite(t, f.Root, prospectiveGatePath, "#!/bin/sh\nexec .bench/gate.sh\n", 0o755)
	}}, prepare...)...)
}

// editBuildInput changes the one Go source of the fixture command, and requires that the
// change moves the build-input digest of the worktree.
func editBuildInput(t *testing.T, f *recordtest.Fixture) {
	t.Helper()
	before := buildInputsIdentity(f.Root)
	f.Write("cmd/bench/main.go", "package main\n\nfunc main() { println() }\n")
	if after := buildInputsIdentity(f.Root); after == before || after == unreadableBaselineRunner {
		t.Fatalf("build-input digest = %s before and %s after the edit, want two readable digests", before, after)
	}
}

// strippedTree answers the HEAD tree of f without its review record, which a chunk-subject
// run grades.
func strippedTree(t *testing.T, f *recordtest.Fixture) string {
	t.Helper()
	stripped, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	return stripped
}

// BG16: a Go build input of the worktree changes after the chunk checkpoint accepts its
// subject. The baseline runner digest reads the private checkout, so the subject stays the
// same and the checkpoint passes.
func TestBaselineReadsSurviveAWorktreeEdit(t *testing.T) {
	f := baselineFixture(t)
	var stderr bytes.Buffer
	if code := RunCheckpointArmed(t, f.Root, []string{"--checkpoint", recordtest.Spec, "--chunk", "1"}, func() { editBuildInput(t, f) }, &stderr); code != 0 {
		t.Fatalf("chunk checkpoint across a worktree edit of a Go build input = %d %s; want its green", code, stderr.String())
	}
}

// BG17: a commit that changes a Go build input lands in the worktree during a background grade.
// The grade reads both baseline members from its private checkout, so it records its green for
// the tree that it graded.
func TestBaselineReadsSurviveABuildInputCommit(t *testing.T) {
	f := baselineFixture(t)
	seedGateCache(t, f)
	stripped := strippedTree(t, f)
	commit := func() {
		editBuildInput(t, f)
		f.Commit("change a Go build input during the run")
	}
	if code, out := gradeBackground(t, f, commit); code != 0 {
		t.Fatalf("background grade across a build-input commit = %d %s; want its green", code, out)
	}
	if got := InspectTree(f.Root, stripped); !got.ReusableGreen {
		t.Fatalf("evidence of the graded tree = %+v; want the green of the run", got)
	}
}

// BG18: the worktree holds a phase manifest that its commit does not hold. The chunk-subject
// run frames the absent manifest of the commit, so its green answers the key of the committed
// schedule. The working file still moves the key of a live-root subject.
func TestBaselineScheduleReadsTheCommit(t *testing.T) {
	f := baselineFixture(t, func(f *recordtest.Fixture) {
		f.Write(".gitignore", ".gate-*\n"+canary.PhaseManifestPath+"\n")
	})
	seedGateCache(t, f)
	stripped := strippedTree(t, f)
	working := filepath.Join(f.Root, filepath.FromSlash(canary.PhaseManifestPath))
	outcomeWrite(t, f.Root, canary.PhaseManifestPath, emptySchedule, 0o644)
	if code, out := gradeBackground(t, f, nil); code != 0 {
		t.Fatalf("background grade over an uncommitted phase manifest = %d %s; want its green", code, out)
	}
	if err := os.Remove(working); err != nil {
		t.Fatal(err)
	}
	if got := InspectTree(f.Root, stripped); !got.ReusableGreen {
		t.Fatalf("evidence under the committed schedule = %+v; want the green of the run", got)
	}
	outcomeWrite(t, f.Root, canary.PhaseManifestPath, emptySchedule, 0o644)
	if got := InspectTree(f.Root, stripped); got.ReusableGreen {
		t.Fatalf("evidence under the working manifest = %+v; want no green, because the working file moves a live-root key", got)
	}
}

// A committed link at the phase manifest enters the chunk-subject identity as the bytes of its
// target, as GR5 holds for a live link. So it gives the identity of a committed manifest with
// those bytes, and two targets with different bytes give two identities.
func TestBaselineScheduleFollowsACommittedLink(t *testing.T) {
	t.Setenv("BENCH_KIT", "")
	const target = ".bench/phases-target.json"
	regular := committedChunkIdentity(t, func(root string) {
		outcomeWrite(t, root, canary.PhaseManifestPath, emptySchedule, 0o644)
	})
	linked := func(body string) string {
		return committedChunkIdentity(t, func(root string) {
			outcomeWrite(t, root, target, body, 0o644)
			if err := os.Symlink(filepath.Base(target), filepath.Join(root, filepath.FromSlash(canary.PhaseManifestPath))); err != nil {
				t.Fatal(err)
			}
		})
	}
	first, second := linked(emptySchedule), linked(`{"phases":[1]}`)
	if regular != first {
		t.Fatalf("identity = %s for the committed manifest and %s for a committed link to its bytes, want one identity", regular, first)
	}
	if first == second {
		t.Fatalf("identity = %s for committed links to two targets, want two identities", first)
	}
}
