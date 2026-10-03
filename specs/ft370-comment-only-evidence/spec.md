# A comment-only correction takes the evidence-only path

Status: staged

Roadmap: FT370

Decision source: named reviewed artifact — roadmap/FT370.md (rule approved by the reviewer 2026-10-02; drain d-4be30106ede9).

Verification log: 2 iteration(s) to accept — spec iteration 1 (Fable xhigh) returned S1 to S9. The coordinator folded eight, S1 through a reviewer decision, and refuted S6. Spec iteration 2 (Sonnet high) confirmed each fold and returned N1, which CG44 and CG45 fold. The tickets pass (Fable xhigh) returned four prose findings, K1 to K4, and the coordinator folded them. The reviewer waived the Sonnet tickets pass because Fable found no blocking defect.

## Problem

A chunk review binds its results to the reviewed chunk tip through the source
digest. The source digest is the Git tree of that tip without the review
record, so every byte of every file is part of it. A one-word comment fix
after the review changes the digest. The checkpoint then refuses with
`stale reviewed source`, and the orchestrator re-records the chunk. That
re-record makes all three axes stale, so the comment fix costs a confirming
round of three axes and one of the two repair cycles.

The worktree-verb-runner build paid this cost in chunk VR-C4. Commit
`c9731327` changed one comment line in a Go test file and no executable line.
The repair took the second and last repair cycle of the chunk, two
verification reruns, and a fresh round of all three axes.

## Solution

`bench gate --checkpoint` accepts a comment-only gap between the reviewed
chunk tip and the graded source. The checkpoint proves the gap itself from the
two Git trees. Each changed path must be an existing regular Go file whose Go
tokens stay the same and which holds no directive-shaped comment. Any other gap keeps the
current refusal and its current message.

The orchestrator commits the correction, and the plan records no assignment
for it, so the plan digest stays the same. The chunk entry keeps its reviewed
frozen pair, so the three axis results and the recorded verification stay
current. Only the issuing axis reaffirms the
correction, with one superseding pass result at that pair. The bounded repair
policy gains one sentence that makes such a correction evidence-only, so it
consumes no repair cycle.

## User stories

Line: current Codex session / high.

The reviewer directs this session to implement all four tickets on 2026-10-03.
The runtime does not expose its exact model identifier.
The run keeps the version 1 completion plan because one session owns every ticket.
Three fresh Sol 6.1 sessions review each chunk at high effort.
The implementation continues without an iteration cap while progress holds.
The bounded repair policy still limits each chunk to two repair cycles.

Implementation-line reason: CG-C2 grades hostile Go inputs at an oracle seam.
The checkpoint and guidance require high effort.

Harder chunks: CG-C2.

### The checkpoint accepts a comment-only gap

1. As an orchestrator, I want a chunk checkpoint to pass a comment-only Go gap, so that a comment fix keeps all three reviews.
2. As an orchestrator, I want the `--complete` checkpoint to pass the same gap, so that final acceptance needs no new chunk review.
3. As a landing broker, I want the prospective completion checkpoint to pass the same gap, so that the landing applies one evidence rule.
4. As an orchestrator, I want an accepted gap to still run the whole gate, so that the gate grades the corrected source.
5. As the issuing axis, I want to reaffirm a correction at the reviewed frozen pair, so that the other two axes stay current.

### The checkpoint fails closed

6. As a reviewer, I want an unproven gap to keep the current refusal and message, so that every executable change still takes review.
7. As a reviewer, I want the checkpoint to classify the graded tree, so that a committed comment fix cannot hide an uncommitted code change.
8. As a reviewer, I want the reviewed tree to come from the chunk tip commit, so that a forged record digest authorizes nothing.

### The Go rule proves one file

9. As an orchestrator, I want a trailing comment edit on a code line to qualify, so that the common comment fix takes this path.
10. As an orchestrator, I want a comment edit on a last line without a newline to qualify, so that a hand-edited file end qualifies.
11. As an orchestrator, I want a whitespace-only layout change in a Go file to qualify, because it leaves every Go token the same.
12. As a reviewer, I want a changed raw-string literal to refuse, so that data that looks like a comment still takes review.
13. As a reviewer, I want a block comment that gains a newline after an operand to refuse, because the newline adds a semicolon.
14. As a reviewer, I want a scanner error on either side to refuse, so that an unterminated block comment cannot pass with equal tokens.
15. As a reviewer, I want any changed file with a directive-shaped comment to refuse, so that a `//go:build`, `//go:embed`, or `//export` edit still takes review.
16. As a reviewer, I want a changed `// +build` line to refuse, so that a legacy build constraint still takes review.
17. As a reviewer, I want a file that imports `"C"` to refuse, because its preamble comment is C source.
18. As a reviewer, I want an example-output comment on either side of a `_test.go` file to refuse, because `go test` compares it.

### The tree rule proves the change set

19. As a reviewer, I want a changed path without the `.go` suffix to refuse, so that shell, Markdown, and data files keep their review.
20. As a reviewer, I want an added, deleted, or renamed path to refuse, so that only an edit of an existing Go file qualifies.
21. As a reviewer, I want a mode change to refuse even when the bytes match, so that an executable-bit change still takes review.
22. As a reviewer, I want a symbolic link or a gitlink entry to refuse, so that no link is read as Go source.
23. As a reviewer, I want a Go blob over the control-record read bound to refuse, so that the classifier reads only bounded bytes.
24. As a reviewer, I want two different trees with an empty change list to refuse, so that a hidden entry cannot pass.
25. As an orchestrator, I want a Go path with a space, a tab, or a glob character to qualify, so that legal names work.
26. As a reviewer, I want a gap of several Go files to qualify only when each file qualifies, so that one code edit refuses.

### The surrounding evidence rules stay strict

27. As a reviewer, I want completion evidence recorded before a comment-only correction to stay stale, so that final verification always names the final source.
28. As a reviewer, I want each later chunk base to keep its exact predecessor tree, so that a correction joins the next delta.
29. As an orchestrator, I want the chain-gap refusal to name comment-only corrections beside record commits, so that its rule text agrees with the implementation phase.
30. As a reviewer, I want each axis result to keep its chunk entry's frozen pair, so that a re-recorded chunk stales those results.

### The guidance states the rule

31. As an orchestrator, I want the bounded repair policy to make a proven comment-only correction evidence-only, so that it costs no repair cycle.
32. As a maintainer, I want an anchor and an independent expectation to pin that policy sentence, so that its removal turns the gate red.
33. As an orchestrator, I want the implementation phase to let comment-only corrections follow the chunk tip, so that it matches the checkpoint.

### One raw tree-change reader

34. As a maintainer, I want the lane and the classifier to share one raw tree-change reader, so that the framing has one source.
35. As a maintainer, I want the lane's change list to stay the same after the reader moves, so that lane selection keeps its behavior.
36. As a maintainer, I want the shared reader to compare two tree IDs that no commit names, so that the classifier compares record-excluded trees.

