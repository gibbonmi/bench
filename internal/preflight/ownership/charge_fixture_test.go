package ownership_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/tickets"
)

func decisionItem(id, row, predicate string, line int, needle string) string {
	loc := fmt.Sprintf(`{"path":%q,"line":%d,"needle":%q,"source":"tip"}`, specPath, line, needle)
	return fmt.Sprintf(`{"id":%q,%s,"source":%s,"predicate":%q,"row":%q}`, id, itemLocation, loc, predicate, row)
}

func seamItem(id, path, fn, state string, rows ...string) string {
	return fmt.Sprintf(`{"id":%q,%s,"path":%q,"package":"example.com/fixture/internal/example","function":%q,"state":%q,"rows":%s}`,
		id, itemLocation, path, fn, state, encodeJSON(rows))
}

func currentManualJSON(id, reason string) string {
	return fmt.Sprintf(`{"id":%q,%s,"reason":%q,"applicability":"current-premise"}`, id, itemLocation, reason)
}

// seededRows is each coverage row of the seeded spec: its line and its exact behavior text.
var seededRows = map[string]struct {
	line     int
	behavior string
}{"PF1": {11, "does x"}, "PF2": {12, "does y"}}

// plannedSeam is one planned seam of both seeded rows, and rowDecisions maps each seeded
// row's exact behavior bytes.
var (
	plannedSeam  = []string{seamItem("seam-1", "internal/example/foo_test.go", "TestFoo", "planned", "PF1", "PF2")}
	rowDecisions = []string{rowDecision("dec-1", "PF1"), rowDecision("dec-2", "PF2")}
)

func rowDecision(id, row string) string {
	r := seededRows[row]
	return decisionItem(id, row, r.behavior, r.line, r.behavior)
}

// stalePremise is a current line-count premise that the seeded foo.go defeats.
var stalePremise = strings.Replace(linesPremise("p1"), `"count":1`, `"count":9`, 1)

// citedSpec is the seeded spec whose PF1 seam cell cites the given functions of
// internal/example/foo_test.go.
func citedSpec(functions string) string {
	return strings.Replace(preflighttest.SpecBody(slug), "| PF1 | 1 | does x | cli seam |", "| PF1 | 1 | does x | `internal/example/foo_test.go` ("+functions+") |", 1)
}

// nativeClaim is one native-future manual claim. located is the row whose spec location the
// claim carries.
type nativeClaim struct {
	id, row, located, chunk, ticket, verification string
	owner                                         ownership.Location
}

// validNative is the native-future claim that the seeded charge repository proves: its own
// row location, the one ticket and chunk that own PF1, the chunk verification, and the
// build phase as the pinned native procedure.
func validNative() nativeClaim {
	return nativeClaim{id: "manual-native", row: "PF1", located: "PF1", chunk: "c1", ticket: "one.md", verification: "tests",
		owner: ownership.Location{Path: chargesource.BuildPhase, Line: 3, Needle: "Run root conformance", Source: "tip"}}
}

// editedNative is validNative changed by edit.
func editedNative(edit func(*nativeClaim)) string {
	claim := validNative()
	edit(&claim)
	return claim.json()
}

func (c nativeClaim) json() string {
	row := seededRows[c.located]
	return encodeJSON(map[string]any{
		"id":            c.id,
		"location":      ownership.Location{Path: specPath, Line: row.line, Needle: fmt.Sprintf("| %s | 1 | %s |", c.located, row.behavior), Source: "spec"},
		"reason":        "future native execution",
		"applicability": "native-future",
		"row":           c.row, "chunk": c.chunk, "ticket": c.ticket, "verification": c.verification,
		"phase_owner": c.owner, "observation": "unobserved",
	})
}

// manualRefusal is one charge repository whose manual claim keeps the charge readiness
// unmet: the repository, the detail its blocker names, and whether the repository also holds
// a valid pending native claim.
type manualRefusal struct {
	name    string
	repo    chargeRepo
	want    string
	pending bool
}

