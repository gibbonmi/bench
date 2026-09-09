package anchors

import "testing"

// TestLocateHonorsSectionScope pins DG28 and DG39: a section kind searches its own
// section body only, and a duplicated owning section — which the evaluator's section
// resolution refuses — locates nothing rather than searching the first duplicate.
func TestLocateHonorsSectionScope(t *testing.T) {
	doc := "# Title\n" +
		"outside needle here\n" +
		"\n" +
		"## Alpha\n" +
		"before\n" +
		"needle here\n" +
		"## Beta\n" +
		"after\n"

	if got := Locate(RequireInSection, "Alpha", "needle here", doc); got != 6 {
		t.Fatalf("needle inside its section: Locate = %d, want line 6", got)
	}
	if got := Locate(RequireInSection, "Alpha", "outside needle", doc); got != 0 {
		t.Fatalf("needle only outside its section: Locate = %d, want 0", got)
	}
	// A whole-file search finds the same needle outside the section, so the section
	// scope — not an absent needle — is what zeroes the row above.
	if got := Locate(Require, "", "outside needle", doc); got != 2 {
		t.Fatalf("whole-file search for the outside needle: Locate = %d, want line 2", got)
	}

	duplicated := "## Alpha\nneedle here\n## Alpha\nneedle here\n"
	if got := Locate(RequireInSection, "Alpha", "needle here", duplicated); got != 0 {
		t.Fatalf("duplicated owning section: Locate = %d, want 0", got)
	}

	if got := Locate(RequireInSection, "Missing", "needle here", doc); got != 0 {
		t.Fatalf("absent section: Locate = %d, want 0", got)
	}
}

// TestLocateMapsCollapsedMatchesToLines pins DG29: a needle split across lines, a needle
// after a multi-line comment, and a case fold that changes a rune's byte length all
// locate the match's first character on its real physical line — never shifted by the
// collapse, the strip, or the fold that found it.
func TestLocateMapsCollapsedMatchesToLines(t *testing.T) {
	t.Run("split across lines", func(t *testing.T) {
		doc := "prefix\none\ntwo suffix\n"
		if got := Locate(Require, "", "one two", doc); got != 2 {
			t.Fatalf("Locate = %d, want line 2 (the first character of \"one\")", got)
		}
	})

	t.Run("after a multi-line comment", func(t *testing.T) {
		doc := "before\n<!-- hidden\nacross lines -->needle right after\nend\n"
		if got := Locate(Require, "", "needle right after", doc); got != 3 {
			t.Fatalf("Locate = %d, want line 3", got)
		}
		if got := Locate(Require, "", "hidden", doc); got != 0 {
			t.Fatalf("needle inside the comment: Locate = %d, want 0", got)
		}
	})

	t.Run("non-ASCII uppercase text", func(t *testing.T) {
		// U+212A KELVIN SIGN lowercases to the ASCII letter k: 3 bytes fold to 1. It sits
		// on the line before the match, not inside it, so a byte-offset implementation runs out
		// of its rune budget two bytes early — before the newline that starts the match's
		// own line — and reports the wrong (earlier) line. A rune-mapped one does not.
		doc := "# Doc\n\n## Notes\nKelvin first\nkelvin note\n## End\n"
		if got := Locate(RequireInSection, "Notes", "kelvin note", doc); got != 5 {
			t.Fatalf("Locate = %d, want line 5", got)
		}
	})

	t.Run("absent needle and absent data", func(t *testing.T) {
		if got := Locate(Require, "", "never appears", "one\ntwo\n"); got != 0 {
			t.Fatalf("absent needle: Locate = %d, want 0", got)
		}
		if got := Locate(Require, "", "anything", ""); got != 0 {
			t.Fatalf("absent data: Locate = %d, want 0", got)
		}
	})

	t.Run("forbid kind locates the violation", func(t *testing.T) {
		if got := Locate(Forbid, "", "banned phrase", "one\nbanned phrase here\n"); got != 2 {
			t.Fatalf("Locate = %d, want line 2 (the forbidden text's own line)", got)
		}
		if got := Locate(Forbid, "", "banned phrase", "one\ntwo\n"); got != 0 {
			t.Fatalf("clean forbid text: Locate = %d, want 0", got)
		}
	})
}
