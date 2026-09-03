# Retro: roadmap-light-path-fixes-3

## Outcome

The spec landed on 2026-09-03 as commit `a7e94f1b` (source pair
`3809acb4..78f7d086`, census 7). Six roadmap rows closed under one gate.
They are FT276, FT279, FT285, FT202, FT264, and FT201. Seven tickets, one
spec fence amendment, one review round, and one repair round make up the
thirteen source commits. All twenty-eight coverage rows close. Six ask-user
findings wait for the reviewer; the retirement removed the pickup file, so
they live in commit `78f7d086` and in one parked idea.

## Gate-stage timings

The landing gate ran six phases. The timings were gofmt 108 ms, vet 870 ms,
test 52 269 ms, race 2 572 ms, system 20 710 ms, and shellcheck 475 ms. Six merge gates before it
ran the same phases in 85 s to 95 s each. The lane on each ticket commit ran
in under 10 s.

## Ticket-versus-spec-slice and delegate performance

Every charge was a ticket file, not a spec slice. Seven ticket charges and one
repair charge ran on Opus at medium effort. Five ticket charges returned
diff-ready with no continuation: spaced, cache, pdeathsig, share, and sweep.
The derive charge took two continuations. The first was a fence extension
for two registry files the ticket did not name. The second was one package
test for a row whose seam could not see its mutation.

The link charge took one continuation for the same row class. Three charges surfaced a shortfall instead of an out-of-fence
edit: the derive registries, the cache HOME exception, and the share
`Getwd` probe that the wrapper row cannot see. Three review axes and one
repair-scoped re-review ran on Opus at medium; each returned cited findings
and a clean tree.

## Coordinator catches

- LQ2 and LQ5 named existing tests that cannot see the row's failure. Both
  needed a new package test; the coordinator probe found LQ2, and the
  delegate self-probe found LQ5.
- The derive delegate changed the gate subject's error posture on an absent
  root. The coordinator restored the explicit existence leg.
- The cache delegate narrowed the refusal to a declared HOME. The coordinator
  recorded it in the spec for veto instead of a silent widening.
- The review preflight found four cache test files outside the fence and one
  missing fixture closure. Both were spec bookkeeping, fixed in one commit.
- The pickup file tripped the prose lane twice on sentence length before it
  committed.

## Repair attribution

| ticket | repair rounds | cause per round |
|---|---|---|
| derive-the-canonical-path-in-one-leaf-package | 2 | ticket-slicing, spec-row |
| refuse-bench-link-in-the-kit-source-checkout | 1 | spec-row |
| route-a-spaced-field-name-to-the-one-line-refusal | 0 | none |
| refuse-an-unheld-or-invalid-build-cache | 0 | none |
| set-pdeathsig-on-builder-children | 0 | none |
| share-the-purity-census-and-count-process-fixtures | 0 | none |
| sweep-cancel-signal-registrations | 0 | none |
| repair-review-round-one | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- The landing's census entry "roadmap-light-path-fixes-3 census: 7 raw calls on the integration assignment, 16 across eight" names the heads and the Bench forms.
  Feeds: new
- Make the block-bench-follow-on hook parse the segment after `--` as an argv, so a quoted pipe inside a child argument is not a follow-on.
  Feeds: new
- Make `bench test --check` print one failures row per diagnostic and name the binary it consulted when a check is unknown.
  Feeds: new
- Make `bench gate-prose <file>` accept a file operand in the first position, because the intent is unambiguous.
  Feeds: none

### Skills

- Add to `craft-spec` the rule that a row which reuses an existing test names the mutation the author ran against it and its observed red.
  Feeds: new

### Process

- Name every test file a ticket may extend in the spec fence, because the review preflight reads the fence and not the ticket headroom.
  Feeds: none
- Run the prose lane on a pickup file before the commit call, because a review artifact pays the same sentence bound as a spec.
  Feeds: none
