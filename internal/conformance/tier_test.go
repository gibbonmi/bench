package conformance

// The tier registry itself: which checks each tier runs, that metadata and bound
// functions cannot drift apart, that the filtered inner run selects real tests, and
// that every executed check leaves one timing line in a stable order.

import (
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

const releaseEvidenceProbeCheck = "release-evidence-probe"

func TestTierMembership(t *testing.T) {
	dev, ship := registry.Names(registry.Dev), registry.Names(registry.Ship)
	if len(dev) == 0 {
		t.Fatal("the dev tier runs no checks at all")
	}
	if slices.Contains(dev, releaseEvidenceProbeCheck) {
		t.Fatalf("the dev tier still runs %s, the ~372 s probe the split moves:\n%s", releaseEvidenceProbeCheck, strings.Join(dev, "\n"))
	}
	if !slices.Contains(ship, releaseEvidenceProbeCheck) {
		t.Fatalf("the ship tier does not run %s, so the probe runs nowhere:\n%s", releaseEvidenceProbeCheck, strings.Join(ship, "\n"))
	}
	for _, name := range dev {
		if !slices.Contains(ship, name) {
			t.Fatalf("the ship tier drops dev check %s; ship green has to reprove everything dev green claims", name)
		}
	}
}

// TestEveryCheckCarriesATier closes the gap Tier's string underlying type leaves open:
// a row whose tier is misspelled or omitted holds "", which RunsAt reads as neither dev
// nor ship, so the check silently stops running on every commit and survives every
// membership assertion — those compare against name lists the untiered check is in or
// out of on both sides.
func TestEveryCheckCarriesATier(t *testing.T) {
	for _, check := range registry.Checks {
		switch check.Tier {
		case registry.Dev, registry.Ship:
		default:
			t.Errorf("registry check %s carries tier %q, which is neither %q nor %q, so no tier executes it on a commit", check.Name, check.Tier, registry.Dev, registry.Ship)
		}
	}
}

func TestRegistryBindsEveryCheck(t *testing.T) {
	for _, check := range registry.Checks {
		if _, bound := conformanceChecks[check.Name]; !bound {
			t.Errorf("registry check %s has no bound function", check.Name)
		}
	}
	for name := range conformanceChecks {
		if !slices.Contains(registry.Names(registry.Ship), name) {
			t.Errorf("bound function %s has no registry row, so it carries no tier", name)
		}
	}
}

func TestDevTierExecutesExactlyDevChecks(t *testing.T) {
	root := gitInitedRoot(t)
	RunConformance(root, NewHarness(t).KitRoot, registry.Dev, "")

	got := timingNames(t, root)
	if want := registry.Names(registry.Dev); !slices.Equal(got, want) {
		t.Fatalf("dev run executed\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if slices.Contains(got, releaseEvidenceProbeCheck) {
		t.Fatalf("a dev run executed %s", releaseEvidenceProbeCheck)
	}
}

func TestTimingLinePerCheck(t *testing.T) {
	root := gitInitedRoot(t)
	RunConformance(root, NewHarness(t).KitRoot, registry.Dev, "")

	lines := registry.ReadTimingLines(root)
	if want := len(registry.Names(registry.Dev)); len(lines) != want {
		t.Fatalf("timing lines = %d, want one per executed check (%d):\n%s", len(lines), want, strings.Join(lines, "\n"))
	}
}

func TestTimingOrderStable(t *testing.T) {
	root := gitInitedRoot(t)
	kitRoot := NewHarness(t).KitRoot

	RunConformance(root, kitRoot, registry.Dev, "")
	first := timingNames(t, root)
	RunConformance(root, kitRoot, registry.Dev, "")
	second := timingNames(t, root)

	if !slices.Equal(first, second) {
		t.Fatalf("timing order differs between runs of one tree:\n%s\nversus\n%s", strings.Join(first, "\n"), strings.Join(second, "\n"))
	}
	if want := registry.Names(registry.Dev); !slices.Equal(first, want) {
		t.Fatalf("timing order is not the registry's order:\n%s\nwant\n%s", strings.Join(first, "\n"), strings.Join(want, "\n"))
	}
}

func TestScopedRunExecutesOnlyTheNamedCheck(t *testing.T) {
	root := gitInitedRoot(t)
	const scope = "line-routing"
	RunConformance(root, NewHarness(t).KitRoot, registry.Dev, scope)

	if got := timingNames(t, root); !slices.Equal(got, []string{scope}) {
		t.Fatalf("scoped run executed\n%s\nwant only %s", strings.Join(got, "\n"), scope)
	}
}

// TestUnknownScopeIsRedAndRunsNothing pins the posture a silent fallback would break:
// a scope naming no check re-pays the full run and hides the drift that renamed it.
// The hostile value is deliberate — the diagnostic quotes the scope, so control bytes
// have to survive as escapes rather than as a mangled line.
func TestUnknownScopeIsRedAndRunsNothing(t *testing.T) {
	root := recordedScopeRoot(t)
	const scope = "no-such-check\x01\n"
	diags := RunConformance(root, NewHarness(t).KitRoot, registry.Dev, scope)

	if len(diags) != 1 || !strings.Contains(diags[0], `"no-such-check\x01\n"`) {
		t.Fatalf("unknown scope: want one diagnostic quoting the scope, got %q", diags)
	}
	if got := timingNames(t, root); len(got) != 0 {
		t.Fatalf("unknown scope left timing lines standing:\n%s", strings.Join(got, "\n"))
	}
}

// TestScopeOutsideTierIsRedAndRunsNothing covers the scope that exists but sits on a
// tier this run does not grade: executing zero checks in silence would read as green
// and leave the fixture reporting a baffling did-not-bite.
func TestScopeOutsideTierIsRedAndRunsNothing(t *testing.T) {
	root := recordedScopeRoot(t)
	diags := RunConformance(root, NewHarness(t).KitRoot, registry.Dev, releaseEvidenceProbeCheck)

	if len(diags) != 1 || !strings.Contains(diags[0], releaseEvidenceProbeCheck) || !strings.Contains(diags[0], string(registry.Dev)) {
		t.Fatalf("tier-mismatched scope: want one diagnostic naming the scope and the tier, got %q", diags)
	}
	if got := timingNames(t, root); len(got) != 0 {
		t.Fatalf("tier-mismatched scope left timing lines standing:\n%s", strings.Join(got, "\n"))
	}
}

// recordedScopeRoot is a graded root whose timing file already holds a completed run's
// lines. Against a pristine root a red posture cannot be told from one that left the
// last run's record standing; against this one, only an emptied file passes.
func recordedScopeRoot(t *testing.T) string {
	t.Helper()
	root := gitInitedRoot(t)
	RunConformance(root, NewHarness(t).KitRoot, registry.Dev, "line-routing")
	if len(registry.ReadTimingLines(root)) == 0 {
		t.Fatal("the seeding run recorded no timing lines, so the posture assertion proves nothing")
	}
	return root
}

// TestFilteredRunSelectsRealTests keeps the skip list from rotting in either
// direction. `go test -list` ignores -skip, so it can only supply the inventory the
// listed names must exist in; compiling the pattern here is the only way to observe
// that it leaves the cheap tests the filtered run exists to keep running.
func TestFilteredRunSelectsRealTests(t *testing.T) {
	h := NewHarness(t)
	probe := runAtCleanEnv(h.KitRoot, "go", "test", "./"+registry.ConformancePackage, "-list", ".*")
	if probe == nil || probe.ExitCode != 0 {
		t.Fatalf("listing the conformance tests failed: %s", formatProbeFailure("go test -list failed", probe, ""))
	}
	inventory := strings.Fields(probe.Stdout)
	for _, name := range registry.InnerSkipTests {
		if !slices.Contains(inventory, name) {
			t.Fatalf("skip list names %s, which is not a test in the conformance package", name)
		}
	}

	pattern, err := regexp.Compile(registry.InnerSkipPattern())
	if err != nil {
		t.Fatalf("compile skip pattern %q: %v", registry.InnerSkipPattern(), err)
	}
	for _, name := range []string{"TestRegistryBindsEveryCheck", "TestPackageCoreAndGuardFixturesBite"} {
		if pattern.MatchString(name) {
			t.Fatalf("skip pattern %q also excludes %s, so the filtered run drops tests the oracle keeps", registry.InnerSkipPattern(), name)
		}
	}
}

// gitInitedRoot is a graded root a timing file can live in: the file sits under the
// root's own git dir, so a bare temp dir gives it nowhere to go.
func gitInitedRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	return root
}

func timingNames(t *testing.T, root string) []string {
	t.Helper()
	var names []string
	for _, line := range registry.ReadTimingLines(root) {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			t.Fatalf("timing line %q is not an index, a check name, and a duration", line)
		}
		names = append(names, fields[1])
	}
	return names
}
