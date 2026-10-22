package probe

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// goldenParent holds the clamp golden of this package run, and TestMain removes it after
// the run. The golden build creates it, not TestMain, because each helper process that
// re-executes this test binary runs TestMain too, and a test can kill a helper before its
// cleanup.
var goldenParent string

// clampBuild is one built clamp repository: the directory that holds it, the Git checkout
// root, and the sealed run binary.
type clampBuild struct {
	dir, root, executable string
}

// clampGolden is the clamp build of this package run. The first test that asks for it builds
// it, and each test takes a private copy. A copy costs less than a build, because a build
// digests the source closure to publish a sealed run binary.
var clampGolden struct {
	once  sync.Once
	build clampBuild
}

// goldenPublishes counts the calls of the fixture publisher, so a test observes the builds.
var goldenPublishes atomic.Int64

func countedPublish(root, staged, executable, manifestDir, version string) error {
	goldenPublishes.Add(1)
	return freshness.Publish(root, staged, executable, manifestDir, version)
}

func TestMain(m *testing.M) {
	code := m.Run()
	if goldenParent != "" {
		if err := os.RemoveAll(goldenParent); err != nil {
			fmt.Fprintln(os.Stderr, err)
			code = 1
		}
	}
	os.Exit(code)
}

// buildClamp builds the clamp repository under dir: a Git checkout in repo, and the sealed run
// binary with its seal and broker manifest in bin.
func buildClamp(t *testing.T, dir string) clampBuild {
	t.Helper()
	root, bin := filepath.Join(dir, "repo"), filepath.Join(dir, "bin")
	for _, path := range []string{root, bin} {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	gittest.Output(t, root, "init", "-q")
	return clampBuild{dir: dir, root: root, executable: probetest.Build(t, root, bin, countedPublish)}
}

// golden builds the clamp golden in the first test that asks for it and answers it.
func golden(t *testing.T) clampBuild {
	t.Helper()
	clampGolden.once.Do(func() {
		parent, err := os.MkdirTemp("", "probe-golden-")
		if err != nil {
			t.Fatal(err)
		}
		goldenParent = parent
		dir, err := filepath.EvalSymlinks(parent)
		if err != nil {
			t.Fatal(err)
		}
		clampGolden.build = buildClamp(t, dir)
	})
	if clampGolden.build.dir == "" {
		t.Fatal("the clamp golden failed to build in an earlier test")
	}
	return clampGolden.build
}

// clampCopy copies the clamp golden into a private directory and binds the copy to the test.
// It answers the copy's root and the home. The checkout, the executable, its seal, and the
// manifest directory move together, because the seal holds digests and no path. The copy's
// broker manifest still names the golden executable, and no probe reader reads it.
func clampCopy(t *testing.T) (string, string) {
	t.Helper()
	base := golden(t)
	dir := scratchDir(t)
	if err := testrepo.CopyTree(dir, base.dir); err != nil {
		t.Fatal(err)
	}
	copied := func(path string) string {
		t.Helper()
		rel, err := filepath.Rel(base.dir, path)
		if err != nil {
			t.Fatal(err)
		}
		return filepath.Join(dir, rel)
	}
	root := copied(base.root)
	return root, probetest.Bind(t, root, copied(base.executable))
}

func TestClampGoldenBuildsOnce(t *testing.T) {
	clampFixture(t)
	clampFixture(t)
	if got := goldenPublishes.Load(); got != 1 {
		t.Fatalf("the fixture publisher ran %d times, want 1", got)
	}
}

func TestClampCopiesAreIsolated(t *testing.T) {
	first, second := clampFixture(t), clampFixture(t)
	writeFixtureFile(t, first.subject, "package probefixture\n", 0o644)
	requireSubjectBytes(t, second, probetest.ClampSource)
}

func TestClampCopyKeepsAFreshSeal(t *testing.T) {
	f := clampFixture(t)
	if _, err := runbinary.Inherit(f.root); err != nil {
		t.Fatalf("the copy's run binary refuses its own root: %v", err)
	}
}
