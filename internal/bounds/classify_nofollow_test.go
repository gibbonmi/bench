package bounds

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClassifyDirNoFollowRefusesLinkAndPreservesDefault(t *testing.T) {
	target := filepath.Join(t.TempDir(), "target")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "linked")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	if got := ClassifyDir(link); got.State != StateEmpty {
		t.Fatalf("ClassifyDir() state = %q, want %q: %s", got.State, StateEmpty, got.Reason)
	}
	got := ClassifyDirNoFollow(link)
	if got.State != StateWrongType {
		t.Fatalf("ClassifyDirNoFollow() state = %q, want %q: %s", got.State, StateWrongType,
			got.Reason)
	}
}
