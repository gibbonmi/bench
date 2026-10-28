# Build cost

Status: staged

Roadmap: FT463

Decision source: reviewer-confirmed current conversation, 2026-10-28, with `roadmap/FT463.md` as the reviewed row

Verification log: 4 iterations to accept — the spec round took 2 (S1 picked a later landing because the rule named a state, and S2 left two live anchors outside the fence), and the ticket round took 2 (T1 left three fixture readers outside ticket 02, and T2 left the command-notes reader outside ticket 09). Both rounds ran on fable / high.

Implementation approval: approved by the reviewer on 2026-10-28 in the current conversation

## Problem

The FT451 build took 16.1 hours of wall time for 16 tickets. Its writers spent 26% of their time on record steps.
Each author ran and recorded each plan requirement after the ticket commit, which took 134 runs and 95.5 minutes in total.
The chunk checkpoint then ran the whole gate on the same subject, so each of those runs repeated a check that the checkpoint grades.
Of the 143 commits in the FT451 landing range, 113 changed only the spec folder or the review record.

The retro has no cost field, and the drain has no route that removes a step. Thus each build can add a step, and no build measures the steps that it pays for.
These sources hold the evidence: `capture/retros/background-grade.md` and `roadmap/FT463.md`.

Smaller costs repeat in each session. The success output of `bench commit` and of `bench probe --mutants` spills past the 10-line bound, and the session then reads the spill page.
A verb's own `--help` does not name `--in`, so a session uses `bench worktree exec`, which bounds the response twice.
A fresh kit worktree refuses each `--in` call until a `bench worktree build`, although the installed binary matches its build inputs.

## Solution

Each retro carries a derived `## Cost` section. The scaffold computes it from git, the landed plan, and the trace store, and it compares two per-ticket values with one budget.
An overrun line says `over budget`, and the drain then proposes one removal or one merge. A breach never makes the gate red.

A chunk requirement in a completion plan always names a probe. The author runs no other plan check after the ticket commit, because the chunk checkpoint grades the subject with the whole gate.
The `check-floor` row of build preflight retires. The witness reds, the probes, the lane pass, and the checkpoint grade stay.

After the ticket commit, the author runs the probe and then `bench record ticket`. That verb stages the retained probe response itself, writes the entries, and prints the exact record commit.
That commit is `bench commit --preflight-build <slug> --base <chunk-base>`, which builds the worktree and runs the build preflight in the same call.

A `bench commit` lane pass prints no check report, and a bitten mutant set prints no mutant report. Each tree-scoped verb's `--help` names `--in`.
A kit worktree target reuses a sealed binary whose seal matches its tree. The guidance names the `--in` forms, and a `Writes:` widen records no learning.

## User stories

Line: opus / medium.
Implementation-line reason: BC-C1 and BC-C3 change the plan grade and the record transaction that every chunk checkpoint reads, so a wrong refusal order or a lost entry stops later builds, this build included. The rows fix exact refusals and exact outputs, and existing record fixtures give strong tests. The other chunks change read-only renderers, guidance, and one executable choice with existing seams, so medium effort fits them.
Harder chunks: BC-C1, BC-C3.

### The retro measures what the build cost

1. As a retro author, I want `bench retro <slug> --scaffold` to print a `## Cost` section derived from git, the landed plan, and the trace store, so that each retro states the build cost without a hand count.
2. As a retro author, I want the cost section to name the landing commit and the build range that it measured, so that a reader can verify the selection.
3. As a retro author, I want the wall time from the first build commit to the landing commit, so that the retro states how long the build took.
4. As a retro author, I want the commits of the build range split into bookkeeping commits and change commits by their changed paths, so that the retro shows how much of the build was record work.
5. As a retro author, I want the ticket count, the longest ticket chain, the chunk count, and the serial checkpoint count of the landed plan, so that each retro measures the graph shape before anyone sets a slicing target.
6. As a retro author, I want the author and repair dispatch counts of the landed version 2 plan, so that the retro shows the repair load.
7. As a retro author, I want the commit, lane, and landing times of the build's assignment from the trace store, so that the retro shows what the checks cost.
8. As a retro author, I want each fact that the tree does not hold to read `unknown`, so that a missing measurement never blocks the draft and never reads as measured.
9. As a retro author, I want the cost section to read the landing commit after `bench spec retire` deleted the spec folder, so that the order of the post-merge tail does not lose the measurement.
10. As a drain operator, I want the per-ticket wall time and the per-ticket bookkeeping commits compared with one budget table, so that an overrun shows without a second copy of the numbers.
11. As a drain operator, I want each line above its budget to end with `over budget`, so that the drain finds each breach by one marker.
12. As a drain operator, I want the drain to propose one removal or one merge of a build step for each `over budget` line, so that the kit loses a step when a build costs too much.
13. As a kit maintainer, I want a budget breach never to make the gate red, so that a costly build still lands on green.
14. As a retro author, I want `bench retro <slug> --body` to refuse a body without the `## Cost` heading, so that every new retro carries the section.
15. As a kit maintainer, I want a pending retro without a `## Cost` heading to stay valid for the gate, so that the drain still reads the retros that exist now.

### The checkpoint grade is the only green evidence for a plain check

16. As an orchestrator, I want the plan grade to refuse a chunk requirement that names no probe, so that no author runs and records a check that the chunk checkpoint grades again.
17. As a spec author, I want a chunk with an empty `verification` list to pass the plan grade, so that a guidance-only chunk needs no placeholder requirement.
18. As a spec author, I want a chunk with no `verification` key to pass the plan grade, so that a chunk without a probe needs no empty list.
19. As a spec author, I want a version 2 ticket that owes no requirement to pass the plan grade, so that only probes bind a ticket.
20. As a ticket author, I want `bench record ticket` to freeze or refreeze the chunk for a ticket that owes no requirement, so that the record step works without a probe.
21. As a reviewer, I want each planned probe requirement still to need a recorded entry that bit and restored, so that the witness reds and the probes stay.
22. As a reviewer, I want the chunk checkpoint to keep its evidence pass and its whole-gate grade, so that the cut removes no gate.
23. As a reviewer, I want a final requirement to stay valid without a probe, so that the complete checkpoint keeps its `bench gate` and coverage evidence.
24. As an orchestrator, I want the `check-floor` row retired from build preflight, so that no plan pays for a floor that the checkpoint grade covers.
25. As an orchestrator, I want a staged spec with a plain chunk requirement to go red on the `completion-plan` row at its build entry, so that the staleness pass removes the requirement there.
26. As a ticket author, I want the build guidance to state that I run only the plan's probes after the ticket commit, so that I do not repeat the checkpoint's checks.
27. As a spec author, I want the completion-plan reference to state the probe-only chunk list, so that a new plan obeys the rule.

### One record step follows the ticket commit

