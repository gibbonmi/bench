# Record ticket

Status: staged

Roadmap: FT405

Decision source: reviewer-confirmed conversation, 2026-10-14

Verification log: 0 iteration(s) to accept — pending review.

## Problem

A ticket author records its chunk evidence in several separate calls today.
It runs `bench record amendment` when the plan digest moved, then `bench record chunk`, then one `bench record verification` call for each requirement.
Each call is its own write, and `replace` rolls back nothing, so a later refusal leaves an earlier entry written.
In the FT379 build, the freeze and amendment order cost three retries.

Each verification excerpt must also exist as a file.
Bench keeps no copy of a response within the response bound, so an author copies each short response into a file by hand.
In the NRO build, an author copied 13 excerpts by hand.

A reviewer cannot bind evidence beside an uncommitted record write.
`bench preflight evidence --check-current` refuses any dirty path, and the record file is dirty after each record write.

Some guidance still names the review phase or no owner as the freezer.
The implementation phase names the last ticket's author, and the roadmap row FT405 still says that the orchestrator owns the freeze.

## Solution

Bench keeps the complete last response of each assignment worktree, with its exit code and its argument list.
`bench record excerpt` stages that response as the excerpt of one planned requirement.
`bench record ticket` then writes every entry that one ticket owes in one write.
That write holds the amendment when the plan digest moved, the freeze for the last ticket, and each verification.
A refusal at any step leaves the record bytes unchanged.

`bench preflight evidence --check-current` reports a dirty review record of its spec and binds the evidence.
Any other dirty path still refuses.

The operating guide, the review phase, the ticket craft, and the reference name the last ticket's author as the freezer.
The implementation phase names `bench record ticket` for the freeze and the verification.

## User stories

Line: opus / high.

Implementation-line reason: The hardest material chunk is RT-TICKET, which composes three record writers into one transaction on the shared review record.
A wrong composition writes a partial record or a second copy of a writer rule, and every later checkpoint reads that record.
The spec is precise, because each composed writer already exists with its refusals and its tests.
Seam uncertainty is low: `write` is the one record transaction, and `runBounded` is the one bounded response owner.
The gate tests are strong, because the record-command tests drive real Git fixtures and read the record bytes.

Harder chunks: RT-TICKET.

### Retained response

1. As a ticket author, I want Bench to keep the complete last response of my worktree, so that I never copy an excerpt by hand.
2. As a ticket author, I want a response within the response bound kept, so that a short response still gives an excerpt.
3. As a ticket author, I want a spilled response kept complete, so that the excerpt is the native output and not the projection.
4. As a ticket author, I want the kept response to hold its exit code and arguments, so that I type no exit code.
5. As a ticket author, I want a `bench record` call to keep the kept response, so that staging reads my check.
6. As a ticket author, I want a `bench worktree` call to keep the kept response, so that an exec call keeps its child's response.
7. As an operator, I want the primary checkout to keep no last response, so that the primary scope holds only its capped spills.
8. As a reviewer, I want a response whose spill failed to clear the kept response, so that incomplete bytes never become evidence.
9. As an operator, I want a failed retention write to leave the response and the exit code unchanged, so that retention never breaks a verb.
10. As an operator, I want a worktree release to remove its kept response and stages, so that no evidence outlives its assignment.

### Staged excerpt

11. As a ticket author, I want `bench record excerpt` to stage the kept response for one planned requirement, so that each requirement has its own excerpt.
12. As a ticket author, I want staging to refuse when no response is kept, so that I learn to run the check first.
13. As a reviewer, I want staging to refuse a kept response whose command differs from the planned command, so that a wrong response never stages.
14. As a reviewer, I want a probe requirement to stage only a complete `bench probe` answer, so that the later import cannot fail on the stage.
15. As a ticket author, I want staging to refuse an unplanned requirement and name the planned ones, so that a typo fails at once.
16. As a ticket author, I want a second stage of one requirement to replace the first, so that the latest run is the excerpt.
17. As a ticket author, I want staging to refuse a response over the excerpt limit, so that the stage fails before the record write.

### Ticket record

18. As a ticket author, I want `bench record ticket` to write every verification my ticket owes, so that one call records my ticket.
19. As the last ticket's author, I want the same call to freeze the chunk, so that the freeze and the verification share one write.
20. As a ticket author, I want the same call to record a needed amendment, so that the amendment order never costs a retry.
21. As a reviewer, I want a refusal at any step to leave the record bytes unchanged, so that no partial record exists.
22. As a ticket author, I want one refusal to name every owed requirement with no stage, so that one retry fixes every gap.
23. As a ticket author, I want an earlier ticket's call to refuse a chunk with no entry, so that the last ticket's author freezes.
24. As a reviewer, I want a call for an earlier ticket to leave the chunk entry unchanged, so that only the last ticket's author freezes.
25. As a repair session, I want a refreeze to keep the base of the existing chunk entry, so that the chunk base stays fixed.
26. As the last ticket's author, I want a new entry to take the predecessor entry's tip as base, so that I type no base.
27. As the first chunk's author, I want a call with no base source to refuse, so that the base is never guessed.
28. As a ticket author, I want a plan change that needs a chunk map to refuse, so that only `bench record amendment` maps chunks.
29. As a reviewer, I want the call to refuse a planned seam cell, so that the ticket form keeps the chunk form's citation rule.
30. As a ticket author, I want the call to refuse a ticket outside the chunk, so that a wrong ticket name fails at once.
31. As a reviewer, I want the call to refuse requirements that name no ticket, so that the owed set comes from a version 2 plan.
32. As a ticket author, I want each evidence ID to join my prefix and the requirement ID, so that I type one prefix.
33. As a reviewer, I want a probe requirement to import the staged probe answer, so that the probe values come from the native answer.
34. As a ticket author, I want the staged excerpts of a written call removed, so that a refreeze never reuses an old excerpt.
35. As a ticket author, I want the staged excerpts kept after a refused call, so that I retry without a new run.
36. As a maintainer, I want the ticket form to compose the existing writer rules, so that no writer rule has two copies.
37. As a ticket author, I want `bench record` help to list the ticket form, so that the form is discoverable.

