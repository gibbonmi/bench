# Ref inventory review record

## RI-C1a author evidence

The plan amendment `2d3535b6` declared the version 2 execution plan on the
`main` tip `72a749a3`, which is the RI-C1a base. The serial order is ticket 6,
then ticket 1. The chunk tip is `c6b2a74e`.

Ticket 6 had a fresh `bench-writer` author on opus at high effort. The author
started at `2d3535b6` and committed `f3ff2543` on a lane pass with the build
preflight green. The delta is two glossary lines in `CONTEXT.md`. The eight
anchors kept their needles and lines. Rows RI57 and RI82 are review-owned, so
they carry no test red; the author reported both as claimed at confidence 9.
The orchestrator read the diff, confirmed a clean tree, and ran the coordinator
probe: an omission of the reader-sweep anchor line under
`docs-currency-workflow`, verdict `bit`, restore `yes`.

Ticket 1 had a fresh `bench-writer` author on opus at high effort from the plan
commit `e3b00c6d`. The author stopped before its commit on a behavioral spec
contradiction: ancestry is transitive, so a ref one commit under an
ancestry-landed ref is itself landed, and RI9 as written was unreachable. The
omission probe of that holder kind was an equivalent mutant. The orchestrator
routed the question to the reviewer through `codex exec` on `gpt-6-astra` at
xhigh reasoning, by user direction on 2026-09-27. The decision was option B:
drop the ancestry-landed holder kind everywhere. The plan commit `e790e56a`
amended story 9, story 63, RI9, RI82, the holder rule, ticket 1, and ticket 6.

The ticket 1 charge evidence `sha256:d03d7a73…` predates that amendment. The
author read the amended spec and ticket bytes from the tree at `e790e56a`,
because the dirty source refused a regenerated charge. The same author session
then committed `2ace6d83` on a lane pass with the build preflight green, and it
reported all 34 rows verified. The named probe, subsumed tested before landed,
bit six cases including RI62 and RI9, restore `yes`. The orchestrator ran the
coordinator probe: an omission of the retained-row skip in the unclaimed apply
loop over `TestCleanUnclaimed`, verdict `bit` with four tests red, restore `yes`.

Ticket 6 reopened for the glossary `holder` sentence. The same author session
committed `c6b2a74e` on a lane pass with the build preflight green. The eight
anchors kept their needles and lines, and the `holder` entry names only a
recorded assignment branch or a unique root.

Two author-discipline notes carry no tree defect. The ticket 1 author wrote one
test block into the pool path with a shell heredoc before it moved to the Edit
tool. It also made one string swap with `sed -i`. The committed bytes passed
the lane and the checks. The ticket 1 charge named `ledger.go` as the
shift-prefix source. The constant lives in `internal/intent/ledger/validate.go`.

### Verification

Each author reran its plan checks at the chunk tip `c6b2a74e` on a clean tree.
The JSON payload holds each result. The `internal/conformance` run skipped the
same three capability tests as the baseline: two need a unix socket and one
needs a character device. The first ticket 1 conformance rerun failed because
the orchestrator's untracked draft of this record was over the prose bound at
that moment. The draft was repaired, and the retained rerun passed.

## RI-C1a chunk review, round 1

The frozen pair is base `72a749a3` and tip `c6b2a74e`. The record commit
`2e460e05` follows the tip, and the shared evidence is `sha256:5f24a383…`.
Each axis ran in a fresh `bench-reviewer` session on opus at high effort, the
conditional review line for an opus implementation. Only the Coverage axis ran
probes, and it left the tree clean.

Raw findings: Standards 6, Spec 2, Coverage 4. Repair targets after the fold:
one ticket 1 repair session with nine items, and two spec wording repairs by the
orchestrator. Two questions went to the reviewer route, `codex exec` on
`gpt-6-astra` at xhigh, and both answers are recorded below as decisions.

## Standards

Findings: 6. Worst: R1.

- R1 (auto-fix, confidence 7): AGENTS.md test-expectation exception, "someone
  must record and demonstrate that red". `clean_classes_test.go:49-50`,
  `clean_unclaimed_test.go:67`, and `clean_unclaimed_test.go:301` hand-spell
  `class=`, ` holder=`, and `retained: content main lacks`, which production
  owns at `clean_classes.go:27-29`. The record holds no red for those
  spellings. The repair records one probe red per spelling.
- R2 (auto-fix, confidence 5): craft-comments register. The comment at
  `clean_classes_test.go:44-46` argues its own case and claims a red no record
  holds. It folds into R1.
- R3 (auto-fix, confidence 5): AGENTS.md one source per fixture harness.
  `squashIntoMain` at `clean_classes_test.go:20-26` repeats `squashLand` at
  `squash_landed_test.go:30-36`.
- R4 (auto-fix, confidence 4): craft-comments one source per fact. Three
  comments restate that a faulted set has no fingerprint and no apply:
  `clean_classes.go:32-33`, `clean_unclaimed.go:33-35`, `clean_unclaimed.go:124-126`.
- R5 (auto-fix, confidence 3): craft-comments contract in one place. The doc
  comment at `clean_classes.go:79-84` restates the holder rule the glossary and
  the spec own; the why clauses stay.
- R6 (no-op, confidence 3): the pool-path heredoc and the `sed -i` swap of the
  ticket 1 author. No tree defect, and the committed bytes passed the lane. The
  orchestrator captured the bypass as a learning.

Advice: the unclaimed row fields are built by hand at four sites, unchanged in
count by the diff. The active-or-cleanup-pending check is the fourth hand-written
one in the package. `reachableFromAHead` reads oddly.

## Spec

Findings: 2. Worst: R7. Every RI-C1a row is delivered, and the holder rule
matches both reviewer decisions.

- R7 (auto-fix, confidence 8): row RI5 read "one commit", but a one-commit
  squash lands by patch containment, so a classifier that stops at cherry prints
  `landed`. The author's two-commit fixture follows story 5. The reviewer route
  confirmed the row now reads "two commits".
- R8 (auto-fix, confidence 4): row RI13 and story 13 read "every unclaimed
  row", but an error row carries no class prefix by RI11. The reviewer route
  confirmed both now read "classified row".

Advice: 32 rows still carry the `planned` seam although their tests exist. The
orchestrator updates the seam cells at the chunk close. The axis also asked
whether a content-landed recorded active branch should hold; that question
became decision D2 below.

## Coverage

Findings: 4. Worst: R9.

