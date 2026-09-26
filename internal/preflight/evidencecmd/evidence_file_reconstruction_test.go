package evidencecmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// The fragment roles below and the fixture's expected patch table are stated by hand from
// the fixture trees, independently of the diff owner and the review collector.

type shapePatch = preflighttest.ShapePatch

var one = preflighttest.OnePatch

// reviewDiff prepares the review over args, retrieves every page through the evidence
// command, and returns the diff sources in shared-row order with their manifest descriptors.
func reviewDiff(t *testing.T, root string, args []string) (identity string, sources []chargeevidence.ManifestSource, joined string) {
	t.Helper()
	identity, _, _ = prepareEvidence(t, args)
	_, retrieved := reconstructEvidence(t, identity, traverseEvidence(t, identity))
	pack := preflighttest.PublishedPack(t, root, identity)
	byID := map[string]chargeevidence.ManifestSource{}
	for _, source := range pack.Manifest().Sources {
		byID[source.ID] = source
		if source.Role == "diff" {
			t.Fatalf("manifest still declares a monolithic diff source %s", source.ID)
		}
	}
	var body strings.Builder
	for _, row := range pack.Metadata().Shared {
		if row.Kind != "diff" {
			continue
		}
		source := byID[row.Source]
		if source.Kind != chargeevidence.KindGenerated || !source.Required {
			t.Fatalf("diff fragment %s = %#v, want a required generated source", row.Source, source)
		}
		sources = append(sources, source)
		body.WriteString(retrieved[row.Source])
	}
	return identity, sources, body.String()
}

// assertFragments checks the fragment descriptors: one prefix, the file paths in patch order,
// and one suffix.
func assertFragments(t *testing.T, sources []chargeevidence.ManifestSource, want []shapePatch) {
	t.Helper()
	var paths, wantPaths []string
	for _, patch := range want {
		for range patch.Patches {
			wantPaths = append(wantPaths, patch.Path)
		}
	}
	if len(sources) < 2 || sources[0].Role != "diff-prefix" || sources[0].Path != "diff-prefix" ||
		sources[len(sources)-1].Role != "diff-suffix" || sources[len(sources)-1].Path != "diff-suffix" {
		t.Fatalf("diff fragments = %#v, want a prefix first and a suffix last", sources)
	}
	for _, source := range sources[1 : len(sources)-1] {
		if source.Role != "diff-file" {
			t.Fatalf("middle fragment %#v is not a file patch", source)
		}
		paths = append(paths, source.Path)
	}
	if !slices.Equal(paths, wantPaths) {
		t.Fatalf("file patch paths =\n%q\nwant\n%q", paths, wantPaths)
	}
}

