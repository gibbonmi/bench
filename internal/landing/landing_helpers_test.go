package landing

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate/authorization"
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

// The composition conflict setups below build one real Git conflict of a named kind on
// two branches, and each returns the destination commit and the source commit.

func changeBoth(path, destinationValue, sourceValue string) func(*testing.T, string, string) (string, string) {
	return func(t *testing.T, root, base string) (string, string) {
		return commitSides(t, root, base, func() { write(t, root, path, destinationValue) }, func() { write(t, root, path, sourceValue) })
	}
}

func modifyDelete(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() { write(t, root, "named", "destination") }, func() { git(t, root, "rm", "-q", "named") })
}

func renameRename(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() { git(t, root, "mv", "named", "destination-name") }, func() { git(t, root, "mv", "named", "source-name") })
}

func fileDirectory(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() { write(t, root, "clash", "file") }, func() { write(t, root, "clash/child", "child") })
}

func modeConflict(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() {
		write(t, root, "named", "destination")
		if err := os.Chmod(filepath.Join(root, "named"), 0o755); err != nil {
			t.Fatal(err)
		}
	}, func() { write(t, root, "named", "source") })
}

func symlinkConflict(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() {
		if err := os.Remove(filepath.Join(root, "named")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("destination-target", filepath.Join(root, "named")); err != nil {
			capability.Capability(t, capability.Symlink, err.Error())
		}
	}, func() { write(t, root, "named", "source") })
}

func gitlinkConflict(t *testing.T, root, base string) (string, string) {
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, nested, "init", "-q", "-b", "main")
	git(t, nested, "config", "user.email", "a@b.c")
	git(t, nested, "config", "user.name", "a")
	write(t, nested, "inside", "base")
	git(t, nested, "add", "inside")
	git(t, nested, "commit", "-qm", "base")
	git(t, root, "add", "nested")
	git(t, root, "commit", "-qm", "add gitlink")
	base = git(t, root, "rev-parse", "HEAD")
	git(t, nested, "checkout", "-qb", "destination")
	write(t, nested, "inside", "destination")
	git(t, nested, "commit", "-am", "destination", "-q")
	destinationNested := git(t, nested, "rev-parse", "HEAD")
	git(t, nested, "checkout", "-qb", "source", baseForNested(t, nested))
	write(t, nested, "inside", "source")
	git(t, nested, "commit", "-am", "source", "-q")
	sourceNested := git(t, nested, "rev-parse", "HEAD")
	return commitSides(t, root, base, func() { git(t, nested, "checkout", "-q", destinationNested); git(t, root, "add", "nested") }, func() { git(t, nested, "checkout", "-q", sourceNested); git(t, root, "add", "nested") })
}

func baseForNested(t *testing.T, root string) string { return git(t, root, "rev-parse", "main") }

func commitSides(t *testing.T, root, base string, destinationChange, sourceChange func()) (string, string) {
	git(t, root, "checkout", "-qb", "destination", base)
	destinationChange()
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "destination")
	destination := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-qb", "source", base)
	sourceChange()
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "source")
	return destination, git(t, root, "rev-parse", "HEAD")
}

type compositionSnapshot struct{ refs, index, status, worktree, mergeHead string }

func compositionState(t *testing.T, root string) compositionSnapshot {
	t.Helper()
	mergeHead, err := os.ReadFile(filepath.Join(root, ".git", "MERGE_HEAD"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return compositionSnapshot{
		refs:      string(gitBytes(t, root, "for-each-ref", "--format=%(refname) %(objectname)")),
		index:     string(gitBytes(t, root, "ls-files", "--stage", "-z")),
		status:    string(gitBytes(t, root, "status", "--porcelain=v2", "--untracked-files=all")),
		worktree:  string(gitBytes(t, root, "diff", "--binary", "HEAD")),
		mergeHead: string(mergeHead),
	}
}

func compositionResult(paths, names []string, content map[string]string) compositionExpectation {
	return compositionExpectation{paths: paths, wantNames: names, wantContent: content}
}

func greenOwner() Owner {
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	return o
}

func assertPublishedTree(t *testing.T, root string, got Result, wantNames []string, wantContent, wantMode map[string]string) {
	t.Helper()
	if got.Tree != git(t, root, "rev-parse", got.Commit+"^{tree}") || git(t, root, "rev-parse", "HEAD") != got.Commit {
		t.Fatal("authorized tree was not the exact published tree")
	}
	rawNames := git(t, root, "ls-tree", "-rz", "--name-only", got.Commit)
	names := strings.Split(rawNames, "\x00")
	if names[len(names)-1] == "" {
		names = names[:len(names)-1]
	}
	sort.Strings(wantNames)
	if !reflect.DeepEqual(names, wantNames) {
		t.Fatalf("published paths = %q, want %q", names, wantNames)
	}
	for path, content := range wantContent {
		if got := git(t, root, "show", got.Commit+":"+path); got != content {
			t.Fatalf("published %s = %q, want %q", path, got, content)
		}
	}
	for path, mode := range wantMode {
		fields := strings.Fields(git(t, root, "ls-tree", got.Commit, "--", path))
		if len(fields) < 4 || fields[0] != mode {
			t.Fatalf("published mode %s = %q, want %s", path, fields, mode)
		}
	}
}
