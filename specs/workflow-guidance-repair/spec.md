# Workflow guidance repair

Status: staged

Decision source: the reviewer-confirmed current conversation, 2026-09-24: the 62-finding guidance audit at `78602e89`, three tree-closed decisions, and the reviewer's scope answer.

Verification log: 2 iteration(s) to accept — iteration 1 found 4 blocking and 8 fold findings. The blocking findings were anchor-output tests outside the fence, the bench-writer agent body, and a README copy. Three Forbid rows also removed the last Require guard of a duty. Iteration 2 confirmed every fold, and the acceptance folded its 3 prose and accounting findings.

## Problem

A read-only audit of the workflow guidance at `78602e89` found 62 defects in 20 guidance files. Some sentences still state the retired single-author model. Other sentences contradict the fresh-author rule, the fast lane decision, or the landing boundary. Some steps name a verb, a path, or a merge that the tooling refuses.

A session that obeys one can read a lane pass as green, send a repair to the wrong session, or stop at a refusal. Other sentences restate a fact that another file owns, so the two copies can drift. The `bench commit` help strings and one ADR consequence state the same stale facts.

## Solution

Each defect sentence changes to agree with its owner. A retired or wrong sentence goes, and the anchor registry refuses its return. A second copy becomes a pointer to its owner, or it goes. A step that the tooling refuses names the route that the tooling accepts.

The `bench commit` help names the fast lane. ADR 0014 states the enforcement that the tree applies. Each fix keeps every budgeted file inside its line budget, and each moved needle moves its canary with it.

## User stories

Line: opus / high.
Implementation-line reason: every chunk rewrites guidance that steers later builds, so the leverage override routes it mid and high. The rows name exact needles, and the anchor registry reds each dropped or returned sentence.
Harder chunks: none.

Delegation skill and discipline:

1. As a ticket author, I want every write delegate to run as `bench-writer`, so that a repair session also has an agent type.
2. As an orchestrator, I want `craft-delegate` to leave author and session changes to `craft-line`, so that a planned fresh author needs no extra user direction.
3. As a reviewer, I want a model change under a transfer trigger to ask me first outside `--delegate`, so that no tier moves silently.
4. As a ticket author, I want every write charge from a spec to carry its coverage rows, so that each author owes the red-then-green duty.
5. As a ticket author, I want every write delegate to treat `Writes:` as an expectation, so that the plan-expansion rule applies to each author.
6. As a ticket author, I want to commit on a lane pass in the shared integration worktree, so that the skill matches the build phase.
7. As an orchestrator, I want the returned-tree coordinator probe limited to a user-directed write delegate, so that I read verdicts and not code.
8. As a repair author, I want my repair fence to be my ticket's `Writes:` line, so that a wider fence needs plan expansion first.
9. As a delegate, I want the `git stash` ban to state the real guard surface, so that I do not rely on a false refusal.
10. As a coordinator, I want every probe to run through `bench probe`, so that no rule sends me to a copy aside or to `cmp`.
11. As a coordinator, I want no build to run in the primary checkout, so that the delegation discipline agrees with the landing boundary.
12. As a coordinator, I want the recovery and digest rules to name public verbs and a merge, so that each step is reachable.
13. As a reader of the delegation guidance, I want each restated rule to live only with its owner, so that no copy drifts.
14. As a reader of the bounded repair policy, I want its scope to name every implementation run, so that it names no retired retained mode.

Phase commands:

15. As a reviewer, I want a coverage-map amendment to update each affected ticket, so that one repair ticket never carries the repairs of several tickets.
16. As a reviewer, I want final-check to retain the orchestrator's integration-verification results, so that it names the performer that the record check requires.
17. As a reviewer, I want a spec-backed red at final check to go to a fresh repair session, so that the orchestrator does not repair.
18. As a ticket author, I want an out-of-fence cause to route through plan expansion or the reviewer's split, so that no coordinator reslices alone.
19. As a ticket author, I want the stops-short tier route to point to the `craft-line` ladder, so that the `--delegate` exception survives.
20. As a reader of the phase commands, I want `.bench/BENCH.md` alone to own the chunk-review timing, so that no command restates it.

Lane and landing:

21. As an operator, I want the guidance to separate the commit lane from the landing gate, so that a lane pass is never green.
22. As an operator, I want `bench help` and `bench commit --help` to name the fast lane, so that the executable help agrees with ADR 0017.
23. As a reader of the reference and ADR 0014, I want the true primary-checkout enforcement, so that no reader expects a commit there.
24. As an operator, I want the reference to drop the landing rebuild claim, so that the landing description matches the stable owner.
25. As an agent, I want the conflict repair to assign the raw merge to the reviewer, so that no step needs a denied merge.
26. As an operator, I want a spec retirement to run in a Bench worktree and land, so that it avoids the primary-checkout refusal.

Retro capture:

27. As a phase closer, I want a tracked retro and its scorecards to commit with the phase close, so that final-check matches the writer.
28. As a phase closer, I want an ignored retro to stay local until the drain, so that such a repository keeps the drain route.
29. As a phase closer, I want to write the retro once from its scaffold through `bench retro`, so that the guidance matches the verb.

One declared line:

30. As a spec author, I want every ticket author to run on the spec's declared line, so that no ceiling moves a tier silently.
31. As an orchestrator, I want each ladder step to obey the `--delegate` tier rule, so that escalation and the top-tier pause read one way.

Owners and references:

32. As a reader of the guide and the agreement, I want each pointer to name the file that holds its fact, so that none misleads.
33. As a session, I want the reference to name `bench gate-prose` as the internal verb that a session runs, so that the prose route is sanctioned.
34. As a reader, I want the reference to own the gate output, census, and handoff facts, so that no other file holds a copy.
35. As a reader of the reference, I want verb grammar left to executable help, so that the reference restates no grammar.
36. As a reader of the reference, I want the plan version stated as the checkpoint grades it, so that the amendment rule has one reading.
37. As a Claude Code user, I want the skill-link statement to match the link rule, so that the `prototype` link is expected.
38. As a drain coordinator, I want the drain and assess commands to name reachable routes, so that each step works as written.
39. As a skill author, I want each craft skill to point to the owner of a borrowed fact, so that no skill contradicts its owner.

Budgets and exclusions:

40. As a reader, I want each budgeted guidance file to stay inside its budget, so that the repair grows no file past its bound.
41. As a reviewer, I want `internal/reviewrecord` unchanged, so that no record-schema change rides with this guidance repair.

Enforcement harnesses:

42. As a maintainer, I want the `bench anchors` tests to read each row's kind from the registry, so that `AGENTS.md` can carry a Forbid row.

## Implementation decisions

