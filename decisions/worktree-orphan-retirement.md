# Worktree orphan retirement (FT148)

## Destination

A worktree cut by a session that later dies must stop being immortal. Today
nothing retires it: `bench worktree release` matches only the exact plaintext
request string that created the assignment, the ledger stores a one-way digest
of that string, and the harness hook derives it from the session id — so once
the creating session is gone its worktrees are structurally unreleasable. The
pool accreted from 2026-07-09 to 2026-07-27, every entry was re-preserved at
every resume sweep, and draining it by hand took a staged script and a full
session.

## Provenance

This map was written in the same session as the spec it compiles — the
highest-bias path in this workflow. Read the `Type:` line on each decision
before trusting it:

- **Closed by reviewer, 2026-07-27 (roadmap row)** — signed off before this
  session, recorded in `ROADMAP.md`'s FT148 row.
- **Closed by reviewer, 2026-07-27 (spec-authoring session)** — put to the
  reviewer during this session and answered.
- **Decided by the author** — no sign-off. Flagged in the spec for veto.

A mid-tier falsification pass on the first draft found that the original #2
(a lease conjunct) was unimplementable, that #5 carried a sign-off the roadmap
row does not give, and that several assertables were unobservable. Those
findings were verified against the tree and are folded in below.

## #1: Which command retires an orphan?

Type: Closed by reviewer, 2026-07-27 (roadmap row)

### Question
Give `release` a request-derivation override so a fresh session can name a dead
session's assignment, or route orphans to `bench worktree clean`?

### Answer
`bench worktree clean`. A request-derivation override is rejected: the request
digest *is* the ownership proof, and deriving it on demand voids the ownership
model the whole lifecycle rests on.

Verified in code: `bench worktree clean` runs `PlanExplicit`
(`internal/worktree/subshell.go`), which never reads `assignment.State`, and
retains on a live lease but not a dead one. It can retire an orphan today with
no change to the cleanup path itself.

Two caveats the first draft overstated. It is a **two-step** command: the bare
form prints a plan and a fingerprint, and removal needs
`--apply <fingerprint>`. And an orphan carrying ignored build output — the
normal state of a worktree a shift ran in — retains under `ignored` and needs
`--discard-ignored`, whose request-less form the kit's own comment says orphans
the assignment (`internal/worktree/ownership.go`, FT93b). So the route out
exists but is not one paste, and the surfaced command must not steer anyone
into the `--discard-ignored` trap.

## #2: What makes an assignment orphaned?

Type: Closed by reviewer, 2026-07-27 (spec-authoring session), replacing a
rejected first answer

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
  it returns for an unreadable or malformed one, so a predicate taking only a
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

