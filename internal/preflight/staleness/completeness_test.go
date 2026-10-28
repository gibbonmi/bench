package staleness_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/preflight/staleness"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// TestCoverageIdentity covers SS28, SS29, SS30, and SS31. One committed pair changes the
// PF1 predicate, removes PF2, and adds the unowned PF4 to the map. It rewrites the Covers of
// one.md to the phantom PF3, adds the ticket two.md, and adds a non-ticket asset. The actual
// audit-spec response grades the pair. The cited-symbols and cited-tests checks stay green, so
// they cannot mask the red ownership rows. The expected drift is the complete list stated here
// by hand, so a lost ticket member, a lost coverage entry, or an asset read as a ticket fails.
func TestCoverageIdentity(t *testing.T) {
	root, slug := preflighttest.SeedAdmitted(t)
	spec, one, two, asset := specPath(slug), "specs/"+slug+"/tickets/one.md", "specs/"+slug+"/tickets/two.md", "specs/"+slug+"/tickets/notes.txt"
	base := preflighttest.RunGit(t, "rev-parse", "HEAD")
	write(t, map[string]string{
		spec: preflighttest.WithChunkTable(strings.Replace(preflighttest.SpecBody(slug), "| PF1 | 1 | does x | cli seam | catches z |\n| PF2 | 1 | does y | cli seam | catches w |\n",
			"| PF1 | 1 | does x again | cli seam | catches z |\n| PF4 | 1 | does z | cli seam | catches v |\n", 1), preflighttest.ChunkTable("PF1", "PF4")),
		one:   preflighttest.TicketDoc("One", "PF1", "PF3"),
		two:   preflighttest.TicketDoc("Two", "PF1"),
		asset: "Covers: PF9\n",
	})
	tip := commitAll(t, "coverage change")

	out, code := audit(t, slug, base, tip)
	if code != 1 || !strings.Contains(out, "spec-drift,red") || !strings.Contains(out, "cited-symbols,green") || !strings.Contains(out, "cited-tests,green") {
		t.Fatalf("audit = (%d), want exit 1, red spec-drift, and green symbol and test checks:\n%s", code, out)
	}
	reds := map[string]string{}
	for _, row := range table(t, out, "checks") {
		reds[fmt.Sprint(row["check"])] = fmt.Sprint(row["detail"])
	}
	modified := func(path string) string {
		return path + "|" + path + "|modified|" + object(t, base, path) + "|" + object(t, tip, path)
	}
	added := func(path string) string {
		return absent + "|" + path + "|added|" + absent + "|" + object(t, tip, path)
	}
	drift := []string{
		spec + "|" + spec + "|coverage PF1 predicate|does x|does x again",
		spec + "|" + absent + "|coverage PF2 predicate|does y|" + absent,
		absent + "|" + spec + "|coverage PF4 predicate|" + absent + "|does z",
		one + "|" + one + "|ticket one.md covers|PF1 PF2|PF1 PF3",
		modified(one),
		absent + "|" + two + "|ticket two.md covers|" + absent + "|PF1",
		added(two),
	}
	for _, test := range []struct{ row, check, detail, drift string }{
		{"SS28", "rows-owned", "PF4", ""},
		{"SS29", "rows-membership", "PF3", ""},
		{"SS30", "", "", drift[0]},
		{"SS30", "", "", drift[1]},
		{"SS30", "", "", drift[2]},
		{"SS31", "", "", drift[3]},
		{"SS31", "", "", drift[4]},
		{"SS31", "", "", drift[5]},
	} {
		t.Run(test.row, func(t *testing.T) {
			if detail, red := reds[test.check]; test.check != "" && (!red || !strings.Contains(detail, test.detail)) {
				t.Errorf("%s red = %t with detail %q, want a red naming %s", test.check, red, detail, test.detail)
			}
			if test.drift != "" && !slices.Contains(driftRows(t, out), test.drift) {
				t.Errorf("drift lacks %q", test.drift)
			}
		})
	}
	// The asset changes as a member of the tickets folder, but it holds no Covers field.
	if got, want := driftRows(t, out), sorted(append(drift, modified(spec), added(asset))...); !slices.Equal(got, want) {
		t.Errorf("drift = %v, want %v", got, want)
	}

	// A ticket object over the control-record bound at either identity does not read, so the
	// audit refuses with no staleness result.
	small, large := preflighttest.TicketDoc("One", "PF1"), preflighttest.TicketDoc("One", "PF1")+overBound()
	write(t, map[string]string{one: large})
	oversized := commitAll(t, "oversized ticket")
	write(t, map[string]string{one: small})
	shrunk := commitAll(t, "ticket shrinks")
	if out, code := audit(t, slug, oversized, shrunk); code != 1 || !strings.Contains(out, "ticket source refused") || strings.Contains(out, "staleness[") {
		t.Errorf("oversized baseline ticket = (%d), want the ticket source refusal and no staleness result:\n%s", code, out)
	}
	// The command reads the tickets of its clean checkout at the tip before it collects, so
	// only Collect itself reaches a tip ticket object that does not read.
	write(t, map[string]string{one: large})
	grown := commitAll(t, "ticket grows")
	var refusal *staleness.Error
	if _, err := staleness.Collect(root, ownership.Source{Base: shrunk, Tip: grown, SpecPath: spec}, ownership.Snapshot{}); !errors.As(err, &refusal) ||
		refusal.Stage != staleness.StageTicket || refusal.Subject != one {
		t.Errorf("oversized tip ticket = %v, want the ticket source refusal of %s", err, one)
	}
}

