---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change the lighter-path threshold in .bench/BENCH.md allows without a spec. An opt-in `--full <spec>` run reaches the final green landing.
---

# /bench-implement-spec — do the work at the seams

## Entry orientation
This is the implementation phase. It starts from an approved spec, or from a change small enough for the lighter-path threshold. It declares the line and works vertical slices at the pre-agreed seams. If there is no spec, the change must fall under that threshold in `.bench/BENCH.md`'s "Right-size the process" paragraph. If it does not fall under that threshold, route to `/bench-write-spec` first.
If a spec-backed run has no `specs/<slug>/tickets/` directory or that directory contains no ticket files, return to `/bench-write-spec`; ticket slicing and approval belong there.

## Exit handoff

Close by reporting the implemented stories, each acceptance row's coverage status, and the landed commits. A reviewed spec-backed build closes through the "Land" section below. A light-path build goes straight to `/bench-final-check`. A build that stops short exits through "When the build stops short" below, which recommends its one durable next action.

## Declare the line, validate the tickets, route the venue

Declare the line before you touch code. `craft-line` owns the template, the tier decision, and the escalation ladder. After the write-spec-phase approval,
create or retain one integration worktree at the reviewed graph commit; that commit is its frozen review base. Enter it and run `bench preflight build <slug>` before the first edit; a red preflight stops the phase.

The reviewer approves the spec and the whole ticket graph once, in `/bench-write-spec`. A ticket inside that approved graph starts without a new approval stop. Run `bench preflight build <slug> --charge --ticket <ticket> --base <base> --source-tip <tip>` and read the prepared evidence identity it returns. Verify the fence, dependency completion, and expected source tip before the ticket's fresh author starts. Each ticket goes to a fresh author session under `.bench/BENCH.md`'s authorship rule. The author charge carries the ticket, its coverage rows, its `Writes:` fence, the evidence identity, and the declared line.

Build action requires a narrow author read: act only after this session reads the metadata and ticket pages and targeted sources. An author does not retrieve every evidence page. `bench preflight evidence <id> --verify` is the separate artifact-integrity check.

Build action requires available required context: a retrieval receipt, a terminal cursor, or another consumer's delivery never replaces what this session reads itself. Build action requires a current-action binding: act only after `bench preflight evidence <id> --check-current` binds the artifact to the current assignment and source pair. Build action requires reviewer approval, which a generated charge, a staged artifact, or a verified artifact never supplies. Build action requires the complete task supplement, which this session writes. The supplement names selectors, model, effort, cap, and mutation, and verified evidence never supplies it.

Follow `.bench/BENCH.md`'s ticket authorship and chunk-review policy. `craft-delegate` governs independent review and brief read-only diagnostic consultation. Neither role receives implementation or repair authorship. Before the first dispatch, the orchestrator takes the plan amendment that `.bench/BENCH-reference.md` states. The orchestrator reads manifests, returns, and verdicts, not code, and it refreshes `bench handoff` at each chunk checkpoint.

## Build

Before each approved slice, identify its acceptance target and existing verification route.
At an approved TDD seam, follow `craft-tdd` for its behavioral-red sequence, minimal compiled setup, and row classifications.

A material action changes behavior, a verification target, or a premise that determines the next action. After each material action, rerun that route. Inspect the result before continuing.
One material action can contain several related edits before the rerun. If a result contradicts the approved behavior or seam, stop. Route it through the wrong-spec path under "When the build stops short."

Work approved tickets in dependency order inside their planned chunks. Use TDD only where `craft-tdd` marks the seam. When the spec carries a coverage map, seed the
harness's native task list from `bench coverage <spec>`, whose `rows[N]{story,behavior,seam}` projection is one task per row. The spec and the ticket are final. Do not evaluate another approach. Implement the ticket as written, run its focused checks, then continue to the commit below.

Apply `craft-line`'s retained implementation continuation policy throughout the ticket graph. At initial review or before post-review repairs, including resumed work, read [the bounded repair policy](../skills/bench-craft-line/references/bounded-repair-policy.md).

