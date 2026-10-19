package preflight

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// Each test in this file drives the command over the derived ownership fence: the union of
// the ticket `Writes:` entries plus the implicit entries. The expected rows and fence blocks
// are stated here independently of the fence owner.

// fenceSlug is the spec slug that seedFence commits.
const fenceSlug = "example"

// seedFence commits spec and each ticket on main, then commits each change on a feature
// branch, and registers the active assignment over the tree.
func seedFence(t *testing.T, spec string, tickets map[string]string, changes ...string) {
	t.Helper()
	root := preflighttest.StartRepo(t)
	preflighttest.MustWriteFile(t, "specs/"+fenceSlug+"/spec.md", spec)
	for name, doc := range tickets {
		preflighttest.MustWriteFile(t, "specs/"+fenceSlug+"/tickets/"+name, doc)
	}
	commitmenttest.SeedAdmission(t, root, "specs/"+fenceSlug+"/spec.md")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	for _, path := range changes {
		preflighttest.MustWriteFile(t, path, "change\n")
	}
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	preflighttest.ActiveAssignment(t, root, root)
}

// withoutFenceSection removes the `## Ownership fences` section from a spec that
// preflighttest wrote, up to the next level-2 heading.
func withoutFenceSection(t *testing.T, spec string) string {
	t.Helper()
	head, rest, found := strings.Cut(spec, "\n## Ownership fences\n")
	_, tail, closed := strings.Cut(rest, "\n## ")
	if !found || !closed {
		t.Fatalf("spec holds no closed fence section:\n%s", spec)
	}
	return head + "\n## " + tail
}

// authorizedRow is the paths-authorized row that Decide grades over the working tree.
func authorizedRow(t *testing.T, mode string) CheckResult {
	t.Helper()
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	facts, failure := Gather(root, mode, fenceSlug, "")
	if failure != nil {
		t.Fatalf("gather %s: %s: %s", mode, failure.Kind, failure.Hint)
	}
	row, ok := checkRow(Decide(facts), "paths-authorized")
	if !ok {
		t.Fatal("verdict holds no paths-authorized row")
	}
	return row
}

// TestCommandSpecWithoutFenceSection covers CR61 and CR62. A spec with no fence section
// reaches a verdict. Its tickets authorize their own writes, and a tickets folder that is
// empty or absent authorizes no production path.
func TestCommandSpecWithoutFenceSection(t *testing.T) {
	const change = "internal/example/foo.go"
	for _, tc := range []struct {
		name    string
		tickets map[string]string
		empty   bool
		wantRed bool
	}{
		{name: "ticket writes", tickets: map[string]string{"one.md": preflighttest.TicketDoc("One", "PF1", "PF2")}},
		{name: "empty tickets folder", empty: true, wantRed: true},
		{name: "absent tickets folder", wantRed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			seedFence(t, withoutFenceSection(t, preflighttest.SpecBody(fenceSlug)), tc.tickets, change)
			if tc.empty {
				if err := os.MkdirAll("specs/"+fenceSlug+"/tickets", 0o755); err != nil {
					t.Fatal(err)
				}
			}
			out, code := Command([]string{"build", fenceSlug})
			if strings.HasPrefix(out, "error: ") {
				t.Fatalf("build without a fence section = (%d), want a verdict:\n%s", code, out)
			}
			row := authorizedRow(t, modeBuild)
			if !tc.wantRed {
				if code != 0 || row.Verdict != verdictGreen {
					t.Fatalf("build = (%d) with paths-authorized %+v, want exit 0 and a green row:\n%s", code, row, out)
				}
				return
			}
			if code != 1 || row.Verdict != verdictRed || row.Detail != "not authorized by any ownership fence: "+change {
				t.Fatalf("build = (%d) with paths-authorized %+v, want exit 1 and a red row naming %s:\n%s", code, row, change, out)
			}
		})
	}
}

// leftoverFence seeds a spec whose leftover fence section lists internal/leftover/ and a
// ticket that writes only internal/ticketonly/. The feature branch changes a path under
// each, the review pickup, and another review record.
func leftoverFence(t *testing.T) {
	t.Helper()
	spec := preflighttest.SpecBody(fenceSlug, "- `internal/leftover/`")
	ticket := preflighttest.WritesTicketDoc("One", []string{"internal/ticketonly/ (new)"}, "PF1", "PF2")
	seedFence(t, spec, map[string]string{"one.md": ticket},
		"internal/ticketonly/a.go", "internal/leftover/b.go", "reviews/example.md", "reviews/other.md")
}

