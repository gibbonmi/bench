package preflight

import (
	"context"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/tickets"
)

// floorModule writes a Go module at a fresh root with one package in each directory, and
// clears BENCH_KIT, so the root is its own kit and declares no conformance entry test.
func floorModule(t *testing.T, dirs ...string) string {
	t.Helper()
	t.Setenv("BENCH_KIT", "")
	root := t.TempDir()
	preflighttest.MustWriteFile(t, filepath.Join(root, "go.mod"), "module example.com/floor\n\ngo 1.25\n")
	for _, dir := range dirs {
		preflighttest.MustWriteFile(t, filepath.Join(root, dir, "p.go"), "package "+path.Base(dir)+"\n")
	}
	return root
}

// kitModule is floorModule with the conformance entry test declared, so root conformance
// applies to the root.
func kitModule(t *testing.T, dirs ...string) string {
	t.Helper()
	root := floorModule(t, dirs...)
	preflighttest.MustWriteFile(t, filepath.Join(root, "internal", "conformance", "entry_test.go"),
		"package conformance\n\nimport \"testing\"\n\nfunc "+registry.RootConformanceTest+"(t *testing.T) {}\n")
	return root
}

// writer is one parsed ticket that writes entries.
func writer(name string, entries ...string) tickets.Ticket {
	return tickets.Ticket{Name: name, Writes: entries}
}

// chunk is one planned chunk over names whose requirements are commands.
func chunk(id string, names []string, commands ...string) reviewrecord.PlannedChunk {
	planned := reviewrecord.PlannedChunk{ID: id, Tickets: names}
	for i, command := range commands {
		planned.Verification = append(planned.Verification, reviewrecord.Requirement{ID: id + "-" + strconv.Itoa(i), Command: command})
	}
	return planned
}

// onePlan is a version 1 plan of one chunk over every ticket of planned.
func onePlan(planned []tickets.Ticket, commands ...string) reviewrecord.Plan {
	var names []string
	for _, ticket := range planned {
		names = append(names, ticket.Name)
	}
	return reviewrecord.Plan{Version: 1, Chunks: []reviewrecord.PlannedChunk{chunk("c1", names, commands...)}}
}

// floorRow gathers the check floor of planned against plan over root, and returns the
// check-floor row that Decide gives for those facts.
func floorRow(t *testing.T, root string, plan reviewrecord.Plan, planned ...tickets.Ticket) CheckResult {
	t.Helper()
	f := baseFacts()
	f.Tickets = planned
	f.CompletionPlan = plan
	f.CheckFloor = gatherCheckFloor(context.Background(), root, planned, plan)
	row, ok := checkRow(Decide(f), checkFloorName)
	if !ok {
		t.Fatalf("Decide gave no %s row", checkFloorName)
	}
	return row
}

// requireFloorRed fails unless the row is red with exactly the omitted items. The detail
// comes from the production builders, so TestCheckFloorDetailNamesTheCommand alone states
// its spelling.
func requireFloorRed(t *testing.T, row CheckResult, omitted ...string) {
	t.Helper()
	requireFloorDetail(t, row, floorRedDetail(omitted))
}

// requireFloorDetail fails unless the row is red with exactly the detail want.
func requireFloorDetail(t *testing.T, row CheckResult, want string) {
	t.Helper()
	if row.Verdict != verdictRed || row.Detail != want {
		t.Fatalf("check-floor row = %+v, want the red detailed %q", row, want)
	}
}

// requireFloorGreen fails unless the row is green.
func requireFloorGreen(t *testing.T, row CheckResult) {
	t.Helper()
	if row.Verdict != verdictGreen || row.Detail != "" {
		t.Fatalf("check-floor row = %+v, want green", row)
	}
}

// floorOf gathers the elements of one ticket over root.
func floorOf(t *testing.T, root string, ticket tickets.Ticket) []string {
	t.Helper()
	facts := gatherCheckFloor(context.Background(), root, []tickets.Ticket{ticket}, onePlan([]tickets.Ticket{ticket}))
	if facts.LoaderError != "" {
		t.Fatalf("loader error: %s", facts.LoaderError)
	}
	return facts.Elements[ticket.Name]
}

func requireFloor(t *testing.T, got []string, want ...string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("floor = %q, want %q", got, want)
	}
}

func TestCheckFloorAddsRootConformanceOnTheKit(t *testing.T) {
	root := kitModule(t)
	for _, kit := range []string{"", root} {
		t.Setenv("BENCH_KIT", kit)
		requireFloor(t, floorOf(t, root, writer("one.md", "docs/guide.md")), testreport.RootConformanceCheckName)
	}
}

