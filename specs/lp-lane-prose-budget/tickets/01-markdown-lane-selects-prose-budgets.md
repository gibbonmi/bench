# The markdown lane class selects the prose-budget check

Blocked by: none
Writes: internal/gate/lane_select.go, internal/gate/lane_select_test.go, internal/gate/lane_run_test.go, projects/benchkit.md
Covers: none

## What to build

A worktree `bench commit` lane passed `.agents/commands/bench-implement-spec.md` at 95 lines.
The prose-budget table of `projects/benchkit.md` gives that file a budget of 93 lines.
On the same tree, `bench test --check guidance-prose-budgets` failed in 14 ms with "prose-budget exceeded".

The `markdown` path class in `internal/gate/lane_select.go` selects only `prose` and `docs-currency-workflow`.
Only a change to `projects/benchkit.md` selects `guidance-prose-budgets`.
Thus a lane commit of an over-budget guidance file passes, and the landing gate then fails.

To reproduce the defect, do these steps in a worktree:

1. Add two lines to `.agents/commands/bench-implement-spec.md`.
2. Add a one-ticket light-path folder that covers the file.
3. Run `bench commit --dry-run` on these paths.

The lane passes with the checks `prose` and `docs-currency-workflow` only.

Add `guidance-prose-budgets` to the checks of the `markdown` class.
Each changed Markdown path then runs the budget check, which costs approximately 14 ms.
The class table stays the one source for the selection.
The `selected by` cell of the profile lane table shows the new class.

## Acceptance

- [ ] The lane of a guidance Markdown change selects `guidance-prose-budgets`.
- [ ] A `bench commit --dry-run` of an over-budget guidance edit refuses.
- [ ] The landing gate is green.