- `.bench/BENCH.md` owns the workflow rules. Every changed sentence agrees with it, and a restated rule becomes a pointer to it.
- Every write delegate runs as `bench-writer`. That covers a fresh ticket author, a repair session, and a user-directed write delegation.
- A ticket author commits its own ticket on a lane pass. The whole-project gate runs once, at `bench worktree land`.
- A spec-backed build has one declared implementation line. A tier move of a fresh ticket author outside `--delegate` asks the reviewer first.
- A tracked retro and its scorecard updates commit with the phase close. An ignored retro stays local until the next reviewer-approved capture drain.
- The retro scaffold owns the retro headings and the calibration table header. `bench retro` writes the retro once and refuses an existing file.
- `.bench/BENCH-reference.md` owns the kit facts that a linked repository needs. This covers the gate output, the census signal, and the handoff verb behavior.
- `projects/benchkit.md` and `AGENTS.md` point to the reference for those facts. `AGENTS.md` keeps its project rules, which include the phase-close commit rule.
- `bench gate-prose` stays internal inventory. The reference names it as the one internal verb that a session runs directly.
- The `bench commit` help names the declared lane, and the gate only when a project declares no lane.
- ADR 0014 changes in place. Its enforcement consequence states that the commit verb refuses the primary checkout. It also states that a file-write guard refuses an agent's write to a tracked path there.
- A retired or wrong sentence takes a Forbid row. An anchored sentence that changes takes a replacement Require row, and its canary names the new bytes.
- A transfer-trigger sentence keeps its bytes. One new bullet applies the tier-move rule to every trigger, so the four trigger canaries stay valid.
- The anchor needle "Before the coordinator reads a probe verdict, the coordinator confirms the" keeps its bytes, so the canary `delegate-probe-mutated-bytes` stays valid.
- A new Require needle sits on one physical line of its guidance file. A Forbid needle can quote old text across a line break, because the matcher collapses whitespace.
- A canary whose needle changes is retargeted, never removed. It names the new Require bytes, or it plants the forbidden bytes back.
- A surviving duty keeps a Require row. A Forbid row that removes a copy or a retired clause never removes the last Require guard of a duty.
- A ticket author and a repair author commit on a lane pass and never run `bench worktree land`. A user-directed write delegate with no ticket returns an uncommitted diff.
- A budgeted file whose line count equals its budget takes line-neutral edits only.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| GR-A / `1-align-the-delegation-skill-with-fresh-authors.md`, `2-align-the-delegation-discipline-with-fresh-authors.md`, `3-route-phase-command-repairs-to-fresh-sessions.md` | The delegation skill, the delegation discipline, and the phase commands route every author, repair, and probe under the fresh-author rule. | GR1, GR2, GR3, GR4, GR5, GR6, GR7, GR8, GR9, GR10, GR11, GR12, GR13, GR14, GR15, GR16, GR17, GR18, GR19, GR20, GR21, GR22, GR23, GR24, GR25, GR26, GR27, GR28, GR29, GR30, GR31, GR32, GR33, GR34, GR35, GR36, GR37, GR38, GR39, GR40, GR41, GR42, GR43, GR44, GR109, GR110, GR111, GR112, GR113, GR115, GR116 | `bench test --check docs-currency-workflow`, `bench test --package ./internal/conformance --run TestRootConformance`, `bench test --package ./internal/anchors`, `bench test --check claude-agent-definitions`, `bench test --check guidance-prose-budgets`, `bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'` | no |
| GR-B / `4-state-the-lane-and-landing-gate-once.md`, `5-capture-the-retro-by-the-tracked-rule.md` | The guidance and the executable help say that the commit runs the lane and the landing runs the gate, and the retro follows its writer. | GR45, GR46, GR47, GR48, GR49, GR50, GR51, GR52, GR53, GR54, GR55, GR56, GR57, GR58, GR59, GR60, GR61, GR62, GR63, GR64, GR65, GR66, GR67, GR117, GR118, GR119, GR120 | `bench test --check docs-currency-workflow`, `bench test --package ./internal/conformance --run TestRootConformance`, `bench test --package ./internal/anchors`, `bench test --package ./cmd/bench --run TestHelpInventoryIsComplete`, `bench test --package ./internal/commit --run TestHelpAdvertisesDryRun` | no |
| GR-C / `6-bind-each-ticket-author-to-the-declared-line.md`, `7-point-the-guides-at-each-fact-owner.md`, `8-correct-the-drain-and-craft-skill-references.md` | Each fact has one owner, every author runs on the declared line, and every remaining reference resolves. | GR68, GR69, GR70, GR71, GR72, GR73, GR74, GR75, GR76, GR77, GR78, GR79, GR80, GR81, GR82, GR83, GR84, GR85, GR86, GR87, GR88, GR89, GR90, GR91, GR92, GR93, GR94, GR95, GR96, GR97, GR98, GR99, GR100, GR101, GR102, GR103, GR104, GR105, GR106, GR107, GR108, GR114, GR121, GR122, GR123, GR124 | `bench test --check docs-currency-workflow`, `bench test --package ./internal/conformance --run TestRootConformance`, `bench test --package ./internal/anchors`, `bench test --package ./cmd/bench --run TestAnchors`, `bench test --check guidance-prose-budgets`, `bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'` | no |

## Testing decisions

- A good test reads the live guidance through the anchor registry. A dropped sentence reds a Require row, and a returned sentence reds a Forbid row.
- The registries under `internal/anchors` are the precedent. The fresh-ticket-authors build that landed at `c58e5aa5` is the nearest precedent.
- The conformance tables in `retained_workflow_test.go` and `implementation_continuation_test.go` pair a diagnostic with a needle. A moved needle updates its table row.
- `checkRecurrenceMaintenanceContract` requires exact drain sentences. The drain rows change its strings and its bite table together.
- `TestAnchorsReportsNeedleLines` and `TestAnchorsReportsAbsentNeedles` build the expected `bench anchors AGENTS.md` output. Ticket 7 makes them take each row's kind from the registry.
- `fixture_bite_test.go` pins exact bytes of `craft-delegate` and `craft-spec` and one `craft-spec` diagnostic. The tickets keep those bytes and that diagnostic.
- `TestWorkflowCadenceAnchorsRejectDeletionAndSwap` and `TestSpecTicketHandoffWorkflowFixturesAreComplete` run those pins. The `docs-currency-workflow` check does not run them, so tickets 1, 6, and 8 and chunks GR-A and GR-C name them.
- `TestHelpInventoryIsComplete` compares the live help with the expected help. `TestHelpAdvertisesDryRun` grows one assertion on the lane wording.
- The `docs-currency-workflow` check and the root conformance test observe the guidance rows. Each changed canary keeps its mutation red.

