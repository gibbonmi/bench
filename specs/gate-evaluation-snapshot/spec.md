# Gate evaluation snapshot

Status: staged

Decision source: reviewer-confirmed current conversation on 2026-08-04, ordering supported by `decisions/gate-budget.md` #15–#17 and `decisions/assets/gate-budget-cpu-wall-census.md`

## Problem

One gate evaluation repeatedly reconstructs and reparses the same Git tree while deciding what the gate must run and what evidence it may reuse. Component, conformance-check, conformance-canary, whole-subject, and stripped-subject identities each obtain local snapshots, and repeated blob reads are not memoized. The focused census recorded 47 materializations and 415 Git children for one public-document mapping row even though every decision answered for one unchanged fixture tree.

Those local snapshots are individually coherent but do not compose into one evaluation-wide fact. Besides process cost, a tree or blob read that moves between identity families can make one decision combine different generations. Any consolidation must preserve the gate's current safety properties: a changed subject still refuses, prospective execution still grades the exact supplied tree with the full inventory, identity failures still widen or refuse in the existing fail-closed direction, and operators still receive the same evidence and oracle verdicts.

## Solution

Make one internal gate evaluation own one immutable parsed tree-and-blob snapshot for its accepted pre-execution generation. Component, conformance-check, conformance-canary, whole-subject, and stripped-subject identity derivations consume that same snapshot. Blob content is memoized by object identity inside the generation, including failures, so no identity family rereads an object another family already requested.

After the gate child finishes, capture a separate post-execution generation and compare the complete subject identities exactly as today. The post generation is never an alias of the pre generation, even when both trees have the same object identity. A working-tree adapter materializes tracked, untracked, and unignored content through the current throwaway-index semantics; a prospective-tree adapter starts from the caller-supplied immutable tree object. Both adapters produce the same snapshot contract, so downstream identity policy does not branch by source kind.

The refactor changes ownership and process count, not policy. It retains current partitions, identities, evidence records, announcements, exit codes, error posture, reduced-scope behavior, exact prospective-tree behavior, and mid-run drift detection. Acceptance counts materializations and Git operations; elapsed time may be reported diagnostically but is not the primary oracle.

## User stories

1. **An ordinary gate decision answers from one coherent pre-execution generation.** On a stable working tree, component, conformance-check, conformance-canary, whole-subject, and stripped-subject identities are derived from the same immutable entry listing and blob cache. The resulting phase/check partitions, reusable-evidence decisions, identity values, and reduced-scope behavior remain semantically equivalent to the current gate.
   `Line: gpt-5.6-terra / medium.` The policy is already pinned, but consolidating every authorization identity behind one owner is correctness-bearing gate logic.

2. **A gate run still refuses when its subject drifts.** The evaluator captures a distinct post-execution generation after the child exits and compares it with the accepted pre-execution subject. A tree, oracle closure, resolution, or passlisted-environment change before execution or during the run keeps the current fail-closed outcome and cannot author green evidence.
   `Line: gpt-5.6-terra / medium.` Existing drift semantics are strong; the risk is accidentally weakening them while removing repeated captures.

3. **A prospective gate grades exactly the supplied immutable tree through the same identity contract.** The prospective adapter reads the named Git tree directly for its pre-execution generation, retains the separate checkout-based post-run drift check, runs the complete applicable inventories, and preserves exact tree/oracle evidence reuse without consulting ordinary component or check slots.
   `Line: gpt-5.6-terra / medium.` The existing prospective seam is established, but an adapter error could silently reintroduce ambient working-tree semantics.

4. **A gate evaluation has a deterministic source-cost ceiling without changing what users observe.** Identity-family breadth does not multiply tree materialization, listing, or blob reads. Source failures remain bounded and fail closed, while stdout, stderr, exit codes, durable verdicts, component/check evidence, and retained stripped evidence remain byte- or field-equivalent where those surfaces are already contractual.
   `Line: gpt-5.6-terra / medium.` The count oracle is deterministic and gate-observable, while equivalence spans several evidence consumers whose authority must not move.