// nativeTicket is the ticket that owns the native rows, and nativeRows are the five rows
// that this spec's own inventory holds as pending native-future entries, each with its
// planned verification.
const nativeTicket = "05-consume-exact-manual-omission.md"

var nativeRows = [][2]string{{"SS47", "ss-native-diagnostic"}, {"SS48", "ss-native-diagnostic"}, {"SS53", "ss-native-clean"}, {"SS54", "ss-native-diagnostic"}, {"SS55", "ss-native-diagnostic"}}

// keptTest is the test file whose two functions the SS28 seam cell cites.
const keptTest = "internal/example/kept_test.go"

// completenessBody is the seeded spec before its inventory and plan. Its coverage map holds
// SS28, whose seam cell cites two tests, and the native rows, whose cells cite none.
func completenessBody(slug string) string {
	rows, ids := "| SS28 | 1 | does x | `"+keptTest+"` (`TestKept`, `TestOther`) | catches z |\n", []string{"SS28"}
	for _, row := range nativeRows {
		rows += "| " + row[0] + " | 1 | native " + row[0] + " | fresh-session adoption | catches w |\n"
		ids = append(ids, row[0])
	}
	body := strings.TrimSuffix(preflighttest.SpecBodyNoInventory(slug), preflighttest.PlanFence("one.md"))
	body = preflighttest.WithChunkTable(body, preflighttest.ChunkTable(ids...))
	return strings.Replace(body, "| PF1 | 1 | does x | cli seam | catches z |\n| PF2 | 1 | does y | cli seam | catches w |\n", rows, 1)
}

// keptSeams are the inventory seams of both SS28 citations, or of the first one only.
func keptSeams(slug string, both bool) []string {
	seams := []string{seamItem(slug, "kept", keptTest, "example.com/internal/example", "TestKept", "existing", "SS28")}
	if both {
		seams = append(seams, seamItem(slug, "other", keptTest, "example.com/internal/example", "TestOther", "existing", "SS28"))
	}
	return seams
}