### Seam diagram

    trigger: `bench test --check docs-currency-workflow`, the gate's test phase
        │
        ▼
    guidance files  ──▶  [ anchor registry: Require and Forbid needles ]  ──▶  diagnostic or pass
                      ◀ tests attach here: a needle row, and a canary that drops or restores its sentence

    trigger: `bench help`, `bench commit --help`
        │
        ▼
    command registry row, commit help text  ──▶  [ help renderers ]  ──▶  help stdout
                      ◀ tests attach here: the complete-inventory comparison and the dry-run help assertion

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GR1 | 1 | `.agents/skills/bench-craft-delegate/SKILL.md` names `bench-writer` as the type for a fresh ticket author, a repair session, and a user-directed write delegation | planned Require needle in `internal/anchors/registry_data.go` | A skill that types only a user-directed write lacks the needle, so the check reds. |
| GR2 | 1 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "delegation runs as `bench-writer`" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the retired sentence matches the needle, so the check reds. |
| GR3 | 1 | `.claude/agents/bench-writer.md` describes the type for a fresh ticket author and a repair session | planned Require needle in `internal/anchors/registry_data.go` | A description that names only a user-directed write lacks the needle, so the check reds. |
| GR4 | 2 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "A change of implementation author or session requires user direction." | planned Forbid needle that replaces the Require row in `internal/anchors/registry_data.go` | A skill that keeps the retired sentence matches the needle, so the check reds. |
| GR5 | 2 | `.agents/skills/bench-craft-delegate/SKILL.md` names `craft-line` as the owner of a change of implementation model or session | planned Require needle in `internal/anchors/registry_data.go` | A skill that drops the pointer lacks the needle, so the check reds. |
| GR6 | 4 | `.agents/skills/bench-craft-delegate/SKILL.md` requires every write charge from a spec to carry its stories' coverage rows | planned Require needle that replaces the user-directed needle, with canary `delegate-coverage-row-charge` retargeted | A skill that scopes the rule to a user-directed write lacks the needle, so the check reds. |
| GR7 | 4 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "A user-directed write-delegation from a spec carries" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the scoped sentence matches the needle, so the check reds. |
| GR8 | 6 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "the coordinator runs `bench commit` per worktree" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the coordinator commit matches the needle, so the check reds. |
| GR9 | 6 | `.agents/skills/bench-craft-delegate/SKILL.md` states that a ticket author commits its ticket on a lane pass | planned Require needle in `internal/anchors/registry_data.go` | A skill that drops the author commit lacks the needle, so the check reds. |
| GR10 | 6 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "Share a worktree only when a delegate's work depends on another's output." | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the dependency-only share rule matches the needle, so the check reds. |
| GR11 | 7 | `.agents/skills/bench-craft-delegate/SKILL.md` scopes the returned-tree coordinator probe to a user-directed write delegate | planned Require needle in `internal/anchors/registry_data.go` | A skill that applies the probe to every ticket lacks the needle, so the check reds. |
| GR12 | 7 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "A ticket delegate returns focused evidence" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the ticket-wide probe sentence matches the needle, so the check reds. |
| GR13 | 9 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "the destructive-git guard refuses it" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the false guard claim matches the needle, so the check reds. |
| GR14 | 9 | `.agents/skills/bench-craft-delegate/SKILL.md` states that the guard refuses only `git stash drop` and `git stash clear` | planned Require needle in `internal/anchors/registry_data.go` | A skill that drops the deny surface lacks the needle, so the check reds. |
| GR15 | 13 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "A diagnostic helper can inspect evidence, but it receives no implementation or repair assignment." | planned Forbid needle that replaces two Require rows, with the `TestImplementationContinuation` row updated | A skill that keeps the second copy matches the needle, so the check reds. |
| GR16 | 13 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "A spec-backed ticket goes to a fresh author session on its integration source" | planned Forbid needle that replaces the Require row in `internal/anchors/registry_retained_workflow.go`, with canary `delegated-per-ticket-author` retargeted | A skill that keeps the second copy matches the needle, so the check reds. |
| GR17 | 3 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` states under "Delegated author transfer" that a model change outside `--delegate` asks the reviewer first | planned RequireInSection needle in `internal/anchors/registry_retained_workflow.go` | A discipline that lets a trigger move the tier silently lacks the needle, so the check reds. |
| GR18 | 5 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "A user-directed write delegate treats `Writes:` as an expectation." | planned Forbid needle that replaces the RequireInSection row, with the `TestRetainedWorkflow` row updated | A discipline that keeps the scoped sentence matches the needle, so the check reds. |
| GR19 | 5 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` states under "In the charge" that every write delegate treats `Writes:` as an expectation | planned RequireInSection needle in `internal/anchors/registry_retained_workflow.go` | A discipline that drops the expectation rule lacks the needle, so the check reds. |
| GR20 | 8 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "A repair fence is the approved chunk union plus the exact paths that the review names." | planned Forbid needle that replaces the RequireInSection row in `internal/anchors/registry_charge_binding.go` | A discipline that keeps the chunk-union fence matches the needle, so the check reds. |
| GR21 | 8 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` states under "In the charge" that a repair fence is the affected ticket's `Writes:` line | planned RequireInSection needle in `internal/anchors/registry_charge_binding.go` | A discipline that drops the ticket fence lacks the needle, so the check reds. |
| GR22 | 10 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "Probe a tracked file that has pending changes with a copy aside." | planned Forbid needle in `internal/anchors/registry_data.go` | A discipline that keeps the copy-aside probe matches the needle, so the check reds. |
| GR23 | 10 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "mutated bytes against the copy aside" | planned Forbid needle in `internal/anchors/registry_data.go` | A discipline that keeps the copy-aside confirmation matches the needle, so the check reds. |
| GR24 | 10 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "restore with `cmp` against the copy aside" | planned Forbid needle that replaces the RequireInSection row in `internal/anchors/registry_data.go` | A discipline that keeps the `cmp` restore matches the needle, so the check reds. |
| GR25 | 11 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "may run in the main checkout" | planned Forbid needle in `internal/anchors/registry_data.go` | A discipline that keeps the main-checkout exception matches the needle, so the check reds. |
| GR26 | 12 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "`bench resume-clean`" | planned Forbid needle in `internal/anchors/registry_data.go` | A discipline that keeps the internal verb matches the needle, so the check reds. |
| GR27 | 12 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "After a rebase changes the source digest" | planned Forbid needle that replaces the RequireInSection row in `internal/anchors/registry_charge_binding.go` | A discipline that keeps the rebase premise matches the needle, so the check reds. |
| GR28 | 12 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` states under "Before the landing" that a `bench worktree merge` that changes the source digest repeats verification | planned RequireInSection needle in `internal/anchors/registry_charge_binding.go` | A discipline that drops the digest rule lacks the needle, so the check reds. |
| GR29 | 13 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "on the integration source after every ticket commit and before the next charge" | planned Forbid needle in `internal/anchors/registry_data.go` | A discipline that keeps the preflight copy matches the needle, so the check reds. |
| GR30 | 13 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` does not contain "a slow tool call is not a failed attempt" | planned Forbid needle in `internal/anchors/registry_data.go` | A discipline that keeps the attempt copy matches the needle, so the check reds. |
| GR31 | 14 | `.agents/skills/bench-craft-line/references/bounded-repair-policy.md` does not contain "This allowance applies to retained, full, delegated, unattended, and light-path implementation runs." | planned Forbid needle that replaces the RequireInSection row, with the `TestImplementationContinuation` row updated | A policy that keeps the retired mode list matches the needle, so the check reds. |
| GR32 | 14 | `.agents/skills/bench-craft-line/references/bounded-repair-policy.md` applies the allowance to every implementation run, the light path included | planned RequireInSection needle in `internal/anchors/registry_retained_workflow.go` | A policy that drops its scope sentence lacks the needle, so the check reds. |
| GR33 | 15 | `.agents/commands/bench-review-implementation.md` does not contain "writes one repair ticket when accepted repairs amend the coverage map" | planned Forbid needle that replaces two RequireInSection rows, with canaries `review-repair-ticket-owner` and `review-repair-ticket-covers` retargeted | A review phase that keeps the single repair ticket matches the needle, so the check reds. |
| GR34 | 15 | `.agents/commands/bench-review-implementation.md` states under "Review modes" that a coverage-map amendment updates each affected ticket's `Covers:` line under the plan-expansion policy | planned RequireInSection needle in `internal/anchors/registry_data.go` | A review phase that drops the per-ticket amendment lacks the needle, so the check reds. |
| GR35 | 16 | `.agents/commands/bench-final-check.md` does not contain "the author's final acceptance and integration command results" | planned Forbid needle in `internal/anchors/registry_data.go` | A phase that keeps the author as final performer matches the needle, so the check reds. |
| GR36 | 16 | `.agents/commands/bench-final-check.md` retains the orchestrator's final `integration-verification` results | planned Require needle in `internal/anchors/registry_data.go` | A phase that names no final performer lacks the needle, so the check reds. |
| GR37 | 17 | `.agents/commands/bench-final-check.md` sends a spec-backed red to a fresh repair session under `.bench/BENCH.md`'s repair rule | planned Require needle in `internal/anchors/registry_data.go` | A phase that lets the orchestrator fix a spec-backed red lacks the needle, so the check reds. |
| GR38 | 18 | `.agents/commands/bench-debug.md` does not contain "The coordinator validates the report and reslices" | planned Forbid needle in `internal/anchors/registry_data.go` | A debug phase that keeps the coordinator reslice matches the needle, so the check reds. |
| GR39 | 18 | `.agents/commands/bench-debug.md` routes an out-of-fence cause through the plan-expansion policy or the reviewer's split | planned Require needle in `internal/anchors/registry_data.go` | A debug phase that drops the route lacks the needle, so the check reds. |
| GR40 | 19 | `.agents/commands/bench-implement-spec.md` does not contain `For a ticket author, raise the effort and resume; a tier move asks the reviewer first.` | planned Forbid needle in `internal/anchors/registry_data.go` | A build phase that keeps the unqualified ladder copy matches the needle, so the check reds. |
| GR41 | 20 | `.agents/commands/bench-implement-spec.md` does not contain "Repeat delegated review only when a later delta or cross-chunk concern invalidates prior evidence." | planned Forbid needle in `internal/anchors/registry_data.go` | A build phase that keeps the review-repeat copy matches the needle, so the check reds. |
| GR42 | 20 | `.agents/commands/bench-review-implementation.md` does not contain "Repeat delegated review only for a later semantic delta" | planned Forbid needle in `internal/anchors/registry_data.go` | A review phase that keeps the review-repeat copy matches the needle, so the check reds. |
| GR43 | 20 | `.agents/commands/bench-review-implementation.md` does not contain "A delegated chunk review starts after every ticket of the chunk reaches the integrated chunk tip." | planned Forbid needle that replaces the Require row in `internal/anchors/registry_retained_workflow.go`, with the `TestRetainedWorkflow` row moved to GR116 | A review phase that keeps the start-rule copy matches the needle, so the check reds. |
| GR44 | 20 | `.agents/commands/bench-review-implementation.md` does not contain "After the last chunk, the orchestrator reconciles overall acceptance and integration before landing." | planned Forbid needle in `internal/anchors/registry_data.go` | A review phase that keeps the reconciliation copy matches the needle, so the check reds. |
| GR45 | 21 | `.agents/skills/bench-craft-synthesis/SKILL.md` does not contain "`bench commit` gates the tree it lands" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the lane-as-gate claim matches the needle, so the check reds. |
| GR46 | 21 | `.agents/skills/bench-craft-synthesis/SKILL.md` takes the prose-only green verdict from the whole-project gate | planned Require needle in `internal/anchors/registry_data.go` | A skill that names no gate source lacks the needle, so the check reds. |
| GR47 | 21 | `.agents/commands/bench-final-check.md` does not contain "This command gates and commits them atomically." | planned Forbid needle in `internal/anchors/registry_data.go` | A phase that keeps the lane-as-gate claim matches the needle, so the check reds. |
| GR48 | 21 | `.agents/commands/bench-final-check.md` does not contain "The commit already is the gate run" | planned Forbid needle in `internal/anchors/registry_data.go` | A phase that keeps the lane-as-gate claim matches the needle, so the check reds. |
| GR49 | 21 | `.agents/commands/bench-final-check.md` does not contain "the oracle run and landing are one command" | planned Forbid needle in `internal/anchors/registry_data.go` | A phase that keeps the one-command claim matches the needle, so the check reds. |
| GR50 | 21 | `.agents/commands/bench-final-check.md` states that `bench worktree land` runs the whole-project gate on work that `bench commit` committed on a lane pass | planned Require needle in `internal/anchors/registry_data.go` | A phase that names no landing gate lacks the needle, so the check reds. |
| GR51 | 22 | `bench help` describes `bench commit` as a lane run, and as a gate run only when the project declares no lane | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A help row that keeps "gate, then commit named paths on green" fails the complete-inventory comparison. |
| GR52 | 22 | `bench commit --help` does not print "gate the exact composed snapshot" | `internal/commit/dry_run_test.go` (`TestHelpAdvertisesDryRun`, extended) | A help text that keeps the gate wording fails the extended assertion. |
| GR53 | 22 | `bench commit --help` names the declared lane on its `--dry-run` line | `internal/commit/dry_run_test.go` (`TestHelpAdvertisesDryRun`, extended) | A help text that drops the lane fails the extended assertion. |
| GR54 | 23 | `.bench/BENCH-reference.md` does not contain "the rule is guidance, not a hook" | planned Forbid needle in `internal/anchors/registry_data.go` | A reference that keeps the unenforced claim matches the needle, so the check reds. |
| GR55 | 23 | ADR 0014 states that the commit verb refuses the primary checkout and that a file-write guard refuses a tracked-path write there | review-owned: no anchor registry row names an ADR | The Spec axis grades the ADR against this spec's decisions. |
| GR56 | 24 | `.bench/BENCH-reference.md` does not contain "A stale Bench executable is rebuilt, and the landing re-runs under it." | planned Forbid needle in `internal/anchors/registry_data.go` | A reference that keeps the rebuild claim matches the needle, so the check reds. |
| GR57 | 25 | `.bench/BENCH-reference.md` names the reviewer as the person who merges the destination into the source worktree with raw Git | planned Require needle in `internal/anchors/registry_data.go` | A reference that assigns the raw merge to an agent lacks the needle, so the check reds. |
| GR58 | 26 | `.agents/commands/bench-final-check.md` runs `bench spec retire <slug>` in a Bench worktree and lands its `spec-retire: <slug>` commit | planned Require needle in `internal/anchors/registry_data.go` | A phase that names no worktree for retirement lacks the needle, so the check reds. |
| GR59 | 27 | `.agents/commands/bench-final-check.md` does not contain "Do not run another gate or commit just to capture the retro" | planned Forbid needle that replaces the Require row in `internal/anchors/registry_data.go` | A phase that keeps the no-commit rule matches the needle, so the check reds. |
| GR60 | 27 | `.agents/commands/bench-final-check.md` does not contain "Capture the retro without another" | planned Forbid needle in `internal/anchors/registry_data.go` | A phase that keeps the report copy of the no-commit rule matches the needle, so the check reds. |
| GR61 | 27 | `.agents/commands/bench-final-check.md` commits a tracked retro and its scorecard updates with the phase close | planned Require needle in `internal/anchors/registry_data.go` | A phase that drops the tracked commit lacks the needle, so the check reds. |
| GR62 | 28 | `.agents/commands/bench-final-check.md` leaves an ignored retro local until the next reviewer-approved capture drain | planned Require needle in `internal/anchors/registry_data.go` | A phase that drops the ignored route lacks the needle, so the check reds. |
| GR63 | 27 | `.bench/BENCH-reference.md` does not contain "owns their reviewed drain and its capture commit" | planned Forbid needle in `internal/anchors/registry_data.go` | A reference that keeps the drain-only commit matches the needle, so the check reds. |
| GR64 | 29 | `.agents/commands/bench-final-check.md` does not contain "rewrite `capture/retros/<spec-slug>.md` in full" | planned Forbid needle that replaces the Require row in `internal/anchors/registry_data.go` | A phase that keeps the rewrite rule matches the needle, so the check reds. |
| GR65 | 29 | `.agents/commands/bench-final-check.md` writes the retro once with `bench retro <slug> --body` after `bench retro <slug> --scaffold` | planned Require needle in `internal/anchors/registry_data.go` | A phase that drops the verb route lacks the needle, so the check reds. |
| GR66 | 29 | `.agents/commands/bench-final-check.md` does not contain "Use these headings exactly" | planned Forbid needle that replaces the nine heading Require rows in `internal/anchors/registry_data.go` | A phase that keeps the heading copy matches the needle, so the check reds. |
| GR67 | 29 | `.agents/commands/bench-final-check.md` does not contain "surface, claim, status, confidence, label, and model, effort, and role" | planned Forbid needle in `internal/anchors/registry_calibration.go` | A phase that keeps the column copy matches the needle, so the check reds. |
| GR68 | 30 | `.agents/skills/bench-craft-line/SKILL.md` does not contain "ceiling, not a binding" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the per-story ceiling matches the needle, so the check reds. |
| GR69 | 30 | `.agents/skills/bench-craft-line/SKILL.md` does not contain "Re-run the decision table per ticket at charge time." | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the per-ticket re-run matches the needle, so the check reds. |
| GR70 | 30 | `.agents/skills/bench-craft-line/SKILL.md` binds every ticket author to the spec's one declared line | planned Require needle in `internal/anchors/registry_data.go` | A skill that drops the binding lacks the needle, so the check reds. |
| GR71 | 30 | `.agents/skills/bench-craft-line/SKILL.md` does not contain "use the highest tier any story needs" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the per-story collapse matches the needle, so the check reds. |
| GR72 | 30 | `.agents/skills/bench-craft-spec/SKILL.md` says under "Template" that the approval table covers the implementation line | planned RequireInSection needle that replaces the "stories and their lines" needle in `internal/anchors/registry_data.go` and keeps its diagnostic | A template that drops the line from the table lacks the needle, so the check reds. |
| GR73 | 30 | `.agents/skills/bench-craft-spec/SKILL.md` does not contain "stories and their lines" | planned Forbid needle in `internal/anchors/registry_data.go` | A template that keeps per-story lines matches the needle, so the check reds. |
| GR74 | 31 | `.agents/skills/bench-craft-line/SKILL.md` makes the step 3 escalation obey the step 2 tier-move rule | planned Require needle in `internal/anchors/registry_data.go` | A ladder whose step 3 escalates without the rule lacks the needle, so the check reds. |
| GR75 | 31 | `.agents/skills/bench-craft-line/SKILL.md` limits the top-tier pause of step 5 to work outside `--delegate` | planned Require needle in `internal/anchors/registry_data.go` | A ladder whose step 5 ignores `--delegate` lacks the needle, so the check reds. |
| GR76 | 32 | `.bench/BENCH.md` does not contain "live in `.bench/BENCH-reference.md`" | planned Forbid needle in `internal/anchors/registry_data.go` | A guide that keeps the missing-list pointer matches the needle, so the check reds. |
| GR77 | 32 | `.bench/BENCH.md` names the command registry as the owner of the plumbing subcommands | planned Require needle in `internal/anchors/registry_data.go` | A guide that drops the owner lacks the needle, so the check reds. |
| GR78 | 32 | `AGENTS.md` does not contain "the four invariants, how the pieces fit" | planned Forbid needle in `internal/anchors/registry_data.go` | An agreement that keeps the wrong content list matches the needle, so the check reds. |
| GR79 | 32 | `AGENTS.md` names `.bench/BENCH-reference.md` as the holder of how the pieces fit and the skills index | planned Require needle in `internal/anchors/registry_data.go` | An agreement that drops the reference pointer lacks the needle, so the check reds. |
| GR80 | 33 | `.bench/BENCH-reference.md` names `bench gate-prose` under "Plumbing subcommands" as the internal verb that a session runs | planned RequireInSection needle in `internal/anchors/registry_data.go` | A reference that drops the exception lacks the needle, so the check reds. |
| GR81 | 34 | `projects/benchkit.md` does not contain "A red run prints one `failures[N]{phase,line}` table" | planned Forbid needle in `internal/anchors/registry_data.go` | A profile that keeps the gate output copy matches the needle, so the check reds. |
| GR82 | 34 | `.bench/BENCH-reference.md` owns the green-run gate output sentence | `internal/anchors/registry_data_test.go` (`TestBoundedGateOutputAnchorTuples`), with its profile tuple moved to the reference | A registry that keeps the tuple on the profile fails the tuple test. |
| GR83 | 34 | `projects/benchkit.md` does not contain "The `census` signal counts raw calls per assignment" | planned Forbid needle that replaces the profile Require row, with `TestCensusDutyAnchorsRedOnRemoval` updated | A profile that keeps the census copy matches the needle, so the check reds. |
| GR84 | 34 | `AGENTS.md` does not contain "`bench handoff` rewrites only the calling worktree's assignment section." | planned Forbid needle that replaces the `AGENTS.md` Require row in `internal/anchors/registry_data.go` | An agreement that keeps the verb copy matches the needle, so the check reds. |
| GR85 | 34 | `.bench/BENCH-reference.md` contains "`bench handoff` rewrites only the calling worktree's assignment section." | planned Require needle moved from `AGENTS.md`, with canary `agents-handoff-section-rule` retargeted | A reference that drops the verb fact lacks the needle, so the check reds. |
| GR86 | 35 | `.bench/BENCH-reference.md` does not contain "`bench handoff [--harness <name>]" | planned Forbid needle in `internal/anchors/registry_data.go` | A reference that keeps the grammar copy matches the needle, so the check reds. |
| GR87 | 35 | `.bench/BENCH-reference.md` does not contain "`bench worktree reset --to <commit> <target>` plans" | planned Forbid needle in `internal/anchors/registry_data.go` | A reference that keeps the grammar copy matches the needle, so the check reds. |
| GR88 | 35 | `.bench/BENCH-reference.md` does not contain "`bench retro <slug> (--body <markdown>" | planned Forbid needle in `internal/anchors/registry_data.go` | A reference that keeps the grammar copy matches the needle, so the check reds. |
| GR89 | 36 | `.bench/BENCH-reference.md` states that the plan amendment makes the authored version 1 fence a version 2 plan before the first dispatch | planned Require needle in `internal/anchors/registry_data.go` | A reference that names only version 1 lacks the needle, so the check reds. |
| GR90 | 37 | `.bench/BENCH-reference.md` does not contain "`.claude/skills/` carries only the" | planned Forbid needle in `internal/anchors/registry_data.go` | A reference that keeps the false link claim matches the needle, so the check reds. |
| GR91 | 37 | `.claude/README.md` states that `.claude/skills/` links every `.agents/skills/` skill that has no same-named command | planned Require needle in `internal/anchors/registry_data.go` | A README that keeps the craft-only claim lacks the needle, so the check reds. |
| GR92 | 38 | `.agents/commands/bench-drain.md` names the `Occurrences:` line in `roadmap/FT<n>.md` for each pending pair | `internal/conformance/recurrence_maintenance_contract_test.go` (`checkRecurrenceMaintenanceContract` and its bite table) | A drain that keeps `ROADMAP.md` fails the updated contract string. |
| GR93 | 38 | `.agents/commands/bench-drain.md` does not contain "`Occurrences:` line in `ROADMAP.md`" | planned Forbid needle in `internal/anchors/registry_data.go` | A drain that keeps the wrong ledger file matches the needle, so the check reds. |
| GR94 | 38 | `.agents/commands/bench-drain.md` does not contain "(the AGENTS.md rule)" | planned Forbid needle in `internal/anchors/registry_data.go` | A drain that keeps the wrong rule owner matches the needle, so the check reds. |
| GR95 | 38 | `.agents/commands/bench-drain.md` names `bench handoff` as the last write of the drain | `internal/conformance/recurrence_maintenance_contract_test.go` (`checkRecurrenceMaintenanceContract` and its bite table) | A drain that keeps a hand-written handoff fails the updated contract string. |
| GR96 | 38 | `.agents/commands/bench-drain.md` dates the `main` handoff section by the file's write time | `internal/conformance/recurrence_maintenance_contract_test.go` (`checkRecurrenceMaintenanceContract` and its bite table) | A drain that dates the whole handoff by write time fails the updated contract string. |
| GR97 | 38 | `.agents/commands/bench-assess.md` does not contain "or into `capture/IDEAS.md`" | planned Forbid needle in `internal/anchors/registry_data.go` | A phase that keeps the hand-append option matches the needle, so the check reds. |
| GR98 | 39 | `.agents/skills/bench-craft-skills/SKILL.md` does not contain "and the phase adapters are not" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the blanket adapter rule matches the needle, so the check reds. |
| GR99 | 39 | `.agents/skills/bench-craft-skills/SKILL.md` points each phase adapter trigger to the invocation-policy account in `.bench/BENCH-reference.md` | planned Require needle in `internal/anchors/registry_data.go` | A skill that drops the pointer lacks the needle, so the check reds. |
| GR100 | 39 | `.agents/skills/bench-craft-cli/SKILL.md` names the ambiguous-name re-query disclosure in its `bench consumers` row | planned Require needle in `internal/anchors/registry_data.go` | A row that keeps the terminal-only claim lacks the needle, so the check reds. |
| GR101 | 39 | `.agents/skills/bench-craft-spec/SKILL.md` does not contain "more than four stories" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the literal count matches the needle, so the check reds. |
| GR102 | 39 | `.agents/skills/bench-craft-spec/SKILL.md` does not contain "Each planned chunk has a stable ID and names its tickets" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the chunk-contract copy matches the needle, so the check reds. |
| GR103 | 39 | `.agents/skills/bench-craft-tdd/references/tests.md` defers to the project's own test-expectation standard | planned Require needle in `internal/anchors/registry_data.go` | A reference that drops the deferral lacks the needle, so the check reds. |
| GR104 | 39 | `.agents/skills/bench-craft-adr/SKILL.md` points its no-paths rule to invariant 3 of `.bench/BENCH.md` | planned Require needle in `internal/anchors/registry_data.go` | A skill that drops the pointer lacks the needle, so the check reds. |
| GR105 | 40 | Every budgeted guidance file that this spec edits stays inside its budget in `projects/benchkit.md` | `bench test --check guidance-prose-budgets` | A file that grows past its budget reds the check. |
| GR106 | 30, 31 | `.agents/skills/bench-craft-line/SKILL.md` keeps the Require needle "Outside `--delegate`, a tier move of a fresh ticket author asks the reviewer first." | existing Require needle in `internal/anchors/registry_data.go` | A ladder edit that drops the step 2 rule removes the needle, so the check reds. |
| GR107 | 42 | `bench anchors AGENTS.md` on the fixture prints each `AGENTS.md` registry row with its registry kind, and a Forbid row reads line 0 | `cmd/bench/anchor_help_test.go` (`TestAnchorsReportsNeedleLines`) | A test that types every row as a planted Require fails on the first Forbid row. |
| GR108 | 42 | `bench anchors AGENTS.md` with dropped or absent needles gives a missing-file diagnostic for each Require row and none for a Forbid row | `cmd/bench/anchor_help_test.go` (`TestAnchorsReportsAbsentNeedles`) | A test that counts one diagnostic per row fails on the first Forbid row. |
| GR109 | 6 | `.claude/agents/bench-writer.md` tells a ticket author or a repair author to commit on a lane pass and not to run `bench worktree land` | planned Require needle in `internal/anchors/registry_data.go` | An agent body that stops every author at an uncommitted diff lacks the needle, so the check reds. |
| GR110 | 1 | `.claude/agents/bench-writer.md` does not contain "The user directed this delegation" | planned Forbid needle in `internal/anchors/registry_data.go` | An agent body that keeps the user-directed premise matches the needle, so the check reds. |
| GR111 | 6 | `.claude/agents/bench-writer.md` does not contain "Stop at a diff that is ready, with the focused checks green." | planned Forbid needle in `internal/anchors/registry_data.go` | An agent body that keeps the uncommitted stop matches the needle, so the check reds. |
| GR112 | 6 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain "stops at diff-ready" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the diff-ready stop matches the needle, so the check reds. |
| GR113 | 6 | `.agents/skills/bench-craft-delegate/SKILL.md` does not contain `Stop at diff ready;` | planned Forbid needle in `internal/anchors/registry_data.go` | A charge example that keeps the uncommitted stop matches the needle, so the check reds. |
| GR114 | 1 | `.claude/README.md` does not contain "`bench-writer` runs a user-directed write delegation" | planned Forbid needle in `internal/anchors/registry_data.go` | A README that keeps the user-directed copy matches the needle, so the check reds. |
| GR115 | 10 | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` states under "Read-only returns" that a shared-worktree reader probes only through `bench probe` and reads its `restored` cell | planned RequireInSection needle in `internal/anchors/registry_data.go` | A discipline that drops the reader restore duty lacks the needle, so the check reds. |
| GR116 | 20 | `.bench/BENCH.md` contains "Every ticket contribution reaches the integrated chunk tip before that chunk's review begins." | planned Require needle in `internal/anchors/registry_retained_workflow.go`, with the `TestRetainedWorkflow` row and canary `delegated-chunk-tip-review` retargeted | A guide that drops the chunk-tip start rule lacks the needle, so the check reds. |
| GR117 | 29 | `.agents/commands/bench-final-check.md` states under "Capture the implementation retro" that the calibration table takes one row per labeled claim | planned RequireInSection needle in `internal/anchors/registry_calibration.go`, with canary `calibration-retro-table-duty` retargeted | A phase that drops the one-row-per-claim duty lacks the needle, so the check reds. |
| GR118 | 21 | `.agents/commands/bench-final-check.md` does not contain "gate-then-commit path" | planned Forbid needle in `internal/anchors/registry_data.go` | A phase that keeps the gate-then-commit claim matches the needle, so the check reds. |
| GR119 | 21 | `.agents/commands/bench-final-check.md` does not contain "then runs the gate and commits only on green" | planned Forbid needle in `internal/anchors/registry_data.go` | A phase that keeps the commit-gates claim matches the needle, so the check reds. |
| GR120 | 27 | `.agents/commands/bench-final-check.md` does not contain "The retro leaves through the next reviewer-approved capture drain." | planned Forbid needle in `internal/anchors/registry_data.go` | A phase that keeps the drain-only exit matches the needle, so the check reds. |
| GR121 | 37 | `.claude/README.md` does not contain "links only the `bench-craft-*` skills" | planned Forbid needle in `internal/anchors/registry_data.go` | A README that keeps the craft-only claim matches the needle, so the check reds. |
| GR122 | 39 | `.agents/skills/bench-craft-adr/SKILL.md` does not contain "No file paths, no code snippets" | planned Forbid needle in `internal/anchors/registry_data.go` | A skill that keeps the invariant copy matches the needle, so the check reds. |
| GR123 | 34 | `projects/benchkit.md` does not contain "A green run prints one `phases[N]{phase,verdict,elapsed_ms}` table" | planned Forbid needle that replaces the profile Require row in `internal/anchors/registry_data.go` | A profile that keeps the green-run copy matches the needle, so the check reds. |
| GR124 | 34 | `projects/benchkit.md` does not contain "The complete phase stream goes to `.logs/gate-<run>.out`" | planned Forbid needle in `internal/anchors/registry_data.go` | A profile that keeps the stream-log copy matches the needle, so the check reds. |

Not covered: story 41 — the reviewed exclusion changes no behavior, and a Won't handle line records it.

### Edge inventory

The canonical edge classes, walked at the anchor seam and the help seam:

- A retired sentence that returns: each retired or wrong sentence takes its own Forbid row.
- A second copy that returns: each removed copy takes a Forbid row, and review grades a paraphrase.
- A needle that wraps across two physical lines: each new Require needle sits on one line, per the prose reference. A Forbid needle can span a break, because the matcher collapses whitespace.
- A canary whose `old` bytes change: the canary names the new bytes or plants the forbidden bytes. No ticket removes a canary directory.
- A Forbid row that replaces the last guard of a surviving duty: the owner sentence of that duty takes a Require row.
- A Forbid row on `AGENTS.md`: the `bench anchors` tests take each row's kind from the registry.
- A conformance table that pairs a diagnostic with a needle: the table row moves with its needle.
- A budgeted file at its budget: the edit stays line-neutral, and GR105 reds growth.
- An imported instruction file: the classifier refuses an agent edit to `.bench/BENCH.md` and `AGENTS.md`, so ticket 7 names a reviewer approval step.
- The kit audience: the commands, the skills, the reference, the guide, and `.claude/` ship to linked repositories. `AGENTS.md`, the profile, and the ADRs serve this repository only.
- A tracked versus an ignored retro path: the capture writer refuses a tracked path in the primary checkout. It writes the primary copy of an ignored path.
- A project with no declared lane: `bench commit` runs the whole-project gate, and GR51 and GR53 name that case.
- Hostile input: this spec changes guidance, anchor rows, and two help strings, so no operator input reaches a new surface.

**Won't handle** lines:

- `internal/reviewrecord` — this spec changes no record field or trigger. GR17 survives.
- The heading "Delegate or retain" in `craft-delegate` — it is a section key for existing anchors, and a rename changes no rule. GR5 survives.
- The done-claim probe under "Verifying the done-claim" in `craft-delegate` — the audit did not name it. GR11 survives.
- The sentence "Neither role receives implementation or repair authorship." in the build phase — the audit did not name it. GR15 survives.
- The scaffold sentence in `AGENTS.md` — it is project content beside the kit pointer. GR65 survives.
- The phrase "`Occurrence:` ledger" in the drain's reconcile section — the audit did not name it, and an anchor pins it. GR92 survives.
- The considered option in ADR 0014 for a default-branch hook — it stays deferred and true. GR55 survives.
- Retired specs, retros, reviews, roadmap rows, decision maps, and the changelog — they are history. GR54 survives.
- Frozen canary fixture copies under `files/` — they are planted content, not live guidance. GR85 survives.

## Ownership fences

- `reviews/workflow-guidance-repair.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.claude/agents/bench-writer.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-line/references/bounded-repair-policy.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-debug.md`
- `.agents/commands/bench-final-check.md`
- `.agents/skills/bench-craft-synthesis/SKILL.md`
- `.bench/BENCH-reference.md`
- `docs/adr/0014-main-receives-writes-only-through-landings.md`
- `cmd/bench/main.go`
- `cmd/bench/help_inventory_test.go`
- `internal/commit/commit.go`
- `internal/commit/dry_run_test.go`
- `.agents/skills/bench-craft-line/SKILL.md`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.bench/BENCH.md`
- `AGENTS.md`
- `projects/benchkit.md`
- `.claude/README.md`
- `cmd/bench/testdata/anchors/fixture-repo/AGENTS.md`
- `cmd/bench/anchor_help_test.go`
- `.agents/commands/bench-drain.md`
- `.agents/commands/bench-assess.md`
- `.agents/skills/bench-craft-skills/SKILL.md`
- `.agents/skills/bench-craft-cli/SKILL.md`
- `.agents/skills/bench-craft-tdd/references/tests.md`
- `.agents/skills/bench-craft-adr/SKILL.md`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_charge_binding.go`
- `internal/anchors/registry_charge_binding_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_front_door.go`
- `internal/anchors/registry_front_door_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/implementation_continuation_test.go`
- `internal/conformance/retained_workflow_test.go`
- `internal/conformance/recurrence_maintenance_contract_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/claude-agent-definitions/`
- `tests/canary/skills-index-command-adapters/`
- `tests/canary/docs-currency-token-diet/`
- `tests/canary/load-validity-metadata/`
- `tests/canary/line-routing/`
- `tests/canary/skill-description-budgets/`
- `tests/canary/row-next-grammar/`

Build preflight requires the closure entries below beside the paths above. They preserve coverage, and they authorize no change to the command inventory or its enforcement.

- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/guidance-prose-budgets/`
- `tests/canary/package-core-guard/`

