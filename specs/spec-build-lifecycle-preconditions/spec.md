# spec-build-lifecycle-preconditions

Status: implemented

Decision source: `ROADMAP.md` FT176, drained and reworded in the
reviewer-approved `/bench-what-next` pass landed at `7c2684b` (2026-08-01),
plus one in-session reviewer ruling the same day: the deadlocked
`reduced-gate-phase-set` residue is **not** restored and its spec stays
retired, so every acceptance fixture here is hermetic.

## Problem

The precondition layer refuses the operations that exist to escape it. A
reviewed spec build whose review returns accepted findings must commit repair
tickets to be able to assign them, and committing them moves the branch tip —
which is the exact condition that makes `assign`, `checkpoint`, `integrate`,
and `review` refuse with a message pointing at `bench spec build promote`.
Promote then refuses too, because it grades the review and the released
assignments *before* it reaches recomposition, and a mid-repair run can meet
neither by construction: its review has accepted findings, and its assignments
cannot release while `integrate` is refusing. `abandon --apply` sits behind the
same gate, so a run that cannot promote cannot be retired either. Separately,
once any spec build has completed on a branch, `start` refuses forever: it asks
for green evidence with no expected prior marker, so the `refs/bench/green/<branch>`
marker its own predecessor left behind reads as a conflict rather than the
benign ancestor it is, and nothing in the kit ever moves or retires that marker.

The observed cost, 2026-08-01: a ten-ticket build deadlocked permanently after
two necessary commits to `main`, and finished as light-path work with the
candidate applied as a patch. Its run record remained `active` and unretirable.

## Solution

An operation that exists to resolve a state never refuses on that state.
Promote reaches recomposition before it grades a composition that recomposition
is about to discard, so the documented repair round becomes walkable:
`promote` recomposes onto the moved tip and asks for a fresh review, then
`assign`, `integrate`, `review`, and `promote` run normally. The escape hatch
stops requiring both the recomposition it exists to escape and a worktree that
is already gone. `start` accepts a green marker it could fast-forward and keeps
refusing a divergent one. Every refusal names the operation it refused. The
lifecycle documentation gains the recomposition step it has never described.

## User stories

1. As a session mid-repair, when I commit repair tickets and the branch tip
   moves, `bench spec build promote <slug>` recomposes the run onto the new tip
   and reports that a fresh review is next, instead of refusing on a review
   and a release state that recomposition discards. Line: `opus` / medium.
   Three gates precede recomposition and one of them is shared with the
   promotion recovery path, so this is a move-or-duplicate decision rather
   than the statement reorder it first appears to be.
2. As a session holding a run I have decided to drop, `bench spec build abandon
   <slug> --apply <fingerprint>` never refuses because the branch tip moved.
   The plan phase already ignores the tip entirely; apply stops disagreeing
   with it. Line: `sonnet` / medium. One exemption at the single refusal site,
   with the plan phase as the behavioral reference.
3. As a session holding a run whose assignment worktrees no longer exist on
   disk, `bench spec build abandon <slug>` still plans and still applies:
   an assignment whose worktree is gone is cleanable, not an ownership fault.
   Identity checks stay fatal; only the liveness probe softens. Line: `opus` /
   medium. Deciding which half of the ownership check may soften is the one
   correctness judgment in this spec, and it crosses two packages: the
   lifecycle's ownership probe and the worktree owner's own absent-path
   refusal both block the plan today.
4. As a session starting a build on a branch where a previous build completed,
   `bench spec build start <slug>` fast-forwards the existing
   `refs/bench/green/<branch>` marker when it is an ancestor of the tip, and
   still refuses a marker that is not. Line: `opus` / medium. No cached routing
   in the profile covers spec-build lifecycle Go, so this is a fresh
   three-signal call: it changes what counts as acceptable green evidence,
   which is the one place in this spec where a wrong answer credits ungraded
   work.
