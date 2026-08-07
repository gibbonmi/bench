# Conformance harness scope

Status: implemented

Decision source: reviewer-confirmed current conversation on 2026-08-07.

## Problem

Bench's conformance canary registry already binds each conformance fixture family
to the one executable check whose diagnostic the family grades, and the
conformance runner already accepts that singular scope. The direct Go fixture-bite
tests ignore both seams: most materialize one mutated fixture and call the complete
dev conformance table. The 2026-08-04 census measured the largest group at 83
subtests and 35.14 seconds because every mutation paid unrelated checks. A wrong or
missing fixture binding can also hide behind that broad run because the expected
diagnostic may appear even when the harness never used the registered owner.

The same package repeats a second over-wide pattern at the shell gate entry. Eight
freshness artifact variants each create a temporary kit module and invoke the shell
entry, whose freshness verifier starts through `go run`. Those variants mostly
exercise artifact verification and freshness-subcommand refusal owned by the
freshness package, not distinct shell routing. The repeated shell setup makes the
test expensive without adding eight independent process-boundary claims.

## Solution

Run every direct conformance fixture bite through the resolved `Fixture.Check`
exported by the canary fixture inventory. Canary alone resolves a fixture-level
`CHECK` over its family's fallback with its shared registry rule; conformance
consumes that decision and calls the existing singular conformance scope. Missing,
unknown, meta, or wrong-tier resolved checks and an empty expectation fail loudly;
no failure falls back to the full table or to an empty successful run. Each fixture
remains a distinct mutated tree with its own expectation, while the timing record
proves that exactly the resolved check executed.

Move the exhaustive freshness-artifact state table to the freshness package's
`Verify` and `Check` seams. Keep representative shell-entry tests for both sides of
the composition: an untrusted artifact refuses before phases with one copy-paste
rebuild action, and a verified legacy binary whose freshness subcommand refuses is
rejected until a freshly published replacement runs phases exactly once. Retain the
static gate-entry ordering check. Do not change production conformance selection,
fixture registration, freshness verification, gate routing, diagnostics, or green
semantics.

## User stories

1. As a gate maintainer, I want every direct conformance fixture bite to execute
   only the resolved check its canary fixture exports, so that each mutation keeps
   its diagnostic bite without regenerating unrelated conformance work.
   Line: `gpt-5.6-terra` / medium. This is test-only gate/conformance work at an
   existing scoped seam, but a false-green fixture would weaken the oracle's
   self-defense.
2. As a gate maintainer, I want freshness artifact variants proved at the
   freshness seam and only representative variants repeated through the shell
   entry, so that the package retains process-boundary routing coverage without
   rebuilding a temporary Go module for every lower-seam state.
   Line: `gpt-5.6-terra` / medium. The seams are existing and gate-observable, but
   the lower and shell assertions must compose without dropping a refusal class.

## Implementation decisions

**Use the canary resolver and singular runner scope.** Canary owns fixture policy:
its shared rule resolves a fixture-level `CHECK` before family fallback, using
`registry.FamilyCheck` only inside canary, then exports the decision as
`Fixture.Check`. This includes overrides such as
`default-branch-refabricated`, which must retain its own resolved check instead of
inheriting its family's fallback. `RunConformance`'s non-empty scope is the execution
seam; conformance consumes `Fixture.Check` and does not call `registry.FamilyCheck`
to reconstruct policy. A test-only resolver accepts the exported check as a function
dependency so its failure postures can be exercised without making production policy
mutable; every real fixture journey supplies `Fixture.Check` directly. The helper
refuses an undiscovered fixture, an unresolved check, an unknown check, a meta check,
a check unavailable at the dev tier, or an empty `EXPECT` before executing anything.
An implementation-time overlay mutation rebinds one existing family to another valid
dev check in the registry source and proves that a real family-fallback fixture's
exported `Fixture.Check` follows the live binding. A private family-to-check table
therefore fails even when it happens to match today's registry.

