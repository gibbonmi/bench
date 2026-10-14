# Bench efficiency and defect assessment

Fix false-success and preservation defects first. Then remove repeated verification, manual evidence work, and conflicting workflow instructions.
Use the approved landing-test-efficiency spec for measured test-cost work.
Keep speculative capabilities behind a demonstrated need.

Status: Assessment and prioritization proposal. No implementation, roadmap edit, commitment change, benchmark, or test-suite run.
Source: `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68`, primary checkout, 2026-10-06.
Consumed by: Reviewer decisions about Bench maintenance and delivery order.
Drift: Recheck affected findings when their cited code, guidance, roadmap rows, or literal dependencies change.
Retire when: A reviewed replacement assessment supersedes this snapshot.

## Evidence and questions

Two Sol 6.1/high readers returned guidance and production-defect assessments.
A third CLI reader exceeded its no-output deadline and returned no usable evidence.
The coordinator completed the CLI assessment from the current source and roadmap.
The coordinator read all 99 open roadmap bodies and checked the sources behind the material findings.
The inventory contains 88 rows outside the parked section and 11 explicitly parked rows.
FT331 also carries a pause inside its body.

This report answers five questions:

1. Which code defects can report false success, lose meaning, or prevent recovery?
2. Which CLI interactions create avoidable calls, reads, or repeated work?
3. Which guide content helps, conflicts, or can move behind a specific task?
4. Which test and maintenance work can reduce cost while preserving failure detection?
5. What is the proposed next action for every open roadmap row?

Source-confirmed means that the inspected control flow establishes the mechanism.
It does not mean that this assessment reproduced the behavior.
Historical observations retain their original evidence limits.
Priorities and proposed cuts are judgments, not measured savings.

## 1. Code defects and recovery

### Failed work can return success

**New owner needed: worktree shell outcome propagation.**
A shell that cannot start calls the release closure with a failed work state.
The closure returns cleanup's exit code, so successful cleanup produces exit zero.
Normal completion discards the shell's wait error and records completed work.
The absent-shell test checks diagnostics and span fields but ignores the returned code.
[Release and completion](../../internal/worktree/subshell.go:73), [missing assertion](../../internal/worktree/subshell_test.go:46).

Preserve the work result through successful cleanup.
Define which result wins when both work and cleanup fail.
Add cases for failed start, nonzero child exit, and failed cleanup.
FT362 overlaps process ownership but does not currently own this exit contract.

**FT354 has a concrete parser defect, beyond duplication.**
The exception validator rejects a second decode only when that decode succeeds.
For `[] junk`, the second decode fails, so the function accepts the document when findings are empty.
The same mechanism accepts `[] 42`.
The caller reads the bytes directly and can mark the vulnerability phase green.
[Validator](../../internal/releasepreflight/vulnerability.go:93), [caller](../../internal/releasepreflight/vulnerability.go:69), [phase result](../../internal/releasepreflight/command.go:215).

Use the existing strict document reader and retain the public diagnostic category.
Its complete-consumption and duplicate-key checks already have an owner.
The current vulnerability test covers finding extraction, not this policy.
[Strict reader](../../internal/jsonfile/decode.go:40), [current test](../../internal/releasepreflight/vulnerability_test.go:8), [partial roadmap owner](../../roadmap/FT354.md:11).

These findings are source-confirmed, with no runtime reproduction.
They warrant focused negative cases before production changes.

### Commit and retirement contracts need preservation fixes

**FT258: pending merges can lose their second parent.**
The commit route passes the current tip to a path-attributed landing.
That owner always supplies one parent to `commit-tree`.
Its request validation checks required fields and commit identity, not pending merge state.
[Commit route](../../internal/commit/commit.go:191), [one-parent creation](../../internal/landing/landing.go:170), [validation](../../internal/landing/attribution.go:14).

The row already records this failure.
First refuse pending merge state with a valid completion route, or implement the decided parent-preserving contract.
The convenience work for automatic path selection need not block that repair.
[FT258](../../roadmap/FT258.md:15).

**FT284: retirement removes decisions before checking surviving consumers.**
The command verifies the implemented marker, then recursively removes the spec folder.
It does not first resolve surviving references or preserve decisions they still need.
The row records a surviving reference to a deleted spec.
This establishes a conditional preservation gap, not a claim that this assessment deleted or lost content.
[Removal path](../../internal/spec/spec.go:269), [recursive targets](../../internal/spec/spec.go:326), [recorded occurrence](../../roadmap/FT284.md:17).

