// Freshness tests for the landing command: stale-executable rebuild, absent build inputs, and sealless refusal.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
)

// landingAssignmentState renders the assignment records a landing would consume, without
// the creation timestamp's address. The pointer differs between two reads of the same
// unchanged state and would report every comparison as a change.
func landingAssignmentState(t *testing.T, root string) string {
	t.Helper()
	assignments, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}
	var state strings.Builder
	for _, a := range assignments {
		fmt.Fprintf(&state, "%s %s %s %s %s\n", assignmentString(a), a.State, a.Branch, a.Worktree, a.Start)
	}
	return state.String()
}

// TestLandCommandRefusesAnUntrustedExecutableBeforeAnyRepositoryProof grades the refusal
// (LF1) and the state it leaves behind (LF2). It also grades the empty manifest that a
// content-sniffing applicability predicate would skip (LF8), and the ordering against a
// destination that would independently refuse (LF9). That last case pins the remedy the
// operator reads to the rebuild command rather than to a later proof's message.
// A stale executable is rebuilt through the sanctioned build and the landing re-runs
// under the fresh binary, before any repository proof. Only a rebuild that fails, or
// a re-run that is still stale, surfaces the owner's message.
func TestLandCommandRebuildsAStaleExecutableBeforeAnyRepositoryProof(t *testing.T) {
	const untrusted = "bench binary is untrusted; rebuild with the sanctioned build"
	for _, tc := range []struct {
		name, manifest string
		setup          func(*testing.T, string)
		rebuildErr     error
		alreadyRebuilt bool
		wantRefusal    string
	}{
		{name: "declared-build-inputs", manifest: "build_script=scripts/go-build.sh\n"},
		{name: "empty-manifest"},
		{name: "destination-would-also-refuse", manifest: "build_script=scripts/go-build.sh\n", setup: func(t *testing.T, root string) {
			mustWrite(t, filepath.Join(root, "dirty"), []byte("dirty\n"), 0o600)
		}},
		{name: "rebuild-fails", manifest: "build_script=scripts/go-build.sh\n", rebuildErr: errors.New("go: command not found"), wantRefusal: untrusted + "; rebuild failed: go: command not found"},
		{name: "still-stale-after-rebuild", manifest: "build_script=scripts/go-build.sh\n", alreadyRebuilt: true, wantRefusal: untrusted},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "landed-untrusted-" + tc.name
			root := newWorktreeRepo(t)
			bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
			if tc.alreadyRebuilt {
				bindEnv(t, rebuiltLandingEnv, "1")
			}
			creation := mustCreate(t, root, request, "untrusted executable")
			commitLandingBuildInputs(t, root, tc.manifest)
			stageLandSpec(t, root, creation.Path)
			base := gitOutput(t, root, "rev-parse", "HEAD")
			commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
			tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
			gitRun(t, root, "update-ref", "refs/bench/green/main", base)
			if tc.setup != nil {
				tc.setup(t, root)
			}
			before := landingAssignmentState(t, root)
			calls, rebuilds := 0, 0
			var rebuiltRoot, rebuiltOutput, execPath string
			var execArgv, execEnv []string
			oldLand, oldVerify, oldRebuild, oldExec := landReviewed, verifyLandingExecutable, rebuildLandingExecutable, reexecLanding
			landReviewed = func(context.Context, landing.ReviewedRequest) (landing.ReviewedResult, error) {
				calls++
				return landing.ReviewedResult{}, errors.New("unexpected landing")
			}
			verifyLandingExecutable = func(string, string) error { return errors.New(untrusted) }
			rebuildLandingExecutable = func(_ context.Context, sourceRoot, output string) error {
				rebuilds++
				rebuiltRoot, rebuiltOutput = sourceRoot, output
				return tc.rebuildErr
			}
			reexecLanding = func(path string, argv, env []string) error {
				execPath, execArgv, execEnv = path, argv, env
				return nil
			}
			t.Cleanup(func() {
				landReviewed, verifyLandingExecutable, rebuildLandingExecutable, reexecLanding = oldLand, oldVerify, oldRebuild, oldExec
			})

			args := landArgs(request, base, tip, creation.Path)
			fresh := filepath.Join(root, "dist", "bench")
			var stdout, stderr bytes.Buffer
			code := LandCommand(root, fresh, args, &stdout, &stderr)
			if tc.wantRefusal != "" {
				if code != 1 || calls != 0 || execPath != "" || stdout.String() != "refused{detail="+tc.wantRefusal+"}\n" || stderr.Len() != 0 {
					t.Fatalf("stale refusal = (%d, calls=%d, exec=%q, stdout=%q, stderr=%q), want the owner's message", code, calls, execPath, stdout.String(), stderr.String())
				}
				if tc.alreadyRebuilt && rebuilds != 0 {
					t.Fatalf("a re-run that is still stale rebuilt again (%d)", rebuilds)
				}
			} else {
				wantArgv := strings.Join(append([]string{fresh, "worktree", "land"}, args...), "\x00")
				rebuiltEnv := false
				for _, item := range execEnv {
					rebuiltEnv = rebuiltEnv || item == rebuiltLandingEnv+"=1"
				}
				if calls != 0 || rebuilds != 1 || rebuiltRoot != root || rebuiltOutput != fresh || execPath != fresh || strings.Join(execArgv, "\x00") != wantArgv || !rebuiltEnv || !strings.Contains(stderr.String(), "rebuilt "+fresh) {
					t.Fatalf("stale rebuild = (calls=%d, rebuilds=%d, root=%q, output=%q, exec=%q, argv=%q, env-marked=%v, stderr=%q)", calls, rebuilds, rebuiltRoot, rebuiltOutput, execPath, execArgv, rebuiltEnv, stderr.String())
				}
			}
			if got := gitOutput(t, root, "rev-parse", "main"); got != base {
				t.Fatalf("refusal moved destination: got %s want %s", got, base)
			}
			if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != base {
				t.Fatalf("refusal moved project-green: got %s want %s", got, base)
			}
			if after := landingAssignmentState(t, root); after != before {
				t.Fatalf("refusal changed assignments:\nbefore: %safter:  %s", before, after)
			}
		})
	}
}

