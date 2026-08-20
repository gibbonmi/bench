# Close the four destructive-git deny gaps

Blocked by: none
Writes: internal/gitguard

## What to build

`git stash drop`, `git stash clear`, `git filter-branch`, and `git rm -rf` are
allowed at HEAD. Each gains its class in the classifier's deny table, so the
core stays the one place a verdict is decided and the shim gains nothing. `git
rm` denies on the recursive-and-forced combination, reading its short-flag
clusters the way `clean` already does, so an ordinary `git rm <path>` still
runs. The safe stash forms — bare `git stash`, `git stash list`, `git stash pop`
— stay allowed.

The new rules resolve their verb the way the existing ones do, so a sought token
supplied only as a flag value or after `--` cannot land on an allow.

## Acceptance

- [ ] each of the four spellings returns a block verdict carrying its class label.
- [ ] bare `git stash`, `git stash list`, `git stash pop`, and `git rm <path>` return an allow verdict.
- [ ] a sought token supplied only as a flag value reaches the destructive verdict, not an allow.