Reviewer disposition: pending sign-off.

## Out of scope

- Publication of `bench gate-prose` in `bench help`: it changes the public inventory and its AXI class — about 6 edits and 3 gate runs.
- The done-claim probe rule for a spec-backed ticket author in `craft-delegate`: the audit did not name it — about 3 edits and 2 gate runs.
- A hook that refuses a commit on the default branch: ADR 0014 keeps it deferred — about 8 edits and 3 gate runs.

## Further notes

### Reviewer decisions

The reviewer closed these decisions on 2026-09-24:

1. The spec covers all 62 audit findings.
2. The spec also covers the `bench commit` help strings and an in-place correction of ADR 0014's enforcement consequence.
3. A tracked retro commits with the phase close, and an ignored retro stays local for `/bench-drain`. The capture-inbox rule of the capture writer is the source.
4. A fresh ticket author runs as `bench-writer` and commits its own ticket on a lane pass. The build phase and the fresh-ticket-authors delegate plan are the source.
5. A build has one declared implementation line. `.bench/BENCH.md` is the source.

### Open reviewer decision

The staged spec `shared-delegate-startup` needs a respec before its build. Its spec says "Repairs remain with the recorded author (SP23)", and its ticket 2 keeps repairs with the recorded author. Its tickets 1 and 3 write canaries that this spec retargets. This spec does not edit that spec.

