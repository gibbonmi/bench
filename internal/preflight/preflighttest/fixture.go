// Package preflighttest owns the preflight command fixture: a throwaway Git repository
// with a conformant spec, its ticket, and the canonical charge sources, plus the charge
// arguments and response decoders the command tests share. Only tests import this
// package. It does not import package preflight, because that package's own internal
// tests import it.
package preflighttest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/markdown"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	toonlib "github.com/toon-format/toon-go"
)

// ResponseBudget is the encoded stdout bound every bounded preflight response obeys. Tests
// state it independently of chargeevidence.ResponseLimit, so a raised limit turns them red.
// A lowered limit leaves them green, because a shorter response still holds under this
// budget. The shipped format reference in package chargeevidence owns that direction,
// because the reference states the limit and its projection compares the shipped bytes.
const ResponseBudget = 48000

// ChargeFixtureAssignment is the assignment ID OwnedAssignment registers.
const ChargeFixtureAssignment = "00000000000000000000000000000001"

// ConformantFence is the ownership fence that SeedConformant's ticket `Writes:` line
// derives, in entry order.
var ConformantFence = []string{
	"internal/example/",
	"reviews/example.md",
	".agents/skills/bench-craft-delegate/",
	".agents/commands/bench-implement-spec.md",
}

// RunGit runs one git command in the working directory and returns its trimmed output.
// It keeps the fixture API that the preflight and gate tests call, and gittest owns the
// execution.
func RunGit(t *testing.T, args ...string) string {
	t.Helper()
	return gittest.Output(t, "", args...)
}

// MustWriteFile writes body at path, creating its parent directories.
func MustWriteFile(t *testing.T, path, body string) {
	t.Helper()
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", dir, err)
		}
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

// TicketDoc renders a grammar-conformant ticket citing rows. Its `Writes:` values derive
// from ConformantFence.
func TicketDoc(title string, covers ...string) string {
	return WritesTicketDoc(title, FenceWrites(ConformantFence), covers...)
}

// FenceWrites derives one `Writes:` entry from each fence entry. The builder has no tree
// access, so each entry carries the (new) marker and writes-resolve stays green.
func FenceWrites(fence []string) []string {
	writes := make([]string, len(fence))
	for i, entry := range fence {
		writes[i] = entry + " (new)"
	}
	return writes
}

// WritesTicketDoc renders a conformant ticket that writes exactly writes and cites rows.
func WritesTicketDoc(title string, writes []string, covers ...string) string {
	return commitmenttest.TicketBody(title, writes, covers...)
}

// PlanCommand and PlanProbe are the chunk probe requirement of the seeded completion plan.
const PlanCommand, PlanProbe = "bench test --package ./...", "omit the example behavior"

// PlanFence renders the bench-completion-plan section the checkpoint reader parses, as one
// chunk over tickets. Every seeded spec carries it, so a green fixture states the plan the
// landing later requires.
func PlanFence(tickets ...string) string {
	return PlanSection(reviewrecord.Plan{
		Version: 1,
		Chunks: []reviewrecord.PlannedChunk{{
			ID:           "c1",
			Tickets:      tickets,
			Verification: []reviewrecord.Requirement{{ID: "tests", Command: PlanCommand, Probe: PlanProbe}},
		}},
		FinalVerification: []reviewrecord.Requirement{{ID: "acceptance", Command: "go test ./..."}},
	})
}

// PlanSection renders plan as the bench-completion-plan section that PlanFence frames.
func PlanSection(plan reviewrecord.Plan) string {
	data, err := json.Marshal(plan)
	if err != nil {
		panic(err)
	}
	return "\n## Completion plan\n\n" + markdown.FenceMarker + "bench-completion-plan\n" + string(data) + "\n" + markdown.FenceMarker + "\n"
}

// SpecBody renders a bootstrap-conformant spec for slug. It has staged status, a valid
// opted-in coverage map declaring rows PF1 and PF2, the complete-empty inventory, a chunk
// table that lists both rows, and the completion plan over the tickets/one.md that every seed
// writes. It holds no fence section, because the ticket `Writes:` lines derive the ownership
// fence.
func SpecBody(slug string) string {
	return specBody(slug, true)
}

// SpecBodyNoInventory is SpecBody with no planning inventory section.
func SpecBodyNoInventory(slug string) string {
	return specBody(slug, false)
}

// WithLeftoverFence appends a leftover `## Ownership fences` section that lists each
// entry to spec. No reader grades the section, so a test uses it to prove that the
// section authorizes nothing.
func WithLeftoverFence(spec string, entries ...string) string {
	var b strings.Builder
	b.WriteString(spec + "\n## Ownership fences\n\n")
	for _, entry := range entries {
		b.WriteString("- `" + entry + "`\n")
	}
	return b.String()
}