func TestCheckFloorOmitsRootConformanceOffTheKit(t *testing.T) {
	root := floorModule(t, "internal/x")
	for _, kit := range []string{"", root, kitModule(t)} {
		t.Setenv("BENCH_KIT", kit)
		requireFloor(t, floorOf(t, root, writer("one.md", "internal/x/p.go")), "./internal/x")
	}
}

func TestCheckFloorAddsCmdBenchForACommandPackage(t *testing.T) {
	root := floorModule(t, "internal/worktree")
	requireFloor(t, floorOf(t, root, writer("one.md", "internal/worktree/p.go")), "./internal/worktree", commandPackageElement)
	// A new file in a command package that holds no package yet still owes the command tests.
	requireFloor(t, floorOf(t, root, writer("one.md", "internal/anchors/next.go (new)")), commandPackageElement)
}

func TestCheckFloorAddsCmdBenchForADirectoryAboveACommandPackage(t *testing.T) {
	root := floorModule(t, "internal/worktree")
	requireFloor(t, floorOf(t, root, writer("one.md", "internal/")), "./internal/worktree", commandPackageElement)
}

func TestCheckFloorAddsCmdBenchForAnEmbedChange(t *testing.T) {
	root := floorModule(t)
	preflighttest.MustWriteFile(t, filepath.Join(root, "internal", "assets", "assets.go"),
		"package assets\n\nimport _ \"embed\"\n\n//go:embed a.txt\nvar A string\n")
	preflighttest.MustWriteFile(t, filepath.Join(root, "internal", "assets", "a.txt"), "asset\n")
	requireFloor(t, floorOf(t, root, writer("one.md", "internal/assets/assets.go")), "./internal/assets", commandPackageElement)
}

func TestCheckFloorOmitsCmdBenchOtherwise(t *testing.T) {
	root := floorModule(t, "internal/x")
	requireFloor(t, floorOf(t, root, writer("one.md", "internal/x/p.go", "docs/guide.md")), "./internal/x")
}

func TestCheckFloorRedsAnOmittedPackage(t *testing.T) {
	t.Run("package", func(t *testing.T) {
		root := floorModule(t, "internal/x", "internal/other")
		one := writer("one.md", "internal/x/p.go")
		requireFloorRed(t, floorRow(t, root, onePlan([]tickets.Ticket{one}, "bench test --package ./internal/other"), one),
			omittedElement("one.md", "./internal/x"))
	})
	t.Run("command package alone", func(t *testing.T) {
		root := floorModule(t, "internal/worktree")
		one := writer("one.md", "internal/worktree/p.go")
		requireFloorRed(t, floorRow(t, root, onePlan([]tickets.Ticket{one}, "bench test --package ./internal/worktree"), one),
			omittedElement("one.md", commandPackageElement))
	})
	t.Run("second ticket", func(t *testing.T) {
		root := floorModule(t, "internal/x", "internal/y")
		one, two := writer("one.md", "internal/x/p.go"), writer("two.md", "internal/y/p.go")
		requireFloorRed(t, floorRow(t, root, onePlan([]tickets.Ticket{one, two}, "bench test --package ./internal/x"), one, two),
			omittedElement("two.md", "./internal/y"))
	})
}

func TestCheckFloorRedsAnUnplannedTicket(t *testing.T) {
	root := floorModule(t, "internal/x", "internal/y")
	one, two := writer("one.md", "internal/x/p.go"), writer("two.md", "internal/y/p.go")
	plan := onePlan([]tickets.Ticket{one}, "bench test --package ./internal/x")
	requireFloorRed(t, floorRow(t, root, plan, one, two), omittedElement("two.md", "./internal/y"))
	// A ticket whose floor is empty stays green, although no chunk names it.
	requireFloorGreen(t, floorRow(t, root, plan, one, writer("three.md", "docs/guide.md")))
}

func TestCheckFloorGreenWhenSatisfied(t *testing.T) {
	root := kitModule(t, "internal/x", "internal/worktree")
	one, two := writer("one.md", "internal/x/p.go"), writer("two.md", "internal/worktree/p.go")
	plan := onePlan([]tickets.Ticket{one, two}, "bench test --check root-conformance",
		"bench test --package ./internal/x", "bench test --package ./internal/worktree", "bench test --package ./cmd/bench")
	requireFloorGreen(t, floorRow(t, root, plan, one, two))
}

