package gocache

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

func cacheProtocol() *protection.Protocol {
	return &protection.Protocol{Validators: protection.Validators{protection.KindCache: ValidateProtection}, ReadLimit: bounds.ControlRecordLimit}
}

// TestProtectionJoinsUnderTheSharedHold opens one generation per cache. A second in-process
// hold shares the descriptor and joins that generation.
func TestProtectionJoinsUnderTheSharedHold(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "go-build")
	first, err := HoldDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Release()
	second, err := HoldDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Release()
	p := cacheProtocol()
	d, err := first.Protect(p)
	canonical, _ := canonicalpath.Resolve(dir)
	described, _ := cacheDirectory.Describe(canonical)
	if err != nil || d.Store != described.Store || ValidateProtection(d) != nil {
		t.Fatalf("cache protection = %+v (%v)", d, err)
	}
	if joined, err := second.Protect(p); err != nil || joined != d {
		t.Fatalf("second hold opened %+v (%v), want %+v", joined, err, d)
	}
}

// TestProtectionRefusesAForeignCache grades each cache identity fact on its own.
func TestProtectionRefusesAForeignCache(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "go-build")
	holder, err := HoldDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer holder.Release()
	d, err := holder.Protect(cacheProtocol())
	if err != nil {
		t.Fatal(err)
	}
	parent := filepath.Dir(d.ResourcePath)
	other := filepath.Join(parent, "other")
	if err := os.Mkdir(other, 0o755); err != nil {
		t.Fatal(err)
	}
	link, parentLink := filepath.Join(t.TempDir(), "link"), filepath.Join(t.TempDir(), "parent")
	if err := errors.Join(os.Symlink(dir, link), os.Symlink(parent, parentLink)); err != nil {
		t.Fatal(err)
	}
	foreignStore, foreignBinding, unlocked, linked, linkedParent := d, d, d, d, d
	described, _ := cacheDirectory.Describe(other)
	foreignStore.Store = described.Store
	foreignBinding.Binding = described.Binding
	// The unlocked directory carries its own binding, so only the missing cache lock refuses it.
	unlocked.ResourcePath, unlocked.Store, unlocked.Binding = other, described.Store, foreignBinding.Binding
	linked.ResourcePath = link
	linked.Store, _ = cacheDirectory.Store(link, linked.Binding)
	// A link in a parent component reaches the real cache, so only the canonical-path check refuses it.
	linkedParent.ResourcePath = filepath.Join(parentLink, filepath.Base(d.ResourcePath))
	linkedParent.Store, _ = cacheDirectory.Store(linkedParent.ResourcePath, linkedParent.Binding)
	for name, foreign := range map[string]protection.Descriptor{"store": foreignStore, "binding": foreignBinding, "lock": unlocked, "link": linked, "parent link": linkedParent} {
		if err := ValidateProtection(foreign); err == nil {
			t.Errorf("%s: the cache validator accepted %+v", name, foreign)
		}
	}
	if err := holder.Release(); err != nil {
		t.Fatal(err)
	}
	if _, err := holder.Protect(cacheProtocol()); err == nil {
		t.Fatal("a released hold opened cache protection")
	}
}

// protectedHold holds dir shared and opens its generation under p.
func protectedHold(t *testing.T, dir string, p *protection.Protocol) (*Holder, protection.Descriptor) {
	t.Helper()
	holder, err := HoldDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = holder.Release() })
	d, err := holder.Protect(p)
	if err != nil {
		t.Fatal(err)
	}
	return holder, d
}

// completedUser registers one user of d and records its proved no-start.
func completedUser(t *testing.T, p *protection.Protocol, d protection.Descriptor) {
	t.Helper()
	run, err := p.Prepare([]protection.Descriptor{d})
	if err == nil {
		err = run.Finish(protection.Terminal{NoStart: true})
	}
	if err != nil {
		t.Fatal(err)
	}
}

// requireShards fails unless the fixture's shards are present exactly when want is true.
func requireShards(t *testing.T, fixture cacheFixture, want bool) {
	t.Helper()
	if got := slices.Contains(entries(t, fixture.dir), "aa"); got != want {
		t.Fatalf("shards present = %t, want %t; entries %q", got, want, entries(t, fixture.dir))
	}
}

// TestProtectionCompletedHold is PL57. The sole holder's release retires a completed
// generation. The retired store names no user, so a clean removes the shards, and the next
// hold opens a fresh generation.
func TestProtectionCompletedHold(t *testing.T) {
	fixture := newCacheFixture(t)
	p := cacheProtocol()
	holder, d := protectedHold(t, fixture.dir, p)
	completedUser(t, p, d)
	if err := holder.Release(); err != nil {
		t.Fatal(err)
	}
	if out, code := cleanWith(cleanEnv(fixture.home), p); code != 0 {
		t.Fatalf("clean of a cache with a retired store = (%q, %d), want exit 0", out, code)
	}
	requireShards(t, fixture, false)
	if _, next := protectedHold(t, fixture.dir, p); next.ResourceID == d.ResourceID {
		t.Fatalf("the next hold joined the released generation %s", d.ResourceID)
	}
}

