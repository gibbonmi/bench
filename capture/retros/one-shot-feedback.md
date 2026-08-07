## Outcome

Promotion run `2934c73613df69abb3a6db88c9049158ce7b1ff654a5462d8af35cefda807420`
landed candidate `effa6ec33a29a187f14e2c36f57f32ff77c30de7` on `main` as
`d1272891603a2df6b94d5b436a5e8a028a475b7e`. The retained promotion verdict is
green for tree `02992a833cf6734bd4aaf610dc5bb7ab006c267e`, recorded at
`2026-08-07T17:00:36Z`. The build shipped the pre-assignment ticket-breakdown
review, retro repair attribution and drain tally, falsification predicate and
partition questions, and exact-predicate grill close. The spec was retired in
`fbd4916` after promotion.

## Gate-stage timings

The retained terminal projection and gate record preserve the green result and
recorded time, but no per-stage elapsed durations. No timing values are
reconstructed or estimated here.

## Ticket-versus-spec-slice and delegate performance

All five assignments were charged from ticket files rather than from an
unbounded spec slice. Four landed without a repair round. The retro-attribution
ticket needed one composed-review repair after its literal anchor could be
satisfied by explanatory prose and its registry placement and ownership claim
needed correction. There was no direct spec-slice delegate cohort in this build,
so the run provides no like-for-like comparison beyond the ticket-sized result.

## Coordinator catches

Acceptance of candidate `9a9c996d` caught that deleting the required
`## Repair attribution` template heading stayed green because a prose copy
shadowed the anchor. The same pass caught ungrouped registry rows placed after a
grouped block and an ownership sentence that ignored the registry's necessary
enforcement copy. Those findings became the bounded
`repair-ra1-anchor-shadowing.md` ticket before final review and promotion.

## Repair attribution

| ticket | repair rounds | causes by round |
|---|---:|---|
| `add-breakdown-review-moment.md` | 0 | none |
| `add-retro-repair-attribution.md` | 1 | `ticket-slicing` |
| `add-falsification-predicate-questions.md` | 0 | none |
| `add-grill-predicate-close.md` | 0 | none |
| `repair-ra1-anchor-shadowing.md` | 0 | none |

## Agent-experience improvements

### Bench CLI

Have `bench spec build status <slug> --full` render the terminal promotion
commit, promotion tree, evidence identity, recorded time, and retained per-stage
durations. This run's command rendered assignments and review receipts but not
the promotion fields the final-check phase must report.

### Skills

Add an exact-occurrence check to the ticket-breakdown review charge when a red
mutation depends on deleting a literal anchor. This would catch explanatory
prose that can satisfy the same needle before assignment.

### Process

Retain per-stage gate durations in the terminal lifecycle record and require the
coordinator's acceptance pass to exercise each promised literal-deletion red
before final review. The first change makes retro timings evidence-backed; the
second moves this run's repair earlier than composition.
