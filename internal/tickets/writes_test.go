package tickets

import "testing"

// TestWritesEntryCover is the `Writes:` grammar contract: the split drops the (new)
// marker and a trailing directory slash, and an entry covers a path only at a `/`
// segment boundary, so `internal/d` never claims `internal/dx`.
func TestWritesEntryCover(t *testing.T) {
	splits := []struct {
		entry, path string
		isNew       bool
	}{
		{"internal/x.go (new)", "internal/x.go", true},
		{"internal/d/", "internal/d", false},
	}
	for _, split := range splits {
		if path, isNew := WritesPath(split.entry); path != split.path || isNew != split.isNew {
			t.Errorf("WritesPath(%q) = %q, %t; want %q, %t", split.entry, path, isNew, split.path, split.isNew)
		}
	}
	covers := []struct {
		entry, path string
		want        bool
	}{
		{"internal/d", "internal/d/a.go", true},
		{"internal/d", "internal/dx/a.go", false},
	}
	for _, cover := range covers {
		if got := Covers(cover.entry, cover.path); got != cover.want {
			t.Errorf("Covers(%q, %q) = %t, want %t", cover.entry, cover.path, got, cover.want)
		}
	}
}