### The orchestrator records a correction

37. As an orchestrator, I want to commit a comment-only correction with no plan assignment, so that the plan digest stays the same.
38. As an orchestrator, I want `bench record completion` to accept the same gap, so that I can record completion after a last-chunk correction.

### Reviewed exclusions

39. As a light-path author, I want a light-path comment fix to keep the current repair accounting, because no checkpoint proves its gap.

40. As a reviewer, I want a hidden gitlink beside a Go comment edit to refuse, so that a configured filter cannot authorize it.

## Implementation decisions

### The classifier seam

A new package, `internal/commentgap`, owns the comment-only verdict. Its one
export is `Prove(root, reviewed, later string) error`. The two operands are
Git tree IDs. `Prove` returns nil only when every difference between the two
trees is a comment-only Go change. Otherwise it returns an error that wraps
one rule sentinel and names the first path that the rule refused.

The package is new because `internal/reviewrecord`, `internal/gate`, and
`internal/git` are crowded directories in `bench structure`. FT358 ticket 7
also adds no file to `internal/reviewrecord`. `internal/reviewrecord` cannot
import `internal/gate`, because `internal/gate` imports it. `internal/git`
imports only `internal/bounds` and `internal/canonicalpath` (observed with
`go list -deps`), so the new edges form no cycle.

### The shared raw tree-change reader

`gate.ComposedChanges` owns the one `git diff --raw --no-renames -z` reader
today, in `internal/gate/lane_select.go`. The build moves that reader to
`internal/git` as `TreeChanges(root, from, to string) ([]TreeChange, error)`.
`TreeChange` keeps the four fields `Status`, `SrcMode`, `DstMode`, and `Path`.
`gate.ComposedChange` becomes an alias of `git.TreeChange`, and
`gate.ComposedChanges` calls `git.TreeChanges` with `base^{tree}`. The gate
wrapper keeps its `gate: composed change list unavailable` text. No test
matches the moved entry-parse error text.

A probe on git 2.43.0 confirmed four facts that the reader relies on. The
`-z` framing keeps a name with a space, a tab, or `*` byte for byte. A
mode-only change arrives as status `M` with two modes and one blob. A tree ID
operand peels through `^{tree}` to itself. `git ls-tree` reads `a*.go` as a
literal path.

The configured reader can hide a gitlink change under `diff.ignoreSubmodules=all`.
The classifier uses `git.TreeChangesIncludingSubmodules` to keep every gitlink visible.
This wrapper shares the raw reader and parser, with `--ignore-submodules=none`.
The lane retains its configured reader.
CG30 refuses a hidden gitlink with `ErrMode`.
CG46 covers a hidden gitlink beside a visible Go comment edit.

### The tree rule

`Prove` applies these rules in this order. The first rule that fails wins.

1. Equal tree IDs return nil.
2. Read the complete list through `git.TreeChangesIncludingSubmodules`. A reader failure wraps `ErrUnreadable`.
3. An empty change list between two different trees wraps `ErrEmptyChanges`.
4. Each change takes the remaining rules in Git path order.
5. A status other than `M` wraps `ErrStatus`. Rename detection is off, so a rename arrives as `D` and `A`.
6. Two different modes, or a mode other than `100644` or `100755`, wrap `ErrMode`.
7. A path without the case-sensitive `.go` suffix wraps `ErrNotGo`.
8. A failed bounded read of either blob through `git.ReadTreeFile` wraps `ErrUnreadable`.
9. The Go rule grades the two blobs.

### The Go rule

The Go rule takes the path and the two blobs. It applies these rules in this
order.

1. Each blob scans through `go/scanner` with comments skipped. Any scanner error on either side wraps `ErrScan`.
2. Two token lists that differ in a token kind or a literal wrap `ErrTokens`. Positions do not take part. The automatic semicolon is a token, so a newline that a block comment gains or loses can change the list.
3. Each blob scans again with `scanner.ScanComments`. A directive-shaped comment on either side wraps `ErrDirective`, even when its text is unchanged.
4. A file whose import list holds the path `C` wraps `ErrCgo`. `go/parser` with `ImportsOnly` reads the list, and a parse failure wraps `ErrScan`.
5. A `_test.go` path with an example-output comment on either side wraps `ErrExampleOutput`.

A directive-shaped comment is one of these:

- a `//` comment whose third byte exists and is not a space or a tab
- a block comment that starts with `/*line `
- a comment that `go/build/constraint.IsPlusBuild` accepts

The first shape is a superset of Go's `//[a-z0-9]+:[a-z0-9]` directive
grammar and of `//line`, `//extern`, and `//export`. `go/ast` exports no
directive parser in Go 1.25, so the rule names the conservative shape and does
not copy an unexported standard-library function.

An example-output comment is a comment whose text, without its markers and
leading whitespace, starts with `output:` or `unordered output:` in any
letter case. `strings.TrimSpace` removes the whitespace, including block-comment newlines.

A whitespace-only layout change passes the token comparison, because every token stays
the same. That widening is flagged decision F1 below.

### The checkpoint integration

`reviewrecord.checkSource` changes at one place: the comparison of the last
matched chunk's source digest with the graded source digest. When the two
digests differ, the checkpoint calls `commentgap.Prove` with the chunk's
source digest and the graded source digest. A nil answer accepts the gap. A
refusal keeps the current error text, `chunk <id>: stale reviewed source: no
chunk review covers <range>`, with its current suffix and range rule. The
checkpoint does not print the rule sentinel.

The earlier per-chunk check recomputes the digest from the chunk tip commit
and compares it with the recorded digest. Only after that check is the
recorded digest safe input, and the integration keeps that order. Both operands are tree objects
that `git.TreeWithoutFile` wrote, so no commit names them (story 36).

The comparison serves `--chunk`, `--complete`, and the landing broker's
prospective completion through `gate.WithCompletion`. The broker passes the
completion source tree as the source tree, and the gate proves the composed
tree against it before this comparison. So the landing accepts the same gap.
`RecordCompletion` also calls `checkSource`, so `bench record completion`
accepts the same gap (CG44) and refuses an unproven gap (CG45).

### The staleness sites

| site | comment-only gap | why |
| --- | --- | --- |
| last chunk digest against the graded digest | passes after the proof | the one site that FT370 changes |
| recorded chunk digest against the digest recomputed from the chunk tip | stays strict | it authenticates the reviewed tree that the proof reads |
| predecessor tip against the next chunk base (chain gap) | stays strict | the next chunk base stays the predecessor tip, so a correction joins the next chunk's delta |
| plan digest | stays strict | the plan digest comparison refuses a spec or ticket edit before the gap comparison runs, and `TestReviewRecordSource` pins it |
| each axis result against its chunk entry's frozen pair | stays strict | the chunk entry keeps the reviewed pair, so the axes stay current without a widening |
| chunk verification against the chunk digest | stays strict | the chunk entry keeps the reviewed digest, so the recorded verification stays current |
| completion evidence and final verification against the graded digest | stays strict | final verification runs after the last correction |
| `bench record verification` source against the chunk source | stays strict | a comment-only correction records no verification |