### Contestable calls

These calls are the author's. Each one is open to reviewer veto:

1. The reference owns the gate output, the census signal, and the handoff verb facts, because a linked repository reads the reference and not this profile.
2. `bench gate-prose` stays internal, and the reference names it as the one internal verb that a session runs.
3. A second copy that the registry pinned takes a Forbid row. A second copy that no row pinned also takes a Forbid row, so every finding has a red-capable row.
4. The drain dating fix keeps write-time dating for the `main` section, because the drain writes its handoff after the landing releases its worktree.
5. The fix for the skill-link claim changes `.claude/README.md` as well, because the reference restates that file.
6. The ADR 0014 correction changes only the two enforcement sentences of its fifth paragraph.

### Bootstrap

This spec's own build runs before its guidance lands, under today's rules. The recommended build is `/bench-implement-spec --full specs/workflow-guidance-repair/spec.md` on opus/high, with one fresh author for each ticket and no tier escalation.

The staged completion plan below stays at version 1. Before the first dispatch, the orchestrator amends it to the version 2 delegate form:

- The plan version becomes 2, with mode `delegate`, an author limit of 1, the run id, and the orchestrator session.
- Each chunk's verification splits into one set for each ticket, and each set names its ticket.
- The final verification names no ticket.

Before ticket 1 and ticket 7, the reviewer grants a one-time permission rule for the imported or adapter files that those tickets edit. The reviewer removes each rule after its ticket commits. The files are `.bench/BENCH.md`, `AGENTS.md`, `.claude/agents/bench-writer.md`, and `.claude/README.md`.

