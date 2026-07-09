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

// writeShim materializes a git-guard shim fixture from a resolver function name,
// an optional in-body comment line, and the ordered (already-quoted) wrapper
// candidates its `for ... in` loop tries. Every shim fixture in this file derives
// from it so the fixtures can't drift from one another the way three hand-pasted
// heredocs would.
func writeShim(t *testing.T, root, funcName, bodyComment string, candidates []string) {
	t.Helper()
	var b strings.Builder
	b.WriteString("#!/usr/bin/env bash\n")
	b.WriteString(funcName + "() {\n")
	b.WriteString("  local root candidate\n")
	b.WriteString("  root=\"$(git rev-parse --show-toplevel 2>/dev/null || true)\"\n")
	b.WriteString("  if [[ -n \"$root\" ]]; then\n")
	if bodyComment != "" {
		b.WriteString("    " + bodyComment + "\n")
	}
	b.WriteString("    for candidate in " + strings.Join(candidates, " ") + "; do\n")
	b.WriteString("      [[ -x \"$candidate\" ]] && { printf '%s\\n' \"$candidate\"; return 0; }\n")
	b.WriteString("    done\n")
	b.WriteString("  fi\n")
	b.WriteString("  command -v bench 2>/dev/null || return 1\n")
	b.WriteString("}\n")
	writeFixtureFile(t, filepath.Join(root, ".bench", "hooks", "block-dangerous-git.sh"), b.String())
}

func TestCheckGuardResolverOrderDriftAnchorsOnRenamedFunction(t *testing.T) {
	root := t.TempDir()
	writeShim(t, root, "resolve_wrapper_v2", "", []string{
		`"$root/.bench/bin/bench.sh"`,
		`"$root/bin/bench.sh"`,
	})
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
	writeShim(t, root, "resolve_wrapper", "", []string{
		`"$root/bin/bench.sh"`,
		`"$root/.bench/bin/bench.sh"`,
	})
	writeFixtureFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), canonicalResolveBenchLib)

	diags := checkGuardResolverOrderDrift(root)
	if !anyContains(diags, "git guard inlined resolver order drifts from .bench/lib/resolve-bench.sh") {
		t.Fatalf("swapped candidate order did not produce the drift diagnostic: %v", diags)
	}
}

func TestCheckGuardResolverOrderDriftAcceptsMatchingOrder(t *testing.T) {
	root := t.TempDir()
	writeShim(t, root, "resolve_wrapper", "", []string{
		`"$root/.bench/bin/bench.sh"`,
		`"$root/bin/bench.sh"`,
	})
	writeFixtureFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), canonicalResolveBenchLib)

	diags := checkGuardResolverOrderDrift(root)
	if len(diags) != 0 {
		t.Fatalf("matching candidate order produced unexpected diagnostics: %v", diags)
	}
}

// TestCheckGuardResolverOrderDriftRedsOnForeignCandidate covers the drift class
// the check exists to catch: a NEW search candidate inserted into the shim loop.
// A foreign path whose tail is "bin/bench.sh" must not alias the kit-wrapper
// token — it must red loudly, naming the offending candidate.
func TestCheckGuardResolverOrderDriftRedsOnForeignCandidate(t *testing.T) {
	root := t.TempDir()
	writeShim(t, root, "resolve_wrapper", "", []string{
		`"$root/.bench/bin/bench.sh"`,
		`"$HOME/.local/bin/bench.sh"`,
		`"$root/bin/bench.sh"`,
	})
	writeFixtureFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), canonicalResolveBenchLib)

	diags := checkGuardResolverOrderDrift(root)
	if !anyContains(diags, "unrecognized wrapper-search candidate") {
		t.Fatalf("foreign wrapper candidate did not produce the unrecognized-candidate diagnostic: %v", diags)
	}
	if !anyContains(diags, "$HOME/.local/bin/bench.sh") {
		t.Fatalf("unrecognized-candidate diagnostic did not name the offending token: %v", diags)
	}
}