### Cancellation needs one bounded owner

**FT362 should include cancellation defects in its scope.**
The vulnerability scanner sends TERM and waits without an escalation deadline.
Prep-release overrides cancellation with TERM and has no kill or drain bound.
A TERM-resistant child can therefore prevent cancellation from completing.
[Scanner](../../internal/releasepreflight/vulnerability.go:33), [prep-release](../../internal/preprelease/preprelease.go:230).

The external phase runner uses `CommandContext`, which can kill its direct child.
Its separate risk is an unbounded drain when descendants retain inherited pipes.
Do not describe that route as lacking direct-child cancellation.
[External runner](../../internal/releasepreflight/command.go:247).

Port these callers through an owner with explicit grace, group termination, drain, and reap behavior.
Keep each caller's evidence-preservation rules.
The broader process census and platform policy remain design work.
[FT362](../../roadmap/FT362.md:3).

### Cheap checks currently defer predictable failures

**FT215: guidance changes do not select every check that reads them.**
Markdown selects prose checks, and registered anchor files also select anchor currency.
The document families do not include guidance.
Yet AXI conformance and guidance budgets read those files.
[Selection](../../internal/gate/lane_select.go:143), [families](../../internal/gate/lane_select.go:171), [AXI reader](../../internal/conformance/axi_query_registry_test.go:54), [registry inputs](../../internal/conformance/registry/checks.go:33).

The final gate still runs those checks.
The defect causes delayed feedback and another repair cycle; it does not demonstrate a bypass of the final gate.
Declare complete check inputs before narrowing any landing contract.
[FT215](../../roadmap/FT215.md:18).

**FT377: unsupported package tools can mutate a supposed inspection subject.**
The package check invokes `npm pack --dry-run --ignore-scripts` without a toolchain-floor check.
The roadmap records repeated Node 18/npm 9 runs that rewrote the binary, seal, and broker manifest.
The current source retains the unguarded invocation.
This assessment did not rerun that destructive case.
[Invocation](../../internal/conformance/package_core_checks_test.go:72), [declared floor](../../package.json:59), [recorded occurrences](../../roadmap/FT377.md:15).

### Markdown parsing is a correctness issue

FT358 should be treated as shared grammar correctness, not merely reduced duplication.
The spec metadata reader recognizes only unindented backtick fences.
A `Status:` example inside a tilde fence can therefore become metadata.
The nearby live-reference reader uses different whitespace treatment.
[Metadata](../../internal/spec/spec.go:69), [reference reader](../../internal/spec/spec.go:45), [staged plan](../../roadmap/FT358.md:3).

The existing staged block-reader spec owns this migration.
Retain tests of each caller's semantics as well as the shared block grammar.

### A degraded guard has a narrower failure than its policy promises

FT366 includes a missing-core refusal defect.
On malformed quoting, the shell fallback splits only on whitespace and line boundaries.
A separator attached to a preceding word can therefore hide the next Git command from command-position recognition.
The fallback then permits a command that its missing-core policy intends to refuse.
[Fallback](../../.bench/hooks/block-dangerous-git.sh:294), [recognition](../../.bench/hooks/block-dangerous-git.sh:127), [missing-core decision](../../.bench/hooks/block-dangerous-git.sh:91).

Repair that refusal before the broader parser consolidation.
The hook is an honest-mistake control, with separate backstops; this finding does not establish an evasion-resistant security boundary.
[Declared boundary](../../.bench/hooks/block-dangerous-git.sh:10), [existing owner](../../roadmap/FT366.md:12).

## 2. CLI and workflow efficiency

### Remove no-progress routes

FT353 describes an action that cannot finish its advertised job.
Status selects unclaimed cleanup whenever unclaimed rows exist.
Cleanup deliberately retains unique refs and emits no apply action when all rows remain retained.
Preservation is correct; the proposed next action is incomplete.
[Status](../../internal/status/status.go:527), [retention](../../internal/worktree/clean_classes.go:47), [terminal plan](../../internal/worktree/clean_unclaimed.go:140).

Route unique refs to an owner and disposition decision.
Do not convert that route into implicit authority to discard or land their work.
[FT353](../../roadmap/FT353.md:10).

### Avoid model work for computable prerequisites

FT375 has a strong cost case.
Implementation always dispatches the staleness pass.
Each reader starts by computing drift and returns immediately when no drift exists.
Compute that fact once before dispatch.
Keep semantic review for actual drift and relevant preflight failures.
[Dispatch](../../.agents/commands/bench-implement-spec.md:21), [early return](../../.agents/skills/bench-implement-spec/references/staleness-pass.md:21).