28. As a ticket author, I want `bench record ticket` to stage the retained `bench probe` response when exactly one owed probe requirement has no stage, so that I make no separate excerpt call.
29. As a ticket author, I want `bench record ticket` to take `--mutant <n>` for that response, so that a mutant-set answer selects its row the way the excerpt form does.
30. As a ticket author, I want `bench record ticket` to refuse a retained response that is not a `bench probe` call and to name the excerpt form, so that a wrong response never becomes probe evidence.
31. As a ticket author, I want `bench record ticket` to refuse when two or more owed requirements have no stage, so that the map from a response to a requirement stays exact.
32. As a ticket author, I want `bench record ticket` to print the exact record commit with the chunk base, so that I make one more call and read no grammar.
33. As a ticket author, I want `bench commit --preflight-build <slug> --base <commit>` to pass that base to the chained build preflight, so that a moved default branch does not route a fold inside a chunk.
34. As a ticket author, I want `bench commit --base` without `--preflight-build` to refuse with the usage line, so that an ignored base never looks accepted.
35. As an orchestrator, I want the build guidance to name the post-ticket sequence of the probe, `bench record ticket`, and the printed record commit, so that the separate post-ticket preflight call goes away.
36. As a ticket author, I want the response of `bench record ticket` to fit in 10 lines, so that it never spills.

### Success output fits the bound

37. As a ticket author, I want a lane pass of `bench commit` to print no check report line, so that the success response fits in 10 lines.
38. As a ticket author, I want the complete check stream of each lane in a lane log, so that a passed report stays available.
39. As a ticket author, I want a lane red to print the failing check's own lines, so that a failure stays diagnosable in the response.
40. As a ticket author, I want a `bench probe --mutants` answer in which every row bit and restored to omit the baseline report and each mutant report, so that the answer stays short.
41. As a ticket author, I want a set with a row that did not bite or did not restore to keep every report, so that a failure stays diagnosable.
42. As a reviewer, I want the short answer to stay importable as a recorded probe entry, so that the probe evidence still records.
43. As a ticket author, I want a bitten set of four mutants to print at most 10 lines, so that the common probe never spills.

### A tree target costs one call

44. As an agent, I want each tree-scoped verb's `--help` to name `--in <label|primary>`, so that I find the tree-target form without the root help.
45. As an agent, I want a repository-scoped verb's `--help` to name no `--in`, so that the help never offers a refused form.
46. As an agent, I want a kit worktree target without a current build of its own to run the primary checkout's published build when that build's seal matches the target tree, so that a fresh worktree needs no build.
47. As an agent, I want such a target to run the running executable when its seal matches the target tree and the primary build does not, so that the installed binary serves a matching tree.
48. As an agent, I want the target's own current build to win over every other matching build, so that a worktree with Go edits runs its own registries.
49. As a reviewer, I want a target whose build inputs no candidate seal matches to keep the current refusal, so that a child never runs registries that differ from its tree.
50. As an agent, I want the guidance to name `bench <verb> --in <target>` for the evidence check, the worktree probe, and the batch gate, so that I pay one response bound and not two.
51. As an orchestrator, I want a `Writes:` widen to record no learning, so that a build pays no journal entry for a path that the widen rules already govern.
52. As a reviewer, I want every other plan or gate expansion still to record a learning, so that the drain still sees a split, a merge, or a gate change.

### Reviewed exclusions

53. As a reviewer, I want the dispatch-fact move, the narrow plan digest, and the charge at the ticket base to wait for a follow-on row, so that this spec stays small.
54. As a reviewer, I want the read-only exec child bound to wait for a follow-on row, so that its child classification gets its own design.
55. As a reviewer, I want the spec to set no slicing target, so that the target waits for cost sections that hold data.

## Implementation decisions

### Decisions closed on 2026-10-28

The reviewer closed these decisions in the current conversation:

- R1: the bookkeeping cut goes first, before FT402-git-cache.
- R5: a slicing target needs data, so this spec sets none.
- Scope: C, D1, D2, the Bench half of D6, the `--in` help, the `--in` guidance, and a widen without a learning. The rest of D5 and the R4 exec bound go to follow-on rows.
- R2: the chunk checkpoint grade is the only green evidence for a plain ticket check. A chunk plan keeps only probe requirements, and `check-floor` retires.
- R3: `bench record ticket` stays the record verb. It stages the retained probe response and prints the record commit, which `bench commit --preflight-build` already chains.
- Cost owner: the derived retro `## Cost` section. The budget starts at the FT451 values per ticket: 60 minutes of wall time and 7.1 bookkeeping commits.
- A removal of a duplicated step names its surviving owner and does not wait for FT231.
- FT341 spec A reopens for one change: a kit worktree target can reuse a sealed binary whose seal matches its tree.

These are the orchestrator's calls: R6 is no, so the cost section reads git and the trace store and no harness transcript. R7 is yes, so fresh repair sessions stay. FT461 and FT420 stay outside this spec.

### Landed work

FT451 already ships the D4 cut. A lane pass of `bench commit` starts a background grade of the chunk subject, and the chunk checkpoint reuses its green. This spec adds nothing for D4.
The lp-lean-author-charge light path already makes the charge name the verb grammar and the batched reads.

### Cost section

The cost section is a fourth derived section of the retro scaffold, beside the timings, the calibration table, and the repair table. The `## Cost` heading joins the canonical heading list directly after `## Outcome`.
The retro parser then demands it in each new body. The gate's retro grade reads only the destination markers, so a pending retro without the heading stays valid.

The landing commit is the newest first-parent commit of the default branch that changes `specs/<slug>/spec.md` to `Status: implemented`. Its tree holds that status, and the tree of its first parent does not. The build range is `<landing>^1..<landing>^2`.
One function in `internal/spec`, beside the status parser, owns this transition rule. No Bench helper names a spec's landing commit now, so the cost renderer must not derive it a second time in its own package.
The scaffold reads the plan and the tickets from the tree of the landing commit, so a retired spec folder loses no fact. A slug with no such commit prints `unknown` in each derived line.

The section derives these facts:

- the landing commit and the build range
- the wall time, from the committer time of the oldest non-merge commit in the range to the committer time of the landing commit
- the commit count of the range, split into bookkeeping commits and change commits
- the ticket count, the longest `Blocked by:` chain, the chunk count, and the longest chain of chunk dependencies
- the author and repair dispatch counts of a version 2 plan, where a repair is an assignment that names a predecessor
- the count and the summed elapsed time of the commit, lane, and landing spans whose trace holds an assignment ID of the plan
- the wall minutes per ticket and the bookkeeping commits per ticket, each with its budget

A bookkeeping commit is a non-merge commit whose every changed path is below `specs/<slug>/` or is `reviews/<slug>.md`. Every other non-merge commit is a change commit.
One table in the cost renderer holds the two budgets. A per-ticket value above its budget ends its line with `over budget`. No guidance file repeats a budget number.

### Probe-only chunk requirements

The plan grade refuses a chunk requirement without a `probe` key, and it names the chunk and the requirement ID. The final verification list keeps plain requirements.
A chunk's `verification` list can be empty or absent. A version 2 plan no longer needs a requirement for each ticket, but each listed chunk requirement still names a ticket of its chunk.

`bench record ticket` accepts a ticket that owes no requirement. It then writes the chunk step and no verification entry.
The shared plan fixtures of `internal/reviewrecord/recordtest` and `internal/preflight/preflighttest` hold plain chunk requirements now. The ticket that adds the refusal moves them, and each test that reads them, to probe requirements.
The checkpoint's evidence pass and its oracle stay as they are. The pass grades only the planned requirements, and those are now probes.

The `check-floor` row, its facts, and its guidance retire. The charge's root conformance line stays, because it names a development check and not a recorded requirement.
The twelve staged specs with plain chunk requirements change at their own build entry. The `completion-plan` row reds there, and the staleness pass of `/bench-implement-spec` removes each plain chunk requirement.

