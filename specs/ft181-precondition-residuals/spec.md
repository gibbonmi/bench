# FT181 precondition residuals

Status: implemented

Decision source: compiled map at
`specs/ft181-precondition-residuals/decisions/ft181-spec-build-preconditions.md`,
with its three structured sources under
`specs/ft181-precondition-residuals/decisions/assets/`. All three assets are
2026-08-03 delegated reads produced hours before this spec; every drift clause
was re-verified against the tree at staging (the only commits since are
doc-only) and all three hold. Nothing was unreadable. Two spec-time questions
— the fast-forward's operation set and the husk byte disposition — were
resolved by the reviewer on 2026-08-03 and are recorded as amendments in the
map's #2 and #4 answers.

## Problem

The spec-build lifecycle's precondition layer fails closed in four places
where the state it refuses is benign. A sibling build's legitimate promotion
makes restart of an abandoned run read as marker tampering, so the abandoned
spec is stuck until hand repair. A worktree reduced to a husk or a dangling
symlink — exactly the states abandon exists to clean up — is classified as an
ownership fault, fatal even for abandon. The escape hatch layered over that
misclassification swallows too much: once a prepared abandon operation exists,
even a forged assignment record rides through. And a run whose branch tip
advances before any checkpoint refuses checkpoint and start outright, forcing
either a promote round-trip of an empty candidate or the standing rule that
capture and `main` commits sequence strictly outside run windows.

## Solution

Each refusal narrows to what it can actually prove. Restart compares against
the live green marker — the same source a fresh start reads — so a sibling's
promotion stops refusing an unrelated restart, while the run's own recorded
base stays in the retained history as evidence. Present-but-not-a-checkout
states classify as liveness, so abandon proceeds and releases ownership
without deleting any bytes; identity faults stay fatal everywhere, and the
blanket prepared-abandon exemption disappears because the narrowed inner
classification is the one source of softening. An empty non-terminal run
whose tip legitimately advanced fast-forwards its recorded base provisionally
on checkpoint and start instead of refusing, retiring the run-window
commit-sequencing rule when it lands.

## User stories

The resolved ids are the `claude` column of the profile's harness × tier
table; a different harness rebinds from that table rather than from these
tokens.

1. As a reviewer restarting an abandoned spec build after a sibling build
   promoted on the same branch, I want `bench spec build start` to succeed
   against the live green marker, so that one build's legitimate promotion
   cannot wedge another's restart. Line: opus / medium. The change swaps the
   `expected` operand of the gate's marker compare-and-swap — authorization
   semantics, which the profile routes at mid effort as gate logic.

2. As an operator abandoning a run whose assignment worktree decayed into a
   husk or a dangling symlink, I want abandon to classify the path as
   not-live and proceed — releasing ownership while leaving the bytes in
   place — so that the cleanup command works exactly where cleanup is needed.
   Line: opus / medium. The classification boundary spans two packages and a
   wrong softening would release a live or foreign checkout, so the semantics
   need the mid tier.

3. As a maintainer of the lifecycle's fail-closed posture, I want a missing
   intent registration classified as liveness and the blanket prepared-abandon
   exemption removed, so that an identity fault can never ride an in-flight
   abandon and the softening rule has one source. Line: opus / medium. The
   deletion is small but it moves the boundary that keeps tampered records
   fatal, which is correctness-critical.

4. As a session whose run's branch tip advanced before any checkpoint, I want
   checkpoint and start to fast-forward the empty run's recorded base
   provisionally instead of refusing, so that a capture or `main` landing
   inside a run window no longer wedges the build. Line: opus / medium. The
   fast-forward's boundary — which operations, which runs, and what promote,
   review, and restart must keep — is the load-bearing half.

## Implementation decisions

**Restart reaches the marker through the existing fallback.** The
terminal-restart branch passes an empty `previousGreen` so `startRun`'s
live-marker fallback fires — the one existing marker-read site — rather than
adding a second read. `run.Base` is not passed anywhere: it survives in the
retained attempt history that restart already appends, which is the map's
"evidence, not comparison operand." The anti-tamper cost is accepted by the
map: with the live marker as `expected`, Bootstrap's ancestor check constrains
the marker against the restart tip only, exactly as it does for a fresh start.
The single-read-site property is not gate-observable — a second read site
producing the same value passes every behavioral row — so it is a
standards-review obligation under the one-source-per-fact rule, not a
coverage row.

