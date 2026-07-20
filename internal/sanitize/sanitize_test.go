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
// its byte-count suffix — the escaping policy moved here from internal/intent. Control
// runes are constructed from their code points and escape tokens from bs, so no fragile
// backslash literal is typed inline.
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