5. As a session reading a refusal, the message names the operation that was
   refused. Today every precondition-gated operation borrows `start`'s wording,
   which already misattributed a `checkpoint` refusal and led a session to call
   parking safe when it was not. Line: `sonnet` / low–medium. Mechanical once
   the operation identity is threaded to the message site.
6. As the teammate who just walked in, the lifecycle documentation describes
   recomposition: when the tip moves, `promote` is the operation that recomposes
   and the review is discarded, so the repair round is
   repair → `promote` → `review` → `assign`…`integrate` → `review` → `promote`.
   The word appears in no kit document today. Line: `fable` / high, the
   profile's standing doc-authoring leverage override.

## Implementation decisions

**Promote's ordering is the whole of story 1, and three gates precede
recomposition, not two.** The shared precondition call and its recomposition
branch move ahead of the clean-review check, the released-assignment loop, and
the retained-promotion-evidence validation — that third one independently
refuses the mid-repair state on both an absent review and an unreleased
assignment, so moving only the first two leaves the deadlock intact. The
evidence validation is also called by the promotion recovery path, which must
keep it; the edit therefore duplicates it into the two paths that need it
rather than relocating a single call. Nothing about recomposition itself
changes:
it already replays the candidate onto the new tip, re-bootstraps green
evidence at that tip, clears the review, and returns a status whose next
action is a fresh review. The two checks that move are not weakened — they
still guard the promotion path — they simply stop guarding the path to
recomposition. The promotion recovery fast path, which resolves an
already-published promotion commit and never consults the precondition layer,
keeps its position ahead of everything.

**A moved tip continues to block `review`.** The decision source asks whether
it should, on the grounds that a review grades the candidate rather than the
branch. It should: recomposition discards the review outright, so a review
taken before recomposition is a wasted round by construction, and the
post-fix flow always reaches recomposition through `promote` first. This is
recorded as a decision rather than a change; it is veto surface.

**The escape hatch is exempt from recomposition, not from identity.**
Story 2 exempts the abandon mutation at the single site that returns the
recomposition refusal. Everything else the shared precondition checks —
branch identity, spec identity, candidate-ref currency, assignment ownership
identity — continues to apply to abandon, and the fingerprint drift check
that makes apply safe is untouched.

**Ownership splits into identity and liveness, across two packages.** Story
3's judgment: the uniqueness and digest checks over assignment ID, path, and
owner-request are *identity* and stay fatal for every operation. Probing that
the worktree still exists and still belongs to this repository is *liveness*;
for the abandon mutation, a failed liveness probe on an unreleased assignment
means the worktree is already gone and there is nothing to release. The
recovery refs the plan already enumerates are what preserve any payload, so
softening the probe drops no work.

Two sites block this, and the story owns both: the lifecycle's own ownership
probe, and the worktree owner's plan, which refuses an absent target path
outright. The worktree owner already has the right mechanism one branch
earlier — a prior cleanup receipt short-circuits the absent-path refusal,
which is why an interrupted release resumes cleanly today. Story 3
generalizes that escape to an assignment the record still holds but the disk
no longer has; it does not invent a new one.

This is the face the decision source does not name, found in the residue. The
apply phase does carry an exemption for a failed ownership probe, but it
requires a prepared abandon operation that only a successful plan produces,
and apply creates that operation *after* its own precondition call. So for a
first apply against a removed worktree with no prior cleanup receipt, the
exemption is structurally unavailable — the escape hatch is unreachable for
exactly the state it exists to escape. An interrupted release, which does
leave a receipt, is already recoverable and stays that way.

**`start` passes the marker it found as the expected prior tip — at the site
that actually runs.** The authorization owner already refuses an expected
marker that is not an ancestor of the tip, already short-circuits when the
marker is the tip, and already fails closed when the marker names a missing
object. The change is to stop passing an empty expectation, which is what
turns a benign ancestor into a conflict. No new ancestor check is introduced.

