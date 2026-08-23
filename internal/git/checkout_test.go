package git

import (
	"path/filepath"
	"testing"
)

func TestIsPrimaryCheckoutSeparatesLinkedWorktree(t *testing.T) {
	root := newRepo(t)
	linked := filepath.Join(t.TempDir(), "linked")
	runGit(t, root, "worktree", "add", "-q", "-b", "linked", linked)

	primary, err := IsPrimaryCheckout(root)
	if err != nil || !primary {
		t.Fatalf("IsPrimaryCheckout(primary) = %t, %v; want true, nil", primary, err)
	}
	primary, err = IsPrimaryCheckout(linked)
	if err != nil || primary {
		t.Fatalf("IsPrimaryCheckout(linked) = %t, %v; want false, nil", primary, err)
	}
}
