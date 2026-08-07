# Gate decision test seam

Status: staged

Decision source: reviewer-confirmed current conversation on 2026-08-06.

## Problem

The gate's exhaustive public-document mapping test proves one decision table by
running the full production engine for every mutation, deletion, addition, and
restoration. In the reviewer-confirmed focused run against `6607236`, the current
19-class fixture took 30.31 seconds under `GOMAXPROCS=2 GOFLAGS=-p=2`; the earlier
21-row form took 35.14 seconds and one representative row launched 415 Git children
while materializing the same fixture tree 47 times. Most of that work does not
exercise a different policy answer. It repeatedly crosses source capture, phase
execution, output projection, and verdict publication to compare the published
check partition with the same decision object that supplied the expected
partition.

That makes the expensive path weak in the place where it claims exhaustive
coverage. A production mapping omission can disappear from both sides of the
comparison, while a projection defect pays the full engine cost for every table
row even though one representative composition test can expose it.

## Solution

Exercise the complete public-document component/check table directly through the
existing read-only gate-decision seam. An independently authored table names every
declared public Markdown file and directory class and its exact expected component
and conformance-check partition. Each file is tested for mutation and absence;
each directory is tested for descendant addition and removal. A test captures one
tree generation for the changed state, asks `scopeComponentsForGeneration` for the
decision, and compares that decision with the independent expectation. It does
not launch gate phases or publish a verdict per row.

Keep the existing representative full-engine contracts that cross the composition
boundary: one mixed public-document change proves source capture, phase selection,
environment projection, operator output, durable verdict evidence, inspection,
and composed-green behavior; the decision-site failure table proves uncertain
source facts execute rather than inherit; the evaluation operation-count tests
prove one accepted generation owns materialization. These controls retain the
integration bite without multiplying it by every policy row.

Production declarations, component identities, check identities, phase selection,
slot authorization, output, and verdict schemas do not change.

## User stories

1. As a gate maintainer, I want every public-document class and change shape
   checked against an independent expected component/check decision at the
   decision seam, while representative full-engine tests retain source-wiring,
   projection, materialization, and fail-closed coverage, so that policy omissions
   red cheaply and integration defects still red at the boundary that owns them.
   Line: `gpt-5.6-terra` / medium. This is gate/conformance logic at an existing
   seam: the edit is bounded, but an incorrect expectation could weaken the
   oracle's tests.

## Implementation decisions

**Use the decision seam that already exists.**
`scopeComponentsForGeneration` is the production decision boundary. It consumes
one accepted `treeGeneration` plus resolution, mode, slot state, and time, and
returns `componentScoping`: eligibility, component identities, executed/skipped
components, check partitions and identity state, and selected phases. Calling it
does not execute a phase or author a verdict. The test may use the existing
`mustScopeComponents` fixture helper when that helper demonstrably captures only
one generation and then calls this seam. No production interface is added merely
for the test.

**Make the expected table independent.** The test expectation explicitly lists
the current public document classes rather than deriving the expected rows from
`ReducedScope`, the conformance registry, `declaredCheckInputPaths`, or
`consumer-payload.json`. The table currently covers these file classes:

- `.bench-notes.md`
- `.bench/BENCH-reference.md`
- `.bench/BENCH.md`
- `.claude/README.md`
- `CHANGELOG.md`
- `DATA_HANDLING.md`
- `README.md`
- `ROADMAP.md`
- `projects/benchkit.md`
- `projects/gl-axi.md`
- `projects/regroup.md`

It also covers these directory classes:

- `.agents/commands/`
- `.agents/skills/`
- `.bench/adapters/`
- `.bench/hooks/`
- `.bench/lib/`
- `capture/`
- `decisions/`
- `specs/`

Each row states the exact expected executed and inherited components and checks,
including check-specific ownership beyond the catch-all set. This duplication is
the repository standard's narrow test-expectation exception: before green, the
implementation must demonstrate that omitting one declared file class, one
declared directory class, one component owner, and one check owner each makes the
independent table red. The table is not exported or reused by production.
The test must also compare the literal membership set with the inventory derived
from production declarations. That set-equality check is the addition/removal
tripwire; it does not derive any row's expected component or check partition.

**One seeded history, one decision per changed state.** The fixture may run one
full green gate to seed reusable component and check evidence. After that seed,
each table case changes the fixture, captures the changed generation once, invokes
the decision seam once, and restores fixture bytes directly. Restoration does not
run another gate. A restored-state decision check may be shared across each
change-shape family when needed to prove the fixture did not corrupt its baseline.
The exhaustive loop contains no `productionGateEngine`, `observeGreenGate`, gate
subprocess, phase recorder, verdict write, or forced-green restoration call.