The chain-gap refusal text changes from `only record commits follow a chunk
tip` to `only record commits and comment-only corrections follow a chunk tip`.
The rest of that message stays.

### The record sequence of a correction

The orchestrator makes the comment edit and commits it on a lane pass. That
authorship is a closed reviewer decision (see Further notes). The plan records
no assignment for the correction, so the plan digest stays the same. A
recorded repair assignment edits the plan in `spec.md`. That edit moves the
plan digest and puts `spec.md` into the gap, so the path never activates.

The orchestrator does not re-record the chunk, and it records no verification
entry. The issuing axis records one superseding pass result at the chunk's
reviewed frozen pair, with `bench record review`. `RecordReview` takes the base, the
tip, and the digest from the chunk entry and does not read HEAD. So the pass
records while HEAD sits at the corrected tip. The checkpoint then proves the
gap, and the gate grades the corrected tree.

If the checkpoint refuses the gap, the correction is not evidence-only. The
finding then takes the normal repair route of `.bench/BENCH.md`. The
re-recorded chunk tip follows the orchestrator's commit, so the repair delta
that the confirming round reads spans that commit.

### The guidance

The bounded repair policy gains one sentence in `Classification and
completion`, directly after the review-record prose definition:

`A comment-only Go correction that the checkpoint accepts is evidence-only despite its source change, and the orchestrator commits it with no plan assignment.`

The phrase `despite its source change` marks the one exception to the
definition before it, which requires an unchanged source identity. The
existing sentences after it then apply unchanged: an evidence-only
correction consumes no repair cycle, and only the issuing axis reaffirms it.
The policy file has no budget row. `bench-review-implementation.md` line 231
already routes a repair through the narrow evidence-only exception, so that
file does not change.

Three unchanged sentences speak of a repair, and none of them applies to an
accepted correction, because the policy makes it evidence-only:

- `.bench/BENCH.md` line 127 sends a post-review repair to a fresh repair session with a plan assignment. It also bars the orchestrator from a repair at final reconciliation.
- `bench-implement-spec.md` line 56 has a repair session retain its verification at the final source.
- `bench-review-implementation.md` line 47 gives a chunk that ends on a repair one confirming round of three axes.

A refused gap is a repair, so those three sentences then apply as written.
`bench-implement-spec.md` line 42 governs a prose-only owner edit, so it does
not apply to a Go comment.

One anchor row pins the new sentence in
`internal/anchors/registry_retained_workflow.go`, beside the three
evidence-only rows. Its diagnostic is `implementation continuation: bounded
repair dropped comment-only evidence`. `TestImplementationContinuation` holds
one independent expectation for it. The expectation follows the exception in
`AGENTS.md`: its named omission is the deletion of the registry row, and the
ticket records that red (row CG38).

`bench-implement-spec.md` line 54 changes one sentence in place, so the file
stays at 80 of its 81 lines. `Only record commits follow the chunk tip.`
becomes `Only record commits and comment-only corrections follow the chunk
tip.` The chunk-chain anchor needle and its expectation in
`TestChunkChainAnchors` take the new text. Their diagnostic becomes `chunk
chain: only record and comment-only commits follow the chunk tip`.

### Bootstrap authority

The checkpoint runs inside the gate binary, which the gate owner builds from
the graded tree. The checkpoint recomputes the reviewed tree from the chunk
tip commit and requires it to equal the recorded digest. It computes the
graded digest from the graded tree. `Prove` reads only those two tree objects
through Git. No record field, author claim, or digest string authorizes the
gap by itself (rows CG7 and CG8). The gate binary that the gate owner builds
from the graded tree is the existing trust root, and FT370 adds no new trust
assumption.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| CG-C1 / `1-move-tree-change-reader.md` | One raw tree-change reader in `internal/git` serves the lane, and the lane's change list stays the same. | CG40, CG41 | `bench test --package ./internal/git`, `bench test --package ./internal/gate` | no |
| CG-C2 / `2-prove-comment-only-gaps.md` | `commentgap.Prove` proves a comment-only Go gap between two trees and names the rule of each refusal. | CG9, CG10, CG11, CG12, CG13, CG14, CG15, CG16, CG17, CG18, CG19, CG20, CG21, CG22, CG23, CG24, CG25, CG26, CG27, CG28, CG29, CG30, CG31, CG32, CG42, CG43, CG46, CG47, CG48, CG49, CG50 | `bench test --package ./internal/commentgap` | yes |
| CG-C3 / `3-accept-proven-gaps-at-checkpoint.md`, `4-state-comment-only-correction-rule.md` | The checkpoint, the completion record, and the landing accept a proven gap, the surrounding evidence stays strict, and the guidance routes the correction. | CG1, CG2, CG3, CG4, CG5, CG6, CG7, CG8, CG33, CG34, CG35, CG36, CG37, CG38, CG39, CG44, CG45 | `bench test --package ./internal/gate`, `bench test --package ./internal/reviewrecord/recordcmd`, `bench test --package ./internal/landing`, `bench test --package ./internal/anchors`, `bench test --package ./internal/conformance`, `bench test --package ./internal/reviewrecord`, `bench test --package ./cmd/bench` | no |

CG-C1 creates the seam that CG-C2 consumes, so its chunk review closes first.
CG-C2 creates the seam that CG-C3 consumes.

## Testing decisions

- A good checkpoint test builds a real fixture repository, makes the gap with real commits, and runs `RunCommand` with `--checkpoint`. It reads the exit code, the output, and the oracle run count. `TestReviewCheckpointLaterSource` in `internal/gate/review_checkpoint_test.go` is the prior art, and it was read in this session.
- The checkpoint fixture gets its Go file through the `prepare` hook of `attachedCheckpointFixture`, before `AddChunk` commits the chunk. A gap test edits that file with `Write` and `Commit`. No `recordtest` file changes.
- A good classifier test drives `Prove` over real tree pairs and the Go rule over byte pairs. Each refusal row asserts its rule sentinel with `errors.Is`, so a row cannot pass through another rule.
- The `--complete` checkpoint row calls `Complete` and then re-points `Completion.SourceDigest` and the final verification to the corrected digest. `recordtest.Complete` pins completion to the last chunk digest, and `f.Verification` rebuilds the final verification at a named digest.
- The landing row adds one case to `TestLandingCompletionEvidence`. Today `fixture(t)` commits only `named` and `foreign`. `attachedCompletionFixture` serves `newCompletionFixture` and the delegated test in `delegated_completion_test.go`. So `attachedCompletionFixture` gains a variadic prepare hook that writes one Go file before the first chunk. `newCompletionFixture` forwards the hook, and the delegated caller stays unchanged. The case commits a comment edit after the second chunk and re-points the completion evidence to the corrected digest.
- `TestReviewCheckpointChainGapNamesTheExpectedBase` asserts only the first two clauses of the chain rule today, so it passes with either rule text. Its second assertion gains the clause `and only record commits and comment-only corrections follow a chunk tip` (CG35).
- The mutation-probe target for CG-C2 replaces `if directive(comment.literal) {` with `if false {`. The directive refusal rows must fail.
- `TestComposedChangesExpandsANamedDirectory`, `TestComposedChangesRepresentsARenameAsDeletionAndAddition`, and `TestComposedChangesCarriesTheSymlinkMode` guard the moved reader without an edit.
- The gate's `test` phase runs every package test above and the root conformance test, so the gate observes each row except CG40.

