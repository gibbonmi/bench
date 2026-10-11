package adopt

import (
	"os"
	"path/filepath"
	"testing"
)

// TestKitClaudeMDMirrorsTheManagedForm keeps the kit's root CLAUDE.md byte-equal to
// benchClaudeMD. The kit is the source that every linked repo receives, so its own
// CLAUDE.md is the managed form and not a hand-kept second copy of it.
func TestKitClaudeMDMirrorsTheManagedForm(t *testing.T) {
	kit, err := os.ReadFile(filepath.Join("..", "..", "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(kit), benchClaudeMD(); got != want {
		t.Fatalf("root CLAUDE.md has drifted from benchClaudeMD()\nCLAUDE.md = %q\nbenchClaudeMD = %q", got, want)
	}
}
