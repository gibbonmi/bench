package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// floorOmission is the check-floor detail of the seeded package that the seeded plan omits.
// The seeded plan runs the whole suite, so the fixture replaces it with a foreign command.
const floorOmission = "completion plan omits a floor check: one.md: ./internal/example (add bench test --package ./internal/example)"

// seedFloorRepo seeds the conformant fixture, folds a Go module into its base, and plans
// a command that satisfies nothing. The ticket's package is then a floor element that no
// requirement satisfies.
func seedFloorRepo(t *testing.T) (root, slug string) {
	t.Helper()
	root, slug = preflighttest.SeedConformant(t)
	t.Setenv("BENCH_KIT", "")
	preflighttest.RunGit(t, "checkout", "-q", "main")
	preflighttest.MustWriteFile(t, "go.mod", "module example.com/floor\n\ngo 1.25\n")
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", strings.Replace(preflighttest.SpecBody(slug),
		"bench test --package ./...", "go test ./...", 1))
	// The approved policy binds the spec bytes, so the edited spec is approved again.
	commitmenttest.SeedAdmission(t, root, "specs/"+slug+"/spec.md")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "module")
	preflighttest.RunGit(t, "checkout", "-q", "feature")
	preflighttest.RunGit(t, "merge", "-q", "--no-edit", "main")
	return root, slug
}

// requireFloorLine fails unless out renders the red check-floor row with detail.
func requireFloorLine(t *testing.T, out, detail string) {
	t.Helper()
	if want := "  check-floor,red,\"" + detail + "\",\"\"\n"; !strings.Contains(out, want) {
		t.Fatalf("output missing the check-floor row %q:\n%s", want, out)
	}
}

func TestCheckFloorCommandRunsAfterTheTicket(t *testing.T) {
	_, slug := seedFloorRepo(t)
	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	requireFloorLine(t, out, floorOmission)
}

func TestCheckFloorCommandRunsAtPlanOnly(t *testing.T) {
	_, slug := seedFloorRepo(t)
	out, code := Command([]string{"build", slug, "--plan-only"})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	requireFloorLine(t, out, floorOmission)
}

func TestCheckFloorCommandRunsInReview(t *testing.T) {
	_, slug := seedFloorRepo(t)
	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	requireFloorLine(t, out, floorOmission)
}

func TestCheckFloorRefusesTheBuildCharge(t *testing.T) {
	root, slug := seedFloorRepo(t)
	out, code := Command(preflighttest.ChargeArgs(t, root, slug))
	want := "error: preflight required: check-floor: " + floorOmission + " — repair check-floor and rerun the exact charge\n"
	if code != 1 || !strings.Contains(out, want) {
		t.Fatalf("build charge = (%d):\n%s\nwant the refusal %q", code, out, want)
	}
}

// TestCheckFloorDetailRendersEscaped reads the rendered red row whole: it stays on one line,
// and it holds no raw control byte.
func TestCheckFloorDetailRendersEscaped(t *testing.T) {
	_, slug := seedFloorRepo(t)
	out, _ := Command([]string{"build", slug, "--plan-only"})
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(line, "  check-floor,red,") {
			continue
		}
		if at := strings.IndexFunc(line, func(r rune) bool { return r < 0x20 || r == 0x7f }); at >= 0 {
			t.Fatalf("check-floor row holds the raw control byte %q at %d: %q", line[at], at, line)
		}
		return
	}
	t.Fatalf("output renders no red check-floor row:\n%s", out)
}