- R9 (ask-user, confidence 9): an unrecorded symbolic ref in a Bench namespace
  is selected, classifies through its target, and the apply's `update-ref -d`
  dereferences it. Probes 4 and 5 deleted a unique root and the checked-out
  `main`. No row decided symrefs. Decision D1 below routes it to an error row
  inside ticket 1's fence, and the landing prune's shared delete is captured
  as an idea.
- R10 (auto-fix, confidence 9): a probe that dropped `class` from the
  fingerprint parts stayed silent, because the RI63 fixture changes both the
  class and the holder. New row RI85 pins a class-only change.
- R11 (auto-fix, confidence 8): a probe that admitted every record state as a
  holder stayed silent. New row RI86 pins a complete record's branch.
- R12 (auto-fix, confidence 8): a probe that limited holders to active records
  stayed silent. New row RI87 pins a cleanup-pending holder.

Probes, each with restore `yes`:

- swap of the fingerprint class part: silent
- swap of the holder filter to active only: silent
- swap of the holder filter to every state: silent
- a symref fixture at a unique root: bit
- a symref fixture at `main`: bit

Advice: `PruneLandedBranches` shares the dereferencing delete, outside this
spec. `LandedInDefault` takes the default branch as a short name.

## RI-C1a decisions and repair cycle 1

The reviewer route decided D1: a symbolic ref produces an error row, never a
silent exclusion, and `DeleteBranchExact` stays unchanged in FT199 (row RI83).
It decided D2: a recorded assignment branch that is landed by content only is
not a holder (rule sentence and row RI84). It confirmed R7 and R8. The plan
commit that follows this record carries RI83 to RI87, the wording repairs, and
the repair session assignment.

One fresh repair session for ticket 1 on opus at high consumes repair cycle 1
of 2. It owns R1 to R5, R9 to R12, and D2. The
implementation command contributed to R10, because the author's probe list
named a fingerprint-binding swap without a field. The orchestrator captured the
per-field probe rule as a learning.

## RI-C1a repair cycle 1

A fresh `bench-writer` repair session `ri-t1-repair-1` on opus at high took R1
to R5, R9 to R12, and D2 from the plan commit `b0701a02`. It committed
`4b80686a` on a lane pass in one attempt. The chunk tip moves to `4b80686a`.
The plan digest moves with the amended spec and ticket. The delta touches the
class function, its table test, the unclaimed planner, and its test.

The class function runs `symbolic-ref --quiet` before `rev-parse`, and a symref
gets the fault `<ref> is a symref to <target>`. Every landed recorded branch is
excluded from the holders, not only a content-landed one. That widening is
behavior-neutral: a ref beneath an ancestry-landed branch is itself landed, and
the landed class wins, so the spec sentence "No landed ref is a holder" holds
for recorded branches too. The fault error is renamed to
`errFaultedUnclaimedRef`, because the old message named only a missing commit.

A new helper `recordedBranch` plants a record in a named state, and RI8 folds
onto it. `squashIntoMain` is gone in favour of `squashLand`.

### Repair probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`.

| File | Mutation | Check or test | Verdict |
|---|---|---|---|
| `internal/worktree/clean_unclaimed.go` | swap: drop the class part of the fingerprint | `TestCleanUnclaimedStaleClassRefusesTheOldPlan` | bit, RI85 |
| `internal/worktree/clean_classes.go` | swap: holder filter admits every state | `TestClassifyUnclaimedRefsOverTheEdgeInventory` | bit, RI86 |
| `internal/worktree/clean_classes.go` | swap: holder filter admits active only | the class table | bit, RI87 |
| `internal/worktree/clean_classes.go` | swap: a landed recorded branch holds | the class table | bit, RI84 twice |
| `internal/worktree/clean_classes.go` | swap: `class=` to `kind=` | `./internal/worktree` | bit, 17 tests |
| `internal/worktree/clean_classes.go` | swap: ` holder=` to ` owner=` | `./internal/worktree` | bit, 6 tests |
| `internal/worktree/clean_classes.go` | swap: the retained detail text | `TestCleanUnclaimedPlanNamesApplyOnlyWhenARowRemoves` | bit |
| `internal/worktree/clean_classes.go` | swap: the symref fault text | `./internal/worktree` | bit, RI83 four times |
| `internal/worktree/clean_classes.go` | swap: the symref check never fires | `./internal/worktree` | bit, RI83 four times |
| `internal/worktree/clean_classes.go` | swap: reverse the recorded-holder order | `./internal/worktree` selected tests | silent |

The three label swaps are the recorded reds that R1 required for the
independent expectations `class=`, ` holder=`, and `retained: content main lacks`.
The last row is the orchestrator's coordinator probe. It is silent because no
test reaches one ref from two recorded holders, so the spec sentence on the
lexically first holder has no row. The code sorts, so no behavior is wrong. The
confirming round's Coverage axis grades that gap.

### Verification after the repair

Each author reran its plan checks at the repair tip `4b80686a` on a clean tree.
The ticket 1 rows carry the repair session as their performer, because the plan
names it as the ticket's effective author. The JSON payload holds each result
after the round 1 rows.

## RI-C1a chunk review, round 2

The confirming round read the repair delta `c6b2a74e..4b80686a` with the
whole chunk as context. The record commit `f9959ec9` follows the tip, and the
shared evidence is `sha256:5cf4cc47…`. Each axis ran in a fresh
`bench-reviewer` session on opus at high effort. Only the Coverage axis ran
probes, and it left the tree clean.

Every round 1 fold is confirmed. Standards confirmed R1 to R6, Spec confirmed
R7, R8, D1, and D2, and Coverage confirmed R9 to R12 with D2.
Each axis also judged the wider holder exclusion behavior-neutral.

Raw findings: Standards 4, Spec 3, Coverage 1. Repair targets after the fold:
two code items for one repair session, three spec items for the orchestrator,
and one reviewer decision.

### Standards, round 2

Findings: 4. Worst: R13.

- R13 (auto-fix, confidence 7): BENCH.md plan-expansion rule. The RI-C1a
  chunk table cell at the spec omits RI83 to RI87, while ticket 1's `Covers:`
  and this record's `acceptance_rows` list them. The orchestrator's plan
  commit owns the fix.
- R14 (auto-fix, confidence 4): craft-comments truth. The comment at
  `clean_classes.go:76-77` gives the content-landed reason for every landed
  recorded branch. The ancestry-landed half has a different reason.
- R15 (auto-fix, confidence 4): AGENTS.md one source per fixture rule. The
  owner-letter-to-ref rule is built at `clean_classes_test.go:178`,
  `clean_set_apply_test.go:392`, and `clean_unclaimed_test.go:318`. One
  helper owns it after the repair.
