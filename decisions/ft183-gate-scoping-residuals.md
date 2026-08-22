# FT183 per-component-scoping review residuals

Status: ready

## Destination

Dispositions for the two faces left by the shipped per-component-gate-scoping
review: the unreachable whole-changeset reduced fallback in `internal/gate`'s
component scoping, and the unbound derivation-source check in
`internal/conformance` that grades `internal/gate`'s registry.

## #1: Retire the whole-changeset reduced fallback

Blocked by: none
Type: Grill

### Question

After per-component scoping, `reducedInheritance` answers only for a kit root
with no `go.mod`. Production never has this shape (review S1; the
comment-truth half already landed at `7936b70`). Remove the path, or keep it
for linked-repo futures?

### Answer

Resolved 2026-08-02: remove it. The removal simplifies `internal/gate` and
retires `reduced_run_test`'s synthetic fixture. With #5's full retirement, it
also removes the reduced-verdict record tests and the mixed-class refusals
that pair reduced against partial. If a linked-repo shape genuinely needs a whole-changeset reduced path, add it
back against that real shape rather than carrying it as an unreachable
branch.

## #2: Observe which derivation a registry row resolves through

Blocked by: none
Type: Research

### Question

The derivation-source check proves an entry is derivation-sourced. It does
not prove the entry binds to its NAMED derivation, so a wrong-but-derived
resolver swap passes (review Sp1, disclosed in `derivation_source_test.go`'s
header). What observation design would expose which function a registry row
actually resolves through, and at what cost?

### Answer

Resolved 2026-08-03: the summary exists at
`decisions/assets/ft183-derivation-binding.md`. It prices five candidates:
function identity, instrumentation, perturbation grid, AST bijection, and
extensional-vs-named-source. It also verifies two resolver swaps empirically
against the derivation-source check. The hand-declared `canary` row's
resolver is exempt from all of that check's grading. Each swap is caught, if
at all, only by an unrelated behavioral expectation. The ruling is #3; the
exemption is #4.

## #3: Which binding mechanism, if any

Blocked by: #2
Type: Grill

### Question

Given the priced candidates in `decisions/assets/ft183-derivation-binding.md`,
which mechanism closes the label↔function gap — or is the disclosed weaker
claim kept and stated honestly?

### Answer

Resolved 2026-08-03: candidate A, function identity. An in-package
`internal/gate` test compares each row's stored resolver
(`reflect.ValueOf(fn).Pointer()`) against the derivation its source label
names. The `Source → function` table is an independently-authored expectation
whose red must be demonstrated: a recorded swap goes red, per the
one-source-per-fact exception. The check also guards the method-expression
assumption pointer equality depends on. No new export is needed: the check
lives beside `component_inputs_test.go`, which already touches `resolve`.

Amended 2026-08-03 after doc review: the check is exhaustive. It fails on
any registry row whose source label has no expectation-table entry. A newly
added row must therefore declare its binding to go green, rather than
silently reopening the ungraded-row gap.

## #4: The hand-declared exemption leaves canary's resolver ungraded

Blocked by: #2
Type: Grill

### Question

Swap B in the research asset shows that the `SourceHandDeclared` branch
exempts the `canary` row from all grading by the derivation-source check. A
resolver swap on this row is invisible to that check; an unrelated
behavioral canary expectation happens to catch the specific swap tried. Does
the chosen mechanism cover the hand-declared row too, or does the exemption
stand documented?

### Answer

Resolved 2026-08-03: bind the canary row too. The identity check covers
`SourceHandDeclared` with the same pointer comparison (its named function is
`canaryInputs`), closing the only row the derivation-source check leaves
entirely ungraded.

## #5: The orphaned Reduced verdict record class

Blocked by: #1
Type: Grill

### Question

`reducedInheritance` is the sole writer of `Reduced` verdicts. So #1's
removal orphans the persisted record fields and their readers: `inherits`,
`validateInheritance`, and the status and prep-release consumers. Retire the
record class with the path, or keep a legacy reader?

### Answer

Resolved 2026-08-03: retire it fully — fields, readers, and record class go
with the path. A legacy on-disk `Reduced` verdict reads as
no-reusable-verdict and forces a fresh gate run, the safe direction; no dead
schema survives.

## Not yet specified

## Spec-writer discretion

- Exact placement and naming of the identity check within `internal/gate`'s
  test files, and the shape of the `Source → function` expectation table.
  These stay open provided the demonstrated-red requirement, the
  method-expression guard, and the refuse-unknown-rows exhaustiveness
  survive.

## Out of scope

- Reintroducing any whole-changeset reduced path for the kit root; #1 closed
  that direction.

## Sources

- Path: `decisions/assets/ft183-derivation-binding.md`
  Supports: #2's summary and the factual premises of #3 and #4. Two resolver swaps verified to pass the derivation-source check. The asset also prices five candidate mechanisms. Produced 2026-08-03 by two read-only research delegates, and corrected after the doc review.
  Drift: re-verify if `internal/gate/component_inputs.go`'s registry shape or `internal/conformance/derivation_source_test.go` changes before the spec reads this map.
