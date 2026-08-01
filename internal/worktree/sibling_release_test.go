package worktree

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// A sibling assignment whose checkpoint is replayed onto a candidate an earlier sibling
// already advanced integrates a tree that differs from its own checkpoint tree. That
// difference is composition, not drift, so both siblings must release.
func TestReleaseProvisionalReleasesSiblingsIntegratedOntoAMovedCandidate(t *testing.T) {
	root := newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	base := gitOutput(t, root, "rev-parse", "HEAD")
	candidateRef := "refs/bench/test-candidate/siblings"
	gitRun(t, root, "update-ref", candidateRef, base)

	first := newSiblingAssignment(t, root, "sibling-first", "first.txt")
	second := newSiblingAssignment(t, root, "sibling-second", "second.txt")
	requireTest(t, first.created.Assignment.Start == base && second.created.Assignment.Start == base,
		"siblings started from %s and %s, want the shared base %s", first.created.Assignment.Start, second.created.Assignment.Start, base)

	firstIntegrated := integrateSibling(t, root, candidateRef, base, first)
	requireTest(t, gitOutput(t, root, "rev-parse", firstIntegrated+"^{tree}") == gitOutput(t, root, "rev-parse", first.checkpoint+"^{tree}"),
		"first sibling integrated at the base must reproduce its checkpoint tree")
	mustNoError(t, ReleaseProvisional(root, first.request, first.created.Path, first.evidence(candidateRef, firstIntegrated)))

	secondIntegrated := integrateSibling(t, root, candidateRef, base, second)
	requireTest(t, gitOutput(t, root, "rev-parse", secondIntegrated+"^{tree}") != gitOutput(t, root, "rev-parse", second.checkpoint+"^{tree}"),
		"second sibling must integrate a tree that differs from its checkpoint tree, or the fixture is not composing")

	mustNoError(t, ReleaseProvisional(root, second.request, second.created.Path, second.evidence(candidateRef, secondIntegrated)))
	if _, err := os.Stat(second.created.Path); !os.IsNotExist(err) {
		t.Fatalf("released sibling checkout remains: %v", err)
	}
}

type siblingAssignment struct {
	request       string
	created       Creation
	checkpointRef string
	checkpoint    string
}

func (s siblingAssignment) evidence(candidateRef, integrated string) ProvisionalEvidence {
	return ProvisionalEvidence{
		Base: s.created.Assignment.Start, CheckpointRef: s.checkpointRef, Checkpoint: s.checkpoint,
		IntegratedRef: candidateRef, Integrated: integrated,
	}
}

func newSiblingAssignment(t *testing.T, root, request, name string) siblingAssignment {
	t.Helper()
	created := mustCreate(t, root, request, request)
	mustWrite(t, filepath.Join(created.Path, name), []byte(name+"\n"), 0o644)
	tree := benchgit.TreeHash(created.Path)
	checkpoint := gitOutput(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench",
		"commit-tree", tree, "-p", created.Assignment.Start, "-m", "checkpoint "+request)
	ref := "refs/bench/test-checkpoint/" + created.Assignment.ID
	gitRun(t, root, "update-ref", ref, checkpoint)
	return siblingAssignment{request: request, created: created, checkpointRef: ref, checkpoint: checkpoint}
}

// integrateSibling mirrors what the spec-build lifecycle does: it replays the sibling's
// base-relative patch onto the current candidate tip and advances the candidate.
func integrateSibling(t *testing.T, root, candidateRef, base string, s siblingAssignment) string {
	t.Helper()
	candidate := gitOutput(t, root, "rev-parse", candidateRef)
	index := filepath.Join(t.TempDir(), "replay-index")
	env := append(os.Environ(), "GIT_INDEX_FILE="+index)
	patch := gitRaw(t, root, nil, nil, "diff", "--binary", "--full-index", "--no-ext-diff", base, s.checkpoint)
	gitRaw(t, root, env, nil, "read-tree", candidate)
	gitRaw(t, root, env, patch, "apply", "--cached", "--whitespace=nowarn")
	tree := strings.TrimSpace(string(gitRaw(t, root, env, nil, "write-tree")))
	commit := gitOutput(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench",
		"commit-tree", tree, "-p", candidate, "-m", "integrate "+s.request)
	gitRun(t, root, "update-ref", candidateRef, commit, candidate)
	return commit
}

func gitRaw(t *testing.T, dir string, env []string, input []byte, args ...string) []byte {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = env
	if input != nil {
		cmd.Stdin = bytes.NewReader(input)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, stderr.String())
	}
	return out
}