### Current binding

38. As a reviewer, I want `--check-current` to bind beside a dirty record of its spec, so that an uncommitted record write blocks nothing.
39. As a reviewer, I want the binding to report the dirty record path, so that I see the uncommitted write.
40. As a reviewer, I want any other dirty path to refuse the binding, so that the working-tree guard keeps its strength.
41. As a reviewer, I want a dirty review record of another spec to refuse, so that only the exact record path relaxes.
42. As a maintainer, I want the charge, the review charge, and the proposal to keep their dirty refusal, so that only the current binding relaxes.
43. As a maintainer, I want the generated evidence reference to name the new field, so that the reference and the schema agree.

### Guidance

44. As a ticket author, I want the implementation phase to name `bench record ticket` for the freeze and the verification, so that I use the one-call form.
45. As a coordinator, I want the operating guide to name the last ticket's author as the freezer, so that no reader assigns it elsewhere.
46. As a reviewer, I want the review phase to name the last ticket's author as the freezer, so that both phases agree.
47. As a reviewer, I want the review phase's record step to list the excerpt and ticket forms, so that the form list stays complete.
48. As a slicer, I want the ticket craft and the reference to name the last ticket's author as freezer, so that they agree.
49. As a coordinator, I want the known-issues row on file excerpts removed, so that no stale workaround stays.
50. As an adopter, I want a changelog entry, so that the release notes name the new forms.

### Exclusions

51. As a reviewer, I want one repair session to keep one ticket, so that the one-session-one-ticket rule stays unchanged.
52. As a maintainer, I want `bench record verification` to keep its file excerpt, so that the orchestrator's final records keep their form.

### Run binding and hostile input

53. As a reviewer, I want a stage bound to its run commit, so that an earlier tip's excerpt never records at the chunk source.
54. As a reviewer, I want a stage from a dirty checkout refused, so that an excerpt never comes from uncommitted bytes.
55. As the last ticket's author, I want a base after a fold to be that fold's commit, so that the base obeys the fold chain.
56. As the first chunk's author, I want the ticket form to create an absent record, so that the first freeze needs no other form.
57. As a ticket author, I want a planned command with shell quotes to match my run, so that a quoted operand stages.
58. As a reviewer, I want a hostile requirement ID to stay inside the assignment scope, so that a stage never writes outside its folder.

## Implementation decisions

### Decision record

The reviewer decided on 2026-10-14 that the chunk freeze belongs to the ticket author.
The implementation phase already states that rule, and CHANGELOG.md records it, so this spec adds no freeze-owner row for that phase.
The roadmap row's premise that the orchestrator owns the freeze is stale; this spec corrects every other text.
The reviewer chose the scope on 2026-10-14: the core forms and the dirty-record report.
The repair session that records for every ticket of its chunk stays out of scope.

### Retained response

`runBounded` is the one owner of a bounded response, so it owns retention.
After `Finish`, it writes the retained response of the assignment scope that `AssignmentScope` answers.
The retained response holds the complete response bytes in arrival order, the exit code, the typed verb word, and its arguments.
It also holds the `HEAD` commit of the worktree and its dirty paths, read after the verb returns.
The dirty paths come from one shared porcelain reader: `git status --porcelain=v1 -z --no-renames` and `ParsePorcelainZ`.
`LandedState` and `Facts` already run that read, so the read moves into one `internal/git` function that `LandedState`, `Facts`, and retention call.

The commit comes from `rev-parse HEAD`.
Both reads run only when retention writes, so the skipped calls below and the primary checkout pay no Git read.

The owner stays the one source of the bytes: `responsebound` gains an accessor that answers the complete bytes of a finished response.
That accessor answers the retained writes of a response within the bound and the spill file bytes of a spilled response.
It answers no bytes for the spill-failed state and for the `passing` state, which is the create-failure route.
The lead identity row is outside the bound and is not part of the retained bytes.

These calls leave the retained response unchanged:

- a call in the primary checkout or outside a repository
- a call of the `record` family
- a call of the `worktree` family
- a call that a retiring leaf runs

A response with no complete bytes clears the retained response.
Retention writes through one durable replacement, and its error changes neither the response nor the exit code.
The retained response lives under the assignment scope of the response store, so `Drop` removes it at release.
Its file names do not match the spill name pattern, so the primary prune never reads them.

### Staged excerpt

`bench record excerpt <slug> --chunk <id> --requirement <id> [--mutant <n>]` stages the retained response for one requirement of the chunk plan at the worktree's `HEAD`.
It refuses with no retained response, with an unplanned requirement, and with bytes over the record excerpt limit.

