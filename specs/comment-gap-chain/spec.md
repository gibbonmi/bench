# The review chain accepts a proved comment gap

Status: staged

Roadmap: FT439

Decision source: named reviewed artifact — `roadmap/FT439.md` at its pinned identity `sha256:8bd6b4fe1bfaa7193b00715fed1d003cb2212467a025a1b13bd5da9bdc6600d9` (the source of outcome `FT439-comment-correction` in `.bench/commitment.json`), which holds the reviewer decision of 2026-10-21.

Verification log: 2 iteration(s) to accept — the spec round and the ticket round each took 1 iteration. The spec round folded ten prose and accounting findings. After that round, the CG-C2 slicing scratch runs found two prose-bound reds in exact guidance sentences, and the coordinator refit both sentences before the ticket round graded them. The ticket round folded three prose findings. Reviewer sign-off 2026-11-01.

## Problem

An orchestrator finds a comment defect after a chunk freeze and commits a comment-only Go correction. The chunk checkpoint accepts that correction, because the covered-digest check proves the comment gap (`internal/reviewrecord/coverage.go:150`). The next chunk then starts from the correction commit, and its checkpoint refuses with `stale review chain gap`. The review-chain check compares the next base digest with the previous chunk digest by equality (`coverage.go:88`), and the fold branch does the same (`internal/reviewrecord/fold.go:115`).

The only recovery today costs a review round. A refreeze keeps the verifications, because each entry binds to its own commit since FT445 (`internal/reviewrecord/check.go:141-151`). But a refreeze stales the three reviews, because each review binds to the frozen pair by equality (`internal/reviewrecord/record.go:158-159`). The other route moves the next base back to the frozen tip, and then the next chunk reviews the comment. This occurred at D-C4 on 2026-10-12 (commit `292326cc5`), at AP-C7 on 2026-10-17 (commits `e9a934b5d`, `1114a9aa4`), and at GP2 on 2026-10-21 (commits `1e3382dee`, `4e8cf2af9`).

The guidance also disagrees with the decision. The bounded repair policy says that a comment-only correction before the green checkpoint still takes a refreeze (`.agents/skills/bench-craft-line/references/bounded-repair-policy.md:49`). The implement phase says that each source change except a record commit or a recorded fold takes a refreeze (`.agents/commands/bench-implement-spec.md:60`). Commit `f17f4b6c9` (2026-10-19) wrote both sentences from the criterion `FT380-plan-commit.moved-tip`, which the reviewer admitted on 2026-10-18 (commit `66a3da81d`).

## Solution

The review-chain check accepts a later chunk base whose digest differs from the covered digest only by a proved comment gap. The proof is the existing `commentgap.Prove` that the covered-digest check uses, so its scope stays the same: Go files only, with equal token streams, and with no directive, cgo, or example-output comment. Each ancestry condition of the chain stays. A comment-only Go correction after a chunk freeze then takes no refreeze, keeps its verifications and its reviews, and keeps the review chain.

The guidance states the same rule. The bounded repair policy and the implement phase take no refreeze for a proved comment-only Go correction. The review phase admits a later chunk base after that correction. The anchor registry forbids the retired refreeze sentences and keeps its row on the unscoped wording of the rule.

## User stories

Line: opus / high.
Implementation-line reason: The hardest material chunk is CG-C2, because it changes guidance that steers every later build, and the leverage override routes guidance at high effort. The spec fixes every new sentence, every retired needle, and every row, so spec precision is exact. CG-C1 changes two comparisons at a known seam: the `follows` helper already holds the predicate, and the fold and covered-digest tests are the precedent. A wrong chain predicate turns the ordinary Go test phase red through planned rows that drive the real checkpoint, so test strength is high for CG-C1. The guidance rows are partly review-owned, which is weaker, and the anchor harness covers only the retired sentences.
Harder chunks: CG-C2.

### The review chain accepts a proved comment gap

