package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"
)

const canonicalResolveBenchLib = `# shellcheck shell=sh
bench_resolve_wrapper() {
  _bench_root=$(git rev-parse --show-toplevel 2>/dev/null) || _bench_root=
  if [ -n "$_bench_root" ]; then
    for _bench_candidate in "$_bench_root/.bench/bin/bench.sh" "$_bench_root/bin/bench.sh"; do
      if [ -x "$_bench_candidate" ]; then
        printf '%s\n' "$_bench_candidate"
        return 0
      fi
    done
  fi
  command -v bench 2>/dev/null || return 1
}
`

func TestCheckGuardResolverOrderDriftAnchorsOnRenamedFunction(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, ".bench", "hooks", "block-dangerous-git.sh"), `#!/usr/bin/env bash
resolve_wrapper_v2() {
  local root candidate
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    for candidate in "$root/.bench/bin/bench.sh" "$root/bin/bench.sh"; do
      [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
    done
  fi
  command -v bench 2>/dev/null || return 1
}
`)
	writeFixtureFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), canonicalResolveBenchLib)

	diags := checkGuardResolverOrderDrift(root)
	if !anyContains(diags, "no resolve_wrapper() function to anchor on") {
		t.Fatalf("renamed resolver function did not produce a missing-anchor diagnostic: %v", diags)
	}
}

func TestCheckGuardResolverOrderDriftAnchorsOnAbsentShim(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), canonicalResolveBenchLib)

	diags := checkGuardResolverOrderDrift(root)
	if !anyContains(diags, ".bench/hooks/block-dangerous-git.sh is missing") {
		t.Fatalf("absent shim file did not produce a missing-anchor diagnostic: %v", diags)
	}
}

func TestCheckGuardResolverOrderDriftDetectsSwappedOrder(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, ".bench", "hooks", "block-dangerous-git.sh"), `#!/usr/bin/env bash
resolve_wrapper() {
  local root candidate
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    for candidate in "$root/bin/bench.sh" "$root/.bench/bin/bench.sh"; do
      [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
    done
  fi
  command -v bench 2>/dev/null || return 1
}
`)
	writeFixtureFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), canonicalResolveBenchLib)

	diags := checkGuardResolverOrderDrift(root)
	if !anyContains(diags, "git guard inlined resolver order drifts from .bench/lib/resolve-bench.sh") {
		t.Fatalf("swapped candidate order did not produce the drift diagnostic: %v", diags)
	}
}

func TestCheckGuardResolverOrderDriftAcceptsMatchingOrder(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, filepath.Join(root, ".bench", "hooks", "block-dangerous-git.sh"), `#!/usr/bin/env bash
resolve_wrapper() {
  local root candidate
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  if [[ -n "$root" ]]; then
    for candidate in "$root/.bench/bin/bench.sh" "$root/bin/bench.sh"; do
      [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
    done
  fi
  command -v bench 2>/dev/null || return 1
}
`)
	writeFixtureFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), canonicalResolveBenchLib)

	diags := checkGuardResolverOrderDrift(root)
	if len(diags) != 0 {
		t.Fatalf("matching candidate order produced unexpected diagnostics: %v", diags)
	}
}

// guardShimRelPath and guardResolveLibRelPath anchor checkGuardResolverOrderDrift:
// the git guard shim deliberately inlines its wrapper resolver rather than
// sourcing the shared lib (sourcing would add a fail-open mode), so this check
// is the only thing keeping the two search orders honest.
const (
	guardShimRelPath       = ".bench/hooks/block-dangerous-git.sh"
	guardResolveLibRelPath = ".bench/lib/resolve-bench.sh"
	guardShimResolverFunc  = "resolve_wrapper"
	guardLibResolverFunc   = "bench_resolve_wrapper"
)

