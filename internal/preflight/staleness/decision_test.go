package staleness_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// Each test in this file drives the actual preflight command. The expected decision of each
// case is stated here by hand, independently of the staleness renderer.

// seedClean commits the complete supported inventory of the completeness fixture, whose audit
// at its unchanged tip is the clean decision. It returns the root, the slug, and the build
// charge arguments, whose tip is element 8.
func seedClean(t *testing.T, manual func(t *testing.T, slug string) []string) (root, slug string, charge []string) {
	t.Helper()
	root, slug = preflighttest.SeedConformant(t)
	var claims []string
	if manual != nil {
		claims = manual(t, slug)
	}
	return root, slug, commitCompleteness(t, root, slug, inventory("declared-complete", keptSeams(slug, true), claims))
}

// decision is the staleness row of one audit response as complete|manual_required.
func decision(t *testing.T, out string) string {
	t.Helper()
	row := table(t, out, "staleness")[0]
	return fmt.Sprint(row["complete"], "|", row["manual_required"])
}

// specChecks is each spec_checks row of one audit response as check|state.
func specChecks(t *testing.T, out string) string {
	t.Helper()
	var rows []string
	for _, row := range table(t, out, "spec_checks") {
		rows = append(rows, fmt.Sprint(row["check"], "|", row["state"]))
	}
	return strings.Join(rows, " ")
}

// breakSeal commits an ignore rule for the published binary and then publishes a binary that
// no seal verifies, so the binary-seal row is red. It returns the new head.
func breakSeal(t *testing.T, root string) string {
	t.Helper()
	write(t, map[string]string{".gitignore": "dist/\n"})
	tip := commitAll(t, "ignore the published binary")
	executable := freshness.PublishedExecutable(root)
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "absent-bench"), executable); err != nil {
		t.Fatal(err)
	}
	return tip
}

// allGreen is the spec_checks rows of an audit whose three mechanical checks are green.
const allGreen = "spec-drift|green cited-symbols|green cited-tests|green"

// TestAuditDecision covers SS41, SS42, SS43, and SS44. Each case runs the actual audit-spec
// form over the clean fixture or one change of it. Only the clean case may omit manual review.
// Each other case keeps the three mechanical checks as stated, so only its named cause holds
// manual review required: an audit red, a readiness red, a commitment refusal, an unsupported
// claim, or a required build row that is not applicable.
func TestAuditDecision(t *testing.T) {
	t.Run("SS41 complete clean evidence", func(t *testing.T) {
		_, slug, charge := seedClean(t, nil)
		out, code := audit(t, slug, charge[8], charge[8])
		if code != 0 || decision(t, out) != "true|false" || specChecks(t, out) != allGreen ||
			!strings.Contains(out, "drift[0]") || !strings.Contains(out, "audit_diagnostics[0]") || !strings.Contains(out, "red=0}") {
			t.Fatalf("clean audit = (%d), want exit 0, complete true, and manual_required false:\n%s", code, out)
		}
	})
	t.Run("SS42 supported audit red", func(t *testing.T) {
		// The charge base precedes the fixture's fenced change, so only spec-drift is red.
		_, slug, charge := seedClean(t, nil)
		out, code := audit(t, slug, charge[6], charge[8])
		if code != 1 || decision(t, out) != "true|true" || specChecks(t, out) != strings.Replace(allGreen, "spec-drift|green", "spec-drift|red", 1) || !strings.Contains(out, "red=0}") {
			t.Fatalf("audit red = (%d), want exit 1, only a red spec-drift row, and manual review:\n%s", code, out)
		}
	})
	t.Run("SS42 supported readiness red", func(t *testing.T) {
		_, slug, _ := seedClean(t, nil)
		body, err := os.ReadFile(specPath(slug))
		if err != nil {
			t.Fatal(err)
		}
		write(t, map[string]string{specPath(slug): strings.Replace(string(body), "| SS28 |", "| SS29 | 1 | unowned | cli seam | catches q |\n| SS28 |", 1)})
		tip := commitAll(t, "an unowned row")
		out, code := audit(t, slug, tip, tip)
		if code != 1 || decision(t, out) != "true|true" || specChecks(t, out) != allGreen || !strings.Contains(out, "rows-owned,red") {
			t.Fatalf("readiness red = (%d), want exit 1, a red rows-owned row, green audit checks, and manual review:\n%s", code, out)
		}
	})
	t.Run("SS42 red seal", func(t *testing.T) {
		// The seal's own rule exempts it only while it is not applicable, never while it is red.
		root, slug, _ := seedClean(t, nil)
		tip := breakSeal(t, root)
		out, code := audit(t, slug, tip, tip)
		if code != 1 || decision(t, out) != "true|true" || specChecks(t, out) != allGreen || !strings.Contains(out, "binary-seal,red") {
			t.Fatalf("red seal = (%d), want exit 1, a red binary-seal row, green audit checks, and manual review:\n%s", code, out)
		}
	})
	t.Run("SS42 commitment refusal", func(t *testing.T) {
		root, slug, charge := seedClean(t, nil)
		preflighttest.OwnedAssignment(t, root, root, intent.StateComplete)
		out, code := audit(t, slug, charge[8], charge[8])
		if code != 1 || decision(t, out) != "true|true" || specChecks(t, out) != allGreen || !strings.Contains(out, "commitment,red") {
			t.Fatalf("commitment refusal = (%d), want exit 1, a red commitment row, green audit checks, and manual review:\n%s", code, out)
		}
	})
	t.Run("SS43 diagnostic only", func(t *testing.T) {
		_, slug, charge := seedClean(t, nativeClaims)
		out, code := audit(t, slug, charge[8], charge[8])
		if code != 0 || decision(t, out) != "false|true" || !strings.Contains(out, "red=0}") {
			t.Fatalf("diagnostic-only audit = (%d), want exit 0, complete false, and manual review:\n%s", code, out)
		}
		// Each unsupported claim stays named in the response.
		subjects := map[string]bool{}
		for _, row := range table(t, out, "audit_diagnostics") {
			subjects[fmt.Sprint(row["subject"])] = true
		}
		for _, row := range nativeRows {
			if !subjects["native-"+row[0]] {
				t.Errorf("audit_diagnostics lacks the unsupported claim native-%s:\n%s", row[0], out)
			}
		}
	})
	t.Run("SS44 not-applicable required row", func(t *testing.T) {
		_, slug, _ := seedClean(t, nil)
		preflighttest.RunGit(t, "rm", "-rq", "specs/"+slug+"/tickets")
		tip := commitAll(t, "no tickets")
		out, code := audit(t, slug, tip, tip)
		if code != 0 || decision(t, out) != "true|true" || specChecks(t, out) != allGreen || !strings.Contains(out, "red=0}") {
			t.Fatalf("absent tickets = (%d), want exit 0, no red row, green audit checks, and manual review:\n%s", code, out)
		}
	})
}