No row claims trusted execution or refusal before execution, so the bootstrap-authority rule adds no hop here.

### Reader sweep

- The sweep searched the exact bytes of every changed sentence across the tree at `78602e89`, except history directories. Each reader takes a row, a fence entry, or a Won't handle line.
- A byte search misses a paraphrase. The review round found three: `.claude/README.md:22-23`, `.claude/agents/bench-writer.md:9`, and `.claude/agents/bench-writer.md:20`. A second sweep for "user-directed write" and "diff ready" added `craft-delegate`'s charge example. Each one now takes a row.
- `internal/conformance/fixture_bite_test.go` pins exact bytes of `craft-delegate` and `craft-spec`, with their line breaks. It also pins the diagnostic of the `craft-spec` approval-paragraph row. The tickets keep those bytes and that diagnostic.
- `cmd/bench/anchor_help_test.go` expects only planted Require rows for `AGENTS.md`, so ticket 7 changes it.
- `retained_workflow_test.go`, `implementation_continuation_test.go`, and `recurrence_maintenance_contract_test.go` pair diagnostics with changed sentences.
- `registry_data_test.go` pins the gate output tuples, the census tuples, and the handoff tuple.
- `cmd/bench/anchor_help_test.go` derives its needles from `testdata/anchors/fixture-repo/AGENTS.md`, so that fixture drops the moved handoff line.
- `cmd/bench/help_inventory_test.go` holds the only other copy of the `bench commit` help row. `internal/commit/commit.go` holds the only copy of the `--dry-run` line.
- These live-mutation canaries name changed bytes: `delegate-coverage-row-charge`, `delegated-per-ticket-author`, `review-repair-ticket-owner`, `review-repair-ticket-covers`, `delegated-chunk-tip-review`, `calibration-retro-table-duty`, and `agents-handoff-section-rule`.
- The canaries `delegate-probe-mutated-bytes` and the four transfer-trigger canaries keep valid bytes by decision.
- The staged spec `shared-delegate-startup` fences some of the same guidance files. The spec that lands second rereads the other's changes.
- The link rule of `bench link` mirrors every skill that has no same-named command into `.claude/skills/`. That rule is the producer for GR91.
- `bench status` dates the `main` handoff section by the file's write time and a worktree section by its branch. That code is the producer for GR96.