One prior audit recorded 180,229 subagent tokens, 98 tool uses, and 577,371 ms.
The measurement labels its charge too broad.
The roadmap states that it found no drift.
These values describe one historical run; they do not predict the cost of every current pass.
[Measurement](/home/mgibs/workspace/bench/.logs/sr-line-comparison/measurements.md:5), [owner](../../roadmap/FT375.md:23).

FT293 and FT369 should make caller, fixture, and registry closure reliable before author dispatch.
Their occurrence ledgers show repeated fence amendments and preserve-and-restore work.
Fix the concrete preparation failures separately from FT369's proposed larger-ticket default.
[FT293](../../roadmap/FT293.md:27), [FT369](../../roadmap/FT369.md:18).

### Use the existing reader and record owners

Prioritize FT125's exact bounded file, section, and artifact reads.
The response owner applies the same size bound to ordinary output and worktree exec output.
A whole-spill read through that route can therefore produce another spill.
Exact pagination would break that read loop without discarding bytes or raising every response's limit.
[Response owner](../../internal/responsebound/owner.go:1), [spill transition](../../internal/responsebound/owner.go:124).

Prioritize FT318's native assignment and repair-record forms where they remove manual JSON edits.
FT304 should project existing execution state rather than create another state store.
FT265 should make sealed capture bodies reachable through one exact read.
[FT125](../../roadmap/FT125.md:15), [FT318](../../roadmap/FT318.md:22), [FT304](../../roadmap/FT304.md:3), [FT265](../../roadmap/FT265.md:7).

These are specific reductions in repeated work.
A general new Git wrapper, universal output rewrite, or command for every raw call has no established benefit here.
FT254 already states that raw-call counts alone do not select commands.
[FT254](../../roadmap/FT254.md:45).

## 3. Does BENCH.md help or hurt?

Its authority, scope, evidence, and isolation rules have clear purposes.
The current problem is that some instructions conflict or expose details before a session needs them.
File size alone does not establish worse model behavior.

The guide has 2,464 whitespace-delimited words, 15,803 characters, and 188 physical lines.
Its workflow section has 1,277 words.
The reference has 4,737 words, and the project profile has 6,214.
These are text counts, not model token counts.

The budget grades physical lines and permits exactly 188 for the guide.
A longer paragraph can therefore increase context cost without increasing budget use.
Choose a word or normalized measure before changing the ceiling.
[Budget](../../projects/benchkit.md:480), [counter](../../internal/conformance/prose_budget_test.go:206), [FT100](../../roadmap/FT100.md:39).

### Correct contradictions first

FT89 identifies conflicting verification instructions.
The delegation skill requires a gate, independent probe, and citation inspection for a returned contribution.
The implementation command tells the orchestrator to read manifests and verdicts, not code.
The reference describes one whole-project landing gate.
Define separate verification routes for ticket authors, independent writers, and read-only readers.
[Delegation](../../.agents/skills/bench-craft-delegate/SKILL.md:113), [orchestration](../../.agents/commands/bench-implement-spec.md:29), [landing](../../.bench/BENCH-reference.md:277).

FT379 has a recorded reviewer decision awaiting consistent guidance.
The scorecard says record-only changes consume no repair cycle.
The current repair policy excludes disposition changes from its evidence-only exception.
The row records six further rounds over record prose without source changes.
[Decision](../../capture/agent-performance/claude-models.md:55), [policy](../../.agents/skills/bench-craft-line/references/bounded-repair-policy.md:46), [occurrence](../../roadmap/FT379.md:14).

Align those rules while retaining checks against false claims and stale evidence.
No cycle exemption excuses an incorrect observation, a missing required check, an unresolved finding, or an invalid source identity.

### Keep, move, and cut candidates

