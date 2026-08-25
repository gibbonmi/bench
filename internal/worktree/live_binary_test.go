package worktree

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/freshness"
)

// newResidueGuardFixture builds a repository whose dist/ is both ignored and declared as
// build output, plus one owned assignment holding an empty dist/. This is the exact shape
// that lets a release reach the residue guard's removal loop. It returns the private
// home the registration lives under, so the fixture binds no process environment.
func newResidueGuardFixture(t *testing.T, request string) (string, Creation, string) {
	t.Helper()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("dist/\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore build output")
	mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[\"dist/\"]}\n"), 0o644)
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, request, "residue guard")
	mustMkdirAll(t, filepath.Join(creation.Path, "dist"), 0o755)
	return root, creation, home
}

// stubbedLiveBinaryJoins returns a seam set whose residue-guard warning sink is the
// returned buffer and whose running-binary resolution answers with running. Each test
// holds its own value, so no two tests share a stub point.
func stubbedLiveBinaryJoins(running string) (joins, *bytes.Buffer) {
	buffer := &bytes.Buffer{}
	j := defaultJoins()
	j.liveBinaryWarnings = buffer
	j.resolveRunningBinary = func() (string, error) { return running, nil }
	return j, buffer
}

// TestResidueGuardWarnsBeforeRemovingTheLiveBinary is H25. Removing the dist/bench the
// wrapper resolved takes down the CLI, the git guard that CLI backs, and the gate's
// BENCH_RUN_BINARY together. Without this test, the guard does that silently. The
// warning has to name the scripts/go-build.sh invocation, because plain `go build` leaves
// the package version the version and upgrade contracts read unstamped.
func TestResidueGuardWarnsBeforeRemovingTheLiveBinary(t *testing.T) {
	t.Parallel()
	const request = "landed-live-binary"
	root, creation, home := newResidueGuardFixture(t, request)
	live := filepath.Join(creation.Path, "dist", "bench")
	mustWrite(t, live, []byte("binary\n"), 0o755)
	j, warnings := stubbedLiveBinaryJoins(live)

	var stdout bytes.Buffer
	code := releaseCommandWith(j, root, home, []string{"--request", request, creation.Path}, &stdout, io.Discard)
	requireTest(t, code == 0, "live-binary release exit=%d stdout=%q", code, stdout.String())

	warned := warnings.String()
	rebuild := freshness.RebuildAction(creation.Path)
	requireTest(t, strings.Contains(warned, live), "warning %q does not name the binary it removed (%s)", warned, live)
	requireTest(t, strings.Contains(warned, rebuild), "warning %q does not name the rebuild invocation (%s)", warned, rebuild)
	requireTest(t, !strings.Contains(warned, "go build "), "warning %q names plain `go build`, which leaves the package version unstamped", warned)
	if _, err := os.Stat(live); !os.IsNotExist(err) {
		t.Errorf("H25 live binary still present after release (stat err=%v); a guard that warns and then skips the removal would pass on the warning text alone", err)
	}
}

// TestIsRunningBinaryFailsSafeWhenResolutionIsUnknown pins the two branches that decide
// the guard's posture when it cannot learn which binary is running. Both answer "warn",
// because nothing has been shown about the candidate. Flipping either to "remove
// silently" is exactly the incident this guard exists to prevent. No fixture reaches
// them through ReleaseCommand, whose stub always resolves.
func TestIsRunningBinaryFailsSafeWhenResolutionIsUnknown(t *testing.T) {
	t.Parallel()
	candidate := filepath.Join(t.TempDir(), "bench")
	mustWrite(t, candidate, []byte("binary\n"), 0o755)

	t.Run("executable unresolvable", func(t *testing.T) {
		t.Parallel()
		j := defaultJoins()
		j.resolveRunningBinary = func() (string, error) { return "", errors.New("no executable") }
		if !isRunningBinary(j, candidate) {
			t.Error("isRunningBinary = false when the running executable is unresolvable, want true (unknown is not absent)")
		}
	})

	t.Run("running path unstattable", func(t *testing.T) {
		t.Parallel()
		j, _ := stubbedLiveBinaryJoins(filepath.Join(t.TempDir(), "vanished", "bench"))
		if !isRunningBinary(j, candidate) {
			t.Error("isRunningBinary = false when the running path cannot be stat'd, want true (unknown is not absent)")
		}
	})

	t.Run("candidate absent", func(t *testing.T) {
		t.Parallel()
		j, _ := stubbedLiveBinaryJoins(candidate)
		if isRunningBinary(j, filepath.Join(t.TempDir(), "absent")) {
			t.Error("isRunningBinary = true for a candidate that does not exist, want false (no live file there to lose)")
		}
	})
}

// TestIsRunningBinaryResolvesThroughASymlink is the profile's "invocation through a
// symlink rather than the real path" class. The wrapper may exec a link, and the guard
// must still recognize the target as the live binary. Two things deliver that today: the
// EvalSymlinks normalization and os.Stat following the link. So no single mutation
// reddens this. It pins the behavior, not either mechanism.
func TestIsRunningBinaryResolvesThroughASymlink(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	real := filepath.Join(dir, "bench")
	mustWrite(t, real, []byte("binary\n"), 0o755)
	link := filepath.Join(dir, "bench-link")
	if err := os.Symlink(real, link); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("cannot create a symlink: %v", err))
	}
	j, _ := stubbedLiveBinaryJoins(link)
	if !isRunningBinary(j, real) {
		t.Error("isRunningBinary = false for the target of the resolved symlink, want true")
	}
}

// TestResidueGuardRemovesForeignBinariesWithoutWarning is H26. A path-shaped predicate
// warns on every checkout that has ever been built, which trains the warning away before
// the one removal that matters arrives. Identity keeps ordinary residue ordinary. The
// space-and-glob name also pins that the guard resolves its candidate exactly rather than
// through a shell-expanded pattern.
func TestResidueGuardRemovesForeignBinariesWithoutWarning(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, residue string }{
		{"plain", "bench"},
		{"space and glob", "be nch[1]*"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "landed-foreign-" + strings.ReplaceAll(tc.name, " ", "-")
			root, creation, home := newResidueGuardFixture(t, request)
			foreign := filepath.Join(creation.Path, "dist", tc.residue)
			mustWrite(t, foreign, []byte("binary\n"), 0o755)
			elsewhere := filepath.Join(t.TempDir(), "bench")
			mustWrite(t, elsewhere, []byte("the binary answering bench\n"), 0o755)
			j, warnings := stubbedLiveBinaryJoins(elsewhere)

			var stdout bytes.Buffer
			code := releaseCommandWith(j, root, home, []string{"--request", request, creation.Path}, &stdout, io.Discard)
			requireTest(t, code == 0, "foreign release exit=%d stdout=%q", code, stdout.String())
			requireTest(t, warnings.Len() == 0, "foreign %s warned: %q", foreign, warnings.String())
		})
	}
}