### Seam diagram

    trigger: `bench gate --checkpoint <spec> (--chunk <id> | --complete)`,
             and the landing broker through `gate.WithCompletion`
        │
        ▼
    reviewed chunk tip commit,  ──▶  [ reviewrecord.checkSource ]  ──▶  accept, or the current
    graded tree                          │                                stale-source refusal
                                         ▼
                              [ commentgap.Prove ] ◀── [ git.TreeChangesIncludingSubmodules, git.ReadTreeFile ]
                      ◀ tests attach here: real fixture repositories, `RunCommand` exit code,
                        output, and oracle run count; tree pairs and byte pairs at `commentgap`

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CG1 | 1, 36 | A `--chunk 1` checkpoint exits 0 when one commit after the reviewed tip changes only a comment line in a `_test.go` file | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointCommentOnlyGap`) | The current digest equality refuses the gap as a stale reviewed source. |
| CG2 | 2 | A `--complete` checkpoint exits 0 when the completion evidence names the corrected source and the gap after the last chunk changes only a Go comment | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointCommentOnlyGap`) | An acceptance placed only on the `--chunk` exit leaves the complete path refusing. |
| CG3 | 3 | A reviewed landing publishes `Status: implemented` when its source tip holds a comment-only Go correction after the last chunk | `internal/landing/completion_evidence_test.go` (`TestLandingCompletionEvidence`) | A proof against the composed tree sees the status transform of the spec and refuses. |
| CG4 | 4 | An accepted comment-only gap raises the oracle run count by one | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointCommentOnlyGap`) | An acceptance that returns before the gate leaves the run count unchanged. |
| CG5 | 5 | A chunk whose Standards result carried a finding passes after a comment-only correction and one superseding Standards pass at the reviewed frozen pair | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointCommentOnlyGap`) | An acceptance that binds the axes to the corrected tree refuses the two unchanged axes. |
| CG6 | 6 | A commit after the reviewed tip that changes a Go statement refuses with `chunk 1: stale reviewed source: no chunk review covers`, and the oracle does not run | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | An acceptance that skips the proof passes the code change. |
| CG7 | 7 | A committed comment-only correction beside an uncommitted Go statement change refuses with `chunk 1: stale reviewed source` | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A proof over the commit range or the HEAD tree passes the uncommitted change. |
| CG8 | 8 | A chunk entry whose recorded digest names the corrected tree while its tip names the reviewed commit refuses with `chunk 1: stale source digest`, and the oracle does not run | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | An acceptance that trusts the recorded digest before the recompute passes the forged entry. |
| CG9 | 9 | A change of only the trailing comment on a line that also holds code is proven | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A rule that requires each changed line to start with `//` refuses the line. |
| CG10 | 10 | A comment edit on a last line with no trailing newline is proven | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A line splitter that requires a final newline misreads the last line. |
| CG11 | 11 | A change of only blank lines and indentation between Go tokens is proven | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A byte comparison of the text without comments refuses the layout change. |
| CG12 | 12 | A raw-string literal whose bytes change from `// one` to `// two` is refused with `ErrTokens` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A line-prefix or pattern comment stripper reads the string bytes as a comment and passes them. |
| CG13 | 13 | A block comment that gains a newline between an operand and its operator is refused with `ErrTokens` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A comment stripper that ignores the inserted semicolon passes the changed program. |
| CG14 | 14 | An edit that leaves a block comment unterminated is refused with `ErrScan` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A token comparison without the scanner error handler sees equal tokens and passes. |
| CG15 | 15 | A `//go:build linux` comment that becomes `//go:build darwin` is refused with `ErrDirective` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A scan that skips comments sees equal tokens and passes the build constraint change. |
| CG16 | 15 | A `//go:embed` comment added above an existing variable declaration is refused with `ErrDirective` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A comparison of changed comments alone misses a directive that only one side holds. |
| CG17 | 15 | A `//export F` comment that becomes `//export G` is refused with `ErrDirective` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A `//go:` prefix rule misses the cgo export directive. |
| CG18 | 16 | A `// +build linux` line that becomes `// +build darwin` is refused with `ErrDirective` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | The directive shape needs a non-space third byte, so only the legacy build-line predicate catches this line. |
| CG19 | 17 | A preamble block comment that changes from `/* #include <stdio.h> */` to `/* #include <stdlib.h> */` in a file that imports `"C"` is refused with `ErrCgo` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A token rule alone passes a change to the C source in the preamble. |
| CG20 | 18 | A comment edit in a `_test.go` file whose two sides both hold an `// Output:` comment is refused with `ErrExampleOutput` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A token rule alone passes a change that `go test` compares with the example output. |
| CG21 | 18 | An `// Output:` comment that only the later side of a `_test.go` file holds is refused with `ErrExampleOutput` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A rule that reads only the reviewed side passes an added expected output. |
| CG22 | 19 | A `notes.md` file whose two sides scan as equal Go tokens is refused with `ErrNotGo` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A classifier without the suffix rule passes the Markdown edit. |
| CG23 | 20 | An added Go file that holds only a comment is refused with `ErrStatus` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A classifier that reads a missing side as empty sees zero tokens on each side and passes. |
| CG24 | 20 | A deleted Go file that held only a comment is refused with `ErrStatus` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A classifier that reads a missing side as empty passes the deletion. |
| CG25 | 20 | A renamed Go file with a comment edit is refused with `ErrStatus` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A rename-detecting change list pairs the two paths and passes the comment edit. |
| CG26 | 21 | A Go file whose mode changes from `100644` to `100755` with the same bytes is refused with `ErrMode` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A blob comparison alone sees equal bytes and passes. |
| CG27 | 22 | A symbolic link named `link.go` whose target text changes is refused with `ErrMode` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A suffix rule alone reads the link target text as Go source. |
| CG28 | 22 | A gitlink entry named `kit.go` whose commit changes is refused with `ErrMode` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A suffix rule alone tries to read a commit ID as a Go blob. |
| CG29 | 23 | A comment-only edit of a Go file larger than `bounds.ControlRecordLimit` is refused with `ErrUnreadable` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | An unbounded blob read passes the oversized file. |
| CG30 | 22 | Two trees that differ only in a gitlink, read under `diff.ignoreSubmodules=all`, are refused with `ErrMode` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A configured reader hides the gitlink and cannot classify its mode. |
| CG31 | 25 | Comment-only edits of `a b*.go` and of a Go file whose name holds a tab are proven | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A newline-framed or C-quoted path list names a path that no tree holds and refuses. |
| CG32 | 26 | A gap with comment-only edits of `a.go` and `b.go` and a statement change in `c.go` is refused with `ErrTokens` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | A classifier that returns after the first proven path passes the later code change. |
| CG33 | 27 | A `--complete` checkpoint refuses with `completion is incomplete or stale` when the completion evidence names the reviewed source and a comment-only correction follows it | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A gap acceptance copied into the completion check passes stale final verification. |
| CG34 | 28 | A second chunk whose base sits after a comment-only correction of the first chunk refuses with `expected base` and the first chunk tip | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A gap acceptance copied into the chain check leaves the correction outside every chunk delta. |
| CG35 | 29 | The chain-gap refusal holds `and only record commits and comment-only corrections follow a chunk tip` | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointChainGapNamesTheExpectedBase`) | The assertion gains the new clause, so a refusal that keeps the old rule text fails it. |
| CG36 | 30 | A chunk entry re-recorded at the corrected tip, with verification at that tip and the three axis results at the reviewed pair, refuses with `chunk 1: stale Standards` | `internal/gate/review_checkpoint_commits_test.go` (`TestReviewCheckpointKeepsStrictEvidence`) | A gap acceptance copied into the axis check passes results that no axis gave for the new pair. |
| CG37 | 31, 37 | The policy section `Classification and completion` holds the new comment-only sentence, and its anchor bites when the sentence is removed | `internal/conformance/implementation_continuation_test.go` (`TestImplementationContinuation`) | Without the anchor, a later edit removes the rule and no check turns red. |
| CG38 | 32 | A deletion of the new registry row reds `TestImplementationContinuation` with `implementation-continuation anchor is absent` | `internal/conformance/implementation_continuation_test.go` (`TestImplementationContinuation`) | An expectation that reads the registry instead of its own copy passes the deletion. |
| CG39 | 33 | The implementation phase's `Land` section holds `Only record commits and comment-only corrections follow the chunk tip.`, and the old sentence fails the anchor | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | An unchanged needle keeps the old sentence, so the guidance and the checkpoint disagree. |
| CG40 | 34 | `internal/gate` holds no raw-diff entry parser after the move, and `gate.ComposedChanges` calls `git.TreeChanges` | review-owned: the Standards axis reads `internal/gate/lane_select.go` | No test can tell a surviving second parser from the shared one. |
| CG41 | 35 | The lane's change list for a commit that names a directory holds the two changed files with their modes | `internal/gate/lane_select_test.go` (`TestComposedChangesExpandsANamedDirectory`) | A moved parser that changes the framing or the field order fails the pinned list. |
| CG42 | 14 | An edit that terminates a block comment that the reviewed side left unterminated is refused with `ErrScan` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A scan that grades only the later blob sees no error and equal tokens, so it passes. |
| CG43 | 18 | An `// Output:` comment that only the reviewed side of a `_test.go` file holds is refused with `ErrExampleOutput` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | A rule that reads only the later side passes a removed expected output. |
| CG44 | 38 | `bench record completion --source <corrected tip>` exits 0 and writes a completion at the corrected digest when a comment-only correction follows the last chunk | `internal/reviewrecord/recordcmd/completion_test.go` (`TestRecordCompletionAcceptsCommentOnlyGap`) | An acceptance placed only in `CheckTrees` leaves `RecordCompletion` refusing, while every checkpoint row stays green. |
| CG45 | 38 | `bench record completion --source <tip>` exits 1 with `stale reviewed source` and leaves the record bytes unchanged when the gap after the last chunk changes a Go statement | `internal/reviewrecord/recordcmd/completion_test.go` (`TestRecordCompletionAcceptsCommentOnlyGap`) | A record path that skips the proof writes a completion over an unreviewed code change. |
| CG46 | 40 | A hidden gitlink change beside a Go comment edit under `diff.ignoreSubmodules=all` refuses with `ErrMode` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | The configured list holds the Go edit, so an empty-list check alone accepts the hidden gitlink. |

