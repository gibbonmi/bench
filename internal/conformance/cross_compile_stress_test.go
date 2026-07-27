//go:build stress

package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// crossCompileMatrix builds the binary for every platform in scripts/release-plan.json and
// reports any target that fails to build. It is heavy (~8s, four serial toolchain builds)
// and portability rarely breaks on an ordinary code change, so it is gold-plating on the
// every-commit gate. It runs only under `go test -tags stress`; the release workflow,
// which cross-compiles every target to ship the platform packages, is the standing
// backstop that catches a real portability break before it can ship.
func crossCompileMatrix(root, buildHelper string) []string {
	if !exists(filepath.Join(root, "scripts", "release-plan.json")) || !exists(buildHelper) {
		return nil
	}
	var diags []string
	matrix, err := platformMatrix(filepath.Join(root, "scripts", "release-plan.json"))
	if err != nil {
		diags = append(diags, "platform matrix unreadable: "+err.Error())
	}
	tmp, err := os.MkdirTemp("", "bench-cross-*")
	if err != nil {
		return append(diags, "cross-compile setup failed: "+err.Error())
	}
	defer os.RemoveAll(tmp)
	for _, target := range matrix {
		env := append(conformanceSubprocessEnv(), "GOOS="+target.Goos, "GOARCH="+target.Goarch)
		probe := runAtEnv(root, env, "bash", buildHelper, root, filepath.Join(tmp, "bench-"+target.Goos+"-"+target.Goarch))
		if probe == nil || probe.ExitCode != 0 {
			diags = append(diags, fmt.Sprintf("cross-compile failed: %s/%s", target.Goos, target.Goarch))
		}
	}
	return diags
}

// TestResidualCheckKeepsCrossCompile is the tripwire for the one thing the toolchain
// split could drop in silence. Cross-compile is a no-op on the dev tier and owns no
// canary fixture, so a subtraction that took the matrix call out with the steps around
// it would leave every other assertion green while ship lost the four-platform matrix.
func TestResidualCheckKeepsCrossCompile(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeFixtureFile(t, filepath.Join(root, "scripts", "release-plan.json"),
		`{"targets":[{"goos":"linux","goarch":"amd64"}]}`+"\n")
	// A helper that refuses every target: the assertion is that the matrix is reached,
	// not that a real cross-compile succeeds, and a real one would cost four builds.
	writeFixtureFile(t, filepath.Join(root, "scripts", "go-build.sh"),
		"#!/usr/bin/env bash\nexit 1\n")

	diags := checkGoToolchain(root)

	if !containsDiagnostic(diags, "cross-compile failed: linux/amd64") {
		t.Fatalf("residual check no longer reaches the cross-compile matrix: %#v", diags)
	}
}