// TestCheckFloorDetailNamesTheCommand states the whole detail as a literal, apart from the
// code that composes it, so a renamed check, flag, or separator turns it red.
func TestCheckFloorDetailNamesTheCommand(t *testing.T) {
	root := kitModule(t, "internal/x", "internal/other")
	one := writer("one.md", "internal/x/p.go")
	row := floorRow(t, root, onePlan([]tickets.Ticket{one}, "bench test --package ./internal/other"), one)
	requireFloorDetail(t, row, "completion plan omits a floor check: "+
		"one.md: root-conformance (add bench test --check root-conformance); "+
		"one.md: ./internal/x (add bench test --package ./internal/x)")
}

// satisfaction grades one element of ticket one.md, which writes entry, against command.
func satisfaction(t *testing.T, root, entry, command string) CheckResult {
	t.Helper()
	one := writer("one.md", entry)
	return floorRow(t, root, onePlan([]tickets.Ticket{one}, command), one)
}

func TestCheckFloorAcceptsTheExactPackage(t *testing.T) {
	root := floorModule(t, "internal/x")
	requireFloorGreen(t, satisfaction(t, root, "internal/x/p.go", "bench test --package ./internal/x"))
}

func TestCheckFloorAcceptsAParentPattern(t *testing.T) {
	root := floorModule(t, "internal/x")
	for _, command := range []string{"bench test --package ./internal/...", "bench test --package ./internal/x/..."} {
		requireFloorGreen(t, satisfaction(t, root, "internal/x/p.go", command))
	}
}

func TestCheckFloorAcceptsTheWholeSuite(t *testing.T) {
	root := floorModule(t, "internal/x")
	preflighttest.MustWriteFile(t, filepath.Join(root, "p.go"), "package floor\n")
	one := writer("one.md", "p.go", "internal/x/p.go", "go.mod")
	requireFloor(t, floorOf(t, root, one), "./.", "./...", "./internal/x")
	requireFloorGreen(t, floorRow(t, root, onePlan([]tickets.Ticket{one}, "bench test --package ./..."), one))
}

func TestCheckFloorKeepsTheSegmentBoundary(t *testing.T) {
	root := floorModule(t, "internal/git", "internal/git2")
	requireFloorRed(t, satisfaction(t, root, "internal/git2/p.go", "bench test --package ./internal/git/..."),
		omittedElement("one.md", "./internal/git2"))
}

func TestCheckFloorKeepsTheReverseBoundary(t *testing.T) {
	root := floorModule(t, "internal/git", "internal/git2")
	requireFloorRed(t, satisfaction(t, root, "internal/git/p.go", "bench test --package ./internal/git2"),
		omittedElement("one.md", "./internal/git"))
}

func TestCheckFloorRejectsAFilteredRun(t *testing.T) {
	root := floorModule(t, "internal/x")
	requireFloorRed(t, satisfaction(t, root, "internal/x/p.go", "bench test --package ./internal/x --run '^TestA$|^TestB$'"),
		omittedElement("one.md", "./internal/x"))
}

func TestCheckFloorRejectsAChangedRun(t *testing.T) {
	root := floorModule(t, "internal/worktree")
	requireFloorRed(t, satisfaction(t, root, "internal/worktree/p.go", "bench test --changed"),
		omittedElement("one.md", "./internal/worktree"), omittedElement("one.md", commandPackageElement))
}

func TestCheckFloorRejectsACompoundCommand(t *testing.T) {
	root := floorModule(t, "internal/x")
	for _, command := range []string{"true && bench test --package ./internal/x", "bench test --package ./internal/x && true"} {
		requireFloorRed(t, satisfaction(t, root, "internal/x/p.go", command), omittedElement("one.md", "./internal/x"))
	}
}

func TestCheckFloorRejectsAnotherExecutable(t *testing.T) {
	root := floorModule(t, "internal/x")
	for _, command := range []string{"go test ./internal/x", "BENCH_KIT=x bench test --package ./internal/x"} {
		requireFloorRed(t, satisfaction(t, root, "internal/x/p.go", command), omittedElement("one.md", "./internal/x"))
	}
}

func TestCheckFloorRejectsARefusedCommand(t *testing.T) {
	root := floorModule(t, "internal/x")
	requireFloorRed(t, satisfaction(t, root, "internal/x/p.go", "bench test --package ./internal/x --bogus"),
		omittedElement("one.md", "./internal/x"))
}

func TestCheckFloorRejectsANonRunningForm(t *testing.T) {
	t.Run("check inventory", func(t *testing.T) {
		requireFloorRed(t, satisfaction(t, floorModule(t, "internal/x"), "internal/x/p.go", "bench test --checks"),
			omittedElement("one.md", "./internal/x"))
	})
	t.Run("fixture listing", func(t *testing.T) {
		requireFloorRed(t, satisfaction(t, kitModule(t), "docs/guide.md", "bench test --check root-conformance --fixtures"),
			conformanceRed())
	})
}

