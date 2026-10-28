# Background grade of lane-passed tips

Status: staged

Roadmap: FT451

Decision source: specs/background-grade/decisions/background-grade.md (ready compiled map, 16 decisions), and the reviewer's late answer of 2026-10-27 that a wait promotes only the I/O class.

Audience: this repository only, because a linked repo starts no background grade.

Verification log: 2 iteration(s) to accept — map-only route, ticket round only. Iteration 1 found a stop signal outside the lifetime owner, a landing wait that contradicted map decisions 6 and 13, too few Forbid rows, an oversized ticket, and bundled seam chunks. The coordinator then moved the census record out of the git-ignored build-log folder.

## Problem

A chunk checkpoint and a landing each run a whole gate after the work is already committed.
A whole gate takes about 4.3 to 6.5 minutes, and the orchestrator waits for it.
The prepared run of FT445 starts that gate earlier, but only at two guidance points, and a plan commit during the run makes it refuse.

## Solution

After each lane pass that publishes a commit in a kit worktree, Bench grades the new tip in the background.
The grade computes the exact subject of the next chunk checkpoint or landing and writes its green to the shared evidence store.
The consumer then finds that green, or waits for the running grade of its key, and runs no second gate.
At most one background grade runs on the machine, at `-p=2 -parallel=2` and at the lowest priority.

## User stories

Line: opus / high.
Implementation-line reason: BG-C6 is the hardest material chunk, because it adds a detached process chain, a machine queue, and a consumer wait across process boundaries. The map fixes every predicate, but the claim and lifetime seams are new, and timing-sensitive tests are weaker than pure fixtures. So the mid tier at high effort is sufficient, with the harder chunks marked.
Harder chunks: BG-C4, BG-C6.

### The consumers find the background green (map decision 3)

1. As an orchestrator, I want a chunk checkpoint of a lane-passed tip to reuse the background green of its subject, so that the checkpoint runs no second gate.
2. As an operator, I want the landing of a spec-less or tickets-only kit worktree to reuse the background green of its subject, so that the landing runs no second gate.
3. As a reviewer, I want the complete checkpoint to keep its own run under the broker, so that a spec landing keeps the evidence that ADR 0031 states.

### A lane pass starts the grade (map decision 4)

4. As an author, I want `bench commit` to start a background grade after a lane pass that publishes, so that no guidance step can be forgotten.
5. As an orchestrator, I want `bench worktree merge` to start a background grade after a merge that publishes, so that a folded tip is graded too.
6. As a linked-repo operator, I want my repository to start no background grade, so that my gate never runs unasked.
7. As an author, I want a dry run, a refused commit, and a merge with nothing to publish to start no grade, so that only a published tip costs a run.
8. As an author, I want `bench commit` to return without waiting for its grade, so that the commit keeps its own speed.

### One derivation of the run-binary source (map decision 5)

9. As a reviewer, I want the grade identity to frame the build inputs of the source root that the run selects, so that a key names the binary that graded.
10. As a reviewer, I want a run that selects no run binary to frame one stable value, so that such a run keeps a deterministic key.

### Claims and the wait (map decision 6 and the late answer)

11. As a consumer, I want a running grade to hold a claim that names its key, its worktree, its owner pid, and the boot id, so that I can find it.
12. As a consumer, I want to wait for a live claim on my own key and then reuse its green, so that I run no duplicate gate.
13. As a consumer, I want my wait to return the grade's I/O class to best-effort and to keep its nice value at 19, so that the grade speeds up within the platform limits.
14. As an orchestrator, I want the wait to count against the consumer's own gate timeout, so that a stuck grade holds the consumer no longer than a gate run.
15. As an orchestrator, I want a checkpoint to wait before it takes the checkout lock, so that the wait holds no lock.
16. As a consumer, I want to run my own gate when the waited grade ends without a green, so that a red or an abort never answers me.
17. As a consumer, I want a queued grade of my key to leave the queue when I start my own run, so that one key is not graded twice.
18. As a consumer, I want to stop and drop the other-key grades of my worktree, and to wait for their exit before I run, so that two runs never write one delta record.
19. As an author, I want a background grade to take no checkout lock and to write no gate cache, so that my foreground gate verbs never refuse because of it.
20. As an operator, I want two foreground runs in one checkout to keep today's refusal, so that FT451 changes no foreground contention.

### The machine bound (map decision 7)

21. As an operator, I want at most one background grade to run at once on the machine, so that background work cannot crowd the host.
22. As an operator, I want a background grade to run its Go test commands at `-p=2 -parallel=2` through a Bench-owned width input, so that the bound holds under the closed gate environment.
23. As a consumer, I want the width input to stay out of the grade key, so that a background green answers my run at another width.
24. As an author, I want a run without the width input to keep its argv, so that every foreground gate stays as it is.
25. As an operator, I want a background grade to run at the lowest CPU and I/O priority, so that foreground work keeps the machine.
26. As an operator, I want grades to deduplicate by grade key, so that a record commit or a second lane pass of one key costs no run.
27. As an operator, I want each worktree to keep at most one queued grade, and a newer tip with a new key to replace it, so that a superseded tip takes no test slot.
28. As an operator, I want a newer tip with a new key to stop the running grade of its worktree, so that a superseded tip keeps no test slot.
29. As an operator, I want a grade whose key another worktree's latest tip still holds to keep running or queued, so that one worktree cannot cancel the work of another.
30. As an operator, I want queued grades to start in the order of their lane passes, so that no worktree waits forever.

### The subject (map decision 8)

31. As an orchestrator, I want the chunk grade to grade the HEAD tree without the review record, so that its key is the key of the chunk checkpoint.
32. As an orchestrator, I want the chunk grade to take the chunk id from the latest chunk entry of the record, and to run whole when none resolves, so that it can take the chunk run decision.
33. As an orchestrator, I want the chunk grade to take the delta set after a test-only change from a whole green, so that a test repair is graded cheaply.
34. As an orchestrator, I want the chunk grade to grade the subject alone, with no evidence pass and no record prose grade, so that the checkpoint keeps those checks for itself.
35. As an orchestrator, I want the chunk grade to run as the executable that the wrapper resolves for the target worktree, so that its key matches the checkpoint's key.
36. As an operator, I want the landing grade to grade the published tree composed on the destination at the start of the grade, so that its key matches the landing's key.
37. As an operator, I want a spec-less landing grade to apply no publish transform, so that its subject is the tip composed on the destination.
38. As an operator, I want the landing grade to run as the authenticated broker with the destination schedule, so that its green can answer the landing.
39. As an author, I want no subject member to read the live worktree, so that my edits, my plan commits, and my probes during a grade leave its verdict valid.
40. As a reviewer, I want a wrong derivation to cost only a lookup miss, so that no grade writes a green under a key that it did not run.

### The signal (map decision 9)

41. As an author, I want a background record of my worktree that names the outcome, the key, the tip, and the log path, so that I can find a background red.
42. As an author, I want `bench status` to show the latest background outcome of each worktree, so that a red reaches me without a search.
43. As an author, I want my next `bench commit` to print the latest background outcome and its log path, so that I see a red in my own flow.
44. As a consumer, I want a background red to give me no refusal, so that a load-coupled red never blocks me.
45. As an operator, I want a background red to write no reader hold, so that a speculative red never sends landings to the whole gate.
46. As an operator, I want a whole background green to write the reader table, as any whole green does, so that the reader set keeps its input.

### The prepared run retires (map decision 10)

47. As an orchestrator, I want `--prepare` to leave the gate grammar, so that one mechanism grades a chunk subject ahead of its checkpoint.
48. As a cold session, I want the guidance to name no prepared run, so that I never start a hand step that no longer exists.
49. As an orchestrator, I want the Land section to state that a chunk checkpoint waits for the background grade of its key, so that I know why the checkpoint can pause.

### The census (map decision 11)