The cost renderer reads a landed plan through a decode-only reader that `internal/reviewrecord` exports beside `checkPlan`. A plan that the new grade refuses, such as the FT451 plan with its plain chunk requirements, still gives its graph and dispatch facts.

### Self-hosting order

This build changes the rules that grade its own plan. The current binary runs the `check-floor` row, which reds a chunk without its floor commands. The new binary refuses a plain chunk requirement.
Thus no single plan passes both binaries, and the order below keeps each step green:

1. At the spec stage, each chunk lists one probe requirement for each ticket and the plain floor commands that the current `check-floor` row names. The spec-stage `--plan-only` preflight and the first charge then pass on the current binary.
2. Ticket 01 lets the plan grade accept a probe-only plan and retires `check-floor` in one commit. The same commit deletes every plain chunk requirement of this plan, so its plan passes on the new binary.
3. Ticket 01's author records with the current binary. Its `bench record ticket` writes the plan-digest amendment itself, because its amendment step runs when the record digest differs from the plan at the tip.
4. Ticket 01's record commit uses plain `bench commit`. The chained form runs the build preflight in its own process, which is the current binary, so it would red on `check-floor`.
5. After ticket 01, the orchestrator installs the worktree build before the next dispatch, as the build guidance states. Each later preflight, charge, and checkpoint runs the new binary.
6. Ticket 02 then refuses a plain chunk requirement and moves the shared plan fixtures to probe requirements. This plan already holds none.

No author runs a plain floor command as a recorded requirement. Ticket 01's own record owes only its probe, because the plan at its commit holds no plain chunk requirement.
Ticket 01 runs its post-ticket preflight through `bench preflight --in <target>` after `bench worktree build <target>`, so that preflight runs the new binary.
Tickets 02 to 06 record with the excerpt form and the chained record commit, because ticket 07's `next=` line does not exist yet. Each ticket after ticket 07 obeys the sequence that tickets 03 and 07 state.
Before the first dispatch, the orchestrator converts this version 1 plan to version 2, as the build guidance states, because the ticket form reads a version 2 plan.

### One record step

`bench record ticket` stages the retained response for an owed probe requirement when exactly one owed requirement has no stage. It applies the checks of the excerpt form to that response, and it removes the stage after the write, as it does now.
A retained response that is not a `bench probe` call refuses with the excerpt form named. Two or more unstaged requirements keep the current refusal, which names the excerpt form.

`bench record ticket` takes `--mutant <n>` with the meaning that the excerpt form gives it. After its table, it prints one `next=` line with the record commit. The line reads `bench commit -m <msg> --preflight-build <slug> --base <chunk-base> reviews/<slug>.md`, with the slug and the chunk base filled.
`bench commit` takes `--base <commit>` only with `--preflight-build`, and the chain passes that base to the build preflight.

### Short success output

The lane writes each check's prefixed stream to the lane log `bench-lane.log` in the git directory of its worktree, and it truncates that log at each run. On a lane red, it prints the failing check's lines from that log.
On a lane pass, it prints no check line. The `lane{...}` line, the background line, and the `committed{...}` line keep their shapes.

A mutant set in which every row bit and restored prints its three tables with `trailing_bytes=0`, and it prints no report. Every other set keeps its current answer.

### Tree targets

The per-verb `--help` of each tree-scoped definition names `[--in <label|primary>]`. The text comes from the same scope field and the same operand as the root help.

A kit worktree target picks its child executable from three candidates in this order: its own published build, the primary checkout's published build, and the running executable. It runs the first candidate that `freshness.VerifyAt` accepts against the target tree.
When no candidate passes, the current refusals stay: `worktree build is missing` for an absent own build, and `worktree build does not match the tree` for every other case.
Each candidate's seal is the same trust root as today. The operator trusts the primary build and the running executable already, and the target's own build stays the candidate-controlled file that it is now.

### Guidance

The build guidance names the post-ticket sequence: the probe, `bench record ticket`, and the record commit that its `next=` line prints. The separate post-ticket `bench worktree exec` preflight step retires.
The known issues, the command notes, and the craft-synthesis skill name the `--in` forms. The known-issues row for the evidence check keeps its no-owner retire condition, which then names the `--in` form. `.bench/BENCH.md` keeps the learning for every plan or gate expansion except a `Writes:` widen.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| BC-C1 / 01-accept-probe-only-plans-and-retire-check-floor.md, 02-refuse-plain-chunk-requirements.md, 03-state-the-probe-only-rule.md | The checkpoint grade is the only green evidence for a plain check | BC23–BC37 | TestPlanProbeOnlyChunks, TestTicketFreezesWithoutRequirements, TestDecideHasNoCheckFloorRow, TestCommandBuildRedsAPlainChunkRequirement, TestBuildCostAnchors | yes |
| BC-C2 / 04-derive-the-retro-cost-section.md, 05-add-trace-times-and-the-drain-rule.md | The retro measures what the build cost | BC1–BC22, BC71, BC72 | TestRetroCostSection, TestRetroCostBudget, TestRetroCostTraceTimes, TestRetroRejectsABodyWithoutCost, TestBuildCostAnchors | no |
| BC-C3 / 06-pass-the-chunk-base-to-the-commit-chain.md, 07-stage-the-probe-in-record-ticket.md | One record step follows the ticket commit | BC38–BC48 | TestCommitChainPassesTheBase, TestParseArgsRefusesABaseWithoutTheChain, TestTicketStagesTheRetainedProbe, TestChunkChainAnchors | yes |
| BC-C4 / 08-print-a-short-lane-pass.md, 09-print-a-short-bitten-mutant-set.md | Success output fits the bound | BC49–BC57 | TestLanePassPrintsNoCheckReport, TestMutantSetSuccessOmitsTheReports, TestExcerptStagesAShortMutantSet | no |
| BC-C5 / 10-name-the-tree-target-in-verb-help.md, 11-reuse-a-matching-sealed-build.md, 12-name-the-tree-target-forms-in-guidance.md | A tree target costs one call | BC58–BC70, BC73 | TestTreeScopedHelpNamesTheTarget, TestRunReusesAMatchingBuild, TestTreeTargetRefusesStaleKitBuild, TestBuildCostAnchors | no |