A prose-only owner edit is an implementation write. Before another session edits it, follow `craft-delegate`'s recorded user-directed author-transfer rules. Run fresh-session adoption after the final committed owner-byte change. A later owner-byte edit invalidates that evidence; rerun it against the new committed tip before review.

Before each ticket, revalidate prepared source identity and required source bytes. Regenerate a charge when either changes. The ticket author runs focused checks and the required mutation probe, then commits the ticket on a green lane. Run the plan's named probe exactly. Do not replace it with a similar mutation. If the named probe cannot run, stop and amend the plan before verification.

An author that takes a regenerated charge reuses exact available source bytes only after the new manifest verifies their membership, role, and requiredness. A matching body digest alone authorizes no reuse. A fresh author makes its own narrow author read from the trusted evidence identity. A transferred final cursor and another consumer's receipt deliver no byte to it.

After each ticket commit, run `bench worktree exec <target> -- bench preflight build <slug>`. A run from the primary checkout grades the wrong tip. After the last ticket in a chunk, freeze the chunk delta and run the three review axes before advancing. `/bench-review-implementation` owns those axes and the review record. The author commits the verification and probe record before the axis dispatch.
When evidence requires an in-scope plan, `Writes:`, or gate expansion, apply `.bench/BENCH.md`'s approved plan-expansion policy before using it.

## Land

Tickets commit green serially in `Blocked by:` order on the retained integration source. After each chunk review, accepted findings go to fresh repair sessions under `.bench/BENCH.md`'s repair rule, and each repair commits with current repair coverage. Start the successor only after that repair closes. Plan commits land before the ticket merge, and a `main` merge lands only before the first chunk. Only record commits follow the chunk tip. The reconciliation commit joins the review delta of the last chunk.

Retain author verification and all three native review results in `reviews/<slug>.md`.
Commit the artifact on its ordinary lane, then run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>` before the successor.

An ordinary lane pass or ordinary green cannot replace this checkpoint. Follow its named record or review action when evidence fails.
After a repair, apply the bounded repair policy's current-evidence rule and its narrow evidence-only exception before retrying.

After the last chunk, reconcile every acceptance row and the integrated behavior. From the clean destination, hand the accepted source to `bench worktree land`; it composes, gates, publishes `Status: implemented`, and releases the source. Then run `/bench-final-check`.

Write the ordinary assessment record before the `bench worktree land` step, and append the landing evidence after it. The review and final-check phases own the other assessment updates. Keep the retained run identity available for those phase boundaries.

## When the build stops short

Report the state: what is done, what remains, the coverage table, and what consumed the cap. Keep landed work landed. Then route by cause:

- Wrong tier: apply `craft-line`'s ladder.
- Wrong spec: route to `/bench-write-spec` with the finding quoted.
- Wrong scope: propose the split for the reviewer to decide.

## `--full <spec>`

This command orchestrates the fresh ticket authors, chunk reviews, the final landing, and `/bench-final-check`. It refuses on a missing or unnamed spec. At every phase boundary it writes the phase reached into `capture/session-handoff.md` and refreshes the pin with `bench handoff`. A green chunk checkpoint and its handoff refresh are not a phase exit. The orchestrator continues into the successor chunk in the same turn, and it stops only on a `craft-line` stop condition. The review phase owns cross-harness review opt-in.

A chunk delta can grow past its approved plan. It touches a path outside the tickets' `Writes:` fences, or it adds a ticket the plan does not list. Such a delta pauses to ask the reviewer before escalating tier.

## `--delegate`

`--delegate` extends a `--full` run with `.bench/BENCH.md`'s delegated run policy. It refuses without `--full`, an approved spec, or an approved ticket graph. Declare the configured model, effort, iteration cap, and author limit before the first dispatch. Record each assignment in the completion plan before you send that author's write charge. A resumed delegated run keeps the recorded identities, source pins, replacement history, and pending obligations.
