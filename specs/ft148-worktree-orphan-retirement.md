# FT148 — worktree orphan retirement

Status: staged

Compiled from `decisions/worktree-orphan-retirement.md`, written in the same
session as this draft from the fix shape signed off on 2026-07-27 in
`ROADMAP.md`'s FT148 row plus a read of the tree. Read that map's **Provenance**
section: each decision is labelled with whether the reviewer closed it or the
author decided it.

**Three calls were made without prior sign-off and put to the reviewer at
approval:** the 7-day window (story 3), the summary cap (story 6), and ledger
compaction (story 7). All three were approved on 2026-07-27; each story records
that below. They are noted rather than removed because they are the parts of this
spec the roadmap row did not decide.

A mid-tier falsification pass on the first draft blocked it. Its findings were
verified against the tree and are folded in; the largest killed the original
safety mechanism outright and is recorded in the map's #2.

## Problem

A worktree cut by a session that later dies is immortal. `bench worktree release`
matches only the exact plaintext request string that created the assignment, the
ledger stores a one-way digest of it, and the harness hook derives that string
from the session id — so once the creating session is gone, its worktrees are
structurally unreleasable rather than accidentally so. Reproduced through the
accused command on 2026-07-27: one `--request` refused all 19 live pool worktrees
identically with "request, assignment, or path mismatch; checkout retained",
including the 12 the tool itself reported `landed=true`.

Three things kept it invisible until the pool had accreted for eighteen days.
Kit prose orders worktree *creation* many times for every retirement
instruction, and `bench worktree release` is named in no guidance file at all —
only in the CLI inventory. Assignments carry no created-at timestamp and no
reaper, and the sweep hard-retains every `active` record. And the sweep
re-preserves every surviving record at every session start, so the reviewer read
a twenty-line "preserved" wall each time with no action attached to it. Draining
the pool by hand took a staged script and a full session.

## Solution

Give an abandoned assignment a name, a verdict, and an attached command.

An assignment records when it was created. An `active` record older than a fixed
window is **orphaned**. `bench resume-clean` reports each orphan with the
`bench worktree clean <path>` that starts its retirement, and compacts the ledger
rows whose trees are already gone and which preserve no work. The summary is
bounded, so the session-start wall stays readable no matter how much residue is
pending.

Nothing is removed unattended: the sweep only reports, and the operator runs the
explicit plan/apply cleanup, which recovers dirty work into a recovery ref before
it removes anything.

Alongside the code, the prose that never mentioned retirement gets it, and each
new duty is pinned by the gate so it cannot be quietly deleted.

**What this does not have.** There is no liveness signal. `Create` writes no
lease, and a lease records a pid that dies the moment the create hook exits — a
request-created worktree outliving its creating process is the design. So
orphanhood is age alone, and safety rests on the window being long, the sweep
only reporting, and cleanup recovering before it removes. The map's #2 records
why the obvious alternatives fail.

## User stories

1. As the intent ledger, I want an assignment to record when it was created, so
   that anything downstream can ask how old it is without guessing from the
   filesystem. The field is optional and lives under the unchanged
   `bench-assignment/v1` schema, so records written before it stay valid and no
   migration runs. Validation rejects a present-but-unparseable value and accepts
   absence.
   Line: `gpt-5.6-terra` / medium (`opus` in Claude Code). This is a
   persisted-schema change whose validation runs on every ledger read, so getting
   the accept-absence half wrong makes every existing ledger unreadable.

2. As `bench worktree create`, I want to stamp each new assignment with its
   creation time, so that every record written from now on can age. The stamp
   must not enter the git worktree lock reason, whose exact string is re-derived
   and compared by both the release path and the explicit cleanup planner.
   Line: `gpt-5.6-terra` / low (`opus`). The write itself is one field at an
   identified construction site, but the lock-reason coupling is the kind a cheap
   pass reasonably misses.

