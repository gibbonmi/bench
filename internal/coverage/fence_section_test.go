package coverage

import "testing"

// TestCommandCheckIgnoresAFenceSection covers CR70. The ticket `Writes:` lines own the
// fence, so `bench coverage --check` grades no leftover `## Ownership fences` section, and
// a section that omits the review pickup is not a fault.
func TestCommandCheckIgnoresAFenceSection(t *testing.T) {
	root, _ := citedSpec(t, "s", "cli seam", "\n## Ownership fences\n\n- `internal/x/`\n")
	t.Chdir(root)

	if out, code := Command([]string{"--check", "specs/s/spec.md"}); code != 0 {
		t.Fatalf("Command = (%d, %q), want exit 0 for a fence section without the review pickup", code, out)
	}
}
