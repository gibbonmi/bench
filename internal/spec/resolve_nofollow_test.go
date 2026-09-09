package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveNoFollowRefusesLinkAndPreservesDefault(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.md")
	if err := os.WriteFile(target, []byte("Status: staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "linked.md")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	if content, _, _, ok, err := Resolve(root, link); err != nil || !ok || string(content) != "Status: staged\n" {
		t.Fatalf("Resolve() = %q, ok=%v, err=%v", content, ok, err)
	}
	_, _, _, ok, err := ResolveNoFollow(root, link)
	if ok || err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("ResolveNoFollow() ok=%v, err=%v", ok, err)
	}
}
