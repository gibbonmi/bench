package conformance

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/guards"
)

func checkGuardHeaderManifests(root string) []string {
	var diags []string
	for _, guard := range []struct {
		name string
		path string
	}{
		{"block-dangerous-git", filepath.Join(root, ".bench", "hooks", "block-dangerous-git.sh")},
		{"block-bench-follow-on", filepath.Join(root, ".bench", "hooks", "block-bench-follow-on.sh")},
		{"block-unsafe-rg", filepath.Join(root, ".bench", "hooks", "block-unsafe-rg.sh")},
		{"check-agent-line", filepath.Join(root, ".bench", "hooks", "check-agent-line.sh")},
		{"stop", filepath.Join(root, ".bench", "hooks", "stop.sh")},
		{"session-start", filepath.Join(root, ".bench", "hooks", "session-start.sh")},
		{"worktree-lifecycle", filepath.Join(root, ".bench", "hooks", "worktree-lifecycle.sh")},
		{"pre-push", filepath.Join(root, "internal", "adopt", "prepush.sh")},
	} {
		path := guard.path
		if !exists(path) {
			continue
		}
		fields, err := guards.HeaderFields(path)
		if err != nil {
			diags = append(diags, fmt.Sprintf("guard %s manifest unreadable", guard.name))
			continue
		}
		for _, key := range fields.MissingRequired() {
			diags = append(diags, fmt.Sprintf("guard %s manifest missing %s key", guard.name, key))
		}
		if guard.name == "session-start" && fields["denies"] != "nothing (informational)" {
			diags = append(diags, "session-start is not classified informational (denies: nothing)")
		}
	}
	return diags
}

// TestGuardHeaderManifestsCoverTheRgGuard covers R41. The rg shim's `# denies:` line is
// the shipped advertisement of its rule, so a shim header without that key must red.
func TestGuardHeaderManifestsCoverTheRgGuard(t *testing.T) {
	const header = "#!/usr/bin/env bash\n# name: block-unsafe-rg\n# boundary: PreToolUse:Bash\n# why: x\nexit 0\n"
	root := throwawayRoot{files: map[string]string{".bench/hooks/block-unsafe-rg.sh": header}}.build(t)

	diags := checkGuardHeaderManifests(root)

	if want := "guard block-unsafe-rg manifest missing denies key"; !containsDiagnostic(diags, want) {
		t.Fatalf("a rg shim without a denies key did not bite with %q:\n%s", want, strings.Join(diags, "\n"))
	}
}
