package testreport

import (
	"strings"
	"testing"
)

// A fixture name in two families is an inventory error, not an empty inventory, so both
// inventory faces refuse it with the inventory diagnostic.
func TestInventoryFacesRefuseDuplicateFixtureName(t *testing.T) {
	root := t.TempDir()
	plantFixture(t, root, "package-core-guard", "twin", "")
	plantFixture(t, root, "line-routing", "twin", "")
	const diagnostic = `"twin" appears in multiple families`
	for _, args := range [][]string{
		{"--checks"},
		{"--check", "package-core-guard", "--fixtures"},
	} {
		output, code := inventoryFace(t, root, args...)
		if code != 1 || !strings.Contains(output, diagnostic) {
			t.Errorf("Command(%q) = %d, %q; want 1 and %s", args, code, output, diagnostic)
		}
	}
}
