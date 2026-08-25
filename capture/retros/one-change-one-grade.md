# Retro: one-change-one-grade

## Outcome

The spec landed on `main` as `9c1580b9` (tree `93ec661b`) from the reviewed
pair base `91443a2a`, source tip `8b2223d9`. A worktree `bench commit` now runs
the declared fast lane; the landing stays the one whole-project gate. Five
tickets and one review pickup commit landed on the integration source. The
review found 10 raw findings across three axes and 2 repair targets, both
closed by ticket 05. The repair-scoped re-review was clean.

## Gate-stage timings

The landing gate took 102.3 s. Its stages were gofmt 0.1 s, vet 1.4 s, test
79.3 s, race 5.1 s, system 10.2 s, and shellcheck 0.5 s. Each of the six
worktree commits paid a full gate of about 100 s to 110 s. The installed broker
predates the lane, and it keeps authority until `bench repair` or the release
install. One commit paid a red gate first. The build paid eight full runs.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized, and no charge carried a spec slice. Tickets 01
and 05 ran at sonnet / low and landed first-pass. Ticket 02 and the ticket-04
conformance follow-up ran at opus / medium and landed first-pass. Ticket 03 ran
at opus / medium and took one repair round. The ticket-04 guidance ran at
fable / high and landed first-pass with prose-clean Markdown.

Tickets 02 and 04 ran in parallel with ticket 01, each in its own worktree off
`main` under a disjoint write fence. The coordinator folded each diff onto the
source by patch in `Blocked by:` order. Ticket 04 split into a guidance charge
and a conformance charge, because the conformance check needs ticket 02's
exported symbol.

## Coordinator catches

- Ticket 01 left the `gate-prose` row out of the exhaustive routing registry.
  The package test passed vacuously, and only the live-root run showed the red.
  The coordinator added the one row inline.
- Ticket 03 asserted the gate-script tally where rows OG14 and OG23 named the
  literal `gate: green`; the coordinator routed the wording to the review.
- Every coordinator probe bit at a site distinct from the delegate's probe.
  The sites were the routing row, the exit code, the placeholder substitution,
  and the named Markdown hand-off. Three more were the OG27 anchor, the profile
  `build` row, and the timeout return.

## Repair attribution

| ticket | repair rounds | cause per round |
|---|---|---|
| 01-grade-named-markdown-paths | 0 | none |
| 02-run-a-declared-lane-in-the-gate | 0 | none |
| 03-commit-through-the-lane-authority | 1 | delegate-error |
| 04-state-the-lane-in-the-rules | 0 | none |
| 05-close-the-review-findings | 0 | none |

## Agent-experience improvements

### Bench CLI

- Give `bench worktree release` a planned discard for a dirty worktree whose diff already landed elsewhere, because the guard refuses an agent's discard.
  Feeds: new
- Make `bench spec --help` print the subcommand inventory instead of an unknown-argument error.
  Feeds: new
- Make the `bench status` gate row name why a verdict is stale when the gated tree equals the work tree after a landing.
  Feeds: new

### Skills

- Make the `craft-delegate` charge name the branch-native census's allowed process seams whenever a test fixture shells out.
  Feeds: none

### Process

- Run independent tickets in parallel worktrees off `main` with disjoint fences, and fold each diff by patch in `Blocked by:` order.
  Feeds: none
- Run `bench repair` or the release install after a landing that changes the broker, so the next build's worktree commits run the lane.
  Feeds: none