1. As a build orchestrator, I want a later chunk base after a comment-only Go correction to pass the review chain, so that the correction takes no refreeze and no new review round.
2. As a build orchestrator, I want the frozen chunk to keep its reviews and its verifications after that correction, so that no recorded evidence goes stale.
3. As a build orchestrator, I want the chain to accept a comment gap after a recorded fold, so that a fold and a comment correction can both follow one chunk.
4. As a reviewer, I want the chain to refuse a base whose gap changes a Go statement, so that no unreviewed behavior enters the chain.
5. As a reviewer, I want the chain to refuse a gap that changes a directive comment, so that a comment with build behavior stays reviewed.
6. As a reviewer, I want the chain to refuse a gap in a non-Go file, so that the proof scope of the chain stays the scope of the covered-digest check.
7. As a reviewer, I want the chain to refuse a base with the gap tree on a parent line that skips the predecessor tip, so that the ancestry guard holds.
8. As a build orchestrator, I want the chain refusal to name the proved comment-only Go correction among the commits that follow a chunk tip, so that the refusal states the current rule.
9. As a build orchestrator, I want the completion checkpoint to accept the same chain, so that the final landing grades one rule.
10. As a kit maintainer, I want the chain, the fold path, and the covered-digest check to read one predicate, so that the proof cannot drift between them.
11. As a build orchestrator, I want a failed ancestry query to keep its refusal with no repair form, so that a broken record stays visible.

### The guidance states the rule

12. As a build orchestrator, I want the bounded repair policy to take no refreeze for a comment-only Go correction before the green checkpoint, so that the policy and the chain agree.
13. As a build orchestrator, I want the implement phase to list the proved comment-only Go correction among the changes that take no refreeze, so that the Land rules agree with the policy.
14. As a reviewer, I want the review phase to admit a later chunk base after a proved comment-only Go correction, so that the base rule agrees with the chain.
15. As a kit maintainer, I want the anchor registry to forbid the retired refreeze sentences, so that the guidance cannot return to the old rule.
16. As a kit maintainer, I want the anchor registry to keep forbidding the unscoped wording `comment-only corrections follow the chunk tip`, so that the guidance names the Go scope and the proof.
17. As a user of a linked repository, I want a changelog entry for the chain rule, so that I know why a comment correction no longer needs a refreeze.
18. As a later orchestrator, I want the Claude scorecard decisions to state no reason that this change makes false, so that a line declaration reads the current chain rule.

### Reviewed exclusions that stay in force

19. As a reviewer, I want a refreeze after a comment-only correction to keep staling the reviews, so that this change alters no review rule.
20. As a build orchestrator, I want the checkpoint after a comment-only Go correction to run the gate again, so that the checkpoint still grades its changed subject.

## Implementation decisions

- One unexported predicate in `internal/reviewrecord` answers whether a later digest follows a covered digest: the two digests are equal, or `commentgap.Prove` accepts the gap from the covered digest to the later digest. The `follows` helper, the equality branch of the chain in `checkSource`, `foldCoversBase`, and the covered-digest check of `checkSource` read that predicate. After the change, `commentgap.Prove` has one call site in `internal/reviewrecord`.
- The equality branch of the chain keeps its ancestry conjunct: the previous chunk tip is an ancestor of the next base. `foldCoversBase` keeps its conjunct: the fold commit is an ancestor of the base. The predicate replaces only the two digest comparisons.
- The chain keeps its order. It tries the previous chunk digest first, then the covering fold, then it refuses with the `FormFold` repair. The refusal sentence becomes: `plan commits land before the ticket merge, a default-branch merge lands before the first chunk or as a recorded fold, and only record commits, recorded folds, and proved comment-only Go corrections follow a chunk tip until its checkpoint; record the fold or review the uncovered delta`.
- The default base of a successor chunk stays the commit that `coveredAfter` answers for its predecessor (`internal/reviewrecord/ticket.go:173-188`). After the change, both that default and a later commit after a proved comment gap pass the chain.
- No record field, no verb grammar, and no refusal form changes. The `bench gate --checkpoint` route and the `bench record completion` route reach the chain through `CheckTrees` and `RecordCompletion` (`internal/reviewrecord/check.go:21`, `internal/reviewrecord/completion.go:37`). The premise census of build preflight reaches it through `CheckTrees` (`internal/preflight/ownership/evidence.go:49`). Each of them gains the rule with no edit.
- The guidance sentences are exact:
  - The bounded repair policy replaces its line 49 with `A comment-only Go correction before the green checkpoint stays evidence-only and takes no refreeze.` It then adds `The review chain accepts a later chunk base after that correction, so the chunk keeps its verifications and its reviews.` Lines 48, 50, and 51 stay.
  - The implement phase replaces the sentence `Each source change except a record commit or a recorded fold moves the chunk tip and takes a refreeze of that chunk.` on line 60 with one sentence: `Each source change except a record commit, a recorded fold, or a proved comment-only Go correction moves the chunk tip and takes a refreeze.` The sentence `This rule applies only before a chunk's green checkpoint and after the last chunk's freeze.` stays after it. The new sentence has 24 words, and the paragraph keeps its 6 sentences, which is the paragraph bound (`internal/prose/prose.go:17-18`). It stays on line 60, because the file is at its 97-line budget (`projects/benchkit.md:553`), and the `guidance-prose-budgets` check observes that budget.
  - The review phase adds `A commit that adds only a proved comment-only Go correction after that base can also be the chunk base.` after its sentence `When recorded folds follow the predecessor tip, the base is the last of them.` "That base" is then the predecessor tip or its last recorded fold. The sentence has 19 words, and the paragraph grows to 6 sentences.
  - The two sentences of the bounded repair policy have 14 and 20 words, and their paragraph grows to 6 sentences. Each of the three edited paragraphs is then at the 6-sentence bound, so the ticket adds no other sentence to them.