`shellcommand.Parse` is the one splitter of a planned command, so a quoted operand folds as the shell folds it.
The stage refuses a planned command unless the parse is lexed, holds only word tokens, and starts with the word `bench`.
For a requirement with no probe, the typed verb word and its arguments must equal the planned words after `bench`.
The comparison uses the typed verb word, so a planned alias matches only a run that typed the same alias.
For a requirement with a probe, the retained call must be `bench probe`, and the bytes must import through the probe reader.
`--mutant` is valid only for a probe requirement, and the stage keeps it for the import.

The stage lives under the assignment scope in one file per key.
The file name is the hex SHA-256 digest of the spec slug, the chunk ID, and the requirement ID, joined by NUL bytes.
No plan string forms a path component, so a `/` or `..` in an ID cannot leave the scope.
The stage copies the retained commit and dirty paths with the bytes.
A second stage of one requirement replaces the first.
The excerpt form writes no review record and no worktree file.

### Ticket record

`bench record ticket <slug> --chunk <id> --ticket <basename> --id-prefix <prefix> --performer <session> --model <model> --effort <effort> --ref <ref> [--base <commit>] [--tip <commit>]` writes one ticket's entries.
`--tip` defaults to `HEAD`.
The owed requirements are the requirements of the chunk plan at the tip whose ticket is the named ticket.

Before the write, the form refuses in this order:

1. a ticket that the chunk plan does not list
2. a plan whose chunk requirements name no ticket
3. a planned seam cell at the tip, when the ticket is the last ticket of the chunk
4. an owed requirement with no staged excerpt, naming every missing requirement
5. a stage whose run commit is not the chunk source, naming the requirement
6. a stage whose run had a dirty path other than the record path of the spec, naming the requirement

Inside one `write` transaction, the form applies these steps in the completion order:

1. When the record plan digest differs from the plan at the tip, it applies the amendment with no chunk map.
2. When the ticket is the last ticket of the chunk, it freezes the chunk at the base and the tip.
3. It appends one verification entry for each owed requirement at the chunk source.

The chunk source is the freeze tip for the last ticket and the chunk entry's tip for an earlier ticket.
An earlier ticket's author therefore reruns its checks at the chunk tip before it stages.

The amendment step refuses a plan change that needs a chunk map with the amendment change function's own error.
The ticket form wraps that error as a repair that names `bench record amendment`.
The change function returns the same error as today, so the amendment form's own refusal text stays unchanged.

The base is `--base` when given, then the base of the existing chunk entry.
Otherwise it is the commit that `coveredAfter` answers for the predecessor chunk's entry in plan order.
`coveredAfter` answers the last recorded fold after that chunk, or the chunk tip when no fold follows it.
This rule matches the chunk base rule of the review phase's step 1, and the checkpoint's fold chain reads the same function.
With none of them, the form refuses and asks for `--base`.
A call for an earlier ticket needs an existing chunk entry, and it never changes that entry.

Each evidence ID is the prefix, a hyphen, and the requirement ID.
The exit code comes from the stage, and a probe requirement imports its staged answer with its staged `--mutant`.
After the write succeeds, the form removes the staged excerpts that it used.

The exported amendment, chunk, and verification writers each keep one change function over the record.
The chunk writer's create function, which owns the version 2 refusal for a new record, becomes one shared function too.
Each exported writer calls `write` with its own functions.
The ticket form calls `write` once, with the chunk create function and the three change functions in sequence.
An earlier ticket's call passes the absent-record refusal as its create function, as the verification writer does.
No refusal text or rule moves into the ticket form.

| site | old rule | new rule |
|---|---|---|
| amendment writer | The change body sits inside `RecordAmendment`. | One change function owns the body, and `RecordAmendment` and the ticket form call it. |
| chunk writer | The change body sits inside `RecordChunk`. | One change function owns the body, and `RecordChunk` and the ticket form call it. |
| verification writer | The change body sits inside `RecordVerification`. | One change function owns the body, and `RecordVerification` and the ticket form call it. |
| chunk create | The create closure sits inside `RecordChunk`. | One create function owns the version 2 refusal, and `RecordChunk` and the ticket form call it. |

`FormTicket` and `FormExcerpt` join the form constants, and the record command registers both forms.
The completion order stays unchanged.
The ticket form's help derives its step list from `CompletionOrder`: the steps before the first review step.

### Current binding

The current binding passes its action to the checkout refusal.
For the check action only, the refusal skips one dirty path: the record path of the selected spec.
The `current` response block gains a `dirty_record` flag field after `current`, true when that path is dirty.
The charge, the review charge, and the proposal keep the strict refusal.

The record path still reaches the changed paths, so `paths-authorized` grades it against the fence, as it does today.
The record's source digest excludes the record path (`SourceDigest`).
The binding compares the source tip and reads each required source from the tip's tree, so a dirty record changes no bound byte.

### Guidance

The implementation phase changes two sentences in place, so its sentence count stays the same.
The freeze sentence and the verification sentence name `bench record ticket`.
The operating guide's chunk sentence names the last ticket's author as the freezer.

The review phase's record step names the last ticket's author as the writer of the frozen pair.
A new sentence in that step lists the excerpt and ticket forms.
The step 6 sentence "Write each completed chunk, verification, review, and amendment entry with `bench record`, not with hand-built JSON." stays byte for byte.

The ticket craft and the reference name the last ticket's author as the freezer.
The known-issues row on file excerpts goes, and the changelog gains one entry.

### Commitment binding

