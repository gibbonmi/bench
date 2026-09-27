package preflighttest

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
)

// fixtureFile is one repository file of a seeded review tree.
type fixtureFile struct{ path, body string }

// reviewCanonical is the canonical review tree every review fixture seeds: the module, the
// conformant spec and ticket, and each charge source. SeedReviewEvidence writes it into the
// worktree, and SeedReviewPair commits it into a pinned base, so both fixtures state it once.
// Extra fence paths join both the spec fence and the ticket's Writes line.
func reviewCanonical(slug string, extraFence ...string) []fixtureFile {
	fenceLines := reviewFenceLines()
	for _, path := range extraFence {
		fenceLines = append(fenceLines, fenceLine(path, "review fixture"))
	}
	return []fixtureFile{
		{"go.mod", "module example.com/review\n\ngo 1.25\n"},
		{"specs/" + slug + "/spec.md", SpecBody(slug, fenceLines...)},
		{"specs/" + slug + "/tickets/one.md", WritesTicketDoc("One", FenceWrites(append(ReviewFence(), extraFence...)), "PF1", "PF2")},
		{chargesource.DelegateSkill, "# Delegation skill\n"},
		{chargesource.DelegateProcedure, "# Delegation procedure\n\nFocused suite: bench test --package ./internal/preflight\n"},
		{chargesource.BuildPhase, "# Build phase\n"},
		{chargesource.ReviewSkill, "# Review skill\n\n## Standards\n\nRules.\n\n## Spec\n\nRequirements.\n\n## Coverage\n\nEdges.\n"},
		{chargesource.ReviewPhase, "# Review phase\n\nUse the three canonical axes.\n"},
	}
}

// pairSlug is the spec slug of every pinned review pair, and pairSpecPath is its spec.
const (
	pairSlug     = "example"
	pairSpecPath = "specs/" + pairSlug + "/spec.md"
)

// ReviewFenceWith is the canonical review tree of a pinned pair with extra fence paths, for
// a pair whose changed paths lie outside the canonical fence. A pair sets it in Base, so both
// commits declare the same fence. SeedReviewPair seeds every base from it.
func ReviewFenceWith(paths ...string) map[string]TreeEntry {
	entries := map[string]TreeEntry{}
	for _, file := range reviewCanonical(pairSlug, paths...) {
		entries[file.path] = Regular(file.body)
	}
	return entries
}

// TreeEntry is one blob of a pinned review tree: its Git file mode and its exact bytes.
type TreeEntry struct{ Mode, Data string }

// Regular is a regular file entry.
func Regular(data string) TreeEntry { return TreeEntry{"100644", data} }

// Executable is an executable file entry.
func Executable(data string) TreeEntry { return TreeEntry{"100755", data} }

// Symlink is a symbolic link entry whose blob is its target.
func Symlink(target string) TreeEntry { return TreeEntry{"120000", target} }

// ReviewPair is one deterministic review fixture. The base commit holds the canonical review
// tree plus Base. The tip commit holds the base tree without Removed, overlaid by Tip. Config
// is repository configuration that the fixture sets after the pinned configuration.
type ReviewPair struct {
	Base, Tip map[string]TreeEntry
	Removed   []string
	Config    [][2]string
}

// pinnedConfig is the repository configuration every pinned review pair sets, so ambient user
// configuration cannot move a stored baseline byte. Local configuration overrides global.
var pinnedConfig = [][2]string{
	{"core.quotePath", "true"}, {"core.abbrev", "12"}, {"core.autocrlf", "false"},
	{"diff.renames", "true"}, {"diff.noprefix", "false"}, {"diff.mnemonicPrefix", "false"},
	{"diff.algorithm", "myers"}, {"diff.context", "3"}, {"diff.interHunkContext", "0"},
	{"diff.indentHeuristic", "true"}, {"diff.suppressBlankEmpty", "false"}, {"diff.relative", "false"},
	{"color.ui", "false"}, {"color.diff", "false"},
}

// pinnedCommitter is the identity and each commit time of a pinned review pair. The stream
// states them, so no process environment and no user configuration reaches the commits.
const pinnedCommitter = "Fixture <fixture@example.com>"

var pinnedTimes = [2]string{"978307200 +0000", "978393600 +0000"}