- The anchor registry changes in `internal/anchors/registry_chunk_chain.go` and its mirror `internal/anchors/registry_chunk_chain_test.go`:
  - It keeps the Forbid row with the needle `comment-only corrections follow the chunk tip` and its mirror row, because that row forbids an unscoped wording that names neither Go nor a proof. The new sentences do not hold that needle.
  - It adds a Forbid row on the bounded repair policy with the needle `takes a refreeze`. The needle also catches a rewrite that drops the word `still`. The diagnostic is `chunk chain: the bounded repair policy restored the refreeze of a comment-only correction`. The mirror row pins the retired sentence `A comment-only correction before the green checkpoint stays evidence-only and still takes a refreeze.`
  - It adds a Forbid row on the implement phase with the needle `Each source change except a record commit or a recorded fold moves the chunk tip`. The diagnostic is `chunk chain: the implement phase restored the refreeze of a comment-only correction`. The mirror row pins the retired sentence `Each source change except a record commit or a recorded fold moves the chunk tip and takes a refreeze of that chunk.`
  - The doc comment of `chunkChainAnchors` gains one short sentence that names the retired refreeze sentences. Its current sentence of about 80 words stays as it is.
- The anchor registry adds no Require row for the new guidance sentences. The admission test refuses a Require row without a label, route, invariant, or command class (`internal/anchors/locate_test.go:93-111`), and these sentences hold none of those.
- `CHANGELOG.md` gets one `### Review chain after a comment-only correction` section under `## [Unreleased]`. Its two old entries at lines 238 and 239 stay as history.
- The two scorecard decisions at `capture/agent-performance/claude-models.md:80-81` lose their false reasons. Line 80 no longer says that the review chain refuses any source commit between two chunks. Line 81 says that a refreeze stales the reviews, in place of the claim that it orphans every verification.

## Implementation chunks