The commitment criteria of FT405 still name `--excerpts <dir>`, a composition of the assignment and verification forms, and the repair-session criterion.
The spec-stage close binds this spec through `bench commitment plan`, and that deliverable binding restates the criteria from this spec.
The repair-session criterion leaves the binding and returns to intake as a Won't handle line.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RT-RETAIN / 01-retain-the-last-bounded-response.md | Bench keeps the last response of each assignment worktree. | RT1–RT8, RT65 | the planned retention tests and the system suite's nested exec test | no |
| RT-EXCERPT / 02-stage-an-excerpt-for-a-requirement.md | `bench record excerpt` stages the retained response for one requirement, and a release removes the stages. | RT9–RT17, RT57–RT59 | the planned excerpt tests and `TestHelpInventoryIsComplete` | no |
| RT-CURRENT / 03-report-a-dirty-record-in-the-current-binding.md | `--check-current` binds beside a dirty review record of its spec and reports it. | RT40–RT46 | `TestEvidenceCurrentBinding`, `TestEvidenceSchemaRegistry`, and the planned binding tests | no |
| RT-EXTRACT / 04-extract-the-record-change-functions.md | Each record writer calls one shared change function, with no change in behavior. | RT39 | `TestChunkNeedsCurrentAmendment` and the existing writer tests | no |
| RT-TICKET / 05-record-one-ticket-in-one-write.md | `bench record ticket` writes the amendment, the freeze, and the owed verifications in one transaction. | RT18–RT38, RT55, RT56, RT60–RT63 | `TestCompletionHelpOrdering`, `TestHelpInventoryIsComplete`, and the planned ticket tests | yes |
| RT-GUIDE / 06-name-the-ticket-author-as-the-freezer.md | The guidance names the last ticket's author as the freezer and the ticket form as the record call. | RT47–RT54, RT64 | `TestChunkChainAnchors` with `registry_chunk_chain.go` and its twin, and `TestRetainedWorkflow` with the `.bench/BENCH.md` needle in `registry_ft311_review_dispatch.go` | no |

Each ticket that creates a seam for a later ticket is its own review chunk.
RT-RETAIN creates the retained response, RT-EXCERPT creates the stage, and RT-EXTRACT creates the change functions.
The slicing split the earlier RT-RETAIN chunk into RT-RETAIN and RT-EXCERPT, and the earlier RT-TICKET chunk into RT-EXTRACT and RT-TICKET.

RT-RETAIN also owns the shared porcelain reader in `internal/git/status.go` and the `LandedState` call site there.
RT-CURRENT shares the command registry files with RT-EXCERPT and RT-TICKET, so it runs between them.
RT-TICKET stays one ticket, because the earlier-ticket route and the base rule run inside the same transaction as the last-ticket route.
RT-GUIDE comes last, because its sentences name the forms that RT-TICKET ships.

The central-property mutation of RT-TICKET makes the ticket form run three separate writes.
RT21 and RT22 turn red under that mutation.