func specBody(slug string, inventory bool) string {
	var b strings.Builder
	b.WriteString("# " + slug + "\n\nStatus: staged\n\n")
	b.WriteString("## User stories\n1. As a, I want b, so c.\n\n")
	b.WriteString("### Acceptance coverage map\n")
	b.WriteString("| row | story | behavior | seam | why it catches the failure |\n")
	b.WriteString("|---|---|---|---|---|\n")
	for _, row := range specRows {
		b.WriteString("| " + row[0] + " | 1 | " + row[1] + " | cli seam | " + row[2] + " |\n")
	}
	if inventory {
		b.WriteString(inventorytest.Section(inventorytest.CompleteEmpty))
	}
	b.WriteString(ChunkTable(specRowIDs()...))
	b.WriteString(PlanFence("one.md"))
	return b.String()
}

// StartRepo starts a fresh git repository at t.TempDir(), makes it the working directory,
// and configures the shared fixture identity.
func StartRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	t.Chdir(root)
	return root
}

// SeedConformant builds the PF1 tracer fixture. A base commit on main carries a
// bootstrap-conformant spec, a ticket citing both declared rows, and the canonical charge
// sources; then a feature branch adds one authorized change. Every check answers green over
// this tree.
func SeedConformant(t *testing.T) (root, slug string) {
	t.Helper()
	slug = "example"
	root = StartRepo(t)
	MustWriteFile(t, "specs/"+slug+"/spec.md", SpecBody(slug))
	MustWriteFile(t, "specs/"+slug+"/tickets/one.md", TicketDoc("One", "PF1", "PF2"))
	MustWriteFile(t, chargesource.DelegateSkill, "# Delegation skill\n")
	MustWriteFile(t, chargesource.DelegateProcedure, "# Delegation procedure\n\nFocused suite: bench test --package ./internal/preflight\n")
	MustWriteFile(t, chargesource.BuildPhase, "# Build phase\n\nRun root conformance before return.\n")
	commitmenttest.SeedAdmission(t, root, "specs/"+slug+"/spec.md")
	RunGit(t, "add", ".")
	RunGit(t, "commit", "-q", "-m", "c0")
	RunGit(t, "checkout", "-q", "-b", "feature")
	MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	RunGit(t, "add", "internal/"+slug+"/foo.go")
	RunGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}

// CommitChargeSource writes files on main in the SeedConformant repository at root, binds
// the spec to the fixture admission, and commits them. It then restarts the feature branch at
// main with one fenced change and returns the charge arguments over that tip.
func CommitChargeSource(t *testing.T, root, slug string, files map[string]string) []string {
	t.Helper()
	RunGit(t, "checkout", "-q", "main")
	for path, body := range files {
		MustWriteFile(t, path, body)
	}
	commitmenttest.SeedAdmission(t, root, "specs/"+slug+"/spec.md")
	RunGit(t, "add", "-A")
	RunGit(t, "commit", "-q", "-m", "charge source")
	RunGit(t, "checkout", "-q", "-B", "feature")
	MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n\nfunc Foo(ctx int) int { return ctx }\n")
	RunGit(t, "add", "-A")
	RunGit(t, "commit", "-q", "-m", "fenced change")
	return ChargeArgs(t, root, slug)
}

// OwnedAssignment registers one assignment in root's ledger, owning the tree at worktree in
// the given lifecycle state, and returns its ID. The record goes in through
// intent.PutAssignment, so the gatherer reads exactly the shape every worktree command
// writes.
func OwnedAssignment(t *testing.T, root, worktree string, state intent.AssignmentState) string {
	t.Helper()
	const owner = "00000000000000000000000000000002"
	err := intent.PutAssignment(root, intent.Assignment{
		Schema:   intent.AssignmentRecordSchema,
		ID:       ChargeFixtureAssignment,
		OwnerID:  owner,
		Request:  intent.RequestDigest("preflight-assignment-target"),
		Label:    "preflight-assignment-target",
		Start:    RunGit(t, "rev-parse", "HEAD"),
		Branch:   intent.AssignmentBranchRef(owner, ChargeFixtureAssignment),
		Worktree: worktree,
		State:    state,
	})
	if err != nil {
		t.Fatalf("PutAssignment: %v", err)
	}
	return ChargeFixtureAssignment
}

// ActiveAssignment registers one active assignment owning the tree at worktree.
func ActiveAssignment(t *testing.T, root, worktree string) string {
	t.Helper()
	id := OwnedAssignment(t, root, worktree, intent.StateActive)
	if root == worktree {
		admitFixture(t, root)
	}
	return id
}

// ChargeArgs activates root's assignment and returns the build charge arguments over the
// seeded ticket, pinned to main and HEAD. The base is element 6 and the tip is element 8.
func ChargeArgs(t *testing.T, root, slug string) []string {
	t.Helper()
	ActiveAssignment(t, root, root)
	return []string{"build", slug, "--charge", "--ticket", "one.md", "--base", RunGit(t, "rev-parse", "main"), "--source-tip", RunGit(t, "rev-parse", "HEAD")}
}

// ReviewArgs is the review preparation argument vector over the seeded review tree, pinned
// to main and HEAD. The base is element 4 and the tip is element 6.
func ReviewArgs(t *testing.T, slug string) []string {
	t.Helper()
	return []string{"review", slug, "--charge",
		"--base", RunGit(t, "rev-parse", "main"), "--source-tip", RunGit(t, "rev-parse", "HEAD")}
}

