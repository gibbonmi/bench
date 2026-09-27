## Outcome

FT337 landed at `9f61740a`. A review charge now publishes the frozen diff as
a prefix, one source for each file patch, and a suffix. An unchanged patch
keeps its page digests across rounds. The implementation phase now states one
order: the author record commit, then the review charge, then the axis
dispatch. A chunk-chain anchor and three workflow fixtures enforce that order.

The build ran in Claude Code. Ticket authors ran on opus at high effort, the
Claude mid binding of the approved line. Every review axis and the RE2 full
control ran on fable at high effort by user direction. RE1 used two repair
cycles. RE2 used two cycles and two reviewer extensions, and then a
destination merge round.

The spec phase took one blocking repair round and one close fold. Opus/high
reviewed the spec twice for USD 2.8950486 in provider cost, and neither pass
ran tests. The first pass found that headerless patch identity was undefined.

The RE2 control comparison found 8 narrow findings in 399474 tokens and 6
control findings in 227918 tokens. The two shapes agreed on 3 findings. The
reviewer kept the narrow shape provisional.

## Gate-stage timings

- landing: commit 9f61740a18ca20993ab6a788397b80c207e70f3a, trace 945b7d211c5de98ae06084cc93e01c2f
- gofmt: 122 ms
- vet: 1182 ms
- test: 127520 ms
- race: 2940 ms
- system: 39661 ms
- shellcheck: 520 ms

## Ticket-versus-spec-slice and delegate performance

Ticket 1 needed two fence expansions: the prose budget row and the seven
fixtures that pin it. Ticket 2 needed four expansions: the version test, the
shared spill helper, the eight paths of the destination merge, and the 13
roadmap fixtures. The ticket 2 focused checks omitted the whole `cmd/bench`
package and the npm pack check, so the RE2 checkpoint went red after a clean
review.

The table holds the labeled claims that a coordinator or axis probe graded.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 1 | an order swap reds the workflow check | verified | 9 | held | claude-opus-5-5 / high / implementation |
| ticket 1 | a later record commit fails current binding | verified | 9 | held | claude-opus-5-5 / high / implementation |
| RE1 cycle 1 R3 | a reverse-order sentence reds the check | verified | 8 | held | claude-opus-5-5 / low / repair |
| RE1 cycle 1 R1 | the Land join is reverted with no word change | verified | 9 | held | claude-opus-5-5 / low / repair |
| RE1 cycle 2 R1 | the cycle 1 R1 work is correct | verified | 9 | held | claude-opus-5-5 / high / repair |
| RE1 cycle 2 R3 | the cycle 1 R3 work is correct | verified | 9 | held | claude-opus-5-5 / high / repair |
| ticket 2 | fragments rebuild every RE5 fixture byte for byte | verified | 8 | held | claude-opus-5-5 / high / implementation |
| ticket 2 | equal-content files stay separate members | verified | 9 | held | claude-opus-5-5 / high / implementation |
| ticket 2 | file selection returns only that file's patch | verified | 8 | held | claude-opus-5-5 / high / implementation |
| RE2 cycle 1 R14 | a trailing space stays path text | verified | 9 | held | claude-opus-5-5 / low / repair |
| RE2 cycle 1 R11 | copy headers are covered at both seams | verified | 9 | held | claude-opus-5-5 / low / repair |
| RE2 cycle 1 R12 | a suppressed blank line is covered | verified | 9 | held | claude-opus-5-5 / low / repair |
| RE2 cycle 2 R16 | one owner builds the canonical entries | verified | 9 | held | claude-opus-5-5 / low / repair |
| RE2 cycle 2 R17 | one owner selects the expected bytes | verified | 8 | held | claude-opus-5-5 / low / repair |
| RE2 cycle 3 | the version test reads the spill file | verified | 9 | held | claude-opus-5-5 / low / repair |
| RE2 cycle 3 | the embed names the baselines literally | verified | 9 | held | claude-opus-5-5 / low / repair |
| RE2 cycle 4 R18 | the version test uses the shared helper | verified | 9 | held | claude-opus-5-5 / low / repair |
| RE2 cycle 4 | the shared helper finds a spill line anywhere | verified | 9 | held | claude-opus-5-5 / low / repair |

The Brier mean is 0.017 over 18 pairs, with 0 abstentions.

## Coordinator catches

- The ticket 1 author ran the charged RE11 probe against the live-tree check, which cannot see a registry omission. The fixture owner proved the row instead.
- The RE2 checkpoint found a response-bound spill and an npm embed glob that the ticket checks never ran.
- The cycle 4 charge named `spilledResponse` as the fold target without a read of it, and the helper could not serve the new shape.
- A `bench worktree merge` from the primary checkout graded prose on paths from the wrong directory, as in FT336. The merge from inside the worktree passed.
- The first landing refused on the destination delta, and the reviewer chose to fence the merged paths.
- In the spec phase, Git fixtures confirmed the inventory mismatch and the missing header shapes before the author defined the descriptors.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| spec authoring | 1 | spec-row |
| 1-order-review-preparation.md | 2 | spec-row, other |
| 2-page-review-diffs-by-file.md | 4 | delegate-error, delegate-error, other, delegate-error |

## Agent-experience improvements

### Bench CLI

- Add a review-record helper that prints the source digest and the plan digest for a tip.
  Feeds: new
- Add an expected verdict and a repeatable package flag to `bench probe`.
  Feeds: new
- Make the prose step of the worktree merge lane read from the composed tree, whatever the caller's directory is.
  Feeds: new

### Skills

- Move the duplicated-facts sweep and one recorded red per independent expectation into the author and repair charges.
  Feeds: new
- Tell a confirming Standards axis to look for duplication in the repair delta first.
  Feeds: new
- Observe headerless and repeated-path output before a spec declares file identity.
  Feeds: none

### Process

- Set a check floor before each ticket commit: root conformance, every written package, and `cmd/bench` for a public response or embed change.
  Feeds: new
- Run the propose-writes preflight for each ticket before graph approval, and fold its fixture paths into the Writes lines.
  Feeds: new
- State one route for a destination merge after the last chunk.
  Feeds: new