3. As the cleanup classifier, I want one predicate that decides whether an
   assignment is orphaned, so that every consumer asks the same question. An
   assignment is orphaned when its state is `active` and it is older than
   `bounds.AssignmentStale`. A record with no creation stamp counts as aged,
   because it was written before the field existed. A record with a future stamp
   is not aged. The predicate is a pure function of assignment and current time,
   following `reclaimable`'s existing precedent of taking `now` as a parameter.
   Line: `gpt-5.6-terra` / high (`opus`). This is the correctness core, and with
   no liveness conjunct the window is the only thing separating a live
   long-running worktree from an orphan verdict.
   **Author's call, approved 2026-07-27:** `bounds.AssignmentStale` is **7 days**. It is not
   configurable — the repo has no configuration surface for lifecycle tunables
   and adding one is a separate decision. Seven rather than one because age now
   carries all the safety weight alone.

4. As a reader of a cleanup plan, I want an orphaned assignment to retain under
   its own reason code rather than the generic `active` one, so that the plan
   says why it declined. This is a label only: `PlanAutomatic` returns at its
   first retain verdict, before it reaches the state branch, so an orphan
   carrying ignored residue or an unexpected lock keeps that more informative
   reason. The resume sweep never reads this label.
   Line: `gpt-5.6-luna` / low (`sonnet`). A reason-code swap on one branch, fully
   observed by the gate, at a seam this spec has already fixed.

5. As a session starting cold, I want `bench resume-clean` to print the command
   that starts each orphan's retirement, so that the residue carries its own next
   action instead of being a wall I have to decode. Each orphan gets one
   `bench worktree clean <path>` line, with the path quoted so a pool path
   containing spaces or glob characters stays a single runnable argument. The
   line is honestly the *first* of two steps — the bare form prints a plan and a
   fingerprint, and removal needs `--apply <fingerprint>`. The emitted line never
   suggests `--discard-ignored`, whose request-less form the codebase documents as
   orphaning the assignment.
   Line: `gpt-5.6-terra` / medium (`opus`). It is rendering, but it emits a
   command a reader will paste, and it is the surface the whole feature is judged
   on.

6. As a session starting cold, I want the resume summary bounded, so that a large
   backlog does not bury the rest of the report. At most three orphan lines and
   three preserved lines print, each group followed by an `and <n> more` line
   naming the true total when the cap bites. Nothing is hidden — the count is
   stated and every record stays listable through `bench worktree list`.
   Line: `gpt-5.6-luna` / low (`sonnet`). Pure rendering behind an assertion the
   gate fully observes.
   **Author's call, approved 2026-07-27:** a scope addition beyond the signed-off split, taken
   because otherwise the build ships and the wall the row was opened about is
   still there until FT98 lands.

7. As the resume sweep, I want an orphaned assignment whose tree is already gone
   and which preserves no work to be deleted from the ledger, so that rows stop
   outliving their trees. The sweep still skips non-orphaned `active` records — a
   live session may own those. An orphaned record that does hold recovery
   metadata is preserved and reported, never deleted.
   Line: `gpt-5.6-terra` / high (`opus`). This is the only story that destroys
   state, and the preserve-if-recovery branch is what stands between it and
   losing work.
   **Author's call, approved 2026-07-27.** `ROADMAP.md` poses ledger compaction
   as "the second thing the row must answer", a question rather than a decision,
   and the signed-off (a)/(b)/(c) split does not contain it. It was proposed here
   because it is the row's own stated question and splitting it would leave the
   stuck row un-drained for another cycle; the reviewer kept it in this spec.

8. As a coordinator accepting a delegate's done-claim, I want `craft-delegate` to
   name releasing the worktree I cut as part of close-out, so that retirement is
   ordered as often as creation is.
   Line: `gpt-5.6-sol` / high (`fable`). Kit guidance prose compounds through
   every session that loads it, which is the leverage override the project
   profile's cached routing already binds to the top tier; this is that cached
   routing, not a bump off it.

9. As a build that stops short, I want `bench-implement-spec`'s stop-short
   section to name who owns retiring the worktree, so that an abandoned build
   leaves an owner rather than an orphan.
   Line: `gpt-5.6-sol` / high (`fable`). Same cached routing — always-loaded
   phase guidance, and the failure it prevents is the one that produced this row.