// SeedReviewEvidence builds the shared review tree: a changed package, a touched consumer,
// an untouched consumer, a deleted symbol, and every canonical review source. A poisoned
// consumer adds a path the full consumer projection cannot represent, so the collector
// declares truncated output. It returns the repository root, the slug, and the review
// preparation arguments.
func SeedReviewEvidence(t *testing.T, poisonedConsumer bool) (root, slug string, args []string) {
	t.Helper()
	slug = "example"
	root = StartRepo(t)
	for _, file := range reviewCanonical(slug) {
		MustWriteFile(t, file.path, file.body)
	}
	MustWriteFile(t, "target/target.go", "package target\n\nfunc Changed() int { return 0 }\nfunc Gone() {}\n")
	MustWriteFile(t, "outside/user.go",
		"package outside\n\nimport \"example.com/review/target\"\n\nfunc Use() int { return target.Changed() }\n")
	MustWriteFile(t, "edited/user.go",
		"package edited\n\nimport \"example.com/review/target\"\n\nfunc Use() int { return target.Changed() }\n")
	if poisonedConsumer {
		MustWriteFile(t, "outside/a\x1b.go", "package outside\n\nimport \"example.com/review/target\"\n\nfunc Poisoned() int { return target.Changed() }\n")
	}
	RunGit(t, "add", ".")
	RunGit(t, "commit", "-q", "-m", "base")
	RunGit(t, "checkout", "-q", "-b", "feature")
	MustWriteFile(t, "target/target.go", "package target\n\nfunc Changed() int { return 1 }\n")
	MustWriteFile(t, "edited/user.go", "package edited\n\nimport \"example.com/review/target\"\n\n// Use is an edited consumer.\nfunc Use() int { return target.Changed() }\n")
	MustWriteFile(t, "notes/a \"quote\" \\ café*.txt", "review π evidence\n")
	RunGit(t, "add", ".")
	RunGit(t, "commit", "-q", "-m", "review source")
	ActiveAssignment(t, root, root)
	return root, slug, ReviewArgs(t, slug)
}

// LegacyCommitted commits every change and repins the charge arguments to the new tip.
func LegacyCommitted(t *testing.T, root, slug, message string) []string {
	t.Helper()
	RunGit(t, "add", "-A")
	RunGit(t, "commit", "-q", "-m", message)
	args := ChargeArgs(t, root, slug)
	args[8] = RunGit(t, "rev-parse", "HEAD")
	return args
}

// DecodeMap decodes one TOON document into its top-level object.
func DecodeMap(t *testing.T, output string) map[string]any {
	t.Helper()
	decoded, err := toonlib.DecodeString(output)
	if err != nil {
		t.Fatalf("decode TOON: %v\n%s", err, output)
	}
	document, ok := decoded.(map[string]any)
	if !ok {
		t.Fatalf("decoded packet = %T, want object", decoded)
	}
	return document
}

// TableRows returns the rows of one decoded table.
func TableRows(t *testing.T, document map[string]any, table string) []map[string]any {
	t.Helper()
	values, ok := document[table].([]any)
	if !ok {
		t.Fatalf("%s = %T, want table", table, document[table])
	}
	rows := make([]map[string]any, len(values))
	for i, value := range values {
		row, ok := value.(map[string]any)
		if !ok {
			t.Fatalf("%s row %d = %T, want object", table, i, value)
		}
		rows[i] = row
	}
	return rows
}

// StoreDir is root's repository-common evidence store directory.
func StoreDir(t *testing.T, root string) string {
	t.Helper()
	common, err := git.CommonDir(root)
	if err != nil {
		t.Fatalf("common dir: %v", err)
	}
	return filepath.Join(common, chargeevidence.StoreName)
}

// PublishedPacks lists the published pack names in root's evidence store.
func PublishedPacks(t *testing.T, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(StoreDir(t, root))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var packs []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), chargeevidence.PackSuffix) {
			packs = append(packs, entry.Name())
		}
	}
	return packs
}

// PublishedPack reads the one published artifact in root's evidence store through the
// strict reader. A test grades the bytes a consumer would retrieve, not a value the
// producer kept in memory.
func PublishedPack(t *testing.T, root, identity string) *chargeevidence.Pack {
	t.Helper()
	if packs := PublishedPacks(t, root); len(packs) != 1 {
		t.Fatalf("published packs = %d, want one", len(packs))
	}
	return IdentifiedPack(t, root, identity)
}

// StagedTemps lists the temporary pack names in root's evidence store.
func StagedTemps(t *testing.T, root string) []string {
	t.Helper()
	entries, _ := os.ReadDir(StoreDir(t, root))
	var temps []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), chargeevidence.TempPrefix) {
			temps = append(temps, entry.Name())
		}
	}
	return temps
}

// AssertNothingPublished proves that a refused preparation left no artifact and no
// temporary pack behind.
func AssertNothingPublished(t *testing.T, root string) {
	t.Helper()
	if packs, temps := PublishedPacks(t, root), StagedTemps(t, root); len(packs) != 0 || len(temps) != 0 {
		t.Fatalf("refusal left packs %v and temporary packs %v", packs, temps)
	}
}