// TestCheckGuardResolverOrderDriftIgnoresBodyComments proves candidate extraction
// reads the operative `for` list, not in-body comments: a canonical-order comment
// above a swapped loop must still red, and a comment above a correct loop must
// stay green.
func TestCheckGuardResolverOrderDriftIgnoresBodyComments(t *testing.T) {
	comment := "# tries .bench/bin/bench.sh then bin/bench.sh"

	t.Run("comment above swapped loop still reds", func(t *testing.T) {
		root := t.TempDir()
		writeShim(t, root, "resolve_wrapper", comment, []string{
			`"$root/bin/bench.sh"`,
			`"$root/.bench/bin/bench.sh"`,
		})
		writeFixtureFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), canonicalResolveBenchLib)

		diags := checkGuardResolverOrderDrift(root)
		if !anyContains(diags, "git guard inlined resolver order drifts from .bench/lib/resolve-bench.sh") {
			t.Fatalf("canonical-order body comment masked a swapped operative loop: %v", diags)
		}
	})

	t.Run("comment above correct loop stays green", func(t *testing.T) {
		root := t.TempDir()
		writeShim(t, root, "resolve_wrapper", comment, []string{
			`"$root/.bench/bin/bench.sh"`,
			`"$root/bin/bench.sh"`,
		})
		writeFixtureFile(t, filepath.Join(root, ".bench", "lib", "resolve-bench.sh"), canonicalResolveBenchLib)

		diags := checkGuardResolverOrderDrift(root)
		if len(diags) != 0 {
			t.Fatalf("body comment above a correct loop produced diagnostics: %v", diags)
		}
	})
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
// mismatch, reds on any unrecognized search candidate (a newly inserted path is
// exactly the drift class this check exists to catch), and reds honestly (never a
// vacuous pass) when a file, a resolver function, or an expected candidate can't
// be found — so renaming the resolver can't amputate the check.
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

	shimOrder, shimUnknown := extractResolverCandidateOrder(shimBody)
	if shimUnknown != "" {
		return []string{fmt.Sprintf("git guard resolver order drift check: %s resolver has an unrecognized wrapper-search candidate %q", guardShimRelPath, shimUnknown)}
	}
	if len(shimOrder) != 3 {
		return []string{fmt.Sprintf("git guard resolver order drift check: %s resolver is missing an expected wrapper-search candidate (found order %v)", guardShimRelPath, shimOrder)}
	}
	libOrder, libUnknown := extractResolverCandidateOrder(libBody)
	if libUnknown != "" {
		return []string{fmt.Sprintf("git guard resolver order drift check: %s resolver has an unrecognized wrapper-search candidate %q", guardResolveLibRelPath, libUnknown)}
	}
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
	guardResolverForLoopRe     = regexp.MustCompile(`(?m)^[[:blank:]]*for[[:blank:]]+[A-Za-z_][A-Za-z0-9_]*[[:blank:]]+in[[:blank:]]+(.+)$`)
	guardResolverQuotedWordRe  = regexp.MustCompile(`"[^"]*"`)
	guardResolverRepoWrapperRe = regexp.MustCompile(`^\$[A-Za-z_][A-Za-z0-9_]*/\.bench/bin/bench\.sh$`)
	guardResolverKitWrapperRe  = regexp.MustCompile(`^\$[A-Za-z_][A-Za-z0-9_]*/bin/bench\.sh$`)
	guardResolverPathFallback  = regexp.MustCompile(`command -v bench`)
)

// extractResolverCandidateOrder returns the operative wrapper-search candidates
// in the order the resolver body tries them, plus the first candidate token it
// does not recognize (empty when all are known). Comment lines are stripped
// before matching so a stale in-body comment can't masquerade as an operative
// line. Candidates come from the `for ... in <list>` word list (each quoted word
// classified by path shape) and the `command -v bench` PATH fallback. Any loop
// candidate that is neither the repo wrapper ($root/.bench/bin/bench.sh) nor the
// kit wrapper ($root/bin/bench.sh) is reported as unrecognized, so a newly
// inserted or substring-aliasing search path reds instead of passing green.
func extractResolverCandidateOrder(body string) (order []string, unrecognized string) {
	stripped := stripShellCommentLines(body)

	type token struct {
		label string
		pos   int
	}
	var tokens []token

	for _, loop := range guardResolverForLoopRe.FindAllStringSubmatchIndex(stripped, -1) {
		listStart, listEnd := loop[2], loop[3]
		if listStart < 0 {
			continue
		}
		list := stripped[listStart:listEnd]
		for _, wLoc := range guardResolverQuotedWordRe.FindAllStringIndex(list, -1) {
			word := list[wLoc[0]:wLoc[1]]
			inner := word[1 : len(word)-1] // strip the surrounding quotes
			pos := listStart + wLoc[0]
			switch {
			case guardResolverRepoWrapperRe.MatchString(inner):
				tokens = append(tokens, token{"repo-wrapper", pos})
			case guardResolverKitWrapperRe.MatchString(inner):
				tokens = append(tokens, token{"kit-wrapper", pos})
			default:
				return nil, inner
			}
		}
	}

	if loc := guardResolverPathFallback.FindStringIndex(stripped); loc != nil {
		tokens = append(tokens, token{"path-fallback", loc[0]})
	}

	sort.Slice(tokens, func(i, j int) bool { return tokens[i].pos < tokens[j].pos })
	order = make([]string, 0, len(tokens))
	for _, tok := range tokens {
		order = append(order, tok.label)
	}
	return order, ""
}

// stripShellCommentLines drops every line whose first non-blank character is `#`,
// so a comment can never be read as an operative resolver line.
func stripShellCommentLines(body string) string {
	lines := strings.Split(body, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimLeft(line, " \t"), "#") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}