**Liveness is decided by shape, not by probe failure.** The checkout probe
classifies by what is at the path, so a probe error can never be mistaken for
absence. Liveness — abandon proceeds, every other mutation refuses: no
filesystem entry; a dangling symlink; a non-directory entry (a regular file,
FIFO, or device node); or a directory with no git metadata entry (the husk).
Fatal identity — every operation refuses: a directory that has a git metadata
entry but whose probe fails (permissions, corruption — it may be a live
checkout we cannot see); a resolvable checkout whose git common directory
differs from the repository's (provably a stranger's); and any `Lstat`
failure other than not-exist. The classification never opens the path: it
stats and runs git against it, so a FIFO cannot block the probe.

**Abandon releases ownership; it never deletes present bytes.** Reviewer
decision 2026-08-03, recorded in the map's #2 amendment: for a
not-live-but-present path the abandonment plan uses a **new non-deleting plan
action** — it performs the same registration cross-check the removed-checkout
path already performs in `internal/worktree` (the live worktree registration
reconciled against the intent ledger, any existing recovery ref re-asserted),
and apply releases the registration and intent entry while leaving the
filesystem entry untouched. It does not inherit the removed-checkout path's
force-removal. The plan carries the leftover path so disposal routes through
the existing size-bounded clean surface; the new action and path join the
plan fingerprint like every existing plan fact. A husk whose registration
cross-check mismatches — registration absent, branch ref, path, or request
disagreeing with the ledger — refuses at the planner, exactly as the
removed-checkout path refuses today.

**The exemption collapses into the classification.** Reclassifying a missing
intent registration as liveness makes the inner ownership check pass abandon
through every genuinely-dead state, so the outer prepared-abandon exemption —
which today re-admits ANY ownership error once a prepared abandon operation
with a recorded result exists — is deleted rather than narrowed. One source
decides what abandon may proceed through. Identity faults are fatal on every
path including a resumed abandon, and the fatal set is enumerated: a
duplicate assignment ID, path, or owner request; an owner-request digest
mismatch; a registration resolving to a different assignment ID or worktree
path; and a foreign checkout.

**The fast-forward is a closed list: checkpoint and start, non-terminal runs
only.** Reviewer decision 2026-08-03, recorded in the map's #4 amendment.
"Empty" is exact: no assignment carries any checkpoint evidence and the
candidate tip still equals the recorded base. For such a run, when the
working tip is a recognized descendant of the recorded base, checkpoint and
start fast-forward the recorded base and candidate — the durable candidate
ref moves by compare-and-swap from its old tip, the same sequence promote's
recomposition uses — and proceed; start then reports status instead of
erroring. The fast-forward is placed behind that closed operation list, never
in the shared precondition return path, because four other consumers depend
on the recompose refusal firing: **promote** routes `errRecompose` into its
recomposition, which runs Bootstrap before moving anything — swallowing it
would skip the gate; **review** must refuse without mutation, since its
receipt is validated only later; **assign** keeps today's refusal (reviewer
decision — start first, then assign); and **restart** of a terminal run fires
only when the refusal surfaces, so a terminal run is never fast-forwarded.
The fast-forward is provisional — no gate runs, because every pre-promote
transition is provisional by the lifecycle's own contract and `promote`
remains the sole green transition. A run with any checkpoint evidence keeps
today's behavior exactly: refusal routing to `promote`.

**Non-ancestor movement stays fatal.** The recognized-advance guard is
untouched: a tip that is not a descendant of the recorded base refuses as
subject mismatch on every path, fast-forward included.

**Existing tests this build re-scopes.** Story 4 changes behavior that four
existing tests pin, and this list is the build's authorization to re-scope
them — each keeps its scenario and gets the new expectation; nothing else
about them may change, and no other existing test may be edited:

- `TestNonAbandonMutationsStillRecomposeOnMovedTip`: the start and checkpoint
  cases now expect fast-forward success; the assign case keeps the refusal.
- `TestLifecycleMutatorsRefuseSharedPreconditionDriftWithoutMutation`, the
  working-advance condition: the checkpoint case now expects fast-forward
  success; the assign case keeps the refusal.