// nativeClaims is one native-future manual claim for each native row, located at that row's
// line of the seeded spec, owned by the native ticket and chunk, and pinned to the build phase.
func nativeClaims(t *testing.T, slug string) []string {
	t.Helper()
	lines := strings.Split(completenessBody(slug), "\n")
	var claims []string
	for _, row := range nativeRows {
		needle := "| " + row[0] + " |"
		line := slices.IndexFunc(lines, func(l string) bool { return strings.HasPrefix(l, needle) }) + 1
		claim, err := json.Marshal(map[string]any{
			"id":       "native-" + row[0],
			"location": map[string]any{"path": specPath(slug), "line": line, "needle": needle, "source": "spec"},
			"reason":   "unobserved native workflow consumption", "applicability": "native-future",
			"row": row[0], "chunk": "SS-C3", "ticket": nativeTicket, "verification": row[1], "observation": "unobserved",
			"phase_owner": map[string]any{"path": chargesource.BuildPhase, "line": 3, "needle": "Run root conformance", "source": "tip"},
		})
		if err != nil {
			t.Fatal(err)
		}
		claims = append(claims, string(claim))
	}
	return claims
}

// commitCompleteness commits the files of completenessFiles as the charge source of the
// active, admitted fixture.
func commitCompleteness(t *testing.T, root, slug string, sections ...string) []string {
	t.Helper()
	return preflighttest.CommitChargeSource(t, root, slug, completenessFiles(slug, sections...))
}

// completenessFiles is the seeded spec with sections as its inventory, with the tickets and
// plan that own every row.
func completenessFiles(slug string, sections ...string) map[string]string {
	plan := reviewrecord.Plan{Version: 1, FinalVerification: []reviewrecord.Requirement{{ID: "acceptance", Command: "go test ./..."}}, Chunks: []reviewrecord.PlannedChunk{
		{ID: "c1", Tickets: []string{"one.md"}, Verification: []reviewrecord.Requirement{{ID: "tests", Command: preflighttest.PlanCommand, Probe: preflighttest.PlanProbe}}},
		{ID: "SS-C3", Tickets: []string{nativeTicket}, Verification: []reviewrecord.Requirement{
			{ID: "ss-native-diagnostic", Command: preflighttest.PlanCommand, Probe: preflighttest.PlanProbe}, {ID: "ss-native-clean", Command: preflighttest.PlanCommand, Probe: preflighttest.PlanProbe}}},
	}}
	var native []string
	for _, row := range nativeRows {
		native = append(native, row[0])
	}
	return map[string]string{
		specPath(slug):                               completenessBody(slug) + strings.Join(sections, "") + preflighttest.PlanSection(plan),
		"specs/" + slug + "/tickets/one.md":          preflighttest.TicketDoc("One", "SS28"),
		"specs/" + slug + "/tickets/" + nativeTicket: preflighttest.WritesTicketDoc("Native", []string{"internal/example/ (new)"}, native...),
		keptTest: "package example\n\nimport \"testing\"\n\nfunc TestKept(t *testing.T) {}\n\nfunc TestOther(t *testing.T) {}\n",
	}
}

// inventory is the planning inventory section that declares completeness over seams and the
// manual claims.
func inventory(completeness string, seams, manual []string) string {
	return inventorytest.Section(inventorytest.Manifest(completeness, nil, seams, nil, manual))
}

