# Worktree orphan retirement (FT148)

Status: ready

## Destination

A worktree cut by a session that later dies must stop being immortal. Today
nothing retires it: `bench worktree release` matches only the exact plaintext
request string that created the assignment, the ledger stores a one-way digest
of that string, and the harness hook derives it from the session id. So once
the creating session is gone its worktrees are structurally unreleasable. The
pool accreted from 2026-07-09 to 2026-07-27, and every entry was re-preserved
at every resume sweep. Draining it by hand took a staged script and a full
session.

## Provenance

This map was written in the same session as the spec it compiles — the
highest-bias path in this workflow. Every ticket uses the canonical `Grill`
type; its decision-specific provenance is:

- **#1 and #4:** Closed by reviewer, 2026-07-27 (roadmap row), signed off
  before this session and recorded in `ROADMAP.md`'s FT148 row.
- **#2:** Closed by reviewer, 2026-07-27 (spec-authoring session), replacing
  a rejected first answer.
- **#3:** Closed by reviewer, 2026-07-27 (spec-authoring session).
- **#5:** Decided by the author, then put to the reviewer and **approved
  2026-07-27**; the roadmap row posed it but did not decide it. The spec flags
  this state-destroying behavior for veto.
- **#6:** Closed by reviewer, 2026-07-27 (spec-authoring session). This scope
  addition goes beyond the roadmap row's split and is flagged in the spec for
  veto.
- **#7:** Closed by reviewer, 2026-07-27 (roadmap row) for content; its seam
  was decided by the author.

A mid-tier falsification pass on the first draft found three faults. The
original #2 (a lease conjunct) was unimplementable, #5 carried a sign-off the
roadmap row does not give, and several assertables were unobservable. Those
findings were verified against the tree and are folded in below.

## #1: Which command retires an orphan?

Blocked by: none
Type: Grill

### Question
Give `release` a request-derivation override so a fresh session can name a dead
session's assignment, or route orphans to `bench worktree clean`?

### Answer
`bench worktree clean`. A request-derivation override is rejected: the request
digest *is* the ownership proof. Deriving it on demand voids the ownership
model the whole lifecycle rests on.

Verified in code: `bench worktree clean` runs `PlanExplicit`
(`internal/worktree/subshell.go`), which never reads `assignment.State`, and
retains on a live lease but not a dead one. It can retire an orphan today with
no change to the cleanup path itself.

Two caveats the first draft overstated. It is a **two-step** command: the bare
form prints a plan and a fingerprint, and removal needs
`--apply <fingerprint>`. An orphan carrying ignored build output — the
normal state of a worktree a shift ran in — retains under `ignored` and needs
`--discard-ignored`. Its request-less form the kit's own comment says orphans
the assignment (`internal/worktree/ownership.go`, FT93b). So the route out
exists but is not one paste, and the surfaced command must not steer anyone
into the `--discard-ignored` trap.

## #2: What makes an assignment orphaned?

Blocked by: none
Type: Grill

### Question
Assignments carry no created-at timestamp, no lease, and no reaper, and
`PlanAutomatic` hard-retains every `active` record. What evidence promotes a
record to `orphaned`?