// unauthorizedLeftover is the paths-authorized detail over leftoverFence: the section
// authorizes nothing, the review pickup is implicit, and no entry covers another review.
const unauthorizedLeftover = "not authorized by any ownership fence: internal/leftover/b.go, reviews/other.md"

// TestCommandBuildIgnoresAFenceSection covers CR72. A path that only a leftover fence
// section lists is unauthorized in build mode.
func TestCommandBuildIgnoresAFenceSection(t *testing.T) {
	leftoverFence(t)
	out, code := Command([]string{"build", fenceSlug})
	if row := authorizedRow(t, modeBuild); code != 1 || row.Verdict != verdictRed || row.Detail != unauthorizedLeftover {
		t.Fatalf("build = (%d) with paths-authorized %+v, want the red %q:\n%s", code, row, unauthorizedLeftover, out)
	}
}

// TestCommandReviewAuthorizesFromTicketWrites covers CR77. Review preflight authorizes a
// ticket-only path and the review pickup, and it reds a path outside every ticket entry.
func TestCommandReviewAuthorizesFromTicketWrites(t *testing.T) {
	leftoverFence(t)
	out, code := Command([]string{"review", fenceSlug})
	if row := authorizedRow(t, modeReview); code != 1 || row.Verdict != verdictRed || row.Detail != unauthorizedLeftover {
		t.Fatalf("review = (%d) with paths-authorized %+v, want the red %q:\n%s", code, row, unauthorizedLeftover, out)
	}
}

// twoTicketSpec is spec with its completion plan over one.md and two.md.
func twoTicketSpec(t *testing.T, spec string) string {
	t.Helper()
	planned := strings.Replace(spec, preflighttest.PlanFence("one.md"), preflighttest.PlanFence("one.md", "two.md"), 1)
	if planned == spec {
		t.Fatalf("spec holds no one-ticket plan:\n%s", spec)
	}
	return planned
}

// requireFence fails unless the published pack of a prepared charge lists want as its
// fence block.
func requireFence(t *testing.T, root, out string, code int, want []string) {
	t.Helper()
	if code != 0 {
		t.Fatalf("charge = (%d), want 0:\n%s", code, out)
	}
	if got := preflighttest.PublishedPack(t, root, preparedIdentity(t, out)).Metadata().Fence; !slices.Equal(got, want) {
		t.Fatalf("fence block =\n%q\nwant\n%q", got, want)
	}
}

// TestChargeFenceListsTheTicketUnion covers CR63. The build charge fence block lists each
// path that the two tickets write, once, in ticket and entry order, with no trailing slash
// and no marker.
func TestChargeFenceListsTheTicketUnion(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	two := preflighttest.WritesTicketDoc("Two", []string{"internal/example (new)", "internal/other/ (new)", "internal/other/x.go (new)", "reviews/example.md (new)"}, "PF2")
	args := preflighttest.CommitChargeSource(t, root, slug, map[string]string{
		"specs/" + slug + "/spec.md":        twoTicketSpec(t, preflighttest.SpecBody(slug)),
		"specs/" + slug + "/tickets/two.md": two,
	})
	out, code := Command(args)
	requireFence(t, root, out, code, []string{
		"internal/example", "reviews/example.md", ".agents/skills/bench-craft-delegate", ".agents/commands/bench-implement-spec.md",
		"internal/other", "internal/other/x.go",
	})
}

// TestReviewChargeFenceListsTheTicketUnion covers CR64. The review charge fence block lists
// each path that the two tickets write, once, in ticket and entry order.
func TestReviewChargeFenceListsTheTicketUnion(t *testing.T) {
	root, slug, _ := preflighttest.SeedReviewEvidence(t, false)
	spec := "specs/" + slug + "/spec.md"
	body, err := os.ReadFile(spec)
	if err != nil {
		t.Fatal(err)
	}
	preflighttest.MustWriteFile(t, spec, twoTicketSpec(t, string(body)))
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/two.md", preflighttest.WritesTicketDoc("Two", []string{"target (new)", "notes/extra.txt (new)", "reviews/example.md (new)"}, "PF2"))
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", "second ticket")
	out, code := Command(preflighttest.ReviewArgs(t, slug))
	requireFence(t, root, out, code, []string{
		"internal/example", "reviews/example.md", ".agents/skills/bench-craft-delegate", ".agents/commands/bench-implement-spec.md",
		"target", "edited", "outside", "notes", ".agents/skills/bench-craft-review", ".agents/commands/bench-review-implementation.md",
		"notes/extra.txt",
	})
}
