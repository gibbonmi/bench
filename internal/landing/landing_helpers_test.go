package landing

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	git(t, root, "init", "-q", "-b", "main")
	git(t, root, "config", "user.email", "a@b.c")
	git(t, root, "config", "user.name", "a")
	write(t, root, "named", "base")
	write(t, root, "foreign", "base")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "base")
	return root
}

func raceFixture(t *testing.T) string {
	t.Helper()
	root := fixture(t)
	write(t, root, "other", "base")
	write(t, root, ".gitignore", "ignored\n")
	git(t, root, "add", "other", ".gitignore")
	git(t, root, "commit", "-qm", "race fixture")
	return root
}

type pathSnapshot struct {
	Status, Index []byte
	Worktree      map[string][]byte
}

func dirtyUnnamedState(t *testing.T, root string) ([]string, pathSnapshot) {
	t.Helper()
	write(t, root, "foreign", "staged\n")
	git(t, root, "add", "foreign")
	write(t, root, "foreign", "staged-plus-unstaged\n")
	write(t, root, "other", "unstaged\n")
	write(t, root, "new", "untracked\n")
	write(t, root, "ignored", "ignored\n")
	paths := []string{"foreign", "other", "new", "ignored"}
	return paths, snapshotPaths(t, root, paths...)
}

func snapshotPaths(t *testing.T, root string, paths ...string) pathSnapshot {
	t.Helper()
	statusArgs := append([]string{"status", "--porcelain=v1", "--ignored", "--"}, paths...)
	indexArgs := append([]string{"ls-files", "--stage", "--"}, paths...)
	snapshot := pathSnapshot{
		Status:   gitBytes(t, root, statusArgs...),
		Index:    gitBytes(t, root, indexArgs...),
		Worktree: make(map[string][]byte, len(paths)),
	}
	for _, path := range paths {
		snapshot.Worktree[path] = mustRead(t, filepath.Join(root, path))
	}
	return snapshot
}

func write(t *testing.T, root, path, value string) {
	t.Helper()
	p := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(value), 0o644); err != nil {
		t.Fatal(err)
	}
}
func mustRead(t *testing.T, p string) []byte {
	t.Helper()
	b, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	return strings.TrimSpace(string(gitBytes(t, root, args...)))
}

func gitBytes(t *testing.T, root string, args ...string) []byte {
	t.Helper()
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	b, e := c.CombinedOutput()
	if e != nil {
		t.Fatalf("git %v: %v: %s", args, e, b)
	}
	return b
}

func gitMode(t *testing.T, root string, args ...string) string {
	t.Helper()
	fields := strings.Fields(git(t, root, args...))
	if len(fields) == 0 {
		t.Fatal("git mode output is empty")
	}
	return fields[0]
}

// The reviewed-landing fixtures below vary one axis: which side edited the spec.
// Every published-tree row shares this builder so a fixture change cannot make two
// rows disagree about what a reviewed landing composes.
const (
	reviewedBaseSpec        = "# X\n\nStatus: staged\n\n## Stories\n\nstory one\nstory two\n"
	reviewedAmendedSpec     = "# X\n\nStatus: staged\n\n## Stories\n\nstory one, amended by the review\nstory two\n"
	reviewedHeadingSpec     = "# X, moved on the destination\n\nStatus: staged\n\n## Stories\n\nstory one\nstory two\n"
	reviewedOverlapSpec     = "# X\n\nStatus: staged\n\n## Stories\n\nstory one, rewritten on the destination\nstory two\n"
	specProvenanceRefusal   = "staged spec bytes are not the reviewed source tip's committed spec"
	reviewedFixtureSpecPath = "specs/x/spec.md"
)

type reviewedLanding struct {
	root, sourceWorktree      string
	base, destination, source string
	specMode                  os.FileMode
}

// newReviewedLanding stages baseSpec at the shared review base, commits sourceSpec on
// the reviewed source branch, and advances the destination with destinationSpec. An
// empty spec string means that side wrote no spec at all. A zero specMode means the
// 0o644 default. The function applies a non-zero mode to the source's spec before the
// commit. The source tip and the request then agree on the mode the composition carries.
func newReviewedLanding(t *testing.T, baseSpec, destinationSpec, sourceSpec string, specMode os.FileMode) reviewedLanding {
	t.Helper()
	if specMode == 0 {
		specMode = 0o644
	}
	root := fixture(t)
	if baseSpec != "" {
		write(t, root, reviewedFixtureSpecPath, baseSpec)
		git(t, root, "add", reviewedFixtureSpecPath)
	}
	git(t, root, "commit", "--allow-empty", "-qm", "review base")
	base := git(t, root, "rev-parse", "HEAD")
	sourceWorktree := filepath.Join(t.TempDir(), "source")
	git(t, root, "worktree", "add", "-qb", "reviewed-source", sourceWorktree, base)
	write(t, sourceWorktree, "reviewed", "source bytes\n")
	if sourceSpec != "" {
		write(t, sourceWorktree, reviewedFixtureSpecPath, sourceSpec)
		if err := os.Chmod(filepath.Join(sourceWorktree, reviewedFixtureSpecPath), specMode); err != nil {
			t.Fatal(err)
		}
	}
	git(t, sourceWorktree, "add", "-A")
	git(t, sourceWorktree, "commit", "-qm", "reviewed source")
	destination := base
	if destinationSpec != "" {
		write(t, root, reviewedFixtureSpecPath, destinationSpec)
		git(t, root, "add", reviewedFixtureSpecPath)
		git(t, root, "commit", "-qm", "destination spec movement")
		destination = git(t, root, "rev-parse", "HEAD")
	}
	return reviewedLanding{
		root: root, sourceWorktree: sourceWorktree, base: base,
		destination: destination, source: git(t, sourceWorktree, "rev-parse", "HEAD"),
		specMode: specMode,
	}
}

func (f reviewedLanding) request(t *testing.T, specBytes, message string) ReviewedRequest {
	t.Helper()
	sourceFingerprint, err := CheckoutFingerprint(f.sourceWorktree)
	if err != nil {
		t.Fatal(err)
	}
	destinationFingerprint, err := CheckoutFingerprint(f.root)
	if err != nil {
		t.Fatal(err)
	}
	return ReviewedRequest{
		Root: f.root, Destination: "refs/heads/main", DestinationBase: f.destination,
		Source: "refs/heads/reviewed-source", SourceTip: f.source, ReviewBase: f.base,
		SourceWorktree: f.sourceWorktree, SourceFingerprint: sourceFingerprint,
		DestinationFingerprint: destinationFingerprint,
		SpecPath:               reviewedFixtureSpecPath, SpecBytes: []byte(specBytes), SpecMode: f.specMode,
		Message: message,
	}
}