| Action | Content | Reason and owner |
|---|---|---|
| Keep | Reviewer authority, acceptance shortfalls, source identity, isolation, scope, final gate | These define what a successful authorized change means. [Guide](../../.bench/BENCH.md:31). |
| Keep | Fresh-author and repair policy until a reviewer changes it | The decision is closed; this assessment supplies no contrary controlled trial. [Authorship](../../.bench/BENCH.md:121). |
| Move | Version-2 plan fields and author-limit details | Load them at implementation entry. Keep the policy and pointer in the guide. [Current location](../../.bench/BENCH.md:121). |
| Move | Ignored-versus-tracked capture routing details | Let capture help or its reference own the lookup. [Current location](../../.bench/BENCH.md:179). |
| Consolidate | Authorship, chunk, and repair restatements in the reference | Keep one policy owner and the reference's actual procedure. [Reference](../../.bench/BENCH-reference.md:275). |
| Replace | Repeated manual no-drift checks | Let preflight compute them once. [FT375](../../roadmap/FT375.md:10). |
| Narrow | Diagnostic ceremony for a known mechanical duplication fold | Keep a failure predicate and regression proof; do not invent several hypotheses. [FT89](../../roadmap/FT89.md:12). |
| Measure | Broad style cuts and changed model or author defaults | FT231 supplies the comparison; size and one favorable sample are insufficient. [FT231](../../roadmap/FT231.md:3). |

FT100 is literally blocked by FT231.
Correctness repairs under FT89 and FT379 do not depend on a broad prose experiment.
Inventory anchor owners before any editorial cut.
[Dependencies](../../ROADMAP.md:283), [anchor disposition](../../roadmap/FT100.md:18).

A guide comparison must test the proposed cuts and retain failure detection and task quality.
An unrelated narrow pilot neither closes FT231 nor waives FT100's dependency.
Any smaller prerequisite milestone requires an explicit amendment to that dependency.

## 4. Tests, cleanup, and scope control

The approved landing-test-efficiency spec already covers the strongest measured test-cost candidates.
Its 14 historical green runs have a median gate duration of 372.581 seconds.
Ordinary tests take 260.313 seconds; system tests take 81.147 seconds.
The adoption and repair fixtures are the largest package observations.
Package durations overlap and cannot be added as expected savings.
[Retained baseline](../../specs/landing-test-efficiency/spec.md:16).

Implement that spec as a distinct efficiency outcome after the narrow correctness repairs.
It retains failure assertions, isolated fixture state, full installation journeys, and independent durable authorization.
It reduces fixture payload, repeated conformance runs, repeated source parsing, and duplicate tree materialization.
No speedup has yet been demonstrated.
[Scope](../../specs/landing-test-efficiency/spec.md:37), [quality constraints](../../specs/landing-test-efficiency/spec.md:101).

The shell defect shows why more tests are not automatically better tests.
A test already exercises the failure but omits the exit-status assertion.
Favor tests that distinguish the named wrong result.
Keep independent expectations when a demonstrated omission needs them.
Do not derive every expected answer from the same implementation being tested.
[Missed assertion](../../internal/worktree/subshell_test.go:46), [expectation rule](../../AGENTS.md:42).

FT365 overlaps the four shared source visitors in the approved spec.
Its independent-expectation audit and remaining scans stay open.
FT360 shares fixture helpers but does not by itself establish lower fixture runtime.
FT255 should control host contention before any new phase overlap.
[FT365](../../roadmap/FT365.md:3), [FT360](../../roadmap/FT360.md:3), [FT255](../../roadmap/FT255.md:18).

FT283 has a specific approved direction for spec-stage lane-only landings.
Its row explicitly says the whole-project contract remains until the change ships.
Finish the safe input and phase contract before using that exception.
Treat FT215's broader Markdown-only proposal and FT314's cross-run evidence reuse as separate decisions.
[FT283](../../roadmap/FT283.md:11), [FT215](../../roadmap/FT215.md:34), [FT314](../../roadmap/FT314.md:3).

### Cleanup claims that need correction

FT364 overstates the absence of callers and tests.
Public npm staged submission is explicitly refused before a lock or registry call.
The fixture adapter still reaches the staged state machine, and adapter tests exist.
Fixture state-machine coverage does not establish public npm staged support.
The unimplemented public adapter is a supported-scope decision, not proof that the entire state machine is dead.
[Refusal and caller](../../internal/publication/command.go:206), [adapter boundary](../../internal/publication/npm_registry.go:158), [fixture tests](../../internal/publication/command_adapter_test.go:316).

FT368's lack of in-tree workflow callers does not establish absence of external users.
Require a consumer decision before removing a public command.
Prefer direct consolidation candidates with known callers: strict JSON, Markdown blocks, fixture runners, and cleanup application.
[FT368](../../roadmap/FT368.md:3).

FT106's premise that no documentation claim is ever checked is too broad.
The tree already grades selected registry and profile facts.
Keep its semantic claim audit and missing-link work, using existing checks as inputs.
[Current checks](../../internal/conformance/tier_test.go:40), [FT106](../../roadmap/FT106.md:3).

