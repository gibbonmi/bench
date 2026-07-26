package contract

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

type FixtureOption func(*fixtureConfig)

type fixtureConfig struct {
	repo      bool
	spacePath bool
}

type Fixture struct {
	t    testing.TB
	Root string
	Env  map[string]string
}

func WithNoRepo() FixtureOption {
	return func(cfg *fixtureConfig) { cfg.repo = false }
}

func WithSpacePath() FixtureOption {
	return func(cfg *fixtureConfig) { cfg.spacePath = true }
}

func KitRoot(t testing.TB) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve kit root: runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func SubjectRoot(t testing.TB) string {
	t.Helper()
	if root := os.Getenv("BENCH_CONTRACT_ROOT"); root != "" {
		abs, err := filepath.Abs(root)
		if err != nil {
			t.Fatalf("resolve BENCH_CONTRACT_ROOT: %v", err)
		}
		return abs
	}
	return KitRoot(t)
}

func NewFixture(t testing.TB, opts ...FixtureOption) Fixture {
	t.Helper()
	cfg := fixtureConfig{repo: true}
	for _, opt := range opts {
		opt(&cfg)
	}
	root := t.TempDir()
	if cfg.spacePath {
		root = filepath.Join(root, "space dir")
		if err := os.MkdirAll(root, 0o755); err != nil {
			t.Fatalf("create space-path fixture: %v", err)
		}
	}
	f := Fixture{t: t, Root: root}
	f.Env = isolatedEnv(t, f.Root)
	if cfg.repo {
		f.Git("init", "-q")
	}
	return f
}

func NewFixtureAt(t testing.TB, root string, env map[string]string) Fixture {
	t.Helper()
	return Fixture{t: t, Root: root, Env: env}
}

func WriteExecutableAbs(t testing.TB, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o755); err != nil {
		t.Fatalf("write executable %s: %v", path, err)
	}
}

func WriteFileAbs(t testing.TB, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func ReadFileAbs(t testing.TB, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func Mkdir(t testing.TB, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func Remove(t testing.TB, path string) {
	t.Helper()
	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("remove %s: %v", path, err)
	}
}

func (f Fixture) WriteFile(path, contents string) {
	f.t.Helper()
	full := filepath.Join(f.Root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		f.t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
	}
	if err := os.WriteFile(full, []byte(contents), 0o644); err != nil {
		f.t.Fatalf("write %s: %v", path, err)
	}
}

func (f Fixture) WriteExecutable(path, contents string) {
	f.t.Helper()
	full := filepath.Join(f.Root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		f.t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
	}
	if err := os.WriteFile(full, []byte(contents), 0o755); err != nil {
		f.t.Fatalf("write executable %s: %v", path, err)
	}
}

// WriteFifo puts a FIFO where a control record would be. A FIFO with no writer never
// yields EOF, so a reader that opens the path before checking its type blocks forever:
// this is the fixture that turns that bug into an expired deadline rather than a hung
// suite. A filesystem without FIFOs skips through capability.Fifo rather than a bare
// t.Skip.
func (f Fixture) WriteFifo(path string) {
	f.t.Helper()
	full := filepath.Join(f.Root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		f.t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
	}
	if err := syscall.Mkfifo(full, 0o644); err != nil {
		capability.Capability(f.t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
	}
}

// WriteUnreadable strips every permission bit from path and proves the strip took, since
// root ignores the mode entirely and would otherwise read straight through the assertion.
// restore is the mode the cleanup puts back — 0o644 for a file, 0o755 for a directory,
// which must be re-entered before the fixture tree can be removed. A host that cannot
// deny itself a read skips through capability.Privilege.
func (f Fixture) WriteUnreadable(path string, restore os.FileMode) {
	f.t.Helper()
	full := filepath.Join(f.Root, filepath.FromSlash(path))
	f.t.Cleanup(func() { _ = os.Chmod(full, restore) })
	if err := os.Chmod(full, 0o000); err != nil {
		capability.Capability(f.t, capability.Privilege, "cannot strip permissions: "+err.Error())
	}
	if fh, err := os.Open(full); err == nil {
		fh.Close()
		capability.Capability(f.t, capability.Privilege, "mode 0o000 is still readable by this user")
	}
}

func (f Fixture) Exists(path string) bool {
	_, err := os.Stat(filepath.Join(f.Root, filepath.FromSlash(path)))
	return err == nil
}

func (f Fixture) ReadFile(path string) string {
	f.t.Helper()
	data, err := os.ReadFile(filepath.Join(f.Root, filepath.FromSlash(path)))
	if err != nil {
		f.t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func noteContractFailure(t testing.TB, msg string) {
	t.Helper()
	t.Cleanup(func() {
		if t.Failed() {
			t.Log(msg)
		}
	})
}

func NoteContractFailure(t testing.TB, msg string) {
	t.Helper()
	noteContractFailure(t, msg)
}

func RunParallel(t *testing.T, name string, fn func(*testing.T)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		fn(t)
	})
}

func skipIfSubjectBenchMissing(t testing.TB) {
	t.Helper()
	skipIfSubjectFileMissing(t, "bin/bench.sh")
}

func SkipIfSubjectBenchMissing(t testing.TB) {
	t.Helper()
	skipIfSubjectBenchMissing(t)
}

func skipIfSubjectFileMissing(t testing.TB, rel string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(SubjectRoot(t), filepath.FromSlash(rel))); err != nil {
		if os.IsNotExist(err) {
			capability.Environment(t, fmt.Sprintf("subject root has no %s", rel))
		}
		t.Fatalf("stat subject %s: %v", rel, err)
	}
}

func SkipIfSubjectFileMissing(t testing.TB, rel string) {
	t.Helper()
	skipIfSubjectFileMissing(t, rel)
}
