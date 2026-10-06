package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestCommandCommittedDeletionResolvesWrites exercises the review command over an
// explicit frozen range. A path deleted in that range remains valid ticket ownership.
func TestCommandCommittedDeletionResolvesWrites(t *testing.T) {
	_, slug, base, tip := seedDeletionReview(t, "", true)

	out, code := Command([]string{"review", slug, "--base", base, "--source-tip", tip})
	if code != 0 {
		t.Fatalf("committed deletion review exit = %d, want 0; output:\n%s", code, out)
	}
	if strings.Contains(out, "writes-resolve,red") {
		t.Fatalf("committed deletion remained unresolved:\n%s", out)
	}
}

// TestCommandCommittedDeletionDoesNotResolveAnotherMissingPath keeps exact-path
// evidence fail closed. Another absent path remains red and is named in the result.
func TestCommandCommittedDeletionDoesNotResolveAnotherMissingPath(t *testing.T) {
	const missing = "internal/example/typo.go"
	_, slug, base, tip := seedDeletionReview(t, missing, true)

	out, code := Command([]string{"review", slug, "--base", base, "--source-tip", tip})
	if code != 1 || !strings.Contains(out, "writes-resolve,red") || !strings.Contains(out, missing) {
		t.Fatalf("unrelated missing path review = (%d, %q), want writes-resolve red naming %s", code, out, missing)
	}
}

func TestCommandCommittedDeletionDoesNotResolveMissingParentDirectory(t *testing.T) {
	const missing = "internal/example"
	_, slug, base, tip := seedDeletionReview(t, missing, true)

	out, code := Command([]string{"review", slug, "--base", base, "--source-tip", tip})
	if code != 1 || !strings.Contains(out, "writes-resolve,red") || !strings.Contains(out, missing) {
		t.Fatalf("missing parent directory review = (%d, %q), want writes-resolve red naming %s", code, out, missing)
	}
}

// TestCommandFrozenRangeWithoutDeletionEvidenceReds proves that a deletion outside
// the selected range does not resolve the same absent path.
func TestCommandFrozenRangeWithoutDeletionEvidenceReds(t *testing.T) {
	_, slug, _, tip := seedDeletionReview(t, "", true)

	out, code := Command([]string{"review", slug, "--base", tip, "--source-tip", tip})
	if code != 1 || !strings.Contains(out, "writes-resolve,red") || !strings.Contains(out, "internal/example/gone.go") {
		t.Fatalf("range after deletion = (%d, %q), want writes-resolve red naming the deleted path", code, out)
	}
}

// TestCommandWorkingTreeDeletionIsNotCommittedDeletion proves that review keeps
// its clean-source prerequisite when the declared path is only deleted locally.
func TestCommandWorkingTreeDeletionIsNotCommittedDeletion(t *testing.T) {
	_, slug, base, tip := seedDeletionReview(t, "", false)

	out, code := Command([]string{"review", slug, "--base", base, "--source-tip", tip})
	if code != 1 || !strings.Contains(out, "source not clean") {
		t.Fatalf("working-tree deletion review = (%d, %q), want source-not-clean refusal", code, out)
	}
}

func seedDeletionReview(t *testing.T, missing string, commitDeletion bool) (root, slug, base, tip string) {
	t.Helper()
	slug = "example"
	root = preflighttest.StartRepo(t)
	fence := append([]string{"internal/example/gone.go"}, preflighttest.ConformantFence[1:]...)
	extraFence := []string{}
	if missing != "" {
		fence = append(fence, missing)
		extraFence = append(extraFence, "- `"+missing+"` (missing-path control)")
	}
	spec := strings.Replace(preflighttest.SpecBody(slug, extraFence...), "`internal/example/`", "`internal/example/gone.go`", 1)
	writes := preflighttest.FenceWrites(fence)
	writes[0] = "internal/example/gone.go"
	if missing != "" {
		writes[len(writes)-1] = missing
	}
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", spec)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", preflighttest.WritesTicketDoc("One", writes, "PF1", "PF2"))
	preflighttest.MustWriteFile(t, "internal/example/gone.go", "package example\n")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "base")
	base = preflighttest.RunGit(t, "rev-parse", "HEAD")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	preflighttest.RunGit(t, "rm", "-q", "internal/example/gone.go")
	if commitDeletion {
		preflighttest.RunGit(t, "commit", "-q", "-m", "delete owned path")
	}
	tip = preflighttest.RunGit(t, "rev-parse", "HEAD")
	return root, slug, base, tip
}
