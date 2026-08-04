package spec

import (
	"slices"
	"strings"
	"testing"
)

func TestParseBuildExposesEveryGrammarOperation(t *testing.T) {
	// buildOperations is the only source of which operations exist, so a row added there
	// without a case here fails below rather than going silently unparsed.
	invocations := map[string][]string{
		"start": nil, "assign": {"--ticket", "one.md", "--request", "request"},
		"checkpoint": {"--assignment", "a", "--evidence", "/receipt"},
		"integrate":  {"--assignment", "a"}, "review": {"--evidence", "/review"},
		"status": {"--full"}, "promote": nil, "abandon": {"--apply", "fingerprint"},
		"reclaim": nil,
	}
	if missing := keysMissingFrom(buildOperations, invocations); len(missing) > 0 {
		t.Errorf("grammar operations with no parse case: %v", missing)
	}
	if stale := keysMissingFrom(invocations, buildOperations); len(stale) > 0 {
		t.Errorf("parse cases naming no grammar operation: %v", stale)
	}
	for operation, flags := range invocations {
		t.Run(operation, func(t *testing.T) {
			got, out, code := ParseBuild(append([]string{operation, "slug"}, flags...))
			if code != 0 || out != "" || got.Operation != operation || got.Slug != "slug" {
				t.Fatalf("ParseBuild = %+v, %q, %d", got, out, code)
			}
		})
	}
	if _, out, code := ParseBuild([]string{"nonesuch", "slug"}); code != 2 || !strings.Contains(out, "unknown argument: nonesuch") {
		t.Fatalf("unknown operation = %q, %d", out, code)
	}
}

// keysMissingFrom reports the keys of from that in does not carry, sorted.
func keysMissingFrom[A, B any](from map[string]A, in map[string]B) []string {
	var missing []string
	for key := range from {
		if _, ok := in[key]; !ok {
			missing = append(missing, key)
		}
	}
	slices.Sort(missing)
	return missing
}

// TestParseBuildNoArgDiagnosticListsExactlyTheGrammarOperations pins the empty-argument
// diagnostic to buildOperations itself, in the declared lifecycle order, so an operation
// added or removed from the table changes the diagnostic without a second list to edit.
func TestParseBuildNoArgDiagnosticListsExactlyTheGrammarOperations(t *testing.T) {
	_, out, code := ParseBuild(nil)
	if code != 2 {
		t.Fatalf("ParseBuild(nil) code = %d, want 2", code)
	}
	want := "usage: bench spec build (missing argument: start|assign|checkpoint|integrate|review|status|promote|abandon|reclaim)\n"
	if out != want {
		t.Fatalf("ParseBuild(nil) = %q, want %q", out, want)
	}
	listed := strings.Split(strings.TrimSuffix(strings.TrimPrefix(out, "usage: bench spec build (missing argument: "), ")\n"), "|")
	if missing := keysMissingFrom(buildOperations, toSet(listed)); len(missing) > 0 {
		t.Errorf("grammar operations missing from no-arg diagnostic: %v", missing)
	}
	if stale := keysMissingFrom(toSet(listed), buildOperations); len(stale) > 0 {
		t.Errorf("no-arg diagnostic names operations absent from buildOperations: %v", stale)
	}
}

// toSet turns a slice into a presence map so keysMissingFrom can compare it against a
// map[string]usage.Grammar.
func toSet(keys []string) map[string]struct{} {
	set := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		set[key] = struct{}{}
	}
	return set
}

// TestParseBuildNoArgDiagnosticIsOrderStable renders the empty-argument diagnostic
// repeatedly: derivation must read a declared order, not Go's randomized map iteration,
// or the operation list would differ run to run.
func TestParseBuildNoArgDiagnosticIsOrderStable(t *testing.T) {
	_, first, _ := ParseBuild(nil)
	for i := 0; i < 20; i++ {
		if _, out, _ := ParseBuild(nil); out != first {
			t.Fatalf("run %d: ParseBuild(nil) = %q, want %q", i, out, first)
		}
	}
}

// The refused token is asserted, not just the exit code: a grammar with no positional
// bound still fails the one-slug check downstream, but blames the first slug rather than
// the excess one.
func TestParseBuildBoundsReclaimToOneSlug(t *testing.T) {
	for _, test := range []struct {
		args   []string
		blames string
	}{
		{[]string{"reclaim"}, "missing argument"},
		{[]string{"reclaim", "one", "two"}, "unknown argument: two"},
		{[]string{"reclaim", "slug", "--force"}, "unknown argument: --force"},
	} {
		t.Run(strings.Join(test.args, " "), func(t *testing.T) {
			if _, out, code := ParseBuild(test.args); code != 2 || !strings.Contains(out, test.blames) {
				t.Fatalf("ParseBuild(%v) = %q, %d, want %q", test.args, out, code, test.blames)
			}
		})
	}
}

func TestParseBuildParsesBothReclaimArms(t *testing.T) {
	for _, args := range [][]string{{"reclaim", "slug"}, {"reclaim", "slug", "--apply", "fingerprint"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			got, out, code := ParseBuild(args)
			if code != 0 || out != "" || got.Operation != "reclaim" || got.Slug != "slug" {
				t.Fatalf("ParseBuild(%v) = %+v, %q, %d", args, got, out, code)
			}
		})
	}
}

func TestParseBuildKeepsFlagValuesAndTerminatedTextLiteral(t *testing.T) {
	got, out, code := ParseBuild([]string{"assign", "slug", "--ticket", "promote", "--request", "--assignment"})
	if code != 0 || out != "" || got.Flags["--ticket"] != "promote" || got.Flags["--request"] != "--assignment" {
		t.Fatalf("flag values = %+v, %q, %d", got, out, code)
	}
	for _, args := range [][]string{{"start", "--", "slug"}, {"assign", "slug", "--", "--ticket"}} {
		if _, out, code := ParseBuild(args); code != 2 || !strings.HasPrefix(out, "usage:") {
			t.Fatalf("terminated text %v = %q, %d", args, out, code)
		}
	}
}
