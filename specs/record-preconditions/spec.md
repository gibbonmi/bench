# Record preconditions before an attempt

Status: staged

Roadmap: FT414

Decision source: Named reviewed artifact roadmap/FT414.md (drain-reviewed row; commitment outcome FT414-preconditions), with the reviewer decision of 2026-11-08 on the bound request token.

Persisted fields: intent-ledger:commitment.bindings.request_token

Verification log: pending

## Problem

An agent learns the rules of the chunk checkpoint, the record verbs, and the landing only from refusals, one at a time.
No read-only verb states the required chunk base, the stale or missing record entries, or the commits that break the review chain.
Orchestrators read the one-line `bench-completion-plan` JSON with `python3` many times in each build.

Authors build probe-result files by hand from spill frames, because `bench probe` and `bench test` write no complete response to a caller file.
A record form that gets two exclusive flags prints only its usage line, so the agent must find the conflict.
`bench preflight build --plan-only` prints only the red rows, so the orchestrator rebuilds the row table by hand.
A landing refusal after a request change names no bound request, because the delivery binding stores only the request digest.

## Solution

Two read-only record forms state the preconditions first.
`bench record status <slug>` prints the whole completion plan and the record state of each chunk as flat TOON tables.
With `--chunk <id>`, it prints the plan rows of that chunk, its required base, each stale or missing entry, and each commit that breaks the chain.
`bench record next <slug>` prints the next required record or plan action as an exact command, with its preconditions.
Both forms derive every finding from the predicates that the chunk checkpoint enforces, so the view and the checkpoint cannot disagree.

`bench probe` and `bench test` take `--out <file>`, which writes the complete response that `bench record verification --probe-result` accepts.
A record form that gets two flags of one alternation group names both flags.
`bench preflight build --plan-only` prints every check row.
`bench commitment start` keeps the plain request token in the delivery binding, and the landing refusals name that token.

## User stories

Line: opus / xhigh.
Implementation-line reason: RP-C1 changes the checkpoint grader into a collector, and RP-C6 crosses the intent, commitment, and worktree owners with a persisted field. The other chunks are exact changes at known seams with covering package tests, so their authors run at medium effort.
Harder chunks: RP-C1, RP-C6.

### The whole plan view

1. As an orchestrator, I want `bench record status <slug>` to print the completion plan as flat TOON tables, so that I parse no JSON with `python3`.
2. As an orchestrator, I want the view to list each chunk with its tickets and requirements, so that I see the chunk graph in one read.
3. As an orchestrator, I want the view to list each ticket with its chunk and its covered rows, so that I map a ticket to its acceptance rows.
4. As an orchestrator, I want the view to list each recorded assignment, so that I read the dispatch history without the spec fence.
5. As an orchestrator, I want the record state of each chunk entry, so that I know which chunks the record holds and which are stale.
6. As an orchestrator, I want the state of each chunk requirement, so that I know which verification each chunk still owes.
7. As an orchestrator, I want the state of each final requirement, so that I know what the completion still owes.
8. As an agent, I want a long view in byte pages with a source identity, so that I read it without a spill file.
9. As an agent, I want a definitive empty state when no record exists, so that an absent record does not look like an error.
10. As an agent, I want the read-only forms to run in any checkout, the primary checkout included, so that a read needs no worktree.
11. As a reviewer, I want the read-only forms to write nothing, so that a read cannot change the record or the plan digest.

### The chunk view

12. As a ticket author, I want the required base of a chunk before I record it, so that `bench record chunk` does not refuse me at the checkpoint.
13. As a ticket author, I want the plan rows of one chunk, so that I see its tickets, requirements, and acceptance rows in one read.
14. As a ticket author, I want each stale or missing entry of the chunk, so that I repair all of them in one pass.
15. As a ticket author, I want each finding to name the form that repairs it, so that I run the exact repair command.
16. As a ticket author, I want each commit that breaks the review chain, so that I know which change needs a fold or a review.
17. As a ticket author, I want a proved comment-only Go commit to stay off the chain list, so that the list names only real breaks.
18. As a ticket author, I want an unknown chunk ID refused with the plan's chunk IDs, so that I correct the call at once.
19. As a reviewer, I want the chunk view and the chunk checkpoint to read one predicate source, so that the view never contradicts a refusal.

### The next action

20. As an orchestrator, I want `bench record next <slug>` to print the next required record or plan action as an exact command, so that I find the order without trial.
21. As an orchestrator, I want the preconditions of that action with their states, so that I fix a failed precondition before the attempt.
22. As an orchestrator, I want a stale plan digest to make the amendment the next action, so that no chunk write refuses on the digest.
23. As an orchestrator, I want a definitive empty state when the record owes nothing, so that I know the completion can run.
24. As a reviewer, I want the next action to come from the checkpoint's own repair, so that the verb never advises an action that the checkpoint refuses.

### The AXI surface

25. As an agent, I want the two read-only forms in the approved AXI query set, so that they keep the query output contract with help rows.

### Exclusive flags

26. As a ticket author, I want a record form that gets two exclusive flags to name both flags, so that I see the conflict without a second attempt.
27. As a ticket author, I want each alternation group of each record form to name its conflict, so that no group falls back to the bare usage line.
28. As a ticket author, I want the conflict refusal to keep exit 2 and the usage line, so that the grammar contract stays the same.
29. As a ticket author, I want an incomplete optional group to keep the plain usage line, so that only a real conflict names a conflict.
47. As a ticket author, I want a call with several conflicting pairs to name one fixed pair, so that the same call always prints the same refusal.

### The complete response file

