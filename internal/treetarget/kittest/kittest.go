// Package kittest writes the smallest kit tree for the tests of a kit worktree target. A kit
// tree declares Bench build inputs, so freshness.DeclaresBuildInputs answers true for it and
// freshness.Digest hashes it. Only tests import this package. It does not import
// internal/treetarget, so the internal tests of that package and the system suite both can
// import it.
package kittest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
)

// Wrapper is the path of the wrapper in a kit tree, relative to its root.
const Wrapper = "bin/bench.sh"

// buildInput is the one build input that the manifest of a kit tree lists.
const buildInput = "scripts/go-build.sh"

// WriteTree writes a kit tree into root: a Go module with a ./cmd/bench package, the
// build-input manifest that lists one input, that input, and the wrapper. It also makes the
// empty directory that the build publishes into, so a test publishes the build with one
// freshness.Publish call.
func WriteTree(t testing.TB, root string) {
	t.Helper()
	for name, body := range map[string]string{
		"go.mod":                      "module benchkit\n\ngo 1.21\n",
		"cmd/bench/main.go":           "package main\n\nfunc main() {}\n",
		freshness.BuildInputsManifest: freshness.BuildInputLine("build_script", buildInput),
		buildInput:                    "#!/bin/sh\n",
		Wrapper:                       "#!/bin/sh\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(freshness.PublishedExecutable(root)), 0o755); err != nil {
		t.Fatal(err)
	}
}

// EditBuildInput changes the listed build input of the kit tree at root, so a build that
// was published before the edit no longer matches the tree.
func EditBuildInput(t testing.TB, root string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(buildInput)), []byte("#!/bin/sh\n# edited\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}
