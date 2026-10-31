---
description: Implement a spec (or a clearly-scoped change) at the pre-agreed seams. Use after /bench-write-spec, or for a change the light path in .bench/BENCH.md allows without a spec. An opt-in `--full <spec>` run reaches the final green landing.
---

# /bench-implement-spec — do the work at the seams

## Entry orientation
This implementation phase starts from an approved spec, or from a change small enough for the light path. It declares the line, starts its committed outcome through `bench commitment start`, and works vertical slices at the pre-agreed seams. If there is no spec, the change must fit the light path in `.bench/BENCH.md`'s "Right-size the process" paragraph. If it does not fit the light path, route to `/bench-write-spec` first. If a spec-backed run has no `specs/<slug>/tickets/` directory or that directory contains no ticket files, return to `/bench-write-spec`; ticket slicing and approval belong there. A light-path change needs no commitment start; `.bench/BENCH.md` owns that rule.

## Exit handoff
Close by reporting the implemented stories, each acceptance row's coverage status, and the landed commits. A reviewed spec-backed build closes through the "Land" section below. A light-path build goes straight to `/bench-final-check`. A build that stops short exits through "When the build stops short" below, which recommends its one durable next action.

## Declare the line, validate the tickets, route the venue

Declare the line before you touch code. `craft-line` owns the template, the tier decision, and the escalation ladder. After the write-spec-phase approval, create or retain one integration worktree at the reviewed graph commit; that commit is its frozen review base. Enter it and run `bench preflight build <slug>` before the first edit. If every red row is `completion-plan`, `writes-landed`, or `map-sources`, the staleness pass below takes the red. Any other red stops the phase.

After that preflight, run `bench preflight build <slug> --audit-spec --base <reviewed graph commit> --source-tip <current tip>`. Then charge the staleness pass that `.agents/skills/bench-implement-spec/references/staleness-pass.md` states, or omit it only under its omission rule. The orchestrator reviews the returns and decides each finding without a reviewer stop. Commit the fixes on the integration source as the enabling plan commit of `.bench/BENCH.md`'s approved plan-expansion policy. This amendment corrects claims, citations, fences, landed work, and any preflight red that the pass took, and it does not reopen the approach. Before the version 2 plan amendment, the preflight reruns green, and the build then continues.

The reviewer approves the spec and the whole ticket graph once, in `/bench-write-spec`, so a ticket in that graph starts without a new approval stop. Before the orchestrator prepares an author charge, it runs `bench consumers` on the symbols that the ticket names. For each caller, test, or fixture that the acceptance reaches, it widens the `Writes:` line under `.bench/BENCH.md`'s approved plan-expansion policy. The plan amendment lists each mutation that a ticket names as a requirement `probe` on its registered check, not only as a helper bite. The orchestrator charges the docs ticket with each spec-versus-code divergence that the build settled, and corrects the spec prose before the docs review.

Run `bench preflight build <slug> --charge --ticket <ticket> --base <base> --source-tip <tip>` and read the prepared evidence identity it returns. Verify the fence, dependency completion, and expected source tip before the ticket's fresh author starts. This charge preparation is the one build preflight of an author or repair dispatch. Prepare that charge at the commit that records the author's assignment, because a charge from an earlier commit fails `bench preflight evidence <id> --check-current`. Each ticket goes to a fresh author session under `.bench/BENCH.md`'s authorship rule.

An author or repair charge carries the ticket, its spec sections, its coverage rows, its `Writes:` fence, the evidence identity, and the declared line. When the repository has `docs/known-issues.md`, the author charge names that path. The charge says that the author changes only what the acceptance needs, and it prints the exact grammar of each Bench verb in its steps. Before the first edit, the author lists every place the change must reach: callers, tests, fixtures, and exports. Code that the author moves or merges keeps its error handling and validation. The author's return names each untouched `Writes:` entry with a one-clause reason.

A witness is one demonstrated red of one independent test expectation. The charge caps the witnesses at one for each coverage row plus one for each production branch that the ticket adds. The author meets the cap with fewer independent expectations, such as one table-driven test, and skips no witness. The author runs the development checks and witnesses of several packages in one `bench test` call that repeats `--package`. The author probes several witnesses of one subject in one `bench probe --mutants` call. The author sends independent reads in one turn, and runs each development and witness test with `bench test --package <expr> --run <regex>`.

Build action requires a narrow author read: act only after this session reads the metadata and ticket pages and targeted sources. An author does not retrieve every evidence page. `bench preflight evidence <id> --verify` is the separate artifact-integrity check.

Build action requires available required context: a retrieval receipt, a terminal cursor, or another consumer's delivery never replaces what this session reads itself. Build action requires a current-action binding: act only after `bench preflight evidence <id> --check-current` binds the artifact to the current assignment and source pair. Build action requires reviewer approval, which a generated charge, a staged artifact, or a verified artifact never supplies. Build action requires the complete task supplement, which this session writes. The supplement names selectors, model, effort, cap, and mutation, and verified evidence never supplies it.