## Implementation decisions

**One evaluation owns generation lifecycle.** The gate execution owner creates, distributes, and retires snapshots. Identity resolvers receive an immutable generation rather than a repository root from which they can independently recapture the tree. No resolver may construct a source adapter or retain a generation beyond its evaluation.

**A generation is an immutable source fact, not a mutable filesystem view.** It carries the source tree identity, the ordered parsed entries, a path index, and memoized blob results keyed by Git object identity. Returned entry and blob views cannot mutate shared state. The first blob failure is cached and returned consistently; retrying inside the generation could combine two source states and is forbidden.

**Working and prospective trees are two adapters to one contract.** The working-tree adapter preserves the current throwaway-index scope: seed from `HEAD` or the empty tree, add tracked plus untracked unignored working-copy content, and never touch the real index. The prospective adapter accepts only the caller-supplied valid tree object and reads it directly; it does not rematerialize that tree through a second throwaway index before identity derivation. Both adapters use the single listing parser and the same entry/blob validity rules.

**There are exactly two parsed generations in a real execution.** The accepted pre-execution generation supplies the plan and every pre-run identity family. A pre-lock candidate may be validated under the execution lock, but validation does not parse another listing or load identity blobs; a mismatch refuses before the gate child starts. The separately captured post-execution generation supplies the complete subject recheck and the stripped-subject identity used to retain or invalidate evidence after the outcome. Optimistic exact-green reuse may return from its one pre generation without creating a post generation because no gate child ran.

**Drift comparison remains complete.** Under-lock validation and the post-run comparison retain tree identity, oracle identity, resolution, closure, and passlisted-environment coverage. Moving snapshot ownership does not narrow `sameSubject`, turn an unavailable identity into an empty identity, or allow a pre-execution generation to stand in for the post-run read.

**All five identity families migrate as one cut.** The enumerated families are component, conformance check, conformance canary, whole subject, and stripped subject. Leaving any one rooted at an independent tree capture would preserve both the coherence defect and family-proportional process growth, so partial migration is not a shippable intermediate state. Build-time ticketing treats this as a wide internal refactor and follows `craft-tickets`' expand–migrate–contract sequence while keeping every landing green.

**Process bounds are formulas over source work, not timings.** For one ordinary real execution, the source layer performs at most three working-tree materializations: the pre candidate, one under-lock validation, and the post generation. Only the pre and post generations parse listings, so they perform at most two recursive tree listings in total. For one prospective real execution, the supplied pre tree requires no working-tree materialization and the checkout-based post check performs at most one; pre and post still perform at most two listings. Within either source kind, each distinct blob object requested by identity policy is read at most once per parsed generation. These ceilings do not grow with the number of identity families or repeated requests for the same blob.

**The existing safety fallbacks remain the policy owners.** Snapshot creation or subject construction failure refuses the subject. Component/check/canary identity failure cannot authorize inherited evidence and therefore executes the affected complete inventory in the existing direction. Malformed entries, non-blob requests, hostile symlinks, escaping targets, and unavailable prospective trees retain their existing errors or widening behavior. No shortened entry set or partially populated identity map may authorize a skip.

**Evidence and oracle schemas do not change.** Policy domains, identity framing and ordering, component/check slot keys, verdict fields, inspection reasons, announcements, and prospective bootstrap evidence remain stable. A successful refactor over an unchanged tree computes the same public identities and records; this scope adds no migration, version bump, or compatibility reader.

**Ownership fence.** One writer owns `internal/gate/`. The source adapters, evaluation owner, identity migrations, focused count/fault tests, and existing semantic controls all live behind that package seam. No production edit outside this fence is authorized by this spec.

## Testing decisions

