# Checkpoint-scoped review

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-10 (the AXI
post-mortem's review-cadence finding: the three-axis review runs once per whole
composition, so a small omission dilutes into a multi-thousand-line diff, every
accepted finding costs a repair-ticket lifecycle plus a whole-composition
re-review, and that cost drove the batching of seven and four independent
findings into single repair tickets), repaired against the seven accepted
findings of the Sol falsification review of this spec's first draft.

## Problem

Semantic review in the spec-build lifecycle is timed for authority, not for
economy. The first three-axis review sees the entire composed candidate after
every ticket has integrated: the AXI run's round 1 confronted the whole
nine-operation migration at once and returned seven findings, round 2 returned
four, and a one-line acceptance-rule gap survived both because it was a tiny
fraction of the diff under review. By the time a finding lands, the fix costs a
repair ticket, a planning commit, assign/checkpoint/integrate, and a fresh
whole-composition review — which is exactly the pressure that produced batch
repair tickets. Per-assignment verification exists (checkpoint receipts carry
focused ticket evidence) but nothing semantic looks at a ticket-sized diff
while the assignment worktree is still open and a finding is still just an
edit.

## Solution

Each assignment gets one advisory three-axis review of its own diff before it
checkpoints, while its worktree is open. A fresh read-only delegate takes the
assignment's fence-scoped diff against its recorded base, the ticket, and the
spec, charged from `craft-review` exactly as today, and returns findings only.
The coordinator dispositions each finding on the spot: fixed by the ticket's
write delegate in its still-open worktree, risk-accepted by the reviewer, or
explicitly deferred to the composed review — never silently dropped. The
done-claim verification then runs against the tree that actually checkpoints,
so a fix never strands stale evidence. The composed review before promote is
unchanged in authority and scope; what changes is what it should have left to
find: composition and interaction defects rather than per-ticket ones.
Repair-ticket assignments get the same pass, which closes the AXI gap where
repair rounds shipped with no small-diff review at all.

## User stories

1. As the reviewer, I want every assignment's diff reviewed on the three axes
   before that assignment checkpoints, with each finding fixed by the write
   delegate in the open worktree, risk-accepted by me, or deferred by name,
   and with done-claim verification binding to the final tree, so a
   ticket-sized defect is caught while it costs an edit instead of a repair
   lifecycle. Line: gpt-5.6-terra / medium. Prose process semantics the gate
   cannot grade.
2. As the reviewer, I want the checkpoint-review outcome recorded per
   assignment — the reviewed subject's identity, verdict, finding count, and
   each finding's disposition — in `capture/session-handoff.md` as the
   checkpoint lands and consolidated into the implementation retro at close,
   so "the pass ran, on this exact tree" is checkable after the fact instead
   of vanishing like the AXI breakdown review did. Line: gpt-5.6-terra /
   medium. Prose and reporting semantics.
3. As the reviewer, I want the composed pre-promote review to keep its full
   three-axis authority while naming the per-checkpoint disposition lines it
   builds on, so earlier passes reduce its discovery load without ever
   narrowing what it may report. Line: gpt-5.6-terra / medium. Prose
   semantics.

## Implementation decisions

- This is a process capability: prose edits only, no new lifecycle operation,
  no CLI change, no retained-state schema. The lifecycle-native retained
  receipt for checkpoint reviews is deliberately not built here — FT184 owns
  checkpoint receipt mechanics and FT200 owns repair-entry plumbing; this
  spec must not pre-empt either with a second evidence convention.
- The pass is advisory and non-blocking: `bench spec build checkpoint` does
  not gate on it, and a finding never overrides the oracle.
- **Ordering.** The cadence per assignment is: the write delegate returns →
  the advisory review runs on the returned tree → findings are dispositioned
  (a fixed disposition sends the write delegate back into its worktree, and
  its new return re-enters this same review-then-disposition step) → the
  done-claim verification of invariant 1 runs against the settled final
  tree → checkpoint. Verification always binds the tree that checkpoints, so
  a review-induced fix can never strand stale done-claim evidence; the review
  delegate reviews a not-yet-verified tree, which is safe because the pass is
  advisory and returns findings only.
- `/bench-implement-spec` owns the step. The coordinator spawns one fresh
  read-only delegate per `craft-delegate`, charged from `craft-review` with
  the three axes over the assignment's fence-scoped diff against its recorded
  base, the ticket file, and the spec rows the ticket covers. The charge
  source stays `craft-review` — no restated axis definitions in the command.
  A harness that cannot spawn a delegate runs the pass inline and flags it,
  the same capability fallback the breakdown review uses.
- **Dispositions.** The set is closed and authority-split: *fixed* (the write
  delegate resumes in its own worktree; the review delegate never edits) and
  *deferred-to-composed-review* (named with the finding) are the
  coordinator's; *risk-accepted* is the reviewer's — the coordinator may
  propose it, and when the reviewer is unreachable under a standing batch
  approval, the acceptance is flagged for veto exactly as the platform's
  batch-approval rule already prescribes, never silently self-granted. A
  finding with no disposition is not a valid close of the step. This
  checkpoint-scoped set is distinct from the command's existing composed
  -review finding-disposition vocabulary and points to it rather than
  merging with it.
