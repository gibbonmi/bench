package preflight

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// summaryFormat is the approved summary line of BO37. The tests write it apart
// from the render, so a change to the rendered form turns them red.
const summaryFormat = "checks{green=%d,not_applicable=%d,red=%d}"

// verdictCounts is one parsed checks{green,not_applicable,red} summary line.
type verdictCounts struct{ green, na, red int }

func (c verdictCounts) total() int { return c.green + c.na + c.red }

// line is the summary line that the render prints for c.
func (c verdictCounts) line() string {
	return fmt.Sprintf(summaryFormat+"\n", c.green, c.na, c.red)
}

// summaryCounts parses the one summary line of a rendered verdict. A verdict
// with no summary line fails the test, so a caller can grade the counts alone.
func summaryCounts(t *testing.T, out string) verdictCounts {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(line, "checks{") {
			continue
		}
		var c verdictCounts
		if _, err := fmt.Sscanf(line, summaryFormat, &c.green, &c.na, &c.red); err != nil {
			t.Fatalf("summary line %q does not parse: %v", line, err)
		}
		return c
	}
	t.Fatalf("output has no checks{...} summary line:\n%s", out)
	return verdictCounts{}
}

// decidedVerdicts gathers the fixture in the working directory and returns each
// check's verdict by name. A green or not-applicable row prints only as a count,
// so a test grades such a row here.
func decidedVerdicts(t *testing.T, mode, slug string, base ...string) map[string]string {
	t.Helper()
	explicitBase := ""
	if len(base) > 0 {
		explicitBase = base[0]
	}
	return pinnedVerdicts(t, mode, slug, explicitBase, "")
}

// pinnedVerdicts is decidedVerdicts with an explicit base and a source-tip pin.
func pinnedVerdicts(t *testing.T, mode, slug, base, pin string) map[string]string {
	t.Helper()
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	facts, failure := GatherPinned(root, mode, slug, base, pin)
	if failure != nil {
		t.Fatalf("gather %s: %s: %s", mode, failure.Kind, failure.Hint)
	}
	verdicts := map[string]string{}
	for _, check := range Decide(facts).Checks {
		verdicts[check.Check] = check.Verdict
	}
	return verdicts
}

// fixtureCounts counts the verdict rows Decide returns for the fixture, so the
// expected summary follows the check registry and not a pasted literal.
func fixtureCounts(t *testing.T, mode, slug string, base ...string) verdictCounts {
	t.Helper()
	return countVerdicts(decidedVerdicts(t, mode, slug, base...))
}

// countVerdicts counts each verdict class in verdicts.
func countVerdicts(verdicts map[string]string) verdictCounts {
	var c verdictCounts
	for _, verdict := range verdicts {
		switch verdict {
		case verdictGreen:
			c.green++
		case verdictNA:
			c.na++
		case verdictRed:
			c.red++
		}
	}
	return c
}

// renderedVerdicts checks that the summary line of out counts the fixture's
// verdicts, and returns those verdicts by name.
func renderedVerdicts(t *testing.T, out, mode, slug string, base ...string) map[string]string {
	t.Helper()
	if got, want := summaryCounts(t, out), fixtureCounts(t, mode, slug, base...); got != want {
		t.Fatalf("summary = %+v, want the fixture's %+v:\n%s", got, want, out)
	}
	return decidedVerdicts(t, mode, slug, base...)
}

// requireVerdict fails for each named check whose verdict is not want.
func requireVerdict(t *testing.T, verdicts map[string]string, want string, checks ...string) {
	t.Helper()
	for _, check := range checks {
		if verdicts[check] != want {
			t.Errorf("%s = %q, want %q", check, verdicts[check], want)
		}
	}
}

// TestPreflightGreenSummaryLine is BO37 and BO38. An all-green build or review
// preflight prints the phase and spec lines and one summary line, and no
// checks[ table. A render that keeps the green rows prints the table again.
func TestPreflightGreenSummaryLine(t *testing.T) {
	for _, mode := range []string{modeBuild, modeReview} {
		t.Run(mode, func(t *testing.T) {
			_, slug := preflighttest.SeedConformant(t)
			want := fixtureCounts(t, mode, slug)
			if want.red != 0 || want.green == 0 {
				t.Fatalf("the conformant fixture is not all green in %s mode: %+v", mode, want)
			}
			t.Logf("%s fixture counts: green=%d not_applicable=%d red=%d", mode, want.green, want.na, want.red)

			out, code := Command([]string{mode, slug})
			if code != 0 {
				t.Fatalf("%s exit = %d, want 0:\n%s", mode, code, out)
			}
			wantOut := fmt.Sprintf("phase: %s\nspec: specs/%s/spec.md\n", mode, slug) + want.line()
			if out != wantOut {
				t.Fatalf("%s output = %q, want %q", mode, out, wantOut)
			}
		})
	}
}

// TestPreflightRedRowsOnly is BO39. A preflight with one red check prints the
// summary line and a checks[1] table that holds the red row with its detail and
// next action, and no green or not-applicable row.
func TestPreflightRedRowsOnly(t *testing.T) {
	_, slug := preflighttest.SeedConformant(t)
	preflighttest.RunGit(t, "checkout", "-q", "main")
	preflighttest.MustWriteFile(t, "unrelated.txt", "advance main\n")
	preflighttest.RunGit(t, "add", "unrelated.txt")
	preflighttest.RunGit(t, "commit", "-q", "-m", "advance main")
	preflighttest.RunGit(t, "checkout", "-q", "feature")
	want := fixtureCounts(t, modeBuild, slug)
	if want.red != 1 {
		t.Fatalf("the stale-base fixture has %d red rows, want 1", want.red)
	}

	out, code := Command([]string{modeBuild, slug})
	if code != 1 {
		t.Fatalf("exit = %d, want 1:\n%s", code, out)
	}
	wantOut := fmt.Sprintf("phase: build\nspec: specs/%s/spec.md\n", slug) + want.line() +
		"checks[1]{check,verdict,detail,next}:\n" +
		"  base-current,red,default branch tip is not an ancestor of HEAD,bench worktree merge --from main <target>\n"
	if out != wantOut {
		t.Fatalf("output = %q, want %q", out, wantOut)
	}
}