FT384 says its production scope matcher is unlocated.
It is now located: candidate authorization calls `tickets.Covers`, while delivery closure uses exact deliverable paths.
Those functions answer different questions.
A difference alone does not prove incorrect authorization or closure.
Require a concrete failing scenario before prescribing one matcher.
[Authorization](../../internal/commitment/repository/candidate.go:232), [closure](../../internal/commitment/delivery.go:67), [FT384](../../roadmap/FT384.md:3).

## 5. Fresh roadmap order

Rank the next bounded action, not every feature accumulated inside a row.
Use consequence, evidence strength, recurrence, cost to unblock later work, and scope.
Existing HIGH, MEDIUM, LOW, and queue positions did not determine this order.
A decision or reproduction can be the next action; the rank does not authorize implementation.

The new shell-exit defect belongs beside the first three repairs below.
FT354's rank applies to its exception-policy defect; its broad quoting and durability work follows later.
The approved landing-test-efficiency spec belongs after the first reliability group, alongside the repeated-work reductions.
It has no separate open roadmap row in this snapshot, and only part of it overlaps FT365.

The exact order inside a group is judgment, not a calculated return-on-investment score.
Literal dependencies remain binding.
Broader experiments do not block small correctness repairs.

### First: reliable results and recoverable operations

| Rank | Owner | Next action and reason |
|---|---|---|
| 1 | [FT258](../../roadmap/FT258.md:1) | Refuse or preserve pending merge parents first; separate automatic path discovery. |
| 2 | [FT377](../../roadmap/FT377.md:1) | Guard the package toolchain floor before inspection can rewrite generated artifacts. |
| 3 | [FT354](../../roadmap/FT354.md:1) | Fix exception-policy trailing input first; then consolidate strict readers without waiting on quoting. |
| 4 | [FT366](../../roadmap/FT366.md:1) | Close the missing-core parser refusal edge; then consolidate wrappers and envelopes. |
| 5 | [FT215](../../roadmap/FT215.md:1) | Close check-input selection gaps before changing lane or landing cost. |
| 6 | [FT349](../../roadmap/FT349.md:1) | Require ticket check evidence and remove completed rows' planned-test exemptions. |
| 7 | [FT342](../../roadmap/FT342.md:1) | Decide one moved-main route; prevent repeated folds, fence repairs, and re-review. |
| 8 | [FT89](../../roadmap/FT89.md:1) | Resolve contradictory verification and mechanical-fix instructions. |
| 9 | [FT379](../../roadmap/FT379.md:1) | Apply the existing record-only correction decision and narrow confirming rounds. |
| 10 | [FT358](../../roadmap/FT358.md:1) | Use the staged block-reader plan to fix inconsistent Markdown interpretation. |
| 11 | [FT341](../../roadmap/FT341.md:1) | Repair terminal-child behavior and label addressing before completing routing migration. |
| 12 | [FT362](../../roadmap/FT362.md:1) | Bound scanner and prep-release cancellation first; then share process ownership. |
| 13 | [FT284](../../roadmap/FT284.md:1) | Preserve live decision references before recursive spec retirement. |

### Next: remove recurrent wasted work

| Rank | Owner | Next action and reason |
|---|---|---|
| 14 | [FT293](../../roadmap/FT293.md:1) | Decide complete ownership closure and useful diagnostics before more author dispatches. |
| 15 | [FT375](../../roadmap/FT375.md:1) | Compute drift and cited-symbol/test state before launching staleness readers. |
| 16 | [FT283](../../roadmap/FT283.md:1) | Implement the decided spec-stage lane contract after its input coverage is sound. |
| 17 | [FT125](../../roadmap/FT125.md:1) | Deliver exact bounded artifact and section reads; validate fewer calls on real tasks. |
| 18 | [FT290](../../roadmap/FT290.md:1) | Complete the staged fixture/count/selection projection; distinguish executed checks from empty success. |
| 19 | [FT255](../../roadmap/FT255.md:1) | Set one test concurrency budget; avoid load-driven retries before phase overlap. |
| 20 | [FT141](../../roadmap/FT141.md:1) | Attribute reds to the exact composed baseline; avoid diagnosing unrelated failures as regressions. |
| 21 | [FT353](../../roadmap/FT353.md:1) | Replace unique-ref cleanup dead ends with an owned disposition and recovery route. |
| 22 | [FT265](../../roadmap/FT265.md:1) | Make sealed capture bodies and primary-aware reads exact and consistently reachable. |
| 23 | [FT317](../../roadmap/FT317.md:1) | Decide capability-blocked evidence values before completion can record them. |
| 24 | [FT318](../../roadmap/FT318.md:1) | Deliver decided completion values and native record forms; depends on FT317 for blocked rows. |
| 25 | [FT344](../../roadmap/FT344.md:1) | Use the supported worktree-build remedy at every freshness refusal. |
| 26 | [FT380](../../roadmap/FT380.md:1) | Check plan changes before dispatch; isolate any later publication-enforcement decision. |
| 27 | [FT369](../../roadmap/FT369.md:1) | Fix caller, fence, and refusal-member preparation; separate the larger-ticket experiment. |
| 28 | [FT98](../../roadmap/FT98.md:1) | Provide verified evidence preservation and precise cleanup plans. |
| 29 | [FT253](../../roadmap/FT253.md:1) | Own the full landing lease and its wait route to reduce conflicting landing attempts. |