- R16 (no-op, confidence 3): the RI83 `identities` closure beside `refsUnder`.
  The closure reads the `%(symref)` column that `refsUnder` lacks, and the
  Coverage axis proved the closure load-bearing at its line 204. So it is not
  a pasted harness.

Advice: a why clause for the RI85 fixture comment, "in the named state" in the
`recordedBranch` comment, the error-variable comment repeats its message, a
package-wide `markState` helper, and the per-record landed-proof cost.

### Spec, round 2

Findings: 3. Worst: R13, shared with Standards. Rows RI5, RI9, RI13, and RI83
to RI87 are delivered, and the renamed error message conflicts with no spelling.

- R17 (ask-user, confidence 5): a dangling symref never reaches the classifier,
  because `for-each-ref` skips a broken ref and `Output` drops stderr. The
  Coverage axis confirmed it by a run, as its C1 below. Decision D3 resolves it.
- R18 (auto-fix, confidence 4): the plan-only cost sentence omitted the landed
  proof the classifier now runs per active or cleanup-pending record. The
  reviewer route confirmed the replacement sentence.

Advice: `LandedInDefault` failing on a recorded branch aborts the plan, and
status then drops its count silently. That is fail-closed and has no row.

### Coverage, round 2

Findings: 1. Worst: R17, the dangling symref, observed by a fixture probe: the
plan exits 0, prints no row for it, and offers the apply. Nothing deletes
through it, and a target that appears later trips the stale refusal.

Probes, each with restore `yes`:

- rerun of the landed-recorded-holder swap: bit, both RI84 cases
- swap: the symref check limited to the assignment prefix: bit, both shift cases
- omission of `PutAssignment` in `recordedBranch`: bit, RI86
- a dereferencing delete inserted before each apply in the RI83 test: bit, four
  cases at the closure line
- a dangling symref fixture: bit, the plan exits 0 with no row (R17)
- a symref chain in the shift namespace: silent, the outer row faults and both
  applies refuse

The axis graded the orchestrator's silent recorded-holder-order probe as
optional advice for a hardening cycle, under the bounded repair policy's
preference rule: the sentence is binding but does not fail, and no row covers
two recorded holders that reach one ref.

### RI-C1a decisions and repair cycle 2

The reviewer route decided D3. The symref rule narrows to a resolving symref,
and a dangling symref is a Won't-handle case. The cost sentence gains the
per-record landed proof. The plan commit that follows this record carries the
chunk table cell, the D3 wording, the cost sentence, and the source-row table.
It also carries the repair session assignment.

One fresh repair session for ticket 1 on opus at high consumes repair cycle 2
of 2. It owns R14 and R15. The allowance is then
exhausted for this chunk, so a later blocking finding returns to the reviewer.
No implementation command change is necessary for this round; R13 is an
orchestrator plan-step gap, captured as a learning.

## RI-C1a repair cycle 2

A fresh `bench-writer` repair session `ri-t1-repair-2` on opus at high took R14
and R15 from the plan commit `b24187e9`. It committed `09611663` on a lane pass
in one attempt. The chunk tip moves to `09611663`, and the plan digest moves
with the amended spec and ticket. The delta is one comment sentence in
`clean_classes.go` and one test helper, `unclaimedBranchRef`, that four sites
now call. The author found a fourth site beyond the three the charge named. It
folded that site the same way, with no line added to the file at its budget.

The charged probe, a shorter owner repeat inside the helper, stayed silent by
construction. Every consumer reads the helper's value for the fixture and the
expectation. The author's second swap moved the ref out of the Bench
namespace and bit thirty tests across every call site, with restore `yes`. The
orchestrator read the diff, confirmed a clean tree, and ran the build preflight
green. The repair allowance of this chunk is now spent.

### Verification after repair 2

Each author reran its plan checks at the repair tip `09611663` on a clean tree.
The ticket 1 rows carry the repair 2 session as their performer. The JSON
payload holds each result after the earlier rows.

## RI-C1a chunk review, round 3

The confirming round read the repair 2 delta `4b80686a..09611663` and the plan
delta `b0701a02..b24187e9`, with the whole chunk as context. The record commit
`1a0ae536` follows the tip, and the shared evidence is `sha256:b898f652…`. Each
axis ran in a fresh `bench-reviewer` session on opus at high effort. The charge
limited each axis to findings above the blocking bar, because the chunk's two
repair cycles are spent. Only the Coverage axis ran probes, and it left the
tree clean.

Every axis returned zero blocking findings. Standards confirmed R13, R14, R15,
R16, and R18, and found no duplicated knowledge in the repair 2 delta. Spec
confirmed R13, R17 with D3, and R18, and re-verified every RI-C1a row at the
tip by test name. Coverage confirmed R15 with four call-site unwrap probes, each
`bit` with restore `yes`, and confirmed that the repair 2 delta changes no
production behavior.

Advice, with no finding id, which the orchestrator carries to the final
reconciliation:

- The cost sentence still omits one ancestry check per recorded holder per
  open ref.
- A branch set to an annotated tag object fails closed at the apply. This
  predates the chunk.
- A fork under two unique roots has no row. The code names the lexically first
  root.
- The RI73 fixture never asserts that its landed row is present.
- A short name ambiguous with a tag drops the branch from the sweep. That
  excludes and never deletes.
- The `holder` glossary entry states the recorded kinds without the
  content-landed exclusion.

RI-C1a is complete at tip `09611663`. Repair cycles consumed: 2 of 2. The
hardening allowance is unused.

## RI-C1b author evidence

The RI-C1b base is the accepted RI-C1a tip `09611663`. Ticket 2 had a fresh
`bench-writer` author on opus at high effort from the plan commit `c42dd50b`.
That commit cited the RI-C1a seams and recorded the assignment. The author committed
`0e8f29cb` on a lane pass in one attempt.

The build preflight at that tip was red on `base-current` alone, because the
light path `precedence-trace` landed on `main` during the RI-C1a close. The reviewer route deferred the `main`
composition to the landing preparation, with its own review round. The
explicit-base preflight with `--base 72a749a3` is green at the tip.

The exported ref list became `CountUnclaimedRefs`, which returns the landed,
subsumed, unique, and faulted counts from the same planner, and both callers
moved to it. The status git row prints the three class details through the
plural helper, omits a zero class, and keeps the dirty and unpushed details. It
routes to the plan command whenever any count is above zero or the planner
fails. The author added a `faulted ref` detail and an `unclaimed refs
unavailable` detail so a fault is never silent.

