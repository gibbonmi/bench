# Admit the review record deletion of a retired spec

Blocked by: none
Writes: internal/commitment/
Covers: none

## What to build

`bench spec retire <slug>` deletes the `specs/<slug>/` folder and the review record `reviews/<slug>.md`. A close worktree has no delivery binding, because the delivered outcome is complete. The candidate authorization reads each path under `specs/` as a planning path, but it reads the review record deletion as a production path. An unbound worktree then refuses the retirement commit with "assignment has no current delivery binding". No current sequence commits a delivered-spec retirement.

Make the candidate authorization read the deletion of `reviews/<slug>.md` as a planning change when the candidate tree holds no `specs/<slug>/` folder. Any other change to a review record stays a production change, so a review record write still needs its binding. Compose the existing owner of the review record path, and do not write a second derivation of that path.

## Acceptance

- [ ] In an unbound worktree, a commit that deletes `specs/<slug>/` and `reviews/<slug>.md` together is admitted.
- [ ] In an unbound worktree, a commit that deletes `reviews/<slug>.md` while `specs/<slug>/` stays in the candidate tree is refused as unbound.
- [ ] In an unbound worktree, a commit that changes the content of `reviews/<slug>.md` is refused as unbound, with or without the spec folder.
- [ ] Dogfood: with this candidate binary, `bench commit --dry-run` of `specs/ft290-test-projection` and `reviews/ft290-test-projection.md` in the `ft290-close` worktree passes.
