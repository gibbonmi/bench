# A chunk repair commits once, and the checkpoint gate alone grades the chunk tip

Status: staged

Roadmap: FT445

Decision source: Named reviewed artifact `roadmap/FT445.md`, with the reviewer approval of 2026-10-23 and the triage decisions of 2026-10-24. FT445 absorbs FT438 (`roadmap/FT438.md`). This spec decides the two open points of the artifact and the engineering seams. The reviewer-named consultants (Opus/high and Fable/medium, with Sonnet/high as the tie-breaker, 2026-10-24) advised those calls. Further notes records each consultant decision and its votes.

Persisted fields: review-record:chunks.verification.source_commit

Verification log: 4 iteration(s) to accept — the spec round took 2 (B1: the prepared run kept the tip-moved guard; B2: the delta operands read the stale record tip), and the ticket round took 2 (three chunks each held a seam ticket and its consumer).

## Problem

One repair cycle of a multi-ticket chunk takes 45 to 60 minutes. The target is 11 to 14 minutes for a one-file test-only repair (`roadmap/FT445.md:5-8`).

Most of that time is the restage. A tip move in the chunk makes every ticket session stage its record again at the new freeze tip. Each session reruns its whole check floor. One restage took 16 to 57 minutes for each tip move and reran 46 commands where 14 were distinct (`roadmap/FT445.md:47`).

Three predicates cause the restage. Each compares an entry with the one tree digest of the chunk:

- `bench record ticket` refuses a staged run whose HEAD is not the chunk source (`internal/reviewrecord/recordcmd/ticket.go:128-131`). It binds an earlier ticket to the tip of the chunk entry (`internal/reviewrecord/ticket.go:109-121`).
- The record writer refuses an entry whose digest differs from the chunk digest (`internal/reviewrecord/change.go:123-125`).
- The checkpoint grades each entry against the chunk digest (`internal/reviewrecord/coverage.go:76-80`, `internal/reviewrecord/check.go:54`).

A verification entry holds a tree digest and no commit (`internal/reviewrecord/record.go:14-38`). So no entry can name the commit that it graded.

Other costs add to the cycle:

- A test-only repair of a delegated build cannot close through `bench record confirmation`. Its assignment record rewrites the spec, and the test-only proof refuses every non-test path (`internal/reviewrecord/confirmation.go:121`, `internal/testgap/testgap.go:50`). FT438 records eight such refusals.
- Each dispatch runs two build preflights of a minute or more: the commit chain and the charge preparation (`roadmap/FT445.md:45`, `cmd/bench/commit_chain.go:42-67`, `internal/preflight/command.go:100-116`).
- The checkpoint gate takes about 4 minutes, and it starts only after the record commit. Twenty gate logs from 2026-10-22 to 2026-10-24 show a typical whole gate of about 243 s. The `test` phase takes 166 to 245 s, and the `system` phase takes 35 to 51 s (`.logs/gate-*.jsonl`).
- A gate green cannot serve the checkpoint after the record commit. The checkpoint oracle frames the checkpoint purpose and the HEAD tip (`internal/gate/checkpoint.go:277-285`), and the subject tree holds the record (`internal/gate/subject.go:67-73`).

## Solution

Each ticket's evidence binds to its own commit. The record keeps that commit, and the checkpoint accepts an entry whose commit is on the chunk line. A repair stages only its own ticket, once, at its repair commit. The whole checkpoint gate at the chunk tip stays the one grade of the tip.

A test-only repair closes through `bench record confirmation` when its plan change is only an assignment record. A plan-fence change, a ticket widening, or any other spec change still takes a confirming round.

The check floor becomes a chunk-level floor. Each repair writes one plan commit, and each dispatch takes one build preflight: the charge preparation.

The orchestrator starts the checkpoint gate when a chunk freezes or a repair commits. The run uses a private checkout of the committed tree without the review record. The chunk checkpoint after the record commit reuses that green, and it still grades the evidence and the record prose on every call.

A test-only delta after a whole green at the reviewed tip takes a delta checkpoint. The delta checkpoint runs gofmt, vet, the changed packages, and root conformance. A delta that touches a tagged or race-table package runs the whole gate. The landing gate always runs whole.

An exec call of a Bench command reuses the worktree's fresh `dist/bench`. The change ships last, after a measurement shows a net saving.

## User stories

Line: opus / high.
Implementation-line reason: The hardest material chunks change the oracle key (RC-C4) and add a reduced verdict class (RC-C6B), and a wrong key makes the oracle lie. The spec is precise, and the existing gate and record fixtures turn a wrong key red. The profile routes gate logic to the mid tier. Seam uncertainty is real in the new `--prepare` form and the delta decision, so the uncertain row of `craft-line` raises the effort to high.
Harder chunks: RC-C2, RC-C2B, RC-C4, RC-C4B, RC-C6, RC-C6B.

### A test-only repair closes without a confirming round

1. As an orchestrator, I want a test-only repair whose delta also holds an assignment-only plan commit to close through `bench record confirmation`, so that a delegated repair does not pay a three-axis confirming round.
2. As a reviewer, I want a change to a plan-fence field in the repair delta to still require a confirming round, so that a changed plan gets a review.
3. As a reviewer, I want a `Writes:` widening in the repair delta to still require a confirming round, so that a wider fence gets a review.
4. As a reviewer, I want any other spec change in the repair delta to still require a confirming round, so that a changed acceptance, such as an added seam cell, gets a review.
5. As a reviewer, I want `bench record confirmation` and the checkpoint to apply one exemption rule, so that the writer and the gate cannot disagree.
6. As a maintainer, I want the known-issues row for this refusal removed and kept out, so that no workaround outlives its issue.

### Evidence binds to its own commit

7. As a ticket author, I want to record my verification and probe entries at my own commit, so that a later ticket's commit does not make me stage again.
8. As a repair session, I want to run each owed check once at my repair commit, so that a tip move does not rerun the whole floor of each ticket.
9. As a repair session of a ticket in any position, I want my record call to refreeze the chunk at my repair commit, so that the checkpoint grades the repaired tip.
10. As an orchestrator, I want a record-only commit to leave every ticket's entries valid, so that the order of record commits costs no restage.
11. As a reviewer, I want the checkpoint to refuse an entry whose commit is not on the chunk line, so that evidence from another chunk or branch cannot pass.
12. As a reviewer, I want the checkpoint to refuse an entry whose digest does not match its commit's tree, so that a relabeled entry cannot pass.
13. As a reviewer, I want each entry to name the commit that it graded, so that the record shows the exact source of each result.
14. As a reviewer, I want another ticket's probe to carry forward to the chunk tip, and the spec to state that this shift applies to probes only, so that the changed guarantee is explicit.
15. As an orchestrator of a record written before this change, I want an entry with no commit to keep today's rule, so that an older record still grades.
16. As a reviewer, I want the ownership premise of reused tests to apply the same commit rule, so that two readers of one record agree.
17. As a maintainer, I want the checkpoint gate at the chunk tip to stay the one grade of the tip, so that carried evidence never replaces the gate run.

### The floor is chunk-level, and each dispatch takes one plan commit and one preflight

18. As a slicer, I want a package that several tickets of a chunk write to need one floor requirement in that chunk, so that the chunk does not run the same package once for each ticket.
19. As a reviewer, I want `check-floor` to name the chunk and the missing element when no requirement of the chunk covers it, so that the gap is clear.
20. As a reviewer, I want a requirement of another chunk to cover nothing, so that a chunk cannot borrow evidence from another chunk.
21. As a maintainer, I want one derivation of requirement ownership, so that the floor and the record agree on who owes each check.
22. As an orchestrator, I want each repair to write one plan commit, with the assignment record first, so that a repair does not pay three commit chains.
23. As an orchestrator, I want each author or repair dispatch to take exactly one build preflight, the charge preparation, so that a dispatch does not grade the same tip twice.
24. As an orchestrator, I want a plan commit that no dispatch follows to keep its full commit chain, so that such a commit still gets a build preflight.

### The checkpoint gate runs beside the review and the confirmation

25. As an orchestrator, I want to start the checkpoint gate when a chunk freezes or a repair commits, so that the gate runs while the review or the confirmation runs.
26. As an orchestrator, I want the chunk checkpoint after the record commit to reuse that green, so that the checkpoint does not run the whole gate again.
27. As a reviewer, I want the reuse to require the same committed tree without the review record and the same oracle, so that any other change runs the gate.
28. As a reviewer, I want the checkpoint to grade the evidence on every call, so that a reused green never skips the record grade.
29. As a reviewer, I want the checkpoint to grade the prose of the review record on every call, so that the record stays graded although the gate key omits it.
30. As an orchestrator, I want a record write or a record commit during the prepared run to leave that run's verdict valid, so that parallel work cannot void the gate.
31. As an orchestrator, I want the guidance to run the coordinator probe before the prepared run, so that the probe does not refuse with `gate-held`.
32. As a reviewer, I want the chunk checkpoint to refuse an uncommitted change, so that the gate grades only committed bytes.
33. As a reviewer, I want the complete checkpoint and the landing to keep their own exact keys, so that a chunk green never authorizes a landing.
34. As a maintainer, I want the current-state documents to state the new chunk-checkpoint key and the delta class, so that a fresh reader sees the decided state.
35. As a reviewer, I want a recorded gate run to stay a native excerpt with its source digest, so that reuse adds no synthesized evidence.

### The selector and the floor see the system suite

36. As a ticket author, I want a changed system-suite test file to select the system check, so that the edited test runs.
37. As a ticket author, I want a production change that the system suite imports to select the system check, so that a reachable system test runs.
38. As a slicer, I want the floor of a ticket that writes the system suite to name `bench test --check system`, so that the plan owes the system run.
39. As a maintainer, I want the tagged-package facts to come from the gate's own system-suite producer, so that no second list drifts.

### A test-only delta takes a delta checkpoint

40. As an orchestrator, I want a test-only delta after a whole green at the reviewed tip to run only gofmt, vet, its changed packages, and root conformance, so that the checkpoint after a test-only repair takes about a minute.
41. As a reviewer, I want the delta package set to come from the test-gap file list and never from the changed-package selector, so that the checkpoint keeps its precondition.
42. As a reviewer, I want a delta with no whole green at its base to run the whole gate, so that a delta never stands alone.
43. As a reviewer, I want a delta with a non-test file to run the whole gate, so that a production change gets the whole grade.
44. As a reviewer, I want a delta that touches a tagged or race-table package to run the whole gate, so that the system and race phases still grade those edits.
45. As a reviewer, I want a delta green to authorize no landing, no complete checkpoint, and no later whole checkpoint, so that a reduced run cannot pass as whole.
46. As a maintainer, I want the delta record to be its own verdict class with a declared reuse refusal, so that the class table states its posture.
47. As an orchestrator, I want `--fresh` to force the whole checkpoint gate, so that I can buy the whole grade on demand.
48. As an orchestrator, I want a green delta checkpoint to open the chunk gap, so that the next chunk starts as after a whole checkpoint.
49. As a reviewer, I want a red delta run to make the checkpoint red, so that a reduced run never hides a failure behind a whole rerun.

### Exec calls reuse the worktree binary

50. As a maintainer, I want the build to measure the run-binary cost before it changes exec, so that a change with no net saving stops.
51. As an orchestrator, I want an exec call of a Bench command to reuse the worktree's fresh `dist/bench`, so that each call does not build a private binary.
52. As an orchestrator, I want a stale `dist/bench` to leave the child to build its own binary and to print `bench worktree build`, so that a stale binary never grades and I know the fix.
53. As a maintainer, I want an exec call of another command to keep today's environment, so that the reuse costs nothing outside Bench.
54. As a reviewer, I want a child to refuse a binary that became stale after the exec check, so that a race cannot grade with a stale binary.

### Scope and evidence