// TestAuditCompleteness covers SS32, SS33, SS34, SS49, SS50, SS57, and SS58. Each case seeds
// one inventory and runs the actual audit-spec form over its unchanged tip, so no drift and
// no red holds the audit open. The complete case proves that the fixture alone admits a
// complete audit, so each other case fails only through the inventory fault it names. The
// expected diagnostics are stated here by hand.
func TestAuditCompleteness(t *testing.T) {
	for _, test := range []struct {
		row      string
		sections func(t *testing.T, slug string) []string
		want     []string
	}{
		{"complete", func(t *testing.T, slug string) []string {
			return []string{inventory("declared-complete", keptSeams(slug, true), nil)}
		}, nil},
		{"SS32 missing", func(*testing.T, string) []string { return nil },
			[]string{"metadata|no bench-planning-facts fence"}},
		{"SS33 partial", func(t *testing.T, slug string) []string {
			return []string{inventory("partial", keptSeams(slug, true), nil)}
		}, []string{"bench-planning-facts|the planning inventory is not complete: metadata: the premise inventory is partial"}},
		{"SS34 omitted citation", func(t *testing.T, slug string) []string {
			return []string{inventory("declared-complete", keptSeams(slug, false), nil)}
		}, []string{"metadata|row SS28 seam omits " + keptTest + " (TestOther)"}},
		{"SS49 duplicate", func(t *testing.T, slug string) []string {
			section := inventory("declared-complete", keptSeams(slug, true), nil)
			return []string{section, section}
		}, []string{"metadata|bench-planning-facts fence appears more than once"}},
		// The empty inventory also omits both citations, and its own diagnostic differs from
		// the missing one.
		{"SS50 present empty", func(*testing.T, string) []string { return []string{inventorytest.Section(inventorytest.CompleteEmpty)} },
			[]string{"metadata|row SS28 seam omits " + keptTest + " (TestKept)", "metadata|row SS28 seam omits " + keptTest + " (TestOther)",
				"bench-planning-facts|the planning inventory declares no claim"}},
		{"SS57 SS58 pending native", func(t *testing.T, slug string) []string {
			return []string{inventory("declared-complete", keptSeams(slug, true), nativeClaims(t, slug))}
		}, []string{"native-SS47|native-future manual claim", "native-SS48|native-future manual claim",
			"native-SS53|native-future manual claim", "native-SS54|native-future manual claim", "native-SS55|native-future manual claim"}},
	} {
		t.Run(test.row, func(t *testing.T) {
			root, slug := preflighttest.SeedConformant(t)
			charge := commitCompleteness(t, root, slug, test.sections(t, slug)...)
			tip := charge[8]
			out, code := audit(t, slug, tip, tip)
			if code != 0 || !strings.Contains(out, "drift[0]") {
				t.Fatalf("audit = (%d), want a diagnostic-only exit 0 with no drift:\n%s", code, out)
			}
			var got []string
			for _, row := range table(t, out, "audit_diagnostics") {
				got = append(got, fmt.Sprint(row["subject"], "|", row["reason"]))
			}
			matches := len(got) == len(test.want)
			for i := 0; matches && i < len(got); i++ {
				matches = strings.HasPrefix(got[i], test.want[i])
			}
			if !matches {
				t.Errorf("diagnostics = %q, want the prefixes %q", got, test.want)
			}
			staleness := table(t, out, "staleness")[0]
			if complete := len(test.want) == 0; staleness["complete"] != complete || !complete && staleness["manual_required"] != true {
				t.Errorf("staleness = %v, want complete %t and manual review unless complete", staleness, complete)
			}
			if !strings.HasPrefix(test.row, "SS57") {
				return
			}
			// The pending native entries keep the audit open although the charge, which grades
			// every required current premise, is ready.
			if out, code := preflight.Command(charge); code != 0 {
				t.Errorf("charge = (%d), want a ready charge beside the pending native entries:\n%s", code, out)
			}
		})
	}
}

// TestAuditReportsNoClosureDiagnostic covers CR13. The complete case's ticket also writes the
// bound renderer package, and no verified kit seal exists. The audit stays complete with no
// diagnostic, because no closure reads the binding registry.
func TestAuditReportsNoClosureDiagnostic(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	const bound = "internal/toon/toon.go"
	files := completenessFiles(slug, inventory("declared-complete", keptSeams(slug, true), nil))
	files["specs/"+slug+"/tickets/one.md"] = preflighttest.WritesTicketDoc("One", append(preflighttest.FenceWrites(preflighttest.ConformantFence), bound), "SS28")
	files[bound] = "package toon\n"
	tip := preflighttest.CommitChargeSource(t, root, slug, files)[8]
	out, code := audit(t, slug, tip, tip)
	if code != 0 || strings.Contains(out, "kit-registry") || table(t, out, "staleness")[0]["complete"] != true {
		t.Fatalf("audit over a bound-package write = (%d), want a complete audit with no kit-registry diagnostic:\n%s", code, out)
	}
}