// TestReviewFileReconstruction is RE5, RE7, and RE8. The retrieved fragments rebuild the
// stored monolithic response byte for byte, and each path's patches equal Git's own patch.
// A later charge holds only its predecessor-to-tip diff. A multi-page Unicode patch rebuilds
// exactly through the page protocol. A body the partition cannot represent publishes nothing.
func TestReviewFileReconstruction(t *testing.T) {
	for _, test := range []struct {
		baseline string
		pair     preflighttest.ReviewPair
		want     []shapePatch
	}{
		// A review charge refuses an empty diff before any collector runs, so the empty
		// baseline belongs to the diff seam alone.
		{"review-documents", preflighttest.ReviewDocuments(), preflighttest.ReviewDocumentPatches()},
		{"review-shapes-renames", preflighttest.ReviewShapes(true), preflighttest.ReviewShapePatches(true)},
		{"review-shapes-no-renames", preflighttest.ReviewShapes(false), preflighttest.ReviewShapePatches(false)},
	} {
		t.Run(test.baseline, func(t *testing.T) {
			root, _, args := preflighttest.SeedReviewPair(t, test.pair)
			identity, sources, joined := reviewDiff(t, root, args)
			if want := string(preflighttest.ReviewBaseline(t, test.baseline)); joined != want {
				t.Fatalf("reconstructed diff differs from the stored baseline:\n%s\nwant\n%s", joined, want)
			}
			assertFragments(t, sources, test.want)
			members := fileMembers(t, root, identity)
			for _, patch := range test.want {
				if got, want := selectedPatch(t, identity, members[patch.Path]), rawPatch(t, args, patch.Pathspecs...); got != want {
					t.Errorf("%q patch =\n%s\nwant\n%s", patch.Path, got, want)
				}
			}
		})
	}
	for _, refusal := range []struct {
		name   string
		config func(t *testing.T) [2]string
		tip    map[string]preflighttest.TreeEntry
		want   string
	}{
		{name: "external diff", want: "diff evidence failed", config: func(t *testing.T) [2]string {
			script := filepath.Join(t.TempDir(), "external.sh")
			if err := os.WriteFile(script, []byte("#!/bin/sh\necho \"external $1\"\n"), 0o755); err != nil {
				t.Fatal(err)
			}
			return [2]string{"diff.external", script}
		}},
		{name: "forced color", want: "evidence", config: func(*testing.T) [2]string { return [2]string{"color.diff", "always"} }},
		{name: "control byte path", want: "unrepresentable", tip: map[string]preflighttest.TreeEntry{"notes/bad\x01.txt": preflighttest.Regular("x\n")}},
	} {
		t.Run(refusal.name+" refuses", func(t *testing.T) {
			pair := preflighttest.ReviewShapes(true)
			for path, entry := range refusal.tip {
				pair.Tip[path] = entry
			}
			if refusal.config != nil {
				pair.Config = append(pair.Config, refusal.config(t))
			}
			root, _, args := preflighttest.SeedReviewPair(t, pair)
			out, code := preflight.Command(args)
			if code != 1 || !strings.Contains(out, refusal.want) || strings.Contains(out, "prepared[") {
				t.Fatalf("%s preparation = (%d), want a refusal containing %q:\n%s", refusal.name, code, refusal.want, out)
			}
			preflighttest.AssertNothingPublished(t, root)
		})
	}
	t.Run("predecessor base", func(t *testing.T) {
		root, _, args := preflighttest.SeedReviewEvidence(t, false)
		commitReview(t, args, map[string]string{"notes/sentinel.txt": "earlier chunk\n"}, "chunk one")
		args[4] = args[6]
		commitReview(t, args, map[string]string{"notes/later.txt": "later chunk\n"}, "chunk two")
		_, sources, joined := reviewDiff(t, root, args)
		assertFragments(t, sources, []shapePatch{one("notes/later.txt")})
		if strings.Contains(joined, "sentinel") || !strings.Contains(joined, ","+args[4]+",explicit pair") {
			t.Fatalf("later charge diff holds the earlier chunk or another base:\n%s", joined)
		}
	})
	t.Run("large Unicode patch", func(t *testing.T) {
		root, _, args := preflighttest.SeedReviewEvidence(t, false)
		commitReview(t, args, map[string]string{"notes/unicode.txt": strings.Repeat("π évidence ✓ 🙂 ", 2000) + "\n"}, "unicode")
		identity, _, _ := prepareEvidence(t, args)
		members := fileMembers(t, root, identity)["notes/unicode.txt"]
		if len(members) != 1 || len(members[0].pages) < 3 {
			t.Fatalf("unicode members = %v, want one patch of at least three pages", members)
		}
		if got, want := selectedPatch(t, identity, members), rawPatch(t, args, "notes/unicode.txt"); got != want {
			t.Fatalf("unicode patch rebuilt %d of %d bytes", len(got), len(want))
		}
	})
}

// TestReviewFileSelectedStream is RE9. Selection by the declared path returns only that
// file's complete patches, for a pure rename, an empty addition and deletion, a binary
// change, a mode change, and a file-to-symlink change, under both rename settings.
func TestReviewFileSelectedStream(t *testing.T) {
	identities := []string{
		"notes/rename-pure.txt", "notes/renamed-pure.txt", "notes/empty-add.txt", "notes/empty-delete.txt",
		"notes/binary.bin", "notes/mode.sh", "notes/type-change",
	}
	for _, test := range []struct {
		renames bool
		cases   int
		absent  []string
	}{
		// Rename detection pairs the pure rename and the empty files, so each pair's base
		// path selects nothing.
		{true, 5, []string{"notes/rename-pure.txt", "notes/empty-delete.txt"}},
		{false, 7, nil},
	} {
		t.Run(fmt.Sprintf("renames=%v", test.renames), func(t *testing.T) {
			root, _, args := preflighttest.SeedReviewPair(t, preflighttest.ReviewShapes(test.renames))
			identity, _, _ := prepareEvidence(t, args)
			members := fileMembers(t, root, identity)
			cases := 0
			for _, patch := range preflighttest.ReviewShapePatches(test.renames) {
				if !slices.Contains(identities, patch.Path) {
					continue
				}
				cases++
				if len(members[patch.Path]) != patch.Patches {
					t.Errorf("%s members = %v, want %d", patch.Path, members[patch.Path], patch.Patches)
					continue
				}
				if got, want := selectedPatch(t, identity, members[patch.Path]), rawPatch(t, args, patch.Pathspecs...); got != want {
					t.Errorf("%s selected patch =\n%s\nwant\n%s", patch.Path, got, want)
				}
			}
			if cases != test.cases {
				t.Errorf("identity cases = %d, want %d", cases, test.cases)
			}
			for _, path := range test.absent {
				if len(members[path]) != 0 {
					t.Errorf("rename base path %s selects %v, want no source", path, members[path])
				}
			}
		})
	}
}