50. As a reviewer, I want the first build ticket to measure the foreground cost of a background gate before any code change, and to stop the build on a failed bound, so that ADR 0024 gets its census.

### The lifecycle (map decision 13)

51. As an operator, I want a landing of another key to stop the grade of its worktree and to drop its queued grade before it grades, so that the landing owns the machine.
52. As an operator, I want `bench worktree release`, `bench worktree clean`, and `bench worktree reclaim` to stop the grades of each worktree that they remove, so that no grade outlives its worktree.
53. As an operator, I want a landing stack to stop the grades of each member worktree before it grades, so that the stack owns the machine.
54. As an operator, I want `bench cache clean` to stop every grade and to drop every queued grade before it cleans, so that optional work never blocks a clean.
55. As an author, I want a stopped grade to write no green and to record `stopped`, so that a stop never leaves a false green.

### The delta record and the reader set (map decisions 14 and 15)

56. As an orchestrator, I want a background delta green to write the delta record of its worktree, so that the next chunk checkpoint of that base and subject reuses it.
57. As an operator, I want a lane pass to start no grade when the tip's landing plan selects the reader set, so that the machine pays no whole run that a reader set replaces.

### The queue runner (map decision 16)

58. As an operator, I want a grade to start as a detached session through the lifetime owner, with a census row, so that it outlives `bench commit` and stays owned.
59. As an operator, I want an ended grade to start the oldest queued grade and then exit, so that no resident process exists.
60. As an operator, I want a lane pass that finds no live running grade to start the oldest queued grade, so that a dead chain resumes.
61. As a consumer, I want a claim to count as ended when its owner is not alive or its boot id differs, so that a dead grade never holds my wait.
62. As an author, I want a timeout, an interrupt, an infrastructure refusal, or a dead owner to record `aborted`, so that each end without a verdict has one name.

### The rules are written down

63. As a cold session, I want the ADRs, the reference, and the profile to state the background grade rules, so that I can act without this spec.

## Implementation decisions

- A new leaf package `internal/gate/backgroundstate` owns the claims, the machine queue and slot, the background records, the liveness rule, and the stop of a grade. It imports no `internal/gate` code, so `internal/gocache`, `internal/worktree`, `internal/status`, and `internal/commit` can import it.
- The stop sends its signal through `internal/subprocess`, because the lifetime closure check refuses a signal from any other package.
- Package `internal/gate` owns the background run: the subject derivation for both consumers, the run mode that writes no gate cache and takes no checkout lock, the width input, the priority, and the plumbing entry that runs one queued grade.
- The plumbing entry is a row of the internal command inventory. It reuses an existing route line of `bin/bench.sh`, because that file is over its line budget.
- `bench commit` and `bench worktree merge` call one launch function of `internal/gate` after they publish. The launch goes through the lifetime owner of `internal/subprocess` and takes a census row in `internal/conformance`.
- The chunk checkpoint waits in the run transaction before the checkout lock. The landing computes its key before the destination lease, waits on a claim of that key, and computes the key again under the lease.
- The landing grade authenticates the broker as `runUnderBroker` does: the primary broker manifest binds its digest, and a mismatch aborts the grade before any oracle.
- `cmd/bench/main.go` and `bin/bench.sh` are over their line budget. `subject.go`, `checkpoint.go`, `commit.go`, and `command_registry.go` have 2 lines of headroom or less. A ticket that edits one of them keeps its line count or moves code out in the same ticket.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| BG-C0 / 01-measure-the-background-census.md | The census measures the foreground cost of one background gate at `-p=2` before any code change, and a failed bound stops the build | BG93, BG94, BG95 | review of the census record | no |
| BG-C1 / 02-derive-the-run-binary-source-once.md, 03-read-the-baseline-members-from-commits.md | The grade identity frames the run's own run-binary source, and it reads its baseline members from commits | BG12-BG18, BG109 | the `internal/gate` package tests | no |
| BG-C2 / 04-add-the-background-width-input.md | A Bench-owned width input sets `-p=2 -parallel=2` outside the grade key | BG19-BG22 | the `internal/gate` package tests | no |
| BG-C3 / 05-record-background-claims-and-outcomes.md | The leaf store records claims, outcomes, and the queue, and it stops a grade through the lifetime owner | BG32-BG35, BG112, BG113 | the `internal/gate/backgroundstate` and `internal/subprocess` package tests, and the `process-lifetime-closure` check | no |
| BG-C4 / 06-grade-a-chunk-subject-in-the-background.md | One background run grades a chunk subject, records `green` or `red`, and a chunk checkpoint reuses its green | BG1, BG23, BG24, BG26-BG29, BG41, BG42, BG106, BG107, BG114 | the `internal/gate`, `cmd/bench`, and `internal/conformance` package tests, and the `bench-sh-routes` and `subcommand-routing` checks | yes |
| BG-C5 / 07-complete-the-background-run-outputs.md | The background run writes the delta record and the reader table, never a hold, runs at the lowest priority, and names each abort | BG25, BG30, BG31, BG36-BG40 | the `internal/gate` package tests | no |
| BG-C6 / 08-start-a-background-grade-after-a-lane-pass.md, 09-queue-and-supersede-background-grades.md, 10-let-a-chunk-checkpoint-wait-on-a-background-grade.md | A lane pass queues one grade, the chain runs the queue, and a chunk checkpoint waits for its key | BG5-BG11, BG43-BG66, BG83, BG108, BG115, BG116 | the `internal/gate`, `internal/commit`, `internal/worktree`, `internal/subprocess`, `internal/conformance`, and `cmd/bench` package tests, and the `process-lifetime-closure` and `system` checks | yes |
| BG-C7 / 11-grade-a-landing-subject-in-the-background.md, 12-let-a-landing-wait-on-a-background-grade.md | A spec-less or tickets-only kit worktree gets a landing grade, and its landing reuses or waits for it | BG2, BG3, BG67-BG77, BG110, BG111 | the `internal/gate`, `internal/worktree`, `internal/landing/published`, and `cmd/bench` package tests | no |
| BG-C8 / 13-stop-background-grades-when-work-goes-away.md, 14-show-background-outcomes.md | Worktree removal and the cache clean stop grades, and the author sees each outcome | BG78-BG82, BG84, BG85, BG99 | the `internal/worktree`, `internal/gocache`, `internal/status`, `internal/commit`, and `cmd/bench` package tests | no |
| BG-C9 / 15-retire-the-prepared-run.md, 16-record-the-background-grade-rules.md | The prepared run leaves the grammar and the guidance, and the ADRs, the reference, and the profile state the new rules | BG4, BG86-BG92, BG96-BG98, BG100-BG105 | the `internal/gate`, `cmd/bench`, and `internal/anchors` package tests, and the `docs-currency-workflow` and `prose` checks | no |

