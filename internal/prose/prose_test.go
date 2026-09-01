package prose

import (
	"strings"
	"testing"
)

func TestGradeNamedResultsExposesSentenceFinding(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".bench/prose-exclusions", "")
	sentence := words(27)
	write(t, root, "docs/notes.md", sentence)

	got := GradeNamedResults(root, []string{"docs/notes.md"})
	if len(got) != 1 {
		t.Fatalf("GradeNamedResults returned %d results, want 1", len(got))
	}
	wantSentence := strings.TrimSpace(sentence)
	if got[0].Path != "docs/notes.md" || got[0].Line != 1 || got[0].Rule != KindSentence || got[0].Count != 27 || got[0].Sentence != wantSentence {
		t.Fatalf("GradeNamedResults result = %+v, want path docs/notes.md, line 1, sentence rule, count 27, and sentence %q", got[0], wantSentence)
	}
}