Safety consequently rests on three things, none of them liveness: the window is
long (7 days, not 24 hours, chosen because it now carries all the weight); the
sweep only ever **reports** (#4); and the explicit cleanup an operator then runs
recovers dirty work into a recovery ref before removing anything.

## #3: How does an unstamped record age?

Type: Closed by reviewer, 2026-07-27 (spec-authoring session)

### Question
The 17 assignment records live in the ledger today predate any `created_at`
field. Absent = never orphaned (fail closed), absent = aged, or backfill the
stamp on first read?

### Answer
Absent = aged. A record with no `created_at` was written before the field
existed, so it is by construction older than any window this repo can set.
Fail-closed would leave today's residue immortal — the exact failure this row
exists to end — and backfilling silently mutates the ledger while delaying the
drain by a full window.

## #4: Does the resume sweep clean orphans, or only report them?

Type: Closed by reviewer, 2026-07-27 (roadmap row)

### Question
Auto-remove on the sweep, or surface a command?

### Answer
Report only. The sweep runs unattended at every session start; auto-removing a
tree that may hold uncommitted work is not a verdict a sweep gets to make alone.
With the lease conjunct gone (#2) this is no longer a preference — it is the
only thing standing between a long-lived legitimate worktree and a destructive
command running against it unattended.

The sweep reports; the operator runs the explicit two-step cleanup.

## #5: What happens to a ledger row whose tree is already gone?

Type: Decided by the author, then put to the reviewer and **approved
2026-07-27** — posed but not decided by the roadmap row.

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

The 16 `recovered` rows are not this spec's to drain: their recovery refs are
retained because FT98's landed proof misses reshaped commits, so payloads main
actually shipped still read as unlanded. That half rides FT98.

## #6: The preserved wall stays after this build — what does the reviewer see?

Type: Closed by reviewer, 2026-07-27 (spec-authoring session). Scope addition
beyond the roadmap row's split; flagged in the spec for veto.

### Question
This build compacts one row. The other 16 keep printing one `preserved` line
each at every session start until FT98 lands. Leave the wall, or bound it?

### Answer
Bound it. At most three orphan lines and three preserved lines print, each group
followed by an `and <n> more` line naming the true total when the cap bites.
Nothing is hidden — the count is stated and every record stays listable through
`bench worktree list`.

## #7: Where does the prose half land, and how does the gate see it?

Type: Closed by reviewer, 2026-07-27 (roadmap row) for content; seam decided by
the author

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
repository; `go test -run` with no match exits 0, so a builder following that
row would have seen green and concluded the row was satisfied.

The reverse sweep `checkColdPickupCLILists` already forces `.bench/BENCH.md` or
`BENCH-reference.md` to name every top-level `bench <cmd>` route, but it matches
top-level commands only, so `worktree release` / `clean` / `recovery` are
unenforced today.

## Handoff

1. **Module boundaries.** `internal/intent` owns the `Assignment.CreatedAt`
   field and its validation. `internal/worktree` owns the orphan predicate
   (`classifier.go`), the `PlanAutomatic` reason label, the resume sweep and its
   summary rendering (`resume.go`, `worktree.go`), and the create-time stamp
   (`ownership.go`). `internal/bounds` owns the window constant. The prose half
   touches `.agents/skills/bench-craft-delegate/SKILL.md`,
   `.agents/commands/bench-implement-spec.md`, `.bench/BENCH.md`, and the
   conformance anchor registry. `internal/gate` and the release path are
   outside — untouched.

2. **Contracts.** `Assignment` gains `created_at` as an optional RFC3339 string,
   omitempty, under the unchanged `bench-assignment/v1` schema — old records stay
   valid, so there is no ledger migration. `ValidateAssignment` rejects a
   present-but-unparseable value and accepts absence. `CreatedAt` must **not**
   enter `lockReason`, whose exact string both `validateCreationBundle` and
   `PlanExplicit` re-derive and compare for every pre-existing locked worktree.
   The orphan predicate is `orphaned(assignment, now) bool` — no lease argument,
   because no lease exists on this path. A new `ReasonOrphaned` joins the
   `CleanupReason` set and the fixed reason order in `renderResumeSummary`.
   `ResumeResult` gains an orphan list beside `Preserved`. No CLI surface, flag,
   or exit code changes.

3. **Deep vs thin.** The predicate is a pure function of `(assignment, now)`,
   following `reclaimable`'s existing precedent of taking `now` as a parameter
   rather than reading the clock, so tests drive age by argument with no clock
   injection machinery. The resume sweep computes orphanhood **directly from the
   ledger record**, not from a plan's reason code — `PlanAutomatic` returns early
   on any retain verdict before it reaches the state branch, so a plan-derived
   reading would silently miss every orphan that also carries ignored build
   output. `PlanAutomatic`'s `ReasonOrphaned` is a labelling improvement layered
   on top, not the source of truth. `bounds.AssignmentStale` is the single source
   of the window.

4. **Black-box assertables.** Through the `internal/worktree` package tests
   against a real temp repo, which is how `worktree_test.go` and
   `ownership_test.go` already work: an assignment younger than the window is not
   orphaned; one older is; one with no stamp is orphaned regardless of age; one
   with a future stamp is not; a record in any non-`active` state is never
   orphaned. The sweep reports an orphan **with ignored residue** — the case
   `PlanAutomatic`'s early return hides. The summary contains a literal
   `bench worktree clean <path>` line per orphan, never contains
   `--discard-ignored`, quotes a path containing a space or glob character, and
   does not let a control byte in a path split its line structure. The summary
   caps at three per group and states the true total. A tree-gone orphaned record
   with no recovery metadata is deleted and counted; the same record with
   recovery metadata is preserved. Two consecutive sweeps produce the same
   summary and delete nothing twice. `bench worktree release` still succeeds on a
   worktree whose lock reason predates the field. The prose duties assert through
   the conformance anchor registry.

5. **Gate attachment.** The `internal/worktree` and `internal/intent` package
   tests run in the gate's conformance phase; the anchor registry runs in the
   same phase as the `docs-currency-workflow` check. Both halves are inside
   `bench gate` — no manual evidence step.

6. **Hostile-input owners.** `ValidateAssignment` owns the malformed
   `created_at`: empty string, non-RFC3339 text, and control bytes are rejected.
   A **future** timestamp is *accepted* by validation and handled by the
   predicate, which treats it as not-aged — the two owners are disjoint, and the
   first draft wrongly assigned it to both. Note the blast radius: validation runs
   on every ledger read, so a rejection turns one malformed record into a total
   ledger outage for every `bench` command. That matches the package's existing
   fail-closed posture and is deliberate, but it is a reviewer-visible call.
   `renderResumeSummary` owns the path it prints: it is a raw line sink with no
   safety predicate today, and this build puts the first attacker-influenced path
   into it, so it owns both shell quoting and line-structure integrity — the
   profile's checklist calls out asserting the *permitted* control bytes, not
   only the refused ones. `sweepOrphanAssignments` owns the
   tree-gone-but-still-registered case, which it already defers to the prune path.

7. **Uncertainty flags.** Three calls originated with the author rather than the
   roadmap row and were approved at spec sign-off on 2026-07-27: the 7-day window
   value (#2), the ledger compaction behavior (#5), and the summary cap (#6). The
   window remains the one to watch during the build — with the lease conjunct
   gone it is the only thing separating a live long-running worktree from an
   orphan verdict, so a story that quietly shortens it is a spec deviation.

8. **Rejected alternatives.** A request-derivation override for `release` (#1 —
   voids the ownership model). A lease-liveness conjunct (#2 — no lease exists on
   this path, and `ProbeLease` cannot express absence anyway). A heartbeat write
   from still-using sessions (#2 — real liveness, but materially larger than the
   signed-off shape). A 24-hour window (#2 — too short once age carries all the
   weight alone). Fail-closed or backfilled unstamped records (#3). Auto-removal
   on the sweep (#4). A ledger schema bump to v3 (item 2 — the field is additive,
   so no migration is needed). Deriving the sweep's orphan verdict from
   `PlanAutomatic`'s reason code (item 3 — the early return hides orphans with
   ignored residue).

9. **Domain watch-outs.** `PlanAutomatic` returns at its first retain verdict,
   *before* the assignment-state branch, so any orphan carrying ignored residue,
   an unexpected lock, or dirty nested state never reaches the reason label —
   build the sweep off the ledger, not off the plan.
   `TestResumeSweepsResidueAndReportsPreserved` (`worktree_test.go`) pins an
   FT93(c) contract that an active, tree-gone, unregistered record **must survive
   the sweep**; after this build it survives for a different reason — the record
   is freshly stamped and therefore not aged — so the test keeps passing while
   the behavior it pins has changed underneath it. Say so, or it reads as
   untouched. `intent.LifecycleEvidence` serializes the whole assignments array
   into both the recovery-retire fingerprint *and* the explicit cleanup
   fingerprint, so adding a field restales every pending plan and every in-flight
   receipt checkpoint comparison; that is the designed plan/apply behavior, but
   story 5's whole workflow is plan/apply, so expect it.
   `automaticFingerprint` composes from named fields rather than the struct and
   is unaffected — keep it that way. `jsonfile.Decode` *requires* a final
   newline rather than tolerating its absence. The command is
   `bench resume-clean`; there is no `bench resume` route, though
   `renderResumeSummary`'s own output prefix says "bench resume:" — a
   pre-existing naming defect this build should not copy into new prose.

Dependency order: single spec. FT98 is a sibling, not a prerequisite.
