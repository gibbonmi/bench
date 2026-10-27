# Repair carry

Status: staged

Roadmap: FT448

Decision source: named reviewed artifact `roadmap/FT448.md`, from drain d-0700c4586eaa with the reviewer decision of 2026-10-25, and the reviewer's late answer of 2026-10-27 that only a test-only repair delta carries.

Verification log: 3 iteration(s) to accept — the spec round took 2: iteration 1 found an unrouted carried-entry writer, an undefined assignment-only delta, and an unswept operating-guide clause. The ticket round accepted in 1.

## Problem

A repair session must run every requirement that its ticket owes, at the repair commit.
The `bench record ticket` verb refuses the call when one owed requirement has no staged excerpt.
Thus a test-only fix reruns the whole system suite, which takes about five minutes.
During the FT445 build on 2026-10-24, each repair reran every owed requirement for a doc-only or test-only fix.

## Solution

The `bench record ticket` verb carries a requirement forward when the repair delta cannot reach it.
The verb writes one carried entry for that requirement, and the repair session runs only the reached requirements.
The changed-package selector of the delta checkpoint decides the reach.
That selector proves reach only for a test-only delta, so another delta reruns every owed requirement.
When a requirement cannot carry, the refusal names it and the reason, so the session knows which checks to run.

## User stories

Line: opus / high.
Implementation-line reason: the carry decision is the hardest material chunk, because it composes four existing owners and a wrong answer keeps a stale green; the seams exist and the fixtures are deterministic, so the mid tier at high effort is sufficient.
Harder chunks: CF-C2.

### Carry an unreached requirement

1. As a repair session, I want the record verb to accept an owed requirement that has no stage when the repair delta cannot reach it, so that I do not rerun a check whose inputs did not change.
2. As a repair session, I want the system-suite requirement to carry across a test-only delta outside the suite, so that the five-minute suite runs only when its inputs change.
3. As a repair session, I want a package requirement to carry when the delta changes only the test files of other packages, so that a focused fix reruns only its own package.
4. As a repair session, I want an assignment-only plan commit in the delta to keep the carry of a package requirement, so that the repair dispatch commit does not force a full rerun.
5. As a repair session, I want an unchanged source to reach no requirement, so that a second record call at the same chunk source carries every passing entry.

### Record the carried entry

6. As a repair session, I want the carried entry to name me as its performer, so that the chunk checkpoint grades it against the current author of the ticket.
7. As a reviewer, I want the carried entry to keep the source commit, the source digest, and the excerpt of the entry that it carries, so that the entry states where the run occurred.
8. As a reviewer, I want the reference of the carried entry to name the evidence ID that it carries, so that I can trace the original run.
9. As a repair session, I want the ticket output to list the carried requirements, so that I know which checks I did not rerun.
10. As a repair session, I want the ticket output to show `none` when the call carries nothing, so that an empty carry is explicit.
11. As a repair session, I want a stage to win over a carry, so that a run that I made is the recorded evidence.
12. As a repair session, I want one call to write the staged entries and the carried entries in owed order, so that one record write holds the whole ticket.
13. As an agent, I want the `bench record ticket` help line to state the carry, so that I find the behavior in the help.
43. As a repair session, I want each carried entry to pass the chunk verification grade, so that a carry never turns the chunk checkpoint red.

### Keep a reached requirement

14. As a reviewer, I want a requirement whose package the delta reaches to need a fresh stage, so that a stale green never stands for changed tests.
15. As a reviewer, I want a package pattern that covers a changed package to count as reached, so that a parent pattern does not hide the change.
16. As a reviewer, I want a delta that changes a non-test path to rerun every owed requirement, so that a doc, production, or plan change keeps its full verification.
17. As a reviewer, I want a delta that changes a ticket file of the spec to rerun every owed requirement, so that a `bench spec widen` commit keeps its full verification.
18. As a reviewer, I want a delta that changes a system-suite test file to rerun every owed requirement, so that the whole-only system phase keeps its rule.
19. As a reviewer, I want a delta that changes a race-table test file to rerun every owed requirement, so that the whole-only race phase keeps its rule.
20. As a reviewer, I want each conformance named check to count as reached, so that the checks that read every test file rerun.
21. As a reviewer, I want the prose named check to count as reached, so that no rule assumes what the prose grader reads.
22. As a reviewer, I want a package run that covers the conformance package to count as reached, so that the conformance rule of the delta set holds for a package requirement.
23. As a reviewer, I want a command that is not one simple `bench test` call to count as reached, so that an unparsed command never carries.
24. As a reviewer, I want a listing form to count as reached, so that a selection that runs no test never carries.
25. As a reviewer, I want a package expression that is not in canonical form to count as reached, so that a second spelling cannot hide a changed package.
26. As a reviewer, I want a quoted package expression that holds a space to match its changed package, so that a hostile directory name keeps the reach.
44. As a reviewer, I want a changed source with no changed test file to reach each conformance check, so that a plan-only change still reruns the checks that read the plan.

