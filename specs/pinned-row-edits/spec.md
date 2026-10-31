# Record a pinned-row decision with its approved plan

Status: staged

Roadmap: FT437

Decision source: named reviewed artifact — the bound criterion `FT437-pinned-rows.done` in `.bench/commitment.json`, read with `roadmap/FT437.md` at its pinned identity `sha256:ed986b1d5df9ff39b3f63d5a4c23c5c7100ac4a116a129a5fb94bced2f91cb80`. The row records the two reviewer decisions of 2026-10-29. Under the interim rule of 2026-10-29, the criterion wins over the row text, and here the two agree.

Verification log: 2 iteration(s) to accept — the spec round and the ticket round each accepted at iteration 1. The ticket round also graded the coordinator's post-acceptance delta, which moved the link walk below the store root and added PE27 to PE29.

Implementation approval: approved by the reviewer on 2026-10-29, including the one-ticket breakdown, the single chunk PE-C1, and the declared implementation line.

## Problem

A drain records each settled row decision in the row's `roadmap/FT<n>.md` file.
When a committed outcome pins that row, the decision cannot reach the file.

- `bench commit` refuses the row edit under the current policy. `protectedCandidate` reads each pinned source at the candidate tree and compares it with the pinned identity (`internal/commitment/repository/candidate.go:239`).
- `bench commitment plan` refuses an input that names the edited identity. `readTree` reads each existing pinned row at the default-branch tip (`internal/commitment/repository/sources.go:98`).
- `bench commit` of an approved policy together with the row edit also refuses. `approvedTransition` reads each plan source through the same `readTree`, so it reads the default-branch tip too (`internal/commitment/repository/candidate.go:211`).

The decision then goes only into an outcome criterion, and the row file keeps its old decision state.
The row records seven such occurrences, and five row edits wait in a saved drain patch.
A stale pinned row also misleads a spec stage that takes the row as its source.

## Solution

A plan input that changes the identity of a pinned row carries the row edit with it.
Plan and approve read that row at the working file of the checkout, so the edit can stay uncommitted until the approval.
The commit and the landing read that row at the candidate tree, so the committed bytes must equal the approved identity.
One approved plan and one commit then record the decision in its row file.

Every other source keeps its current read.
A source whose identity the proposal keeps still reads the default-branch tip, so a worktree edit of that row still cannot pass as current.
A pinned-row edit with no exact approved plan still refuses.

## User stories

Line: opus / medium.
Implementation-line reason: the one chunk moves the read location of an authority check. The spec fixes each read location and each refusal, the seam is the existing `readTree` with its fixtures, and package tests red each wrong location. The project `Lines` section routes enforcement logic to mid effort, because a wrong admission is the worst class of defect.
Harder chunks: none.

### Plan and approve read a changed pinned row at the working file

1. As a drain operator, I want `bench commitment plan` to accept the new identity of a pinned row when the working row file holds those bytes, so that the plan can carry an uncommitted row edit.
2. As a drain operator, I want `bench commitment approve` to accept that plan while the working row file still holds those bytes, so that the approval stages the policy beside the row edit.
3. As a drain operator, I want one plan to carry the new identities of several pinned rows, so that one drain batch records all of its pinned-row decisions in one plan.
4. As a reviewer, I want the plan to refuse a new identity that the working row file does not hold, so that a plan cannot name bytes that nobody wrote.
5. As a reviewer, I want the plan to refuse a new identity that only the checkout HEAD holds, so that the read location is the working file and not the HEAD commit.
6. As a reviewer, I want the approval to refuse when the working row file changed after the plan, so that the approval cannot stage a policy for bytes that the checkout no longer holds.
7. As a reviewer, I want the plan to refuse a working row path that is absent, empty, malformed, oversized, linked, or not a regular file, so that a hostile checkout cannot redirect or block the read.
8. As a drain operator, I want the working-file identity to drop the `Occurrences:` ledger line as the tree read does, so that the plan identity equals the identity that the commit reads for the same bytes.
9. As a reviewer, I want the plan to refuse a working row with a malformed occurrence ledger, so that the working read keeps the ledger rule of the tree read.
10. As a reviewer, I want a refused plan to write nothing, so that a refusal leaves no receipt and no changed checkout file.