The chunk graph has BC-C2 and BC-C3 after BC-C1, and BC-C5 after every other chunk. BC-C4 needs no other chunk.
BC-C2 and BC-C3 write no shared file, so they can author at the same time. The chunk-chain anchors of BC-C3 stay in `internal/anchors/registry_chunk_chain.go`, and the other anchors of this spec go to `internal/anchors/registry_build_cost.go`.
At slicing, each chunk also lists the plain floor commands that the current `check-floor` row names, as the self-hosting order states. Ticket 01's commit deletes them.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft463-build-cost-full-20261028","orchestrator_session":"claude:bc-orchestrator-20261028","author_limit":3,"assignments":{"01-accept-probe-only-plans-and-retire-check-floor.md":[{"session":"claude:bc-t01-author-20261028","assignment":"65a705264b8f1022b3ed0107187e8c56","model":"opus","effort":"high","source":"fe44357351b06f65518a80532f5f2b3574487807","native_ref":"claude-agent:bench-writer-high:bc-t01-author-20261028"}],"02-refuse-plain-chunk-requirements.md":[{"session":"claude:bc-t02-author-20261028","assignment":"65a705264b8f1022b3ed0107187e8c56","model":"opus","effort":"high","source":"4fafbadfdc47b69c93a475ffbe773567e24e47c2","native_ref":"claude-agent:bench-writer-high:bc-t02-author-20261028"}],"03-state-the-probe-only-rule.md":[{"session":"claude:bc-t03-author-20261028","assignment":"5499d93b2a59d3fc6ab8be589c71a812","model":"opus","effort":"high","source":"ba8a0568a3896dfb5c7d58b6e3d5705aec8fb570","native_ref":"claude-agent:bench-writer-high:bc-t03-author-20261028"}],"04-derive-the-retro-cost-section.md":[],"05-add-trace-times-and-the-drain-rule.md":[],"06-pass-the-chunk-base-to-the-commit-chain.md":[],"07-stage-the-probe-in-record-ticket.md":[],"08-print-a-short-lane-pass.md":[],"09-print-a-short-bitten-mutant-set.md":[],"10-name-the-tree-target-in-verb-help.md":[],"11-reuse-a-matching-sealed-build.md":[],"12-name-the-tree-target-forms-in-guidance.md":[]}},"chunks":[{"id":"BC-C1","tickets":["01-accept-probe-only-plans-and-retire-check-floor.md","02-refuse-plain-chunk-requirements.md","03-state-the-probe-only-rule.md"],"verification":[{"id":"probe-01","command":"bench probe internal/reviewrecord/plan.go --swap 'if len(items) == 0 && owned == nil {' --with 'if len(items) == 0 {' --package ./internal/reviewrecord --run TestPlanProbeOnlyChunks","probe":"Refuse an empty chunk list again. The probe-only plan with an empty chunk list must then refuse, so the plan test fails.","ticket":"01-accept-probe-only-plans-and-retire-check-floor.md"},{"id":"probe-02","command":"bench probe internal/reviewrecord/plan.go --swap 'if owned != nil && item.Probe == \"\" {' --with 'if false {' --package ./internal/reviewrecord --run TestPlanProbeOnlyChunks","probe":"Disable the plain chunk requirement refusal. The plan with a plain chunk requirement must then pass, so the plan test fails.","ticket":"02-refuse-plain-chunk-requirements.md"}]},{"id":"BC-C2","tickets":["04-derive-the-retro-cost-section.md","05-add-trace-times-and-the-drain-rule.md"],"verification":[{"id":"probe-04","command":"bench probe internal/retros/retros.go --omit 'CostHeading,' --package ./internal/roadmap --run TestRetroRejectsABodyWithoutCost","probe":"Remove the Cost heading from the canonical heading list. The body without a Cost section must then pass, so the refusal test fails.","ticket":"04-derive-the-retro-cost-section.md"},{"id":"probe-05","command":"bench probe internal/roadmap/retro_cost.go --swap 'if ids[span.Attributes[otelrecord.AttrAssignmentID]] {' --with 'if true {' --package ./internal/roadmap --run TestRetroCostTraceTimes","probe":"Count every span of the record. The span of another assignment must then reach the sums, so the trace test fails.","ticket":"05-add-trace-times-and-the-drain-rule.md"}]},{"id":"BC-C3","tickets":["06-pass-the-chunk-base-to-the-commit-chain.md","07-stage-the-probe-in-record-ticket.md"],"verification":[{"id":"probe-06","command":"bench probe cmd/bench/commit_chain.go --omit 'args = append(args, evidencecmd.FlagBase, outcome.Base)' --package ./cmd/bench --run TestCommitChainPassesTheBase","probe":"Drop the base from the chained preflight arguments. The chain test must then see no base.","ticket":"06-pass-the-chunk-base-to-the-commit-chain.md"},{"id":"probe-07","command":"bench probe internal/reviewrecord/recordcmd/ticket.go --swap 'if len(unstaged) == 1 && retained.Verb == probe.Verb {' --with 'if false {' --package ./internal/reviewrecord/recordcmd --run TestTicketStagesTheRetainedProbe","probe":"Disable the staging of the retained probe. The ticket call must then refuse with the excerpt form, so the staging test fails.","ticket":"07-stage-the-probe-in-record-ticket.md"}]},{"id":"BC-C4","tickets":["08-print-a-short-lane-pass.md","09-print-a-short-bitten-mutant-set.md"],"verification":[{"id":"probe-08","command":"bench probe internal/gate/lane.go --swap 'laneWriters(log, log, diagnostics)' --with 'laneWriters(req.Stdout, req.Stderr, diagnostics)' --package ./internal/gate --run TestLanePassPrintsNoCheckReport","probe":"Stream the check reports to the response again. The pass response must then hold check lines, so the lane test fails.","ticket":"08-print-a-short-lane-pass.md"},{"id":"probe-09","command":"bench probe internal/probe/set.go --swap 'if s.bitten() {' --with 'if false {' --package ./internal/probe --run TestMutantSetSuccessOmitsTheReports","probe":"Keep the reports of a bitten set. The answer must then hold a mutant_report table, so the set test fails.","ticket":"09-print-a-short-bitten-mutant-set.md"}]},{"id":"BC-C5","tickets":["10-name-the-tree-target-in-verb-help.md","11-reuse-a-matching-sealed-build.md","12-name-the-tree-target-forms-in-guidance.md"],"verification":[{"id":"probe-10","command":"bench probe cmd/bench/tree_scope.go --swap 'if definition.Scope != scopeTree {' --with 'if true {' --package ./cmd/bench --run TestTreeScopedHelpNamesTheTarget","probe":"Give no definition the tree-target text. The verb help must then omit --in, so the help test fails.","ticket":"10-name-the-tree-target-in-verb-help.md"},{"id":"probe-11","command":"bench probe internal/treetarget/build.go --swap 'if freshness.VerifyAt(root, root, candidate) == nil {' --with 'if candidate == freshness.PublishedExecutable(root) && freshness.VerifyAt(root, root, candidate) == nil {' --package ./internal/treetarget --run TestRunReusesAMatchingBuild","probe":"Accept only the target's own build. A target without its own build must then refuse, so the reuse test fails.","ticket":"11-reuse-a-matching-sealed-build.md"},{"id":"probe-12","command":"bench probe .bench/BENCH.md --omit ' except a `Writes:` widen' --check docs-currency-workflow","probe":"Restore the learning for every widen. The Forbid anchor on the retired wording must then fail.","ticket":"12-name-the-tree-target-forms-in-guidance.md"}]}],"final_verification":[{"id":"whole-project","command":"bench gate"},{"id":"coverage","command":"bench coverage --check specs/build-cost/spec.md"}]}
```

## Testing decisions

A good test drives the real verb or the real owner with an exact fixture, and it compares an exact output or an exact refusal. It never compares two real runs, because their elapsed times differ.
The cost section uses canned git fixtures and a canned trace record. Its arithmetic rows state their inputs.

- The retro scaffold tests in `internal/roadmap` are the precedent for the cost section. They already drive `scaffoldBody` with a canned home and root.
- The plan grade tests attach at `checkPlan` through `ReadPlan` on a committed fixture tree. The record fixtures in `internal/reviewrecord/recordtest` build the chunk and the record.
- The `bench record ticket` tests attach at `recordcmd.Command`, with a retained response that `responsebound` writes, as the excerpt tests do now.
- The commit chain tests replace one chain step, as `cmd/bench/commit_chain_test.go` does now.
- The lane tests run `RunLane` on a fixture checkout, as `internal/gate/lane_run_test.go` does now.
- The mutant-set tests use the probe test fixtures of `internal/probe`.
- The tree-target tests use the marker scripts and `freshness.Publish` of `internal/treetarget/build_test.go`.
- The guidance rows are anchors. The `docs-currency-workflow` check grades them against the real files.

The gate observes every row through its package tests and the `docs-currency-workflow` check. No row adds a gate check.

### Seam diagram

    trigger: bench retro <slug> --scaffold
        │
        ▼
    git history + landed plan + trace record  ──▶  [ cost renderer in internal/roadmap ]  ──▶  ## Cost lines
                      ◀ tests attach here: a canned repository and home, then an exact compare of the printed lines

    trigger: bench record ticket <slug> … after bench probe
        │
        ▼
    retained probe response + plan at the tip  ──▶  [ recordcmd ticket form ]  ──▶  record entry + next= record commit
                      ◀ tests attach here: recordcmd.Command on a record fixture, then the record bytes and the response lines

    trigger: bench <verb> --in <label>
        │
        ▼
    target tree + candidate executables  ──▶  [ treetarget childExecutable ]  ──▶  child argv or refusal
                      ◀ tests attach here: marker scripts sealed by freshness.Publish, then the marker that ran

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| BC1 | 1 | The scaffold prints the `## Cost` heading directly after the `## Outcome` heading, followed by derived lines | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | The current scaffold prints no Cost heading |
| BC2 | 2 | For a fixture landing L, the first cost line names L and the range `L^1..L^2` | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | A section that measured the newest landing of another slug names another commit |
| BC72 | 2 | When the landing L of the slug comes before a landing of another slug and then the retirement of the slug, the first cost line still names L | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | A rule that takes the newest commit whose tree holds the status names the other landing |
| BC3 | 3 | For a range whose oldest non-merge commit is 960 minutes older than L, the wall line reads 960 min | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | A start at the first parent of L or at a folded default-branch commit gives another number |
| BC4 | 4 | For a range of three commits that touch only `specs/<slug>/` or `reviews/<slug>.md` and two commits that touch `internal/x.go`, the commits line reads 5 total, 3 bookkeeping, and 2 change | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | A subject-based count or a count of folded commits differs |
| BC5 | 4 | A commit that touches `specs/<slug>/spec.md` and `internal/x.go` counts as a change commit | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | An any-path rule counts it as bookkeeping |
| BC6 | 5 | For a plan with tickets a, b, c, and d, where b is blocked by a and c by b, and chunks X of a and b, Y of c, and Z of d, the graph line reads 4 tickets, chain 3, 3 chunks, and 2 serial checkpoints | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | A count of tickets or of chunks in place of the longest chain gives 4 or 3 |
| BC7 | 6 | For a version 2 plan with five assignments of which two name a predecessor, the dispatch line reads 3 authors and 2 repairs | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | A count of every assignment as an author gives 5 |
| BC8 | 8 | A version 1 plan prints `unknown` for both dispatch counts | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | A zero reads as a measured count |
| BC9 | 8 | A slug with no implemented landing in history prints `unknown` in each derived cost line and exits 0 | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | A refusal blocks the draft, and a zero misleads it |
| BC10 | 9 | After a later commit deletes `specs/<slug>/`, the scaffold prints the same graph line from the tree of L | planned TestRetroCostSection in internal/roadmap/retro_cost_test.go | A read of the working tree prints unknown |
| BC11 | 10 | For 20 tickets and 1,300 wall minutes, the wall line reads 65.0 min per ticket against 60 min and ends with `over budget` | planned TestRetroCostBudget in internal/roadmap/retro_cost_test.go | A missing comparison prints no marker |
| BC12 | 11 | For 20 tickets and 140 bookkeeping commits, the bookkeeping line reads 7.0 per ticket against 7.1 and carries no `over budget` | planned TestRetroCostBudget in internal/roadmap/retro_cost_test.go | A marker on every line, or a marker at or below the budget, marks it |
| BC13 | 11 | For 20 tickets and 144 bookkeeping commits, the bookkeeping line reads 7.2 per ticket and ends with `over budget` | planned TestRetroCostBudget in internal/roadmap/retro_cost_test.go | A wall-only comparison leaves this breach without a marker |
| BC14 | 10 | `bench-drain.md` holds neither budget number | planned TestBuildCostAnchors in internal/anchors/registry_build_cost_test.go | A second copy of a budget drifts from the renderer table |
| BC15 | 14 | `bench retro <slug> --body` refuses a body without `## Cost`, names that heading, and writes no file | planned TestRetroRejectsABodyWithoutCost in internal/roadmap/retro_test.go | A heading list without Cost accepts the body |
| BC16 | 1 | The scaffold draft with the Cost section passes the retro parser | `internal/roadmap/retro_scaffold_test.go` (`TestRetroScaffoldParses`) | A Cost heading at another position fails the parser |
| BC17 | 15 | The retro grade returns no diagnostic for a pending retro without `## Cost` | planned TestImprovementMarkersIgnoreTheCostHeading in internal/retros/recommendations_test.go | A grade that called the parser reds each pending retro |
| BC18 | 13 | The retro grade returns no diagnostic for a retro whose cost lines end with `over budget` | planned TestImprovementMarkersIgnoreTheCostHeading in internal/retros/recommendations_test.go | A grade that reads the marker reds a costly build |
| BC19 | 7 | For a canned record with two commit spans, one lane span, and one landing span in traces of the plan's assignment, the trace line states those counts and their summed milliseconds | planned TestRetroCostTraceTimes in internal/roadmap/retro_cost_test.go | A sum over every span of the record counts another build |
| BC20 | 7 | A commit span of another assignment adds nothing to the trace line | planned TestRetroCostTraceTimes in internal/roadmap/retro_cost_test.go | A filter on the seam alone counts it |
| BC21 | 8 | An absent trace record prints `unknown` in the trace line | planned TestRetroCostTraceTimes in internal/roadmap/retro_cost_test.go | A zero reads as measured |
| BC22 | 12 | Review-owned: Section 4 of `bench-drain.md` holds the sentence "When a retro's `## Cost` section marks a line `over budget`, propose one removal or one merge of a build step for that line." | review-owned | A drain that lists the breach and proposes no removal passes |
| BC71 | 13 | Review-owned: Section 4 of `bench-drain.md` holds the sentence "A budget breach never makes the gate red." | review-owned | A drain that blocks a landing on cost passes |
| BC23 | 16 | The plan grade refuses a chunk requirement without `probe` and names the chunk ID and the requirement ID | planned TestPlanProbeOnlyChunks in internal/reviewrecord/plan_probe_test.go | The current grade accepts it |
| BC24 | 17 | The plan grade accepts a chunk whose `verification` list is empty | `internal/reviewrecord/plan_probe_test.go` (`TestPlanProbeOnlyChunks`) | The current grade refuses the empty inventory |
| BC25 | 18 | The plan grade accepts a chunk without a `verification` key | `internal/reviewrecord/plan_probe_test.go` (`TestPlanProbeOnlyChunks`) | A decode that demands the key refuses it |
| BC26 | 19 | A version 2 plan in which one ticket owes no requirement passes the plan grade | `internal/reviewrecord/plan_probe_test.go` (`TestPlanProbeOnlyChunks`) | The current grade refuses it as a ticket without a requirement |
| BC27 | 19 | A version 2 chunk requirement that names a ticket outside its chunk still refuses | `internal/reviewrecord/plan_probe_test.go` (`TestPlanProbeOnlyChunks`) | A relaxed grade that drops the ticket check accepts it |
| BC28 | 23 | The plan grade accepts a final requirement without `probe` | `internal/reviewrecord/plan_probe_test.go` (`TestPlanProbeOnlyChunks`) | A refusal in both lists breaks the complete checkpoint |
| BC29 | 20 | `bench record ticket` for a ticket that owes no requirement exits 0, freezes the chunk at the tip, and prints an empty `verifications` cell | `internal/reviewrecord/recordcmd/ticket_probe_test.go` (`TestTicketFreezesWithoutRequirements`) | The current verb refuses because the requirements name no ticket |
| BC30 | 21 | A chunk checkpoint whose planned probe requirement has no entry refuses with the verification form named | `internal/reviewrecord/plan_probe_test.go` (`TestPlanProbeOnlyChunks`) | A grade that skips probe requirements too passes it |
| BC31 | 22 | Review-owned: the diff leaves the evidence pass and the oracle selection of the chunk checkpoint unchanged in `internal/gate/checkpoint.go` | review-owned | A cut that also skips the evidence pass removes a gate |
| BC32 | 24 | `bench preflight build` prints no `check-floor` row for any plan | `internal/preflight/charge_cells_test.go` (`TestDecideHasNoCheckFloorRow`) | A kept row reds each probe-only plan |
| BC33 | 25 | A staged-spec fixture with one plain chunk requirement reds the `completion-plan` row of `bench preflight build` and names the requirement | planned TestCommandBuildRedsAPlainChunkRequirement in internal/preflight/completion_plan_test.go | A preflight that skips the plan grade stays green |
| BC34 | 26 | Review-owned: The Build section of `bench-implement-spec.md` holds the sentence "After the ticket commit, the author runs only the probes of the plan, because the chunk checkpoint grades every other check." | review-owned | The old rule survives without the new one |
| BC35 | 26 | `bench-implement-spec.md` holds neither "then runs each verification requirement once at that commit" nor "`check-floor`" | planned TestBuildCostAnchors in internal/anchors/registry_build_cost_test.go | A surviving sentence keeps the recorded runs and routes a retired row |
| BC36 | 27 | Review-owned: `completion-plan.md` holds the sentence "Each chunk requirement names a `probe`, and a chunk can list no requirement." | review-owned | A plan author keeps writing plain chunk requirements |
| BC37 | 24 | `completion-plan.md` holds no "`check-floor`" | planned TestBuildCostAnchors in internal/anchors/registry_build_cost_test.go | The reference routes a retired row |
| BC38 | 33 | `bench commit --preflight-build s --base B` passes `--base B` and `--source-tip <published>` to the preflight step | planned TestCommitChainPassesTheBase in cmd/bench/commit_chain_test.go | A chain that drops the flag passes no base |
| BC39 | 34 | `bench commit --base B -m m p` without `--preflight-build` prints the usage line at exit 2 and publishes nothing | planned TestParseArgsRefusesABaseWithoutTheChain in internal/commit/commit_test.go | A parser that accepts it ignores the base |
| BC40 | 33 | `bench commit --help` names `--base <commit>` in the `--preflight-build` group | planned TestCommitChainPassesTheBase in cmd/bench/commit_chain_test.go | A flag that the help omits sends the author to the code |
| BC41 | 28 | With one owed probe requirement unstaged and a retained `bench probe` single-form response that bit, `bench record ticket` writes that requirement's entry with the response bytes and removes no other stage | planned TestTicketStagesTheRetainedProbe in internal/reviewrecord/recordcmd/ticket_probe_test.go | The current verb refuses with the excerpt form |
| BC42 | 29 | With `--mutant 2` over a retained mutant-set response, the written entry records the mutation of row 2 | planned TestTicketStagesTheRetainedProbe in internal/reviewrecord/recordcmd/ticket_probe_test.go | A verb that reads row 1 records the wrong mutation |
| BC43 | 30 | A retained `bench test` response refuses, names `bench record excerpt`, and leaves the record bytes unchanged | planned TestTicketStagesTheRetainedProbe in internal/reviewrecord/recordcmd/ticket_probe_test.go | A verb that stages any response records a test run as probe evidence |
| BC44 | 31 | With two owed probe requirements unstaged, the verb refuses, names the excerpt form, and leaves the record bytes unchanged | planned TestTicketStagesTheRetainedProbe in internal/reviewrecord/recordcmd/ticket_probe_test.go | A verb that maps the response to the first requirement records the wrong probe |
| BC45 | 32 | After the ticket table, the verb prints `next=bench commit -m <msg> --preflight-build <slug> --base <chunk-base> reviews/<slug>.md` with the slug and the chunk base of the written entry | planned TestTicketStagesTheRetainedProbe in internal/reviewrecord/recordcmd/ticket_probe_test.go | A missing line or a line without the base sends the author to the help |
| BC46 | 36 | The response of a ticket call that writes three entries is at most 10 lines | planned TestTicketStagesTheRetainedProbe in internal/reviewrecord/recordcmd/ticket_probe_test.go | A response that lists each entry on its own lines spills |
| BC47 | 35 | Review-owned: The Build section of `bench-implement-spec.md` holds the sentence "After the ticket commit, run the probe, then `bench record ticket`, then the record commit that its `next=` line prints." | review-owned | The author keeps the separate preflight and excerpt calls |
| BC48 | 35 | `bench-implement-spec.md` holds no "`bench worktree exec <target> -- bench preflight build`" | planned TestChunkChainAnchors in internal/anchors/registry_chunk_chain_test.go | The retired exec preflight step survives |
| BC49 | 37 | A lane pass of `bench commit` over the prose and docs-currency-workflow checks prints no line that starts with `[prose]` or `[docs-currency-workflow]` | planned TestLanePassPrintsNoCheckReport in internal/gate/lane_report_test.go | The current lane streams each report to the response |
| BC50 | 38 | After that pass, `bench-lane.log` in the git directory holds the `[prose]` and `[docs-currency-workflow]` lines | planned TestLanePassPrintsNoCheckReport in internal/gate/lane_report_test.go | A pass that discards the reports loses them |
| BC51 | 39 | A lane red prints the failing check's own lines | `internal/gate/lane_run_test.go` (`TestLaneRelaysARedCheckOwnLines`) | A lane that hides every report hides the failure |
| BC52 | 38 | A second lane run truncates the log, so the log holds no line of the first run | planned TestLanePassPrintsNoCheckReport in internal/gate/lane_report_test.go | An appending log grows without bound and mixes runs |
| BC53 | 40 | A mutant set in which every row bit and restored prints the three tables, `trailing_bytes=0`, and no `mutant_report` table | planned TestMutantSetSuccessOmitsTheReports in internal/probe/set_report_test.go | The current answer prints each report |
| BC54 | 41 | A set with one silent row prints the report of each started mutant | planned TestMutantSetSuccessOmitsTheReports in internal/probe/set_report_test.go | A rule that drops reports for every set hides the failure |
| BC55 | 41 | A set with a failed restore keeps the preserved table | planned TestMutantSetSuccessOmitsTheReports in internal/probe/set_report_test.go | A rule keyed on the bite alone drops the preserved copy |
| BC56 | 42 | `bench record excerpt --mutant 1` stages a canned short answer, and `bench record ticket` writes a probe entry that bit and restored | planned TestExcerptStagesAShortMutantSet in internal/reviewrecord/recordcmd/excerpt_short_test.go | An import that demands the reports refuses the short answer |
| BC57 | 43 | A bitten set of four mutants prints at most 10 lines | planned TestMutantSetSuccessOmitsTheReports in internal/probe/set_report_test.go | A table layout of more than one line for each row spills |
| BC58 | 44 | `bench commit --help` names the `--in` flag with the operand that the root help shows | planned TestTreeScopedHelpNamesTheTarget in cmd/bench/tree_scope_test.go | The current verb help prints no tree target |
| BC59 | 44 | For each tree-scoped definition and leaf in the command registry, its `--help` output names the `--in` flag with that operand | planned TestTreeScopedHelpNamesTheTarget in cmd/bench/tree_scope_test.go | A hard-coded verb list misses a definition |
| BC60 | 45 | For each repository-scoped definition in the command registry, its `--help` output names no `--in` | planned TestTreeScopedHelpNamesTheTarget in cmd/bench/tree_scope_test.go | A rule that adds the text to every verb offers a refused form |
| BC61 | 46 | A kit target without its own build runs the primary checkout's sealed build when that seal matches the target tree | planned TestRunReusesAMatchingBuild in internal/treetarget/build_test.go | The current rule refuses with `worktree build is missing` |
| BC62 | 47 | A kit target without its own build runs the running executable when its seal matches and the primary build's seal does not | planned TestRunReusesAMatchingBuild in internal/treetarget/build_test.go | A rule that stops after the primary build refuses |
| BC63 | 48 | A kit target whose own build matches runs its own build when the primary build also matches | planned TestRunReusesAMatchingBuild in internal/treetarget/build_test.go | A rule that prefers the primary build runs other registries |
| BC64 | 49 | A kit target with an edited build input and no matching candidate refuses `worktree build does not match the tree` and starts no child | `internal/systemtest/tree_target_test.go` (`TestTreeTargetRefusesStaleKitBuild`) | A reuse without the seal check starts a stale child |
| BC65 | 49 | A kit target without its own build whose candidates all differ refuses `worktree build is missing` and starts no child | planned TestRunReusesAMatchingBuild in internal/treetarget/build_test.go | A reuse that falls back to the wrapper starts a stale child |
| BC66 | 50 | `docs/known-issues.md` holds "Run it as `bench preflight --in <target> evidence <id> --check-current`." | planned TestBuildCostAnchors in internal/anchors/registry_build_cost_test.go | The double-bound exec form survives |
| BC73 | 50 | `docs/known-issues.md` holds no "names the exec form" | planned TestBuildCostAnchors in internal/anchors/registry_build_cost_test.go | The retire column still routes the reader to the exec form |
| BC67 | 50 | `.bench/BENCH-reference.md` holds "For a worktree, run `bench probe --in <target> ...`." and no "`bench worktree exec <target> -- bench probe ...`" | planned TestBuildCostAnchors in internal/anchors/registry_build_cost_test.go | The stale claim that the probe takes no target survives |
| BC68 | 50 | The craft-synthesis skill holds "Take the green verdict from the whole-project gate. That gate is the landing's gate, or `bench gate --in <target>` for a batch that waits for approval." | planned TestBuildCostAnchors in internal/anchors/registry_build_cost_test.go | The exec gate form survives |
| BC69 | 51 | `.bench/BENCH.md` holds no "record every plan or gate expansion with `bench learning`" | planned TestBuildCostAnchors in internal/anchors/registry_build_cost_test.go | The learning for each widen survives |
| BC70 | 52 | Review-owned: `.bench/BENCH.md` still holds "Before dispatch, record every plan or gate expansion" | review-owned | A deletion of the whole clause drops the learning for a split or a gate change |