### Then: consolidate existing machinery

| Rank | Owner | Next action and reason |
|---|---|---|
| 30 | [FT106](../../roadmap/FT106.md:1) | Refresh stale claims and verify concrete links; reuse current mechanical checks. |
| 31 | [FT172](../../roadmap/FT172.md:1) | Decide roadmap identity and grammar before building another execution view. |
| 32 | [FT304](../../roadmap/FT304.md:1) | Project existing run and roadmap state after FT172; keep observation separate from control. |
| 33 | [FT307](../../roadmap/FT307.md:1) | Show changed-path and package headroom; reproduce remaining untracked-file concerns. |
| 34 | [FT168](../../roadmap/FT168.md:1) | Add the probe forms repeatedly needed for existing verification; preserve restoration. |
| 35 | [FT333](../../roadmap/FT333.md:1) | Support current tickets-only evidence without requiring a synthetic spec. |
| 36 | [FT312](../../roadmap/FT312.md:1) | Automate review invocation and venue lifecycle as independent slices. |
| 37 | [FT387](../../roadmap/FT387.md:1) | Finish verified output preservation on supported harness paths; refuse unsupported replacement. |
| 38 | [FT388](../../roadmap/FT388.md:1) | Retire the coordinator immediately after FT387; this is closure, not a build. |
| 39 | [FT231](../../roadmap/FT231.md:1) | Run one narrow comparison of an expensive workflow; defer the experiment portfolio. |
| 40 | [FT100](../../roadmap/FT100.md:1) | Measure and cut unnecessary guidance after FT231; resolve FT89 first. |
| 41 | [FT326](../../roadmap/FT326.md:1) | Make anchor mutations robust to harmless wrapping without weakening omission detection. |
| 42 | [FT343](../../roadmap/FT343.md:1) | Decide test-hook policy and remove unused hooks with caller evidence. |
| 43 | [FT373](../../roadmap/FT373.md:1) | Use the bounded standard-library duplicate check after concrete correctness and rework fixes. |
| 44 | [FT365](../../roadmap/FT365.md:1) | Complete remaining scanner and expectation work; avoid duplicating landing-test-efficiency. |
| 45 | [FT360](../../roadmap/FT360.md:1) | Apply the staged shared-fixture plan where it removes divergent harness knowledge. |
| 46 | [FT363](../../roadmap/FT363.md:1) | Share cleanup-set application while preserving mode-specific decisions and fault behavior. |
| 47 | [FT302](../../roadmap/FT302.md:1) | Promote only current repeated Git facts; deepen diff policy in a separate slice. |
| 48 | [FT217](../../roadmap/FT217.md:1) | Share adoption decisions with unchanged public behavior and fixture assertions. |
| 49 | [FT254](../../roadmap/FT254.md:1) | Prioritize stdin, timeout, and recovery ergonomics; defer convenience faces without call evidence. |
| 50 | [FT244](../../roadmap/FT244.md:1) | Provide lifecycle-owned scratch storage to avoid untracked source contamination. |
| 51 | [FT378](../../roadmap/FT378.md:1) | Use one concise known-issue owner and retirement trigger; avoid another policy manual. |
| 52 | [FT381](../../roadmap/FT381.md:1) | Decide only demonstrated delegation gaps; preserve the existing bounded repair policy. |
| 53 | [FT384](../../roadmap/FT384.md:1) | Reproduce the alleged scope mismatch using the now-located matchers before designing a replacement. |
| 54 | [FT207](../../roadmap/FT207.md:1) | Reproduce malformed-admin hangs, then route mutations through the existing refusal owner. |
| 55 | [FT350](../../roadmap/FT350.md:1) | Prove the ref race through a sanctioned route before changing transaction guarantees. |
| 56 | [FT140](../../roadmap/FT140.md:1) | Revalidate and split the mixed residual bundle; drop superseded or unavailable claims. |
| 57 | [FT348](../../roadmap/FT348.md:1) | Decide disposition of the remaining recovered streams and close settled history. |
| 58 | [FT308](../../roadmap/FT308.md:1) | Design a race-safe stale-lock sweep; never unlink a lock merely because it looks empty. |
| 59 | [FT208](../../roadmap/FT208.md:1) | Improve current refusal diagnostics; keep its expensive fixture-boundary decision separate. |
| 60 | [FT260](../../roadmap/FT260.md:1) | Provide historical replay or transfer only for recurring unsupported workflows. |