**Scope direct bite tests, not full-table controls.** The fixture-bite groups for
load validity, skills/index adapters, documentation/workflow guidance, coverage
maps, line routing, package/core guards, data handling, decision maps, and example
agreement use the scoped helper wherever they materialize a canary fixture and ask
for its expected diagnostic. Tests whose subject is the
full conformance table itself remain unscoped: registry completeness, absent or
unbound family diagnostics, ordered-set selection, timing order, and the public
`TestRootConformance` entry.

The registered conformance families currently reached by those fixture-bite owners
are `load-validity-metadata`, `skills-index-command-adapters`,
`docs-currency-token-diet`, `workflow-guidance-anchors`,
`coverage-map-validation`, `line-routing`, `package-core-guard`,
`data-handling-derivation`, `decision-map-integrity`, and `example-agreement`. The
helper derives this membership from the discovered fixtures; this list fixes the
review quantifier but is not executable policy. `injected-ports` remains a
canary-sweep-owned family and has no direct `RunConformance` fixture-bite journey to
migrate.

**Prove the executed identity structurally.** A scoped bite still asserts its
fixture's independently authored `EXPECT` diagnostic. It also reads the
conformance timing record and asserts that the run recorded exactly the registered
check, ignoring the duration. This makes an unscoped helper, a hard-coded check, an
empty run, or a merged group observable without using wall-clock thresholds.
Multiple fixtures bound to one check still run once per mutated tree; the helper
does not batch or merge their subjects.

**Keep two ownership fences.** Fixture scoping owns only the conformance test
harness and its direct fixture-bite callers:
`internal/conformance/fixture_bite_test.go`,
`internal/conformance/data_handling_test.go`,
`internal/conformance/decision_map_integrity_test.go`, and
`internal/conformance/example_agreement_test.go`. Freshness test movement owns only
`internal/freshness/freshness_test.go` and
`internal/conformance/gate_entry_test.go`. Production files, the executable check
registry, canary fixture contents, and gate scripts are outside both fences.

**Split freshness proof by its real owner.** Artifact and seal absence,
malformation, unreadability, executable-digest mismatch, and verified
freshness-subcommand failure are enumerated at `freshness.Verify` or
`freshness.Check`, whichever first owns the refusal. The lower tests assert the
same untrusted-binary diagnostic class and exactly one copy-paste-safe rebuild
action. Existing symlink, special-file, hostile-path, content-digest, and
publication controls remain their independent owners rather than being copied into
the moved table.

**Retain representative shell composition.** The shell gate entry continues to
prove current-source verification happens before `gate-phases`; a missing or
untrusted artifact emits one rebuild action and never schedules phases; a verified
legacy binary whose `freshness-check` refuses is not trusted; and a published
replacement reaches phases exactly once per invocation. The representative entry
tests run from a nested root containing spaces and glob characters. No elapsed-time
assertion defines acceptance.

This is an ordinary test-only build, not a wide refactor. `craft-tickets` may slice
the two disjoint ownership fences into independently green tracer tickets. The
fixture ticket must keep helper and migrated callers together because a helper with
no callers leaves the over-wide tests unchanged, while caller edits without the
helper do not compile. The freshness ticket keeps the lower-seam additions and
shell-matrix contraction together because removing the broad shell table before
its refusal classes land at the lower seam strands coverage red.

## Testing decisions

- The fixture execution seam is `RunConformance` with the exact non-empty scope from
  the discovered canary `Fixture.Check`; canary alone resolves fixture `CHECK` over
  family fallback, and conformance does not reconstruct that policy.
- The fixture tests attach above materialization: each drives a real mutated tree,
  asserts its existing `EXPECT`, and checks the timing record for exactly one
  executed registered check.
- The freshness state seam is `Verify` for artifact/seal trust and `Check` for the
  verified binary's freshness-subcommand result.