// manualRefusals is each current-premise claim, unproved current premise beside pending
// native work, and invalid native ownership that keeps the charge readiness unmet. A padded
// row cell keeps its native claim pending, because the coverage parser trims the cell and
// reports its line.
func manualRefusals() []manualRefusal {
	inventory := func(premises, seams []string, claim string) string {
		return manifestJSON("declared-complete", premises, seams, nil, []string{claim})
	}
	native := func(edit func(*nativeClaim)) chargeRepo {
		return chargeRepo{facts: inventory(nil, nil, editedNative(edit))}
	}
	valid := validNative().json()
	pad := func(text string) string {
		return strings.Replace(text, "| PF1 | 1 | does x |", "| PF1  | 1 | does x |", 1)
	}
	twoOwners := chargeRepo{facts: inventory(nil, nil, valid), plan: preflighttest.PlanFence("one.md", "two.md"), extra: secondTicket("PF1")}
	foreign, delegated := native(func(c *nativeClaim) { c.chunk = "c2" }), native(func(c *nativeClaim) { c.verification = "tests-two" })
	foreign.plan, foreign.extra = splitPlan(), secondTicket("PF2")
	delegated.plan, delegated.extra = delegatedPlan(), secondTicket("PF2")
	return []manualRefusal{
		{name: "OC89 OC94 current-premise claim", repo: chargeRepo{facts: inventory(nil, nil, currentManualJSON("manual-current", "unproved current assumption"))}, want: "manual-current"},
		{name: "OC94 unproved current premise beside pending native work", repo: chargeRepo{facts: inventory([]string{stalePremise}, nil, valid)}, want: "p1: ", pending: true},
		{name: "OC94 OC95 padded row cell keeps its native claim pending", repo: chargeRepo{spec: pad(preflighttest.SpecBody(slug)), facts: inventory([]string{stalePremise}, nil, pad(valid))}, want: "p1: ", pending: true},
		{name: "OC95 supported test seam on the native row", repo: chargeRepo{spec: citedSpec("`TestFoo`"), facts: inventory(nil, []string{seamItem("seam-1", "internal/example/foo_test.go", "TestFoo", "existing", "PF1")}, valid)}, want: "declares a supported test seam"},
		{name: "OC95 another row's location", repo: native(func(c *nativeClaim) { c.located = "PF2" }), want: "is not the location of row PF1"},
		{name: "OC95 unowned ticket", repo: native(func(c *nativeClaim) { c.ticket = "two.md" }), want: "owned by ticket two.md alone"},
		{name: "OC95 two owning tickets", repo: twoOwners, want: "owned by ticket one.md alone"},
		{name: "OC95 foreign chunk", repo: foreign, want: "chunk c2 does not own ticket one.md"},
		{name: "OC95 missing chunk verification", repo: native(func(c *nativeClaim) { c.verification = "missing-req" }), want: "declares no verification missing-req"},
		{name: "OC95 another ticket's verification", repo: delegated, want: "declares no verification tests-two"},
		{name: "OC95 missing phase owner", repo: native(func(c *nativeClaim) { c.owner.Path = ".agents/commands/absent.md" }), want: "absent.md"},
		{name: "OC95 another procedure as phase owner", repo: native(func(c *nativeClaim) {
			c.owner = ownership.Location{Path: chargesource.DelegateSkill, Line: 1, Needle: "# Delegation skill", Source: "tip"}
		}), want: "is not the native implementation procedure"},
	}
}

// chargeRepo is one conformant build repository. spec is the spec body before its inventory
// and plan, facts is its planning inventory, plan is its completion plan section, and extra
// holds further committed files. An empty spec or plan takes the seeded one.
type chargeRepo struct {
	spec, facts, plan string
	extra             map[string]string
}

// seedChargeRepo seeds a conformant build repository whose main branch carries specBody,
// facts, and extra files, and whose feature branch edits internal/example/foo.go.
func seedChargeRepo(t *testing.T, specBody, facts string, extra map[string]string) (string, []string) {
	t.Helper()
	return chargeRepo{spec: specBody, facts: facts, extra: extra}.seed(t)
}

// seed commits the spec files and extra as the charge source of the shared fixture, which
// owns every commit step.
func (c chargeRepo) seed(t *testing.T) (string, []string) {
	t.Helper()
	files := specFiles(t, c.spec, c.facts, c.plan)
	for path, body := range c.extra {
		files[path] = body
	}
	root, _ := preflighttest.SeedConformant(t)
	return root, preflighttest.CommitChargeSource(t, root, slug, files)
}

// writeSpec writes the files of specFiles.
func writeSpec(t *testing.T, specBody, facts, plan string) {
	t.Helper()
	for path, body := range specFiles(t, specBody, facts, plan) {
		preflighttest.MustWriteFile(t, path, body)
	}
}

// specFiles is the spec of specBody with facts as its inventory and plan as its completion
// plan, and the test file whose functions the seam cases cite.
func specFiles(t *testing.T, specBody, facts, plan string) map[string]string {
	t.Helper()
	if specBody == "" {
		specBody = preflighttest.SpecBody(slug)
	}
	if plan == "" {
		plan = preflighttest.PlanFence("one.md")
	}
	base := ownershiptest.Cut(t, ownershiptest.Cut(t, specBody, preflighttest.PlanFence("one.md")), fenced(completeEmpty))
	return map[string]string{
		specPath:                       base + fenced(facts) + plan,
		"internal/example/foo_test.go": "package example\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\nfunc TestBar(t *testing.T) {}\n",
	}
}