10. As a session reading the CLI inventory cold, I want `.bench/BENCH.md` to name
    the worktree retirement subcommands, so that the route out of the pool is
    discoverable where the route in already is. The existing cold-pickup sweep
    checks top-level `bench <cmd>` routes only, so these subcommands are
    undocumented and unenforced today.
    Line: `gpt-5.6-sol` / high (`fable`). A short edit, but it lands in the
    always-loaded operating guide the cached routing binds to the top tier.

11. As the gate, I want each new prose duty pinned by exact phrase, so that
    deleting or paraphrasing it turns the build red instead of passing silently.
    The `require(<file>, <phrase>)` registry in `checkWorkflowAnchors`
    (`internal/conformance/docs_workflow_helpers_test.go`) already pins dozens of
    kit-prose obligations, including one on `craft-delegate`; each duty from
    stories 8, 9, and 10 gets a row there.
    Line: `gpt-5.6-terra` / medium (`opus`). The project profile routes gate and
    conformance logic to mid effort, because a wrong oracle is the worst class of
    bug in a kit whose premise is that the gate is the oracle.

12. As a future reader of the resume sweep, I want its doc comment to stop citing
    a spec that no longer exists, so that the pointer resolves. The comment cites
    `specs/worktree-orphan-reconcile.md`, retired earlier; the build is already
    inside that function for story 7.
    Line: `gpt-5.6-luna` / low (`sonnet`). A dead reference in a comment on a
    function this spec is already editing.

## Implementation decisions

**Schema.** `Assignment` gains `created_at` as an optional RFC3339 string with
`omitempty`, under the unchanged `bench-assignment/v1` record schema and the
unchanged ledger schema. Absence is valid, so records written before this build
stay readable and there is no migration step. A ledger schema bump was rejected
for the same reason.

**Validation rejects, and that is load-bearing in both directions.**
`ValidateAssignment` runs on every ledger read, so rejecting a malformed
`created_at` turns one bad record into a total ledger outage for every `bench`
command. That matches the package's existing fail-closed posture and is the
deliberate choice, but it is a reviewer-visible consequence rather than a
detail. A future timestamp is *accepted* here and handled by the predicate
instead; the two owners are disjoint.

**The stamp stays out of the lock reason.** `lockReason` composes a fixed set of
named assignment fields into a string that both `validateCreationBundle` and
`PlanExplicit` re-derive and compare against the live git worktree lock.
Serializing the whole record into it would make every worktree created before
this build fail both its own release and its own explicit cleanup. The field set
stays as it is.

**Fingerprints restale, everywhere.** `LifecycleEvidence` serializes the entire
assignments array into the recovery-retire fingerprint *and* the explicit cleanup
fingerprint, so adding a field invalidates every pending plan and every in-flight
receipt checkpoint comparison. That is the designed plan/apply staleness
behavior, not a regression — but story 5's whole workflow is plan/apply, so
expect it rather than discover it. `automaticFingerprint` composes from named
fields rather than the struct and is unaffected; it stays that way.

**Orphan is age alone.** State `active`, older than the window. No liveness
conjunct, because none is available: `Create` writes no lease, `ProbeLease`
cannot distinguish an absent lease from an unreadable one, and a lease records a
pid that dies with the create hook. The map's #2 records the full argument and
the rejected alternatives.

**The sweep reads the ledger, not the plan.** `PlanAutomatic` returns at its
first retain verdict — ignored residuals, unexpected lock, dirty nested state,
unknown lease — *before* it reaches the assignment-state branch. An orphan
carrying ignored build output is the normal state of a worktree a shift ran in,
so a sweep that derived orphanhood from the plan's reason code would silently
report nothing for exactly the population this row is about. The sweep computes
the predicate directly from the ledger record. Story 4's reason label is layered
on top and is never the source of truth.

**The window is a fixed constant.** `bounds.AssignmentStale` at 7 days, beside
`LeaseStale`. Not configurable; the value is flagged for veto.

**Report, never reap.** The sweep reports; removal stays behind the explicit
path-addressed `bench worktree clean`, which never consults assignment state and
recovers dirty work into a recovery ref before removing. With no liveness signal
this is not a preference — it is the only thing between a long-lived legitimate
worktree and an unattended destructive command.

