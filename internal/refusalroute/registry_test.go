package refusalroute

import (
	"slices"
	"strings"
	"testing"
)

// TestRegistryFaceNamesAreUnique runs the uniqueness walk over the declared inventory and
// over injected lists, so the walk both passes the real registry and bites on a
// duplicate. It covers RR04.
func TestRegistryFaceNamesAreUnique(t *testing.T) {
	if err := uniqueNames(inventory); err != nil {
		t.Fatalf("declared inventory: %v", err)
	}
	if err := uniqueNames([]Face{{Name: "first"}, {Name: "second"}}); err != nil {
		t.Fatalf("distinct names refused: %v", err)
	}
	err := uniqueNames([]Face{{Name: "twice-declared"}, {Name: "second"}, {Name: "twice-declared"}})
	if err == nil || !strings.Contains(err.Error(), "twice-declared") {
		t.Fatalf("duplicate name passed the walk or went unnamed: %v", err)
	}
}

// TestUnregisteredFaceRefusesWithAReviewerRoute asks the one constructor for a name that
// no face declares. The bookkeeping fault must come back as a refusal with a reviewer
// route rather than a panic or an empty route. It covers RR07.
func TestUnregisteredFaceRefusesWithAReviewerRoute(t *testing.T) {
	got := New("no-such-face", Facts{Sentence: "an observed cause", Paths: []string{"a.go"}})
	if want := "refusal face no-such-face is unregistered"; got.Sentence != want || got.Error() != want {
		t.Fatalf("sentence %q, error %q; want %q", got.Sentence, got.Error(), want)
	}
	if got.Face.Authority != Reviewer {
		t.Fatalf("authority %q; want %q", got.Face.Authority, Reviewer)
	}
	if !strings.HasPrefix(got.Route, "reviewer: ") || strings.TrimSpace(strings.TrimPrefix(got.Route, "reviewer: ")) == "" {
		t.Fatalf("route %q; want a reviewer route with a step", got.Route)
	}
}

// TestConstructorCarriesTheRegisteredFace asks the constructor for a face of an injected
// inventory. The declared sentence fills an empty observed one, an observed sentence
// keeps its cause, and the paths and the rendered route travel on the refusal.
func TestConstructorCarriesTheRegisteredFace(t *testing.T) {
	faces := []Face{{
		Verb:      Land,
		Name:      "destination-not-clean",
		Sentence:  "landing destination is not clean",
		Authority: Agent,
		Route:     []Step{Instruction(Text("commit the work")), Command(Composed("rerun"))},
	}}
	facts := Facts{Paths: []string{"a.go", "b.go"}, Values: map[string]string{"rerun": "bench worktree land"}}

	got := newIn(faces, "destination-not-clean", facts)
	if got.Face.Name != "destination-not-clean" || got.Sentence != "landing destination is not clean" {
		t.Fatalf("face %q, sentence %q; want the declared face and its sentence", got.Face.Name, got.Sentence)
	}
	if !slices.Equal(got.Paths, facts.Paths) {
		t.Fatalf("paths %q; want %q", got.Paths, facts.Paths)
	}
	if want := "commit the work; then bench worktree land"; got.Route != want {
		t.Fatalf("route %q; want %q", got.Route, want)
	}

	facts.Sentence = "landing destination is not clean: a.go is modified"
	if got := newIn(faces, "destination-not-clean", facts); got.Sentence != facts.Sentence {
		t.Fatalf("sentence %q; want the observed %q", got.Sentence, facts.Sentence)
	}
}
