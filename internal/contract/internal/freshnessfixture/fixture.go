// Package freshnessfixture builds contract subjects with fresh or deliberately stale seals.
package freshnessfixture

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
)

// PublishedSubject returns a minimal selected subject with a matching executable and seal.
func PublishedSubject(t testing.TB, output string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "subject [glob]* path")
	files := map[string]string{
		"go.mod":                  "module freshnesssubject\n\ngo 1.25\n",
		"cmd/bench/main.go":       "package main\n\nfunc main() {}\n",
		"scripts/go-build.sh":     "#!/usr/bin/env bash\n",
		"scripts/go-build.inputs": "build_script=scripts/go-build.sh\npackage_version=package.json\ngo_requirements=internal/releaseevidence/requirements.json\n",
		"package.json":            "{}\n",
		"internal/releaseevidence/requirements.json": "{}\n",
		"bin/bench.sh": "#!/bin/sh\nexec \"$(dirname \"$0\")/../dist/bench\" \"$@\"\n",
		"dist/staged":  "#!/bin/sh\nprintf '%s\\n' '" + output + "'\n",
	}
	for path, contents := range files {
		mode := os.FileMode(0o644)
		if path == "bin/bench.sh" || path == "dist/staged" || path == "scripts/go-build.sh" {
			mode = 0o755
		}
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(contents), mode); err != nil {
			t.Fatal(err)
		}
	}
	if err := freshness.Publish(root, filepath.Join(root, "dist", "staged"), filepath.Join(root, "dist", "bench")); err != nil {
		t.Fatalf("publish subject: %v", err)
	}
	return root
}

// StaleSubject returns a minimal selected subject whose executable predates its sources.
func StaleSubject(t testing.TB, output string) string {
	t.Helper()
	root := PublishedSubject(t, output)
	if err := os.WriteFile(filepath.Join(root, "cmd", "bench", "main.go"), []byte("package main\n\nfunc main() { println(\"changed\") }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}
