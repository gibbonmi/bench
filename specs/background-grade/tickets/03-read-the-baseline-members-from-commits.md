# Read the baseline members from commits

Blocked by: 02-derive-the-run-binary-source-once.md
Writes: internal/gate/prospective.go, internal/gate/subject.go, internal/gate/evaluation.go, internal/gate/engine.go, internal/gate/checkpoint.go, internal/gate/baseline_commit.go (new), internal/gate/baseline_commit_test.go (new), internal/gate/grade_identity_test.go, CHANGELOG.md
Covers: BG16, BG17, BG18

## What to build

A prospective subject reads its graded tree from a private checkout.
It also reads two members from the live identity root, which is the worktree for each checkpoint.
The baseline schedule enters the grade identity from that root (`internal/gate/subject.go:140`, `internal/gate/prospective.go:75-88`).
The baseline runner digest enters the oracle digest from that root (`internal/gate/prospective.go:103-116`).
The run compares its subject before and after the oracle, so an edit of the worktree during the run refuses it (`internal/gate/run_settle.go:16-39`).
A commit that changes only the spec, the tickets, or the review record moves neither member, because `sameSubject` compares no identity member (`internal/gate/run_settle.go:37-39`). A commit that changes a Go build input moves the baseline runner digest and refuses the run.

Read both baseline members from the commit that the identity root names, never from its working files.
The run already materializes that commit, minus the review record, in its private checkout (`internal/gate/engine.go:55-61`).
For the chunk-subject path, pass that checkout as the baseline root: it feeds the `baseline` argument (`internal/gate/engine.go:65-66`), the baseline schedule, and the baseline runner digest.
Select it from `executeChunkCheckpoint` through a context key, as `deltaSourceKey` does. Every other caller keeps the live root.
Canonicalize the checkout root as `internal/gate/subject.go:66` does, so that each relative path keeps its frame.
A clean checkout then gives the same key as before, so a current green keeps answering.
Do not change the GR3 to GR6 cases of `internal/gate/baseline_schedule_identity_test.go`, because the builder still takes a directory. Add one committed-symlink case to `baseline_commit_test.go` that matches GR5.

First reproduce the refusal: a subject-alone chunk run across a commit that changes a Go build input in its worktree refuses today.
The live manifest reads of `internal/gate/delta_checkpoint.go` are run decisions, not key members, and ticket 07 owns them.

`subject.go` has 1 line of headroom under its grant, and `checkpoint.go` has 2. Keep each line count, or move code into `baseline_commit.go`.

## Acceptance

- [ ] A chunk-subject run keeps its subject when a Go build input of the worktree changes during the oracle.
- [ ] A chunk-subject run across a commit that changes a Go build input in its worktree records its green.
- [ ] The baseline schedule member frames the manifest bytes of the identity commit, not the working file.
- [ ] A clean checkout gives the same grade key before and after this change.
- [ ] `bench structure --growth` reports no growth of `internal/gate/subject.go` or `internal/gate/checkpoint.go`.
