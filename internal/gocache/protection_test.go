package gocache

import (
	"os"
	"path/filepath"
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
	if err != nil || d.Store != filepath.Join(canonical, ProtectionStore) || ValidateProtection(d) != nil {
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
	other := filepath.Join(t.TempDir(), "go-build")
	if err := os.Mkdir(other, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(dir, link); err != nil {
		t.Fatal(err)
	}
	foreignStore, foreignBinding, unlocked, linked := d, d, d, d
	foreignStore.Store = filepath.Join(other, ProtectionStore)
	foreignBinding.Binding, _ = cacheIdentity(other)
	unlocked.ResourcePath, unlocked.Store = other, filepath.Join(other, ProtectionStore)
	linked.ResourcePath, linked.Store = link, filepath.Join(link, ProtectionStore)
	for name, foreign := range map[string]protection.Descriptor{"store": foreignStore, "binding": foreignBinding, "lock": unlocked, "link": linked} {
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