55. As the reviewer, I want the final reconciliation to measure one test-only repair cycle against the 45-to-60-minute baseline, so that the speed claim has evidence.
56. As a maintainer, I want the guidance files to stay inside their line budgets, so that the edits do not grow the loaded prose.
57. As a maintainer, I want the FT438 row to retire with FT445, so that the board holds one owner for the outcome.
58. As a slicer, I want no cap on the tickets of a chunk, so that a chunk stays one coherent outcome now that the restage is gone.

## Implementation decisions

### The FT438 exemption (RC-C1)

`reviewrecord` owns the exemption at the one call site of the test-only proof (`internal/reviewrecord/confirmation.go:121`). `bench record confirmation` and the checkpoint both reach that site through `FindReviewedPair`, so the writer and the gate share one rule.

The rule applies when the spec file is the only non-test path in the delta. It decodes the plan fence of the reviewed spec and of the repaired spec. It clears `execution.assignments` in both and re-encodes both through the payload, locate, and splice path that `RecordAssignment` uses (`internal/reviewrecord/assignment.go:65-77`). It then requires the two spec documents to be byte-equal. The other changed paths go through `testgap.Prove` unchanged.

Any change confined to the assignment block passes, because the reviewer decided on 2026-10-21 that the exemption covers the whole assignment block. A change to any other fence field, a change outside the fence, and a ticket file change refuse. A ticket file stays a non-test path for `testgap.Prove`.

`testgap` gains one exported function that returns the proven test-file list. `Prove` composes it, so the delta checkpoint (RC-C6) and the confirmation read one list.

The FT438 sentence names the review record among the changes that still take a confirming round. `SourceDigest` already strips `reviews/<slug>.md` (`internal/reviewrecord/plan.go:201-210`), so a record change never reaches the proof. This spec follows the tree and flags the clause for reviewer veto.

The chunk removes the FT438 row of `docs/known-issues.md`. A Forbid anchor on `FT438` in that file follows the FT405 precedent (`internal/anchors/registry_ticket_passes.go:15`).

### Commit-bound evidence (RC-C2, RC-C2B)

Each chunk verification entry gains the field `source_commit`. It holds the full commit id of the staged run's HEAD (`internal/responsebound/retain.go:21-27`). The parser refuses a value that is not a full object id. `DisallowUnknownFields` makes an older executable refuse a record with the field (`internal/reviewrecord/parse.go:232`), so this build is a declared build.

One exported `reviewrecord` predicate decides whether an entry is current for a chunk. With a commit C, the entry is current when C strictly descends from the chunk base, C is an ancestor of the chunk tip, and the entry digest equals `SourceDigest` of C's tree. The ancestry query is reflexive, so the predicate refuses C equal to the chunk base by a separate test, as `FindReviewedPair` does (`internal/reviewrecord/confirmation.go:95-98`). The command, performer, role, outcome, exit, and probe-bite rules stay as they are (`internal/reviewrecord/check.go:54-68`).

The predicate does not prove that an entry's commit follows its ticket's latest commit. The performer rule already makes a repair session record each requirement that its ticket owes, because the owed performer is the ticket's latest assignee (`internal/reviewrecord/delegated.go:198-211`). An entry that a repair session stages at an earlier commit needs deliberate contrivance, and the whole checkpoint gate still grades the tip. This spec accepts that posture.

An entry with no commit keeps today's rule: its digest equals the chunk digest. The JSON decoder cannot tell an empty value from an absent one, so an empty value takes the same rule. This legacy rule keeps the two landed records and this build's own RC-C1 entries valid.

These callers adopt the predicate:

- the checkpoint (`internal/reviewrecord/coverage.go:76-80`);
- the record writer (`internal/reviewrecord/change.go:115-125`);
- the stage check of `bench record ticket` (`internal/reviewrecord/recordcmd/ticket.go:120-137`);
- the reused-test premise of ownership evidence (`internal/preflight/ownership/evidence.go:63,208`).

The writer binds the entries to the staged run's commit, not to the chunk tip. A record-only commit never moves the chunk tip (`.agents/commands/bench-implement-spec.md:60`). So a staged run at a record-only descendant of the tip binds the tip's commit.

Any ticket's record call at a commit that descends from the chunk tip refreezes the chunk at that commit. The rule that only the last ticket freezes (`internal/reviewrecord/ticket.go:98-108,134-143,177-181`) goes away. The first record call of a chunk creates the chunk entry.

The dirty-path refusal stays (`internal/reviewrecord/recordcmd/ticket.go:133-137`).

The whole checkpoint gate at the chunk tip stays the grade of the tip. A carried probe proves that its test bit at its own commit, not at the tip. The FT445 artifact accepts this shift for probes (`roadmap/FT445.md:23-24`), and it applies to nothing else.

The shared fixture `recordtest` stops re-retaining entries on a refreeze (`internal/reviewrecord/recordtest/confirmation.go:20-26`).

### The chunk-level floor, one plan commit, and one preflight (RC-C3)

`check-floor` grades the union of a chunk's floor elements against the union of the chunk's requirements. The red detail names the chunk and the element.

`preflight.ticketRequirements` goes away (`internal/preflight/check_floor.go:106-123`). `reviewrecord.TicketRequirements` stays the one ownership derivation (`internal/reviewrecord/ticket.go:39-58`).

Each requirement still names one owner ticket, and a probe requirement stays with its own ticket. `requirementsValid` keeps the rule that every ticket owes one requirement (`internal/reviewrecord/plan.go:193-196`). That rule keeps one source-bound result for each ticket author, and a ticket can meet it with its probe alone.

A repair runs each requirement that its ticket owns once, at its own commit.

The guidance states the repair order. `bench record assignment` comes first, because it refuses a dirty spec or ticket (`internal/reviewrecord/assignment.go:105-117`). Then `bench spec widen` and any floor edit follow, and one plain `bench commit` carries them all.

The charge preparation (`bench preflight build <slug> --charge`) is the one build preflight of an author or repair dispatch. A plan commit that a dispatch follows uses `bench commit` without `--preflight-build`. A plan commit that no dispatch follows keeps `bench commit --preflight-build <slug>`. The commit chain code does not change.

### The prepared checkpoint gate and the chunk-checkpoint key (RC-C4, RC-C4B)

The chunk checkpoint grades the HEAD commit, not the working tree. Before it runs, it refuses a working tree that differs from HEAD with the existing `checkpoint-dirty-checkout` face. The complete form already does this (`internal/gate/complete_checkpoint.go:31-37`).

The chunk-checkpoint subject is the HEAD tree without `reviews/<slug>.md`. `reviewrecord.SourceDigest` derives it, so the record exclusion has one source. The checkpoint runs that tree in a private checkout through the prospective artifact owner (`internal/gate/engine.go:53-66,159-172`). The repair worktree is the storage and identity root.

The chunk-checkpoint oracle is the base oracle of that subject. The purpose JSON and the HEAD tip leave the chunk-checkpoint oracle (`internal/gate/checkpoint.go:277-285`). The complete form keeps its own published-tree oracle and its completion frames.

The tip-moved guard and the evidence pass still run on every chunk-checkpoint call. The evidence pass reads the record from the HEAD tree, not from the stripped subject.

The checkpoint grades the prose of `reviews/<slug>.md` on every call. It reads the record blob and the exclusion-policy blob from the HEAD tree and grades them through `prose.NewGraderFromIndex` and `GradeBytes`. The staged form of `bench gate-prose` is the precedent (`internal/gate/gate_prose.go:132-176`). A HEAD tree with no `.bench/prose-exclusions` has no prose rule, so the checkpoint grades no record prose there. The spec adds no second splitter and no exclusion row.

A new flag starts the prepared run: `bench gate --checkpoint <spec-path> --chunk <id> --prepare`. It runs the oracle on the same subject that the chunk checkpoint grades. It skips the evidence pass, the record prose grade, the tip-moved guard, and the chunk gap. It writes the green to the evidence store under the subject key, and that key holds only the stripped tree that the run captured at its start. A moved HEAD cannot misattribute that green, because a checkpoint at the new HEAD computes its own key (RC37, RC38). A `--prepare` with `--complete`, or with no `--checkpoint`, is a usage error at exit 2. The help line comes from the flag constants (`internal/gate/checkpoint.go:64-69`). The refusal routes reuse the existing gate faces.

A record write or commit during the prepared run cannot void it. The post-run capture reads the private checkout (`internal/gate/evaluation.go:53-62`), and the subject excludes the record. The prepared run has no tip-moved guard, which today resolves HEAD in the worktree on every capture (`internal/gate/evaluation.go:78-127`, `internal/gate/checkpoint.go:249-263`). The fast lane holds no gate lock (`internal/gate/lane.go:116-127`), and no `bench record` form reads it. Only `bench probe` refuses on the gate lock (`internal/probe/set.go:223`), so the guidance runs the coordinator probe before the prepared run.

A chunk checkpoint while a prepared run holds the lock refuses with the existing message (`internal/gate/run_transaction.go:90-96`). The guidance waits on the run's PID before the checkpoint.

The guidance starts a prepared run at two points: when a chunk freezes, beside the review axes, and when a repair commits, beside the confirmation. An axis that runs `bench probe` in the integration worktree meets the same `gate-held` refusal. So the prepared run at a freeze starts only when each axis probes in its own worktree.

A recorded verification of a gate command stays a staged native excerpt with its digest and its source digest (`internal/reviewrecord/recordcmd/stage.go:56-82`, `internal/reviewrecord/parse.go:212-213`). Reuse writes no verification entry.

The landing stays whole and exact. FT314-gate-reuse keeps the landing reuse policy.

These current-state documents change: `.bench/BENCH-reference.md` (lines 339, 547-559, and Gate output), `projects/benchkit.md` (lines 469-473), `.agents/skills/bench-craft-line/references/bounded-repair-policy.md:62`, and ADR 0002 posture 5 ("The subject binds the working tree").

### The system suite in the selector and the floor (RC-C5)

The changed-package index loads `./...` with no tags (`internal/testreport/selection.go:173-201`), so it omits `internal/systemtest`. A changed system-suite file then refuses as "not in a current package" (`internal/testreport/selection.go:233-236`).

The index adds the system-suite package from `gate.SystemSuite` (`internal/gate/phases.go:80-87`). `testreport` already imports `gate`, so this adds no edge. A selected system-suite package maps to `bench test --check system`. The floor reads the same index (`internal/testreport/floor.go:11-25`), so a ticket that writes the system suite owes the system check.

No skip rule reads the selector today (`internal/gate/lane_select.go:132-164`, `projects/benchkit.md:361-374`). This chunk lands before the delta checkpoint, as the FT445 precondition requires.

### The delta checkpoint (RC-C6, RC-C6B)

The chunk checkpoint and the prepared run share one run decision, in this order:

1. Reuse a whole green of the current subject key (RC-C4).
2. Reuse a delta green of the same subject and base, when the delta is still eligible.
3. Run the delta set, when the delta is eligible.
4. Run the whole gate.

`--fresh` skips steps 1 to 3.

The delta has two operands. The base is the record-stripped tree of the reviewed tip. The subject is the record-stripped tree that the run grades.

The reviewed tip comes from the three current axis results of the chunk in the HEAD record. One `reviewrecord` derivation finds it, with the pair checks of `FindReviewedPair`: one pair for all three axes, a tip other than the chunk base, a tip that descends from the chunk base, a tip other than the target, and a tip that is an ancestor of the target (`internal/reviewrecord/confirmation.go:74-112`). The derivation takes the subject's HEAD commit as its target, and it never reads the chunk tip or the chunk digest of the record. At a repair commit, the record still names the earlier freeze tip, so a derivation from the record's chunk tip finds the axis results at the chunk pair and refuses (`internal/reviewrecord/confirmation.go:92-93`). `FindReviewedPair` composes the same derivation with the chunk tip as its target, so the confirmation and the delta share one rule.

A delta is eligible when every condition below holds:

- The test-gap proof, with the RC-C1 exemption, proves the change from the base to the subject test-only.
- A whole green exists for the base under the prospective base oracle of the base subject. The read uses the prospective inspection (`internal/gate/engine.go:90-102`) with the non-expiring window (`internal/gate/engine.go:108-141`), because the freshness window is 60 minutes.
- No changed test file sits in the system-suite package (`gate.SystemSuite`) or in a race-table package (`raceTests`, `internal/gate/gate_go.go:32-40`).

The delta record's `base_tree` and `tree` are exactly the two operands. Step 2 of the run decision reuses a delta green only when both fields equal the current operands, and it checks the base-green condition again.

A test-only change keeps the baseline runner identity, because the freshness digest reads only the non-test build inputs (`internal/freshness/freshness_buildinputs.go:34-48`). So the base subject and the delta subject do not differ in their runner identity.

The delta set runs gofmt, vet, `go test` of the changed packages, and the root-conformance entry. Gofmt and vet run with the argv of their gate phases over the whole root (`internal/gate/phases.go:107-110`). The changed packages are the directories of the files that the test-gap list names. A `_test.go` file always builds with the package of its directory. The gate cannot import `testreport`, because `testreport` imports `gate` (`internal/testreport/floor.go:7`). So no checkpoint path can call the changed-package selector.

Root conformance runs every dev-tier registry check, the meta checks included. The gate runs it through its own entry (`internal/gate/phases.go:147-164`, `rootConformanceEnv`). That is what "conformance" means in FT445 change 7.

A green delta writes a new verdict record class shaped like the lane record (`internal/gate/lane.go:52-55,329-345`). Its fields are `schema`, `delta`, `base_tree`, `tree`, `recorded_at`, and `run_binary`, with no `oracle` and no `status`. The class row declares a non-empty reuse refusal. Each general reuse reader refuses the class (`internal/gate/composed_green.go:33`, `internal/gate/verdict.go:120-129`). Only step 2 of the chunk-checkpoint run decision reads `bench-last-delta`.

The record goes to its own file in the worktree git dir, `bench-last-delta`. It never goes to the evidence store or to `bench-last-gate`. So the landing (`internal/gate/engine.go:70-88`), the complete checkpoint, and a later whole checkpoint read the store and find no delta green. The complete form never enters this run decision, and it keeps its published-tree run.

A red delta makes the checkpoint red, with no fallback to the whole gate. A green delta opens the chunk gap as a whole chunk checkpoint does (`internal/gate/gate.go:213-215`).

The gate output prints the delta phases table and the line `gate: green (delta checkpoint)`. `.bench/BENCH-reference.md` Gate output owns that account. `projects/benchkit.md` (lines 339-348, 387-389, 469-472) names the delta schedule as the one reduced schedule and the delta class as the one reduced record.

### Run-binary reuse in exec (RC-C7)

The first ticket measures the wall time of `freshness.Digest`, of one run-binary selection, and the count of exec calls that select a binary in one test-only repair. If the measured net saving for each repair is not positive, the build stops this chunk and returns to the spec.

When the base name of the child's `argv[0]` is `bench`, `exec` runs `freshness.VerifyAt(worktree, worktree, dist/bench)` (`internal/freshness/freshness_verify.go:23-40`). On a pass, it sets `BENCH_RUN_BINARY` to the worktree's `dist/bench`. On a refusal, it omits the variable, the child builds its private selection as today, and exec prints `bench worktree build` on stderr. For any other `argv[0]`, exec runs no freshness check and keeps today's environment (`internal/worktree/exec.go:244-252`).

The child's `ReuseOrOwn`, `Inherit`, and `Validate` sequence checks the source digest again (`internal/runbinary/runbinary.go:155-183`). A binary that became stale after the exec check refuses, and the refusal names `bench worktree build`.

The worktree assignment owns `dist/bench`, and `Close` does nothing for an inherited selection (`internal/runbinary/runbinary.go:232-235`). The spec adds no new process owner. The comment at `internal/worktree/exec.go:217-229` changes, because `Inherit` with a named root checks the source digest.

`uncertain`: the wrapper may treat an inherited `BENCH_RUN_BINARY` as its own resolution (`internal/systemtest/owner_environment_test.go:38`). The measurement ticket reads that route before the second ticket starts.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RC-C1 / 01-exempt-the-assignment-block.md | A test-only repair with an assignment-only plan commit closes through the confirmation | RC1, RC2, RC3, RC4, RC5, RC6, RC7 | Confirmation refusals, checkpoint confirmation, known-issues anchor | no |
| RC-C2 / 02-grade-entries-at-their-own-commit.md | The checkpoint and the evidence readers grade each entry at its own commit | RC8, RC13, RC14, RC15, RC17, RC18, RC19, RC20, RC21, RC22, RC23 | Checkpoint evidence, ownership premise, parser | yes |
| RC-C2B / 03-record-entries-at-the-own-commit.md | The record writer binds each entry to its own commit, with no restage | RC9, RC10, RC11, RC12, RC16, RC24 | Record ticket, record verification, guidance | yes |
| RC-C3 / 04-grade-the-floor-per-chunk.md, 05-one-plan-commit-and-preflight-per-dispatch.md | Chunk-level floor, one plan commit for each repair, and one preflight for each dispatch | RC25, RC26, RC27, RC28, RC29, RC30, RC31, RC32, RC33, RC84 | Check floor, assignment order, guidance budgets | no |
| RC-C4 / 06-refuse-a-dirty-chunk-checkout.md, 07-key-the-chunk-checkpoint-without-the-record.md | The chunk checkpoint grades the committed tree without the record, and a record-only commit reuses its green | RC39, RC40, RC41, RC44, RC46, RC48 | Dirty checkout, reuse, evidence pass, record prose, landing key | yes |
| RC-C4B / 08-prepare-the-checkpoint-gate.md | The prepared gate runs beside the review or confirmation, and the checkpoint reuses its green | RC34, RC35, RC36, RC37, RC38, RC42, RC43, RC45, RC47, RC49, RC50, RC51, RC52 | Prepare, usage, freshness window, guidance | yes |
| RC-C5 / 09-list-the-system-suite-in-the-selector.md | The selector and the floor select the system check | RC53, RC54, RC55, RC56, RC57 | Changed selection, writes floor | no |
| RC-C6 / 10-run-the-delta-checkpoint.md | A test-only delta takes a delta checkpoint that no landing accepts | RC58, RC59, RC60, RC61, RC62, RC63, RC64, RC65, RC66, RC67, RC71, RC72, RC73 | Delta checkpoint, delta operands, eligibility | yes |
| RC-C6B / 11-record-and-reuse-the-delta-green.md | The delta green has its own record class, and only the chunk checkpoint reuses it | RC68, RC69, RC70, RC74, RC75, RC82, RC83 | Verdict class guard, complete checkpoint, prepared delta | yes |
| RC-C7 / 12-measure-the-run-binary-cost.md, 13-reuse-the-worktree-binary-in-exec.md | Exec calls of Bench commands reuse a fresh worktree binary | RC76, RC77, RC78, RC79, RC80, RC81 | Exec environment, run-binary inherit | no |

RC-C1 comes first, because change 3 needs the exemption. A seam chunk closes its review before its consumer chunk starts (`craft-tickets`): RC-C2 before RC-C2B, RC-C4 before RC-C4B, and RC-C6 before RC-C6B. RC-C2 and RC-C2B land before RC-C3, because the chunk-level floor assumes that each entry is graded at its own commit. RC-C4 and RC-C4B land before RC-C6, because the delta checkpoint shares the prepared run decision and the subject key. RC-C5 lands before RC-C6, as precondition 3 requires. RC-C7 ships last, and its measurement can stop it.

### Completion plan