### Carry only a passing current entry

27. As a reviewer, I want a requirement with no prior entry to need a stage, so that a first record still runs every check.
28. As a reviewer, I want a requirement whose last entry failed to need a stage, so that a failed result never carries.
29. As a reviewer, I want the last entry of the requirement to be the carry candidate, so that the carry reads the entry that the checkpoint grades.
30. As a reviewer, I want an entry with no source commit to stay ineligible, so that the delta always has a defined start.
31. As a reviewer, I want an entry that is off the chunk line after the freeze to stay ineligible, so that a moved chunk base reruns its checks.
32. As a reviewer, I want a probe requirement to carry only a probe that bit and restored, so that the probe obligation survives the carry.
33. As a reviewer, I want a probe requirement to carry its prior probe result, so that the checkpoint grades the planned probe on the carried entry.
34. As a maintainer, I want the carry and the checkpoint to read one pass predicate, so that the two never disagree about a passing entry.

### Refuse what does not carry

35. As a repair session, I want the refusal to name each unstaged requirement that cannot carry and its reason, so that I run exactly those checks.
36. As a repair session, I want the refusal to omit each requirement that can carry, so that I do not run a check that the carry covers.
37. As a repair session, I want a refusal to leave the record bytes and every stage unchanged, so that I can stage the named checks and call again.

### One source for each rule

38. As a maintainer, I want the build preflight and the record verb to parse a plan command with one parser, so that the two never disagree about a selection.
39. As a maintainer, I want the delta checkpoint and the record verb to read one delta selector, so that the carry rule follows the checkpoint rule.

### State the rule

40. As an orchestrator, I want the implement phase to state that a repair session runs only the reached checks, so that a repair charge does not demand a full rerun.
41. As a reviewer, I want each retired full-rerun sentence forbidden, so that it does not return.
42. As a user, I want the changelog to state the carry, so that I know that the verb changed.
45. As an orchestrator, I want the operating guide to state that a repair session reruns only the reached checks, so that every session reads one repair rule.

## Implementation decisions

### Closed decisions

- The reviewer decided on 2026-10-27 that only a test-only repair delta carries. A doc-only delta, a production delta, a `bench spec widen` delta, and a fold delta rerun every owed requirement. Package tests, conformance checks, and the system suite read repository documents at run time, so no package rule proves that a doc change is safe.
- The FT448 row names the changed-package selector of the delta checkpoint as the reach rule. This spec adds no second selector.

### The carry decision

The carry decision is one deep module, `internal/reviewrecord/carry`. Its one entry takes the root, the spec, the chunk entry after the freeze, and one owed requirement. It answers the carry candidate, or the reason that the requirement cannot carry.

The decision applies these steps in order. The first step that fails gives the reason.

1. The candidate is the last entry of the requirement in the chunk entry. The candidate must pass the shared pass predicate, hold a source commit, and be current for the chunk entry under `EntryCurrent`. Otherwise the reason is `no current passing entry`.
2. When the source digest of the candidate equals the source digest of the chunk entry, the source did not change. Then the requirement carries.
3. The test-only file list of the review record grades the delta from the candidate digest to the chunk digest, with its exemption for the assignment block of the spec. If that list refuses the delta, the reason is `the delta is not test-only`.
4. The gate answers the delta reach of the file list: the package operand of each changed test file, and the conformance package operand. An empty file list gives only the conformance operand. A changed test file in the system suite or in a race-table package makes the delta whole-only, and the reason is `the delta changes a whole-only test file`.
5. The requirement command is reached unless one of two forms holds. The reason is then `the delta reaches the command`.
   - The shared parser reads one simple `bench test` call in the check form, with the target `system`, that runs a test. A `--run` filter keeps this form.
   - The shared parser reads one simple `bench test` call in the package form that runs a test and is not `--changed`. Each package expression has canonical form, and no expression covers a reach operand.

A canonical package expression is `./...`, `./.`, `./<path>`, or `./<path>/...`. Here `<path>` is a repository path under the rule of the ticket `Writes:` paths, and it holds no `...`.

The carry adds the conformance package operand to every changed source. The conformance checks read the tree, and the delta set of the gate runs root conformance wherever the kit runs it. The carry rule is thus at least as strict as the gate. This rule also makes each named check other than `system` a reached check.