The reviewer route accepted both details as row RI88. The plan commit
`f1204694` carries the sentence, the row, the ticket 2 `Covers:` line, and the
RI-C1b seam citations. The system route test gained a landed case that runs the
printed apply command. The 43-ref plan time is about 550 ms over three runs, and
a full `bench status` over that fixture is about 850 ms. The author measured a
scratch repository with a worktree build, and no test asserts a bound.

### Author probe verdicts

Each probe ran through `bench probe` on `internal/status/status.go`, and each
restore reads `yes`.

| Mutation | Check or test | Verdict |
|---|---|---|
| swap: a zero class prints | the producible table | bit, 7 cases including RI27 |
| omission: the subsumed detail | the producible table | bit, RI24 |
| swap: the dirty path wins over the plan route | the producible table | bit, RI26 |
| omission: the faulted detail | the producible table | bit, the mixed case (RI88) |

The system landed case took a copy-aside mutation, because the system suite is
not a probe target. The route test made status route a landed ref to `git
push` and went red, and `cmp` confirmed the restore. The faulted-detail probe
was silent until the author added the symref fixture to the mixed case.

### Verification

The author reran its three plan checks at the tip `f1204694` on a clean tree.
The JSON payload holds each result under the RI-C1b chunk.

The author kept two spellings of the plan command, the status action literal
and `unclaimedReplan`, because `axi.InvocationArgument` has no exported string
renderer. The round 1 review graded that keep decision as R19 below.

## RI-C1b chunk review, round 1

The frozen pair is base `09611663` and tip `f1204694`. The record commit
`c1240476` follows the tip, and the shared evidence is `sha256:13d67d1f…`.
Each axis ran in a fresh `bench-reviewer` session on opus at high effort. Only
the Coverage axis ran probes, and it left the tree clean.

Raw findings: Standards 5, Spec 2, Coverage 3. Repair targets after the fold:
one ticket 2 repair session with seven items, three new coverage rows, and one
spec sentence by the orchestrator. Two questions went to the reviewer route,
and the answers are decisions D4 and D5 below.

### Standards

Findings: 5. Worst: R19.

- R19 (ask-user, confidence 6): AGENTS.md one source per fact. The status
  action literal at `status.go:133` and `unclaimedReplan` build the same plan
  command. `axi` has no string renderer, but the worktree package can export
  the spelling without an import cycle. Decision D5 folds it.
- R20 (auto-fix, confidence 6): finding discipline, evidence only in a
  delegate return. The keep decision for R19 was absent from this record. The
  paragraph above now holds it.
- R21 (auto-fix, confidence 7): craft-comments aging. The `appendGit` doc
  comment names parameters that do not exist and says nothing of the class
  counts the body now carries.
- R22 (auto-fix, confidence 5): one source per fact. The class set is
  enumerated four more times: the count fields, `Rows()`, the count switch,
  and the hand-typed detail table. The repair returns ordered name and count
  pairs from the class constants and renders them once.
- R23 (auto-fix, confidence 4): repeated switches. The count switch sends an
  unknown class to `Unique` while the action switch sends it to
  `discard-remove`. The repair matches `classUnique` explicitly.

Advice: the deleted landed-posture test is covered by the RI24 case and the
landed system subtest. The `--apply` shape check in the system test has no
recorded red but follows precedent. The promised-labels list omits the two new
details.

### Spec

Findings: 2. Worst: R24.

- R24 (auto-fix, confidence 7): "the counts match the plan". A resolving
  symref whose target is unique is counted as a faulted ref and again as a
  unique branch. The subtraction at `status.go:652` removes only subsumed and
  unique. Traced: `1 faulted ref, 2 unique branches` where `1 unique branch`
  is right. New row RI89 pins it.
- R25 (ask-user, confidence 6): "Status shows a nonzero `<n> faulted ref`
  count". A blob-tip Bench ref fails `LandedState` first, and the row reads
  `git state unavailable` with `git status`. Decision D4 resolves it.

Advice: the 43-ref timing lives in this record, not in the ticket file; the
`unclaimed refs unavailable` detail had no test.

### Coverage

Findings: 3. Worst: R25, confirmed by a fixture probe that turned the mixed
case into a blob ref.

- R25 (ask-user, confidence 9): as above, with the probe evidence.
- R24 (auto-fix, confidence 8): as above, confirmed by a fixture probe that
  pointed the symref at a unique ref: `2 unique branches`.
- R26 (auto-fix, confidence 8): the planner-failure detail and route have no
  test, and a faulted-only set has no test. Three mutations stayed silent:
  `Rows()` without `Faulted`, the route without the planner error, and the
  detail omitted. New rows RI90 and RI91 pin them.

Probes, each with restore `yes`:

- `Rows()` without `Faulted`: silent
- the route without the planner error: silent
- the detail omitted: silent
- the mixed case with a blob ref: bit (R25)
- the mixed case with a symref to a unique ref: bit (R24)

Advice: on a planner failure the Bench refs fall into the unique-branch count;
the dashboard width over eight details is unmeasured.

### RI-C1b decisions and repair cycle 1

The reviewer route decided D4. `appendGit` tolerates a `LandedState` error, it
still counts the unclaimed rows, and it routes to the plan with the details.
The row keeps `git state unavailable` beside them. It decided D5: the worktree
package exports the plan command spelling, and status and ticket 5 read it.
The system test literal stays as an independent expectation with its recorded
red. It supplied rows RI89, RI90, and RI91.

The plan commit that follows this record carries the D4 sentence, the three
rows, and the ticket 2 `Covers:` line. It also carries the chunk table cell
and the repair session assignment.

One fresh repair session for ticket 2 on opus at high consumes repair cycle 1
of 2. It owns R21 to R26 and D4 and D5. No implementation command change is
necessary; R20 is an orchestrator record gap, and the record now carries the
keep decision.

## RI-C1b repair cycle 1

A fresh `bench-writer` repair session `ri-t2-repair-1` on opus at high took R21
to R26, D4, and D5 from the plan commit `80b31424`. It committed `f0a8a3c6` on
a lane pass in one attempt, and the explicit-base build preflight is green at
that tip. `appendGit` no longer returns early when `LandedState` fails. The
unique-branch count now works by ref identity, so a faulted symref counts once.
`UnclaimedPlanCommand` is the one exported spelling of the plan command, and the
status action table reads it. The class names come from one enumeration built
on the class constants, and an unknown class is an error, not a unique count.

The author's probes each bit with restore `yes`:

- `Faulted` dropped from the row total: RI91 and the blob case
- the planner-error route removed: RI90
- the detail omitted: RI90
- the identity check removed: three cases with RI89
- the early return restored: the blob-tip case
- a modifier flag swapped in the exported spelling: six status cases, and both
  system subtests through a copy-aside run
