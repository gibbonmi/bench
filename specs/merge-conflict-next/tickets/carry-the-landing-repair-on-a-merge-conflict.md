# Carry the landing's conflict repair on a merge conflict

Blocked by: none
Writes: internal/worktree/merge.go, internal/worktree/land_refusal.go, internal/worktree/merge_test.go

## What to build

When `bench worktree merge --from <incoming> <target>` refuses a composition
conflict, the `refused{...}` record carries a `next=` field. Before this ticket that
field was empty, and only the landing's conflict refusal composed the hand repair.
The merge verb's `next=` is the landing's repair up to and including the commit
step. Its text is the merge step, the clause
`(bench worktree merge refuses this conflict; resolve it by hand)`, and
`; then bench commit`. The merge step is `git -C '<path>' merge '<incoming>'`, where
`<incoming>` is the full commit the verb resolved from `--from`.

Both refusals render that prefix from one producer. The landing appends its
review-and-land tail to it,
so the landing's `next=` stays byte-identical (WL18 in
`internal/worktree/land_surface_test.go`). A worktree path that is not line-safe
takes the pointer form the producer already uses.

The decision source is `decisions/worktree-exec-comfort.md` #8. The light-path item
is that the merge conflict refusal carries the same `next=` the landing conflict
refusal composes. Everything else about the refusal stays as WM9 pins it. That is
exit 1, `refused{` on stdout alone, `detail=composition conflict: <kind>`, and the
`refusal_paths` table. The branch tip, the checkout, and the lane record stay
unchanged.

## Acceptance

- [ ] MC1: a merge that refuses a textual conflict on a line-safe worktree path
      prints the `next=` prefix inside the `refused{` record, with `<incoming>` as
      the full commit id. Seam: `internal/worktree/merge_test.go`, a new row beside
      WM9 that asserts the exact text. Red when the merge path leaves `next` empty.
- [ ] MC2: the same refusal keeps exit 1, `detail=composition conflict: textual`,
      `paths_total=1`, the conflicting path, and an unchanged tip, HEAD, and lane
      tally. Seam: WM9 extended with a `next=` presence check. Red when the change
      alters the other fields or writes to the tree.
- [ ] MC3: the landing's conflict `next=` is unchanged. Seam: WL18 passes unchanged,
      and the clause `refuses this conflict` appears once in non-test source under
      `internal/worktree/`. Review-owned: the reviewer runs
      `rg -c 'refuses this conflict' internal/worktree/` and expects one non-test hit.