### The narrowed story 16 of `spec-stage-binding`

11. As a reviewer, I want each pinned source whose identity the proposal keeps to read the default-branch tip, so that a worktree edit of that row cannot pass as current.
12. As a reviewer, I want a new outcome source at a path that main holds and that no current source binds to read the default-branch tip, so that only an approved change of a pinned row reads the working file.
13. As a spec author, I want a new deliverable to keep its checkout HEAD read at plan, so that the spec-stage binding route does not change.

### The commit and the landing read a changed pinned row at the candidate tree

14. As a drain operator, I want `bench commit` to accept the approved policy, the board, and the row edit in one commit, so that one commit records the decision in its row file.
15. As a reviewer, I want that commit to refuse when it names the approved policy without the row edit, so that the policy cannot land beside a stale row.
16. As a reviewer, I want that commit to refuse when the committed row bytes differ from the approved identity, so that an edit after the approval cannot ride on the approval.
17. As a reviewer, I want a pinned-row edit with no policy change to keep its refusal, so that the row still changes only through the outcome that pins it.
18. As a drain operator, I want the landing to admit the composed tree that holds the approved policy and the row edit, so that the decision reaches main.
19. As a reviewer, I want the landing to refuse a composed tree whose row bytes differ from the approved identity, so that the publication grades the tree that it publishes.

### Records

20. As a teammate who walks in cold, I want the changelog to state the new read locations, so that I learn how a pinned-row decision lands after the row file retires.

### One component-link walk

21. As a drain operator, I want the working read to accept a checkout that sits below a linked directory, so that a repository under a linked home or temporary directory can record a pinned-row decision.
22. As a maintainer, I want the retro write and the working read to share one component-link walk, so that the two link refusals cannot drift.

## Implementation decisions

**One predicate names a changed source.**
The commitment package gains one predicate beside `NewSource` (`internal/commitment/delivery.go:197`).
It is true for an outcome source of the proposal when a current outcome source binds the same id and path, and no current outcome source binds the exact same identity.
It composes `namesSource` (`internal/commitment/delivery.go:203`).
A deliverable does not use it: `NewDeliverable` keeps its rule (`internal/commitment/delivery.go:179`).

**One function keeps the read-location decision.**
`readTree` stays the one owner of where each source reads.
The table below states its answers after this change.

| source | plan and approve | commit and landing |
| --- | --- | --- |
| a new deliverable | checkout HEAD | candidate tree |
| a new outcome source whose path the default-branch tip lacks | checkout HEAD | candidate tree |
| a changed outcome source | working file of the store root | candidate tree |
| every other source | default-branch tip | default-branch tip |

Only the third row is new.
At plan and approve the store root is the invoking checkout, because the command builds the store from the checkout top level (`internal/commitment/commitcmd/command.go:110`).
The author chooses how the location carries the working-file form; the decision stays in one function.

**The working read composes existing guards.**
The working read refuses a symbolic link in each path component from the store root down to the file.
A link above the store root does not refuse, because a checkout can sit below a linked directory, such as the temporary directory of macOS.
For that reason the read does not use `bounds.RefuseLinks`, which walks each ancestor up to the file-system root (`internal/bounds/classify.go:147`).
The read also does not use `bounds.OpenNoFollow` (`internal/bounds/open.go:16`), because that open calls `bounds.RefuseLinks` first.
`refuseSymlinkComponents` holds the needed walk for a retro write today (`internal/roadmap/retro.go:76`).
The ticket moves that walk into `internal/bounds` once, and the retro write and the working read both call it.
The moved walk names the linked component in a neutral text, because it now serves a read and a write.
Its refusal wraps `bounds.ErrLinkRefused` (`internal/bounds/classify.go:144`), so that one sentinel marks each link refusal of the package.
No test pins the full old text. `TestRetroRefusesSymlinkedDestinationComponents` requires the words `symbolic link`, and the neutral text keeps them, so that test needs no edit.
The read then classifies the file with `bounds.ClassifyNoFollow` (`internal/bounds/classify.go:135`), which binds the control-record limit that the tree read also binds (`internal/git/run.go:154`).
Only a parsed classification passes, and each other state refuses.

