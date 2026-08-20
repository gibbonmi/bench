package worktree

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
)

// newResidueGuardFixture builds a repository whose dist/ is both ignored and declared as
// build output, plus one owned assignment holding an empty dist/ — the exact shape that
// lets a release reach the residue guard's removal loop.
func newResidueGuardFixture(t *testing.T, request string) (string, Creation) {
	t.Helper()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("dist/\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore build output")
	mustMkdirAll(t, filepath.Join(root, ".bench"), 0o755)
	mustWrite(t, filepath.Join(root, ".bench", "build-outputs.json"), []byte("{\"schema\":1,\"paths\":[\"dist/\"]}\n"), 0o644)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	creation := mustCreate(t, root, request, "residue guard")
	mustMkdirAll(t, filepath.Join(creation.Path, "dist"), 0o755)
	return root, creation
}

// captureLiveBinaryWarnings redirects the guard's warning sink for one test.
func captureLiveBinaryWarnings(t *testing.T) *bytes.Buffer {
	t.Helper()
	previous := liveBinaryWarnings
	buffer := &bytes.Buffer{}
	liveBinaryWarnings = buffer
	t.Cleanup(func() { liveBinaryWarnings = previous })
	return buffer
}

// stubRunningBinary answers the guard's "which binary did the wrapper resolve" question
// with path, standing in for the exec the wrapper performs outside a test process.
func stubRunningBinary(t *testing.T, path string) {
	t.Helper()
	previous := resolveRunningBinary
	resolveRunningBinary = func() (string, error) { return path, nil }
	t.Cleanup(func() { resolveRunningBinary = previous })
}

// TestResidueGuardWarnsBeforeRemovingTheLiveBinary is H25. Removing the dist/bench the
// wrapper resolved takes the CLI, the git guard that CLI backs, and the gate's
// BENCH_RUN_BINARY down together, and the guard did it silently. The warning has to name
// the scripts/go-build.sh invocation, because plain `go build` leaves the package version
// the version and upgrade contracts read unstamped.
func TestResidueGuardWarnsBeforeRemovingTheLiveBinary(t *testing.T) {
	const request = "landed-live-binary"
	root, creation := newResidueGuardFixture(t, request)
	live := filepath.Join(creation.Path, "dist", "bench")
	mustWrite(t, live, []byte("binary\n"), 0o755)
	warnings := captureLiveBinaryWarnings(t)
	stubRunningBinary(t, live)

	var stdout bytes.Buffer
	code := ReleaseCommand(root, []string{"--request", request, creation.Path}, &stdout, io.Discard)
	requireTest(t, code == 0, "live-binary release exit=%d stdout=%q", code, stdout.String())

	warned := warnings.String()
	rebuild := freshness.RebuildAction(creation.Path)
	requireTest(t, strings.Contains(warned, live), "warning %q does not name the binary it removed (%s)", warned, live)
	requireTest(t, strings.Contains(warned, rebuild), "warning %q does not name the rebuild invocation (%s)", warned, rebuild)
	requireTest(t, !strings.Contains(warned, "go build "), "warning %q names plain `go build`, which leaves the package version unstamped", warned)
}

// TestResidueGuardRemovesForeignBinariesWithoutWarning is H26. A path-shaped predicate
// warns on every checkout that has ever been built, which trains the warning away before
// the one removal that matters arrives; identity keeps ordinary residue ordinary. The
// space-and-glob name also pins that the guard resolves its candidate exactly rather than
// through a shell-expanded pattern.
func TestResidueGuardRemovesForeignBinariesWithoutWarning(t *testing.T) {
	for _, tc := range []struct{ name, residue string }{
		{"plain", "bench"},
		{"space and glob", "be nch[1]*"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "landed-foreign-" + strings.ReplaceAll(tc.name, " ", "-")
			root, creation := newResidueGuardFixture(t, request)
			foreign := filepath.Join(creation.Path, "dist", tc.residue)
			mustWrite(t, foreign, []byte("binary\n"), 0o755)
			elsewhere := filepath.Join(t.TempDir(), "bench")
			mustWrite(t, elsewhere, []byte("the binary answering bench\n"), 0o755)
			warnings := captureLiveBinaryWarnings(t)
			stubRunningBinary(t, elsewhere)

			var stdout bytes.Buffer
			code := ReleaseCommand(root, []string{"--request", request, creation.Path}, &stdout, io.Discard)
			requireTest(t, code == 0, "foreign release exit=%d stdout=%q", code, stdout.String())
			requireTest(t, warnings.Len() == 0, "foreign %s warned: %q", foreign, warnings.String())
		})
	}
}