The shared pass predicate holds when the state is `completed`, the outcome is `pass`, and the exit code is present and zero. The command must equal the requirement command. For a probe requirement, the probe mutation must equal the planned probe, and `ProbeBites` must report a bite and a passed restore. The predicate does not read the performer or the role. The chunk verification grade adds the performer and the role of the plan, and it keeps its refusal messages and their order.

### Shared owners

- The plan-command parser and the package-cover predicate move from the build preflight into `internal/testreport`. The check floor and the carry decision both call them.
- The gate exports the delta reach beside the delta set. The delta checkpoint and the carry decision read the same whole-only rule and the same package operands.
- The review record exports the shared pass predicate and the carry candidate. The chunk verification grade and the carry candidate read the same predicate.

### The ticket form

The ticket form keeps its order of steps. A missing stage no longer refuses before the write. Instead, the verification step decides each unstaged owed requirement against the chunk entry after the freeze. A stage always wins over a carry.

When the call tip does not descend from the chunk tip, the chunk entry stays as it is. The carry then compares the candidate with the chunk tip, which is the source that the checkpoint grades.

The verb writes each carried entry through the one verification writer of the review record, as it writes a staged entry. The writer then derives the role, the command, the state, and the outcome, and it refuses a duplicate evidence ID. The call to the writer holds these values:

- The source is the source commit of the candidate. The binding of the writer keeps an ancestor commit on the chunk line as it is.
- The exit code and the probe result are those of the candidate.
- The evidence ID, the performer, the model, and the effort come from the call, as for a staged entry.
- The native reference is `carried:` followed by the evidence ID of the candidate. Its excerpt and its digest are those of the candidate.

The record schema does not change. A carried probe entry keeps the excerpt and the digest of the candidate, so the evidence reader of the build preflight still finds the test subject in the latest entry.

The ticket output row gains a `carried` cell after the `verifications` cell. The cell lists the carried requirement IDs in owed order, or `none`.

When an unstaged requirement cannot carry, the verb refuses through the excerpt form and writes nothing. The refusal keeps its current text, `these owed requirements have no staged excerpt: <ids>`, and lists only the requirements that cannot carry. The text then adds `; carry: <id> (<reason>)` for each of them, in owed order. A repair session can call the verb before it runs any check, and the refusal names the checks to run.

A refusal of the amendment step or the chunk step now comes before the missing-stage refusal, because the carry decision needs the frozen chunk entry. This spec makes no order promise between those refusals, the stale-stage refusal, and the carry refusal.

### Guidance

The Land section of the implement phase replaces one sentence. The old sentence is ``A repair session runs each check that its ticket owes once at its repair commit, and its `bench record ticket` call refreezes the chunk at that commit.`` The new sentence states that a repair session runs each check that its repair delta reaches at its repair commit. It also states that the `bench record ticket` call carries the rest and refreezes the chunk at that commit. The edit keeps the line count of that file.

The operating guide replaces the clause `that session reruns the ticket's verification with current repair coverage`. The new clause states that the session reruns each ticket check that its repair delta reaches, with current repair coverage. The edit keeps the line count of that file.

Both files are prose subjects, with a bound of 25 words for each sentence and 6 sentences for each paragraph. The implement-phase paragraph already holds 6 sentences, so its replacement stays one sentence of 25 words or fewer. The operating-guide clause can become a second sentence, because its paragraph then holds 6 sentences.