- **Evidence.** One line per assignment, written into
  `capture/session-handoff.md` in the same breath as the checkpoint operation
  and consolidated into the implementation retro at close: assignment ID, the
  reviewed subject's identity (the returned tree the checkpoint binds),
  verdict, finding count, and each disposition. The handoff file is durable
  on disk the moment it is written, so a session dying before the retro
  leaves the witness — the failure mode that erased the recovery-discard
  review. Mid-run the file is written, not committed, respecting the standing
  caution about capture commits inside an active lifecycle; a line whose
  subject identity does not match what checkpointed is stale and does not
  witness that checkpoint. A checkpoint reported with no line is the visible
  signature of a skipped pass.
- Repair-ticket assignments take the identical pass — the step keys off
  "assignment about to checkpoint," not off which round produced the ticket.
- `/bench-review-implementation` (the composed review) adds one sentence of
  posture, not a scope change: it receives the per-checkpoint disposition
  lines as inputs, may re-open anything including risk-accepts, and weights
  its attention toward cross-fence composition — the class the per-ticket
  passes cannot see, per `craft-spec`'s composition-degenerate rule.
- `/bench-final-check`'s retro instructions name the per-assignment
  disposition lines as retro content, consolidating them from the handoff
  file; that command therefore sits inside this spec's fence.
- `craft-review` gains the checkpoint moment in its charge-source list (the
  skill already serves "a delegate's returned work"; the addition names the
  pre-checkpoint diff as a standard charge target and its bounded size as
  the point), keeping one source for what a review hunts and cites.
- **Line economics.** The review delegate runs at the project cache's
  review-axis routing — mid tier, medium effort, per `projects/benchkit.md`'s
  cached lines — with `craft-line`'s kit-guidance leverage override applying
  when the diff is kit prose. No cheap-tier default: the cache, not this
  spec, owns routing.
- The pass does not replace the done-claim verification of invariant 1; it
  precedes it in the settled ordering above and reviews meaning, not
  doneness.

## Testing decisions

- Every changed surface is guidance prose; the gate cannot grade its
  semantics. Red signals below are honest reviewer-graded reads anchored to
  observed current state plus fixed-string absence observations, not pattern
  sweeps.
- The dogfood oracle is the next spec-backed build (the ticket-bundle-refusal
  implementation, four tickets): each of its assignments must produce a
  subject-bound disposition line in the handoff file before checkpointing,
  and its retro must consolidate them. A build that promotes without them is
  the observable failure.

### Seam diagram

    trigger: a write delegate returns; assignment not yet checkpointed
        │
        ▼
    fence-scoped diff + ticket + covered spec rows
        ──▶ [ fresh read-only three-axis delegate (craft-review charge) ] ──▶ findings only
        ──▶ [ dispositions: fixed (write delegate resumes) | risk-accepted (reviewer) | deferred ]
        ──▶ [ done-claim verification of the settled tree ] ──▶ checkpoint
                      ◀ capture/session-handoff.md carries one subject-bound line per
                        assignment; the retro consolidates them; the composed
                        pre-promote review consumes them

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| CR1 | 1 | `/bench-implement-spec`'s cadence contains the review-then-verify ordering: on a write delegate's return, one fresh read-only delegate reviews the fence-scoped diff on the three axes (charged from `craft-review` by pointer, findings only, inline-fallback for a delegate-less harness); a fixed disposition re-enters the same step on the delegate's new return; done-claim verification runs against the settled tree, then checkpoint; the step applies to repair-ticket assignments identically. | `.agents/commands/bench-implement-spec.md` delegate/checkpoint sections | reviewer-graded; observed today: the command's cadence binds verification to the delegate's first return and goes directly to checkpoint, and the fixed string `craft-review` appears nowhere in the file | Without the ordering in the owning command, either the pass exists only in conversation (the breakdown-review failure) or a review-induced fix strands verification evidence taken from a superseded tree. |
| CR2 | 1 | The checkpoint-scoped disposition set is closed and authority-split in the same step: fixed (write delegate, own worktree; the review delegate never edits) and deferred-to-composed-review are the coordinator's; risk-accepted is the reviewer's, batch-approval flagged-for-veto when unreachable; a finding with no disposition is not a valid close; the set is named as distinct from the command's existing composed-review finding-disposition vocabulary. | `.agents/commands/bench-implement-spec.md` | reviewer-graded; observed today: the command's `Finding disposition` section (its lines 337 and 384) covers composed-review findings only — no checkpoint-scoped set, no authority split, no worktree-fix route exists for the pre-checkpoint moment | An open set decays to silent drops, and a coordinator-owned risk-accept lets every defect ship relabeled — the platform's ships-what authority is the reviewer's. |
| CR3 | 2 | Each assignment's line — assignment ID, reviewed-subject identity binding the checkpointed tree, verdict, finding count, dispositions — is written to `capture/session-handoff.md` with the checkpoint and consolidated into the retro by `/bench-final-check`'s instructions; a line whose subject does not match the checkpointed tree does not witness it; a checkpoint with no line is the named signature of a skipped pass. | `.agents/commands/bench-implement-spec.md` reporting plus `.agents/commands/bench-final-check.md` retro instructions | reviewer-graded; observed today: the fixed strings `checkpoint-review` and `disposition` appear in neither command's reporting/retro sections, and no artifact durable at review time exists — the retro is created only after terminal promotion | The AXI breakdown review vanished for lack of a durable witness; subject binding is what stops a stale line from satisfying the evidence for a different tree. |
| CR4 | 3 | `/bench-review-implementation` states that the composed review consumes the per-checkpoint disposition lines by their subject identities, may re-open anything including risk-accepts, and weights cross-fence composition; its three-axis authority and scope are unchanged. | `.agents/commands/bench-review-implementation.md` | reviewer-graded; observed today: the composed review's inputs are the diff, spec, and standards sources only — no mention of prior per-assignment review state | Without the consuming side, disposition lines are write-only and deferred findings silently expire — the composed review must know what was deferred to it and for which tree. |
| CR5 | 1 | `craft-review` names the pre-checkpoint assignment diff as a standard charge target alongside the existing ones, keeping the skill the single charge source. | `.agents/skills/bench-craft-review/SKILL.md` charge-source list | reviewer-graded; observed today: the skill's list names the phase review, a delegate's returned work, a PR, and self-review — no checkpoint moment | If the charge lives only in the command, the axis definitions fork the day either file is edited alone — the one-source rule this repo enforces. |

