# Apply the landed plan and settle records

Blocked by: plan-the-landed-set-under-one-fingerprint.md
Writes: internal/worktree/clean_landed.go, internal/worktree/clean_landed_apply_test.go, internal/worktree/clean_landed_hostile_test.go, internal/worktree/resume.go

## What to build

`bench worktree clean --landed --apply <fp>` re-selects and re-plans the set, recomputes
the set fingerprint (same tuple, flags, and tag as the plan) and refuses with the existing
stale-fingerprint error, mutating nothing, unless it equals the given value; then for each
removable row in id order it re-plans the row, recomputes its seven-field tuple, refuses
and stops if the tuple differs from the validated one, and otherwise runs the existing
per-path cleanup transaction with the fresh per-path fingerprint and a terminal that
marks the assignment `complete` so the transaction's completion deletes the record.
Retained rows are reported unchanged. Proven-landed branches are deleted with their trees
as per-path clean does; `--discard-branch` changes no bulk outcome beyond `detail`;
`--discard-ignored` and `--full` reach every row. Exit 0 when every removable row
completed. Demo: apply, then `bench worktree list`.

## Acceptance

- [ ] `(covers LR8)` Apply on LR7's pool removes both clean trees (`removed`), deletes their proven-landed branches, leaves the dirty tree and its `retain` row, exits 0, and `bench worktree list` immediately shows no row for the removed assignments and still shows the dirty one.
- [ ] `(covers LR9)` A new landed row, an uncommitted file, a lease going live, or a fast-forward that changes a planned HEAD OID while staying landed each make `--apply <fp>` exit 1 with the stale-fingerprint error, remove nothing, settle nothing.
- [ ] `(covers LR10)` An unparseable-lease tree plans `retain` (`uncertain`) with its per-path remedy, counts `landed` in the summary, and is present after apply.
- [ ] `(covers LR13)` Only the second of two removable rows has undeclared residue: bare retains it (`ignored`); `--discard-ignored` plans it `discard-remove`, `--full` widens its preview; apply removes both and deletes both branches without `--discard-branch`; with `--discard-branch` the outcome is identical and `detail` names the assertion; a squash-landed unprovable branch is not selected.
- [ ] Before any mutation, a fingerprint with one hex digit altered, or the same rows planned under a different flag set, is refused with the stale-fingerprint error, nothing removed, nothing settled.
- [ ] The LR17 spaced/glob tree is removed by apply; the LR18 ESC tree and every LR19 special-shaped path are present after apply.
