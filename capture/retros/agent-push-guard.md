# Retro — agent-push-guard

## Outcome

`agent-push-guard` landed as `f2399156` from the reviewed source tip
`571a5f1f` over the frozen base `5044b03d`. The diff held 25 files, 827
insertions, and 73 deletions across 13 commits, and the landing released the
integration source with `census=9`.

The Bash guard now allows a plain push to any branch but the default. It
keeps the denial for a forced push, a deletion, the three broadcast forms,
and an unresolved destination. A bare push reads `push.default` and the
checked-out branch through one owner in the git package. The reference guide
states the rule in its hook-layer list, and an anchor with a canary fixture
pins the sentence.

The review round added five rows. A redirected repository and an `xargs`
prefix fail closed. `@` and a `heads/` prefix resolve the way git resolves
them.

All 46 acceptance rows resolve to a real test. Two build decisions and two
review repairs are recorded in the spec for reviewer veto. One Standards
finding stays open in `reviews/agent-push-guard.md`.

## Gate-stage timings

The landing gate, in milliseconds:

| phase | verdict | elapsed |
|---|---|---|
| gofmt | green | 94 |
| vet | green | 877 |
| test | green | 65512 |
| race | green | 5707 |
| system | green | 24701 |
| shellcheck | green | 574 |

One whole-project gate ran, at the landing. Every merge and commit before it
ran the lane only. The `gitguard` package test grew from 4 seconds to 19
seconds, because the bare-push timeout row pays five 3-second stub calls.

## Ticket-versus-spec-slice and delegate performance

Four ticket charges ran as four Opus delegates, three in sibling worktrees
and one on the integration source. At most two ran at a time under the
test-parallel cap. One fence-extension continuation, two repair charges, three
review axes, and one repair-scoped re-review followed. Every delegate ran at
Opus low or medium.

Every ticket charge verified its premise with citations before an edit. Two
found the premise wrong. The verdict delegate found that its prescribed
self-probe could not red, because the first free arg is never a protected
name. It added the row that made the probe bite. The wire delegate found that
the probe bound row PG31 names does not exist in the destination probes. It
reported the gap instead of a silent pass.

Three of four ticket charges landed their behavior first-pass. The verdict
charge needed one fence extension for two test files the ticket's `Writes:`
line did not name. The repair charge for the review findings landed all
seven targets in one pass with its own probe bitten.

## Coordinator catches

The coordinator ran an independent mutation probe for every accepted
done-claim, at a site and a kind distinct from the delegate's own. Six probes
ran. Five bit as expected. One came back silently green.

A swap of the production checked-out adapter in the guard subcommand left
the subcommand test green. The junction test graded its own copy of the
wiring. One allow row for `git push origin HEAD` closed it.

The coordinator also caught the twice-derived checked-out mapping the wire
delegate reported, and collapsed it into `git.CheckedOutName` before the
review round. Two review preflights went red on fence closure and row
ownership. Each was a ticket-file omission, and each took one commit.

## Repair attribution

| ticket | rounds | cause per round |
|---|---|---|
| resolve-the-bare-push-destination | 0 | none |
| rewrite-the-push-verdict-for-refspecs | 1 | ticket-slicing |
| state-the-push-rule-in-the-reference-guide | 1 | ticket-slicing |
| wire-the-guard-facts-in-the-subcommand | 1 | spec-row |
| repair-the-review-findings | 1 | ticket-slicing |

The verdict ticket's round went to two test files inside the spec fence that
the ticket did not name. The guide ticket's round went to the fixture
directory its `Writes:` line named only by file. The wire ticket's round went
to the checked-out mapping the spec's implementation decision described
twice. The repair ticket's round went to the fixture pins and registries its
fence had to close.

## Agent-experience improvements

### Bench CLI

- The landing's census entry records 9 raw calls in one worktree, all
  delegate reads and one coordinator probe edit. Give `bench gate-prose` a
  lone file operand and `bench canary` a single-fixture form.
  Feeds: new
- Make `bench worktree release` say that a landed branch releases through
  `bench worktree clean --landed`, because the refusal names only the
  release form again.
  Feeds: new

### Skills

- Make `craft-tickets` require a ticket's `Writes:` line to name the
  fixture directory, not only its files, when the ticket adds a canary
  fixture.
  Feeds: new
- Make `craft-spec` require a coverage row that leans on a bound or a
  timeout to name the owning symbol in its seam column.
  Feeds: new

### Process

- Write the repair ticket that cites the amended rows before the
  repair-scoped re-review starts, so the review preflight is green on the
  first run.
  Feeds: none
- Probe a production adapter through the test that grades the production
  binary, not through a junction test that composes its own copy.
  Feeds: none
