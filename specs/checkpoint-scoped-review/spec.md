# Checkpoint-scoped review

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-10 (the AXI
post-mortem's review-cadence finding: the three-axis review runs once per whole
composition, so a small omission dilutes into a multi-thousand-line diff, every
accepted finding costs a repair-ticket lifecycle plus a whole-composition
re-review, and that cost drove the batching of seven and four independent
findings into single repair tickets).

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
spec, charged from `craft-review` exactly as today. Findings are dispositioned
by the coordinator on the spot: fixed by the same delegate in the same open
worktree, risk-accepted, or explicitly deferred to the composed review — never
silently dropped. The composed review before promote is unchanged in authority
and scope; what changes is what it should have left to find: composition and
interaction defects rather than per-ticket ones. Repair-ticket assignments get
the same pass, which closes the AXI gap where repair rounds shipped with no
small-diff review at all.

## User stories

1. As the reviewer, I want every assignment's diff reviewed on the three axes
   before that assignment checkpoints, with each finding fixed in the open
   worktree, risk-accepted, or deferred by name, so a ticket-sized defect is
   caught while it costs an edit instead of a repair lifecycle. Line:
   gpt-5.6-terra / medium. Prose process semantics the gate cannot grade.
2. As the reviewer, I want the checkpoint-review outcome recorded per
   assignment — verdict, finding count, and each finding's disposition — in
   the coordinator's build report and the implementation retro, so "the pass
   ran" is checkable after the fact instead of vanishing like the AXI
   breakdown review did. Line: gpt-5.6-terra / medium. Prose and reporting
   semantics.
3. As the reviewer, I want the composed pre-promote review to keep its full
   three-axis authority while naming the per-checkpoint receipts it builds on,
   so earlier passes reduce its discovery load without ever narrowing what it
   may report. Line: gpt-5.6-terra / medium. Prose semantics.

## Implementation decisions

- This is a process capability: prose edits only, no new lifecycle operation,
  no CLI change, no retained-state schema. The lifecycle-native retained
  receipt for checkpoint reviews is deliberately not built here — FT184 owns
  checkpoint receipt mechanics and FT200 owns repair-entry plumbing; this spec
  must not pre-empt either with a second evidence convention.
- The pass is advisory and non-blocking: `bench spec build checkpoint` does
  not gate on it, and a red finding never overrides the oracle. Its timing is
  after the write-delegate reports done and before the checkpoint operation,
  because that is the last moment a finding is an edit in an open fence rather
  than a repair ticket.
- `/bench-implement-spec` owns the step: after a delegate's done-claim is
  verified and before checkpointing, the coordinator spawns one fresh
  read-only delegate per `craft-delegate`, charged from `craft-review` with
  the three axes over the assignment's fence-scoped diff against its recorded
  base, the ticket file, and the spec's rows that the ticket covers. The
  charge source stays `craft-review` — no restated axis definitions in the
  command. A harness that cannot spawn a delegate runs the pass inline and
  flags it, the same capability fallback the breakdown review uses.
- Finding dispositions are closed-set: fixed (the write delegate resumes in
  the same worktree and the fix rides the same checkpoint), risk-accepted
  (named in the build report), or deferred-to-composed-review (named with the
  finding). The coordinator records the disposition line per assignment in
  the build report and the retro; a checkpoint reported with no review line
  is the visible signature of a skipped pass.
- Repair-ticket assignments take the identical pass — the step keys off
  "assignment about to checkpoint," not off which round produced the ticket.
- `/bench-review-implementation` (the composed review) adds one sentence of
  posture, not a scope change: it receives the per-checkpoint disposition
  lines as inputs, may re-open anything including risk-accepts, and weights
  its attention toward cross-fence composition — the class the per-ticket
  passes cannot see, per `craft-spec`'s composition-degenerate rule.
- `craft-review` gains the checkpoint moment in its charge-source list (the
  skill already serves "a delegate's returned work"; the addition names the
  pre-checkpoint diff as a standard charge target and its bounded size as the
  point), keeping one source for what a review hunts and cites.
- Line economics: the review delegate is read-only over a ticket-sized diff —
  cheap tier by default (`BENCH_*_CHEAP`), bumped to mid only when the ticket
  carries prose or authority semantics, per `craft-line`'s cached routings.
- The pass does not replace the delegate done-claim verification in invariant
  1 (gate-facing evidence and `git status` attribution); it runs after it and
  reviews meaning, not doneness.

## Testing decisions

- Every changed surface is guidance prose; the gate cannot grade its
  semantics. Red signals below are therefore honest reviewer-graded reads
  anchored to observed current state — the absence of any checkpoint-review
  step in the named sections today — plus fixed-string absence observations,
  not pattern sweeps.
- The dogfood oracle is the next spec-backed build (the ticket-bundle-refusal
  implementation, four tickets): each of its assignments must produce a
  disposition line before checkpointing, and its retro must carry them. A
  build that promotes without them is the observable failure.