CG-C1 ships the chain rule in Go with its tests. The guidance in CG-C2 states that rule, so CG-C1 lands first, and its review closes before CG-C2 starts.
CG-C2 changes the guidance, the anchor registry, the changelog, and the scorecard.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| CG-C1 / 01-chain-proves-comment-gap.md | The review chain, the fold path, and the covered-digest check read one predicate, so a later base after a proved comment gap passes, and every other gap still refuses | CG01, CG02, CG03, CG04, CG05, CG06, CG07, CG08, CG09, CG10, CG11, CG12, CG13, CG14, CG15, CG16, CG17 | TestChainAcceptsAProvedCommentGap, TestFoldChainAcceptsARecordedFoldBetweenChunks, TestChainRefusesAFailedAncestryQuery, TestReviewRecordSource, TestRecordCompletionAcceptsCommentOnlyGap, TestReviewCheckpointChainAcceptsACommentGap, TestReviewCheckpointChainGapNamesTheExpectedBase, TestReviewCheckpointCommentOnlyGap, TestReviewCheckpointKeepsStrictEvidence, TestDeltaCheckpointRunsWholeForANonTestDelta | no |
| CG-C2 / 02-guidance-states-comment-gap.md | The bounded repair policy, the implement phase, and the review phase state the rule, the anchor registry forbids the retired refreeze sentences, and the changelog and scorecard follow | CG18, CG19, CG20, CG21, CG22, CG23, CG24, CG25, CG26 | TestChunkChainAnchors, TestRegistryNeedlesDistinct | yes |

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft439-build-20261101","orchestrator_session":"session_01JK5CPcwpGroiAPWzzXxytc","author_limit":1,"assignments":{"01-chain-proves-comment-gap.md":[{"session":"ft439-author-01","assignment":"58abfa6badfad14b8718236a3ee8a691","model":"opus","effort":"medium","source":"1d08e3dd749824d76d22fa22e752cae34b53a944","native_ref":"claude-agent:bench-writer-medium:ft439-author-01"}],"02-guidance-states-comment-gap.md":[{"session":"ft439-author-02","assignment":"58abfa6badfad14b8718236a3ee8a691","model":"opus","effort":"high","source":"4815d5be3d1108024bb377f5b76ef982ca6eba22","native_ref":"claude-agent:bench-writer-high:ft439-author-02"}]}},"chunks":[{"id":"CG-C1","tickets":["01-chain-proves-comment-gap.md"],"verification":[{"id":"chain-gap","command":"bench test --package ./internal/reviewrecord/... --run 'Chain|Fold|ReviewRecordSource|CommentOnlyGap'","probe":"Make foldCoversBase compare the fold digest with the base digest by equality, and record that the fold-then-comment case of TestChainAcceptsAProvedCommentGap turns red.","ticket":"01-chain-proves-comment-gap.md"},{"id":"checkpoint","command":"bench test --package ./internal/gate --run 'ReviewCheckpoint|DeltaCheckpointRunsWholeForANonTestDelta'","probe":"Restore the equality comparison of the base digest with the previous chunk digest in checkSource, and record that TestReviewCheckpointChainAcceptsACommentGap turns red.","ticket":"01-chain-proves-comment-gap.md"}]},{"id":"CG-C2","tickets":["02-guidance-states-comment-gap.md"],"verification":[{"id":"anchors","command":"bench test --package ./internal/anchors","probe":"Delete the Forbid row with the needle 'takes a refreeze' from internal/anchors/registry_chunk_chain.go, and record that TestChunkChainAnchors turns red on its mirror row.","ticket":"02-guidance-states-comment-gap.md"},{"id":"guidance","command":"bench test --check docs-currency-workflow","probe":"Restore the sentence 'Each source change except a record commit or a recorded fold moves the chunk tip and takes a refreeze of that chunk.' in the Land section of the implement phase, and record that the docs-currency-workflow check turns red.","ticket":"02-guidance-states-comment-gap.md"}]}],"final_verification":[{"id":"owner-packages","command":"bench test --package ./internal/reviewrecord/... --package ./internal/commentgap --package ./internal/gate --package ./internal/anchors --package ./internal/landing --package ./internal/preflight/..."},{"id":"guidance","command":"bench test --check docs-currency-workflow"},{"id":"prose","command":"bench test --check prose-mechanics"},{"id":"coverage","command":"bench coverage --check comment-gap-chain"}]}
```

## Testing decisions

- A good chain test builds real commits with the `recordtest` fixture and calls `reviewrecord.CheckSource` for chunk 2. The comment gap edits one committed Go file between the chunk 1 tip and the chunk 2 base. The fold tests are the precedent: `foldedAfterChunk1` and `foldRepair` in `internal/reviewrecord/fold_test.go:16-33`. The new test `TestChainAcceptsAProvedCommentGap` lives in that file, which has 168 lines.
- A good checkpoint test runs `RunCommand` with `--checkpoint` over the real oracle fixture and reads the exit code and the retained record. `TestReviewCheckpointCommentOnlyGap` with `seedCommentGap` and `correctCommentGap` is the precedent (`internal/gate/review_checkpoint_commits_test.go:104-178`). The new test `TestReviewCheckpointChainAcceptsACommentGap` lives in that file, which has 246 lines.
- The `chain base` case of `TestReviewCheckpointKeepsStrictEvidence` changes its gap from a comment edit to a Go statement edit, so it keeps the `expected base` refusal. The comment-gap form of that case moves to `TestReviewCheckpointChainAcceptsACommentGap` as an acceptance.
- `TestReviewCheckpointChainGapNamesTheExpectedBase` pins the new rule clause of the refusal sentence.
- Each expected refusal and each expected exit code is a literal in the test. No expectation calls the predicate under test.
- The anchor rows extend `TestChunkChainAnchors` with one mirror row for each new Forbid row. The mirror row of the kept Forbid row stays. `TestChunkChainAnchors` evaluates only temporary trees (`internal/anchors/anchor_harness_test.go:52-54`), so its probe mutates the registry, and the guidance probe runs the `docs-currency-workflow` check on the tree. `TestRegistryNeedlesDistinct` grades the new needles (`internal/anchors/anchor_harness_diagnostics_test.go:224-253`).
- The ordinary Go test phase of the gate observes CG01 to CG11, CG13 to CG17, and CG21 to CG23. The `docs-currency-workflow` check also observes CG21 to CG23, because it evaluates the `AfterImplementSpec` anchor group on the tree (`internal/conformance/docs_workflow_helpers_test.go:36`). CG12, CG18 to CG20, and CG24 to CG26 are review-owned.

### Seam diagram

    trigger: bench gate --checkpoint <spec> (--chunk <id> | --complete), bench record completion, the preflight premise census
        │
        ▼
    record, tree, tip  ──▶  [ reviewrecord.checkSource: chain loop over the chunks ]  ──▶  nil or Repair{Form: FormFold}
                               │  base digest vs previous chunk digest ──▶ [ follows predicate ] ◀── commentgap.Prove
                               │  base digest vs covering fold digest   ──▶ [ follows predicate ]
                               │  requested digest vs covered digest    ──▶ [ follows predicate ]
                               ◀ tests attach here: fold_test.go calls CheckSource for chunk 2;
                                 review_checkpoint_commits_test.go runs RunCommand --checkpoint and reads the exit code and the record

    trigger: an edit of the guidance
        │
        ▼
    bounded-repair-policy.md, bench-implement-spec.md  ──▶  [ anchors.EvaluateGroup(AfterImplementSpec) ]  ──▶  chunk-chain diagnostics
                               ◀ tests attach here: TestChunkChainAnchors restores each retired sentence in a copy

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CG01 | 1 | `bench gate --checkpoint <spec> --chunk 2` exits 0 when the chunk 2 base is a commit after the chunk 1 tip whose only change is a Go comment edit | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointChainAcceptsACommentGap`) | The equality branch refuses that base with `stale review chain gap`, so the checkpoint exits 1 |
| CG02 | 2 | After that green checkpoint, the retained chunk 1 entry has the frozen Base, Tip, SourceDigest, Reviews, and Verification | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointChainAcceptsACommentGap`) | A route that refreezes chunk 1 changes the Tip and the SourceDigest, and the test compares each field with the frozen entry |
| CG03 | 9 | `bench gate --checkpoint <spec> --complete` exits 0 over two chunks when the chunk 2 base follows the chunk 1 tip by a Go comment edit | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointChainAcceptsACommentGap`) | The complete check walks every chunk through the same chain, so an equality branch refuses chunk 2 |
| CG04 | 1 | `reviewrecord.CheckSource` for chunk 2 returns nil when the chunk 2 base follows the chunk 1 tip by a Go comment edit | `internal/reviewrecord/fold_test.go` (`TestChainAcceptsAProvedCommentGap`) | The equality branch returns the `FormFold` repair for that base |
| CG05 | 3 | `reviewrecord.CheckSource` for chunk 2 returns nil when the chunk 2 base follows a recorded fold after chunk 1 by a Go comment edit | `internal/reviewrecord/fold_test.go` (`TestChainAcceptsAProvedCommentGap`) | `foldCoversBase` compares the fold digest by equality, so the chain refuses that base |
| CG06 | 4 | `reviewrecord.CheckSource` refuses with the `FormFold` repair and `stale review chain gap` when the gap changes a Go statement | `internal/reviewrecord/fold_test.go` (`TestChainAcceptsAProvedCommentGap`) | A predicate that accepts each descendant base passes this case |
| CG07 | 5 | `reviewrecord.CheckSource` refuses with `stale review chain gap` when the gap changes a comment into a `//go:generate` directive | `internal/reviewrecord/fold_test.go` (`TestChainAcceptsAProvedCommentGap`) | A predicate that compares the files without comments, in place of `commentgap.Prove`, passes this case |
| CG08 | 6 | `reviewrecord.CheckSource` refuses with `stale review chain gap` when the gap edits only a comment line in a committed Markdown file | `internal/reviewrecord/fold_test.go` (`TestChainAcceptsAProvedCommentGap`) | A predicate that accepts a change in a file that is not Go passes this case |
| CG09 | 7 | `reviewrecord.CheckSource` refuses with `stale review chain gap` when the chunk 2 base holds the comment-gap tree on a parent line that skips the chunk 1 tip, and the chunk 2 tip is a two-parent commit of that base and the chunk 1 tip | `internal/reviewrecord/fold_test.go` (`TestChainAcceptsAProvedCommentGap`) | The two-parent tip passes the pair guard at `coverage.go:45-53`, so only the chain ancestry conjunct refuses, and a change that drops that conjunct passes this case |
| CG10 | 4 | The `chain base` case of the strict checkpoint test refuses with `expected base` when the gap changes a Go statement | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A checkpoint that accepts each descendant base exits 0 on this case |
| CG11 | 8 | The chain refusal holds `only record commits, recorded folds, and proved comment-only Go corrections follow a chunk tip until its checkpoint` | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointChainGapNamesTheExpectedBase`) | The old refusal sentence lacks that clause, so the substring test fails |
| CG12 | 10 | `commentgap.Prove` has one call site in `internal/reviewrecord`, and that call site is the shared predicate | review-owned: `rg -n 'commentgap\.Prove\(' internal/reviewrecord` answers one line | A second inline proof can keep its own scope when the predicate changes |
| CG13 | 10 | The covered-digest check keeps its acceptance of a Go comment gap and its refusal of a statement gap | `internal/reviewrecord/recordcmd/completion_test.go` (`TestRecordCompletionAcceptsCommentOnlyGap`) | A move of the covered-digest check onto a predicate with the wrong scope turns the comment case or the statement case red |
| CG14 | 11 | A failed ancestry query of the chunk base or of the fold commit refuses with no repair form | `internal/reviewrecord/fold_test.go` (`TestChainRefusesAFailedAncestryQuery`) | In the fold-commit case, a fold predicate that answers before the ancestry query hides the failed query; the pair guard catches the chunk-base case first |
| CG15 | 7 | The unrecorded-fold case and the off-fold case keep the `FormFold` refusal | `internal/reviewrecord/fold_test.go` (`TestFoldChainAcceptsARecordedFoldBetweenChunks`) | A predicate that accepts the fold merge delta, or a change that drops the fold ancestry conjunct, passes these cases |
| CG16 | 19 | A refrozen chunk with its reviews at the earlier pair refuses with `chunk 1: stale Standards` and routes to the confirmation form | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A change that also relaxes the review pair check exits 0 on the `moved chunk` case |
| CG17 | 20 | A delta checkpoint after a comment-only Go source commit runs the whole gate | `internal/gate/delta_checkpoint_test.go` (`TestDeltaCheckpointRunsWholeForANonTestDelta`) | A change that treats a proved comment gap as an unchanged subject skips the whole run |
| CG18 | 12 | The bounded repair policy holds `A comment-only Go correction before the green checkpoint stays evidence-only and takes no refreeze.` | review-owned: the sentence in `.agents/skills/bench-craft-line/references/bounded-repair-policy.md` | The old sentence names a refreeze, and the reviewer reads the replaced line |
| CG19 | 13 | The Land section of the implement phase holds `Each source change except a record commit, a recorded fold, or a proved comment-only Go correction moves the chunk tip and takes a refreeze.` | review-owned: the sentence in `.agents/commands/bench-implement-spec.md` | The old sentence lists only record commits and recorded folds, and the reviewer reads the replaced line |
| CG20 | 14 | The review phase holds `A commit that adds only a proved comment-only Go correction after that base can also be the chunk base.` | review-owned: the sentence in `.agents/commands/bench-review-implementation.md` | The old base rule admits only the predecessor tip, and the reviewer reads the added sentence |
| CG21 | 15 | The anchor registry reds the bounded repair policy when it holds `takes a refreeze` | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | Without the Forbid row, the mirror row that restores the retired sentence raises no diagnostic |
| CG22 | 15 | The anchor registry reds the implement phase when it holds `Each source change except a record commit or a recorded fold moves the chunk tip` | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | Without the Forbid row, the mirror row that restores the retired sentence raises no diagnostic |
| CG23 | 16 | The anchor registry reds the implement phase when it holds `comment-only corrections follow the chunk tip` | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A removal of the Forbid row leaves its mirror row with no diagnostic, so the unscoped wording can return |
| CG24 | 17 | `CHANGELOG.md` holds one `### Review chain after a comment-only correction` section under `## [Unreleased]` | review-owned: the section in `CHANGELOG.md` | A missing entry leaves a linked repository with no record of the changed rule |
| CG25 | 18 | The scorecard decision on chunk bases names no review-chain refusal of each source commit between two chunks | review-owned: line 80 of `capture/agent-performance/claude-models.md` | The kept reason is false after CG-C1 |
| CG26 | 18 | The scorecard decision on comment-only corrections gives the review staleness of a refreeze as its reason | review-owned: line 81 of `capture/agent-performance/claude-models.md` | The kept reason names orphaned verifications, which FT445 made false |

