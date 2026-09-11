---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change the lighter-path threshold in .bench/BENCH.md lets you build without a spec. Declares its line and uses TDD at seams. An opt-in `--full <spec>` run carries the spec through implementation, review, the final green landing, and terminal reporting — resuming from capture/session-handoff.md.
---

# /bench-implement-spec — do the work at the seams

## Entry orientation
This is the implementation phase. It starts from an approved spec, or from a change
small enough for the lighter-path threshold. It declares the line and works vertical
slices at the pre-agreed seams. If there is no spec, the change must fall under that
threshold in `.bench/BENCH.md`'s "Right-size the process" paragraph. If it does not
fall under that threshold, route to `/bench-write-spec` first.

If a spec-backed run has no `specs/<slug>/tickets/` directory or that directory contains no
ticket files, return to `/bench-write-spec`; ticket slicing and approval belong there.
## Exit handoff

Close by reporting the implemented stories, each acceptance row's coverage status, and
the landed commits. A reviewed spec-backed build closes after its chunk reviews and final reconciliation through `bench worktree land`, then `/bench-final-check`. A light-path build goes straight
to `/bench-final-check`. A build that stops short exits through "When the build stops
short" below, which recommends its one durable next action.
## Declare the line, validate the tickets, route the venue

Declare the line before you touch code. `craft-line` owns the template, the tier decision,
and the escalation ladder. After the write-spec-phase approval,
create or retain one integration worktree at the reviewed graph commit; that
commit is its frozen review base. Enter it and run `bench preflight build <slug>`
before the first edit; a red preflight stops the phase.

After reviewer approval of the selected spec and ticket, run `bench preflight build <slug> --charge --ticket <ticket> --base <base> --source-tip <tip> --full` and retrieve its complete mechanical inputs. Require a complete task-specific supplement. A generated charge or staged artifact does not supply approval. The supplement names selectors, model, effort, cap, and mutation. Verify the fence, dependency completion, and expected source tip before the retained author starts the ticket.

Follow `.bench/BENCH.md`'s retained implementation and chunk-review policy. `craft-delegate` governs independent review and brief read-only diagnostic consultation. Neither role receives implementation or repair authorship.

## Build

Work approved tickets in dependency order inside their planned chunks. Use TDD only where `craft-tdd` marks the seam. When the spec carries a coverage map, seed the
harness's native task list from `bench coverage <spec>`, whose
`rows[N]{story,behavior,seam}` projection is one task per row.

Before each ticket, revalidate prepared source identity and required source bytes. Regenerate a charge when either changes. The retained author runs focused checks and the required mutation probe, then commits the ticket on a green lane. After the last ticket in a chunk, freeze the chunk delta and run the three review axes before advancing.

## Land

Tickets commit green serially in `Blocked by:` order on the retained integration source. After each chunk review, accepted findings return to the retained author and commit there with current repair coverage. Start the successor only after that repair closes.

After the last chunk, reconcile every acceptance row and the integrated behavior. Repeat delegated review only when a later delta or cross-chunk concern invalidates prior evidence. From the clean destination, hand the accepted source to `bench worktree land`; it composes, gates, publishes `Status: implemented`, and releases the source. Then run `/bench-final-check`.

## When the build stops short

Report the state: what is done, what remains, the coverage table, and what
consumed the cap. Keep landed work landed. Then route by cause:

- Wrong tier: escalate one `craft-line` rung and resume.
- Wrong spec: route to `/bench-write-spec` with the finding quoted.
- Wrong scope: propose the split for the reviewer to decide.

## `--full <spec>`

This command orchestrates retained implementation, chunk reviews, the final landing,
and `/bench-final-check`. It refuses on a missing or unnamed spec. At every phase
boundary it writes the phase reached into `capture/session-handoff.md` and refreshes
the pin with `bench handoff`. Outside the kit-guidance set, a diff large enough to hide
bugs pauses to ask the reviewer before escalating tier or adding a cross-harness
falsification pass. Tier escalation is an offer; the review phase owns the standing pass.