**Ledger compaction is bounded by preserved work.** The sweep deletes an orphaned
tree-gone record only when it holds no recovery metadata. One holding recovery
metadata is preserved and reported, exactly as `recovered` records are today.

**A pinned contract changes meaning underneath its test.**
`TestResumeSweepsResidueAndReportsPreserved` pins the FT93(c) contract that an
active, tree-gone, unregistered record must survive the sweep. After this build
it survives for a different reason — the record is freshly stamped, so it is not
aged. The test keeps passing while what it pins has moved; the build says so in
the test's comment rather than leaving it reading as untouched.

**Out of this build by design.** The sixteen `recovered` records whose recovery
refs are retained forever are FT98's: their payloads read as unlanded because
that landed proof misses reshaped commits. This build bounds how loudly they
print; it does not touch the refs or the proof.

**Prose is gate-held.** Each new duty gets a `require(<file>, <phrase>)` row in
the conformance anchor registry, so the prose half carries a real red signal.

## Testing decisions

A good test here drives a real temporary repository through the package's
exported entry points and reads the resulting plan, summary, or ledger — never an
internal field. That is how `worktree_test.go`, `ownership_test.go`, and
`lifecycle_test.go` already work, and those are the prior art. The predicate
takes `now` as a parameter, so age is driven by argument rather than by sleeping
or injecting a clock; `reclaimable` sets that precedent.

Gate command: `bench gate` (the project gate). Both halves land inside its
conformance phase — the `internal/worktree` and `internal/intent` package tests,
and the `docs-currency-workflow` check that runs the anchor registry. There is no
manual evidence step.

### Seam diagram

**Seam A — the worktree lifecycle package, driven against a real temp repo.**
Covers stories 2 through 7 and 12.

    trigger: bench resume-clean (session start), or bench worktree clean <path>
        │
        ▼
    assignment record ──▶ [ orphaned(assignment, now)        ] ──▶ retain + ReasonOrphaned
    now               ──▶ [   ├─ PlanAutomatic (label only)  ] ──▶ summary lines + cap
    on-disk tree      ──▶ [   └─ sweepOrphanAssignments      ] ──▶ ledger row deleted
    ledger            ──▶ [        (reads the record direct) ]      or preserved
                             ◀ tests attach here: build a temp repo, create an
                               assignment, set its stamp, call PlanAutomatic /
                               ConservativeCleanup, assert the plan's reason, the
                               rendered summary text, and the surviving ledger rows

**Seam B — assignment record validation.** Covers story 1.

    trigger: every ledger read, and PutAssignment
        │
        ▼
    Assignment{CreatedAt} ──▶ [ ValidateAssignment ] ──▶ nil, or a rejection
                                ◀ tests attach here: table of records — absent,
                                  valid RFC3339, empty string, non-timestamp text,
                                  control bytes, future — asserting accept or reject