Follow `.bench/BENCH.md`'s ticket authorship and chunk-review policy. `craft-delegate` governs independent review and brief read-only diagnostic consultation. Neither role receives implementation or repair authorship. Before the first dispatch, the orchestrator takes the plan amendment that `.bench/BENCH-reference.md` states. The orchestrator reads manifests, returns, and verdicts, not code, and it refreshes `bench handoff` at each chunk checkpoint and at the final reconciliation. The orchestrator runs `bench probe` only in a tree where no author runs checks, because an author check can read the probe's mutant.

## Build

Before each approved slice, identify its acceptance target and existing verification route. At an approved TDD seam, follow `craft-tdd` for its behavioral-red sequence, minimal compiled setup, and row classifications.

A material action changes behavior, a verification target, or a premise that determines the next action. After each material action, rerun that route. Inspect the result before continuing. One material action can contain several related edits before the rerun. If a result contradicts the approved behavior or seam, stop. Route it through the wrong-spec path under "When the build stops short."

Work approved tickets in dependency order inside their planned chunks. Use TDD only where `craft-tdd` marks the seam. When the spec carries a coverage map, seed the
harness's native task list from `bench coverage <spec>`, whose `rows[N]{story,behavior,seam}` projection is one task per row. The spec and the ticket are final. Do not evaluate another approach. Implement the ticket as written, run its focused checks, then continue to the commit below.

Apply `craft-line`'s retained implementation continuation policy throughout the ticket graph. At initial review or before post-review repairs, including resumed work, read [the bounded repair policy](../skills/bench-craft-line/references/bounded-repair-policy.md).

A prose-only owner edit is an implementation write. Before another session edits it, follow `craft-delegate`'s recorded user-directed author-transfer rules. Run fresh-session adoption after the final committed owner-byte change. A later owner-byte edit invalidates that evidence; rerun it against the new committed tip before review.

Before the first ticket dispatch of a declared build, run `bench worktree build --install <target>` from the primary checkout. After a Go change or a fold reaches that build, run the install again before the next dispatch or checkpoint. The Command Notes of `.bench/BENCH-reference.md` define a declared build and describe the install.

Before each ticket, revalidate prepared source identity and required source bytes. Regenerate a charge when either changes. The author commits the ticket on a green lane. After the ticket commit, the author runs only the probes of the plan, because the chunk checkpoint grades every other check. Run the plan's named probe exactly, and do not replace it with a similar mutation. If the named probe cannot run, stop and amend the plan before verification.

An author that takes a regenerated charge reuses exact available source bytes only after the new manifest verifies their membership, role, and requiredness. A matching body digest alone authorizes no reuse. A fresh author makes its own narrow author read from the trusted evidence identity. A transferred final cursor and another consumer's receipt deliver no byte to it.

The ticket commit rewrites each covered `planned` seam cell into the citation form, so that `bench record ticket` accepts the ticket. After the ticket commit, run the probe, then `bench record ticket`, then the record commit that its `next=` line prints. The first-ticket charge of each chunk names `bench record ticket --base <chunk base>`, and step 1 of `/bench-review-implementation` states that base. After the last ticket in a chunk, freeze the chunk delta and run the three review axes before advancing. The last ticket's author performs that freeze, as the Land section states. `/bench-review-implementation` owns those axes and the review record.

The Land section states how the authors of a chunk stage and commit the verification and probe record. Prepare the review charge from that record commit. The sequence is the author record commit, then the review charge, then the axis dispatch. The review charge says that an untouched `Writes:` entry is not a defect. An edit outside the acceptance's reach is a finding, even when its path is in `Writes:`. Before you widen a `Writes:` line or expand the plan or its gate coverage, apply `.bench/BENCH.md`'s approved plan-expansion policy.

A build dispatches no light-path fix while it runs. A fix that the build's own outcome needs lands as the build's own commit under `.bench/BENCH.md`'s fix-don't-park rule.

## Land

Tickets commit green serially in `Blocked by:` order on the retained integration source. After each chunk review, accepted findings go to fresh repair sessions under `.bench/BENCH.md`'s repair rule, and each repair commits with current repair coverage. Start the successor only after that repair closes. Each source change except a record commit or a recorded fold moves the chunk tip and takes a refreeze of that chunk. This rule applies only before a chunk's green checkpoint and after the last chunk's freeze. A successor-only plan commit therefore waits for that checkpoint.

Plan commits land before the ticket merge, and a `main` merge lands before the first chunk or as a recorded fold. A plan commit that a dispatch follows uses plain `bench commit`, and a plan commit that no dispatch follows uses `bench commit --preflight-build <slug>`. A plan commit after the green checkpoint lands after the chunk-gap fold of `main` and joins the successor's review delta. A coverage-cell citation and its cited test land in one commit. A review's prose-only spec corrections go in one commit before the repair's refreeze, or in the next plan commit before a ticket dispatch.

A build folds `main` with `bench worktree merge --from main <target>` before the first chunk, or in a chunk gap before the first dispatch of the next chunk. The Command Notes of `.bench/BENCH-reference.md` state how a chunk gap opens and closes. Each fold after the first chunk takes `bench record fold` before the next commit. After the last freeze, `main` moves only through the reviewer's landing override, and a recorded fold covers the completion source without a refreeze. A build never folds `main` between two tickets of a chunk. Inside a chunk, the post-ticket preflight passes the chunk's `--base` and the current `--source-tip`, so a moved `main` does not route a fold.

