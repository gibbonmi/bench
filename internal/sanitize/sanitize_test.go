package sanitize

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

// bs is a single backslash, built from a raw literal so the escape-token expectations
// below never rely on a backslash escape sequence surviving transport intact.
const bs = `\`

// containsControl reports whether s carries any raw C0 control byte or DEL — the thing a
// sanitizer must never leave behind.
func containsControl(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}

// TestPreviewBoundariesAndControls pins Preview's escaping, its bounds.PreviewRuneLimit cap, and
// its byte-count suffix. Control runes are constructed from their code points and
// escape tokens from bs, so no fragile backslash literal is typed inline.
func TestPreviewBoundariesAndControls(t *testing.T) {
	if got := Preview(""); got != "" {
		t.Errorf("Preview(empty) = %q", got)
	}
	if got := Preview("a"); got != "a" {
		t.Errorf("Preview(a) = %q", got)
	}

	in := "line" + string(rune(0x0a)) + string(rune(0x1b)) + string(rune(0x07))
	got := Preview(in)
	if containsControl(got) {
		t.Errorf("Preview left a raw control byte: %q", got)
	}
	for _, tok := range []string{"line", bs + "n", bs + "u001b", bs + "u0007"} {
		if !strings.Contains(got, tok) {
			t.Errorf("Preview(%q) = %q, missing escape token %q", in, got, tok)
		}
	}

	backslash := Preview("back" + bs + "slash")
	if backslash != "back"+bs+bs+"slash" {
		t.Errorf("Preview did not double a literal backslash: %q", backslash)
	}

	atCap := strings.Repeat("é", bounds.PreviewRuneLimit)
	if got := Preview(atCap); got != atCap {
		t.Errorf("Preview capped a %d-rune string: %q", bounds.PreviewRuneLimit, got)
	}
	overCap := strings.Repeat("é", bounds.PreviewRuneLimit+1)
	if got := Preview(overCap); got != atCap+"… (482 bytes)" {
		t.Errorf("Preview(%d runes) = %q", bounds.PreviewRuneLimit+1, got)
	}
}

// TestControlsEscapesWithoutCapping pins the distinction from Preview: Controls escapes
// every control rune but never truncates and never appends a byte suffix, so a long
// control-bearing string round-trips its whole content in escaped form.
func TestControlsEscapesWithoutCapping(t *testing.T) {
	long := strings.Repeat("é", 200)
	if got := Controls(long); got != long {
		t.Errorf("Controls dropped or truncated a 200-rune string: len=%d", len([]rune(got)))
	}
	in := "a" + string(rune(0x1b)) + "[31mb" + string(rune(0x0a)) + "c" + string(rune(0x09)) + "d"
	want := "a" + bs + "u001b" + "[31mb" + bs + "n" + "c" + bs + "t" + "d"
	if got := Controls(in); got != want {
		t.Errorf("Controls(%q) = %q, want %q", in, got, want)
	}
	if containsControl(Controls("x" + string(rune(0x1b)) + "y")) {
		t.Error("Controls left a raw control byte in its output")
	}
}

// TestStripRemovesControlsAndEscapesNothing pins the third duty against the other two.
// An ESC leaves no trace, a tab survives because a cell may be tab-aligned, and a
// backslash stays single: the sink's own encoder escapes it, and escaping here first
// would double what the encoder then doubles again. A non-ASCII rune passes whole, which
// is what separates a rune-wise filter from a byte-wise one.
func TestStripRemovesControlsAndEscapesNothing(t *testing.T) {
	in := "a" + string(rune(0x1b)) + "[31mb" + string(rune(0x0a)) + "c" + string(rune(0x09)) + "d"
	want := "a[31mbc" + "\t" + "d"
	if got := Strip(in); got != want {
		t.Errorf("Strip(%q) = %q, want %q", in, got, want)
	}
	if got := Strip("back" + bs + "slash"); got != "back"+bs+"slash" {
		t.Errorf("Strip escaped a literal backslash: %q", got)
	}
	if got := Strip("caf" + "é" + " ☃"); got != "café ☃" {
		t.Errorf("Strip corrupted a non-ASCII rune: %q", got)
	}
	// DEL and the C1 controls sit at or above U+0020, so they pass. This filter guards
	// what the cell encoder refuses, and the encoder refuses only below U+0020.
	if got := Strip(string(rune(0x7f)) + string(rune(0x85))); got != string(rune(0x7f))+string(rune(0x85)) {
		t.Errorf("Strip removed a rune at or above U+0020: %q", got)
	}
	// Carriage return is Strip's own exclusion rather than the encoder's. The encoder
	// escapes it, so toon.Representable passes it, and Strip still drops it: a cell
	// that carries one would break the row it sits on.
	if got := Strip("a" + string(rune(0x0d)) + "b"); got != "ab" {
		t.Errorf("Strip kept a carriage return: %q", got)
	}
}

// TestStripKeepsOnlyTabBelowSpace sweeps every C0 byte to pin the composed result:
// tab is the one rune below U+0020 a cell may carry. Strip asks the cell encoder's own
// predicate which byte it refuses, so this sweep is what catches a widened cell policy
// arriving through that composition.
func TestStripKeepsOnlyTabBelowSpace(t *testing.T) {
	for b := byte(0); b < 0x20; b++ {
		in := "x" + string([]byte{b}) + "y"
		want := "xy"
		if b == '\t' {
			want = "x\ty"
		}
		if got := Strip(in); got != want {
			t.Errorf("Strip(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestPreformattedPreservesLayoutWhitespace pins the <pre>-panel variant. Newline and
// tab pass through verbatim, so multi-line layout survives. Carriage return and every
// other control rune still escape through the same \uXXXX mechanism Controls uses. A
// raw C0 byte must never reach the output.
func TestPreformattedPreservesLayoutWhitespace(t *testing.T) {
	in := "a" + string(rune(0x0a)) + "b" + string(rune(0x09)) + "c" + string(rune(0x0d)) + "d" + string(rune(0x07)) + "e"
	got := Preformatted(in)
	want := "a" + "\n" + "b" + "\t" + "c" + bs + "u000d" + "d" + bs + "u0007" + "e"
	if got != want {
		t.Errorf("Preformatted(%q) = %q, want %q", in, got, want)
	}
	if strings.ContainsRune(got, rune(0x0d)) || strings.ContainsRune(got, rune(0x07)) {
		t.Error("Preformatted left a raw non-layout control byte in its output")
	}
}
