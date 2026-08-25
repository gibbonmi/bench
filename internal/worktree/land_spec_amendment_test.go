// Spec-amendment landing tests: an in-range amendment publishes, and a resume completes an amended source.
package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
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
	return amended
}

func TestLandCommandPublicLandsAnInRangeSpecAmendment(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-land-spec-amendment"
	root, creation, base, _, tally, _ := publicLandingFixture(t, request, "", "")
	amended := landingSpecAmendment(t, creation.Path)
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land the amended source", creation.Path)
	cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
	if code := exitCode(cmd.Run()); code != 0 || !strings.Contains(stdout.String(), "worktree=released}") {
		t.Fatalf("amended landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	requirePublishedSpec(t, root, published, amended)
	parents := strings.Fields(gitOutput(t, root, "rev-list", "--parents", "-n", "1", published))
	if len(parents) != 3 || parents[1] != base || parents[2] != tip {
		t.Fatalf("published parents = %q, want destination %s and source %s", parents, base, tip)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v", got, err)
	}
}

func TestResumeLandCommandPublicCompletesAnAmendedSourceLanding(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "public-resume-spec-amendment"
	root, creation, base, _, tally, _ := publicLandingFixture(t, request, "private/output", "dist/")
	amended := landingSpecAmendment(t, creation.Path)
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	land := func(args ...string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
		cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}
	code, stdout, stderr := land("--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land the amended source", creation.Path)
	if code != 3 || !strings.Contains(stdout, "worktree=incomplete:release") {
		t.Fatalf("interrupted amended landing = (%d, %q, %q)", code, stdout, stderr)
	}
	published := gitOutput(t, root, "rev-parse", "main")
	requirePublishedSpec(t, root, published, amended)
	if err := os.Remove(filepath.Join(creation.Path, "private", "output")); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr = land("--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path)
	if code != 0 || !strings.Contains(stdout, "worktree=released}") || stderr != "" {
		t.Fatalf("amended resume = (%d, %q, %q)", code, stdout, stderr)
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != published {
		t.Fatalf("resume republished: main=%s, want %s", got, published)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
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