Not covered: story 53 — the follow-on row for the rest of D5 owns it, and this spec changes no dispatch fact, digest, or charge base.
Not covered: story 54 — the follow-on row for the R4 exec bound owns it, and this spec changes no exec bound.
Not covered: story 55 — the spec sets no target, and the measurement under Further notes is data, not a check.

### Edge inventory

The walk covers the cost section, the plan grade, the ticket record, the commit chain, the lane log, the mutant set, the verb help, and the child executable.
The cost section serves this repository and every linked repository that writes retros. The other changes serve the Bench executable in every linked repository, except the kit worktree reuse, which serves only this repository.

Each in-scope edge has a row:

- An absent landing, a version 1 plan, an absent trace record, and a retired spec folder each print `unknown` or read the landing tree (BC8, BC9, BC10, BC21).
- A commit with mixed paths counts as change (BC5). A fold commit stays outside the range (BC3).
- A per-ticket value at or below its budget carries no marker (BC12).
- An empty and an absent chunk list are two rows (BC24, BC25). A chunk requirement outside its chunk still refuses (BC27).
- A non-probe response and two unstaged requirements each refuse, and the record bytes stay unchanged (BC43, BC44).
- A base without the chain refuses (BC39).
- A lane red keeps its lines (BC51), and a second run truncates the log (BC52).
- A silent row and a failed restore keep the reports and the preserved table (BC54, BC55).
- A repository-scoped verb's help gets no text (BC60).
- A changed build input and an all-different candidate set keep the two refusals (BC64, BC65).

