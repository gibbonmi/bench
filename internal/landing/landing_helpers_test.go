package landing

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/gittest"
)

func fixture(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	write(t, root, "named", "base")
	write(t, root, "foreign", "base")
	gittest.MustRun(t, root, "add", ".")
	gittest.MustRun(t, root, "commit", "-qm", "base")
	return root
}

func raceFixture(t *testing.T) string {
	t.Helper()
	root := fixture(t)
	write(t, root, "other", "base")
	write(t, root, ".gitignore", "ignored\n")
	gittest.MustRun(t, root, "add", "other", ".gitignore")
	gittest.MustRun(t, root, "commit", "-qm", "race fixture")
	return root
}

type pathSnapshot struct {
	Status, Index []byte
	Worktree      map[string][]byte
}

func dirtyUnnamedState(t *testing.T, root string) ([]string, pathSnapshot) {
	t.Helper()
	write(t, root, "foreign", "staged\n")
	gittest.MustRun(t, root, "add", "foreign")
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
		Status:   gittest.OutputBytes(t, root, statusArgs...),
		Index:    gittest.OutputBytes(t, root, indexArgs...),
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
func gitMode(t *testing.T, root string, args ...string) string {
	t.Helper()
	fields := strings.Fields(gittest.Output(t, root, args...))
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
		gittest.MustRun(t, root, "add", reviewedFixtureSpecPath)
	}
	gittest.MustRun(t, root, "commit", "--allow-empty", "-qm", "review base")
	base := gittest.Output(t, root, "rev-parse", "HEAD")
	sourceWorktree := filepath.Join(t.TempDir(), "source")
	gittest.MustRun(t, root, "worktree", "add", "-qb", "reviewed-source", sourceWorktree, base)
	write(t, sourceWorktree, "reviewed", "source bytes\n")
	if sourceSpec != "" {
		write(t, sourceWorktree, reviewedFixtureSpecPath, sourceSpec)
		if err := os.Chmod(filepath.Join(sourceWorktree, reviewedFixtureSpecPath), specMode); err != nil {
			t.Fatal(err)
		}
	}
	gittest.MustRun(t, sourceWorktree, "add", "-A")
	gittest.MustRun(t, sourceWorktree, "commit", "-qm", "reviewed source")
	destination := base
	if destinationSpec != "" {
		write(t, root, reviewedFixtureSpecPath, destinationSpec)
		gittest.MustRun(t, root, "add", reviewedFixtureSpecPath)
		gittest.MustRun(t, root, "commit", "-qm", "destination spec movement")
		destination = gittest.Output(t, root, "rev-parse", "HEAD")
	}
	return reviewedLanding{
		root: root, sourceWorktree: sourceWorktree, base: base,
		destination: destination, source: gittest.Output(t, sourceWorktree, "rev-parse", "HEAD"),
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
	return commitSides(t, root, base, func() { write(t, root, "named", "destination") }, func() { gittest.MustRun(t, root, "rm", "-q", "named") })
}

func renameRename(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() { gittest.MustRun(t, root, "mv", "named", "destination-name") }, func() { gittest.MustRun(t, root, "mv", "named", "source-name") })
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
			capability.Unavailable(t, capability.Symlink, err.Error())
		}
	}, func() { write(t, root, "named", "source") })
}

func gitlinkConflict(t *testing.T, root, base string) (string, string) {
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	gittest.MustRun(t, nested, "init", "-q", "-b", "main")
	for _, args := range gittest.IdentityConfig() {
		gittest.MustRun(t, nested, args...)
	}
	write(t, nested, "inside", "base")
	gittest.MustRun(t, nested, "add", "inside")
	gittest.MustRun(t, nested, "commit", "-qm", "base")
	gittest.MustRun(t, root, "add", "nested")
	gittest.MustRun(t, root, "commit", "-qm", "add gitlink")
	base = gittest.Output(t, root, "rev-parse", "HEAD")
	gittest.MustRun(t, nested, "checkout", "-qb", "destination")
	write(t, nested, "inside", "destination")
	gittest.MustRun(t, nested, "commit", "-am", "destination", "-q")
	destinationNested := gittest.Output(t, nested, "rev-parse", "HEAD")
	gittest.MustRun(t, nested, "checkout", "-qb", "source", baseForNested(t, nested))
	write(t, nested, "inside", "source")
	gittest.MustRun(t, nested, "commit", "-am", "source", "-q")
	sourceNested := gittest.Output(t, nested, "rev-parse", "HEAD")
	return commitSides(t, root, base, func() {
		gittest.MustRun(t, nested, "checkout", "-q", destinationNested)
		gittest.MustRun(t, root, "add", "nested")
	}, func() {
		gittest.MustRun(t, nested, "checkout", "-q", sourceNested)
		gittest.MustRun(t, root, "add", "nested")
	})
}

func baseForNested(t *testing.T, root string) string {
	return gittest.Output(t, root, "rev-parse", "main")
}

