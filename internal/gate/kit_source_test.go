package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
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

// TestKitSourceCheckoutResolvesARelativeKit pins the absolute compare (TT51). A BENCH_KIT
// of "." names the current directory, so a compare that kept the relative spelling would
// answer false for the root that the current directory is.
func TestKitSourceCheckoutResolvesARelativeKit(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	t.Setenv("BENCH_KIT", ".")

	if !KitSourceCheckout(root) {
		t.Fatalf("KitSourceCheckout(%q) with BENCH_KIT=. = false, want true", root)
	}
}

// TestKitSourceCheckoutResolvesARelativeKitAgainstTheWorkingDirectory pins the other half
// of the absolute compare. A BENCH_KIT of "." names the current directory and not the
// root, so a compare that joined the kit onto the root would count every root as the kit.
func TestKitSourceCheckoutResolvesARelativeKitAgainstTheWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	t.Chdir(t.TempDir())
	t.Setenv("BENCH_KIT", ".")

	if KitSourceCheckout(root) {
		t.Fatalf("KitSourceCheckout(%q) with BENCH_KIT=. and another working directory = true, want false", root)
	}
}

// TestKitSourceCheckoutAnswersFalseOnAResolveError pins the resolve-error rule. A deleted
// working directory leaves "." with no absolute spelling, so both resolves fail, and two
// failed resolves must not compare equal as two empty spellings.
func TestKitSourceCheckoutAnswersFalseOnAResolveError(t *testing.T) {
	gone := filepath.Join(t.TempDir(), "gone")
	if err := os.Mkdir(gone, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(gone)
	if err := os.Remove(gone); err != nil {
		capability.Environment(t, "cannot remove the working directory: "+err.Error())
	}
	if _, err := filepath.Abs("."); err == nil {
		capability.Environment(t, "a deleted working directory still has an absolute spelling on this platform")
	}
	t.Setenv("BENCH_KIT", ".")

	if KitSourceCheckout(".") {
		t.Fatal("KitSourceCheckout(\".\") with a deleted working directory = true, want false")
	}
}