**One normalization serves both reads.**
The tree read normalizes the bytes with `roadmap.RequirementBytes` and then takes `commitment.Identity` (`internal/commitment/repository/sources.go:39`).
The working read uses the same step, factored once, so the two reads cannot compute two identities for the same bytes.
`RequirementBytes` drops only the one `Occurrences:` ledger line of a row file (`internal/roadmap/occurrences.go:181`). An `Occurrence:` line stays inside the identity.

**The refusal texts and faces do not change.**
A read failure keeps the form `commitment source "<id>" refused: <cause>` (`internal/commitment/repository/sources.go:190`).
An identity mismatch keeps `validateIdentity`, which names the proposed identity and then the read identity (`internal/commitment/repository/sources.go:195`).
A plan refusal keeps the `commitment-plan-input` face, and an approval refusal keeps `commitment-decision` through `replan` (`internal/commitment/repository/admission.go:107`).
A commit refusal from `approvedTransition` keeps the `commit-handback` face (`internal/commit/commit.go:216`).
The route-bypass check scans this package, so no new error literal holds a next command (`internal/conformance/registry/route_bypass.go:16`).

**The plan identity does not change.**
`BuildPlan` keeps the old identity of a changed row as an open source of the current policy (`internal/commitment/authority.go:72`).
The proposal does not name that old binding, so it is not a changed source, and it reads the default-branch tip.
It passes at plan, approve, commit, and landing until the landing publishes.
No change reaches `BuildPlan`, the plan identity, or `protectedCandidate`.
Under a policy transition, `protectedCandidate` already reads the candidate policy at the candidate tree (`internal/commitment/repository/candidate.go:75`).

**Story 16 narrows and stays.**
The `changed-row` subtest of `TestHeadBindingKeepsTheDefaultTipForEveryOtherSource` pins the old rule (`internal/commitment/head_binding_test.go:146`).
That subtest becomes the witness for story 11: a proposal that keeps the published identity plans while the worktree holds an edit of the row.
The comment above `readTree` and the comment above that test state the new rule.

## Implementation chunks

One ticket delivers the whole outcome, because one function serves both read sides.
A plan-side ticket alone would leave an approved plan that no commit can record, and a commit-side ticket alone has no approved plan to grade.
The move of the link walk is not its own prefactor ticket, because this ticket is its sole new consumer, and the walk has one existing caller.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| PE-C1 / 1-read-changed-pinned-rows.md | An approved plan and one commit record a pinned-row decision, and every other source keeps its read | PE01, PE02, PE03, PE04, PE05, PE06, PE07, PE08, PE09, PE10, PE11, PE12, PE13, PE14, PE15, PE16, PE17, PE18, PE19, PE20, PE21, PE22, PE23, PE24, PE25, PE26, PE27, PE28, PE29 | TestPinnedRowPlanReadsTheWorkingFile, TestPinnedRowPlanRefusesBytesTheWorkingFileLacks, TestPinnedRowApproveRereadsTheWorkingFile, TestPinnedRowPlanRefusesAHostileWorkingFile, TestPinnedRowPlanNormalizesTheOccurrenceLedger, TestHeadBindingKeepsTheDefaultTipForEveryOtherSource, TestHeadBindingReadsANewRowAtHEAD, TestHeadBindingRefusesAnUnheldDeliverable, TestCommitmentApprovedRowEdit, TestCommitmentProtectedScope, TestAdmitPublicationGradesAChangedRowAtTheComposedTree, TestRetroRefusesSymlinkedDestinationComponents | no |

### Completion plan