Retain author verification and all three native review results in `reviews/<slug>.md`. After the last ticket of a chunk, that ticket's author freezes the chunk and records the chunk entry with `bench record ticket`. Until the review record holds sibling chunks, a delegated run freezes its concurrent chunks in plan order. `bench record ticket` stages the last probe response, so an author runs `bench record excerpt` only for an earlier probe of the same ticket, before its next Bench call. `bench record ticket` binds each entry to its run's commit. The author stages with only the record path dirty, and one record commit on the ordinary lane follows the last author's record.

Before a repair dispatch, run `bench record assignment`, then `bench spec widen`, then one `bench commit`. A repair session runs each check that its repair delta reaches at its repair commit, and `bench record ticket` carries the rest and refreezes the chunk. Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`. This checkpoint waits for the live background grade of its key and can pause, and it reuses the green only when that grade ends green. An ordinary lane pass or ordinary green cannot replace this checkpoint, and when evidence fails, follow its named record or review action. After a repair, apply the bounded repair policy's current-evidence rule and its narrow evidence-only exception before retrying.

Before the orchestrator prepares the charge of a confirming round, the review record holds each red that the repair's `Witness return:` field reports. `bench record verification --probe-result <file>` retains the red of a planned probe requirement. The orchestrator writes each other red as prose in `reviews/<slug>.md` outside the fenced record, then commits the record. Before the initial review charge, the record also holds each witness red that an author reports. The orchestrator writes those reds the same way, in the same plan commit as the next assignment.

After the last chunk, reconcile every acceptance row and the integrated behavior. The last chunk runs no chunk checkpoint gate, because the complete checkpoint that `/bench-final-check` owns grades every chunk again. The last chunk keeps its chunk entry, its verification entries, and its three review entries. From the clean destination, hand the accepted source to `bench worktree land`; it composes, gates, publishes `Status: implemented`, and releases the source. Then run `/bench-final-check`.

Write the ordinary assessment record before the `bench worktree land` step, and append the landing evidence after it. The review and final-check phases own the other assessment updates. Keep the retained run identity available for those phase boundaries.

The orchestrator keeps one build log at the primary checkout's `capture/build-log/<slug>.log`. For each reviewer direction that saves a step, it appends one line that names the direction and the step that it removed. The retro author reads that log. At each chunk checkpoint, the orchestrator reports the wall time, the bookkeeping commits against the change commits, and each step that repeats an earlier record. In the same update, it names the cut that removes each such cost.

## When the build stops short

Report the state: what is done, what remains, the coverage table, and what consumed the cap. Keep landed work landed. For a declared build, follow the restore steps in the Command Notes of `.bench/BENCH-reference.md`. Then route by cause:

- Wrong tier: apply `craft-line`'s ladder.
- Wrong spec: route to `/bench-write-spec` with the finding quoted.
- Wrong scope: propose the split for the reviewer to decide.
- Commitment refusal: report the decision that the refusal names; `.bench/BENCH.md` owns the commitment rule.

## `--full <spec>`

This command orchestrates the fresh ticket authors, chunk reviews, the final landing, and `/bench-final-check`. It refuses on a missing or unnamed spec. At every phase boundary it writes the phase reached into `capture/session-handoff.md` and refreshes the pin with `bench handoff`. A green chunk checkpoint and its handoff refresh are not a phase exit. The orchestrator continues into the successor chunk in the same turn, and it stops only on a `craft-line` stop condition. The review phase owns cross-harness review opt-in.

A chunk delta can grow past its approved plan. It touches a path that requires the reviewer's decision under `.bench/BENCH.md`'s plan-expansion policy, or it adds a ticket the plan does not list. Such a delta pauses to ask the reviewer before escalating tier. For each other path outside the fence, widen the `Writes:` line with `bench spec widen` without a pause.

A `--full` run accepts these options. It takes `--delegate` and `--auto-approve` by default, and `--no-delegate` or `--no-auto-approve` turns one off. `craft-line`'s "Delegated author lines" owns the default author limit, review line, and consultation line.
- `--delegate` extends a `--full` run with `.bench/BENCH.md`'s delegated run policy. It refuses without `--full`, an approved spec, or an approved ticket graph. Declare the configured model, effort, iteration cap, and author limit before the first dispatch. Record each assignment in the completion plan with `bench record assignment`. The "Declare the line, validate the tickets, route the venue" section states when to prepare that author's charge. A resumed delegated run keeps the recorded identities, source pins, replacement history, and pending obligations.
- `--reviewer <model> <effort>` sets the line of every review axis in the run.
- `--consultant <model> <effort>` sets the line of every read-only diagnostic consultation in the run.
- `--auto-approve` is the reviewer's approval, given in advance, for `Writes:` widenings, spec and ticket expansions, and repair rounds past the bounded repair allowance. Under it, the orchestrator fixes a found defect at once and does not defer it.
