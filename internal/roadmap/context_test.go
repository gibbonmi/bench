package roadmap

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestContextBodyLimitBoundaries(t *testing.T) {
	for _, n := range []int{4095, 4096, 4097} {
		got, full, truncated := limited(strings.Repeat("x", n), false)
		if full != n || truncated != (n > 4096) || len(got) != min(n, 4096) {
			t.Fatalf("n=%d got bytes=%d full=%d truncated=%v", n, len(got), full, truncated)
		}
	}
	s := strings.Repeat("x", 4095) + "€"
	got, full, truncated := limited(s, false)
	if full != 4098 || !truncated || !utf8.ValidString(got) || len(got) != 4095 {
		t.Fatalf("UTF-8 boundary got bytes=%d full=%d truncated=%v valid=%v", len(got), full, truncated, utf8.ValidString(got))
	}
	got, full, truncated = limited(s, true)
	if got != s || full != 4098 || truncated {
		t.Fatal("--full did not preserve body")
	}
}