**Assert the whole policy answer, not incidental identities.** Exact assertions
cover eligibility, executed and inherited reusable components, executed and
inherited conformance checks, and any narrowed phase/check-selection projection
present in `componentScoping`. Content hashes and timestamps are inputs to the
answer, not stable expectations. Ordered public forms use their canonical order;
map-backed private fields are normalized only for comparison.

**Retain representative composition owners.**
`TestMixedCheckPartitionProjectsExactOutputAndVerdict` remains the full-engine
proof that a public-document decision reaches the phase environment, operator
output, durable check evidence, `Inspect`, and `ComposedGreen` through one
aggregate conformance invocation. `TestDecisionSiteFailsClosed` remains the
full-engine proof that unavailable or uncertain decision inputs cannot authorize
inheritance. The generation/materialization operation-count contracts remain the
source-wiring bound. They may be tightened only if needed to state the existing
contract; they are not cloned into the matrix.

**Keep one ownership fence.** The implementation outcome owns the exhaustive
matrix and its local fixture helpers in `internal/gate/check_slots_test.go`.
Existing representative tests in other gate test files are regression controls,
not new writers. Production files are outside the fence.

Ownership fence: `internal/gate/check_slots_test.go`.

## Testing decisions

- The policy seam is `scopeComponentsForGeneration`, driven from a kit-shaped
  fixture with one seeded reusable-green history and one captured generation per
  changed state.
- The expected mapping is literal and independently authored. Deriving `want`
  from the returned `componentScoping`, the production registries, reduced scope,
  declared input resolvers, or consumer payload is a failing implementation.
- The matrix is exhaustive over the 11 current file classes and eight current
  directory classes. Its row-count and membership assertions make additions or
  removals in either production declaration surface require an explicit reviewed
  test expectation.
- The matrix's process boundary is structural and measurable: after the single
  seed, its cases invoke the decision seam only. A recorder or operation counter
  proves there are no gate-engine executions or verdict publications in the loop
  and at most one tree-generation capture per changed state.
- The full-engine seam is the existing mixed-partition projection test. It remains
  one representative public-document mutation, because its job is composition,
  not table enumeration.
- The fail-closed seam is the existing decision-site failure table; the source
  materialization seam is the existing evaluation operation-count suite.
- The feature gate is `bench gate`.

### Seam diagram: exhaustive policy decision

```text
trigger: public-document mapping matrix
    |
    v
changed fixture --> [ capture one tree generation ] --> [ gate decision seam ]
                                                               |
                                                               v
                                                   componentScoping answer
                                                               ^
                                                               |
tests attach: independent class table compares exact component/check partitions;
              recorder proves no phase execution or verdict publication per row
```

### Seam diagram: representative full-engine composition

```text
trigger: one ROADMAP.md mutation after a seeded green
    |
    v
working tree --> [ source capture -> decision -> selected phases ] --> phase env
                                                                    + output
                                                                    + verdict evidence
                                                                    + inspection
                                                                    + composed green
                     ^ existing mixed-partition full-engine test attaches here
```

These seams do not collapse into one test. The decision seam catches a missing or
wrong mapping without orchestration noise. The full-engine seam catches a correct
decision that is lost or distorted while materializing phases, environment,
output, or durable evidence. Either degenerate implementation can pass one seam
and must fail the other.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | the independent matrix contains exactly the 11 current file classes and eight current directory classes named above | decision seam | TDD-able: replace the production-derived inventory in the current test with the literal expectation, then demonstrate that omitting one declared file and one declared directory from production each leaves a mismatch before restoring them | A table derived from the same declarations would silently shrink with an omission; independent membership cannot. |
| 1 | mutating every declared file class yields its exact executed/inherited component and check partitions | decision seam | TDD-able: the present test compares the engine record with the decision object's own partition and has no independent expected component/check table | A wrong owner or catch-all classification differs from the literal row even when production and projection agree with each other. |
| 1 | deleting every declared file class yields its exact fail-closed or scoped partition, as specified by that row | decision seam | TDD-able: each current delete case crosses the full engine but still computes `want` from the same decision | Absence can follow a different identity/error path from mutation; a mutation-only implementation cannot satisfy both rows. |
| 1 | adding a Markdown descendant beneath every declared directory class yields its exact component/check partition | decision seam | TDD-able: the current directory loop derives its expected partition from the returned decision | A resolver that watches only existing descendants or misses one directory owner fails the independently named addition case. |
| 1 | removing an existing Markdown descendant beneath every declared directory class yields its exact component/check partition | decision seam | TDD-able: demonstrate one removed directory owner makes the literal row red before restoring production | Addition-only or existence-only directory tracking cannot satisfy the removal row. |
| 1 | each changed-state row captures at most one generation and invokes no full gate, phase, or verdict publication after the single seed | decision seam operation recorder | Observed red: the current 30.31-second test calls `observeGreenGate` for each changed state and a forced production engine after each restoration | A rewrite that keeps hidden engine calls or repeated source materialization violates the count even if wall time happens to be fast. |
| 1 | restoring fixture bytes returns the decision to the seeded partition without running a restorative gate | decision seam | TDD-able: assert the restored decision and unchanged slot-store bytes after each change-shape family | A case that mutates evidence, leaks generated state, or depends on a forced-green reset cannot return to the same answer. |
| 1 | one mixed public-document change reaches selected checks, inherited checks, phase environment, operator output, durable evidence, inspection, and composed green | full-engine composition seam | Already covered by `TestMixedCheckPartitionProjectsExactOutputAndVerdict`, retained unchanged as a regression control | Direct decision tests cannot catch a correct answer dropped between the decision and its public/durable projections. |
| 1 | uncertainty in decision source capture or identity resolution executes affected work and never authorizes inheritance | full-engine fail-closed seam | Already covered by `TestDecisionSiteFailsClosed`, retained as a regression control | A fast decision test over valid generations alone cannot prove the engine refuses unsafe reuse when its source is unavailable. |
| 1 | one evaluation generation bounds tree materialization, listings, and blob reads | evaluation source seam | Already covered by the gate-evaluation operation-count controls, retained as regression controls | Moving the matrix cannot conceal a regression that rematerializes the same accepted source inside one production evaluation. |
| 1 | production phase selection, check selection, output, slot records, and verdict schema are unchanged | package and feature gate | Already covered by the existing `internal/gate` suite and `bench gate`; the implementation fence contains no production file | A test-only speedup cannot earn green by weakening or changing the oracle it observes. |