30. As a ticket author, I want `bench probe --out <file>` to write the complete probe response, so that `--probe-result` imports it with no edit.
31. As a ticket author, I want `bench test --out <file>` to write the complete test response, so that a record excerpt holds native bytes.
32. As an agent, I want stdout to keep its bounded response with `--out`, so that the call costs no more context.
33. As a ticket author, I want an existing `--out` path refused before the verb runs, so that no evidence file is overwritten.
34. As a ticket author, I want a symbolic link at the `--out` path refused, so that the write cannot leave the named path.
35. As a ticket author, I want a relative `--out` path to resolve against my working directory, also under `--in`, so that the file lands where I named it.
36. As a ticket author, I want a named refusal when the response has no complete copy, so that a partial file never looks complete.
37. As an agent, I want the help rows of both verbs to show `--out`, so that I discover the option from `bench help`.

### The plan-only row table

38. As an orchestrator, I want `bench preflight build --plan-only` to print every check row, so that I do not rebuild the table by hand.
39. As an orchestrator, I want the ordinary build verdict to keep its red-only table, so that its help stays true.
40. As an orchestrator, I want the plan-only exit codes to stay the same, so that the spec-stage close reads the same verdict.

### The bound request token

41. As an orchestrator, I want `bench commitment start` to keep the plain request token in the delivery binding, so that a refusal can name it.
42. As an orchestrator, I want the landing refusal of an assignment whose binding holds another request to name the bound token, so that I land with the bound token.
43. As an orchestrator, I want the unknown-request refusal to name the bound token of the candidate assignment, so that I do not reauthorize and orphan the binding.
44. As an orchestrator, I want a binding written before the field to name its request digest and say the token is unknown, so that an old ledger still reads.
45. As a reviewer, I want the ledger to refuse a binding whose token does not match its request digest, so that a refusal never names a false token.
46. As a reviewer, I want the light-path fallback of an unbound commit to stay the same, so that the new sentence changes no admission.
48. As a reviewer, I want a sibling binding to carry no copied token, so that every ledger read stays valid and no refusal names the token of another request.

## Implementation decisions

### One finding collector for the checkpoint and the views

Package `reviewrecord` gains a collector that walks the same predicates as the chunk checkpoint and yields every finding.
Each finding carries the chunk ID, the entry, the state (`missing` or `stale`), the error, and the `Repair` value when a form repairs it.
`checkSource` keeps its signature and returns the error of the first finding, so every checkpoint caller keeps its refusal text and order.
Within one chunk, the identity group stops at its first failure, because the later identity predicates read the earlier results.
The identity group is ancestry, source digest, plan identity, and acceptance rows.
The collector then evaluates the verification, review, chain, and mapping groups of that chunk, and it continues with the next chunk.

The required base of a chunk comes from `chunkBase` and `coveredAfter`, the derivation that `bench record ticket` uses today.
A chain gap lists each commit on the gap range whose source digest does not follow its first parent's digest under `digestFollows`.
An uncovered reviewed source lists each such commit on the range from the covered tip to the requested tip.
A default-branch merge that no recorded fold covers takes the reason `unrecorded fold`, from the `uncoveredFold` helper.
No other code derives a chain, staleness, or base predicate. A second derivation is a defect under the one-source standard.

### The read-only record forms

The `bench record` family gains the forms `status` and `next`. Their form-name constants live in package `reviewrecord` beside the other form names.
The constant `"next"` cannot live in `recordcmd`, because RR49 reds a bare `next` literal in a scanned write-verb package.
The read-only forms skip the primary-checkout refusal and read the tree that the call targets.
Every other record form keeps that refusal.
Each read-only form reads the record and the plan at HEAD of its tree, and it writes no file and no ledger entry.

The grammar is `bench record status <slug> [--chunk <id>] [--offset <byte-offset> --source <identity>]` and `bench record next <slug>`.
The status view pages its projection through `responsebound.PageRead`, as `bench worktree list --spec` does.
A read exits 0 when it prints findings, because a stale entry is a state and not a failure.
A read exits 1 for an invalid record, an unreadable plan, or an unknown chunk, and it exits 2 for a grammar error.

The output tables are fixed:

| call | table | fields |
|---|---|---|
| status | `plan[1]` | `version`, `plan_digest`, `record` |
| status | `chunks[N]` | `chunk`, `tickets`, `requirements`, `entry`, `base`, `tip` |
| status | `tickets[N]` | `ticket`, `chunk`, `covers` |
| status | `verification[N]` | `chunk`, `requirement`, `command`, `probe`, `ticket`, `state` |
| status | `final_verification[N]` | `requirement`, `command`, `state` |
| status | `assignments[N]` | `ticket`, `session`, `assignment`, `model`, `effort`, `source`, `trigger` |
| status --chunk | `chunk[1]` | `chunk`, `required_base`, `base`, `tip`, `entry` |
| status --chunk | `tickets[N]` | `ticket`, `covers` |
| status --chunk | `verification[N]` | `requirement`, `command`, `probe`, `ticket`, `state` |
| status --chunk | `findings[N]` | `entry`, `state`, `detail`, `repair` |
| status --chunk | `chain[N]` | `commit`, `range`, `reason` |
| next | `action[N]` | `form`, `command` |
| next | `preconditions[N]` | `precondition`, `state` |

The `record` cell of `plan[1]` holds `absent`, `current`, or `stale`.
Each `entry` and each `state` cell holds `missing`, `current`, or `stale`, and each precondition `state` holds `holds` or `fails`.
Every read ends with the shared `help[N]{cmd,why}:` block that `axi.RenderHelp` emits.

### The next action

`bench record next` selects one checkpoint mode at HEAD.
A record entry covers a planned chunk directly or through the amendment mapping of `mappedIDs`.
If a planned chunk has no covering record entry, it runs the collector in chunk mode on the first such chunk in plan order.
If the record covers every planned chunk, it runs the collector in completion mode.
It takes the first finding in checkpoint order, and it never skips a finding.
A first finding with a `Repair` prints that form's command.
The command cell is the repair form's usage line with the known values, from the same `suffix(known(...))` composition that a record refusal prints today.
A first finding with no `Repair` prints `action[0]` and the finding in `preconditions` with the state `fails`.
An absent record file prints the `record chunk` command of the first planned chunk.
The collector seeds an absent record through `createRecord`, as the chunk write does, so the seed starts at the current plan digest.
A record that owes nothing prints `action[0]` and `preconditions[0]`.