// SeedReviewPair builds the pinned pair in a fresh repository at t.TempDir(), makes that
// repository the working directory, checks out the tip, and registers the active assignment.
// Equal pairs produce equal commit identities, so a stored response stays comparable. It
// returns the root, the slug, and the review preparation arguments over the pair.
func SeedReviewPair(t *testing.T, pair ReviewPair) (root, slug string, args []string) {
	t.Helper()
	slug = pairSlug
	root = t.TempDir()
	t.Chdir(root)
	RunGit(t, "init", "-q", "--object-format=sha1", "-b", "main")
	base := ReviewFenceWith()
	for path, entry := range pair.Base {
		base[path] = entry
	}
	tip := map[string]TreeEntry{}
	for path, entry := range base {
		tip[path] = entry
	}
	for _, path := range pair.Removed {
		delete(tip, path)
	}
	for path, entry := range pair.Tip {
		tip[path] = entry
	}
	var stream strings.Builder
	writeCommit(&stream, "refs/heads/main", 1, pinnedTimes[0], "base", base)
	writeCommit(&stream, "refs/heads/feature", 2, pinnedTimes[1], "review source", tip)
	importer := exec.Command("git", "fast-import", "--quiet")
	importer.Stdin = strings.NewReader(stream.String())
	if out, err := importer.CombinedOutput(); err != nil {
		t.Fatalf("git fast-import: %v\n%s", err, out)
	}
	RunGit(t, "checkout", "-q", "-f", "feature")
	for _, setting := range append(append([][2]string{}, pinnedConfig...), pair.Config...) {
		RunGit(t, "config", setting[0], setting[1])
	}
	ActiveAssignment(t, root, root)
	return root, slug, ReviewArgs(t, slug)
}

// writeCommit appends one fast-import commit that replaces the whole tree with entries.
func writeCommit(b *strings.Builder, ref string, mark int, when, message string, entries map[string]TreeEntry) {
	fmt.Fprintf(b, "commit %s\nmark :%d\ncommitter %s %s\ndata %d\n%s\n", ref, mark, pinnedCommitter, when, len(message), message)
	if mark > 1 {
		fmt.Fprintf(b, "from :%d\n", mark-1)
	}
	b.WriteString("deleteall\n")
	paths := make([]string, 0, len(entries))
	for path := range entries {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		entry := entries[path]
		fmt.Fprintf(b, "M %s inline %s\ndata %d\n%s\n", entry.Mode, quoteImportPath(path), len(entry.Data), entry.Data)
	}
}

// quoteImportPath renders a path in the C-style quoting fast-import reads, so any byte a
// hostile path holds reaches the tree unchanged.
func quoteImportPath(path string) string {
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(path); i++ {
		switch c := path[i]; {
		case c == '"' || c == '\\':
			b.WriteByte('\\')
			b.WriteByte(c)
		case c < 0x20 || c == 0x7f:
			fmt.Fprintf(&b, "\\%03o", c)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// DiffArgs is the diff collector invocation over the pair that review arguments pin.
func DiffArgs(reviewArgs []string) []string {
	return []string{"--base", reviewArgs[4], "--source-tip", reviewArgs[6], "--full"}
}

// RawPatch is Git's own patch for pathspecs over the frozen pair that review arguments pin,
// read in the working directory, so a test compares a stored patch with the bytes Git prints.
func RawPatch(t *testing.T, reviewArgs []string, pathspecs ...string) string {
	t.Helper()
	args := append([]string{"--literal-pathspecs", "diff", reviewArgs[4], reviewArgs[6], "--"}, pathspecs...)
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return string(out)
}

// IdentifiedPack reads the published artifact of identity through the strict reader when the
// store holds several, as a series of preparations over one repository does.
func IdentifiedPack(t *testing.T, root, identity string) *chargeevidence.Pack {
	t.Helper()
	var refusals []string
	for _, name := range PublishedPacks(t, root) {
		data, err := os.ReadFile(filepath.Join(StoreDir(t, root), name))
		if err != nil {
			t.Fatal(err)
		}
		pack, err := chargeevidence.Read(data, identity)
		if err == nil {
			return pack
		}
		refusals = append(refusals, err.Error())
	}
	t.Fatalf("no published pack holds %s; read refusals: %v", identity, refusals)
	return nil
}

// storedBaselines holds each stored full diff response. The files are embedded, so a test
// reads them from any working directory and under a trimmed build path.
//
//go:embed testdata/review-documents.toon testdata/review-empty.toon testdata/review-shapes-no-renames.toon testdata/review-shapes-renames.toon
var storedBaselines embed.FS

// ReviewBaseline reads one stored full diff response. Each baseline is the output of the
// monolithic producer before file partitioning, captured once over a pinned pair, so a test
// compares current bytes with an expectation no current owner derived.
func ReviewBaseline(t *testing.T, name string) []byte {
	t.Helper()
	data, err := storedBaselines.ReadFile("testdata/" + name + ".toon")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// Expected is the full diff response the case's pair must rebuild: the stored baseline where
// one exists, else the current response of command over the pair that review arguments pin.
// The caller passes the diff command, so this package does not import the diff owner.
func (c ReviewCase) Expected(t *testing.T, command func([]string) (string, int), reviewArgs []string) string {
	t.Helper()
	if c.Stored {
		return string(ReviewBaseline(t, c.Name))
	}
	out, code := command(DiffArgs(reviewArgs))
	if code != 0 {
		t.Fatalf("command response = (%d):\n%s", code, out)
	}
	return out
}
