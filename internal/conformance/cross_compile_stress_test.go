//go:build stress

package conformance

import (
	"fmt"
	"os"
	"path/filepath"
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