The blocking empty expectation is the one the fresh-start path hands to the
run bootstrap, not the one in the start-completion helper that the decision
source names. The completion helper's empty expectation is real but nearly
unreachable: the fresh path reaches it only with green already established,
and the one caller that does not has already pinned the tip to the recorded
base. The restart-after-terminal path in the same function already passes the
prior base, so the asymmetry is visible in one function. Story 4 fixes the
fresh-start site because that is where the observed refusal comes from, and
fixes the completion helper in the same change so the two stop disagreeing.

**Refusal messages take the operation from the caller.** The working-subject
resolver produces two messages that hardcode `start`; both become
operation-named. The recomposition refusal keeps pointing at `promote`,
because `promote` genuinely is the operation that recomposes.

**Exactly one existing test asserts behavior this spec changes.** The promotion
test that asserts an unreleased run is refused before recomposition encodes the
current contract; its expectation inverts here as authored behavior change
under this spec, never as a check weakened to reach green — flagged for
reviewer veto in the approval table. The test pinning the recomposition
refusal wording keeps its expectation, because that message continues to point
at `promote`. The two other promote-refusal assertions run with an unmoved tip
and survive the reorder unchanged; the changed refusal literals are asserted in
no test today, which is itself why story 5 attaches at the CLI.

**An unrecognized head move still refuses `abandon`.** Story 2 exempts the
recomposition refusal only. A branch that was rebased or amended fails the
recognized-advance check earlier, with a subject-mismatch message, and keeps
failing it: that is identity drift, not recomposition, and an escape hatch
that acted on a record whose branch no longer contains its own base would
clean up against the wrong history. Recorded as a decision because the
decision source's phrasing leaves it open; it is veto surface.

## Testing decisions

Good tests here drive the public `Service` operations against a hermetic
repository and observe refusals, run-record state, and git refs — not internal
predicates. The package's in-package harness already builds a temp repo, a
started run, assignments, checkpoints, integrations, and reviews; every story
composes it rather than adding a second fixture family.

Seams receiving tests, all existing:

- The `internal/specbuild` `Service` operations, through the package's
  table-driven precondition harness that crosses every mutator against every
  drift condition. A new condition row exercises every mutator for free; this
  is where stories 1–3 and story 4's success path attach.
- The recomposition and abandon test files, for the ordering and escape-hatch
  behaviors that need a fully-composed run rather than a table row.
- `internal/worktree`'s own plan tests, for story 3's second site: the absent
  target path is refused there, not in the lifecycle package.
- `internal/gate/authorization`, for story 4's refusal direction. This is a
  deliberate seam split: the lifecycle package reaches green evidence only
  through a gate-owner interface whose in-package double is a bare
  compare-and-swap with no ancestry rule, so a divergent marker cannot be
  observed there — a test written against the double would assert the double.
  The owner's ancestor refusal has no test today and gets one here.
- The runtime contract suite that runs the real `bin/bench.sh` as a subprocess,
  for story 5, for one end-to-end abandon of a tip-moved run, and for story 4
  end to end across a real process boundary.
- The conformance anchor mechanism, for story 6's documented sentences. The
  anchor targets are the spec-build lifecycle lookup table in the kit's
  reference file and the implement-spec command; both exist today.

The gate seam that observes this feature is the dev tier's Go test and
conformance phases; the runtime contract rows also ride the contract phase.

