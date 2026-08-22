//go:build stress

package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// crossCompileMatrix builds the binary for every platform in scripts/release-plan.json.
// It reports any target that fails to build. It is heavy, at about 8s across four serial
// toolchain builds, and portability rarely breaks on an ordinary code change. It runs
// only under `go test -tags stress`. The release workflow cross-compiles every target to
// ship the platform packages. That workflow is the standing backstop that catches a real
// portability break before it can ship.
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

// TestResidualCheckKeepsCrossCompile grades that the residual check drives a matrix that
// works. It runs only under `-tags stress`, where the matrix is not the dev tier's no-op.
// prep-release's ship `-run` filter executes it. The dev tier's
// TestResidualCheckCallsCrossCompileMatrix grades a different fact: that the call site is
// still there. That fact is all a tier whose matrix returns nil can observe.
func TestResidualCheckKeepsCrossCompile(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, "go.mod"), "module fixture\n\ngo 1.25\n")
	writeFixtureFile(t, filepath.Join(root, "scripts", "release-plan.json"),
		`{"targets":[{"goos":"linux","goarch":"amd64"}]}`+"\n")
	// This helper refuses every target. The assertion is that the matrix is reached, not
	// that a real cross-compile succeeds. A real cross-compile would cost four builds.
	writeFixtureFile(t, filepath.Join(root, "scripts", "go-build.sh"),
		"#!/usr/bin/env bash\nexit 1\n")

	diags := checkGoToolchain(root)

	if !containsDiagnostic(diags, "cross-compile failed: linux/amd64") {
		t.Fatalf("residual check no longer reaches the cross-compile matrix: %#v", diags)
	}
}