### Edge inventory

The edge classes at the chain predicate, with their dispositions:

- A gap that changes a Go statement: CG06 and CG10.
- A gap that adds or changes a directive comment: CG07. The cgo and example-output refusals of `commentgap.Prove` carry over with no new row, because `TestProveGoGap` grades them (`internal/commentgap/commentgap_test.go:8-45`).
- A gap in a file that is not Go: CG08.
- A gap that adds, deletes, or renames a file, or changes a mode, a symlink, or a gitlink: the tree rules of `commentgap.Prove` refuse each one, and `TestProveTreeGap` grades them (`internal/commentgap/trees_test.go:12-78`).
- A base with the gap tree off the predecessor line: CG09 and CG15.
- A comment gap after a recorded fold: CG05.
- A failed ancestry query: CG14.
- A record-only commit in the gap: the source digest removes the record file, so the digests are equal and the chain accepts them as today (`internal/reviewrecord/source_test.go:43-53`).
- A comment gap before the first chunk: not applicable, because the first chunk has no predecessor, and the chain runs only from the second chunk.

The profile's hostile-input classes, walked:

- Paths with spaces or glob characters, control bytes in git text, and quoted diff-header paths: not applicable, because the predicate compares tree object IDs, and `commentgap.Prove` reads tree changes, not patch headers.
- A file whose last line lacks a newline: `commentgap.Prove` compares token streams, so a final newline edit is a layout change that it accepts.
- An absent file versus an empty file: an added or deleted file refuses through the tree rules above.
- Special files and dangling symlinks: a symlink or a gitlink refuses through the tree rules above. The predicate reads no path from disk.
- A command whose write changes the facts it reports: not applicable, because the chain check writes no record.

