# Land the spec and its deliverable binding at the spec-stage close

Status: staged

Roadmap: FT395

Decision source: named reviewed artifact `roadmap/FT395.md` at `0bdc11c5`, with the three FT395 criteria of `.bench/commitment.json` and the reviewer direction of 2026-10-11 in this phase's invocation

Audience: every repository that links the kit

Verification log: pending

## Problem

A spec-stage close lands its staged spec, but no step binds that spec as the deliverable of its outcome.
The Exit handoff of `.agents/commands/bench-write-spec.md` (lines 13-16) names no landing and no binding step.
`.bench/BENCH.md:144` lets a closing phase bind a deliverable only for a row whose `Next:` is `ticket`.
So the next session reads `bench commitment start --outcome <id> --request <request> --deliverable <path>` in `bench status`, and the start refuses.

The outlook fills `--deliverable` only from an approved binding in the default-branch policy (`internal/commitment/outlook.go:100`, `:109-121`, `internal/commitment/commitcmd/command.go:181-196`).
`Store.Start` refuses a deliverable that the policy does not approve (`internal/commitment/repository/admission.go:31-40`).

Each recent spec stage then needed a second landing and a second gate for the binding:

| spec | spec landing | binding landing | gap |
| --- | --- | --- | --- |
| spec-stage-grader-trace (FT376) | `ad410244` | `dad3721f` | 7m24s |
| complete-checkpoint-closure (FT392) | `a9506ce5` | `ea3e7867` | 33m32s |
| refusal-route-registry (FT393) | `9a98b054` | `9c228393` | 17m38s |
| build-directions | `a897b693` | `dcf7fad0` | 6m09s |
| commit-merge-parent (FT258-merge-parent) | `66af483b` | `342975e3` | 8m41s |

The second landing exists because each binding check reads the default-branch tip.
`Store.Plan` and `Store.Approve` validate every source at that tip (`internal/commitment/repository/sources.go:74-93`, `:105-118`, `internal/commitment/repository/repository.go:104`, `:111`, `:183`).
The approval check of `bench commit` and of `bench worktree land` does the same (`internal/commitment/repository/candidate.go:39-45`, `:160-188`).
A spec that only the spec-stage worktree holds therefore cannot bind until a landing puts it on main.

The handoff cannot show this state.
`bench handoff` writes one `Spec` and `Spec status` pair for each live spec and no binding fact (`internal/handoff/facts.go:83-100`).
The FT376 handoff called the build ready while the staged spec was off main and unbound.

A binding commit inside a build range also turns the build preflight red.
`paths-authorized` authorizes only the fence entries, `capture`, and the spec folder (`internal/preflight/paths_authorized.go:24-47`, `internal/preflight/fence_writes.go:71-77`).
No implicit entry covers `.bench/commitment.json`, so the FT393 binding commit turned the preflight red.

## Solution

A spec-stage close lands its staged spec and its deliverable binding in one spec-less landing.
The close derives the binding from the spec's `bench commitment inventory` row, which carries the spec path, its identity, and its `Roadmap:` value.
Plan and approve read a new deliverable at the worktree's committed HEAD. The approval check of the commit and the landing reads it in the candidate tree.
After the landing, `bench status` prints the complete `bench commitment start` command with the spec path.

`bench handoff` names the binding state of each live spec on one `Spec binding` line, read from the default-branch policy.
The spec-stage close passes the start command of `bench status` to the handoff as its next command.
Each spec's fence and one ticket's `Writes:` line own `.bench/commitment.json`, so the plan-only preflight stays green after the binding commit.

## User stories

Line: opus / xhigh
Implementation-line reason: SB1 is the hardest material chunk. It changes the approval check of the commit and landing owner and the revision rule of plan and approve. So it crosses owner seams, and real Git landings prove it. A repair round there costs over 100k tokens, so the scorecard routes its tickets to Opus at xhigh. SB2 is a small derived field, and SB3 is exact guidance, so their authors can lower effort in session under `craft-line`.
Harder chunks: SB1

### The one-landing binding

1. As a spec author, I want plan to bind a spec that my HEAD holds and main lacks, so that both land together.
2. As a spec author, I want `bench commitment approve` to stage that plan in my worktree, so that I commit the binding beside the spec.
3. As a spec author, I want `bench commit` to accept the approved binding in my worktree, so that the binding joins the reviewed source.
4. As a spec author, I want one spec-less landing to publish the spec and its binding, so that one gate closes the spec stage.
5. As a spec author, I want my worktree inventory to list the staged spec that HEAD holds, so that I need no hand digest.
6. As a spec author, I want that inventory row to carry the spec's `Roadmap:` value as its scope, so that the row names the obligation.
7. As a spec author, I want plan to bind a tickets-only folder that HEAD holds, so that a ticket folder binds under the same rule.
8. As the next session, I want `bench status` after the landing to print the start command with the spec path, so that I need no lookup.
9. As the next session, I want `bench commitment start` to admit the bound spec after the landing, so that the build starts on the first try.

### The refusals and the unchanged reads

10. As a reviewer, I want plan to refuse a staged spec that only the working tree holds, so that a binding names committed bytes.
11. As a reviewer, I want plan to refuse a new binding of a spec that HEAD lacks, so that each new binding reads one tree.
12. As a reviewer, I want plan to refuse a spec at HEAD that is not staged, so that only planning work binds.
13. As a reviewer, I want `bench commit` to refuse the binding after the spec changes, so that the binding names the reviewed bytes.
14. As a reviewer, I want the landing to refuse a composed tree that changes the bound spec, so that main holds no stale binding.
15. As a reviewer, I want a policy change without an exact approved plan to stay refused, so that the one-landing route needs the same approval.
16. As a reviewer, I want each roadmap-row source to keep the default-branch tip, so that a worktree edit of a row cannot pass as current.
17. As a spec author in an older worktree, I want plan to keep each binding on main, so that the rule adds no refusal.
18. As the reviewer, I want `bench status` to show `--deliverable <path>` until the binding lands, so that a worktree binding never reads as live.
19. As an agent, I want `bench commitment start` to refuse the spec before the binding lands, so that start reads only the published policy.
20. As the reviewer, I want the printed start command to keep the `<request>` placeholder, so that status never invents a request.
21. As an agent, I want an eligible outcome with no binding to keep its current status output, so that this delivery changes no route.

### The handoff binding state

22. As the next session, I want each live spec in a handoff section to carry one `Spec binding` line, so that I see its binding.
23. As the next session, I want that line to name the outcome that main binds to the spec, so that I know which outcome starts.
24. As the next session, I want that line to read `unbound` when main binds no outcome to the spec, so that I see the gap.
25. As the next session, I want a binding that only the worktree holds to read `unbound`, so that the handoff states the truth of main.
26. As the next session, I want two binding outcomes to appear in policy order, so that no binding goes missing.
27. As an agent with no policy, I want the line to read `unbound`, so that a kit without a commitment still writes its handoff.
28. As an agent, I want an unreadable policy to give `binding unknown (policy unreadable)` and a written handoff, so that a commitment fault never blocks the handoff.
29. As a maintainer, I want the Shape text to name the `Spec binding` label, so that the Shape describes the section grammar.
30. As a maintainer, I want the `Spec` and `Spec status` lines to stay adjacent, so that the readers of the pair keep their bytes.

### The close guidance

31. As a spec author, I want the Exit handoff to land the spec and its binding together, so that both reach main first.
32. As a spec author, I want the Exit handoff to bind from the inventory row after the last spec edit, so that the identity holds.
33. As a spec author, I want the Exit handoff to pass the start command to `bench handoff --next`, so that the handoff pins that command.
34. As the reviewer, I want `.bench/BENCH.md` to let a spec-stage close bind its spec on my sign-off, so that nothing else rides on it.
35. As a spec author, I want the slicing checks to put `.bench/commitment.json` in the fence and one `Writes:` line, so that the preflight stays green.
36. As the reviewer, I want this spec to obey that rule, so that its own preflight shows the rule.
37. As a teammate, I want `AGENTS.md` to pin the binding state in a phase-close handoff, so that a cold reader checks it.
38. As a teammate, I want `AGENTS.md` to run the spec-stage handoff from the primary checkout after its landing, so that the handoff reads main.
39. As a spec author, I want the template's `Roadmap:` comment to name the binding it feeds, so that the derivation has its input.
40. As the reviewer, I want each guidance edit to fit its current prose budget, so that no budget rises.
41. As a maintainer, I want each new pinned sentence to carry an anchor and a recorded mutation red, so that the gate keeps its bite.
42. As a reviewer, I want the Exit handoff to keep the planning-only staged-spec sentence, so that delivery still waits for the start.

## Implementation decisions

### Terms