Two Forbid anchors refuse the return of the retired bytes. The changelog states the carry under one new heading.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| CF-C1 / 01-share-the-plan-command-parser.md, 02-share-the-pass-predicate.md | The check floor parses plan commands through one shared parser and cover predicate, and the chunk grade and the carry candidate read one pass predicate | CF1, CF2, CF3, CF30 | check floor tests in `internal/preflight` and the pass-predicate agreement test | no |
| CF-C2 / 03-decide-the-carry.md | One deep module decides whether an owed requirement carries and gives the reason when it cannot | CF4, CF5, CF6, CF7, CF8, CF9, CF10, CF11, CF12, CF13, CF14, CF15, CF16, CF17, CF18, CF19, CF20, CF21, CF22, CF23, CF24, CF25, CF26, CF27, CF28, CF29, CF31, CF32, CF51 | carry decision tests, carry reach tests, and the delta checkpoint tests | yes |
| CF-C3 / 04-record-carried-entries.md, 05-state-the-repair-carry.md | The ticket form writes carried entries that the checkpoint accepts, reports them, and names each check to run, and the guidance states the rule | CF33, CF34, CF35, CF36, CF37, CF38, CF39, CF40, CF41, CF42, CF43, CF44, CF45, CF46, CF47, CF48, CF49, CF50 | ticket form carry tests, the help inventory test, and the chunk chain anchor test | no |

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft448-repair-carry-full-20261027","orchestrator_session":"claude:rc-orchestrator-20261027","author_limit":3,"assignments":{"01-share-the-plan-command-parser.md":[{"session":"claude:rc-t01-author-20261027","assignment":"a2a41b7c586b6bd6d96de31844ff5b7b","model":"opus","effort":"high","source":"ee17bd346dc30724c596cbf92da2f84cbd000db0","native_ref":"claude-agent:bench-writer-high:rc-t01-author-20261027"}],"02-share-the-pass-predicate.md":[{"session":"claude:rc-t02-author-20261027","assignment":"a2a41b7c586b6bd6d96de31844ff5b7b","model":"opus","effort":"high","source":"8986e0c8cfa9bc156954e6b1c469cee0734ea353","native_ref":"claude-agent:bench-writer-high:rc-t02-author-20261027"}],"03-decide-the-carry.md":[{"session":"claude:rc-t03-author-20261027","assignment":"a2a41b7c586b6bd6d96de31844ff5b7b","model":"opus","effort":"high","source":"40f4f7130504226734876e16b3e3f3910e66b7c0","native_ref":"claude-agent:bench-writer-high:rc-t03-author-20261027"},{"session":"claude:rc-t03-repair1-20261027","assignment":"a2a41b7c586b6bd6d96de31844ff5b7b","model":"opus","effort":"high","source":"2713e9bb87c6557cba127be0cd58968a1f32e7c0","native_ref":"claude-agent:bench-writer-high:rc-t03-repair1-20261027","predecessor":"claude:rc-t03-author-20261027","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:rc-t03-author-20261027","preserved":"22dfdadc07c0f973b0e19b9737ce238f565cdbd6"}],"04-record-carried-entries.md":[{"session":"claude:rc-t04-author-20261027","assignment":"a2a41b7c586b6bd6d96de31844ff5b7b","model":"opus","effort":"high","source":"1b6d5fe06919c9100c337116d405e9910f9b3ad9","native_ref":"claude-agent:bench-writer-high:rc-t04-author-20261027"}],"05-state-the-repair-carry.md":[{"session":"claude:rc-t05-author-20261027","assignment":"a2a41b7c586b6bd6d96de31844ff5b7b","model":"opus","effort":"high","source":"d109312d6985bcabad1456edf92e603f88f6f356","native_ref":"claude-agent:bench-writer-high:rc-t05-author-20261027"}]}},"chunks":[{"id":"CF-C1","tickets":["01-share-the-plan-command-parser.md","02-share-the-pass-predicate.md"],"verification":[{"id":"testreport","command":"bench test --package ./internal/testreport","ticket":"01-share-the-plan-command-parser.md"},{"id":"preflight","command":"bench test --package ./internal/preflight","ticket":"01-share-the-plan-command-parser.md"},{"id":"record","command":"bench test --package ./internal/reviewrecord","ticket":"02-share-the-pass-predicate.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"01-share-the-plan-command-parser.md"},{"id":"conformance-02","command":"bench test --check root-conformance","ticket":"02-share-the-pass-predicate.md"}]},{"id":"CF-C2","tickets":["03-decide-the-carry.md"],"verification":[{"id":"carry","command":"bench test --package ./internal/reviewrecord/carry","ticket":"03-decide-the-carry.md"},{"id":"fixtures","command":"bench test --package ./internal/reviewrecord/recordtest","ticket":"03-decide-the-carry.md"},{"id":"gate","command":"bench test --package ./internal/gate","ticket":"03-decide-the-carry.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"03-decide-the-carry.md"}]},{"id":"CF-C3","tickets":["04-record-carried-entries.md","05-state-the-repair-carry.md"],"verification":[{"id":"recordcmd","command":"bench test --package ./internal/reviewrecord/recordcmd","ticket":"04-record-carried-entries.md"},{"id":"fixtures","command":"bench test --package ./internal/reviewrecord/recordtest","ticket":"04-record-carried-entries.md"},{"id":"help","command":"bench test --package ./cmd/bench","ticket":"04-record-carried-entries.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"04-record-carried-entries.md"},{"id":"anchors","command":"bench test --package ./internal/anchors","ticket":"05-state-the-repair-carry.md"},{"id":"prose","command":"bench test --check prose","ticket":"05-state-the-repair-carry.md"},{"id":"conformance-05","command":"bench test --check root-conformance","ticket":"05-state-the-repair-carry.md"}]}],"final_verification":[{"id":"whole-project","command":"bench gate"},{"id":"coverage","command":"bench coverage --check specs/repair-carry/spec.md"}]}
```

## Testing decisions

- A good test drives the carry decision or the ticket form with a real Git fixture. It commits a real delta and reads the real record. No test stubs the delta, the plan, or the record.
- The carry decision tests use the `recordtest` fixtures. Each test commits a prior entry, then a delta, and then reads the decision for one requirement. A legacy entry with no source commit comes from a hand-built record that the fixture saves.
- The ticket form tests drive `recordcmd.Command` over a delegated fixture. They combine the repaired-ticket precedent `repairedEarlierTicket` with a repair dispatch through `recordtest.Reassign`.
- The pass-predicate agreement test follows the precedent of `TestProbeChecksAgree`. It grades each entry variant through the carry candidate and through the chunk verification grade, with the plan performer, and it expects one verdict.
- The grade row uses a one-chunk delegated fixture with `Reassign("1.md")`, because the grade stops at the first chunk that lacks reviews. The test commits the record, because `rr.Check` reads the record from the commit tree. The grade runs the verification step of a chunk before its review step. So a review-form repair of the carried chunk itself proves that its verification step passed.
- The gate observes the feature at the chunk checkpoint, which grades each carried entry as an ordinary chunk verification entry.

### Seam diagram

    trigger: a repair session calls `bench record ticket` after its repair commit
        │
        ▼
    owed requirements + stages + record  ──▶  [ ticket form ]  ──▶  record write or refusal
                                                   │
                                                   ▼ one unstaged requirement + frozen chunk entry
                                             [ carry decision ]  ──▶  candidate, or reason
                                                   ◀ tests attach here: fixture commits a prior entry and a delta, the test reads the decision
        ◀ tests attach here: `recordcmd.Command` drives the verb, the test reads the record bytes, the output row, the stages, and the chunk grade

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CF1 | 38 | The check floor refuses a compound plan command through the shared parser | `internal/preflight/check_floor_test.go` (`TestCheckFloorRejectsACompoundCommand`) | A moved parser that drops the one-simple-command rule lets the floor accept a compound line |
| CF2 | 38 | The check floor keeps the segment boundary of a package pattern through the shared cover predicate | `internal/preflight/check_floor_test.go` (`TestCheckFloorKeepsTheSegmentBoundary`) | A moved cover predicate that matches a prefix accepts a sibling package |
| CF3 | 38 | The build preflight declares no plan-command parser and no package-cover predicate of its own | review-owned: Standards runs `rg -n "func planSelection\|func packageCovers" internal/preflight` and expects no hit | A surviving copy is a second source for the selection rule |
| CF4 | 2 | A test-only delta in `./internal/b` carries the `bench test --check system` requirement | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that treats every check as reached refuses the system requirement |
| CF5 | 3 | A test-only delta in `./internal/b` carries the requirement `bench test --package ./internal/a` | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that treats every package run as reached refuses the requirement |
| CF6 | 4 | A delta that holds an assignment-only plan commit and a test change in `./internal/b` carries the requirement `bench test --package ./internal/a` | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that grades the raw tree delta refuses the repair dispatch commit |
| CF7 | 5 | A candidate whose source digest equals the chunk digest carries the `bench test --check root-conformance` requirement | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that adds the conformance operand before it compares the digests refuses an unchanged source |
| CF8 | 14 | A test change in `./internal/a` refuses the requirement `bench test --package ./internal/a` with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision that carries every package run keeps a stale green |
| CF9 | 15 | A test change in `./internal/a` refuses the requirement `bench test --package ./internal/...` with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | An equality-only cover test misses the parent pattern |
| CF10 | 16 | A change to `README.md` refuses each requirement with the reason `the delta is not test-only` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision that reads only Go packages carries across a doc change |
| CF11 | 17 | A change to a ticket file of the spec refuses each requirement with the reason `the delta is not test-only` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | An exemption widened to the spec folder carries across a widen commit |
| CF12 | 18 | A change to a system-suite test file refuses the system requirement with the reason `the delta changes a whole-only test file` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision without the whole-only rule carries a changed suite |
| CF13 | 19 | A change to a race-table test file refuses each requirement with the reason `the delta changes a whole-only test file` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision that copies only the system-suite member of the rule misses the race table |
| CF14 | 20 | A test-only delta in `./internal/b` refuses the `bench test --check root-conformance` requirement with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision that carries every named check skips the conformance rerun |
| CF15 | 21 | A test-only delta in `./internal/b` refuses the `bench test --check prose` requirement with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision that carries each non-conformance check skips the prose rerun |
| CF16 | 22 | A test-only delta in `./internal/b` refuses the requirement `bench test --package ./internal/conformance` with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A reach without the conformance operand carries the conformance package |
| CF17 | 23 | A test-only delta in `./internal/b` refuses the requirement `go test ./internal/a` with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision that reads the package words of a foreign command carries it |
| CF18 | 23 | A test-only delta in `./internal/b` refuses the requirement `bench test --package ./internal/a && true` with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision that splits the line on words carries a compound line |
| CF19 | 44 | A delta that holds only an assignment-only plan commit refuses the `bench test --check root-conformance` requirement with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision that reads an empty file list as an unchanged source skips the conformance rerun |
| CF20 | 24 | A test-only delta in `./internal/b` refuses the requirement `bench test --check system --fixtures` with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision that reads only the check target carries a listing of the system check |
| CF21 | 25 | A test-only delta in `./internal/b` refuses the requirement `bench test --package ./internal/a/` with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A decision without the canonical-form rule carries a second spelling |
| CF22 | 2 | A test-only delta in `./internal/b` carries the requirement `bench test --check system --run TestA` | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that requires an unfiltered system run refuses a filtered one |
| CF23 | 27 | A requirement with no prior entry refuses with the reason `no current passing entry` | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that treats an absent entry as carried writes evidence for a run that never occurred |
| CF24 | 28 | A requirement whose last entry failed refuses with the reason `no current passing entry` | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that ignores the outcome carries a failure |
| CF25 | 29 | A requirement whose earlier entry passed and whose last entry failed refuses with the reason `no current passing entry` | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that searches for any passing entry carries a result that the checkpoint does not grade |
| CF26 | 30 | A requirement whose last entry is a hand-built legacy entry with no source commit refuses with the reason `no current passing entry` | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that reads an empty source commit computes a delta from no start |
| CF27 | 31 | A requirement whose last entry is off the chunk line after the freeze refuses with the reason `no current passing entry` | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision without the currency test carries an entry that the checkpoint refuses |
| CF28 | 32 | A probe requirement whose last probe did not bite refuses with the reason `no current passing entry` | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A decision that ignores the probe carries a failed probe obligation |
| CF29 | 33 | A carried probe requirement answers the candidate with its probe outcome, exit, and restore | `internal/reviewrecord/carry/carry_test.go` (`TestCarryDecision`) | A candidate without its probe makes the checkpoint refuse the carried entry |
| CF30 | 34 | The carry candidate and the chunk verification grade give one verdict for each entry variant under the plan performer | `internal/reviewrecord/carry_candidate_test.go` (`TestCarryAndCheckpointAgree`) | A second pass predicate drifts from the grade on a failed, empty, command, or probe variant |
| CF31 | 39 | The delta checkpoint still runs the whole gate for a system-suite test change through the shared selector | `internal/gate/delta_checkpoint_test.go` (`TestDeltaCheckpointRunsWholeForTheSystemSuite`) | An export that breaks the whole-only rule turns this test red |
| CF32 | 39 | No package outside the gate states the whole-only rule of the delta set | review-owned: Standards runs `rg -n "SystemSuite\(\|raceTests" internal/reviewrecord` and expects no hit | A copied rule is a second source for the delta selector |
| CF33 | 1 | Under a repair dispatch, a call that stages only the reached requirement exits 0 and writes a carried entry for the unstaged system requirement | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketCarriesTheUnreachedRequirement`) | A verb that still demands every stage refuses the call |
| CF34 | 6 | The carried entry names the repair session that the plan names as the author of the ticket | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketCarriesTheUnreachedRequirement`) | A carried entry that copies the prior performer fails the checkpoint owner test |
| CF35 | 7 | The carried entry keeps the source commit, the source digest, and the excerpt of the candidate | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketCarriesTheUnreachedRequirement`) | An entry bound to the repair commit claims a run that never occurred |
| CF36 | 8 | The native reference of the carried entry is `carried:` followed by the evidence ID of the candidate | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketCarriesTheUnreachedRequirement`) | An entry with the call reference loses the trace to the original run |
| CF37 | 9 | The `carried` cell of the ticket output lists each carried requirement ID | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketCarriesTheUnreachedRequirement`) | An output without the cell hides which checks did not rerun |
| CF38 | 10 | The `carried` cell of the ticket output is `none` when the call carries nothing | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketCarriesNothingWhenEachRequirementIsStaged`) | An empty cell reads as a missing value |
| CF39 | 11 | A staged requirement that the delta does not reach records from its stage at the repair commit | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketCarriesNothingWhenEachRequirementIsStaged`) | A verb that prefers the carry drops the run that the session made |
| CF40 | 12 | The call writes the staged entry and the carried entry in owed order | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketCarriesTheUnreachedRequirement`) | A verb that appends the carried entries after the staged entries breaks the owed order |
| CF41 | 35 | The refusal names each unstaged requirement that cannot carry with its reason after `; carry:` | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketNamesTheChecksToRun`) | A refusal without the reasons sends the session to rerun every check |
| CF42 | 36 | The refusal omits each unstaged requirement that can carry | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketNamesTheChecksToRun`) | A refusal that lists every unstaged requirement asks for a run that the carry covers |
| CF43 | 37 | A carry refusal leaves the record bytes and every staged excerpt unchanged | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestTicketNamesTheChecksToRun`) | A partial write leaves a carried entry beside a refused requirement |
| CF44 | 13 | The help line of `bench record ticket` states the carried entries | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A stale help line hides the carry from an agent |
| CF45 | 40 | The Land section of the implement phase states that a repair session runs each check that the repair delta reaches | review-owned: Spec axis reads the repair sentence of the Land section | A charge that follows the old sentence demands a full rerun |
| CF46 | 41 | A Forbid anchor refuses the implement-phase needle `runs each check that its ticket owes once at its repair commit` | planned TestChunkChainAnchors in internal/anchors/registry_chunk_chain_test.go | Without the anchor, the retired full-rerun sentence can return |
| CF47 | 42 | The changelog states the carry of `bench record ticket` and its test-only bound | review-owned: Spec axis reads the new changelog heading | A user does not learn that the verb changed |
| CF48 | 43 | On a one-chunk fixture, after the carry call under a repair dispatch and a record commit, `rr.Check` of the chunk refuses only through a review-form repair of that chunk | `internal/reviewrecord/recordcmd/ticket_carry_test.go` (`TestCarriedEntryPassesTheChunkGrade`) | A carried entry with an empty state, outcome, or exit code, or with the prior performer, fails the chunk grade |
| CF49 | 41 | A Forbid anchor refuses the operating-guide needle `reruns the ticket's verification with current repair coverage` | planned TestChunkChainAnchors in internal/anchors/registry_chunk_chain_test.go | Without the anchor, the retired operating-guide clause can return |
| CF50 | 45 | The operating guide states that a repair session reruns each ticket check that its repair delta reaches | review-owned: Spec axis reads the repair sentence of the operating guide | A session that loads the old clause demands a full rerun |
| CF51 | 25 | A test-only delta in `./internal/ab` refuses the requirement `bench test --package ./internal/a...` with the reason `the delta reaches the command` | `internal/reviewrecord/carry/carry_reach_test.go` (`TestCarryReach`) | A canonical rule that accepts a mid-path wildcard carries a package that `go test` runs |

