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

// TestParseWriteRefusesAnUncleanPath covers CR19 and CR20. An entry whose path has a ".."
// segment, an absolute path, an unclean spelling, or no path past the repository root is a
// grammar error, because a lexical cover of such a spelling can authorize a file outside the
// fence or the repository. A clean entry still parses with its markers and directory
// spelling, and a ".." inside one segment name is no ".." segment.
func TestParseWriteRefusesAnUncleanPath(t *testing.T) {
	for _, entry := range []string{
		"internal/x/../../cmd/y.go",
		"../outside.txt",
		"internal/a.go/.. (deleted)",
		"/etc/passwd",
		"/abs/dir/",
		"internal//x.go",
		"./internal/x.go (new)",
		"internal/./x.go",
		"internal/x/.",
		"/",
		"/ (new)",
		"./",
		".",
	} {
		if write, err := ParseWrite(entry); err == nil {
			t.Errorf("ParseWrite(%q) = %+v, want an unclean-path refusal", entry, write)
		}
	}
	for _, entry := range []string{"internal/x.go", "internal/d/", "internal/x.go (new)", "internal/gone.go (deleted)", "a..b/x.go"} {
		if _, err := ParseWrite(entry); err != nil {
			t.Errorf("ParseWrite(%q) = %v, want a clean entry to parse", entry, err)
		}
	}
}

// TestRepositoryPathRefusesTheRoot pins the root rule: neither spelling of the repository root
// names a path below it.
func TestRepositoryPathRefusesTheRoot(t *testing.T) {
	for p, want := range map[string]bool{"": false, ".": false, "internal": true} {
		if got := RepositoryPath(p); got != want {
			t.Errorf("RepositoryPath(%q) = %t, want %t", p, got, want)
		}
	}
}