This version 2 plan records the delegate execution, the author verification, and the final checks.

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft437-build-20261031","orchestrator_session":"session_01JK5CPcwpGroiAPWzzXxytc","author_limit":1,"assignments":{"1-read-changed-pinned-rows.md":[{"session":"ft437-author-01","assignment":"76a3920f758a0cd21a81050992211eb9","model":"opus","effort":"medium","source":"fba398122a3165ab735ed509aa85b05136ea716b","native_ref":"claude-agent:bench-writer-medium:ft437-author-01"}]}},"chunks":[{"id":"PE-C1","tickets":["1-read-changed-pinned-rows.md"],"verification":[{"id":"plan-approve","command":"bench test --package ./internal/commitment --run 'PinnedRow|HeadBinding'","probe":"Make readTree return the default-branch tip for a changed outcome source, and record that TestPinnedRowPlanReadsTheWorkingFile turns red.","ticket":"1-read-changed-pinned-rows.md"},{"id":"hostile-read","command":"bench test --package ./internal/commitment --run 'PinnedRowPlanRefusesAHostileWorkingFile'","probe":"Replace the guarded working read with a plain os.ReadFile, and record that TestPinnedRowPlanRefusesAHostileWorkingFile turns red.","ticket":"1-read-changed-pinned-rows.md"},{"id":"commit","command":"bench test --package ./internal/commit --run 'CommitmentApprovedRowEdit|CommitmentProtectedScope'","probe":"Make approvedTransition skip the source read of a changed outcome source, and record that TestCommitmentApprovedRowEdit turns red.","ticket":"1-read-changed-pinned-rows.md"},{"id":"landing","command":"bench test --package ./internal/commitment/repository --run 'AdmitPublication|AuthorizeCandidate'","probe":"Make approvedTransition read the default-branch tip for a changed outcome source, and record that TestAdmitPublicationGradesAChangedRowAtTheComposedTree turns red.","ticket":"1-read-changed-pinned-rows.md"},{"id":"retro-walk","command":"bench test --package ./internal/roadmap --run 'Retro'","probe":"Remove the call of the moved component walk from the retro write, and record that TestRetroRefusesSymlinkedDestinationComponents turns red.","ticket":"1-read-changed-pinned-rows.md"}]}],"final_verification":[{"id":"commitment-packages","command":"bench test --package ./internal/commitment/..."},{"id":"commit-package","command":"bench test --package ./internal/commit"},{"id":"walk-owners","command":"bench test --package ./internal/bounds --package ./internal/roadmap"},{"id":"coverage","command":"bench coverage --check pinned-row-edits"}]}
```

## Testing decisions

- A good test drives the store or the command with a real repository and a real worktree, and it observes the verdict and the files.
- The plan and approve rows drive `Store.Plan` and `Store.Approve` from an owned worktree of a fixture repository. `commitmenttest.SeedProtected` pins two rows (`internal/commitment/commitmenttest/repo.go:168`), and `commitmenttest.Assignment` makes the owned worktree (`internal/commitment/commitmenttest/repo.go:152`). The precedent is `internal/commitment/head_binding_test.go`.
- The commit rows drive `bench commit` through `runCommand` on `planningCommitRepo` (`internal/commit/commitment_test.go:76`). The precedent is `TestCommitmentProposalIdentity` (`internal/commit/commitment_test.go:100`).
- The landing rows drive `Store.AdmitPublication` over a composed tree. The precedent is `TestAdmitPublicationGradesANewBindingAtTheComposedTree` (`internal/commitment/repository/head_binding_test.go:128`).
- Each expected identity comes from `commitment.Identity` over the bytes that the test writes, never from the code under test.
- The ordinary Go test phase of the gate observes every row except the review-owned changelog row.

### Seam diagram

    trigger: bench commitment plan / approve (store root = the invoking checkout)
             bench commit / bench worktree land (candidate tree = the graded tree)
        │
        ▼
    current policy + proposal ──▶ [ readTree: one read location per source ] ──▶ identity compare ──▶ admit or refuse
                                      ◀ tests attach here: Store.Plan, Store.Approve,
                                        bench commit through runCommand, Store.AdmitPublication

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| PE01 | 1 | Plan accepts the new identity of a pinned row while the worktree holds the edit only in its working file | `internal/commitment/pinned_row_test.go` (`TestPinnedRowPlanReadsTheWorkingFile`) | The tip read of today refuses with `changed from`, and a HEAD read refuses because HEAD keeps the published bytes |
| PE02 | 3 | Plan accepts one proposal that changes the identities of both seeded pinned rows while both working files hold the edits | `internal/commitment/pinned_row_test.go` (`TestPinnedRowPlanReadsTheWorkingFile`) | A read that admits only the first changed source refuses the second |
| PE03 | 4 | Plan refuses a new identity that the working row file does not hold with `commitment source "FT1" changed from <proposed> to <working>` | `internal/commitment/pinned_row_test.go` (`TestPinnedRowPlanRefusesBytesTheWorkingFileLacks`) | An implementation that passes a changed source without a read admits bytes that no file holds |
| PE04 | 5 | Plan refuses a new identity that the checkout HEAD holds while the working file holds the published bytes | `internal/commitment/pinned_row_test.go` (`TestPinnedRowPlanRefusesBytesTheWorkingFileLacks`) | A HEAD read admits the committed edit that the working file reverted |
| PE05 | 10 | A refused plan of a changed pinned row leaves the milestone state of the root and the worktree unchanged | `internal/commitment/pinned_row_test.go` (`TestPinnedRowPlanRefusesBytesTheWorkingFileLacks`) | A kept-rule witness: plan reads each source before its ledger transaction, and this row reds if a change moves the read after the receipt write |
| PE06 | 2 | Approve writes the proposal bytes into the worktree policy file while the row edit stays uncommitted | `internal/commitment/pinned_row_test.go` (`TestPinnedRowApproveRereadsTheWorkingFile`) | The tip read of today refuses the approval, and a HEAD read refuses it too |
| PE07 | 6 | Approve refuses a plan whose working row file changed after the plan with the changed-source text and leaves the worktree policy file unchanged | `internal/commitment/pinned_row_test.go` (`TestPinnedRowApproveRereadsTheWorkingFile`) | A kept-rule witness: approve reads each source again on each call, and this row reds if the approval trusts the plan-time read |
| PE08 | 7 | Plan refuses a changed pinned row whose working path is absent with `commitment source "FT1" refused:` | `internal/commitment/pinned_row_hostile_test.go` (`TestPinnedRowPlanRefusesAHostileWorkingFile`) | A read that falls back to the tip or to HEAD for a missing file admits a deleted row |
| PE09 | 7 | Plan refuses a zero-byte working row file even when the proposal names the identity of no bytes | `internal/commitment/pinned_row_hostile_test.go` (`TestPinnedRowPlanRefusesAHostileWorkingFile`) | A plain file read admits an empty control record, which ADR 0010 calls a failure |
| PE10 | 7 | Plan refuses a working row file that holds invalid UTF-8 even when the proposal names the identity of those bytes | `internal/commitment/pinned_row_hostile_test.go` (`TestPinnedRowPlanRefusesAHostileWorkingFile`) | A plain file read admits bytes that the classification calls malformed |
| PE11 | 7 | Plan refuses a working row file one byte over the control-record limit even when the proposal names the identity of those bytes | `internal/commitment/pinned_row_hostile_test.go` (`TestPinnedRowPlanRefusesAHostileWorkingFile`) | An unbounded read admits a row that the bounded tree read refuses |
| PE12 | 7 | Plan refuses a working row path that is a live link to a file with the proposed bytes or a dangling link | `internal/commitment/pinned_row_hostile_test.go` (`TestPinnedRowPlanRefusesAHostileWorkingFile`) | A read that follows the link admits bytes outside the checkout |
| PE13 | 7 | Plan refuses a working row path that is a FIFO and returns without a writer on the FIFO | `internal/commitment/pinned_row_hostile_test.go` (`TestPinnedRowPlanRefusesAHostileWorkingFile`) | A read that opens the path before its type check blocks in open(2) |
| PE14 | 7 | Plan refuses a working row whose `roadmap` directory is a link to a directory that holds the proposed bytes | `internal/commitment/pinned_row_hostile_test.go` (`TestPinnedRowPlanRefusesAHostileWorkingFile`) | A read that checks only the final component admits bytes behind a linked ancestor |
| PE15 | 8 | Plan accepts the identity of the working row bytes without their `Occurrences:` ledger line | `internal/commitment/pinned_row_test.go` (`TestPinnedRowPlanNormalizesTheOccurrenceLedger`) | A raw-byte identity differs from the tree identity of the same bytes, so it refuses the identity that the commit reads |
| PE16 | 9 | Plan refuses a working row file with two `Occurrences:` lines with `has a malformed occurrence ledger` | `internal/commitment/pinned_row_test.go` (`TestPinnedRowPlanNormalizesTheOccurrenceLedger`) | A working read without the ledger step admits a ledger that the tree read refuses |
| PE17 | 11 | Plan accepts a proposal that keeps the published identity of a pinned row while the worktree HEAD and working file hold an edit of that row | `internal/commitment/head_binding_test.go` (`TestHeadBindingKeepsTheDefaultTipForEveryOtherSource`) | A working-file read for every source refuses the kept identity with `changed from` |
| PE18 | 12 | Plan refuses a new outcome source at a path that main holds and that no current source binds | `internal/commitment/head_binding_test.go` (`TestHeadBindingReadsANewRowAtHEAD`) | A working-file read for each path that main holds admits an unpinned row edit |
| PE19 | 13 | Plan refuses a new deliverable that only the working tree holds | `internal/commitment/head_binding_test.go` (`TestHeadBindingRefusesAnUnheldDeliverable`) | A working-file read for deliverables admits an uncommitted spec |
| PE20 | 14 | `bench commit` of the approved policy, the board, and the row edit exits 0 and adds one commit whose tree holds the new row bytes | `internal/commit/commitment_test.go` (`TestCommitmentApprovedRowEdit`) | The tip read of today in `approvedTransition` refuses the commit |
| PE21 | 15 | `bench commit` that names the approved policy and the board without the row edit exits 1 with `commitment source "FT1" changed from <new> to <published>` and leaves HEAD in place | `internal/commit/commitment_test.go` (`TestCommitmentApprovedRowEdit`) | An implementation that passes each changed source without a read admits the stale row |
| PE22 | 16 | `bench commit` whose row bytes differ from the approved identity exits 1 with `commitment source "FT1" changed from <new> to <rewritten>` and leaves HEAD in place | `internal/commit/commitment_test.go` (`TestCommitmentApprovedRowEdit`) | A commit that trusts the approval without a candidate read admits an edit after the approval |
| PE23 | 17 | `bench commit` of a pinned-row edit with no policy change refuses and leaves HEAD in place | `internal/commit/commitment_test.go` (`TestCommitmentProtectedScope`) | A kept-rule witness: with no policy change, `protectedCandidate` reads the current policy at the candidate tree, and this row reds if the change weakens that read |
| PE24 | 18 | `Store.AdmitPublication` admits the composed tree that holds the approved policy and the row edit | `internal/commitment/repository/head_binding_test.go` (`TestAdmitPublicationGradesAChangedRowAtTheComposedTree`) | The tip read of today in `approvedTransition` refuses the landing |
| PE25 | 19 | `Store.AdmitPublication` refuses a composed tree whose row bytes differ from the approved identity with the changed-source text | `internal/commitment/repository/head_binding_test.go` (`TestAdmitPublicationGradesAChangedRowAtTheComposedTree`) | A landing that trusts the commit-time read admits bytes that it never graded |
| PE26 | 20 | The Unreleased section of `CHANGELOG.md` holds one entry that names the working-file read at plan and approve and the candidate-tree read at commit and landing | review-owned: the Unreleased section of CHANGELOG.md | A missing entry leaves the read rule only in code after the row file retires |
| PE27 | 21 | Plan accepts the new identity of a pinned row when the store root of the worktree is a path through a linked parent directory | `internal/commitment/pinned_row_test.go` (`TestPinnedRowPlanReadsTheWorkingFile`) | A link walk up to the file-system root, as `bounds.RefuseLinks` does, refuses this checkout |
| PE28 | 22 | `bench retro` still refuses a retro path whose destination component is a symbolic link | `internal/roadmap/retro_test.go` (`TestRetroRefusesSymlinkedDestinationComponents`) | A kept-rule witness: the moved walk must keep the refusal of the retro write |
| PE29 | 22 | The tree holds one root-relative component-link walk, in `internal/bounds`, and `internal/roadmap` holds no copy | review-owned: `rg -n "refuseSymlinkComponents" internal` returns no hit after the move | A surviving copy is a second source for one link rule |

### Edge inventory

These behaviors serve the kit and every linked repository, because the commitment verbs ship.
The working read reads one file, so no directory read has an absent-versus-empty pair.
PE08 and PE09 cover the absent-versus-empty pair of that file.
PE12 and PE13 cover the live link, the dangling link, and the special file of the hostile-input checklist.
PE14 covers a linked ancestor directory inside the checkout, and PE27 covers a linked directory above the checkout.
A path with spaces or glob characters cannot occur, because the policy parse refuses a source path that is not one clean relative path (`internal/commitment/parse.go:232`).
A last line without a trailing newline keeps its bytes, because the ledger step splits and joins on the newline only.
No test swaps a package variable.

Won't handle: a line-ending conversion or a clean or smudge filter that makes the working bytes differ from the blob — the commit and the landing read the candidate tree again and refuse the mismatch.
Won't handle: a plan from the primary checkout that reads the working file of the primary checkout — `bench commitment approve` refuses the primary checkout, so no policy stages from it.
Won't handle: an untracked working row file — the commit reads the candidate tree of the named tracked paths, so an untracked edit cannot land.
Won't handle: a proposal that changes one of two identities that bind one row path — `protectedCandidate` refuses that commit at the candidate tree. No current outcome shares a path.
Won't handle: the held row edits in the saved drain patch — the next `/bench-drain` applies them through this route.
Won't handle: the `rejectSymlinkComponents` walk of the freshness build inputs (`internal/freshness/freshness_buildinputs.go:211`) — it walks from the volume root, as `bounds.RefuseLinks` does, so it is not a copy of the moved root-relative walk.
Won't handle: the interim spec-stage rule for a pinned row — only `roadmap/FT437.md` holds it, and the promote-then-delete of this spec removes that file.

## Out of scope

- A working-file identity column in `bench commitment inventory` is a separate capability. Today the plan refusal names the working identity. Estimate: 4 edits, 1 gate run.
- A dedicated refusal face for a commit whose candidate row differs from its approval is a separate capability. Estimate: 3 edits, 1 gate run.
- A plan-time refusal of one proposal that binds one row path at two identities is a separate capability. Estimate: 3 edits, 1 gate run.

## Further notes

### Source-sentence-to-row table

| source sentence | rows |
| --- | --- |
| Criterion: "An approved plan that carries a pinned row's new identity, together with the row edit, records the decision in one plan and one commit" | PE01, PE02, PE06, PE20, PE24 |
| Criterion: "a pinned-row edit with no exact approved plan still refuses" | PE21, PE22, PE23, PE25 |
| Row: "Plan and approve read a source whose identity the proposal changes at the checkout's working file." | PE01, PE03, PE04, PE05, PE06, PE07, PE08, PE09, PE10, PE11, PE12, PE13, PE14, PE15, PE16 |
| Row: "The commit and the landing read it at the candidate tree." | PE20, PE21, PE22, PE24, PE25 |
| Row: "the spec narrows story 16 of `spec-stage-binding` (SB10) and does not remove it" | PE17, PE18, PE19, PE23 |
| Row: "until this row ships, `/bench-write-spec` entry reads a pinned row's bound outcome criteria as its decision source" | none — Won't handle, because the rule retires with the row file |

### Flagged additions

- PE08 to PE14 go beyond the source sentences. The hostile-input checklist requires them for a new working-file read.
- PE15 and PE16 go beyond the source sentences. They keep the plan identity equal to the identity that the commit reads.
- PE26 follows the tree convention for a user-facing change. The "One-landing row admission" entry of `CHANGELOG.md` is the precedent for a change of the plan and approve reads.
- PE05, PE07, and PE23 are kept-rule witnesses, not new promises. Each one pins a guarantee that the change must not weaken.
- PE27 to PE29 go beyond the source sentences. The coordinator added them after the spec review accepted, because a walk up to the file-system root would refuse each checkout below a linked directory. The move reuses the one existing walk instead of a second copy.
- No guidance edit is planned. `.agents/commands/bench-drain.md:107` already routes a protected row through `bench commitment plan`. `.bench/BENCH.md:153` and ADR 0028 state that a pinned row changes only through its outcome, and each statement stays true.

### Reader sweep

- Go readers of the read location: `readTree` and its three callers `validateDeliverables`, `validateSources`, and `approvedTransition`; `Store.Plan` and `Store.approve` call `validateSources`; `authorizeCandidate` calls `approvedTransition` for `bench commit` and for the landing.
- Adjacent readers that keep the default-branch tip: `Store.Inventory`, `validateDeliverable`, and `completionEvidence`. None of them reads a changed source of a proposal.
- Guidance readers: `.bench/BENCH.md:153`, `docs/adr/0028-the-commitment-gates-spec-implementations.md:17`, and `.agents/commands/bench-drain.md:107`. Each one stays true.
- No `bench anchors` needle sits on these lines, and no line holds a repo-only path beside a shipped claim word.
- `docs/known-issues.md` holds no row on commitment sources.

### Pre-review proof checklist

- Cited symbols: `readTree` (`internal/commitment/repository/sources.go:98`), `sourceAt` (`internal/commitment/repository/sources.go:27`), `validateSourceAt` (`internal/commitment/repository/sources.go:187`), `validateIdentity` (`internal/commitment/repository/sources.go:195`), `NewSource` (`internal/commitment/delivery.go:197`), `NewDeliverable` (`internal/commitment/delivery.go:179`), `namesSource` (`internal/commitment/delivery.go:203`), `BuildPlan` (`internal/commitment/authority.go:48`), `approvedTransition` (`internal/commitment/repository/candidate.go:190`), `protectedCandidate` (`internal/commitment/repository/candidate.go:239`), `Store.Plan` (`internal/commitment/repository/repository.go:108`), `Store.approve` (`internal/commitment/repository/repository.go:165`), `Store.AdmitPublication` (`internal/commitment/repository/publication.go:21`), `roadmap.RequirementBytes` (`internal/roadmap/tree_validation.go:32`), `bounds.ClassifyNoFollow` (`internal/bounds/classify.go:135`), `bounds.RefuseLinks` (`internal/bounds/classify.go:147`), and `refuseSymlinkComponents` (`internal/roadmap/retro.go:76`).
- Import edges: the production code of `internal/commitment/repository` gains `internal/bounds`. `go list -deps ./internal/bounds` names no `internal/commitment` package, so the edge makes no cycle. `internal/roadmap` already imports `internal/bounds` (`internal/roadmap/tree.go:11`).
- Source-row clauses and occurrences: the table above quotes each clause. The row records seven occurrences, from the 2026-10-18 drain to the 2026-10-29 drain re-plan (`roadmap/FT437.md:23`).
- Promised field labels: none. No record field changes, so no `Persisted fields:` line applies.
- Changed-function callers: at 4bea51fc8, `readTree` has callers at `candidate.go:211`, `sources.go:152`, and `sources.go:165`; `validateSources` has callers at `repository.go:129` and `repository.go:217`; `approvedTransition` has its caller at `candidate.go:67`; `sourceAt` has callers at `closure.go:40`, `repository.go:62`, `sources.go:18`, `sources.go:61`, `sources.go:82`, `sources.go:133`, and `sources.go:188`; and `NewSource` has its caller at `sources.go:99`. No test calls `sourceAt` directly. `refuseSymlinkComponents` has one caller, at `internal/roadmap/retro.go:50`.
- Copy survival: the working read and the tree read share one normalization step. PE15 turns red when the working read omits that step. PE29 names the copy-survival check of the moved walk.
- Rendered-shape readers: the texts `changed from` and `refused:` keep their producers. The moved walk changes the text `destination component is a symbolic link`, and a repository-wide `rg` finds only its producer at `internal/roadmap/retro.go:88`.
- Pin operators: the plan rows compare the exact error text where the precedent does (`internal/commitment/head_binding_test.go:163`). The commit rows use `strings.Contains` on the standard error, as `refuseCommit` does (`internal/commit/commitment_test.go:86`).
- Entry reads: `Store.Plan` and `Store.approve` read the checkout HEAD and the working file under `Store.Root`. The tests set `Store.Root` to the fixture worktree, so no ambient read remains.
- Derived expectations: each identity comes from `commitment.Identity` over the literal bytes that the test writes. PE15 derives its identity from the literal body without the ledger line, not from `RequirementBytes`.
- Consolidated rules: the read-location table in Implementation decisions is the consolidation table. Its only changed cells are the third row, and PE01 to PE16, PE20 to PE22, PE24, and PE25 cover them.
- Quantified obligations: "a source whose identity the proposal changes" quantifies over each plan source. PE02 covers two changed rows in one plan, and PE17 to PE19 cover each kind of source that keeps its read.
- Workflow-step writes: no new workflow step. Plan still writes only its receipt, and approve still writes only the policy file and the board (`internal/commitment/repository/replacement.go:19`).
