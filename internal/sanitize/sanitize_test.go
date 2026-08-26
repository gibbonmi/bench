package sanitize

import (
	"strings"
	"testing"
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

// TestPreviewBoundariesAndControls pins Preview's escaping, its 120-code-point cap, and
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

	if got := Preview(strings.Repeat("é", 120)); got != strings.Repeat("é", 120) {
		t.Errorf("Preview capped a 120-rune string: %q", got)
	}
	if got := Preview(strings.Repeat("é", 121)); got != strings.Repeat("é", 120)+"… (242 bytes)" {
		t.Errorf("Preview(121 runes) = %q", got)
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