The chunks run serially in the order of the table, because they share `CHANGELOG.md`, the background state package, and the run transaction.
BG-C0 comes first, because the census must pass before any code change.
BG-C1 to BG-C4 each hold a seam that a later chunk consumes, so each review closes before any consumer ticket starts.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","review_mode":"unified","run_id":"ft451-background-grade-full-20261027","orchestrator_session":"claude:bg-orchestrator-20261027","author_limit":3,"assignments":{"01-measure-the-background-census.md":[{"session":"claude:bg-t01-author-20261027","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"dfc5e6d75782f91dfb4c5a7aa8308ca84ed25b25","native_ref":"claude-agent:bench-writer-high:bg-t01-author-20261027"},{"session":"claude:bg-t01-repair1-20261027","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"e51f889c4ef3e7abb9eccefd1b4fa7b57e2e08de","native_ref":"claude-agent:bench-writer-high:bg-t01-repair1-20261027","predecessor":"claude:bg-t01-author-20261027","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:bg-t01-author-20261027","preserved":"7bdc95c72c6ef570407a4105cfd8d3597b69e989"}],"02-derive-the-run-binary-source-once.md":[{"session":"claude:bg-t02-author-20261027","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"ac53a09920a8d9e716b2aeca7c141f392afd15dc","native_ref":"claude-agent:bench-writer-high:bg-t02-author-20261027"},{"session":"claude:bg-t02-repair1-20261027","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"fd96df22cf42d26a1160fe130a7af9ba660b522d","native_ref":"claude-agent:bench-writer-high:bg-t02-repair1-20261027","predecessor":"claude:bg-t02-author-20261027","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:bg-t02-author-20261027","preserved":"2f788c7e089aeb26a394761866c52b08caaba7a7"}],"03-read-the-baseline-members-from-commits.md":[{"session":"claude:bg-t03-author-20261027","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"1a0bc5f8d3c3bade273f633461bffff38959b9c0","native_ref":"claude-agent:bench-writer-high:bg-t03-author-20261027"}],"04-add-the-background-width-input.md":[{"session":"claude:bg-t04-author-20261027","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"2931581854e9db7c7d72dcfcdee77ee30b8154f5","native_ref":"claude-agent:bench-writer-high:bg-t04-author-20261027"}],"05-record-background-claims-and-outcomes.md":[{"session":"claude:bg-t05-author-20261027","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"medium","source":"d6c1b9f6031e6684ac62dbfe2438e7e25b4cfc47","native_ref":"claude-agent:bench-writer-medium:bg-t05-author-20261027"},{"session":"claude:bg-t05-repair1-20261027","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"medium","source":"40ed1d68afb5b8c40011eb3a867e55d5216e4cec","native_ref":"claude-agent:bench-writer-medium:bg-t05-repair1-20261027","predecessor":"claude:bg-t05-author-20261027","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:bg-t05-author-20261027","preserved":"1ab1a028e"}],"06-grade-a-chunk-subject-in-the-background.md":[{"session":"claude:bg-t06-author-20261028","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"36319642fa9daae27ab039ed8999dab102ff5a13","native_ref":"claude-agent:bench-writer-high:bg-t06-author-20261028"},{"session":"claude:bg-t06-repair1-20261028","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"6dfa4408f6d0a0a97e747fe4e8808e1f6bbaa955","native_ref":"claude-agent:bench-writer-high:bg-t06-repair1-20261028","predecessor":"claude:bg-t06-author-20261028","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:bg-t06-author-20261028","preserved":"c05c87a8d42e7c13e44c8181f2cb43dddd7e1eb7"}],"07-complete-the-background-run-outputs.md":[{"session":"claude:bg-t07-author-20261028","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"medium","source":"b210690fa0ae9d27b9ddedd556ab6b614276f889","native_ref":"claude-agent:bench-writer-medium:bg-t07-author-20261028"},{"session":"claude:bg-t07-repair1-20261028","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"medium","source":"f95899567fbb20bd93722b98988c34bbd2516e94","native_ref":"claude-agent:bench-writer-medium:bg-t07-repair1-20261028","predecessor":"claude:bg-t07-author-20261028","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:bg-t07-author-20261028","preserved":"3cf91fd013916354df370e3930d4e26a400f01fd"}],"08-start-a-background-grade-after-a-lane-pass.md":[{"session":"claude:bg-t08-author-20261028","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"28d89b56064f6e04a3989eaf07290a996d3d17d6","native_ref":"claude-agent:bench-writer-high:bg-t08-author-20261028"}],"09-queue-and-supersede-background-grades.md":[{"session":"claude:bg-t09-author-20261028","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"3d17c259ae263767556f8ec7bd7cc580d3afe045","native_ref":"claude-agent:bench-writer-high:bg-t09-author-20261028"},{"session":"claude:bg-t09-repair1-20261028","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"d194f6ed2d19302310bbe6b2639e4fb508fd41d0","native_ref":"claude-agent:bench-writer-high:bg-t09-repair1-20261028","predecessor":"claude:bg-t09-author-20261028","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:bg-t09-author-20261028","preserved":"a58ef0a78ef140161482bf0adbf9e9ac4864a503"}],"10-let-a-chunk-checkpoint-wait-on-a-background-grade.md":[{"session":"claude:bg-t10-author-20261028","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"high","source":"9acd4dad595328ba6858049d5f4bf74c42ad4ce5","native_ref":"claude-agent:bench-writer-high:bg-t10-author-20261028"}],"11-grade-a-landing-subject-in-the-background.md":[{"session":"claude:bg-t11-author-20261028","assignment":"c7bb50cb85b54be3caa07403fc532593","model":"opus","effort":"medium","source":"669734be4ca17501f1bf05ca9d41c882a3403f5d","native_ref":"claude-agent:bench-writer-medium:bg-t11-author-20261028"}],"12-let-a-landing-wait-on-a-background-grade.md":[],"13-stop-background-grades-when-work-goes-away.md":[],"14-show-background-outcomes.md":[],"15-retire-the-prepared-run.md":[],"16-record-the-background-grade-rules.md":[]}},"chunks":[{"id":"BG-C0","tickets":["01-measure-the-background-census.md"],"verification":[{"id":"conformance","command":"bench test --check root-conformance","ticket":"01-measure-the-background-census.md"}]},{"id":"BG-C1","tickets":["02-derive-the-run-binary-source-once.md","03-read-the-baseline-members-from-commits.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"02-derive-the-run-binary-source-once.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"02-derive-the-run-binary-source-once.md"},{"id":"gate-03","command":"bench test --package ./internal/gate","ticket":"03-read-the-baseline-members-from-commits.md"},{"id":"conformance-03","command":"bench test --check root-conformance","ticket":"03-read-the-baseline-members-from-commits.md"}]},{"id":"BG-C2","tickets":["04-add-the-background-width-input.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"04-add-the-background-width-input.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"04-add-the-background-width-input.md"}]},{"id":"BG-C3","tickets":["05-record-background-claims-and-outcomes.md"],"verification":[{"id":"state","command":"bench test --package ./internal/gate/backgroundstate","ticket":"05-record-background-claims-and-outcomes.md"},{"id":"subprocess","command":"bench test --package ./internal/subprocess/...","ticket":"05-record-background-claims-and-outcomes.md"},{"id":"lifetime","command":"bench test --check process-lifetime-closure","ticket":"05-record-background-claims-and-outcomes.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"05-record-background-claims-and-outcomes.md"}]},{"id":"BG-C4","tickets":["06-grade-a-chunk-subject-in-the-background.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate/...","ticket":"06-grade-a-chunk-subject-in-the-background.md"},{"id":"command","command":"bench test --package ./cmd/bench","ticket":"06-grade-a-chunk-subject-in-the-background.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"06-grade-a-chunk-subject-in-the-background.md"},{"id":"routes","command":"bench test --check bench-sh-routes","ticket":"06-grade-a-chunk-subject-in-the-background.md"},{"id":"routing","command":"bench test --check subcommand-routing","ticket":"06-grade-a-chunk-subject-in-the-background.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"06-grade-a-chunk-subject-in-the-background.md"}]},{"id":"BG-C5","tickets":["07-complete-the-background-run-outputs.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate/...","ticket":"07-complete-the-background-run-outputs.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"07-complete-the-background-run-outputs.md"}]},{"id":"BG-C6","tickets":["08-start-a-background-grade-after-a-lane-pass.md","09-queue-and-supersede-background-grades.md","10-let-a-chunk-checkpoint-wait-on-a-background-grade.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate/...","ticket":"08-start-a-background-grade-after-a-lane-pass.md"},{"id":"commit","command":"bench test --package ./internal/commit","ticket":"08-start-a-background-grade-after-a-lane-pass.md"},{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"08-start-a-background-grade-after-a-lane-pass.md"},{"id":"subprocess","command":"bench test --package ./internal/subprocess/...","ticket":"08-start-a-background-grade-after-a-lane-pass.md"},{"id":"command","command":"bench test --package ./cmd/bench","ticket":"08-start-a-background-grade-after-a-lane-pass.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"08-start-a-background-grade-after-a-lane-pass.md"},{"id":"lifetime","command":"bench test --check process-lifetime-closure","ticket":"08-start-a-background-grade-after-a-lane-pass.md"},{"id":"system","command":"bench test --check system","ticket":"08-start-a-background-grade-after-a-lane-pass.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"08-start-a-background-grade-after-a-lane-pass.md"},{"id":"gate-09","command":"bench test --package ./internal/gate/...","ticket":"09-queue-and-supersede-background-grades.md"},{"id":"subprocess-09","command":"bench test --package ./internal/subprocess/...","ticket":"09-queue-and-supersede-background-grades.md"},{"id":"conformance-09","command":"bench test --check root-conformance","ticket":"09-queue-and-supersede-background-grades.md"},{"id":"gate-10","command":"bench test --package ./internal/gate/...","ticket":"10-let-a-chunk-checkpoint-wait-on-a-background-grade.md"},{"id":"conformance-10","command":"bench test --check root-conformance","ticket":"10-let-a-chunk-checkpoint-wait-on-a-background-grade.md"}]},{"id":"BG-C7","tickets":["11-grade-a-landing-subject-in-the-background.md","12-let-a-landing-wait-on-a-background-grade.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate/...","ticket":"11-grade-a-landing-subject-in-the-background.md"},{"id":"published","command":"bench test --package ./internal/landing/published","ticket":"11-grade-a-landing-subject-in-the-background.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"11-grade-a-landing-subject-in-the-background.md"},{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"12-let-a-landing-wait-on-a-background-grade.md"},{"id":"command","command":"bench test --package ./cmd/bench","ticket":"12-let-a-landing-wait-on-a-background-grade.md"},{"id":"gate-12","command":"bench test --package ./internal/gate/...","ticket":"12-let-a-landing-wait-on-a-background-grade.md"},{"id":"conformance-12","command":"bench test --check root-conformance","ticket":"12-let-a-landing-wait-on-a-background-grade.md"},{"id":"landing","command":"bench test --package ./internal/landing","ticket":"11-grade-a-landing-subject-in-the-background.md"}]},{"id":"BG-C8","tickets":["13-stop-background-grades-when-work-goes-away.md","14-show-background-outcomes.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"13-stop-background-grades-when-work-goes-away.md"},{"id":"gocache","command":"bench test --package ./internal/gocache","ticket":"13-stop-background-grades-when-work-goes-away.md"},{"id":"command","command":"bench test --package ./cmd/bench","ticket":"13-stop-background-grades-when-work-goes-away.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"13-stop-background-grades-when-work-goes-away.md"},{"id":"status","command":"bench test --package ./internal/status","ticket":"14-show-background-outcomes.md"},{"id":"commit","command":"bench test --package ./internal/commit","ticket":"14-show-background-outcomes.md"},{"id":"conformance-14","command":"bench test --check root-conformance","ticket":"14-show-background-outcomes.md"},{"id":"state","command":"bench test --package ./internal/gate/backgroundstate","ticket":"14-show-background-outcomes.md"}]},{"id":"BG-C9","tickets":["15-retire-the-prepared-run.md","16-record-the-background-grade-rules.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"15-retire-the-prepared-run.md"},{"id":"command","command":"bench test --package ./cmd/bench","ticket":"15-retire-the-prepared-run.md"},{"id":"anchors","command":"bench test --package ./internal/anchors","ticket":"15-retire-the-prepared-run.md"},{"id":"guidance","command":"bench test --check docs-currency-workflow","ticket":"15-retire-the-prepared-run.md"},{"id":"prose","command":"bench test --check prose","ticket":"15-retire-the-prepared-run.md"},{"id":"conformance","command":"bench test --check root-conformance","ticket":"15-retire-the-prepared-run.md"},{"id":"anchors-16","command":"bench test --package ./internal/anchors","ticket":"16-record-the-background-grade-rules.md"},{"id":"guidance-16","command":"bench test --check docs-currency-workflow","ticket":"16-record-the-background-grade-rules.md"},{"id":"prose-16","command":"bench test --check prose","ticket":"16-record-the-background-grade-rules.md"},{"id":"conformance-16","command":"bench test --check root-conformance","ticket":"16-record-the-background-grade-rules.md"}]}],"final_verification":[{"id":"whole-project","command":"bench gate"},{"id":"coverage","command":"bench coverage --check specs/background-grade/spec.md"}]}
```

## Testing decisions

- A good test drives a real Git fixture through the gate's own entry points. It reads the evidence store, the delta record, the background record, and the claim bytes, and it counts oracle runs.
- The background run tests reuse the checkpoint fixtures of `internal/gate`: `recordtest`, `outcomeRuns`, and the armed checkpoint runner. A fixture gate script writes its own nice value and I/O class into a witness file.
- The launch and queue tests drive real processes, and they observe a session id, a pid, and a claim file. No test reads elapsed time as its only oracle.
- The worktree tests use the landing and stack fixtures of `internal/worktree`. Each new worktree test runs in parallel, or it raises `worktreeSerialCeiling` in the same ticket.
- The gate observes the feature at the chunk checkpoint and at the landing, which reuse a background green under the grade key.

### Seam diagram

    trigger: `bench commit` or `bench worktree merge` publishes a lane-passed tip in a kit worktree
        │
        ▼
    tip + worktree  ──▶  [ launch ]  ──▶  machine queue  ──▶  [ background run ]  ──▶  store green, delta record, background record
                             ◀ tests attach here: a real commit, then the queue file, the session id, and the claim
                                                                      ◀ tests attach here: the run entry over a checkpoint fixture, then the store, the records, and the oracle count
        ▼
    chunk checkpoint or landing  ──▶  [ lookup, then a wait on a live claim ]  ──▶  reuse, or its own run
                                          ◀ tests attach here: a held claim fixture, then the oracle count, the lock state, and the I/O class

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| BG1 | 1 | A chunk checkpoint after a background chunk green of its tip prints the reuse line and runs no oracle | `internal/gate/background_grade_test.go` (`TestChunkCheckpointReusesABackgroundGreen`) | A background green under any other key leaves the checkpoint a second run |
| BG2 | 2 | The landing of a spec-less kit worktree after its background green runs no oracle | planned TestSpecLessLandingReusesTheBackgroundGreen in internal/worktree/land_background_test.go | A landing grade keyed off the landing subject leaves the landing a second run |
| BG3 | 2 | The landing of a tickets-only kit worktree after its background green runs no oracle | planned TestTicketsOnlyLandingReusesTheBackgroundGreen in internal/worktree/land_background_test.go | A grade that skips the folder close keys a tree that the landing never grades |
| BG4 | 3 | A complete checkpoint after a background chunk green of the same source runs its own oracle | planned TestCompleteCheckpointIgnoresABackgroundChunkGreen in internal/gate/complete_checkpoint_test.go | A chunk green that answers the complete checkpoint skips the published-tree grade |
| BG5 | 4 | A `bench commit` that publishes in a kit worktree queues one grade of its tip | `internal/commit/background_test.go` (`TestCommitStartsABackgroundGrade`) | A commit without the launch call leaves every consumer to run |
| BG6 | 5 | A `bench worktree merge` of kind `merge` queues one grade of the new tip | `internal/worktree/merge_background_test.go` (`TestMergeStartsABackgroundGrade`) | A merge without the launch call leaves a folded tip ungraded |
| BG7 | 6 | A `bench commit` in a linked repo queues no grade | `internal/commit/background_test.go` (`TestLinkedRepoCommitStartsNoGrade`) | A launch without the kit predicate runs a linked repo's gate unasked |
| BG8 | 7 | A `bench commit --dry-run` queues no grade | `internal/commit/background_test.go` (`TestDryRunStartsNoGrade`) | A launch before the publish check grades an unpublished tree |
| BG9 | 7 | A `bench commit` that the lane refuses queues no grade | `internal/commit/background_test.go` (`TestRefusedCommitStartsNoGrade`) | A launch on every exit grades a tree that no branch holds |
| BG10 | 7 | A `bench worktree merge` of kind `current` queues no grade | `internal/worktree/merge_background_test.go` (`TestCurrentMergeStartsNoGrade`) | A launch on every merge kind grades an unchanged tip again |
| BG11 | 8 | `bench commit` exits 0 while its background grade still holds its claim | `internal/commit/background_test.go` (`TestCommitReturnsBeforeItsGrade`) | A launch that waits for the grade adds a whole gate to each commit |
| BG12 | 9 | An ordinary kit run whose `BENCH_KIT` names another checkout frames the build inputs of that checkout | `internal/gate/run_binary_source_test.go` (`TestIdentityFramesTheRunSelection`) | The presence rule frames the graded tree while the run builds from the other kit |
| BG13 | 9 | No grade identity code reads the presence of `.bench/gate-prospective.sh` to pick a kind | review-owned: Standards runs `rg -n "sourceKind" internal/gate` and expects no hit | A surviving second derivation can drift from the run's selection |
| BG14 | 10 | A run with no run-binary selection frames one stable identity value | `internal/gate/run_binary_source_test.go` (`TestNoSelectionFramesAStableValue`) | A missing case frames a path-dependent value and moves the key |
| BG15 | 9 | One clean kit tree gives one run-binary member to its ordinary run and to its prospective run | `internal/gate/run_binary_source_test.go` (`TestOrdinaryAndProspectiveFrameOneSource`) | A derivation that splits by evaluation mode changes the keys of current greens |
| BG16 | 39 | A chunk-subject run keeps its subject when a Go build input of the worktree changes during the oracle | `internal/gate/baseline_commit_test.go` (`TestBaselineReadsSurviveAWorktreeEdit`) | A baseline runner read from the working tree moves the oracle digest |
| BG17 | 39 | A chunk-subject run across a commit that changes a Go build input in its worktree records its green | `internal/gate/baseline_commit_test.go` (`TestBaselineReadsSurviveABuildInputCommit`) | A baseline runner digest read from the moved worktree refuses the run with a changed subject |
| BG18 | 39 | The baseline schedule member frames the manifest bytes of the identity commit, not the working file | `internal/gate/baseline_commit_test.go` (`TestBaselineScheduleReadsTheCommit`) | A working-file read lets an uncommitted manifest edit move the key |
| BG19 | 22 | With the width input set, each Go test phase argv holds `-p=2` and `-parallel=2` | `internal/gate/width_test.go` (`TestWidthInputSetsTheTestWidths`) | An ambient `GOFLAGS` never reaches the closed phase environment |
| BG20 | 23 | The width input leaves the grade key of one subject unchanged | `internal/gate/width_test.go` (`TestWidthInputLeavesTheGradeKey`) | A width that enters the key makes a background green answer nothing |
| BG21 | 24 | Without the width input, no Go test phase argv holds `-p` or `-parallel` | `internal/gate/width_test.go` (`TestNoWidthInputKeepsTheArgv`) | A default width changes every foreground gate |
| BG22 | 22 | A phase child does not inherit the width input | `internal/gate/width_test.go` (`TestPhaseChildrenDropTheWidthInput`) | A nested gate in a test inherits the background width |
| BG23 | 31 | A background chunk grade writes its green under the key of the HEAD tree without the review record | `internal/gate/background_grade_test.go` (`TestBackgroundChunkGradeKeysTheStrippedTree`) | A grade of the raw HEAD tree keys a tree that the checkpoint never grades |
| BG24 | 32 | A background chunk grade with no chunk entry in the record runs the whole schedule | `internal/gate/background_grade_test.go` (`TestBackgroundGradeRunsWholeWithNoChunkEntry`) | A missing fallback aborts the first grade of a build |
| BG25 | 33, 56 | A test-only change after a whole green runs the delta set and writes the delta record of the worktree | `internal/gate/background_outcome_test.go` (`TestBackgroundGradeTakesTheDelta`) | A grade without the chunk decision runs a whole gate for each test repair |
| BG26 | 34 | A record whose evidence the checkpoint refuses still yields a background green of the subject | `internal/gate/background_grade_test.go` (`TestBackgroundGradeSkipsTheEvidencePass`) | A grade that runs the evidence pass refuses before the record commit lands |
| BG27 | 35 | The grade key of a chunk grade frames the digest of the executable that the wrapper resolves for the target worktree | `internal/gate/background_grade_test.go` (`TestBackgroundGradeRunsAsTheTargetExecutable`) | A grade under the launcher's executable keys a member that the checkpoint never frames |
| BG28 | 19 | A background grade leaves the gate cache of its worktree byte for byte unchanged | `internal/gate/background_grade_test.go` (`TestBackgroundGradeWritesNoGateCache`) | A pending cache from the grade makes each foreground lookup refuse |
| BG29 | 19 | A foreground `bench gate` in the worktree takes the checkout lock while a background grade runs | `internal/gate/background_grade_test.go` (`TestForegroundTakesTheLockDuringAGrade`) | A grade that holds the checkout lock makes the foreground refuse |
| BG30 | 25 | The oracle of a background grade runs at nice 19 | `internal/gate/background_outcome_test.go` (`TestBackgroundOracleRunsAtNiceNineteen`) | A grade at the default priority competes with foreground work |
| BG31 | 25 | On Linux, the oracle of a background grade runs in the idle I/O class | `internal/gate/background_outcome_test.go` (`TestBackgroundOracleRunsInTheIdleClass`) | A grade at the default I/O class competes with foreground disk work |
| BG32 | 11 | A running grade holds a claim that names its key, its worktree, its owner pid, and the boot id | `internal/gate/backgroundstate/claim_test.go` (`TestClaimNamesKeyWorktreeOwnerAndBoot`) | A claim without one field cannot answer a consumer or a liveness test |
| BG33 | 61 | A claim whose owner pid is not alive is not live | `internal/gate/backgroundstate/claim_test.go` (`TestClaimOfADeadOwnerIsNotLive`) | A stale claim holds a consumer until its timeout |
| BG34 | 61 | A claim whose boot id differs from the current boot is not live | `internal/gate/backgroundstate/claim_test.go` (`TestClaimOfAnotherBootIsNotLive`) | A reused pid after a reboot passes for a live owner |
| BG35 | 41 | The background record names the outcome, the grade key, the tip, and the log path | `internal/gate/backgroundstate/record_test.go` (`TestBackgroundRecordNamesTheOutcome`) | A record without the log path hides the red from the author |
| BG36 | 45 | A red background landing grade writes no reader hold | `internal/gate/background_outcome_test.go` (`TestBackgroundRedWritesNoHold`) | A hold from a background red sends every landing to the whole gate |
| BG37 | 46 | A whole background green writes the reader table of its tree | `internal/gate/background_outcome_test.go` (`TestBackgroundGreenWritesTheReaderTable`) | A grade that skips the table leaves the reader set without its input |
| BG38 | 62 | A grade whose oracle passes the gate timeout records `aborted` | `internal/gate/background_outcome_test.go` (`TestTimedOutGradeRecordsAborted`) | A timeout recorded as red hides that no verdict exists |
| BG39 | 62 | An interrupted grade records `aborted` | `internal/gate/background_outcome_test.go` (`TestInterruptedGradeRecordsAborted`) | An interrupt with no record leaves the claim to look live |
| BG40 | 62 | A grade whose subject is unavailable records `aborted` and runs no oracle | `internal/gate/background_outcome_test.go` (`TestRefusedGradeRecordsAborted`) | An infrastructure refusal with no record looks like a run in progress |
| BG41 | 40 | A grade with a wrong chunk id writes its green only under the key that it ran | `internal/gate/background_grade_test.go` (`TestWrongChunkOnlyMisses`) | A grade that records under the consumer's key writes a false green |
| BG42 | 39 | A background grade across an edit of a worktree Go file records its green | `internal/gate/background_grade_test.go` (`TestBackgroundSubjectIgnoresAWorktreeEdit`) | A subject that reads the live worktree loses its verdict to the edit |
| BG43 | 58 | A launched grade runs in a session other than the session of the launching verb | `internal/gate/background_launch_test.go` (`TestLaunchedGradeLeavesTheVerbSession`) | A child in the verb's group dies with the harness call |
| BG44 | 58 | The process-lifetime census names the launch caller | the `process-lifetime-closure` check over the launch census row | An uncensused owner call reds the closure check |
| BG45 | 21 | A second lane pass while one grade runs starts no second oracle | `internal/gate/background_launch_test.go` (`TestSecondGradeStartsNoSecondOracle`) | A launch without the slot runs two whole gates at once |
| BG46 | 30 | Queued grades start in the order of their lane passes | `internal/gate/background_queue_test.go` (`TestQueueServesLanePassOrder`) | A queue that serves the newest entry first starves an old worktree |
| BG47 | 26 | A lane pass whose key holds a fresh green queues nothing | `internal/gate/background_queue_test.go` (`TestFreshGreenKeyStartsNothing`) | A launch without the store lookup grades a record commit again |
| BG48 | 26 | A lane pass whose key has a live claim queues nothing | `internal/gate/background_queue_test.go` (`TestRunningKeyStartsNothing`) | A launch without the claim lookup runs one key twice |
| BG49 | 26 | A lane pass whose key is already queued queues nothing | `internal/gate/background_queue_test.go` (`TestQueuedKeyQueuesNothing`) | A launch without the queue lookup queues one key twice |
| BG50 | 27 | A newer tip with a new key replaces the queued grade of its worktree | `internal/gate/background_queue_test.go` (`TestNewerTipReplacesTheQueuedGrade`) | A queue that appends keeps a superseded tip in line |
| BG51 | 28 | A newer tip with a new key stops the running grade of its worktree | `internal/gate/background_queue_test.go` (`TestNewerTipStopsTheRunningGrade`) | A grade that runs on keeps a superseded tip in the slot |
| BG52 | 29 | A running grade whose key is the latest tip key of another worktree keeps running | `internal/gate/background_queue_test.go` (`TestWantedKeyKeepsRunning`) | A stop by worktree alone cancels the work that another worktree needs |
| BG53 | 29 | A queued grade whose key is the latest tip key of another worktree stays queued | `internal/gate/background_queue_test.go` (`TestWantedKeyStaysQueued`) | A replace by worktree alone drops the work that another worktree needs |
| BG54 | 59 | An ended grade starts the oldest queued grade and then exits | `internal/gate/background_queue_test.go` (`TestEndedGradeStartsTheOldestQueued`) | A chain that exits first leaves the queue stalled |
| BG55 | 60 | A lane pass that finds no live running grade starts the oldest queued grade | `internal/gate/background_queue_test.go` (`TestLanePassResumesADeadChain`) | A dead chain leaves the queue stalled until a reboot |
| BG56 | 62 | A grade whose owner died records `aborted` when the next lane pass reads its claim | `internal/gate/background_queue_test.go` (`TestDeadOwnerRecordsAborted`) | A dead grade with no record looks like a run in progress |
| BG57 | 12 | A chunk checkpoint that finds a live claim on its key waits, then reuses the green, with one oracle run in total | `internal/gate/background_wait_test.go` (`TestChunkCheckpointWaitsForItsKey`) | A checkpoint that refuses or runs at once pays a second gate |
| BG58 | 13 | When a consumer starts to wait, the waited grade's I/O class becomes best-effort | `internal/gate/background_wait_test.go` (`TestWaitPromotesTheIOClass`) | A wait with no promotion leaves the grade in the idle class |
| BG59 | 13 | During the wait, the waited grade's nice value stays 19 | `internal/gate/background_wait_test.go` (`TestWaitKeepsTheNiceValue`) | A promotion that lowers the nice value fails with a permission error on Linux |
| BG60 | 14 | A wait that passes the consumer's gate timeout ends the consumer with exit 124 | `internal/gate/background_wait_test.go` (`TestWaitCountsAgainstTheGateTimeout`) | A wait outside the timeout holds the consumer without a bound |
| BG61 | 15 | A waiting checkpoint holds no checkout lock | `internal/gate/background_wait_test.go` (`TestWaitHoldsNoCheckoutLock`) | A wait under the lock makes every other foreground verb refuse |
| BG62 | 16, 44 | A waited grade that ends red makes the checkpoint run its own oracle | `internal/gate/background_wait_test.go` (`TestCheckpointRunsAfterABackgroundRed`) | A checkpoint that refuses on the background red gives a red authority |
| BG63 | 17 | A queued grade of the checkpoint's key leaves the queue when the checkpoint starts its own run | `internal/gate/background_wait_test.go` (`TestQueuedSameKeyLeavesTheQueue`) | A queued duplicate grades the key again after the checkpoint |
| BG64 | 18 | Before its own run, a checkpoint stops the other-key grade of its worktree and waits for its exit | `internal/gate/background_wait_test.go` (`TestCheckpointStopsOtherKeyGrades`) | Two runs write the delta record of one worktree at once |
| BG65 | 18 | Before its own run, a checkpoint drops the other-key queued grade of its worktree | `internal/gate/background_wait_test.go` (`TestCheckpointDropsOtherKeyQueue`) | A queued superseded grade starts after the checkpoint |
| BG66 | 20 | A second foreground run in one checkout refuses with `gate execution already in progress` | `internal/gate/run_failure_outcomes_test.go` (`TestGateRunRefusesInFlightExecution`) | A wait added to the foreground pair changes today's refusal |
| BG67 | 36 | A background landing grade writes its green under the key of the published tree composed on the destination | planned TestLandingGradeKeysThePublishedComposition in internal/gate/background_landing_test.go | A grade of the raw tip keys a tree that the landing never grades |
| BG68 | 37 | A spec-less landing grade keys the tip composed on the destination, with no transform | planned TestSpecLessLandingGradeHasNoTransform in internal/gate/background_landing_test.go | A transform on a spec-less tip keys a tree that the landing never grades |
| BG69 | 38 | The grade key of a landing grade frames the digest of the installed broker | planned TestLandingGradeRunsAsTheBroker in internal/gate/background_landing_test.go | A grade under the candidate executable keys a member that the landing never frames |
| BG70 | 38 | A broker whose digest differs from its manifest makes the landing grade record `aborted` and run no oracle | planned TestUnauthenticatedBrokerAbortsTheGrade in internal/gate/background_landing_test.go | A grade under an unauthenticated broker writes a green that the landing trusts |
| BG71 | 38 | A landing grade frames the phase schedule of the destination commit | planned TestLandingGradeTakesTheDestinationSchedule in internal/gate/background_landing_test.go | A candidate schedule lets the tip omit the checks that grade it |
| BG72 | 57 | A tip whose landing plan against the destination at the lane pass selects the reader set queues no grade | planned TestReaderSetLandingStartsNoGrade in internal/gate/background_landing_test.go | A whole background run replaces a shorter reader set |
| BG73 | 57 | A tip whose landing plan selects the whole gate queues one landing grade | planned TestWholeGateLandingStartsAGrade in internal/gate/background_landing_test.go | A skip that reads the wrong plan leaves a whole landing ungraded |
| BG74 | 40 | A landing grade with a wrong deliverable writes its green only under the key that it ran | planned TestWrongDeliverableOnlyMisses in internal/gate/background_landing_test.go | A grade that records under the landing's key writes a false green |
| BG75 | 12 | A landing that finds a live claim on its own key waits before it takes the destination lease | planned TestLandingWaitsBeforeTheLease in internal/worktree/land_background_test.go | A wait under the lease blocks every other landing |
| BG76 | 12 | After the wait, the landing takes its lease and reuses the green with no oracle run | planned TestLandingReusesTheGreenAfterTheWait in internal/worktree/land_background_test.go | A landing that skips the second lookup pays a second gate |
| BG77 | 51 | A landing of another key stops the grade of its worktree and drops its queued grade before it grades | planned TestLandingOfAnotherKeyStopsTheGrade in internal/worktree/land_background_test.go | A grade beside the landing gate doubles the machine load |
| BG78 | 52 | `bench worktree release` stops the grade of its worktree before it removes the worktree | planned TestReleaseStopsTheGrade in internal/worktree/background_stop_test.go | A grade outlives its worktree and its record |
| BG79 | 52 | `bench worktree clean` stops the grade of each worktree that it removes | planned TestCleanStopsTheGrade in internal/worktree/background_stop_test.go | A grade outlives a cleaned worktree |
| BG80 | 52 | `bench worktree reclaim` stops the grade of each worktree that it removes | planned TestReclaimStopsTheGrade in internal/worktree/background_stop_test.go | A grade outlives a reclaimed worktree |
| BG81 | 53 | A landing stack stops the grade of each member worktree before it grades | planned TestStackLandingStopsMemberGrades in internal/worktree/background_stop_test.go | A member grade runs beside the stack gate |
| BG82 | 54 | `bench cache clean` stops every grade and drops every queued grade, then cleans | planned TestCacheCleanStopsBackgroundGrades in internal/gocache/background_clean_test.go | A clean that refuses under a grade blocks for a whole low-priority run |
| BG83 | 55 | A stopped grade writes no green, and its background record names `stopped` | `internal/gate/background_queue_test.go` (`TestStoppedGradeRecordsStopped`) | A stop that records a partial run as green writes a false verdict |
| BG84 | 42 | `bench status` shows the latest background outcome of each worktree | planned TestStatusShowsTheBackgroundOutcome in internal/status/status_background_test.go | A status that reads only the gate cache hides each background red |
| BG85 | 43 | The next `bench commit` prints the latest background outcome and its log path | planned TestCommitPrintsTheLastBackgroundOutcome in internal/commit/background_test.go | A commit that prints nothing leaves the red for a search |
| BG99 | 41 | A worktree with a red background outcome and a green gate cache still shows the gate cache as its verdict | planned TestStatusKeepsTheGateVerdictBesideABackgroundRed in internal/status/status_background_test.go | A dashboard that reads the background record as a verdict gives a red authority |
| BG86 | 47 | `bench gate --checkpoint <spec> --chunk 1 --prepare` exits 2 with the usage line | planned TestCheckpointPrepareIsAUsageError in cmd/bench/gate_route_test.go | A surviving flag keeps a second mechanism |
| BG87 | 47 | The gate help row names no `--prepare` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A stale usage line advertises a removed flag |
| BG88 | 48 | The Land section of `/bench-implement-spec` holds no `--prepare` token | planned Forbid anchor in internal/anchors/registry_chunk_chain.go, `docs-currency-workflow` check | A surviving freeze sentence tells an orchestrator to run a removed flag |
| BG89 | 48 | The Land section holds no sentence that starts a prepared run after a repair commit | planned Forbid anchor in internal/anchors/registry_chunk_chain.go, `docs-currency-workflow` check | A surviving repair sentence orders a run that no longer exists |
| BG90 | 48 | The Land section holds no sentence that waits on the PID of a prepared run | planned Forbid anchor in internal/anchors/registry_chunk_chain.go, `docs-currency-workflow` check | A surviving sentence names a process that never starts |
| BG91 | 48 | The bounded repair policy holds no sentence that says a prepared chunk run excludes the review record | planned Forbid anchor in internal/anchors/registry_chunk_chain.go, `docs-currency-workflow` check | A surviving sentence describes a run that no longer exists |
| BG92 | 49 | The Land section states that a chunk checkpoint waits for the background grade of its key | review-owned: the Spec axis reads the Land section | A silent pause looks like a stuck checkpoint |
| BG93 | 50 | The census record holds five alone and five beside wall times of the foreground workload on four pinned cores | review-owned: the reviewer reads specs/background-grade/census.md | A census with fewer runs cannot show the spread |
| BG94 | 50 | The census record shows nice 19 and the idle I/O class for the proxy | review-owned: the reviewer reads specs/background-grade/census.md | A proxy at the default priority measures another bound |
| BG95 | 50 | The verdict line says pass only when the beside median is at most 1.2 times the alone median and no run is red | review-owned: the reviewer recomputes the medians from specs/background-grade/census.md | A pass on a failed bound adopts the overlap without its census |
| BG96 | 63 | An ADR states the background grade and amends ADR 0002 posture 5, ADR 0017, ADR 0024, and ADR 0031 | review-owned: the Spec axis reads the ADRs | A stale posture says that the gate serializes every run on one checkout lock |
| BG97 | 63 | `.bench/BENCH-reference.md` states the claim, the wait, the background record, and the queue | review-owned: the Spec axis reads the reference | A cold session cannot find why a checkpoint paused |
| BG98 | 63 | The profile's note on a running gate names a foreground gate | review-owned: the Spec axis reads `projects/benchkit.md` | A stale note forbids every commit while a background grade runs |
| BG100 | 48 | The bounded repair policy holds no sentence that orders the prepared run between the probe and the checkpoint | planned Forbid anchor in internal/anchors/registry_chunk_chain.go, `docs-currency-workflow` check | A surviving sentence orders a step that no longer exists |
| BG101 | 48 | `.bench/BENCH-reference.md` holds no `--prepare` token | planned Forbid anchor in internal/anchors/registry_chunk_chain.go, `docs-currency-workflow` check | A surviving reference sentence documents a removed flag |
| BG102 | 48 | `.bench/BENCH-reference.md` holds no phrase `prepared run` | planned Forbid anchor in internal/anchors/registry_chunk_chain.go, `docs-currency-workflow` check | A surviving reference sentence documents a removed run |
| BG103 | 48 | `.bench/BENCH-reference.md` holds no sentence that lists what a prepared run skips | planned Forbid anchor in internal/anchors/registry_chunk_chain.go, `docs-currency-workflow` check | A surviving list sentence describes a removed run |
| BG104 | 48 | `projects/benchkit.md` holds no sentence about a prepared chunk run | planned Forbid anchor in internal/anchors, `docs-currency-workflow` check | A surviving profile sentence describes a removed run |
| BG105 | 48 | ADR 0002 holds no sentence about a prepared chunk run | planned Forbid anchor in internal/anchors, `docs-currency-workflow` check | A surviving posture sentence describes a removed run |
| BG106 | 41 | A background run that ends green records `green` in the background record | `internal/gate/background_grade_test.go` (`TestBackgroundGreenRecordsGreen`) | A record that names only failures hides a finished green |
| BG107 | 41 | A background run that ends red records `red` in the background record | `internal/gate/background_grade_test.go` (`TestBackgroundRedRecordsRed`) | A red recorded as aborted hides that the oracle refused the tree |
| BG108 | 55 | After a stopped grade exits, no protection registration of that grade stays unresolved | `internal/gate/background_queue_test.go` (`TestStoppedGradeReleasesItsResources`) | An unresolved registration makes a later cache clean refuse |
| BG109 | 9 | The grade key of one clean kit fixture equals the key that the base commit's code gives | `internal/gate/run_binary_source_test.go` (`TestKitKeyIsUnchanged`) | A derivation change that moves the key makes every current green answer nothing |
| BG110 | 12 | A landing whose key differs from the running grade of its worktree does not wait | planned TestLandingOfAnotherKeyDoesNotWait in internal/worktree/land_background_test.go | A landing that waits on another key spends its timeout on a grade that it then stops |
| BG111 | 17 | A queued grade of the landing's key leaves the queue when the landing starts its own run | planned TestLandingDropsItsQueuedKey in internal/worktree/land_background_test.go | A queued duplicate grades the key again after the landing |
| BG112 | 18 | No code in `internal/gate/backgroundstate` sends a signal itself | the `process-lifetime-closure` check over the new package | A signal from the leaf package skips the lifetime owner and its descendant teardown |
| BG113 | 18 | The stop operation returns only after the owner process has exited | `internal/gate/backgroundstate/stop_test.go` (`TestStopWaitsForTheOwnerExit`) | A stop that returns at the signal lets two runs write one delta record |
| BG114 | 22 | A background chunk grade runs its oracle with the width input set | `internal/gate/background_grade_test.go` (`TestBackgroundGradeSetsTheWidthInput`) | A grade that drops the width input runs at the foreground width and crowds the host |
| BG115 | 35 | The grade key that a lane pass queues frames the executable that the wrapper resolves for the target worktree | `internal/gate/background_queue_test.go` (`TestLanePassKeyFramesTheTargetExecutable`) | A key framed by the launching executable misses its dedup and its claim lookup |
| BG116 | 18 | Before its own run, a checkpoint stops the other-key grade of its worktree even when another worktree wants that key | `internal/gate/background_wait_test.go` (`TestCheckpointStopsAWantedOtherKeyGrade`) | A kept grade and the checkpoint write one delta record at once |

### Edge inventory

The audience is this repository, because a linked repo starts no grade.
The absent store, the absent queue, and the absent claim each read as empty, and each one starts a grade.
A present but unreadable claim, queue, or record retains the grade and starts none, as ADR 0030 retains a resource.

- An unreadable claim file attaches to BG40, because the grade cannot prove its subject and records `aborted`.
- A worktree path with spaces attaches to the fixtures of BG5 and BG78, which use the existing spaced paths of `internal/worktree`.
- A pid that a new process reuses attaches to BG34.
- A tip whose tree holds no review record holds no open build, so its grade runs no chunk subject; a record with no chunk entry attaches to BG24.

Won't handle: a new grade when the destination moves — the landing then runs its own gate, as it does today, so the landing stays correct.
Won't handle: a background grade of the complete checkpoint subject — the complete checkpoint keeps its own run, and its green answers the landing.
Won't handle: a background grade in a linked repo — the kit keeps every background rule, and a parked idea holds the opt-in.
Won't handle: a CPU priority promotion — Linux refuses an unprivileged lower nice value, and the I/O promotion still applies.
Won't handle: the GR44 edge between two foreground runs — BG66 keeps today's refusal for that pair.

## Out of scope

- A new grade when the destination moves: 9 edits, 3 gate runs.
- A background grade of the complete checkpoint subject: 6 edits, 2 gate runs.
- A linked-repo opt-in declaration: 7 edits, 2 gate runs.
- The FT255 shared pool, the FT448 rerun cost, the FT372 nested-gate family, a landing stack as a consumer, a landing that reads a delta green, and a longer freshness window each belong to their own row or decision.

## Further notes

### Gap log

| class | entry |
|---|---|
| (a) map discretion | The claim store, the queue and slot storage, and the priority mechanism of each platform follow map decisions 6, 7, 13, and 16. |
| (a) map discretion | The form of the width input follows the precedent of the delta schedule entry, which is a Bench input outside the grade key. |
| (a) map discretion | The rule that classifies a kit worktree as an open build of a staged spec and derives its landing deliverable costs only a lookup miss when it is wrong, as BG41 and BG74 prove. |
| (a) map discretion | The stable identity value of a run with no run-binary selection is the author's choice under BG14. |
| (a) map discretion | The wording of the background outcome line and of its status row is the author's choice under BG84 and BG85. |
| (a) map discretion | Ticket 16 picks a new ADR or amendments of ADR 0002, ADR 0017, ADR 0024, and ADR 0031. |
| (b) seam or test choice | The state lives in the leaf package `internal/gate/backgroundstate`, so `internal/gocache` can stop grades with no import of `internal/gate`. |
| (b) seam or test choice | The stop of ticket 05 sends its signal through `internal/subprocess`, because the lifetime closure check refuses a signal from any other package. |
| (b) seam or test choice | The plumbing entry reuses an existing route line of `bin/bench.sh`, which is over its line budget. |
| (c) missing fact | Map unknown U1: a plan commit that changes only the spec, the tickets, or the review record moves no baseline member, because the settle compares no identity member. A commit that changes a Go build input moves the baseline runner digest, and ticket 03 reproduces that refusal through BG17 before its change. |
| (c) missing fact | Map unknown U3: ticket 01 measures it. |
| (c) missing fact | Map unknown U5: no record counts how many kit landings take the whole gate. BG-C4 still ships, because a miss costs only the run that the landing pays today. |
| (c) missing fact | Map unknown U6: the idle I/O class was accepted on this host on 2026-10-27, and ticket 01 records whether it takes effect. |
| (d) reviewer decision | none: the reviewer answered the CPU promotion question on 2026-10-27. |

### Overlaps with staged specs

The staged `repair-carry` spec of FT448 writes `.agents/commands/bench-implement-spec.md` and `CHANGELOG.md` in its ticket 05.
FT448 lands first in the approved order, so ticket 15 of this spec edits the Land section after that change.

### Flagged additions

- The background record of a worktree is an addition that map decision 9 implies. BG35 covers it.
- The Forbid anchors of BG88 to BG91 and BG100 to BG105 follow the retirement rule of the slicing checks. They are not in the map.

### Pre-review proof checklist

- Cited symbols: `runBinarySource`, `buildSubjectForTree`, `frameBaselineSchedule`, `hashProspectivePreparation`, `buildInputsIdentity`, `executeSubjectOwning`, `executeChunkCheckpoint`, `gradesSubjectAlone`, `runUnderBroker`, `published.Tree`, `deltaPlan`, `writeDeltaGreen`, `holdWholeGate`, `leaseLanding`, `KitForRoot`, `worktreeSerialCeiling`, and `processLifetimeCensus` resolve at `main` commit `0dbfa0147`.
- Import edges: the new leaf package imports no `internal/gate` code. `internal/gocache` imports only that leaf. `internal/commit` and `internal/worktree` already import `internal/gate`.
- Source-row clauses and occurrences: each map decision has a story group, and the four `roadmap/FT451.md` occurrences map to BG17, BG12, BG28, and the BG-C6 wait rows.
- Promised field labels: the outcomes `green`, `red`, `aborted`, and `stopped`. The claim fields are the key, the worktree, the owner pid, and the boot id.
- Changed-function callers: `runBinarySource` has three callers, in `gate.go`, `subject.go`, and `run_transaction.go`. `leaseLanding` has one caller, in `land.go`. `executeSubjectOwning` has callers in `engine.go` and `run_transaction.go`.
- Copy survival: BG13.
- Rendered-shape readers: the gate help row has one pin, in `TestHelpInventoryIsComplete`. The cache clean help row ends with "refuses under a live run", and ticket 13 owns its pin in the same test file.
- Pin operators: the help test compares whole lines.
- Entry reads: none beyond the existing gate entries.
- Derived expectations: each key expectation derives from the fixture tree through the consumer's own entry, not from the background run under test.
- Consolidated rules: BG13 retires the second kind derivation.
- Quantified obligations: each worktree has at most one queued grade, and the machine runs at most one grade.
- Workflow-step writes: the background run writes the store, the delta record, the reader table, and the background record. None of them is in a source digest that a later checkpoint compares.