### Seam diagram

    trigger: session runs `bench spec build promote <slug>` on a moved tip
        │
        ▼
    run record + working subject
        │
        ▼
    [ recovery fast path ] ──▶ published promotion, unchanged
        │ (not applicable)          (keeps its own evidence validation)
        ▼
    [ shared precondition + recompose branch ]  ──▶  recomposed run, next = review
        │ (tip unmoved)                                   ▲
        ▼                                                 │ MOVES HERE (story 1)
    [ clean review + released assignments                 │
      + retained-evidence validation ]  ──▶  gate + publish
      (all three gates sit after the move; the third also
       refuses the mid-repair state, so it cannot stay put)
                      ◀ tests attach here: compose a reviewed run, advance the
                        working branch, call Promote, assert recomposition
                        happened and no refusal was returned

    trigger: session runs `bench spec build abandon <slug> [--apply <fp>]`
        │
        ▼
    plan: ownership identity + liveness ──▶ worktree owner plan ──▶ inventory
        │                          ▲                     ▲
        │  liveness softens (story 3, site 1)            │ absent path softens
        ▼                                                  (story 3, site 2)
    apply: fingerprint match ──▶ [ shared precondition ] ──▶ cleanup
                                       ▲
                                       │ recompose exempt (story 2)
                      ◀ tests attach here: build a run with a moved tip and a
                        removed worktree directory; assert plan and apply both
                        succeed and the recovery refs survive

    trigger: session runs `bench spec build start <slug>` after a prior build
        │
        ▼
    working subject ──▶ [ run bootstrap ] ──▶ gate owner bootstrap(expected)
                              │                      │
                              │ passes found marker  ▼
                              └──────────────▶ ancestor: fast-forward
                                               divergent: refuse
                      ◀ success attaches in the lifecycle package: plant an
                        ancestor marker, assert start succeeds and the marker
                        ends at the tip. The refusal attaches one layer down,
                        at the authorization owner, whose ancestor rule the
                        in-package gate double does not model.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | a reviewed run whose branch tip advanced recomposes on `promote` and returns no error, with next = review | recomposition test file, fully-composed run | observed red planned | the current ordering refuses at the clean-review check; a reorder that misses the released-assignment check still refuses |
