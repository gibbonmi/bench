package spec

import (
	"bytes"
	"reflect"
	"testing"
)

// TestFenceTokensTrailingNewlineParity pins the parser against a spec whose final
// fence line lacks a trailing newline. FenceTokens splits on "\n" directly, so the
// unterminated last element from strings.Split still carries its content. A scanner
// keyed on a trailing "\n" per token instead would drop it.
func TestFenceTokensTrailingNewlineParity(t *testing.T) {
	terminated := []byte("## Ownership fences\n\n- `internal/example/`\n")
	unterminated := bytes.TrimSuffix(terminated, []byte("\n"))
	if bytes.Equal(terminated, unterminated) {
		t.Fatal("fixture invalid: terminated fixture already lacks a trailing newline")
	}

	got, _ := FenceTokens(unterminated)
	want, _ := FenceTokens(terminated)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FenceTokens(unterminated) = %#v, want %#v (same as terminated)", got, want)
	}
	if !reflect.DeepEqual(want, []string{"internal/example/"}) {
		t.Fatalf("fixture invalid: terminated form parsed to %#v, want [\"internal/example/\"]", want)
	}
}

// TestFenceTokensSkipParentheticalAnnotation pins the rule that parenthetical prose
// annotates rather than authorizes, including a parenthetical that spans lines.
func TestFenceTokensSkipParentheticalAnnotation(t *testing.T) {
	content := []byte("## Ownership fences\n\n- `internal/a/` (see `internal/b/`\n  and `internal/c/`)\n- `internal/d/`\n")
	got, declared := FenceTokens(content)
	if !declared {
		t.Fatal("FenceTokens reported no declared section for content that opens one")
	}
	if want := []string{"internal/a/", "internal/d/"}; !reflect.DeepEqual(got, want) {
		t.Errorf("FenceTokens = %#v, want %#v", got, want)
	}
}

// TestFenceTokensEndTheSectionAtAnyHeading pins the section bound. A level-2-or-deeper
// heading closes the section, so an entry written under a subsection is outside it.
// internal/coverage's review-pickup check depends on this bound.
func TestFenceTokensEndTheSectionAtAnyHeading(t *testing.T) {
	content := []byte("## Ownership fences\n\n### Paths\n\n- `internal/a/`\n")
	got, declared := FenceTokens(content)
	if !declared {
		t.Fatal("FenceTokens reported no declared section for content that opens one")
	}
	if len(got) != 0 {
		t.Errorf("FenceTokens = %#v, want no tokens below the subsection heading", got)
	}
}

// TestFenceTokensReportAnAbsentSection separates a spec that declares no fences from
// one whose declared section is empty. The two are different bootstrap answers.
func TestFenceTokensReportAnAbsentSection(t *testing.T) {
	if _, declared := FenceTokens([]byte("# spec\n\nStatus: staged\n")); declared {
		t.Error("FenceTokens reported a declared section for a spec that has none")
	}
	if tokens, declared := FenceTokens([]byte("## Ownership fences\n\n## Next\n")); !declared || len(tokens) != 0 {
		t.Errorf("FenceTokens = %#v, declared %v; want an empty declared section", tokens, declared)
	}
}