| CG47 | 24 | Different trees whose only difference is an empty directory tree refuse with `ErrEmptyChanges` | `internal/commentgap/commentgap_test.go` (`TestProveTreeGap`) | Git omits the empty directory from its change list, so this row pins the empty-list refusal. |

| CG48 | 15 | Moving an unchanged embed directive from one variable to another refuses with `ErrDirective` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | Equal token and directive-text lists cannot prove that directive attachment stays unchanged. |
| CG49 | 15 | An ordinary comment edit beside an unchanged directive refuses with `ErrDirective` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | The conservative rule refuses directive-bearing files, even when only another comment changes. |
| CG50 | 18 | Multiline block example output on either side refuses with `ErrExampleOutput` | `internal/commentgap/commentgap_test.go` (`TestProveGoGap`) | Go reads output after a leading newline, which a spaces-and-tabs trim misses. |

Not covered: story 39 — the light path has no checkpoint, so no row can observe its repair accounting.

### Edge inventory

The walk of the hostile-input checklist in `projects/benchkit.md`:

| class | disposition |
| --- | --- |
| path with spaces or glob characters | CG31 |
| control bytes in git-sourced text | The checkpoint does not print the classifier error, and the reader keeps raw bytes. CG31 covers a tab. |
| C-quoted paths in a patch header | The reader uses `-z` raw entries, not a patch. CG31 covers the framing. |
| bytes a sink permits | No new sink: the refusal text does not change. |
| numeric-looking TOON cells | No TOON output changes. |
| whitespace predicate on hand-edited text | The Go scanner owns the whitespace class. CG11 covers the accepted side, and CG13 covers a newline that changes a token. |
| a command whose own write changes a reported fact | The checkpoint writes no record. |
| a path read out of a file | The reviewed tree comes from the recorded tip commit and its recompute. CG8 covers a forged digest. |
| last line without a trailing newline | CG10 |
| absent versus empty file | CG23 and CG24 cover absent sides. An empty file on both sides gives equal tokens, and the gate grades its compile. |
| special files | CG27 and CG28 |
| dangling or live symlink | CG27: the reader takes the tree entry mode, and no link is followed. |
| symbolic refs | No ref changes. |
| multi-word arguments | No shell entry. |
| lifecycle guidance through raw Git | The guidance edits name `bench record review` and the checkpoint only. |
| JSON-escaped operators | No hook or guard input. |
| operator runs | No hook or guard input. |
| flag value read as a positional | No new grammar. |
| grammar token quoted in prose | No new grammar. |
| non-ASCII whitespace in Markdown | No Markdown rule. A non-ASCII space outside a Go comment is a scanner error, which CG14's rule refuses. |
| required tool missing | The checkpoint already needs `git`. |
| invocation through a symlink | No new entry point. |
| every shipped surface | CG1 covers `bench gate --checkpoint`, CG3 covers the landing broker, and CG44 covers `bench record completion`. |
| destructive worktree state | No worktree change. |
| a two-sided merge of a rewritten path | No composition change. |
| state reloaded by a fresh process | CG1 to CG8 read the saved record through `RunCommand`. |
| cwd deeper than the root | The gate resolves its root as today. |
| non-TTY stdin | No prompt. |
| WSL fsync stalls | No new durable write. |
| unterminated delimiter | CG14 |
| temporary root with a symlink component | Tests compare tree IDs, not paths. |
| a check inside the fast lane checkout | The proof runs in the checkpoint, not in a lane check. |