**Seam C — the kit-prose anchor registry.** Covers stories 8 through 11.

    trigger: the gate's conformance phase, docs-currency-workflow check
        │
        ▼
    kit markdown files ──▶ [ checkWorkflowAnchors: require(file, phrase) ] ──▶ diagnostics
                             ◀ tests attach here: add a row per duty; deleting or
                               paraphrasing the prose makes the check red

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | an assignment with no `created_at` validates and round-trips through the ledger unchanged | Seam B plus a ledger write/read in a temp repo | already covered by construction — the field does not exist pre-build, so this is a regression guard, not a red-first row; it must be written before story 1's field lands and stay green across it | an implementation that makes the field required rejects every one of the 17 records live today, and because validation runs on every ledger read that is a total outage for every `bench` command |
| 1 | a present-but-unparseable `created_at` is rejected | Seam B | `go test ./internal/intent -run TestAssignmentCreatedAtRejectsMalformed` observed red before build, with an empty string, non-timestamp text, and a control byte as subcases | a field parsed leniently at use-time reads as the zero time, which the predicate treats as maximally old, so a malformed record becomes an orphan verdict |
| 1 | a future `created_at` is accepted by validation, not rejected | Seam B | `go test ./internal/intent -run TestAssignmentCreatedAtAcceptsFuture` observed red before build | the predicate owns clock skew and can only do so if the record survives validation; a validator that rejects it makes story 3's future-stamp row unreachable and turns a skewed host into a ledger outage |
| 2 | a newly created assignment carries a stamp | Seam A, `Create` against a temp repo | `go test ./internal/worktree -run TestCreateStampsAssignment` observed red before build | without the stamp every new record falls to the absent-is-aged rule and is orphan-eligible from the moment the window passes, regardless of use |
| 2 | release and explicit cleanup both still succeed on a worktree whose lock reason predates the field | Seam A, release and `PlanExplicit` against a fixture assignment | already covered by construction — `lockReason` composes five named fields and no faithful implementation reddens this; it is the regression guard for the coupling, written before story 2 lands | an implementation that serializes the whole record into `lockReason` makes both `validateCreationBundle` and `PlanExplicit` re-derive a different string, and every existing worktree becomes both unreleasable and uncleanable — this row's own defect, reintroduced one layer down |
| 3 | an active assignment younger than the window is not orphaned | Seam A, predicate called directly with a fixed `now` | `go test ./internal/worktree -run TestOrphanedRequiresAge` observed red before build | a predicate that ignores age reports every active worktree as abandoned, and with no liveness conjunct nothing else would catch it |
| 3 | an active assignment older than the window is orphaned | Seam A, predicate with `now` past the window | `go test ./internal/worktree -run TestOrphanedOnAgedRecord` observed red before build | this is the whole feature; without it the predicate is a constant `false` and the leak is unchanged |
| 3 | an assignment with no stamp is orphaned | Seam A, predicate with an unstamped record | `go test ./internal/worktree -run TestOrphanedTreatsAbsentStampAsAged` observed red before build | a fail-closed reading of absence leaves today's 17 records immortal, so the build ships and the residue that motivated it never drains |
| 3 | a stamp in the future is not aged | Seam A, predicate with `now` before the stamp | `go test ./internal/worktree -run TestOrphanedRejectsFutureStamp` observed red before build | a subtraction that ignores sign makes clock skew read as enormous age, so a worktree created minutes ago on a skewed host becomes an orphan |
| 3 | an assignment that is not `active` is never orphaned | Seam A, predicate over `cleanup-pending`, `recovered`, `complete` | `go test ./internal/worktree -run TestOrphanedOnlyActiveState` observed red before build | a state-blind predicate double-claims a record already on the recovery path, so the sweep offers a cleanup command for work being preserved |
| 4 | `PlanAutomatic` labels an orphan `orphaned` where it would otherwise say `active` | Seam A, plan against a temp repo with an aged assignment and a clean tree | `go test ./internal/worktree -run TestPlanAutomaticLabelsOrphaned` observed red before build | without it the plan's own explanation still says the record is merely active, contradicting the sweep line that just told the reader to clean it |
| 4 | `PlanAutomatic` keeps the earlier retain reason when one fires first | Seam A, aged assignment carrying ignored residue | `go test ./internal/worktree -run TestPlanAutomaticKeepsEarlierRetainReason` observed red before build | the early return means the state branch is unreachable here; an implementation that "fixes" that by hoisting the orphan label above the safety branches loses the reason that actually blocks removal |
| 5 | the summary emits one `bench worktree clean <path>` line per orphan | Seam A, `ConservativeCleanup` then the rendered summary | `go test ./internal/worktree -run TestResumeSummaryNamesCleanCommand` observed red before build | a summary that only counts orphans leaves the reviewer with the same undecodable wall the row was opened about |
| 5 | an orphan carrying ignored residue still gets its summary line | Seam A, aged assignment with ignored build output | `go test ./internal/worktree -run TestResumeSummaryReportsOrphanWithIgnoredResidue` observed red before build | this is the normal state of a worktree a shift ran in; a sweep that read `PlanAutomatic`'s reason instead of the ledger reports nothing for exactly this population while every other row stays green |
| 5 | the emitted line never contains `--discard-ignored` | Seam A, same summary | `go test ./internal/worktree -run TestResumeSummaryNeverSuggestsDiscardIgnored` observed red before build | the codebase documents that flag's request-less form as orphaning the assignment, so an emitted remedy naming it manufactures the next generation of this defect |
| 5 | a pool path containing a space or a glob character emits as one runnable argument | Seam A, summary over an assignment at such a path | `go test ./internal/worktree -run TestResumeSummaryQuotesHostilePaths` observed red before build | an unquoted path splits into two arguments or glob-expands when pasted, so the emitted command fails or names a different tree |
| 5 | a control byte in a path cannot split the summary's line structure | Seam A, summary over a path containing a newline and an ESC | `go test ./internal/worktree -run TestResumeSummaryPreservesLineStructure` observed red before build, asserting the emitted line count | `renderResumeSummary` is a raw line sink with no safety predicate and this build puts the first attacker-influenced path into it; quoting stops argument splitting but not an embedded newline forging a whole extra summary line |
| 6 | the summary caps orphan and preserved lines and states the true total | Seam A, summary over more records than the cap | `go test ./internal/worktree -run TestResumeSummaryCapsListings` observed red before build, asserting both the line count and the `and <n> more` total | a cap without the count reads as "that is all of them", which is the one way bounding the output could mislead rather than help |
| 7 | an orphaned, tree-gone record with no recovery metadata is deleted and counted as reconciled | Seam A, `ConservativeCleanup` then a ledger read | `go test ./internal/worktree -run TestSweepCompactsOrphanedActiveResidue` observed red before build | today's sweep skips every `active` record unconditionally, which is exactly why one row survived the manual drain |
| 7 | an orphaned, tree-gone record holding recovery metadata is preserved, not deleted | Seam A, same sweep with recovery metadata present | `go test ./internal/worktree -run TestSweepPreservesOrphanedRecoveryRecords` observed red before build | this is the only destructive path in the build; a sweep that drops the preserve branch deletes the ledger's pointer to preserved work and orphans the recovery refs permanently |
| 7 | a freshly stamped active record with a missing tree is still skipped | Seam A, the existing FT93(c) fixture, re-read | already covered — `TestResumeSweepsResidueAndReportsPreserved` asserts this record must survive, and after this build it survives because it is not aged rather than because it is active; the build annotates that shift in the test | a sweep that compacts on tree-absence alone races a session between `worktree add` and its first write, and the existing test would keep passing while the contract it pins had silently changed |
| 7 | a tree-gone record that is still a registered worktree is left to the prune path | Seam A, the existing `isRegisteredWorktree` branch | already covered — the branch and its assertions in `worktree_test.go` | a new orphan branch placed above that check compacts a record git is about to prune, producing a ledger and a registration that disagree |
| 7 | two consecutive sweeps produce the same summary and delete nothing twice | Seam A, `ConservativeCleanup` run twice against one tree | `go test ./internal/worktree -run TestSweepIsIdempotent` observed red before build, comparing both summaries | the summary reprints at every session start and story 6 exists to bound it, so a sweep whose second run reports different counts falsifies the artifact the reviewer actually reads |
| 8 | `craft-delegate` orders releasing a coordinator-cut worktree at done-claim acceptance | Seam C | `BENCH_CONFORMANCE_ROOT=$PWD BENCH_CONFORMANCE_CHECK=docs-currency-workflow go test ./internal/conformance -run TestRootConformance` observed red before the prose lands, via a new `require` row | prose with no anchor row is deletable in one edit, which is how retirement guidance went missing in the first place |
| 9 | `bench-implement-spec`'s stop-short names the worktree retirement owner | Seam C | the same scoped conformance run, observed red before the prose lands, via a new `require` row | same failure mode: an unpinned duty survives exactly until the next rewrite of that section |
| 10 | `.bench/BENCH.md` names `bench worktree release`, `clean`, and `recovery` | Seam C | the same scoped conformance run, observed red before the prose lands, via one `require` row per subcommand | the existing cold-pickup sweep matches top-level routes only, so without these rows the subcommands stay undocumented and nothing notices |
| 11 | the new `require` rows are real obligations, not vacuous ones | Seam C, the conformance fixture-bite sweep | already covered — the registry's existing bite fixtures assert that each anchor can red | a row whose phrase never appears, or appears in a file the walk does not reach, grades green forever and pins nothing |
| 12 | the sweep's doc comment cites no non-existent spec | Seam C, the stale-reference sweep in the gate's docs fragment | already covered — the existing dangling-reference sweep over repository markdown and source comments | a dead pointer sends the next reader to a file git deleted, which is the failure the promote-then-delete rule exists to prevent |