- **new deliverable source**: a source binding that a deliverable of the proposed policy names. No deliverable of the current default-branch policy names that id, path, and identity.
- **checkout HEAD**: the commit that `git rev-parse HEAD` names in the store root.
- **candidate tree**: the tree that `authorizeCandidate` grades. It is the composed candidate of `bench commit`, or the composed tree of a landing.
- **binding state**: the ids of the outcomes whose undelivered deliverable binding in the default-branch policy names a spec path.

The spec uses these terms only in its own text. It adds no CONTEXT.md entry.

### Closed decisions

The FT395 criteria (`.bench/commitment.json:479-493`) state the outcome:

- `.close-lands`: "A spec-stage close lands its staged spec and its deliverable binding on main before the handoff, and the handoff names the binding state."
- `.status-command`: "bench status prints the complete bench commitment start command, with the deliverable that the spec's Roadmap line and path derive."
- `.preflight-green`: "A ticket Writes line owns .bench/commitment.json, so the build preflight stays green after the binding."

The reviewer directed on 2026-10-11 that `bench commitment start` follows the spec stage with the bound deliverable, and that the handoff pins that command.
On the same date, the reviewer decided that the close passes that command to `bench handoff --next` (E7).

The coordinator locked the engineering decisions below. Each one is open to reviewer veto at sign-off.

- D1 One landing. A new deliverable source binds from the tree that publishes it. Reason: the 2026-10-10 occurrence names the second landing and gate as the defect (`roadmap/FT395.md:18`). `protectedCandidate` already validates each unsettled roadmap-row source at the candidate tree (`internal/commitment/repository/candidate.go:240-244`).
- D2 Derivation. The close derives the binding from the inventory row of the staged spec. No new verb or flag joins the commitment grammar. Reason: `Store.Inventory` reads the status, the `Roadmap:` value, and the identity of each spec from the same bytes (`internal/commitment/repository/sources.go`, `stagedSpecs`).
- D3 Complete command. "Complete" means that the deliverable path is filled. The `--request <request>` placeholder stays. Reason: `withCommand` keeps it by design, because only the worker names a request (`internal/commitment/commitcmd/command.go:181-182`), and no request exists at status time.
- D4 Handoff binding state. `bench handoff` writes one derived `Spec binding` line after each `Spec status` line. Reason: a separate label keeps one fact per field, and no code reads the `Spec status` value (`internal/handoff/facts.go:83-100`).
- D5 Writes rule. Every spec's fence and exactly one ticket's `Writes:` line own `.bench/commitment.json`. Reason: the close commits the binding inside the spec-stage range, and the plan-only preflight grades that range (`internal/preflight/paths_authorized.go:24-47`). No preflight row requires a `Writes:` path to change: `writes-resolve` passes when the path exists, and no closure row pins `.bench/commitment.json` (`internal/preflight/decision.go:201-222`, `internal/preflight/ticket_rows.go:13-27`, `internal/anchors/registry_ticket_passes.go:64-65`). `internal/preflight/decision_test.go:179-183` pins the `paths-authorized` red of an unfenced path.
- D6 Guidance. The write-spec Exit handoff owns the close order. `.bench/BENCH.md` owns the binding authority, and `AGENTS.md` owns the pin list of a phase-close handoff. Reason: one source per fact.
- D7 An eligible outcome with no approved deliverable keeps its current status output. Reason: the decision source states the status output after the binding, not a new route.
- D8 Bootstrap. This spec's own close binds through the current two-landing route, because the one-landing route ships with this spec.
- D9 Line: `opus / xhigh`, for the reason in the user stories.

### The revision rule of a deliverable

`Store.Plan` validates each new deliverable source at the checkout HEAD.
`Store.Approve` and its replay validate each new deliverable source at the checkout HEAD.
`Store.Inventory` lists each staged spec that the checkout HEAD holds, and it reads each identity at that HEAD.
The approval check `approvedTransition` validates each new deliverable source of the plan at the candidate tree.

Every other source keeps the default-branch tip.
That rule covers each roadmap-row source and each deliverable that the current policy already binds.
`Store.Start` and `Store.readyFor` keep the default-branch tip for every deliverable (`internal/commitment/repository/admission.go:41`, `internal/commitment/repository/readiness.go:60`).
The outlook keeps the default-branch policy (`internal/commitment/repository/outlook.go:15-32`).
A binding therefore counts for start, readiness, and status only after it lands.

A new deliverable source that HEAD does not hold refuses, even when main holds it.
A worktree that predates the spec folds main first, through `bench worktree merge --from main`.
The refusal keeps the current sentence `commitment source "<id>" refused:`.

The HEAD read reuses `sourceAt` and the staged check of `validateDeliverable`.
The tickets-only rule of `ticketsOnlyAt` applies at the same HEAD.
`validateDeliverable` keeps its default-tip behavior for its two callers outside plan and approve.
One owner decides which tree each source reads, so plan, approve, and the approval check cannot drift.
The inventory reads the same checkout HEAD that plan reads for a new deliverable.

`.bench/commitment.json` is a planning path (`internal/commitment/model.go:112-113`), so a commit that changes only the policy ends at `internal/commitment/repository/candidate.go:99-100` with no light-path grading.
Each later commit in the same worktree grades the transition again, because `branchReference` returns the tip while the base policy differs (`internal/commitment/repository/candidate.go:131-144`).
SB13 therefore holds for the whole life of the worktree, not only for its first commit.

### The approval check

`approvedTransition` receives the candidate tree beside the default-branch revision.
For each source of the expected plan, it validates a new deliverable source at the candidate tree and every other source at the revision.
The receipt match, the plan identity, and the refusal `candidate policy has no exact approval` stay unchanged.

`bench commit` grades the composed candidate of the named paths, so an edit of the spec after approve fails the identity check at commit.
`bench worktree land` grades the composed tree, so a landing whose composition changes the spec refuses before publication.
A spec-less landing carries no delivery, so it takes the transition path and never the closure path (`internal/commitment/repository/publication.go:39-52`).

### The printed command

The outlook, `withCommand`, and the status row stay unchanged.
After the landing, `soleDeliverable` returns the bound spec path, so the command cell carries it.
The `--request <request>` placeholder stays, as `cmd/bench/commitment_test.go` (`TestCommitmentNextProjection`) pins it.
Before the landing, the command keeps `--deliverable <path>`.
The outlook comment at `internal/commitment/outlook.go:37-40` says that a reader never selects work that a start refuses. An eligible outcome with no binding contradicts that comment, and Further notes records it.

### The handoff binding state

`pins` writes one `Spec binding` line directly after the `Spec status` line of each live spec.
`handoffdoc` gains the label constant `LabelSpecBinding` with the value `Spec binding`.
The parser already reads any label line as a generic field (`internal/handoffdoc/document.go:381`), so the parser takes no edit beside the constant.
The value lists the id of each outcome whose undelivered deliverable binding in the default-branch policy names the spec path.
The ids follow policy order, with one space between two ids.

The value is `unbound` when no outcome binds the path.
The value is `unbound` when the repository has no policy, or when no default branch resolves, as the outlook reads that state (`internal/commitment/repository/outlook.go:17-19`).
The value is `binding unknown (policy unreadable)` when the policy does not read or does not parse.
That phrase joins the unknown phrases of `internal/handoff/facts.go:44-53`, and the handoff still writes.

The handoff reads the policy once per run through `commitrepo.Store.Policy`.
`errDefaultUnresolved` is private (`internal/commitment/repository/outlook.go:10-19`), so the repository package exports one test for an unresolved default branch, and the outlook and the handoff both call it.
`internal/commitment` owns the rule that maps a path to its binding outcomes, beside `Unsettled`, so the handoff restates no settlement rule.
The new import edges are `internal/handoff` to `internal/commitment` and to `internal/commitment/repository`.
`internal/handoff` already depends on both packages through `internal/status`, and neither package depends on `internal/handoff`.

The Shape text of `internal/handoff/text.go` names the new label beside `Spec` and `Spec status`.
`TestHandoffShapeNamesTheSectionGrammar` gains the label in its token list.

### The close order

The Exit handoff of `.agents/commands/bench-write-spec.md` states this order after sign-off:

1. Write the verification log into the spec and the scorecard updates, and commit them.
2. After the last spec edit, copy the spec's `bench commitment inventory` row into one new deliverable binding of the current policy. Put it on the outcome whose sources hold the row's scope, with that scope as the obligation. A spec with no `Roadmap:` line binds to the outcome that `bench status` names.
3. Run `bench commitment plan` and `bench commitment approve` with the sign-off as the decision, and commit the binding.
4. Run `bench preflight build <slug> --plan-only` at that tip.
5. Land spec-less once.
6. From the primary checkout, run `bench handoff --next` with the start command that `bench status` prints.