**Won't handle** lines:

- A comment edit outside a Go file — shell, Markdown, JSON, YAML, `.mjs`, `.cjs`, TOON, canary fixtures, and `#`-comment configuration keep the repair path (CG22).
  Their comment bytes are read as data, or the tree has no scanner for them.
- A comment edit on the light path — the light path has no checkpoint that proves the gap, so it keeps the light-path repair allowance.
- Completion evidence recorded before a comment-only correction — the orchestrator reruns final verification and `bench record completion` (CG33, CG44).
  That rerun needs no review and consumes no repair cycle.
- A chunk base after a comment-only correction — the next base stays the predecessor tip, so the correction joins its delta (CG34).
- A file with a directive-shaped comment, a cgo file, or a `_test.go` file with an example-output comment takes the repair path.
  CG19, CG20, CG21, CG48, CG49, and CG50 pin these refusals.
- A Go blob over the control-record bound — the existing repair path takes the edit (CG29 pins the refusal).
- A comment edit that adds or removes the space after `//` — the rule reads `//TODO` as directive-shaped, so the edit takes the repair path (CG15).

## Ownership fences

- `internal/git/tree.go`
- `internal/gate/lane_select.go`
- `internal/commentgap/`
- `internal/reviewrecord/coverage.go`
- `internal/reviewrecord/recordcmd/completion_test.go`
- `internal/reviewrecord/recordcmd/verification_test.go`
- `internal/reviewrecord/recordcmd/chunk_test.go`
- `internal/gate/review_checkpoint_commits_test.go`
- `internal/landing/completion_evidence_test.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/conformance/implementation_continuation_test.go`
- `.agents/skills/bench-craft-line/references/bounded-repair-policy.md`
- `.agents/commands/bench-implement-spec.md`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration`
- `tests/canary/workflow-guidance-anchors/delegated-entry-refusals`
- `tests/canary/workflow-guidance-anchors/delegated-resumption-contents`
- `tests/canary/workflow-guidance-anchors/dg-25`
- `tests/canary/workflow-guidance-anchors/dg-26`
- `tests/canary/workflow-guidance-anchors/dg-29`
- `tests/canary/workflow-guidance-anchors/dg-29-verification-target`
- `tests/canary/workflow-guidance-anchors/dg-30`
- `tests/canary/workflow-guidance-anchors/dg-31`
- `tests/canary/workflow-guidance-anchors/dg-31-contradiction-trigger`
- `tests/canary/workflow-guidance-anchors/implement-spec-adoption-freshness`
- `tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding`
- `tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer`
- `tests/canary/workflow-guidance-anchors/implement-spec-entry-validation`
- `tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness`
- `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception`
- `tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-prose-owner-transfer`
- `tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper`
- `tests/canary/workflow-guidance-anchors/implement-spec-red-preflight-route`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-omitted`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-order`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-reversed`
- `tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight`
- `tests/canary/workflow-guidance-anchors/implement-spec-write-delegation`
- `tests/canary/workflow-guidance-anchors/line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/prepared-build-approval`
- `tests/canary/workflow-guidance-anchors/prepared-build-freshness`
- `CHANGELOG.md`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `reviews/ft370-comment-only-evidence.md`

## Out of scope

- A comment-only predicate for shell scripts, with a shell lexer that keeps the guard header, the shebang, and the `# shellcheck` directives strict. Estimate: 14 edits, 3 gate runs.
- A comment-only path for light-path work, which needs a light-path checkpoint first. Estimate: 10 edits, 3 gate runs.
- The lane reader's exposure to `diff.ignoreSubmodules`, which hides a gitlink change from lane selection. Estimate: 3 edits, 1 gate run.
- A split of the overloaded term "source identity" across the glossary, the policy, and the code. Estimate: 8 edits, 1 gate run.

## Further notes

### Closed reviewer decision: build order

The reviewer decided on 2026-10-03 that FT370 builds before the staged FT358
spec `markdown-block-reader`. FT358's build starts after the FT370 landing,
and it absorbs the FT370 diff with one `bench worktree merge --from main`.

### Closed reviewer decision: the orchestrator makes the correction

The reviewer decided on 2026-10-03 that the orchestrator commits a
comment-only correction. The plan records no assignment for it, so the plan
digest stays the same. Review iteration 1 found that a recorded repair
assignment edits the plan in `spec.md` and so blocks the comment-only path.
The plan already names the orchestrator session, so the author stays
auditable. A refused gap takes the normal repair route.

### Overlap with FT358

The planned FT370 writes and their line counts at `16efdcb4`:

| path | lines | FT358 writer |
| --- | --- | --- |
| `internal/git/tree.go` | 197 | none |
| `internal/gate/lane_select.go` | 370 | none |
| `internal/commentgap/` (new) | 0 | none |
| `internal/reviewrecord/coverage.go` | 258 | none |
| `internal/reviewrecord/recordcmd/completion_test.go` | 113 | none |
| `internal/gate/review_checkpoint_commits_test.go` | 48 | none |
| `internal/landing/completion_evidence_test.go` | 173 | none |
| `internal/anchors/registry_retained_workflow.go` | 396 | ticket 11 |
| `internal/anchors/registry_chunk_chain.go` | 42 | none |
| `internal/anchors/registry_chunk_chain_test.go` | 30 | none |
| `internal/conformance/implementation_continuation_test.go` | 308 | none |
| `.agents/skills/bench-craft-line/references/bounded-repair-policy.md` | 68 | none |
| `.agents/commands/bench-implement-spec.md` | 80 | none |
| `internal/anchors/registry_debug_loop.go` | 61 | none |
| `internal/anchors/registry_data.go` | 481 | ticket 11 |
| `internal/anchors/registry_data_test.go` | 1288 | ticket 11 |
| `internal/anchors/registry_ft311_preparation.go` | 51 | none |
| `cmd/bench/command_registry.go` | 387 | tickets 2, 3, 4, 6, 8, and 11 |
| `cmd/bench/command_registry_test.go` | 794 | tickets 2, 3, 4, 6, 8, and 11 |
| `cmd/bench/help_inventory_test.go` | 314 | tickets 2, 3, 4, 6, 8, and 11 |
| `internal/conformance/axi_query_registry_test.go` | 445 | tickets 2, 3, 4, 6, 8, and 11 |
| `internal/conformance/subcommand_routing_table_test.go` | 86 | tickets 2, 3, 4, 6, 8, and 11 |
| the 29 `tests/canary/workflow-guidance-anchors/` fixture directories in the fence | fixtures | none |