### Edge inventory

Walked against the project profile's hostile-input checklist for shell CLIs.

- **Paths with spaces or glob characters** — coverage row, story 5.
- **Control bytes a sink permits but cannot survive** — the profile's own entry,
  and the one this build newly exposes: `renderResumeSummary` is a line-structured
  sink and this is the first attacker-influenced path to reach it. Coverage row,
  story 5, asserting the emitted line count rather than only rejection.
- **Absent file vs present-but-empty** — for the ledger, absence and an empty
  file are already distinct in `readPath` and unchanged here. For `created_at`,
  absent and empty-string are distinct and both have rows under story 1.
- **Hand-edited file with no trailing newline** — **Won't handle**: the ledger is
  machine-written JSON read through `jsonfile.Decode`, which *requires* a final
  newline and errors without one. That behavior is unchanged by this build.
- **Malformed or hostile parsed value** — coverage row, story 1.
- **Clock skew and a future timestamp** — two rows: accepted at validation
  (story 1), not aged at the predicate (story 3).
- **Re-run idempotency** — coverage row, story 7: two consecutive sweeps,
  summaries compared.
- **A command whose own write changes a fact it reports** — the sweep deletes
  ledger rows and reports counts derived from the ledger. The idempotency row is
  the profile's required "assert repeated application in the tracked
  configuration" for this surface.