Plan's obligation rule still guards the binding of a spec with no `Roadmap:` line (`internal/commitment/authority.go:120-135`).
The binding comes after the last spec edit, because the identity is the digest of the spec bytes (`internal/commitment/repository/sources.go:24-40`).
The Exit handoff keeps the planning-only staged-spec sentence that `internal/anchors/registry_commitment.go:42` pins.

`.bench/BENCH.md:144` gains the spec-stage grant before its binding-only limit sentence.
The proposed text is: "A spec-stage close binds its approved spec the same way and cites my sign-off as the decision."
The pinned sentences at `.bench/BENCH.md:142` and `:144` stay byte-identical.

`AGENTS.md` changes two sentences of its Phase-close handoff.
The pin list adds the binding state beside the spec path and status.
The venue sentence at `AGENTS.md:45` gains the spec-stage exception: that close runs `bench handoff` from the primary checkout after its landing releases the worktree.
The venue joins the existing sentence, because that paragraph already holds six sentences with its bold label (`internal/prose/parse.go:298`).

The command's close order also names that venue, because the command ships to every linked repository and `AGENTS.md` does not.
The pinned sentence "A drain runs it from the primary checkout after its landing." stays.

The slicing checks of `craft-tickets` gain one slicing rule.
The proposed text is: "The ownership fence and exactly one ticket's `Writes:` line name `.bench/commitment.json`. The spec-stage close commits the binding inside the range that the plan-only preflight grades."

The `Roadmap:` comment of the `craft-spec` template states that the close binds the spec to the outcome whose sources hold that row.
The comment stays on its one line, so the skill keeps its line count.

### The prose budgets

`.agents/commands/bench-write-spec.md` holds 87 lines against its 87-line budget, and `.bench/BENCH.md` holds 181 against 192 (`projects/benchkit.md:488-501`).
`.agents/skills/bench-craft-spec/SKILL.md` holds 155 against 155.
No budget rises, because a raised budget weakens a check.
`guidance-prose-budgets` grades the result of each edit.

The Exit handoff paragraph of the command already holds five sentences, and a paragraph holds at most six (`internal/prose/prose.go:17-18`).
The close order therefore takes a new block after that paragraph.
To pay for that block, the edit joins wrapped prose lines into one physical line each, and it changes no word.

Anchors collapse whitespace (`internal/anchors/locate.go:177`), so a join keeps every anchor.
A canary `MUTATE.json` `old` string or a Go mirror string is byte-exact (`internal/canary/mutation.go:68-78`).
Such strings hold the newlines at command lines 21-22, 32-35, 50-51, 61-62, 68-69, 79-81, and 86-87, so the edit joins none of them.
The joins at command lines 9-11, 20-21, 22-23, 28-29, 51-52, and 63-65 free eight lines, and the block uses at most eight.

The edit keeps the two `Require` sentences that `internal/anchors/registry_data.go:336-337` pin on command line 16.
The edit adds no needle that `internal/anchors/registry_data.go:335`, `:358-359`, or `internal/anchors/registry_ft311_preparation.go:50` forbid.
The `.bench/BENCH.md` grant extends the paragraph of line 144 from three sentences to four, so it adds no line.

### The anchors

The new Exit handoff sentences and the new `.bench/BENCH.md` sentence join the commitment anchor family in `internal/anchors/registry_commitment.go`.
Each one takes one mutation mirror in `TestCommitmentGuidance`, which records its red.
The new slicing rule joins the slicing-checks anchors in `internal/anchors/registry_ticket_passes.go`, with its recorded red in `TestTicketSlicingPasses`.
Each new sentence avoids the forbidden needles of the write-spec command, such as "## Handoff", "map's Handoff", and "bench preflight evidence".

### The unchanged surfaces

The commitment grammar, its help rows, and its refusal faces stay unchanged.
The refusal registry gains no face.
`internal/commit/commit.go` stays unchanged.
The roadmap row and the commitment criteria stay unchanged, and the completion landing settles FT395 as it settles any outcome.

### The changelog entry

One typed section under `## [Unreleased]` of `CHANGELOG.md` states the one-landing close and the `Spec binding` line (`.agents/skills/bench-craft-synthesis/SKILL.md:74`).
The last ticket writes it, because that ticket follows both behaviors.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| SB1 / 1-bind-a-new-deliverable-from-its-publishing-tree.md | A new deliverable binds from the tree that publishes it, so one spec-less landing publishes a staged spec and its binding, and `bench status` then prints the complete start command | SB01-SB22 | the `internal/commitment`, `internal/commitment/repository`, `internal/commitment/commitmenttest`, `internal/commit`, `internal/worktree`, `internal/anchors`, `cmd/bench`, and `internal/conformance` package tests | yes |
| SB2 / planned: 2-name-the-binding-state-in-the-handoff.md | `bench handoff` names the binding state of each live spec | SB23-SB30, SB44 | the `internal/handoff`, `internal/handoffdoc`, `internal/commitment`, `internal/commitment/repository`, `internal/commitment/commitmenttest`, `internal/anchors`, `cmd/bench`, and `internal/conformance` package tests | no |
| SB3 / planned: 3-state-the-one-landing-close.md | The guidance states the one-landing close, the binding grant, the handoff pin, the Writes rule, and the template comment | SB31-SB43 | the `internal/anchors`, `cmd/bench`, and `internal/conformance` package tests, and the `docs-currency-workflow`, `guidance-prose-budgets`, `prose-mechanics`, and `skills-index-command-adapters` checks | no |

SB1 and SB2 both write files of `internal/commitment`, so the serial order also keeps their shared helper files apart.
SB3 describes the shipped behavior of both, so it waits for the SB1 and SB2 checkpoints.
The serial order is SB1, SB2, and then SB3.
The SB1 ticket owns `.bench/commitment.json` in its `Writes:` line under D5.
SB17-SB19 add only a journey test over the seam that SB12-SB15 open, so they are acceptance rows of the same ticket (`.agents/skills/bench-craft-tickets/SKILL.md:11-14`).
The orchestrator reconciles the final acceptance after the SB3 checkpoint.

### Completion plan