Each precondition row comes from a predicate that the action form's write path calls.
The rows are the plan digest from `currentPlan`, the required base from `chunkBase`, and the chunk entry from `ChunkEntry`.
A form shows only the preconditions that its own write path reads.

### AXI approval

The `record` registry row changes from the mutation exemption to approved children.
Each approved child joins `approvedAXIQueries`, the craft-cli approved-query table, the AXI query list in `projects/benchkit.md`, and the AXI envelope cases of `cmd/bench`.
The status form joins with its producer, and the next form joins with its own producer.
The reviewer decided on 2026-11-08 that the craft-cli skill takes a per-file prose budget row of 122 lines in `projects/benchkit.md`.

### Exclusive flag conflicts

The record family derives a conflict from the form layout that `admits` reads.
Two present flags conflict when no admitted flag set holds both of them.
The refusal prints `error: conflicting flags — <first> and <second> exclude each other`, then the form's usage line, at exit 2.
A flag set that fails `together` with no conflicting pair keeps the bare usage line, and so does a `valid` failure.
The layout is the one source of the pairs, so no code lists them by hand.
Today the layout gives ten pairs: nine in the verification form and one in the review form.
The verification pairs are `--chunk` with `--final`, and each of `--excerpt`, `--probe-outcome`, `--probe-exit-code`, and `--probe-restore` with each of `--probe-result` and `--mutant`.
The review pair is `--chunk` with `--fold`.
When a call holds several conflicting pairs, the refusal names the first pair in the order of flag declaration.
The conflict helper lives in `refusal.go`, because `command.go` has too little line headroom for it.

### The complete response file

The command registry gains one capability field for the `--out` option, and only the `probe` and `test` rows set it.
The dispatcher reads `--out <file>` as the first argument after the verb, or as the first argument after `--in <label>`.
The help insertion reads the same field, as the tree-target insertion reads the scope field.
The verbs' own grammars do not change.

The dispatcher resolves a relative path against the caller's working directory before the tree-target route, and it passes the absolute path to the child.
Before the verb runs, the dispatcher refuses at exit 1 when the path exists, when its parent directory is absent, or when the path holds a control byte.
After the owner finishes, the dispatcher writes the owner's complete output with `O_EXCL` and `O_NOFOLLOW`, at mode 0644.
The complete output is the one that the response retention and the spill file hold, from one exported owner accessor.
When no complete output exists, the dispatcher writes no file and exits 1.
It then prints `error: incomplete response — --out wrote no file, because the response has no complete copy` after the response.
A write failure after a zero verb exit also exits 1. A nonzero verb exit stays the verb's exit.
Stdout keeps its bounded response, so the closed bound-exemption set does not change.

### The plan-only row table

`renderChecks` takes the operation kind. The plan-only kind prints every check row in check order, and every other kind keeps the red-only table.
The count line and the exit rule stay the same in both modes.
The plan-only help description states that it prints every check row.
A full table holds more lines than the public response bound, so the response spills, and the spill file holds the complete table.
The response bound does not change.

### The bound request token

`DeliveryBinding` gains `request_token`, which `Store.Start` writes from its `--request` value.
`ValidateCommitment` refuses a binding whose present token is not line-safe, or whose `RequestDigest` differs from the binding's request digest.
A binding without the field stays valid.
`RegisterSibling` copies a binding to a sibling with another request digest, so the copy clears `request_token`.
`DATA_HANDLING.md` gains a row for `commitment.bindings.request_token`, because the commitment-record conformance check derives every ledger field.
The row states that the token is the opaque request label that `bench worktree list` already prints.

One `intent` function answers the bound request of an assignment from the ledger. Both landing refusals read it.
It is a function and not a method, because `DeliveryBinding` is a type alias in package `intent`.
When `readyFor` finds a binding of the owner assignment with another request digest, its unbound error adds one sentence.
With a token, the sentence is `; the delivery binding of assignment <id> holds the request token <token>`.
Without a token, the sentence is `; the delivery binding of assignment <id> holds the request digest <digest>, and its token is unknown`.
The error still wraps `errUnbound`, so the light-path fallback does not change.
On the light path, the production-path clause comes before the binding sentence.
The full text is then `assignment has no current delivery binding for production path "<path>"; the delivery binding of assignment <id> holds the request token <token>`.
The unknown-request refusal of a landing can name one candidate assignment in its `observed` field.
When that assignment holds a binding, the landing caller adds the same sentence to the refusal detail. Its route does not change.
The shared identity component does not change, so the release refusal keeps its text.

The field makes this spec a declared build. Each record verb of the build therefore runs on an installed candidate.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RP-C1 / 01-collect-record-findings.md, 02-show-the-whole-plan-status.md, 03-show-one-chunk-status.md | The collector, and `bench record status` with the whole plan and the chunk view | RP1–RP22, RP28 | collector differential, status view, chain list, AXI registry | yes |
| RP-C2 / 04-print-the-next-record-action.md | `bench record next` with its preconditions | RP23–RP27, RP51–RP53, RP58 | next action and preconditions, AXI registry | no |
| RP-C3 / 05-name-exclusive-flag-conflicts.md | Named conflicts for each alternation group | RP29–RP32, RP54, RP55, RP57 | conflict refusals | no |
| RP-C4 / 06-write-the-complete-response.md | `--out` on `bench probe` and `bench test` | RP33–RP41 | out file and import round trip | no |
| RP-C5 / 07-print-every-plan-only-row.md | The full plan-only row table | RP42–RP44 | preflight render | no |
| RP-C6 / 08-name-the-bound-request-token.md | The persisted bound token and the two landing refusals | RP45–RP50, RP56 | ledger, admission, sibling, and landing refusals | yes |