| 1 | a run with accepted review findings and unreleased assignments recomposes on `promote` — the exact mid-repair state, not merely a clean one | recomposition test file | observed red planned | this is the deadlock; a fix tested only against a clean review passes while the real case still hangs |
| 1 | after recomposition the review is cleared, the base is the new tip, and the candidate ref advanced to the replayed commit | run-record assertion after Promote | observed red planned | a reorder that returns early without recomposing looks green on the error-free assertion alone |
| 1 | with the tip unmoved, `promote` still refuses an absent review, a stale-candidate review, a review with accepted findings, and any unreleased or unintegrated assignment (four cases) | recomposition + promotion test files | already covered; re-asserted after the move | moving the checks past the precondition call must not delete them — the degenerate fix is deletion |
| 1 | with the tip unmoved, `promote` still refuses on each retained-evidence fault: drifted candidate ref, drifted review binding, incomplete checkpoint fields, drifted checkpoint ref, integration outside candidate ancestry (five cases) | promotion test file | observed red planned per case | this is the third gate; duplicating it into two paths is the edit, and dropping it from the promotion path is the invisible half of the degenerate fix |
| 1 | the promotion recovery fast path is unchanged for an already-published promotion commit | existing publish test | already covered | a reorder placed above the recovery branch would re-gate a terminal run |
| 2 | `abandon --apply` with a valid fingerprint succeeds when the branch tip has advanced | abandon test file | observed red planned | today it returns the recomposition refusal; the escape hatch must never require the state it escapes |
| 2 | `abandon --apply` still refuses a drifted fingerprint on a moved tip, and still refuses on branch, spec, and candidate identity drift (four cases) | abandon test file | observed red planned per case | exempting the whole precondition rather than the recomposition branch is the cheap wrong fix |
| 3 | `abandon` plans and applies for a run whose unreleased assignment worktree directory has been removed | abandon test file with the directory deleted | observed red planned | the residue's exact shape; the plan fails the liveness probe today, so apply's existing exemption is unreachable |
| 3 | the worktree owner's own plan returns a cleanable plan, not a mismatch refusal, for an absent target with no prior cleanup receipt | `internal/worktree` plan test | observed red planned | the second blocking site; a fix confined to the lifecycle package leaves this refusal in place and the story red |
| 3 | an interrupted release that did leave a cleanup receipt still resumes through the receipt path unchanged | existing worktree resume test | already covered | the softening generalizes the receipt escape; it must not replace it |
| 3 | assignment identity faults — duplicate ID, duplicate path, duplicate owner-request, owner-request digest mismatch (four cases) — still refuse `abandon` | abandon test file | observed red planned per case | softening the whole ownership check instead of the liveness half would accept a forged record |
| 3 | the recovery refs for a removed worktree survive the abandon and still name their assignment | ref assertion after apply | observed red planned | a cleanup that drops recovery refs destroys the payload the plan promised was preserved |
| 3 | a non-abandon mutation still refuses on a failed liveness probe | precondition table, new condition row | observed red planned | the exemption is per-mutation; widening it lets `integrate` write into a vanished worktree |
| 4 | `start` succeeds on a fresh run when `refs/bench/green/<branch>` is a strict ancestor of the tip, and the marker ends at the tip | start test file, planted marker, fresh `Start` | observed red planned | the live blocker; an implementation that ignores the marker entirely also passes unless the marker's final position is asserted, and one that patches only the start-completion helper never executes on this path |
| 4 | the authorization owner refuses an expected marker that is not an ancestor of the tip | `internal/gate/authorization` unit test | observed red planned — this refusal has no test today | the cheap wrong fix passes the marker unconditionally; this cannot be observed in the lifecycle package, whose gate double is a bare compare-and-swap with no ancestry rule |
| 4 | the in-package gate double is not taught the ancestry rule | review of the double's definition | not TDD-able — a structural constraint on the build, not a behavior | upgrading the double to fake the rule would make the divergent-marker case assert the fixture instead of the product |
| 4 | `start` on a divergent marker refuses end to end through the real CLI | runtime contract suite | observed red planned | closes the gap between the owner's unit refusal and what an operator actually meets |
| 4 | `start` is unchanged when the marker equals the tip and when no marker exists (two cases) | start test file | already covered; re-asserted | the short-circuit and the zero-expectation path must both survive the change |
| 5 | each precondition-gated operation — `assign`, `checkpoint`, `integrate`, `review`, `promote`, `abandon --apply` — names itself, not `start`, in the dirty-checkout refusal (six cases enumerated) | runtime contract suite, real CLI subprocess | observed red planned per operation | one shared literal satisfies a single-operation assertion; only enumerating every operation catches it |
| 5 | the `review` case reaches the precondition refusal rather than the receipt refusal, driven by a valid three-axis receipt bound to the current candidate | runtime contract suite | observed red planned | `review` validates its receipt first, so a stub receipt returns the wrong error and the row passes without testing anything |
| 5 | the same six operations name themselves in the no-working-branch refusal | runtime contract suite | observed red planned | the resolver has two hardcoded messages; fixing one is the cheap half-fix |
| 5 | the recomposition refusal continues to name `promote` | precondition test | already covered; wording re-asserted | over-correcting every message to the calling operation would point a refused `assign` at itself instead of the fix |
| 6 | the spec-build lifecycle lookup table and the implement-spec command each state that a moved tip recomposes through `promote` and discards the review | conformance anchor over the documented sentences | observed red planned by deleting the sentence | an undocumented recomposition step is what let a session believe the repair path was walkable |
| edge of 1 | a tip advance that is not a recognized ancestor still refuses with the subject-mismatch message, not recomposition | precondition table | already covered | the reorder must not turn an unrecognized head move into a silent replay |
| edge of 2 | `abandon --apply` on a rebased or amended branch still refuses on the unrecognized head move | abandon test file | observed red planned | the recorded decision that the exemption is for recomposition only; cleaning up against a history that no longer contains the run's base acts on the wrong tree |
| edge of 3 | a worktree path that exists but belongs to another repository still refuses `abandon` | abandon test file | observed red planned | absent and foreign are distinct; conflating them lets cleanup act on a stranger's checkout |
| edge of 3 | an assignment path containing spaces or glob characters is planned and applied unchanged | abandon test file | observed red planned | the profile's first hostile-input class; the worktree pool already produces such paths |
| edge of 4 | a `refs/bench/green/<branch>` marker naming a commit absent from the object store refuses rather than fast-forwarding | start test file | observed red planned | an unreadable marker must fail closed; the ancestor probe reports false for a missing object and would otherwise read as divergent-but-plausible |
| edge of 5 | a refusal reaching the CLI as an operator string survives a branch name carrying a control byte | runtime contract suite | observed red planned | the profile's control-byte class; the message now interpolates an operation and is read by a TOON-adjacent sink |
| edge of all | every changed refusal and success path is observed through a fresh CLI process reloading the run record, not only in-process | runtime contract suite | observed red planned | the FT164 process-boundary class: unit-level green has hidden serialization defects in this exact package before |