This version-1 completion plan is a planning verification inventory, not a dispatch plan.
Its commands are future obligations, not executed proof.
The ticket names are planned, and the slicing replaces them with the sliced names.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "SB1",
      "tickets": ["1-bind-a-new-deliverable-from-its-publishing-tree.md"],
      "verification": [
        {"id": "commitment-packages", "command": "bench test --package ./internal/commitment --package ./internal/commitment/repository --package ./internal/commitment/commitcmd --package ./internal/commitment/commitmenttest"},
        {"id": "commit-package", "command": "bench test --package ./internal/commit"},
        {"id": "worktree-package", "command": "bench test --package ./internal/worktree"},
        {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
        {"id": "command-package", "command": "bench test --package ./cmd/bench"},
        {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
        {"id": "root-conformance", "command": "bench test --check root-conformance"},
        {"id": "git-plumbing-owner", "command": "bench test --check git-plumbing-owner"},
        {"id": "ordinary-build-census", "command": "bench test --check ordinary-build-census"}
      ]
    },
    {
      "id": "SB2",
      "tickets": ["2-name-the-binding-state-in-the-handoff.md"],
      "verification": [
        {"id": "handoff-packages", "command": "bench test --package ./internal/handoff --package ./internal/handoffdoc"},
        {"id": "commitment-packages", "command": "bench test --package ./internal/commitment --package ./internal/commitment/repository --package ./internal/commitment/commitcmd --package ./internal/commitment/commitmenttest"},
        {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
        {"id": "command-package", "command": "bench test --package ./cmd/bench"},
        {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
        {"id": "root-conformance", "command": "bench test --check root-conformance"}
      ]
    },
    {
      "id": "SB3",
      "tickets": ["3-state-the-one-landing-close.md"],
      "verification": [
        {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
        {"id": "command-package", "command": "bench test --package ./cmd/bench"},
        {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
        {"id": "guidance", "command": "bench test --check docs-currency-workflow"},
        {"id": "budgets", "command": "bench test --check guidance-prose-budgets"},
        {"id": "prose", "command": "bench test --check prose-mechanics"},
        {"id": "skills-index", "command": "bench test --check skills-index-command-adapters"},
        {"id": "root-conformance", "command": "bench test --check root-conformance"}
      ]
    }
  ],
  "final_verification": [
    {"id": "coverage", "command": "bench coverage --check specs/spec-stage-binding/spec.md"},
    {"id": "commitment-packages", "command": "bench test --package ./internal/commitment --package ./internal/commitment/repository --package ./internal/commitment/commitcmd --package ./internal/commitment/commitmenttest"},
    {"id": "commit-package", "command": "bench test --package ./internal/commit"},
    {"id": "worktree-package", "command": "bench test --package ./internal/worktree"},
    {"id": "handoff-packages", "command": "bench test --package ./internal/handoff --package ./internal/handoffdoc"},
    {"id": "anchors-package", "command": "bench test --package ./internal/anchors"},
    {"id": "command-package", "command": "bench test --package ./cmd/bench"},
    {"id": "conformance-package", "command": "bench test --package ./internal/conformance"},
    {"id": "guidance", "command": "bench test --check docs-currency-workflow"},
    {"id": "budgets", "command": "bench test --check guidance-prose-budgets"},
    {"id": "prose", "command": "bench test --check prose-mechanics"},
    {"id": "skills-index", "command": "bench test --check skills-index-command-adapters"},
    {"id": "root-conformance", "command": "bench test --check root-conformance"},
    {"id": "git-plumbing-owner", "command": "bench test --check git-plumbing-owner"},
    {"id": "ordinary-build-census", "command": "bench test --check ordinary-build-census"}
  ]
}
```

## Testing decisions

A good test binds through real Git and reads the result through the public owner.
It commits a staged spec in a linked worktree, runs plan and approve there, and grades the candidate or the composed tree.
A refusal row reads the error sentence and an unchanged ledger, policy, or ref.
No expectation calls the code under test: each identity comes from `commitment.Identity` over the fixture bytes, and each path comes from the fixture writes.

These seams receive tests:

- `internal/commitment/head_binding_test.go` (new): `Store.Plan`, `Store.Approve`, `Store.Inventory`, and `Store.Start` over a `commitmenttest.Repo` with an assignment worktree from `commitmenttest.Assignment`. `TestCommitmentInventoryIdentity` and `TestCommitmentDeliverableTypes` are the precedents.
- `internal/commitment/repository/head_binding_test.go` (new): `Store.AuthorizeCandidate` and `Store.AdmitPublication` over a worktree candidate and a hand-built composed tree. `TestAdmitPublicationFrozenIdentity` and `TestStaleBaseCandidate` are the precedents.
- `internal/worktree/land_specless_test.go`: one journey that commits a staged spec and its approved binding in a planning worktree, lands it spec-less, and reads `bench commitment show` before and after. `planningLandingFixture` with an adopted policy and `TestCommitmentPlanningBootstrap` are the precedents.
- `internal/handoff/render_test.go`: the `Spec binding` line over a policy on main. `TestSectionRendersTheAssignmentPins` and `TestSectionRendersOnePairPerLiveSpec` are the precedents.
- `internal/conformance/commitment_guidance_test.go` and `internal/anchors/registry_ticket_passes_test.go`: one mutation mirror for each new pinned sentence.

No system test joins the build, because the package journeys cross the public seams with real Git.
The gate observes the feature through the package suites and the guidance checks.

### Seam diagram

    trigger: the spec-stage close in a planning worktree
        │
        ▼
    staged spec at HEAD ──▶ [ Store.Inventory ── Store.Plan ── Store.Approve: new deliverable at HEAD ]
        │
        ▼
    bench commit ──▶ [ authorizeCandidate ── approvedTransition: new deliverable at the candidate tree ]
        │
        ▼
    bench worktree land (spec-less) ──▶ [ AdmitPublication ── approvedTransition at the composed tree ]
        │
        ▼
    main holds the spec and its binding ──▶ [ Outlook ── withCommand ] ──▶ the complete start command
        │
        ▼
    bench handoff ──▶ [ pins ── Spec binding from the default-branch policy ]
        ◀ tests attach at Store.Plan, Approve, Inventory, and Start in a linked worktree, at
          AuthorizeCandidate and AdmitPublication with real trees, at the spec-less landing
          journey, and at the handoff section render

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SB01 | 1 | In a linked worktree whose HEAD holds a staged spec that main lacks, `Store.Plan` accepts a policy that binds that spec | planned TestHeadBindingPlansAStagedSpecThatHEADHolds in internal/commitment/head_binding_test.go | The current plan reads the deliverable at the default-branch tip, so it refuses with `commitment source ... refused`. |
| SB02 | 2 | `Store.Approve` of that plan writes the bound policy into the worktree's `.bench/commitment.json` | planned TestHeadBindingPlansAStagedSpecThatHEADHolds in internal/commitment/head_binding_test.go | The current approve validates the plan sources at the default-branch tip, so it refuses the plan. |
| SB03 | 5 | In that worktree, `Store.Inventory` lists a `deliverable` row for the HEAD-held spec | planned TestHeadBindingInventoryListsAHeadHeldSpec in internal/commitment/head_binding_test.go | The current inventory skips a spec that the default-branch revision does not hold. |
| SB04 | 6 | The inventory row of the HEAD-held spec carries the spec's `Roadmap:` value as its scope | planned TestHeadBindingInventoryListsAHeadHeldSpec in internal/commitment/head_binding_test.go | A HEAD row that skips the metadata read of its own bytes leaves the scope empty. |
| SB05 | 5 | The inventory identity of the HEAD-held spec equals the identity that `Store.Plan` accepts for that spec | planned TestHeadBindingInventoryListsAHeadHeldSpec in internal/commitment/head_binding_test.go | An inventory that reads another revision than plan gives a digest that plan refuses. |
| SB06 | 7 | In a linked worktree whose HEAD holds a tickets-only folder that main lacks, `Store.Plan` accepts a policy that binds that folder | planned TestHeadBindingPlansATicketsOnlyFolder in internal/commitment/head_binding_test.go | A HEAD rule that skips `ticketsOnlyAt` refuses the folder as a source file. |
| SB07 | 10 | `Store.Plan` in the worktree refuses a staged spec that only the working tree holds | planned TestHeadBindingRefusesAnUnheldDeliverable in internal/commitment/head_binding_test.go | A plan that reads the working tree binds bytes that no commit holds. |
| SB08 | 11 | `Store.Plan` in a worktree whose HEAD lacks a spec that main holds refuses a new binding of that spec | planned TestHeadBindingRefusesAnUnheldDeliverable in internal/commitment/head_binding_test.go | A fallback to the default-branch tip makes one new binding read two trees. |
| SB09 | 12 | `Store.Plan` in the worktree refuses a HEAD-held spec that carries no `Status: staged` line, with the sentence `deliverable "<path>" is not staged` | planned TestHeadBindingRefusesAnUnheldDeliverable in internal/commitment/head_binding_test.go | A HEAD read that skips the staged check binds a draft or an implemented spec. |
| SB10 | 16 | `Store.Plan` in a worktree that commits a changed roadmap row refuses a plan that names the changed row identity | planned TestHeadBindingKeepsTheDefaultTipForEveryOtherSource in internal/commitment/head_binding_test.go | A rule that moves every source to HEAD lets a worktree row edit pass as current. |
| SB11 | 17 | `Store.Plan` in a worktree whose HEAD predates another outcome's binding on main accepts a policy that keeps that binding | planned TestHeadBindingKeepsTheDefaultTipForEveryOtherSource in internal/commitment/head_binding_test.go | A rule that reads every deliverable at HEAD refuses the retained binding that HEAD lacks. |
| SB12 | 3 | In the worktree, `Store.AuthorizeCandidate` admits a candidate tree that holds the staged spec and its approved binding | planned TestAuthorizeCandidateAdmitsAHeadHeldBinding in internal/commitment/repository/head_binding_test.go | The current approval check validates the plan sources at the default-branch tip, so it refuses. |
| SB13 | 13 | `Store.AuthorizeCandidate` refuses a candidate tree whose spec bytes differ from the approved identity while HEAD keeps the approved bytes | planned TestAuthorizeCandidateRefusesAChangedBoundSpec in internal/commitment/repository/head_binding_test.go | A check that reads the deliverable at HEAD and not at the candidate tree admits an edit named after approve. |
| SB14 | 4 | `Store.AdmitPublication` admits a spec-less publication whose composed tree holds the staged spec and its approved binding | planned TestAdmitPublicationGradesANewBindingAtTheComposedTree in internal/commitment/repository/head_binding_test.go | The current landing validates the plan sources at the default-branch tip, so it refuses. |
| SB15 | 14 | `Store.AdmitPublication` refuses a composed tree whose spec bytes differ from the approved identity | planned TestAdmitPublicationGradesANewBindingAtTheComposedTree in internal/commitment/repository/head_binding_test.go | A landing that reads the source tip and not the composed tree publishes a binding with a stale identity. |
| SB16 | 15 | `bench commit` refuses a policy candidate whose transition has no exact approved plan | `internal/commit/commitment_test.go` (`TestCommitmentProposalIdentity`) | A one-landing change that skips the receipt match fails this pinned `exact approval` refusal. |
| SB17 | 4 | A spec-less `bench worktree land` of the planning worktree publishes the staged spec and its binding in one landing | planned TestCommitmentSpecStageLandsSpecAndBinding in internal/worktree/land_specless_test.go | The current landing refuses the binding, because main lacks the spec. |
| SB18 | 8 | After that landing, `bench commitment show` in the primary checkout prints the start command with `--deliverable specs/<slug>/spec.md` | planned TestCommitmentSpecStageLandsSpecAndBinding in internal/worktree/land_specless_test.go | Without the binding on main, the command keeps `--deliverable <path>`. |
| SB19 | 18, 21 | Before that landing, `bench commitment show` in the planning worktree prints the start command with `--deliverable <path>` | planned TestCommitmentSpecStageLandsSpecAndBinding in internal/worktree/land_specless_test.go | The planning worktree's policy file holds the binding, so an outlook that reads that file or HEAD shows a binding that main lacks. |
| SB20 | 9 | `bench commitment start` admits a staged spec that the default-branch policy binds | `internal/commitment/deliverable_test.go` (`TestCommitmentDeliverableTypes`) | A HEAD rule that breaks the shared default-tip read of `validateDeliverable` fails this pinned admission. SB21 owns a start read that moves to HEAD. |
| SB21 | 19 | Before the landing, `Store.Start` from the planning worktree refuses the HEAD-bound spec with `is not approved for outcome` | planned TestHeadBindingStartRefusesBeforeTheLanding in internal/commitment/head_binding_test.go | A start that reads the worktree policy or HEAD admits a binding that main lacks. |
| SB22 | 20 | The bound start command keeps the `--request <request>` placeholder | `cmd/bench/commitment_test.go` (`TestCommitmentNextProjection`) | A change that fills the request from the checkout fails this pinned command. |
| SB23 | 22, 23 | With a binding on main, `bench handoff` writes `Spec binding: <outcome>` on the line after the spec's `Spec status` line | planned TestSectionRendersTheSpecBinding in internal/handoff/render_test.go | The current pins write no binding line. |
| SB24 | 24 | `bench handoff` writes `Spec binding: unbound` for a live spec that the default-branch policy does not bind | planned TestSectionRendersTheSpecBinding in internal/handoff/render_test.go | A writer that omits the unbound state lets a reader read silence as a binding. |
| SB25 | 25 | `bench handoff` writes `Spec binding: unbound` when only the worktree's policy binds the spec | planned TestSectionRendersTheSpecBinding in internal/handoff/render_test.go | A reader of the worktree policy reports a binding that main lacks, as the FT376 handoff did. |
| SB26 | 26 | `bench handoff` writes both outcome ids in policy order with one space between them when two outcomes bind the spec | planned TestSectionRendersTheSpecBinding in internal/handoff/render_test.go | A first-match reader drops the second binding. |
| SB27 | 27 | `bench handoff` writes `Spec binding: unbound` in a repository whose default branch holds no policy | planned TestSectionRendersAnUnavailableBinding in internal/handoff/render_test.go | A reader that treats an absent policy as a fault refuses the handoff of a kit without a commitment. |
| SB28 | 28 | With an unparseable policy on main, `bench handoff` writes `Spec binding: binding unknown (policy unreadable)` | planned TestSectionRendersAnUnavailableBinding in internal/handoff/render_test.go | A reader that returns the policy error refuses the whole handoff. |
| SB44 | 27 | In a primary checkout whose default branch does not resolve, `bench handoff` writes `Spec binding: unbound` | planned TestSectionRendersAnUnavailableBinding in internal/handoff/render_test.go | A reader that treats the unresolved default branch as a policy fault writes `binding unknown (policy unreadable)` or refuses the handoff. |
| SB29 | 29 | The Shape text names the `Spec binding` label | `internal/conformance/handoff_single_source_test.go` (`TestHandoffShapeNamesTheSectionGrammar`) | Once the token list holds the label, a Shape that omits it reds this test. |
| SB30 | 30 | Each live spec keeps its `Spec` line directly before its `Spec status` line | `internal/handoff/render_test.go` (`TestSectionRendersOnePairPerLiveSpec`) | A binding line placed between the pair breaks this pinned pair. |
| SB31 | 31 | The write-spec Exit handoff states that one spec-less landing carries the staged spec and its binding | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | The new mirror removes the sentence and requires its anchor diagnostic. |
| SB32 | 32 | The write-spec Exit handoff binds from the spec's `bench commitment inventory` row after the last spec edit | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | The new mirror removes the sentence and requires its anchor diagnostic. |
| SB33 | 33 | The write-spec Exit handoff passes the start command that `bench status` prints to `bench handoff --next` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | The new mirror removes the sentence and requires its anchor diagnostic. |
| SB34 | 34 | `.bench/BENCH.md` lets a spec-stage close bind its approved spec with the sign-off as the decision | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | The new mirror removes the sentence and requires its anchor diagnostic. |
| SB35 | 34 | `.bench/BENCH.md` keeps the binding-only limit sentence after the new grant | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | The existing `binding-only limit` mirror reds when an edit drops or rewords that sentence. |
| SB36 | 35 | The slicing checks require the fence and exactly one ticket `Writes:` line to name `.bench/commitment.json` | `internal/anchors/registry_ticket_passes_test.go` (`TestTicketSlicingPasses`) | The new slicing anchor reds when the rule leaves the file. |
| SB37 | 36 | The ownership fence of this spec and the `Writes:` line of its first ticket name `.bench/commitment.json` | review-owned: the ticket review reads the fence, and the plan-only preflight `fence-writes` row grades the union | A fence without the path lets the close's binding commit red `paths-authorized`. |
| SB38 | 37 | `AGENTS.md` requires the phase-close handoff to pin the binding state | review-owned: the SB3 review reads the Phase-close handoff paragraph | No anchor pins that paragraph, so only review sees a dropped pin. |
| SB39 | 38 | `AGENTS.md` runs the spec-stage handoff from the primary checkout after its landing | review-owned: the SB3 review reads the Phase-close handoff paragraph | No anchor pins that paragraph, so only review sees the missing venue. |
| SB40 | 39 | The `Roadmap:` comment of the `craft-spec` template states that the close binds the spec to the outcome whose sources hold that row | review-owned: the SB3 review reads the template line | No anchor pins the template comment, so only review sees a stale comment. |
| SB41 | 40 | Each edited guidance file stays within its `guidance-prose-budgets` limit | `internal/conformance/prose_budget_test.go` (`TestGuidanceProseBudgetsHoldOnTheLiveTree`) | A line added to a full file reds this live-tree check. |
| SB42 | 41 | Each new pinned sentence has one anchor row and one mutation mirror that logs its observed red | review-owned: the SB3 review reads the anchor rows and the red log of the author | A sentence without a mirror keeps an anchor whose bite no run proved. |
| SB43 | 42 | The write-spec Exit handoff keeps the planning-only staged-spec sentence | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | The existing `staged spec` mirror reds when an edit drops or rewords that sentence. |

### Edge inventory

The profile's hostile-input checklist applies to the spec path, the `Roadmap:` value, and the handoff line.
The HEAD read reuses `sourceAt` and `git.ReadTreeFile`, so it serves each path that the default-tip read serves today.
The `Spec binding` value holds only outcome ids, and the policy parser limits an id to `^[A-Za-z][A-Za-z0-9._-]*$` (`internal/commitment/parse.go:224-229`).
A control byte or a newline therefore cannot reach the single-line sink.
SB27 and SB28 cover the absent-versus-empty pair: an absent policy reads `unbound`, and an empty policy file fails the parse.

No test swaps a package variable.
The refusal rows cover a fault at plan (SB07-SB09), at commit (SB13), and at the landing (SB15).

- **Won't handle**: a spec path with a glob character — `git.ReadTreeFile` lists the path as a pathspec, as the default-tip read of the start does today.
- **Won't handle**: a control byte in a `Roadmap:` value — the inventory TOON sink refuses it today, and the HEAD-held row shares that sink.
- **Won't handle**: a primary checkout whose HEAD is not the default-branch tip — `bench commit` refuses the primary checkout, so no binding from there reaches a landing.
- **Won't handle**: an unborn HEAD — `Store.Plan` already refuses a repository whose default branch does not resolve.
- **Won't handle**: a moved policy on main between plan and landing — `Store.Approve` refuses a changed predecessor, and `TestCommitmentStaleApproval` pins that refusal.
- **Won't handle**: a spec stage for an outcome that is not next — its `Spec binding` line names the outcome, and `bench status` names the next one.
- **Won't handle**: a ticket author who edits `.bench/commitment.json` under its `Writes:` line — `bench commit` refuses a policy change without an exact approved plan (SB16).
- **Won't handle**: a live spec whose binding a delivery already settled — a delivery flips its spec to implemented, so that spec is no longer live.
- **Won't handle**: a status route for an eligible outcome with no binding — D7 keeps the current output, and Out of scope prices the route.

## Ownership fences

This fence is the union of the three ticket Writes lines, which the plan-only preflight grades.
The coordinator owns the review pickup, which the coverage check requires in a declared fence section.
Ticket 1 alone names the commitment policy file (D5).

- `reviews/spec-stage-binding.md`
- `.agents/commands/bench-write-spec.md`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/skills/bench-craft-tickets/references/slicing-checks.md`
- `.bench/BENCH.md`
- `.bench/commitment.json`
- `AGENTS.md`
- `CHANGELOG.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/testdata/anchors/fixture-repo/AGENTS.md`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_craft_research.go`
- `internal/anchors/registry_craft_research_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/commitment/commitcmd/command.go`
- `internal/commitment/commitmenttest/repo.go`
- `internal/commitment/delivery.go`
- `internal/commitment/delivery_test.go`
- `internal/commitment/head_binding_test.go`
- `internal/commitment/repository/admission.go`
- `internal/commitment/repository/candidate.go`
- `internal/commitment/repository/closure.go`
- `internal/commitment/repository/head_binding_test.go`
- `internal/commitment/repository/outlook.go`
- `internal/commitment/repository/publication.go`
- `internal/commitment/repository/readiness.go`
- `internal/commitment/repository/repository.go`
- `internal/commitment/repository/sources.go`
- `internal/commitment/store_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/commitment_guidance_test.go`
- `internal/conformance/fixture_bite_test.go`
- `internal/conformance/handoff_single_source_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/handoff/facts.go`
- `internal/handoff/facts_test.go`
- `internal/handoff/handoff.go`
- `internal/handoff/render.go`
- `internal/handoff/render_test.go`
- `internal/handoff/sections_test.go`
- `internal/handoff/text.go`
- `internal/handoffdoc/document.go`
- `internal/handoffdoc/document_test.go`
- `internal/worktree/land_specless_test.go`
- `internal/worktree/parallel_census_test.go`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/dogfood-referent-shipped`
- `tests/canary/docs-currency-token-diet/missing-cli-inventory`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/load-validity-metadata/readme-shared-rule-drift`
- `tests/canary/load-validity-metadata/shared-rule-drift`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/workflow-guidance-anchors/acceptance-coverage-anchor`
- `tests/canary/workflow-guidance-anchors/agents-system-suite-route`
- `tests/canary/workflow-guidance-anchors/bootstrap-authority-before-after-softening`
- `tests/canary/workflow-guidance-anchors/bootstrap-authority-rule-deletion`
- `tests/canary/workflow-guidance-anchors/capture-sink-anchor`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/command-handoff-anchor`
- `tests/canary/workflow-guidance-anchors/core-absent-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/core-absent-category-context`
- `tests/canary/workflow-guidance-anchors/core-absent-category-oracle`
- `tests/canary/workflow-guidance-anchors/core-absent-category-setup`
- `tests/canary/workflow-guidance-anchors/core-absent-category-work`
- `tests/canary/workflow-guidance-anchors/core-absent-gate-authority`
- `tests/canary/workflow-guidance-anchors/core-absent-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/core-absent-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/core-absent-retained-integration-source`
- `tests/canary/workflow-guidance-anchors/core-absent-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/core-absent-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/core-absent-skills-guidance`
- `tests/canary/workflow-guidance-anchors/core-absent-upgrade-route`
- `tests/canary/workflow-guidance-anchors/core-cli-bounded-complete-response`
- `tests/canary/workflow-guidance-anchors/core-cli-call-pair`
- `tests/canary/workflow-guidance-anchors/core-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/craft-spec-empty-or-invalid-fence`
- `tests/canary/workflow-guidance-anchors/craft-spec-exact-literal-fence`
- `tests/canary/workflow-guidance-anchors/craft-spec-fence-after-slice`
- `tests/canary/workflow-guidance-anchors/craft-spec-further-notes-bare`
- `tests/canary/workflow-guidance-anchors/craft-spec-reader-sweep-name`
- `tests/canary/workflow-guidance-anchors/craft-spec-red-signal-grammar-retired`
- `tests/canary/workflow-guidance-anchors/craft-spec-template-verification-log`
- `tests/canary/workflow-guidance-anchors/decision-map-asset-path`
- `tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review`
- `tests/canary/workflow-guidance-anchors/delegated-green-predecessor`
- `tests/canary/workflow-guidance-anchors/delegated-opt-in-entry`
- `tests/canary/workflow-guidance-anchors/delegated-pending-predecessor-stop`
- `tests/canary/workflow-guidance-anchors/delegated-prerequisite-checkpoint`
- `tests/canary/workflow-guidance-anchors/delegated-repair-ownership`
- `tests/canary/workflow-guidance-anchors/dg-10`
- `tests/canary/workflow-guidance-anchors/dg-11`
- `tests/canary/workflow-guidance-anchors/dg-12`
- `tests/canary/workflow-guidance-anchors/dg-13`
- `tests/canary/workflow-guidance-anchors/dg-14`
- `tests/canary/workflow-guidance-anchors/dg-15`
- `tests/canary/workflow-guidance-anchors/dg-15-executable-red`
- `tests/canary/workflow-guidance-anchors/dg-15-mandatory-red`
- `tests/canary/workflow-guidance-anchors/dg-15-mandatory-red-capitalized`
- `tests/canary/workflow-guidance-anchors/dg-15-mandatory-red-emphasis`
- `tests/canary/workflow-guidance-anchors/dg-9`
- `tests/canary/workflow-guidance-anchors/edge-inventory-anchor`
- `tests/canary/workflow-guidance-anchors/fix-dont-park-section-relocated`
- `tests/canary/workflow-guidance-anchors/light-path-inline`
- `tests/canary/workflow-guidance-anchors/slicing-in-write-spec`
- `tests/canary/workflow-guidance-anchors/spec-retire-roadmap-row`
- `tests/canary/workflow-guidance-anchors/staged-command-sweep-anchor`
- `tests/canary/workflow-guidance-anchors/story-line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/structured-phase-progress-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-cross-pointers-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-decision-map-lifecycle-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-light-path-anchor`
- `tests/canary/workflow-guidance-anchors/write-spec-acceptance-round-exit`
- `tests/canary/workflow-guidance-anchors/write-spec-approval-table-emit`
- `tests/canary/workflow-guidance-anchors/write-spec-artifact-authorization`
- `tests/canary/workflow-guidance-anchors/write-spec-authorization-boundary`
- `tests/canary/workflow-guidance-anchors/write-spec-branch-split`
- `tests/canary/workflow-guidance-anchors/write-spec-contract-rederivation`
- `tests/canary/workflow-guidance-anchors/write-spec-conversation-authorization`
- `tests/canary/workflow-guidance-anchors/write-spec-conversation-fork`
- `tests/canary/workflow-guidance-anchors/write-spec-decision-source`
- `tests/canary/workflow-guidance-anchors/write-spec-degenerate-cheapness`
- `tests/canary/workflow-guidance-anchors/write-spec-fence-approval`
- `tests/canary/workflow-guidance-anchors/write-spec-frozen-base-and-tip-review`
- `tests/canary/workflow-guidance-anchors/write-spec-handoff-anchor`
- `tests/canary/workflow-guidance-anchors/write-spec-identified-coverage-default`
- `tests/canary/workflow-guidance-anchors/write-spec-late-uncertainty`
- `tests/canary/workflow-guidance-anchors/write-spec-map-sources`
- `tests/canary/workflow-guidance-anchors/write-spec-materiality-exit`
- `tests/canary/workflow-guidance-anchors/write-spec-ownership-fences`
- `tests/canary/workflow-guidance-anchors/write-spec-phase-ownership`
- `tests/canary/workflow-guidance-anchors/write-spec-post-slicing-handoff`
- `tests/canary/workflow-guidance-anchors/write-spec-promise-inflation-guard`
- `tests/canary/workflow-guidance-anchors/write-spec-reader-sweep-sequence`
- `tests/canary/workflow-guidance-anchors/write-spec-ready-map-authorization`
- `tests/canary/workflow-guidance-anchors/write-spec-review-made-conditional`
- `tests/canary/workflow-guidance-anchors/write-spec-review-tier-escalated`
- `tests/canary/workflow-guidance-anchors/write-spec-slicing-sign-off`
- `tests/canary/workflow-guidance-anchors/write-spec-ticket-approval-fields`
- `tests/canary/workflow-guidance-anchors/write-spec-ticket-breakdown-charge`
- `tests/canary/workflow-guidance-anchors/write-spec-unflagged-addition-removal`
- `tests/canary/workflow-guidance-anchors/write-spec-unique-row-id`
- `tests/canary/workflow-guidance-anchors/write-spec-verification-learning`
- `tests/canary/workflow-guidance-anchors/write-spec-verification-learning-contents`
- `tests/canary/workflow-guidance-anchors/write-spec-verification-log`

## Out of scope

- A derived plan form of `bench commitment plan` that builds the binding from a spec path: 7 edits, 2 gate runs.
- A status route that names the binding step for an eligible outcome with no binding, with the comment fix: 5 edits, 2 gate runs.
- A checkout-HEAD rule for roadmap-row sources, which the open drain learning asks for: 4 edits, 1 gate run.
- A preflight row that refuses a bound spec whose `Writes:` lines omit `.bench/commitment.json`: 4 edits, 1 gate run.

These estimates are planning figures, not approved priorities.

## Further notes

Evidence status: source-backed proposal, with no Bench runtime claim.
Remaining uncertainty: future red-to-green evidence and the ticket review.

### Flagged additions

Each item below goes beyond the literal decision source, and the reviewer can veto it at sign-off:

- E1: a new deliverable source reads only the checkout HEAD, so a worktree that predates the spec folds main first (SB08). The spec review judged this strict rule safer than a fallback to the tip. A fallback at plan forces the same fallback into the approval check. A candidate tree that lacks the spec could then publish a binding without its spec.
- E2: `bench commitment inventory` reads staged specs at the checkout HEAD, not the default-branch tip (SB03-SB05).
- E3: the `Spec binding` value format, its policy-order join, and the phrase `binding unknown (policy unreadable)` (SB23-SB28).
- E4: `AGENTS.md` moves the spec-stage handoff to the primary checkout after its landing (SB39).
- E5: the `Roadmap:` comment of the `craft-spec` template (SB40). It changes no behavior, and its removal costs nothing.
- E6: the new anchors and mirrors of the guidance sentences (SB31-SB36).
- E7: the close passes the start command to `bench handoff --next` (SB33). The reviewer chose this form on 2026-10-11 over a next command of `/bench-implement-spec specs/<slug>/spec.md` with the start command in the State body. `bench handoff` already derives the same placeholder command as its next command when the commitment row leads the board (`internal/handoff/handoff.go:84-90`).

### Source sentence to row

| source sentence | rows |
| --- | --- |
| `roadmap/FT395.md:3-4`: "A spec-stage handoff called a run ready while the staged spec was off `main` and unbound, so `bench commitment start` refused." | the defect before this fix; SB17, SB18, and SB23 prove the fix |
| `roadmap/FT395.md:4-5`: "The spec-stage close lands the staged spec and its deliverable binding on `main` before the handoff." | SB01, SB02, SB12, SB14, SB17, SB31 |
| `roadmap/FT395.md:5`: "The handoff names the binding state." | SB23, SB24, SB25, SB38 |
| `roadmap/FT395.md:7-8`: "The binding derives from the spec's `Roadmap:` line and its path, so `bench status` prints the complete `bench commitment start` command." | SB03, SB04, SB05, SB18, SB22, SB32 |
| `roadmap/FT395.md:8-9`: "A ticket `Writes:` line owns `.bench/commitment.json`, so the build preflight stays green after the binding." | SB36, SB37 |
| criterion `.close-lands` | SB17, SB23 |
| criterion `.status-command` | SB18, SB22 |
| criterion `.preflight-green` | SB36, SB37 |
| reviewer direction: the start follows the spec stage with the bound deliverable | SB20, SB33 |
| reviewer direction: the handoff pins that command | SB33 |

The occurrences at `roadmap/FT395.md:14-18` support the problem statement.
The 2026-10-05 FT390 occurrence is the unbound status command (SB19).
The 2026-10-06 FT376 and 2026-10-07 FT392 occurrences are the refused start and the separate binding commit (SB17, SB18).
The 2026-10-07 FT393 occurrence is the `paths-authorized` red (SB36, SB37).
The 2026-10-10 build-directions occurrence is the second landing and gate (SB17).

### Pre-review proof checklist

- `Cited symbols`: each symbol below resolves in the tree.
  - `Store.Plan`, `Store.Approve`, `Store.Inventory`, `Store.Policy`, `validateSources`, `defaultPolicy`, and `policyAt` resolve in `internal/commitment/repository/repository.go`.
  - `validateDeliverable`, `validateDeliverables`, `sourceAt`, `ticketsOnlyAt`, `sourceRevision`, and `validateSourceAt` resolve in `internal/commitment/repository/sources.go`.
  - `authorizeCandidate`, `approvedTransition`, `branchReference`, and `protectedCandidate` resolve in `internal/commitment/repository/candidate.go`.
  - `AdmitPublication` resolves at `internal/commitment/repository/publication.go:19`.
  - `Store.Start` resolves at `internal/commitment/repository/admission.go:16`, and `readyFor` resolves in `internal/commitment/repository/readiness.go`.
  - `soleDeliverable` resolves at `internal/commitment/outlook.go:109`, and `withCommand` at `internal/commitment/commitcmd/command.go:183`.
  - `Unsettled` resolves at `internal/commitment/delivery.go:127`.
  - `pins` and `liveSpecs` resolve in `internal/handoff/facts.go`, and `ShapeSection` in `internal/handoff/text.go`.
- `Import edges`: `internal/handoff` gains direct edges to `internal/commitment` and `internal/commitment/repository`. `go list -f '{{.Imports}}' ./internal/handoff` shows neither today, and `go list -f '{{join .Deps " "}}' ./internal/handoff` already lists both through `internal/status`. `go list -f '{{join .Deps "\n"}}' ./internal/commitment/repository` lists no `internal/handoff`, `internal/handoffdoc`, or `internal/status`, so the edges add no cycle.
- `Source-row clauses and occurrences`: the source table above quotes each clause, and the paragraph after it maps every occurrence.
- `Promised field labels`: `Spec binding`, `unbound`, `binding unknown (policy unreadable)`, and the refusals `commitment source "<id>" refused:`, `deliverable "<path>" is not staged` (`internal/commitment/repository/sources.go:89`), and `is not approved for outcome`.
- `Changed-function callers`: `bench consumers --production` lists each caller below.
  - `validateDeliverable`: `internal/commitment/repository/admission.go:41` (`Store.Start`), `internal/commitment/repository/readiness.go:60` (`Store.readyFor`), and `internal/commitment/repository/sources.go:98` (`validateDeliverables`)
  - `validateDeliverables`: `internal/commitment/repository/repository.go:104` (`Store.Plan`)
  - `validateSources`: `internal/commitment/repository/repository.go:111` (`Store.Plan`) and `:183` (`Store.Approve`)
  - `approvedTransition`: `internal/commitment/repository/candidate.go:67` (`authorizeCandidate`)
  - `Store.Inventory`: `internal/commitment/commitcmd/command.go:212` (`inventory`)
  - `pins`: `internal/handoff/facts.go:71` (`collect`)
- `Copy survival`: none. No new owner replaces a copy.
- `Rendered-shape readers`: the handoff section gains one label line. No code reads the `Spec` or `Spec status` values, and `TestSectionRendersOnePairPerLiveSpec` counts `Spec: ` lines, which `Spec binding: ` does not match. The inventory row shape stays unchanged.
- `Pin operators`: SB29 uses `strings.Contains` over the Shape text. SB30 uses `strings.Contains` over the adjacent pair.
- `Entry reads`: the HEAD read runs below `Store.Plan`, `Store.Approve`, and `Store.Inventory`. Real linked worktrees grade each entry in SB01-SB11. No entry takes an injected revision.
- `Derived expectations`: each identity comes from `commitment.Identity` over the fixture bytes, each path from the fixture writes, and each outcome id from the fixture policy. No expectation calls the code under test.
- `Consolidated rules`: the table below lists the one consolidated rule.
- `Quantified obligations`: "each new deliverable source" maps to SB01, SB06, SB12, and SB14. "Every other source" maps to SB10 and SB11. "Each live spec" maps to SB23 and SB30. "Each new pinned sentence" maps to SB31-SB36 and SB42.
- `Workflow-step writes`: the close commits the binding before the plan-only preflight. That preflight grades the range that holds `.bench/commitment.json`, so the fence and one `Writes:` line hold it (SB37).

| site | old rule | new rule | row |
| --- | --- | --- | --- |
| `Store.Plan` | Every source reads the default-branch tip. | A new deliverable source reads the checkout HEAD, and every other source reads the tip. | SB01, SB08, SB10, SB11 |
| `Store.Approve` | Every plan source reads the default-branch tip. | A new deliverable source reads the checkout HEAD, and every other source reads the tip. | SB02 |
| `Store.Inventory` | Staged specs that the tip holds, read at the tip. | Staged specs that the checkout HEAD holds, read at that HEAD. | SB03, SB05 |
| `approvedTransition` | Every plan source reads the default-branch tip. | A new deliverable source reads the candidate tree, and every other source reads the tip. | SB12, SB13, SB14, SB15 |
| `Store.Start` and `readyFor` | Every deliverable reads the default-branch tip. | Unchanged. | SB20, SB21 |

### Reader inventory

These readers change their result, and each change is intended:

- `internal/commitment/commitcmd/command.go:211-219` prints the inventory, which now lists a HEAD-held staged spec.
- `internal/commitment/outlook.go:100` returns the bound path after the one landing, so `bench status`, `bench roadmap`, `bench dashboard`, and the session-start status show the complete command sooner.
- `internal/handoff/facts.go:83-100` writes the new `Spec binding` line.

These readers stay unchanged:

- `internal/status/delivery.go:19-41` cuts the `bench commitment ` prefix and keeps the rest.
- `internal/commitment/repository/readiness.go:40-83` and `internal/commitment/repository/admission.go:30-44` keep the default-branch tip.
- `internal/preflight/command.go:105-107` and `internal/preflight/charge_pack.go:116-117` read readiness only.
- `internal/spec/spec.go:69-88` keeps its `Roadmap:` grammar. The inventory reads the parsed value.
- `internal/status/handoff.go:56-98` reads only the section key and the `Worktree tip` field.
- No `.mjs` script, hook, or workflow file reads the commitment policy, the inventory, or the handoff labels.

The pinned strings stay unchanged by design:

- The start command and its placeholders: `cmd/bench/commitment_test.go:26`, `internal/status/delivery_test.go:36-37`, `internal/worktree/land_identity_test.go:211`, and `internal/commit/commitment_route_test.go:72-79`.
- The start help row: `cmd/bench/help_inventory_test.go:85`.
- The outlook row: `cmd/bench/commitment_command_test.go:32` and `internal/roadmap/commitment_blocks_test.go:29-32`.

### Observations

The comment at `internal/commitment/outlook.go:37-40` says: "A reader therefore never selects work that a start refuses."
`eligible` checks no deliverable approval, while `Store.Start` refuses an unapproved deliverable (`internal/commitment/repository/admission.go:31-40`).
D7 keeps this behavior, and Out of scope prices the route that would end the contradiction.

All five binding commits of the problem table changed only `.bench/commitment.json`.
The landing of each binding needed an exact approved plan receipt (`internal/commitment/repository/candidate.go:160-188`).
Only `acba6bed` names its plan in its message.

Each existing plan and approve test roots at the primary repository, where HEAD equals the tip (`internal/commitment/commitmenttest/repo.go:134-146`).
The one worktree-rooted `AuthorizeCandidate` test uses a legacy adoption that binds no deliverable (`internal/commitment/repository/publication_test.go:201-212`, `internal/commitment/commitmenttest/admission.go:76-79`).
So the static search finds no existing fixture that the HEAD rule turns red. The SB1 author runs the commitment packages before the first production edit to confirm it, and lists each red as a posture change. A red in a file outside the ticket's `Writes:` line is a plan expansion under `.bench/BENCH.md`, not an inline fix.

The slicing research found three files at or near the 400-line structure limit.
`internal/handoffdoc/document.go` holds 400 lines, `internal/conformance/handoff_single_source_test.go` holds 399, and `internal/commitment/repository/repository.go` holds 385.
The lane's growth check reds a file over the limit that gains a line (`internal/structure/structure.go:308-331`).
Each ticket therefore keeps these files at 400 lines or fewer.
`internal/worktree/parallel_census_test.go:1061` pins the package test count at 773, so each new top-level `internal/worktree` test raises that pin in the same ticket.

The reviewer stated on 2026-10-11 that light-path tickets of FT231 land on main during this build.
Each one can move a cited line and add a `CHANGELOG.md` section, so the build finds each guidance edit by its quoted text.
`bench worktree merge` refuses a textual `CHANGELOG.md` conflict, and the raw merge belongs to the reviewer.
When the reviewer is away, one commit holds back the branch's own section, the fold runs, and one commit restores the section on top.
The inventory light path used that route at `fa9f149c`.

### Not verified

- Whether `spec.CommitTree` and `ticketsOnlyAt` resolve a tree object id as they resolve a commit. `protectedCandidate` already reads roadmap rows at a tree id, but no current test binds a tickets-only folder at a tree id.
- Whether the landing gate runs a commitment check other than `AdmitPublication` over the composed tree.
- The exact path of the FT393 `paths-authorized` red. The tree keeps only the occurrence line at `roadmap/FT395.md:17`.
- The executed census of existing tests under the HEAD rule. A static search finds none that would turn red, and Observations states the evidence.
- Whether `roadmap.ProjectSequence` leaves `ROADMAP.md` byte-identical for an approve that only adds a binding (`internal/commitment/repository/repository.go:287`).
- Whether a test outside `internal/handoff` pins the exact bytes of a handoff section by a form other than the label text. A tree-wide search for `Spec status`, `Recorded base`, and `Worktree tip:` finds no such test outside `internal/handoff`, `internal/handoffdoc`, and `internal/conformance`.

### Validation plan

The build first observes SB01 red at `Store.Plan`, before the first production edit.
Each independent expectation then records one red through a named omission or a behavioral swap.

The central author probe swaps the candidate tree for the default-branch revision in the new-deliverable branch of `approvedTransition`. SB12 and SB14 must fail through `bench probe`.
The coordinator probe omits the new-deliverable branch of the plan rule, so a new deliverable falls back to the default-branch tip. SB01 must fail.
The two probes differ in kind and in site, and each one must report a restored subject.

Run `bench test --package ./internal/commitment --package ./internal/commitment/repository --package ./internal/commitment/commitcmd` for SB1.
Run `bench test --package ./internal/commit`, `bench test --package ./internal/worktree`, and `bench test --package ./cmd/bench` for the SB1 journey and pins.
Run `bench test --package ./internal/handoff --package ./internal/handoffdoc`, the commitment packages, and the `cmd/bench` package for SB2.
Run the `internal/anchors`, `cmd/bench`, and `internal/conformance` package tests and the `docs-currency-workflow`, `guidance-prose-budgets`, `prose-mechanics`, and `skills-index-command-adapters` checks for SB3.
The coordinator runs the whole gate only at the landing.

### Bootstrap of this spec

This spec's own close binds through the current two-landing route, because the one-landing route ships with this spec.
The close lands the spec and its tickets spec-less, and a second planning worktree then plans, approves, commits, and lands the binding.
The handoff then pins the start command that `bench status` prints for FT395.

### Review chronology

- SR1: the spec review ran on Fable at xhigh effort, by reviewer direction of 2026-10-11, over commit `a0cc7a54`. It accepted in iteration 1 with eight findings, and none of them blocked.
- The coordinator folded the eight findings into this revision. R1 names the freed command lines and the kept anchors. R2 names the SB09 refusal sentence. R3, R4, R5, and R6 move settled facts into the evidence. R7 became the decision E7, which the reviewer closed on 2026-10-11. R8 keeps E5 flagged for veto.
- Slicing research: two blind readers for each chunk, one on Haiku at xhigh and one on Sonnet at high. The coordinator folded five findings, and none changes a row or a chunk. F1 corrects the package overlap of SB1 and SB2, and F2 replaces the line joins that broke three canary needles. F3 exports the unresolved-default test for the handoff, and F4 names the one `CHANGELOG.md` entry. F5 adds the `cmd/bench` check floor to SB2 and SB3, and the commitment packages to SB2.
- SB1 slicing: one fresh slicer on Opus at high effort. The coordinator merged its second ticket into the first, because that ticket added only a journey test over the first ticket's seam. SB19 moves to the planning worktree, where a wrong outlook can turn it red, and story 18 drops "on main" to match. SB20 names the default-tip read that it pins, because SB21 owns a start read that moves to HEAD. The planned SB2 and SB3 tickets take the numbers 2 and 3.
- SB2 slicing: one fresh slicer on Opus at high effort, with no re-cut. SB44 joins the map, because the spec promised `unbound` for an unresolved default branch and no row graded it. The primary checkout owns the main section whatever the default branch is (`internal/handoff/handoff.go:140-147`), so the case is reachable.
- SB3 slicing: one fresh slicer on Opus at xhigh effort, by reviewer direction, with no re-cut. The coordinator added the `internal/anchors` and `commitmenttest` floors to the completion plan. It also put "After the last spec edit" into step 2 of the close order. The command and `AGENTS.md` both name the handoff venue, because the command ships to linked repositories and `AGENTS.md` does not.
- Main moved twice during slicing. The guidance-trim landing moved the cited lines of `.bench/BENCH.md` and `AGENTS.md`. The inventory-revision-metadata light path, by reviewer direction of 2026-10-11, makes the inventory read the status and the scope from the identity bytes. So the Won't handle line for a working-tree `Roadmap:` edit left the edge inventory.
- TR1: the ticket review ran on Fable at xhigh effort, by reviewer direction, over commit `513e3642`. It accepted in iteration 1 with nine findings, and none of them blocked. The coordinator folded T1-T9. They fix the line citations after the light path, the size counts, and the venue sentence's word count. They also name the Shape prose test, ticket 2's edge, the out-of-line red rule, the journey obligation, and the SB44 fixture.
