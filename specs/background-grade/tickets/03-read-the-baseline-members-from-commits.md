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
On 2026-10-26, a plan commit during the prepared GR-C2 gate made that gate refuse with a changed subject.

Read both baseline members from the commit that the identity root names, never from its working files.
A clean checkout then gives the same key as before, so a current green keeps answering.

First reproduce the GR-C2 refusal: a chunk-subject run across a plan commit in its worktree refuses today.
Record which member moved in the commit message of this ticket.

`subject.go` and `checkpoint.go` are at or over their line budget. Keep each line count, or move code into `baseline_commit.go`.

## Acceptance

- [ ] A chunk-subject run keeps its subject when a Go build input of the worktree changes during the oracle.
- [ ] A chunk-subject run across a plan commit in its worktree records its green.
- [ ] The baseline schedule member frames the manifest bytes of the identity commit, not the working file.
- [ ] A clean checkout gives the same grade key before and after this change.
- [ ] `bench structure --growth` reports no growth of `internal/gate/subject.go` or `internal/gate/checkpoint.go`.