Not covered: story 26 — a correct parse and a word-splitting parse both refuse the quoted expression, because the split expression is not canonical; the `shellcommand` lexer tests pin a single-quoted word that holds a space, and the shared parser reads that lexer.

### Edge inventory

The audience is every repository that links the kit, because each delegated build calls `bench record ticket`. In a repository without a Go module, no delta is test-only, so every requirement reruns. In a linked repository, the system suite and root conformance do not run. The conservative rules still apply there, and they only add reruns.

The canonical edge classes and the hostile-input checklist attach to these rows:

- A deleted, renamed, mode-changed, or `testdata` test file makes the delta not test-only, through the existing test-only file list. This disposition attaches to CF10.
- A `--changed` run parses to the `./...` target, which covers every reach operand. It is thus reached, and CF9 covers the parent cover rule.
- A requirement with no prior entry, a failed entry, an entry with no source commit, and an entry off the chunk line attach to CF23, CF24, CF26, and CF27.
- An absent record makes the first chunk entry. No prior entry exists, so CF23 covers it.
- A stage and a carry that are both possible attach to CF39.

Won't handle: a doc-only, production, widen, or fold repair delta carries nothing — the reviewer decided this on 2026-10-27, and the full rerun keeps each such repair correct.
Won't handle: the chunk checkpoint does not re-derive the reach of a carried entry — `bench record verification` already writes an entry at any commit on the chunk line, so the record verbs stay the one enforcement of the rerun.
Won't handle: a repair that reaches the requirements of another ticket of the same chunk does not refresh them — the ticket form already leaves the entries of other tickets unchanged, and the chunk review covers that delta.
Won't handle: the final verification list does not carry — the ticket form writes only chunk entries, and the orchestrator runs the final list once.

