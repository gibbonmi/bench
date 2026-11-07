package conformance

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/guards"
)

// checkGuardHeaderManifests grades the manifest header of every hook-script guard that
// guards.HookScripts finds on disk, plus the pre-push hook source. A new hook therefore
// cannot escape the check.
func checkGuardHeaderManifests(root string) []string {
	all := append(guards.HookScripts(root), guards.HookScript{Name: "pre-push", Path: filepath.Join(root, "internal", "adopt", "prepush.sh")})
	var diags []string
	for _, guard := range all {
		path := guard.Path
		if !exists(path) {
			continue
		}
		fields, err := guards.HeaderFields(path)
		if err != nil {
			diags = append(diags, fmt.Sprintf("guard %s manifest unreadable", guard.Name))
			continue
		}
		for _, key := range fields.MissingRequired() {
			diags = append(diags, fmt.Sprintf("guard %s manifest missing %s key", guard.Name, key))
		}
		if guard.Name == "session-start" && fields["denies"] != "nothing (informational)" {
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

// TestGuardHeaderManifestsGradeEveryHookOnDisk proves that the check reads the guard
// hooks from disk. A hook that no list names, with no manifest header, must red.
func TestGuardHeaderManifestsGradeEveryHookOnDisk(t *testing.T) {
	root := throwawayRoot{files: map[string]string{".bench/hooks/unlisted-guard.sh": "#!/usr/bin/env bash\nexit 0\n"}}.build(t)

	diags := checkGuardHeaderManifests(root)

	if want := "guard unlisted-guard manifest missing denies key"; !containsDiagnostic(diags, want) {
		t.Fatalf("a hook without a manifest header escaped the check, want %q:\n%s", want, strings.Join(diags, "\n"))
	}
}
