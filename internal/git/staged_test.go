package git

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
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

// TestIsWorkTreeTopAcceptsARelativeRoot proves the comparison absolutizes both operands:
// `.` inside the top names the top, and `.` inside a subdirectory still does not. Git
// answers an absolute top, so a comparison against the operand as spelled refuses every
// relative root.
func TestIsWorkTreeTopAcceptsARelativeRoot(t *testing.T) {
	root := initRepo(t)
	plant(t, root, "docs/one.md", "one\n")
	runGit(t, root, "add", "-A")

	t.Run("the top itself", func(t *testing.T) {
		t.Chdir(root)
		if !IsWorkTreeTop(".") {
			t.Fatal(`IsWorkTreeTop(".") = false at the top, want true`)
		}
		index, err := ReadStagedIndex(".")
		if err != nil {
			t.Fatalf(`ReadStagedIndex(".") at the top: %v`, err)
		}
		if !equalStrings(stagedPaths(index), []string{"docs/one.md"}) {
			t.Fatalf("staged paths = %q, want the one staged file", stagedPaths(index))
		}
	})

	t.Run("a subdirectory", func(t *testing.T) {
		t.Chdir(filepath.Join(root, "docs"))
		if IsWorkTreeTop(".") {
			t.Fatal(`IsWorkTreeTop(".") = true inside a subdirectory, want false`)
		}
	})
}

// indexBlobBytes drains one index blob through the stream form and closes it. A blob under
// the caller's bound reaches memory whole, so an assertion over its bytes reads as it did
// before the stream form.
func indexBlobBytes(t *testing.T, root, path string) ([]byte, error) {
	t.Helper()
	stream, err := IndexBlobReader(root, path)
	if err != nil {
		return nil, err
	}
	body, readErr := io.ReadAll(stream)
	if closeErr := stream.Close(); closeErr != nil {
		return nil, closeErr
	}
	return body, readErr
}

// TestIndexBlobReadsTheIndexNotTheWorkingFile proves the blob read answers the staged
// bytes for an ordinary path and for a framing-sensitive one, and that a path the index
// does not hold refuses under the `git show :<path> in <root>` text every caller reports.
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
		body, err := indexBlobBytes(t, root, want.path)
		if err != nil {
			t.Fatalf("IndexBlobReader %q: %v", want.path, err)
		}
		if string(body) != want.body {
			t.Errorf("IndexBlobReader %q = %q, want %q", want.path, body, want.body)
		}
	}
	_, err := indexBlobBytes(t, root, "docs/absent.md")
	if err == nil {
		t.Fatal("IndexBlobReader accepted a path the index does not hold")
	}
	if want := "git show :docs/absent.md in " + root + ": "; !strings.HasPrefix(err.Error(), want) {
		t.Fatalf("missing-path error = %q, want the %q prefix", err, want)
	}
}

// TestIndexBlobReaderBoundsAnOversizedBlob proves the read stops at the caller's bound: a
// blob three times the control-record limit answers the oversized verdict at the limit,
// and the close returns rather than leaving the child git writing into the pipe. A form
// that buffered the whole blob first would hold every byte before the bound applied.
func TestIndexBlobReaderBoundsAnOversizedBlob(t *testing.T) {
	root := initRepo(t)
	plant(t, root, "docs/big.md", strings.Repeat("a", int(bounds.ControlRecordLimit*3)))
	runGit(t, root, "add", "-A")

	stream, err := IndexBlobReader(root, "docs/big.md")
	if err != nil {
		t.Fatalf("IndexBlobReader: %v", err)
	}
	read := bounds.Read(stream, bounds.ControlRecordLimit)
	if read.Status != bounds.ReadOversized {
		t.Fatalf("read status = %q, want %q", read.Status, bounds.ReadOversized)
	}
	if got := int64(len(read.Data)); got != bounds.ControlRecordLimit {
		t.Fatalf("bytes kept = %d, want the limit %d", got, bounds.ControlRecordLimit)
	}

	closed := make(chan error, 1)
	go func() { closed <- stream.Close() }()
	window := bounds.TestDeadline(0)
	select {
	case err := <-closed:
		if err != nil {
			t.Fatalf("close after a bounded read: %v", err)
		}
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("the bounded blob read to reap its git child", window))
	}
}