The hostile-input classes reach these surfaces:

- A hostile slug reaches the cost section only after `retros.ValidSlug` accepts it, as now.
- A changed path with a newline or a quote reaches the bookkeeping rule. The rule reads `git diff-tree -z` output, so each path is one record.
- A retained response is a file under the Bench home that a Bench verb wrote. The ticket form applies the excerpt form's bound and its import checks to it.
- A candidate executable that is a link or a special file fails `freshness.VerifyAt`, because the seal read follows no link.

**Won't handle** lines:

- The twelve staged specs with plain chunk requirements — each one's build-entry staleness pass removes them when the `completion-plan` row reds.
- A cost line for delegate model time or for each role — that needs harness transcripts, and R6 is no.
- The spans of a background grade — such a grade runs in its own trace without an assignment ID, and the landing span already holds the time that a landing waited for it.
- The `bench worktree exec <target> -- bench gate` help row — FT341 spec B owns the refusal of a Bench child under exec and the removal of that row.
- The charge's root conformance line in the delegation discipline — it names a development check, and the anchor "The root conformance pass is `bench test --check root-conformance`." stays.
- The commit and merge success output of a project with no declared lane — the gate prints that output, and the gate response keeps its current shape.

## Out of scope

- Dispatch facts move to the review record, the plan digest covers only the approved plan, a dispatch takes no plan commit, and the charge binds the ticket base. This is about 12 edits and 4 gate runs.
- A read-only exec child gets a 60-line, 8 KB bound. This is about 5 edits and 2 gate runs.
- A graph-depth rule and parallel dispatch by default, after FT420. This is about 6 edits and 2 gate runs, and it needs cost sections from several builds first.
- FT461, the close landing as a second whole gate, stays its own row.