The authored version 1 plan states the verification obligations. The approved build records its fresh authors through the version 2 amendment before dispatch. From RC-C4 on, each chunk checkpoint of this build runs with `--fresh`, because the candidate under build owns the new key.

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft445-repair-cycle-full-20261024","orchestrator_session":"claude:rc-orchestrator-20261024","author_limit":3,"assignments":{"01-exempt-the-assignment-block.md":[{"session":"claude:rc-t01-author-20261024","assignment":"aa7435c8d0031009239fa5fb68b92b67","model":"opus","effort":"high","source":"183dd4f2e22796bb88e30b7ac126bf4cef59830d","native_ref":"claude-agent:bench-writer-high:rc-t01-author-20261024"},{"session":"claude:rc-t01-repair1-20261024","assignment":"aa7435c8d0031009239fa5fb68b92b67","model":"opus","effort":"high","source":"c01f9b8cb1fd92b5ecbad6a302f4b493cc6827bc","native_ref":"claude-agent:bench-writer-high:rc-t01-repair1-20261024","predecessor":"claude:rc-t01-author-20261024","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:rc-t01-author-20261024","preserved":"c01f9b8cb1fd92b5ecbad6a302f4b493cc6827bc"}],"02-grade-entries-at-their-own-commit.md":[{"session":"claude:rc-t02-author-20261024","assignment":"c8cfd1b42d799126bbf9a6aff35d1ab3","model":"opus","effort":"xhigh","source":"b410a3f3dbafd5fe91b7b0cf7887f30eb99cc3ca","native_ref":"claude-agent:bench-writer-xhigh:rc-t02-author-20261024"},{"session":"claude:rc-t02-repair1-20261024","assignment":"c8cfd1b42d799126bbf9a6aff35d1ab3","model":"opus","effort":"xhigh","source":"5bd2a526d67074a4bd07024f52886fcd5801323c","native_ref":"claude-agent:bench-writer-xhigh:rc-t02-repair1-20261024","predecessor":"claude:rc-t02-author-20261024","trigger":"user-directed","stopped":"claude-agent:bench-writer-xhigh:rc-t02-author-20261024","preserved":"5bd2a526d67074a4bd07024f52886fcd5801323c"}],"03-record-entries-at-the-own-commit.md":[{"session":"claude:rc-t03-author-20261024","assignment":"5a03d130da674f0ad7b251f37506e5fc","model":"opus","effort":"xhigh","source":"0f637ca8879c83c03422af7b9f83516537160a1a","native_ref":"claude-agent:bench-writer-xhigh:rc-t03-author-20261024"},{"session":"claude:rc-t03-repair1-20261024","assignment":"5a03d130da674f0ad7b251f37506e5fc","model":"opus","effort":"xhigh","source":"5fd813ea21d2a90b71a59b04575e64a9c72a8572","native_ref":"claude-agent:bench-writer-xhigh:rc-t03-repair1-20261024","predecessor":"claude:rc-t03-author-20261024","trigger":"user-directed","stopped":"claude-agent:bench-writer-xhigh:rc-t03-author-20261024","preserved":"5fd813ea21d2a90b71a59b04575e64a9c72a8572"}],"04-grade-the-floor-per-chunk.md":[{"session":"claude:rc-t04-author-20261024","assignment":"d5312d8cb94033f49954bd1df32dc1a3","model":"opus","effort":"high","source":"7aa9256c802a5f6cfe4d9e1bb8a3dcf1ef21aa3c","native_ref":"claude-agent:bench-writer-high:rc-t04-author-20261024"},{"session":"claude:rc-t04-repair1-20261024","assignment":"d5312d8cb94033f49954bd1df32dc1a3","model":"opus","effort":"high","source":"44a4e21286b329632d7ea6eeaa5a3e0030960753","native_ref":"claude-agent:bench-writer-high:rc-t04-repair1-20261024","predecessor":"claude:rc-t04-author-20261024","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:rc-t04-author-20261024","preserved":"44a4e21286b329632d7ea6eeaa5a3e0030960753"}],"05-one-plan-commit-and-preflight-per-dispatch.md":[{"session":"claude:rc-t05-author-20261024","assignment":"78e692ab9f542d168c6f7257aa532f4d","model":"opus","effort":"high","source":"1d02223fff82def1d12c0c284f0a550914b2ffd6","native_ref":"claude-agent:bench-writer-high:rc-t05-author-20261024"}],"06-refuse-a-dirty-chunk-checkout.md":[{"session":"claude:rc-t06-author-20261024","assignment":"28903206f6c64deb4ceb82370c9083ff","model":"opus","effort":"xhigh","source":"b77e4be404e76eef13b3e0f714ee80cce08bd2d9","native_ref":"claude-agent:bench-writer-xhigh:rc-t06-author-20261024"}],"07-key-the-chunk-checkpoint-without-the-record.md":[{"session":"claude:rc-t07-author-20261024","assignment":"2c6d8ab766ce223656ec23a8c9909518","model":"opus","effort":"xhigh","source":"0b97628459ab7f45c84b760459cc812a19fd4885","native_ref":"claude-agent:bench-writer-xhigh:rc-t07-author-20261024"},{"session":"claude:rc-t07-repair1-20261024","assignment":"2c6d8ab766ce223656ec23a8c9909518","model":"opus","effort":"xhigh","source":"a91c31f295b834be59377652e75988f877ebe436","native_ref":"claude-agent:bench-writer-xhigh:rc-t07-repair1-20261024","predecessor":"claude:rc-t07-author-20261024","trigger":"user-directed","stopped":"claude-agent:bench-writer-xhigh:rc-t07-author-20261024","preserved":"a91c31f295b834be59377652e75988f877ebe436"}],"08-prepare-the-checkpoint-gate.md":[{"session":"claude:rc-t08-author-20261024","assignment":"675e85cf4f6cc4a6822b618debd57b6b","model":"opus","effort":"xhigh","source":"1f6485e7cc148bbf23546f69221bfca25b1ce6b0","native_ref":"claude-agent:bench-writer-xhigh:rc-t08-author-20261024"}],"09-list-the-system-suite-in-the-selector.md":[],"10-run-the-delta-checkpoint.md":[],"11-record-and-reuse-the-delta-green.md":[],"12-measure-the-run-binary-cost.md":[],"13-reuse-the-worktree-binary-in-exec.md":[]}},"chunks":[{"id":"RC-C1","tickets":["01-exempt-the-assignment-block.md"],"verification":[{"id":"record","command":"bench test --package ./internal/reviewrecord","ticket":"01-exempt-the-assignment-block.md"},{"id":"record-command","command":"bench test --package ./internal/reviewrecord/recordcmd","ticket":"01-exempt-the-assignment-block.md"},{"id":"gate","command":"bench test --package ./internal/gate","ticket":"01-exempt-the-assignment-block.md"},{"id":"testgap","command":"bench test --package ./internal/testgap","ticket":"01-exempt-the-assignment-block.md"},{"id":"anchors","command":"bench test --package ./internal/anchors","ticket":"01-exempt-the-assignment-block.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"01-exempt-the-assignment-block.md"},{"id":"floor","command":"bench test --package ./internal/reviewrecord/recordtest --package ./cmd/bench","ticket":"01-exempt-the-assignment-block.md"},{"id":"exemption-probe","command":"bench test --package ./internal/reviewrecord/recordcmd","probe":"clear the whole plan fence instead of the assignment block","ticket":"01-exempt-the-assignment-block.md"}]},{"id":"RC-C2","tickets":["02-grade-entries-at-their-own-commit.md"],"verification":[{"id":"record","command":"bench test --package ./internal/reviewrecord","ticket":"02-grade-entries-at-their-own-commit.md"},{"id":"record-command","command":"bench test --package ./internal/reviewrecord/recordcmd","ticket":"02-grade-entries-at-their-own-commit.md"},{"id":"record-fixture","command":"bench test --package ./internal/reviewrecord/recordtest","ticket":"02-grade-entries-at-their-own-commit.md"},{"id":"ownership","command":"bench test --package ./internal/preflight/ownership","ticket":"02-grade-entries-at-their-own-commit.md"},{"id":"gate","command":"bench test --package ./internal/gate","ticket":"02-grade-entries-at-their-own-commit.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"02-grade-entries-at-their-own-commit.md"},{"id":"floor","command":"bench test --package ./internal/landing --package ./internal/worktree --package ./cmd/bench","ticket":"02-grade-entries-at-their-own-commit.md"},{"id":"line-probe","command":"bench test --package ./internal/gate","probe":"drop the ancestry test of the entry commit","ticket":"02-grade-entries-at-their-own-commit.md"}]},{"id":"RC-C2B","tickets":["03-record-entries-at-the-own-commit.md"],"verification":[{"id":"record","command":"bench test --package ./internal/reviewrecord","ticket":"03-record-entries-at-the-own-commit.md"},{"id":"record-command","command":"bench test --package ./internal/reviewrecord/recordcmd","ticket":"03-record-entries-at-the-own-commit.md"},{"id":"record-fixture","command":"bench test --package ./internal/reviewrecord/recordtest","ticket":"03-record-entries-at-the-own-commit.md"},{"id":"gate","command":"bench test --package ./internal/gate","ticket":"03-record-entries-at-the-own-commit.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"03-record-entries-at-the-own-commit.md"},{"id":"floor","command":"bench test --package ./internal/anchors --package ./cmd/bench","ticket":"03-record-entries-at-the-own-commit.md"},{"id":"writer-probe","command":"bench test --package ./internal/reviewrecord/recordcmd","probe":"bind an earlier ticket to the chunk tip","ticket":"03-record-entries-at-the-own-commit.md"}]},{"id":"RC-C3","tickets":["04-grade-the-floor-per-chunk.md","05-one-plan-commit-and-preflight-per-dispatch.md"],"verification":[{"id":"preflight","command":"bench test --package ./internal/preflight","ticket":"04-grade-the-floor-per-chunk.md"},{"id":"record-command","command":"bench test --package ./internal/reviewrecord/recordcmd","ticket":"04-grade-the-floor-per-chunk.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"04-grade-the-floor-per-chunk.md"},{"id":"floor","command":"bench test --package ./internal/preflight/preflighttest --package ./internal/reviewrecord --package ./cmd/bench","ticket":"04-grade-the-floor-per-chunk.md"},{"id":"floor-probe","command":"bench test --package ./internal/preflight","probe":"grade each ticket against its own requirements only","ticket":"04-grade-the-floor-per-chunk.md"},{"id":"guidance-conformance","command":"bench test --check root-conformance","ticket":"05-one-plan-commit-and-preflight-per-dispatch.md"},{"id":"guidance-floor","command":"bench test --package ./internal/anchors --package ./cmd/bench","ticket":"05-one-plan-commit-and-preflight-per-dispatch.md"}]},{"id":"RC-C4","tickets":["06-refuse-a-dirty-chunk-checkout.md","07-key-the-chunk-checkpoint-without-the-record.md"],"verification":[{"id":"dirty-gate","command":"bench test --package ./internal/gate","ticket":"06-refuse-a-dirty-chunk-checkout.md"},{"id":"dirty-conformance","command":"bench test --check root-conformance","ticket":"06-refuse-a-dirty-chunk-checkout.md"},{"id":"dirty-floor","command":"bench test --package ./internal/refusalroute --package ./internal/reviewrecord/recordtest --package ./internal/worktree --package ./cmd/bench","ticket":"06-refuse-a-dirty-chunk-checkout.md"},{"id":"gate","command":"bench test --package ./internal/gate","ticket":"07-key-the-chunk-checkpoint-without-the-record.md"},{"id":"record","command":"bench test --package ./internal/reviewrecord","ticket":"07-key-the-chunk-checkpoint-without-the-record.md"},{"id":"command","command":"bench test --package ./cmd/bench","ticket":"07-key-the-chunk-checkpoint-without-the-record.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"07-key-the-chunk-checkpoint-without-the-record.md"},{"id":"system","command":"bench test --check system","ticket":"07-key-the-chunk-checkpoint-without-the-record.md"},{"id":"floor","command":"bench test --package ./internal/reviewrecord/recordcmd","ticket":"07-key-the-chunk-checkpoint-without-the-record.md"},{"id":"key-probe","command":"bench test --package ./internal/gate","probe":"frame the HEAD tip into the chunk-checkpoint oracle","ticket":"07-key-the-chunk-checkpoint-without-the-record.md"},{"id":"prose-probe","command":"bench test --package ./internal/gate","probe":"skip the record prose grade after a reuse","ticket":"07-key-the-chunk-checkpoint-without-the-record.md"}]},{"id":"RC-C4B","tickets":["08-prepare-the-checkpoint-gate.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"08-prepare-the-checkpoint-gate.md"},{"id":"command","command":"bench test --package ./cmd/bench","ticket":"08-prepare-the-checkpoint-gate.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"08-prepare-the-checkpoint-gate.md"},{"id":"system","command":"bench test --check system","ticket":"08-prepare-the-checkpoint-gate.md"},{"id":"floor","command":"bench test --package ./internal/refusalroute --package ./internal/anchors","ticket":"08-prepare-the-checkpoint-gate.md"},{"id":"tip-probe","command":"bench test --package ./internal/gate","probe":"keep the tip-moved guard in the prepared run","ticket":"08-prepare-the-checkpoint-gate.md"}]},{"id":"RC-C5","tickets":["09-list-the-system-suite-in-the-selector.md"],"verification":[{"id":"testreport","command":"bench test --package ./internal/testreport","ticket":"09-list-the-system-suite-in-the-selector.md"},{"id":"preflight","command":"bench test --package ./internal/preflight","ticket":"09-list-the-system-suite-in-the-selector.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"09-list-the-system-suite-in-the-selector.md"},{"id":"system","command":"bench test --check system","ticket":"09-list-the-system-suite-in-the-selector.md"},{"id":"floor","command":"bench test --package ./internal/conformance --package ./internal/gate --package ./internal/reviewrecord --package ./internal/reviewrecord/recordcmd","ticket":"09-list-the-system-suite-in-the-selector.md"},{"id":"index-probe","command":"bench test --package ./internal/testreport","probe":"drop the system-suite package from the index","ticket":"09-list-the-system-suite-in-the-selector.md"}]},{"id":"RC-C6","tickets":["10-run-the-delta-checkpoint.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"10-run-the-delta-checkpoint.md"},{"id":"testgap","command":"bench test --package ./internal/testgap","ticket":"10-run-the-delta-checkpoint.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"10-run-the-delta-checkpoint.md"},{"id":"system","command":"bench test --check system","ticket":"10-run-the-delta-checkpoint.md"},{"id":"floor","command":"bench test --package ./internal/gate/authorization --package ./internal/reviewrecord","ticket":"10-run-the-delta-checkpoint.md"},{"id":"tag-probe","command":"bench test --package ./internal/gate","probe":"admit a system-suite file to the delta","ticket":"10-run-the-delta-checkpoint.md"},{"id":"operand-probe","command":"bench test --package ./internal/gate","probe":"take the delta base from the record chunk tip","ticket":"10-run-the-delta-checkpoint.md"}]},{"id":"RC-C6B","tickets":["11-record-and-reuse-the-delta-green.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"11-record-and-reuse-the-delta-green.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"11-record-and-reuse-the-delta-green.md"},{"id":"system","command":"bench test --check system","ticket":"11-record-and-reuse-the-delta-green.md"},{"id":"floor","command":"bench test --package ./internal/git --package ./internal/status --package ./internal/testreport","ticket":"11-record-and-reuse-the-delta-green.md"},{"id":"store-probe","command":"bench test --package ./internal/gate","probe":"write the delta green to the evidence store","ticket":"11-record-and-reuse-the-delta-green.md"}]},{"id":"RC-C7","tickets":["12-measure-the-run-binary-cost.md","13-reuse-the-worktree-binary-in-exec.md"],"verification":[{"id":"measure-conformance","command":"bench test --check root-conformance","ticket":"12-measure-the-run-binary-cost.md"},{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"13-reuse-the-worktree-binary-in-exec.md"},{"id":"runbinary","command":"bench test --package ./internal/runbinary","ticket":"13-reuse-the-worktree-binary-in-exec.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"13-reuse-the-worktree-binary-in-exec.md"},{"id":"system","command":"bench test --check system","ticket":"13-reuse-the-worktree-binary-in-exec.md"},{"id":"floor","command":"bench test --package ./internal/freshness --package ./internal/treetarget --package ./internal/usage --package ./cmd/bench","ticket":"13-reuse-the-worktree-binary-in-exec.md"},{"id":"stale-probe","command":"bench test --package ./internal/worktree","probe":"pass a stale worktree binary to the child","ticket":"13-reuse-the-worktree-binary-in-exec.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/chunk-repair-cycle/spec.md"},{"id":"record","command":"bench test --package ./internal/reviewrecord"},{"id":"record-command","command":"bench test --package ./internal/reviewrecord/recordcmd"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"conformance","command":"bench test --check root-conformance"},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

- A good test drives the real record writer, the real checkpoint, and the real gate entry over a disposable repository. It observes exit codes, refusal text, record bytes, and evidence-store files. It reads no private field.
- The record seams use the `recordtest` fixture and its `Stage`, `Run`, and chunk helpers (`internal/reviewrecord/recordcmd/stage_test.go:52-189`). The checkpoint seams use the gate checkpoint fixtures (`internal/gate/review_checkpoint_test.go`, `internal/gate/complete_checkpoint_test.go`).
- A reuse row runs the prepared gate and the checkpoint as two process calls, so the second call reloads the evidence store. In-memory reuse cannot pass it.
- A row that compares a run against a reuse uses the phase rows and the reuse line, never elapsed times.
- A delta row plants a red in a package outside the delta set. The delta must stay green there, and the whole run must go red. That proves the delta ran only its set.
- The gate seam observes every change of this spec: the checkpoint gate, the prepared run, and the delta checkpoint all run through `bench gate`.
- Prior art: `TestReviewCheckpointReuse` (reuse refusals), `TestReviewCheckpointConfirmationRefusals` (test-only proof), `TestTicketRefusals` (stage refusals), `TestCheckFloorReadsOnlyTheOwnChunk` (floor scope), `TestVerdictRecordClassRegistryMatchesExpectation` (class table), `TestExecChildOwnsRatherThanInheritsItsRunBinary` (exec environment).
- No fixed wall-clock threshold enters the gate. The final reconciliation measures the repair cycle as review-owned evidence.

### Seam diagram

```text
trigger: a ticket author or a repair session
    │
    ▼
staged run at commit C ──▶ [ bench record ticket ] ──▶ entry {source_commit: C, source_digest}
                              ◀ tests: recordtest fixture, record bytes, refusal text

trigger: the orchestrator after a freeze or a repair commit
    │
    ▼
HEAD tree minus record ──▶ [ bench gate --checkpoint --chunk --prepare ] ──▶ evidence store green
                              ◀ tests: second-process checkpoint reuse, record commit during the run

trigger: the orchestrator before a successor chunk
    │
    ▼
HEAD commit ──▶ [ chunk checkpoint: reuse | delta | whole ] ──▶ evidence pass + record prose + chunk gap
                   ◀ tests: planted red outside the delta, store and cache files, landing inspection

trigger: bench worktree exec with a Bench argv[0]
    │
    ▼
dist/bench + seal ──▶ [ exec freshness check ] ──▶ BENCH_RUN_BINARY or a private build
                        ◀ tests: child environment, stale-binary refusal
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RC1 | 1 | A repair delta of one test file plus an assignment-only plan commit records a confirmation. | `internal/reviewrecord/recordcmd/confirmation_test.go` (`TestRecordConfirmationAcceptsAnAssignmentOnlyPlanCommit`). | Today the spec path fails the test-only proof, so the call refuses. |
| RC2 | 2 | A repair delta that changes a plan-fence field outside the assignment block refuses the confirmation. | `internal/reviewrecord/recordcmd/confirmation_test.go` (`TestRecordConfirmationRefusesAPlanFenceChange`). | An exemption that clears the whole fence accepts the changed field. |
| RC3 | 3 | A repair delta that widens a ticket `Writes:` line refuses the confirmation. | `internal/reviewrecord/recordcmd/confirmation_test.go` (`TestRecordConfirmationRefusesATicketWidening`). | An exemption that admits every spec-folder path accepts the ticket file. |
| RC4 | 4 | A repair delta that changes a seam cell outside the plan fence refuses the confirmation. | `internal/reviewrecord/recordcmd/confirmation_test.go` (`TestRecordConfirmationRefusesASpecProseChange`). | An exemption that compares only the decoded fence accepts the changed prose. |
| RC5 | 5 | The chunk checkpoint accepts the assignment-only confirmation from RC1. | `internal/gate/review_checkpoint_confirmation_test.go` (`TestReviewCheckpointConfirmationAcceptsAnAssignmentOnlyPlanCommit`). | An exemption in the writer alone leaves the checkpoint refusing the same delta. |
| RC6 | 5 | A repair delta with a non-test Go file still refuses the confirmation. | `internal/reviewrecord/recordcmd/confirmation_test.go` (`TestRecordConfirmationRefusals`). | An exemption that skips the test-gap proof accepts the production file. |
| RC7 | 6 | The known-issues file holds no row that names FT438. | Gate check: the `docs-currency-workflow` anchors with a Forbid row on `FT438` in docs/known-issues.md. | A kept or restored row matches the forbidden needle. |
| RC8 | 7 | The checkpoint accepts an earlier ticket's entry at that ticket's commit after a later ticket moves the tip. | `internal/gate/review_checkpoint_test.go` (`TestCheckpointAcceptsAnEntryAtAnEarlierTicketCommit`). | Today the entry digest differs from the chunk digest, so the checkpoint calls it stale. |
| RC9 | 7 | An earlier ticket's `bench record ticket` at its own commit writes its entries after the tip moved past it. | `internal/reviewrecord/recordcmd/ticket_test.go` (`TestTicketRecordsAtItsOwnCommitOnTheChunkLine`). | Today the call refuses because the staged run is not at the chunk source. |
| RC10 | 8 | A repair session's record call writes entries for its own ticket only and leaves the other entries unchanged. | `internal/reviewrecord/recordcmd/ticket_test.go` (`TestRepairRecordLeavesOtherTicketsEntries`). | A writer that restages the chunk rewrites the other tickets' entries. |
| RC11 | 9 | A non-last ticket's record call at a commit that descends from the chunk tip refreezes the chunk at that commit. | `internal/reviewrecord/recordcmd/ticket_test.go` (`TestNonLastTicketRefreezesAtItsRepairCommit`). | Today only the last ticket can freeze, so the call refuses. |
| RC12 | 10 | A staged run at a record-only descendant of the chunk tip binds the tip commit and leaves the tip unchanged. | `internal/reviewrecord/recordcmd/ticket_test.go` (`TestTicketStagedAfterARecordCommitBindsTheTip`). | A writer that binds HEAD literally refreezes the chunk on a record commit. |
| RC13 | 11 | The checkpoint refuses an entry whose commit is not an ancestor of the chunk tip. | `internal/gate/review_checkpoint_test.go` (`TestCheckpointRefusesAnEntryOffTheChunkLine`). | A predicate without the tip ancestry test accepts a commit from another branch. |
| RC14 | 11 | The checkpoint refuses an entry whose commit is the chunk base. | `internal/gate/review_checkpoint_test.go` (`TestCheckpointRefusesAnEntryAtTheChunkBase`). | A predicate that uses the reflexive ancestry query alone accepts the base, which is the tip of the earlier chunk. |
| RC15 | 12 | The checkpoint refuses an entry whose digest differs from the digest of its commit's tree. | `internal/gate/review_checkpoint_test.go` (`TestCheckpointRefusesARelabeledEntryDigest`). | A predicate that trusts the recorded commit alone accepts a relabeled entry. |
| RC16 | 13 | A written entry carries `source_commit` equal to the staged run's HEAD. | `internal/reviewrecord/recordcmd/ticket_test.go` (`TestTicketWritesTheSourceCommit`). | A writer that omits the field drops every new entry to the legacy rule. |
| RC17 | 14 | The checkpoint accepts another ticket's probe entry at its own commit after a repair moved the tip. | `internal/gate/review_checkpoint_test.go` (`TestCheckpointCarriesAnotherTicketsProbeForward`). | A predicate that demands each probe at the tip refuses the carried probe. |
| RC18 | 15 | The checkpoint accepts an entry with no `source_commit` at the chunk digest. | `internal/gate/review_checkpoint_test.go` (`TestCheckpointKeepsTheLegacyEntryRule`). | A predicate that requires the new field refuses every older record. |
| RC19 | 15 | The checkpoint refuses an entry with no `source_commit` at a digest other than the chunk digest. | `internal/gate/review_checkpoint_test.go` (`TestCheckpointRefusesAStaleLegacyEntry`). | A legacy rule that accepts any digest lets a stale older entry pass. |
| RC20 | 16 | The reused-test premise accepts an entry at an earlier commit on the chunk line. | `internal/preflight/ownership/premises_test.go` (`TestReusedTestPremiseAcceptsAnEarlierTicketCommit`). | Today the premise requires the chunk tip and the chunk digest. |
| RC21 | 16 | The reused-test premise refuses an entry whose commit is not on the chunk line. | `internal/preflight/ownership/premises_test.go` (`TestReusedTestPremiseRefusesAnEntryOffTheChunkLine`). | A premise that drops the tip test without the ancestry test accepts any commit. |
| RC22 | 17 | The checkpoint goes red on a red planted at the chunk tip when every entry is current. | `internal/gate/review_checkpoint_test.go` (`TestCheckpointRunsTheGateWithCarriedEntries`). | A checkpoint that skips the gate when all entries carry forward misses the planted red. |
| RC23 | 13 | The record parser refuses a `source_commit` that is not a full object id. | `internal/reviewrecord/record_test.go` (`TestParseRefusesAMalformedSourceCommit`). | A parser with no form rule accepts a short or hostile commit text. |
| RC24 | 7 | A staged run with a dirty path other than the record still refuses. | `internal/reviewrecord/recordcmd/ticket_test.go` (`TestTicketRefusals`). | Removing the dirty-path predicate with the source predicate accepts an uncommitted change. |
| RC25 | 18 | `check-floor` is green when one ticket's requirement covers a package that another ticket of the chunk writes. | `internal/preflight/check_floor_test.go` (`TestCheckFloorGradesTheChunkUnion`). | The per-ticket filter reds the uncovered ticket today. |
| RC26 | 19 | A `check-floor` red names the chunk and the omitted element. | `internal/preflight/check_floor_test.go` (`TestCheckFloorNamesTheChunkElement`). | A detail that keeps the ticket name points at a ticket that owes nothing. |
| RC27 | 20 | A requirement of another chunk covers no element of this chunk. | `internal/preflight/check_floor_test.go` (`TestCheckFloorReadsOnlyTheOwnChunk`). | A union over the whole plan lets one chunk borrow another chunk's run. |
| RC28 | 20 | A `--changed` run satisfies no floor element. | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsAChangedRun`). | A union that admits every `bench test` form accepts a selector run. |
| RC29 | 21 | The preflight package holds no requirement-ownership derivation of its own. | Review-owned: `bench consumers reviewrecord.TicketRequirements` lists the only ownership derivation. No gate check can red it, because the chunk-union grade calls no ownership filter, so a kept `ticketRequirements` has no caller, and no gate phase reds an uncalled unexported function. | A kept `ticketRequirements` survives as a second rule that a later caller can revive. |
| RC30 | 22 | `bench record assignment` refuses while a ticket file differs from HEAD. | `internal/reviewrecord/recordcmd/assignment_test.go` (`TestAssignmentSourceRefusals`). | The repair order depends on this refusal, and a lost refusal records an assignment over an unreviewed widening. |
| RC31 | 23 | The implement phase names the charge preparation as the one build preflight of a dispatch. | Review-owned: the RC-C3 review reads .agents/commands/bench-implement-spec.md. | A guidance that keeps the chain for a dispatching plan commit runs two preflights. |
| RC32 | 24 | The implement phase keeps `bench commit --preflight-build <slug>` for a plan commit that no dispatch follows. | Review-owned: the RC-C3 review reads .agents/commands/bench-implement-spec.md. | A guidance that drops every chain leaves such a commit with no build preflight. |
| RC33 | 56 | Each edited guidance file stays inside its line budget. | Gate check: `bench test --check guidance-prose-budgets`. | A guidance edit that grows past its budget reds the check. |
| RC34 | 25 | `--prepare` runs the whole gate on the record-stripped HEAD tree and writes a green to the evidence store. | planned TestCheckpointPrepareRunsTheGateWithoutEvidence in internal/gate/checkpoint_prepare_test.go. | A prepare that runs the evidence pass refuses before the record commit. |
| RC35 | 25 | `--prepare` leaves the chunk gap closed. | planned TestCheckpointPrepareOpensNoChunkGap in internal/gate/checkpoint_prepare_test.go. | A prepare that opens the gap lets a successor start before the evidence grade. |
| RC36 | 26 | A chunk checkpoint in a second process after a record-only commit reuses the prepared green and runs no phase. | planned TestChunkCheckpointReusesThePreparedGreen in internal/gate/checkpoint_prepare_test.go. | Today the tip and the purpose in the oracle force a new run. |
| RC37 | 27 | A spec change after the prepared run makes the chunk checkpoint run the gate. | planned TestChunkCheckpointRunsAfterASpecChange in internal/gate/checkpoint_prepare_test.go. | A key that strips the spec with the record reuses a green for a changed plan. |
| RC38 | 27 | A source change other than the record after the prepared run makes the chunk checkpoint run the gate. | planned TestChunkCheckpointRunsAfterASourceChange in internal/gate/checkpoint_prepare_test.go. | A key without the tree reuses a green for changed code. |
| RC39 | 28 | A reused green with a stale verification entry refuses through the evidence pass. | `internal/gate/checkpoint_prepare_test.go` (`TestChunkCheckpointReuseKeepsTheEvidencePass`). | A reuse path that returns before the evidence pass accepts stale evidence. |
| RC40 | 29 | A record with a sentence over the prose bound reds the chunk checkpoint on a reused green. | `internal/gate/checkpoint_prepare_test.go` (`TestChunkCheckpointGradesTheRecordProseOnReuse`). | A record grade only inside the gate run never sees the record on a reuse. |
| RC41 | 29 | The same record reds the chunk checkpoint on a fresh run. | `internal/gate/checkpoint_prepare_test.go` (`TestChunkCheckpointGradesTheRecordProseOnARun`). | A record grade only on the reuse path misses the record, because the run grades the stripped tree. |
| RC42 | 30 | A record write and a record commit in the worktree during the prepared run leave its green. | planned TestCheckpointPrepareSurvivesARecordCommit in internal/gate/checkpoint_prepare_test.go. | A prepare that captures the working tree after the run returns exit 1, and a prepare that keeps the tip-moved guard refuses on the moved HEAD. |
| RC43 | 31 | The implement phase runs the coordinator probe before the prepared run and waits on the run's PID before the checkpoint. | Review-owned: the RC-C4B review reads .agents/commands/bench-implement-spec.md and the bounded repair policy. | A guidance that overlaps the probe and the run meets the `gate-held` refusal. |
| RC44 | 32 | The chunk checkpoint refuses a working tree that differs from HEAD with the `checkpoint-dirty-checkout` face. | `internal/gate/checkpoint_prepare_test.go` (`TestChunkCheckpointRefusesADirtyCheckout`). | A checkpoint that grades HEAD and ignores the working tree passes an uncommitted change silently. |
| RC45 | 33 | A complete checkpoint after a prepared chunk green runs the gate on its published tree. | planned TestCompleteCheckpointIgnoresAPreparedChunkGreen in internal/gate/complete_checkpoint_test.go. | A complete form that takes the chunk key reuses a green over the wrong tree. |
| RC46 | 33 | Prospective landing evidence still keys to the tree and the baseline runner identity. | `internal/gate/prospective_owner_test.go` (`TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity`). | A shared key change that reaches the landing reuses a green across a changed baseline. |
| RC47 | 34 | The reference, the profile, the bounded repair policy, and ADR 0002 state the chunk-checkpoint key and the prepared run. | Review-owned: the RC-C4B review reads each named document. | A stale sentence tells a fresh reader that record bytes still enter the key. |
| RC48 | 35 | A plan requirement whose command runs the gate accepts only a staged native excerpt with a matching digest. | `internal/reviewrecord/recordcmd/stage_test.go` (`TestGateRequirementKeepsTheNativeExcerpt`). | A reuse path that writes a synthesized entry passes with no native output. |
| RC49 | 25 | `bench gate --help` names `--prepare` in its usage line. | planned TestGateUsageNamesPrepare in internal/gate/checkpoint_prepare_test.go. | A flag missing from the usage line is a dead key for a fresh reader. |
| RC50 | 25 | `--prepare` with `--complete` is a usage error at exit 2. | planned TestCheckpointPrepareRefusesComplete in cmd/bench/gate_route_test.go. | A prepare of the complete form writes a chunk-shaped green under the complete purpose. |
| RC51 | 25 | `--prepare` without `--checkpoint` is a usage error at exit 2. | planned TestCheckpointPrepareRefusesNoCheckpoint in cmd/bench/gate_route_test.go. | A bare prepare runs an ordinary gate under a name that promises a checkpoint subject. |
| RC52 | 26 | A prepared green past the freshness window makes the chunk checkpoint run the gate. | planned TestChunkCheckpointRunsAfterThePreparedGreenExpires in internal/gate/checkpoint_prepare_test.go. | A reuse with no window check accepts an expired green. |
| RC53 | 36 | A changed system-suite test file selects `bench test --check system`. | planned TestChangedSelectionSelectsTheSystemSuite in internal/testreport/selection_test.go. | Today the selector refuses the path as outside every current package. |
| RC54 | 37 | A production change in a package that the system suite imports selects `bench test --check system`. | planned TestChangedSelectionReachesTheSystemSuiteByImport in internal/testreport/selection_test.go. | An index without the tagged package never reaches the suite through imports. |
| RC55 | 38 | The floor of a ticket that writes the system suite names `bench test --check system`. | planned TestWritesFloorNamesTheSystemCheck in internal/testreport/floor_test.go. | Today the floor adds no element for a system-suite path. |
| RC56 | 38 | The floor and the changed selection agree on a system-suite path. | `internal/testreport/floor_test.go` (`TestWritesFloorAgreesWithChangedSelection`). | A fix in one reader leaves the other reader with the old answer. |
| RC57 | 39 | The tagged operands of the index come from `gate.SystemSuite`. | Review-owned: `bench consumers gate.SystemSuite` lists the index as a consumer. RC53 derives its expected operands from `gate.SystemSuite`, so a second list with other values reds RC53. A second list with equal values reds no test, so review owns that case. | A second tag list in the selector drifts from the gate's system phase. |
| RC58 | 40 | An eligible test-only delta runs gofmt and vet over the whole root, `go test` of the changed packages, and root conformance, and no other phase. | planned TestDeltaCheckpointRunsTheDeltaSet in internal/gate/delta_checkpoint_test.go. | A red planted outside the delta set stays green only when the delta set alone runs. |
| RC59 | 41 | The delta package set equals the directories of the test-gap file list, a path with a space included. | planned TestDeltaCheckpointPackagesComeFromTheTestGapList in internal/gate/delta_checkpoint_test.go. | A set from another source adds or drops a package against the proven list. |
| RC60 | 41 | No checkpoint path calls the changed-package selector. | Review-owned: Go refuses the import, because `internal/testreport` imports `internal/gate`. | The compile fails on any direct call, and review reads any subprocess route. |
| RC61 | 42 | A test-only delta with no whole green at the reviewed tip runs the whole gate. | planned TestDeltaCheckpointNeedsAWholeBaseGreen in internal/gate/delta_checkpoint_test.go. | A delta that skips the base test grades a test-only delta over an ungraded base. |
| RC62 | 42 | A whole base green older than the freshness window still makes the delta eligible. | planned TestDeltaCheckpointReadsTheBaseWithoutExpiry in internal/gate/delta_checkpoint_test.go. | A base read through the expiring window drops every delta after 60 minutes. |
| RC63 | 43 | A delta with a non-test file runs the whole gate. | planned TestDeltaCheckpointRunsWholeForANonTestDelta in internal/gate/delta_checkpoint_test.go. | A delta that skips the test-gap proof grades a production change with the reduced set. |
| RC64 | 43 | A delta whose only non-test change is an assignment-only plan commit stays eligible. | planned TestDeltaCheckpointAcceptsAnAssignmentOnlyPlanCommit in internal/gate/delta_checkpoint_test.go. | A delta proof without the RC-C1 exemption sends every delegated repair to the whole gate. |
| RC65 | 44 | A delta that touches the system-suite package runs the whole gate. | planned TestDeltaCheckpointRunsWholeForTheSystemSuite in internal/gate/delta_checkpoint_test.go. | A delta without the tag rule leaves the edited system test unrun. |
| RC66 | 44 | A delta that touches a race-table package runs the whole gate. | planned TestDeltaCheckpointRunsWholeForARacePackage in internal/gate/delta_checkpoint_test.go. | A delta without the race rule leaves the race sentinels unrun. |
| RC67 | 45 | After a green delta, the landing inspection is not green and the evidence store holds no new entry. | planned TestDeltaGreenAuthorizesNoLanding in internal/gate/delta_checkpoint_test.go. | A delta written to the store under a whole key authorizes the landing. |
| RC68 | 45 | A complete checkpoint after a delta-only chunk runs the whole gate. | planned TestCompleteCheckpointIgnoresADeltaGreen in internal/gate/complete_checkpoint_test.go. | A complete form that reads the delta file passes on a reduced run. |
| RC69 | 45 | A green delta leaves `bench-last-gate` unchanged, and a later whole chunk checkpoint does not reuse the delta. | planned TestDeltaGreenLeavesTheGateCacheUntouched in internal/gate/delta_checkpoint_test.go. | A delta in the one-slot cache overwrites the whole green that `Inspect` and the composed-green reader report. |
| RC70 | 46 | The class table holds the delta class with a non-empty reuse refusal and no oracle field. | `internal/gate/verdict_registry_guard_test.go` (`TestVerdictRecordClassRegistryMatchesExpectation`). | A class row with an empty refusal or an oracle field changes the pinned table. |
| RC71 | 47 | `--fresh` on an eligible delta runs the whole gate. | planned TestDeltaCheckpointFreshRunsWhole in internal/gate/delta_checkpoint_test.go. | A fresh flag that the run decision ignores gives no way to buy the whole grade. |
| RC72 | 48 | A green delta checkpoint opens the chunk gap. | planned TestDeltaCheckpointOpensTheChunkGap in internal/gate/delta_checkpoint_test.go. | A delta path that returns before the gap write strands the successor chunk. |
| RC73 | 49 | A red delta run makes the checkpoint red and records no green. | planned TestDeltaCheckpointRedStaysRed in internal/gate/delta_checkpoint_test.go. | A fallback to the whole gate after a red delta hides a failure that the delta found. |
| RC74 | 40 | A prepared run at a repair commit, before the record commit refreezes the chunk, writes a delta green for an eligible delta, and the checkpoint after the record commit reuses it. | planned TestCheckpointPrepareTakesTheDelta in internal/gate/delta_checkpoint_test.go. | A derivation that reads the record's chunk tip finds the axis results at the chunk pair, so the prepared run always runs whole. |
| RC75 | 34 | The reference and the profile name the delta schedule and the delta class as the one reduced schedule and record. | Review-owned: the RC-C6B review reads .bench/BENCH-reference.md and projects/benchkit.md. | A kept sentence that forbids every reduced schedule contradicts the decided state. |
| RC76 | 50 | The first RC-C7 ticket records the digest cost, the selection cost, and the exec calls of one repair, and the chunk stops when the net saving is not positive. | Review-owned: the RC-C7 review reads the measurement record. | A build that skips the measurement ships a change that can cost more than it saves. |
| RC77 | 51 | An exec call of a `bench` child with a fresh `dist/bench` sets `BENCH_RUN_BINARY` to that path. | planned TestExecPassesTheFreshWorktreeBinary in internal/worktree/exec_test.go. | Today exec strips the variable from every child. |
| RC78 | 52 | An exec call of a `bench` child with a stale `dist/bench` omits the variable and prints `bench worktree build`. | planned TestExecOmitsAStaleWorktreeBinary in internal/worktree/exec_test.go. | An exec that passes the binary without the freshness check grades with a stale binary. |
| RC79 | 53 | An exec call of a non-Bench child runs no freshness check and keeps today's environment. | planned TestExecLeavesANonBenchChildUnchanged in internal/worktree/exec_test.go. | An exec that checks every child pays the digest cost on each shell command. |
| RC80 | 54 | A child that inherits a binary that became stale refuses and names `bench worktree build`. | planned TestInheritedStaleBinaryNamesTheRebuild in internal/runbinary/runbinary_test.go. | A refusal with no route leaves the operator with no fix after the race. |
| RC81 | 55 | The final reconciliation records the wall time of one test-only repair cycle against the 45-to-60-minute baseline. | Review-owned: the final reconciliation of the build. | A faster green run alone does not show that the cycle met its target. |
| RC82 | 43 | A prepared run at a repair commit whose record names an earlier chunk tip, and whose change from the reviewed tip holds a non-test file, runs the whole gate and writes no delta green. | planned TestCheckpointPrepareRunsWholeForANonTestRepair in internal/gate/delta_checkpoint_test.go. | A proof over the reviewed tip and the record's earlier chunk tip misses the new production file and writes a delta green for it. |
| RC83 | 42 | After a new review round moves the reviewed tip, the chunk checkpoint does not reuse a delta green whose `base_tree` is the old base. | planned TestDeltaGreenNeedsTheCurrentBase in internal/gate/delta_checkpoint_test.go. | A reuse that matches the subject tree alone accepts a delta over a base that no current review holds. |
| RC84 | 22 | The implement phase orders the repair plan edits: `bench record assignment`, then `bench spec widen` and any floor edit, then one `bench commit`. | Review-owned: the RC-C3 review reads .agents/commands/bench-implement-spec.md. | A guidance that widens first meets the dirty-ticket refusal of RC30 and splits the plan commit. |

Not covered: story 57 — the FT438 row retires in the spec-retire commit after the landing, outside the build.

Not covered: story 58 — the decision adds no behavior, and no code or guidance states a chunk ticket cap today.

### Edge inventory

The walk applies the profile's hostile-input checklist to the record, the gate, the selector, and exec.

- Absent versus empty: an entry with no `source_commit` and an entry with an empty value take the legacy rule (RC18, RC19). An absent assignment block and an empty assignment list both clear to the same fence in the RC-C1 comparison.
- A malformed or hostile commit text in the record refuses at the parser (RC23).
- A command whose own write changes a fact it reports: the record commit moves HEAD, and the key excludes the tip, so the tracked record commit is the tested case (RC36, RC42).
- State written by one process and read by a fresh one: the prepared run and the checkpoint are two processes (RC36, RC74).
- A path with a space: the delta package set keeps a test file path with a space (RC59).
- A dirty checkout: the chunk checkpoint refuses it (RC44), and the stage check keeps its dirty-path refusal (RC24).
- A symbolic link at the record path: the prose grade reads the tree blob and refuses a link entry through the existing grader classification.
- A tagged package: the system-suite package enters the index (RC53) and leaves the delta (RC65).
- Two concurrent gate runs in one worktree: the existing lock refuses the second run (`internal/gate/run_transaction.go:90-96`).
- A probe during the prepared run: the probe refuses with `gate-held`. The guidance orders the coordinator probe first, and an axis probe runs in its own worktree (RC43).
- An expired prepared green: the checkpoint runs the gate (RC52). An expired base green still serves the delta (RC62).
- A race between the exec check and the child: the child refuses (RC80).

**Won't handle:** A file-overlap skip of checks — FT445 puts it out of scope, and the whole checkpoint gate at the tip stays the surviving grade.

**Won't handle:** A general carry-forward engine for other evidence kinds — FT445 puts it out of scope, and the commit-bound entry of RC-C2 is the only carried evidence.

**Won't handle:** Landing reuse changes — FT314-gate-reuse owns them, and the landing's exact prospective key stays the surviving caller (RC46).

**Won't handle:** A delta checkpoint for a non-test delta — the whole gate is the surviving caller (RC63).

**Won't handle:** A delta checkpoint that adds the system or race run — the whole gate covers those packages (RC65, RC66).

**Won't handle:** A run-binary store that spans processes — the worktree's sealed `dist/bench` is the surviving holder (RC77).

**Won't handle:** An exclusion row for `reviews/` in the prose policy — the checkpoint record grade is the surviving caller (RC40, RC41).

**Won't handle:** Completion verification entries with a commit field — the complete checkpoint keeps its tree binding and stays the surviving caller.

**Won't handle:** A cap on the tickets of a chunk — the coherent-chunk rule of `.bench/BENCH.md` stays the surviving guidance.

**Won't handle:** A chunk checkpoint that waits for a held lock — the existing refusal and the guidance wait on the PID stay the surviving route (RC43).

## Out of scope

| Separate capability | Estimated work | Reason |
|---|---|---|
| FT314 landing reuse under the complete authorization key | 5 edits, 2 gate runs | Its own outcome owns the landing key. |
| A pre-publication build preflight for each plan commit (FT380-residual) | 4 edits, 2 gate runs | The reviewer narrowed its criterion to one preflight before each dispatch on 2026-10-24. |
| The standalone `bench record chunk` refusal after an assignment record (`docs/known-issues.md:13`) | 2 edits, 1 gate run | `bench record ticket` already applies the amendment step, and only the standalone form keeps the refusal. |
| A file-overlap skip of owed checks | 6 edits, 2 gate runs | FT445 excludes it. |

## Further notes

### Consultant decisions (2026-10-24)

The reviewer named Opus/high and Fable/medium as consultants, with Sonnet/high as the tie-breaker.

| Decision | Choice | Votes |
|---|---|---|
| Open point 1: a chunk ticket cap | No cap. A chunk stays one coherent outcome. | Opus 7, Fable 7, agreed |
| Open point 2: one build preflight for each dispatch | The charge preparation is the one run for an author or repair dispatch. | Opus 6, Fable 7, agreed |
| C1a: the chunk-checkpoint key | The tree without the record, plus the base oracle, with no purpose and no tip. | Opus 7, Fable 7, agreed |
| C1b: the record prose | The checkpoint grades it on every call, with no exclusion row. | Opus 7 (every call), Fable 7 (after a reuse). The coordinator chose every call, because the prepared run grades a tree without the record. |
| C1c: where the prepared run runs | A private prospective checkout, after the coordinator probe. | Opus 6, Fable 6, agreed |
| C2a: "conformance" in change 7 | Root conformance alone. | Opus 7, Fable 8, agreed |
| C2b: a tagged or race-table delta | The whole gate runs. | Opus 7, Fable 7, agreed |
| C2c: the delta record | A lane-shaped class in its own file, never in the store or the gate cache. | Split. Sonnet chose Fable's shape and the lane precedent's own file, 7. |
| C3: exec run-binary reuse | The worktree's sealed `dist/bench`, after a measurement. | Split. Sonnet chose Fable's option, 6. |
| T1: chunk order of changes 6 and 7 | Change 7 before change 6. | Split. Sonnet chose 5, 7, 6, 6. |

### Gap answers

- G1: The fast lane holds only the shared build-cache lock (`internal/gate/lane.go:116-127`). The gate lock opens only in the gate run (`internal/gate/run_transaction.go:85`). Only `bench probe` reads it (`internal/probe/probe.go:207`, `internal/probe/set.go:223`). A kit worktree commit runs the lane (`.bench/BENCH-reference.md:341-345`).
- G2: The new flag joins the gate flag constants and parser (`internal/gate/checkpoint.go:64-104`) and the validator (`internal/gate/checkpoint.go:47-62`). The existing faces in `internal/refusalroute/faces_gate.go` serve its refusals. The usage line is `CommandUsage` (`internal/gate/gate.go:220`), and `cmd/bench/help_inventory_test.go:128` holds a hand-typed copy of the gate row. The RC-C4B ticket changes both.
- G3: `review-record:chunks.verification.source_commit`. The grammar is in `internal/spec/spec.go:129-169`, and the record root is `internal/reviewrecord/record.go:88-97`.
- G4: The grader grades blob bytes through `prose.NewGraderFromIndex` and `GradeBytes` (`internal/gate/gate_prose.go:132-176`). `GradeNamed` reads the working tree (`internal/prose/subject.go:140-160`), so the checkpoint does not use it.
- G5: The base green is the evidence-store entry (`internal/gate/engine.go:143-152`) of the reviewed tip's stripped tree under the prospective base oracle. `FindReviewedPair` supplies the reviewed tip (`internal/reviewrecord/confirmation.go:74-126`).
- G6: No in-flight record lacks the field. The FT360-fixtures worktree holds the same two landed records as `main`, `reviews/ft297-build-conformance.md` and `reviews/needle-distinctness.md`. The legacy rule covers them and this build's own RC-C1 entries.
- G7: No. The census tests parse their own directory only (`internal/systemtest/parallel_census_test.go:354`, `internal/systemtest/serial_census_test.go:88`). The delta refusal set stays the system-suite package and the race-table packages.

### Flagged additions

Each item below goes beyond the literal decision source. Each cites its research basis, as the `docs/known-issues.md` row on unresearched spec seams requires.

1. The field `source_commit`: change 1 needs a commit binding, and no record field holds one (`internal/reviewrecord/record.go:32-38`; Q3 finding b).
2. The legacy entry rule: the two landed records (`reviews/ft297-build-conformance.md`, `reviews/needle-distinctness.md`) and this build's own early entries have no commit, and the parser refuses an unknown field (`internal/reviewrecord/parse.go:232`; G6).
3. The record-only descendant rule (RC12): the FT402 occurrence stranded four tickets on a record commit (`capture/build-log/landing-gate-long-poles.log:16`).
4. The `--prepare` flag: change 5 needs a gate start that skips the evidence pass and the tip-moved guard before the record commit (`internal/gate/checkpoint.go:249-285`, `internal/gate/evaluation.go:53-62`; Q4 finding a; consultant C1c).
5. The chunk checkpoint grades HEAD in a private checkout and refuses a dirty checkout: the prepared run and the checkpoint need equal oracles (consultant C1c; the complete form precedent at `internal/gate/complete_checkpoint.go:31-37`).
6. The record prose grade in the checkpoint pass: the prose check grades `reviews/` (`.bench/prose-exclusions:1-12`; coordinator fact 2).
7. A prepared run when a chunk freezes: change 7 needs a whole green at the reviewed tip, and today no run writes a green that a later key can find, because the checkpoint key frames the purpose and the tip (`internal/gate/checkpoint.go:277-285`; Q4 opus finding d).
8. The exported test-file list of `testgap`: the delta needs the proven list, and `Prove` returns only an error (`internal/testgap/testgap.go:24`).
9. The delta package set by directory: the gate cannot import `testreport` (`internal/testreport/floor.go:7`).
10. The exec reuse limit to a `bench` child: exec runs any argv directly (`internal/worktree/exec.go:118-130`), and the freshness check costs a `go list -deps` run (`internal/freshness/freshness_buildinputs.go:55`; Sonnet's narrowing of T3).
11. An axis probe in its own worktree when a prepared run starts at a freeze: `bench probe` refuses while the gate lock of its worktree is held (`internal/probe/set.go:223`).
12. `--fresh` forces the whole gate over an eligible delta (RC71): the existing flag already forces a run (`internal/gate/checkpoint.go:65,83-84`).
13. A red delta takes no fallback to the whole gate (RC73): a fallback would hide the red that the delta found (consultant C2b).
14. The stale-binary hint and the route text of an inherited refusal name `bench worktree build` (RC78, RC80): an inherited stale binary refuses and does not rebuild (`internal/runbinary/runbinary.go:155-168`; Sonnet F5).
15. The one reviewed-tip derivation with a target commit: at a repair commit, the record names the earlier freeze tip, so the record's chunk tip cannot serve as the delta operand (`internal/reviewrecord/confirmation.go:92-93`; spec review B2).
16. The record prose grade skips a HEAD tree with no `.bench/prose-exclusions`: no consumer repo receives that file (`.bench/consumer-payload.json`), and the prose package defines a tree with no policy as ungradeable (`internal/prose/subject.go:89-97`). The consultants decided this on 2026-10-24 (RC-C4 round r1).

### Contradictions flagged for reviewer veto

1. FT438 names "the review record" among the changes that still take a confirming round. `SourceDigest` already strips the record, so the clause has no effect. The spec follows the tree.
2. Change 4 and open point 2 conflict with the criterion `FT380-residual.before-publication` (`.bench/commitment.json:1361-1362`). That criterion had every plan commit run the build preflight before it publishes. The spec-stage commitment plan narrowed it to the one preflight run before a dispatch (sign-off decision 2).
3. Change 5 conflicts in intent with the staged spec `landing-test-efficiency`: "No verdict format, token, schema, or reuse policy changes" (`specs/landing-test-efficiency/spec.md:167`) and "Won't handle: Broader verdict reuse" (line 277). That spec also fences `internal/gate/engine.go` (line 309). Neither spec binds an outcome that runs first, but a later build of that spec must fold this one.
4. The profile states that reduced-scope records are never authored (`projects/benchkit.md:469-473`) and that no stripped-subject phase schedule exists (lines 387-389). FT445 change 7 reverses both for the delta checkpoint. ADR 0002 posture 5 says that the subject binds the working tree. The reviewer approved the changes in FT445, and RC47 and RC75 update the text.
5. `TestReviewCheckpointReuse` (`internal/gate/review_checkpoint_test.go:117-125`) pins that an ordinary green never satisfies a checkpoint. Under the new key, any whole green over the same stripped subject and oracle can serve, because the evidence pass still runs. The RC-C4 ticket rewrites that pin.
6. FT445 change 6 reverses a deliberate strip in exec (`internal/worktree/exec.go:217-229`, `TestExecChildOwnsRatherThanInheritsItsRunBinary`). RC76 can stop it.

### Reviewer decisions at sign-off (2026-10-24)

1. The reviewer approved this spec and its 13 tickets, and vetoed no flagged contradiction.
2. The spec-stage commitment plan narrows `FT380-residual.before-publication` to one build preflight before each dispatch. A plan commit that no dispatch follows keeps `--preflight-build`.
3. The same plan records `FT380-plan-commit.preflight-build` as superseded by FT445. It also states in `FT405.verb` that any ticket's record call can refreeze the chunk.
4. `internal/gate/` takes three new files over its directory budget: `delta_checkpoint.go`, `delta_checkpoint_test.go`, and `checkpoint_prepare_test.go`.
5. From RC-C3 on, this build records each ticket at its own commit, and no ticket restages. By then, the installed candidate holds the RC-C2B record rules.

### Bootstrap authority

The installed broker runs the landing. It builds the gate executable from the prospective tree and runs the whole exact gate (`projects/benchkit.md:49-56`). A broken chunk-checkpoint key in the candidate therefore cannot authorize the landing. During this build, each chunk checkpoint from RC-C4 on runs with `--fresh`, so the candidate does not reuse a green under its own new key. RC-C2 is a declared build, so `bench worktree build --install <target>` runs before its first dispatch (`.bench/BENCH-reference.md:398-433`).

### Source sentence to row

| Source clause (`roadmap/FT445.md`, `roadmap/FT438.md`) | Rows or disposition |
|---|---|
| Each ticket's evidence binds to its own commit tip, and the restage goes away. | RC8, RC9, RC11, RC12, RC13, RC14, RC15, RC16, RC23 |
| The whole checkpoint gate alone grades the chunk tip. | RC22 |
| The repair session runs each owed check once, at its own commit. | RC10, RC24 |
| A test-only repair closes through the test-only confirmation form. | RC1, RC5 |
| Each repair writes one plan commit, and only the last plan commit runs the preflight build. | RC30, RC31, RC32, RC84 |
| The gate starts when the repair commits, in parallel with the confirmation. | RC34, RC35, RC42, RC43, RC49, RC50, RC51 |
| The checkpoint reuses that green when the reuse key excludes the record path. | RC36, RC37, RC38, RC39, RC40, RC41, RC44, RC45, RC52 |
| Exec calls in one session reuse one run binary. | RC76, RC77, RC78, RC79, RC80 |
| A test-only delta after a whole-gate green takes a delta checkpoint. | RC58, RC61, RC62, RC63, RC64, RC74, RC82, RC83 |
| The delta checkpoint runs the changed packages, root-conformance, conformance, vet, and gofmt. | RC58, RC59 |
| The landing gate always runs whole. | RC46, RC67, RC68 |
| The repaired ticket reruns its own probe. Another ticket's probe carries forward. | RC10, RC17 |
| Every test-only judgment uses the test-gap proof. | RC1, RC63 |
| The checkpoint never uses the changed-package selector. | RC60 |
| The gate run stays as a native excerpt with its source digest. | RC48 |
| The selector lists tagged packages before any skip rule depends on it. | RC53, RC54, RC55, RC56, RC57 |
| The per-ticket check floor moves to a chunk-level floor. | RC25, RC26, RC27, RC28, RC29 |
| The FT438 exemption is the first chunk. | RC-C1 order |
| The FT438 row retires when FT445 lands. | Story 57, not covered; the spec-retire commit |
| A file-overlap skip and a general carry-forward engine stay out of scope. | Won't handle lines |
| Open point: the chunk cap of four tickets. | Story 58, not covered; consultant decision |
| Open point: one preflight build for an author dispatch. | RC31 |
| FT438: the exemption covers the assignment block only. | RC1, RC2 |
| FT438: a plan-fence change, a `Writes:` widening, the review record, or a seam cell still takes a confirming round. | RC2, RC3, RC4; the record clause is contradiction 1 |
| Occurrence: a non-last ticket's repair could not refreeze the chunk. | RC11 |
| Occurrence: a record commit before every ticket restaged stranded four tickets. | RC12 |
| Occurrence: each dispatch ran two preflight builds. | RC31 |

### Pre-review proof checklist

- `Cited symbols`: each cited symbol resolves at `a6243061d`: `FindReviewedPair`, `SourceDigest`, `TreeWithoutFile`, `ReadTreeFile`, `Prove`, `TicketRequirements`, `ticketRequirements`, `RecordAssignment`, `ExecutionInProgress`, `executeTreeWithOwner`, `openProspectiveArtifacts`, `inspectEvidenceWindowed`, `SystemSuite`, `raceTests`, `NewGraderFromIndex`, `GradeBytes`, `VerifyAt`, `Inherit`.
- `Import edges`: `gate` already imports `reviewrecord` (`internal/gate/checkpoint.go:18`). `testreport` already imports `gate` (`internal/testreport/floor.go`). The new `gate` to `testgap` edge has no cycle, because `testgap` imports neither `gate` nor `testreport`.
- `Source-row clauses and occurrences`: the source-sentence table above.
- `Promised field labels`: `source_commit`, `schema`, `delta`, `base_tree`, `tree`, `recorded_at`, `run_binary`, `bench-last-delta`, `--prepare`, `gate: green (delta checkpoint)`.
- `Changed-function callers`: `reviewrecord.TicketRequirements` has two callers (`internal/reviewrecord/recordcmd/ticket.go:48`, `internal/reviewrecord/ticket.go:91`). `preflight.ticketRequirements` has one (`internal/preflight/check_floor.go:136`). `testgap.Prove` has one (`internal/reviewrecord/confirmation.go:121`). `FindReviewedPair` has two (`internal/reviewrecord/write.go:272`, `internal/reviewrecord/confirmation.go:136`).
- `Copy survival`: RC29 for the ownership derivation and RC57 for the tag list.
- `Rendered-shape readers`: the restage refusal text "the staged run is at" is pinned at `internal/reviewrecord/recordcmd/ticket_test.go:146,155`. The RC-C2B ticket rewrites both.
- `Pin operators`: RC70 uses the exact class-table equality of `TestVerdictRecordClassRegistryMatchesExpectation`.
- `Entry reads`: none.
- `Derived expectations`: none. Each planned row states an independent expectation.
- `Consolidated rules`: the commit predicate replaces the digest equality at four sites (`coverage.go:76-80`, `change.go:115-125`, `recordcmd/ticket.go:120-137`, `ownership/evidence.go:63,208`). Each site maps to RC8, RC9, RC20, and RC21.
- `Quantified obligations`: every chunk verification entry, every chunk-checkpoint call, and every delta checkpoint are covered by RC8 to RC22, RC36 to RC44, RC58 to RC74, RC82, and RC83.
- `Workflow-step writes`: the record commit changes only `reviews/<slug>.md`, which `SourceDigest` strips, so the prepared key holds. A plan commit changes the spec, which the key keeps, so RC37 forces a run.

`bench anchors` ran for `projects/benchkit.md`, `docs/known-issues.md`, and the bounded repair policy. The research round read the anchors of `.agents/commands/bench-implement-spec.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`, and `.agents/commands/bench-review-implementation.md`. No anchor pins a rule sentence that this spec edits. New wording avoids the Forbid needles of `internal/anchors/registry_chunk_chain.go:18-24` and `internal/anchors/registry_data.go:156`. The profile forbids gate-output sentences (`bench anchors projects/benchkit.md`), so the delta output account goes to `.bench/BENCH-reference.md` only.

### Guidance edits and budgets

`.agents/commands/bench-implement-spec.md` is at 95 of 95 lines (`projects/benchkit.md:500`). The lane grades `guidance-prose-budgets` at each commit, so each ticket keeps the file at 95 lines or fewer at its own commit (RC33). Each paragraph of the file is one line, so a removed sentence frees no line. Each ticket that edits the file keeps its paragraph count and edits inside paragraphs.

Three anchored command sentences stay verbatim: the charge command at line 21, the post-ticket preflight at line 52, and the line 68 sentence that `require-in-section Land` pins: "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`." The other sentences of line 68 change. The line 62 clause "no dispatch starts until its `commit-chain` line reads green" also changes under RC31.

`.bench/BENCH.md` (190 of 192 lines) and `craft-tickets` (97 of 100) need no edit. "The last ticket's author freezes" stays true, because the last ticket's record call still freezes the chunk before the review. `.bench/BENCH-reference.md` and the bounded repair policy have no line budget, so they carry the command detail.

### Caller dispositions

| Caller or test | Disposition |
|---|---|
| `TestTicketRefusals` cases "stage from another commit" and "later stage from another commit" | Rewrite to the off-line refusal (RC9, RC13). |
| `TestTicketRefreezeKeepsTheBase`, `TestTicketEarlierTicketKeepsTheEntry` | Keep the base rule. Rewrite the earlier-ticket case to RC11. |
| `TestReviewCheckpointConfirmationRefusals` case "verification at the reviewed digest" | Invert to RC17. |
| `TestReviewCheckpointReuse` first case | Rewrite under contradiction 5. Keep its record-mutation case as RC39. |
| `TestCheckFloorReadsOnlyTheOwnTicket` | Replace with RC25. |
| `TestExecChildOwnsRatherThanInheritsItsRunBinary`, `TestExecChildDropsWrapperRouting` | Rewrite to RC77 and RC79. |
| `recordtest` `Refreeze` | Stop re-retaining entries. |
| `internal/worktree/delegated_integration_test.go:203-209` relabel case | Keep. It must still refuse under RC15. |
| `TestChunkCheckpointGradesTheCheckoutTree` | Keep. Its committed fixture passes the dirty-checkout rule. |