- **External behavior a good test exercises.** Drive ordinary `Execute` and prospective `ExecuteTree` in temporary Git repositories, then compare selection, output, exit, inspection, durable verdict, component/check slots, stripped evidence, and exact-tree behavior with the existing controls. Mutate the tree during the child run to prove the second generation still detects drift.
- **Primary seam.** The internal gate-evaluation owner is the highest seam that observes all five identity families and can deterministically count adapter operations. Its tests inject working-tree and prospective-tree sources that return real Git objects while recording materialization, listing, and blob requests.
- **Lower seam.** The immutable snapshot contract receives focused parser, blob-memoization, and fault tests. This second seam is justified because malformed listing data and repeated object-read failures cannot be induced deterministically through the public gate without replacing Git itself.
- **Prior art.** Existing subject, component-identity, check-identity, stripped-subject, component-decision, drift, and prospective tests own the semantic expectations. The single-listing-parser test remains the structural guard against a second parser. New tests compose those owners rather than copying their fixture harnesses.
- **Gate seam.** `bench gate` observes the package tests through the ordinary `test` component, while the current gate/canary controls continue to prove selection and evidence behavior. This scope changes no gate policy or gate script.
- **Central mutation probes.** Reintroduce an independent source capture in one identity family: the operation-count test must exceed its ceiling. Reuse the pre generation for the post comparison: the mid-run drift test must green the child but red the action. Route prospective input through the working-tree adapter: the exact unpublished-tree test or prospective materialization bound must fail.