### Answer
**Age alone**: the assignment is in state `active`, and it is older than
`bounds.AssignmentStale`, a fixed **7-day** constant beside `LeaseStale`. A
record with no creation stamp counts as aged (see #3). A record in any other
state is already on a retirement path and is never orphaned.

The first draft added a lease-liveness conjunct and argued that the conjunction
was what made the verdict safe. That is false, on two independent grounds found
by the falsification pass and verified:

- `ProbeLease` returns `LeaseUnknown` for an **absent** lease, the same verdict
  it returns for an unreadable or malformed one. A predicate taking only a
  three-valued `LeaseState` cannot separate "no liveness claim" from "cannot
  tell". `listLease` already works around this by stat-ing outside `ProbeLease`
  and returning a fourth value.
- More fundamentally, **no lease exists on this path at all.** `Create` writes
  the owner marker, the ledger record, and the git lock — never a lease. Leases
  are written only by `Claim` via `Acquire`, the `bench shift` pool path. And a
  lease records a *pid*: `bench worktree create` runs as a short-lived hook
  process that exits immediately, so a lease it wrote would read dead within
  milliseconds. A request-created worktree outliving its creating process is the
  design, not a fault.

There is therefore no liveness signal available for request-created assignments
and no cheap one to add. A heartbeat that a still-using session touches was
considered and rejected as materially larger than the signed-off shape.

Safety consequently rests on three things, none of them liveness. The window is
long (7 days, not 24 hours), chosen because it now carries all the weight.
The sweep only ever **reports** (#4). The explicit cleanup an operator then
runs recovers dirty work into a recovery ref before removing anything.

## #3: How does an unstamped record age?

Blocked by: none
Type: Grill

### Question
The 17 assignment records live in the ledger today predate any `created_at`
field. Absent = never orphaned (fail closed), absent = aged, or backfill the
stamp on first read?

### Answer
Absent = aged. A record with no `created_at` was written before the field
existed, so it is by construction older than any window this repo can set.
Fail-closed would leave today's residue immortal — the exact failure this row
exists to end. Backfilling silently mutates the ledger, while delaying the
drain by a full window.

## #4: Does the resume sweep clean orphans, or only report them?

Blocked by: none
Type: Grill

### Question
Auto-remove on the sweep, or surface a command?

### Answer
Report only. The sweep runs unattended at every session start. Auto-removing a
tree that may hold uncommitted work is not a verdict a sweep gets to make
alone. With the lease conjunct gone (#2) this is no longer a preference. It is the
only thing standing between a long-lived legitimate worktree and a
destructive command running against it unattended.

The sweep reports; the operator runs the explicit two-step cleanup.

## #5: What happens to a ledger row whose tree is already gone?

Blocked by: none
Type: Grill

### Question
After the manual drain, 17 records survive with no tree on disk — 16
`recovered` plus one stuck `active`. Nothing compacts an active record whose
tree is missing, because `sweepOrphanAssignments` skips `active` on the
principle that a live session owns it. What compacts them?

### Answer
Extend the same sweep: an `active` record that is orphaned by #2, has no tree on
disk, is not a registered worktree, and holds no recovery metadata is deleted.
An orphaned record that *does* hold recovery metadata is preserved and reported,
exactly like the `recovered` ones.

`ROADMAP.md` frames this as "the second thing the row must answer" — a question,
not a closed decision. The signed-off (a)/(b)/(c) split does not contain it. It
is the only state-destroying behavior in the build and it carries no sign-off,
so the spec flags it.

The 16 `recovered` rows are not this spec's to drain. Their recovery refs are
retained because FT98's landed proof misses reshaped commits, so payloads main
actually shipped still read as unlanded. That half rides FT98.

## #6: The preserved wall stays after this build — what does the reviewer see?

Blocked by: none
Type: Grill

### Question
This build compacts one row. The other 16 keep printing one `preserved` line
each at every session start until FT98 lands. Leave the wall, or bound it?

### Answer
Bound it. At most three orphan lines and three preserved lines print, each group
followed by an `and <n> more` line naming the true total when the cap bites.
Nothing is hidden — the count is stated and every record stays listable through
`bench worktree list`.

## #7: Where does the prose half land, and how does the gate see it?

Blocked by: none
Type: Grill

### Question
Kit prose orders worktree *creation* many times for every retirement
instruction, and `bench worktree release` is named in no guidance file at all.
Three prose edits are signed off. Prose is not gate-observable, so how does the
gate hold them?

### Answer
The `require(<file>, <phrase>)` registry in `checkWorkflowAnchors`
(`internal/conformance/docs_workflow_helpers_test.go`), which already pins
dozens of kit-prose obligations by exact phrase, including one on
`craft-delegate`. It is bound into the conformance sweep as the
`docs-currency-workflow` registry check.

The observable red command is a scoped conformance run:

    BENCH_CONFORMANCE_ROOT=$PWD BENCH_CONFORMANCE_CHECK=docs-currency-workflow \
      go test ./internal/conformance -run TestRootConformance

The first draft named `TestDocsWorkflowContracts`, which exists nowhere in the
repository. `go test -run` with no match exits 0, so a builder following that
row would have seen green and concluded the row was satisfied.

The reverse sweep `checkColdPickupCLILists` already forces `.bench/BENCH.md` or
`BENCH-reference.md` to name every top-level `bench <cmd>` route. It matches
top-level commands only, so `worktree release` / `clean` / `recovery` are
unenforced today.

## Not yet specified

## Spec-writer discretion

## Out of scope

- Draining recovered rows and their recovery references; FT98 owns that retained payload work.
- Any gate, release-path, or worktree behavior outside FT148's orphan retirement and reported cleanup scope.

## Sources