- **Interrupt mid-loop** — **Won't handle**: this build adds no multi-step
  transaction. Compaction is a single `DeleteAssignment` under the existing
  ledger lock, and the cleanup transaction it reports on is unchanged.
- **Dangling symlink where a file is expected** — **Won't handle**: this build
  reads no new file. The stamp comes from the ledger the sweep already reads.
- **Special files in a discovered path** — **Won't handle**: no new discovery
  walk is introduced.
- **Required tool missing from PATH** — **Won't handle**: no new external tool.
- **Destructive worktree state** — the preserve-if-recovery branch is the guard,
  with its own row under story 7. Foreign registrations, reused paths, and the
  primary checkout are unchanged upstream of this build.
- **Non-TTY stdin** — **Won't handle**: nothing added here prompts.
- **cwd deeper than the repo root** — **Won't handle**: every read resolves
  through the existing `intent.Address` common-directory query, already
  root-independent.
- **Invocation through every shipped surface** — **Won't handle**: no new
  subcommand or route, so existing routing coverage still holds.
- **Host-backed filesystems under I/O pressure** — **Won't handle**: no new
  `fsync` or durable-write path; compaction reuses the existing locked
  read-modify-write.

## Out of scope

- **FT98's landed proof for reshaped commits.** The sixteen `recovered` records
  keep their recovery refs because payloads main actually shipped still read as
  unlanded. A separate capability — proving landedness by content across a
  rebase — and already its own roadmap row. Estimate: 6 edits, 3 gate runs.
- **A liveness heartbeat for request-created assignments.** Real liveness would
  need something a still-using session writes on a schedule, plus a staleness
  posture and a write on a hot path. It is the only thing that would let the
  window shorten safely, and it is materially larger than this row.
  Estimate: 8 edits, 4 gate runs.
- **A configurable orphan window.** Exposing lifecycle tunables needs a
  configuration surface, a precedence order against the environment, and a
  validation posture — a separate capability the repo has never needed.
  Estimate: 5 edits, 2 gate runs.
- **Renaming `renderResumeSummary`'s "bench resume:" output prefix.** The summary
  names a command that does not exist; the route is `bench resume-clean`. A
  pre-existing defect this build must not copy into new prose, but fixing the
  emitted string touches assertions outside this spec's seams.
  Estimate: 2 edits, 1 gate run.
- **FT149, the branch-delete guard label that quotes `-D` for a `-d`.** Adjacent
  in the same drain but an unrelated defect in a different guard.
  Estimate: 2 edits, 1 gate run.
