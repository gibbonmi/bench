# FT183 #2 research — observing which derivation a registry row resolves through

Produced 2026-08-03 by two read-only research delegates (mid tier); supports the
FT183 map's decision ticket #2. Citations are to the tree at `843a4b7`.

## The gap, made concrete

A registry row is `{component, source Source, resolve func(*inputResolver)(...)}`
(`internal/gate/component_inputs.go:87-91`); the source label and the resolver
function sit side by side in a literal with nothing binding them. The exported
surfaces (`ResolveComponentInputs`, `ComponentInputSources`) drop the resolver
value, so `internal/conformance/derivation_source_test.go` can only grade
`(paths, digests)` predicates — its header discloses this (`:26-27`).

Two swaps were run in a scratch copy of `HEAD` and verified to pass the
conformance check:

- **Swap A**: vet's resolver changed from `(*inputResolver).moduleClosure` to
  `(*inputResolver).contractInputs` (`component_inputs.go:104`) — conformance
  green. It was caught only by an unrelated behavioral expectation
  (`internal/gate/component_decision_test.go:340`, vet must skip on an
  agent-Markdown edit); a future registry entry carries no equivalent.
- **Swap B**: the `canary` row's resolver changed from `canaryInputs` to
  `shellcheckInputs` (`component_inputs.go:109`) — conformance green. The
  `SourceHandDeclared` branch (`derivation_source_test.go:163-170`) exempts the
  row from all of the derivation-source check's grading. Corrected 2026-08-03
  after doc review: "green everywhere" overclaimed — this specific swap goes
  red in `internal/gate/component_decision_test.go` (`TestAgentMarkdownEdit-
  RunsConsumersAndSkipsToolchain`, canary must run on an agent-Markdown edit,
  and `shellcheckInputs` drops the agent-Markdown roots). Same posture as
  Swap A: the resolver is unobserved by the derivation-source check, caught
  only by a hand-written behavioral expectation no future row is guaranteed
  to carry.

Derivation outputs nest (`buildClosure` ⊂ `moduleClosure` ⊂ `contractInputs`),
which is why superset swaps pass the "gained ≥1 path" clause
(`derivation_source_test.go:180-182`).

## Candidate mechanisms

| | Proves | Survives behaviorally-identical functions | Production change | Test cost | Main fragility |
|---|---|---|---|---|---|
| **A** function identity (`reflect.ValueOf(fn).Pointer()`) | label↔function, exactly | yes | none if in-package (`component_inputs_test.go:181-192` already touches `resolve`); ~5 lines of new export if in `internal/conformance` | ~35 lines | depends on rows staying method *expressions* — method values/closures break pointer equality silently; needs a `Source → function` table |
| **B** instrumented derivations | named derivation fired at runtime | yes | worst: every derivation instrumented + per-row resolver seam; no repo precedent for test-only production hooks; memoization (`:159-191`) and `contractInputs`→`moduleClosure` composition force fresh-resolver-per-row | medium | forget-to-instrument is silent green; needs the same table |
| **C** per-derivation perturbation grid | behavioral shape only | no | none (pure `internal/conformance` test) | ~40 lines | nesting means `build` is separable only by *negative* membership; reintroduces the per-component table the design deliberately avoided (`derivation_source_test.go:49-51`) |
| **D** AST bijection over the registry literal | textual pairing; bijection variant needs no table | yes (textual) | none; precedented pattern (`internal/conformance/cross_compile_default_test.go:23-53`) | ~100 lines | breaks on any registry-literal refactor; residual escape: a swap to a function no other row uses |
| **E** extensional equality vs the *named* source | full-set behavioral match | no | none; generalizes an existing pattern (build already pinned to `freshness.BuildInputs` at `component_inputs_test.go:87-105`, shellcheck to `shellcheckFiles` at `:213-222`) | low-medium | in-package only (named sources unexported); couples test to derivation internals |

Cross-cutting: mechanisms A, B, and C each need a `Source → expected derivation`
table. This table restates a fact the registry already carries. The
one-source-per-fact standard admits this restatement only under the
independently-authored-expectation exception. That exception applies because
the table's independence is what turns a named swap red, and someone must
record and demonstrate that red. D's bijection variant and E avoid the table.

Existing seams a check can reuse: in-package registry enumeration and direct
`resolve` invocation (`component_inputs_test.go:181-192`); the AST-tripwire
pattern in `internal/conformance`; the label↔prose binding checker
(`component_scope_binding_test.go:47-86`). No function-identity observation
exists anywhere in the tree today.