- The process seam is the real shell gate entry, retained for one untrusted-artifact
  journey and the verified-legacy-to-replacement journey.
- Existing registry/meta, singular-scope, ordered-selection, hostile artifact,
  publication, and gate-entry ordering tests remain regression controls.
- Both stories use TDD at their named seams. Add one failing behavior row at a time;
  do not weaken, delete, or rewrite an existing expectation to obtain green.
- The feature gate is `bench spec build promote conformance-harness-scope`; promotion
  is the sole whole-project gate and implemented-status author for this lifecycle.

### Seam diagram: registered fixture bite

```text
trigger: direct Go fixture-bite subtest
    |
    v
fixture name --> [ canary fixture inventory ] --> family + fixture CHECK
                                                   |
                                                   v
                         [ canary shared family-fallback rule ] --> Fixture.Check
                                                                           |
                                                   mutated tree ------------v
                                      [ singular RunConformance ]
                                                   |
                                      diagnostics + timing identity
                                                   ^
                                                   |
tests attach: existing EXPECT bites; exactly the registered check is recorded
```

### Seam diagram: freshness proof and shell composition

```text
trigger: freshness artifact-state table
    |
    v
root + executable --> [ Verify / Check ] --> trusted or one rebuild refusal
                            ^
                            | lower tests attach per refusal class

trigger: representative .bench/gate.sh invocation
    |
    v
nested subject --> [ go-run freshness entry ] --> [ verified bench gate-phases ]
                            |                                  |
                            +-- refusal: no phases             +-- replacement: once
                            ^
                            | shell composition tests attach here
```