The cheapest wrong implementations are explicit. Computing `want` from
`componentScoping` or any production declaration fails the independent membership
and owner-omission demonstrations. Moving only file cases leaves the directory
addition/removal rows red. Keeping a hidden full gate behind a helper fails the
operation count. Deleting all engine coverage fails the mixed-projection and
fail-closed rows. Adding a test-only production interface violates the ownership
fence and the no-production-change row.

### Edge inventory

- Error path — resolved by the file-deletion rows and the retained full-engine
  decision-site failure table; unsafe source or identity errors execute rather
  than inherit.
- Empty or absent input — absence is resolved for every file class by deletion
  and for every directory class by descendant removal. A present empty Markdown
  file is not a separate mapping state because identities are content-sensitive
  but ownership is path-class-sensitive; ordinary mutation already changes its
  identity without changing its owner set.
- Boundary values — resolved by exact matrix cardinality, the first/last ordered
  component and check expectations, and the single-seed/one-capture operation
  bounds.
- Malformed input — **Won't handle**: this test movement introduces no parser or
  serialized input. Malformed tree/blob sources remain owned by the retained
  capture and fail-closed controls.
- Interrupted or partial state — **Won't handle**: the decision seam is read-only
  and adds no process or transaction. Existing engine cancellation and durable
  verdict tests continue to own interruption behavior.
- Re-run idempotency — resolved by restored-state decisions and unchanged
  slot-store bytes without a restorative gate.
- Process lifecycle — resolved by the structural no-engine/no-phase/no-verdict
  count after the one seed; the matrix creates no child gate lifecycle to reap.
- Hostile environment — **Won't handle** at this seam: no environment variable,
  shell argv, terminal, network, or external path enters the mapping table.
  Existing working-tree capture tests own spaces, glob characters, control bytes,
  symlinks, special files, non-ASCII names, and missing tools.
- A command whose write changes a fact it reports — **Won't handle**: no command
  or production write is added. The fixture mutation is test setup, and slot-store
  immutability is asserted explicitly.
- Host-backed filesystem I/O pressure — resolved within scope by eliminating all
  post-seed engine executions and bounding generation capture per row. Absolute
  disk throughput and cache warmth remain machine-dependent and are not used as
  the oracle.

## Out of scope

- Scoping each conformance fixture to its registered check is
  `conformance-harness-scope`, blocked on this lifecycle — approximately 5 edits
  and 3 gate runs across the harness, fixture owners, and representative shell
  controls.
- Re-running the cold workload census after demand reduction is decision-map #20
  — 0 product edits and at least 3 instrumented gate runs so cold, warm, and
  representative compositions are not conflated.
- Adding a gate-wide memory/concurrency admission budget is decision-map #8,
  downstream of that census — approximately 10 edits and 5 gate runs across pool
  ownership, phase weights, nested-run propagation, diagnostics, and canaries.
- Changing component/check declarations, reuse authorization, phase composition,
  stripped-subject ordering, canary scheduling, or verdict schemas is a production
  oracle capability — at least 6 edits and 3 gate runs even for a bounded policy
  change. This spec deliberately leaves all of it unchanged.
- Removing or weakening checks to improve runtime is not an alternative
  capability. The source requires retaining the oracle and moving exhaustive
  assertions to the seam that owns them.