- `classUnique` dropped from the enumeration: three cases

The author found one edge: a non-repository fails both reads, and the exact
`git unavailable` case expects `git state unavailable` with `git status`. The
author kept that row and routed only rows the planner still finds. The
reviewer route decided D6: a repository with a planner error beside a Git state
failure still routes to the plan, a non-repository keeps today's row, and row
RI92 pins a blob-tip ref beside an unreadable ledger. A second repair session
implements D6 as repair cycle 2 of 2.

## RI-C1b repair cycle 2

A fresh `bench-writer` repair session `ri-t2-repair-2` on opus at high took D6
and RI92 from the plan commit `f56fad23`. It committed `c6d2cfbf` on a lane
pass in one attempt, and the explicit-base build preflight is green at that
tip. `appendGit` drops a planner error only when `git.CommonDir` also fails,
which means the directory is not a repository. The RI92 case prints
`git state unavailable, unclaimed refs unavailable` with the plan action, and
the exact `git unavailable` case is unchanged. The producible file stays at 422
lines through two shared fixture closures.

The author's probe swapped the distinction back to repair 1's behavior and bit
the RI92 case alone. Its counter-probe forced the route on and bit the
`git unavailable` case alone. Both restores read `yes`. The orchestrator ran the
coordinator probe: an omission of the faulted-detail append, verdict `bit` on
three cases, restore `yes`. The chunk tip moves to `c6d2cfbf`, and the plan
digest moves with the amended spec and ticket. The repair allowance of this
chunk is now spent.

### Verification after repair 2

The author reran its three plan checks at the repair tip `c6d2cfbf` on a clean
tree. The JSON payload holds each result after the round 1 rows.

## RI-C1b chunk review, round 2

The confirming round read the two repair deltas `f1204694..c6d2cfbf` with the
whole chunk as context. The record commit `6e93138e` follows the tip, and the
shared evidence is `sha256:d671257b…`. Each axis ran in a fresh
`bench-reviewer` session on opus at high effort. The charge limited each axis
to findings above the blocking bar, because the chunk's two repair cycles are
spent. Only the Coverage axis ran probes, and it left the tree clean.

Every axis returned zero blocking findings. Standards confirmed R19 to R23 with
D4, D5, and D6, and found no duplicated knowledge in the repair deltas. Spec
confirmed R24 to R26 with D4, D5, and D6, and re-verified all eleven RI-C1b
rows at the tip by test name. Coverage confirmed the same folds with two
fixture-source probes and one producer probe of a new kind and site, each
`bit` with restore `yes`. It also reran the orchestrator's faulted-detail
omission probe with the same result.

Advice, with no finding id, which the orchestrator carries to the final
reconciliation:

- Ticket 5 does not yet tell its author to read `UnclaimedPlanCommand`. The
  RI-C2b plan commit adds that line.
- An unborn repository, or one with no resolvable default branch, now routes
  to a plan that fails with the same cause. The spec sentence requires the
  route, and no row pins it.
- A dirty path beside a blob-tip ref loses its dirty count, because the Git
  state read fails as a whole.
- `UnclaimedPlanCommand` re-spells the `worktree clean` words and the
  modifier order that `cleanArguments` owns, an honest repetition under D5.
- `rev-parse --absolute-git-dir` is spelled at four older sites in the status
  test package.
- The `unreadableLedger` closure writes an empty ledger, and the git dir it
  writes equals the common dir only in a primary checkout.

RI-C1b is complete at tip `c6d2cfbf`. Repair cycles consumed: 2 of 2. The
hardening allowance is unused.

