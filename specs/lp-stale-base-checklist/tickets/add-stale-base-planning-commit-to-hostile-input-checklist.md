# Add the stale-base planning commit to the hostile-input checklist

Blocked by: none
Writes: projects/benchkit.md
Covers: none

## What to build

The benchkit hostile-input checklist gains one edge class. A commit that changes the commitment policy grades its tree against the tip of `main`. When the worktree base trails `main`, the diff shows each newer production change of `main` in reverse, and the commit refuses. The class tells a spec author to fold `main` before the plan. The class also tells the author to assert the refusal and the admission.

## Acceptance

- [ ] `rg -n "base trails" projects/benchkit.md` finds the new class inside the hostile-input checklist.
- [ ] The new class names the fold of `main` before the plan as the remedy.
- [ ] The new class tells the spec to assert the refusal on the stale base and the admission after the fold.