## Out of scope

- A checkpoint grade that re-derives the reach of each carried entry: 6 edits, 2 gate runs.
- A carry across a doc-only delta with a reach rule for documents: 8 edits, 3 gate runs.
- A refresh of the entries of other tickets that a repair delta reaches: 5 edits, 2 gate runs.

## Further notes

### Reader sweep

- The ticket form refusal text has two test pins, which compare with a substring test. One is the missing-stage case of `TestTicketRefusals`. The other is the `quoted` case in the stage tests of the record command. The new text keeps that prefix, and neither requirement has a prior entry, so each pin stays green.
- The help line of the ticket form has one pin, in `TestHelpInventoryIsComplete`.
- The ticket output header has no reader outside the ticket form.
- The full-rerun rule has two statements in guidance: one sentence in the Land section of the implement phase, and one clause in the operating guide. The review phase, the delegation discipline, and the reference do not state it. No anchor pins either statement.
- The evidence reader of the build preflight reads the latest entry of a requirement and its probe excerpt. A carried entry keeps the excerpt of its candidate, so that reader is unchanged.
- The plan-command parser and the package-cover predicate have one caller each, in the check floor.
- The whole-only rule of the delta set has one caller, the delta plan of the gate.

### Interaction with FT455

FT455 records a reviewer decision of 2026-10-27: the test-only predicate is to class ticket `Writes:` lines, the plan fence, and the paths of a sibling chunk as plan paths. The carry reads that same predicate. If FT455 ships first, a widen commit no longer stops the carry, and row CF11 changes with it. This spec flags that interaction for reviewer veto and adds no second predicate.