### Seam diagrams

    trigger: ordinary gate execution
        │
        ▼
    working tree ──▶ [ working-tree source adapter ] ──▶ immutable pre generation
                                                              │
                                                              ▼
                         [ gate evaluation owner: subject + component/check/canary decisions ]
                                                              │
                                                   gate child │ outcome
                                                              ▼
    working tree ──▶ [ working-tree source adapter ] ──▶ distinct post generation ──▶ drift/evidence result
                         ◀ tests attach here: inject source, count operations, observe records and output

    trigger: ExecuteTree(root, supplied tree)
        │
        ▼
    supplied Git tree ──▶ [ prospective-tree source adapter ] ──▶ immutable pre generation
                                                                        │
                                                                        ▼
                                                        [ same gate evaluation owner ]
                                                                        │
    prospective checkout ──▶ [ post source capture ] ──▶ distinct post generation
                         ◀ tests attach here: exact tree, full inventories, post-run drift, zero pre materializations

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | On one stable ordinary tree, component, conformance-check, conformance-canary, whole-subject, and stripped-subject identities consume one accepted pre generation | gate-evaluation owner with recording working-tree source | observed in the 2026-08-04 census: one mapping row launched 415 Git children and materialized the same fixture tree 47 times | A facade that leaves any enumerated family on its old root-based resolver remains semantically green but exceeds the shared-generation operation bound |
| 1 | Component partitions and exact identity values remain unchanged for build, gofmt, vet, test, race, conformance-suite, contract, shellcheck, and canary | existing component identity and component decision suites plus new before/after golden capture | already covered by the current component identity/decision controls; the implementation must retain their literal expected identities and partitions rather than regenerating expectations | An implementation that omits one component input or changes framing can reduce process count while wrongly inheriting evidence |
| 1 | Conformance check identities retain subject, tier, meta flag, implementation closure, declared inputs, invocation schema, and owning canary families; canary identities retain shared, bound, and unattributable implementation behavior | existing check-identity and check-slot suites | already covered by the current check identity and slot controls | These fields are the authorization preimage; the rows catch a cheaper snapshot port that hashes only file content |
| 1 | Stripped identity ignores exactly declared reduced-scope paths and moves on every undeclared edit, while whole-subject identity moves on either class | existing stripped-subject suite driven through the evaluation owner | already covered by `TestStrippedIdentityIgnoresAllowlistedEdit` and `TestStrippedIdentityMovesOnUnlistedEdit`; wiring through the owner must keep both controls green | Computing stripped identity from another generation can preserve each local hash test yet authorize the wrong ancestor during a real evaluation |
| 2 | A gate child that exits zero after changing a tracked or untracked unignored file returns the existing operational drift failure and authors no ready green or reusable evidence | public `Execute` drift fixture plus durable-record assertions | already covered for tracked drift by `TestExecutionLockAndDriftFailClosed`; add the untracked case before migration | Reusing the pre generation for the post check makes a drifting child look stable and would turn its zero exit into authority |
| 2 | Tree, oracle closure, resolution, and passlisted-environment movement before execution or during the child all refuse in the existing fail-closed direction | gate-evaluation owner with source/oracle mutation points | already covered by the gate fault ledger and subject-mutation controls; retain one row per enumerated class through the new owner | Comparing only tree IDs preserves the common case but misses an external gate or environment change whose oracle identity moved |
| 2 | A stable real execution creates a post generation distinct from the pre generation even when both resolve the same tree and blobs | recording source plus generation-identity assertion | not TDD-able before the generation seam exists; the first build ticket must scaffold an inert source interface, then observe the distinct-generation assertion red before wiring behavior | Pointer or cache reuse across the run would make the drift test incapable of observing a change after child execution |
| 3 | `ExecuteTree` grades the exact unpublished supplied tree, leaves the ordinary checkout untouched, and preserves ordinary-green prospective bootstrap evidence | existing prospective integration suite | already covered by `TestExecuteTreeBuildsExactUnpublishedBenchkitSource` and `TestOrdinaryGreenRemainsProspectiveBootstrapEvidence` | Falling back to the ambient working tree or changing evidence keys breaks an existing public journey |
| 3 | Prospective execution ignores ordinary component/check slots and runs the complete applicable component and conformance-check inventories | existing prospective full-inventory control | already covered by `TestExecuteTreeIgnoresStoredCheckSlotsAndRunsFullConformanceInventory`; retain the component-full control beside it | Sharing a snapshot contract must not be mistaken for sharing ordinary narrowing policy |
| 3 | The prospective pre generation performs zero working-tree materializations, lists the supplied tree once, and reads each requested distinct blob at most once; its post check materializes at most once | gate-evaluation owner with recording prospective source and real supplied tree | observed red in the census-derived current call graph: prospective subject construction routes the checkout through repeated working-tree `TreeHash` calls | A prospective adapter that merely wraps the current checkout path keeps semantics green while retaining accidental materialization |
| 4 | One ordinary real execution performs at most three working-tree materializations and two parsed listings total, independent of the five identity families | real Git command recorder around one stable gate evaluation | observed red in the focused census: the row's four gate-engine evaluations performed 47 materializations and 24 `ls-tree` calls, exceeding derived ceilings of 12 and 8 | The fixed ceiling catches independent captures even if all resulting identities happen to match |
| 4 | Within each generation, two paths sharing one blob object and repeated requests from different identity families cause one blob read for that object; a read error is memoized and returned identically | snapshot contract with recording blob source | observed red in the census: 84 `cat-file` children included repeated blob reads; add the same-object and same-error assertions first | Memoizing by path or retrying failures still multiplies processes and can combine different object availability within one generation |
| 4 | Snapshot materialization, listing parse, blob read, non-blob request, hostile symlink, escaping target, and unavailable prospective-tree failures never produce a reusable skip or a shortened authoritative identity map | snapshot fault table through gate evaluation | already covered in parts by component/check fail-closed tests; the adapter fault cases must be added before migration and retain full-execution or refusal assertions | An always-empty snapshot is the cheapest wrong implementation: it removes work and can make unchanged evidence appear valid unless each failure posture is asserted |
| 4 | Stable-tree stdout/stderr, action and gate exits, inspection reasons, verdict fields, component/check evidence, and stripped evidence remain equivalent | ordinary and prospective integration controls with literal output and decoded-record comparisons | already covered across gate, component-decision, check-slot, verdict, and prospective suites; consolidate assertions at the evaluation journey without deleting their current bite | A refactor can compute correct identities while silently changing the operator or durable evidence contract |