### Flagged additions

These additions go past the audit table:

- The `.claude/README.md` fix beside the reference fix for the skill-link claim.
- The contract-string changes in `recurrence_maintenance_contract_test.go` for the drain rows.
- The tuple moves in `registry_data_test.go` for the gate output and census rows.
- The fixture line removal in `cmd/bench/testdata/anchors/fixture-repo/AGENTS.md`, and the GR79 needle planted there.
- The `craft-delegate` charge example, re-keyed to one ticket of its spec, and its Forbid row GR113.
- New Require rows for files that no anchor row names today: the agent definition, the Claude README, `craft-skills`, `craft-cli`, `craft-synthesis`, `craft-adr`, and the tests reference.

### Source-sentence-to-row table

| audit # | defect site at `78602e89` | rows |
|---|---|---|
| 1 | `bench-review-implementation.md:49` | GR33, GR34 |
| 2 | `BENCH-reference.md:205-207` | GR54 |
| 3 | `delegation-discipline.md:109-111` | GR25 |
| 4 | `bench-craft-synthesis/SKILL.md:64-66` | GR45, GR46 |
| 5 | `bench-final-check.md:179-184` | GR48, GR49, GR50, GR119 |
| 6 | `bench-craft-line/SKILL.md:57` | GR68, GR69, GR70 |
| 7 | `bench-craft-delegate/SKILL.md:17` | GR4, GR5 |
| 8 | `bench-craft-delegate/SKILL.md:102` | GR13, GR14 |
| 9 | `bench-final-check.md:16-18` | GR35, GR36 |
| 10 | `bench-drain.md:125-127` | GR92, GR93 |
| 11 | `BENCH-reference.md:209-210` | GR56 |
| 12 | `bench-final-check.md:161-164`, `BENCH-reference.md:45-48` | GR59, GR60, GR61, GR62, GR63, GR120 |
| 13 | `bench-final-check.md:39-41` | GR47, GR50, GR118 |
| 14 | `BENCH-reference.md:225-226` | GR57 |
| 15 | `bench-craft-delegate/SKILL.md:28` | GR1, GR2, GR3, GR109, GR110, GR111, GR114 |
| 16 | `bench-craft-delegate/SKILL.md:99-100` | GR8, GR9, GR112, GR113 |
| 17 | `delegation-discipline.md:120-123` | GR22, GR23 |
| 18 | `bench-implement-spec.md:68` | GR40 |
| 19 | `bench-craft-skills/SKILL.md:20-22` | GR98, GR99 |
| 20 | `bench-craft-delegate/SKILL.md:53` | GR6, GR7 |
| 21 | `bench-craft-delegate/SKILL.md:97` | GR10 |
| 22 | `bench-craft-line/SKILL.md:81-82` | GR74, GR106 |
| 23 | `bench-final-check.md:88-91` | GR64, GR65 |
| 24 | `delegation-discipline.md:95-98` | GR17 |
| 25 | `delegation-discipline.md:89` | GR20, GR21 |
| 26 | `bench-craft-line/SKILL.md:87-89` | GR75 |
| 27 | `bench-final-check.md:207-210` | GR37 |
| 28 | `bench-drain.md:273` | GR94 |
| 29 | `AGENTS.md:9-11` | GR78, GR79, GR107, GR108 |
| 30 | `BENCH-reference.md:68-69`, `.claude/README.md:8` | GR90, GR91, GR121 |
| 31 | `bench-assess.md:24` | GR97 |
| 32 | `bench-final-check.md:95-113`, `:125` | GR66, GR67, GR117 |
| 33 | `bench-craft-delegate/SKILL.md:13` | GR15 |
| 34 | `BENCH-reference.md:373-384`, `benchkit.md:475-486` | GR81, GR82, GR123, GR124 |
| 35 | `bench-review-implementation.md:45` | GR42 |
| 36 | `bench-implement-spec.md:60` | GR41 |
| 37 | `delegation-discipline.md:207` | GR27, GR28 |
| 38 | `delegation-discipline.md:106` | GR26 |
| 39 | `delegation-discipline.md:206` | GR29 |
| 40 | `bench-craft-line/SKILL.md:62` | GR71 |
| 41 | `bench-craft-spec/SKILL.md:155` | GR72, GR73 |
| 42 | `bench-craft-delegate/SKILL.md:19` | GR16 |
| 43 | `bench-craft-cli/SKILL.md:84` | GR100 |
| 44 | `BENCH.md:17-18` | GR76, GR77 |
| 45 | `BENCH-reference.md:39` | GR83 |
| 46 | `bench-debug.md:156-157` | GR38, GR39 |
| 47 | `bench-review-implementation.md:58` | GR43, GR116 |
| 48 | `bench-drain.md:294-299` | GR95, GR96 |
| 49 | `AGENTS.md:116` | GR80 |
| 50 | `AGENTS.md:77-81` | GR84, GR85, GR107, GR108 |
| 51 | `delegation-discipline.md:177-179` | GR30 |
| 52 | `bench-craft-spec/SKILL.md:45` | GR101 |
| 53 | `delegation-discipline.md:166-167` | GR24, GR115 |
| 54 | `bench-final-check.md:56-59` | GR58 |
| 55 | `delegation-discipline.md:36` | GR18, GR19 |
| 56 | `bounded-repair-policy.md:6` | GR31, GR32 |
| 57 | `BENCH-reference.md:275`, `:255`, `:282` | GR86, GR87, GR88 |
| 58 | `BENCH-reference.md:318-319` | GR89 |
| 59 | `bench-craft-tdd/references/tests.md:11-15` | GR103 |
| 60 | `bench-craft-spec/SKILL.md:70` | GR102 |
| 61 | `bench-craft-delegate/SKILL.md:70-71` | GR11, GR12 |
| 62 | `bench-craft-adr/SKILL.md:43` | GR104, GR122 |
| extra a | `cmd/bench/main.go:132`, `internal/commit/commit.go:225` | GR51, GR52, GR53 |
| extra b | ADR 0014 fifth paragraph | GR55 |