### Source-sentence-to-row table

| source sentence | rows |
|---|---|
| A repair session reruns every plan requirement of its ticket at the repair commit. | problem statement; CF45, CF46, CF49, CF50 |
| The record verb refuses a repair record that lacks an owed requirement. | CF41, CF42, CF43 |
| The record verb accepts a repair entry that carries forward each requirement whose packages the repair delta cannot reach. | CF4, CF5, CF6, CF7, CF22, CF33 to CF40, CF48 |
| The changed-package selector of the delta checkpoint decides the reach. | CF8 to CF21, CF31, CF32, CF51 |
| The repair reruns the rest. | CF8 to CF21, CF23 to CF28, CF41 |

### Pre-review proof checklist

- Cited symbols: `TicketRequirements`, `RecordTicket`, `freezeTip`, `ownedStages`, `ownedRuns.calls`, `checkVerification`, `EntryCurrent`, `ProbeBites`, `RunBinding`, `changeVerification`, `TestOnlyFiles`, `deltaPackages`, `deltaPhases`, `planSelection`, `packageCovers`, `Prepare`, `Reassign`, `repairedEarlierTicket`, `TestProbeChecksAgree`, `Check`, and `Unstage` resolve in the tree at `48f9241df`.
- Import edges: `recordcmd` already imports `testreport`. The gate does not import `recordcmd`, `testreport`, or `preflight`. The `shellcommand` package imports no internal package. The new carry package imports the review record, the gate, and `testreport`, and none of them imports it.
- Source-row clauses and occurrences: the table above quotes each clause. The one occurrence is the FT445 repairs of 2026-10-24.
- Promised field labels: the `carried` cell, the `carried:` reference prefix, and the four reasons `no current passing entry`, `the delta is not test-only`, `the delta changes a whole-only test file`, and `the delta reaches the command`.
- Changed-function callers: `ownedStages` and `ownedRuns.calls` have one caller, the ticket form. `checkVerification` has three callers: the chunk grade, the completion grade, and `TestProbeChecksAgree`. `deltaPackages` has one caller, the delta plan. `planSelection` and `packageCovers` have one caller each, in the check floor.
- Copy survival: CF3 and CF32.
- Rendered-shape readers: the old refusal text has two hits, in `TestTicketRefusals` and in the `quoted` stage test. The old help line has one hit, in `TestHelpInventoryIsComplete`. The output header has no hit outside the ticket form. Each retired guidance needle has one hit, in its own file.
- Pin operators: the refusal test compares with a substring test, and the help test compares whole lines.
- Entry reads: none. The carry decision takes the root, the spec, the chunk entry, and the requirement as values.
- Derived expectations: the tests derive each reason and each package operand from the fixture delta, not from the code under test.
- Consolidated rules: none.
- Quantified obligations: each owed requirement without a stage takes one decision. Each carried requirement takes one entry.
- Workflow-step writes: the carried entry joins the chunk entry in the one record write of the ticket form. The record path is outside every source digest, so the write moves no digest that a later checkpoint compares.

### Flagged additions

- The `carried` cell of the ticket output is an addition beyond the source. It supports story 9, and CF37 and CF38 cover it.
- The refusal reasons are an addition beyond the source. They support story 35, and CF41 covers them.
- The operating-guide edit is an addition that the spec review found. It supports story 45, and CF49 and CF50 cover it.
- The rule that a changed source with no changed test file reaches each conformance check settles a case that the spec review found. It supports story 44, and CF19 covers it.
