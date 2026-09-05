package gate

import (
	"os"
	"path/filepath"
	"testing"
)

// TestKitSourceCheckoutMatchesThroughASymlinkSpelling pins the symlink-only compare. A
// reader reaches the kit repo by one spelling and BENCH_KIT names another, so a compare
// that read the two strings would answer false for the kit's own checkout.
func TestKitSourceCheckoutMatchesThroughASymlinkSpelling(t *testing.T) {
	home := t.TempDir()
	kit := filepath.Join(home, "kit")
	other := filepath.Join(home, "consumer")
	for _, dir := range []string{kit, other} {
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	link := filepath.Join(home, "kit-link")
	if err := os.Symlink(kit, link); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_KIT", kit)

	if !KitSourceCheckout(link) {
		t.Fatalf("KitSourceCheckout(%q) = false, want true", link)
	}
	if KitSourceCheckout(other) {
		t.Fatalf("KitSourceCheckout(%q) = true, want false", other)
	}
	if KitSourceCheckout("") {
		t.Fatal("KitSourceCheckout(\"\") = true, want false")
	}
}