The one shared path that FT370 edits is
`internal/anchors/registry_retained_workflow.go`. FT370 adds one row there, so
the file reaches 397 of its 400 lines before FT358 ticket 11 writes it. The
other shared paths are build preflight closure of ticket 4, and FT370 expects
no edit there. The slicing pass needs these exact paths for the fence.

FT358 ticket 7 writes `internal/reviewrecord/parse.go` (275 lines),
`internal/reviewrecord/write.go` (318 lines),
`internal/reviewrecord/record_test.go` (114 lines), and
`internal/reviewrecord/recordtest/fixture.go` (366 lines). FT370 writes none of
these four files. So FT370 keeps each contract of ticket 7:

- FT370 adds no file to `internal/reviewrecord`.
- FT370 adds no use of the `fenceMarker` constant.
- FT370 adds no string literal with a run of three backticks to `internal/reviewrecord`.
- FT370 grows no file in ticket 7's `Writes:` line.

FT370 moves neither ticket 7 citation, `parse.go:244` or
`recordtest/fixture.go:81`. The cleanest FT370 seam needs none of the four
files. The gate tests make the gap through the `prepare` hook and the
fixture's `Write` and `Commit` methods.

### Flagged decisions

- **F1. Whitespace-only layout passes.** The Go rule compares tokens, so a change of only blank lines and indentation passes beside a comment edit. The Go tokens are the executable content of the file. The reviewer accepted F1 on 2026-10-03, because a token-equal change has no executable line, which is the source's reason for the rule.
- **F2. Records of the correction.** The chunk keeps its reviewed pair, and only the issuing axis records a superseding pass. `bench record verification` refuses a source that differs from the chunk source, and the source claims no verification. The reviewer accepted F2 on 2026-10-03, because the checkpoint proof and the gate replace re-verification. The authorship half of F2 is closed above.

### Flagged additions

Each addition below is not a sentence of the decision source. The reviewer can
veto each one.

- The new package `internal/commentgap` and its ten rule sentinels.
- The move of the raw tree-change reader to `internal/git` (CG40, CG41).
- The empty-change-list refusal (CG47).
- The whitespace-only acceptance (CG11, F1).
- The edit of `bench-implement-spec.md` line 54 and its anchor (CG39).
- The chain-gap message text (CG35).
- The rows that pin unchanged strictness: CG33, CG34, and CG36.

### Source trace

| source sentence | rows |
| --- | --- |
| "A comment edit changes the source digest. Thus a one-word comment fix makes every review axis stale and costs a full repair cycle." | CG1, CG5 |
| "A comment-only correction has no executable line, no finding, no observation, and no verification claim." | CG9 to CG32, CG42, CG43, and F2 |
| "It takes the evidence-only path." | CG5, CG37 |
| "The bounded repair policy gains one clause." | CG37, CG38 |
| "The checkpoint accepts a comment-only gap between the reviewed tip and the current tip." | CG1 to CG8 |
| Occurrence: VR-C4, commit `c9731327`, one comment line in `internal/worktree/land_fixtures_test.go` | CG1 uses the same shape |
| Occurrence: each confirming-round finding was a real defect, and the cost came from the tip move | CG5 |

### Reader sweep

Readers of the decision fact "the reviewed source is current":

- `reviewrecord.checkSource` holds every staleness comparison. Its production callers are `reviewrecord.CheckTrees`, through `gate.applyCheckpoint`, and `reviewrecord.RecordCompletion`, through `bench record completion`.
- `gate.WithCompletion` routes the landing broker into the same comparison.
- `preflight.completionEvidenceRow` renders the source digest and the record state. It renders no staleness verdict, so it does not change.
- `recordcmd` renders the digests of each written entry. It renders no verdict, so it does not change.
- `reviewrecord.RecordVerification` refuses a source that differs from the chunk source. It stays strict, because a correction records no verification.

Guidance readers, found by key terms and in normalized form:

- the bounded repair policy, lines 46 to 52
- `bench-implement-spec.md`, lines 54 to 58
- `bench-review-implementation.md`, lines 42 to 48, 95, 98, 229, and 231
- `.bench/BENCH-reference.md`, lines 375 to 383
- `projects/benchkit.md`, lines 38 to 47
- `CONTEXT.md`

Only the policy and line 54 change. Line 95 already sets each later chunk base
to the accepted predecessor tip. Line 98 stays true because the chain stays
strict. `.bench/BENCH.md` line 127, `bench-implement-spec.md` line 56, and
`bench-review-implementation.md` line 47 speak of a repair. The guidance
section states why none of them applies to an accepted correction.

Readers of the raw tree-change reader:

- `internal/gate/authorization/authorization.go` line 227
- the lane class functions `proseSubject` and `changeClasses`
- the tests at `lane_select_test.go` lines 72, 95, and 118
- the tests at `lane_run_test.go` lines 171 and 233

Shipped-surface claim words: none. The guidance edits name no repository-only
path.

### Pre-review proof checklist