### Pre-review proof checklist

- `Cited symbols`: `reviewrecord.Triggers` (`internal/reviewrecord/delegated.go:44`), `retros.RequiredHeadings` (`internal/retros/retros.go:53`), `retros.CalibrationHeader` (`internal/retros/retros.go:36`), `gate.LaneForCommit` (`internal/gate/lane_select.go:104`), `checkRecurrenceMaintenanceContract` (`internal/conformance/recurrence_maintenance_contract_test.go:16`), `claudeSkillsEntries` (`internal/adopt/link.go:212`), `ignoredHandoffAgeRows` (`internal/status/handoff.go:79`), `TestBoundedGateOutputAnchorTuples` (`internal/anchors/registry_data_test.go:55`), `TestCensusDutyAnchorsRedOnRemoval` (`internal/anchors/registry_data_test.go:260`), `TestHelpInventoryIsComplete` (`cmd/bench/help_inventory_test.go:40`), `TestHelpAdvertisesDryRun` (`internal/commit/dry_run_test.go:58`).
- `Cited symbols`, continued: `anchorsFixtureNeedles` (`cmd/bench/anchor_help_test.go:25`), `TestAnchorsReportsNeedleLines` (`cmd/bench/anchor_help_test.go:40`), `TestAnchorsReportsAbsentNeedles` (`cmd/bench/anchor_help_test.go:65`), `anchorKindName` (`cmd/bench/anchors_command.go:71`).
- `Import edges`: none.
- `Source-row clauses and occurrences`: the source-sentence-to-row table above lists every finding and its defect site.
- `Promised field labels`: `Covers:`, `Writes:`, `Occurrences:`, `Line:`, the `bench probe` cell `restored`, the role `integration-verification`, and the trigger `user-directed`.
- `Changed-function callers`: `checkRecurrenceMaintenanceContract` has two callers, `checkOccurrenceLedgerAndMaintenance` and `TestRecurrenceMaintenanceContractCheckBites`. `anchorsFixtureNeedles` has two callers, the two `bench anchors` tests that GR107 and GR108 name. No production function changes its signature.
- `Copy survival`: a Forbid row reds when a removed copy survives. The rows are GR15, GR16, GR29, GR30, GR41, GR42, GR43, GR44, and GR66. The other rows are GR67, GR81, GR83, GR84, GR86, GR87, GR88, GR101, GR102, GR114, GR122, GR123, and GR124.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"GR-A","tickets":["1-align-the-delegation-skill-with-fresh-authors.md","2-align-the-delegation-discipline-with-fresh-authors.md","3-route-phase-command-repairs-to-fresh-sessions.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow"},{"id":"conformance","command":"bench test --package ./internal/conformance --run TestRootConformance"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"agents","command":"bench test --check claude-agent-definitions"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"},{"id":"fixture-bite","command":"bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'"}]},{"id":"GR-B","tickets":["4-state-the-lane-and-landing-gate-once.md","5-capture-the-retro-by-the-tracked-rule.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow"},{"id":"conformance","command":"bench test --package ./internal/conformance --run TestRootConformance"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"help","command":"bench test --package ./cmd/bench --run TestHelpInventoryIsComplete"},{"id":"commit-help","command":"bench test --package ./internal/commit --run TestHelpAdvertisesDryRun"}]},{"id":"GR-C","tickets":["6-bind-each-ticket-author-to-the-declared-line.md","7-point-the-guides-at-each-fact-owner.md","8-correct-the-drain-and-craft-skill-references.md"],"verification":[{"id":"workflow","command":"bench test --check docs-currency-workflow"},{"id":"conformance","command":"bench test --package ./internal/conformance --run TestRootConformance"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"anchors-cli","command":"bench test --package ./cmd/bench --run TestAnchors"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"},{"id":"fixture-bite","command":"bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/workflow-guidance-repair/spec.md"},{"id":"workflow","command":"bench test --check docs-currency-workflow"},{"id":"conformance","command":"bench test --package ./internal/conformance --run TestRootConformance"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"budgets","command":"bench test --check guidance-prose-budgets"}]}
```
