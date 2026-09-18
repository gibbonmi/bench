// Package preflighttest owns the preflight command fixture: a throwaway Git repository
// with a conformant spec, its ticket, and the canonical charge sources, plus the charge
// arguments and response decoders the command tests share. Only tests import this
// package. It does not import package preflight, because that package's own internal
// tests import it.
package preflighttest

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
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

// ConformantFence is the ownership fence SeedConformant's spec declares, in document order.
var ConformantFence = []string{
	"internal/example/",
	"reviews/example.md",
	".agents/skills/bench-craft-delegate/",
	".agents/commands/bench-implement-spec.md",
}

// reviewFenceExtra are the fence entries SeedReviewEvidence's spec adds to ConformantFence,
// each beside the annotation its line carries. The seeded spec lines and ReviewFence both
// derive from this table, so the fixture states its own fence once.
var reviewFenceExtra = []struct{ path, annotation string }{
	{"target/", "review fixture"},
	{"edited/", "review fixture"},
	{"outside/", "review fixture"},
	{"notes/", "review fixture"},
	{".agents/skills/bench-craft-review/", "review instructions"},
	{".agents/commands/bench-review-implementation.md", "review phase"},
}

// ReviewFence is the complete ownership fence SeedReviewEvidence's spec declares, in
// document order.
func ReviewFence() []string {
	fence := append([]string{}, ConformantFence...)
	for _, entry := range reviewFenceExtra {
		fence = append(fence, entry.path)
	}
	return fence
}

// reviewFenceLines renders the extra fence entries as the spec's own document lines.
func reviewFenceLines() []string {
	lines := make([]string, len(reviewFenceExtra))
	for i, entry := range reviewFenceExtra {
		lines[i] = "- `" + entry.path + "` (" + entry.annotation + ")"
	}
	return lines
}

// RunGit runs one git command in the working directory and returns its trimmed output.
func RunGit(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
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
// from ConformantFence, so the fence and the ticket union stay equal.
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
	return "# " + title + "\n\n" +
		"Blocked by: none\n" +
		"Writes: " + strings.Join(writes, ", ") + "\n" +
		"Covers: " + strings.Join(covers, ", ") + "\n\n" +
		"## What to build\n\nBuild it.\n\n" +
		"## Acceptance\n\n- [ ] It is built.\n"
}

// PlanFence renders the bench-completion-plan section the checkpoint reader parses, as one
// chunk over tickets. Every seeded spec carries it, so a green fixture states the plan the
// landing later requires. The section heading closes the ownership-fence section above it,
// so no payload byte reads as a fence token.
func PlanFence(tickets ...string) string {
	plan := reviewrecord.Plan{
		Version: 1,
		Chunks: []reviewrecord.PlannedChunk{{
			ID:           "c1",
			Tickets:      tickets,
			Verification: []reviewrecord.Requirement{{ID: "tests", Command: "go test ./..."}},
		}},
		FinalVerification: []reviewrecord.Requirement{{ID: "acceptance", Command: "go test ./..."}},
	}
	data, err := json.Marshal(plan)
	if err != nil {
		panic(err)
	}
	return "\n## Completion plan\n\n```bench-completion-plan\n" + string(data) + "\n```\n"
}

// SpecBody renders a bootstrap-conformant spec for slug. It has staged status, a valid
// opted-in coverage map declaring rows PF1 and PF2, and the ConformantFence entries plus
// extraFenceLines. It closes with the completion plan the landing's checkpoint reads, over
// the tickets/one.md every conformant seed writes.
func SpecBody(slug string, extraFenceLines ...string) string {
	var b strings.Builder
	b.WriteString("# " + slug + "\n\nStatus: staged\n\n")
	b.WriteString("## User stories\n1. As a, I want b, so c.\n\n")
	b.WriteString("### Acceptance coverage map\n")
	b.WriteString("| row | story | behavior | seam | why it catches the failure |\n")
	b.WriteString("|---|---|---|---|---|\n")
	b.WriteString("| PF1 | 1 | does x | cli seam | catches z |\n")
	b.WriteString("| PF2 | 1 | does y | cli seam | catches w |\n")
	b.WriteString("\n## Ownership fences\n\n")
	b.WriteString("- `internal/" + slug + "/` (implementation)\n")
	b.WriteString("- `reviews/" + slug + ".md` (review pickup)\n")
	b.WriteString("- `.agents/skills/bench-craft-delegate/` (delegation instructions)\n")
	b.WriteString("- `.agents/commands/bench-implement-spec.md` (build instructions)\n")
	for _, line := range extraFenceLines {
		b.WriteString(line + "\n")
	}
	b.WriteString(PlanFence("one.md"))
	return b.String()
}

// StartRepo starts a fresh git repository at t.TempDir(), makes it the working directory,
// and configures an author identity.
func StartRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Chdir(root)
	RunGit(t, "init", "-q", "-b", "main")
	RunGit(t, "config", "user.email", "t@example.com")
	RunGit(t, "config", "user.name", "t")
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
	RunGit(t, "add", ".")
	RunGit(t, "commit", "-q", "-m", "c0")
	RunGit(t, "checkout", "-q", "-b", "feature")
	MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	RunGit(t, "add", "internal/"+slug+"/foo.go")
	RunGit(t, "commit", "-q", "-m", "c1")
	return root, slug
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
	return OwnedAssignment(t, root, worktree, intent.StateActive)
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
	MustWriteFile(t, "go.mod", "module example.com/review\n\ngo 1.25\n")
	MustWriteFile(t, "specs/"+slug+"/spec.md", SpecBody(slug, reviewFenceLines()...))
	MustWriteFile(t, "specs/"+slug+"/tickets/one.md", WritesTicketDoc("One", FenceWrites(ReviewFence()), "PF1", "PF2"))
	MustWriteFile(t, chargesource.DelegateSkill, "# Delegation skill\n")
	MustWriteFile(t, chargesource.DelegateProcedure,
		"# Delegation procedure\n\nFocused suite: bench test --package ./internal/preflight\n")
	MustWriteFile(t, chargesource.BuildPhase, "# Build phase\n")
	MustWriteFile(t, chargesource.ReviewSkill,
		"# Review skill\n\n## Standards\n\nRules.\n\n## Spec\n\nRequirements.\n\n## Coverage\n\nEdges.\n")
	MustWriteFile(t, chargesource.ReviewPhase, "# Review phase\n\nUse the three canonical axes.\n")
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
	packs := PublishedPacks(t, root)
	if len(packs) != 1 {
		t.Fatalf("published packs = %d, want one", len(packs))
	}
	data, err := os.ReadFile(filepath.Join(StoreDir(t, root), packs[0]))
	if err != nil {
		t.Fatal(err)
	}
	pack, err := chargeevidence.Read(data, identity)
	if err != nil {
		t.Fatalf("read published pack: %v", err)
	}
	return pack
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