### Edge inventory

- Error path — a review delegate that fails or returns nothing is a skipped
  pass: the coordinator records it as such in the assignment's line and
  proceeds or retries; it never blocks the checkpoint (CR3's missing-line
  rule is what surfaces a silent skip).
- Empty or absent input — a clean review (zero findings) still produces its
  subject-bound line ("clean"); an assignment with an empty diff cannot
  checkpoint anyway under existing lifecycle rules, so the pass never sees
  one.
- Boundary values — a one-line ticket diff gets the same pass at the cached
  review routing; the cost floor is one read-only delegate.
- Malformed input — a disposition outside the closed set, or a line missing
  its subject identity, is not a valid close of the step (CR2, CR3); the
  composed review re-opens anything ambiguous (CR4).
- Interrupted or partial state — an assignment interrupted anywhere in the
  review-fix-verify loop re-enters at the review step on resume: the
  returned tree may have moved, and the line binds to the tree that finally
  checkpoints, never to a stale read. A session dying after the line is
  written leaves the handoff-file witness on disk (CR3).
- Re-run idempotency — re-running the pass on an unchanged tree may produce
  different advisory findings (reviews are not deterministic); the recorded
  line is the last run before checkpoint, bound to that subject.
- Process-boundary lifecycle — the handoff file survives the session; the
  retro consolidates at close; nothing depends on session memory (CR3).
- Hostile environment — the review delegate is read-only and worktree-scoped
  per `craft-delegate`; it returns findings only and cannot write into the
  assignment fence, so a malicious or confused "fix" cannot ride the review
  itself — fixes are the write delegate's, in its own fence.
- Command self-observation — the pass reviews a diff it cannot mutate; line
  recording changes the handoff file, never the reviewed subject.
- Special files and dangling symlinks — no new discovered-path reader: every
  input is an explicit path the coordinator already holds.

**Won't handle:** a lifecycle-native retained receipt or new `bench` operation
for checkpoint reviews — FT184 (checkpoint receipt mechanics) and FT200
(repair-entry plumbing) own that seam; revisit graduation there once the
prose cadence has dogfood evidence. Blocking `checkpoint` or `integrate` on a
review verdict — the pass is advisory by the same rule that keeps every
review advisory. Second-guessing the composed review's authority — CR4 is
posture, not scope. Committing the handoff line mid-run — the standing
capture-commit caution stands; the witness is the written file. Sub-bound
batching detection promised in the ticket-bundle-refusal spec — this
capability makes per-ticket independence *visible* (each assignment reviews
alone); deciding on it stays with the composed review and the reviewer.

## Ownership fences

- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-final-check.md`
- `.agents/skills/bench-craft-review/SKILL.md`

## Out of scope

- Lifecycle-native checkpoint-review receipts — FT184/FT200 (roadmap-owned;
  ~12 edits, 1 promotion gate when graduated).
- The cross-harness falsification pass and its placement — FT158.
- Any change to `bench spec build` behavior, the gate, or retained state —
  0 edits by design; this capability must prove itself in prose before it
  earns machinery.