- `TestStartResumeAndConflictsDoNotDuplicateRun`, the moved-tip subtest: start
  now fast-forwards and reports status without duplicating the run.
- `TestRecompositionErrorIsStable` drives assign and stays green unchanged —
  it becomes the assign-refusal regression control.

## Testing decisions

- The primary seam is the `specbuild.Service` public lifecycle surface
  (`Start`, `Abandon`/`ApplyAbandon`, `Checkpoint`, `Promote`, `Review`),
  driven against real temporary git repositories exactly as the package's
  existing tests do. Prior art: `start_test.go`'s real-authorization and
  recording gates, `abandon_test.go`'s removed-worktree and forged-identity
  fixtures.
- Story 1 is asserted at two depths: with the real authorization gate
  (sibling promotion then restart must succeed end-to-end) and with the
  recording gate (the `expected` value Bootstrap receives is the live marker,
  not the recorded base) — the second is what pins the mechanism if the
  refusal message ever changes.
- Story 2's classification is one table test over the decayed-path shapes —
  husk, dangling symlink, FIFO, regular file, unreadable-with-metadata — so
  a lazy classifier cannot pass by special-casing one shape. The worktree
  half tests `internal/worktree`'s abandon planning directly; prior art is
  `recovery_retry_test.go`'s removed-directory plan and apply cases.
- Story 3 is asserted from both directions: a missing registration lets
  abandon proceed, and every enumerated identity class refuses abandon even
  when a prepared abandon operation with a recorded result exists — the state
  the deleted exemption used to admit.
- Story 4 drives a fresh-process reload after the fast-forward (prior art:
  the recomposed-attempt reload test in `promotion_recompose_test.go`), so
  the moved base survives serialization, not just the mutating process's
  memory. Its boundary controls use fixtures with real checkpoint evidence
  and a genuinely terminal run, because the empty fixtures the package has
  today cannot distinguish the boundary.
- The feature gate is `bench gate`; the `test` and `race` phases grade both
  packages, and no new gate phase or CLI subcommand is added.

### Seam diagram: restart marker source

```text
trigger: bench spec build start <slug>  (run is Terminal, abandon completed)
    |
    v
recorded run --> [ Start terminal branch ] --previousGreen: ""--> [ startRun ]
                                                                      |
                                              live refs/bench/green/<branch>
                                                                      |
                                                                      v
                                                        [ gate.Bootstrap CAS ]
   ^ tests attach here: real repo, sibling promotion moves marker+branch,
     restart must succeed; recording gate asserts expected == live marker
```

### Seam diagram: liveness classification and abandon

```text
assignment path --> [ checkout probe ] --> absent .................. liveness
                        |                  dangling symlink ........ liveness (new)
                        |                  non-directory entry ..... liveness (new)
                        |                  dir, no git metadata .... liveness (new)
                        |                  dir with metadata,
                        |                    probe fails ........... identity, fatal
                        |                  foreign checkout ........ identity, fatal
                        v
        [ ownedAssignments ] --liveness + op==abandon--> proceed
                        |                                   |
                        |                                   v
                        |                [ worktree abandon plan: new
                        |                  non-deleting action — reconcile
                        |                  registration, re-assert recovery
                        |                  refs, release; bytes left in
                        |                  place, leftover path in the plan ]
                        +--identity--> refuse (every op, no exemption)
   ^ tests attach here: one table test over the path shapes through Abandon,
     ApplyAbandon, and Checkpoint; planner fixtures for the cross-check
```

### Seam diagram: empty-run fast-forward

```text
trigger: bench spec build checkpoint | start   (closed list; run non-terminal)
    |
    v
run state --> [ empty? no checkpoint evidence        --no--> errRecompose
                and candidate == base ]                      (unchanged, to promote)
    | yes, and tip is recognized descendant
    v
[ fast-forward: candidate ref CAS old->tip, base := tip ] --> mutation proceeds
                                                              (no gate runs)
unchanged around the list: promote --errRecompose--> recomposition + Bootstrap
                           review | assign refuse; terminal runs restart
   ^ tests attach here: empty-run fixtures for the list; checkpointed,
     terminal, review, promote, and assign fixtures pin the boundary
```

### Acceptance coverage map

