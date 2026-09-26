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
func reviewCanonical(slug string) []fixtureFile {
	return []fixtureFile{
		{"go.mod", "module example.com/review\n\ngo 1.25\n"},
		{"specs/" + slug + "/spec.md", SpecBody(slug, reviewFenceLines()...)},
		{"specs/" + slug + "/tickets/one.md", WritesTicketDoc("One", FenceWrites(ReviewFence()), "PF1", "PF2")},
		{chargesource.DelegateSkill, "# Delegation skill\n"},
		{chargesource.DelegateProcedure, "# Delegation procedure\n\nFocused suite: bench test --package ./internal/preflight\n"},
		{chargesource.BuildPhase, "# Build phase\n"},
		{chargesource.ReviewSkill, "# Review skill\n\n## Standards\n\nRules.\n\n## Spec\n\nRequirements.\n\n## Coverage\n\nEdges.\n"},
		{chargesource.ReviewPhase, "# Review phase\n\nUse the three canonical axes.\n"},
	}
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
	slug = "example"
	root = t.TempDir()
	t.Chdir(root)
	RunGit(t, "init", "-q", "--object-format=sha1", "-b", "main")
	base := map[string]TreeEntry{}
	for _, file := range reviewCanonical(slug) {
		base[file.path] = Regular(file.body)
	}
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

// renameEditBody is the ten-line body the edited rename keeps nine lines of.
var renameEditBody = "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n"

// hostilePaths are the added paths whose bytes need quoting, escaping, or pathspec care:
// spaces, glob characters, a quote, a backslash, non-ASCII bytes, and the tab, newline, and
// return bytes the current renderer permits.
var hostilePaths = []string{
	"notes/space name.txt", "notes/glob*[?].txt", "notes/quote\"d.txt", "notes/back\\slash.txt",
	"notes/café π.txt", "notes/tab\there.txt", "notes/new\nline.txt", "notes/carriage\rreturn.txt",
}

// ReviewShapes is the RE5 edge inventory as one pinned pair: added, modified, deleted, and
// renamed files, a pure rename, an empty addition and deletion, a mode-only change, a binary
// change, a committed symlink, a file-to-symlink change, a missing final newline, content
// that resembles patch headers, and hostile paths. renames sets diff.renames for the pair.
func ReviewShapes(renames bool) ReviewPair {
	pair := ReviewPair{
		Base: map[string]TreeEntry{
			"notes/modify.txt":       Regular("one\ntwo\nthree\n"),
			"notes/delete.txt":       Regular("delete me\n"),
			"notes/rename-pure.txt":  Regular("pure rename body\nline two\nline three\n"),
			"notes/rename-edit.txt":  Regular(renameEditBody),
			"notes/empty-delete.txt": Regular(""),
			"notes/mode.sh":          Regular("#!/bin/sh\n"),
			"notes/binary.bin":       Regular("\x00\x01binary"),
			"notes/type-change":      Regular("regular\n"),
			"notes/no-newline.txt":   Regular("last"),
			"notes/lookalike.txt":    Regular("a\n"),
		},
		Removed: []string{
			"notes/delete.txt", "notes/rename-pure.txt", "notes/rename-edit.txt", "notes/empty-delete.txt",
		},
		Tip: map[string]TreeEntry{
			"notes/modify.txt":       Regular("one\ntwo changed\nthree\n"),
			"notes/renamed-pure.txt": Regular("pure rename body\nline two\nline three\n"),
			"notes/renamed-edit.txt": Regular(strings.Replace(renameEditBody, "five", "FIVE", 1)),
			"notes/mode.sh":          Executable("#!/bin/sh\n"),
			"notes/binary.bin":       Regular("\x00\x02binary"),
			"notes/type-change":      Symlink("modify.txt"),
			"notes/no-newline.txt":   Regular("last line changed"),
			"notes/lookalike.txt":    Regular("a\ndiff --git a/x b/x\n--- a/x\n+++ b/x\nrename to y\n@@ -1 +1 @@\n"),
			"notes/add.txt":          Regular("added\n"),
			"notes/empty-add.txt":    Regular(""),
			"notes/link":             Symlink("add.txt"),
		},
		Config: [][2]string{{"diff.renames", fmt.Sprint(renames)}},
	}
	for _, path := range hostilePaths {
		pair.Tip[path] = Regular("hostile\n")
	}
	return pair
}

// ShapePatch is one expected file identity: its declared path, the pathspecs whose raw Git
// patch equals every source of that path joined, and its number of patch sources.
type ShapePatch struct {
	Path      string
	Pathspecs []string
	Patches   int
}

// OnePatch is an expected identity whose raw patch is the patch of its own path.
func OnePatch(path string) ShapePatch { return ShapePatch{path, []string{path}, 1} }

// RenamedPatch is an expected rename identity: the tip path, found by both spellings.
func RenamedPatch(from, to string) ShapePatch { return ShapePatch{to, []string{from, to}, 1} }

// typeChangePatches is the file-to-symlink identity: two patches for one path.
var typeChangePatches = ShapePatch{"notes/type-change", []string{"notes/type-change"}, 2}

// shapesWithRenames is the patch order of ReviewShapes with rename detection enabled.
var shapesWithRenames = []ShapePatch{
	OnePatch("notes/add.txt"), OnePatch("notes/back\\slash.txt"), OnePatch("notes/binary.bin"), OnePatch("notes/café π.txt"),
	OnePatch("notes/carriage\rreturn.txt"), OnePatch("notes/delete.txt"),
	RenamedPatch("notes/empty-delete.txt", "notes/empty-add.txt"), OnePatch("notes/glob*[?].txt"), OnePatch("notes/link"),
	OnePatch("notes/lookalike.txt"), OnePatch("notes/mode.sh"), OnePatch("notes/modify.txt"), OnePatch("notes/new\nline.txt"),
	OnePatch("notes/no-newline.txt"), OnePatch("notes/quote\"d.txt"),
	RenamedPatch("notes/rename-edit.txt", "notes/renamed-edit.txt"),
	RenamedPatch("notes/rename-pure.txt", "notes/renamed-pure.txt"), OnePatch("notes/space name.txt"),
	OnePatch("notes/tab\there.txt"), typeChangePatches,
}

// shapesWithoutRenames is the patch order of ReviewShapes with rename detection disabled.
var shapesWithoutRenames = []ShapePatch{
	OnePatch("notes/add.txt"), OnePatch("notes/back\\slash.txt"), OnePatch("notes/binary.bin"), OnePatch("notes/café π.txt"),
	OnePatch("notes/carriage\rreturn.txt"), OnePatch("notes/delete.txt"), OnePatch("notes/empty-add.txt"),
	OnePatch("notes/empty-delete.txt"), OnePatch("notes/glob*[?].txt"), OnePatch("notes/link"), OnePatch("notes/lookalike.txt"),
	OnePatch("notes/mode.sh"), OnePatch("notes/modify.txt"), OnePatch("notes/new\nline.txt"), OnePatch("notes/no-newline.txt"),
	OnePatch("notes/quote\"d.txt"), OnePatch("notes/rename-edit.txt"), OnePatch("notes/rename-pure.txt"),
	OnePatch("notes/renamed-edit.txt"), OnePatch("notes/renamed-pure.txt"), OnePatch("notes/space name.txt"),
	OnePatch("notes/tab\there.txt"), typeChangePatches,
}

// ReviewShapePatches is the expected patch table of ReviewShapes under one rename setting,
// stated by hand from the pair's trees and Git's path order.
func ReviewShapePatches(renames bool) []ShapePatch {
	if renames {
		return shapesWithRenames
	}
	return shapesWithoutRenames
}

// ReviewDocuments is a pinned pair whose diff holds only the spec and the review record.
func ReviewDocuments() ReviewPair {
	return ReviewPair{Tip: map[string]TreeEntry{
		"specs/example/spec.md": Regular(SpecBody("example", reviewFenceLines()...) + "\n## Further notes\n\nEdited.\n"),
		"reviews/example.md":    Regular("# Review record\n"),
	}}
}

// ReviewDocumentPatches is the expected patch table of ReviewDocuments.
func ReviewDocumentPatches() []ShapePatch {
	return []ShapePatch{OnePatch("reviews/example.md"), OnePatch("specs/example/spec.md")}
}

// DiffArgs is the diff collector invocation over the pair that review arguments pin.
func DiffArgs(reviewArgs []string) []string {
	return []string{"--base", reviewArgs[4], "--source-tip", reviewArgs[6], "--full"}
}

// RawGit runs one git command in the working directory and returns its exact standard
// output, so a test compares a stored patch with the bytes Git itself prints.
func RawGit(t *testing.T, args ...string) string {
	t.Helper()
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
	for _, name := range PublishedPacks(t, root) {
		data, err := os.ReadFile(filepath.Join(StoreDir(t, root), name))
		if err != nil {
			t.Fatal(err)
		}
		if pack, err := chargeevidence.Read(data, identity); err == nil {
			return pack
		}
	}
	t.Fatalf("no published pack holds %s", identity)
	return nil
}

// storedBaselines holds each stored full diff response. The files are embedded, so a test
// reads them from any working directory and under a trimmed build path.
//
//go:embed testdata/*.toon
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