**Won't handle** lines:

- A comment-only correction in a file that is not Go, such as a shell file — the decision source scopes the proof to the covered-digest check, which is Go-only, and the open learning on a comment-only finding in a shell file stays uncommitted intake in `capture/learnings.md`; the surviving in-scope caller is the covered-digest check, which refuses the same gap.
- A refreeze after a comment-only correction keeps staling the reviews — the decision routes no refreeze, and `checkReviews` still guards a moved pair for every other change (CG16).
- A layout-only Go change counts as a proved comment gap — `commentgap.Prove` already accepts it at the covered-digest check (`internal/commentgap/commentgap_test.go:15`), and the chain reads the same proof.
- The default successor base stays the predecessor's covered commit — the chain accepts that base as today, and `chunkBase` keeps its one caller (`internal/reviewrecord/ticket.go:173-188`).
- The delivered criterion `FT380-plan-commit.moved-tip` keeps its text — it is delivered history of a closed outcome (`.bench/commitment.json:966-967`, delivered at `.bench/commitment.json:2143-2149`), and the FT439 criterion governs the comment-only correction from this landing.

## Out of scope

- A comment-only correction in a file that is not Go takes the evidence-only route. This needs a proof for each file kind and its own spec: about 6 edits, 2 gate runs.

## Further notes

