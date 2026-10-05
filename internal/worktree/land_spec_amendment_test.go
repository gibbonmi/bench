// Spec-amendment landing tests: an in-range amendment publishes, and a resume completes an amended source.
package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/spec"
)

// landingSpecAmendment commits an in-range spec amendment on the reviewed source and
// returns the amended bytes. The fixture spec fences only owned.txt, so nothing but the
// implicit specs/<slug>/ authorization lets this commit through the source preflight.
func landingSpecAmendment(t *testing.T, source string) []byte {
	t.Helper()
	rel := filepath.Join("specs", "x", "spec.md")
	body, err := os.ReadFile(filepath.Join(source, rel))
	if err != nil {
		t.Fatal(err)
	}
	amended := bytes.Replace(body, []byte("1. Land source."), []byte("1. Land source, as the review amended it."), 1)
	if bytes.Equal(amended, body) {
		t.Fatal("spec amendment fixture changed nothing")
	}
	commitInWorktree(t, source, rel, string(amended), "amend the spec in range")
	refreshLandingEvidence(t, source, gitOutput(t, source, "merge-base", "main", "HEAD"))
	return amended
}

func TestLandCommandPublicLandsAnInRangeSpecAmendment(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-spec-amendment"
	f := publicLandingFixture(t, request, "", "")
	amended := landingSpecAmendment(t, f.creation.Path)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", f.base, "--source-tip", tip, "--spec", "x", "-m", "land the amended source", f.creation.Path)
	cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
	if code := exitCode(cmd.Run()); code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("amended landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	requirePublishedSpec(t, f.root, published, amended)
	parents := strings.Fields(gitOutput(t, f.root, "rev-list", "--parents", "-n", "1", published))
	if len(parents) != 3 || parents[1] != f.base || parents[2] != tip {
		t.Fatalf("published parents = %q, want destination %s and source %s", parents, f.base, tip)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
}

func TestResumeLandCommandPublicCompletesAnAmendedSourceLanding(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-resume-spec-amendment"
	f := publicLandingFixture(t, request, "private/output", "dist/")
	amended := landingSpecAmendment(t, f.creation.Path)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	land := func(args ...string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
		cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}
	code, stdout, stderr := land("--request", request, "--base", f.base, "--source-tip", tip, "--spec", "x", "-m", "land the amended source", f.creation.Path)
	if code != 3 || !strings.Contains(stdout, "worktree=incomplete:release") {
		t.Fatalf("interrupted amended landing = (%d, %q, %q)", code, stdout, stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	requirePublishedSpec(t, f.root, published, amended)
	if err := os.Remove(filepath.Join(f.creation.Path, "private", "output")); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr = land("--resume", published, "--request", request, "--base", f.base, "--source-tip", tip, "--spec", "x", f.creation.Path)
	if code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") || stderr != "" {
		t.Fatalf("amended resume = (%d, %q, %q)", code, stdout, stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != published {
		t.Fatalf("resume republished: main=%s, want %s", got, published)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("resume reran the gate: tally=%q error=%v", got, err)
	}
}

func requirePublishedSpec(t *testing.T, root, published string, staged []byte) {
	t.Helper()
	want, err := spec.Implemented(staged)
	if err != nil {
		t.Fatal(err)
	}
	got, err := git.Raw("-C", root, "show", published+":specs/x/spec.md")
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("published spec = %q (%v), want %q", got, err, want)
	}
}

// Adoption lists an existing run with its approved scope instead of a delivery
// binding. The listed run lands that scope and nothing beyond it, and an unlisted run with
// no binding lands nothing. The fixture spec is the approved deliverable of a rowless
// outcome, so each scope lists it and the landing carries its delivery fact. A run whose
// scope omits that spec, or that adoption did not list, is not authorized to record that
// fact, and its refusal names the scope or the missing binding.
func TestCommitmentLegacyContinuation(t *testing.T) {
	t.Parallel()
	const deliverable = "specs/x/spec.md"
	for _, row := range []struct {
		name  string
		scope []string
		want  string
	}{
		{name: "listed-scope", scope: []string{"owned.txt", "reviews/x.md", deliverable}},
		{name: "beyond-scope", scope: []string{"reviews/x.md", deliverable}, want: "legacy continuation scope excludes"},
		// The scope entry "owned" is a string prefix of the source's owned.txt, not its
		// directory, so the sibling stays outside the scope.
		{name: "sibling-prefix", scope: []string{"owned", "reviews/x.md", deliverable}, want: "legacy continuation scope excludes"},
		{name: "directory-scope", scope: []string{"owned.txt", "reviews", deliverable}},
		{name: "outside-directory", scope: []string{"reviews", deliverable}, want: "legacy continuation scope excludes"},
		// A scope that omits the approved deliverable does not authorize its delivery fact.
		{name: "scope-without-deliverable", scope: []string{"owned.txt", "reviews/x.md"}, want: "legacy continuation scope excludes \"" + deliverable + "\""},
		{name: "unlisted", want: "assignment has no current delivery binding"},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			request := "land-legacy-continuation-" + row.name
			f := publicLandingFixture(t, request, "", "")
			err := intent.Transact(f.root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
				state := intent.CommitmentState{}
				if row.scope != nil {
					state.Continuations = []intent.LegacyContinuation{{Assignment: f.creation.Assignment.ID, Request: f.creation.Assignment.Request, Scope: row.scope}}
				}
				ledger.Commitment = &state
				return ledger, true, nil
			}, nil)
			if err != nil {
				t.Fatal(err)
			}
			r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
			if row.want == "" {
				if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released") {
					t.Fatalf("listed continuation = (%d, %q, %q), want a published release", r.exit, r.stdout, r.stderr)
				}
				return
			}
			if r.exit != 1 || !strings.Contains(r.stdout, "refused{detail=commitment: ") || !strings.Contains(r.stdout, row.want) {
				t.Fatalf("continuation landing = (%d, %q, %q), want a commitment refusal naming %q", r.exit, r.stdout, r.stderr, row.want)
			}
			requireIdentityRefusalState(t, f.root, f.creation.Path, f.tally, 0)
		})
	}
}
