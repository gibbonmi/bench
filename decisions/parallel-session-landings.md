# Parallel-session landings

Status: shaping

## Destination

Decide how independent Bench sessions can prepare, grade, maintain, and land
changes concurrently without mixed authorship, review churn from unrelated tip
movement, or weaker gate authority. The result covers ordinary path-scoped
commits and the reviewed spec-build lifecycle, including ticket eligibility,
run-state mutation, recomposition, evidence assembly, and the boundary between
maintenance and publish authority.

## #1: What subjects and cleanliness predicates currently bind ordinary commits and spec-build landings?

Blocked by: none
Type: Research

### Question

Trace the exact subject construction, cleanliness checks, and committed-ticket
predicates for `bench commit`, lifecycle entry, assignment, prospective gating,
and promotion. Record which existing identities could safely distinguish an
owned diff from unrelated working-tree edits.

### Answer

Ordinary commits gate a complete working-tree Git tree only after refusing all
dirty paths outside the named set. The gate already supports an unpublished
exact Git tree as a full prospective subject, and spec promotion uses that
surface. There is no current patch-plus-inputs subject. Assignment separately
requires the selected ticket's working copy and index to equal its committed
`HEAD` bytes. One current exception needs a deliberate disposition:
`bench commit --spec` flips the spec after its gate. Evidence and citations:
`decisions/assets/parallel-session-landings-research.md`.

## #2: What currently invalidates a run after tip movement, and what overlap identities already exist?

Blocked by: none
Type: Research

### Question

Trace the predicates behind recomposition, held-review invalidation, primary
checkout ownership, ticket fences, and ancestry changes. Record whether the
tree already carries enough identity to distinguish overlapping from unrelated
tip movement.

### Answer

Any recognized advance from the recorded base triggers one whole-run
recomposition predicate without inspecting changed paths or fences.
Recomposition clears the complete held review. The run already retains ticket
digests, charged rows, ownership fences, assumptions, candidate identity, and
receipt digests, but no owner compares those identities across runs. The
observed assignment-worktree refusal is a consequence of invocation-root and
spec/ticket identity, not an explicit primary-checkout predicate. Evidence and
citations: `decisions/assets/parallel-session-landings-research.md`.

## #3: What are the actual run-state, evidence, and promote atomicity boundaries?

Blocked by: none
Type: Research

### Question

Trace spec-build state storage and locking, checkpoint and review evidence
schemas, coordinator-derived digests, and every side effect of `promote`.
Record which operations are logically separable in the current implementation
and which must remain atomic to preserve authority.

### Answer

State locking is per spec slug, while the contended branch, green ref, Git
objects, and worktree effects are shared. State replacement is durable, but
lifecycle transitions are journaled multi-step transactions protected by
identity checks and compare-and-swap. The CLI forwards hand-built checkpoint
and review receipt paths; it emits no receipt skeleton. Recomposition and
publication are already separate internal branches of `promote`, but share one
public authority surface. Recomposition still performs exact-green bootstrap,
replay, candidate compare-and-swap, and review invalidation. Evidence and
citations: `decisions/assets/parallel-session-landings-research.md`.

## #4: What is the smallest unit that must be isolated for a green landing?

Blocked by: #1
Type: Grill

### Question

Choose whether landing authority remains a whole-checkout lock or moves to an
owned prospective diff composed with the exact committed base and every gate
input that can affect its verdict. Define the fail-closed response when that
composition cannot be attributed.

### Answer

Isolate an exact prospective composition and the destination-ref
compare-and-swap, not the whole checkout. The composition starts from the
committed base, applies only the attributed patch and any lifecycle-owned
transition, and is graded as one exact Git tree against the gate's complete
derived oracle closure. Publication proceeds only while the destination ref
still equals the expected base. Missing attribution, subject identity, oracle
identity, or compare-and-swap ancestry fails closed. This is a per-composition
landing boundary, not a diff-scoped weakening of the whole-project gate.

## #5: Which tip movements must force recomposition and re-review?

