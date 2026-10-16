# Give each structure line grant 15 percent headroom

Blocked by: none
Writes: .bench/structure.budgets
Covers: none

## What to build

Each line-count grant in `.bench/structure.budgets` is an exact cap today. Thus each new line in a granted file needs a new reviewer decision, and the author finds the cap only at commit time. The reviewer accepted headroom of about 15 percent on each grant (learning L327, 2026-10-16). The headroom is over the actual line count of the file, not over the existing grant.

Set each existing file line-count grant to the actual line count of the file at the ticket base, multiplied by 1.15 and rounded up. Do not set a grant lower than its existing value. Do not add a grant for a file that has no grant. Do not change a directory file-count grant or the default limits.

Reword the header comment of the file to state the headroom rule. The header also states that a directory file-count grant stays an exact cap. Remove each grant comment that states a re-ask rule, because the header is the one source of that rule.

## Acceptance

- [ ] Each file line-count grant is at least 1.15 times the line count of the file at the ticket base, rounded up.
- [ ] No file line-count grant is lower than its value at the ticket base.
- [ ] The directory file-count grants `.bench/`, `internal/gate/`, `internal/worktree/`, and `internal/shift/` keep their values.
- [ ] The file contains no grant for a path that had no grant at the ticket base.
- [ ] The header comment states the headroom rule and the exact directory cap, and no grant comment states a re-ask rule.
- [ ] `bench structure` reports no new violation, and the structure package tests stay green.
