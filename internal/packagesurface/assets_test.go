package packagesurface

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRequiredPackAssetsIncludeFollowOnGuard(t *testing.T) {
	const want = ".bench/hooks/block-bench-follow-on.sh"
	for _, asset := range RequiredPackAssets {
		if asset == want {
			return
		}
	}
	t.Fatalf("RequiredPackAssets omit %q", want)
}

// embedFixtureRoot builds the smallest tree the pack-asset derivations walk: a
// module root with the two source directories they descend, and one internal
// source whose //go:embed names a sibling file. The sibling is what a
// root-relative resolution would miss, so the fixture is the discriminating case
// for both derivations.
func embedFixtureRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{"cmd", filepath.Join("internal", "adopt")} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	source := "package adopt\n\nimport _ \"embed\"\n\n//go:embed prepush.sh\nvar prepushHook string\n"
	for path, body := range map[string]string{
		filepath.Join("internal", "adopt", "link_hook.go"): source,
		filepath.Join("internal", "adopt", "prepush.sh"):   "#!/bin/sh\n",
	} {
		if err := os.WriteFile(filepath.Join(root, path), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// TestRequiredBuildPackAssetsCarryEmbedTargets pins the composed pack list to the
// embed rule. The list is what the npm git install packages, so an embed target the
// list omits ships a binary that cannot build.
func TestRequiredBuildPackAssetsCarryEmbedTargets(t *testing.T) {
	const want = "internal/adopt/prepush.sh"
	assets, err := RequiredBuildPackAssets(embedFixtureRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, asset := range assets {
		if asset == want {
			return
		}
	}
	t.Fatalf("RequiredBuildPackAssets = %v, want it to carry %q", assets, want)
}

// TestEmbedTargetsResolveAgainstTheSourceDirectory pins the exported derivation the
// path-aware lane consumes. A pattern is directory-relative, so a resolution against
// the module root names a path no checkout carries and the lane grades the wrong file.
func TestEmbedTargetsResolveAgainstTheSourceDirectory(t *testing.T) {
	targets, err := EmbedTargets(embedFixtureRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"internal/adopt/prepush.sh"}
	if !reflect.DeepEqual(targets, want) {
		t.Fatalf("EmbedTargets = %v, want %v", targets, want)
	}
}