// TestLandCommandSkipsTheFreshnessProofWithoutDeclaredBuildInputs is LF3. A repository
// that declares no Go build inputs — every linked repository — never consults the owner.
// A substituted seam that fails the test on any call proves non-consultation while the
// landing runs to completion around it.
func TestLandCommandSkipsTheFreshnessProofWithoutDeclaredBuildInputs(t *testing.T) {
	request := "landed-no-build-inputs"
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	creation := mustCreate(t, root, request, "no declared build inputs")
	stageLandSpec(t, root, creation.Path)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "owned.txt", "owned\n", "owned")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	restore := stubLandJoins(t, base, tip)
	defer restore()
	releaseLandingAssignment = func(string, []string, io.Writer, io.Writer) int { return 0 }
	oldVerify := verifyLandingExecutable
	verifyLandingExecutable = func(string, string) error {
		t.Error("landing consulted the freshness owner without declared build inputs")
		return nil
	}
	t.Cleanup(func() { verifyLandingExecutable = oldVerify })

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, filepath.Join(root, "dist", "bench"), landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released}") {
		t.Fatalf("linked-repository landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
}

// TestLandCommandPublicRefusesASeallessExecutable is LF5: the real owner, reached through
// the land surface. The copied executable has no adjacent seal, which the owner refuses
// before it needs a Go toolchain in the fixture.
func TestLandCommandPublicRefusesASeallessExecutable(t *testing.T) {
	request := "public-land-sealless"
	root, creation, base, tip, tally := publicLandingFixture(t, request, "", "")
	commitLandingBuildInputs(t, root, "build_script=scripts/go-build.sh\n")
	published := gitOutput(t, root, "rev-parse", "main")
	sealless := filepath.Join(t.TempDir(), "bench")
	built, err := os.ReadFile(testRunBinary(t))
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, sealless, built, 0o755)

	var stdout, stderr bytes.Buffer
	cmd := descendant(t, sealless, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path)
	cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
	code := exitCode(cmd.Run())
	if code != 1 || !strings.HasPrefix(stdout.String(), "refused{detail=") || !strings.Contains(stdout.String(), sealless) || stderr.Len() != 0 {
		t.Fatalf("sealless landing = (%d, %q, %q), want a refusal naming %s", code, stdout.String(), stderr.String(), sealless)
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != published {
		t.Fatalf("sealless refusal moved destination: got %s want %s", got, published)
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("sealless refusal ran the gate: %v", err)
	}
}
