package gittest

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestKitCopyPreservesTheVisibleWorkingTree(t *testing.T) {
	root := RepoOnBranch(t, "main")
	visible := map[string]struct {
		body string
		mode os.FileMode
	}{
		".gitignore":                   {"dist/\nbin/bench-broker.manifest\n*.ignored\n", 0o644},
		"tracked":                      {"tracked bytes", 0o644},
		"tracked.ignored":              {"tracked despite the ignore pattern", 0o644},
		"bin/tool":                     {"\x00\xff\x01", 0o755},
		"untracked file\nwith newline": {"visible working bytes", 0o600},
	}
	for name, file := range visible {
		writeKitCopyFile(t, root, name, file.body, file.mode)
	}
	run(t, root, "add", "-f", ".gitignore", "tracked", "tracked.ignored", "bin/tool")
	run(t, root, "commit", "-qm", "fixture")
	ignored := []string{"dist/generated", "bin/bench-broker.manifest", "secret.ignored", ".git/operator-record"}
	for _, name := range ignored {
		writeKitCopyFile(t, root, name, "operator bytes", 0o600)
	}
	links := map[string]string{"link": "tracked", "dangling-link": "absent"}
	for name, target := range links {
		if err := os.Symlink(target, filepath.Join(root, name)); err != nil {
			t.Fatal(err)
		}
	}

	copyRoot := KitCopy(t, root)
	if copyRoot == root {
		t.Fatal("kit copy returned the source root")
	}
	for name, file := range visible {
		path := filepath.Join(copyRoot, name)
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, []byte(file.body)) {
			t.Fatalf("copy file %q = %q, %v", name, got, err)
		}
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != file.mode {
			t.Fatalf("copy file mode %q = %v, %v; want %v", name, info, err, file.mode)
		}
	}
	for _, name := range append(ignored, "dist") {
		if _, err := os.Lstat(filepath.Join(copyRoot, name)); !os.IsNotExist(err) {
			t.Fatalf("copy contains excluded path %q: %v", name, err)
		}
	}
	for name, want := range links {
		if got, err := os.Readlink(filepath.Join(copyRoot, name)); err != nil || got != want {
			t.Fatalf("copy link %q = %q, %v; want %q", name, got, err, want)
		}
	}
	if info, err := os.Stat(filepath.Join(copyRoot, ".git")); err != nil || !info.IsDir() {
		t.Fatalf("copy has no private git directory: %v, %v", info, err)
	}
	if got := output(t, copyRoot, "rev-list", "--count", "HEAD"); got != "1" {
		t.Fatalf("copy commits = %q, want one", got)
	}
	if got := output(t, copyRoot, "status", "--porcelain"); got != "" {
		t.Fatalf("copy has uncommitted files: %s", got)
	}
}

func writeKitCopyFile(t *testing.T, root, name, body string, mode os.FileMode) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), mode); err != nil {
		t.Fatal(err)
	}
}