Blocked by: #2, #4
Type: Grill

### Question

Choose whether any ancestry change remains invalidating or whether a run may
retain its composition and held review when the landed change is proven
disjoint from its owned fences and gate-relevant inputs. Define the conservative
fallback for missing or ambiguous overlap evidence.

### Answer

Every destination-tip movement refreshes the prospective composition because
the exact tree to be graded has changed. It does not automatically invalidate
the held review. Review survives only when recomposition proves that the
run-owned patch and every recorded review input are identical; ownership-fence
disjointness alone is not sufficient evidence of unchanged semantics. A
changed input, conflict, missing dependency identity, or ambiguous comparison
clears the review and requires a new one. The final prospective tree still
receives the authoritative whole-project gate regardless of review reuse.

## #6: Who may mutate a run, and at what concurrency granularity?

Blocked by: #2, #3
Type: Grill

### Question

Choose whether a run remains single-coordinator, admits multiple coordinators
through operation-level compare-and-swap state, or decomposes into independently
owned records with a derived run projection. Include whether lifecycle mutations
may originate from assignment worktrees when run identity is explicit.

### Answer

Keep one authoritative run record and allow any coordinator to submit an
operation against an explicit run identity and expected revision. Mutations
serialize at the per-run compare-and-swap boundary; independent assignment
work and separate runs remain concurrent. Invocation checkout is not authority:
an operation may originate from an assignment worktree or another checkout
when the recorded project root, run, revision, and assignment ownership all
resolve exactly. Slow external effects use prepared-operation journals outside
the state lock, then compare-and-swap finalize or remain recoverable. Do not
split authority across independently owned records and a derived projection.

## #7: Which evidence must the CLI derive rather than coordinators assemble?

Blocked by: #3, #6
Type: Grill

### Question

Choose the supported receipt-construction surface for checkpoint and review
evidence, including which facts the CLI derives from the assignment/run and
which externally observed claims remain caller-supplied and attributable.

### Answer

The CLI derives and seals every repository or run fact: run identity and
revision, assignment and ownership, base and candidate, tree and patch
identity, ticket digest, charged rows, fence, assumptions, axes, timestamps,
and receipt digests. Callers supply only attributable external observations:
checks performed and their results/log references, independent probes, and
review findings and dispositions. Coordinators neither reconstruct the schema
nor calculate its hashes.

The surface is AXI-conformant. Successes and errors return minimal TOON on
stdout with honest exit codes, definitive empty/no-op states, bounded previews
and a `--full` escape hatch where receipt bodies are large. Every return appends
contextual `help[]` entries containing concrete next-step helper command
templates. Those helpers carry forward known run, revision, assignment, and
evidence identifiers, while leaving only genuinely unknown observation values
as placeholders. The gate covers the TOON shape, structured errors, exit codes,
and contextual helper returns.

## #8: Should recomposition be a maintenance operation separate from publish?

Blocked by: #3, #5, #6
Type: Grill

### Question

Choose whether an unprivileged, non-publishing recompose operation may update a
run onto an eligible tip while `promote` retains the sole gate, commit, status,
and project-green transition. Define which evidence recompose invalidates.

### Answer

Add a separate non-publishing `recompose` maintenance operation. Any
identity-proven coordinator may use it against an explicit expected run
revision. It may verify exact-green evidence for the new base, replay the
candidate, compare-and-swap the candidate ref, and update recoverable run
state. It may not gate the candidate, create or publish the landing commit,
advance the destination branch or project-green ref, mark the spec implemented,
or terminate the run; `promote` remains the sole owner of those effects.

Recomposition always clears prior promotion evidence. Checkpoint evidence
survives only while its exact assignment inputs remain unchanged, and review
survives only under #5's exact patch-and-review-input rule. Missing exact-green
base evidence or ambiguous identity refuses without mutation and returns the
AXI `help[]` command template needed to establish or inspect the evidence.

## #9: What subject makes a newly authored ticket assignable?

Blocked by: #1, #4, #5
Type: Grill

### Question