// checkGuardResolverOrderDrift compares the ordered wrapper-search candidates the
// git guard shim's inlined resolve_wrapper() tries against the shared lib's
// bench_resolve_wrapper(): repo wrapper (.bench/bin/bench.sh) -> kit wrapper
// (bin/bench.sh) -> PATH fallback (command -v bench). It reds on any order
// mismatch, and reds honestly (never a vacuous pass) when a file, a resolver
// function, or an expected candidate can't be found — so renaming the resolver
// can't amputate the check.
func checkGuardResolverOrderDrift(root string) []string {
	shimText, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(guardShimRelPath)))
	if err != nil {
		return []string{fmt.Sprintf("git guard resolver order drift check: %s is missing; cannot verify wrapper search order", guardShimRelPath)}
	}
	libText, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(guardResolveLibRelPath)))
	if err != nil {
		return []string{fmt.Sprintf("git guard resolver order drift check: %s is missing; cannot verify wrapper search order", guardResolveLibRelPath)}
	}

	shimBody, ok := extractResolverFunctionBody(string(shimText), guardShimResolverFunc)
	if !ok {
		return []string{fmt.Sprintf("git guard resolver order drift check: %s has no %s() function to anchor on", guardShimRelPath, guardShimResolverFunc)}
	}
	libBody, ok := extractResolverFunctionBody(string(libText), guardLibResolverFunc)
	if !ok {
		return []string{fmt.Sprintf("git guard resolver order drift check: %s has no %s() function to anchor on", guardResolveLibRelPath, guardLibResolverFunc)}
	}

	shimOrder := extractResolverCandidateOrder(shimBody)
	if len(shimOrder) != 3 {
		return []string{fmt.Sprintf("git guard resolver order drift check: %s resolver is missing an expected wrapper-search candidate (found order %v)", guardShimRelPath, shimOrder)}
	}
	libOrder := extractResolverCandidateOrder(libBody)
	if len(libOrder) != 3 {
		return []string{fmt.Sprintf("git guard resolver order drift check: %s resolver is missing an expected wrapper-search candidate (found order %v)", guardResolveLibRelPath, libOrder)}
	}

	if !slices.Equal(shimOrder, libOrder) {
		return []string{fmt.Sprintf("git guard inlined resolver order drifts from %s (shim order %v, lib order %v)", guardResolveLibRelPath, shimOrder, libOrder)}
	}
	return nil
}

// extractResolverFunctionBody returns the text between a `name() {` signature
// (anchored to its own line) and the next line that is a bare `}` closing brace.
// Both resolvers in this repo close on a line of their own, so this doesn't need
// to track brace nesting for the inline `{ ...; }` group inside the shim's loop.
func extractResolverFunctionBody(text, funcName string) (string, bool) {
	sigRe := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(funcName) + `\s*\(\)\s*\{\s*$`)
	sigLoc := sigRe.FindStringIndex(text)
	if sigLoc == nil {
		return "", false
	}
	rest := text[sigLoc[1]:]
	closeRe := regexp.MustCompile(`(?m)^\}\s*$`)
	closeLoc := closeRe.FindStringIndex(rest)
	if closeLoc == nil {
		return "", false
	}
	return rest[:closeLoc[0]], true
}

var (
	guardResolverRepoWrapperRe = regexp.MustCompile(`\.bench/bin/bench\.sh`)
	guardResolverKitWrapperRe  = regexp.MustCompile(`bin/bench\.sh`)
	guardResolverPathFallback  = regexp.MustCompile(`command -v bench`)
)

// extractResolverCandidateOrder returns the labels of the three known wrapper
// candidates in the order they first appear in a resolver function body. Note
// "bin/bench.sh" is a substring of ".bench/bin/bench.sh": the repo-wrapper
// occurrences are masked out before searching for the kit-wrapper token so the
// two can't be confused.
func extractResolverCandidateOrder(body string) []string {
	type token struct {
		label string
		pos   int
	}
	var tokens []token

	if loc := guardResolverRepoWrapperRe.FindStringIndex(body); loc != nil {
		tokens = append(tokens, token{"repo-wrapper", loc[0]})
	}

	masked := guardResolverRepoWrapperRe.ReplaceAllStringFunc(body, func(s string) string {
		return strings.Repeat("#", len(s))
	})
	if loc := guardResolverKitWrapperRe.FindStringIndex(masked); loc != nil {
		tokens = append(tokens, token{"kit-wrapper", loc[0]})
	}

	if loc := guardResolverPathFallback.FindStringIndex(body); loc != nil {
		tokens = append(tokens, token{"path-fallback", loc[0]})
	}

	sort.Slice(tokens, func(i, j int) bool { return tokens[i].pos < tokens[j].pos })
	order := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		order = append(order, tok.label)
	}
	return order
}