The seams do not collapse. Lower freshness tests cannot prove shell ordering or that
phases stay unstarted, while shell-only variants cannot distinguish lower artifact
classes without repeatedly compiling the same entry. Likewise an `EXPECT` alone
cannot prove the harness used the registered check, and a scope identity alone
cannot prove the mutated fixture still bites.

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| CH1 | 1 | Every direct fixture-bite invocation in the ten enumerated conformance families obtains its resolved check from canary `Fixture.Check`; canary alone resolves fixture `CHECK` over family fallback, including the `default-branch-refabricated` override, and conformance does not reconstruct the policy | canary fixture discovery/resolution and the conformance consumer seam | Red mutation: in an overlay, rebind one existing family in the registry source to another valid ordinary dev check, run that family's real fixture without `CHECK`, and require its exported `Fixture.Check` and timing identity to follow the rebound check; also compare `default-branch-refabricated` with its fixture-level override before restoring the source | An empty scope, direct conformance `registry.FamilyCheck` lookup, hard-coded check, or private duplicate family table ignores the exported decision or live fallback binding; a family-only resolver loses the fixture-level override. |
| CH2 | 1 | Each mutated fixture tree runs exactly once under its resolved ordinary dev check and records no other check | singular `RunConformance` scope plus timing record | Observed in the current tree: the fixture callers pass `""`, and their timing records contain the complete dev table rather than one check | An unscoped call, empty successful run, hard-coded wrong check, or second execution changes the recorded identity/count. |
| CH3 | 1 | Every migrated fixture retains its existing independently authored `EXPECT` diagnostic | scoped fixture-bite journey | Already covered by the existing fixture-bite subtests; run each unchanged expectation through the new scope before removing the broad call | Correct scoping to a check that no longer emits the promised diagnostic goes red instead of turning a speedup into a vacuous green. |
| CH4 | 1 | An undiscovered fixture, an unbound family whose canary resolution exports no `Fixture.Check`, an unknown check, a meta check, a wrong-tier check, or an empty `EXPECT` refuses before any check runs | fixture helper consuming the exported resolved check | TDD-able: drive the test-only resolver with a missing fixture, empty and invalid exported checks, and an empty expectation; the current broad callers ignore those states or accept an empty substring | Fail-closed resolved-check and expectation validation preserve the user-visible unbound-family refusal and prevent a missing canary decision, stale registry row, or vacuous expectation from falling back to expensive or empty green execution. |
| CH5 | 1 | Multiple fixtures sharing one registered check remain separate mutated-tree runs, while timing state is cleared between runs | fixture loop plus timing writer | Already covered in subject count by the existing per-fixture subtests; add repeated same-check fixtures and assert one fresh timing identity for each | Batching fixtures can hide one mutation behind another, and uncleared timing can attribute a prior check to the current tree. |
| CH6 | 1 | Full-table registry, family-integrity, ordered-selection, timing-order, and public entry controls remain unscoped | existing conformance integration tests | Already covered by `TestRootConformance`, registry/meta controls, and singular/ordered selection tests; retain them outside the migrated helper | Scoping the oracle's own completeness controls would let the registry omit checks or families while every narrowed fixture stayed green. |
| FR1 | 2 | Missing executable; missing, unreadable, malformed complete, and malformed partial seals; and executable-digest mismatch each refuse at the lower freshness seam with one rebuild action | `freshness.Verify` | Partly covered by `TestVerifyRefusesUntrustedArtifactStates`; the current shell matrix is the only grouped owner for the remaining enumerated classes | The lower table identifies every artifact trust failure without paying shell/module setup and prevents one removed class from hiding behind a representative case. |
| FR2 | 2 | A verified executable whose `freshness-check` exits nonzero refuses through `freshness.Check`, suppresses untrusted child output, and emits one rebuild action | `freshness.Check` | TDD-able: the lower suite has no direct `Check` refusal test; today only the shell legacy-binary journey observes it | A Verify-only implementation would trust a sealed binary that cannot prove its own current command surface. |
| FR3 | 2 | The shell entry refuses an untrusted artifact before `gate-phases`, from a nested hostile path, with exactly one stable rebuild action | real `.bench/gate.sh` entry | Already covered by `TestGateEntryRefusesUnverifiedBinaryBeforeGatePhases`, retained as the representative pre-verification integration | Direct freshness calls cannot prove shell ordering, process exit normalization, or that phases were never scheduled. |
| FR4 | 2 | The shell entry rejects a verified legacy binary and, after publication of a replacement, runs replacement phases exactly once on repeated invocations | real `.bench/gate.sh` entry | Already covered by `TestGateEntryRejectsLegacyBeforeRunningOldTableAndRunsReplacementOnce`, retained as the representative verified-binary integration | This is the composition degenerate: every lower refusal can pass while the shell bypasses `Check` or schedules the stale binary. |
| FR5 | 2 | Gate-entry source continues to invoke current-source verification before its one `gate-phases` handoff | static conformance check plus shell integrations | Already covered by `checkGateEntryContract` and its fixture controls | A test-only migration cannot earn speed by deleting or reordering the production verification call. |
| OP1 | 1, 2 | The fixture table has one check execution per mutated tree and the exhaustive freshness table has no shell entry or temporary-module setup | timing identity and lower freshness seams | Observed in the 2026-08-04 census: 83 fixture subtests ran the full check table and eight freshness variants paid cold shell/Go setup | Structural operation ownership catches hidden broad work even when host timing is noisy or cached. |
| OP2 | 1, 2 | Production registry contents, selection behavior, diagnostics, freshness policy, gate routing, timing format, and verdict semantics are unchanged | ownership fences plus package and feature gate | Already covered by existing conformance, freshness, contract, canary, and gate controls; production paths are outside both fences | A demand-reduction test refactor cannot pass by weakening the oracle or changing the behavior the tests grade. |

The cheapest wrong implementations are explicit. Replacing `""` with a hard-coded
check or having conformance call `registry.FamilyCheck` passes a family-fallback
fixture today but fails CH1 when the live registry binding moves; resolving every
fixture only from its family also loses `default-branch-refabricated`'s `CHECK`
override. Calling the right check and dropping the `EXPECT` passes timing but fails
CH3. Merging every same-check fixture into one tree passes scope identity but fails
CH5. Deleting the
shell freshness matrix without moving all refusal classes fails FR1/FR2. Moving all
variants lower and deleting shell integration fails FR3/FR4. Retaining only lower
and shell tests while reordering `.bench/gate.sh` fails FR5.