### Completion plan

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "RT-RETAIN",
      "tickets": ["01-retain-the-last-bounded-response.md"],
      "verification": [
        {"id": "retention", "command": "bench test --package ./cmd/bench --package ./internal/responsebound/... --package ./internal/git --package ./internal/gate --package ./internal/handoff --package ./internal/roadmap --package ./internal/shift --package ./internal/status --package ./internal/worktree --package ./internal/conformance"},
        {"id": "retain-root", "command": "bench test --check root-conformance"},
        {"id": "nested-exec", "command": "bench test --check system --run TestExecNested"},
        {"id": "record-skip", "command": "bench test --package ./cmd/bench", "probe": "drop the record-family skip from retention"}
      ]
    },
    {
      "id": "RT-EXCERPT",
      "tickets": ["02-stage-an-excerpt-for-a-requirement.md"],
      "verification": [
        {"id": "excerpt-stage", "command": "bench test --package ./internal/reviewrecord/... --package ./internal/responsebound/... --package ./internal/shellcommand --package ./cmd/bench --package ./internal/conformance"},
        {"id": "excerpt-root", "command": "bench test --check root-conformance"},
        {"id": "command-match", "command": "bench test --package ./internal/reviewrecord/recordcmd", "probe": "compare only the first planned word with the retained call"}
      ]
    },
    {
      "id": "RT-CURRENT",
      "tickets": ["03-report-a-dirty-record-in-the-current-binding.md"],
      "verification": [
        {"id": "binding", "command": "bench test --package ./internal/preflight/... --package ./internal/chargeevidence --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance"},
        {"id": "binding-root", "command": "bench test --check root-conformance"},
        {"id": "strict-charge", "command": "bench test --package ./internal/preflight/...", "probe": "skip the record path for every preparation action"}
      ]
    },
    {
      "id": "RT-EXTRACT",
      "tickets": ["04-extract-the-record-change-functions.md"],
      "verification": [
        {"id": "extract-writers", "command": "bench test --package ./internal/reviewrecord/..."},
        {"id": "extract-root", "command": "bench test --check root-conformance"}
      ]
    },
    {
      "id": "RT-TICKET",
      "tickets": ["05-record-one-ticket-in-one-write.md"],
      "verification": [
        {"id": "record-writers", "command": "bench test --package ./internal/reviewrecord/... --package ./cmd/bench --package ./internal/responsebound --package ./internal/conformance"},
        {"id": "record-root", "command": "bench test --check root-conformance"},
        {"id": "one-write", "command": "bench test --package ./internal/reviewrecord/...", "probe": "run the ticket form as three separate writes"}
      ]
    },
    {
      "id": "RT-GUIDE",
      "tickets": ["06-name-the-ticket-author-as-the-freezer.md"],
      "verification": [
        {"id": "guidance-anchors", "command": "bench test --package ./internal/conformance --package ./internal/anchors --package ./cmd/bench"},
        {"id": "guidance-root", "command": "bench test --check root-conformance"},
        {"id": "guidance-prose", "command": "bench test --check prose-mechanics"},
        {"id": "guidance-budgets", "command": "bench test --check guidance-prose-budgets"}
      ]
    }
  ],
  "final_verification": [
    {"id": "coverage", "command": "bench coverage --check specs/record-ticket/spec.md"},
    {"id": "record-and-binding", "command": "bench test --package ./internal/reviewrecord/... --package ./internal/preflight/... --package ./internal/chargeevidence --package ./internal/responsebound/... --package ./internal/git --package ./cmd/bench"},
    {"id": "guidance", "command": "bench test --package ./internal/conformance --package ./internal/anchors"},
    {"id": "system", "command": "bench test --check system --run TestExecNested"}
  ]
}
```

## Testing decisions

- A good record test drives `recordcmd.Command` on a real Git fixture and reads the record bytes and the exit code.
- The retention tests drive one bounded call through the command dispatcher in a fixture pool worktree, then read the retained response.
- The excerpt tests stage a retained response and read the staged bytes and the refusal text.
- The ticket tests stage excerpts, run the ticket form, and compare the record bytes before and after each refusal.
- The binding tests reuse the `TestEvidenceCurrentBinding` fixture and dirty one path for each case.
- The guidance rows use the anchor registry and its independent twin.
- The prior art is `TestChunkNeedsCurrentAmendment` for the record order. The nested exec test of the system suite is the prior art for RT5, which needs a real child process.
- The gate observes the feature through the `test` phase and the anchor checks.

### Seam diagram

    trigger: the ticket author after its ticket commit
        │
        ▼
    bench test / bench probe ──▶ [ runBounded: retain the last response ]
        │
        ▼
    bench record excerpt ──▶ [ stage: requirement → retained bytes, exit, argument list ]
        │
        ▼
    bench record ticket ──▶ [ one write: amendment → chunk → verifications ]
                               ◀ tests attach here: recordcmd.Command on a fixture repository, record bytes and refusal text

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RT1 | 1, 2 | A `bench test` call within the bound in an assignment worktree leaves a retained response with its exact bytes. | planned | A retention that runs only on a spill leaves no retained response and the row reds. |
| RT2 | 3 | A spilled `bench test` response leaves a retained response that equals the spill file bytes and holds no `spilled{` frame. | planned | A retention that keeps the printed projection holds the frame and the row reds. |
| RT3 | 4 | The retained response holds the verb's exit code and its argument list after `bench`. | planned | A retention that drops the exit code reads as zero for a failing call and the row reds. |
| RT4 | 5 | A `bench record` call after a `bench test` call leaves the `bench test` response retained. | planned | A retention without the record-family skip replaces the response and the row reds. |
| RT5 | 6 | A `bench worktree exec` call in the worktree leaves its child's `bench test` response retained. | planned | A retention without the worktree-family skip replaces the child's response and the row reds. The test needs a real child process, so it sits in the system suite that `bench test --check system` and the gate's system phase run. |
| RT6 | 7 | A bounded call in the primary checkout writes no retained response. | planned | A retention that writes the primary scope leaves a file there and the row reds. |
| RT7 | 8 | A response in the spill-failed state removes the earlier retained response. | planned | A retention that skips a failed response leaves the earlier response and the row reds. |
| RT8 | 9 | A retention write that fails leaves the printed response and the exit code unchanged. | planned | A dispatcher that returns the retention error changes the exit code and the row reds. |
| RT9 | 10 | The release of a worktree removes its retained response and its staged excerpts. | planned | A retention outside the assignment scope survives the drop and the row reds. |
| RT10 | 11 | `bench record excerpt` stages the retained `bench test` response for a planned requirement whose command matches. | planned | A stage that copies nothing leaves no staged excerpt and the row reds. |
| RT11 | 12 | `bench record excerpt` refuses with no retained response. | planned | A stage that stages an empty excerpt passes and the row reds. |
| RT12 | 13 | `bench record excerpt` refuses a retained response whose argument list differs from the planned command. | planned | A stage without the command comparison stages the wrong response and the row reds. |
| RT13 | 14 | `bench record excerpt` refuses a probe requirement whose retained call is not `bench probe`. | planned | A stage that skips the probe-call check stages a test response and the row reds. |
| RT14 | 14 | `bench record excerpt` refuses a mutant-set probe answer without `--mutant`. | planned | A stage that skips the import stages an answer the ticket form cannot import and the row reds. |
| RT15 | 15 | `bench record excerpt` refuses an unplanned requirement and names the planned requirements. | planned | A stage that keys by any name stages a typo and the row reds. |
| RT16 | 16 | A second stage of one requirement replaces the first staged bytes. | planned | A stage that keeps the first file holds the old bytes and the row reds. |
| RT17 | 17 | `bench record excerpt` refuses a retained response over the record excerpt limit. | planned | A stage without the bound stages bytes that the record write refuses later and the row reds. |
| RT18 | 18, 19 | For the last ticket, `bench record ticket` writes the chunk entry at the tip and one verification entry per owed requirement. | planned | A form that writes only the verifications leaves no chunk entry and the row reds. |
| RT19 | 20 | With a moved plan digest, the form writes one amendment entry before the chunk entry. | planned | A form without the amendment step refuses with the stale plan digest and the row reds. |
| RT20 | 20 | With an unchanged plan digest, the form writes no amendment entry. | planned | A form that always amends refuses with the unchanged-plan refusal and the row reds. |
| RT21 | 21 | When the last owed verification refuses, the record bytes equal the bytes before the call. | planned | A form that runs three writes leaves the amendment and the chunk entry and the row reds. |
| RT22 | 21 | When the chunk step refuses after a moved plan digest, the record bytes equal the bytes before the call. | planned | A form that writes the amendment first leaves the amendment entry and the row reds. |
| RT23 | 22 | The form refuses when two owed requirements have no staged excerpt, and the refusal names both requirements. | planned | A form that stops at the first missing requirement names one requirement and the row reds. |
| RT24 | 23 | A call for an earlier ticket with no chunk entry refuses and names the last ticket of the chunk. | planned | A form that freezes for any ticket writes the entry and the row reds. |
| RT25 | 24 | A call for an earlier ticket with a chunk entry leaves the chunk tip and base unchanged. | planned | A form that refreezes for every ticket moves the tip to the call's tip and the row reds. |
| RT26 | 25 | A refreeze for the last ticket keeps the base of the existing chunk entry. | planned | A form that derives the base from the predecessor moves the base and the row reds. |
| RT27 | 26 | A new entry for a later chunk takes the predecessor chunk entry's tip as its base when no recorded fold follows that entry. | planned | A form that takes the plan base or `HEAD^` writes another base and the row reds. |
| RT28 | 27 | A call for the first chunk with no `--base` and no chunk entry refuses and asks for `--base`. | planned | A form that guesses a base writes an entry and the row reds. |
| RT29 | 28 | A plan change that removes a recorded chunk refuses, and the refusal names `bench record amendment`. | planned | A form that maps the chunk itself writes an amendment and the row reds. |
| RT30 | 29 | The form refuses a planned seam cell at the tip for the last ticket and names the rows. | planned | A form without the planned-cell check freezes the chunk and the row reds. |
| RT31 | 30 | The form refuses a ticket that the chunk plan does not list. | planned | A form that reads every requirement of the chunk records the wrong set and the row reds. |
| RT32 | 31 | The form refuses a plan whose chunk requirements name no ticket. | planned | A form that treats every requirement as owed records another ticket's entries and the row reds. |
| RT33 | 32 | Each verification entry ID is the prefix, a hyphen, and the requirement ID. | planned | A form that uses the prefix alone refuses the second entry as a duplicate ID and the row reds. |
| RT34 | 33 | A probe requirement records the probe outcome, focused exit, and restore of the staged answer. | planned | A form that records the probe as absent refuses with the probe mismatch and the row reds. |
| RT35 | 34 | After a written call, the staged excerpts of its owed requirements are absent. | planned | A form that keeps the stage leaves the files and the row reds. |
| RT36 | 35 | After a refused call, the staged excerpts of its owed requirements are present. | planned | A form that removes the stage before the write removes the files and the row reds. |
| RT37 | 36 | The amendment, chunk, and verification change functions each have one definition that the exported writer and the ticket form call. | review-owned | Review compares the call sites, because a copied writer body passes every behavior test. |
| RT38 | 37 | `bench record` help lists the ticket form and the excerpt form with their grammar. | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | The inventory reds when a registered form has no help row. |
| RT39 | 36 | A chunk entry with a stale plan digest still refuses through the chunk form. | `internal/reviewrecord/recordcmd/ordering_test.go` (`TestChunkNeedsCurrentAmendment`) | A refactor that drops the stale check from the chunk change function passes the stale chunk and the row reds. |
| RT40 | 38, 39 | `--check-current` with a dirty record of its spec binds and reports `dirty_record` true. | planned | A binding that keeps the strict refusal refuses and the row reds. |
| RT41 | 39 | `--check-current` on a clean checkout reports `dirty_record` false. | planned | A binding that reports true for every call reds the clean case. |
| RT42 | 40 | `--check-current` with a dirty untracked Go file still refuses with the dirty checkout refusal. | `internal/preflight/evidencecmd/evidence_modes_test.go` (`TestEvidenceCurrentBinding`) | A binding that skips every dirty path passes the case and the row reds. |
| RT43 | 40 | `--check-current` with a dirty record and a dirty Go file refuses. | planned | A binding that skips the check when the record is dirty passes the Go file and the row reds. |
| RT44 | 41 | `--check-current` with a dirty record of another spec refuses. | planned | A binding that skips every path under `reviews/` passes the case and the row reds. |
| RT45 | 42 | The build charge with a dirty record of its spec still refuses with the dirty checkout refusal. | planned | A relaxation in the shared refusal passes the charge and the row reds. |
| RT46 | 43 | The generated evidence reference lists `dirty_record` in the `current` block. | `internal/chargeevidence/format_test.go` (`TestEvidenceSchemaRegistry`) | A schema change without the regenerated reference reds the registry comparison. |
| RT47 | 44 | The implementation phase's Land section states that the last ticket's author freezes the chunk with `bench record ticket`. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | The changed anchor and its independent twin red when the sentence is absent. |
| RT48 | 44 | The implementation phase's Land section states that each ticket author writes its verification entries with `bench record ticket`. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | The changed anchor and its independent twin red when the sentence is absent. |
| RT49 | 45 | The operating guide states that the last ticket's author freezes the chunk delta before the three axes run. | `internal/conformance/retained_workflow_test.go` (`TestRetainedWorkflow`) | The changed needle and its family bite red when the sentence is absent. |
| RT50 | 46 | The review phase's step 6 states that the last ticket's author writes the frozen pair with `bench record ticket`. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A new anchor and its independent twin red when the sentence is absent. |
| RT51 | 48 | The ticket craft states that the last ticket's author freezes the chunk base and the current tip. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A new anchor and its independent twin red when the sentence is absent. |
| RT52 | 48 | The ticket craft and the reference hold no `review freezes` text. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A Forbid anchor on `review freezes` in each file reds a restored sentence. |
| RT53 | 49 | The known-issues file holds no row that names FT405. | planned | A row that survives the change names FT405 and the row reds. |
| RT54 | 48 | The reference states that the last ticket's author freezes each chunk's base and current tip. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A new anchor and its independent twin red when the sentence is absent. |
| RT55 | 55 | For a later chunk whose predecessor entry has one recorded fold after it, the ticket form freezes with that fold's commit as base. | planned | A form that takes the predecessor tip writes a base that the checkpoint refuses as a chain gap and the row reds. |
| RT56 | 56 | For the last ticket of the first chunk with an absent record and `--base`, the ticket form creates the record with the chunk entry and the verifications. | planned | A form that passes the absent-record refusal as its create function refuses and the row reds. |
| RT57 | 57 | `bench record excerpt` stages a run whose arguments equal a planned command with a single-quoted operand. | planned | A stage that splits on spaces keeps the quotes and refuses the match and the row reds. |
| RT58 | 13 | `bench record excerpt` refuses a planned command whose first word is not `bench`, a planned command that holds a pipe operator, and a planned command that the lexer cannot read. | planned | A stage that drops the first word compares a foreign command, and a stage that checks only the first word accepts the pipe or the unlexed text, so the row reds. |
| RT59 | 58 | A requirement ID of `../x` stages a file inside the assignment scope and nowhere else. | planned | A stage that joins the ID into the path writes outside the scope and the row reds. |
| RT60 | 53 | The ticket form refuses a stage whose run commit is not the chunk source and names the requirement. | planned | A form that ignores the run commit records an excerpt from an earlier tip and the row reds. |
| RT61 | 54 | The ticket form refuses a stage whose run had a dirty Go file and names the requirement. | planned | A form that ignores the dirty paths records an uncommitted run and the row reds. |
| RT62 | 54 | The ticket form accepts a stage whose run had only the spec's record path dirty. | planned | A form that refuses every dirty path refuses a run after a record write and the row reds. |
| RT63 | 37 | The ticket form's help lists the amendment, chunk, and verification steps in the order of `CompletionOrder`. | planned | A help text that types its own step list keeps the old order when the completion order changes, and the row's comparison with `CompletionOrder` reds. |
| RT64 | 47 | The review phase's step 6 states that `bench record excerpt` stages an excerpt and `bench record ticket` writes one ticket's entries. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A new anchor and its independent twin red when the sentence is absent. |
| RT65 | 8 | A response that takes the create-failure route clears the retained response. | planned | A retention that treats the passing state as complete keeps the earlier response and the row reds. |

Not covered: story 50 — the changelog entry is review-owned.
Not covered: story 51 — the one-session-one-ticket rule stays unchanged, and `TestDelegatedDistinctTicketAuthors` keeps grading it.
Not covered: story 52 — the verification form stays unchanged, and its existing tests keep grading it.

### Edge inventory

The audience is every repository that links the kit, because the forms and the retention ship in the binary.

- Paths with spaces: the stage keys by the slug, the chunk ID, and the requirement ID. The record grammar refuses control characters in each flag.
- Control bytes: the record command refuses a control character in a flag, and the excerpt bound refuses invalid UTF-8.
- Absent versus empty: RT11 covers an absent retained response, and the excerpt bound refuses an empty one.
- Symbolic links: the response store never follows a symbolic link at a store level, and the stage reads through the same no-follow classification.
- A command whose write changes a fact it reports: the ticket form writes only the record file, and every source digest excludes that file.
- State reloaded by a fresh process: each ticket row runs the form in a fresh `recordcmd.Command` call on saved bytes.
- A partial transaction: RT21 and RT22 cover a refusal after an earlier step.
- Concurrent calls in one worktree: the last retention write wins, and RT12 refuses a stage of the wrong command.
- A hand-edited stage: the ticket form imports a probe stage through the probe reader, and the verification writer computes each digest.

**Won't handle** one repair session for every ticket of its chunk — it breaks the one-session-one-ticket rule, and each repair session uses the ticket form.
**Won't handle** a retained-response read in `bench record verification` — the ticket form covers the author route, and final records keep their file form.
**Won't handle** retention of a tree-target or exempt call — neither runs through `runBounded`, and the author's checks are bounded calls.
**Won't handle** retention of more than one response — the stage copies each response before the next run, and the ticket form reads only the stages.
**Won't handle** a run binding on `bench record verification --excerpt <file>` — that form stays unchanged and keeps its file trust.
**Won't handle** a ticket form for final verification — the final list names no ticket, and the orchestrator keeps `bench record verification --final`.

## Ownership fences

The slicing writes the ticket entries of this section, and the review pickup is the one entry before the slicing.

- `reviews/record-ticket.md`

## Out of scope

- One repair session that records for every ticket of its chunk: 9 edits, 2 gate runs.
- A retained-response read in `bench record verification`: 5 edits, 2 gate runs.
- A history of more than one retained response: 7 edits, 2 gate runs.

## Further notes

### Flagged additions

- Staging compares the retained argument list with the planned command (RT12). A stage of the wrong response otherwise records evidence for another command.
- Staging and the ticket form refuse a probe answer that does not import (RT13, RT14). The later import otherwise fails inside the transaction.
- The form names every missing requirement in one refusal (RT23). The FT379 retries came from one gap at a time.
- The form removes its used stages after a write (RT35). A refreeze otherwise reuses an excerpt from an earlier tip.
- The base derives from the record (RT26, RT27, RT28). The FT379 build passed a wrong base once.
- The stage binds to its run commit and refuses a dirty run (RT60, RT61, RT62). The orchestrator decided this on 2026-10-14, and it stays open to reviewer veto.
- The retained response holds `HEAD` and the dirty paths, so each retained call runs one porcelain status read and one `rev-parse`.
- The deliverable binding drops the FT405 repair-session criterion, which narrows the commitment. The reviewer's 2026-10-14 scope choice deferred that change, and the drop stays open to reviewer veto.
- A spill-failed response clears the retained response (RT7). The incomplete bytes otherwise stage as evidence.

### Source sentence to row

| source sentence | rows |
|---|---|
| FT405: `bench record ticket` records the plan amendment when the digest moved. | RT19, RT20, RT29 |
| FT405: it freezes the chunk at `HEAD`. | RT18, RT24–RT28, RT30 |
| FT405: it records one verification entry for each planned requirement from a folder of excerpts named by requirement id. | RT18, RT23, RT31–RT36 |
| FT405: the implement-spec guidance then gives the freeze and the verification records to the ticket author. | RT47–RT52, RT54 |
| FT405: the verb or the verification form can read an excerpt from the last response of the worktree. | RT1–RT17 |
| FT405: `bench preflight evidence --check-current` reports the dirty paths instead of a refusal. | RT40–RT46 |
| FT405: one repair session records verification for every ticket of its chunk. | Won't handle; reviewer decision 2026-10-14 |
| reviewer decision 2026-10-14: no partial record write. | RT21, RT22 |

The excerpts come from staged retained responses in place of a folder in the worktree, because a folder in the worktree makes the checkout dirty.

### Pre-review proof checklist

- `Cited symbols`: `coveredAfter`, `shellcommand.Parse`, `LandedState`, and `ParsePorcelainZ` resolve in the tree. `Facts`, `write`, `replace`, `RecordAmendment`, `RecordChunk`, `RecordVerification`, `currentPlan`, `CompletionOrder`, `refusePlanned`, `importedProbe`, `runBounded`, `AssignmentScope`, `Drop`, `preparationCheckoutRefusal`, `currentEvidenceCommand`, `chargeevidence.Current`, and `ResponseBlocks` also resolve in the tree.
- `Import edges`: `cmd/bench` already imports `internal/responsebound`; `recordcmd` imports `internal/responsebound` for the stage. The slicing runs `go list` before it locks the edge.
- `Source-row clauses and occurrences`: the source-sentence table quotes each clause; the occurrences are FT349, NRO, FT217, and the FT379 freeze retries.
- `Promised field labels`: forms `ticket` and `excerpt`; flags `--ticket`, `--id-prefix`, `--requirement`, `--mutant`, `--base`, and `--tip`; response field `dirty_record`.
- `Changed-function callers`: the slicing runs `bench consumers` for `RecordAmendment`, `RecordChunk`, `RecordVerification`, `runBounded`, and `preparationCheckoutRefusal`.
- `Copy survival`: RT37 is review-owned, because a copied change body passes every behavior row.
- `Rendered-shape readers`: the `current` header needle in `evidence_modes_test.go` and the generated evidence reference change with the new field.
- `Pin operators`: the anchor rows use exact-needle `RequireInSection`, `Require`, and `Forbid` matches.
- `Entry reads`: none.
- `Derived expectations`: RT38 keeps the inventory test's own form list.
- `Consolidated rules`: the table under Ticket record gives each writer's old and new rule, and RT37 covers the changed cells.
- `Quantified obligations`: "every owed requirement" takes RT18, RT23, and RT33; "any other dirty path" takes RT42 to RT44.
- `Workflow-step writes`: the ticket form writes only the record file, and the stage and the retained response live outside the worktree.

### Paths-authorized settlement

`preparedAttempts` gathers the source snapshot paths, which include tracked and untracked working-tree paths (`internal/diff/range.go:243-258`, `internal/preflight/preparation.go:36`).
A dirty record path therefore reaches `ChangedPaths`, and `paths-authorized` grades it (`internal/preflight/paths_authorized.go:24-33`).
The implicit entries hold only the capture folder and the spec folder (`internal/preflight/fence_writes.go:71-77`).
So the record path is authorized only through the spec's fence, which lists the review pickup under the craft-spec fence rule.
This spec keeps that behavior and adds no implicit entry.

### Not verified

- The exact prose budget of the implementation phase file after the two in-place sentence changes.
- Whether any other reader parses the `current` block by column count, beyond the header needle in `evidence_modes_test.go`.
- The test name of the known-issues row check; RT53 stays planned until the slicing names it.