The semantic degenerate is a shared snapshot used only by whole-subject identity while component, check, canary, and stripped identities still recapture independently; all local identity tests pass, but the family-enumeration and operation-count rows red. The drift degenerate is one snapshot retained across the child run; every stable-tree row passes, but the distinct-generation and mid-run mutation rows red. The prospective degenerate is the working-tree adapter pointed at a disposable checkout; exact-tree semantics may pass, but the zero-pre-materialization row reds.

### Edge inventory

- **Error path** — covered by source materialization, listing parse, blob read, non-blob, symlink, escaping-target, unavailable prospective-tree, pre-execution drift, mid-run drift, gate red, and evidence-persistence controls. A partial snapshot never authorizes evidence.
- **Empty or absent input** — an unborn or empty working repository retains the current empty-index fallback; a declared file absent from the snapshot remains distinct from a present empty blob and refuses where its current identity family refuses. An empty prospective tree is valid only when the existing prospective gate can resolve its oracle.
- **Boundary values** — zero entries, one entry, two paths sharing one blob object, one blob requested by all applicable identity families, and the complete current benchkit tree are covered. The operation bound is constant in family count and linear only in distinct requested blob objects.
- **Malformed and hostile input** — raw spaces, glob characters, tabs/newlines representable by Git's NUL-delimited listing, invalid listing records, invalid metadata, special files rejected by working-tree materialization, dangling and escaping symlinks, gitlinks, and control-bearing diagnostics retain the current bounded refusal posture. The parser continues to consume raw NUL-delimited paths from one source.
- **Interrupted or partial state** — cancellation during source capture returns no generation; cancellation during the gate retains the existing interrupted evidence behavior; process death never publishes a partially populated cache. Durable recovery for snapshot construction is **Won't handle** because generations are evaluation-local and safe to discard.
- **Re-run idempotency** — every invocation creates fresh generation ownership. A prior run's in-memory entries, blobs, or failures are never reused; durable verdict and component/check slots remain the only cross-run evidence surfaces.
- **Process-boundary lifecycle** — ordinary and prospective public entries create their own evaluation in each process. A child gate cannot mutate the parent's pre generation, and the parent must capture post state after the child exits rather than accept an in-memory child claim.
- **Hostile environment** — missing or failing Git, an unreadable object database, cwd below the repository root, a symlinked repository root, a held gate lock, and passlisted environment drift retain current errors and lock behavior. PATH command recorders are test instrumentation only and do not become production policy.
- **Elapsed-time budgets** — **Won't handle as the acceptance oracle**: scheduler and host I/O variance make wall time nondeterministic. The count assertions are authoritative; timings may remain census diagnostics.
- **Cross-evaluation snapshot reuse** — **Won't handle**: a durable parsed-tree/blob cache would add invalidation and lifecycle authority. This scope owns only one evaluation's generations.

## Out of scope

- **`gate-decision-test-seam`.** Moving the exhaustive component/check mapping matrix to the decision seam and retaining representative full-engine bites is the second ordered capability. Derived breadth: resolver-facing test seam, matrix migration, representative engine controls, and fixture/count updates; `6 edits, 2 gate runs`.
- **`conformance-harness-scope`.** Passing each fixture its registered check and moving freshness variants to the lower seam is the third ordered capability. Derived breadth: harness routing, registry attachment, fixture migration, representative entry controls, and conformance canary updates; `8 edits, 3 gate runs`.
- **Ship-tier proof deduplication.** Release-only package enumeration, inherited dev conformance/vet proof, one artifact build, and race-proof ownership change publication evidence policy and require separate shaping; `12+ edits, 4+ gate runs`.
- **Token-pool or reserve implementation and pricing.** Cross-process token transport, reclaim, grant splitting, reserve selection, instrumentation, and the post-reduction census remain downstream of all three demand-reduction specs; `15+ edits, 5+ gate runs` before pricing evidence is complete.
- **FT87 publication timeout behavior.** Replacing the unreachable-port default-client wait with a bounded publication-owned timeout is an independent network-lifecycle capability; `5 edits, 2 gate runs`.