### Seams beyond the decision source

The decision source names the review-chain check and the bounded repair policy. The spec adds these seams, each with its citation and research finding:

- `foldCoversBase` in the fold branch of the chain (`internal/reviewrecord/fold.go:113-119`): the Q1 readers found that it compares the fold digest by equality, on the same chain path as `coverage.go:88`.
- The covered-digest check moves onto the shared predicate (`internal/reviewrecord/coverage.go:150`): the Q1 readers found that it states the same proof inline, and the one-source rule collapses the copies.
- The implement phase sentence (`.agents/commands/bench-implement-spec.md:60`): the Q2 readers found that it requires a refreeze for each source change except a record commit or a recorded fold.
- The review phase base rule (`.agents/commands/bench-review-implementation.md:99`): the Q2 readers found that it admits only the predecessor tip.
- The anchor registry rows (`internal/anchors/registry_chunk_chain.go:20`, `internal/anchors/registry_chunk_chain_test.go:22`): the Q2 readers found a Forbid row on an unscoped wording of the rule that FT439 restores. The spec keeps that row and adds two Forbid rows in the same registry.
- The scorecard decisions (`capture/agent-performance/claude-models.md:80-81`): the Q2 reader A found that both reasons become false.

### Flagged additions

- The covered-digest check moves onto the shared predicate. Its behavior stays, and CG13 guards it.
- The registry adds two Forbid rows for the retired refreeze sentences, with two mirror rows.
- The bounded repair policy names `Go` in its line 49, which narrows the old sentence to the scope that the checkpoint proves. Lines 48 and 50 already name `Go`.
- The scorecard decisions at lines 80 and 81 change their reasons.