Choose whether ticket bytes must remain committed on the shared destination tip
before assignment, may be committed on a run-owned planning subject, or may be
sealed as immutable run evidence without first moving the shared tip. Preserve
reviewability and refusal after ticket mutation.

### Answer

Newly authored tickets become assignable from an immutable run-owned planning
Git subject without first landing on the shared destination branch. `start`
composes the explicitly named staged spec and ticket files into a content-
addressed run tree/ref. Assignment reads and verifies ticket bytes from that
subject, not the coordinator's working copy, index, or destination `HEAD`.

A ticket change requires an explicit compare-and-swap planning revision. The
revision refuses once any dependent assignment, checkpoint, integration, or
review evidence exists; a new run is then the safe replacement. The planning
subject remains inspectable and diffable through the CLI, while run-state JSON
stores its identity rather than duplicating ticket bytes.

## #10: Should ordinary commits and spec-build coordination ship as one scope?

Blocked by: #4, #5, #6, #7, #8, #9
Type: Grill

### Question

After the authority and state choices are known, choose whether one spec owns
both ordinary per-diff landing and multi-coordinator spec builds, or whether the
reviewed boundary is two separately shippable scopes with an explicit ordering.

### Answer

Ship two ordered, independently green specs. The first introduces the exact
prospective landing substrate: compose named paths against an expected base,
run the authoritative whole-project gate on that exact tree, and compare-and-
swap publish it; ordinary `bench commit` adopts this boundary. The second adds
multi-coordinator spec builds on that proven primitive: run revisions,
run-owned planning subjects, AXI receipt helpers, evidence-preserving
recomposition, and checkout-independent lifecycle operations. The second scope
literally depends on the first; neither absorbs optional parallel hardening.

## #11: How does an interrupted mutation recover or abandon safely?

Blocked by: #10
Type: Grill

### Question

Choose the recovery authority and terminal behavior when a coordinator dies
after preparing an operation, during an external effect, or after a branch/ref
compare-and-swap but before run-state finalization. Define when another
coordinator may resume, finalize, abandon, or preserve recovery payloads.

### Answer

— (open; awaiting #10)

## #12: Which gate evidence may cross composed landing subjects?

Blocked by: #10
Type: Grill

### Question

Choose whether independently composed trees may share retained gate evidence,
and distinguish exact-subject reuse, component/check inheritance, exact-green
base authorization for recomposition, and the one project-green transition.

### Answer

— (open; awaiting #10)

## #13: What AXI status surface exposes concurrent runs and conflicts?

Blocked by: #11, #12
Type: Grill

### Question

Choose the minimal default and detail views for active runs, revisions,
prepared operations, tip drift, conflicts, recovery state, and actionable next
commands without turning session-start context into a full run dump.

### Answer

— (open; awaiting #11 and #12)

## Not yet specified

## Spec-writer discretion

- Reversible internal names and file layout after the reviewer fixes the state
  ownership and authority boundaries.

## Out of scope

- Weakening `bench prep-release` or allowing a partial dev verdict to satisfy
  its exact-tree full-green precondition.
- Granting agents publish authority through a maintenance command.
- Making mixed-authorship commits the normal resolution for concurrent work.
- Changing the closed portability decision across supported harnesses.

## Sources

- URL: `https://axi.md/`
  Supports: #7's TOON, structured stdout error, exit-code, truncation, contextual `help[]`, and per-subcommand help requirements, checked 2026-08-03.
  Drift: mutable external specification; re-verify before spec authoring or AXI conformance changes.
- Path: `decisions/assets/parallel-session-landings-research.md`
  Supports: local-code answers and citations for research tickets #1, #2, and #3.
  Drift: code-derived on 2026-08-03; re-run the cited tests and re-check owners before spec authoring if the lifecycle changes.
- Path: `capture/parallel-session-friction.md`
  Supports: the seven observed refusals and the 2026-08-03 injected-interface-junctions build context that opened this map.
  Drift: session evidence; re-check refusal owners and current tests before relying on its diagnosis.
