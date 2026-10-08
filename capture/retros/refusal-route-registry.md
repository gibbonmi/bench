## Outcome

FT393 landed at `088cd613` from source `9c228393..0b3dabc0`, and the landing gate was green. Every write verb now routes each refusal through one shared refusal-route registry. The six verbs are the landing, the merge, the reset, the commit, the gate checkpoint, and the commitment verb. Each face declares its verb, its authority, and its route, and `bench recovery` prints the registry as the recovery matrix. A conformance check refuses a route literal that a write verb composes outside the registry. Thirteen tickets in seven chunks covered 68 acceptance rows, and the final reconciliation cited a test for each row except the review-owned RR54.

## Gate-stage timings

- landing: commit 088cd61324629f893a3116633d2acc2b50ebbc8a, trace 10f6d35783d345e8932c52946fbc23d5
- gofmt: 159 ms
- vet: 1396 ms
- test: 276705 ms
- race: 3000 ms
- system: 77880 ms
- shellcheck: 531 ms

## Ticket-versus-spec-slice and delegate performance

An Opus/high orchestrator dispatched 13 fresh Opus/high authors, one author successor, 18 repair sessions, and 5 verification-only transfers. Sonnet/high ran every review axis by reviewer direction: 21 first-round axes and 34 confirming axes. Fable/high ran one read-only consultation. Every author committed inside its fence or stopped at the fence with a precise report. Every chunk needed at least one repair cycle, and RR-C1b, RR-C3, and RR-C5 needed two.

The review axes found real defects in each chunk. The worst were authority errors. One agent-clearable cause routed to the reviewer handback, and one rerun cause sat on a reviewer face. One route was a usage error for its second printing verb. The Coverage axes found the most gaps, and their silent probes proved each one.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| unknown | unknown | unknown | unknown | unknown | unknown |

Brier mean: unknown. The orchestrator declared no confidence for a claim before its outcome.

### The 300k context rule

Mid-build, the reviewer made a standing rule. A delegate near 300k tokens finishes its ticket, and a fresh recorded session takes later work.

| measure | before the rule (10 sessions) | after the rule (25 sessions) |
|---|---|---|
| mean peak tokens | 256k | 156k (188k for authors and repairs) |
| sessions past 300k | 4 (largest 419k) | 3, each finishing its own ticket |
| mean wall time | 24 minutes | 14 minutes |
| blocked returns per session | 1.1 | 0.2 |

Five verification-only transfers averaged 28k tokens, and every record passed on its first run. Read these gains with care, because the debug-step rule and the lessons charge started at the same time. The two worst late defects, C4-2 and C5-S1, came from two of the three sessions that passed 300k. Those sessions also held the largest tickets.

## Coordinator catches

- `bench commitment start` refused at first, because no landed plan bound the deliverable. The binding landed at `9c228393`, and a learning asks each spec stage to bind its deliverable.
- A record amendment after a chunk update was refused as stale. The orchestrator then recorded the amendment first, then the chunk, then the verification, then the review.
- Repair charges first lacked the `/bench-debug` step. After the reviewer named it a requirement, every repair stated its diagnosed cause before its fix.
- The orchestrator charged one wrong design: plants derived from the scanned package list. The repair session proved that design silent and kept a separate plant table with recorded reds.
- A ticket deleted a test file, and `writes-resolve` went red, because the Writes grammar has no deletion marker. The entry carries `(new)`, flagged for reviewer veto.
- The final reconciliation found 55 seam cells that still read "planned". The orchestrator rewrote them as citations, and three axes checked every citation against the tree.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-create-the-shared-refusal-route-registry.md | 1 | coverage |
| 02-move-the-landing-faces-into-the-shared-registry.md | 1 | one-source, coverage |
| 03-prove-each-agent-route-passes-the-wired-guards.md | 2 | one-source, coverage |
| 04-render-the-recovery-matrix-from-the-registry.md | 1 | coverage |
| 05-route-each-merge-refusal-through-the-registry.md | 1 | one-source, coverage |
| 06-give-the-red-source-fold-an-exit.md | 1 | coverage |
| 07-route-each-reset-refusal-through-the-registry.md | 1 | defect, coverage |
| 08-route-the-commit-exit-3-to-the-reset-plan.md | 1 | one-source |
| 09-route-each-commit-refusal-through-the-registry.md | 2 | authority, coverage |
| 10-route-each-checkpoint-refusal-by-its-cause.md | 1 | authority, routing, coverage |
| 11-route-the-commitment-policy-refusals-through-faces.md | 2 | route grammar, one-source |
| 12-print-the-commitment-verb-routes-from-the-registry.md | 2 | authority, coverage |
| 13-refuse-a-route-literal-outside-the-registry.md | 1 | one-source, coverage |

## Agent-experience improvements

### Bench CLI

- Add a deletion marker to the ticket Writes grammar, and make `writes-resolve` and the landing fence both read it.
  Feeds: new
- Make `bench probe` print the first failing test name in its verdict row, and the exit code of the mutated run.
  Feeds: new
- Make `bench probe` accept a line or an occurrence anchor, so that a swap can target one of two identical lines.
  Feeds: new
- Make `bench structure --growth` report the directory file cap and the headroom of a file near its budget.
  Feeds: new
- Add a `bench record` form that appends a version 2 author assignment, so that no orchestrator edits the plan with a script.
  Feeds: new

### Skills

- Make `/bench-implement-spec` charge every ticket author with the recurring review lessons: the shared walk, the authority rule, typed selection, hostile values, and no leftovers.
  Feeds: new
- Make `/bench-write-spec` write each seam cell as a citation that the ticket's author fills, not as "planned".
  Feeds: new
- Make `/bench-write-spec` name each verb that prints a shared face, so that each route fits each verb's grammar.
  Feeds: new

### Process

- Keep the 300k context rule, with the verification-only transfer as its cheap form, and size each ticket to finish under 300k.
  Feeds: new
- Before a ticket adds a printed route to a shared funnel, sweep every caller of that funnel.
  Feeds: new