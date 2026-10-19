# State the plan-commit rules in the implement and review phases

Blocked by: 01-grade-the-chunk-table-against-the-map.md
Writes: .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-line/references/bounded-repair-policy.md, .agents/skills/bench-craft-spec/references/completion-plan.md, internal/reviewrecord/coverage.go, internal/gate/review_checkpoint_commits_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, projects/benchkit.md, CHANGELOG.md
Covers: PC22, PC23, PC24, PC25, PC26, PC27, PC28, PC29, PC30, PC31, PC32, PC33, PC34, PC35, PC36, PC39

## What to build

The implement phase's Land section states the plan-commit rules, and the review phase, the repair policy, the completion-plan reference, and the chain-gap refusal agree with them. The spec's decision "The guidance" and the reviewer answer of 2026-10-19 fix the rules. Write each rule in ASD-STE100 prose under `.agents/skills/bench-craft-spec/references/ste-prose.md`.

In `.agents/commands/bench-implement-spec.md`, replace the clause `Only record commits, recorded folds, and comment-only corrections follow the chunk tip.` in the serial-commit paragraph with the moved-tip rule (PC23): before the chunk's green checkpoint, each source change except a record commit or a recorded fold moves the chunk tip and takes a refreeze. That paragraph has 5 sentences now, so it holds at most one more. Then add one new paragraph after it. The new paragraph states, in at most 6 sentences of at most 25 words each:

- each plan commit uses `bench commit --preflight-build <slug>`, and no dispatch starts until its `commit-chain` line reads green (PC22);
- a red `check-floor` row blocks the next charge and the review charge (PC29);
- a plan commit after the green checkpoint lands in the chunk gap after the chunk-gap fold of `main`, and it joins the successor's review delta, so a successor-only plan commit waits for that checkpoint (PC24, PC39);
- after the last chunk's freeze, each source change except a record commit or a recorded fold takes a refreeze of that chunk (PC25);
- a coverage-cell citation and its cited test land in one commit (PC27);
- a review's prose-only spec corrections go in one commit before the repair's refreeze, or into the next plan commit before a ticket dispatch (PC28).

The new text must match no Forbid row of `internal/anchors/registry_chunk_chain.go` whose file is the implement phase. These needles are the likely collisions: ``review delta of the last chunk``, ``the last chunk refreezes``, ``When the orchestrator freezes a chunk after its last ticket``, ``The fold at the start of the last chunk comes before its first ticket dispatch.``, ``folds `main` only before the first chunk, at the start of the last chunk``, and ``--chunk <id>` before the successor``. Keep the fold sentence that begins `A build folds` byte-identical (PC36): its Require row stays. The fold paragraph already has 6 sentences, so add nothing to it.

Raise the implement-phase row of `projects/benchkit.md` "Guidance prose budgets" from 93 to 95 in the same commit, and state the reason in the commit message: the plan-commit rules are the phase's own build rules (PC34). If the rules do not fit one paragraph inside the sentence bounds, stop and report before you edit the budget; a value above 95 is a reviewer decision.

In `.agents/commands/bench-review-implementation.md` step 1, rewrite the base sentence at line 97 (PC26). The successor base is the last recorded source tip of its predecessor, never a record commit; it is the last recorded fold after that tip, or that tip when no fold follows it. State that a plan commit after the predecessor's green checkpoint sits in the successor's delta. Fold the sentence at line 98, ``A later plan commit is never a chunk base.``, into this rewrite, so the step 1 paragraph keeps at most 6 sentences of at most 25 words each. The new text must match no Forbid row of `internal/anchors/registry_chunk_chain.go` or `internal/anchors/registry_retained_workflow.go` whose file is the review phase; ``review delta of the last chunk`` and ``so the review chain refuses a `main` merge between two chunks`` are the likely collisions.

In `.agents/skills/bench-craft-line/references/bounded-repair-policy.md`, add one sentence after line 47: a comment-only correction before the green checkpoint stays evidence-only and still takes a refreeze (PC32). In `.agents/skills/bench-craft-spec/references/completion-plan.md`, name the `chunk-rows` row of build preflight beside the `check-floor` sentence: it compares the chunk table with the coverage map in both directions (PC33).

In `internal/reviewrecord/coverage.go:102`, replace the clause `only record commits, recorded folds, and comment-only corrections follow a chunk tip` with the exact clause `only record commits and recorded folds follow a chunk tip until its checkpoint` (PC31). Update the one exact-text assertion in `TestReviewCheckpointChainGapNamesTheExpectedBase` (`internal/gate/review_checkpoint_commits_test.go:49`). The substrings `stale review chain gap` and `chain gap` that other tests assert stay.

Add one Forbid row for `.agents/commands/bench-implement-spec.md` to `internal/anchors/registry_chunk_chain.go` with the exact needle `comment-only corrections follow the chunk tip` (PC30). Its mirror test `internal/anchors/registry_chunk_chain_test.go` holds an independent expectation of each row, so add the mirror rule and record its red under the expectation exception of `AGENTS.md`.

Add one `CHANGELOG.md` entry under `[Unreleased]` for the new `chunk-rows` row and the new chain-gap clause. Keep `roadmap/FT380.md` unchanged (PC35): the delivery closes no roadmap row.

The staged `guidance-weight` spec also writes `.agents/commands/bench-implement-spec.md`, its budget row, and `internal/anchors/registry_chunk_chain.go`. The later of the two builds takes the other's landed text through its staleness pass.

## Acceptance

- [ ] The Land section requires `bench commit --preflight-build <slug>` for each plan commit before the next dispatch.
- [ ] The Land section states the moved-tip rule before the green checkpoint, the chunk-gap rule after it, and the refreeze after the last chunk's freeze, each with the record-commit and recorded-fold exception.
- [ ] The Land section states the citation rule, the prose-correction batch rule, the fold-before-plan-commit rule, and that a red `check-floor` row blocks the next charge and the review charge.
- [ ] The implement phase no longer contains `comment-only corrections follow the chunk tip`, and the new Forbid row reds a restored clause.
- [ ] Review step 1 states that the successor base is the last recorded source tip of its predecessor, never a record commit.
- [ ] The chain-gap refusal contains `only record commits and recorded folds follow a chunk tip until its checkpoint`.
- [ ] The repair policy and the completion-plan reference state their new sentences.
- [ ] The implement phase stays inside its budget row of 95 lines, and every PC-C2 completion-plan command passes.
