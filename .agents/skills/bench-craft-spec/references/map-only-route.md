# Map-only route

## Eligibility

The reviewer selects the map-only route for a ready compiled map.
The full route of `/bench-write-spec` stays the default.
The route accepts only a ready compiled map that `bench maps` reports with no diagnostic.
A reviewer-confirmed conversation and a named reviewed artifact never take the route.
The reviewer judges that each decision is an exact predicate, and no check grades predicate quality.

## Route steps

1. Compile and move the map under the `Entry contract` section of [`/bench-write-spec`](../../../commands/bench-write-spec.md#entry-contract).
2. Charge one fresh slicer on `craft-line`'s ticket slicing line.
   The route has no research fan-out, no spec fork, no spec review, and no paired slicing readers.
3. The slicer writes the thin spec's stories and coverage rows from the decision predicates first, so the row IDs exist before the tickets.
4. The slicer cuts the tickets under `craft-tickets` and [the slicing checks](../../bench-craft-tickets/references/slicing-checks.md), and each ticket's `Covers:` line names declared rows of the thin spec.
   The slicer does its own tree research under the slicing checks' `Fence-expansion evidence` classes.
5. The slicer records the gap log in the thin spec's `Further notes`.
   Each entry takes one class: (a) map discretion, (b) a seam or test choice, (c) a missing fact, or (d) a reviewer decision.
6. After every ticket exists, the coordinator completes the chunk table and the completion plan under an interim fence that holds only the review pickup.
   The coordinator then commits the thin spec and the tickets.
   The coordinator runs `--propose-writes` for each ticket on the committed thin spec, with its chunk table and completion plan.
   The proposal tolerates only the closure-check reds and the `fence-writes` red.
7. One ticket review runs on `craft-line`'s ticket review line, with an iteration cap of 2.
8. After the ticket review accepts, the coordinator writes the fence section, runs the plan-only preflight, and presents the approval table.
9. Each class (d) entry is a decision row in the approval table.
   These rows are not late clarification questions, so the two-question limit does not apply to them.
10. The verification log counts only the ticket round, and its note names the map-only route.

The spec-stage close in the `Exit handoff` section of [`/bench-write-spec`](../../../commands/bench-write-spec.md#exit-handoff) applies unchanged.

## Stop and return

The slicer stops when a map decision is not an exact predicate.
The slicer also stops when the map leaves a code seam open.
The stop returns to the reviewer with a recommendation for the full route, and the slicer writes no further ticket.

## Thin spec

The coverage check and the build preflight read the thin spec's header fields, stories, coverage map, chunk table, `bench-completion-plan` fence, and ownership fences.
The thin spec also holds the gap log for the sign-off, and no check reads it.

- The header holds `Status: staged`, one `Decision source:` line that names the compiled map, a `Verification log:` line, and the `Line:` fields.
- `## User stories` holds at least one story for each map decision.
- `### Acceptance coverage map` includes the `row` column.
- `## Implementation chunks` holds the chunk table and the `bench-completion-plan` fence.
- `## Ownership fences` lists `reviews/<slug>.md` until step 8 writes the full fence.
- `## Further notes` holds the gap log.

`Problem` and `Solution` stay short, and the thin spec omits any section that no check reads and the map already settles.
