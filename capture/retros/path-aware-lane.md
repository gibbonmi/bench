# Retro: path-aware-lane

## Outcome

FT215 landed at `fc3a4046` on 2026-08-30 from the source pair `7fa3023` to
`a5b7a90b`, eight ticket and repair commits over 22 files. The kit's worktree
commit now grades the raw diff between the expected base tree and the composed
tree. Each composed change carries one or more path classes, and the classes
select checks from the declared lane. The four document families run their
registry checks through `bench test --check <name>` inside the lane. The lane
line reads `lane{outcome=pass,checks=<names>,classes=<classes>}`. The last
three source commits ran through the worktree's own binary and paid only the
checks their paths selected.

The spec flipped to `Status: implemented` at the landing.
`reviews/path-aware-lane.md` landed with five `ask-user` findings for the
reviewer: the structure grant, the mid-build fence amendment, the empty change
list on a merge, a real-build row for the document hop, and the unknown-path
cost. The spec awaits retirement after those decisions.

## Gate-stage timings

The landing gate ran gofmt in 98 ms and vet in 787 ms. The test phase took
48.4 s, race 4.7 s, system 14.9 s, and shellcheck 452 ms. The two sibling folds
through `bench worktree merge` each ran the same six phases in about 74 s. A
selective lane commit on Markdown alone ran in about 1.4 s. That time excludes
the run binary build.

## Ticket-versus-spec-slice and delegate performance

Five ticket charges ran on Opus, one at low and four at medium effort, and all
five landed first-pass on behavior. Each returned red-to-green evidence per
row and a self-probe that bit. Two charges reported an out-of-fence file
instead of an edit: the live-tree test inventory in `tier_test.go`, and the
`lane_record_test.go` call site of a deleted fixture. No charge received a
spec slice; the ticket files carried the rows, so the comparison has one arm.

The three review axes at Opus medium returned 17 raw findings that collapsed
to six repair targets and five reviewer decisions. One repair charge closed
the six targets, and its continuation closed the one citation defect the
scoped re-review found.

## Coordinator catches

The coordinator ran five independent probes at distinct sites and kinds, and
every probe bit. The review preflight reddened on the `tier_test.go` fence
gap. The landing reddened on the review pickup file outside the fence. Each
took one spec amendment.

The prose lane refused four sentences in the
coordinator's own review artifact and repair ticket. The Coverage axis's worst
finding was refuted by two live runs of the document hop. One changelog line
ran past the file's wrap width after the repair.

## Repair attribution

| ticket | repair rounds | cause per round |
|---|---|---|
| derive-composed-changes | 0 | none |
| export-embed-target-derivation | 1 | spec-row |
| select-kit-lane-checks | 1 | spec-row |
| add-document-family-checks | 1 | ticket-slicing |
| record-the-selection-in-guidance | 1 | spec-row |
| repair-review-findings | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- Make `bench worktree exec <target> --script <file>` run a multi-step probe inside the venue, because each probe cycle now costs one call per step. The census entry `path-aware-lane census 241` carries the heads.
  Feeds: new
- Make `bench gate-prose` default its root to the resolved worktree when it runs under `bench worktree exec`, because every prose-lane call spells the same absolute path.
  Feeds: new
- Make `bench worktree land` keep the census record until the landing report prints, because the release deleted it before the heads were read.
  Feeds: new

### Skills

- Make `craft-tickets` add the review pickup file `reviews/<slug>.md` to the spec fence when it slices a spec, because the landing refused it as unauthorized.
  Feeds: new
- Make `craft-spec` fence the live-tree test inventory for a ticket that adds a live-tree test, because PL37's charge reported it out of fence.
  Feeds: new

### Process

- Run `bench coverage --check` and a test-name resolution after a repair renames a test, because the structural check passed a row that cited a deleted name.
  Feeds: none
- Run the review preflight from inside the integration worktree, because the primary checkout derives the wrong tip.
  Feeds: none