```bench-review-record
{
  "version": 2,
  "spec": "specs/ref-inventory/spec.md",
  "plan_digest": "sha256:74f8fb4685ed7e515b51757c2cd9d36d59a79ba8659977672e3b3aad337a118a",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RI-C1a",
      "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
      "tip": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
      "plan_digest": "sha256:084df113f531adbc7d196def8c112c45a1925608d2a1e45099e9428dbb5ba82b",
      "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
      "acceptance_rows": [
        "RI1", "RI2", "RI3", "RI4", "RI5", "RI6", "RI7", "RI8", "RI9", "RI10", "RI11", "RI12", "RI13", "RI14", "RI15",
        "RI17", "RI18", "RI19", "RI20", "RI21", "RI22", "RI23", "RI25", "RI28", "RI55", "RI57", "RI59", "RI60", "RI61",
        "RI62", "RI63", "RI66", "RI72", "RI73", "RI81", "RI82", "RI83", "RI84", "RI85", "RI86", "RI87"
      ],
      "verification": [
        {
          "id": "ri-c1a-6-anchors-r1",
          "performer": "claude:bench-writer/ri-t6-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t6-author-20260927/6-anchors@c6b2a74e",
            "digest": "sha256:80be124f957a7a1d26cd487fc038f3d197227ddf17c6bd7937649af23f57371d",
            "excerpt": "github.com/gibbonmi/bench/internal/anchors,pass,937"
          },
          "requirement": "6-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-6-conformance-r1",
          "performer": "claude:bench-writer/ri-t6-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t6-author-20260927/6-conformance@c6b2a74e",
            "digest": "sha256:1fad3031dc2da3fcafd69f0df5d0adb2237368d60772997c3ec6c354ef6b8682",
            "excerpt": "github.com/gibbonmi/bench/internal/conformance,pass,43922"
          },
          "requirement": "6-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-worktree-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-worktree@c6b2a74e",
            "digest": "sha256:df17c9b605a27341a7b94b29354dcd35fa2bd39ba05dc2fb28244fa164e9f9b8",
            "excerpt": "github.com/gibbonmi/bench/internal/worktree,pass,52743"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-status-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-status@c6b2a74e",
            "digest": "sha256:8c89907f1a191473cd2b9f9ad6edefca805768ef41cc8727ebbc472a99bd11cb",
            "excerpt": "github.com/gibbonmi/bench/internal/status,pass,12750"
          },
          "requirement": "1-status",
          "command": "bench test --package ./internal/status",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-cmd-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-cmd@c6b2a74e",
            "digest": "sha256:59f6131664735a60e2038849a565fb2922d7ed0a58afd1e3e560d13be7606bd0",
            "excerpt": "github.com/gibbonmi/bench/cmd/bench,pass,8319"
          },
          "requirement": "1-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-conformance-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-conformance@c6b2a74e",
            "digest": "sha256:0d336fac346b0b7813079779d353daffd315eafee9b0af3a05b5d16e61e83f97",
            "excerpt": "github.com/gibbonmi/bench/internal/conformance,pass,36775"
          },
          "requirement": "1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-system-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-system@c6b2a74e",
            "digest": "sha256:ab56672c31a42be0c763cb5a3e5345a5cec9edfcfd6b6e17bb7cabafd68f2196",
            "excerpt": "github.com/gibbonmi/bench/internal/systemtest,pass,38997"
          },
          "requirement": "1-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-6-anchors-r2",
          "performer": "claude:bench-writer/ri-t6-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t6-author-20260927/6-anchors@4b80686a",
            "digest": "sha256:c84530d71fb19801d588927aa8ef66079c5a2daba98b66dd63331f2ccfccfae5",
            "excerpt": "github.com/gibbonmi/bench/internal/anchors,pass,944"
          },
          "requirement": "6-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-6-conformance-r2",
          "performer": "claude:bench-writer/ri-t6-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t6-author-20260927/6-conformance@4b80686a",
            "digest": "sha256:b9783166abd55b57c3bdaef654ac5c3a30d0497f90ee00187629a1c7c0558ab0",
            "excerpt": "github.com/gibbonmi/bench/internal/conformance,pass,45257"
          },
          "requirement": "6-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-worktree-r2",
          "performer": "claude:bench-writer/ri-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-1-20260927/1-worktree@4b80686a",
            "digest": "sha256:b6e44140a56bbd7085eab52aefe0d8a7d27716f1ff36d64f53b57ee678ce54c6",
            "excerpt": "github.com/gibbonmi/bench/internal/worktree,pass,54141"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-status-r2",
          "performer": "claude:bench-writer/ri-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-1-20260927/1-status@4b80686a",
            "digest": "sha256:8389f3cb6d0972958c2486d661f3e31150fd1ee3175e3b89c7de0e4d7d250b13",
            "excerpt": "github.com/gibbonmi/bench/internal/status,pass,13117"
          },
          "requirement": "1-status",
          "command": "bench test --package ./internal/status",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-cmd-r2",
          "performer": "claude:bench-writer/ri-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-1-20260927/1-cmd@4b80686a",
            "digest": "sha256:d6f2f1dd6410089b6bee4226d3b64239be8a9f3794e98634033f91fbe1dc40b2",
            "excerpt": "github.com/gibbonmi/bench/cmd/bench,pass,9261"
          },
          "requirement": "1-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-conformance-r2",
          "performer": "claude:bench-writer/ri-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-1-20260927/1-conformance@4b80686a",
            "digest": "sha256:b6fed8a9faa808070a69525526d8da361d5ae7ca429374468ff2d71296c51c1e",
            "excerpt": "github.com/gibbonmi/bench/internal/conformance,pass,35468"
          },
          "requirement": "1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-system-r2",
          "performer": "claude:bench-writer/ri-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-1-20260927/1-system@4b80686a",
            "digest": "sha256:5feb1076feb2d6497662939683e4eb3d53a85856a22eeeb4cfc8f40fcd192ae6",
            "excerpt": "github.com/gibbonmi/bench/internal/systemtest,pass,39178"
          },
          "requirement": "1-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-6-anchors-r3",
          "performer": "claude:bench-writer/ri-t6-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t6-author-20260927/6-anchors@09611663",
            "digest": "sha256:1f8bc1486d8676cbba75a6dd778eba2e2d8804a9621decf4661ce520ca987db1",
            "excerpt": "github.com/gibbonmi/bench/internal/anchors,pass,1028"
          },
          "requirement": "6-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-6-conformance-r3",
          "performer": "claude:bench-writer/ri-t6-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t6-author-20260927/6-conformance@09611663",
            "digest": "sha256:67296c993506fc96922a0b06c2ef35a5238fe0790edd039577ddb6916aaf99f4",
            "excerpt": "github.com/gibbonmi/bench/internal/conformance,pass,45688"
          },
          "requirement": "6-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-worktree-r3",
          "performer": "claude:bench-writer/ri-t1-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-2-20260927/1-worktree@09611663",
            "digest": "sha256:6c46822e35b5d9bbbda2fcfd347e3675bfb4922e9dd2602c824feed0641150fd",
            "excerpt": "github.com/gibbonmi/bench/internal/worktree,pass,54085"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-status-r3",
          "performer": "claude:bench-writer/ri-t1-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-2-20260927/1-status@09611663",
            "digest": "sha256:cc377cf3b46495a9ffd89ad57364fca21f2a0f454f6627cc118e85216b79dfe8",
            "excerpt": "github.com/gibbonmi/bench/internal/status,pass,13324"
          },
          "requirement": "1-status",
          "command": "bench test --package ./internal/status",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-cmd-r3",
          "performer": "claude:bench-writer/ri-t1-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-2-20260927/1-cmd@09611663",
            "digest": "sha256:4be54fbb06d7ef5b258e6d4f3153a9a3ac0de9dec4a6c4a6c90c40da9c514106",
            "excerpt": "github.com/gibbonmi/bench/cmd/bench,pass,8823"
          },
          "requirement": "1-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-conformance-r3",
          "performer": "claude:bench-writer/ri-t1-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-2-20260927/1-conformance@09611663",
            "digest": "sha256:68c5ed52e85e8ef11ed62b3117f78ef753e47a9181570d9de06a85872b11b68c",
            "excerpt": "github.com/gibbonmi/bench/internal/conformance,pass,35351"
          },
          "requirement": "1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-system-r3",
          "performer": "claude:bench-writer/ri-t1-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-repair-2-20260927/1-system@09611663",
            "digest": "sha256:aa8ec50038fcd6ce1c0a45133ccec9a797ebe641a8a8185926494ab07d8be78e",
            "excerpt": "github.com/gibbonmi/bench/internal/systemtest,pass,39299"
          },
          "requirement": "1-system",
          "command": "bench test --check system",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "ri-c1a-r1-standards",
          "performer": "claude:bench-reviewer/ri-c1a-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-standards@c6b2a74e",
            "digest": "sha256:d6f3966b440ea17dcb22ec20b97b50cb877761e05779ba5204d1264ed001e961",
            "excerpt": "Standards: 6 findings. Worst: the new tests spell class=, holder=, and retained: content main lacks independently of the production constants with no recorded red."
          },
          "axis": "Standards",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb",
          "finding_ids": ["R1", "R2", "R3", "R4", "R5", "R6"],
          "supersedes": []
        },
        {
          "id": "ri-c1a-r1-spec",
          "performer": "claude:bench-reviewer/ri-c1a-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-spec@c6b2a74e",
            "digest": "sha256:e41ea229c8be77d3e2137c89f35bab643cf354dfdc9efa96af7473a9dbf8bf87",
            "excerpt": "Spec: 2 findings. Worst: row RI5 reads one commit, but a one-commit squash lands by patch containment; the two-commit fixture follows story 5."
          },
          "axis": "Spec",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb",
          "finding_ids": ["R7", "R8"],
          "supersedes": []
        },
        {
          "id": "ri-c1a-r1-coverage",
          "performer": "claude:bench-reviewer/ri-c1a-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-coverage@c6b2a74e",
            "digest": "sha256:a09849151032402867990fe921a2563d7548a14599b378b9ca010743682b501f",
            "excerpt": "Coverage: 4 findings. Worst: a symbolic ref in a Bench namespace makes the bulk apply delete its target, observed on a unique root and on the checked-out main."
          },
          "axis": "Coverage",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb",
          "finding_ids": ["R9", "R10", "R11", "R12"],
          "supersedes": []
        },
        {
          "id": "ri-c1a-r2-standards",
          "performer": "claude:bench-reviewer/ri-c1a-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-standards-r2@4b80686a",
            "digest": "sha256:e6e9a0142212fe9c526af9e4166f2bcb6a9e60341c024bc3add50b7c32b52e6f",
            "excerpt": "Standards round 2: 4 findings, six folds confirmed. Worst: the RI-C1a chunk table cell omits RI83 to RI87 after the plan expansion."
          },
          "axis": "Standards",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "4b80686a12ea46afa0e6c0ba4fe9aab2e00935f6",
          "finding_ids": ["R13", "R14", "R15", "R16"],
          "supersedes": ["ri-c1a-r1-standards"]
        },
        {
          "id": "ri-c1a-r2-spec",
          "performer": "claude:bench-reviewer/ri-c1a-spec-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-spec-r2@4b80686a",
            "digest": "sha256:54109ca17603b30d4a18e9fb88fcbdf2ddddc1d0943f84348f7419032b4236a8",
            "excerpt": "Spec round 2: 3 findings, R7, R8, D1, D2 confirmed. Worst: the chunk table cell omits RI83 to RI87, and a dangling symref is skipped by for-each-ref."
          },
          "axis": "Spec",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "4b80686a12ea46afa0e6c0ba4fe9aab2e00935f6",
          "finding_ids": ["R13", "R17", "R18"],
          "supersedes": ["ri-c1a-r1-spec"]
        },
        {
          "id": "ri-c1a-r2-coverage",
          "performer": "claude:bench-reviewer/ri-c1a-coverage-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "a6b40aeaceecdcdce4aeef5b455b95b43cbd5449",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-coverage-r2@4b80686a",
            "digest": "sha256:82d71f9cc5717dd192fd495bf8e6519cf1b7ff4119be37dfef68630944648505",
            "excerpt": "Coverage round 2: 1 finding, R9 to R12 and D2 confirmed. Worst: a dangling symref in a Bench namespace is excluded silently, against the RI83 wording."
          },
          "axis": "Coverage",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "4b80686a12ea46afa0e6c0ba4fe9aab2e00935f6",
          "finding_ids": ["R17"],
          "supersedes": ["ri-c1a-r1-coverage"]
        },
        {
          "id": "ri-c1a-r3-standards",
          "performer": "claude:bench-reviewer/ri-c1a-standards-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-standards-r3@09611663",
            "digest": "sha256:e66920a0027f6cc0fdf894cfbe376e8270ea4e1484945d8bfc26559b382b97df",
            "excerpt": "Standards round 3: 0 blocking findings. R13, R14, R15, R16, and R18 confirmed; the repair 2 delta adds no duplicated knowledge."
          },
          "axis": "Standards",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "finding_ids": [],
          "supersedes": ["ri-c1a-r2-standards"]
        },
        {
          "id": "ri-c1a-r3-spec",
          "performer": "claude:bench-reviewer/ri-c1a-spec-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-spec-r3@09611663",
            "digest": "sha256:7b22718f680d7295995d5124cfb02b38204aa2cf9f033f43c34d68c47de84b58",
            "excerpt": "Spec round 3: 0 blocking findings. R13, R17 with D3, and R18 confirmed; every RI-C1a row re-verified at 09611663."
          },
          "axis": "Spec",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "finding_ids": [],
          "supersedes": ["ri-c1a-r2-spec"]
        },
        {
          "id": "ri-c1a-r3-coverage",
          "performer": "claude:bench-reviewer/ri-c1a-coverage-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8cb477fca58ec20c7e531218bc6ed542a9e7391a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-coverage-r3@09611663",
            "digest": "sha256:6f5607871dde38df253baad1cc70c50847ec67fca327955f29049751deb7e906",
            "excerpt": "Coverage round 3: 0 blocking findings. R15 confirmed by four call-site unwrap probes; the repair 2 delta changes no production behavior."
          },
          "axis": "Coverage",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "finding_ids": [],
          "supersedes": ["ri-c1a-r2-coverage"]
        }
      ]
    },
    {
      "id": "RI-C1b",
      "base": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
      "tip": "c6d2cfbf66d4b82f284e008ab063266bf61c4a23",
      "plan_digest": "sha256:74f8fb4685ed7e515b51757c2cd9d36d59a79ba8659977672e3b3aad337a118a",
      "source_digest": "039b389cfbca79ef22bfd8e3092c47270ee89388",
      "acceptance_rows": ["RI24", "RI26", "RI27", "RI88", "RI89", "RI90", "RI91", "RI92"],
      "verification": [
        {
          "id": "ri-c1b-2-status-r1",
          "performer": "claude:bench-writer/ri-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "05b8019186bf6d6b724705534ec913486ac49945",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t2-author-20260927/2-status@f1204694",
            "digest": "sha256:29f24d6fa9e391ef2be584c359f14dd8e618cb5b0cd28cd261d85bfaafcb5a7b",
            "excerpt": "github.com/gibbonmi/bench/internal/status,pass,13921"
          },
          "requirement": "2-status",
          "command": "bench test --package ./internal/status",
          "exit_code": 0
        },
        {
          "id": "ri-c1b-2-worktree-r1",
          "performer": "claude:bench-writer/ri-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "05b8019186bf6d6b724705534ec913486ac49945",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t2-author-20260927/2-worktree@f1204694",
            "digest": "sha256:87c4a1aa0a6619187f2ee5ab8d2ff295a9b18090e3395688cc006f85f930d577",
            "excerpt": "github.com/gibbonmi/bench/internal/worktree,pass,47162"
          },
          "requirement": "2-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "ri-c1b-2-system-r1",
          "performer": "claude:bench-writer/ri-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "05b8019186bf6d6b724705534ec913486ac49945",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t2-author-20260927/2-system@f1204694",
            "digest": "sha256:ad44035df9e790e200d8d60bb48bdc142bd582426679189d7fc6677539475004",
            "excerpt": "github.com/gibbonmi/bench/internal/systemtest,pass,39666"
          },
          "requirement": "2-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "ri-c1b-2-status-r2",
          "performer": "claude:bench-writer/ri-t2-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b389cfbca79ef22bfd8e3092c47270ee89388",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t2-repair-2-20260927/2-status@c6d2cfbf",
            "digest": "sha256:0b1a35a60e68fe999d5a95cc44aa61b9154c8acd60e471846aba7afc74508ae6",
            "excerpt": "github.com/gibbonmi/bench/internal/status,pass,14673"
          },
          "requirement": "2-status",
          "command": "bench test --package ./internal/status",
          "exit_code": 0
        },
        {
          "id": "ri-c1b-2-worktree-r2",
          "performer": "claude:bench-writer/ri-t2-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b389cfbca79ef22bfd8e3092c47270ee89388",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t2-repair-2-20260927/2-worktree@c6d2cfbf",
            "digest": "sha256:4474e4bc8de376b0b5ce7cf2fd48ed1ede575a95d8af48b892440eaff2affa08",
            "excerpt": "github.com/gibbonmi/bench/internal/worktree,pass,52407"
          },
          "requirement": "2-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "ri-c1b-2-system-r2",
          "performer": "claude:bench-writer/ri-t2-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b389cfbca79ef22bfd8e3092c47270ee89388",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t2-repair-2-20260927/2-system@c6d2cfbf",
            "digest": "sha256:3aece8da14e66f5bf31be8a95184f88186bb96c1b8fcdfcd0bc6b7f5ef37ae59",
            "excerpt": "github.com/gibbonmi/bench/internal/systemtest,pass,44100"
          },
          "requirement": "2-system",
          "command": "bench test --check system",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "ri-c1b-r1-standards",
          "performer": "claude:bench-reviewer/ri-c1b-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "05b8019186bf6d6b724705534ec913486ac49945",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1b-standards@f1204694",
            "digest": "sha256:5f525a693e2763eb1fb542791e79f79dee68d61ba1ac32a25588a95f97501d3f",
            "excerpt": "Standards: 5 findings. Worst: the plan command is spelled in the status action table and built by unclaimedReplan, two sources for one fact."
          },
          "axis": "Standards",
          "base": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "tip": "f12046941a0ad1f2812e82de0bf4e9c61303079b",
          "finding_ids": ["R19", "R20", "R21", "R22", "R23"],
          "supersedes": []
        },
        {
          "id": "ri-c1b-r1-spec",
          "performer": "claude:bench-reviewer/ri-c1b-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "05b8019186bf6d6b724705534ec913486ac49945",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1b-spec@f1204694",
            "digest": "sha256:bb100f40da676efbbef661ae83ecdd5054db9f48d24b3e05c9a23dc121e4af36",
            "excerpt": "Spec: 2 findings. Worst: a faulted symref whose target is unique is counted as a faulted ref and again as a unique branch."
          },
          "axis": "Spec",
          "base": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "tip": "f12046941a0ad1f2812e82de0bf4e9c61303079b",
          "finding_ids": ["R24", "R25"],
          "supersedes": []
        },
        {
          "id": "ri-c1b-r1-coverage",
          "performer": "claude:bench-reviewer/ri-c1b-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "05b8019186bf6d6b724705534ec913486ac49945",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1b-coverage@f1204694",
            "digest": "sha256:1d3bcdac79cdc43c1a29c3dce82f5fda38c851a65e732be73c71a984ad4b0048",
            "excerpt": "Coverage: 3 findings. Worst: a blob-tip Bench ref makes LandedState fail first, so the row reads git state unavailable and never routes to the plan."
          },
          "axis": "Coverage",
          "base": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "tip": "f12046941a0ad1f2812e82de0bf4e9c61303079b",
          "finding_ids": ["R24", "R25", "R26"],
          "supersedes": []
        },
        {
          "id": "ri-c1b-r2-standards",
          "performer": "claude:bench-reviewer/ri-c1b-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b389cfbca79ef22bfd8e3092c47270ee89388",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-c1b-standards-r2@c6d2cfbf",
            "digest": "sha256:f56429cd44782253c4a633cb017d015b6b0f2284544bb5f87af530418b769e80",
            "excerpt": "Standards round 2: 0 blocking findings. R19 to R23, D4, D5, and D6 confirmed; the repair deltas add no duplicated knowledge."
          },
          "axis": "Standards",
          "base": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "tip": "c6d2cfbf66d4b82f284e008ab063266bf61c4a23",
          "finding_ids": [],
          "supersedes": ["ri-c1b-r1-standards"]
        },
        {
          "id": "ri-c1b-r2-spec",
          "performer": "claude:bench-reviewer/ri-c1b-spec-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b389cfbca79ef22bfd8e3092c47270ee89388",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-c1b-spec-r2@c6d2cfbf",
            "digest": "sha256:cddde126f0db26bf6e5f13172a423e469eb62b8b691d792e9993f58955c3a2cc",
            "excerpt": "Spec round 2: 0 blocking findings. R24, R25, R26, D4, D5, and D6 confirmed; all eleven RI-C1b rows delivered at c6d2cfbf."
          },
          "axis": "Spec",
          "base": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "tip": "c6d2cfbf66d4b82f284e008ab063266bf61c4a23",
          "finding_ids": [],
          "supersedes": ["ri-c1b-r1-spec"]
        },
        {
          "id": "ri-c1b-r2-coverage",
          "performer": "claude:bench-reviewer/ri-c1b-coverage-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "039b389cfbca79ef22bfd8e3092c47270ee89388",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-c1b-coverage-r2@c6d2cfbf",
            "digest": "sha256:36a8cfad50e20decad31f0a8cf39a2cafa8f30f90acab45adbd6d28ff3d6a437",
            "excerpt": "Coverage round 2: 0 blocking findings. R24 to R26 and D6 confirmed by fixture and producer probes; the repair deltas violate no row."
          },
          "axis": "Coverage",
          "base": "09611663f8a49bd5f037b24ba9385eaa2e42be41",
          "tip": "c6d2cfbf66d4b82f284e008ab063266bf61c4a23",
          "finding_ids": [],
          "supersedes": ["ri-c1b-r1-coverage"]
        }
      ]
    }
  ],
  "completion": {
    "state": "pending"
  }
}
```