`CHANGELOG.md` gains one entry for each chunk, from the last ticket of that chunk.
The build runs the chunks in table order on one integration source.
RP-C2 waits for RP-C1, because it consumes the collector, the status form, and the AXI child approval.
RP-C3 shares `command.go` and `refusal.go` with RP-C1 and RP-C2, and RP-C4 and RP-C5 share the exact help inventory, so these chunks do not author at the same time.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"RP-C1","tickets":["01-collect-record-findings.md","02-show-the-whole-plan-status.md","03-show-one-chunk-status.md"],"verification":[{"id":"collector","command":"bench test --package ./internal/reviewrecord","probe":"return only the first finding from the collector"},{"id":"status","command":"bench test --package ./internal/reviewrecord/recordcmd","probe":"omit the chain table rows"},{"id":"axi","command":"bench test --check axi-query-registry","probe":"drop status from the approved record children"}]},{"id":"RP-C2","tickets":["04-print-the-next-record-action.md"],"verification":[{"id":"next","command":"bench test --package ./internal/reviewrecord/recordcmd","probe":"skip the plan digest precondition"},{"id":"axi-next","command":"bench test --check axi-query-registry","probe":"drop next from the approved record children"}]},{"id":"RP-C3","tickets":["05-name-exclusive-flag-conflicts.md"],"verification":[{"id":"conflicts","command":"bench test --package ./internal/reviewrecord/recordcmd","probe":"print the bare usage line for a conflicting pair"}]},{"id":"RP-C4","tickets":["06-write-the-complete-response.md"],"verification":[{"id":"out-file","command":"bench test --package ./cmd/bench","probe":"write the projected stdout in place of the complete output"},{"id":"owner","command":"bench test --package ./internal/responsebound","probe":"return the head lines from the complete output accessor"}]},{"id":"RP-C5","tickets":["07-print-every-plan-only-row.md"],"verification":[{"id":"plan-only","command":"bench test --package ./internal/preflight","probe":"keep the red filter for the plan-only kind"}]},{"id":"RP-C6","tickets":["08-name-the-bound-request-token.md"],"verification":[{"id":"ledger","command":"bench test --package ./internal/intent/...","probe":"accept a token whose digest differs"},{"id":"admission","command":"bench test --package ./internal/commitment/...","probe":"omit the token write in Start"},{"id":"landing","command":"bench test --package ./internal/worktree","probe":"omit the bound token sentence from the unknown-request detail"},{"id":"sibling","command":"bench test --package ./internal/worktree","probe":"keep the source token in the sibling copy"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/record-preconditions/spec.md"},{"id":"packages","command":"bench test --package ./..."},{"id":"axi","command":"bench test --check axi-query-registry"},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

- A good test drives the public verb or the exported owner and observes printed TOON, exit codes, and file bytes.
- The record views use the `recordtest` fixture harness: `NewLinked`, `Prepare`, `RecordChunk`, `Verification`, and `RecordFold`.
- The collector keeps the existing checkpoint tests green with no edit, and one differential test compares `CheckSource` with the first collected finding.
- The `--out` rows run the real dispatcher through `Command.Run` and import the file through `bench record verification --probe-result`.
- The landing rows reuse the reauthorization fixtures in `internal/worktree` and the publication fixtures in `internal/commitment/repository`.
- The AXI rows run the `axi-query-registry` conformance check.

### Seam diagram

    trigger: an orchestrator or a ticket author before a record write
        │
        ▼
    plan + record at HEAD ──▶ [ reviewrecord collector ] ──▶ findings ──▶ status / next / checkpoint
                                  ◀ tests attach here: recordtest fixtures; compare the printed tables and CheckSource's first error

    trigger: bench probe|test --out <file>
        │
        ▼
    verb output ──▶ [ responsebound owner ] ──▶ bounded stdout + complete file ──▶ record verification --probe-result
                        ◀ tests attach here: Command.Run in a fixture repo; read the file bytes and import them

    trigger: bench commitment start, then bench worktree land
        │
        ▼
    --request token ──▶ [ ledger binding ] ──▶ readyFor / identity refusal ──▶ refusal text
                            ◀ tests attach here: ledger fixtures with a reauthorized assignment

### Acceptance coverage map

All seams marked planned are future executable evidence. No implementation red is claimed in this planning phase.

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| RP1 | 1, 2 | The status view prints one `chunks` row per planned chunk, with its tickets and requirement IDs. | planned TestRecordStatusPrintsThePlanChunks in internal/reviewrecord/recordcmd/status_test.go | A view that skips a chunk or a requirement prints fewer rows than the fixture plan holds. |
| RP2 | 3 | The status view prints one `tickets` row per planned ticket, with its chunk and covered rows. | planned TestRecordStatusPrintsTheTickets in internal/reviewrecord/recordcmd/status_test.go | A ticket with no row or a wrong chunk cell fails the exact table comparison. |
| RP3 | 4 | The status view prints one `assignments` row per recorded assignment of a version 2 plan. | planned TestRecordStatusPrintsTheAssignments in internal/reviewrecord/recordcmd/status_test.go | A view that reads only the version 1 keys prints `assignments[0]` for a dispatched fixture. |
| RP4 | 5 | A chunk entry whose tip tree digest differs from its source digest shows `stale` in the `entry` cell. | planned TestRecordStatusMarksAStaleChunk in internal/reviewrecord/recordcmd/status_test.go | A view that reports presence alone prints `current` for the edited fixture. |
| RP5 | 5 | A planned chunk with no record entry shows `missing` in the `entry` cell. | planned TestRecordStatusMarksAMissingChunk in internal/reviewrecord/recordcmd/status_test.go | A view that lists only recorded chunks omits the row. |
| RP6 | 6 | A chunk requirement with no verification entry shows `missing` in its `verification` row. | planned TestRecordStatusMarksAMissingRequirement in internal/reviewrecord/recordcmd/status_test.go | A view that prints the plan without the record state prints an empty `state` cell. |
| RP7 | 6 | A verification entry that `EntryCurrent` rejects shows `stale` in its `verification` row. | planned TestRecordStatusMarksAStaleRequirement in internal/reviewrecord/recordcmd/status_test.go | A view that checks only presence prints `current` for the stale fixture entry. |
| RP8 | 7 | The status view prints one `final_verification` row per final requirement with its state. | planned TestRecordStatusPrintsTheFinalList in internal/reviewrecord/recordcmd/status_test.go | A view that omits the final list prints no `final_verification` table. |
| RP9 | 8 | A status projection larger than the response bound prints its first page and a continuation cursor. | planned TestRecordStatusPagesALongPlan in internal/reviewrecord/recordcmd/status_test.go | An unpaged view spills, and the test finds a `spilled{` frame in place of the cursor. |
| RP10 | 8 | A continued page with a stale source identity refuses with `source changed since the cursor`. | planned TestRecordStatusRefusesAStaleCursor in internal/reviewrecord/recordcmd/status_test.go | A view that ignores the identity prints bytes of the changed projection. |
| RP11 | 9 | With no record file, the status view prints `record` as `absent` and every chunk `entry` as `missing`, at exit 0. | planned TestRecordStatusWithNoRecord in internal/reviewrecord/recordcmd/status_test.go | A view that reads the record through the write path refuses with `holds no record` at exit 1. |
| RP12 | 10 | `bench record status` in a primary checkout prints the view at exit 0. | planned TestRecordStatusRunsInThePrimaryCheckout in internal/reviewrecord/recordcmd/status_test.go | A form that keeps the primary-checkout refusal prints that refusal at exit 1. |
| RP13 | 10 | `bench record chunk` in a primary checkout still prints the primary-checkout refusal. | `internal/reviewrecord/recordcmd/refusal_test.go` (`TestRecordRefusesThePrimaryCheckout`) | A change that drops the refusal for every form lets a write run in the primary checkout. |
| RP14 | 11 | A status read leaves the record bytes, the spec bytes, and `git status` unchanged. | planned TestRecordReadsWriteNothing in internal/reviewrecord/recordcmd/status_test.go | A read that routes through the write transaction rewrites the record and changes its bytes. |
| RP15 | 12 | The chunk view prints the `required_base` that `chunkBase` derives from the predecessor's covered tip. | planned TestRecordStatusNamesTheRequiredBase in internal/reviewrecord/recordcmd/status_chunk_test.go | A view that prints the stored base prints the wrong commit for a chunk whose predecessor has a recorded fold. |
| RP16 | 13 | The chunk view prints only the tickets, requirements, and covered rows of the named chunk. | planned TestRecordStatusChunkPrintsItsPlanRows in internal/reviewrecord/recordcmd/status_chunk_test.go | A view that prints the whole plan for `--chunk` prints rows of a second fixture chunk. |
| RP17 | 14, 15 | With two stale entries in one chunk, the `findings` table holds both, each with its repair command. | planned TestRecordStatusListsEveryFinding in internal/reviewrecord/recordcmd/status_chunk_test.go | A view built on the first-error checker prints one finding for the two-fault fixture. |
| RP18 | 16 | A non-comment commit between the predecessor tip and the chunk base appears in the `chain` table. | planned TestRecordStatusListsAChainBreak in internal/reviewrecord/recordcmd/status_chunk_test.go | A view that prints only the gap sentence names no commit hash. |
| RP19 | 16 | An unrecorded default-branch merge on the uncovered range appears with the reason `unrecorded fold`. | planned TestRecordStatusNamesAnUnrecordedFold in internal/reviewrecord/recordcmd/status_chunk_test.go | A view that treats the merge as an ordinary commit prints another reason. |
| RP20 | 17 | A proved comment-only Go commit on the gap range does not appear in the `chain` table. | planned TestRecordStatusSkipsACommentOnlyCommit in internal/reviewrecord/recordcmd/status_chunk_test.go | A view that compares raw digests lists the comment-only commit. |
| RP21 | 18 | `--chunk` with an unplanned ID refuses at exit 1 and names each planned chunk ID. | planned TestRecordStatusRefusesAnUnknownChunk in internal/reviewrecord/recordcmd/status_chunk_test.go | A view that prints an empty chunk accepts the typo at exit 0. |
| RP22 | 19, 24 | For each refusal fixture of the reviewrecord package, the first collected finding's error equals the `CheckSource` error. | planned TestCollectFirstFindingIsCheckSource in internal/reviewrecord/collect_test.go | A collector with its own copy of a predicate gives a different first finding on one fixture. |
| RP23 | 20 | With a recorded first chunk and an unrecorded second chunk, `bench record next` prints the `record chunk` command for the second chunk. | planned TestRecordNextNamesTheNextChunk in internal/reviewrecord/recordcmd/next_test.go | A verb that walks the forms in a fixed order without the record prints another form. |
| RP24 | 20 | The `command` cell equals the repair command that a `bench record chunk` refusal prints for the same fixture. | planned TestRecordNextCommandMatchesTheRefusal in internal/reviewrecord/recordcmd/next_test.go | A second command composition drifts from the refusal text on the slug or a flag value. |
| RP25 | 21 | The next chunk action lists the `required base` precondition with the `chunkBase` commit. | planned TestRecordNextListsTheChunkPreconditions in internal/reviewrecord/recordcmd/next_test.go | A verb that prints the action alone prints `preconditions[0]`. |
| RP26 | 22 | After a plan commit, `bench record next` prints the `record amendment` command and a `fails` plan digest precondition. | planned TestRecordNextPutsTheAmendmentFirst in internal/reviewrecord/recordcmd/next_test.go | A verb that skips `currentPlan` advises the chunk write that refuses on the stale digest. |
| RP27 | 23 | A record that owes nothing prints `action[0]` and `preconditions[0]` at exit 0. | planned TestRecordNextWithNothingOwed in internal/reviewrecord/recordcmd/next_test.go | A verb that always proposes a form prints a completion command for a completed record. |
| RP51 | 24 | In chunk mode on the second chunk, a first finding with no `Repair` prints `action[0]` and that finding as a `fails` precondition. | planned TestRecordNextStopsAtAFindingWithNoRepair in internal/reviewrecord/recordcmd/next_test.go | In a fixture with a stale first-chunk source digest and an unrecorded second chunk, a verb that skips to a repair advises the refused chunk write. |
| RP58 | 11 | A next read leaves the record bytes, the spec bytes, and `git status` unchanged. | planned TestRecordNextWritesNothing in internal/reviewrecord/recordcmd/next_test.go | A next read that seeds the absent record through the write transaction creates the record file. |
| RP52 | 20 | With no record file, `bench record next` prints the `record chunk` command of the first planned chunk. | planned TestRecordNextWithNoRecordNamesTheFirstChunk in internal/reviewrecord/recordcmd/next_test.go | A verb that reads an absent record as an empty plan digest advises the amendment, which refuses on an absent record. |
| RP28 | 25 | The `record` row declares the approved child `status`, and the AXI registry, guidance table, and profile list agree. | `bench test --check axi-query-registry` | A registry change without the guidance or profile edit reds the check. |
| RP53 | 25 | The `record` row declares the approved child `next`, and the AXI registry, guidance table, and profile list agree. | `bench test --check axi-query-registry` | A `next` form with no approval, or an approval with no guidance row, reds the check. |
| RP29 | 26, 27 | `--excerpt` with `--probe-result` refuses at exit 2 and names both flags. | planned TestRecordConflictNamesExcerptAndProbeResult in internal/reviewrecord/recordcmd/conflict_test.go | The bare usage line holds no `conflicting flags` text. |
| RP30 | 27 | Verification `--chunk` with `--final` refuses at exit 2 and names both flags. | planned TestRecordConflictNamesChunkAndFinal in internal/reviewrecord/recordcmd/conflict_test.go | A conflict list written for one form misses this group. |
| RP31 | 27 | Review `--chunk` with `--fold` refuses at exit 2 and names both flags. | planned TestRecordConflictNamesChunkAndFold in internal/reviewrecord/recordcmd/conflict_test.go | A conflict list written for one form misses this group. |
| RP32 | 28, 29 | An assignment call with `--predecessor` and no `--trigger` prints the bare usage line at exit 2 with no conflict text. | planned TestRecordIncompleteGroupKeepsTheUsageLine in internal/reviewrecord/recordcmd/conflict_test.go | A rule that calls every `together` failure a conflict names a false pair. |
| RP54 | 27 | Verification `--probe-outcome` with `--probe-result` refuses at exit 2 and names both flags. | planned TestRecordConflictNamesAScalarAndProbeResult in internal/reviewrecord/recordcmd/conflict_test.go | A hand list of the three named pairs misses this layout pair and prints the bare usage line. |
| RP55 | 27 | Verification `--excerpt` with `--mutant` refuses at exit 2 and names both flags. | planned TestRecordConflictNamesExcerptAndMutant in internal/reviewrecord/recordcmd/conflict_test.go | A rule that reads only the alternation and not the `--mutant` admission misses this pair. |
| RP57 | 47 | A verification call with `--chunk`, `--final`, `--excerpt`, and `--probe-result` names `--chunk` and `--final`. | planned TestRecordConflictNamesTheFirstPair in internal/reviewrecord/recordcmd/conflict_test.go | A rule that walks a map in random order names a different pair across runs. |
| RP33 | 30 | A spilled `bench probe --out <file>` call writes a file that `bench record verification --probe-result` imports at exit 0. | planned TestProbeOutFileImports in cmd/bench/out_file_test.go | A file that holds the projected stdout carries a `spilled{` frame, and the import refuses it. |
| RP34 | 31 | A spilled `bench test --out <file>` call writes the same bytes as its spill file. | planned TestTestOutFileHoldsTheCompleteResponse in cmd/bench/out_file_test.go | A file that holds the head and tail lines differs from the spill file. |
| RP35 | 32 | With `--out`, stdout still holds the bounded projection with its `spilled{` frame. | planned TestOutKeepsTheBoundedStdout in cmd/bench/out_file_test.go | A design that exempts the call prints the complete response on stdout. |
| RP36 | 33 | An existing `--out` path refuses at exit 1 before the verb runs, and the file keeps its bytes. | planned TestOutRefusesAnExistingPath in cmd/bench/out_file_test.go | A write with truncation replaces the planted bytes, and a late check runs the probe first. |
| RP37 | 34 | A symbolic link at the `--out` path refuses, and the link target keeps its bytes. | planned TestOutRefusesASymlink in cmd/bench/out_file_test.go | A write without `O_NOFOLLOW` writes through the link into the target. |
| RP38 | 35 | Under `--in <label>`, a relative `--out` path writes the file under the caller's working directory, not in the target tree. | planned TestOutResolvesAgainstTheCallerDirectory in cmd/bench/out_file_test.go | A child that resolves the path in the target tree writes an untracked file there. |
| RP39 | 36 | When the owner holds no complete output, `--out` writes no file, prints `error: incomplete response`, and exits 1. | planned TestOutRefusesAnIncompleteResponse in cmd/bench/out_file_test.go | A write of the partial bytes leaves a file that looks complete. |
| RP40 | 37 | The `bench help` rows of `probe` and `test` show `[--out <file>]`, and no other row does. | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A hand-edited help row drifts from the registry field and fails the exact inventory. |
| RP41 | 30, 31 | The complete output accessor returns the same bytes for retention and for `--out`. | planned TestCompleteOutputHasOneSource in internal/responsebound/owner_test.go | A second accessor that reads only the in-memory lines loses the spilled bytes. |
| RP42 | 38 | A green plan-only verdict prints a `checks` table with one row per check, green and not-applicable rows included. | planned TestPlanOnlyPrintsEveryRow in internal/preflight/verdict_summary_test.go | The red-only filter prints no table for a green verdict. |
| RP43 | 39 | A green ordinary build verdict prints only the count line. | `internal/preflight/verdict_summary_test.go` (`TestPreflightGreenSummaryLine`) | A change that removes the filter for every kind prints a table here. |
| RP44 | 40 | A plan-only verdict exits 0 when green and 1 when a row is red. | `internal/preflight/commitment_test.go` (`TestCommitmentPlanOnly`) | A render change that alters the exit rule breaks the existing exit assertions. |
| RP45 | 41 | After `bench commitment start --request <token>`, the ledger binding holds `request_token` equal to the token. | planned TestStartKeepsTheRequestToken in internal/commitment/repository/admission_token_test.go | A Start that writes only the digest leaves the field empty. |
| RP46 | 42 | A landing after a reauthorize to a new token refuses with `holds the request token <original token>`. | planned TestLandingNamesTheBoundToken in internal/commitment/repository/admission_token_test.go | The current unbound text names no token, and the test finds no token needle. |
| RP47 | 42, 46 | The refusal of RP46 still satisfies `errors.Is` with the unbound error. | planned TestBoundTokenSentenceKeepsTheUnboundError in internal/commitment/repository/admission_token_test.go | A new error kind breaks the light-path fallback that reads `errUnbound`. |
| RP48 | 43 | An unknown-request landing whose one candidate assignment holds a binding names the bound token in its detail. | planned TestUnknownRequestNamesTheBoundToken in internal/worktree/land_reauthorization_test.go | The current detail is `request token matches no assignment`, and only the `observed` field names `assignment:<id>`. |
| RP49 | 44 | A binding with no `request_token` refuses with `holds the request digest <digest>, and its token is unknown`. | planned TestLegacyBindingNamesItsDigest in internal/commitment/repository/admission_token_test.go | A refusal that prints an empty token prints no digest. |
| RP50 | 45 | `ValidateCommitment` refuses a binding whose `request_token` digest differs from its `request`. | planned TestBindingTokenMustMatchItsDigest in internal/intent/ledger/commitment_test.go | A validator that checks only line safety accepts the mismatched fixture. |
| RP56 | 48 | A sibling binding that `RegisterSibling` writes carries no `request_token`, and the ledger reads back valid. | `internal/worktree/tree_target_test.go` (`TestCommitmentSiblingBinding`) | A copy that keeps the source token fails the RP50 digest check, and the ledger read refuses. |

### Edge inventory

- An empty plan chunk list is not reachable, because the plan parser refuses a plan with no final verification and every build plan holds chunks.
- A version 1 plan has no `execution`, so the `assignments` table prints `assignments[0]` (RP3 covers the version 2 side).
- A record with an invalid fence refuses at exit 1 with the existing record error line, because the read shares `ReadTree`.
- An `--out` parent directory that does not exist refuses before the verb runs; the RP36 test file also holds this case.
- An `--out` path with a control byte refuses before the verb runs; the RP36 test file also holds this case.
- A `--chunk` value with a control character keeps the existing single-line refusal of the record family.

**Won't handle** lines:

- The tip-current `next` cell of `bench preflight build` — `TestDecideOtherRedsCarryNoNext` pins it empty under WF36, and the plan-only table of RP42 still prints the row and its detail.
- The `bench record excerpt` refusal route — FT401 owns the refusal routes, and `bench record next` names the excerpt form only through a checkpoint repair.
- The record imports and a frame check on `--excerpt` — FT318 owns the record imports, and `--out` gives `--excerpt` a complete file.
- Sibling-chunk record bases — FT420 owns concurrent chunk bases, and `required_base` keeps the serial derivation that `bench record ticket` uses.
- TOON field selectors and a record reader by ID — FT125 owns them, and the paged status view serves the reads of this spec.
- The `bench commit` change set — FT258 owns it, and this spec adds no commit candidate form.
- A reauthorize that carries the binding to the new request — FT422 owns that binding family, and the refusals of RP46 and RP48 name the bound token instead.
- A containment rule for the `--out` path inside a tracked tree — the `bench preflight evidence --to <dir>` export has none (internal/chargeevidence/export.go:140-190), and RP38 keeps a relative path out of the target tree.
- The unknown-request refusal of `bench worktree release` — it shares the identity producer, and the landing caller adds the sentence, so `bench worktree land` keeps the RP48 behavior.

## Out of scope

- A reauthorize that rewrites the delivery binding: 4 edits, 2 gate runs, under FT422.
- A tip-current remedy in the preflight `next` cell: 3 edits, 1 gate run, as a WF36 reversal that needs a reviewer decision.
- A frame check on `--excerpt` imports: 2 edits, 1 gate run, under FT318.

## Further notes

### Flagged additions

These additions go beyond the literal row text. Each one is open to reviewer veto.

1. The status view pages through `responsebound.PageRead`, which adds `--offset` and `--source` to the form (RP9, RP10).
2. The read-only forms skip the primary-checkout refusal (RP12).
3. The unknown-request identity refusal also names the bound token (RP48). The row names one landing refusal, and FT419 reached both sites.
4. `ValidateCommitment` checks the token against the request digest (RP50).
5. A failed `--out` write after a zero verb exit changes the exit to 1 (RP39).
6. The unknown-chunk refusal of the status view names each planned chunk ID (RP21), as the ticket refusal names the chunk tickets (internal/reviewrecord/ticket.go:46).
7. A call with several conflicting pairs names the first pair in flag declaration order (RP57).
8. `RegisterSibling` clears the copied token (RP56). Without it, a sibling copy fails the RP50 check (slicing research, internal/commitment/repository/readiness.go:138).

### Seams beyond the decision source

| seam | path:line | research finding |
|---|---|---|
| finding collector | internal/reviewrecord/coverage.go:17-173 | Slice B: `checkSource` returns only its first refusal, and no function lists every missing or stale entry. |
| required base | internal/reviewrecord/ticket.go:174-188 | Slice B: `chunkBase` is the only required-base rule, and `bench record chunk` stores any base. |
| AXI children | cmd/bench/main.go:117 | Slice A: `record` is a mutation exemption, and a read child needs `axiApprovedChildren`. |
| RR49 form name | internal/conformance/registry/route_bypass.go:24 | Slice A: recordcmd is scanned, and a bare `next` literal is a route fault. |
| paged read | internal/worktree/related.go:85-96 | Slice C: `PageRead` is the precedent for a complete read-only projection. |
| `--out` capability | cmd/bench/census_output.go:20-34 | Slice D: the owner holds the complete output that retention reads. |
| complete accessor | internal/responsebound/owner.go:254-271 | Slice D: `completeOutput` serves retention and is unexported. |
| conflict derivation | internal/reviewrecord/recordcmd/command.go:136-181 | Slice B: `admits` and `together` own the alternation groups, and the refusal prints the bare usage line. |
| plan-only render | internal/preflight/command.go:220-250 | Slice C: `redRows` filters every kind, and plan-only shares the renderer. |
| bound token | internal/intent/ledger/commitment.go:46-53 | Slice D: `DeliveryBinding.Request` holds a digest, and the plain token lives only on the assignment. |
| unbound refusal | internal/commitment/repository/readiness.go:47-84 | Slice D: `readyFor` skips a binding with another request digest and raises `errUnbound`. |

### Source sentence to row

| source sentence (roadmap/FT414.md) | rows |
|---|---|
| A read-only `bench record status` lists the required chunk base. | RP15 |
| It lists each stale or missing entry. | RP4–RP7, RP17 |
| It lists each commit that breaks the chain. | RP18–RP20 |
| A `bench record next` prints the next required record or plan action with its preconditions. | RP23–RP27, RP51, RP52 |
| `bench probe` and `bench test` take an `--out` option that writes the complete response. | RP33–RP41 |
| `bench record status` also prints the plan rows of one chunk. | RP16, RP21 |
| A record verb that receives two exclusive flags names the conflict. | RP29–RP32, RP54, RP55, RP57 |
| A read-only view prints the whole completion plan as TOON. | RP1–RP3, RP8–RP11 |
| `bench preflight build --plan-only` prints the full row table. | RP42–RP44 |
| A landing refusal for a request mismatch names the token that the commitment binding holds. | RP45–RP50, RP56 |
| This row owns only the precondition projection. | Won't handle lines for FT318, FT258, and FT401 |

### Pre-review proof checklist

- `Cited symbols`: `checkSource`, `CheckSource`, `chunkBase`, `coveredAfter`, `digestFollows`, `uncoveredFold`, `EntryCurrent`, `ChunkEntry`, `currentPlan`, `CompletionOrder`, `admits`, `together`, `renderChecks`, `redRows`, `readyFor`, `errUnbound`, `Store.Start`, `ValidateCommitment`, `PageRead`, and `completeOutput` resolve in the tree at the spec base.
- `Import edges`: `recordcmd` already imports `responsebound` through `stage.go`, and `intent` already reads `ledger`. No new package edge is planned. The slicing research confirms each edge with `go list`.
- `Source-row clauses and occurrences`: the table above maps each clause. The FT405 tip-current occurrence maps to a Won't handle line. The FT293 excerpt-route occurrence maps to the FT401 Won't handle line. The `preflight evidence --check-current` occurrence has no clause in the row body, and it maps to no row.
- `Promised field labels`: the output table in Implementation decisions names each table and field.
- `Changed-function callers`: `checkSource` has the callers `CheckSource`, `RecordCompletion` in `completion.go`, `CheckTrees` through `check.go`, and the tests in `source_test.go` and `fold_test.go`. It reaches `checkCompletion` in completion mode. `renderChecks` has one caller, `verdictCommand`. `readyFor` has the callers `authorizeCandidate` and `Ready`. The slicer reruns `bench consumers` for each.
- `Copy survival`: RP22 reds if the collector keeps a predicate copy, and RP41 reds if a second complete-output accessor survives.
- `Rendered-shape readers`: the plan-only table change reaches `verdict_summary_test.go`, `commitment_test.go`, and `charge_test.go`. The help change reaches `help_inventory_test.go`. The landing text change reaches `land_reauthorization_test.go`, `identity_component_test.go`, and `publication_test.go`.
- `Pin operators`: RP1–RP21 compare exact tables. RP29–RP31 compare a substring and the exit code. RP33 compares the import exit code. RP34 compares bytes for equality.
- `Entry reads`: the read-only forms read HEAD through `git`. Their tests drive linked fixture repositories, so no ambient read needs an internal form.
- `Derived expectations`: RP22 derives its expectation from `CheckSource`, which is the production oracle that the refactor keeps. RP24 derives its expectation from the printed refusal of the same fixture.
- `Consolidated rules`: none.
- `Quantified obligations`: "each alternation group" is ten layout pairs today, and RP29–RP31, RP54, and RP55 cover each distinct rule path. "Every check row" is the `Decide` list of the plan-only kind (RP42).
- `Workflow-step writes`: none. The two new forms write nothing, and `--out` writes one caller file outside the record and the plan digests.

### Line calibration

Expected repair rounds: 2 for RP-C1, 1 for RP-C6, and 1 for each other chunk / confidence 6.
The scorecard routes a code ticket that crosses owner seams to xhigh, and RP-C1 and RP-C6 meet that condition.
Under `craft-line`, a harder chunk keeps the declared effort, and the author of every other chunk runs the declared model at medium effort.