- Cited symbols: each symbol resolves at `16efdcb4`. In `internal/reviewrecord`: `checkSource`, `CheckSource`, `CheckTrees`, `SourceDigest`, `checkReviews`, `checkVerification`, `checkCompletion`, `RecordChunk`, `RecordVerification`, and `RecordReview`. In `internal/gate`: `applyCheckpoint`, `WithCompletion`, `completionTree`, `ComposedChanges`, `ComposedChange`, `parseComposedChange`, `attachedCheckpointFixture`, and `runCheckpoint`. In `internal/git`: `TreeWithoutFile`, `ReadTreeFile`, and `ReadControlBlob`. In the standard library: `scanner.ScanComments` and `constraint.IsPlusBuild`.
- Import edges: `internal/reviewrecord` gains `internal/commentgap`. `internal/commentgap` imports `internal/git`, `go/scanner`, `go/token`, `go/parser`, and `go/build/constraint`. `internal/git` depends only on `internal/bounds` and `internal/canonicalpath`.
- Source-row clauses and occurrences: the source trace table above.
- Promised field labels: `TreeChange{Status, SrcMode, DstMode, Path}` and the two diagnostics named under the guidance.
- Promised sentinels: `ErrUnreadable`, `ErrEmptyChanges`, `ErrStatus`, `ErrMode`, `ErrNotGo`, `ErrScan`, `ErrTokens`, `ErrDirective`, `ErrCgo`, and `ErrExampleOutput`.
- Changed-function callers: `checkSource` has the callers `CheckSource`, whose callers are tests, and `CheckTrees`, whose callers are `gate.applyCheckpoint` and `Check`, whose callers are tests. `RecordCompletion` also calls `checkSource`, and the `recordcmd` completion handler calls `RecordCompletion`. `gate.ComposedChanges` has the production caller `internal/gate/authorization/authorization.go` line 227 and the test callers named in the reader sweep.
- Copy survival: CG40 is review-owned, because no test can tell a surviving parser from the shared one.
- Rendered-shape readers: the old needle `only record commits follow a chunk tip` hits only `internal/reviewrecord/coverage.go` line 81. The old needle `Only record commits follow the chunk tip.` hits `bench-implement-spec.md` line 54, `registry_chunk_chain.go` line 18, and `registry_chunk_chain_test.go` line 13. No canary fixture holds either needle.

### Proposed glossary entry

The FT358 precedent lands glossary terms in a separate coordinator commit. The
proposed entry for `CONTEXT.md`:

- **comment-only gap** — a gap after a reviewed chunk tip in which each
  changed path is an existing regular Go file. Each file keeps its Go tokens
  and holds no directive-shaped comment. The checkpoint proves the gap from the two
  trees. Not "comment fix", not "cosmetic change", not "trivial delta" —
  comment-only gap.

`CONTEXT.md` defines **source identity** as the SHA-256 of one exact source
body. The bounded repair policy, `reviewrecord` error text, and
`internal/diff/range.go` use the term for a review source. FT370 does not use
the term, and the split is out of scope.

### Sources read in the authoring session

Read in full or in the named lines:

- `roadmap/FT370.md`
- `internal/reviewrecord`: `coverage.go`, `check.go`, `record.go`, `plan.go` lines 150 to 163, `write.go` lines 100 to 318, and `recordtest/fixture.go` lines 1 to 340
- `internal/gate`: `checkpoint.go`, `completion.go`, `lane_select.go`, `review_checkpoint_test.go`, `review_checkpoint_commits_test.go`, and `lane_select_test.go` lines 55 to 130
- `internal/git/tree.go` and `internal/landing/completion_evidence_test.go`
- `internal/preflight/review.go` lines 270 to 321
- `internal/anchors`: `registry_retained_workflow.go` lines 370 to 396, `registry_chunk_chain.go`, and `registry_chunk_chain_test.go`
- `internal/conformance`: `implementation_continuation_test.go` lines 170 to 308 and `git_plumbing_owner_test.go` lines 1 to 60
- the bounded repair policy and `bench-implement-spec.md`
- `bench-review-implementation.md` lines 36 to 65 and 200 to 246
- `CONTEXT.md` lines 340 to 473 and `projects/benchkit.md`
- the FT358 fence and ticket 7

The author compared the `bench anchors` output for the policy and for
`bench-implement-spec.md` with each anchor claim above.

Not read: `internal/gate/run_outcomes_test.go`, which holds `outcomeFixture`
and `outcomeRuns`. The landing `fixture` helper and the `bench anchors` output
for `bench-review-implementation.md` are also unread. This spec does not edit
that command file.

### In-scope implementation expansion

Ticket 2 also writes `internal/git/tree.go` to add the complete-list wrapper.
The wrapper shares the existing reader and parser.
The reviewer approves the complete reader and CG30's `ErrMode` result on 2026-10-03.
The classifier refuses the mixed hidden-gitlink case through CG46.
CG47 preserves the empty-list guarantee with an empty directory tree.
All chunk IDs and dependencies stay unchanged.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"CG-C1","tickets":["1-move-tree-change-reader.md"],"verification":[{"id":"git","command":"bench test --package ./internal/git"},{"id":"gate","command":"bench test --package ./internal/gate"}]},{"id":"CG-C2","tickets":["2-prove-comment-only-gaps.md"],"verification":[{"id":"commentgap","command":"bench test --package ./internal/commentgap"},{"id":"git","command":"bench test --package ./internal/git"},{"id":"gate","command":"bench test --package ./internal/gate"}]},{"id":"CG-C3","tickets":["3-accept-proven-gaps-at-checkpoint.md","4-state-comment-only-correction-rule.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"recordcmd","command":"bench test --package ./internal/reviewrecord/recordcmd"},{"id":"landing","command":"bench test --package ./internal/landing"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord"},{"id":"cmd-bench","command":"bench test --package ./cmd/bench"}]}],"final_verification":[{"id":"coverage-check","command":"bench coverage --check specs/ft370-comment-only-evidence/spec.md"},{"id":"commentgap","command":"bench test --package ./internal/commentgap"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"recordcmd","command":"bench test --package ./internal/reviewrecord/recordcmd"},{"id":"landing","command":"bench test --package ./internal/landing"},{"id":"anchors","command":"bench test --package ./internal/anchors"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"git","command":"bench test --package ./internal/git"},{"id":"reviewrecord","command":"bench test --package ./internal/reviewrecord"},{"id":"cmd-bench","command":"bench test --package ./cmd/bench"}]}
```

### Approved conservative comment rules

The reviewer approves refusal of changed files with any directive-shaped comment on 2026-10-03.
This rule covers unchanged directives whose attachment can change program behavior.
The reviewer also approves recognition of multiline block example output.
CG48 to CG50 pin these decisions, and the directive-refusal omission is the named probe.

### Integration fixture and release note

Ticket 3 moves `recorded` beside the chunk fixture helpers and adds a prepare hook.
This keeps the verification test file within its existing 400-line limit.
The completion fixture forwards that hook to seed Go before the reviewed chunk.
Ticket 4 adds the required typed entry to `CHANGELOG.md`.
The C3 and final verification lists also run the reviewrecord and command packages.

### Candidate dogfood evidence

On 2026-10-03, a deterministic root-authored adapter drove a real shift in a disposable Go repository with the candidate kit.
The task corrected integer addition, and the existing Go test graded the result.
This run checked the CLI, gate, and Stop hook; the acceptance tests checked the comment-gap behavior.

The first run, `bench/shift-20261003-195354`, failed because the fixture lacked its Go gate environment declaration.
No iteration committed.
The fixture then declared its environment and tool inputs.

The second run, `bench/shift-20261003-195441`, completed with exit 0 and one committed iteration.
The Stop hook returned 2 while the real Go test failed and returned 0 after the correction.
The shift reused that exact green verdict and met its completion predicate.
The native output is retained in the review record.