## Further notes

### Measured graph shape

Before it set any slicing target, the spec measured the 23 specs that a `spec-retire` commit removed from 2026-10-12 to 2026-10-28. The script read each spec and its tickets at the parent of the retirement commit.
It computed the longest `Blocked by:` chain and the longest chain of chunk dependencies.

| measure | value |
|---|---|
| tickets | 161 |
| tickets on the longest chain, summed | 138 (86%) |
| chunks | 108 |
| chunks on the longest chain, summed | 99 (92%) |
| FT451 | 16 tickets, chain 10, 9 chunks, 7 serial |
| landing-gate-long-poles | 16 tickets, chain 14, 10 chunks, 9 serial |
| chunk-repair-cycle | 12 tickets, chain 12, 10 chunks, 10 serial |

Recent builds are almost fully serial, so a parallel author limit saves little. This spec's own graph has 12 tickets, a longest chain of 5, 5 chunks, and 3 serial checkpoints.
The FT451 landing range `55c699286..4cd081de8` holds 143 non-merge commits, of which 113 are bookkeeping by the changed-path rule. Its wall time is 964.8 minutes.
Thus FT451 measures 60.3 minutes and 7.06 bookkeeping commits per ticket against the budgets.

### Source and caller inventory

The changed owners and their callers are these:

- `retros.RequiredHeadings` serves `scaffoldBody` and the anchor in `internal/anchors/registry_retained_workflow.go`. `retros.Parse` serves the retro writer.
- `reviewrecord.checkPlan` serves `ReadPlan`, `RecordAssignment`, and the plan readers in `internal/preflight/plan.go`, `internal/preflight/ownership/premises.go`, and `internal/gate/chunk_gap.go`.
- `reviewrecord.requirementsValid` serves `checkPlan` alone. `reviewrecord.TicketRequirements` serves `RecordTicket` and the ticket form.
- `checkFloorCheck` serves the row list in `internal/preflight/decision.go`, and `gatherCheckFloor` serves the facts of `internal/preflight/plan.go`.
- `commitChainCommand` serves the `commit` definition, and `commit.parseArgs` serves `commit.Run`.
- `runLane` serves `RunLane`, which `bench commit` and `bench worktree merge` call. The lane line renders in `internal/gate/authorization/authorization.go`.
- `setRun.answer` serves each mutant-set path of `internal/probe`.
- `treeTargetInsertion` serves the root help. `childExecutable` and `worktreeBuild` serve `treetarget.Run`.

These files pin the surfaces that the spec changes:

- `internal/anchors/registry_chunk_chain.go:12`, pinned by `internal/anchors/registry_chunk_chain_test.go:14`, requires the exec preflight sentence that BC48 forbids. Ticket 07 replaces that entry.
- `internal/anchors/registry_retained_workflow.go:162` requires the exec gate sentence that BC68 replaces. Ticket 12 replaces that entry.
- `internal/preflight/charge_cells_test.go:22` and `internal/preflight/decision_test.go:57` and `:535` assert the `check-floor` row. Ticket 01 changes them and deletes `check_floor.go`, `check_floor_test.go`, and `check_floor_command_test.go`.
- The shared plan fixtures are `internal/reviewrecord/recordtest/fixture.go:194` and `internal/preflight/preflighttest/fixture.go:102`. Their plain requirement `additional` is read in `internal/reviewrecord/recordcmd/verification_test.go`, `internal/reviewrecord/recordcmd/ticket_test.go`, and `internal/preflight/ownership/premises_test.go`. Ticket 02 moves them.
- The retro body fixtures are `internal/retros/retros_test.go`, `internal/retros/recommendations_test.go`, and `internal/retros/testdata/eligible.md`. Ticket 04 adds the Cost heading to each valid body.

The build reruns `bench consumers` on each changed function before each charge, unexported functions included.

### Pre-review proof checklist

- Cited symbols: each symbol above resolves in the named file at base 887fd6c93. The landing-commit owner in `internal/spec` is new.
- Import edges: `go list -deps ./internal/reviewrecord` holds no `internal/roadmap`, so the new edge from `internal/roadmap` to `internal/reviewrecord` adds no cycle. `internal/spec` already imports `internal/git`, and `go list -deps ./internal/git` holds no `internal/spec`.
- Source-row clauses and occurrences: the source-to-row table below.
- Promised field labels: `## Cost`, `over budget`, `bench-lane.log`, `next=`, and `--base`.
- Changed-function callers: the inventory above.
- Copy survival: BC14 reds a second copy of a budget, BC35 and BC37 red the retired floor wording, BC48 reds the retired exec step, and BC73 reds the exec form in the known-issues row. The two live Require anchors on the retired forms change in tickets 07 and 12.
- Rendered-shape readers: the `[<check>]` prefix lines appear in `internal/gate/lane_run_test.go` and `internal/gate/run_stream_test.go`. Only the lane file reads lane lines, and BC51 keeps its red case. The `mutant_report` table appears in `internal/probe/set_test.go` (`TestMutantSetPrintsEachReport`), which ticket 08 adjusts for the bitten case.
- Pin operators: the output rows compare exact lines, and the bound rows compare a line count with `<=` 10.
- Entry reads: the cost renderer takes the home and the root as arguments, as `scaffoldBody` does now.
- Derived expectations: each fixture states its inputs, and no expectation reads the renderer's budget table.
- Consolidated rules: none.
- Quantified obligations: BC59 and BC60 enumerate the command registry, and BC33 grades each plain chunk requirement of a plan.
- Workflow-step writes: the ticket form writes the review record and removes its stage, as now. The record commit then runs the build preflight at the published commit with the chunk base.

### Source-to-row table

| source sentence | rows |
|---|---|
| The retro scaffold derives a cost section from git and the trace store. | BC1–BC10, BC19–BC21, BC72 |
| A budget overrun makes the drain propose one removal or one merge, and it never makes the gate red. | BC11–BC14, BC17, BC18, BC22, BC71 |
| An author records no requirement runs. | BC23–BC30, BC32–BC37 |
| The witness reds, the `--mutants` probe, the lane pass, and the chunk checkpoint grade stay. | BC30, BC31 |
| One record verb after the ticket commit replaces the separate record steps. | BC38–BC48 |
| A widen writes no learning. | BC69, BC70 |
| The success output of `bench commit` and of `bench probe --mutants` fits in 10 lines. | BC49–BC57 |
| Per-verb help names `--in` (reviewer scope, 2026-10-28). | BC58–BC60 |
| Guidance names `bench <verb> --in` (reviewer scope, 2026-10-28). | BC66–BC68, BC73 |
| A kit worktree target reuses a sealed binary that matches its tree (FT341 spec A reopened, 2026-10-28). | BC61–BC65 |

### Flagged additions

- `bench commit --base` is new grammar. The record commit needs it, because the chained preflight must grade inside the chunk with the chunk base.
- The lane log `bench-lane.log` is a new file. The short lane pass needs a place for the passed reports.
- BC13 adds the over-budget case for bookkeeping commits beside the wall case.

### Staged-spec claims

`specs/guidance-weight/` names the `check-floor` row in its spec and its ticket 2. Its build-entry staleness pass corrects that claim.
FT380 names a `check-floor` face, which this spec retires. The next drain reconciles that row.
