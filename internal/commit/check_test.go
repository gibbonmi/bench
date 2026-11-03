package commit

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// errChainRed is the refusal of the stub chain check.
var errChainRed = errors.New("the chain's preflight is red on a.txt")

// redCandidateCheck is the stub chain check: it reds a candidate commit whose a.txt holds red.
func redCandidateCheck(root, candidate string) error {
	if body, err := git.Raw("-C", root, "show", candidate+":a.txt"); err != nil || strings.TrimSpace(string(body)) == "red" {
		return errors.Join(errChainRed, err)
	}
	return nil
}

// chainRedFixture is the commit-chain-red route fixture. The caller's chain check reds the
// candidate until the walk repairs a.txt, and the rerun that the route prints keeps the
// chain flags and publishes.
func chainRedFixture(rerun func(label string, flags ...string) string) commitFaceFixture {
	flags := []string{PreflightBuildFlag, "commit-fixture"}
	return commitFaceFixture{
		face: faceChainRed, flags: flags, exit: 1,
		build: func(t *testing.T) commitSet {
			f := assignedCommitSet(t, faceChainRed, gateScript("exit 0"))
			f.check = redCandidateCheck
			mustWrite(t, filepath.Join(f.checkout, "a.txt"), "red\n", 0o644)
			return f
		},
		suffix: rerun(faceChainRed, PreflightBuildFlag, sanitize.ShellQuote("commit-fixture")),
		carry: map[int]func(*testing.T, commitSet){0: func(t *testing.T, f commitSet) {
			mustWrite(t, filepath.Join(f.checkout, "a.txt"), "fixed\n", 0o644)
		}},
	}
}

// faceSteps checks the route that stderr prints against the registered face name and
// returns its steps. The first step must be the face's own, which tells apart two faces of
// one step shape.
func faceSteps(t *testing.T, name, stderr string) []string {
	t.Helper()
	next, printed := printedNext(stderr)
	if !printed {
		t.Fatalf("stderr = %q, want one next= route of the registered face %s", stderr, name)
	}
	face := registeredFace(t, name)
	steps := routetest.Steps(t, face, next, "", nil)
	if first := face.Route[0].Filled(func(string) string { return "" }); steps[0] != first {
		t.Fatalf("%s next = %q, want the first step %q", name, next, first)
	}
	return steps
}

// PP13: a red check refuses the commit at exit 1 with the commit-chain-red route, whose rerun
// keeps the caller's --preflight-build, and HEAD does not move. The check receives the
// checkout and the candidate commit. A check error that wraps another face's refusal still
// prints the chain-red route, because every check refusal routes as commit-chain-red.
func TestRunCheckedRefusesARedCheck(t *testing.T) {
	quoted := sanitize.ShellQuote
	for _, tc := range []struct {
		name string
		wrap func(error) error
	}{
		{"plain", func(err error) error { return err }},
		{"wraps a raised face", func(err error) error {
			return fmt.Errorf("chain: %w", refusalroute.Raised{Name: faceHandback, Err: err})
		}},
		{"wraps a tip-moved refusal", func(err error) error {
			return fmt.Errorf("chain: %w", landing.TipMovedError{Err: err})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := assignedCommitSet(t, "chain", gateScript("exit 0"))
			mustWrite(t, filepath.Join(f.checkout, "a.txt"), "red\n", 0o644)
			before := gittest.Output(t, f.checkout, "rev-parse", "HEAD")
			var roots, candidates []string
			check := func(root, candidate string) error {
				roots, candidates = append(roots, root), append(candidates, candidate)
				if err := redCandidateCheck(root, candidate); err != nil {
					return tc.wrap(err)
				}
				return nil
			}
			code, stdout, stderr := runChecked(t, f.checkout, check, PreflightBuildFlag, "commit-fixture", "-m", "m", "--", "a.txt")
			steps := faceSteps(t, faceChainRed, stderr)
			want := callerCommit(quoted("chain"), PreflightBuildFlag, quoted("commit-fixture"), "-m", quoted("m"), "--", quoted("a.txt"))
			if code != 1 || steps[len(steps)-1] != want || !strings.Contains(stderr, errChainRed.Error()) {
				t.Fatalf("commit = (%d, %q, %q), want exit 1, the chain-red sentence, and a next= route that repairs the chain and ends with %q", code, stdout, stderr, want)
			}
			if head := gittest.Output(t, f.checkout, "rev-parse", "HEAD"); head != before {
				t.Fatalf("HEAD moved from %s to %s", before, head)
			}
			if top := gittest.Output(t, f.checkout, "rev-parse", "--show-toplevel"); len(roots) != 1 || roots[0] != top || candidates[0] == before {
				t.Fatalf("check calls = roots %q candidates %q, want one call at %s with a candidate other than HEAD", roots, candidates, top)
			}
		})
	}
}

// PP14: a lane red refuses with the lane's own route before any commit object exists, so a
// supplied check never runs.
func TestRunCheckedLaneRedWinsOverTheCheck(t *testing.T) {
	f := assignedCommitSet(t, "lane", func(t *testing.T, root string) {
		gateScript("exit 0")(t, root)
		mustWrite(t, filepath.Join(root, ".bench", "phases.json"), `{"lane":[{"name":"check","argv":["sh","-c","exit 1"]}]}`, 0o644)
	})
	before := gittest.Output(t, f.checkout, "rev-parse", "HEAD")
	calls := 0
	check := func(string, string) error { calls++; return nil }
	code, stdout, stderr := runChecked(t, f.checkout, check, PreflightBuildFlag, "commit-fixture", "-m", "m", "--", "a.txt")
	faceSteps(t, faceRed, stderr)
	if code != 1 || calls != 0 {
		t.Fatalf("commit = (%d, %q, %q) with %d check calls, want exit 1, the lane's commit-red route, and no check call", code, stdout, stderr, calls)
	}
	if head := gittest.Output(t, f.checkout, "rev-parse", "HEAD"); head != before {
		t.Fatalf("HEAD moved from %s to %s", before, head)
	}
}