// TestReleaseRunsTheRetiredCheck is RS28 and RS11. Only the last uncontended release of a
// cache hold runs the absent-resource check, once, and the hold's own store already holds no
// generation then. A release while another holder in this process or another process still
// holds the cache runs no check.
func TestReleaseRunsTheRetiredCheck(t *testing.T) {
	saved := discardRetiredCaches
	t.Cleanup(func() { discardRetiredCaches = saved })
	var store string
	var retired []bool
	// A protocol with no validator recovers no generation and changes no store, so the
	// recorder only observes whether the hold's generation retired.
	observer := &protection.Protocol{ReadLimit: bounds.ControlRecordLimit}
	discardRetiredCaches = func(int64) error {
		retired = append(retired, protection.NoGeneration(observer.Recover(store, nil)))
		return nil
	}
	fixture := newCacheFixture(t)
	p := cacheProtocol()

	stop := startHolderProcess(t, fixture.home)
	holder, d := protectedHold(t, fixture.dir, p)
	store = d.Store
	if err := holder.Release(); err != nil {
		t.Fatal(err)
	}
	stop()
	if len(retired) != 0 {
		t.Fatalf("a contended release ran the check %d times, want none", len(retired))
	}

	first, _ := protectedHold(t, fixture.dir, p)
	last, _ := protectedHold(t, fixture.dir, p)
	if err := first.Release(); err != nil {
		t.Fatal(err)
	}
	if len(retired) != 0 {
		t.Fatalf("a release beside another in-process holder ran the check %d times, want none", len(retired))
	}
	if err := last.Release(); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(retired, []bool{true}) {
		t.Fatalf("the last uncontended release ran the check with retired generations %v, want one call after the retirement", retired)
	}
}

// TestProtectionUnreadableProtection is PL56. A pending user keeps its sealed generation,
// and unreadable records leave the exclusive cleaner uncertain, so go clean never starts.
func TestProtectionUnreadableProtection(t *testing.T) {
	fixture := newCacheFixture(t)
	p := cacheProtocol()
	holder, d := protectedHold(t, fixture.dir, p)
	if _, err := p.Prepare([]protection.Descriptor{d}); err != nil {
		t.Fatal(err)
	}
	if err := holder.Release(); err != nil {
		t.Fatal(err)
	}
	err := filepath.WalkDir(d.Store, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		if info, err := entry.Info(); err != nil || info.Size() == 0 {
			return err
		}
		return os.WriteFile(path, []byte("{"), 0o600)
	})
	if err != nil {
		t.Fatal(err)
	}
	if out, code := cleanWith(cleanEnv(fixture.home), p); code != 1 || !strings.HasPrefix(out, "error: cache protected — ") {
		t.Fatalf("clean over unreadable protection = (%q, %d), want the protection refusal", out, code)
	}
	requireShards(t, fixture, true)
}

// TestProtectionAbsentCache is PL76. An installed protocol keeps the absent-cache answer, so
// the clean creates no directory and no store.
func TestProtectionAbsentCache(t *testing.T) {
	home := t.TempDir()
	if out, code := cleanWith(cleanEnv(home), cacheProtocol()); code != 0 {
		t.Fatalf("clean of an absent cache = (%q, %d), want exit 0", out, code)
	}
	if _, err := os.Lstat(cacheDir(home)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("cache after the clean = %v, want it absent", err)
	}
}

// TestProtectionEmptyCache is PL77. Another process holds the cache, so the release keeps
// the open generation; its only user is complete, so the clean still removes the shards.
func TestProtectionEmptyCache(t *testing.T) {
	fixture := newCacheFixture(t)
	stop := startHolderProcess(t, fixture.home)
	p := cacheProtocol()
	holder, d := protectedHold(t, fixture.dir, p)
	completedUser(t, p, d)
	if err := holder.Release(); err != nil {
		t.Fatal(err)
	}
	stop()
	if _, err := os.Lstat(d.Store); err != nil {
		t.Fatalf("protection store after a contended release = %v, want it present", err)
	}
	if out, code := cleanWith(cleanEnv(fixture.home), p); code != 0 {
		t.Fatalf("clean of an empty protected cache = (%q, %d), want exit 0", out, code)
	}
	requireShards(t, fixture, false)
}
