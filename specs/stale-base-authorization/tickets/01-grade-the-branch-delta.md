# Grade the commit candidate by the branch's own change

Blocked by: none
Writes: internal/commitment/repository/
Covers: none

## What to build

`bench commit` authorizes the candidate tree against the default-branch tip. When main moves after the branch point, each main change reads as a change that the branch makes. A stale branch that never touched those paths is then refused. The refusals are "assignment has no current delivery binding", "candidate policy has no exact approval", and "protected recommended sequence changed". Their printed routes do not recover the branch, and a dirty checkout cannot fold main. The FT290 close and the ft362 planning branch hit this deadlock on 2026-10-07.

Make the commit authorization grade only the change of the branch, from the merge base of the checkout and the default branch. Authority stays with the default-branch tip: the current policy, its approvals, and its protected rows. A branch that does not change the commitment policy has no policy transition. Then the production paths, the recommended sequence, the protected row owners, and the unsettled sources grade against the merge base.

A branch that changes the policy keeps the current grading against the tip. A publication keeps the current grading against the tip, because the landing composes main into its tree. If the merge base does not resolve, the authorization uses the tip as the reference.

## Acceptance

- [ ] After main lands a production file, an unbound branch from the earlier main commits a planning-only path.
- [ ] After main approves a policy edit, a branch from the earlier main commits a light-path change.
- [ ] After main closes a delivery, a branch from the earlier main commits a light-path change. The closure changes the policy, the recommended sequence, and the row owners.
- [ ] After main pins a new roadmap row, a branch from the earlier main commits a light-path change.
- [ ] A stale branch that edits the policy is refused with "candidate policy has no exact approval".
- [ ] A stale branch that changes the recommended sequence is refused with "protected recommended sequence changed".
- [ ] A stale branch that deletes a pinned row owner is refused with "candidate changes protected commitment".
- [ ] A stale branch that writes a production path with no ticket is refused as unbound.
- [ ] The existing light-path, publication, and commit tests stay green.