### Later: bounded measurement and optional workflow expansion

| Rank | Owner | Next action and reason |
|---|---|---|
| 61 | [FT314](../../roadmap/FT314.md:1) | Consider cross-run green reuse only with the complete authorization key and measured repetition. |
| 62 | [FT332](../../roadmap/FT332.md:1) | Add constrained stress runs after budgeting shared test resources and assigning a result owner. |
| 63 | [FT232](../../roadmap/FT232.md:1) | Run the existing repair pilot and finish its report before adding a detector. |
| 64 | [FT386](../../roadmap/FT386.md:1) | Trial a spec-size advisory; avoid another hard threshold without quality evidence. |
| 65 | [FT352](../../roadmap/FT352.md:1) | Obtain reliable context measurements before changing author or spec-stage session policy. |
| 66 | [FT204](../../roadmap/FT204.md:1) | Add a bounded session/census query for demonstrated consumers; reuse current records. |
| 67 | [FT222](../../roadmap/FT222.md:1) | Keep routing stable until comparable measurements support a default change. |
| 68 | [FT243](../../roadmap/FT243.md:1) | Set recurring maintenance only after showing this assessment's recommendations repay their cost. |
| 69 | [FT130](../../roadmap/FT130.md:1) | Verify capture interference for the actual gate subject before adding queue coordination. |
| 70 | [FT371](../../roadmap/FT371.md:1) | Decide a safe scratch-Git probe route when a current investigation needs it. |
| 71 | [FT368](../../roadmap/FT368.md:1) | Decide public consumers before speculative cuts; retain independent test bite. |
| 72 | [FT364](../../roadmap/FT364.md:1) | Correct the caller/test premise, then decide unsupported staged-release scope after FT142. |
| 73 | [FT383](../../roadmap/FT383.md:1) | Take the small fixture deduplication during nearby maintenance; no standalone campaign. |
| 74 | [FT235](../../roadmap/FT235.md:1) | Keep directory renaming low; current labels already support human navigation. |
| 75 | [FT200](../../roadmap/FT200.md:1) | Revisit broader landing admission after FT342 and phase-scoped contracts settle. |

### Last: adoption and new capabilities

| Rank | Owner | Next action and reason |
|---|---|---|
| 76 | [FT142](../../roadmap/FT142.md:1) | Revalidate release residuals before release machinery changes or external qualification. |
| 77 | [FT374](../../roadmap/FT374.md:1) | Decide needed non-Go equivalents for an actual adoption target; validate toolchain requirements. |
| 78 | [FT305](../../roadmap/FT305.md:1) | Build durable execution after cancellation, landing, phase, and evidence contracts settle. |
| 79 | [FT306](../../roadmap/FT306.md:1) | Qualify a pinned release and external adoption after FT305 and the release requirements. |
| 80 | [FT275](../../roadmap/FT275.md:1) | Add seam tracing where a real diagnostic consumer earns its runtime and policy cost. |
| 81 | [FT287](../../roadmap/FT287.md:1) | Expand AXI only where a concrete command benefits; scoped consistency already has an owner. |
| 82 | [FT241](../../roadmap/FT241.md:1) | Add versioned behavioral promises after the evidence record and its consumers stabilize. |
| 83 | [FT101](../../roadmap/FT101.md:1) | Keep the declared multi-context adoption trigger; do not use it as a local speed shortcut. |
| 84 | [FT351](../../roadmap/FT351.md:1) | Build the approved visual report after higher-value inspection and correctness work. |
| 85 | [FT334](../../roadmap/FT334.md:1) | Add deadline batching after reliable admission, cost records, and durable execution. |
| 86 | [FT240](../../roadmap/FT240.md:1) | Preserve the FT231 three-arm dependency and require a relevant monorepo trial. |
| 87 | [FT346](../../roadmap/FT346.md:1) | Commission a topology trial only if its distinct question remains worth the cost. |
| 88 | [FT331](../../roadmap/FT331.md:1) | Keep paused with FT347; no advisor adoption or paid resumption is authorized. |

