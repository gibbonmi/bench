# Retro — structural-refactor-pass

## Outcome

`structural-refactor-pass` landed as `e86d344d` from the reviewed source tip
`39698bca` over the frozen base `8eea2d15`. The build ran 34 commits and 10
folds on the integration source, and the landing released it with
`census=7`.

Part 1 folded the six survey findings into one owner each. The `bench
worktree` family reads a leaf table with one root-resolution site. The
landing package owns the lane-to-owner constructor, and it states the
tickets-only rule once over two tree readers. The gate package owns the
kit-source predicate and the kit directory, and the worktree hub no longer
imports adopt. The intent package exports the canonical-path match the
operand resolution reads. The refresh package sits beside its two consumers
with one start-ref entry point.

Part 2 added the growth ratchet the reviewer chose over a whole-tree split.
`bench structure --growth <base>` reds only an over-budget file that gained
lines, and the fast lane runs it on every commit and merge. The census of the
104 over-budget files at the base is the restructure backlog. The whole-tree
debt count fell from 114 issues to 112.

All 63 acceptance rows resolve to a test or a recorded run. Twenty-four
decisions (a) to (x2) are recorded in the spec for veto. No review finding
stays open.

## Gate-stage timings

The landing gate, in milliseconds:

| phase | verdict | elapsed |
|---|---|---|
| gofmt | green | 102 |
| vet | green | 1088 |
| test | green | 67903 |
| race | green | 2863 |
| system | green | 27860 |
| shellcheck | green | 593 |

Eleven whole-project gates ran: ten fold gates and the landing. Every one was
green on its first run. The fold gates took 100,820 to 132,370 ms, a mean of
108,411 ms, and the landing took 103,503 ms.

Cost. This session ran 776 assistant turns, summed from its transcript. It
used 20,914 input tokens, 1,853,421 output tokens, 4,102,701 cache-creation
tokens, and 467,734,638 cache-read tokens. No Codex session took part.

Twenty delegate charges and three resumes reported 2,411,333 tokens in total.
The four census charges used 482,139, the ticket slicer 116,140, and the spec
review 186,522. The nine ticket and split charges used 901,051, the three
review axes 512,206, the repair 123,646, and the re-review 104,413.

The census heads at the landing were `ls=3 sed=2 cd=1 rg=1`. The ten
worktrees held 31 raw calls before the release deleted them.

## Ticket-versus-spec-slice and delegate performance

Eight ticket charges and one split charge ran as nine Opus delegates in
sibling worktrees, two at a time under the test-parallel cap. Four Opus
medium research charges censused the over-budget files before the spec
locked. One Opus medium charge sliced the tickets. One Sonnet xhigh round
reviewed the spec and the tickets. Three Opus medium axes reviewed the
diff, one Opus medium charge repaired it, and one Opus medium re-review
graded the repair.

Every ticket charge verified its premise with citations. Three delegates
reported a spec or charge contradiction instead of an edit outside it. The
kit-source charge grew a file decision (o) forbade, and it offered the
revert. The lane-check charge moved its test to a new file, because the
named neighbour would have crossed 400 lines. The growth-mode charge named
the unspecified flag pair. Two delegates named a pin the charge got wrong
and followed the tree.

Five of eight ticket charges landed their behavior first-pass. The operand
and tickets-only charges each took one continuation, because the
coordinator's probe came back silently green and a row was missing. The
kit-source charge took one continuation for the two-line revert. The
growth-mode charge took the review repair, because the spec's own decision
named the wrong query.

## Coordinator catches

The coordinator ran an independent mutation probe on every accepted
done-claim, at a site and a kind distinct from the delegate's own. Eleven
probes ran and nine bit. Two came back silently green, on the row-side
symlink spelling of the operand match and on the post-reconcile resume of a
close. Both became rows, and the re-probes bit.

The first dogfood run of the new verb over the real artifact caught the one
real growth of the pass. Ticket 2 had pushed the landing file from 388 to
410 lines. The pass took the split the backlog named, and the growth run
against main then printed the ok line.

The review round returned 13 findings, which collapsed to seven repair
targets. One blocked: the lane's growth check never fired, because the
private checkout keeps HEAD and the query diffed `base..HEAD`. The Coverage
axis proved it with a live dry run and an independent repro. The
repair-scoped re-review returned seven findings and one blocker, the
unrecorded dogfood, which the coordinator ran and recorded.

## Repair attribution

| ticket | rounds | cause per round |
|---|---|---|
| declare-the-worktree-leaf-table | 0 | none |
| construct-the-lane-owner-once | 1 | spec-row |
| move-the-kit-source-predicate-into-gate | 1 | other |
| own-the-operand-path-match-in-intent | 1 | spec-row |
| state-the-tickets-only-rule-once | 1 | spec-row |
| move-the-refresh-package-beside-its-consumers | 0 | none |
| add-the-structure-growth-mode | 1 | spec-row |
| run-the-growth-check-in-the-fast-lane | 0 | none |
| repair-the-review-findings | 0 | none |

The lane-owner round went to a file the spec let cross the soft limit. The
kit-source round went to a charge sentence that contradicted decision (o).
The operand round went to a row that pinned the caller's spelling and not
the record's. The tickets-only round went to a resume fixture that broke
before the reconcile. The growth-mode round went to a query the spec spelled
as `base..HEAD`.

## Agent-experience improvements

### Bench CLI

- The landing's census entry records 7 raw calls with heads `ls=3 sed=2 cd=1 rg=1`, and the pass's ten worktrees held 31.
  Feeds: new
- Add a production-or-test column to `bench consumers` and a `--cwd <dir>` option to `bench worktree exec`, so a delegate reads and runs without `sed` or `cd`.
  Feeds: new
- Make the pool-path guard name the matched form and the label of the blocked worktree, so a delegate given a path finds its exec route.
  Feeds: new
- Print the composed tree's lane line in the commit message trailer or the intent record, so a re-review can read which checks a commit ran.
  Feeds: new

### Skills

- Make `craft-tickets` require a lane-check ticket to prove its check through the real lane over a composed tree, never through a shell stand-in.
  Feeds: new
- Make `craft-spec` require the current value of a bench signal in the first status update when a spec inherits a closed decision about that signal.
  Feeds: none

### Process

- Run a new verb over the real artifact at the first boundary after its fold, because that run caught the pass's one real growth.
  Feeds: none
- Record every dogfood run in the spec before the repair-scoped re-review starts, because an unrecorded run is a blocking finding.
  Feeds: none
- Put a new addition in a new file when the touched file is over its budget, because the ratchet reds the growth at commit.
  Feeds: none
