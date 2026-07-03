# workflow exits — routes for the paths that currently dead-end

## Problem

Three workflow paths have no defined exit. A build that hits its token cap or
ends with stories unmet has nowhere to route — invariant #2 says "stop and
report," and no phase says what happens next. A spec superseded by a newer one
stays live in `specs/`, reading as a second source of truth (the kit-audit
retirement was done by hand, as a one-off). And `/bench-debug` tells the agent
to add the repro test to the gate and run `bench shift` without saying the test
must be **committed first** — a shift iteration that ends red rolls the worktree
back, destroying an uncommitted repro test.

## Solution

Three prose edits, one per dead end, each pinned by a gate anchor:

- `/bench-implement-spec` gains a "When the build stops short" section routing a
  capped or unmet build by cause (tier, spec, or scope), and its exit handoff
  names that route.
- `/bench-write-spec` gains a retirement step codifying the existing precedent:
  a superseded spec gets a **Superseded by** header plus the historical markers
  so the gate's currency and coverage checks skip it.
- `/bench-debug` states the ordering: commit the repro test into the gate before
  launching the fix shift, with the rollback rationale.

## User stories

1. As a reviewer whose build hit its cap or ended unmet, I want the phase to
   exit through a defined route — state reported, green work kept, cause named,
   next command recommended — so a stopped build is a decision point, not an
   abandoned worktree.
   Line: claude-fable-5 / high. Command guidance prose under the profile's
   leverage override.
2. As a teammate walking in cold, I want a superseded spec visibly retired the
   moment its successor is written, so `specs/` never carries two live sources
   for one feature.
   Line: claude-fable-5 / high. Same edit class.
3. As an agent on the bug path, I want `/bench-debug` to order the repro-test
   commit before the shift, so the first red iteration cannot destroy the loop
   the whole phase exists to build.
   Line: claude-fable-5 / high. Same edit class.
4. As a kit dev, I want each new rule pinned by a prose anchor, so the routes
   survive future edits to the three phase files.
   Line: inline on the session model, low effort. require_anchor pattern.

## Implementation decisions

- Stops-short routing lives in the implement phase (the phase that owns the
  build loop), not in `BENCH.md` — the invariant keeps its one-line "stop and
  report"; the phase owns what happens next. Routes: wrong tier → re-declare
  one tier up per the `craft-line` ladder and resume; wrong spec → back to
  `/bench-write-spec` with the finding; wrong scope → propose a split, reviewer
  decides. Committed green work stays; nothing is squash-finished to fake
  completion.
- Spec retirement codifies the kit-audit precedent (retired header + historical
  markers) as a `/bench-write-spec` step rather than inventing a new mechanism;
  the gate's existing `command-currency` and `coverage-map` opt-outs are the
  enforcement-side halves it composes with.
- The debug ordering is stated where the shift is launched ("How it meets the
  rest of Bench") and names the exception explicitly: the repro-test commit
  deliberately precedes green because it strengthens the oracle — the shift's
  own commits stay on-green as always.

## Testing decisions

- Same seam as the craft-skill specs: the docs-contract anchor list in the
  gate; each anchor observed red before its phase text lands.
- Gate command: `bench gate`.

### Seam diagram

    trigger: bench gate (every shift iteration, every final-check)
        │
        ▼
    .agents/commands/bench-implement-spec.md ──▶ [ docs-contract        ] ──▶ exit 0/1 +
    .agents/commands/bench-write-spec.md     ──▶ [ anchor checks        ]     attributed
    .agents/commands/bench-debug.md          ──▶ [ (require_anchor)     ]     stderr lines
                              ◀ tests attach here: run `bench gate` with the rule
                                absent (red, targeted message), then present (green)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | implement phase carries the stops-short route | docs-contract anchor seam | `bench gate` run red with the anchor added before the phase text | the anchor greps the needle in the file; absence is the reported failure |
| 2 | write-spec phase carries the retirement step | docs-contract anchor seam | same red run (all three anchors added together, observed red together) | same mechanism, distinct per-file message |
| 3 | debug phase orders the repro commit before the shift | docs-contract anchor seam | same red run | same mechanism, distinct per-file message |
| 4 | the three anchors run on every gate | docs-contract anchor seam | already covered once landed — red-capability observed in the story-1 red run | an anchor deleted later goes red at the next run with its file+needle message |

### Edge inventory

- error path — an anchor needle reworded away in a future edit: covered by the
  anchors themselves (story 4).
- empty/absent — a phase file missing entirely: `require_anchor` errs distinctly
  on the missing file, already covered by its helper.
- re-run idempotency / boundary / hostile environment — **Won't handle:** not
  meaningful for prose edits to command files.
- interrupted/partial state — **Won't handle:** a mid-edit tree is graded at the
  next gate run.

## Out of scope

- **Mechanical enforcement of the stops-short route** (e.g. the shift loop
  detecting a cap and emitting the route itself) — a `bench shift` behavior
  change on the loop seam, its own capability. Estimate: ~20 edits, ~10 gate
  runs.
- **A `bench retire <spec>` subcommand** — parser-candidate territory; parked
  by `decisions/second-wave-parsers.md` #1's evidence rule (no observed
  recurring ad-hoc assembly yet; the write-spec step is the manual convention
  that would generate that evidence). Estimate: ~8 edits, ~4 gate runs.