### Source-sentence-to-row table

| source sentence (`roadmap/FT439.md`) | rows |
| --- | --- |
| The covered-digest check accepts a proved comment gap. | CG13 |
| The review-chain check compares the next chunk base with the previous chunk digest, or with a covering fold digest, by equality. | CG04, CG05 |
| A comment gap therefore breaks the chain and forces a new review round. | CG01, CG03 |
| The bounded repair policy still requires a refreeze for a comment-only correction before the green checkpoint. | CG18, CG21 |
| Since FT445, each verification entry binds to its own commit, so a refreeze keeps the recorded verifications of the chunk. | CG02, CG26 |
| The review-chain check accepts a proved comment gap, as the covered-digest check does. | CG01, CG03, CG04, CG05, CG06, CG07, CG08, CG09, CG11, CG12 |
| The bounded repair policy takes no refreeze for it. | CG18, CG19, CG20, CG22, CG23 |
| This decision supersedes the 2026-10-18 rule for a comment-only correction. | CG19, CG22, CG23, CG24 |
| Criterion `FT439-comment-correction.done`: A comment-only correction after a chunk freeze keeps its verifications and its review chain. | CG01, CG02 |

### Pre-review proof checklist

- `Cited symbols`: `checkSource`, `CheckSource`, `CheckTrees`, `RecordCompletion`, `follows`, `coveringFold`, `coveredAfter`, `foldCoversBase`, `chunkBase`, `commentgap.Prove`, `checkReviews`, `EntryCurrent`, `chunkChainAnchors`, `admissionFaults`, and each test name in the map resolve in the tree at `4f9637904`.
- `Import edges`: none. `internal/reviewrecord` already imports `internal/commentgap` (`internal/reviewrecord/coverage.go:9`, `internal/reviewrecord/fold.go:7`).
- `Source-row clauses and occurrences`: the table above quotes each clause. The occurrences are D-C4 (`292326cc5`), AP-C7 (`e9a934b5d`, `1114a9aa4`), and GP2 (`1e3382dee`, `4e8cf2af9`).
- `Promised field labels`: none.
- `Changed-function callers`: `bench consumers` at `4f9637904` answers `checkSource` from `CheckTrees` (`check.go:21`), `RecordCompletion` (`completion.go:37`), and `CheckSource` (`coverage.go:15`); `foldCoversBase` from `checkSource` (`coverage.go:94`); `follows` from `RecordFold` (`fold.go:64`) and `coveringFold` (`fold.go:93`). The test callers of `CheckSource` are in `fold_test.go` and `source_test.go`.
- `Copy survival`: CG12 fails if an inline `commentgap.Prove` call survives beside the shared predicate.
- `Rendered-shape readers`: the old needle `only record commits and recorded folds follow a chunk tip until its checkpoint` has three hits: `internal/reviewrecord/coverage.go:99`, `internal/gate/review_checkpoint_commits_test.go:99`, and `CHANGELOG.md:238`. Ticket 01 writes the first two, and the changelog line stays as history. `dist/bench` embeds the old text, and a rebuild replaces it.
- `Pin operators`: CG11 uses `strings.Contains` on the refusal text. CG21 and CG22 use the Forbid substring match of `anchors.EvaluateGroup`.
- `Entry reads`: none.
- `Derived expectations`: none. Each expectation is a literal.
- `Consolidated rules`: the three digest comparisons consolidate onto one predicate.

  | site | old rule | new rule | row |
  | --- | --- | --- | --- |
  | chain branch, `coverage.go:88` | equal digests | equal or proved comment gap | CG04 |
  | fold branch, `fold.go:115` | equal digests | equal or proved comment gap | CG05 |
  | covered-digest check, `coverage.go:150` | equal or proved comment gap, inline | the same rule through the predicate | CG13 |
  | `follows`, `fold.go:138-141` | proved comment gap, which includes equal trees | the same rule through the predicate | CG15 |

- `Quantified obligations`: "each ancestry condition stays" covers two conjuncts, `coverage.go:89` and `fold.go:118`, which CG09 and CG15 grade. "The guidance states the same rule" covers three files, which CG18, CG19, and CG20 grade.
- `Workflow-step writes`: none. The change adds no workflow step.

### Facts that differ from the research inputs

- The anchor admission test refuses a Require row without a label, route, invariant, or command class (`internal/anchors/locate_test.go:93-111`), so no Require row pins the new guidance sentences. The two Forbid rows for the retired sentences guard the change.
