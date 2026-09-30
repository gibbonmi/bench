# Pin the escape rule of the shared worktree refusal printer

Blocked by: none
Writes: internal/sanitize/sanitize.go, internal/worktree/target_refusal_test.go (new), reviews/tree-targets.md
Covers: none

## What to build

The shared refusal printer of the five target-taking worktree verbs prints a line that passes `sanitize.LineSafe` as it is. It escapes a line that fails the predicate through `sanitize.Controls`. This rule is the approved C4-P5 decision of the tree-targets build.

The doc of `sanitize.LineSafe` says that a caller that fails the predicate emits a pointer instead of the value. That sentence contradicts the approved printer. Make the doc name both approved responses: a caller emits a pointer, or a line-level printer escapes the whole line through `sanitize.Controls`. Keep the rest of the doc.

A table test in `internal/worktree` pins the escape rule at the seam of the shared printer. Today the only pin is in `internal/treetarget`, through the exported build refusal. The new test calls the printer through an in-package path. It asserts each exact stderr line for a plain label, a label with a backslash, and a detail or label with a control byte.

This ticket is the third TT-C4 repair cycle, which the reviewer approved on 2026-09-30 as an extension of the repair allowance. Record that extension and the work it permits in `reviews/tree-targets.md`.

## Acceptance

- [ ] The doc of `sanitize.LineSafe` names the pointer response and the escape response, and it no longer states the pointer as the only response.
- [ ] A test in `internal/worktree` wants `next=bench worktree build 'a\b'` as typed for a label with a backslash.
- [ ] The same test wants each line escaped, with no raw control byte on stderr, when the detail or the label holds a control byte.
- [ ] A probe that always escapes, and a probe that never escapes, each red the new test.
- [ ] `reviews/tree-targets.md` records the reviewer extension to three repair cycles for TT-C4 and names this ticket as its permitted work.