// conformanceRed is the red item of the unsatisfied root-conformance element.
func conformanceRed() string {
	return omittedElement("one.md", testreport.RootConformanceCheckName)
}

func TestCheckFloorAcceptsTheRootConformanceCheck(t *testing.T) {
	root := kitModule(t)
	requireFloorGreen(t, satisfaction(t, root, "docs/guide.md", "bench test --check root-conformance"))
}

func TestCheckFloorRejectsAScopedConformanceCheck(t *testing.T) {
	root := kitModule(t)
	requireFloorRed(t, satisfaction(t, root, "docs/guide.md", "bench test --check docs-currency-workflow"), conformanceRed())
}

func TestCheckFloorRejectsTheConformancePackage(t *testing.T) {
	root := kitModule(t)
	requireFloorRed(t, satisfaction(t, root, "docs/guide.md", "bench test --package ./internal/conformance"), conformanceRed())
}

func TestCheckFloorReadsOnlyTheOwnChunk(t *testing.T) {
	root := floorModule(t, "internal/x", "internal/y")
	one, two := writer("one.md", "internal/x/p.go"), writer("two.md", "internal/y/p.go")
	plan := reviewrecord.Plan{Version: 1, Chunks: []reviewrecord.PlannedChunk{
		chunk("c1", []string{"one.md"}, "bench test --package ./internal/x", "bench test --package ./internal/y"),
		chunk("c2", []string{"two.md"}, "bench test --package ./internal/x"),
	}}
	requireFloorRed(t, floorRow(t, root, plan, one, two), omittedElement("two.md", "./internal/y"))
}

func TestCheckFloorReadsOnlyTheOwnTicket(t *testing.T) {
	root := floorModule(t, "internal/x", "internal/y")
	one, two := writer("one.md", "internal/x/p.go"), writer("two.md", "internal/y/p.go")
	planned := chunk("c1", []string{"one.md", "two.md"},
		"bench test --package ./internal/x", "bench test --package ./internal/y", "bench test --package ./internal/x")
	for i, owner := range []string{"one.md", "one.md", "two.md"} {
		planned.Verification[i].Ticket = owner
	}
	plan := reviewrecord.Plan{Version: 2, Execution: &reviewrecord.Execution{Mode: "delegate"}, Chunks: []reviewrecord.PlannedChunk{planned}}
	if !plan.Delegated() {
		t.Fatal("fixture invalid: the plan is not delegated")
	}
	requireFloorRed(t, floorRow(t, root, plan, one, two), omittedElement("two.md", "./internal/y"))
}

// omittingFacts is a fact set whose floor element no requirement satisfies.
func omittingFacts() Facts {
	f := baseFacts()
	f.CompletionPlan = onePlan(f.Tickets, "go test ./...")
	f.CheckFloor = CheckFloorFacts{Elements: map[string][]string{"one.md": {"./internal/example"}}}
	return f
}

func TestCheckFloorDefersToTheCompletionPlan(t *testing.T) {
	f := omittingFacts()
	requireFloorRed(t, checkFloorCheck(f), omittedElement("one.md", "./internal/example"))
	f.CompletionPlanError = "missing or invalid completion plan: no record"
	if row, _ := checkRow(Decide(f), checkFloorName); row.Verdict != verdictNA || row.Detail != "" {
		t.Fatalf("check-floor under an invalid plan = %+v, want not-applicable", row)
	}
}

func TestCheckFloorNotApplicableWithoutTickets(t *testing.T) {
	f := omittingFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = false
	if row, _ := checkRow(Decide(f), checkFloorName); row.Verdict != verdictNA || row.Detail != "" {
		t.Fatalf("check-floor in a build without tickets = %+v, want not-applicable", row)
	}
}

func TestCheckFloorRedsALoaderFailure(t *testing.T) {
	root := floorModule(t, "internal/x")
	preflighttest.MustWriteFile(t, filepath.Join(root, "go.mod"), "module example.com/floor\n\nnot a directive\n")
	one := writer("one.md", "internal/x/p.go")
	_, err := testreport.WritesFloor(context.Background(), root, one.Writes)
	if err == nil {
		t.Fatal("fixture invalid: the package loader read a malformed go.mod")
	}
	row := floorRow(t, root, onePlan([]tickets.Ticket{one}, "bench test --package ./internal/x"), one)
	if want := "package index not readable: " + err.Error(); row.Verdict != verdictRed || row.Detail != want {
		t.Fatalf("check-floor row = %+v, want the red detailed %q", row, want)
	}
}
