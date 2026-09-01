package canary

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

func TestRestoreMutationFixtureReinstatesBaseAndRemovesOverlay(t *testing.T) {
	source := t.TempDir()
	if err := os.WriteFile(filepath.Join(source, "owned.txt"), []byte("clean\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fixture := t.TempDir()
	files := filepath.Join(fixture, filesDirName)
	if err := os.MkdirAll(files, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"owned.txt": "overlay\n", "added.txt": "added\n"} {
		if err := os.WriteFile(filepath.Join(files, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(fixture, "BASE"), []byte("owned.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := t.TempDir()
	if err := materializeMutationFixture(source, fixture, dst); err != nil {
		t.Fatalf("materialize fixture: %v", err)
	}
	if err := RestoreMutationFixture(source, fixture, dst); err != nil {
		t.Fatalf("restore fixture: %v", err)
	}
	if got, err := os.ReadFile(filepath.Join(dst, "owned.txt")); err != nil || string(got) != "clean\n" {
		t.Fatalf("restored base = %q, %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(dst, "added.txt")); !os.IsNotExist(err) {
		t.Fatalf("overlay remains after restore: %v", err)
	}
}

func TestRestoreMutationFixtureUsesSourceForOverlayPaths(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ fixture, rel string }{
		{filepath.Join(root, "tests", "canary", "load-validity-metadata", "codex-hooks-broken"), filepath.Join(".codex", "hooks.json")},
	} {
		dst := t.TempDir()
		if err := MaterializeMutationFixture(root, test.fixture, dst); err != nil {
			t.Fatal(err)
		}
		if err := RestoreMutationFixture(root, test.fixture, dst); err != nil {
			t.Fatal(err)
		}
		want, sourceErr := os.ReadFile(filepath.Join(root, test.rel))
		got, restoredErr := os.ReadFile(filepath.Join(dst, test.rel))
		if sourceErr != nil || restoredErr != nil || string(got) != string(want) {
			t.Fatalf("%s restored %s differs from source: source=%v restored=%v", test.fixture, test.rel, sourceErr, restoredErr)
		}
	}
}

// TestRestoreMutationFixtureRefusesDestinationEqualToRoot pins the refusal that keeps a
// restore off the real tree. A dst equal to root sends the overlay walk at the source, and
// that walk removes a directory the overlay names, so the refusal fires before the walk.
func TestRestoreMutationFixtureRefusesDestinationEqualToRoot(t *testing.T) {
	root, fixture := refusalSubject(t)
	// Each spelling resolves to root, so the refusal reads resolved absolute paths rather
	// than the caller's string.
	for _, dst := range []string{root, root + string(filepath.Separator), filepath.Join(root, "child", "..")} {
		assertRefusedAndRootIntact(t, root, fixture, dst)
	}
}

// TestRestoreMutationFixtureRefusesSymlinkedRootSpelling pins the refusal against a root
// reached through a symbolic link. The link and the real directory hold one tree, so a
// comparison of unresolved paths would let the walk run against the source.
func TestRestoreMutationFixtureRefusesSymlinkedRootSpelling(t *testing.T) {
	root, fixture := refusalSubject(t)
	link := filepath.Join(t.TempDir(), "root-link")
	if err := os.Symlink(root, link); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symbolic links unavailable: %v", err))
	}
	assertRefusedAndRootIntact(t, link, fixture, root)
	assertRefusedAndRootIntact(t, root, fixture, link)
}

// refusalSubject builds a source tree and a fixture whose overlay names a directory the
// source owns. A restore into the source removes that directory, so the directory reports
// whether the refusal ran before the walk.
func refusalSubject(t *testing.T) (root, fixture string) {
	t.Helper()
	root = t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "owned.txt"), []byte("clean\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "scratch"), 0o755); err != nil {
		t.Fatal(err)
	}
	fixture = t.TempDir()
	files := filepath.Join(fixture, filesDirName)
	if err := os.MkdirAll(filepath.Join(files, "scratch"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"owned.txt": "overlay\n", filepath.Join("scratch", "leftover.txt"): "added\n"} {
		if err := os.WriteFile(filepath.Join(files, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(fixture, "BASE"), []byte("owned.txt\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, fixture
}

// assertRefusedAndRootIntact requires the refusal and then requires that root keeps the
// directory and the file the walk would have changed.
func assertRefusedAndRootIntact(t *testing.T, root, fixture, dst string) {
	t.Helper()
	err := RestoreMutationFixture(root, fixture, dst)
	// The tree reads first. A restore that ran leaves its damage in root, and that damage
	// is the property under test.
	info, statErr := os.Stat(filepath.Join(root, "scratch"))
	if statErr != nil || !info.IsDir() {
		t.Fatalf("restore with root %q and dst %q removed the scratch directory: %v", root, dst, statErr)
	}
	if got, readErr := os.ReadFile(filepath.Join(root, "owned.txt")); readErr != nil || string(got) != "clean\n" {
		t.Fatalf("restore with root %q and dst %q rewrote owned.txt: %q, %v", root, dst, got, readErr)
	}
	if err == nil {
		t.Fatalf("restore with root %q and dst %q returned no error", root, dst)
	}
	if !strings.Contains(err.Error(), "RestoreMutationFixture refuses dst == root") {
		t.Fatalf("restore with root %q and dst %q error = %v, want the dst == root refusal", root, dst, err)
	}
}

// TestMaterializeMutationFixtureAnchorRefusals pins the anchor refusal messages. The
// anchor evaluator matches under collapsed whitespace, so a needle that wraps across a
// line in the target passes the evaluator and fails the byte-exact materializer. That
// case names the wrap; the other misses keep the plain message.
func TestMaterializeMutationFixtureAnchorRefusals(t *testing.T) {
	for _, test := range []struct {
		name     string
		body     string
		old      string
		wantWrap bool
	}{
		{"wrapped anchor names the wrap", "alpha\nbeta gamma\n", "alpha beta", true},
		{"absent anchor keeps the plain message", "alpha beta\n", "delta", false},
		{"repeated anchor keeps the plain message", "alpha\nalpha\n", "alpha", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := t.TempDir()
			files := filepath.Join(fixture, filesDirName)
			if err := os.MkdirAll(files, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(files, "target.txt"), []byte(test.body), 0o644); err != nil {
				t.Fatal(err)
			}
			mutate := fmt.Sprintf("[{\"path\": \"target.txt\", \"old\": %q, \"new\": \"replaced\"}]", test.old)
			if err := os.WriteFile(filepath.Join(fixture, "MUTATE.json"), []byte(mutate), 0o644); err != nil {
				t.Fatal(err)
			}
			err := materializeMutationFixture(t.TempDir(), fixture, t.TempDir())
			if err == nil {
				t.Fatal("materialize accepted an anchor that does not occur exactly once")
			}
			if !strings.Contains(err.Error(), "did not occur exactly once") {
				t.Fatalf("refusal = %v, want the exactly-once message", err)
			}
			if got := strings.Contains(err.Error(), "line wrap"); got != test.wantWrap {
				t.Fatalf("refusal = %v, wrap hint = %v, want %v", err, got, test.wantWrap)
			}
		})
	}
}