// compositionExpectation is what one composition row publishes, and what it leaves in the
// index and the working tree.
type compositionExpectation struct {
	paths        []string
	wantNames    []string
	wantContent  map[string]string
	wantMode     map[string]string
	wantIndex    map[string]string
	wantWorktree map[string]string
}

func commitSides(t *testing.T, root, base string, destinationChange, sourceChange func()) (string, string) {
	gittest.MustRun(t, root, "checkout", "-qb", "destination", base)
	destinationChange()
	gittest.MustRun(t, root, "add", "-A")
	gittest.MustRun(t, root, "commit", "-qm", "destination")
	destination := gittest.Output(t, root, "rev-parse", "HEAD")
	gittest.MustRun(t, root, "checkout", "-qb", "source", base)
	sourceChange()
	gittest.MustRun(t, root, "add", "-A")
	gittest.MustRun(t, root, "commit", "-qm", "source")
	return destination, gittest.Output(t, root, "rev-parse", "HEAD")
}

// deletedSide is the settle-row side value that removes the row's path. No side writes
// a NUL byte, so the sentinel cannot collide with content a row means to publish.
const deletedSide = "\x00deleted"

// settleSide builds one side of a settle row: deletedSide removes the path, and every
// other value writes itself there.
func settleSide(t *testing.T, root, path, value string) func() {
	return func() {
		if value == deletedSide {
			if err := os.Remove(filepath.Join(root, filepath.FromSlash(path))); err != nil {
				t.Fatal(err)
			}
			return
		}
		write(t, root, path, value)
	}
}

type compositionSnapshot struct{ refs, index, status, worktree, mergeHead string }

func compositionState(t *testing.T, root string) compositionSnapshot {
	t.Helper()
	mergeHead, err := os.ReadFile(mergeStatePath(t, root, "MERGE_HEAD"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return compositionSnapshot{
		refs:      string(gittest.OutputBytes(t, root, "for-each-ref", "--format=%(refname) %(objectname)")),
		index:     string(gittest.OutputBytes(t, root, "ls-files", "--stage", "-z")),
		status:    string(gittest.OutputBytes(t, root, "status", "--porcelain=v2", "--untracked-files=all")),
		worktree:  string(gittest.OutputBytes(t, root, "diff", "--binary", "HEAD")),
		mergeHead: string(mergeHead),
	}
}

func compositionResult(paths, names []string, content map[string]string) compositionExpectation {
	return compositionExpectation{paths: paths, wantNames: names, wantContent: content}
}

// assertCompositionRow observes what one composition row left outside the published tree:
// the untouched foreign index and worktree bytes, and each row's own index and worktree
// expectations. A gitlink publishes its commit, so that row reads the tree entry.
func assertCompositionRow(t *testing.T, root string, got Result, want compositionExpectation) {
	t.Helper()
	if have := gittest.Output(t, root, "show", ":foreign"); have != "foreign-index" {
		t.Fatalf("foreign index = %q, want foreign-index", have)
	}
	if have := string(mustRead(t, filepath.Join(root, "foreign"))); have != "foreign-worktree" {
		t.Fatalf("foreign worktree = %q, want foreign-worktree", have)
	}
	for path, content := range want.wantIndex {
		if mode := want.wantMode[path]; mode == "160000" {
			if have := strings.Fields(gittest.Output(t, root, "ls-tree", got.Commit, "--", path))[2]; have != content {
				t.Fatalf("published gitlink %s = %q, want %q", path, have, content)
			}
			continue
		}
		if have := gittest.Output(t, root, "show", ":"+path); have != content {
			t.Fatalf("index %s = %q, want %q", path, have, content)
		}
	}
	for path, content := range want.wantWorktree {
		if have := string(mustRead(t, filepath.Join(root, filepath.FromSlash(path)))); have != content {
			t.Fatalf("worktree %s = %q, want %q", path, have, content)
		}
	}
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
	if got.Tree != gittest.Output(t, root, "rev-parse", got.Commit+"^{tree}") || gittest.Output(t, root, "rev-parse", "HEAD") != got.Commit {
		t.Fatal("authorized tree was not the exact published tree")
	}
	rawNames := gittest.Output(t, root, "ls-tree", "-rz", "--name-only", got.Commit)
	names := strings.Split(rawNames, "\x00")
	if names[len(names)-1] == "" {
		names = names[:len(names)-1]
	}
	sort.Strings(wantNames)
	if !reflect.DeepEqual(names, wantNames) {
		t.Fatalf("published paths = %q, want %q", names, wantNames)
	}
	for path, content := range wantContent {
		if got := gittest.Output(t, root, "show", got.Commit+":"+path); got != content {
			t.Fatalf("published %s = %q, want %q", path, got, content)
		}
	}
	for path, mode := range wantMode {
		fields := strings.Fields(gittest.Output(t, root, "ls-tree", got.Commit, "--", path))
		if len(fields) < 4 || fields[0] != mode {
			t.Fatalf("published mode %s = %q, want %s", path, fields, mode)
		}
	}
}
