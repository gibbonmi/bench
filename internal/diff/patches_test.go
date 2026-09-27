package diff

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestReviewFileReconstruction is RE5 at the diff seam. Over each pinned pair, the command
// still prints the stored monolithic response where one exists, the file patch partition
// rebuilds the command response byte for byte, and each path's patches equal Git's own
// patch for that path.
func TestReviewFileReconstruction(t *testing.T) {
	for _, test := range preflighttest.ReviewCases() {
		t.Run(test.Name, func(t *testing.T) {
			_, _, args := preflighttest.SeedReviewPair(t, test.Pair)
			want := test.Expected(t, Command, args)
			response, code := Command(preflighttest.DiffArgs(args))
			if code != 0 {
				t.Fatalf("command response = (%d):\n%s", code, response)
			}
			if response != want {
				t.Fatalf("command response =\n%s\nwant the expected response\n%s", response, want)
			}
			snapshot, out, code := PairPatches(preflighttest.DiffArgs(args))
			if code != 0 {
				t.Fatalf("partition = (%d):\n%s", code, out)
			}
			joined := string(snapshot.Prefix)
			var paths []string
			bodies := map[string]string{}
			for _, patch := range snapshot.Patches {
				joined += string(patch.Body)
				paths = append(paths, patch.Path)
				bodies[patch.Path] += string(patch.Body)
			}
			if joined += string(snapshot.Suffix); joined != response {
				t.Fatalf("partition rebuilt\n%s\nwant the command response\n%s", joined, response)
			}
			for _, patch := range test.Patches {
				if got, raw := patch.Selected(bodies), preflighttest.RawPatch(t, args, patch.Pathspecs...); got != raw {
					t.Errorf("%q patches =\n%s\nwant\n%s", patch.Path, got, raw)
				}
			}
			if want := preflighttest.ShapePaths(test.Patches); !slices.Equal(paths, want) {
				t.Fatalf("patch paths =\n%q\nwant\n%q", paths, want)
			}
		})
	}
}

// TestPartitionPatchIdentity pins the file identity rule for header shapes: the tip path of a
// surviving file and the base path of a deleted one, read from explicit rename or copy paths, the
// `diff --git` paths, or the side headers, with header-like content kept as content.
func TestPartitionPatchIdentity(t *testing.T) {
	for _, test := range []struct {
		name, body string
		inventory  []string
		want       string
	}{
		{"pure rename", "diff --git a/old b/new\nsimilarity index 100%\nrename from old\nrename to new\n", []string{"new", "old"}, "new"},
		{"quoted rename", "diff --git \"a/o\\tld\" b/new name\nsimilarity index 100%\nrename from \"o\\tld\"\nrename to new name\n", []string{"new name", "o\tld"}, "new name"},
		{"empty addition", "diff --git a/e f b/e f\nnew file mode 100644\nindex 0000000..e69de29\n", []string{"e f"}, "e f"},
		{"empty deletion", "diff --git \"a/caf\\303\\251\" \"b/caf\\303\\251\"\ndeleted file mode 100644\nindex e69de29..0000000\n", []string{"café"}, "café"},
		{"mode change", "diff --git a/a b/a b/a b/a\nold mode 100644\nnew mode 100755\n", []string{"a b/a"}, "a b/a"},
		{"binary", "diff --git a/b.bin b/b.bin\nnew file mode 100644\nindex 0000000..1111111\nBinary files /dev/null and b/b.bin differ\n", []string{"b.bin"}, "b.bin"},
		{"deleted with sides", "diff --git a/gone b/gone\ndeleted file mode 100644\nindex 1111111..0000000\n--- a/gone\n+++ /dev/null\n@@ -1 +0,0 @@\n-diff --git a/fake b/fake\n", []string{"gone"}, "gone"},
		{"copy", "diff --git a/src b/dst\nsimilarity index 100%\ncopy from src\ncopy to dst\n", []string{"dst", "src"}, "dst"},
		{"trailing-space sides and a bare blank context line", "diff --git a/x  b/x \nindex 1..2 100644\n--- a/x \t\n+++ b/x \t\n@@ -1,3 +1,3 @@\n a\n\n-b\n+c\n", []string{"x "}, "x "},
		{"header-like content", "diff --git a/x b/x\nindex 1..2 100644\n--- a/x\n+++ b/x\n@@ -1 +1,3 @@\n a\n+rename to y\n+--- a/z\n\\ No newline at end of file\n", []string{"x"}, "x"},
	} {
		t.Run(test.name, func(t *testing.T) {
			patches, err := partitionPatches([]byte(test.body), test.inventory)
			if err != nil || len(patches) != 1 || patches[0].Path != test.want || string(patches[0].Body) != test.body {
				t.Fatalf("partition = %#v, %v; want one verbatim patch of %q", patches, err, test.want)
			}
		})
	}
}

// TestPartitionPatchRefusals proves that a body the partition cannot represent refuses as a
// whole: colored or external output, a torn body, and a patch set that disagrees with the
// frozen inventory in either direction.
func TestPartitionPatchRefusals(t *testing.T) {
	valid := "diff --git a/x b/x\nindex 1..2 100644\n--- a/x\n+++ b/x\n@@ -1 +1 @@\n-a\n+b\n"
	for _, test := range []struct {
		name, body string
		inventory  []string
		want       string
	}{
		{"forced color", "\x1b[1mdiff --git a/x b/x\x1b[m\n", []string{"x"}, "does not start a patch"},
		{"external diff", "external x\n", []string{"x"}, "does not start a patch"},
		{"torn body", strings.TrimSuffix(valid, "\n"), []string{"x"}, "line feed"},
		{"side header without tip", "diff --git a/x b/x\n--- a/x\n", []string{"x"}, "tip side"},
		{"content before a hunk", "diff --git a/x b/x\nindex 1..2\n+b\n", []string{"x"}, "neither a hunk"},
		{"foreign content line", valid + "stray\n", []string{"x"}, "not patch content"},
		{"path outside the inventory", valid, []string{"y"}, "does not hold"},
		{"inventory path without a patch", valid, []string{"x", "y"}, "has no patch"},
	} {
		t.Run(test.name, func(t *testing.T) {
			patches, err := partitionPatches([]byte(test.body), test.inventory)
			if err == nil || !strings.Contains(err.Error(), test.want) || patches != nil {
				t.Fatalf("partition = %#v, %v; want a whole refusal containing %q", patches, err, test.want)
			}
		})
	}
}