### Seam diagram

    trigger: a write delegate's done-claim verifies, assignment not yet checkpointed
        │
        ▼
    fence-scoped diff + ticket + covered spec rows
        ──▶ [ fresh read-only three-axis delegate (craft-review charge) ] ──▶ findings
        ──▶ [ coordinator disposition: fix in worktree | risk-accept | defer ] ──▶ disposition line
                      ◀ the build report and retro carry one line per assignment;
                        the composed pre-promote review consumes those lines

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| CR1 | 1 | `/bench-implement-spec`'s delegate cadence contains the checkpoint-review step: after done-claim verification and before `checkpoint`, one fresh read-only delegate reviews the assignment's fence-scoped diff on the three axes, charged from `craft-review` by pointer, with the inline-fallback rule for a delegate-less harness; the step applies to repair-ticket assignments identically. | `.agents/commands/bench-implement-spec.md` delegate/checkpoint sections | reviewer-graded; observed today: the command's checkpoint cadence goes from done-claim verification directly to the checkpoint operation, and the fixed string `craft-review` does not appear in the command file | Without a named step in the command that owns the cadence, the pass exists only in conversation — the exact failure mode the AXI breakdown review demonstrated. |
| CR2 | 1 | The disposition set is closed and named in the same step: fixed-in-worktree (fix rides the same checkpoint), risk-accepted, or deferred-to-composed-review; a finding with no disposition is not a valid close of the step. | `.agents/commands/bench-implement-spec.md` | reviewer-graded; observed today: no disposition vocabulary exists in the command | An open disposition set decays to silent drops; the closed set is what makes "every finding went somewhere" checkable. |
| CR3 | 2 | The build report and the implementation retro carry one line per assignment — verdict, finding count, dispositions — and the command names the missing line as the signature of a skipped pass. | `.agents/commands/bench-implement-spec.md` reporting section and the retro instructions it points to | reviewer-graded; observed today: per-round reporting names "checkpointed, or the missed case" with no review evidence | Evidence-free advisory passes are unverifiable after the fact — the breakdown-review lesson; the per-assignment line is the cheapest durable witness short of FT184's receipt mechanics. |
| CR4 | 3 | `/bench-review-implementation` states that the composed review consumes the per-checkpoint disposition lines, may re-open anything including risk-accepts, and weights cross-fence composition; its three-axis authority and scope are unchanged. | `.agents/commands/bench-review-implementation.md` | reviewer-graded; observed today: the composed review's inputs are the diff, spec, and standards sources only — no mention of prior per-assignment review state | Without the consuming side, disposition lines are write-only and deferred findings silently expire — the composed review must know what was deferred to it. |
| CR5 | 1 | `craft-review` names the pre-checkpoint assignment diff as a standard charge target alongside the existing ones, keeping the skill the single charge source. | `.agents/skills/bench-craft-review/SKILL.md` charge-source list | reviewer-graded; observed today: the skill's list names the phase review, a delegate's returned work, a PR, and self-review — no checkpoint moment | If the charge lives only in the command, the axis definitions fork the day either file is edited alone — the one-source rule this repo enforces. |

### Edge inventory

- Error path — a review delegate that fails or returns nothing is a skipped
  pass: the coordinator records it as such and proceeds or retries; it never
  blocks the checkpoint (CR2's closed set includes only findings, not
  delegate health — the missing line rule in CR3 is what surfaces it).
- Empty or absent input — a clean review (zero findings) still produces its
  disposition line ("clean"); an assignment with an empty diff cannot
  checkpoint anyway under existing lifecycle rules, so the pass never sees
  one.
- Boundary values — a one-line ticket diff gets the same pass at cheap tier;
  the cost floor is one read-only delegate, which is the point of the tier
  routing decision.
- Malformed input — a disposition outside the closed set is not a valid close
  of the step (CR2); the composed review re-opens anything ambiguous (CR4).
- Interrupted or partial state — an assignment interrupted between review and
  checkpoint re-runs the pass on resume: the diff may have moved, and the
  disposition line binds to what was actually checkpointed, not to a stale
  read.
- Re-run idempotency — re-running the pass on an unchanged diff may produce
  new advisory findings (reviews are not deterministic); the disposition line
  reflects the last run before checkpoint, and that is the recorded one.
- Process-boundary lifecycle — dispositions survive in the build report and
  retro, which are durable artifacts; nothing depends on session memory
  (CR3).
- Hostile environment — the review delegate is read-only and worktree-scoped
  per `craft-delegate`; it cannot write into the assignment fence, so a
  malicious or confused "fix" cannot ride the review itself — fixes are the
  write delegate's, in its own fence.
- Command self-observation — the pass reviews a diff it cannot mutate;
  disposition recording changes the build report, never the reviewed subject.
- Special files and dangling symlinks — no new discovered-path reader: every
  input is an explicit path the coordinator already holds.

**Won't handle:** a lifecycle-native retained receipt or new `bench` operation
for checkpoint reviews — FT184 (checkpoint receipt mechanics) and FT200
(repair-entry plumbing) own that seam; revisit graduation there once the
prose cadence has dogfood evidence. Blocking `checkpoint` or `integrate` on a
review verdict — the pass is advisory by the same rule that keeps every
review advisory. Second-guessing the composed review's authority — CR4 is
posture, not scope. Sub-bound batching detection promised in the
ticket-bundle-refusal spec — this capability makes per-ticket independence
*visible* (each assignment reviews alone); deciding on it stays with the
composed review and the reviewer.

## Ownership fences

- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/skills/bench-craft-review/SKILL.md`

## Out of scope

- Lifecycle-native checkpoint-review receipts — FT184/FT200 (roadmap-owned;
  ~12 edits, 1 promotion gate when graduated).
- The cross-harness falsification pass and its placement — FT158.
- Any change to `bench spec build` behavior, the gate, or retained state —
  0 edits by design; this capability must prove itself in prose before it
  earns machinery.
