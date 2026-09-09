package git

import (
	"os"
	"path/filepath"
	"testing"
)

// hostilePaths are the framing-sensitive names the index must carry whole. Every non-`-z`
// git form C-quotes each of them, so a reader that split on whitespace or trusted the
// quoting would lose or corrupt every one.
var hostilePaths = []string{"docs/a b.md", "docs/new\nline.md", "docs/q\"uote.md", "docs/tab\t.md"}

// initRepo makes an empty repository with a commit identity and no commit.
func initRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init", "-q", "-b", "main")
	runGit(t, root, "config", "user.email", "bench@example.invalid")
	runGit(t, root, "config", "user.name", "bench test")
	return root
}

// plant writes a fixture at the repository-relative path rel and makes its parents.
func plant(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("make directory for %q: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %q: %v", rel, err)
	}
}

func stagedPaths(index StagedIndex) []string {
	out := make([]string, 0, len(index.Staged))
	for _, entry := range index.Staged {
		out = append(out, entry.Path)
	}
	return out
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestReadStagedIndexFramesHostilePaths proves the NUL framing: a space, a newline, a
// double quote, and a tab in a staged path survive the listing whole.
func TestReadStagedIndexFramesHostilePaths(t *testing.T) {
	root := initRepo(t)
	runGit(t, root, "commit", "-q", "--allow-empty", "-m", "base")
	for _, rel := range hostilePaths {
		plant(t, root, rel, "body\n")
	}
	runGit(t, root, "add", "-A")

	index, err := ReadStagedIndex(root)
	if err != nil {
		t.Fatalf("ReadStagedIndex: %v", err)
	}
	if got := stagedPaths(index); !equalStrings(got, hostilePaths) {
		t.Fatalf("staged paths = %q, want %q", got, hostilePaths)
	}
	if got := len(index.Entries); got != len(hostilePaths) {
		t.Fatalf("index entries = %d, want %d", got, len(hostilePaths))
	}
}

// TestReadStagedIndexReportsModesAndOmitsDeletions proves the mode column and the
// `--diff-filter=ACMRT` selection: a symbolic link carries mode 120000, and a staged
// deletion is no entry at all.
func TestReadStagedIndexReportsModesAndOmitsDeletions(t *testing.T) {
	root := initRepo(t)
	plant(t, root, "docs/gone.md", "gone\n")
	plant(t, root, "docs/kept.md", "kept\n")
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "-q", "-m", "base")
	runGit(t, root, "rm", "-q", "docs/gone.md")
	if err := os.Symlink("kept.md", filepath.Join(root, "docs", "link.md")); err != nil {
		t.Fatalf("plant a symbolic link: %v", err)
	}
	runGit(t, root, "add", "-A")

	index, err := ReadStagedIndex(root)
	if err != nil {
		t.Fatalf("ReadStagedIndex: %v", err)
	}
	modes := map[string]string{}
	for _, entry := range index.Staged {
		modes[entry.Path] = entry.Mode
	}
	if _, present := modes["docs/gone.md"]; present {
		t.Fatalf("staged entries = %v, want no deletion", modes)
	}
	if modes["docs/link.md"] != "120000" {
		t.Fatalf("link mode = %q, want 120000", modes["docs/link.md"])
	}
	for _, entry := range index.Staged {
		if entry.Path == "docs/link.md" && entry.IsRegularFile() {
			t.Fatal("a symbolic-link entry reads as a regular file")
		}
	}
}

// TestReadStagedIndexStagesEveryEntryOnAnUnbornBranch proves the unborn rule: with no HEAD
// to diff against, every index entry is staged.
func TestReadStagedIndexStagesEveryEntryOnAnUnbornBranch(t *testing.T) {
	root := initRepo(t)
	plant(t, root, "docs/one.md", "one\n")
	plant(t, root, "docs/two.md", "two\n")
	runGit(t, root, "add", "-A")

	index, err := ReadStagedIndex(root)
	if err != nil {
		t.Fatalf("ReadStagedIndex: %v", err)
	}
	want := []string{"docs/one.md", "docs/two.md"}
	if got := stagedPaths(index); !equalStrings(got, want) {
		t.Fatalf("staged paths = %q, want %q", got, want)
	}
}

// TestReadStagedIndexRefusesAnyRootThatIsNotAWorkingTreeTop covers both refusals: a
// directory outside every repository, and a directory inside one that is not its top.
func TestReadStagedIndexRefusesAnyRootThatIsNotAWorkingTreeTop(t *testing.T) {
	if _, err := ReadStagedIndex(t.TempDir()); err == nil {
		t.Fatal("ReadStagedIndex accepted a plain directory")
	}
	root := initRepo(t)
	plant(t, root, "docs/one.md", "one\n")
	runGit(t, root, "add", "-A")
	if _, err := ReadStagedIndex(filepath.Join(root, "docs")); err == nil {
		t.Fatal("ReadStagedIndex accepted a subdirectory as the root")
	}
	if !IsWorkTreeTop(root) {
		t.Fatalf("IsWorkTreeTop(%q) = false, want true", root)
	}
}

// TestIndexBlobReadsTheIndexNotTheWorkingFile proves the blob read answers the staged
// bytes for an ordinary path and for a framing-sensitive one.
func TestIndexBlobReadsTheIndexNotTheWorkingFile(t *testing.T) {
	root := initRepo(t)
	const hostile = "docs/new\nline.md"
	plant(t, root, "docs/one.md", "staged\n")
	plant(t, root, hostile, "staged hostile\n")
	runGit(t, root, "add", "-A")
	plant(t, root, "docs/one.md", "working\n")
	plant(t, root, hostile, "working hostile\n")

	for _, want := range []struct{ path, body string }{
		{"docs/one.md", "staged\n"},
		{hostile, "staged hostile\n"},
	} {
		body, err := IndexBlob(root, want.path)
		if err != nil {
			t.Fatalf("IndexBlob %q: %v", want.path, err)
		}
		if string(body) != want.body {
			t.Errorf("IndexBlob %q = %q, want %q", want.path, body, want.body)
		}
	}
	if _, err := IndexBlob(root, "docs/absent.md"); err == nil {
		t.Fatal("IndexBlob accepted a path the index does not hold")
	}
}