func writePredicateSpec(t *testing.T, root, specBody, facts string) {
	t.Helper()
	writeSpec(t, specBody, facts, "")
	commitmenttest.SeedAdmission(t, root, specPath)
	commit(t, "predicate spec")
}

// finalTests is the final verification of every plan these cases seed.
var finalTests = []reviewrecord.Requirement{{ID: "acceptance", Command: "go test ./..."}}

// splitPlan places one.md in chunk c1 and two.md in chunk c2.
func splitPlan() string {
	tests := []reviewrecord.Requirement{{ID: "tests", Command: preflighttest.PlanCommand}}
	return preflighttest.PlanSection(reviewrecord.Plan{Version: 1, FinalVerification: finalTests, Chunks: []reviewrecord.PlannedChunk{
		{ID: "c1", Tickets: []string{"one.md"}, Verification: tests}, {ID: "c2", Tickets: []string{"two.md"}, Verification: tests}}})
}

// delegatedPlan is a version 2 plan whose one chunk holds one.md and two.md, each with its
// own verification.
func delegatedPlan() string {
	return preflighttest.PlanSection(reviewrecord.Plan{Version: 2, FinalVerification: finalTests,
		Execution: &reviewrecord.Execution{Mode: "delegate", RunID: "run-1", OrchestratorSession: "orchestrator", AuthorLimit: 1,
			Assignments: map[string][]reviewrecord.Assignment{"one.md": {}, "two.md": {}}},
		Chunks: []reviewrecord.PlannedChunk{{ID: "c1", Tickets: []string{"one.md", "two.md"}, Verification: []reviewrecord.Requirement{
			{ID: "tests-one", Command: preflighttest.PlanCommand, Ticket: "one.md"}, {ID: "tests-two", Command: preflighttest.PlanCommand, Ticket: "two.md"}}}}})
}

// secondTicket is the committed two.md ticket that covers rows.
func secondTicket(rows ...string) map[string]string {
	return map[string]string{"specs/" + slug + "/tickets/two.md": preflighttest.WritesTicketDoc("Two", []string{"internal/example/ (new)"}, rows...)}
}

func assertChargeRefused(t *testing.T, root string, args []string, want ...string) {
	t.Helper()
	out, code := preflight.Command(args)
	if code == 0 {
		t.Fatalf("charge succeeded, want refusal:\n%s", out)
	}
	for _, needle := range want {
		if !strings.Contains(out, needle) {
			t.Fatalf("charge refusal = %q, want substring %q", out, needle)
		}
	}
	preflighttest.AssertNothingPublished(t, root)
}

// assertChargePublished runs the charge, requires its one published pack, and returns the
// evidence handle it prints.
func assertChargePublished(t *testing.T, root string, args []string) string {
	t.Helper()
	out, code := preflight.Command(args)
	if code != 0 {
		t.Fatalf("charge = (%d), want 0:\n%s", code, out)
	}
	if packs := preflighttest.PublishedPacks(t, root); len(packs) != 1 {
		t.Fatalf("published packs = %v, want 1", packs)
	}
	return preflighttest.TableRows(t, preflighttest.DecodeMap(t, out), "prepared")[0]["evidence"].(string)
}

// checkCurrent runs the actual current-evidence route over one prepared handle.
func checkCurrent(identity string) (string, int) {
	return preflight.Command([]string{"evidence", identity, "--check-current"})
}

// collectChargeSnapshot collects the premise snapshot of the seeded spec against every
// seeded ticket, between main and HEAD.
func collectChargeSnapshot(t *testing.T, root string) ownership.Snapshot {
	t.Helper()
	specBytes, err := os.ReadFile(filepath.Join(root, specPath))
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "specs", slug, "tickets")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	var parsed []tickets.Ticket
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		ticket, _ := tickets.ParseTicket(name, data, names, "PF")
		parsed = append(parsed, ticket)
	}
	source := ownership.Source{
		SpecPath: specPath,
		Base:     preflighttest.RunGit(t, "rev-parse", "main"),
		Tip:      preflighttest.RunGit(t, "rev-parse", "HEAD"),
	}
	return ownership.Collect(root, source, ownership.ParseManifest(specBytes), parsed)
}

// assertNotReady fails the test unless the seeded snapshot's charge readiness is unmet with a
// blocker that names want.
func assertNotReady(t *testing.T, root, want string) {
	t.Helper()
	ready := collectChargeSnapshot(t, root).ChargeReadiness
	if ready.InventoryComplete && ready.CurrentComplete || !strings.Contains(ready.Blocker, want) {
		t.Fatalf("charge readiness = %+v, want unmet with a blocker naming %q", ready, want)
	}
}
