package gate

import (
	"reflect"
	"testing"
)

// expectedResolvers is the Source → derivation binding the identity check grades the
// registry against: an independently authored second statement of what the registry
// already carries, admissible under the one-source-per-fact exception because its
// independence is what makes a resolver swap red. That red is demonstrated by the two
// recorded swaps in decisions/assets/ft183-derivation-binding.md — vet's resolver to
// contractInputs, canary's to shellcheckInputs — each observed failing this check on
// 2026-08-03, where both pass the derivation-source conformance check. The
// hand-declared canary row is bound here precisely because that check exempts it.
var expectedResolvers = map[Source]func(*inputResolver) ([]string, []string, error){
	SourceBuildClosure:                              (*inputResolver).buildClosure,
	SourceModuleTestClosure:                         (*inputResolver).moduleClosure,
	SourceModuleTestClosureWithSealAndAgentMarkdown: (*inputResolver).contractInputs,
	SourceShellcheckArgv:                            (*inputResolver).shellcheckInputs,
	SourceHandDeclared:                              (*inputResolver).canaryInputs,
}

// TestRegistryRowsResolveThroughTheirNamedDerivation closes the label↔function gap the
// derivation-source check leaves open: that check proves a row is derivation-sourced,
// not that it resolves through the function its source label names. Exhaustiveness is
// per row — a row whose source has no expectation refuses, so a new registry row must
// declare its binding here to go green rather than silently reopening the gap.
func TestRegistryRowsResolveThroughTheirNamedDerivation(t *testing.T) {
	t.Parallel()
	for _, declaration := range componentInputDeclarations() {
		want, named := expectedResolvers[declaration.source]
		if !named {
			t.Fatalf("component %q declares source %q with no expected derivation; add its binding to expectedResolvers", declaration.component, declaration.source)
		}
		if reflect.ValueOf(declaration.resolve).Pointer() != reflect.ValueOf(want).Pointer() {
			t.Fatalf("component %q resolves through a function its source label %q does not name", declaration.component, declaration.source)
		}
	}
}

// TestResolverPointerIdentityDiscriminates guards the assumption the check above
// depends on: pointer equality identifies a derivation only while registry rows and
// table entries stay method expressions with unmerged bodies. A bound method value
// wraps the method in fresh code, so it must not share the expression's pointer, and
// the table's derivations must all be pairwise distinct — a linker that merged two
// bodies would leave the comparison unable to tell a swap from the real binding.
func TestResolverPointerIdentityDiscriminates(t *testing.T) {
	t.Parallel()
	expression := reflect.ValueOf((*inputResolver).moduleClosure).Pointer()
	if expression == 0 {
		t.Fatal("method-expression pointer is zero; pointer identity observes nothing")
	}
	if again := reflect.ValueOf((*inputResolver).moduleClosure).Pointer(); again != expression {
		t.Fatal("the same method expression yields two pointers; pointer identity is unstable")
	}
	if reflect.ValueOf((&inputResolver{}).moduleClosure).Pointer() == expression {
		t.Fatal("a bound method value shares the method expression's pointer; a rewritten row would grade as its named derivation")
	}
	seen := map[uintptr]Source{}
	for source, fn := range expectedResolvers {
		pointer := reflect.ValueOf(fn).Pointer()
		if prior, dup := seen[pointer]; dup {
			t.Fatalf("sources %q and %q name one code pointer; the table cannot discriminate their derivations", prior, source)
		}
		seen[pointer] = source
	}
}