### Parked rows: retain their evidence triggers

| Owner | Proposed disposition |
|---|---|
| [FT347](../../roadmap/FT347.md:1) | Keep parked until the November 3 reminder and an explicit response; no automatic research restart. |
| [FT6](../../roadmap/FT6.md:1) | Keep each candidate behind its stated observed-demand or failure trigger. |
| [FT24](../../roadmap/FT24.md:1) | Keep pending a deny-capable upstream delegation hook; no live upstream check was made here. |
| [FT261](../../roadmap/FT261.md:1) | Require the exact untracked-spec preflight reproduction; current source facts already include untracked paths. |
| [FT327](../../roadmap/FT327.md:1) | Require a current landing refusal on doctor-published output. |
| [FT328](../../roadmap/FT328.md:1) | Require the evidence-store read refusal under the declared read-only sandbox. |
| [FT330](../../roadmap/FT330.md:1) | Require the checkpoint to emit the wrong recovery command for the named failure. |
| [FT335](../../roadmap/FT335.md:1) | Keep only the unresolved full-gate-selection face parked; the caller-root face already shipped. |
| [FT355](../../roadmap/FT355.md:1) | Require a leftover tied to the crash test or its gate process before assigning cause. |
| [FT372](../../roadmap/FT372.md:1) | Require the local-capture test failure again before treating shared state as its cause. |
| [FT385](../../roadmap/FT385.md:1) | Require a recorded ambiguous-prefix refusal; do not infer that cause from any test retry. |

### Ordering constraints and unresolved scope

The five literal edges are FT100 after FT231, FT240 after FT231, FT304 after FT172, FT306 after FT305, and FT388 after FT387.
A dependency's required deliverable must finish before its dependent action starts.
The table ranks proposals and leaves the current commitment unchanged.
[Dependency source](../../ROADMAP.md:283).

The active FT290 commitment is existing work, not evidence that FT290 must rank first.
Changing that commitment is a separate reviewer decision.
The report does not stop, replace, or redirect its active assignment.
[Committed outcomes](../../.bench/commitment.json:1).

Several rows combine decisions, repairs, and optional features.
Examples include FT258, FT354, FT369, FT140, and FT231.
Give their first bounded action its own acceptance and completion condition.
Do not make a narrow confirmed defect wait for every proposed feature in its parent row.

## Unknowns and validation plan

No new elapsed-time comparison, full gate, regression suite, process-hang reproduction, or publication operation ran.
No controlled evidence establishes the effect of a shorter guide or a different author/reviewer topology.
External CLI consumers, unsupported-host behavior, and paid model cost remain unknown.
The suspected ref races and parked flakes retain their existing reproduction conditions.

Local logs are mutable and can expire.
The report retains the numeric observations it uses.
Recheck citations against the pinned source before implementation if the tree has advanced.

Before each proposed fix:

1. Reproduce the named wrong result at the smallest public or composition boundary.
2. Confirm that the new assertion fails for that result.
3. Preserve the existing correct result, failure diagnostics, and cleanup behavior.
4. Run the affected checks and the project's required publication checks.
5. Compare actual accepted-change cost, including review, repair, verification, and supervision.

For the test-efficiency spec, retain each named omission proof and measure during required verification.
For guide cuts, first fix contradictions and define one narrow FT231 comparison.
For deletion proposals, establish consumers and replacement guarantees before removing the surface.

## Verification record

- [x] The recommendation states its scope and evidence limits.
- [x] Findings separate source facts, historical observations, unknowns, and proposals.
- [x] Each material finding has local source citations.
- [x] The coordinator reopened the sources behind the central defect and workflow claims.
- [x] All 99 open rows have one primary disposition: 88 ranked and 11 parked.
- [x] The proposed order retains the five literal dependency edges.
- [x] No implementation, roadmap priority, delivery commitment, or user-owned file changed.
- [x] Independent synthesis comments identified dependency, guard-priority, and evidence-wording refinements; the coordinator applied them.
- [x] Document checks passed for all row identities, literal dependencies, citation targets, and prose.

The second synthesis pass supplied comments but no complete final verdict before the coordinator ended it.
No formal acceptance or implementation-quality verdict is claimed.