// TestAuditApplicability covers SS45 and SS46. Only the audit-spec form prints a staleness
// decision. The ordinary verdict, the plan-only form, the review verdict, the build charge,
// and the current-evidence binding over the same clean source print none, so no other
// operation grants a manual-review omission. After a predecessor ships its own planned test,
// the ordinary preflight stays usable, and the audit-spec form over the same pair reports it.
func TestAuditApplicability(t *testing.T) {
	t.Run("SS46 no other operation decides", func(t *testing.T) {
		root, slug, charge := seedClean(t, nil)
		if out, code := audit(t, slug, charge[8], charge[8]); code != 0 || decision(t, out) != "true|false" {
			t.Fatalf("clean audit = (%d), want the omission decision:\n%s", code, out)
		}
		prepared, code := preflight.Command(charge)
		if code != 0 {
			t.Fatalf("charge = (%d):\n%s", code, prepared)
		}
		identity := table(t, prepared, "prepared")[0]["evidence"].(string)
		for _, args := range [][]string{
			{"build", slug, "--base", charge[6], "--source-tip", charge[8]},
			{"build", slug, "--plan-only", "--base", charge[6], "--source-tip", charge[8]},
			{"review", slug, "--base", charge[6], "--source-tip", charge[8]},
			charge,
			{"evidence", identity, "--check-current"},
		} {
			out, code := preflight.CommandWithVersion("staleness-test")(args)
			if code != 0 || strings.Contains(out, "staleness[") || strings.Contains(out, "manual_required") {
				t.Errorf("%v = (%d), want exit 0 and no staleness decision:\n%s", args, code, out)
			}
		}
		if packs := preflighttest.PublishedPacks(t, root); len(packs) != 1 {
			t.Errorf("published packs = %d, want the one charge", len(packs))
		}
	})
	t.Run("SS45 ordinary after-ticket preflight", func(t *testing.T) {
		_, slug := preflighttest.SeedAdmitted(t)
		const planned = "internal/example/planned_test.go"
		write(t, map[string]string{specPath(slug): inventorySpec(slug, nil, []string{
			seamItem(slug, "planned", planned, "example.com/internal/example", "TestPlanned", "planned", "PF1"),
		})})
		base := commitAll(t, "planned test")
		write(t, map[string]string{planned: "package example\n\nimport \"testing\"\n\nfunc TestPlanned(t *testing.T) {}\n"})
		tip := commitAll(t, "the predecessor ships its planned test")
		out, code := preflight.CommandWithVersion("staleness-test")([]string{"build", slug, "--base", base, "--source-tip", tip})
		if code != 0 || strings.Contains(out, "already-shipped") || strings.Contains(out, "staleness[") {
			t.Fatalf("ordinary preflight = (%d), want exit 0 with no initial planned-test audit:\n%s", code, out)
		}
		if out, code := audit(t, slug, base, tip); code != 1 || !strings.Contains(out, "test TestPlanned already-shipped") || decision(t, out) == "true|false" {
			t.Fatalf("audit over the same pair = (%d), want the already-shipped drift and manual review:\n%s", code, out)
		}
	})
}