### Edge inventory

- **Error path** — every refusal this spec touches or must preserve: story 1
  four-case row, story 2 four-case row, story 3 identity and liveness rows,
  story 4 divergent-marker row.
- **Empty/absent input** — no green marker, no review, no assignments, an
  absent worktree directory: story 3 and story 4 rows.
- **Boundary** — the marker exactly at the tip versus one commit behind it;
  a tip advance that is an ancestor versus one that is not: story 4 and
  edge-of-1 rows.
- **Malformed input** — forged assignment identity, a marker naming a missing
  object, a drifted abandon fingerprint: story 2, story 3, and edge-of-4 rows.
- **Interrupted or partial state** — a run holding prepared integrate
  operations and cleanup-pending assignments is the residue's own shape and is
  the story 3 fixture; a recomposition interrupted between replay and record
  save leaves the prior candidate ref, which the candidate-currency
  precondition already refuses on the next call.
- **Re-run idempotency** — `promote` called twice on a moved tip recomposes
  once and then reports review-next; `abandon --apply` with the same
  fingerprint is already idempotent through its completed-operation
  short-circuit. Both are asserted in the story 1 and story 2 rows.
- **Hostile environment** — assignment paths with spaces or glob characters,
  a foreign repository at the worktree path, control bytes in a branch name:
  edge-of-3 and edge-of-5 rows.
- **Won't handle:** a run record hand-edited into a state no operation can
  produce — the identity checks refuse it, and repairing arbitrary corrupt
  records is not what an escape hatch is for.
- **Won't handle:** concurrent operations on one run from two processes — the
  per-slug lock already serializes them and this spec does not change it.
- **Won't handle:** retiring or garbage-collecting the `refs/bench/green/<branch>`
  marker. Story 4 makes a stale marker harmless rather than fatal, which is
  what unblocks `start`; deciding a retention policy for the marker is a
  separate capability with its own evidence question.

## Out of scope

- **A standalone `bench spec build recompose` operation.** The decision source
  offers it as an alternative to reordering promote's checks. Reordering is
  strictly smaller and needs no new porcelain, help text, adapter row, or
  documented operation. Exposing recomposition as its own operation is a real
  capability with its own contract — roughly **8 edits, 3 gate runs**.
- **A retention or retirement policy for `refs/bench/green/<branch>`.** Story 4
  removes the harm; deciding who prunes the marker and when needs measurement
  this spec does not have — roughly **4 edits, 2 gate runs** once decided.
- **The ticket-authoring conventions (FT164) and the ticket dependency-edge
  parser (FT174).** Both are open roadmap rows with their own owners. This
  spec's own tickets follow the live corpus shape rather than the stale
  template; see the approval table.
- **Restoring or resurrecting the `reduced-gate-phase-set` residue.** Reviewer
  ruling 2026-08-01: fixtures here are hermetic and nothing is un-retired.