### Edge inventory

- **Error path** — CH4 covers discovery and resolved-check failures after canary
  resolution; FR1 and FR2 cover artifact verification and verified-command failures;
  FR3 and FR4 cover their shell refusal projections. No failure falls back to full,
  empty, or phase-running success.
- **Empty or absent input** — an absent fixture, an unresolved canary
  `Fixture.Check`, and a present empty `EXPECT` are CH4; missing executable and seal
  are FR1. The scoped helper adds the empty-expectation refusal because the current
  reader trims to an empty substring, which otherwise matches any unrelated
  diagnostic vacuously.
- **Boundary values** — zero executed checks is CH2/CH4, one is CH2, two fixtures
  sharing one check is CH5, a fixture-level `CHECK` override and a family-fallback
  fixture are CH1, and the complete current direct-fixture family set is CH1.
  Full-tier execution stays covered separately by CH6.
- **Malformed input** — malformed complete and partial seals are FR1. Fixture files,
  mutations, and expectations introduce no new parser; their existing malformed
  JSON, frontmatter, decision-map, and coverage-map fixtures retain CH3.
- **Interrupted or partial state** — **Won't handle**: these test-only helpers add no
  production state, child lifetime, lock, or transaction. Existing freshness
  publication interruption and gate process-group controls remain unchanged.
- **Re-run idempotency** — CH5 requires timing to clear between repeated fixture
  runs; FR4 repeats the replacement shell entry and observes one phase handoff per
  invocation. Lower freshness checks are read-only.
- **Process-boundary lifecycle** — FR3 and FR4 cross the real shell and executable
  process seams. CH1–CH5 deliberately remain in-process because they grade the
  existing singular check interface rather than gate routing.
- **Hostile environment** — representative shell runs retain nested paths with
  spaces and glob characters; existing freshness tests retain quoted rebuild
  actions, symlinked artifacts and ancestors, FIFOs, sockets, unreadable artifacts,
  missing tools, and content-not-mtime verification. Ambient conformance selection
  variables cannot override the explicit scope argument.
- **A command whose write changes a fact it reports** — **Won't handle**: no command
  or production write is introduced. Timing is test evidence cleared at each run,
  not a durable authorization record.
- **Unquoted arguments, control bytes, missing trailing newline, non-TTY stdin, and
  shipped-surface routing** — **Won't handle** at the narrowed fixture seam: it adds
  no CLI parsing or text serialization. Existing gate-entry and freshness hostile
  path controls continue to grade the only shell surface this scope retains.
- **Elapsed-time threshold** — **Won't handle as an acceptance oracle**: cache warmth
  and host contention make wall time nondeterministic. OP1's executed-check and
  process-seam structure is authoritative; the downstream census measures the gain.

## Out of scope

- Re-running the post-reduction cold workload census and focused-versus-in-gate
  span probes is the next FT171 research capability: approximately 1 evidence-asset
  edit and 3 full gate measurement runs.
- Pricing and implementing the cross-process token pool, grant split, reclaim,
  reserved-headroom constant, and diagnostics remains FT171's downstream scheduler
  capability: at least 15 edits and 5 gate runs.
- Deduplicating ship-tier core, conformance, vet, race, and artifact proof changes
  release-evidence policy and remains separate shaping: at least 12 edits and 4 gate
  runs.
- Bounding FT87's unreachable-network waits is an independent network-lifecycle
  capability: approximately 5 edits and 2 gate runs.
- Changing the executable check registry, conformance selection semantics, canary
  fixture batching, gate phase scheduling, verdict reuse, or scoped-gating policy is
  excluded. Those are distinct production capabilities, not the rest of this
  test-harness reduction.
