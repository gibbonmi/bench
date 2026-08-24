package diff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The explicit pair is immutable: it reports base-to-tip and nothing from the live
// checkout, so the output cannot change when the checkout moves after the tip.
func TestSourceTipPairReportsTheImmutableRange(t *testing.T) {
	root, base, feature := seedCompatibilityRepo(t)
	out, code := Command([]string{"--base", base, "--source-tip", feature})
	if code != 0 {
		t.Fatalf("pair diff exit = %d:\n%s", code, out)
	}
	for _, want := range []string{base, feature, "committed.txt", "explicit pair"} {
		if !strings.Contains(out, want) {
			t.Fatalf("pair diff omitted %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "untracked.txt") {
		t.Fatalf("pair view leaked live-checkout facts:\n%s", out)
	}
	if err := os.WriteFile(filepath.Join(root, "later.txt"), []byte("moved\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	again, code := Command([]string{"--base", base, "--source-tip", feature})
	if code != 0 || again != out {
		t.Fatalf("pair view changed with the checkout (exit %d):\n%s", code, again)
	}
}

// The tip never comes from the checkout implicitly, and the pair never mixes with
// --commit; a misuse of the grammar exits 2 with the usage line, and a bad value is a
// structured error that blames the flag the caller got wrong.
func TestSourceTipRefusals(t *testing.T) {
	_, base, feature := seedCompatibilityRepo(t)
	for _, args := range [][]string{
		{"--source-tip", feature},
		{"--commit", feature, "--source-tip", feature},
	} {
		out, code := Command(args)
		if code != 2 || !strings.HasPrefix(out, "usage:") {
			t.Fatalf("Command(%v) = (%d,%q), want the usage refusal", args, code, out)
		}
	}
	out, code := Command([]string{"--base", base, "--source-tip", "missing"})
	if code != 1 || !strings.HasPrefix(out, "error:") || !strings.Contains(out, "--source-tip") {
		t.Fatalf("unresolvable tip = (%d,%q), want an error blaming --source-tip", code, out)
	}
	out, code = Command([]string{"--base", feature, "--source-tip", base})
	if code != 1 || !strings.HasPrefix(out, "error:") || !strings.Contains(out, "ancestor") {
		t.Fatalf("non-ancestor base = (%d,%q), want the ancestor error", code, out)
	}
}
