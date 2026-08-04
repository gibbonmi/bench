package specbuild

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func workingSubjectRefusal(t *testing.T, fixture preconditionFixture, op mutation) string {
	t.Helper()
	specPath := filepath.Join(fixture.root, "specs", "build demo", "spec.md")
	if _, err := fixture.service.preconditions(op, "build demo", specPath, nil, "", ""); err != nil {
		return err.Error()
	}
	t.Fatalf("preconditions(%q) accepted an unusable working subject", op)
	return ""
}
func dirtyCheckout(t *testing.T, fixture preconditionFixture) {
	t.Helper()
	write(t, filepath.Join(fixture.root, "tracked.txt"), "dirty\n")
}
func detachedCheckout(t *testing.T, fixture preconditionFixture) {
	t.Helper()
	git(t, fixture.root, "checkout", "--detach", "-q")
}

// strictCheckoutMutations deliberately repeats the production policy: omitting a
// strict operation or mutating it to provisional must leave this dirty-checkout
// oracle red instead of shrinking the tested family with the implementation.
var strictCheckoutMutations = []mutation{mutationStart, mutationPromote, mutationAbandon}

// operationWordIn reads the whole rendered operation, so a word carrying a flag fails rather than passing on its first token.
func operationWordIn(t *testing.T, message string) string {
	t.Helper()
	const prefix = "spec build "
	end := strings.Index(message, " requires ")
	if !strings.HasPrefix(message, prefix) || end < len(prefix) {
		t.Fatalf("refusal %q does not name an operation between %q and %q", message, prefix, " requires ")
	}
	return message[len(prefix):end]
}

func TestDirtyCheckoutRefusalNamesOperation(t *testing.T) {
	for _, op := range strictCheckoutMutations {
		t.Run(string(op), func(t *testing.T) {
			fixture := newPreconditionFixture(t, false)
			dirtyCheckout(t, fixture)
			message := workingSubjectRefusal(t, fixture, op)
			if want := "spec build " + string(op) + " requires a clean working checkout"; !strings.Contains(message, want) {
				t.Fatalf("dirty refusal = %q, want %q", message, want)
			}
			if op != mutationStart && strings.Contains(message, "spec build start") {
				t.Fatalf("dirty refusal for %q borrows start's wording: %q", op, message)
			}
		})
	}
}
func TestNoWorkingBranchRefusalNamesOperation(t *testing.T) {
	for _, op := range lifecycleMutations {
		t.Run(string(op), func(t *testing.T) {
			fixture := newPreconditionFixture(t, false)
			detachedCheckout(t, fixture)
			message := workingSubjectRefusal(t, fixture, op)
			if want := "spec build " + string(op) + " requires a checked-out working branch"; message != want {
				t.Fatalf("no-branch refusal = %q, want %q", message, want)
			}
		})
	}
}
func TestEmptyMutationNeverRendersAnonymously(t *testing.T) {
	for _, condition := range []struct {
		name  string
		apply func(*testing.T, preconditionFixture)
	}{{"dirty checkout", dirtyCheckout}, {"no working branch", detachedCheckout}} {
		t.Run(condition.name, func(t *testing.T) {
			fixture := newPreconditionFixture(t, false)
			condition.apply(t, fixture)
			message := workingSubjectRefusal(t, fixture, mutation(""))
			if word := operationWordIn(t, message); strings.TrimSpace(word) == "" {
				t.Fatalf("zero token rendered an anonymous refusal: %q", message)
			}
		})
	}
}
func TestRecompositionRefusalStillNamesPromote(t *testing.T) {
	fixture := newPreconditionFixture(t, true)
	advanceWorking(t, fixture.root)
	_, _, err := fixture.service.Assign(t.Context(), "build demo", "one.md", "next request")
	if err == nil {
		t.Fatal("Assign accepted a moved working tip")
	}
	if !strings.Contains(err.Error(), "bench spec build promote") || strings.Contains(err.Error(), "assign") {
		t.Fatalf("recomposition refusal = %q, want it to point at promote", err)
	}
}
func TestOperationWordMatchesSubcommandVerb(t *testing.T) {
	verbs := specBuildSubcommandVerbs(t)
	for _, op := range lifecycleMutations {
		t.Run(string(op), func(t *testing.T) {
			fixture := newPreconditionFixture(t, false)
			detachedCheckout(t, fixture)
			word := operationWordIn(t, workingSubjectRefusal(t, fixture, op))
			if word != string(op) || !verbs[word] {
				t.Fatalf("operation word %q for %q is not the subcommand verb; parsed verbs: %v", word, op, verbs)
			}
		})
	}
}

// TestLifecycleFamilyMatchesSubcommandVerbs pins membership against the shipped
// CLI instead of the production list every per-operation test iterates, so an
// omitted operation cannot shrink implementation and coverage together.
func TestLifecycleFamilyMatchesSubcommandVerbs(t *testing.T) {
	// Two verbs sit outside the preconditioned family, and each exclusion is asserted
	// against the shipped usage block so it cannot outlive the verb it names. `status`
	// inspects a run and mutates nothing. `reclaim` acts only on refs the run record
	// already proves dead, and it reaches runs whose spec is retired and whose candidate
	// promotion already deleted, so the identity checks preconditions makes would put
	// exactly the residue it exists for out of reach; the disposition filter, not the
	// checkout envelope, is what keeps a live run's refs safe from it.
	unpreconditioned := []string{"status", "reclaim"}
	parsed := specBuildSubcommandVerbs(t)
	want := map[string]bool{}
	for verb := range parsed {
		want[verb] = true
	}
	for _, verb := range unpreconditioned {
		if !parsed[verb] {
			t.Fatalf("bin/bench.sh no longer documents %q; the exclusion here is stale", verb)
		}
		delete(want, verb)
	}
	got := map[string]bool{}
	for _, op := range lifecycleMutations {
		if got[string(op)] {
			t.Fatalf("lifecycleMutations repeats %q", op)
		}
		got[string(op)] = true
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("lifecycleMutations = %v, want every mutating subcommand verb %v", got, want)
	}
}

// specBuildSubcommandVerbs reads the verbs from the shipped usage block so this suite cannot drift from the CLI it quotes.
func specBuildSubcommandVerbs(t *testing.T) map[string]bool {
	t.Helper()
	usage, err := os.ReadFile(filepath.Join("..", "..", "bin", "bench.sh"))
	if err != nil {
		t.Fatal(err)
	}
	verbs := map[string]bool{}
	for _, line := range strings.Split(string(usage), "\n") {
		if fields := strings.Fields(line); len(fields) > 3 && fields[0] == "bench" && fields[1] == "spec" && fields[2] == "build" {
			verbs[fields[3]] = true
		}
	}
	if len(verbs) == 0 {
		t.Fatal("bin/bench.sh lists no bench spec build subcommands")
	}
	return verbs
}