Rows marked TDD-able start red at build time; the sibling-restart,
husk-abandon, missing-registration, and empty-run-checkpoint cases are all
refusals today, observed in the FT176 review and re-derived in the map's
assets. Rows naming an existing test are regression controls and say so; rows
introducing a control that is green at introduction say that instead of
claiming a red.

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | after a sibling promotion advances branch and marker, restart of an abandoned run succeeds | Service.Start, real authorization gate | TDD-able: the fixture refuses today with the marker-conflict message. | This is the exact reported failure; nothing else exercises restart under a moved marker. |
| 1 | on restart, Bootstrap's expected value is the live marker, and the recorded base survives in retained history | Service.Start, recording gate | TDD-able: the recording gate observes `run.Base` today; the history assertion pins the evidence half. | An implementation that deletes the history entry passes the end-to-end row while destroying the anti-tamper evidence the map kept. |
| 1 | restart with the marker absent bootstraps fresh, as a fresh start does | Service.Start, real authorization gate | TDD-able: today the absent-marker restart refuses because a non-empty expected demands a marker. | Routing through the fallback must inherit the fallback's absent-marker semantics, not add a fourth behavior. |
| 1 | a fresh (non-restart) start keeps its existing marker semantics | Service.Start | Already covered by the fresh-path marker tests in `start_test.go`, which run unchanged. | The regression control that the restart change touched only the terminal branch. |
| 1 | an abandoned empty run whose tip advanced still restarts — the fast-forward never fires on a terminal run | Service.Start on a terminal run | Composition control, green once story 1 lands: red only if story 4's fast-forward converts the recompose refusal restart depends on into success. | Restart fires only when the refusal surfaces; an unbounded fast-forward silently turns restart into a status report and CASes a terminal run's candidate. |
| 2 | abandon proceeds for a husk (directory present, git metadata gone) and releases registration and intent entry | Service.Abandon/ApplyAbandon | TDD-able: the husk fixture refuses as ownership today. | The primary reported failure of face 2; only a husk fixture exercises the metadata-gone branch. |
| 2 | abandon proceeds for a dangling symlink and for a FIFO at the assignment path, and the probe returns without opening either | Service.Abandon/ApplyAbandon, table test | TDD-able: both shapes refuse as ownership today; the FIFO fixture has no writer, so an implementation that opens it hangs instead of returning. | A classifier keyed on directory-ness alone passes the husk row and fails the symlink; one that opens the path blocks on the FIFO — the profile's special-file class. |
| 2 | a directory carrying git metadata whose probe fails stays fatal for every operation, abandon included | Service.Abandon on an unreadable checkout | Green at introduction (the state refuses today); red against the degenerate that softens every probe failure to absence. | The cheapest wrong classifier — probe failed, call it absent — would release a live checkout we merely cannot read; this fixture is what kills it. |
| 2 | the husk's bytes are untouched after apply, and the plan's new non-deleting action names the leftover path | worktree abandon planner | TDD-able: the planner has no present-but-not-a-checkout route today; assert byte-identical husk content after apply and the action and path in the plan. | An implementation that routes the husk through the existing removed-checkout action force-deletes uncommitted bytes no recovery ref holds — the discard the reviewer refused. |
| 2 | the planner refuses a husk whose registration cross-check mismatches: registration absent, branch ref, path, or request disagreeing with the ledger | worktree abandon planner | TDD-able: red against the blind-release degenerate; the removed-checkout path's mismatch refusals are the prior art. | Softening liveness without the cross-check lets abandon release a registration that belongs to different work. |
| 2 | existing recovery refs survive an abandon over a not-live-but-present path | worktree abandon planner | TDD-able against the planner; prior art is the removed-directory recovery-survival test. | Extending the removed-checkout path must inherit its recovery guarantee, not just its permission to proceed. |
| 2 | a foreign checkout at the assignment path still refuses abandon | Service.Abandon | Already covered by the foreign-checkout refusal test, which runs unchanged. | The regression control that softening stopped at "cannot prove" and never reached "provably someone else's." |
| 2 | every non-abandon mutation still refuses a husk or dangling symlink | Service.Checkpoint on the decayed fixtures | Green at introduction (non-abandon ops refuse these states today and keep refusing); red against a classifier that softens the shared probe for all operations. | A classification moved at the shared probe could accidentally soften checkpoint's write path into a dead checkout. |
| 3 | abandon proceeds when the intent registration is gone while the run record still lists the assignment | Service.Abandon/ApplyAbandon | TDD-able: delete the ledger entry; the state refuses as ownership today. | The reclassification is the half of face 3 that makes the exemption deletable; nothing else exercises not-found. |
| 3 | every enumerated identity class — duplicate ID, duplicate path, duplicate owner request, digest mismatch, registration resolving to a different assignment, foreign checkout — refuses abandon even with a prepared abandon operation and recorded result present | Service.ApplyAbandon resumed mid-apply, table over the classes | TDD-able: construct the prepared-op state, then forge each class; today the blanket exemption admits all of them. | This is the exact hole the deletion closes; one representative forgery would let an implementation reject it while still swallowing another class — the enumeration is the row. |
| 3 | an interrupted abandon apply (some worktrees released, then a crash) completes on re-run | Service.ApplyAbandon re-entry | Already covered by the apply re-entry and journal reconcile tests, which run unchanged. | The regression control that deleting the exemption did not re-wedge the mid-apply resume the exemption existed for. |
| 4 | on an empty non-terminal run whose tip advanced, checkpoint proceeds after fast-forwarding base and candidate | Service.Checkpoint | TDD-able: start, land a commit, checkpoint refuses with the recompose message today. | The reported wedge, on the op that hits it first in practice. |
| 4 | start on the same state reports status instead of erroring, with base and candidate advanced | Service.Start | TDD-able: the same fixture errors on start today. | Start's resume path is a distinct call site; fixing checkpoint alone leaves `start` still erroring on an in-flight run. |
| 4 | the record agrees with the moved base and candidate after a fresh service reload | Service.Checkpoint, fresh Service instance | TDD-able: reload and assert base, candidate tip, and ref agreement. | An in-memory-only fast-forward passes the mutating process and corrupts the next session's identity checks. |
| 4 | a fast-forward against a candidate ref that moved externally refuses without mutating the record | Service.Checkpoint after moving the durable ref | New guard, green only under the specified CAS; red against an unconditional ref update. | "Moves by compare-and-swap" is otherwise a prose promise — persistence rows cannot distinguish CAS from overwrite. |
| 4 | a run with real checkpoint evidence and a moved tip still refuses checkpoint with the promote route | Service.Checkpoint on a checkpointed fixture | Green at introduction: a new fixture with actual checkpoint evidence — the package's existing moved-tip tests all use empty runs, so none pins this boundary. | The cheapest wrong implementation fast-forwards every descendant tip; only a genuinely checkpointed fixture goes red on it. |
| 4 | promote on an empty moved-tip run still routes through recomposition and its Bootstrap | Service.Promote | Already covered by the empty-recomposition success and bootstrap-refusal tests in `promotion_recompose_test.go`, which run unchanged. | A fast-forward reachable from promote's precondition call would skip the one gate run that makes recomposition green — the regression control is load-bearing. |
| 4 | review on an empty moved-tip run still refuses without mutation | Service.Review | Green at introduction: new fixture asserting refusal and untouched state — no existing test drives review across a moved tip. | Review's receipt is validated after preconditions, so a fast-forward there would mutate base and candidate for a receipt that then fails. |
| 4 | assign on an empty moved-tip run still refuses with the promote route | Service.Assign | Already covered by `TestRecompositionErrorIsStable`, which runs unchanged as the assign control. | The reviewer closed the list at checkpoint and start; this row pins the exclusion. |
| 4 | a non-ancestor tip still refuses as subject mismatch, empty run or not | Service.Checkpoint after a forced non-ff move | Green at introduction (the mismatch refusal is today's behavior); red against a fast-forward keyed on "tip differs" instead of "recognized descendant". | An unguarded fast-forward silently adopts a rewritten history. |
| 4 | the fast-forward runs no gate | Service.Checkpoint, counting gate | TDD-able: assert zero Bootstrap calls across the fast-forward. | A gate call here would make a provisional transition project-green authority, violating the lifecycle's promote-only contract. |
| 4 | a second mutation after the fast-forward is a no-op advance (idempotent) | Service.Checkpoint twice | TDD-able: assert one ref move and stable state across the second call. | A fast-forward that re-fires on equal tips would churn the candidate ref every mutation. |

Degenerate implementations, pinned per story: leaving `run.Base` in place (1)
fails both sibling rows; passing the live marker but dropping the retained
history (1) fails the evidence half. Softening every probe failure to absence
(2) fails the unreadable-with-metadata fixture; softening by
directory-existence (2) fails the dangling-symlink and FIFO rows; opening the
path (2) hangs on the FIFO; deleting husk bytes or reusing the removal action
(2) fails the bytes-untouched row; releasing without the cross-check (2)
fails the planner-mismatch row. Deleting the exemption without reclassifying
not-found (3) fails the missing-registration row; reclassifying while
swallowing any single identity class (3) fails the enumerated table.
Fast-forwarding at the shared precondition seam (4) fails the promote
Bootstrap control, the review no-mutation row, the assign control, and the
terminal-restart composition row; fast-forwarding any moved tip (4) fails the
checkpointed-fixture and non-ancestor rows; an unconditional ref write (4)
fails the CAS-drift guard; an in-memory advance (4) fails the fresh-reload
row.

### Edge inventory

- Error path — resolved by the marker-absent restart row, the
  unreadable-with-metadata row, the foreign-checkout control, and the
  non-ancestor row.
- Empty or absent input — resolved by the absent-versus-present distinction
  (absent is already covered by the removed-worktree tests, which run
  unchanged) and the marker-absent row.
- Boundary values — resolved by the exact "empty" boundary: the
  checkpointed-fixture row, the terminal-run composition row, and the
  idempotent second-mutation row.
- Malformed input — resolved by the enumerated identity-class table; a
  hand-edited run record stays fatal on every path.
- Interrupted or partial state — resolved by the mid-apply re-entry control.
  **Won't handle:** a crash between the fast-forward's candidate-ref move and
  the record save wedges the run until hand repair — the same ref-then-save
  exposure promote's recomposition has today, reused unchanged. Stated
  plainly: abandon cannot rescue that state either (the candidate identity
  check refuses before abandon's exemption is reached), and the window now
  opens on checkpoint and start fast-forwards rather than on promote alone.
  Closing it is a lifecycle-journal capability, not this build.
- Re-run idempotency — resolved by the second-mutation row and the apply
  re-entry control; restart re-entry keeps its existing prepared-start
  resume, which runs unchanged.
- Process-boundary lifecycle — resolved by the fresh-reload row; the
  recomposed-attempt reload test is the prior art and regression control.
- Hostile environment — resolved by the dangling-symlink and FIFO rows (the
  profile's dangling-symlink and special-file classes). **Won't handle:**
  device nodes are asserted through the same non-directory branch as the FIFO
  rather than a device fixture — creating one needs privileges the test
  environment does not have; the FIFO is the shape that can actually block,
  and it has the fixture. **Won't handle:** control bytes in git-sourced text
  — this build renders no new git-sourced value; every message names paths
  already rendered today.
- A command whose write changes a fact it reports — **Won't handle:** no
  command here reports tree state; lifecycle mutations report run state they
  re-read after writing, which the fresh-reload row covers.
- Concurrent second writer during a mutation — **Won't handle** beyond what
  the per-slug lock already provides; cross-process race hardening is the
  lifecycle's existing posture, unchanged.

## Out of scope

- The FT98 preserve-then-discard primitive: disposing of the husk bytes this
  spec deliberately leaves in place — its own capability, already a roadmap
  row (est. 12+ edits, 4+ gate runs there).
- A recomposition path for a *checkpointed* run whose tip moved without
  paying promote — refused by the source: promote's recomposition is the
  replay owner and stays so.
- Extending the fast-forward to `assign` — refused by the reviewer
  (2026-08-03): start first, then assign.
- Retiring the run-window commit-sequencing rule from `ROADMAP.md`'s FT181
  row text: happens at spec retirement with the row, per the retirement
  convention — 0 edits now.
- Any new `bench` CLI subcommand or gate phase. The abandonment plan's new
  action and leftover-path field are in scope — they are the story 2
  contract — and flow through the existing plan rendering unchanged
  otherwise; refusal-wording and `help[]` work is FT173's.
